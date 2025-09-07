export type TraqMessageSearchParams = {
  word?: string
  after?: string
  before?: string
  in?: string
  to?: string[]
  from?: string[]
  citation?: string
  stampId?: string
  bot?: boolean
  hasURL?: boolean
  hasAttachments?: boolean
  hasImage?: boolean
  hasVideo?: boolean
  hasAudio?: boolean
  limit?: number
  offset?: number
  sort?: string
}

function buildQuery(params: Record<string, unknown>): string {
  const u = new URLSearchParams()
  for (const [k, v] of Object.entries(params)) {
    if (v == null) continue
    if (Array.isArray(v)) {
      for (const e of v) {
        u.append(k, String(e))
      }
    } else {
      u.append(k, String(v))
    }
  }
  return u.toString()
}

// テスト用途では明示的に /api/v1/traq/messages を叩く
export async function searchTraqMessages(params: TraqMessageSearchParams): Promise<unknown> {
  const query = buildQuery(params as Record<string, unknown>)
  const url = `/api/v1/traq/messages${query ? `?${query}` : ''}`

  // 既存の apiClient を経由する必要があれば差し替え可能だが、ここでは直接 fetch
  const res = await fetch(url, { credentials: 'same-origin' })
  if (!res.ok) {
    const text = await res.text()
    throw new Error(text || `HTTP ${res.status}`)
  }
  return res.json()
}
