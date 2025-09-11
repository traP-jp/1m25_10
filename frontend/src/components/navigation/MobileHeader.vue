<template>
  <header :class="$style.header">
    <div :class="$style.logoContainer">
      <router-link to="/">
        <img :class="$style.logo" src="/dummyLogo.png" alt="Application logo" />
      </router-link>
    </div>
    <div :class="$style.userProfileContainer">
      <div :class="$style.userProfileMobile">
        <UserProfile />
      </div>
    </div>
  </header>
</template>

<script lang="ts" setup>
import UserProfile from '../sidebar/UserProfile.vue'
</script>

<style lang="scss" module>
.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background-color: #ffffff;
  border-bottom: 1px solid #e2e5e9;
  padding: 12px 16px;
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  z-index: 999;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  height: 80px; /* ヘッダーの高さを固定 */
}

.logoContainer {
  display: flex;
  align-items: center;
}

.logo {
  width: 40px;
  height: 40px;
}

/*
  Remove tap highlight color.
  Remove focus outline for mouse users, but keep it for keyboard users for accessibility.
*/
.logoContainer :global(a) {
  -webkit-tap-highlight-color: transparent;
  outline: none;
  box-shadow: none;
}

.logoContainer :global(a:focus-visible) {
  outline: revert; /* Restore outline for keyboard navigation */
  box-shadow: revert;
}

.logoContainer :global(a) > img {
  display: block; /* avoid extra inline spacing */
}

.userProfileContainer {
  display: flex;
  align-items: center;
}

.userProfileMobile {
  /* モバイル用UserProfileのスタイル調整 */
  :global(.userProfile) {
    max-width: 200px;
  }
}
</style>
