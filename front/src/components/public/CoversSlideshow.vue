<template lang="html">
  <transition mode="out-in" name="cover-fade">
    <img :key="coverPicture" :src="coverPicture"
         alt="movie cover" class="movie-cover">
  </transition>
</template>

<script>
import Picture1 from './../../static/pictures/movie-cover-1.jpg';
import Picture2 from './../../static/pictures/movie-cover-2.jpg';
import Picture3 from './../../static/pictures/movie-cover-3.jpg';

export default {
  data() {
    return {
      interval: '',
      covers: [Picture1, Picture2, Picture3],
      coverPicture: Picture1,
      pictureIndex: 0,
    };
  },
  mounted() {
    this.updatePictureWithInterval();
  },
  beforeDestroy() {
    clearInterval(this.interval);
  },
  methods: {
    updatePictureWithInterval() {
      this.interval = setInterval(() => {
        this.pictureIndex += 1;
        if (this.pictureIndex === this.covers.length) this.pictureIndex = 0;
        this.coverPicture = this.covers[this.pictureIndex];
      }, 5000);
    },
  },
};
</script>

<style scoped lang="sass">
  .movie-cover
    width: 100%
    height: 100%
    object-fit: cover

  .cover-fade-enter-active, .cover-fade-leave-active
    transition: opacity 0.75s

  .cover-fade-enter, .cover-fade-leave-active
    opacity: 0
</style>
