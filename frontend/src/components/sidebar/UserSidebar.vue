<script lang="ts" setup>
import { onMounted, onBeforeUnmount, ref } from 'vue'
import SidebarItem from './SidebarItem.vue'
import SidebarToggle from './SidebarToggle.vue'
import { sidebarItems } from './sidebarConfig'
import { useSidebarStore } from '../../stores/sidebarStore'

const sidebarStore = useSidebarStore()
const resizeHandle = ref<HTMLElement | null>(null)
const isDragging = ref(false)

// ドラッグ開始
const startResize = (e: MouseEvent) => {
  e.preventDefault()
  isDragging.value = true
  sidebarStore.setResizing(true)
  document.body.style.cursor = 'col-resize'
  document.body.style.userSelect = 'none'

  document.addEventListener('mousemove', onMouseMove)
  document.addEventListener('mouseup', stopResize)
}

// ドラッグ中
const onMouseMove = (e: MouseEvent) => {
  if (!isDragging.value) return
  const newWidth = e.clientX
  sidebarStore.setWidth(newWidth)
}

// ドラッグ終了
const stopResize = () => {
  if (!isDragging.value) return
  isDragging.value = false
  sidebarStore.setResizing(false)
  document.body.style.cursor = ''
  document.body.style.userSelect = ''

  document.removeEventListener('mousemove', onMouseMove)
  document.removeEventListener('mouseup', stopResize)
}

// クリーンアップ
onBeforeUnmount(() => {
  document.removeEventListener('mousemove', onMouseMove)
  document.removeEventListener('mouseup', stopResize)
})

onMounted(() => {
  // ストアの初期化
  sidebarStore.initialize()
})
</script>

<template>
  <div :class="[$style.sidebar, { [$style.compact]: sidebarStore.isCompact }]">
    <SidebarToggle />
    <div :class="$style.logoContainer">
      <img
        v-if="!sidebarStore.isCompact"
        :class="$style.logo"
        src="/dummyLogo.png"
        alt="Application logo"
      />
      <img v-else :class="$style.logoCompact" src="/dummyLogo.png" alt="Application logo" />
    </div>
    <nav :class="$style.navigation">
      <SidebarItem
        v-for="item in sidebarItems"
        :key="item.title"
        :icon="item.icon"
        :title="item.title"
        :path="item.path"
      />
    </nav>
    <div
      v-if="sidebarStore.isOpen"
      ref="resizeHandle"
      :class="$style.resizeHandle"
      @mousedown="startResize"
    />
  </div>
</template>

<style lang="scss" module>
.sidebar {
  position: relative;
  width: 100%;
  height: 100vh;
  display: flex;
  flex-direction: column;
  gap: 16px;
  transition: all 0.25s ease;
}

.logoContainer {
  display: flex;
  justify-content: center;
  width: 100%;
}

.logo {
  width: 98px;
  height: 98px;
  flex-shrink: 0;
  transition: all 0.25s ease;
}

.logoCompact {
  width: 48px;
  height: 48px;
  flex-shrink: 0;
  transition: all 0.25s ease;
}

.navigation {
  display: flex;
  flex-direction: column;
  gap: 8px;
  padding-top: 30px;
  padding-left: 25px;
}

.sidebar.compact {
  align-items: center;
}

.sidebar.compact .logoContainer {
  margin-top: 56px;
  padding-left: 0;
}

.sidebar.compact .navigation {
  padding-left: 0;
}

.resizeHandle {
  position: absolute;
  /* Ensure the handle spans the full sidebar height and is positioned at the far right */
  top: 0;
  right: 0;
  bottom: 0;
  width: 8px; /* slightly wider for easier grabbing */
  height: auto;
  cursor: col-resize;
  background: transparent;
  transition: background-color 0.2s ease;
  z-index: 10;

  &:hover {
    background-color: rgba(0, 91, 172, 0.3);
  }

  &:active {
    background-color: rgba(0, 91, 172, 0.5);
  }
}

@media (prefers-reduced-motion: reduce) {
  .sidebar,
  .logo,
  .logoCompact {
    transition: none;
  }
}
</style>
