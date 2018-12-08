const Hapi = require('hapi');
const jwt = require('jsonwebtoken')
const sock = require('socket.io')

const fetch = require('node-fetch');

const { addTorrentQueue, removeFirstQueue, changeStatus, getStatus, getQueueList, setIO } = require('./client-torrent/torrentQueue.js');
const { videoToHls, addToStream } = require('./hlsTranscoder/hls.js')
const { updateMovie, deleteFailMovie } = require('./mongo.js')

const port = 4002
const server = Hapi.server({
	port,
  host: '0.0.0.0',
});
const jwtSecret = process.env.jwtSecret

server.route({
  config: {
    // cors: {
		//		CROS Handled by Apache
    //     origin: ['http://localhost:8080'],
    //     additionalHeaders: ['cache-control', 'x-requested-with']
    // }
  },
	method: 'POST',
	path: '/torrent/add',
	handler: async (request, h) => {
		if (!request.payload.token || request.payload.token.length < 1)
			return h.response( { message: 'invalid token'}).code(403)

		var torrent;
		try {
			torrent = jwt.verify(request.payload.token, jwtSecret)
		} catch (e) {
			return h.response({ message: 'invalid token' }).code(403)
		}

		let add = await updateMovie(torrent.id, 0, null, null, true)

		if (add && add.code)
			return h.response(add.message).code(add.code)
		else if (add == 1)
			addTorrentQueue(torrent.poster, torrent.duration, torrent.magnet, torrent.title, torrent.id, videoToHls, addToStream)
		else if (add == -1) {
			return h.response({
				'message': 'movie does no exists'
			}).code(403)
		}
		else if (add == -2) {
			return h.response({
				'message': 'movie already downloaded or added to list'
			}).code(403)
		}

		return h.response()
	}
});

server.route({
	config: {
    },
	method: 'GET',
	path: '/torrent/list',
	handler: (request, h) => {
		return getQueueList()
	}
});

const io = sock(server.listener)
setIO(io)
io.on('connect', (socket) => {
})

const init = async () => {
	await server.start();
	console.log(`Server is running on port ${port}...`);
	await deleteFailMovie(true)
};

process.on('unhandledRejection', err => {

});

init();


// Title
require('./jobs.js')
