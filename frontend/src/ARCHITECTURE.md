# Frontend Architecture（OpenAPI仕様準拠）

このプロジェクトでは、OpenAPI仕様に基づいてAPIクライアントと型定義を実装し、axiosを使用してHTTP通信を行っています。

## 🔄 最新の更新（OpenAPI準拠）

- **axios**を使用したHTTPクライアントに変更
- OpenAPI仕様（`/openapi.yaml`）に基づいた型定義
- エラーハンドリングの改善（レスポンスインターセプター）
- 認証トークンサポート（リクエストインターセプター）
- ベースURL設定（デフォルト: `http://localhost:3001`）

## ディレクトリ構造

```
src/
├── types/          # 型定義（OpenAPI準拠）
│   ├── album.ts    # アルバム関連の型（Album, AlbumItem, CreateAlbumRequest など）
│   ├── user.ts     # ユーザー関連の型（将来の拡張用）
│   ├── image.ts    # 画像関連の型（Image, ImageDetail, Post）
│   ├── common.ts   # 共通の型（ApiError, LoadingState）
│   └── index.ts    # 型定義のエクスポート
├── services/       # API関連のサービス（axios使用）
│   ├── apiClient.ts    # axiosベースのAPIクライアント
│   ├── albumService.ts # アルバムAPI（GET, POST, PUT, DELETE）
│   ├── userService.ts  # ユーザーAPI（将来の拡張用）
│   ├── imageService.ts # 画像API（GET /images, GET /images/{id}）
│   └── index.ts        # サービスのエクスポート
└── stores/         # Pinia ストア（状態管理のみ）
    ├── albumStore.ts   # アルバム状態管理
    ├── userStore.ts    # ユーザー状態管理
    ├── imageStore.ts   # 画像状態管理
    └── counter.ts      # カウンター（サンプル）
```

## OpenAPI仕様対応

### アルバム関連エンドポイント

- `GET /albums` - アルバム一覧取得（フィルタリング対応）
- `POST /albums` - アルバム作成
- `GET /albums/{id}` - アルバム詳細取得
- `PUT /albums/{id}` - アルバム更新
- `DELETE /albums/{id}` - アルバム削除

### 画像関連エンドポイント

- `GET /images` - 画像一覧取得
- `GET /images/{id}` - 画像詳細取得

### 型定義の特徴

- UUID形式のID
- ISO 8601形式の日時
- 必須/オプショナルプロパティの明確な区別
- OpenAPIスキーマとの完全な互換性

## 技術的な改善点

### 1. axiosの活用

```typescript
// インターセプターでエラーハンドリング
this.client.interceptors.response.use(
  (response) => response,
  (error: AxiosError<ApiError>) => {
    // 統一されたエラー形式で返却
    const apiError: ApiError = {
      error: error.response?.data?.error || 'Unknown error',
      message: error.response?.data?.message || error.message,
    }
    return Promise.reject(apiError)
  },
)
```

### 2. 型安全性の向上

```typescript
// OpenAPI仕様に基づく正確な型定義
export interface Album {
  id: string // UUID
  title: string
  description: string
  creator: string // UUID
  images: string[] // 画像UUIDの配列
  created_at: string // ISO 8601 format
  updated_at: string // ISO 8601 format
}
```

### 3. 効率的なキャッシング

```typescript
// アルバム詳細をキャッシュして無駄な API 呼び出しを削減
state: () => ({
  albums: [] as AlbumItem[], // 一覧用
  albumDetails: {} as Record<string, Album>, // 詳細キャッシュ
})
```

## 使用例

```typescript
// 型定義のインポート
import type { Album, CreateAlbumRequest } from '@/types'

// サービスの直接使用
import { albumService } from '@/services'

// ストアの使用
import { useAlbumStore } from '@/stores/albumStore'

// コンポーネント内での使用
const albumStore = useAlbumStore()

// フィルタリング付きでアルバムを取得
await albumStore.fetchAlbums({
  creator_id: 'user-uuid',
  limit: 10,
})

// アルバム詳細を取得（キャッシュされる）
const album = await albumStore.fetchAlbum('album-uuid')
```

## 利点

1. **OpenAPI準拠**: APIスキーマと型定義が完全に同期
2. **axios活用**: より柔軟なHTTP通信とエラーハンドリング
3. **関心の分離**: 型定義、API操作、状態管理が明確に分離
4. **再利用性**: サービスはコンポーネントから直接使用可能
5. **保守性**: 変更の影響範囲が限定的
6. **型安全性**: TypeScriptの利点を最大限活用
7. **パフォーマンス**: 詳細データのキャッシング機能
8. **テスト容易性**: 各レイヤーが独立してテスト可能
