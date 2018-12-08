const { blocksPerPiece, BLOCK_LEN, blockLen } = require('./tools.js');

class Queue {
	constructor(torrent) {
		this._torrent = torrent;
		this._queue = [];
		this.choked = true;
		this.pieceQueuLength = []
		this.dequeNumber = 0;
		this.pieceFinish = false;
	}

	queue(pieceIndex) {
		const nBlocks = blocksPerPiece(this._torrent, pieceIndex);
		this.pieceQueuLength.push({ index: pieceIndex, blockNumber: nBlocks})
		for (let i = 0; i < nBlocks; i++) {
			const pieceBlock = {
				index: pieceIndex,
				begin: i * BLOCK_LEN,
				length: blockLen(this._torrent, pieceIndex, i)
			};
			this._queue.push(pieceBlock);
		} 
	}

	deque() {
		this.dequeNumber += 1
		if (this.pieceQueuLength[0] && this.pieceQueuLength[0].blockNumber === this.dequeNumber)
		{
			this.pieceFinish = true
			this.dequeNumber = 0;
			this.pieceQueuLength.shift()
		}
		return this._queue.shift();
	};

	havePieceInQueue(pieceIndex) {
		let have = false;
		this.pieceQueuLength.forEach((element) => {
			if (element.index === pieceIndex)
				have = true
		})
		return have
	}
	
	peek() {
		return this._queue[0];
	};
	
	length() {
		return this._queue.length;
	};
}

module.exports = { Queue }