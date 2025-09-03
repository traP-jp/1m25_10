<script setup lang="ts">
import { ref } from 'vue'

interface Props {
  url: string
  alt?: string
}
const props = defineProps<Props>()
const isSelected = ref(false)

const toggleSelection = () => {
  isSelected.value = !isSelected.value
}
</script>

<template>
  <div :class="[$style.imageCard, { [$style.selected]: isSelected }]">
    <img :src="props.url" :alt="props.alt" :class="$style.image" />
    <div :class="[$style.overlay, { [$style.selected]: isSelected }]">
      <div :class="$style.overlayContent">
        <div
          :class="[$style.checkButton, { [$style.selected]: isSelected }]"
          @click="toggleSelection"
        >
          <svg viewBox="0 0 24 24" fill="none"></svg>
        </div>
      </div>
    </div>
  </div>
</template>

<style lang="scss" module>
.imageCard {
  width: 196px;
  height: 196px;
  overflow: hidden;
  position: relative;
  background-color: var(--Light-UI-Tertiary, #ced6db);
  transition: all 0.3s ease;

  &.selected {
    .image {
      transform: scale(0.8);
      border-radius: 8px;
    }
  }
}

.image {
  width: 100%;
  height: 100%;
  object-fit: cover;
  display: block;
  transition: transform 0.3s ease;
}

.overlay {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(206, 214, 219, 0.3);
  display: flex;
  justify-content: flex-start;
  align-items: flex-start;
  opacity: 0;
  transition: opacity 0.3s ease;
  padding: 8px;
}

.imageCard:hover .overlay,
.overlay.selected {
  opacity: 1;
}

.overlayContent {
  text-align: center;
}

.checkButton {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background-color: transparent;
  border: 2px solid #ced6db;
  display: flex;
  justify-content: center;
  align-items: center;
  cursor: pointer;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  transition: all 0.2s ease;

  svg {
    width: 24px;
    height: 24px;
    color: #005bac;
  }

  &:hover {
    transform: scale(1.1);
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  }

  &.selected {
    background-color: #005bac;
    svg {
      color: white;
    }
  }
}
</style>
