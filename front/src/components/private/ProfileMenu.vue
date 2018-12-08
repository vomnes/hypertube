<template>
  <section>
    <div @click="toogleProfileMenu" class="header-profile-picture" title="My profile">
      <img :src="profile_picture"
           class="adapt-picture-circle">
    </div>
    <div v-if="showMenu">
      <div class="header-info">
        <div class="header-info-header">
          <h2>{{$t('profileLabel')}}</h2>
          <loading-animation v-if="loading" class="loadingProfile"/>
          <div class="displayFlex" v-else>
            <button style="" :class="{'hidden': hideSaveButton && editMode}"
                    @click="enableEdit"
                    ref="saveEditButton"
                    class="header-info-btn">{{editMode ? $t('saveLabel') : $t('editLabel')}}
            </button>

            <button @click="cancelEdit" class="header-info-btn"
                    ref="cancelButton"
                    style="background-color: #F0544E" v-if="editMode">{{ $t('cancelLabel') }}
            </button>
          </div>
          <img @click="logout" alt="logout" class="logout" src="../../static/icons/logout.svg">
        </div>
        <div class="header-info-body">
          <div>
            <img :class="{ clickable: editMode }" :src="profile_picture"
                 @click="changePhoto"
                 alt="profile image"
                 class="header-info-photo">
            <button @click="savePhoto" class="header-info-btn spacing-top-btn"
                    style="left: calc( 50% - (57.95px/2))"
                    v-if="photoChanged && editMode && !loadingChangePhoto">{{$t('saveLabel')}}
            </button>
            <loading-animation v-if="loadingChangePhoto" class="loadingProfile"/>
          </div>
          <div class="header-info-subbody">
            <p>{{ username }}</p>
            <input :class="[canSaveEmail ? 'normal' : 'error', 'spacing']"
                   :placeholder="$t('emailLabel')"
                   @input="checker"
                   autocomplete="email"
                   name="email" type="email"
                   v-if="email && editMode && !oauth"
                   v-model="email">
            <img
              :src="typeOauth(oauth)"
              alt="oauth picture"
              v-if="oauth">
            <p v-if="email && !editMode">{{ email }}</p>
            <p v-if="!editMode">{{ firstname }}&nbsp;{{ lastname }}</p>
            <input :class="[canSaveFirstName ? 'normal' : 'error', 'spacing']"
                   :placeholder="$t('firstnameLabel')"
                   @input="checker"
                   name="firstname"
                   v-if="editMode"
                   v-model="firstname">
            <input :class= "canSaveLastName ? 'normal' : 'error'"
                   :placeholder="$t('lastnameLabel')"
                   @input="checker"
                   name="lastname"
                   v-if="editMode"
                   v-model="lastname">
            <button @click="passwordEditing = !passwordEditing"
                    class="header-info-btn spacing-top-btn"
                    v-if="editMode && !oauth">{{ $t('editPWLabel') }}
            </button>
            <div class="password-container" v-if="!oauth && editMode && passwordEditing">
              <p>{{$t('changePasswordLabel')}}</p>
              <form @submit.prevent="updatePassword"  >
                <input type="text" autocomplete="username" style="visibility: hidden; position: absolute">
                <input :placeholder="$t('passwordLabel')"
                      pattern="(?=^.{4,}$)((?!.*\s)((?=(.*\d){1,})|(?=(.*[@#$%^&+=]){1,})))^.*$"
                      type="password"
                      autocomplete="new-password"
                      v-model="password">
                <input :placeholder="$t('confirmPasswordLabel')"
                      pattern="(?=^.{4,}$)((?!.*\s)((?=(.*\d){1,})|(?=(.*[@#$%^&+=]){1,})))^.*$"
                      type="password"
                      autocomplete="new-password"
                      v-model="new_password">
                <input :placeholder="$t('confirmPasswordLabel')"
                      pattern="(?=^.{4,}$)((?!.*\s)((?=(.*\d){1,})|(?=(.*[@#$%^&+=]){1,})))^.*$"
                      type="password"
                      autocomplete="new-password"
                      v-model="new_rePassword">
                <button type="submit" v-if="!loadingChangePassword"
                        class="header-info-btn spacing-top-btn">{{$t('saveLabel')}}
                </button>
                <loading-animation v-else class="loadingProfile"/>
              </form>
            </div>

            <div>
              <p style="display: inline-block">{{$t('languageLabel')}}</p>
              <img :src="getLanguageIcon()" class="language-icon">
              <div v-if="editMode">
                <img @click="setCurLanguage('en')" alt="uk"
                     class="language-icon clickable" src="../../static/icons/united-kingdom.svg">
                <img @click="setCurLanguage('fr')" alt="france"
                     class="language-icon clickable" src="../../static/icons/france.svg">
                <img @click="setCurLanguage('it')" alt="italy"
                     class="language-icon clickable" src="../../static/icons/italy.svg">
              </div>
            </div>
          </div>
        </div>
        <div class="arrow-close">
          <img @click="toogleProfileMenu" alt="close" class="arrow-close-img"
               src="../../static/icons/arrow.svg">
        </div>
      </div>
    </div>


  </section>
</template>

<script>
import swal from 'sweetalert';
import auth from '../../static/js/auth';
import lang from '../../static/js/language';
import profile from '../lib/api/index';
import modal from '../lib/modal';
import api from '../lib/api';
import LoadingAnimation from './module/LoadingAnimation';
import language from './../../static/js/language.js';

