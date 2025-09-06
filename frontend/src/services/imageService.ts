// 画像関連のAPIサービス（OpenAPI仕様準拠）

import { apiClient } from './apiClient'
import type { Image, ImageDetail, GetImagesParams } from '@/types'

export class ImageService {
  // 全画像取得（GET /images）
  async getImages(searchQuery?: string): Promise<Image[]> {
    const queryParams: GetImagesParams = {}

    if (searchQuery && searchQuery.trim()) {
      queryParams.word = searchQuery.trim()
    }

    return apiClient.get<Image[]>('/images', queryParams as Record<string, unknown>)
  }

  // 特定の画像詳細取得（GET /images/{id}）
  async getImageDetail(imageId: string): Promise<ImageDetail> {
    return apiClient.get<ImageDetail>(`/images/${imageId}`)
  }
}

export const imageService = new ImageService()
