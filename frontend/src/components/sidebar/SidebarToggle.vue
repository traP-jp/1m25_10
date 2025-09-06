<script lang="ts" setup>
import { useSidebarStore } from '../../stores/sidebarStore'

const sidebarStore = useSidebarStore()
</script>

<template>
  <div>
    <button
      v-if="sidebarStore.isOpen"
      :class="[$style.toggleButton, $style.closeButton]"
      @click="sidebarStore.closeSidebar()"
      aria-label="サイドバーを閉じる"
    >
      <svg :class="$style.icon" viewBox="0 0 24 24" fill="none">
        <path
          d="M6 6L18 18M6 18L18 6"
          stroke="currentColor"
          stroke-width="2"
          stroke-linecap="round"
        />
      </svg>
    </button>

    <button
      v-else
      :class="[$style.toggleButton, $style.openButton]"
      @click="sidebarStore.openSidebar()"
      aria-label="サイドバーを開く"
    >
      <svg :class="$style.icon" viewBox="0 0 24 24" fill="none">
        <path
          d="M3 6h18M3 12h18M3 18h18"
          stroke="currentColor"
          stroke-width="2"
          stroke-linecap="round"
        />
      </svg>
    </button>
  </div>
</template>

<style lang="scss" module>
.toggleButton {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 40px;
  height: 40px;
  border: none;
  border-radius: 8px;
  background-color: rgba(255, 255, 255, 0.8);
  color: #333;
  cursor: pointer;
  transition: all 0.2s ease;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  z-index: 1001;

  &:hover {
    background-color: rgba(255, 255, 255, 0.95);
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
    color: #005bac;
  }

  &:active {
    transform: scale(0.95);
  }
}

.closeButton {
  position: absolute;
  top: 16px;
  right: 19px;
  z-index: 1002;
}

.openButton {
  position: fixed;
  top: 16px;
  left: 20px;
  z-index: 1002;
}

.icon {
  width: 20px;
  height: 20px;
}

@media (prefers-reduced-motion: reduce) {
  .toggleButton {
    transition: none;
  }
}
</style>
