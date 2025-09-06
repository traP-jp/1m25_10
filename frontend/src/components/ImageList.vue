<template>
  <div :class="$style.imageListContainer">
    <!-- ローディング状態 -->
    <div v-if="loading" :class="$style.loading">
      <p>画像を読み込み中...</p>
    </div>

    <!-- エラー状態 -->
    <div v-else-if="error" :class="$style.error">
      <p>エラーが発生しました: {{ error }}</p>
      <button @click="$emit('retry')" :class="$style.retryButton">再試行</button>
    </div>

    <!-- 画像一覧 -->
    <div v-else-if="images.length > 0" :class="$style.imageGrid">
      <ImageCard
        v-for="image in images"
        :key="image.id"
        :url="getImageUrl(image)"
        :alt="image.post.content"
        :is-selected="isImageSelected(image.id)"
        @toggle-selection="() => $emit('toggleSelection', image.id)"
      />
    </div>

    <!-- 画像がない場合 -->
    <div v-else :class="$style.empty">
      <p v-if="hasSearchQuery">検索条件に一致する画像がありません</p>
      <p v-else>まだ画像がありません</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import ImageCard from '@/components/ImageCard.vue'
import type { Image } from '@/types'

interface Props {
  images: Image[]
  loading?: boolean
  error?: string | null
  selectedImageIds?: Set<string>
  hasSearchQuery?: boolean
  getImageUrl: (image: Image) => string
}

interface Emits {
  retry: []
  toggleSelection: [imageId: string]
}

const props = withDefaults(defineProps<Props>(), {
  loading: false,
  error: null,
  selectedImageIds: () => new Set(),
  hasSearchQuery: false,
})

defineEmits<Emits>()

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

// レスポンシブデザイン
@media (max-width: 768px) {
  .imageGrid {
    grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
    gap: 16px;
  }
}

@media (max-width: 480px) {
  .imageGrid {
    grid-template-columns: repeat(auto-fill, minmax(120px, 1fr));
    gap: 12px;
  }
}
</style>
