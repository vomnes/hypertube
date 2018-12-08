<template>
  <div class="width-form">
    <p style="font-size: 32px;">
      <span style="font-weight: 300">HYPER</span><span style="font-weight: 600">TUBE</span>
    </p>
    <p style="margin-bottom: 40px; font-style: italic">{{ $t('subtitle') }}</p>
    <form @submit.prevent="login">
      <input type="text" autocomplete="username" style="visibility: hidden; position: absolute">
      <input :placeholder="$t('usernameLabel')" autoComplete="username"
             autofocus name="username" required
             autocomplete="username"
             type="text" v-model="user.login">
      <input :placeholder="$t('passwordLabel')"
             autocomplete="current-password" name="password" required
             type="password" v-model="user.password">
      <button type="submit">{{ $t('btnLabel') }}</button>
    </form>
    <div class="footerForm">
      <router-link :to="{ name: 'recover' }" class="sublink" style="float: left">
        {{ $t('forgotLabel') }}
      </router-link>
      <router-link :to="{ name: 'register' }" class="sublink" style="float: right">
        {{ $t('createLabel') }}
      </router-link>
    </div>
    <oauth-links></oauth-links>
  </div>
</template>

<script>
import OauthLinks from './OauthLinks';
import api from '../lib/api';
import modal from '../lib/modal';
import auth from '../../static/js/auth';

export default {
  components: {
    OauthLinks,
  },
  data() {
    return {
      user: {
        login: null,
        password: null,
      },
    };
  },
  props: ['lang'],
  watch: {
    lang() {
      this.$i18n.locale = this.lang;
    },
  },
  methods: {
    login() {
      api.login({
        username: this.user.login,
        password: this.user.password,
      })
        .then((res) => {
          if (res.status >= 500) {
            modal.open({ content: this.$t('errorServer'), type: 'warning' });
            throw new Error('Server side error');
          } else if (res.status >= 400) {
            res.json()
              .then((data) => {
                modal.open({
                  content: this.$t(modal.errorFormatLanguage(data.error, res.status)),
                  type: 'error',
                });
              });
          } else {
            res.json()
              .then((data) => {
                if (data && data.token) {
                  auth.storeToken(data.token);
                  this.$router.push('/gallery');
                } else {
                  throw new Error('Login call body response is invalid');
                }
              });
          }
        })
        .catch((err) => {
          modal.open({ content: this.$t('serverConnectionFailed'), type: 'warning' });
          console.error('Error - Message:', err.message);
        });
    },
  },
};
</script>

<style scoped lang="sass">
  @import "Form"

  .width-form
    position: relative
    top: calc(50% - (257px / 2))
    margin: auto
    width: 25vw

  @media screen and (max-width: 640px)
    p
      color: #F8F9FD
      text-shadow: 1px 1px 3px rgba(150, 150, 150, 1)

    .width-form
      position: relative
      width: 80%
      top: -36%
      height: 100%
</style>
