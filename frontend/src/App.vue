<script setup lang="ts">
import { RouterView } from 'vue-router'
import Sidebar from './components/sideBar/UserSidebar.vue';
</script>

<template>
  <div :class="$style.app">
    <aside :class="$style.sidebar">
      <Sidebar />
    </aside>
    <main :class="$style.main">
      <RouterView />
    </main>
  </div>
</template>

<style lang="scss" module>
.app {
  display: flex;
  flex-direction: row;
  align-items: start;
  width: 100%;
  height: 100dvh;
  gap: 16px;
}
.sidebar {
  display: flex;
  width: 365px;
  height: 100%;
  padding: 40px;
  flex-direction: column;
  align-items: flex-start;
  gap: 16px;
  flex-shrink: 0;
  position: fixed;
  top: 0;
  bottom: 0;
  background-color: #F0F2F5;
  /*色はあとでやる*/
}
.main {
  width: calc(100% - 365px);
  margin-left: auto;
}

</style>
