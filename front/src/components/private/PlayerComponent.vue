<template>
  <div id="modal-video-player" @click.self="$emit('updatePlayerVue')">
    <video
      controls="true"
      crossorigin="true"
      id="video-player"
      autoplay="true"
    >
      <track v-for="(el, index) in subs" :key="index" :src="el.path" :srclang="el.lang"/>
    </video>
  </div>
</template>

<script>
import Hls from 'hls.js';
import api from '../lib/api';
import auth from '../../static/js/auth'

export default {
  props: {
    moviePath: String,
    sub: Array
  },

  data() {
    return {
      subs: []
    }
  }
  ,

  mounted() {
    const videoPlayer = document.getElementById('video-player');
    const source = process.env.API_DOMAIN_NAME + this.moviePath;
    if (this.sub) {
      this.subs = this.sub.map(item => {
        let lang = 'En'
        if (item.path.includes('en.vtt')) lang = 'En'
        if (item.path.includes('fr.vtt')) lang = 'Fr'
        if (item.path.includes('it.vtt')) lang = 'It'

        return {
          lang,
          path: process.env.API_DOMAIN_NAME + item.path
        }
      });
    }


    if (Hls.isSupported()) {
      const hls = new Hls();
      hls.loadSource(source);
      hls.attachMedia(videoPlayer);
      hls.on(Hls.Events.MANIFEST_PARSED, () => {
        videoPlayer.play().catch(() => null);
        this.updateViewMovie();
      });
    } else if (videoPlayer.canPlayType('application/vnd.apple.mpegurl')) {
      videoPlayer.src = source;
      videoPlayer.addEventListener('loadedmetadata', () => {
        videoPlayer.play().catch(() => null);
        this.updateViewMovie();
      });
    }
  },

  methods: {
    updateViewMovie() {
      api.updateViewMovie(this.$route.params.id, auth.token())
        .then()
        .catch(() => null);
    },
  },
};
</script>

<style scoped lang="sass">
  @import "PlayerComponent.sass"
</style>
