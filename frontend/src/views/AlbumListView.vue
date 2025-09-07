<template>
  <div :class="$style.albumListView">
    <header :class="$style.header">
      <div :class="$style.titleSection">
        <h1 :class="$style.title">アルバム一覧</h1>
        <p :class="$style.subtitle">すべてのアルバム ({{ albumStore.albumCount }}個)</p>
      </div>
    </header>

    <!-- ローディング状態 -->
    <div v-if="albumStore.loading" :class="$style.loading">
      <p>アルバムを読み込み中...</p>
    </div>

    <!-- エラー状態 -->
    <div v-else-if="albumStore.error" :class="$style.error">
      <p>エラーが発生しました: {{ albumStore.error }}</p>
      <button @click="retryLoad" :class="$style.retryButton">再試行</button>
    </div>

    <!-- アルバム一覧 -->
    <div v-else-if="albumStore.albums.length > 0" :class="$style.albumGrid">
      <router-link
        v-for="album in albumStore.albums"
        :key="album.id"
        :to="`/albums/${album.id}`"
        :class="$style.albumLink"
      >
        <AlbumCard
          :id="album.id"
          :title="album.title"
          :image-urls="[]"
          :created-at="new Date().toISOString()"
          :creator="{ id: album.creator, name: album.creator }"
        />
      </router-link>
    </div>

    <!-- アルバムが存在しない場合 -->
    <div v-else :class="$style.empty">
      <p>まだアルバムがありません</p>
      <router-link to="/" :class="$style.backButton">ホームに戻る</router-link>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useAlbumStore } from '@/stores/albumStore'
import AlbumCard from '@/components/AlbumCard.vue'

const albumStore = useAlbumStore()

// アルバム一覧を読み込み
const loadAlbums = async () => {
  await albumStore.fetchAlbums()
}

// 再読み込み
const retryLoad = () => {
  albumStore.clearError()
  loadAlbums()
}

onMounted(() => {
  loadAlbums()
})
</script>

<style lang="scss" module>
.albumListView {
  padding: 24px;
  max-width: 1200px;
  margin: 0 auto;
  min-height: 100vh;
}

.header {
  margin-bottom: 32px;
}

.titleSection {
  margin-bottom: 20px;
}

.title {
  font-size: 2rem;
  font-weight: 700;
  color: #1a1a1a;
  margin: 0 0 8px 0;
}

.subtitle {
  font-size: 1rem;
  color: #666;
  margin: 0;
}

.loading,
.error,
.empty {
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

.albumGrid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 20px;
  justify-items: center;
}

.albumLink {
  text-decoration: none;
  color: inherit;
  width: 100%;
  max-width: 350px;

  &:hover {
    transform: translateY(-2px);
    transition: transform 0.2s ease-in-out;
  }
}

// レスポンシブデザイン
@media (max-width: 768px) {
  .albumListView {
    padding: 16px;
  }

  .title {
    font-size: 1.5rem;
  }

  .albumGrid {
    grid-template-columns: 1fr;
    gap: 16px;
  }
}

@media (max-width: 480px) {
  .albumListView {
    padding: 12px;
  }

  .title {
    font-size: 1.3rem;
  }
}
</style>
