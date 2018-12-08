<template>
  <section id="mainAppSection">
    <layout v-if="$route.meta.requiresAuth"></layout>
    <router-view ref="router"></router-view>
  </section>
</template>

<script>
import Layout from './components/private/Layout';

export default {
  components: {
    Layout,
  },
};
</script>

<style lang="sass">
  @import "static/globalreset"
  #mainAppSection
    display: flex
    flex-direction: column
    height: 100%

  .contentRouter
    position: relative
    width: 95%
    height: 100%
    overflow: auto
    margin-left: auto
    margin-right: auto
    display: block
    margin-top: 4em
</style>
