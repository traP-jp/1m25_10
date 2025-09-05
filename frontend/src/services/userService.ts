// ユーザー関連のAPIサービス

import { apiClient } from './apiClient'
import type { User, CreateUserRequest, UpdateUserRequest } from '@/types'

export class UserService {
  // 全ユーザー取得
  async fetchUsers(): Promise<User[]> {
    return apiClient.get<User[]>('/users')
  }

  // 特定のユーザー取得
  async fetchUser(userId: string): Promise<User> {
    return apiClient.get<User>(`/users/${userId}`)
  }

  // ユーザー作成
  async createUser(userData: CreateUserRequest): Promise<User> {
    return apiClient.post<User>('/users', userData)
  }

  // ユーザー更新
  async updateUser(userId: string, userData: UpdateUserRequest): Promise<User> {
    return apiClient.put<User>(`/users/${userId}`, userData)
  }

  // ユーザー削除
  async deleteUser(userId: string): Promise<void> {
    return apiClient.delete<void>(`/users/${userId}`)
  }

  // ユーザー認証
  async authenticate(email: string, password: string): Promise<User> {
    return apiClient.post<User>('/auth/login', { email, password })
  }

  // ユーザー登録
  async register(userData: CreateUserRequest): Promise<User> {
    return apiClient.post<User>('/auth/register', userData)
  }

  // ログアウト
  async logout(): Promise<void> {
    return apiClient.post<void>('/auth/logout')
  }

  // 現在のユーザー情報取得
  async getCurrentUser(): Promise<User> {
    return apiClient.get<User>('/auth/me')
  }
}

export const userService = new UserService()
