const crypto = require('crypto');
const bignum = require('bignum');
const bencode = require('bencode');

let id = null;

const genId = () => {
	if (!id) {
		id = crypto.randomBytes(20);
		Buffer.from('-ND0001-').copy(id, 0);
	}
	return id;
};

const infoHash = torrent => {
	const info = bencode.encode(torrent.info);
	return crypto.createHash('sha1').update(info).digest();
};

const size = torrent => {
	let size = torrent.info.files ?
		torrent.info.files.map(file => file.length).reduce((a, b) => a + b) :
		torrent.info.length;

	return bignum.toBuffer(size, {
		size: 8
	});
};

const BLOCK_LEN = Math.pow(2, 14)

const pieceLen = (torrent, pieceIndex) => {
	const totalLength = bignum.fromBuffer(size(torrent)).toNumber();
	const pieceLength = torrent.info['piece length'];
	const lastPieceLength = totalLength % pieceLength;
	const lastPieceIndex = Math.floor(totalLength / pieceLength)
	return lastPieceIndex === pieceIndex ? lastPieceLength : pieceLength;
}

const blocksPerPiece = (torrent, pieceIndex) => {
	const pieceLength = pieceLen(torrent, pieceIndex);
	return Math.ceil(pieceLength / BLOCK_LEN)
}

const blockLen = (torrent, pieceIndex, blockIndex) => {
	const pieceLength = pieceLen(torrent, pieceIndex);
	const lastPieceLength = pieceLength % BLOCK_LEN;
	const lastPieceIndex = Math.floor(pieceLength / BLOCK_LEN);

	return blockIndex === lastPieceIndex ? lastPieceLength : BLOCK_LEN;
}

module.exports = { genId, infoHash, size, BLOCK_LEN, pieceLen, blocksPerPiece, blockLen };