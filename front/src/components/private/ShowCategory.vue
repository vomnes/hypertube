<template lang="html">
  <section class="contentRouter" ref="categoryList">
    <div class="category-content" ref="category">
      <p class="box-title category-title">{{ $t(categoryName.toLowerCase() + 'Label').toUpperCase() }}</p>
      <div class="gallery-pictures">
        <div :key="index" class="item" v-for="(picture, index) in movieList">
          <router-link :to="`/movie/${picture.id}`">
            <div class="picture-frame picture-size">
              <img :alt="`picture ${index}`"
                   :src="resizeImage(picture.url)"
                   :title="picture.title"
                   class="adapt-picture">
            </div>
            <watched :watched="picture.isWatched" :status="picture.status" :type="2"></watched>
            <rating :value="picture.rating"></rating>
          </router-link>
        </div>
      </div>
      <loading-animation class="loader" v-if="isLoadingData"/>
    </div>
  </section>
</template>

<script>
import LoadingIcon from './../../static/icons/three-dots.svg';
import PlayIcon from './../../static/icons/triangle.svg';
import Rating from './module/Rating';
import api from '../lib/api';
import auth from '../../static/js/auth';
import scroll from '../../static/js/scroll';
import LoadingAnimation from './module/LoadingAnimation';
import Watched from './Watched';

export default {
  components: {
    Rating,
    LoadingAnimation,
    Watched,
  },
  data() {
    return {
      categoryName: '',
      movieList: [],
      fullLoaded: false,
      offset: 0,
      numberItems: 45,
      isLoadingData: false,
    };
  },
  methods: {
    iconPath(state) {
      if (state) {
        if (state === 'loading') {
          return LoadingIcon;
        } else if (state === 'play') {
          return PlayIcon;
        }
      }
      return null;
    },
    getCategory(category, offset, numberItems) {
      const jwt = auth.token();
      if (!jwt || this.fullLoaded) {
        return;
      }
      api.getCategory({
        category,
        offset,
        numberItems,
        language: (this.$i18n && this.$i18n.locale) ? this.$i18n.locale : 'en',
      }, jwt)
        .then((res) => {
          if (res.status >= 500) {
            throw new Error('Server side error');
          } else if (res.status >= 400) {
            throw new Error('Client side error');
          } else {
            res.json()
              .then((data) => {
                if (data) {
                  if (data.status && data.status === 'No (more) data') {
                    if (this.movieList.length === 0) {
                      this.$router.push('/gallery?status=no_category');
                      return;
                    }
                    this.fullLoaded = true;
                  } else {
                    this.categoryName = category.toUpperCase();
                    data.forEach((item) => {
                      this.movieList.push({
                        id: item.id,
                        title: item.title,
                        url: item.poster,
                        rating: item.rating,
                        isWatched: item.is_watched,
                        status: item.status,
                      });
                    });
                  }
                  this.isLoadingData = false;
                }
              });
          }
        })
        .catch(() => null);
    },
    onScroll() {
      if (!this.isLoadingData && this.$refs.category &&
          scroll.isScrollYMax(this.$refs.categoryList)) {
        this.isLoadingData = true;
        setTimeout(() => {
          this.offset += 45;
          this.getCategory(this.$route.params.name, this.offset, this.numberItems);
        }, 1000);
      }
    },
    resizeImage(img) {
      // Update imdb picture size by updating the url 182*268
      return img.replace('_V1_', '_V1_UX182_CR0,0,182,268_AL_');
    },
  },
  mounted() {
    this.getCategory(this.$route.params.name, this.offset, this.numberItems);
    this.$refs.categoryList.addEventListener('scroll', this.onScroll);
  },
  beforeDestroy() {
    this.$refs.categoryList.removeEventListener('scroll', this.onScroll);
  },
};
</script>

<style scoped lang="sass">
  @import "ShowCategory"
</style>
