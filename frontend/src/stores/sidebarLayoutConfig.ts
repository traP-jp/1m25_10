// サイドバーの共通設定（読み取り専用）
export const SIDEBAR_CONFIG = {
  MIN_WIDTH: 80,
  MAX_WIDTH: 500,
  DEFAULT_WIDTH: 300,
  COMPACT_THRESHOLD: 150,
  STORAGE_KEY: 'sidebar-config',
} as const

export type SidebarConfig = typeof SIDEBAR_CONFIG
