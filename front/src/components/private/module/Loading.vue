<template lang="html">
  <div :style="{ background: transparent === true ? 'none': '#F8F9FD' }" class="loading"
       v-if="active">
    <img alt="logo" class="loading-logo" src="../../../static/logo.svg">
    <div class="animation">
      <loading-animation></loading-animation>
    </div>
  </div>
</template>

<script>
import LoadingAnimation from './LoadingAnimation';

export default {
  props: {
    active: Boolean,
    transparent: Boolean,
  },
  components: {
    LoadingAnimation,
  },
};
</script>

<style scoped lang="sass">
  .loading-logo
    position: fixed
    top: calc(50% - 100px)
    left: calc(50% - 50px)
    width: 77.5px
    min-width: 75px
    height: auto
    background: none

  .loading
    width: 100%
    height: 100%
    position: fixed
    top: 0
    left: 0

  .animation
    position: fixed
    top: calc(50% - 20px)
    left: calc(50% - 42px)
    width: auto
    height: 10px
</style>
