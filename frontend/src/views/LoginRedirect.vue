<template>
  <div class="login-redirect">
    <p v-if="state === 'checking'">Checking session...</p>
    <p v-else-if="state === 'redirecting-login'">Redirecting to login...</p>
    <p v-else>Redirecting...</p>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'

type State = 'checking' | 'redirecting-login' | 'redirecting-next'
const state = ref<State>('checking')

function getNextPath(): string {
  const params = new URLSearchParams(window.location.search)
  const next = params.get('next')
  return next && next.startsWith('/') ? next : '/'
}

onMounted(async () => {
  const next = getNextPath()
  try {
    const resp = await fetch('/api/auth/me', { credentials: 'include' })
    if (resp.ok) {
      // 既にログイン済みなら next に進む
      state.value = 'redirecting-next'
      window.location.replace(next)
      return
    }
  } catch {
    // 無視してログインへ
  }

  // 未ログイン: OAuth フロー開始。callback は /login ではなく next にする
  state.value = 'redirecting-login'
  const callback = next
  window.location.href = `/api/auth/request?callback=${encodeURIComponent(callback)}`
})
</script>

<style scoped>
.login-redirect {
  padding: 24px;
}
</style>
