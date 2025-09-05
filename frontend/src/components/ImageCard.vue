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
  <div
    :class="[$style.imageCard, { [$style.selected]: isSelected }]"
    tabindex="0"
    @keydown.enter="toggleSelection"
    @keydown.space.prevent="toggleSelection"
  >
    <img :src="props.url" :alt="props.alt" :class="$style.image" />
    <div :class="[$style.overlay, { [$style.selected]: isSelected }]">
      <div :class="$style.overlayContent">
        <div
          :class="[$style.checkButton, { [$style.selected]: isSelected }]"
          @click.stop="toggleSelection"
        >
          <svg viewBox="0 0 24 24" fill="none" :class="$style.checkIcon">
            <path
              d="M5 13l4 4L19 7"
              stroke="currentColor"
              stroke-width="3"
              stroke-linecap="round"
              stroke-linejoin="round"
            />
          </svg>
        </div>
      </div>
    </div>
    <div v-if="isSelected" :class="$style.selectedBorder"></div>
  </div>
</template>

<style lang="scss" module>
.imageCard {
  width: 196px;
  height: 196px;
  overflow: hidden;
  position: relative;
  background-color: var(--Light-UI-Tertiary, #ced6db);
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.15s cubic-bezier(0.4, 0, 0.2, 1);
  user-select: none;

  &:hover {
    transform: translateY(-2px);
    box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12);
  }

  &:focus {
    outline: none;
    box-shadow: 0 0 0 4px rgba(0, 91, 172, 0.5);
  }

  &.selected {
    transform: scale(0.95);

    .image {
      transform: scale(0.85);
      border-radius: 6px;
    }
  }
}

.image {
  width: 100%;
  height: 100%;
  object-fit: cover;
  display: block;
  transition: transform 0.15s cubic-bezier(0.4, 0, 0.2, 1);
  border-radius: 0;
}

.overlay {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: linear-gradient(135deg, rgba(0, 91, 172, 0.1) 0%, rgba(0, 91, 172, 0.05) 100%);
  display: flex;
  justify-content: flex-start;
  align-items: flex-start;
  opacity: 0;
  transition: opacity 0.12s ease-out;
  padding: 12px;

  &.selected {
    opacity: 1;
    background: linear-gradient(135deg, rgba(0, 91, 172, 0.15) 0%, rgba(0, 91, 172, 0.08) 100%);
  }
}

.imageCard:hover .overlay {
  opacity: 0.7;
}

.overlayContent {
  text-align: center;
}

.checkButton {
  width: 28px;
  height: 28px;
  border-radius: 50%;
  background-color: rgba(255, 255, 255, 0.95);
  border: 2px solid rgba(0, 91, 172, 0.3);
  display: flex;
  justify-content: center;
  align-items: center;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.15);
  transition: all 0.12s cubic-bezier(0.4, 0, 0.2, 1);
  backdrop-filter: blur(8px);

  &:hover {
    transform: scale(1.05);
    border-color: rgba(0, 91, 172, 0.5);
    box-shadow: 0 4px 16px rgba(0, 91, 172, 0.2);
  }

  &.selected {
    background: linear-gradient(135deg, #005bac 0%, #004a94 100%);
    border-color: #005bac;
    transform: scale(1.1);
    box-shadow: 0 4px 16px rgba(0, 91, 172, 0.4);
  }
}

.checkIcon {
  width: 16px;
  height: 16px;
  color: rgba(0, 91, 172, 0.6);
  opacity: 0;
  transform: scale(0.5) rotate(-10deg);
  transition: all 0.15s cubic-bezier(0.68, -0.55, 0.265, 1.55);

  .checkButton.selected & {
    color: white;
    opacity: 1;
    transform: scale(1) rotate(0deg);
  }
}

.selectedBorder {
  position: absolute;
  top: -2px;
  left: -2px;
  right: -2px;
  bottom: -2px;
  border: 3px solid #005bac;
  border-radius: 10px;
  pointer-events: none;
  animation: borderPulse 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

@keyframes borderPulse {
  0% {
    transform: scale(0.95);
    opacity: 0;
  }
  50% {
    transform: scale(1.02);
    opacity: 0.8;
  }
  100% {
    transform: scale(1);
    opacity: 1;
  }
}
</style>
