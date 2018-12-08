const net = require('net');

const { chokeHandler, unchokeHandler, haveHandler, bitfieldHandler, pieceHandler } = require('./handlerDownload.js');
const { Queue } = require('./queue.js')
const {
	genId,
	infoHash
} = require('./tools.js');

const download = (infoMovie, peer, torrent, pieces, files, videoToHlsCallback, addToStreamCallback, changeStatusCallback, getPeers) => {
	const {
		port,
		ip
	} = peer;
	if (peerExist(port, ip, pieces._globalPeerList))
		pieces._peerIpClone++
	const socket = net.Socket();

	socket.setTimeout(10000);
	socket.on('timeout', () => {
		socket.destroy();
	});

	socket.on('error', () => {
		socket.destroy();
	});
	socket.connect(port, ip, () => {
		socket.write(buildHandshake(torrent));
		const queue = new Queue(torrent);
		onWholeMsg(infoMovie, socket, pieces, queue, files, torrent, peer, videoToHlsCallback, addToStreamCallback, changeStatusCallback, getPeers);
	});
};

const onWholeMsg = (infoMovie, socket, pieces, queue, files, torrent, peer, videoToHlsCallback, addToStreamCallback, changeStatusCallback, getPeers) => {
	let savedBuf = Buffer.alloc(0);
	let handshake = true;

	socket.on('data', (recvBuf) => {
		const msgLen = () => handshake ? savedBuf.readUInt8(0) + 49 : savedBuf.readInt32BE(0) + 4;
		savedBuf = Buffer.concat([savedBuf, recvBuf]);
		while (savedBuf.length >= 4 && savedBuf.length >= msgLen()) {
			let msg = savedBuf.slice(0, msgLen());
			msgHandler(infoMovie, msg, socket, pieces, queue, files, torrent, peer, videoToHlsCallback, addToStreamCallback, changeStatusCallback, getPeers)
			savedBuf = savedBuf.slice(msgLen());
			handshake = false;
		}
	})
};

const msgHandler = (infoMovie, msg, socket, pieces, queue, files, torrent, peer, videoToHlsCallback, addToStreamCallback, changeStatusCallback, getPeers) => {
	if (isHandshake(msg)) {
		socket.write(buildInterested());
		if (peerExist(peer.port, peer.ip, pieces._globalPeerList))
			pieces._peerIpClone++
		pieces._globalPeerList.push({
			port: peer.port,
			ip: peer.ip,
			socket,
			queue
		});
	} else {
		const m = parseMsg(msg);
		if (m.id === 0) chokeHandler(socket);
		if (m.id === 1) unchokeHandler(socket, pieces, queue);
		if (m.id === 4) haveHandler(torrent, socket, pieces, queue, m.payload);
		if (m.id === 5) bitfieldHandler(socket, pieces, queue, m.payload);
		if (m.id === 7) pieceHandler(infoMovie, socket, pieces, queue, torrent, files, m.payload, videoToHlsCallback, addToStreamCallback, changeStatusCallback, getPeers);
	}
};

const isHandshake = (msg) => {
	return msg.length === msg.readUInt8(0) + 49 &&
		msg.toString('utf8', 1, 20) === 'BitTorrent protocol';
};

const peerExist = (ip, port, list) => {
	if (list)
	{
		for (let i = 0; i < list.length; i++) {
			if (list[i] && list[i].ip && list[i].port && list[i].ip === ip && list[i].port === port)
				return true
		}

	}
	return false
}

const parseMsg = (msg) => {
	const id = msg.length > 4 ? msg.readInt8(4) : null;
	let payload = msg.length > 5 ? msg.slice(5) : null;

	if (id === 6 || id === 7 || id === 8) {
		const rest = payload.slice(8);
		payload = {
			index: payload.readInt32BE(0),
			begin: payload.readInt32BE(4)
		};
		payload[id === 7 ? 'block' : 'length'] = rest;
	}

	return {
		size: msg.readInt32BE(0),
		id: id,
		payload: payload
	};
}

const buildHandshake = (torrent) => {
	const buf = Buffer.alloc(68)

	buf.writeUInt8(19, 0)
	buf.write('BitTorrent protocol', 1)

	buf.writeUInt32BE(0, 20)
	buf.writeUInt32BE(0, 24)

	infoHash(torrent).copy(buf, 28)
	genId().copy(buf, 48);
	return buf;
}

const buildInterested = () => {
	const buf = Buffer.alloc(5);
	buf.writeUInt32BE(1, 0);
	buf.writeUInt8(2, 4);
	return buf;
}


module.exports = { download };
