<template>
  <div>
    <div class="form-side">
      <router-link :to="{ name: 'index' }">
        <div class="logo" title="Hypertube"></div>
      </router-link>
      <router-view :lang="curLanguage"></router-view>
    </div>
    <div class="picture-side">
      <covers-slideshow></covers-slideshow>
      <div class="color-filter"></div>
      <language-list
        :language="getCurLanguage()"
        @languageChanged="updateTranslation"
      ></language-list>
    </div>
  </div>
</template>

<script>
import LanguageList from './LanguageList';
import CoversSlideshow from './CoversSlideshow';
import LoginForm from './LoginForm';
import { getCurLanguage } from '../../static/js/language';
import modal from '../lib/modal';

const queryString = require('query-string');

export default {
  data() {
    return {
      curLanguage: typeof (Storage) !== 'undefined' ? getCurLanguage() : 'en',
    };
  },
  components: {
    LanguageList,
    LoginForm,
    CoversSlideshow,
  },
  methods: {
    getCurLanguage,
    updateTranslation(event) {
      if (typeof (Storage) !== 'undefined') {
        sessionStorage.setItem('language', event);
        this.curLanguage = event;
      }
    },
    handleURL() {
      const parsed = queryString.parse(location.search);
      if (parsed.status === 'account_created') {
        modal.open({ content: this.$t('accountCreated'), type: 'success' });
      } else if (parsed.status === 'reset_password_email_sent') {
        modal.open({ content: this.$t('resetPasswordEmailSent'), type: 'success' });
      } else if (parsed.status === 'password_reset') {
        modal.open({ content: this.$t('passwordReseted'), type: 'success' });
      } else if (parsed.status === 'failed_reset_password') {
        modal.open({ content: this.$t('failedResetPassword'), type: 'error' });
      }
    },
  },
  mounted() {
    this.handleURL();
  },
  updated() {
    this.handleURL();
  },
};
</script>

<style scoped lang="sass">
  @import "Index"
</style>
