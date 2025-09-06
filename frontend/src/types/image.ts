// 画像関連の型定義（OpenAPI仕様準拠）

export interface Post {
  id: string // UUID
  content: string
}

export interface Image {
  id: string // UUID
  creator: string // UUID
  post: Post
}

export interface ImageDetail {
  id: string // UUID
  creator: string // UUID
  post: Post
}

// 画像検索用のクエリパラメータ
export interface GetImagesParams {
  word?: string // 検索キーワード
}
