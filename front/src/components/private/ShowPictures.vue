<template lang="html">
  <div class="line-pictures">
    <router-link
      :to="`/category/${title.toLowerCase()}`"
      v-if="link">
      <div class="box-title hover-title">
        <p
          :class="`${size}-title`"
          :title="`See more ${title.toLowerCase()} movies`">
          {{ $t(title.toLowerCase() + 'Label').toUpperCase() }}
        </p>
        <img alt="more of this category" class="arrow" src="../../static/icons/category-arrow.svg"
             :class="`${size}-arrow`">
        <loading-animation
          :class="`loading-pictures-${size}`"
          v-if="isLoading"/>
      </div>
    </router-link>
    <div class="box-title" v-else>
      <p
        :class="`${size}-title`"
        v-if="title !== 'empty'">{{ $t(title.toLowerCase() + 'Label').toUpperCase() }}</p>
      <loading-animation
        :class="`loading-pictures-${size}`"
        v-if="isLoading"/>
    </div>
    <div v-if="!touchScreen" class="pictures" ref="pictures"
         :id="`pictures-${title.toLowerCase()}`"
         v-dragscroll
         v-on:dragscrollend="isClickable = true"
         v-on:dragscrollmove="isClickable = false"
         >
      <div class="picture"
           :key="index"
           v-for="(element, index) in list"
           v-if="list.length">

        <div :class="[`${size}-frame`, 'picture-frame']"
             :title="element.title">
          <router-link
            :to="setLink(element.id)"
            v-if="element.id && isClickable">
            <img class="adapt-picture crusor-drag"
              :alt="`picture ${title}`"
              :src="resizeImage(element.url)">
          </router-link>
          <img class="adapt-picture crusor-drag"
               :alt="`picture ${title}`"
               :src="resizeImage(element.url)"
               v-else>
              <watched :watched="element.isWatched" :status="element.status" :type="3"></watched>
        </div>
        <rating
          :style="{ bottom: '-9px' }"
          :value="element.rating"></rating>
      </div>
      <div class="picture" v-if="!list.length">
        <div :class="[`${size}-frame`, 'picture-frame']"></div>
      </div>
    </div>

    <div v-else class="pictures" ref="pictures"
         :id="`pictures-${title.toLowerCase()}`"
         >
      <div class="picture"
           :key="index"
           v-for="(element, index) in list"
           v-if="list.length">
        <div :class="[`${size}-frame`, 'picture-frame']"
             :title="element.title">
          <router-link
            :to="setLink(element.id)"
            v-if="element.id && isClickable">
            <img class="adapt-picture crusor-drag"
              :alt="`picture ${title}`"
              :src="resizeImage(element.url)">
          </router-link>
          <img class="adapt-picture crusor-drag"
               :alt="`picture ${title}`"
               :src="resizeImage(element.url)"
               v-else>
        </div>
        <rating
          :style="{ bottom: '-9px' }"
          :value="element.rating"></rating>
      </div>
      <div class="picture" v-if="!list.length">
        <div :class="[`${size}-frame`, 'picture-frame']"></div>
      </div>
    </div>
  </div>
</template>

<script>
import { dragscroll } from 'vue-dragscroll';
import Rating from './module/Rating';
import LoadingAnimation from './module/LoadingAnimation';
import scroll from '../../static/js/scroll';
import Watched from './Watched'


export default {
  components: {
    Rating,
    LoadingAnimation,
    Watched
  },
  props: {
    title: String,
    list: Array,
    link: Boolean,
    size: {
      type: String,
      default: 'small', // tiny - small - medium - big
    },
    loadMore: Function,
    updateDataAvailable: Boolean,
  },
  directives: {
    dragscroll,
  },
  data() {
    if (this.updateDataAvailable) {
      return {
        isLoading: false,
        isClickable: true,
        offset: 10,
        numberItems: 10,
        touchScreen: false,
      };
    }
    return {
      isLoading: false,
      isClickable: true,
    };
  },
  created: function() {
		if ('ontouchstart' in document.documentElement)
			this.touchScreen = true
		else
			this.touchScreen = false
	},
  methods: {
    setLink(id) {
      if (id) return `/movie/${id}`;
      return null;
    },
    onScroll() {
      if (scroll.isScrollXMax(this.$refs.pictures)) {
        this.isLoading = true;
        setTimeout(() => {
          this.$emit('loadMore', this.title.toLowerCase(), this.offset, this.numberItems, (this.$i18n && this.$i18n.locale) ? this.$i18n.locale : 'en');
          this.offset += 10;
          this.isLoading = false;
        }, 1000);
      }
    },
    createRefScroll() {
      if (this.updateDataAvailable) {
        const picturesRef = document.getElementById(`pictures-${this.title.toLowerCase()}`);
        if (picturesRef) {
          picturesRef.addEventListener('scroll', this.onScroll);
        }
      }
    },
    resizeImage(img) {
      // Update imdb picture size by updating the url 182*268
      if (this.size === 'medium') {
        return img.replace('_V1_', '_V1_UX182_CR0,0,182,268_AL_');
      }
      return img;
    },
  },
  mounted() {
    this.createRefScroll();
  },
  updated() {
    this.createRefScroll();
  },
  destroyed() {
    this.createRefScroll();
  },
};
</script>

<style scoped lang="sass">
  @import "ShowPictures"
</style>
