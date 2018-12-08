<template lang="html">
  <section class="contentRouter">
    
    <cover-side
      :director="details.director"
      :duration="details.duration"
      :genres="details.genres"
      :rating="details.rating"
      :movieObject="movie"
      :synopsis="details.synopsis"
      :title="details.title"
      :trailer="details.trailer"
      :url="details.url"
      :year="details.year"
      :isWatched="details.isWatched"
      @openTorrents="openTorrents"
      @updatePlayerVue="updatePlayerVue"
    >
    </cover-side>
    <div class="movie-other-data">
      <show-pictures
        :list="casting.list"
        :size="casting.size"
        :title="casting.title">
      </show-pictures>
      <gallery-pictures
        :pictures="listPictures">
      </gallery-pictures>
      <show-pictures
        :list="similarMovies.list"
        :size="similarMovies.size"
        :title="similarMovies.title">
      </show-pictures>
      <comments
        :commentLoading="commentLoading"
        :comments="comments"
        :filmId="details.imdb"
      ></comments>
      <loading :active="globalLoading"></loading>
    </div>
    <torrent-list
      :change="changeStatus"
      :movieID="details.imdb"
      @closeTorrents="closeTorrents"
      v-if="torrentOpened"></torrent-list>
    <player-component
      @updatePlayerVue="updatePlayerVue"
      v-if="videoPlayerActive"
      :moviePath="movie.path"
      :sub="movie.subtitles"
    ></player-component>
  </section>
</template>

<script>
// import Layout from '../Layout';
import ShowPictures from '../ShowPictures';
import Comments from './Comments';
import CoverSide from './CoverSide';
import GalleryPictures from './GalleryPictures';
import TorrentList from './TorrentList';
import Loading from '../module/Loading';
import api from '../../lib/api';
import auth from '../../../static/js/auth';
import time from '../../../static/js/time';
import PlayerComponent from '../PlayerComponent';
import moment from 'moment';
import swal from 'sweetalert';

