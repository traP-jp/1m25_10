// 環境変数の集中管理
// 優先順位はViteの仕様に従い、OSで定義された VITE_* が .env.* より優先される

const rawApiBaseUrl = import.meta.env.VITE_API_BASE_URL as string | undefined
const rawMockEnabled = import.meta.env.VITE_MOCK_ENABLED as string | undefined

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
