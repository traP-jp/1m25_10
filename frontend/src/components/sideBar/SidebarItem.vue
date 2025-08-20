<template>
  <router-link :to="path" :class="[$style.sidebarItem, { [$style['is-active']]: isActive }]">
    <img :src="icon" :alt="title + ' icon'" />
    <span :class="$style.title">
      {{ title }}
    </span>
  </router-link>
</template>

<script lang="ts" setup>
import { useRoute } from 'vue-router'
import { computed } from 'vue'

const props = defineProps<{
  icon: string
  title: string
  path: string
}>()

const route = useRoute()
const isActive = computed(() => route.path === props.path)
</script>

<style lang="scss" module>
.sidebarItem {
  display: flex;
  padding: 8px 16px;
  align-items: center;
  gap: 16px;
  color: #333;
  border-radius: 32px;
  transition: background-color 0.3s;
  text-decoration: none;
}

.sidebarItem:hover {
  background-color: #e2e5e9;
  /*色はあとでやる*/
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
