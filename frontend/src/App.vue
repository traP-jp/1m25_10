<script setup lang="ts">
import { ref } from 'vue'
import { RouterView } from 'vue-router'
import Sidebar from './components/sidebar/UserSidebar.vue'
import BottomNavigation from './components/navigation/BottomNavigation.vue'
import MobileHeader from './components/navigation/MobileHeader.vue'
import { useWindowSize } from './composables/useWindowSize'

// TODO: 将来的に Pinia の uiStore に移行予定
// ex. `const ui = useUiStore(); const isSidebarOpen = storeToRefs(ui).isSidebarOpen`
const isSidebarOpen = ref(true)
const { isMobile } = useWindowSize()
</script>

<template>
  <!-- デスクトップレイアウト -->
  <div v-if="!isMobile" :class="$style.app" :data-sidebar-open="isSidebarOpen">
    <aside :class="$style.sidebar">
      <Sidebar />
    </aside>
    <main :class="$style.main">
      <RouterView />
    </main>
  </div>

  <!-- モバイルレイアウト -->
  <div v-else :class="$style.mobileApp">
    <MobileHeader />
    <main :class="$style.mobileMain">
      <RouterView />
    </main>
    <BottomNavigation />
  </div>
</template>

<style lang="scss" module>
/* デスクトップレイアウト */
.app {
  box-sizing: border-box;
  display: flex;
  flex-direction: row;
  align-items: start;
  width: 100%;
  min-height: 100svh;
  gap: 16px;
  --sidebar-width: clamp(240px, 25vw, 365px);
  padding-inline-start: var(--sidebar-width);
}

.sidebar {
  box-sizing: border-box;
  display: flex;
  width: var(--sidebar-width);
  height: 100%;
  padding: 40px;
  flex-direction: column;
  align-items: flex-start;
  gap: 16px;
  flex-shrink: 0;
  position: fixed;
  left: 0;
  top: 0;
  bottom: 0;
  background-color: #f0f2f5;
  /* TODO: Set color later */
  transition: transform 0.25s ease;
  will-change: transform;
  z-index: 1000;
}

.main {
  width: 100%;
  background: #f8f9fa;
}

.app[data-sidebar-open='false'] .sidebar {
  transform: translateX(calc(-1 * var(--sidebar-width)));
}

.app[data-sidebar-open='false'] {
  padding-inline-start: 0;
}

/* モバイルレイアウト */
.mobileApp {
  display: flex;
  flex-direction: column;
  width: 100%;
  min-height: 100svh;
  background: #f8f9fa;
}

.mobileMain {
  flex: 1;
  width: 100%;
  padding-top: 100px; /* ヘッダーの高さ分 */
  padding-bottom: 76px; /* ボトムナビの高さ分 + セーフエリア */
  background: #f8f9fa;
  min-height: calc(100svh - 144px); /* フルスクリーンの高さからヘッダーとナビを除く */
}

@media (prefers-reduced-motion: reduce) {
  .sidebar {
    transition: none;
  }
}
</style>
