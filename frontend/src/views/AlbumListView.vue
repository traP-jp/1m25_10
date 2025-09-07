<template>
  <div :class="$style.albumListView">
    <header :class="$style.header">
      <div :class="$style.titleSection">
        <h1 :class="$style.title">Albums</h1>
        <p :class="$style.subtitle">すべてのアルバム ({{ albumStore.albumCount }}件)</p>
      </div>

      <div :class="$style.controls">
        <button :disabled="albumStore.loading" :class="$style.refreshButton" @click="reload">
          {{ albumStore.loading ? '読み込み中…' : '再読み込み' }}
        </button>
      </div>
    </header>

    <div v-if="albumStore.error" :class="$style.errorBox">
      <span>エラー: {{ albumStore.error }}</span>
      <button :class="$style.retryButton" @click="reload">再試行</button>
    </div>

    <div v-if="albumStore.loading && albums.length === 0" :class="$style.loading">読み込み中…</div>

    <div v-else-if="albums.length === 0" :class="$style.empty">まだアルバムがありません。</div>

    <div v-else :class="$style.grid">
      <AlbumCard
        v-for="album in albums"
        :key="album.id"
        :id="album.id"
        :title="album.title"
        :image-urls="album.imageUrls"
        :created-at="album.createdAt"
        :creator="album.creator"
        @click="goToAlbum(album.id)"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import AlbumCard from '@/components/AlbumCard.vue'
import { useAlbumStore } from '@/stores/albumStore'
import { generateImageUrl } from '@/config/env'

const router = useRouter()
const albumStore = useAlbumStore()

type AlbumCardVM = {
  id: string
  title: string
  imageUrls: string[]
  createdAt: string
  creator: { id: string; name: string; avatarUrl?: string }
}

const getCreatorData = (creatorId: string): AlbumCardVM['creator'] => {
  // 本実装ではユーザー情報APIが無いため、IDのみから暫定で表示
  return { id: creatorId, name: creatorId, avatarUrl: undefined }
}

const albums = computed<AlbumCardVM[]>(() => {
  return albumStore.albums.map((a) => {
    const detail = albumStore.albumDetails[a.id]
    const imageUrls = detail?.images ? detail.images.map((id) => generateImageUrl(id)) : []
    const createdAt = detail?.created_at ?? new Date().toISOString()

    return {
      id: a.id,
      title: a.title,
      imageUrls,
      createdAt,
      creator: getCreatorData(a.creator),
    }
  })
})

const reload = async () => {
  await albumStore.fetchAlbums()
  // 一覧に出ている分の詳細を並行で取得（画像/作成日時のため）
  await Promise.allSettled(albumStore.albums.map((a) => albumStore.fetchAlbum(a.id)))
}

const goToAlbum = (albumId: string) => {
  router.push({ name: 'album', params: { id: albumId } })
}

onMounted(async () => {
  if (albumStore.albums.length === 0) {
    await albumStore.fetchAlbums()
  }
  // 初回に一覧分の詳細を取得
  if (albumStore.albums.length > 0) {
    await Promise.allSettled(albumStore.albums.map((a) => albumStore.fetchAlbum(a.id)))
  }
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

.controls {
  display: flex;
  gap: 8px;
  align-items: center;
}

.refreshButton {
  padding: 8px 12px;
  background-color: #007bff;
  color: #fff;
  border: none;
  border-radius: 6px;
  font-size: 14px;
  cursor: pointer;

  &:disabled {
    background-color: #6c757d;
    cursor: not-allowed;
  }

  &:hover:not(:disabled) {
    background-color: #0056b3;
  }
}

.grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 20px;
}

.loading,
.empty {
  color: #666;
}

.errorBox {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  border: 1px solid #f5c2c7;
  background: #f8d7da;
  color: #842029;
  border-radius: 6px;
  margin-bottom: 16px;
}

.retryButton {
  padding: 6px 10px;
  background: #dc3545;
  color: #fff;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

@media (max-width: 768px) {
  .albumListView {
    padding: 16px;
  }
}
</style>
