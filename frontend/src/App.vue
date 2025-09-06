<script setup lang="ts">
import { RouterView } from 'vue-router'
import UserSidebar from './components/sidebar/UserSidebar.vue'
import SidebarToggle from './components/sidebar/SidebarToggle.vue'
import { useSidebarStore } from './stores/sidebarStore'

const sidebarStore = useSidebarStore()
</script>

<template>
  <div :class="$style.app" :data-sidebar-open="sidebarStore.isOpen">
    <aside
      v-if="sidebarStore.isOpen"
      :class="$style.sidebar"
      :style="{ width: `${sidebarStore.effectiveWidth}px` }"
    >
      <UserSidebar />
    </aside>
    <main :class="$style.main" :style="{ marginLeft: `${sidebarStore.effectiveWidth}px` }">
      <SidebarToggle v-if="!sidebarStore.isOpen" />
      <RouterView />
    </main>
  </div>
</template>

<style lang="scss" module>
.app {
  box-sizing: border-box;
  display: flex;
  flex-direction: row;
  align-items: start;
  width: 100%;
  min-height: 100svh;
  transition: all 0.25s ease;
}
.sidebar {
  position: fixed;
  left: 0;
  top: 0;
  bottom: 0;
  z-index: 1000;
  height: 100vh;
  background-color: #f0f2f5;
  transition: all 0.25s ease;
  will-change: transform, width;
}
.main {
  width: 100%;
  min-height: 100vh;
  transition: margin-left 0.25s ease;
  position: relative;
}

@media (prefers-reduced-motion: reduce) {
  .app,
  .sidebar,
  .main {
    transition: none;
  }
}
</style>
