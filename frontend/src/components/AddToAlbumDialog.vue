<template>
  <div v-if="isVisible" class="dialog-overlay" @click="handleOverlayClick">
    <div class="dialog" @click.stop>
      <div class="dialog-header">
        <h3>アルバムに追加</h3>
        <button type="button" class="close-button" @click="$emit('close')">×</button>
      </div>

      <div class="dialog-content">
        <div class="selected-images-info">
          <p>選択された画像: {{ selectedImageCount }}枚</p>
        </div>

        <!-- アルバム一覧 -->
        <div class="album-list-section">
          <h4>追加先のアルバムを選択</h4>

          <!-- ローディング状態 -->
          <div v-if="isLoadingAlbums" class="loading-state">アルバム一覧を読み込み中...</div>

          <!-- エラー状態 -->
          <div v-else-if="albumError" class="error-state">
            <p>{{ albumError }}</p>
            <button @click="loadAlbums" class="retry-button">再試行</button>
          </div>

          <!-- アルバム一覧 -->
          <div v-else-if="albums.length > 0" class="album-list">
            <div
              v-for="album in albums"
              :key="album.id"
              class="album-item"
              :class="{ selected: selectedAlbumId === album.id }"
              @click="selectAlbum(album.id)"
            >
              <div class="album-info">
                <h5>{{ album.title }}</h5>
                <span class="album-meta">作成者: {{ album.creator }}</span>
              </div>
              <div class="album-selection">
                <input
                  type="radio"
                  :value="album.id"
                  v-model="selectedAlbumId"
                  :id="`album-${album.id}`"
                />
              </div>
            </div>
          </div>

          <!-- アルバムが存在しない場合 -->
          <div v-else class="empty-state">
            <p>アルバムがありません</p>
            <p class="hint">先にアルバムを作成してください</p>
          </div>
        </div>

        <!-- エラーメッセージ -->
        <div v-if="errorMessage" class="error-message">
          {{ errorMessage }}
        </div>

        <div class="dialog-actions">
          <button type="button" @click="$emit('close')" :disabled="isAdding" class="cancel-button">
            キャンセル
          </button>
          <button
            type="button"
            @click="handleAddToAlbum"
            :disabled="!selectedAlbumId || isAdding"
            class="add-button"
          >
            {{ isAdding ? '追加中...' : 'アルバムに追加' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import type { AlbumItem } from '@/types'
import { albumService } from '@/services'

interface Props {
  isVisible: boolean
  selectedImageCount: number
}

interface Emits {
  (e: 'close'): void
  (e: 'add-to-album', data: { albumId: string }): Promise<void>
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const selectedAlbumId = ref('')
const isAdding = ref(false)
const errorMessage = ref('')
const albums = ref<AlbumItem[]>([])
const isLoadingAlbums = ref(false)
const albumError = ref('')

const handleOverlayClick = () => {
  if (!isAdding.value) {
    emit('close')
  }
}

const selectAlbum = (albumId: string) => {
  selectedAlbumId.value = albumId
}

const handleAddToAlbum = async () => {
  if (!selectedAlbumId.value || isAdding.value) return

  isAdding.value = true
  errorMessage.value = ''

  try {
    await emit('add-to-album', {
      albumId: selectedAlbumId.value,
    })
    // 成功時は親コンポーネントがダイアログを閉じる
  } catch (error) {
    errorMessage.value = error instanceof Error ? error.message : 'アルバムへの追加に失敗しました'
  } finally {
    isAdding.value = false
  }
}

// アルバム一覧を取得
const loadAlbums = async () => {
  isLoadingAlbums.value = true
  albumError.value = ''

  try {
    albums.value = await albumService.getAlbums()
  } catch (error) {
    albumError.value = error instanceof Error ? error.message : 'アルバム一覧の取得に失敗しました'
  } finally {
    isLoadingAlbums.value = false
  }
}

// ダイアログが閉じられたときにフォームをリセット
const resetForm = () => {
  selectedAlbumId.value = ''
  isAdding.value = false
  errorMessage.value = ''
}

// props.isVisibleが変更されたときの処理
watch(
  () => props.isVisible,
  (newValue) => {
    if (newValue) {
      resetForm()
      loadAlbums()
    }
  },
)
</script>

<style scoped lang="scss">
.dialog-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.dialog {
  background: white;
  border-radius: 8px;
  width: 90%;
  max-width: 600px;
  max-height: 90vh;
  overflow: hidden;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.15);
}

.dialog-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1.5rem;
  border-bottom: 1px solid #e0e0e0;

  h3 {
    margin: 0;
    color: #333;
    font-size: 1.25rem;
  }
}

.close-button {
  background: none;
  border: none;
  font-size: 1.5rem;
  cursor: pointer;
  color: #666;
  padding: 0;
  width: 30px;
  height: 30px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 4px;

  &:hover {
    background-color: #f0f0f0;
    color: #333;
  }
}

.dialog-content {
  padding: 1.5rem;
  max-height: calc(90vh - 140px);
  overflow-y: auto;
}

.selected-images-info {
  background-color: #f8f9fa;
  padding: 1rem;
  border-radius: 4px;
  margin-bottom: 1.5rem;

  p {
    margin: 0;
    color: #666;
    font-size: 0.9rem;
  }
}

.album-list-section {
  margin-bottom: 1.5rem;

  h4 {
    margin: 0 0 1rem 0;
    color: #333;
    font-size: 1rem;
  }
}

.loading-state {
  text-align: center;
  padding: 2rem;
  color: #666;
}

.error-state {
  text-align: center;
  padding: 2rem;

  p {
    color: #dc3545;
    margin-bottom: 1rem;
  }
}

.retry-button {
  background-color: #6c757d;
  color: white;
  border: none;
  padding: 0.5rem 1rem;
  border-radius: 4px;
  cursor: pointer;

  &:hover {
    background-color: #5a6268;
  }
}

.album-list {
  max-height: 300px;
  overflow-y: auto;
  border: 1px solid #e0e0e0;
  border-radius: 4px;
}

.album-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem;
  border-bottom: 1px solid #f0f0f0;
  cursor: pointer;
  transition: background-color 0.2s;

  &:last-child {
    border-bottom: none;
  }

  &:hover {
    background-color: #f8f9fa;
  }

  &.selected {
    background-color: #e3f2fd;
    border-color: #90caf9;
  }
}

