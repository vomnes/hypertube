<template lang="html">
  <div>
    <div class="header">
      <router-link to='/gallery'>
        <img alt="logo" class="logo" src="../../static/logo.svg">
      </router-link>
      <router-link class="header-title" to="/gallery">
        <span style="font-weight: 300">HYPER</span><span style="font-weight: 600">TUBE</span>
      </router-link>
      <div id="buttonsHeaderContainer">
        <router-link to="/search">
          <img alt="search"
            class="search"
            src="../../static/icons/magnifying-glass.svg"
            title="Search for a movie">
        </router-link>
        <download-menu></download-menu>
        <profile-menu></profile-menu>
      </div>
    </div>
    <div class="headerFilterSpace" v-if="headerFilterSpace"></div>
    <div @click="closePopup" class="blackModal" v-if="mobileBlackModal"></div>
  </div>
</template>

<script>
import ProfileMenu from './ProfileMenu';
import DownloadMenu from './DownloadMenu';

export default {
  components: {
    ProfileMenu,
    DownloadMenu,
  },
  mounted() {
    this.$_bus.$on('headerSpaceToogle', (value = !this.headerFilterSpace) => {
      this.headerFilterSpace = value;
    });
    this.$_bus.$on('mobileHeaderBlack', (value = !this.mobileBlackModal) => {
      this.mobileBlackModal = value;
    });
  },
  data() {
    return {
      headerFilterSpace: false,
      mobileBlackModal: false,
    };
  },
  methods: {
    closePopup() {
      this.$_bus.$emit('downloadMenuHide');
      this.$_bus.$emit('profileMenuHide');
      this.mobileBlackModal = false;
    },
    toogleHeaderSpace() {
      this.headerFilterSpace = !headerFilterSpace;
    },
    showMenuDownload() {
      if (this.curComponent === 'DownloadMenu') {
        this.showMenu = !this.showMenu;
        return;
      }
      this.curComponent = 'DownloadMenu';
      this.showMenu = true;
    },
  },
};
</script>

<style scoped lang="sass">
  @import "Layout"
</style>
