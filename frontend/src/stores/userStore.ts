import { defineStore } from 'pinia'
import type { User } from '@/types'

// 注意: 現在のOpenAPI仕様にはユーザー関連のエンドポイントが定義されていません
// このStoreは将来の拡張のために保持していますが、現在は使用されていません

export const useUserStore = defineStore('user', {
  state: () => ({
    users: [] as User[],
    currentUser: null as User | null,
    authenticatedUser: null as User | null,
    loading: false,
    error: null as string | null,
  }),

  getters: {
    // ユーザー数
    userCount: (state) => state.users.length,

    // 名前でユーザーを検索
    getUserByName: (state) => {
      return (name: string) =>
        state.users.find((user) => user.name.toLowerCase().includes(name.toLowerCase()))
    },

    // メールでユーザーを検索
    getUserByEmail: (state) => {
      return (email: string) => state.users.find((user) => user.email === email)
    },

    // 認証済みかどうか
    isAuthenticated: (state) => state.authenticatedUser !== null,
  },

  actions: {
    // 注意: 以下のメソッドは現在のOpenAPI仕様では使用できません
    // 将来ユーザー関連のエンドポイントが追加された際に有効化してください

    // 全ユーザー取得
    async fetchUsers() {
      this.loading = true
      this.error = null

      try {
        // TODO: OpenAPI仕様にユーザーエンドポイントが追加されたら有効化
        // const users = await userService.fetchUsers()
        // this.users = users
        console.warn('User endpoints are not available in current OpenAPI specification')
      } catch (error) {
        this.error = error instanceof Error ? error.message : 'Unknown error occurred'
        console.error('Error fetching users:', error)
      } finally {
        this.loading = false
      }
    },

    // エラーをクリア
    clearError() {
      this.error = null
    },

    // 現在のユーザーをクリア
    clearCurrentUser() {
      this.currentUser = null
    },

    // 認証ユーザーをクリア
    clearAuthenticatedUser() {
      this.authenticatedUser = null
    },
  },
})
