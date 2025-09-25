<template>
  <div :class="$style.imageListContainer">
    <div v-if="loading" :class="$style.loading">
      <p>画像を読み込み中...</p>
    </div>

    <div v-else-if="error" :class="$style.error">
      <p>エラーが発生しました: {{ error }}</p>
      <button @click="$emit('retry')" :class="$style.retryButton">再試行</button>
    </div>

    <div v-else-if="images.length > 0" :class="$style.imageGrid">
      <ImageCard
        v-for="image in images"
        :key="image.id"
        :url="getImageUrl(image)"
        alt=""
        :is-selected="isImageSelected(image.id)"
        @toggle-selection="() => $emit('toggleSelection', image.id)"
        @open-detail="$emit('openDetail', image.id)"
      />

      <div
        v-if="hasMore"
        :class="[$style.loadMoreCard, { [$style.loading]: loadingMore }]"
        @click="$emit('loadMore')"
      >
        <div :class="$style.loadMoreContent">
          <span v-if="!loadingMore">Load more</span>
          <span v-else>Loading...</span>
        </div>
      </div>
    </div>

    <div v-else :class="$style.empty">
      <p v-if="hasSearchQuery">検索条件に一致する画像がありません</p>
      <p v-else>まだ画像がありません</p>
    </div>
    <!-- Always-present sentinel for IntersectionObserver -->
    <div ref="sentinelRef" :class="$style.sentinel" aria-hidden="true"></div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount, watch, nextTick } from 'vue'
import ImageCard from '@/components/ImageCard.vue'
import type { Image } from '@/types'

// Infinite scroll margin/threshold in pixels
const INFINITE_SCROLL_MARGIN_PX = 400

interface Props {
  images: Image[]
  loading?: boolean
  error?: string | null
  selectedImageIds?: Set<string>
  hasSearchQuery?: boolean
  hasMore?: boolean
  loadingMore?: boolean
  getImageUrl: (image: Image) => string
}

interface Emits {
  retry: []
  toggleSelection: [imageId: string]
  loadMore: []
  openDetail: [imageId: string]
}

const props = withDefaults(defineProps<Props>(), {
  loading: false,
  error: null,
  selectedImageIds: () => new Set(),
  hasSearchQuery: false,
  hasMore: false,
  loadingMore: false,
})

const emit = defineEmits<Emits>()

const sentinelRef = ref<HTMLElement | null>(null)
let observer: IntersectionObserver | null = null

const RESIZE_DEBOUNCE_MS = 120

const resizeTimer = ref<number | undefined>(undefined)

const onResize = () => {
  if (resizeTimer.value !== undefined) {
    window.clearTimeout(resizeTimer.value)
  }
  resizeTimer.value = window.setTimeout(() => {
    maybeAutoLoad()
  }, RESIZE_DEBOUNCE_MS)
}

const setupObserver = () => {
  if (observer) {
    observer.disconnect()
    observer = null
  }

  if (!props.hasMore) return

  observer = new IntersectionObserver(
    (entries) => {
      const entry = entries[0]
      if (
        entry &&
        entry.isIntersecting &&
        props.hasMore &&
        !props.loadingMore &&
        !props.loading &&
        !props.error
      ) {
        emit('loadMore')
      }
    },
    {
      root: null,
      rootMargin: `${INFINITE_SCROLL_MARGIN_PX}px 0px`,
      threshold: 0,
    },
  )

  if (sentinelRef.value) {
    observer.observe(sentinelRef.value)
  }
}

const isInViewport = (el: HTMLElement, rootMarginPx = INFINITE_SCROLL_MARGIN_PX) => {
  const rect = el.getBoundingClientRect()
  const viewH = window.innerHeight || document.documentElement.clientHeight
  return rect.top <= viewH + rootMarginPx
}

const maybeAutoLoad = () => {
  if (!props.hasMore || props.loading || props.loadingMore || props.error) return
  const el = sentinelRef.value
  const doc = document.scrollingElement || document.documentElement
  const fitsViewport = doc
    ? doc.scrollHeight <= window.innerHeight + INFINITE_SCROLL_MARGIN_PX
    : false
  if ((el && isInViewport(el)) || fitsViewport) {
    emit('loadMore')
  }
}

onMounted(() => {
  nextTick(() => {
    setupObserver()
    maybeAutoLoad()
  })
  window.addEventListener('resize', onResize)
})

onBeforeUnmount(() => {
  if (observer) observer.disconnect()
  window.removeEventListener('resize', onResize)
  if (resizeTimer.value !== undefined) {
    window.clearTimeout(resizeTimer.value)
    resizeTimer.value = undefined
  }
})

watch(
  () => [props.hasMore, props.loadingMore, props.loading, props.error, props.images.length],
  () => {
    // 状態変化に応じて監視を張り直す
    setupObserver()
    nextTick(() => maybeAutoLoad())
  },
)

// 画像が選択されているかチェック
const isImageSelected = (imageId: string): boolean => {
  return props.selectedImageIds.has(imageId)
}
</script>

<style lang="scss" module>
.imageListContainer {
  width: 100%;
}

.imageGrid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(196px, 1fr));
  gap: 20px;
  justify-items: center;
}

.loading {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 200px;

  p {
    font-size: 1.1rem;
    color: #666;
  }
}

.error {
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  min-height: 200px;

  p {
    font-size: 1.1rem;
    color: #e74c3c;
    margin-bottom: 16px;
  }
}

.retryButton {
  padding: 10px 20px;
  background-color: #005bac;
  color: white;
  border: none;
  border-radius: 6px;
  font-size: 1rem;
  cursor: pointer;
  transition: background-color 0.2s;

  &:hover {
    background-color: #004a94;
  }

  &:active {
    transform: translateY(1px);
  }
}

.empty {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 300px;

  p {
    font-size: 1.2rem;
    color: #999;
  }
}

.loadMoreCard {
  width: 196px;
  height: 196px;
  border: 2px dashed #ccc;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.2s ease;
  background-color: #f9f9f9;

  &:hover {
    border-color: #007bff;
    background-color: #f0f8ff;
  }

  &.loading {
    cursor: not-allowed;
    opacity: 0.6;
    border-color: #999;
  }
}

.loadMoreContent {
  text-align: center;
  color: #666;
  font-size: 1rem;
  font-weight: 500;
}

.sentinel {
  width: 100%;
  height: 1px;
}

// レスポンシブデザイン
@media (max-width: 768px) {
  .imageGrid {
    grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
    gap: 16px;
  }

  .loadMoreCard {
    width: 150px;
    height: 150px;
  }
}

@media (max-width: 480px) {
  .imageGrid {
    grid-template-columns: repeat(auto-fill, minmax(120px, 1fr));
    gap: 12px;
  }

  .loadMoreCard {
    width: 120px;
    height: 120px;
    font-size: 0.9rem;
  }
}
</style>
