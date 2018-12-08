<template>
  <div class="width-form">
    <p style="font-size: 32px;">
      <span style="font-weight: 300">HYPER</span><span style="font-weight: 600">TUBE</span>
    </p>
    <p style="margin-bottom: 40px; font-style: italic">{{ $t('subtitle') }}</p>
    <form @submit.prevent="register">
      <input type="text" autocomplete="username" style="visibility: hidden; position: absolute">
      <div @click="uploadPicture" class="custom-file-upload">
        <input @change="createPicture" accept=".jpg, .jpeg, .png"
               ref="picture" type="file"/>
        <p>{{ $t('addPicLabel') }}</p>
        <img :src="user.pictureBase64" alt="profile picture uploaded"
             class="adapt-picture" v-show="user.pictureBase64">
      </div>
      <input :placeholder="$t('usernameLabel')" :title="$t('usernameStandard')"
             autofocus name="username" pattern="^[a-zA-Z0-9\.\-_ÆÐƎƏƐƔĲŊŒẞÞǷȜæðǝəɛɣĳŋœĸſßþƿȝĄƁÇĐƊĘĦĮƘŁØƠŞȘŢȚŦŲƯY̨Ƴąɓçđɗęħįƙłøơşșţțŧųưy̨ƴÁÀÂÄǍĂĀÃÅǺĄÆǼǢƁĆĊĈČÇĎḌĐƊÐÉÈĖÊËĚĔĒĘẸƎƏƐĠĜǦĞĢƔáàâäǎăāãåǻąæǽǣɓćċĉčçďḍđɗðéèėêëěĕēęẹǝəɛġĝǧğģɣĤḤĦIÍÌİÎÏǏĬĪĨĮỊĲĴĶƘĹĻŁĽĿNŃN̈ŇÑŅŊÓÒÔÖǑŎŌÕŐỌØǾƠŒĥḥħıíìiîïǐĭīĩįịĳĵķƙĸĺļłľŀŉńn̈ňñņŋóòôöǒŏōõőọøǿơœŔŘŖŚŜŠŞȘṢẞŤŢṬŦÞÚÙÛÜǓŬŪŨŰŮŲỤƯẂẀŴẄǷÝỲŶŸȲỸƳŹŻŽẒŕřŗſśŝšşșṣßťţṭŧþúùûüǔŭūũűůųụưẃẁŵẅƿýỳŷÿȳỹƴźżžẓ]{1,254}$"
             autocomplete="username"
             type="text" v-model="user.username">
      <input :placeholder="$t('emailLabel')" name="email"
            autocomplete="email"
             type="email" v-model="user.email">
      <div class="firstname-lastname">
        <input :placeholder="$t('firstnameLabel')"
               :title="$t('firstnameStandard')" pattern="^[a-zA-Z0-9\-ÆÐƎƏƐƔĲŊŒẞÞǷȜæðǝəɛɣĳŋœĸſßþƿȝĄƁÇĐƊĘĦĮƘŁØƠŞȘŢȚŦŲƯY̨Ƴąɓçđɗęħįƙłøơşșţțŧųưy̨ƴÁÀÂÄǍĂĀÃÅǺĄÆǼǢƁĆĊĈČÇĎḌĐƊÐÉÈĖÊËĚĔĒĘẸƎƏƐĠĜǦĞĢƔáàâäǎăāãåǻąæǽǣɓćċĉčçďḍđɗðéèėêëěĕēęẹǝəɛġĝǧğģɣĤḤĦIÍÌİÎÏǏĬĪĨĮỊĲĴĶƘĹĻŁĽĿNŃN̈ŇÑŅŊÓÒÔÖǑŎŌÕŐỌØǾƠŒĥḥħıíìiîïǐĭīĩįịĳĵķƙĸĺļłľŀŉńn̈ňñņŋóòôöǒŏōõőọøǿơœŔŘŖŚŜŠŞȘṢẞŤŢṬŦÞÚÙÛÜǓŬŪŨŰŮŲỤƯẂẀŴẄǷÝỲŶŸȲỸƳŹŻŽẒŕřŗſśŝšşșṣßťţṭŧþúùûüǔŭūũűůųụưẃẁŵẅƿýỳŷÿȳỹƴźżžẓ]{1,254}$"
               autocomplete="given-name"
               required style="min-width: 0"
               type="text" v-model="user.firstname">
        <div style="padding-left: 5px; padding-right: 5px"></div>
        <input :placeholder="$t('lastnameLabel')"
               :title="$t('lastnameStandard')"
               autocomplete="family-name"
               pattern="^[a-zA-Z0-9\-ÆÐƎƏƐƔĲŊŒẞÞǷȜæðǝəɛɣĳŋœĸſßþƿȝĄƁÇĐƊĘĦĮƘŁØƠŞȘŢȚŦŲƯY̨Ƴąɓçđɗęħįƙłøơşșţțŧųưy̨ƴÁÀÂÄǍĂĀÃÅǺĄÆǼǢƁĆĊĈČÇĎḌĐƊÐÉÈĖÊËĚĔĒĘẸƎƏƐĠĜǦĞĢƔáàâäǎăāãåǻąæǽǣɓćċĉčçďḍđɗðéèėêëěĕēęẹǝəɛġĝǧğģɣĤḤĦIÍÌİÎÏǏĬĪĨĮỊĲĴĶƘĹĻŁĽĿNŃN̈ŇÑŅŊÓÒÔÖǑŎŌÕŐỌØǾƠŒĥḥħıíìiîïǐĭīĩįịĳĵķƙĸĺļłľŀŉńn̈ňñņŋóòôöǒŏōõőọøǿơœŔŘŖŚŜŠŞȘṢẞŤŢṬŦÞÚÙÛÜǓŬŪŨŰŮŲỤƯẂẀŴẄǷÝỲŶŸȲỸƳŹŻŽẒŕřŗſśŝšşșṣßťţṭŧþúùûüǔŭūũűůųụưẃẁŵẅƿýỳŷÿȳỹƴźżžẓ]{1,254}$"
               required style="min-width: 0"
               type="text" v-model="user.lastname">
      </div>
      <input :placeholder="$t('passwordLabel')" :title="$t('passwordStandard')" name="password"
             pattern="(?=^.{4,}$)((?!.*\s)((?=(.*\d){1,})|(?=(.*[@#$%^&+=]){1,})))^.*$" required
             type="password"
             autocomplete="new-password"
             v-model="user.password">
      <input :placeholder="$t('confirmPasswordLabel')"
             :title="$t('passwordStandard')" name="password"
             pattern="(?=^.{4,}$)((?!.*\s)((?=(.*\d){1,})|(?=(.*[@#$%^&+=]){1,})))^.*$" required
             type="password"
             autocomplete="new-password"
             v-model="user.rePassword">
      <button type="submit">{{ $t('createLabel') }}</button>
    </form>
    <div class="footerForm">
      <router-link :to="{ name: 'recover' }" class="sublink" style="float: left">
        {{ $t('forgotLabel') }}
      </router-link>
      <router-link :to="{ name: 'index'}" class="sublink" style="float: right">
        {{ $t('accessLabel') }}
      </router-link>
    </div>
    <oauth-links></oauth-links>
  </div>
