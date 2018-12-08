<template lang="html">
  <div class="movie-informations">
    <div class="movie-cover">
      <img :src="url" class="adapt-picture">
      <watched :watched="isWatched" :type="1"></watched>
      <button :title="$t('playMovie')"
              @click="$emit('updatePlayerVue')"
              class="status-movie"
              name="play-movie"
              v-if="(movieObject.status === 2 && movieObject.stream === true) || (movieObject.status === 1 && movieObject.stream === true)">
        <div class="play-triangle"></div>
      </button>
      <button :title="$t('isLoadingMovie')"
              class="status-movie no-opacity"
              v-if="movieObject.status === 0 || (movieObject.status === 1 && movieObject.stream === false)">
        <p class="dots">...</p>
      </button>
      <button :title="$t('loadMovie')"
              @click="$emit('openTorrents')"
              class="status-movie no-opacity"
              v-if="movieObject.status === -1">
        <img alt="download" class="download" src="../../../static/icons/download-white.svg">
      </button>
    </div>
    <div class="movie-data">
      <h1 class="movie-title">{{ title }}</h1>
      <p class="movie-synopsis">{{ synopsis }}</p>
      <ul class="movie-list-details">
        <li><span style="font-weight: 600">{{ $t('genreLabel') }}:</span> 
              <span v-for="(genre, index) in genres">
                <span>{{$t(genre.toLowerCase() + 'Label')}}</span><span v-if="genres[index + 1]">, </span>
            </span>
        </li>
        <li><span style="font-weight: 600">{{ $t('ratingLabelLC') }}:</span> {{ rating }}/10</li>
        <li><span style="font-weight: 600">{{$t('prodYearLabelLC')}}:</span> {{ year }}</li>
        <li><span style="font-weight: 600">{{ $t('durationLabel') }}:</span> {{ duration }}</li>
      </ul>
      <a :href="trailer.url" rel="noopener noreferrer" target="_blank" v-if="trailer">
        <button :title="`${trailer.title} - ${trailer.quality}p`"
                class="play-trailer"
                name="play-trailer">
        {{ $t('trailerLabel') }}
        </button>
      </a>
      <p class="movie-synopsis-responsive">
        {{ synopsisProcessed }}
        <span @click="fullText = !fullText" class="see-more"
              v-if="synopsis.length > 300">{{fullText ? $t('seeLessLabel') : $t('seeMoreLabel') }}
        </span>
      </p>
    </div>
  </div>
</template>

<script>
import Watched from '../Watched'
export default {
  components: {
    Watched,
  },
  props: {
    title: String,
    url: String,
    synopsis: String,
    genres: Array,
    rating: Number,
    year: Number,
    duration: String,
    trailer: Object,
    movieObject: {
      stream: Boolean,
      status: Number,
      path: String
    },
    isWatched: Boolean
  },
  data() {
    return {
      fullText: false,
    };
  },
  computed: {
    synopsisProcessed() {
      if (this.fullText || this.synopsis.length < 300) return this.synopsis;
      return `${this.synopsis.substring(0, 300)}...`;
    },
  }
};
</script>

<style scoped lang="sass">
  @import "CoverSide"
</style>
