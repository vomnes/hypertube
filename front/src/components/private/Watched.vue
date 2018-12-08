<template  lang="html">
  <div class="infoContainer watched" :class="{'small': (type == 1), 'big': (type == 2), 'medium': (type == 3)}" v-if="watched || (status && (status == 'waiting' || status == 'downloading' || status == 'ready'))">
    <p v-if="watched" class="textWatched" >{{ $t('watchedLabelUC') }}</p>
    <div v-if="status && (status == 'waiting' || status == 'downloading' || status == 'ready')" class="statusContainer">
      <img v-if="status == 'waiting' || status == 'downloading'" class="iconStatus" src="../../static/icons/three-dots.svg">
      <img v-if="status == 'ready'" class="iconStatus triangle" src="../../static/icons/triangle.svg">
    </div>
  </div>
</template>

<script>

export default {
  props: {
    watched: Boolean,
    type: Number,
    status: String
  },
  data() {
    return {
    };
  },
  created: function() {
	},
  methods: {

    },
};
</script>

<style scoped lang="sass">
$watch-height: 20px

.purple
  background: linear-gradient(90deg, #9461F7 0%, #7C4CF3 100%)
  box-shadow: 0 0 15px 0 rgba(0, 0, 0, 0.15)

.infoContainer
  display: flex
  align-items: center
  justify-content: center
  width: auto
  position: absolute;

.watched
  width: auto
  min-height: $watch-height
  border-radius: 6.5px
  font-family: "Source Sans Pro", serif
  font-size: 10px
  // padding: 0.25em
  font-weight: 600
  line-height: 13px
  color: white
  text-align: center
  @extend.purple

.textWatched
  margin: 0.5em

.iconStatus
  position: relative
  height: 15px
  width: 15px

.triangle
  height: 9px

.textWatched
  margin: 0.5em

.statusContainer
  height: 20px
  display: flex
  align-items: center
  justify-content: center
  width: 20px
  border-radius: 100%
  border-radius: 10.5px
  @extend.purple


.small
  top: 3%
  left: 3%

.medium
  top: 2%
  left: 6%

.big
  top: 7%
  left: 10%

</style>
