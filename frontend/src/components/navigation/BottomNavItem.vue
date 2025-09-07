<template>
  <router-link :to="path" :class="[$style.navItem, { [$style['is-active']]: isActive }]">
    <div :class="$style.iconContainer">
      <img v-if="icon" :src="icon" :alt="title + ' icon'" :class="$style.icon" />
    </div>
    <span :class="$style.title">{{ title }}</span>
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
.navItem {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 6px 8px;
  text-decoration: none;
  color: #666;
  transition: all 0.2s ease;
  min-width: 50px;
  flex: 1;
  max-width: 90px;
  border-radius: 8px;
  position: relative;
}

.iconContainer {
  width: 28px;
  height: 28px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 2px;
  transition: transform 0.2s ease;
}

.icon {
  width: 24px;
  height: 24px;
  transition: filter 0.2s ease;
}

.title {
  font-size: 10px;
  font-weight: 400;
  text-align: center;
  line-height: 1.1;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 100%;
  transition: color 0.2s ease;
}

.navItem:hover {
  color: #005bac;
  background-color: rgba(0, 91, 172, 0.05);
}

.navItem:hover .iconContainer {
  transform: translateY(-1px);
}

.navItem:active {
  transform: scale(0.95);
  background-color: rgba(0, 91, 172, 0.1);
}

.navItem.is-active {
  color: #005bac;
  background-color: rgba(0, 91, 172, 0.1);
}

.navItem.is-active .title {
  font-weight: 600;
}

.navItem.is-active .iconContainer {
  transform: scale(1.1);
}

/* より小さな画面向けの調整 */
@media (max-width: 375px) {
  .navItem {
    padding: 4px 6px;
    min-width: 45px;
  }

  .iconContainer {
    width: 24px;
    height: 24px;
  }

  .icon {
    width: 20px;
    height: 20px;
  }

  .title {
    font-size: 9px;
  }
}

/* タッチデバイス用の調整 */
@media (pointer: coarse) {
  .navItem {
    min-height: 44px; /* タッチターゲットの最小サイズ */
  }

  .navItem:hover {
    background-color: transparent; /* タッチデバイスではホバー効果を無効化 */
  }
}
</style>
