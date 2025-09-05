// 画像関連のAPIサービス（OpenAPI仕様準拠）

import { apiClient } from './apiClient'
import type { Image, ImageDetail } from '@/types'

export class ImageService {
  // 全画像取得（GET /images）
  async getImages(): Promise<Image[]> {
    return apiClient.get<Image[]>('/images')
  }

  // 特定の画像詳細取得（GET /images/{id}）
  async getImageDetail(imageId: string): Promise<ImageDetail> {
    return apiClient.get<ImageDetail>(`/images/${imageId}`)
  }
}

export const imageService = new ImageService()
