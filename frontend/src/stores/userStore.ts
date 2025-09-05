import { defineStore } from 'pinia'
import type { Me } from '@/types/user'

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
