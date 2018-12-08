FROM node
RUN echo "deb http://www.deb-multimedia.org jessie main non-free" >> /etc/apt/sources.list
RUN echo "deb-src http://www.deb-multimedia.org jessie main non-free" >> /etc/apt/sources.list
RUN echo "deb http://httpredir.debian.org/debian/ jessie-backports main" >> /etc/apt/sources.list
RUN apt-get update && apt-get install deb-multimedia-keyring -y --force-yes && apt-get update
RUN apt-get install ffmpeg -y
# RUN apt-get update
# RUN apt-get install -y git build-essential gcc make yasm autoconf automake cmake libtool checkinstall wget software-properties-common pkg-config libmp3lame-dev libunwind-dev zlib1g-dev
# RUN wget https://www.ffmpeg.org/releases/ffmpeg-4.0.2.tar.gz
# RUN tar -xzf ffmpeg-4.0.2.tar.gz
# RUN cd ffmpeg-4.0.2 && ./configure --enable-gpl --enable-libmp3lame --enable-decoder=mjpeg,png --enable-encoder=png && make && make install