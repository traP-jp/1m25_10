// サイドバーの共通設定（読み取り専用）
export const SIDEBAR_CONFIG = {
  MIN_WIDTH: 100,
  MAX_WIDTH: 500,
  DEFAULT_WIDTH: 300,
  COMPACT_THRESHOLD: 250,
  STORAGE_KEY: 'sidebar-config',
  // テキスト測定用の設定
  ICON_WIDTH: 40, // アイコンの幅
  ICON_GAP: 12, // アイコンとテキストの間隔
  ITEM_HORIZONTAL_PADDING: 32, // SidebarItemの左右パディング合計 (16px * 2)
  EXTRA_PADDING: 16, // 余裕をもたせるためのパディング
} as const

export type SidebarConfig = typeof SIDEBAR_CONFIG
