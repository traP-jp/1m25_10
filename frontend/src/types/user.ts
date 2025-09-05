// ユーザー関連の型定義
// 注意: 現在のOpenAPI仕様にはユーザー関連のエンドポイントが定義されていません
// これらの型は将来の拡張のために保持しています

export interface User {
  id: string
  name: string
  email: string
}

export interface CreateUserRequest {
  name: string
  email: string
}

export interface UpdateUserRequest {
  name?: string
  email?: string
}