</template>

<script>
import CoversSlideshow from './CoversSlideshow';
import LanguageList from './LanguageList';
import OauthLinks from './OauthLinks';
import { getCurLanguage } from '../../static/js/language';
import api from '../lib/api';
import modal from '../lib/modal';

export default {
  components: {
    CoversSlideshow,
    LanguageList,
    OauthLinks,
  },
  data() {
    return {
      user: {
        pictureBase64: null,
        username: null,
        firstname: null,
        lastname: null,
        password: null,
        rePassword: null,
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
    getCurLanguage,
    updateTranslation(event) {
      if (typeof (Storage) !== 'undefined') {
        sessionStorage.setItem('language', event);
        this.curLanguage = event;
        this.$i18n.locale = this.curLanguage;
      }
    },
    uploadPicture() {
      this.$refs.picture.click();
    },
    createPicture(e) {
      const file = e.target.files[0];
      if (file && file.size > 5242880) {
        modal.open({ content: this.$t('errorSizeImage'), type: 'error' });
        return;
      }
      const reader = new FileReader();
      reader.onloadend = () => {
        this.user.pictureBase64 = reader.result;
      };
      if (file) {
        if (/\.(jpe?g|png|gif)$/i.test(file.name)) {
          reader.readAsDataURL(file);
        }
      }
    },
    register() {
      if (this.user.pictureBase64 === null) {
        modal.open({ content: this.$t('noPicture'), type: 'error' });
        return;
      }
      if (this.user.password !== this.user.rePassword) {
        modal.open({ content: this.$t('identicalPassword'), type: 'error' });
        return;
      }
      api.register({
        username: this.user.username,
        email: this.user.email,
        lastname: this.user.lastname,
        firstname: this.user.firstname,
        password: this.user.password,
        rePassword: this.user.rePassword,
        picture_base64: this.user.pictureBase64,
      }).then((res) => {
        if (res.status >= 500) {
          modal.open({ content: this.$t('errorServer'), type: 'warning' });
        } else if (res.status >= 400) {
          res.json().then((data) => {
            modal.open({
              content: this.$t(modal.errorFormatLanguage(data.error, res.status)),
              type: 'error',
            });
          });
        } else {
          this.$router.push('/?status=account_created');
        }
      }).catch(() => {
        modal.open({ content: this.$t('serverConnectionFailed'), type: 'warning' });
        // console.error('Error - Message:', err.message);
      });
    },
  },
};
</script>

<style scoped lang="sass">
  @import "Form"
  @import "Index"

  .adapt-picture
    position: absolute
    width: 100%
    height: 100%
    object-fit: cover
    top: 0
    left: 0

  input[type="file"]
    display: none

  .custom-file-upload
    position: relative
    display: inline-block
    background-color: white
    font-size: 15px
    font-family: 'Source Sans Pro', sans-serif
    border-radius: 2px
    box-shadow: 0 0 10px 0 rgba(0, 0, 0, 0.1)
    border: unset
    cursor: pointer
    width: fit-content
    padding-left: 10px
    padding-right: 10px
    height: 100px
    color: #D6D9E1
    line-height: 100px
    margin-bottom: 10px

  .firstname-lastname
    display: flex
    justify-content: space-between
    @media screen and (max-width: 640px)
      flex-wrap: wrap

  .width-form
    position: relative
    top: calc(50% - (678px / 2))
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
