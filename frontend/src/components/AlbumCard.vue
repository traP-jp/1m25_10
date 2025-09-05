<template>
  <div :class="$style.albumCard">
    <div :class="$style.imageContainer">
      <div
        v-if="imageUrls && imageUrls.length > 0"
        :class="[$style.images, $style[`grid${Math.min(imageUrls.length, 4)}`]]"
      >
        <!-- 最大4枚の画像を表示 -->
        <div
          v-for="(imageUrl, index) in displayImages"
          :key="`${imageUrl}-${index}`"
          :class="[$style.imageItem, $style[getImageItemClass(imageUrls.length, index)]]"
        >
          <img
            :src="imageUrl"
            :alt="`Album ${title} image ${index + 1}`"
            :class="$style.image"
            @error="handleImageError"
            loading="lazy"
          />
        </div>
        <!-- 画像が4枚より多い場合の残り枚数表示 -->
        <div
          v-if="imageUrls.length > 4"
          :class="$style.moreImages"
          :aria-label="`${imageUrls.length - 4} more images`"
        >
          +{{ imageUrls.length - 4 }}
        </div>
      </div>
      <div v-else :class="$style.noImage">
        <div :class="$style.noImageIcon" aria-hidden="true">📷</div>
        <div :class="$style.noImageText">No Images</div>
      </div>
    </div>
    <div :class="$style.content">
      <h3 :class="$style.title">{{ title }}</h3>
      <div :class="$style.meta">
        <div :class="$style.leftMeta">
          <div :class="$style.creator">
            <img
              v-if="creator.avatarUrl"
              :src="creator.avatarUrl"
              :alt="`${creator.name}'s avatar`"
              :class="$style.creatorAvatar"
              @error="handleAvatarError"
            />
            <div v-else :class="$style.creatorAvatarPlaceholder">
              {{ creator.name.charAt(0).toUpperCase() }}
            </div>
            <span :class="$style.creatorName">{{ creator.name }}</span>
          </div>
        </div>
        <div :class="$style.rightMeta">
          <span :class="$style.imageCount">{{ imageUrls.length }} images</span>
          <span :class="$style.date">{{ formatDate(createdAt) }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

interface Props {
  id: string
  title: string
  description?: string // オプショナルに変更
  imageUrls: string[]
  createdAt: string
  creator: {
    id: string
    name: string
    avatarUrl?: string
  }
}

const props = defineProps<Props>()

// 表示する画像（最大4枚）
const displayImages = computed(() => props.imageUrls.slice(0, 4))

// 画像読み込みエラー時の処理
const handleImageError = (event: Event) => {
  const img = event.target as HTMLImageElement
  if (img) {
    img.src = '/dummyAlbumsIcon.png' // フォールバック画像
  }
}

// アバター読み込みエラー時の処理
const handleAvatarError = (event: Event) => {
  const img = event.target as HTMLImageElement
  if (img) {
    // アバターエラー時は非表示にして、プレースホルダーを表示
    img.style.display = 'none'
  }
}

// 画像アイテムのクラスを決定する関数
const getImageItemClass = (totalImages: number, index: number): string => {
  const count = Math.min(totalImages, 4)

  if (count === 1) return 'single'
  if (count === 2) return 'half'
  if (count === 3) {
    return index === 0 ? 'thirdMain' : 'thirdSub'
  }
  return 'quarter'
}

// 日付をフォーマットする関数
const formatDate = (dateString: string): string => {
  const date = new Date(dateString)
  return date.toLocaleDateString('ja-JP', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
  })
}
</script>

<style lang="scss" module>
.albumCard {
  background: white;
  border-radius: 8px;
  border: 1px solid var(--Light-UI-Tertiary, #ced6db);
  overflow: hidden;
  transition: all 0.15s cubic-bezier(0.4, 0, 0.2, 1);
  cursor: pointer;
  max-width: 300px;
  user-select: none;
  display: flex;
  flex-direction: column;
  height: 330px; /* 固定高さを設定 */

  &:hover,
  &:focus-visible {
    transform: translateY(-2px);
    box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12);
    border-color: rgba(0, 91, 172, 0.2);
  }

  &:focus-visible {
    outline: 2px solid rgba(0, 91, 172, 0.5);
    outline-offset: 2px;
  }
}

