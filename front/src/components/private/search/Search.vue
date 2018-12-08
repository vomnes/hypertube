<template lang="html">
  <section class="contentRouter" ref="searchList">
    <div class="search-area">
      <input class="search-field" :placeholder="$t('searchLabel')"
             type="text" v-model="fields.search">
      <img @click="$_bus.$emit('headerSpaceToogle'); toggleParameters = !toggleParameters"
           alt="settings"
           class="settings-icon"
           src="../../../static/icons/settings-icon.svg">
      <transition mode="out-in" name="fade">
        <div class="filters" v-if="toggleParameters">
          <div class="filter-area">
            <slider :max="10"
                    :min="0"
                    :scale="0.1"
                    :title="$t('ratingLabelUC')"
                    :type="'rating'"
                    :value="[fields.rating.min, fields.rating.max]"
                    @updateRange="updateRange"
                    class="rating-filter">
            </slider>
            <slider :max="new Date(Date.now()).getFullYear()"
                    :min="1900"
                    :scale="1"
                    :title="$t('prodYearLabelUC')"
                    :type="'year'"
                    :value="[fields.year.min, fields.year.max]"
                    @updateRange="updateRange"
                    class="year-filter">
            </slider>
          </div>
          <div class="filter-area types">
            <p :style="typeSelected('ready')" @click="switchType('ready')"
               title="List only downloaded movies">{{ $t('downloadedLabel') }}</p>
            <p :style="typeSelected('downloading')" @click="switchType('downloading')"
               class="middle-type" title="List only loading movies">{{ $t('loadingLabel') }}</p>
            <p :style="typeSelected('not downloaded')"
               @click="switchType('not downloaded')"
               title="List only not downloaded movies">{{$t('notDownloadedLabel')}}</p>
          </div>
        </div>
      </transition>
    </div>
    <div class="search-response-area" ref="results" v-show="!newSearch">
      <router-link :to="`/movie/${movie.id}`" v-for="(movie, index) in list" :key="index" >
        <div class="data-box">
          <div class="picture-frame picture-size">
            <img :alt="`cover picture ${index}`"
                 :src="resizeImage(movie.url)" class="adapt-picture">
                 <watched :watched="movie.iswatched" :status="movie.status" :type="2"></watched>
          </div>
          <div class="text-area">
            <div>
              <h1 class="title">{{ movie.title }}</h1>
              <p class="rating">{{ $t('rating') }}: {{ movie.rating }}</p>
              <p class="rating">{{ $t('durationLabel') }}: {{ movie.duration }}</p>
              <p class="rating">{{ $t('genresLabel') }}:
                <span v-for="(genre, index) in movie.genres" :key="index">
                    <span>{{$t(genre.toLowerCase() + 'Label')}}</span><span v-if="movie.genres[index + 1]">, </span>
                </span>
              </p>
              <p class="rating">{{ $t('yearLabel') }}: {{ movie.year }}</p>
              <p class="rating">{{ $t('watchedLabelLC') }}: {{ movie.iswatched ?
                $t('yesLabel') : $t('noLabel') }}</p>
            </div>
          </div>
        </div>
      </router-link>
      <h1 v-if="!list || list.length < 1" class="search-empty">{{ $t('noResultLabel') }}</h1>
    </div>
    <loading :active="newSearch" :transparent="true"></loading>
    <loading-animation class="loader" v-if="isLoadingData"/>
  </section>
</template>

<script>
import Slider from './Slider';
import api from '../../lib/api';
import auth from '../../../static/js/auth';
import base64 from '../../../static/js/base64';
import Loading from '../module/Loading';
import LoadingAnimation from '../module/LoadingAnimation';
import scroll from '../../../static/js/scroll';
import time from '../../../static/js/time';
import Watched from '../Watched'

