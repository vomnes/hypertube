const MongoClient = require('mongodb').MongoClient;

const host = process.env.MONGO_DB_HOST;
const url = 'mongodb://' + host +  ':27017';
const dbName = process.env.MONGO_DB_NAME;

const getConnectionCollection = async (ignoreError=false) => {
	try {
		let client = await MongoClient.connect(url, { useNewUrlParser: true });
		return { client, collection : client.db(dbName).collection('movies') }
	}
	catch (e) {
			return null;
	}
}
const updateMovie = async (id, status=null, path=null, stream=null, add=false, subsInfo=null) => {
	try {
		let connection = await getConnectionCollection();
		if (!connection)
			return { code: 500, message: { "message": "unable to connect to database" } }
		let exist = await connection.collection.find({ _id : id }, {}).limit(1).toArray()
		if (!exist[0])
		{
			connection.client.close()
			return -1
		}
		else {
			let exist = await connection.collection.find({ _id: id, 'video.status': { $gt: -1 } }, {}).limit(1).toArray()
			if (exist[0] && add) {
				connection.client.close()
				return -2
			}
			let update = {
				video: {
				}
			};
			if (path)
				await connection.collection.updateOne({ _id: id }, { $set: { 'video.path': path } })
			if (stream)
				await connection.collection.updateOne({ _id: id }, { $set: { 'video.stream': stream } })
			if (status !== null)
				await connection.collection.updateOne({ _id: id }, { $set: { 'video.status': status } })
			if (subsInfo)
			{
				await connection.collection.updateOne({ _id: id }, { $set: { 'video.subs': subsInfo } })
			}
			await connection.client.close()
			return 1
		}
	}
	catch (e) {
	}
}

const deleteFailMovie = async (autoretry = false) => {
	let connection = await getConnectionCollection(true);
	if (!connection && autoretry == false)
		return { code: 500, message: { "message": "unable to connect to database" } }
	else if (!connection && autoretry == true)
	{
		while (!connection)
		{
			connection = await getConnectionCollection(true);
		}
	}

	try {
		let update = {
			video: {
				status: -1,
				path: null,
				stream: false
			}
		}
		await connection.collection.updateMany({ 'video.status': { $gt: -1 } }, {$set: update})
		connection.client.close();
	}
	catch (e) {
	}

}

module.exports = { getConnectionCollection, updateMovie, deleteFailMovie };
