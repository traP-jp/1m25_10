<template>
  <nav :class="$style.bottomNav">
    <BottomNavItem
      v-for="item in navigationItems"
      :key="item.title"
      :icon="item.icon"
      :title="item.title"
      :path="item.path"
    />
  </nav>
</template>

<script lang="ts" setup>
import BottomNavItem from './BottomNavItem.vue'
import { sidebarItems } from '../sidebar/sidebarConfig'

// モバイルでは "Test" は表示しない
const navigationItems = sidebarItems.filter((item) => item.title !== 'Test')
</script>

<style lang="scss" module>
.bottomNav {
  display: flex;
  justify-content: space-around;
  align-items: center;
  background-color: #ffffff;
  border-top: 1px solid #e2e5e9;
  padding: 6px 0 calc(6px + env(safe-area-inset-bottom));
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  z-index: 1000;
  box-shadow: 0 -2px 8px rgba(0, 0, 0, 0.1);
  height: 60px; /* 固定の高さ */
  box-sizing: border-box;
}

/* iPhone X以降のセーフエリアを考慮 */
@supports (padding: max(0px)) {
  .bottomNav {
    padding-bottom: max(6px, env(safe-area-inset-bottom));
  }
}
</style>
