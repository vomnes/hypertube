const fs = require('fs');


const { updateMovie } = require('../mongo.js')

const chokeHandler = (socket) => {
	socket.end();
};

const unchokeHandler = ( socket, pieces, queue ) => {
	queue.choked = false;
	requestPiece(socket, pieces, queue);
};

const haveHandler = ( torrent, socket, pieces, queue, payload ) => {
	const pieceIndex = payload.readUInt32BE(0);
	let bitfield = pieces.getPeerBitfield(socket)
	if (typeof bitfield === null)
		bitfield = new Array(pieces._globalTotalPiecesNumber).fill(0)
	if (bitfield && bitfield[pieceIndex])
		bitfield[pieceIndex] = 1;

	pieces.updatePeerList(socket, bitfield)
};

const requestPiece = ( socket, pieces, queue ) => {
	if (queue.choked)
		return null;
	while (queue.length()) {
		const pieceBlock = queue.deque();
		if (pieces.needed(pieceBlock)) {
			socket.write(buildRequest(pieceBlock));
			pieces.addRequested(pieceBlock);
			break;
		}
	}
}

const bitfieldHandler = ( socket, pieces, queue, payload ) => {
	let bitfield = new Array(pieces._globalTotalPiecesNumber).fill(0)
	payload.forEach(( byte, i ) => {
		for (let j = 0; j < 8; j++) {
			if (byte % 2)
				bitfield[i * 8 + 7 - j] = 1
			byte = Math.floor(byte / 2);
		}
	});
	let peer = pieces.updatePeerList(socket, bitfield);
	askPieces(peer, pieces);
};

var timeoutforData = null;
var autoRetryCount = 0;
var stopError = false;

const timeoutDownloadHandler = (torrent, infoMovie, pieces, files, videoToHlsCallback, addToStreamCallback, changeStatusCallback, getPeers ) => {
	if ((pieces && pieces.isDone && pieces.isDone()) || stopError)
		return
	autoRetryCount++
	if(pieces) {
		cleanConnections(infoMovie, pieces, files)
		getPeers(infoMovie, torrent, pieces, files, videoToHlsCallback, addToStreamCallback, changeStatusCallback)
	}

	clearTimeout(timeoutforData)
	timeoutforData = setTimeout(timeoutDownloadHandler, 10000, infoMovie, pieces, files, videoToHlsCallback, addToStreamCallback, getPeers);
}

const pieceHandler = (infoMovie, socket, pieces, queue, torrent, files, pieceResp, videoToHlsCallback, addToStreamCallback, changeStatusCallback, getPeers) => {
	stopError = false;

	if (timeoutforData !== null)
		clearTimeout(timeoutforData)

	timeoutforData = setTimeout(timeoutDownloadHandler, 5000, torrent, infoMovie, pieces, files, videoToHlsCallback, addToStreamCallback, changeStatusCallback, getPeers);
	process.stdout.write('progress: ' + pieces.getPercentDone() + '%\r');

	changeStatusCallback('Download', pieces.getPercentDone())
	if (pieces.isDone())
	{
		socket.end()
		return;
	}
	pieces.addReceived(pieceResp);

	let indexOfReceivedPiece = pieces._globalIndexList.indexOf(pieceResp.index);

	if (indexOfReceivedPiece >= 0 && pieces.isPieceIsReceived(pieceResp.index))
	{
		pieces._arrayOfData.push({ index: pieceResp.index, data: Buffer.concat(pieces._dataBuffer[pieceResp.index]) });
		pieces._dataBuffer[pieceResp.index] = [];
		pieces._globalIndexList.splice(indexOfReceivedPiece, 1);

		if (pieces._globalIndexList.length === 0) {
			let dataArray = [];
			pieces._arrayOfData
				.sort((a, b) => a.index - b.index)
				.forEach(({ data }) => dataArray.push(data));

			let buffToPush = Buffer.concat(dataArray);

			files.forEach(({ path, fd, begin, end, length, bigger, fullPath, filename}) => {
				if (end >= pieces._dataWritted && begin <= pieces._dataWritted + buffToPush.length) {
					let startOffset = 0;
					while (begin > pieces._dataWritted && begin > startOffset && startOffset < buffToPush.length && begin > (pieces._dataWritted + startOffset))
						startOffset++;

					let endOffset = startOffset;
					while ((endOffset - startOffset + pieces._dataWritted) < end && endOffset < end && endOffset < buffToPush.length && (endOffset - startOffset + 1) < length)
						endOffset++;

						const newBuff = buffToPush.slice(startOffset, endOffset + 1);
						if ( bigger && pieces._firstWrite) {
							pieces._firstWrite = false;
							let filename = infoMovie.folderName;
							videoToHlsCallback(infoMovie, fullPath, newBuff, filename, pieces, files, changeStatusCallback)
						}
						else if (bigger)
							addToStreamCallback(newBuff)
				}
			});

			pieces._dataWritted += buffToPush.length;
			pieces._offset += pieces._range;

			for (let x = pieces._offset; x < pieces._offset + pieces._range; x++) {
				if (x < pieces._globalTotalPiecesNumber)
					pieces._globalIndexList.push(x)
			}

			pieces._arrayOfData = [];
		}
		pieces._globalPeerList.forEach(element => {
			askPieces(element, pieces)
		});
	}

	if (pieces.isDone()) {
		let time = new Date() - pieces._timeStarted;
		clearTimeout(timeoutforData)
		timeoutforData = -1
		addToStreamCallback(null)
		changeStatusCallback('Download', 100)
		cleanConnections(infoMovie, pieces, files)

		console.log('DONE with ', pieces._globalPeerList.length, " peers, ", pieces._peerIpClone, " ip duplication,", autoRetryCount, "auto retry in", Math.round(time / 60000), "minutes !");
	} else {
		requestPiece(socket, pieces, queue);
	}
}

const askPieces = (peerElement, pieces) => {
	pieces._globalIndexList.forEach((element) => {
		if (peerElement.bitfield && peerElement.bitfield[element] === 1 && !peerElement.queue.length()) {
			peerElement.queue.queue(element);
		}
	});
	requestPiece(peerElement.socket, pieces, peerElement.queue);
}

const cleanConnections = (infoMovie, pieces, files, finish = true) => {
	if (pieces._globalPeerList && pieces._globalPeerList.length > 0)
	pieces._globalPeerList.forEach(({socket}) => socket.destroy());
	pieces._globalPeerList = []
	if (pieces._globalUdpTrackerList && pieces._globalUdpTrackerList.length > 0)
		pieces._globalUdpTrackerList.forEach((socket) => socket.close());
	pieces._globalUdpTrackerList = []
	if (pieces._globalDHTListener)
		pieces._globalDHTListener.destroy();
	try {
		if (files && files.length > 0)
			files.forEach(({ fd}) => fd != null ? fs.closeSync(fd) : null);
	} catch (e) { }
	if (!finish)
	{
		console.log("download fail go delete in db")
		stopError = true
		updateMovie(infoMovie.id, -1, null, false)
	}

}


const buildRequest = (payload) => {
	const buf = Buffer.alloc(17);
	buf.writeUInt32BE(13, 0);
	buf.writeUInt8(6, 4);
	buf.writeUInt32BE(payload.index, 5);
	buf.writeUInt32BE(payload.begin, 9);
	buf.writeUInt32BE(payload.length, 13);
	return buf;
};

module.exports = {
	haveHandler,
	requestPiece,
	pieceHandler,
	bitfieldHandler,
	unchokeHandler,
	chokeHandler,
	cleanConnections,
};
