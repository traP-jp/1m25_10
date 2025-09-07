import { defineStore } from 'pinia'
import { imageService, AlbumService } from '@/services'
import { generateImageUrl } from '@/config/env'
import type { Image, ImageDetail } from '@/types'

const albumService = new AlbumService()

export const useImageStore = defineStore('image', {
  state: () => ({
    images: [] as Image[],
    imageDetails: {} as Record<string, ImageDetail>,
    currentImage: null as ImageDetail | null,
    selectedImageIds: new Set<string>(),
    loading: false,
    loadingMore: false,
    error: null as string | null,
    hasMore: true,
    currentOffset: 0,
    currentSearchQuery: undefined as string | undefined,
    pageSize: 20,
  albumChance: false, // アルバムチャンスフィルタの有効/無効
  }),

  getters: {
    imageCount: (state) => state.images.length,
    selectedImageCount: (state) => state.selectedImageIds.size,
    selectedImages: (state) => {
      return state.images.filter((image) => state.selectedImageIds.has(image.id))
    },
    getImageDetail: (state) => {
      return (imageId: string) => state.imageDetails[imageId]
    },
    getImageUrl: () => {
      return (image: Image): string => {
        // 環境設定に基づいて画像URLを生成
        return generateImageUrl(image.id)
      }
    },
  },

  actions: {
    async fetchImages(searchQuery?: string) {
      this.loading = true
      this.error = null
      this.currentOffset = 0
      this.currentSearchQuery = searchQuery
      this.hasMore = true

      try {
        console.log('Search query:', searchQuery)
        const images = await imageService.getImages(
          searchQuery,
          this.pageSize,
          0,
          { albumChance: this.albumChance },
        )
        this.images = images

        if (images.length < this.pageSize) {
          this.hasMore = false
        } else {
          this.currentOffset = this.pageSize
        }
      } catch (error) {
        this.error = error instanceof Error ? error.message : 'Unknown error occurred'
        console.error('Error fetching images:', error)
      } finally {
        this.loading = false
      }
    },

    async loadMoreImages() {
      if (!this.hasMore || this.loadingMore) {
        return
      }

      this.loadingMore = true
      this.error = null

      try {
        let nextOffset = this.currentOffset
        let moreImages: Image[] = []

        // 空ウィンドウをスキップしながら取得（最大数回のセーフガード）
        for (let attempt = 0; attempt < 5; attempt++) {
          moreImages = await imageService.getImages(
            this.currentSearchQuery,
            this.pageSize,
            nextOffset,
            { albumChance: this.albumChance },
          )

          if (moreImages.length > 0 || !this.albumChance) {
            // 画像が返ってきた、またはアルバムチャンスでないなら終了
            this.currentOffset = nextOffset
            break
          }

          // 0件かつアルバムチャンス => バックエンド側がフィルタして0件の可能性
          // ウィンドウをひとつ進めて再試行
          nextOffset += this.pageSize
        }

        if (moreImages.length === 0) {
          this.hasMore = false
        } else {
          this.images.push(...moreImages)
          this.currentOffset += moreImages.length

          if (moreImages.length < this.pageSize) {
            this.hasMore = false
          }
        }
      } catch (error) {
        this.error = error instanceof Error ? error.message : 'Unknown error occurred'
        console.error('Error loading more images:', error)
      } finally {
        this.loadingMore = false
      }
    },

    setAlbumChance(enabled: boolean) {
      this.albumChance = enabled
    },

    // 画像詳細を取得（必要な時のみ）
    async fetchImageDetailIfNeeded(imageId: string): Promise<ImageDetail> {
      // 既にキャッシュされている場合はそれを返す
      if (this.imageDetails[imageId]) {
        return this.imageDetails[imageId]
      }

      // キャッシュされていない場合は取得
      return this.fetchImageDetail(imageId)
    },

    async fetchImageDetail(imageId: string): Promise<ImageDetail> {
      this.loading = true
      this.error = null

      try {
        const imageDetail = await imageService.getImageDetail(imageId)
        this.currentImage = imageDetail
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

    clearError() {
      this.error = null
    },

    clearCurrentImage() {
      this.currentImage = null
    },

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

        const createdAlbum = await albumService.createAlbum(albumData)
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
