<template lang="html">
  <section class="contentRouter" ref="movieList">
    <show-pictures
      :link="popular.link"
       :list="popular.list"
      :size="popular.size"
      :title="popular.title"
      :updateDataAvailable="true"
      @loadMore="getCategory">
    </show-pictures>
    <div :key="key" :ref="key" v-for="(category, key) in list">
      <show-pictures
        :link="category.link"
        :list="category.list"
        :size="category.size"
        :title="category.title"
        :updateDataAvailable="true"
        @loadMore="getCategory">
      </show-pictures>
      <div style="position: relative; height: 25px;"></div>
    </div>
  </section>
</template>

<script>
import ShowPictures from './ShowPictures';
import api from '../lib/api';
import auth from '../../static/js/auth';
import string from '../../static/js/string';
import modal from '../lib/modal';

const queryString = require('query-string');

const categories = [
  'thriller', 'horror', 'history', 'documentary', 'sport', 'biography',
  'sci-Fi', 'musical', 'music', 'adventure', 'romance', 'crime', 'drama',
  'western', 'action', 'comedy', 'war', 'animation', 'fantasy', 'family',
  'mystery',
];

export default {
  components: {
    ShowPictures,
  },
  mounted() {
    this.getCategory('popular', 0, 10);
    this.$refs.movieList.addEventListener('scroll', this.loadCategories);
    window.addEventListener('resize', this.loadCategories);
    const parsed = queryString.parse(location.search);
    if (parsed.status === 'no_movie') {
      modal.open({ content: this.$t('noMovie'), type: 'error' });
    } else if (parsed.status === 'no_category') {
      modal.open({ content: this.$t('noCategory'), type: 'error' });
    }
  },
  methods: {
    getCorrectLabel(category, addLabel=false) {
      if (addLabel)
        return this.$t(category.toLowerCase() + 'Label');
      return this.$t(category);
    },

    getCategory(category, offset, numberItems, language) {
      if ((category === 'popular' && this.popular && this.popular.fullLoaded) ||
          (category !== 'popular' && this.list && this.list[category] && this.list[category].fullLoaded)) {
        return;
      }
      api.getCategory({
        category,
        offset,
        numberItems,
        language: (this.$i18n && this.$i18n.locale) ? this.$i18n.locale : 'en',
      },
      auth.token())
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
                    if (category === 'popular') {
                      this[category].fullLoaded = true;
                    } else {
                      this.list[category].fullLoaded = true;
                    }
                  } else {
                    data.forEach((item) => {
                      if (category === 'popular') {
                        this[category].list.push({
                          id: item.id,
                          title: item.title,
                          url: item.poster,
                          isWatched: item.is_watched,
                          status: item.status,
                        });
                      } else {
                        this.list[category].list.push({
                          id: item.id,
                          title: item.title,
                          url: item.poster,
                          rating: item.rating,
                          isWatched: item.is_watched,
                          status: item.status,
                        });
                      }
                    });
                  }
                }
              });
          }
        })
        .catch(() => null);
    },
    generateListCategories() {
      const categoryList = {};
      categories.forEach((category) => {
        categoryList[category] = {
          title: category,
          list: [],
          size: 'medium',
          link: true,
          fullLoaded: false,
        };
      });
      return categoryList;
    },
    loadCategories() {
      categories.forEach((category) => {
        if (this.isVisible(this.$refs[category][0]) &&
            string.inArray(category, this.loadedCategories) === false) {
          this.loadedCategories.push(category);
          this.getCategory(category, 0, 15);
        }
      });
    },
    isVisible(element) {
      const cTop = this.$refs.movieList.scrollTop;
      const cBottom = cTop + this.$refs.movieList.clientHeight;

      const eTop = element.offsetTop;
      const eBottom = eTop + element.clientHeight;

      const isTotal = (eTop >= cTop && eBottom <= cBottom);
      const isPartial = ((eTop < cTop && eBottom > cTop) || (eBottom > cBottom && eTop < cBottom));

      return (isTotal || isPartial);
    },
  },
  data() {
    return {
      loadedCategories: [],
      popular: {
        title: 'Popular',
        list: [],
        size: 'big',
        link: true,
        fullLoaded: false,
      },
      list: this.generateListCategories(),
    };
  },
  updated() {
    this.loadCategories();
  },
  beforeDestroy() {
    this.$refs.movieList.removeEventListener('scroll', this.loadCategories);
    window.removeEventListener('resize', this.loadCategories);
  },
};
</script>

<style lang="sass" scoped>

</style>