export default {
  components: {
    LoadingAnimation,
  },
  data() {
    return {
      editMode: false,
      passwordEditing: false,
      canSaveUsername: true,
      canSaveFirstName: true,
      canSaveLastName: true,
      canSaveEmail: true,
      username: '',
      usernameRef: '',
      profile_picture: '',
      firstname: '',
      firstnameRef: '',
      lastname: '',
      lastnameRef: '',
      email: '',
      emailRef: '',
      locale: '',
      oauth: '',
      password: '',
      new_password: '',
      new_rePassword: '',
      photoChanged: false,
      showMenu: false,
      loading: false,
      loadingChangePhoto: false,
      loadingChangePassword: false,
    };
  },
  methods: {
    cancelEdit(mouseEvent, save = false) {
      this.editMode = false;
      if (!save) {
        this.username = this.usernameRef;
        this.firstname = this.firstnameRef;
        this.lastname = this.lastnameRef;
        this.email = this.emailRef;
        this.$i18n.locale = this.locale;
      } else {
        this.usernameRef = this.username;
        this.firstnameRef = this.firstname;
        this.lastnameRef = this.lastname;
        this.emailRef = this.email;
        if (this.locale !== this.$i18n.locale) { lang.setCurLanguage(this.$i18n.locale); }
        this.locale = this.$i18n.locale;
      }
    },
    toogleProfileMenu() {
      this.$_bus.$emit('downloadMenuHide');
      if (this.showMenu) {
        this.$_bus.$emit('mobileHeaderBlack', false);
      } else {
        this.$_bus.$emit('mobileHeaderBlack', true);
      }
      this.showMenu = !this.showMenu;
    },
    typeOauth(type) {
      if (type === 'gplus') return 'data:image/svg+xml;utf8;base64,PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0iaXNvLTg4NTktMSI/Pgo8IS0tIEdlbmVyYXRvcjogQWRvYmUgSWxsdXN0cmF0b3IgMTguMC4wLCBTVkcgRXhwb3J0IFBsdWctSW4gLiBTVkcgVmVyc2lvbjogNi4wMCBCdWlsZCAwKSAgLS0+CjwhRE9DVFlQRSBzdmcgUFVCTElDICItLy9XM0MvL0RURCBTVkcgMS4xLy9FTiIgImh0dHA6Ly93d3cudzMub3JnL0dyYXBoaWNzL1NWRy8xLjEvRFREL3N2ZzExLmR0ZCI+CjxzdmcgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIiB4bWxuczp4bGluaz0iaHR0cDovL3d3dy53My5vcmcvMTk5OS94bGluayIgdmVyc2lvbj0iMS4xIiBpZD0iQ2FwYV8xIiB4PSIwcHgiIHk9IjBweCIgdmlld0JveD0iMCAwIDQ1NS43MyA0NTUuNzMiIHN0eWxlPSJlbmFibGUtYmFja2dyb3VuZDpuZXcgMCAwIDQ1NS43MyA0NTUuNzM7IiB4bWw6c3BhY2U9InByZXNlcnZlIiB3aWR0aD0iNTEycHgiIGhlaWdodD0iNTEycHgiPgo8cGF0aCBzdHlsZT0iZmlsbDojREQ0QjM5OyIgZD0iTTAsMHY0NTUuNzNoNDU1LjczVjBIMHogTTI2NS42NywyNDcuMDM3Yy03Ljc5Myw1MS4xOTQtNDUuOTYxLDgwLjU0My05NS4zNzYsODAuNTQzICBjLTU1LjUzMSwwLTEwMC41NTItNDUuMDIxLTEwMC41NTItMTAwLjU1MmMwLTU1LjUxNyw0NS4wMjEtMTAwLjUzOCwxMDAuNTUyLTEwMC41MzhjMjYuODYyLDAsNTAuMzk5LDkuNTg2LDY3LjUzMSwyNi4yMjYgIGwtMjguODU3LDI4Ljg1N2MtOS43NzMtOS44NDYtMjMuMTQ3LTE1LjA5NC0zOC42NzQtMTUuMDk0Yy0zMi42ODgsMC01OS4xODksMjcuODc0LTU5LjE4OSw2MC41NDggIGMwLDMyLjcwMywyNi41MDEsNTkuNzY4LDU5LjE4OSw1OS43NjhjMjcuMzk3LDAsNDguMTQ0LTEzLjI0Myw1NC4xMjktMzkuNzU4aC01NC4xMjl2LTQwLjM4aDk1LjEzMSAgYzEuMTQyLDYuNTA2LDEuNzIsMTMuMzE1LDEuNzIsMjAuMzdDMjY3LjE0NCwyMzQuMDI1LDI2Ni42MzgsMjQwLjY5LDI2NS42NywyNDcuMDM3eiBNMzg2LjQxOSwyMzQuNTE3aC0zNS4yMzN2MzUuMjE4SDMyNi4xNiAgdi0zNS4yMThoLTM1LjIzM3YtMjUuMDQxaDM1LjIzM3YtMzUuMjMzaDI1LjAyNnYzNS4yMzNoMzUuMjMzVjIzNC41MTd6Ii8+CjxnPgo8L2c+CjxnPgo8L2c+CjxnPgo8L2c+CjxnPgo8L2c+CjxnPgo8L2c+CjxnPgo8L2c+CjxnPgo8L2c+CjxnPgo8L2c+CjxnPgo8L2c+CjxnPgo8L2c+CjxnPgo8L2c+CjxnPgo8L2c+CjxnPgo8L2c+CjxnPgo8L2c+CjxnPgo8L2c+Cjwvc3ZnPgo=';
      if (type === 'facebook') return 'data:image/svg+xml;utf8;base64,PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0iaXNvLTg4NTktMSI/Pgo8IS0tIEdlbmVyYXRvcjogQWRvYmUgSWxsdXN0cmF0b3IgMTguMC4wLCBTVkcgRXhwb3J0IFBsdWctSW4gLiBTVkcgVmVyc2lvbjogNi4wMCBCdWlsZCAwKSAgLS0+CjwhRE9DVFlQRSBzdmcgUFVCTElDICItLy9XM0MvL0RURCBTVkcgMS4xLy9FTiIgImh0dHA6Ly93d3cudzMub3JnL0dyYXBoaWNzL1NWRy8xLjEvRFREL3N2ZzExLmR0ZCI+CjxzdmcgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIiB4bWxuczp4bGluaz0iaHR0cDovL3d3dy53My5vcmcvMTk5OS94bGluayIgdmVyc2lvbj0iMS4xIiBpZD0iQ2FwYV8xIiB4PSIwcHgiIHk9IjBweCIgdmlld0JveD0iMCAwIDQ1NS43MyA0NTUuNzMiIHN0eWxlPSJlbmFibGUtYmFja2dyb3VuZDpuZXcgMCAwIDQ1NS43MyA0NTUuNzM7IiB4bWw6c3BhY2U9InByZXNlcnZlIiB3aWR0aD0iNTEycHgiIGhlaWdodD0iNTEycHgiPgo8cGF0aCBzdHlsZT0iZmlsbDojM0E1NTlGOyIgZD0iTTAsMHY0NTUuNzNoMjQyLjcwNFYyNzkuNjkxaC01OS4zM3YtNzEuODY0aDU5LjMzdi02MC4zNTNjMC00My44OTMsMzUuNTgyLTc5LjQ3NSw3OS40NzUtNzkuNDc1ICBoNjIuMDI1djY0LjYyMmgtNDQuMzgyYy0xMy45NDcsMC0yNS4yNTQsMTEuMzA3LTI1LjI1NCwyNS4yNTR2NDkuOTUzaDY4LjUyMWwtOS40Nyw3MS44NjRoLTU5LjA1MVY0NTUuNzNINDU1LjczVjBIMHoiLz4KPGc+CjwvZz4KPGc+CjwvZz4KPGc+CjwvZz4KPGc+CjwvZz4KPGc+CjwvZz4KPGc+CjwvZz4KPGc+CjwvZz4KPGc+CjwvZz4KPGc+CjwvZz4KPGc+CjwvZz4KPGc+CjwvZz4KPGc+CjwvZz4KPGc+CjwvZz4KPGc+CjwvZz4KPGc+CjwvZz4KPC9zdmc+Cg==';
      if (type === 'telegram') return 'data:image/svg+xml;utf8;base64,PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0iaXNvLTg4NTktMSI/Pgo8IS0tIEdlbmVyYXRvcjogQWRvYmUgSWxsdXN0cmF0b3IgMTguMC4wLCBTVkcgRXhwb3J0IFBsdWctSW4gLiBTVkcgVmVyc2lvbjogNi4wMCBCdWlsZCAwKSAgLS0+CjwhRE9DVFlQRSBzdmcgUFVCTElDICItLy9XM0MvL0RURCBTVkcgMS4xLy9FTiIgImh0dHA6Ly93d3cudzMub3JnL0dyYXBoaWNzL1NWRy8xLjEvRFREL3N2ZzExLmR0ZCI+CjxzdmcgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIiB4bWxuczp4bGluaz0iaHR0cDovL3d3dy53My5vcmcvMTk5OS94bGluayIgdmVyc2lvbj0iMS4xIiBpZD0iQ2FwYV8xIiB4PSIwcHgiIHk9IjBweCIgdmlld0JveD0iMCAwIDQ1NS43MzEgNDU1LjczMSIgc3R5bGU9ImVuYWJsZS1iYWNrZ3JvdW5kOm5ldyAwIDAgNDU1LjczMSA0NTUuNzMxOyIgeG1sOnNwYWNlPSJwcmVzZXJ2ZSIgd2lkdGg9IjUxMnB4IiBoZWlnaHQ9IjUxMnB4Ij4KPGc+Cgk8cmVjdCB4PSIwIiB5PSIwIiBzdHlsZT0iZmlsbDojNjFBOERFOyIgd2lkdGg9IjQ1NS43MzEiIGhlaWdodD0iNDU1LjczMSIvPgoJPHBhdGggc3R5bGU9ImZpbGw6I0ZGRkZGRjsiIGQ9Ik0zNTguODQ0LDEwMC42TDU0LjA5MSwyMTkuMzU5Yy05Ljg3MSwzLjg0Ny05LjI3MywxOC4wMTIsMC44ODgsMjEuMDEybDc3LjQ0MSwyMi44NjhsMjguOTAxLDkxLjcwNiAgIGMzLjAxOSw5LjU3OSwxNS4xNTgsMTIuNDgzLDIyLjE4NSw1LjMwOGw0MC4wMzktNDAuODgybDc4LjU2LDU3LjY2NWM5LjYxNCw3LjA1NywyMy4zMDYsMS44MTQsMjUuNzQ3LTkuODU5bDUyLjAzMS0yNDguNzYgICBDMzgyLjQzMSwxMDYuMjMyLDM3MC40NDMsOTYuMDgsMzU4Ljg0NCwxMDAuNnogTTMyMC42MzYsMTU1LjgwNkwxNzkuMDgsMjgwLjk4NGMtMS40MTEsMS4yNDgtMi4zMDksMi45NzUtMi41MTksNC44NDcgICBsLTUuNDUsNDguNDQ4Yy0wLjE3OCwxLjU4LTIuMzg5LDEuNzg5LTIuODYxLDAuMjcxbC0yMi40MjMtNzIuMjUzYy0xLjAyNy0zLjMwOCwwLjMxMi02Ljg5MiwzLjI1NS04LjcxN2wxNjcuMTYzLTEwMy42NzYgICBDMzIwLjA4OSwxNDcuNTE4LDMyNC4wMjUsMTUyLjgxLDMyMC42MzYsMTU1LjgwNnoiLz4KPC9nPgo8Zz4KPC9nPgo8Zz4KPC9nPgo8Zz4KPC9nPgo8Zz4KPC9nPgo8Zz4KPC9nPgo8Zz4KPC9nPgo8Zz4KPC9nPgo8Zz4KPC9nPgo8Zz4KPC9nPgo8Zz4KPC9nPgo8Zz4KPC9nPgo8Zz4KPC9nPgo8Zz4KPC9nPgo8Zz4KPC9nPgo8Zz4KPC9nPgo8L3N2Zz4K';
      if (type === '42') return 'data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHZpZXdCb3g9IjAgMCA5NjEgOTYxIj48dGl0bGU+QXNzZXQgNDwvdGl0bGU+PGcgaWQ9IkxheWVyXzIiIGRhdGEtbmFtZT0iTGF5ZXIgMiI+PGcgaWQ9IkxheWVyXzEtMiIgZGF0YS1uYW1lPSJMYXllciAxIj48cmVjdCB3aWR0aD0iOTYxIiBoZWlnaHQ9Ijk2MSIgc3R5bGU9ImZpbGw6IzIzMWYyMCIvPjxnIGlkPSJzdmdfMiIgZGF0YS1uYW1lPSJzdmcgMiI+PHBvbHlnb24gaWQ9InBvbHlnb243IiBwb2ludHM9IjU4MC42IDM1Mi41IDcyMC41IDIxMi4yIDU4MC42IDIxMi4yIDU4MC42IDM1Mi41IiBzdHlsZT0iZmlsbDojZmZmIi8+PHBvbHlnb24gaWQ9InBvbHlnb24xMSIgcG9pbnRzPSI4NjAuOCA0OTIuNSA3MjAuNSA2MzIuNCA4NjAuOCA2MzIuNCA4NjAuOCA0OTIuNSIgc3R5bGU9ImZpbGw6I2ZmZiIvPjxwb2x5Z29uIGlkPSJwb2x5Z29uNSIgcG9pbnRzPSIxMDAuMiA2MDUuOCAzODAuNCA2MDUuOCAzODAuNCA3NDYuMiA1MjAuMiA3NDYuMiA1MjAuMiA0OTIuNSAyNDAuNSA0OTIuNSA1MjAuMiAyMTIuMiAzODAuNCAyMTIuMiAxMDAuMiA0OTIuNSAxMDAuMiA2MDUuOCIgc3R5bGU9ImZpbGw6I2ZmZiIvPjxwb2x5Z29uIGlkPSJwb2x5Z29uOSIgcG9pbnRzPSI3MjAuNSAzNTIuNSA1ODAuNiA0OTIuNSA1ODAuNiA2MzIuNCA3MjAuNSA2MzIuNCA3MjAuNSA0OTIuNSA4NjAuOCAzNTIuNSA4NjAuOCAyMTIuMiA3MjAuNSAyMTIuMiA3MjAuNSAzNTIuNSIgc3R5bGU9ImZpbGw6I2ZmZiIvPjwvZz48L2c+PC9nPjwvc3ZnPg==';
      return '';
    },
    enableEdit() {
      if (this.editMode) {
        if (this.saveProfile()) this.editMode = !this.editMode;
      } else {
        this.editMode = !this.editMode;
      }
    },
    canSave() {
      return this.canSaveFirstName && this.canSaveUsername &&
          this.canSaveLastName && this.canSaveEmail;
    },
    somethingToSave() {
      return !!(this.firstname !== this.firstnameRef ||
          this.lastname !== this.lastnameRef ||
          this.email !== this.emailRef ||
          this.locale !== this.$i18n.locale ||
          this.profile_picture !== this.profile_pictureRef);
    },
    saveProfile() {
      let langChange = false;
      if (this.canSave() && this.somethingToSave()) {
        if (this.locale !== this.$i18n.locale) {
          langChange = true;
        }
        this.cancelEdit(null, true);
        this.loading = true;
        profile.editData(auth.token(), {
          firstname: this.firstname,
          lastname: this.lastname,
          email: this.email,
          locale: this.$i18n.locale,
        })
          .then((res) => {
            this.loading = false;
            if (res.status === 200) {
              this.$_bus.$emit('langChange', this.locale);
            } else if (res.status >= 500) {
              modal.open({ content: this.$t('errorServer'), type: 'warning' });
            } else if (res.status >= 400) {
              res.json()
                .then((data) => {
                  modal.open({
                    content: this.$t(modal.errorFormatLanguage(data.error, res.status)),
                    type: 'error',
                  });
                });
            }
          })
          .catch(() => null);
      }
    },
    updatePassword() {
      const reg = new RegExp('(?=^.{4,}$)((?!.*\\s)((?=(.*\\d){1,})|(?=(.*[@#$%^&+=]){1,})))^.*$');
      if (reg.test(this.password) && reg.test(this.new_password) && reg.test(this.new_rePassword)) {
        if (this.new_password !== this.new_rePassword) {
          swal({ text: 'Error password not matching', icon: 'error' });
        } else {
          this.loadingChangePassword = true;
          profile.updatePassword(auth.token(), {
            password: this.password,
            new_password: this.new_password,
            new_rePassword: this.new_rePassword,
          }).then((res) => {
            this.loadingChangePassword = false;
            if (res.status >= 500) {
              modal.open({ content: this.$t('errorServer'), type: 'warning' });
              return null;
            } else if (res.status >= 400) {
              res.json()
                .then((data) => {
                  modal.open({
                    content: this.$t(modal.errorFormatLanguage(data.error, res.status)),
                    type: 'error',
                  });
                  return null;
                });
            }
            if (res.status === 200) {
              this.passwordEditing = false;
              this.password = '';
              this.new_password = '';
              this.new_rePassword = '';
            }
          })
            .catch(() => null);
        }
      } else {
        swal({ text: 'Check password rules', icon: 'error' });
      }
    },
    changePhoto() {
      if (this.editMode) {
        const input = document.createElement('input');
        input.type = 'file';
        input.accept = '.jpg, .jpeg, .png';
        input.click();
        input.onchange = () => {
          const file = input.files[0];
          if (file && file.size > 5242880) {
            modal.open({ content: this.$t('errorSizeImage'), type: 'error' });
            return;
          }
          const reader = new FileReader();
          reader.onloadend = () => {
            this.profile_picture = reader.result;
            this.photoChanged = true;
          };
          reader.readAsDataURL(file);
        };
      }
    },
    savePhoto() {
      if (this.photoChanged) {
        this.loadingChangePhoto = true;
        profile.updatePhoto(auth.token(), { picture_base64: this.profile_picture }).then((res) => {
          this.loadingChangePhoto = false;
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
            this.photoChanged = false;
          }
        })
          .catch(() => null);
      }
    },
    checker(event) {
      if (event.target.name === 'username') {
        this.canSaveUsername = RegExp('^[a-zA-Z0-9\.\-_ÆÐƎƏƐƔĲŊŒẞÞǷȜæðǝəɛɣĳŋœĸſßþƿȝĄƁÇĐƊĘĦĮƘŁØƠŞȘŢȚŦŲƯY̨Ƴąɓçđɗęħįƙłøơşșţțŧųưy̨ƴÁÀÂÄǍĂĀÃÅǺĄÆǼǢƁĆĊĈČÇĎḌĐƊÐÉÈĖÊËĚĔĒĘẸƎƏƐĠĜǦĞĢƔáàâäǎăāãåǻąæǽǣɓćċĉčçďḍđɗðéèėêëěĕēęẹǝəɛġĝǧğģɣĤḤĦIÍÌİÎÏǏĬĪĨĮỊĲĴĶƘĹĻŁĽĿNŃN̈ŇÑŅŊÓÒÔÖǑŎŌÕŐỌØǾƠŒĥḥħıíìiîïǐĭīĩįịĳĵķƙĸĺļłľŀŉńn̈ňñņŋóòôöǒŏōõőọøǿơœŔŘŖŚŜŠŞȘṢẞŤŢṬŦÞÚÙÛÜǓŬŪŨŰŮŲỤƯẂẀŴẄǷÝỲŶŸȲỸƳŹŻŽẒŕřŗſśŝšşșṣßťţṭŧþúùûüǔŭūũűůųụưẃẁŵẅƿýỳŷÿȳỹƴźżžẓ]+$')
          .test(event.target.value) && event.target.length >= 1 && event.test.length <= 254;
      }
      if (event.target.name === 'firstname') {
        this.canSaveFirstName = RegExp('^[a-zA-Z0-9\-ÆÐƎƏƐƔĲŊŒẞÞǷȜæðǝəɛɣĳŋœĸſßþƿȝĄƁÇĐƊĘĦĮƘŁØƠŞȘŢȚŦŲƯY̨Ƴąɓçđɗęħįƙłøơşșţțŧųưy̨ƴÁÀÂÄǍĂĀÃÅǺĄÆǼǢƁĆĊĈČÇĎḌĐƊÐÉÈĖÊËĚĔĒĘẸƎƏƐĠĜǦĞĢƔáàâäǎăāãåǻąæǽǣɓćċĉčçďḍđɗðéèėêëěĕēęẹǝəɛġĝǧğģɣĤḤĦIÍÌİÎÏǏĬĪĨĮỊĲĴĶƘĹĻŁĽĿNŃN̈ŇÑŅŊÓÒÔÖǑŎŌÕŐỌØǾƠŒĥḥħıíìiîïǐĭīĩįịĳĵķƙĸĺļłľŀŉńn̈ňñņŋóòôöǒŏōõőọøǿơœŔŘŖŚŜŠŞȘṢẞŤŢṬŦÞÚÙÛÜǓŬŪŨŰŮŲỤƯẂẀŴẄǷÝỲŶŸȲỸƳŹŻŽẒŕřŗſśŝšşșṣßťţṭŧþúùûüǔŭūũűůųụưẃẁŵẅƿýỳŷÿȳỹƴźżžẓ]+$')
          .test(event.target.value) && event.target.value.length >= 1 && event.target.value.length <= 254;
      }
      if (event.target.name === 'lastname') {
        this.canSaveLastName = RegExp('^[a-zA-Z0-9\-ÆÐƎƏƐƔĲŊŒẞÞǷȜæðǝəɛɣĳŋœĸſßþƿȝĄƁÇĐƊĘĦĮƘŁØƠŞȘŢȚŦŲƯY̨Ƴąɓçđɗęħįƙłøơşșţțŧųưy̨ƴÁÀÂÄǍĂĀÃÅǺĄÆǼǢƁĆĊĈČÇĎḌĐƊÐÉÈĖÊËĚĔĒĘẸƎƏƐĠĜǦĞĢƔáàâäǎăāãåǻąæǽǣɓćċĉčçďḍđɗðéèėêëěĕēęẹǝəɛġĝǧğģɣĤḤĦIÍÌİÎÏǏĬĪĨĮỊĲĴĶƘĹĻŁĽĿNŃN̈ŇÑŅŊÓÒÔÖǑŎŌÕŐỌØǾƠŒĥḥħıíìiîïǐĭīĩįịĳĵķƙĸĺļłľŀŉńn̈ňñņŋóòôöǒŏōõőọøǿơœŔŘŖŚŜŠŞȘṢẞŤŢṬŦÞÚÙÛÜǓŬŪŨŰŮŲỤƯẂẀŴẄǷÝỲŶŸȲỸƳŹŻŽẒŕřŗſśŝšşșṣßťţṭŧþúùûüǔŭūũűůųụưẃẁŵẅƿýỳŷÿȳỹƴźżžẓ]+$')
          .test(event.target.value) && event.target.value.length >= 1 && event.target.value.length <= 254;
      }
      if (event.target.name === 'email') {
        this.canSaveEmail = RegExp('^[a-zA-Z0-9.!#$%&\'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$')
          .test(event.target.value);
      }
    },
    getLanguageIcon() {
      if (this.$i18n.locale === 'en') return 'data:image/svg+xml;utf8;base64,PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0iaXNvLTg4NTktMSI/Pgo8IS0tIEdlbmVyYXRvcjogQWRvYmUgSWxsdXN0cmF0b3IgMTkuMC4wLCBTVkcgRXhwb3J0IFBsdWctSW4gLiBTVkcgVmVyc2lvbjogNi4wMCBCdWlsZCAwKSAgLS0+CjxzdmcgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIiB4bWxuczp4bGluaz0iaHR0cDovL3d3dy53My5vcmcvMTk5OS94bGluayIgdmVyc2lvbj0iMS4xIiBpZD0iTGF5ZXJfMSIgeD0iMHB4IiB5PSIwcHgiIHZpZXdCb3g9IjAgMCA1MTIgNTEyIiBzdHlsZT0iZW5hYmxlLWJhY2tncm91bmQ6bmV3IDAgMCA1MTIgNTEyOyIgeG1sOnNwYWNlPSJwcmVzZXJ2ZSIgd2lkdGg9IjUxMnB4IiBoZWlnaHQ9IjUxMnB4Ij4KPGNpcmNsZSBzdHlsZT0iZmlsbDojRjBGMEYwOyIgY3g9IjI1NiIgY3k9IjI1NiIgcj0iMjU2Ii8+CjxnPgoJPHBhdGggc3R5bGU9ImZpbGw6IzAwNTJCNDsiIGQ9Ik01Mi45MiwxMDAuMTQyYy0yMC4xMDksMjYuMTYzLTM1LjI3Miw1Ni4zMTgtNDQuMTAxLDg5LjA3N2gxMzMuMTc4TDUyLjkyLDEwMC4xNDJ6Ii8+Cgk8cGF0aCBzdHlsZT0iZmlsbDojMDA1MkI0OyIgZD0iTTUwMy4xODEsMTg5LjIxOWMtOC44MjktMzIuNzU4LTIzLjk5My02Mi45MTMtNDQuMTAxLTg5LjA3NmwtODkuMDc1LDg5LjA3Nkg1MDMuMTgxeiIvPgoJPHBhdGggc3R5bGU9ImZpbGw6IzAwNTJCNDsiIGQ9Ik04LjgxOSwzMjIuNzg0YzguODMsMzIuNzU4LDIzLjk5Myw2Mi45MTMsNDQuMTAxLDg5LjA3NWw4OS4wNzQtODkuMDc1TDguODE5LDMyMi43ODRMOC44MTksMzIyLjc4NCAgIHoiLz4KCTxwYXRoIHN0eWxlPSJmaWxsOiMwMDUyQjQ7IiBkPSJNNDExLjg1OCw1Mi45MjFjLTI2LjE2My0yMC4xMDktNTYuMzE3LTM1LjI3Mi04OS4wNzYtNDQuMTAydjEzMy4xNzdMNDExLjg1OCw1Mi45MjF6Ii8+Cgk8cGF0aCBzdHlsZT0iZmlsbDojMDA1MkI0OyIgZD0iTTEwMC4xNDIsNDU5LjA3OWMyNi4xNjMsMjAuMTA5LDU2LjMxOCwzNS4yNzIsODkuMDc2LDQ0LjEwMlYzNzAuMDA1TDEwMC4xNDIsNDU5LjA3OXoiLz4KCTxwYXRoIHN0eWxlPSJmaWxsOiMwMDUyQjQ7IiBkPSJNMTg5LjIxNyw4LjgxOWMtMzIuNzU4LDguODMtNjIuOTEzLDIzLjk5My04OS4wNzUsNDQuMTAxbDg5LjA3NSw4OS4wNzVWOC44MTl6Ii8+Cgk8cGF0aCBzdHlsZT0iZmlsbDojMDA1MkI0OyIgZD0iTTMyMi43ODMsNTAzLjE4MWMzMi43NTgtOC44Myw2Mi45MTMtMjMuOTkzLDg5LjA3NS00NC4xMDFsLTg5LjA3NS04OS4wNzVWNTAzLjE4MXoiLz4KCTxwYXRoIHN0eWxlPSJmaWxsOiMwMDUyQjQ7IiBkPSJNMzcwLjAwNSwzMjIuNzg0bDg5LjA3NSw4OS4wNzZjMjAuMTA4LTI2LjE2MiwzNS4yNzItNTYuMzE4LDQ0LjEwMS04OS4wNzZIMzcwLjAwNXoiLz4KPC9nPgo8Zz4KCTxwYXRoIHN0eWxlPSJmaWxsOiNEODAwMjc7IiBkPSJNNTA5LjgzMywyMjIuNjA5aC0yMjAuNDRoLTAuMDAxVjIuMTY3QzI3OC40NjEsMC43NDQsMjY3LjMxNywwLDI1NiwwICAgYy0xMS4zMTksMC0yMi40NjEsMC43NDQtMzMuMzkxLDIuMTY3djIyMC40NHYwLjAwMUgyLjE2N0MwLjc0NCwyMzMuNTM5LDAsMjQ0LjY4MywwLDI1NmMwLDExLjMxOSwwLjc0NCwyMi40NjEsMi4xNjcsMzMuMzkxICAgaDIyMC40NGgwLjAwMXYyMjAuNDQyQzIzMy41MzksNTExLjI1NiwyNDQuNjgxLDUxMiwyNTYsNTEyYzExLjMxNywwLDIyLjQ2MS0wLjc0MywzMy4zOTEtMi4xNjd2LTIyMC40NHYtMC4wMDFoMjIwLjQ0MiAgIEM1MTEuMjU2LDI3OC40NjEsNTEyLDI2Ny4zMTksNTEyLDI1NkM1MTIsMjQ0LjY4Myw1MTEuMjU2LDIzMy41MzksNTA5LjgzMywyMjIuNjA5eiIvPgoJPHBhdGggc3R5bGU9ImZpbGw6I0Q4MDAyNzsiIGQ9Ik0zMjIuNzgzLDMyMi43ODRMMzIyLjc4MywzMjIuNzg0TDQzNy4wMTksNDM3LjAyYzUuMjU0LTUuMjUyLDEwLjI2Ni0xMC43NDMsMTUuMDQ4LTE2LjQzNSAgIGwtOTcuODAyLTk3LjgwMmgtMzEuNDgyVjMyMi43ODR6Ii8+Cgk8cGF0aCBzdHlsZT0iZmlsbDojRDgwMDI3OyIgZD0iTTE4OS4yMTcsMzIyLjc4NGgtMC4wMDJMNzQuOTgsNDM3LjAxOWM1LjI1Miw1LjI1NCwxMC43NDMsMTAuMjY2LDE2LjQzNSwxNS4wNDhsOTcuODAyLTk3LjgwNCAgIFYzMjIuNzg0eiIvPgoJPHBhdGggc3R5bGU9ImZpbGw6I0Q4MDAyNzsiIGQ9Ik0xODkuMjE3LDE4OS4yMTl2LTAuMDAyTDc0Ljk4MSw3NC45OGMtNS4yNTQsNS4yNTItMTAuMjY2LDEwLjc0My0xNS4wNDgsMTYuNDM1bDk3LjgwMyw5Ny44MDMgICBIMTg5LjIxN3oiLz4KCTxwYXRoIHN0eWxlPSJmaWxsOiNEODAwMjc7IiBkPSJNMzIyLjc4MywxODkuMjE5TDMyMi43ODMsMTg5LjIxOUw0MzcuMDIsNzQuOTgxYy01LjI1Mi01LjI1NC0xMC43NDMtMTAuMjY2LTE2LjQzNS0xNS4wNDcgICBsLTk3LjgwMiw5Ny44MDNWMTg5LjIxOXoiLz4KPC9nPgo8Zz4KPC9nPgo8Zz4KPC9nPgo8Zz4KPC9nPgo8Zz4KPC9nPgo8Zz4KPC9nPgo8Zz4KPC9nPgo8Zz4KPC9nPgo8Zz4KPC9nPgo8Zz4KPC9nPgo8Zz4KPC9nPgo8Zz4KPC9nPgo8Zz4KPC9nPgo8Zz4KPC9nPgo8Zz4KPC9nPgo8Zz4KPC9nPgo8L3N2Zz4K';
      if (this.$i18n.locale === 'fr') return 'data:image/svg+xml;utf8;base64,PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0iaXNvLTg4NTktMSI/Pgo8IS0tIEdlbmVyYXRvcjogQWRvYmUgSWxsdXN0cmF0b3IgMTkuMC4wLCBTVkcgRXhwb3J0IFBsdWctSW4gLiBTVkcgVmVyc2lvbjogNi4wMCBCdWlsZCAwKSAgLS0+CjxzdmcgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIiB4bWxuczp4bGluaz0iaHR0cDovL3d3dy53My5vcmcvMTk5OS94bGluayIgdmVyc2lvbj0iMS4xIiBpZD0iTGF5ZXJfMSIgeD0iMHB4IiB5PSIwcHgiIHZpZXdCb3g9IjAgMCA1MTIgNTEyIiBzdHlsZT0iZW5hYmxlLWJhY2tncm91bmQ6bmV3IDAgMCA1MTIgNTEyOyIgeG1sOnNwYWNlPSJwcmVzZXJ2ZSIgd2lkdGg9IjUxMnB4IiBoZWlnaHQ9IjUxMnB4Ij4KPGNpcmNsZSBzdHlsZT0iZmlsbDojRjBGMEYwOyIgY3g9IjI1NiIgY3k9IjI1NiIgcj0iMjU2Ii8+CjxwYXRoIHN0eWxlPSJmaWxsOiNEODAwMjc7IiBkPSJNNTEyLDI1NmMwLTExMC4wNzEtNjkuNDcyLTIwMy45MDYtMTY2Ljk1Ny0yNDAuMDc3djQ4MC4xNTVDNDQyLjUyOCw0NTkuOTA2LDUxMiwzNjYuMDcxLDUxMiwyNTZ6Ii8+CjxwYXRoIHN0eWxlPSJmaWxsOiMwMDUyQjQ7IiBkPSJNMCwyNTZjMCwxMTAuMDcxLDY5LjQ3MywyMDMuOTA2LDE2Ni45NTcsMjQwLjA3N1YxNS45MjNDNjkuNDczLDUyLjA5NCwwLDE0NS45MjksMCwyNTZ6Ii8+CjxnPgo8L2c+CjxnPgo8L2c+CjxnPgo8L2c+CjxnPgo8L2c+CjxnPgo8L2c+CjxnPgo8L2c+CjxnPgo8L2c+CjxnPgo8L2c+CjxnPgo8L2c+CjxnPgo8L2c+CjxnPgo8L2c+CjxnPgo8L2c+CjxnPgo8L2c+CjxnPgo8L2c+CjxnPgo8L2c+Cjwvc3ZnPgo=';
      if (this.$i18n.locale === 'it') return 'data:image/svg+xml;utf8;base64,PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0iaXNvLTg4NTktMSI/Pgo8IS0tIEdlbmVyYXRvcjogQWRvYmUgSWxsdXN0cmF0b3IgMTkuMC4wLCBTVkcgRXhwb3J0IFBsdWctSW4gLiBTVkcgVmVyc2lvbjogNi4wMCBCdWlsZCAwKSAgLS0+CjxzdmcgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIiB4bWxuczp4bGluaz0iaHR0cDovL3d3dy53My5vcmcvMTk5OS94bGluayIgdmVyc2lvbj0iMS4xIiBpZD0iTGF5ZXJfMSIgeD0iMHB4IiB5PSIwcHgiIHZpZXdCb3g9IjAgMCA1MTIgNTEyIiBzdHlsZT0iZW5hYmxlLWJhY2tncm91bmQ6bmV3IDAgMCA1MTIgNTEyOyIgeG1sOnNwYWNlPSJwcmVzZXJ2ZSIgd2lkdGg9IjUxMnB4IiBoZWlnaHQ9IjUxMnB4Ij4KPGNpcmNsZSBzdHlsZT0iZmlsbDojRjBGMEYwOyIgY3g9IjI1NiIgY3k9IjI1NiIgcj0iMjU2Ii8+CjxwYXRoIHN0eWxlPSJmaWxsOiNEODAwMjc7IiBkPSJNNTEyLDI1NmMwLTExMC4wNzEtNjkuNDcyLTIwMy45MDYtMTY2Ljk1Ny0yNDAuMDc3djQ4MC4xNTVDNDQyLjUyOCw0NTkuOTA2LDUxMiwzNjYuMDcxLDUxMiwyNTZ6Ii8+CjxwYXRoIHN0eWxlPSJmaWxsOiM2REE1NDQ7IiBkPSJNMCwyNTZjMCwxMTAuMDcxLDY5LjQ3MiwyMDMuOTA2LDE2Ni45NTcsMjQwLjA3N1YxNS45MjNDNjkuNDcyLDUyLjA5NCwwLDE0NS45MjksMCwyNTZ6Ii8+CjxnPgo8L2c+CjxnPgo8L2c+CjxnPgo8L2c+CjxnPgo8L2c+CjxnPgo8L2c+CjxnPgo8L2c+CjxnPgo8L2c+CjxnPgo8L2c+CjxnPgo8L2c+CjxnPgo8L2c+CjxnPgo8L2c+CjxnPgo8L2c+CjxnPgo8L2c+CjxnPgo8L2c+CjxnPgo8L2c+Cjwvc3ZnPgo=';
      return '';
    },
    setCurLanguage(lang) {
      this.$i18n.locale = lang;
    },
    logout() {
      window.document.location.href = '/';
      auth.logOut();
    },
  },
  mounted() {
    this.$_bus.$on('profileMenuHide', () => {
      this.showMenu = false;
    });
    api.getProfile(auth.token())
      .then((res) => {
        if (res.status >= 500) {
          throw new Error('Server side error');
        } else if (res.status >= 400) {
          auth.logOut();
          window.document.location.href = '/';
        } else {
          res.json()
            .then((data) => {
              if (data) {
                if (data.username) {
                  this.username = data.username;
                  this.usernameRef = data.username;
                }
                if (data.email) {
                  this.email = data.email;
                  this.emailRef = data.email;
                }
                if (data.firstname) {
                  this.firstname = data.firstname;
                  this.firstnameRef = data.firstname;
                }
                if (data.lastname) {
                  this.lastname = data.lastname;
                  this.lastnameRef = data.lastname;
                }
                if (data.oauth) {
                  this.oauth = data.oauth;
                }
                if (data.profile_picture) {
                  this.profile_picture = data.profile_picture;
                  this.profile_pictureRef = data.profile_picture;
                }
                if (data.locale) {
                  this.locale = data.locale;
                  if (this.locale !== this.$i18n.locale) {
                    lang.setCurLanguage(data.locale);
                  }
                  this.$i18n.locale = data.locale;
                }
              }
            });
        }
      })
      .catch((err) => {
        console.error('Error - Message:', err.message);
      });
  },
  computed: {
    hideSaveButton() {
      const haveModif = (this.username !== this.usernameRef ||
          this.firstname !== this.firstnameRef ||
          this.lastname !== this.lastnameRef ||
          this.email !== this.emailRef ||
          this.locale !== this.$i18n.locale);
      return (!haveModif || !this.canSave());
    },
  },
};
</script>

<style scoped lang="sass">
  @import "ProfileMenu"
</style>
