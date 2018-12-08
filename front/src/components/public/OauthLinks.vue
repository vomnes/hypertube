<template lang="html">
  <div>
    <hr :data-content="$t('alternativeLabel')" class="hr-text">
    <div style="display: flex; justify-content: center">
      <a :href="`${ apiDomain }/api/auth/42`">
        <div class="auth-icon-container" style="background-color: #9461F7" title="42">
          <img alt="42 Logo" class="auth-icon" src="../../static/icons/42_Logo.svg">
        </div>
      </a>
      <a :href="`${ apiDomain }/api/auth/facebook`">
        <div class="auth-icon-container" style="background-color: #7C4CF3" title="Facebook">
          <img alt="Facebook Logo"
               class="auth-icon" src="../../static/icons/F_Logo.svg">
        </div>
      </a>
      <a :href="`${ apiDomain }/api/auth/gplus`">
        <div class="auth-icon-container" style="background-color: #6535CB" title="Google">
          <img alt="Google Logo"
               class="auth-icon" src="../../static/icons/G_Logo.svg">
        </div>
      </a>
    </div>
  </div>
</template>

<script>
import config from '../lib/api/config';

export default {
  data() {
    return {
      apiDomain: config.apiDomain,
    };
  },
};
</script>

<style scoped lang="sass">
  .hr-text
    margin-top: 40px
    clear: both
    line-height: 1em
    position: relative
    outline: 0
    border: 0
    color: black
    text-align: center
    height: 1.5em
    opacity: .5

    &:before
      content: ''
      position: absolute
      background: #384064
      left: 37%
      top: 50%
      width: 26%
      height: 1px

    &:after
      font-family: "Source Sans Pro", sans-serif
      font-size: 15px
      font-weight: 300
      content: attr(data-content)
      position: relative
      display: inline-block
      padding: 0 .5em
      line-height: 1.5em
      color: black
      background-color: #fcfcfa

  .auth-icon
    height: 20px
    width: 20px
    margin-top: 10px

  .auth-icon-container
    user-select: none
    width: 40px
    height: 40px
    margin: 5px
    border-radius: 3px
    box-shadow: 0 0 10px 0 rgba(0, 0, 0, 0.1)

    &:active
      transform: scale(0.90)
      transition: linear 0.05s
</style>
