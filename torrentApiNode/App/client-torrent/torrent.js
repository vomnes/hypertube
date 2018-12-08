const fs = require('fs');
const path = require('path');
const bencode = require('bencode');
const fetchMetadata = require('bep9-metadata-dl');
const magnet = require('magnet-uri');

const { getFileStructure, getFd } = require('./fileHandler.js');
const { getPeers } = require('./tracker.js');
const { Pieces } = require('./pieces.js');
const { updateMovie } = require('../mongo.js')

const downloadFolder = '../downloads'

require ('wtfnode')

const startTorrent = async (infoMovie, magnetLink, videoToHlsCallback, addToStreamCallback, removeFirstQueueCallback, changeStatusCallback, torrentPath = null) => {
	var torrent;

	if (magnetLink)
	{
		let infoHash = magnet.decode(magnetLink).infoHash
		try {
			changeStatusCallback('Search')
			torrent = await fetchMetadata(infoHash, {  maxConns: 200, fetchTimeout: 30000, socketTimeout: 30000 })
			changeStatusCallback('Connecting')
		}
		catch (error) {
			removeFirstQueueCallback('magnet timeout')
			updateMovie(infoMovie.id, -1, null, false)
			return
		}
	}
	else {
		torrent = bencode.decode(fs.readFileSync(torrentPath));
	}

	const piecesInstance = new Pieces(torrent);
	piecesInstance._timeStarted = new Date();

	let files = [];

	const fileStructure = getFileStructure(torrent);

	if (torrent.info.files !== undefined) {
		const directoryPath = downloadFolder + '/' + torrent.info.name.toString('utf8').replace(/[^a-zA-Z0-9-_\.\\/]/g, '') + '/';
		fileStructure.forEach((fileElement) => {
			files.push({
				fd: null,
				path: fileElement.path.replace(/[^a-zA-Z0-9-_\.\\/]/g, ''),
				begin: fileElement.begin,
				end: fileElement.end,
				length: fileElement.length,
				bigger: fileElement.bigger ? true : false,
				fullPath: directoryPath + fileElement.path.replace(/[^a-zA-Z0-9-_\.\\/]/g, '')
			});
		});
	} else {
		fileStructure.forEach((fileElement) => {
			files.push({
				fd: null, // getFd(fileElement, downloadFolder),
				path: fileElement.path.replace(/[^a-zA-Z0-9-_\.\\/]/g, ''),
				begin: fileElement.begin,
				end: fileElement.end,
				length: fileElement.length,
				bigger: fileElement.bigger ? true : false,
				fullPath: downloadFolder.replace(/[^a-zA-Z0-9-_\.\\/]/g, '') + fileElement.path.replace(/[^a-zA-Z0-9-_\.\\/]/g, '')
			});
		});
	}

	getPeers(infoMovie, torrent, piecesInstance, files, videoToHlsCallback, addToStreamCallback, changeStatusCallback)
}


module.exports = { startTorrent }







// addTorrentQueue("magnet:?xt=urn:btih:55c3737267ac920ed18eca2fb1a94ddf18a12397&dn=Eminem+-+Kamikaze+%282018%29+Mp3+%28320kbps%29+%5BHunter%5D&tr=udp%3A%2F%2Ftracker.leechers-paradise.org%3A6969&tr=udp%3A%2F%2Fzer0day.ch%3A1337&tr=udp%3A%2F%2Fopen.demonii.com%3A1337&tr=udp%3A%2F%2Ftracker.coppersurfer.tk%3A6969&tr=udp%3A%2F%2Fexodus.desync.com%3A6969");

// magnetLink = "magnet:?xt=urn:btih:4b47cdc9f1557ab72f0eb82bfdf375a4920d94a1&dn=Back+Channel+by+Stephen+L.+Carter+EPUB&tr=udp%3A%2F%2Ftracker.leechers-paradise.org%3A6969&tr=udp%3A%2F%2Fzer0day.ch%3A1337&tr=udp%3A%2F%2Fopen.demonii.com%3A1337&tr=udp%3A%2F%2Ftracker.coppersurfer.tk%3A6969&tr=udp%3A%2F%2Fexodus.desync.com%3A6969"
// magnetLink = "magnet:?xt=urn:btih:71702053a0beea9cdc3f4647871f99b15e87b483&dn=Vengeance+in+Venice+by+Philip+Gwynne+Jones+EPUB&tr=udp%3A%2F%2Ftracker.leechers-paradise.org%3A6969&tr=udp%3A%2F%2Fzer0day.ch%3A1337&tr=udp%3A%2F%2Fopen.demonii.com%3A1337&tr=udp%3A%2F%2Ftracker.coppersurfer.tk%3A6969&tr=udp%3A%2F%2Fexodus.desync.com%3A6969"
// magnetLink = "magnet:?xt=urn:btih:ea932e60674f2d01b0f00db7c654d00e93518d4b&dn=The+Better+Angels+of+Our+Nature+by+Steven+Pinker+%28.epub%29%2B&tr=udp%3A%2F%2Ftracker.leechers-paradise.org%3A6969&tr=udp%3A%2F%2Fzer0day.ch%3A1337&tr=udp%3A%2F%2Fopen.demonii.com%3A1337&tr=udp%3A%2F%2Ftracker.coppersurfer.tk%3A6969&tr=udp%3A%2F%2Fexodus.desync.com%3A6969"
// magnetLink = "magnet:?xt=urn:btih:506ed434b58547cae7ab03a0318fb5278a6836ef&dn=My+Thoughts+Exactly+by+Lily+Allen+EPUB&tr=udp%3A%2F%2Ftracker.leechers-paradise.org%3A6969&tr=udp%3A%2F%2Fzer0day.ch%3A1337&tr=udp%3A%2F%2Fopen.demonii.com%3A1337&tr=udp%3A%2F%2Ftracker.coppersurfer.tk%3A6969&tr=udp%3A%2F%2Fexodus.desync.com%3A6969"
// magnetLink = "magnet:?xt=urn:btih:253ceaa1f94c4c0df7deff2f67863d609157621b&dn=Moonlight+%282016%29+1080p+BluRay+-+6CH+-+2GB+-+ShAaNiG&tr=udp%3A%2F%2Ftracker.leechers-paradise.org%3A6969&tr=udp%3A%2F%2Fzer0day.ch%3A1337&tr=udp%3A%2F%2Fopen.demonii.com%3A1337&tr=udp%3A%2F%2Ftracker.coppersurfer.tk%3A6969&tr=udp%3A%2F%2Fexodus.desync.com%3A6969"
// magnetLink = "magnet:?xt=urn:btih:3F2A441E2C4B84F25E44403328AEFFC432A15AE7&dn=Ready%20Player%20One.2018.HDRip.X264.AC3-EVO%5bN1C%5d&tr=udp%3a%2f%2ftracker.leechers-paradise.org%3a6969&tr=udp%3a%2f%2fzer0day.ch%3a1337&tr=udp%3a%2f%2fopen.demonii.com%3a1337&tr=udp%3a%2f%2ftracker.coppersurfer.tk%3a6969&tr=udp%3a%2f%2fexodus.desync.com%3a6969"
// addTorrent(null, '../download', magnetLink);
