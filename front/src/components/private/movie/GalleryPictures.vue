<template lang="html">
  <div class="gallery-pictures" v-if="pictures.length">
    <h2 class="box-title">{{ $t('picturesLabel') }}</h2>
    <div class="box-gallery-pictures">
      <a :href="picture" :key="index" v-for="(picture, index) in pictures">
        <img :alt="`picture${index}`" :src="picture" class="gallery-picture">
      </a>
    </div>
  </div>
</template>

<script>
export default {
  props: {
    pictures: Array,
  },
};
</script>

<style scoped lang="sass">
  img a
    color: rgba(0, 0, 0, 0)
    text-decoration: none

    .gallery-pictures
      position: relative
      margin-top: 5px
      margin-bottom: 5px

    .box-gallery-pictures
      position: relative
      margin-top: 8px
      text-align: center

    .gallery-picture
      position: relative
      height: 8vw
      width: auto
      border-radius: 10px
      background-color: #D8D8D8
      box-shadow: 0 0 10px 0 rgba(52, 36, 84, 0.1)
      margin-right: 5px

    @media screen and (max-width: 640px)

      .gallery-picture
        height: 58px

      .box-gallery-pictures
        white-space: nowrap
        overflow-x: auto
        overflow-y: hidden
        vertical-align: middle
        padding-bottom: 5px
</style>
