const fs = require('fs');
const { spawn } = require('child_process');
const { Readable } = require('stream');

const { addTorrentQueue, removeFirstQueue, changeStatus, getStatus, getQueueList } = require('../client-torrent/torrentQueue.js');
const { cleanConnections } = require('../client-torrent/handlerDownload.js');
const { updateMovie } = require('../mongo.js');
const { downloadSubs } = require('./subtitle.js');
const moviePath = '/movies/'
const movieDir = "/torrent/public" + moviePath

var piecesCurrent = null;
var filesCurrent = null;

var buffStream = null;

const addToStream = (buffer) => {
	buffStream.push(buffer)
}

const videoToHls = (infoMovie, fromPath, buffer, folderName, pieces, files, changeStatusCallback) => {
	try {

		var start = new Date()
		var stderr = null

		piecesCurrent = pieces;
		filesCurrent = files;

		buffStream = new Readable({
			read() {}
		})

		buffStream.push(buffer)

		if (!fs.existsSync(movieDir))
			fs.mkdirSync(movieDir);

		if (!fs.existsSync(movieDir + folderName))
			fs.mkdirSync(movieDir + folderName);

			const ffmpeg = spawn('ffmpeg', ['-i', '-', '-hls_list_size', '0',
				'-f', 'hls', movieDir + folderName + '/video.m3u8'
			]);

		buffStream.pipe(ffmpeg.stdin)

		ffmpeg.on('exit', (code, signal) => {
			if (code == 0)
			{
				let time = new Date() - start;
				updateMovie(infoMovie.id, 2)
				changeStatus('Done')
				removeFirstQueue()
			}
			else
			{
				cleanConnections(infoMovie, piecesCurrent, filesCurrent, false);
				piecesCurrent = null;
				filesCurrent = null;
				changeStatusCallback('Error')
				removeFirstQueue('transcoding error')
			}
		})

		ffmpeg.stderr.on('data', (data) => {
			try {
				let regex = /time=([0-9][0-9]):([0-9][0-9]):([0-9][0-9]).([0-9][0-9])/;
				let time = regex.exec(data)
				if (time)
				{
					let minute = 0
					minute += parseFloat(time[1] * 60);
					minute += parseFloat(time[2]);
					minute += parseFloat(time[3] / 60);
					minute += parseFloat(time[4] / 60000);
					minute = parseFloat(minute).toFixed(2)

					let pourcentage = Math.floor((minute / infoMovie.duration * 100));
					let download = getStatus()
					if (download.status == 'Transcode' || (download.status == 'Download' && download.percentage == 100))
						changeStatus('Transcode', pourcentage)
					if (!infoMovie.first)
					{
						if (fs.existsSync(movieDir + folderName + '/video5.ts'))
						{
							infoMovie.first = true;
							updateMovie(infoMovie.id, 1, moviePath + folderName + '/video.m3u8', true)
							changeStatus(null, null, moviePath + folderName + '/video.m3u8', true)
						}
					}
				}

				stderr = data
			}
			catch (error) {
			}
		});

		ffmpeg.stdout.on('data', (data) => {
		});

		ffmpeg.stdin.on('error', (data) => {
		});

	}
	catch (error) {
	}
}

module.exports = { videoToHls, addToStream }
