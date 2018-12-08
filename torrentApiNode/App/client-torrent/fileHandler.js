const fs = require('fs');

const listContain = ( list, item ) => {
	let check = false;
	list.forEach(element => (element === item) ? check = true : null );
	return check;
}

const getFileStructure = ( torrent ) => {
	let fileStructure = [];
	let currentLen = 0;

	if (torrent.info.files !== undefined) {
		var biggerIndex = null;
		torrent.info.files.forEach(({ length, path }, index) => {
			if (path.length > 1) {
				let directory = '', file = '';
				path.forEach(( buff, index ) => {
					if (index === path.length - 1)
						file = `${directory}/${buff.toString('utf8')}`;
					else
						directory += `/${buff.toString('utf8')}`;
				});

				if (!listContain(fileStructure, file)) {
					const fileToPush = {
						path:   file,
						begin: currentLen === 0 ? 0 : currentLen + 0,
						end: currentLen === 0 ? 0 + length : currentLen  + length,
						length: length,
					}

					if (fileToPush.end > 1)
						fileToPush.end--
					currentLen += length;
					let index = fileStructure.push(fileToPush);
					if (!fileStructure[biggerIndex] || (fileStructure[biggerIndex].length < fileToPush.length))
						biggerIndex = index - 1
				}
			} else {
				const file = path[0].toString('utf8');
				if (!listContain(fileStructure, file)) {
					const fileToPush = {
						path:   file,
						begin: currentLen === 0 ? 0 : currentLen + 0,
						end: currentLen === 0 ? 0 + length : currentLen + 0 + length,
						length
					}
					if (fileToPush.end > 1)
						fileToPush.end--
					currentLen += length;
					let index = fileStructure.push(fileToPush);
					if (!fileStructure[biggerIndex] || 	(fileStructure[biggerIndex] && fileStructure[biggerIndex].length < fileToPush.length))
						biggerIndex = index - 1
				}
			}
		});
		fileStructure[biggerIndex].bigger = true
	}
	else {
		fileStructure.push({
			path: `/${torrent.info.name.toString('utf8')}`,
			begin: currentLen === 0 ? 0 : currentLen + 1,
			end: currentLen === 0 ? 0 + torrent.info.length : currentLen + 0 + torrent.info.length,
			length: torrent.info.length,
			bigger: true
		});
		currentLen += torrent.info.length;
	}
	return fileStructure;
};


const getFd = ( file, toPath ) => {
	const pathElements = file.path.split('/');
	let totalPath = toPath;
	let fd = null;

	pathElements.forEach((element, index) => {
		if (element !== '.' && index !== pathElements.length - 1 && element) {
			totalPath += `/${element}`;
		}
	});

	return fd;
};

module.exports = { getFileStructure, getFd };
