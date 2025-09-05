// 共通の型定義（OpenAPI仕様準拠）

// エラーレスポンス
export interface ApiError {
  error: string
  message?: string
}

// 共通のstate型
export interface LoadingState {
  loading: boolean
  error: string | null
}
