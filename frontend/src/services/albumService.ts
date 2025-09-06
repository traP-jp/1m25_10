// アルバム関連のAPIサービス（OpenAPI仕様準拠）

import { apiClient } from './apiClient'
import type {
  Album,
  AlbumItem,
  CreateAlbumRequest,
  UpdateAlbumRequest,
  GetAlbumsParams,
} from '@/types'

export class AlbumService {
  // 全アルバム取得（GET /albums）
  async getAlbums(params?: GetAlbumsParams): Promise<AlbumItem[]> {
    const queryParams: Record<string, unknown> = {}

    if (params?.creator_id) queryParams.creator_id = params.creator_id
    if (params?.before_date) queryParams.before_date = params.before_date
    if (params?.after_date) queryParams.after_date = params.after_date
    if (params?.limit) queryParams.limit = params.limit
    if (params?.offset) queryParams.offset = params.offset

    return apiClient.get<AlbumItem[]>('/albums', queryParams)
  }

  // 特定のアルバム取得（GET /albums/{id}）
  async getAlbumDetail(albumId: string): Promise<Album> {
    return apiClient.get<Album>(`/albums/${albumId}`)
  }

  // アルバム作成（POST /albums）
  async createAlbum(albumData: CreateAlbumRequest): Promise<Album> {
    return apiClient.post<Album>('/albums', albumData)
  }

  // アルバム更新（PUT /albums/{id}）
  async updateAlbum(albumId: string, albumData: UpdateAlbumRequest): Promise<Album> {
    return apiClient.put<Album>(`/albums/${albumId}`, albumData)
  }

  // 選択した画像をアルバムに追加
  async addImagesToAlbum(albumId: string, imageIds: string[]): Promise<Album> {
    // まず現在のアルバム詳細を取得
    const currentAlbum = await this.getAlbumDetail(albumId)

    // 既存の画像IDのSetを作成して重複を避ける
    const existingImageIds = new Set(currentAlbum.images)

    // 新しい画像IDを追加（重複を除外）
    const newImageIds = imageIds.filter((id) => !existingImageIds.has(id))
    const allImageIds = [...currentAlbum.images, ...newImageIds]

    // 統合された画像リストでアルバムを更新
    const updateData: UpdateAlbumRequest = {
      title: currentAlbum.title,
      description: currentAlbum.description,
      images: allImageIds,
    }

    return this.updateAlbum(albumId, updateData)
  }

  // アルバム削除（DELETE /albums/{id}）
  async deleteAlbum(albumId: string): Promise<void> {
    return apiClient.delete<void>(`/albums/${albumId}`)
  }

  // 特定のユーザーのアルバム取得（convenience method）
  async getAlbumsByCreator(creatorId: string): Promise<AlbumItem[]> {
    return this.getAlbums({ creator_id: creatorId })
  }
}

export const albumService = new AlbumService()
