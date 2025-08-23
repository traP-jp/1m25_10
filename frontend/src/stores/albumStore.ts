import { defineStore } from 'pinia'
import { albumService } from '@/services'
import type {
  Album,
  AlbumItem,
  CreateAlbumRequest,
  UpdateAlbumRequest,
  GetAlbumsParams,
} from '@/types'

export const useAlbumStore = defineStore('album', {
  state: () => ({
    albums: [] as AlbumItem[],
    albumDetails: {} as Record<string, Album>, // アルバム詳細をキャッシュ
    currentAlbum: null as Album | null,
    loading: false,
    error: null as string | null,
  }),

  getters: {
    // アルバム数
    albumCount: (state) => state.albums.length,

    // 特定のユーザーのアルバム
    albumsByCreator: (state) => {
      return (creatorId: string) => state.albums.filter((album) => album.creator === creatorId)
    },

    // キャッシュされたアルバム詳細を取得
    getAlbumDetail: (state) => {
      return (albumId: string) => state.albumDetails[albumId]
    },
  },

  actions: {
    // 全アルバム取得
    async fetchAlbums(params?: GetAlbumsParams) {
      this.loading = true
      this.error = null

      try {
        const albums = await albumService.getAlbums(params)
        this.albums = albums
      } catch (error) {
        this.error = error instanceof Error ? error.message : 'Unknown error occurred'
        console.error('Error fetching albums:', error)
      } finally {
        this.loading = false
      }
    },

    // 特定のアルバム取得
    async fetchAlbum(albumId: string) {
      this.loading = true
      this.error = null

      try {
        const album = await albumService.getAlbumDetail(albumId)
        this.currentAlbum = album

        // アルバム詳細をキャッシュ
        this.albumDetails[albumId] = album

        // アルバムリストも更新（すでに存在する場合）
        const existingIndex = this.albums.findIndex((a) => a.id === albumId)
        if (existingIndex !== -1) {
          // AlbumItemとして必要なプロパティのみを更新
          this.albums[existingIndex] = {
            id: album.id,
            title: album.title,
            creator: album.creator,
          }
        }

        return album
      } catch (error) {
        this.error = error instanceof Error ? error.message : 'Unknown error occurred'
        console.error('Error fetching album:', error)
        throw error
      } finally {
        this.loading = false
      }
    },

    // アルバム作成
    async createAlbum(albumData: CreateAlbumRequest) {
      this.loading = true
      this.error = null

      try {
        const newAlbum = await albumService.createAlbum(albumData)

        // アルバム詳細をキャッシュ
        this.albumDetails[newAlbum.id] = newAlbum

        // アルバムリストに追加（AlbumItemとして）
        const albumItem: AlbumItem = {
          id: newAlbum.id,
          title: newAlbum.title,
          creator: newAlbum.creator,
        }
        this.albums.push(albumItem)

        return newAlbum
      } catch (error) {
        this.error = error instanceof Error ? error.message : 'Unknown error occurred'
        console.error('Error creating album:', error)
        throw error
      } finally {
        this.loading = false
      }
    },

    // アルバム更新
    async updateAlbum(albumId: string, albumData: UpdateAlbumRequest) {
      this.loading = true
      this.error = null

      try {
        const updatedAlbum = await albumService.updateAlbum(albumId, albumData)

        // アルバム詳細を更新
        this.albumDetails[albumId] = updatedAlbum

        // アルバムリストを更新
        const index = this.albums.findIndex((album) => album.id === albumId)
        if (index !== -1) {
          this.albums[index] = {
            id: updatedAlbum.id,
            title: updatedAlbum.title,
            creator: updatedAlbum.creator,
          }
        }

        if (this.currentAlbum?.id === albumId) {
          this.currentAlbum = updatedAlbum
        }

        return updatedAlbum
      } catch (error) {
        this.error = error instanceof Error ? error.message : 'Unknown error occurred'
        console.error('Error updating album:', error)
        throw error
      } finally {
        this.loading = false
      }
    },

    // アルバム削除
    async deleteAlbum(albumId: string) {
      this.loading = true
      this.error = null

      try {
        await albumService.deleteAlbum(albumId)

        // ローカルデータから削除
        this.albums = this.albums.filter((album) => album.id !== albumId)
        delete this.albumDetails[albumId]

        if (this.currentAlbum?.id === albumId) {
          this.currentAlbum = null
        }
      } catch (error) {
        this.error = error instanceof Error ? error.message : 'Unknown error occurred'
        console.error('Error deleting album:', error)
        throw error
      } finally {
        this.loading = false
      }
    },

    // 特定のユーザーのアルバム取得
    async fetchAlbumsByCreator(creatorId: string) {
      this.loading = true
      this.error = null

      try {
        const albums = await albumService.getAlbumsByCreator(creatorId)
        this.albums = albums
      } catch (error) {
        this.error = error instanceof Error ? error.message : 'Unknown error occurred'
        console.error('Error fetching albums by creator:', error)
      } finally {
        this.loading = false
      }
    },

    // エラーをクリア
    clearError() {
      this.error = null
    },

    // 現在のアルバムをクリア
    clearCurrentAlbum() {
      this.currentAlbum = null
    },
  },
})
