const fs = require('fs');

const { startTorrent } = require('./torrent.js')
const { updateMovie } = require('../mongo.js');
const { downloadSubs } = require('../hlsTranscoder/subtitle.js');
// const { videoToHls, addToStream } = require('../hlsTranscoder/hls.js')

const moviePath = '/movies/'
const movieDir = "/torrent/public" + moviePath

var torrentQueue = []
var io = null

const setIO = (new_io) => io = new_io

const addTorrentQueue = (poster, duration, magnetLink, title, id, videoToHlsCallback, addToStreamCallback, torrentPath = null) => {
	var exist = torrentQueue.some( (el) => el.id === id );
	if (exist)
		return
	torrentQueue.push({
		title,
		id,
		poster,
		magnetLink,
		torrentPath,
		videoToHlsCallback,
		addToStreamCallback,
		duration,
		status: 'Waiting',
		percentage: -1,
		stream: false,
		path: null,
	})
	io.emit('add', { title, id, duration, status: 'Waiting', percentage: -1, stream: false, path: null, poster })
	if (torrentQueue.length < 2)
		handleQueue()
}

const handleQueue = () => {
	if (torrentQueue[0]) {
		let info = {
			title : torrentQueue[0].title,
			id : torrentQueue[0].id,
			duration : torrentQueue[0].duration,
			folderName: torrentQueue[0].title.replace(/[\s,/\\]+/g, "")
		}
		startTorrent(info, torrentQueue[0].magnetLink, torrentQueue[0].videoToHlsCallback, torrentQueue[0].addToStreamCallback, removeFirstQueue, changeStatus, torrentQueue[0].torrentPath)
		if (!fs.existsSync(movieDir))
			fs.mkdirSync(movieDir);

		if (!fs.existsSync(movieDir + info.folderName))
			fs.mkdirSync(movieDir + info.folderName);

		downloadSubs(info.id, movieDir + info.folderName).then((subsInfo) => {
			updateMovie(info.id, null, null, null, null, subsInfo)
		})
	}
}

const removeFirstQueue = (error = null) => {
	io.emit('delete', { id : torrentQueue[0].id, error })
	torrentQueue.splice(0, 1)
	handleQueue()
}

const changeStatus = (status = null, percentage = null, path= null, streamStart = null) => {
	let sendEvent = false
	if ((percentage != null && percentage != torrentQueue[0].percentage) || percentage == -1 || (streamStart != null && torrentQueue[0].stream != streamStart) || (status != null && torrentQueue[0].status != status))
		sendEvent = true
	if (torrentQueue[0] && status != null)
		torrentQueue[0].status = status
	if (torrentQueue[0] && percentage != null)
		torrentQueue[0].percentage = percentage
	if (torrentQueue[0] && path != null)
		torrentQueue[0].path = path
	if (torrentQueue[0] && streamStart == true)
		torrentQueue[0].stream = true

	if (sendEvent) {
			io.emit('update', {
				title: torrentQueue[0].title,
				id: torrentQueue[0].id,
				duration: torrentQueue[0].duration,
				status: torrentQueue[0].status,
				percentage: torrentQueue[0].percentage,
				stream: torrentQueue[0].stream,
				path: torrentQueue[0].path,
				poster: torrentQueue[0].poster
			})
	}
}

const getStatus = () => {
	return {
		status: torrentQueue[0].status,
		percentage: torrentQueue[0].percentage
	}
}

const getQueueList = () => {
	return torrentQueue
}

module.exports = {
	addTorrentQueue,
	removeFirstQueue,
	changeStatus,
	getStatus,
	getQueueList,
	setIO
};