export default {
  components: {
    // Layout,
    ShowPictures,
    Comments,
    CoverSide,
    GalleryPictures,
    TorrentList,
    Loading,
    PlayerComponent,
  },
  watch: {
    '$route.params.id': function () {
      this.loadData();
    },
  },
  created() {
    this.loadData();
    this.$_bus.$on('langChange', this.langChangeHandler);
    this.$_bus.$on('refreshMovieStatus', this.refreshMovieStatusHandler);
  },
  data() {
    return {
      globalLoading: true,
      videoPlayerActive: false,
      torrentOpened: false,
      movie: {
        stream: false,
        status: -1,
        path: null,
      },
      details: {
        title: '',
        url: '',
        synopsis: '',
        rating: 10,
        year: 0,
        duration: '0',
        genres: [],
        imdb: '',
      },
      casting: {
        title: 'CASTING',
        list: [],
        size: 'small',
      },
      similarMovies: {
        title: 'SIMILAR',
        list: [],
        size: 'small',
      },
      listPictures: [],
      comments: [],
      commentLoading: true,
      errorServer: false,
      interval: []
    };
  },
  methods: {
    langChangeHandler(data) {
        this.loadData()
    },
    refreshMovieStatusHandler(data) {
      if (!this.movie.stream && data.stream)
        this.loadData();
      else if (this.movie.status !== data.status)
        this.loadData();
    },
    openTorrents() {
      this.torrentOpened = true;
    },

    closeTorrents() {
      this.torrentOpened = false;
    },

    updatePlayerVue() {
      this.videoPlayerActive = !this.videoPlayerActive;
    },

    changeStatus() {
      this.movie.status = 0;
    },

    loadData() {
      // Get data for the movie
      api.getMovie({
        filmID: this.$route.params.id,
        language: (this.$i18n && this.$i18n.locale) ? this.$i18n.locale : 'en',
      },
      auth.token())
        .then((res) => {
          this.globalLoading = true;
          if (res.status >= 500) {
            throw new Error('Server side error');
          } else if (res.status >= 400) {
            this.$router.push('/gallery?status=no_movie');
          } else {
            res.json()
              .then((data) => {
                if (data) {
                  this.casting.list = [];
                  if (data.casting) {
                    data.casting.forEach((person) => {
                      this.casting.list.push({
                        title: `${person.name} - ${person.character}`,
                        url: person.picturePath,
                      });
                    });
                  }
                  this.similarMovies.list = [];
                  if (data.similar) {
                    data.similar.forEach((similar) => {
                      this.similarMovies.list.push({
                        id: similar.id,
                        title: similar.title,
                        url: similar.coverPath,
                      });
                    });
                  }

                  this.listPictures = [];
                  if (data.images) this.listPictures = data.images;
                  this.details.title = data.title;
                  this.details.url = data.cover;
                  this.details.duration = time.getHoursFromMins(data.duration);
                  this.details.rating = data.rating;
                  this.details.year = new Date(data.production_year).getFullYear();
                  this.details.synopsis = data.resume;
                  this.details.genres = data.genres;
                  this.details.trailer = data.trailer;
                  this.details.imdb = data.id_imdb;
                  this.details.isWatched = data.is_watched;
                  
                  if (data.movie.status === 2 || data.movie.status === 1)
                    this.movie = { ...data.movie, stream: true };
                  else
                    this.movie = { ...data.movie, stream: false };

                  this.getComments(data);
                  this.globalLoading = false;
                }
              });
          }
        })
        .catch(err => {
          swal({ text: this.$t('errorServer'), icon: 'error' })
          this.commentLoading = false;
          this.globalLoading = false;
        });
    },

    updateComments(comment) {
      this.comments.unshift({ ...comment, createdat: moment(comment.createdat).locale(this.$i18n.locale).fromNow()});
      let self = this;
      let I = setInterval(function() {
        const id = self.comments.map(item => item.id).indexOf(comment.id);
        self.comments[id].createdat = moment(comment.createdat).locale(self.$i18n.locale).fromNow();
      }, 1000);
      this.interval.push(I);
    },

    getComments(data) {
      this.destroyInterval();
      this.comments = [];
      api.getFilmComments(
        data.id_imdb,
        auth.token(),
      )
      .then((res) => {
        if (res.status === 200) {
          res.json()
          .then((data) => {
            if (data.status !== 'No comments') {
              let reverseData = data.reverse();
              this.comments = reverseData.map(el => ({
                  ...el,
                  createdat: moment(el.createdat).locale(this.$i18n.locale).fromNow()
              }));

              let self = this;
              reverseData.forEach((item, index) => {
                let I = setInterval(() => {
                  const id = self.comments.map(el => el.id).indexOf(item.id);
                  const lang = self.$i18n.locale;
                  self.comments[id].createdat = moment(item.createdat).locale(lang).fromNow()
                }, 1000);
                this.interval.push(I);
              });
            }
            this.commentLoading = false;
          })
          .catch(error => {
            this.commentLoading = false;
            swal({ text: this.$t('errorServer'), icon: 'error' });
          });
        } else {
          swal({ text: this.$t('errorServer'), icon: 'error' });
          this.commentLoading = false;
        }
      })
      .catch((err) => {
        swal({ text: this.$t('errorServer'), icon: 'error' })
        this.commentLoading = false;
      });
    },

    destroyInterval() {
      this.interval.forEach(item =>{
        clearInterval(item);
      });

      this.interval = [];
    }
  },
  destroyed() {
    this.$_bus.$off('langChange', this.langChangeHandler);
    this.$_bus.$off('refreshMovieStatus', this.refreshMovieStatusHandler);
    this.destroyInterval();
  }
};
</script>

<style scoped lang="sass">
  .movie-other-data
    position: absolute
    top: 0  
    left: 355px
    width: calc(100% - 355px)

  @media screen and (max-width: 640px)

    .movie-other-data
      position: relative
      width: 100%
      left: 0
      margin-top: 10px
</style>

<style lang="sass">
  .box-title
    font-family: "Source Sans Pro", serif
    font-size: 13px
    color: #8F848F
    font-weight: 600
    line-height: 20px
    letter-spacing: 0.5px
    margin-left: 8px
</style>