.album-info {
  flex: 1;

  h5 {
    margin: 0 0 0.5rem 0;
    color: #333;
    font-size: 0.95rem;
  }
}

.album-meta {
  font-size: 0.8rem;
  color: #999;
}

.album-selection {
  margin-left: 1rem;
}

.empty-state {
  text-align: center;
  padding: 2rem;
  color: #666;

  .hint {
    font-size: 0.9rem;
    color: #999;
    margin-top: 0.5rem;
  }
}

.error-message {
  background-color: #f8d7da;
  color: #721c24;
  padding: 1rem;
  border-radius: 4px;
  margin-bottom: 1.5rem;
  font-size: 0.9rem;
  border: 1px solid #f5c6cb;
}

.dialog-actions {
  display: flex;
  gap: 1rem;
  justify-content: flex-end;

  button {
    padding: 0.75rem 1.5rem;
    border: none;
    border-radius: 4px;
    font-size: 1rem;
    cursor: pointer;
    transition: all 0.2s;

    &:disabled {
      cursor: not-allowed;
      opacity: 0.6;
    }
  }
}

.cancel-button {
  background-color: #6c757d;
  color: white;

  &:hover:not(:disabled) {
    background-color: #5a6268;
  }
}

.add-button {
  background-color: #28a745;
  color: white;

  &:hover:not(:disabled) {
    background-color: #218838;
  }
}
</style>
