import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { SIDEBAR_CONFIG } from './sidebarLayoutConfig'

export interface SidebarState {
  width: number
  isOpen: boolean
  isCompact: boolean
}

export const useSidebarStore = defineStore('sidebar', () => {
  // 状態
  const width = ref<number>(SIDEBAR_CONFIG.DEFAULT_WIDTH)
  const isOpen = ref(true)
  const isResizing = ref(false)

  // 計算
  const isCompact = computed(() => {
    if (!isOpen.value) return false
    return width.value < SIDEBAR_CONFIG.COMPACT_THRESHOLD
  })

  const effectiveWidth = computed(() => {
    if (!isOpen.value) return 0
    return Math.max(SIDEBAR_CONFIG.MIN_WIDTH, Math.min(SIDEBAR_CONFIG.MAX_WIDTH, width.value))
  })

  // アクション
  const setWidth = (newWidth: number) => {
    const clampedWidth = Math.max(
      SIDEBAR_CONFIG.MIN_WIDTH,
      Math.min(SIDEBAR_CONFIG.MAX_WIDTH, newWidth),
    )
    width.value = clampedWidth
    saveToStorage()
  }

  const toggleSidebar = () => {
    isOpen.value = !isOpen.value
    saveToStorage()
  }

  const openSidebar = () => {
    isOpen.value = true
    saveToStorage()
  }

  const closeSidebar = () => {
    isOpen.value = false
    saveToStorage()
  }

  const setResizing = (resizing: boolean) => {
    isResizing.value = resizing
  }

  // ローカルストレージ
  const saveToStorage = () => {
    if (typeof window !== 'undefined') {
      const state: SidebarState = {
        width: width.value,
        isOpen: isOpen.value,
        isCompact: isCompact.value,
      }
      localStorage.setItem(SIDEBAR_CONFIG.STORAGE_KEY, JSON.stringify(state))
    }
  }

  const loadFromStorage = () => {
    if (typeof window !== 'undefined') {
      try {
        const stored = localStorage.getItem(SIDEBAR_CONFIG.STORAGE_KEY)
        if (stored) {
          const state: SidebarState = JSON.parse(stored)
          width.value = Math.max(
            SIDEBAR_CONFIG.MIN_WIDTH,
            Math.min(SIDEBAR_CONFIG.MAX_WIDTH, state.width),
          )
          isOpen.value = state.isOpen
        }
      } catch (error) {
        console.warn('Failed to load sidebar state from localStorage:', error)
      }
    }
  }

  const initialize = () => {
    loadFromStorage()
  }

  return {
    // 状態
    width,
    isOpen,
    isResizing,
    isCompact,
    effectiveWidth,

    // アクション
    setWidth,
    toggleSidebar,
    openSidebar,
    closeSidebar,
    setResizing,
    saveToStorage,
    initialize,
  }
})
