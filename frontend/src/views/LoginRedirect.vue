<template>
  <div class="login-redirect">
    <p v-if="state === 'checking'">Checking session...</p>
    <p v-else-if="state === 'redirecting-login'">Redirecting to login...</p>
    <p v-else>Redirecting...</p>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'

type State = 'checking' | 'redirecting-login' | 'redirecting-next'
const state = ref<State>('checking')
const router = useRouter()

function isSafeInternalPath(p: string | null | undefined): p is string {
  if (!p) return false
  if (p[0] !== '/') return false
  if (p.length > 1 && (p[1] === '/' || p[1] === '\\')) return false
  if (p.includes('\\')) return false
  for (let i = 0; i < p.length; i++) {
    const code = p.charCodeAt(i)
    if (code < 0x20 || code === 0x7f) return false
  }
  try {
    const u = new URL(p, window.location.origin)
    if (u.origin !== window.location.origin) return false
  } catch {
    return false
  }
  return true
}

function getNextPath(): string {
  const params = new URLSearchParams(window.location.search)
  const next = params.get('next')
  return isSafeInternalPath(next) ? next : '/'
}

onMounted(async () => {
  const next = getNextPath()
  try {
    const resp = await fetch('/api/auth/me', { credentials: 'include' })
    if (resp.ok) {
      // 既にログイン済みなら next に進む
      state.value = 'redirecting-next'
      await router.replace(next)
      return
    }
  } catch {
    // 無視してログインへ
  }

  // 未ログイン: OAuth フロー開始。callback は /login ではなく next にする
  state.value = 'redirecting-login'
  const callback = getNextPath()
  window.location.href = `/api/auth/request?callback=${encodeURIComponent(callback)}`
})
</script>

<style scoped>
.login-redirect {
  padding: 24px;
}
</style>
