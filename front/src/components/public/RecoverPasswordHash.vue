<template>
  <div class="form-side">
    <div class="logo"></div>
    <div class="width-form">
      <p style="font-size: 32px;">
        <span style="font-weight: 300">HYPER</span><span style="font-weight: 600">TUBE</span>
      </p>
      <p style="margin-bottom: 40px; font-style: italic">{{ $t('subtitle') }}</p>
      <form @submit.prevent="recoverPassword">
        <input :placeholder="$t('passwordLabel')" :title="$t('passwordStandard')" autofocus
               name="password"
               pattern="(?=^.{4,}$)((?!.*\s)((?=(.*\d){1,})|(?=(.*[@#$%^&+=]){1,})))^.*$" required
               type="password"
               autocomplete="new-password"
               v-model="password">
        <input :placeholder="$t('rePasswordLabel')" :title="$t('passwordStandard')" autofocus
               name="password"
               pattern="(?=^.{4,}$)((?!.*\s)((?=(.*\d){1,})|(?=(.*[@#$%^&+=]){1,})))^.*$" required
               type="password"
               autocomplete="new-password"
               v-model="rePassword">
        <button type="submit">{{ $t('updatePasswordLabel') }}</button>
      </form>
    </div>
  </div>
</template>

<script>
import api from '../lib/api';
import modal from '../lib/modal';

export default {
  data() {
    return {
      password: null,
      rePassword: null,
    };
  },
  props: ['lang'],
  watch: {
    lang() {
      this.$i18n.locale = this.lang;
    },
  },
  methods: {
    recoverPassword() {
      if (this.password !== this.rePassword) {
        modal.open({ content: this.$t('identicalPassword'), type: 'error' });
        return;
      }
      api.forgotPasswordChange({
        randomToken: this.$route.params.hash,
        password: this.password,
        rePassword: this.rePassword,
      })
        .then((res) => {
          if (res.status >= 500) {
            modal.open({ content: this.$t('errorServer'), type: 'warning' });
          } else if (res.status >= 400) {
            if (res.status === 400) {
              this.$router.push('/?status=failed_reset_password');
            } else {
              res.json()
                .then((data) => {
                  modal.open({
                    content: this.$t(modal.errorFormatLanguage(data.error, res.status)),
                    type: 'error',
                  });
                });
            }
          } else {
            this.$router.push('/?status=password_reset');
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
  @import "Index"
  @import "Form"
  @import "LanguageList"

  .width-form
    position: relative
    top: calc(50% - (202px / 2))
    margin: auto
    width: 25vw

  .movie-cover
    width: 100%
    height: 100%
    object-fit: cover

  .cover-fade-enter-active, .cover-fade-leave-active
    transition: opacity 0.75s

  .cover-fade-enter, .cover-fade-leave-active
    opacity: 0

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
