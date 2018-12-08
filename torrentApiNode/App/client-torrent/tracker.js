const bignum = require('bignum');
const bencode = require('bencode');
const net = require('net');
const dgram = require('dgram');
const crypto = require('crypto');
const urlParser = require('url').parse;
const querystring = require('querystring').escape;
const fetch = require('node-fetch');
const http = require('http');
const DHT = require('bittorrent-dht');


const { genId, infoHash, size } = require('./tools.js')
const { download } = require('./download.js');

const buildAnnounceReq = (connId, torrent, port = 6881) => {
	const buf = Buffer.allocUnsafe(98);

	// connection id
	connId.copy(buf, 0);
	// action
	buf.writeUInt32BE(1, 8);
	// transaction id
	crypto.randomBytes(4).copy(buf, 12);
	// info hash
	infoHash(torrent).copy(buf, 16);
	// peerId
	genId().copy(buf, 36);
	// downloaded
	Buffer.alloc(8).copy(buf, 56);
	// left
	size(torrent).copy(buf, 64);
	// uploaded
	Buffer.alloc(8).copy(buf, 72);
	// event
	buf.writeUInt32BE(0, 80);
	// ip address
	buf.writeUInt32BE(0, 80);
	// key
	crypto.randomBytes(4).copy(buf, 88);
	// num want
	buf.writeInt32BE(-1, 92);
	// port
	buf.writeUInt16BE(port, 96);

	return buf;
}

const buildConnReq = () => {
	const buf = Buffer.alloc(16);

	// connection id
	buf.writeUInt32BE(0x417, 0);
	buf.writeUInt32BE(0x27101980, 4);

	// action
	buf.writeUInt32BE(0, 8);

	// transaction id
	crypto.randomBytes(4).copy(buf, 12);

	return buf;
};

const respType = (resp) => {
	const action = resp.readUInt32BE(0);
	if (action === 0)
		return 'connect';
	else if (action === 1)
		return 'announce';
};

const parseConnResp = (resp) => {
	return {
		action: resp.readUInt32BE(0),
		transactionId: resp.readUInt32BE(4),
		connectionId: resp.slice(8)
	};
};

const httpparseAnnounceResp = (resp) => {
	function group(iterable, groupSize) {
		let groups = [];
		for (let i = 0; i < iterable.length; i += groupSize) {
			groups.push(iterable.slice(i, i + groupSize));
		}
		return groups;
	}

	try {
		if (bencode.decode(resp)['failure reason'] !== undefined) {
			return [];
		}

	} catch (e) {
		return [];
	}

	let peersBuff = bencode.decode(resp).peers;
	let peersArray = group(peersBuff.slice(0), 6).map(address => ({
		ip: address.slice(0, 4).join('.'),
		port: address.readUInt16BE(4)
	}));

	return peersArray;
};

const parseAnnounceResp = (resp) => {
	function group(iterable, groupSize) {
		let groups = [];
		for (let i = 0; i < iterable.length; i += groupSize) {
			groups.push(iterable.slice(i, i + groupSize));
		}
		return groups;
	}

	return {
		action: resp.readUInt32BE(0),
		transactionId: resp.readUInt32BE(4),
		leechers: resp.readUInt32BE(8),
		seeders: resp.readUInt32BE(12),
		peers: group(resp.slice(20), 6).map(address => {
			return {
				ip: address.slice(0, 4).join('.'),
				port: address.readUInt16BE(4)
			}
		})
	}
}

