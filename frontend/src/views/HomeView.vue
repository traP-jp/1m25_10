<template>
  <div :class="$style.homeView">
    <header :class="$style.header">
      <div :class="$style.titleSection">
        <h1 :class="$style.title">Home</h1>
        <p :class="$style.subtitle">すべての画像 ({{ imageStore.imageCount }}枚)</p>
      </div>

      <!-- 検索・フィルタ -->
      <div :class="[$style.controls, { [$style.hasSelection]: imageStore.selectedImageCount > 0 }]">
        <div :class="$style.searchContainer">
          <div :class="$style.searchInputWrapper">
            <input
              v-model="searchQuery"
              type="text"
              placeholder="画像を検索... (Enterで実行)"
              :class="$style.searchInput"
              @keyup.enter="performSearch"
            />
            <button
              v-if="hasSearchQuery"
              type="button"
              @click="clearSearch"
              :class="$style.clearSearchButton"
              title="検索をクリア"
            >
              ×
            </button>
          </div>
          <label :class="$style.switchLabel">
            <input type="checkbox" v-model="albumChance" @change="onToggleAlbumChance" />
            <span>アルバムチャンス</span>
          </label>
          <button
            type="button"
            @click="performSearch"
            :class="$style.searchButton"
            :disabled="imageStore.loading"
          >
            {{ imageStore.loading ? '検索中...' : '検索' }}
          </button>
        </div>

        <!-- 選択された画像がある場合の操作ボタン -->
        <div
          :class="[
            $style.selectionActions,
            { [$style.visible]: imageStore.selectedImageCount > 0 },
          ]"
        >
          <button
            @click="imageStore.deselectAllImages"
            :class="$style.clearSelectionButton"
            title="選択を解除"
          >
            ×
          </button>
          <span :class="$style.selectionCount"> {{ imageStore.selectedImageCount }}枚選択中 </span>
          <button @click="showCreateAlbumDialog" :class="$style.createAlbumButton">
            アルバムを作成
          </button>
          <button @click="showAddToAlbumDialog" :class="$style.addToAlbumButton">
            アルバムに追加
          </button>
        </div>
      </div>
    </header>

    <!-- アルバム作成ダイアログ -->
    <CreateAlbumDialog
      :is-visible="showAlbumDialog"
      :selected-image-count="imageStore.selectedImageCount"
      @close="closeDialog"
      @create="createAlbum"
    />

    <!-- アルバムに追加ダイアログ -->
    <AddToAlbumDialog
      :is-visible="showAddAlbumDialog"
      :selected-image-count="imageStore.selectedImageCount"
      @close="closeAddDialog"
      @add-to-album="addToAlbum"
    />

    <!-- 画像一覧 -->
    <ImageList
      :images="images"
      :loading="imageStore.loading"
      :error="imageStore.error"
      :selected-image-ids="imageStore.selectedImageIds"
      :has-search-query="!!searchQuery"
      :has-more="imageStore.hasMore"
      :loading-more="imageStore.loadingMore"
      :get-image-url="imageStore.getImageUrl"
      @retry="retryLoad"
      @toggle-selection="imageStore.toggleImageSelection"
      @load-more="loadMoreImages"
      @open-detail="openDetail"
    />

    <!-- 画像詳細 右側パネル -->
    <ImageDetailSidePanel :visible="detailVisible" :image-id="detailImageId" @close="closeDetail" />
  </div>
</template>

<script setup lang="ts">
import { computed, ref, onMounted } from 'vue'
import { useImageStore } from '@/stores/imageStore'
import { albumService } from '@/services'
import ImageList from '@/components/ImageList.vue'
import CreateAlbumDialog from '@/components/CreateAlbumDialog.vue'
import AddToAlbumDialog from '@/components/AddToAlbumDialog.vue'
import ImageDetailSidePanel from '@/components/ImageDetailSidePanel.vue'

const imageStore = useImageStore()

const searchQuery = ref('')
const images = computed(() => imageStore.images)
const hasSearchQuery = computed(() => searchQuery.value.trim() !== '')
const albumChance = ref(imageStore.albumChance)

const performSearch = () => {
  imageStore.setAlbumChance(albumChance.value)
  imageStore.fetchImages(searchQuery.value.trim() || undefined)
}

const clearSearch = () => {
  searchQuery.value = ''
  imageStore.fetchImages()
}

const showAlbumDialog = ref(false)
const showAddAlbumDialog = ref(false)

const showCreateAlbumDialog = () => {
  showAlbumDialog.value = true
}

const showAddToAlbumDialog = () => {
  showAddAlbumDialog.value = true
}

const closeDialog = () => {
  showAlbumDialog.value = false
}

const closeAddDialog = () => {
  showAddAlbumDialog.value = false
}

const createAlbum = async (data: { title: string; description: string }): Promise<void> => {
  const result = await imageStore.createAlbumFromSelectedImages(
    data.title,
    data.description || undefined,
  )

  alert(`アルバム「${result.title}」を作成しました！（${result.imageCount}枚の画像）`)
  closeDialog()
}

