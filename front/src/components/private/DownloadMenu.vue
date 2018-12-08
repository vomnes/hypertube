<template>
  <div>
    <img
      @click="toogleDownloadMenu"
      alt="download"
      class="download"
      src="../../static/icons/download.svg"
      title="Downloading movies"
    >

    <div v-if="showMenuDownload">
      <div class="header-info">
        <div class="header-info-loading" v-if="loading">
          <loading-animation class="header-info-loading-animation"></loading-animation>
        </div>
        <div class='menu-content' v-else>
          <div class="header-info-header" style="margin-bottom: 20px">
            <div class="header-info-title">
              <p class="header-info-title-1">{{$t('downloadLabel')}}</p>
              <p class="header-info-title-2">({{queue.length}})</p>
            </div>
          </div>

          <div class="download-container" v-if="queue.length > 0">
            <div style="display:flex;">
              <router-link :to="'/movie/' + queue[0].id">
                <img :src="resizeImage(queue[0].poster)" class="film-picture">
              </router-link>
              <div class="data-container">
                <p class="data-p">{{queue[0].title}}</p>
                <div class="bar-container">
                  <!-- `Download` || `Search` || `Connecting` ou encore `Error`(mais qui reste pas longtemps car du coup il est virée de la file) -->
                  <div :style="{ width: queue[0].percentage + '%'}" :class="{ bar: queue[0].status === 'Download', 'bar-hls': queue[0].status === 'Transcode' }"></div>
                  <div class="bar-text">{{ getTextInfo(queue[0]) }}</div>
                </div>
              </div>
            </div>
          </div>
          <div class="empty-download" v-else>
            {{$t('emptyDlLabel')}}
          </div>


          <!-- THERE -->
          <div class="film-display" v-if="queue.length > 1">
            <show-pictures
              :list="queue.map(({ id, title, poster, }) =>
              ({  id, title, url: poster })).slice(1, queue.length)"
              :size="'tiny'"
              :title="'empty'"
            >
            </show-pictures>
          </div>

          <div @click="toogleDownloadMenu" class="close-btn">
            <img class="close-btn-icon" src='../../static/icons/arrow.svg'/>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import swal from 'sweetalert';
import api from '../lib/api';
import loadingAnimation from './module/LoadingAnimation';
import ShowPictures from './ShowPictures';

export default {
  components: {
    loadingAnimation,
    ShowPictures,
  },
  data() {
    return {
      showMenuDownload: false,
      loading: true,
      queue: [],
    };
  },
  methods: {
    toogleDownloadMenu() {
      this.$_bus.$emit('profileMenuHide');
      if (this.showMenuDownload) {
        this.$_bus.$emit('mobileHeaderBlack', false);
      } else {
        this.$_bus.$emit('mobileHeaderBlack', true);
      }
      this.showMenuDownload = !this.showMenuDownload;
    },

    getTextInfo(element) {
      if (element.status === 'Search')
        return 'Search..';
      else if (element.status === 'Connecting')
        return 'Connecting..';
      else if (element.status === 'Download' || element.status === 'Transcode')
        return element.percentage + '%';
      else
        return 'Error server.'
    },
    resizeImage(img) {
      // Update imdb picture size by updating the url 182*268
      return img.replace('_V1_', '_V1_UX182_CR0,0,182,268_AL_');
    },

    $_add(data) {
      this.queue.push(data);
    },

    $_delete(data) {
      let title = this.queue[0].title
      const position = this.queue.map(item => item.id)
        .indexOf(data.id);
      this.queue.splice(position, 1);

      if (data.error) {
        swal({ title: `Error during the download of ${title}`, text: data.error, icon: 'error' })
        this.$_bus.$emit('refreshMovieStatus', { status: -1, path: null, stream: false });
      }
    },

    $_update(data) {
      const updatedData = this.queue.map((item) => {
        if (item.id === data.id) {
          return data;
        }
        return item;
      });

      this.queue = updatedData;
      if (data.id === this.queue[0].id && this.$route.params.id !== undefined && this.$route.params.id === this.queue[0].id) {
        if (this.queue[0].stream)
          this.$_bus.$emit('refreshMovieStatus', { status: 1, path: this.queue[0].path, stream: true });
      }
    },

  },
  created() {
    api.torrentList()
      .then((res) => {
        if (res.status >= 500) {
          throw new Error('Server side error');
        } else if (res.status >= 400) {
          throw new Error('Client side error');
        } else {
          return res.json();
        }
      })
      .then((result) => {
        this.loading = false;
        this.queue = result;
      })
      .catch(err => console.error('Error - Message:', err.message));
  },
  mounted() {
    this.$_bus.$on('downloadMenuHide', () => {
      this.showMenuDownload = false;
    });
    this.$_bus.$on('downloadMenuShow', () => {
      if (this.showMenuDownload == false)
        this.toogleDownloadMenu()
    });

    this.$_socket.on('add', this.$_add)
    this.$_socket.on('delete', this.$_delete)
    this.$_socket.on('update', this.$_update)
  },
};
</script>

<style scoped lang="sass">
  @import "DownloadMenu"
</style>
