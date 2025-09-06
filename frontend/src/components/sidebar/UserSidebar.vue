<script lang="ts" setup>
import SidebarItem from './SidebarItem.vue'
import UserProfile from './UserProfile.vue'
import { sidebarItems } from './sidebarConfig'
</script>

<template>
  <div :class="$style.sidebar">
    <div :class="$style.logoContainer">
      <img :class="$style.logo" src="/dummyLogo.png" alt="Application logo" />
    </div>
    <div :class="$style.navigationItems">
      <SidebarItem
        v-for="item in sidebarItems"
        :key="item.title"
        :icon="item.icon"
        :title="item.title"
        :path="item.path"
      />
    </div>
    <div :class="$style.userProfileContainer">
      <UserProfile />
    </div>
  </div>
</template>

<style lang="scss" module>
.sidebar {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
}

.logoContainer {
  display: flex;
  justify-content: center;
  width: 100%;
  margin-bottom: 16px;
}

.logo {
  width: 98px;
  height: 98px;
  flex-shrink: 0;
}

.navigationItems {
  display: flex;
  flex-direction: column;
  gap: 16px;
  flex: 1;
}

.userProfileContainer {
  margin-top: auto;
  padding-top: 16px;
}
</style>
