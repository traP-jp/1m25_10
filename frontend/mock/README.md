# Frontend API Mock

フロントエンド開発用のAPIモックサーバーです。JSON Serverを使用して、traPhoto(仮) APIのモックを提供します。以下AI生成。

## 🚀 使用方法

### モックサーバーの起動

```bash
# モックサーバーのみ起動
npm run mock

# フロントエンドとモックサーバーを同時に起動
npm run dev:mock
```

### エンドポイント

モックサーバーは `http://localhost:3001` で動作します。

#### 画像関連
- `GET /images` - 全ての画像を取得
- `GET /images/{id}` - IDで画像詳細を取得

#### アルバム関連
- `GET /albums` - 全てのアルバムを取得（フィルタリング対応）
- `GET /albums/{id}` - IDでアルバム詳細を取得
- `POST /albums` - 新しいアルバムを作成
- `PUT /albums/{id}` - アルバムを更新
- `DELETE /albums/{id}` - アルバムを削除

## 📖 API仕様書

OpenAPI 3.0仕様書: [`openapi.yaml`](/openapi.yaml)

オンラインビューアーで表示：
- [Swagger Editor](https://editor.swagger.io/) にYAMLファイルをコピー&ペースト
- [Redoc](https://redocly.github.io/redoc/) でプレビュー

## 📁 ファイル構成

```
openapi.yaml              # OpenAPI仕様書
frontend/
├── mock/
│   ├── db.json           # モックデータ
│   ├── server.js         # カスタムサーバー設定
│   ├── assets/           # 静的ファイル（画像など）
│   └── README.md         # このファイル
└── src/                  # フロントエンドソースコード
```

## 🔧 設定

### 環境変数

`.env.development` と `.env.production` で環境ごとの設定を管理：

```bash
# 開発環境
VITE_API_BASE_URL=http://localhost:3001
VITE_MOCK_ENABLED=true

# 本番環境
VITE_API_BASE_URL=https://api.your-domain.com
VITE_MOCK_ENABLED=false
```

### クエリパラメータ（アルバム検索）

`GET /albums` エンドポイントでは以下のクエリパラメータをサポート：

- `creator_id` - 作成者のUUIDでフィルタ
- `before_date` - この日付より前に作成されたアルバムでフィルタ
- `after_date` - この日付より後に作成されたアルバムでフィルタ
- `limit` - 返すアルバムの最大数（最大100、デフォルト20）
- `offset` - スキップするアルバム数（デフォルト0）

例：
```
GET /albums?creator_id=550e8400-e29b-41d4-a716-446655441001&limit=10
```

## 📝 データ形式

### Image
```typescript
interface Image {
  id: string        // UUID
  creator: string   // UUID
  post: {
    id: string      // UUID
    content: string
  }
}
```

### Album
```typescript
interface Album {
  id: string          // UUID
  title: string
  description: string
  creator: string     // UUID
  images: string[]    // 画像UUIDの配列
  created_at: string  // ISO 8601形式
  updated_at: string  // ISO 8601形式
}
```

### AlbumItem (一覧表示用)
```typescript
interface AlbumItem {
  id: string        // UUID
  title: string
  creator: string   // UUID
}
```

### CreateAlbumRequest
```typescript
interface CreateAlbumRequest {
  title: string
  description: string
  images?: string[] // 画像UUIDの配列（オプション）
}
```

## 🔄 データの永続化

JSON Serverはファイルベースなので、APIを通じた変更は `mock/db.json` に保存されます。開発時にデータをリセットしたい場合は、このファイルを元の状態に戻してください。

### サンプルデータ

現在のdb.jsonには以下のデータが含まれています：
- **画像**: 7件のサンプル画像データ（すべてUUID形式のID）
- **アルバム**: 3件のサンプルアルバム（夏の思い出、開発記録、日常スナップ）

すべてのIDはUUID v4形式で統一されており、OpenAPI仕様に完全準拠しています。

## 🎯 次のステップ

1. **画像アップロード**: `mock/assets/` に画像ファイルを追加してファイルアップロード機能をテスト
2. **バックエンド連携**: 実際のバックエンドAPIが完成したら環境変数でエンドポイントを切り替え
3. **エラーハンドリング**: 特定の条件でエラーレスポンスを返すロジックを追加
4. **認証機能**: 認証が必要な場合は認証ヘッダーの検証を追加
5. **ページネーション**: アルバム一覧でページネーション機能をテスト

## 🐛 トラブルシューティング

### ポート競合
モックサーバーのポート（3001）が使用中の場合：
```bash
# server.jsでポートを変更
const PORT = process.env.PORT || 3002
```

### CORS エラー
ブラウザでCORSエラーが発生する場合は、`server.js` にCORS設定を追加してください。

### UUID形式エラー
APIエンドポイントのパスパラメータはUUID形式である必要があります：
- 正しい例: `/images/550e8400-e29b-41d4-a716-446655440001`
- 間違った例: `/images/1` や `/images/album-001`

### データ型エラー
- すべてのIDはUUID形式の文字列
- 日付は ISO 8601 形式（例：`2024-07-01T00:00:00Z`）
- アルバムの`images`フィールドは文字列の配列

## 📚 参考資料

- [JSON Server公式ドキュメント](https://github.com/typicode/json-server)
- [OpenAPI 3.0仕様](https://spec.openapis.org/oas/v3.0.3)
- [UUID v4仕様](https://tools.ietf.org/html/rfc4122)
