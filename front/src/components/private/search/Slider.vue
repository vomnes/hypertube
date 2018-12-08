<template lang="html">
  <div class="filter">
    <h1>{{ title }}</h1>
    <vue-slider
      class="slider-position"
      ref="slider"
      @drag-end="updateRange"
      v-bind="range"
      v-model="range.value">
    </vue-slider>
  </div>
</template>

<script>
import VueSlider from 'vue-slider-component';

export default {
  components: {
    VueSlider,
  },
  props: {
    type: String,
    title: String,
    value: Array,
    min: Number,
    max: Number,
    scale: Number,
  },
  methods: {
    updateRange() {
      this.$emit('updateRange', this.range.value, this.type);
    },
  },
  data() {
    return {
      range: {
        value: this.value,
        min: this.min,
        max: this.max,
        interval: this.scale,
        width: '90%',
        useKeyboard: true,
        show: true,
        speed: 0,
        bgStyle: {
          backgroundColor: '#F8F9FD',
          borderRadius: '2.5px',
          height: '12px',
        },
        tooltipStyle: {
          backgroundColor: '#FFFFFF00',
          borderColor: '#FFFFFF00',
          color: '#4A4350',
          fontFamily: 'Source Sans Pro',
          top: '13px',
          position: 'relative',
          fontSize: '8px',
          fontWeight: 'bold',
          lineHeight: '10px',
        },
        processStyle: {
          backgroundColor: '#65E3B9',
        },
        sliderStyle: {
          backgroundColor: '#4A4350',
          position: 'relative',
          top: '2.5px',
          boxShadow: 'none',
        },
        focusStyle: {
          boxShadow: 'none',
        },
      },
    };
  },
};
</script>

<style scoped lang="sass">
  @import "Search"
  .filter
    position: relative
    width: 50%
    height: 70px
    margin-top: 10px
    margin-bottom: 10px
    display: block
    @extend .box-layout

    & h1
      @extend .font-style
      font-size: 10px
      font-weight: 600
      line-height: 13px
      margin-top: 5px
      margin-left: 10px

  .slider-position
    position: absolute
    bottom: 16px
    left: 50%
    transform: translate(-50%, 0)

  @media screen and (max-width: 640px)
    .filter
      height: 50px
      margin-top: 6px
      margin-bottom: 6px

      & h1
        font-size: 8px

    .slider-position
      bottom: 0
</style>
