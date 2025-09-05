// axiosベースのAPIクライアント

import axios, { type AxiosInstance, type AxiosResponse, AxiosError } from 'axios'
import type { ApiError } from '@/types'

class ApiClient {
  private client: AxiosInstance

  constructor(baseURL: string = 'http://localhost:3001') {
    this.client = axios.create({
      baseURL,
      timeout: 10000,
      headers: {
        'Content-Type': 'application/json',
      },
    })

    // リクエストインターセプター
    this.client.interceptors.request.use(
      (config) => {
        // 必要に応じて認証トークンなどを追加
        return config
      },
      (error) => {
        return Promise.reject(error)
      },
    )

    // レスポンスインターセプター
    this.client.interceptors.response.use(
      (response: AxiosResponse) => {
        return response
      },
      (error: AxiosError<ApiError>) => {
        if (error.response) {
          // サーバーからエラーレスポンスが返された場合
          const apiError: ApiError = {
            error: error.response.data?.error || `HTTP ${error.response.status}`,
            message: error.response.data?.message || error.message,
          }
          return Promise.reject(apiError)
        } else if (error.request) {
          // リクエストが送信されたがレスポンスがない場合
          const apiError: ApiError = {
            error: 'Network Error',
            message: 'No response from server',
          }
          return Promise.reject(apiError)
        } else {
          // リクエスト設定でエラーが発生した場合
          const apiError: ApiError = {
            error: 'Request Error',
            message: error.message,
          }
          return Promise.reject(apiError)
        }
      },
    )
  }

  async get<T>(endpoint: string, params?: Record<string, unknown>): Promise<T> {
    const response = await this.client.get<T>(endpoint, { params })
    return response.data
  }

  async post<T>(endpoint: string, data?: unknown): Promise<T> {
    const response = await this.client.post<T>(endpoint, data)
    return response.data
  }

  async put<T>(endpoint: string, data?: unknown): Promise<T> {
    const response = await this.client.put<T>(endpoint, data)
    return response.data
  }

  async delete<T>(endpoint: string): Promise<T> {
    const response = await this.client.delete<T>(endpoint)
    return response.data
  }

  async upload<T>(endpoint: string, formData: FormData): Promise<T> {
    const response = await this.client.post<T>(endpoint, formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    })
    return response.data
  }
}

export const apiClient = new ApiClient()
