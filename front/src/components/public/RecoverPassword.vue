<template>
  <div class="width-form">
    <p style="font-size: 32px;">
      <span style="font-weight: 300">HYPER</span><span style="font-weight: 600">TUBE</span>
    </p>
    <p style="margin-bottom: 40px; font-style: italic">{{ $t('subtitle') }}</p>
    <form @submit.prevent="sendEmail">
      <input type="text" autocomplete="username" style="visibility: hidden; position: absolute">
      <input :placeholder="$t('emailLabel')" autofocus name="email"
              autocomplete="email"
             require type="email" v-model="email">
      <button type="submit">{{ $t('restoreLabel') }}</button>
    </form>
    <div class="footerForm">
      <a @click="goBack" class="sublink" style="float: left">
        {{ $t('backLabel') }}
      </a>
    </div>
  </div>
</template>

<script>
import CoversSlideshow from './CoversSlideshow';
import LanguageList from './LanguageList';
import api from '../lib/api';
import modal from '../lib/modal';

export default {
  components: {
    CoversSlideshow,
    LanguageList,
  },
  data() {
    return {
      email: null,
    };
  },
  props: ['lang'],
  watch: {
    lang() {
      this.$i18n.locale = this.lang;
    },
  },
  methods: {
    updateTranslation(event) {
      if (typeof (Storage) !== 'undefined') {
        sessionStorage.setItem('language', event);
        this.curLanguage = event;
        this.$i18n.locale = this.curLanguage;
      }
    },
    sendEmail() {
      if (this.email == null) {
        modal.open({ content: this.$t('emptyEmailField'), type: 'warning' });
        return;
      }
      api.forgotPasswordSendEmail({
        email: this.email,
      })
        .then((res) => {
          if (res.status >= 500) {
            modal.open({ content: this.$t('errorServer'), type: 'warning' });
          } else if (res.status >= 400) {
            res.json()
              .then((data) => {
                modal.open({
                  content: this.$t(modal.errorFormatLanguage(data.error, res.status)),
                  type: 'error',
                });
              });
          } else {
            this.$router.push('/?status=reset_password_email_sent');
          }
        })
        .catch(() => {
          modal.open({ content: this.$t('serverConnectionFailed'), type: 'warning' });
        });
    },
    goBack() {
      window.history.back();
    },
  },
};
</script>

<style scoped lang="sass">
  @import "Form"
  @import "Index"

  input[type="file"]
    display: none

  .width-form
    position: relative
    top: calc(50% - (202px / 2))
    margin: auto
    width: 25vw

  @media screen and (max-width: 640px)
    p
      color: #F8F9FD
      text-shadow: 1px 1px 3px rgba(150, 150, 150, 1)

    .width-form
      position: relative
      width: 80%
      top: -44%
      height: 100%
</style>
