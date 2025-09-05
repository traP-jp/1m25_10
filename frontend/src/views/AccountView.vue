<script setup lang="ts">
import { onMounted, ref } from 'vue'
import type { Me } from '@/types/user'

const me = ref<Me | null>(null)
const loading = ref(true)
const error = ref<string | null>(null)

async function fetchMe() {
  loading.value = true
  error.value = null
  try {
    const resp = await fetch('/api/auth/me', { credentials: 'include' })
    if (resp.ok) {
      me.value = await resp.json()
    } else if (resp.status === 401) {
      me.value = null
    } else {
      error.value = `failed: ${resp.status}`
    }
  } catch (e) {
    error.value = String(e)
  } finally {
    loading.value = false
  }
}

function toLogin() {
  // 現在のパスをcallbackに付与
  const cb = location.pathname + location.search + location.hash
  location.href = `/api/auth/request?callback=${encodeURIComponent(cb)}`
}

async function logout() {
  await fetch('/api/auth/logout', { method: 'POST', credentials: 'include' })
  await fetchMe()
}

onMounted(fetchMe)
</script>

<template>
  <div class="container">
    <div style="display: flex; align-items: baseline; gap: 16px;">
      <h1>OAuth</h1>
      <p style="color: darkgray; margin: 0;">For debugging. Will require login at startup in the future.</p>
    </div>

    <div v-if="loading">Loading...</div>
    <div v-else>
      <div v-if="me">
        <p>
          Logged in as: <strong>{{ me.name }}</strong>
        </p>
        <button @click="logout">Logout</button>
      </div>
      <div v-else>
        <p>You are not logged in.</p>
        <button @click="toLogin">Login with traQ</button>
      </div>
      <p v-if="error" class="error">{{ error }}</p>
    </div>
  </div>
</template>

<style scoped>
.container {
  padding: 16px;
}
.error {
  color: #c00;
}
button {
  padding: 8px 12px;
  border: 1px solid #ccc;
  background: #fff;
  border-radius: 6px;
  cursor: pointer;
}
button:hover {
  background: #f5f5f5;
}
</style>
