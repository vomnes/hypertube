<template>
  <div @click="toggleList = !toggleList" class="lang-area" title="Change language">
    <ul class="lang-data">
      <li>{{ currentLanguage | uppercase }}</li>
      <li :key="index"
          @click="selectLang(lang)"
          class="lang-select"
          v-for="(lang, index) in listLanguages"
          v-if="toggleList">
        {{ lang | uppercase }}
      </li>
    </ul>
    <img alt="list languages" class="lang-list"
         src="../../static/icons/language-icon.svg">
  </div>
</template>

<script>
export default {
  props: {
    language: String,
  },
  data() {
    return {
      currentLanguage: this.language,
      toggleList: false,
      languageList: ['en', 'fr', 'it'],
    };
  },
  computed: {
    listLanguages() {
      return this.languageList.filter(lang => lang !== this.currentLanguage);
    },
  },
  methods: {
    selectLang(lang) {
      this.currentLanguage = lang;
      this.$emit('languageChanged', this.currentLanguage);
    },
  },
};
</script>

<style scoped lang="sass">
  @import "LanguageList"
</style>
