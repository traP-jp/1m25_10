export interface Post {
  id: string // UUID
  content: string
}

// 画像の基本情報（IDのみ）
export interface Image {
  id: string // UUID
}

// 画像の詳細情報（クリック時に取得）
export interface ImageDetail {
  id: string // UUID
  creator: string // UUID
  post: Post
}

// GET /images のレスポンス型
export interface GetImagesResponse {
  hits: string[] // UUID[]
}

// 画像検索用のクエリパラメータ
export interface GetImagesParams {
  word?: string // 検索キーワード
  limit?: number // 取得件数 (デフォルト: 20)
  offset?: number // オフセット (デフォルト: 0)
}