const getPeers = (infoMovie, torrent, pieces, files, videoToHlsCallback, addToStreamCallback, changeStatusCallback) => {
	let announceListUrl = [];

	if (torrent['announce-list']) {
		torrent['announce-list'].forEach(arrayOfBuff => {
			arrayOfBuff.forEach(buff => {
				let url = null;
				try {
					url = urlParser(buff.toString('utf8'));
				} catch (e) {
				}

				if (url !== null)
					announceListUrl.push(url)
			});
		});
	} else if (torrent['announce']){
		announceListUrl.push(urlParser(torrent.announce.toString('utf8')));
	}


	announceListUrl.forEach(announce => {
		const { protocol } = announce;
		if (protocol === 'udp:') {
			udpGetList(infoMovie, torrent, pieces, announce, files, videoToHlsCallback, addToStreamCallback, changeStatusCallback, getPeers);
		} else if (protocol === 'http:') {
			tcpGetList(infoMovie, torrent, pieces, announce, files, videoToHlsCallback, addToStreamCallback, changeStatusCallback, getPeers);
		}
	})
	dhtGetList(infoMovie, torrent, pieces, files, videoToHlsCallback, addToStreamCallback, changeStatusCallback, getPeers);
};

const dhtGetList = (infoMovie, torrent, pieces, files, videoToHlsCallback, addToStreamCallback, changeStatusCallback, getPeersCallback) => {
		const dht = new DHT();
		dht.listen(20000, () => {});
		dht.on('peer', ( peer, infoHash, from ) => {
			let peerFormat = {
				port: peer.port,
				ip: peer.host
			}
			download(infoMovie, peerFormat, torrent, pieces, files, videoToHlsCallback, addToStreamCallback, changeStatusCallback, getPeers)
		});
		dht.lookup(infoHash(torrent));
		pieces._globalDHTListener = dht;
}

const udpGetList = (infoMovie, torrent, pieces, urlObject, files, videoToHlsCallback, addToStreamCallback, changeStatusCallback, getPeers) => {
	const socket = dgram.createSocket('udp4');
	const buff = buildConnReq();
	try {
		socket.send(buff, 0, buff.length, urlObject.port, urlObject.hostname, () => { });
		pieces._globalUdpTrackerList.push(socket)
		socket.on('message', response => {
			if (respType(response) === 'connect') {
				const connResp = parseConnResp(response);
				const announceReq = buildAnnounceReq(connResp.connectionId, torrent);
				socket.send(announceReq, 0, announceReq.length, urlObject.port, urlObject.hostname, () => { });
			} else if (respType(response) === 'announce') {
				const announceRep = parseAnnounceResp(response);
				announceRep.peers.forEach(peer => download(infoMovie, peer, torrent, pieces, files, videoToHlsCallback, addToStreamCallback, changeStatusCallback, getPeers))
			}
		});
	} catch (error) {}
};

const cleanTcpVar = (content) => {
	return content.replace(/.{2}/g, function (m) {
		var v = parseInt(m, 16);
		if (v <= 127) {
			m = encodeURIComponent(String.fromCharCode(v));
			if (m[0] === '%')
				m = m.toLowerCase();
		} else
			m = '%' + m;
		return m;
	});
}

const tcpGetList = (infoMovie, torrent, pieces, urlObject, files, videoToHlsCallback, addToStreamCallback, changeStatusCallback, getPeers) => {

	const info = bencode.encode(torrent.info);
	let h = crypto.createHash('sha1').update(info).digest('hex');
	h = cleanTcpVar(h)

	let id = genId().toString('hex');
	id = cleanTcpVar(id)

	const params = {
		info_hash: h,
		peer_id: id,
		port: urlObject.port
	};

	const query = `?info_hash=${params.info_hash}&peer_id=${params.peer_id}&port=${params.port}`;

	const opt = {
		hostname: urlObject.hostname,
		path: urlObject.path + query,
		port: urlObject.port
	};

	http.get(opt, res => {
		res.on('data', data => httpparseAnnounceResp(data).forEach((peer) => download(infoMovie, peer, torrent, pieces, files, videoToHlsCallback, addToStreamCallback, changeStatusCallback, getPeers)))
	})
	.on('error', e => {});
}


module.exports = { buildAnnounceReq, buildConnReq, respType, parseConnResp, parseAnnounceResp, getPeers };
