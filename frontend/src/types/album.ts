// アルバム関連の型定義（OpenAPI仕様準拠）

export interface Album {
  id: string // UUID
  title: string
  description: string
  creator: string // UUID
  images: string[] // 画像UUIDの配列
  created_at: string // ISO 8601 format
  updated_at: string // ISO 8601 format
}

export interface AlbumItem {
  id: string // UUID
  title: string
  creator: string // UUID
}

export interface CreateAlbumRequest {
  title: string
  description: string
  images?: string[] // 画像UUIDの配列
}

export interface UpdateAlbumRequest {
  title?: string
  description?: string
  images?: string[] // 画像UUIDの配列
}

// クエリパラメータ
export interface GetAlbumsParams {
  creator_id?: string // UUID
  before_date?: string // ISO 8601 date-time
  after_date?: string // ISO 8601 date-time
  limit?: number // 1-100, default: 20
  offset?: number // minimum: 0, default: 0
}
