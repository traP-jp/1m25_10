import { defineStore } from 'pinia'
import { imageService, AlbumService } from '@/services'
import { generateImageUrl } from '@/config/env'
import type { Image, ImageDetail } from '@/types'

const albumService = new AlbumService()

export const useImageStore = defineStore('image', {
  state: () => ({
    images: [] as Image[],
    imageDetails: {} as Record<string, ImageDetail>, // 画像詳細をキャッシュ
    currentImage: null as ImageDetail | null,
    selectedImageIds: new Set<string>(), // 選択された画像ID
    loading: false,
    error: null as string | null,
  }),

  getters: {
    // 画像数
    imageCount: (state) => state.images.length,

    // 選択された画像数
    selectedImageCount: (state) => state.selectedImageIds.size,

    // 選択された画像一覧
    selectedImages: (state) => {
      return state.images.filter((image) => state.selectedImageIds.has(image.id))
    },

    // 特定のユーザーの画像
    imagesByCreator: (state) => {
      return (creatorId: string) => state.images.filter((image) => image.creator === creatorId)
    },

    // キャッシュされた画像詳細を取得
    getImageDetail: (state) => {
      return (imageId: string) => state.imageDetails[imageId]
    },

    // 画像URLを生成
    getImageUrl: () => {
      return (image: Image): string => {
        // 環境設定に基づいて画像URLを生成
        return generateImageUrl(image.id)
      }
    },
  },

  actions: {
    // 全画像取得（検索クエリ対応）
    async fetchImages(searchQuery?: string) {
      this.loading = true
      this.error = null

      try {
        console.log('Search query:', searchQuery) // デバッグ用ログ
        const images = await imageService.getImages(searchQuery)
        this.images = images
      } catch (error) {
        this.error = error instanceof Error ? error.message : 'Unknown error occurred'
        console.error('Error fetching images:', error)
      } finally {
        this.loading = false
      }
    },

    // 特定の画像詳細取得
    async fetchImageDetail(imageId: string) {
      this.loading = true
      this.error = null

      try {
        const imageDetail = await imageService.getImageDetail(imageId)
        this.currentImage = imageDetail

        // 画像詳細をキャッシュ
        this.imageDetails[imageId] = imageDetail

        return imageDetail
      } catch (error) {
        this.error = error instanceof Error ? error.message : 'Unknown error occurred'
        console.error('Error fetching image detail:', error)
        throw error
      } finally {
        this.loading = false
      }
    },

    // エラーをクリア
    clearError() {
      this.error = null
    },

    // 現在の画像をクリア
    clearCurrentImage() {
      this.currentImage = null
    },

    // 画像選択機能
    selectImage(imageId: string) {
      this.selectedImageIds.add(imageId)
    },

    deselectImage(imageId: string) {
      this.selectedImageIds.delete(imageId)
    },

    toggleImageSelection(imageId: string) {
      if (this.selectedImageIds.has(imageId)) {
        this.selectedImageIds.delete(imageId)
      } else {
        this.selectedImageIds.add(imageId)
      }
    },

    deselectAllImages() {
      this.selectedImageIds.clear()
    },

    // バルク操作
    async createAlbumFromSelectedImages(albumTitle: string, albumDescription?: string) {
      if (this.selectedImageIds.size === 0) {
        throw new Error('画像が選択されていません')
      }

      this.loading = true
      this.error = null

      try {
        const selectedIds = Array.from(this.selectedImageIds)

        const albumData = {
          title: albumTitle,
          description: albumDescription || '',
          images: selectedIds,
        }

        console.log('Creating album with data:', albumData)

        // 実際のアルバム作成APIを呼び出し
        const createdAlbum = await albumService.createAlbum(albumData)

        // 選択をクリア
        this.selectedImageIds.clear()

        return {
          id: createdAlbum.id,
          title: createdAlbum.title,
          description: createdAlbum.description,
          imageCount: selectedIds.length,
        }
      } catch (error) {
        this.error = error instanceof Error ? error.message : 'Unknown error occurred'
        console.error('Error creating album from selected images:', error)
        throw error
      } finally {
        this.loading = false
      }
    },
  },
})
