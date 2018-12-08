<template>
  <router-link :to="{ name: 'index' }" title="Hypertube">
    <div class="center">
      <div class="logo"></div>
      <h1><span style="font-weight: 300; font-size: 30px;">HYPER</span><span style="font-weight: 600; font-size: 30px;">TUBE</span></h1>
      <h2>404 - Page not found</h2>
    </div>
  </router-link>
</template>

<script>
export default {};
</script>

<style scoped lang="sass">
  .center
    position: absolute
    width: 250px
    height: auto
    top: calc(50% - 125px)
    left: calc(50% - 125px)
    font-family: 'Source Sans Pro', sans-serif
    color: #384064
    text-align: center
    line-height: 1.5
  .logo
    position: relative
    left: 70px
    width: 100px
    height: 100px
    background-image: url('../static/logo.svg')
    background-size: cover
    -moz-background-size: cover
</style>
