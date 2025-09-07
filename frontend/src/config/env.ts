// 環境変数の集中管理
// 優先順位はViteの仕様に従い、OSで定義された VITE_* が .env.* より優先される

const rawApiBaseUrl = import.meta.env.VITE_API_BASE_URL as string | undefined
const rawMockEnabled = import.meta.env.VITE_MOCK_ENABLED as string | undefined
const rawImageBaseUrl = import.meta.env.VITE_IMAGE_BASE_URL as string | undefined
const rawImageFormat = import.meta.env.VITE_IMAGE_FORMAT as string | undefined
const rawImageSize = import.meta.env.VITE_IMAGE_SIZE as string | undefined
const rawRequireLogin = import.meta.env.VITE_REQUIRE_LOGIN as string | undefined

function toBoolean(v: string | undefined): boolean | undefined {
  if (v == null) return undefined
  const s = String(v).trim().toLowerCase()
  if (s === 'true') return true
  if (s === 'false') return false
  return undefined
}

export const env = {
  MODE: import.meta.env.MODE,
  BASE_URL: import.meta.env.BASE_URL,
  PROD: import.meta.env.PROD,
  DEV: import.meta.env.DEV,
  SSR: import.meta.env.SSR,
  VITE_API_BASE_URL: rawApiBaseUrl,
  VITE_MOCK_ENABLED: toBoolean(rawMockEnabled),
  VITE_IMAGE_BASE_URL: rawImageBaseUrl,
  VITE_IMAGE_FORMAT: rawImageFormat,
  VITE_IMAGE_SIZE: rawImageSize,
  VITE_REQUIRE_LOGIN: toBoolean(rawRequireLogin) ?? false,
}

// APIのベースURL決定（優先度）
// 1. VITE_API_BASE_URL が設定されていればそれを採用
// 2. VITE_MOCK_ENABLED が true ならモックサーバ(3001)
// 3. それ以外は相対 /api （Viteのproxy経由でバックエンドへ）
export function resolveApiBaseUrl(): string {
  if (env.VITE_API_BASE_URL && env.VITE_API_BASE_URL.trim() !== '') {
    return env.VITE_API_BASE_URL
  }
  if (env.VITE_MOCK_ENABLED === true) {
    return 'http://localhost:3001'
  }
  return '/api'
}

// 画像URL設定の解決
export interface ImageUrlConfig {
  baseUrl: string
  format: string
  size: string
  useTraqFiles: boolean // traQファイルエンドポイントを使用するか
}

export function resolveImageUrlConfig(): ImageUrlConfig {
  // 本番環境かどうかの判定
  const isProduction = env.PROD || (env.VITE_API_BASE_URL?.includes('trap.show') ?? false)

  // カスタムベースURLが設定されている場合
  if (env.VITE_IMAGE_BASE_URL && env.VITE_IMAGE_BASE_URL.trim() !== '') {
    return {
      baseUrl: env.VITE_IMAGE_BASE_URL,
      format: env.VITE_IMAGE_FORMAT || '',
      size: env.VITE_IMAGE_SIZE || '',
      useTraqFiles: isProduction, // 本番環境ならtraQファイル形式
    }
  }

  // 本番環境: traQ files API を使用
  if (isProduction) {
    return {
      baseUrl: 'https://1m25-10-dev.trap.show/api/v1',
      format: '',
      size: '',
      useTraqFiles: true,
    }
  }

  // 開発環境: Picsumを使用
  return {
    baseUrl: 'https://picsum.photos',
    format: env.VITE_IMAGE_FORMAT || '', // Picsumでは拡張子は不要
    size: env.VITE_IMAGE_SIZE || '400',
    useTraqFiles: false,
  }
}

// 画像URLを生成
export function generateImageUrl(imageId: string, config?: Partial<ImageUrlConfig>): string {
  const defaultConfig = resolveImageUrlConfig()
  const finalConfig = { ...defaultConfig, ...config }

  if (finalConfig.useTraqFiles) {
    // traQ files API形式: https://1m25-10-dev.trap.show/api/v1/traq/files/{id}/thumbnail
    return `${finalConfig.baseUrl}/traq/files/${imageId}/thumbnail`
  } else {
    // Picsum形式: https://picsum.photos/seed/{id}/400
    const formatSuffix = finalConfig.format ? `.${finalConfig.format}` : ''
    return `${finalConfig.baseUrl}/seed/${imageId}/${finalConfig.size}${formatSuffix}`
  }
}
