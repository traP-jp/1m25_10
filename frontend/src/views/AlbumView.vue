<template>
  <div :class="$style.albumView">
    <!-- アルバムが読み込み中の場合 -->
    <div v-if="albumStore.loading" :class="$style.loading">
      <p>アルバムを読み込み中...</p>
    </div>

    <!-- エラーが発生した場合 -->
    <div v-else-if="albumStore.error" :class="$style.error">
      <p>エラーが発生しました: {{ albumStore.error }}</p>
      <button @click="retryLoad" :class="$style.retryButton">再試行</button>
    </div>

    <!-- アルバムが見つからない場合 -->
    <div v-else-if="!album" :class="$style.notFound">
      <p>アルバムが見つかりません</p>
      <router-link to="/albums" :class="$style.backButton">アルバム一覧に戻る</router-link>
    </div>

    <!-- アルバム詳細表示 -->
    <div v-else>
      <!-- ナビゲーション（パンくずリスト風） -->
      <nav :class="$style.breadcrumb">
        <router-link to="/albums" :class="$style.breadcrumbLink">Albums</router-link>
        <span :class="$style.breadcrumbSeparator">/</span>
        <span :class="$style.breadcrumbCurrent">{{ album.title }}</span>
      </nav>

      <header :class="$style.header">
        <div :class="$style.titleSection">
          <h1 :class="$style.title">{{ album.title }}</h1>
          <p v-if="album.description" :class="$style.description">{{ album.description }}</p>
          <div :class="$style.albumMeta">
            <span :class="$style.imageCount">{{ album.images.length }}枚の画像</span>
            <span :class="$style.creator">作成者: {{ album.creator }}</span>
            <span :class="$style.createdAt">{{ formatDate(album.created_at) }}</span>
          </div>
        </div>
      </header>

      <!-- 画像一覧 -->
      <ImageList
        :images="albumImages"
        :loading="false"
        :error="null"
        :selected-image-ids="new Set()"
        :has-search-query="false"
        :has-more="false"
        :loading-more="false"
        :get-image-url="getImageUrl"
        @retry="() => {}"
        @toggle-selection="() => {}"
        @load-more="() => {}"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { useAlbumStore } from '@/stores/albumStore'
import { useImageStore } from '@/stores/imageStore'
import ImageList from '@/components/ImageList.vue'
import type { Image } from '@/types'

const route = useRoute()
const albumStore = useAlbumStore()
const imageStore = useImageStore()

const albumId = computed(() => route.params.id as string)

// 現在のアルバム
const album = computed(() => albumStore.currentAlbum)

// アルバムの画像をImage型の配列に変換
const albumImages = computed((): Image[] => {
  if (!album.value) return []

  return album.value.images.map((imageId) => ({
    id: imageId,
    creator: album.value!.creator,
    post: {
      id: `post-${imageId}`,
      content: `Image ${imageId} in album ${album.value!.title}`,
    },
  }))
})

// 画像URLを取得する関数
const getImageUrl = (image: Image): string => {
  return imageStore.getImageUrl(image)
}

// 日付をフォーマットする関数
const formatDate = (dateString: string): string => {
  const date = new Date(dateString)
  return date.toLocaleDateString('ja-JP', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  })
}

// アルバムを読み込み
const loadAlbum = async () => {
  if (albumId.value) {
    // 前のアルバムをクリア
    albumStore.clearCurrentAlbum()
    await albumStore.fetchAlbum(albumId.value)
  }
}

// 再読み込み
const retryLoad = () => {
  albumStore.clearError()
  loadAlbum()
}

// ルートのIDが変更された時にアルバムを再読み込み
watch(
  albumId,
  () => {
    loadAlbum()
  },
  { immediate: true },
)

onMounted(() => {
  loadAlbum()
})
</script>

<style lang="scss" module>
.albumView {
  padding: 24px;
  max-width: 1200px;
  margin: 0 auto;
  min-height: 100vh;
}

.loading,
.error,
.notFound {
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  min-height: 50vh;
  text-align: center;

  p {
    font-size: 1.1rem;
    color: #666;
    margin-bottom: 1rem;
  }
}

.retryButton,
.backButton {
  padding: 10px 20px;
  background-color: #007bff;
  color: white;
  border: none;
  border-radius: 4px;
  text-decoration: none;
  cursor: pointer;
  transition: background-color 0.2s;
  display: inline-block;

  &:hover {
    background-color: #0056b3;
  }
}

.header {
  margin-bottom: 32px;
  padding-bottom: 24px;
  border-bottom: 1px solid #e9ecef;
}

.titleSection {
  margin-bottom: 0;
}

.title {
  display: none; // パンくずリストがタイトルの役割を果たすため非表示
}

.description {
  font-size: 1rem;
  color: #666;
  margin: 0 0 12px 0;
  line-height: 1.5;
}

.albumMeta {
  display: flex;
  flex-wrap: wrap;
  gap: 16px;
  font-size: 0.9rem;
  color: #777;
}

.imageCount {
  font-weight: 600;
  color: #007bff;
}

.creator,
.createdAt {
  color: #666;
}

// ナビゲーション（パンくずリスト）のスタイル
.breadcrumb {
  margin-bottom: 20px;
  padding: 0;
  font-size: 2rem;
  font-weight: 700;
  color: #1a1a1a;
}

.breadcrumbLink {
  color: #666;
  text-decoration: none;
  transition: color 0.2s ease;
  font-weight: 700;
  font-size: 1.5rem; // メインタイトルより少し小さく

  &:hover {
    color: #333;
    text-decoration: underline;
  }
}

.breadcrumbSeparator {
  margin: 0 0.5rem;
  color: #666;
  font-weight: 700;
  font-size: 1.5rem; // リンクと同じサイズ
}

.breadcrumbCurrent {
  color: #1a1a1a;
  font-weight: 700;
  // 長いタイトルの場合は省略
  max-width: 300px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  display: inline-block;
  vertical-align: bottom;
}

// レスポンシブデザイン
@media (max-width: 768px) {
  .albumView {
    padding: 16px;
  }

  .header {
    margin-bottom: 24px;
  }

  .breadcrumb {
    font-size: 1.5rem;
  }

  .breadcrumbLink,
  .breadcrumbSeparator {
    font-size: 1.2rem; // タブレットサイズでも少し小さく
  }

  .breadcrumbCurrent {
    max-width: 150px;
  }

  .albumMeta {
    flex-direction: column;
    gap: 8px;
  }
}

@media (max-width: 480px) {
  .albumView {
    padding: 12px;
  }

  .breadcrumb {
    font-size: 1.3rem;
  }

  .breadcrumbLink,
  .breadcrumbSeparator {
    font-size: 1rem; // モバイルサイズでも少し小さく
  }

  .breadcrumbCurrent {
    max-width: 120px;
  }
}
</style>
