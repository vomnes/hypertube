<template lang="html">
  <div :title="$t('rating')" class="rating" v-if="value">
    <p>{{ value }}</p>
  </div>
</template>

<script>
export default {
  props: {
    value: Number,
  },
};
</script>

<style scoped lang="sass">
  .rating
    position: absolute
    height: 20px
    width: 81px
    border-radius: 5px
    background-color: #F8F9FD
    box-shadow: 0 0 10px 0 rgba(52, 36, 84, 0.1)
    bottom: 0
    text-align: center
    left: 50%
    transform: translate(-50%, 0)

    & p
      color: rgba(74, 67, 80, 0.6)
      font-family: "Source Sans Pro", serif
      font-size: 10px
      font-weight: bold
      line-height: 13px
      margin-top: 3px
      margin-bottom: 4px

  @media screen and (max-width: 640px)
    .rating
      width: 50px
      height: 15px
      left: 45%

      & p
        font-size: 8px
        line-height: 8px
</style>
