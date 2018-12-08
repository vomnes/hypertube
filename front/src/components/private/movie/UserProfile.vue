<template lang="html">
  <div @click.self="toggleUserProfil" class="user-profil-bg">
    <div class="user-profil-modal">
      <img :src="userProfil.picture" class="user-profil-modal-img">
      <div class="user-profil-modal-text">
        <p class="user-profil-modal-fullname">{{ userProfil.fullname }}</p>
        <img class="user-profil-modal-lgn" src="../../../static/icons/france.svg" v-if="userProfil.local === 'fr'" />
        <img class="user-profil-modal-lgn" src="../../../static/icons/italy.svg" v-else-if="userProfil.local === 'it'" />
        <img class="user-profil-modal-lgn" src="../../../static/icons/united-kingdom.svg" v-else-if="userProfil.local === 'en'" />
      </div>
    </div>
  </div>
</template>

<script>

export default {
  props: ['userData', 'toggleUserProfil', 'userProfil'],
};
</script>

<style scoped lang="sass">
  @import "UserProfile"
</style>
