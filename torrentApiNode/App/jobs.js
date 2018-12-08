const { getConnectionCollection } = require('./mongo');
const path = require('path');
const fs = require('fs');
const exec = require('child_process').exec;

const init = async () => {
  let connection = null;
  while (!connection)
    connection = await getConnectionCollection(true)

  console.log('INITIALIZE CRON JOB')

  var date = new Date();
  date.setDate(date.getDate() - 30);

  const data = await connection.collection.find({ 'video.status': 2 }).toArray()

  if (data.length === 0) {
    console.log('Nothing to clean');
    return;
  } else {
    console.log(`${data.length} movies cleaned, they are too old`);
  }


  let moviesToClean = [];
  data.forEach(item => {
    let toClean = false;
    item.watchedby.forEach(user => {
      if (new Date(user.watchedat) < date)
        toClean = true;
    });

    if (toClean)
      moviesToClean.push(item);
  });

  console.log(moviesToClean.length, ' items to clean')

  moviesToClean.forEach(item => {
    exec('rm -rf ' + '/torrent/public/movies/' + item.video.path.split('/')[2], (err, stdout, stderr) => {
        console.log('FILM DELETED');
        connection.collection.updateOne({ _id: item._id }, { $set: { 'video.status': -1, 'video.stream': false, 'video.path': null }})
    });
  });
}

init();