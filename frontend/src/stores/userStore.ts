import { defineStore } from 'pinia'
import type { Me } from '@/types/user'
import type { User } from '@/types'

// 注意: 現在のOpenAPI仕様にはユーザー関連のエンドポイントが定義されていません
// このStoreは将来の拡張のために保持していますが、現在は使用されていません

export const useUserStore = defineStore('user', {
  state: () => ({
    me: null as Me | null,
    loading: false as boolean,
  }),
  actions: {
    async fetchMe() {
      this.loading = true
      try {
        const resp = await fetch('/api/auth/me', { credentials: 'include' })
        if (resp.ok) {
          this.me = await resp.json()
        } else if (resp.status === 401) {
          this.me = null
        }
      } finally {
        this.loading = false
      }
    },
  },
})
