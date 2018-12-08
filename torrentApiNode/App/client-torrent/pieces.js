const { blocksPerPiece, BLOCK_LEN } = require('./tools.js');

var count = 0
class Pieces {
	constructor(torrent) {
		this.torrent = torrent;
		function buildPiecesArray() {
			let piecesNb = torrent.info.pieces.length / 20;
			const arr = new Array(piecesNb).fill(null);
			return arr.map((_, i) => new Array(blocksPerPiece(torrent, i)).fill(false));
		}

		this._requested = buildPiecesArray();
		this._received = buildPiecesArray();
		this._dataBuffer = buildPiecesArray();
		this._dataWritted = 0;
		this._arrayOfData = [];

		this._globalTotalPiecesNumber = torrent.info.pieces.length / 20;
		this._range = 2;
		this._offset = 0;

		function fillRange() {
			let range = [];

			const _range = 2;
			for (let x = 0; x < _range; x++) {
				if (x < (torrent.info.pieces.length / 20))
					range.push(x)
			}

			return range;
		}

		this._globalIndexList = fillRange();
		this._globalPeerList = [];
		this._globalUdpTrackerList = [];
		this._globalDHTListener = null;
		this._peerIpClone = 0;
		this._firstWrite = true;
		this._timeStarted = null;
		this.count = 0;
	}

	addRequested(pieceBlock) {
		const blockIndex = pieceBlock.begin / BLOCK_LEN;
		this._requested[pieceBlock.index][blockIndex] = true;
	}

	addReceived(pieceBlock) {
		const blockIndex = pieceBlock.begin / BLOCK_LEN;
		this._received[pieceBlock.index][blockIndex] = true;
		this._dataBuffer[pieceBlock.index][blockIndex] = pieceBlock.block;
		count++
	}
	isPieceIsReceived(pieceIndex) {
		return this._received[pieceIndex].every(block => block);
	}

	needed(pieceBlock) {
		if (this._requested.every(blocks => blocks.every(i => i))) {
			this._requested = this._received.map(blocks => blocks.slice());
		}
		const blockIndex = pieceBlock.begin / BLOCK_LEN;
		return !this._received[pieceBlock.index][blockIndex];
	}

	neededPiece(pieceIndex) {
		this._received[pieceIndex].forEach(el => {
			if (el) nbBlock++;
		});
		return this._received[pieceIndex].every(block => block);
	}

	isDone() {
		return this._received.every(blocks => blocks.every(i => i));
	}

	getPercentDone() {
		const downloaded = this._received.reduce((totalBlocks, blocks) => {
			return blocks.filter(i => i).length + totalBlocks;
		}, 0);

		const total = this._received.reduce((totalBlocks, blocks) => {
			return blocks.length + totalBlocks;
		}, 0);

		const percent = Math.floor(downloaded / total * 100);
		return percent;
	}

	updatePeerList( socket, bitfield ) {
		let index = null;
		let peerTEMP = null;

		this._globalPeerList.forEach((peer, i) => {
			if (peer.socket === socket)
			{
				peerTEMP = peer
				index = i;
			}
		});

		this._globalPeerList[index].bitfield = bitfield;
		return peerTEMP
	}

	getPeerBitfield( socket ) {
		let index = null;

		this._globalPeerList.forEach((peer, i) => {
			if (peer.socket === socket)
				index = i;
		});

		return this._globalPeerList[index].bitfield ? this._globalPeerList[index].bitfield : null;
	}
}

module.exports = { Pieces };