const addToAlbum = async (data: { albumId: string }): Promise<void> => {
  try {
    const selectedImageIds = Array.from(imageStore.selectedImageIds)
    const updatedAlbum = await albumService.addImagesToAlbum(data.albumId, selectedImageIds)

    alert(`${selectedImageIds.length}枚の画像をアルバム「${updatedAlbum.title}」に追加しました！`)
    imageStore.deselectAllImages()
    closeAddDialog()
  } catch (error) {
    console.error('Failed to add to album:', error)
    throw error
  }
}

const retryLoad = () => {
  imageStore.clearError()
  imageStore.fetchImages(searchQuery.value || undefined)
}

const loadMoreImages = () => {
  imageStore.loadMoreImages()
}

onMounted(() => {
  imageStore.fetchImages()
})

const onToggleAlbumChance = () => {
  imageStore.setAlbumChance(albumChance.value)
  imageStore.fetchImages(searchQuery.value.trim() || undefined)
}

// 画像詳細モーダルの制御
const detailVisible = ref(false)
const detailImageId = ref<string | null>(null)
const openDetail = (imageId: string) => {
  detailImageId.value = imageId
  detailVisible.value = true
}
const closeDetail = () => {
  detailVisible.value = false
  detailImageId.value = null
}
</script>
<style lang="scss" module>
.homeView {
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

.switchLabel {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  margin-left: 12px;
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
  gap: 0;
  align-items: center;
  padding: 16px;
  background-color: #f8f9fa;
  border-radius: 8px;
  border: 1px solid #e9ecef;
  min-height: 64px;
  transition: gap 0.3s ease;

  &.hasSelection {
    gap: 16px;
  }
}

.searchContainer {
  flex: 1;
  display: flex;
  gap: 8px;
  align-items: center;
  min-width: 0;
}

.searchInputWrapper {
  flex: 1;
  position: relative;
  display: flex;
  align-items: center;
}

.searchInput {
  width: 100%;
  padding: 8px 12px;
  padding-right: 36px;
  border: 1px solid #ddd;
  border-radius: 6px;
  font-size: 14px;

  &:focus {
    outline: none;
    border-color: #005bac;
    box-shadow: 0 0 0 3px rgba(0, 91, 172, 0.1);
  }
}

.clearSearchButton {
  position: absolute;
  right: 8px;
  top: 50%;
  transform: translateY(-50%);
  background: none;
  border: none;
  color: #666;
  font-size: 18px;
  width: 20px;
  height: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  border-radius: 50%;
  transition: all 0.2s;

  &:hover {
    background-color: #f0f0f0;
    color: #333;
  }

  &:active {
    background-color: #e0e0e0;
  }
}

.searchButton {
  padding: 8px 16px;
  background-color: #007bff;
  color: white;
  border: none;
  border-radius: 6px;
  font-size: 14px;
  cursor: pointer;
  transition: background-color 0.2s;
  white-space: nowrap;

  &:hover:not(:disabled) {
    background-color: #0056b3;
  }

  &:disabled {
    background-color: #6c757d;
    cursor: not-allowed;
  }
}

.selectionActions {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 8px 0;
  background-color: transparent;
  border-radius: 0;
  border: 1px solid transparent;
  opacity: 0;
  max-width: 0;
  overflow: hidden;
  transition: all 0.3s ease;
  pointer-events: none;
  white-space: nowrap;

  &.visible {
    opacity: 1;
    max-width: 500px;
    padding: 8px 16px;
    background-color: #e3f2fd;
    border-radius: 6px;
    border: 1px solid #90caf9;
    pointer-events: auto;
  }
}

.clearSelectionButton {
  background: none;
  border: none;
  color: #666;
  font-size: 18px;
  width: 20px;
  height: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  border-radius: 50%;
  transition: all 0.2s;
  margin-right: -4px;

  &:hover {
    background-color: #f0f0f0;
    color: #333;
  }

  &:active {
    background-color: #e0e0e0;
  }
}

.selectionCount {
  font-size: 14px;
  font-weight: 500;
  color: #1565c0;
}

.createAlbumButton {
  padding: 6px 12px;
  background-color: #007bff;
  color: white;
  border: none;
  border-radius: 4px;
  font-size: 12px;
  cursor: pointer;
  transition: background-color 0.2s;

  &:hover {
    background-color: #0056b3;
  }
}

.addToAlbumButton {
  padding: 6px 12px;
  background-color: #28a745;
  color: white;
  border: none;
  border-radius: 4px;
  font-size: 12px;
  cursor: pointer;
  transition: background-color 0.2s;

  &:hover {
    background-color: #218838;
  }
}

// レスポンシブデザイン
@media (max-width: 768px) {
  .homeView {
    padding: 16px;
  }

  .controls {
    flex-direction: column;
    align-items: stretch;
    min-height: auto;
  }

  .searchContainer {
    flex-direction: column;
    gap: 8px;
  }

  .selectionActions {
    &.visible {
      max-width: none;
    }
  }

  .title {
    font-size: 1.5rem;
  }
}

@media (max-width: 480px) {
  .controls {
    padding: 12px;
  }

  .searchButton {
    width: 100%;
  }

  .searchInput {
    font-size: 16px; // iOS Safariでズームを防ぐ
  }
}
</style>
