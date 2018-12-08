FROM git.harknet.ml:5005/becrespi/hypertube/ffmpeg
# RUN echo "deb http://www.deb-multimedia.org jessie main non-free" >> /etc/apt/sources.list
# RUN echo "deb-src http://www.deb-multimedia.org jessie main non-free" >> /etc/apt/sources.list
# RUN echo "deb http://httpredir.debian.org/debian/ jessie-backports main" >> /etc/apt/sources.list
# RUN apt-get update && apt-get install deb-multimedia-keyring -y --force-yes && apt-get update
# RUN apt-get install ffmpeg -y
EXPOSE 4002
RUN mkdir -p /torrent/public/movies
RUN mkdir -p /torrent/api/
COPY ./App/ /torrent/api/
RUN npm i  --prefix /torrent/api
RUN apt-get update --fix-missing && apt-get install cron -y --fix-missing
ADD ./crontab /etc/cron.d/update-db-cron
RUN crontab /etc/cron.d/update-db-cron
ENTRYPOINT cron && node /torrent/api/server.js
