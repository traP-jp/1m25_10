<template>
  <router-link :to="path" :class="[$style.sidebarItem, { [$style['is-active']]: isActive }]">
    <!-- fixed icon slot: preserves space even when icon is absent -->
    <span :class="$style.iconSlot">
      <img v-if="icon" :src="icon" :alt="title + ' icon'" :class="$style.icon" />
    </span>

    <span :class="$style.title">
      {{ title }}
    </span>
  </router-link>
</template>

<script lang="ts" setup>
import { useRoute } from 'vue-router'
import { computed } from 'vue'

const props = defineProps<{
  icon?: string
  title: string
  path: string
}>()

const route = useRoute()
const isActive = computed(() => route.path === props.path)
</script>

<style lang="scss" module>
.sidebarItem {
  display: grid;
  grid-template-columns: 40px 1fr; /* fixed slot for icon + flexible content */
  align-items: center;
  padding: 8px 16px;
  gap: 12px;
  color: #333;
  border-radius: 32px;
  transition: background-color 0.3s;
  text-decoration: none;
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
}

.sidebarItem.is-active {
  background-color: #d8e0eb;
}

.sidebarItem.is-active .title {
  color: #005bac;
  font-weight: bold;
}
</style>
