import { defineStore } from 'pinia'
import { imageService, AlbumService } from '@/services'
import { generateImageUrl, env } from '@/config/env'
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
  albumChance: env.VITE_ALBUM_CHANCE_DEFAULT === true, // アルバムチャンス既定値（env）
  scanOffset: 0, // アルバムチャンス時のスキャン用オフセット
  totalHits: undefined as number | undefined, // traQのtotalHits（参考）
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
      this.scanOffset = 0
      this.currentSearchQuery = searchQuery
      this.hasMore = true

      try {
        console.log('Search query:', searchQuery)
        let usedOffset = 0
        let gathered: Image[] = []
        let totalHits: number | undefined = undefined

        if (!this.albumChance) {
          const res = await imageService.getImages(searchQuery, this.pageSize, 0)
          gathered = res.images
          totalHits = res.totalHits
          usedOffset = 0
        } else {
          // アルバムチャンス: 空ウィンドウをスキップして最初の非空ウィンドウを探す
          const maxAttempts = 50 // セーフガード
          let attempts = 0
          let upperBound = Number.MAX_SAFE_INTEGER
          while (attempts < maxAttempts) {
            const res = await imageService.getImages(searchQuery, this.pageSize, usedOffset, {
              albumChance: true,
            })
            totalHits = res.totalHits
            if (typeof totalHits === 'number') upperBound = totalHits
            if (res.images.length > 0) {
              gathered = res.images
              break
            }
            // 空だったので次のウィンドウへ
            usedOffset += this.pageSize
            if (usedOffset >= upperBound) break
            attempts++
          }
        }

        this.images = gathered
        this.totalHits = totalHits
        this.currentOffset = usedOffset + gathered.length
        this.scanOffset = usedOffset + this.pageSize
        this.hasMore = gathered.length >= this.pageSize && (this.totalHits == null || this.currentOffset < this.totalHits)
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
        if (!this.albumChance) {
          const res = await imageService.getImages(
            this.currentSearchQuery,
            this.pageSize,
            this.currentOffset,
          )
          this.totalHits = res.totalHits
          const more = res.images
          if (more.length === 0) {
            this.hasMore = false
          } else {
            this.images.push(...more)
            this.currentOffset += more.length
            if (more.length < this.pageSize) this.hasMore = false
          }
        } else {
          // アルバムチャンス: 空ウィンドウをスキャン
          let usedOffset = Math.max(this.scanOffset, this.currentOffset)
          const upperBound = this.totalHits ?? Number.MAX_SAFE_INTEGER
          const maxAttempts = 50
          let attempts = 0
          let found: Image[] = []
          while (usedOffset < upperBound && attempts < maxAttempts) {
            const res = await imageService.getImages(
              this.currentSearchQuery,
              this.pageSize,
              usedOffset,
              { albumChance: true },
            )
            this.totalHits = res.totalHits
            if (res.images.length > 0) {
              found = res.images
              break
            }
            usedOffset += this.pageSize
            attempts++
          }

          this.scanOffset = usedOffset + this.pageSize
          if (found.length === 0) {
            this.hasMore = false
          } else {
            this.images.push(...found)
            this.currentOffset = usedOffset + found.length
            if (found.length < this.pageSize) this.hasMore = false
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
