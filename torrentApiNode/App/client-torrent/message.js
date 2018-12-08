const {
	genId,
	infoHash
} = require('./tools.js');

const {
	download
} = require('./download.js');

const buildKeepAlive = () => Buffer.alloc(4)

const buildChoke = () => {
	const buf = Buffer.alloc(5);
	buf.writeUInt32BE(1, 0);
	buf.writeUInt8(0, 4);
	return buf;
}

const buildUnchoke = () => {
	const buf = Buffer.alloc(5);
	buf.writeUInt32BE(1, 0);
	buf.writeUInt8(1, 4);
	return buf;
}



const buildUninterested = () => {
	const buf = Buffer.alloc(5);
	buf.writeUInt32BE(1, 0);
	buf.writeUInt8(3, 4);
	return buf;
}

const buildHave = ( payload ) => {
	const buf = Buffer.alloc(9);
	buf.writeUInt32BE(5, 0);
	buf.writeUInt8(4, 4);
	buf.writeUInt32BE(payload, 5);
	return buf;
}

const buildBitfield = ( payload ) => {
	const buf = Buffer.alloc(14);
	buf.writeUInt32BE(payload.length + 1, 0);
	buf.writeUInt8(5, 4);
	bitfield.copy(buf, 5);
	return buf;
}

const buildPiece = ( payload ) => {
	const buf = Buffer.alloc(payload.block.length + 13);
	buf.writeUInt32BE(payload.block.length + 9, 0);
	buf.writeUInt8(7, 4);
	buf.writeUInt32BE(payload.index, 5);
	buf.writeUInt32BE(payload.begin, 9);
	payload.block.copy(buf, 13);
	return buf;
};

const buildCancel = ( payload ) => {
	const buf = Buffer.alloc(17);
	buf.writeUInt32BE(13, 0);
	buf.writeUInt8(8, 4);
	buf.writeUInt32BE(payload.index, 5);
	buf.writeUInt32BE(payload.begin, 9);
	buf.writeUInt32BE(payload.length, 13);
	return buf;
};

const buildPort = ( payload ) => {
	const buf = Buffer.alloc(7);
	buf.writeUInt32BE(3, 0);
	buf.writeUInt8(9, 4);
	buf.writeUInt16BE(payload, 5);
	return buf;
};

module.exports = {
	buildKeepAlive,
	buildChoke,
	buildUnchoke,
	buildUninterested,
	buildHave,
	buildBitfield,
	buildPiece,
	buildCancel,
	buildPort,
};