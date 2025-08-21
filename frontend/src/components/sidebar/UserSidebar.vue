<script lang="ts" setup>
import SidebarItem from './SidebarItem.vue'
import { sidebarItems } from './sidebarConfig'
</script>

<template>
  <div :class="$style.sidebar">
    <div :class="$style.logoContainer">
      <img :class="$style.logo" src="/dummyLogo.png" alt="Application logo" />
    </div>
    <SidebarItem
      v-for="item in sidebarItems"
      :key="item.title"
      :icon="item.icon"
      :title="item.title"
      :path="item.path"
    />
  </div>
</template>

<style lang="scss" module>
.sidebar {
  width: 100%;
  display: flex;
  flex-direction: column;
  gap: 16px;
}
.logo {
  width: 98px;
  height: 98px;
  flex-shrink: 0;
}
.logoContainer {
  display: flex;
  justify-content: center;
  width: 100%;
}
</style>
