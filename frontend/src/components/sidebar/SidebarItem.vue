<template>
  <div :class="$style.itemWrapper">
    <router-link
      :to="path"
      :class="[
        $style.sidebarItem,
        { [$style['is-active']]: isActive, [$style.compact]: isCompact },
      ]"
      :title="isCompact ? title : undefined"
    >
      <!-- fixed icon slot: preserves space even when icon is absent -->
      <span :class="$style.iconSlot">
        <img v-if="icon" :src="icon" :alt="title + ' icon'" :class="$style.icon" />
      </span>

      <span v-if="!isCompact" :class="$style.title">
        {{ title }}
      </span>
    </router-link>

    <div v-if="isCompact" :class="$style.tooltip">
      {{ title }}
    </div>
  </div>
</template>

<script lang="ts" setup>
import { useRoute } from 'vue-router'
import { computed } from 'vue'
import { useSidebarStore } from '../../stores/sidebarStore'

const props = defineProps<{
  icon?: string
  title: string
  path: string
}>()

const route = useRoute()
const sidebarStore = useSidebarStore()

const isActive = computed(() => route.path === props.path)
const isCompact = computed(() => sidebarStore.isCompact)
</script>

<style lang="scss" module>
.itemWrapper {
  position: relative;
}

.sidebarItem {
  display: grid;
  align-items: center;
  padding: 8px 16px;
  gap: 12px;
  color: #333;
  border-radius: 32px;
  transition: background-color 0.3s;
  text-decoration: none;
  width: 100%;
  box-sizing: border-box;
}

.sidebarItem:not(.compact) {
  grid-template-columns: 40px 1fr; /* fixed slot for icon + flexible content */
  max-width: calc(100% - 32px);
}

.sidebarItem.compact {
  grid-template-columns: 40px; /* icon only */
  justify-content: center;
  padding: 8px;
}

.iconSlot {
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.icon {
  width: 32px;
  height: 32px;
}

.sidebarItem:hover {
  background-color: #e2e5e9;
  /* TODO: Set color later */
}

.title {
  color: var(--Light-UI-Primary, #49535b);
  font-family: Inter;
  font-size: 32px;
  font-style: normal;
  font-weight: 400;
  line-height: normal;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.sidebarItem.is-active {
  background-color: #d8e0eb;
}

.sidebarItem.is-active .title {
  color: #005bac;
  font-weight: bold;
}

/* ツールチップ */
.tooltip {
  position: absolute;
  left: 100%;
  top: 50%;
  transform: translateY(-50%);
  margin-left: 8px;
  background-color: rgba(0, 0, 0, 0.8);
  color: white;
  padding: 8px 12px;
  border-radius: 6px;
  font-size: 14px;
  white-space: nowrap;
  pointer-events: none;
  opacity: 0;
  transition: opacity 0.2s ease;
  z-index: 1000;
}

.itemWrapper:hover .tooltip {
  opacity: 1;
}

/* ツールチップの矢印 */
.tooltip::before {
  content: '';
  position: absolute;
  top: 50%;
  left: -4px;
  transform: translateY(-50%);
  width: 0;
  height: 0;
  border-top: 4px solid transparent;
  border-bottom: 4px solid transparent;
  border-right: 4px solid rgba(0, 0, 0, 0.8);
}
</style>