export default {
  components: {
    Slider,
    Loading,
    LoadingAnimation,
    Watched
  },
  data() {
    return {
      fields: {
        search: null,
        rating: {
          min: 0,
          max: 10,
        },
        year: {
          min: 1900,
          max: new Date(Date.now()).getFullYear(),
        },
        typeMovies: null,
      },
      toggleParameters: false,
      list: [],
      offset: 0,
      numberItems: 25,
      fullLoaded: false,
      newSearch: false,
      isLoadingData: false,
    };
  },
  methods: {
    handleLongSynopsis(content) {
      if (content.length < 150) {
        return content;
      }
      return `${content.substring(0, 150)}...`;
    },
    switchType(type) {
      if (type === this.fields.typeMovies) {
        this.fields.typeMovies = null;
        return;
      }
      this.fields.typeMovies = type;
    },
    typeSelected(selected) {
      return this.fields.typeMovies === selected ? { fontWeight: '600' } : {};
    },
    updateRange(value, type) {
      if (type !== 'year' && type !== 'rating') return;
      if (value && value.length !== 2) return;
      this.fields[type] = {
        min: value[0],
        max: value[1],
      };
    },
    runSearch(isUpdated) {
      if (!isUpdated && this.fullLoaded) {
        return;
      }
      if (isUpdated) {
        this.newSearch = true;
      }
      const optionsBase64 = base64.objectToBase64({
        search: this.fields.search,
        rating: {
          max: this.fields.rating.max,
          min: this.fields.rating.min,
        },
        year: {
          max: this.fields.year.max,
          min: this.fields.year.min,
        },
        status: this.fields.typeMovies,
      });
      api.searchMovies({
        optionsBase64,
        offset: this.offset,
        numberItems: this.numberItems,
        language: (this.$i18n && this.$i18n.locale) ? this.$i18n.locale : 'en',
      },
      auth.token())
        .then((res) => {
          if (res.status >= 500) {
            throw new Error('Server side error');
          } else if (res.status >= 400) {
            throw new Error('Client side error');
          } else {
            res.json()
              .then((data) => {
                if (data) {
                  if (isUpdated) {
                    this.clearData();
                    this.list = [];
                  }
                  if (data.status && data.status === 'No (more) data') {
                    this.fullLoaded = true;
                  } else {
                    data.forEach((item) => {
                      this.list.push({
                        id: item.id,
                        title: item.title,
                        url: item.poster,
                        rating: item.rating,
                        genres: item.genres,
                        year: item.year,
                        duration: time.getHoursFromMins(item.duration),
                        iswatched: item.is_watched,
                        status: item.status,
                      });
                    });
                    if (data.length < this.numberItems) {
                      this.fullLoaded = true;
                    }
                  }
                }
              });
            this.newSearch = false;
            this.isLoadingData = false;
          }
        })
        .catch(() => {
          // Clear loading flags
          this.newSearch = false;
          this.isLoadingData = false;
        });
    },
    onFieldUpdated() {
      if (this.$_searchTimeout) {
        window.clearTimeout(this.$_searchTimeout);
        this.$_searchTimeout = null;
      }
      this.$_searchTimeout = setTimeout(() => {
        this.runSearch(true);
        this.$_searchTimeout = null;
      }, 500);
    },
    clearData() {
      this.offset = 0;
      this.numberItems = 25;
      this.fullLoaded = false;
    },
    onScroll() {
      if (!this.fullLoaded && !this.isLoadingData
          && this.$refs.results && scroll.isScrollYMax(this.$refs.searchList)) {
        this.isLoadingData = true;
        setTimeout(() => {
          this.offset += this.numberItems;
          this.runSearch(false);
        }, 1000);
      }
    },
    resizeImage(img) {
      // Update imdb picture size by updating the url 182*268
      return img.replace('_V1_', '_V1_UX182_CR0,0,182,268_AL_');
    },
  },
  mounted() {
    this.$refs.searchList.addEventListener('scroll', this.onScroll);
    this.runSearch(true);
  },
  watch: {
    fields: {
      handler() {
        this.onFieldUpdated();
      },
      deep: true,
    },
  },
  beforeDestroy() {
    this.$_bus.$emit('headerSpaceToogle', false);
    this.$refs.searchList.removeEventListener('scroll', this.onScroll);
  },
};
</script>

<style scoped lang="sass">
  @import "Search"
  @import "Result"
</style>
