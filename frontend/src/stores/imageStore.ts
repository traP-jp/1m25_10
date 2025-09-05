import { defineStore } from 'pinia'
import { imageService } from '@/services'
import type { Image, ImageDetail } from '@/types'

export const useImageStore = defineStore('image', {
  state: () => ({
    images: [] as Image[],
    imageDetails: {} as Record<string, ImageDetail>, // 画像詳細をキャッシュ
    currentImage: null as ImageDetail | null,
    loading: false,
    error: null as string | null,
  }),

  getters: {
    // 画像数
    imageCount: (state) => state.images.length,

    // 特定のユーザーの画像
    imagesByCreator: (state) => {
      return (creatorId: string) => state.images.filter((image) => image.creator === creatorId)
    },

    // キャッシュされた画像詳細を取得
    getImageDetail: (state) => {
      return (imageId: string) => state.imageDetails[imageId]
    },
  },

  actions: {
    // 全画像取得
    async fetchImages() {
      this.loading = true
      this.error = null

      try {
        const images = await imageService.getImages()
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
  },
})