.imageContainer {
  position: relative;
  width: 100%;
  height: 200px;
  background: var(--Light-UI-Tertiary, #ced6db);
  overflow: hidden;
}

.images {
  position: relative;
  width: 100%;
  height: 100%;
  display: grid;
  gap: 2px;

  // デフォルト: 4枚以上の場合
  grid-template-columns: 1fr 1fr;
  grid-template-rows: 1fr 1fr;
}

// 画像数に応じたグリッドレイアウト
.grid1 {
  grid-template-columns: 1fr;
  grid-template-rows: 1fr;
}

.grid2 {
  grid-template-columns: 1fr 1fr;
  grid-template-rows: 1fr;
}

.grid3 {
  grid-template-columns: 1fr 1fr;
  grid-template-rows: 1fr 1fr;
}

.grid4 {
  grid-template-columns: 1fr 1fr;
  grid-template-rows: 1fr 1fr;
}

.imageItem {
  position: relative;
  overflow: hidden;
}

.single {
  grid-column: 1 / -1;
  grid-row: 1 / -1;
}

.half {
  grid-row: 1;

  &:nth-child(1) {
    grid-column: 1;
  }
  &:nth-child(2) {
    grid-column: 2;
  }
}

.thirdMain {
  grid-column: 1;
  grid-row: 1 / -1;
}

.thirdSub {
  grid-column: 2;

  &:nth-child(2) {
    grid-row: 1;
  }
  &:nth-child(3) {
    grid-row: 2;
  }
}

.image {
  width: 100%;
  height: 100%;
  object-fit: cover;
  display: block;
  transition: transform 0.15s cubic-bezier(0.4, 0, 0.2, 1);

  &:hover {
    transform: scale(1.05);
  }
}

.moreImages {
  position: absolute;
  bottom: 8px;
  right: 8px;
  background: linear-gradient(135deg, rgba(0, 91, 172, 0.9) 0%, rgba(0, 74, 148, 0.9) 100%);
  color: white;
  padding: 4px 8px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 600;
  backdrop-filter: blur(8px);
  box-shadow: 0 2px 8px rgba(0, 91, 172, 0.3);
}

.noImage {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  color: #666;
}

.noImageIcon {
  font-size: 48px;
  margin-bottom: 8px;
  opacity: 0.6;
}

.noImageText {
  font-size: 14px;
  font-weight: 500;
  opacity: 0.8;
}

.content {
  padding: 16px;
  display: flex;
  flex-direction: column;
  flex: 1; /* 残りのスペースを全て使用 */
}

.title {
  margin: 0 0 12px 0;
  font-size: 18px;
  font-weight: 600;
  color: #333;
  line-height: 1.3;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  line-clamp: 2;
  overflow: hidden;
  flex-shrink: 0; /* タイトルの縮小を防ぐ */
}

.meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 8px;
  margin-top: auto; /* これによりmetaが下端に固定される */
}

.leftMeta {
  flex: 1;
  min-width: 0; // フレックスアイテムの縮小を許可
}

.rightMeta {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 2px;
  font-size: 12px;
  color: #999;
}

.creator {
  display: flex;
  align-items: center;
  gap: 8px;
  min-width: 0;
}

.creatorAvatar {
  width: 24px;
  height: 24px;
  border-radius: 50%;
  object-fit: cover;
  border: 1px solid rgba(0, 0, 0, 0.1);
}

.creatorAvatarPlaceholder {
  width: 24px;
  height: 24px;
  border-radius: 50%;
  background: linear-gradient(135deg, rgba(0, 91, 172, 0.8) 0%, rgba(0, 74, 148, 0.8) 100%);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: 600;
}

.creatorName {
  font-size: 13px;
  font-weight: 500;
  color: #555;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  min-width: 0;
}

.imageCount {
  font-weight: 600;
  color: rgba(0, 91, 172, 0.8);
}

.date {
  font-weight: 400;
}
</style>
