import { ref, onMounted, onUnmounted } from 'vue'

export function useWindowSize() {
  const width = ref(window.innerWidth)
  const height = ref(window.innerHeight)

  // モバイル: 1024px未満、デスクトップ: 1024px以上
  const isMobile = ref(width.value < 1024)
  const isTablet = ref(width.value >= 768 && width.value < 1024)
  const isDesktop = ref(width.value >= 1024)

  function updateSize() {
    width.value = window.innerWidth
    height.value = window.innerHeight
    isMobile.value = width.value < 1024
    isTablet.value = width.value >= 768 && width.value < 1024
    isDesktop.value = width.value >= 1024
  }

  onMounted(() => {
    window.addEventListener('resize', updateSize)
  })

  onUnmounted(() => {
    window.removeEventListener('resize', updateSize)
  })

  return {
    width,
    height,
    isMobile,
    isTablet,
    isDesktop,
  }
}
