<template lang="html">
  <div @click.self="toggleVisiblity" class="black-back">
    <div class="box">
      <h1 class="title">{{ $t('sourcesLabel') }}</h1>
      <img :title="$t('close')" @click="toggleVisiblity"
           alt="close" class="close" src="../../../static/icons/close.svg">
      <div class="list">
        <div :key="index" class="torrent" v-for="(torrent, index) in list">
          <h1>{{ formatString(torrent.title) }}</h1>
          <img
            @click="downloadTorrent(torrent, index)"
            alt="download"
            src="../../../static/icons/download.svg"
            title="Download through this torrent"
          >
          <div class="details">
            <div style="backgroundColor: #65E3B9">
              <p>Seed: {{ torrent.seed }}</p>
            </div>
            <div style="backgroundColor: #9461F7; marginLeft: 10px; marginRight: 10px">
              <p>Leech: {{ torrent.peer - torrent.seed }}</p>
            </div>
            <div style="backgroundColor: #4A4350">
              <p>{{ Math.round(torrent.size/1000000) }} mb</p>
            </div>
          </div>
        </div>
      </div>
      <p class="no-torrents" v-if="noData">{{ $t('noTorrents') }}</p>
      <loading-animation class="loader" v-if="isLoading"/>
    </div>
  </div>
</template>

<script>
import swal from 'sweetalert';
import api from '../../lib/api';
import auth from '../../../static/js/auth';
import LoadingAnimation from './../module/LoadingAnimation';

export default {
  props: {
    movieID: String,
    change: Function,
  },
  components: {
    LoadingAnimation,
  },
  data() {
    return {
      isLoading: false,
      noData: false,
      list: [],
    };
  },
  methods: {

    toggleVisiblity() {
      this.$emit('closeTorrents');
    },

    getTorrentsMagnets() {
      this.isLoading = true;
      api.torrents({
        imdb: this.movieID,
      },
      auth.token())
        .then((res) => {
          if (res.status !== 200) {
            throw new Error('Error status ', res.status);
          } else {
            res.json()
              .then((data) => {
                if (data) {
                  if (data.message === 'No torrents') {
                    this.noData = true;
                  } else {
                    this.noData = false;
                    this.list = data;
                  }
                } else {
                  this.noData = true;
                }
                this.isLoading = false;
              });
          }
        })
        .catch(() => {
          this.isLoading = false;
          this.nodata = true;
        });
    },

    formatString(str) {
      if (!str || str.length < 55) return str;
      return `${str.substring(0, 55)}...`;
    },

    // Download the choosen torrent
    downloadTorrent(torrent, index) {
      this.status = true;
      this.toggleVisiblity();
      api.downloadMovies(torrent.jwt)
        .then((res) => {
          if (res.status === 200) {
            this.change();
            this.$_bus.$emit('downloadMenuShow');
          } else {
            throw new Error('Cannot download this movie')
          }
        })
        .catch(error => {
          swal({ text: this.$t('errorServer'), icon: 'error' })
        });
    },
  },
  mounted() {
    this.noData = false;
    this.getTorrentsMagnets();
  },
};
</script>

<style scoped lang="sass">
  @import "Over"
</style>
