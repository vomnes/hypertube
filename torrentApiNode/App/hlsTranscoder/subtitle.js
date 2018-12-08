const OS = require('opensubtitles-api');
const { spawn } = require('child_process');
var srt2vtt = require('srt-to-vtt')
var fs = require('fs')
var request = require('request');


const getSubs = async ( imdbid, folder ) => {
	try {
		let subsInfo = []
		const OpenSubtitles = new OS({
			useragent: "SolEol v1.21",
			ssl: true
		});

		let subs = await OpenSubtitles.search({
			imdbid,
			sublanguageid: "fre,eng,ita",
			extensions: ['srt']
		 })

		if (subs.en)
		{
			saveSub(subs.en.url, folder, subs.en.langcode)
			subsInfo.push({
				lang: subs.en.langcode,
				path: `${folder.split("/torrent/public").pop()}/${subs.en.langcode}.vtt`
			})
		}
		if (subs.fr)
		{
			saveSub(subs.fr.url, folder, subs.fr.langcode)
			subsInfo.push({
				lang: subs.fr.langcode,
				path: `${folder.split("/torrent/public").pop()}/${subs.fr.langcode}.vtt`
			})
		}
		if (subs.it)
		{
			saveSub(subs.it.url, folder, subs.it.langcode)
			subsInfo.push({
				lang: subs.it.langcode,
				path: `${folder.split("/torrent/public").pop()}/${subs.it.langcode}.vtt`
			})
		}
		return subsInfo
	}
	catch (error) {
		return -1
	}
}

const saveSub = async (path, folder, name) => {
	request(path)
		.pipe(srt2vtt())
		.pipe(fs.createWriteStream(`${folder}/${name}.vtt`))
}

const downloadSubs = async (imdbid, folder) => {
	let countTry = 0;
	let subsInfo = await getSubs(imdbid, folder)
	while (subsInfo == -1 && countTry < 30)
	{
		subsInfo = await getSubs(imdbid, folder)
		countTry++
	}
	return subsInfo
}

module.exports = { downloadSubs }
