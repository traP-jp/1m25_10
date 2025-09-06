<template>
  <div class="oauth-debug-test">
    <div class="container">
      <h1>OAuth Debug Test</h1>

      <div class="test-section">
        <h2>Authentication Status</h2>
        <div v-if="loading">Loading...</div>
        <div v-else>
          <div v-if="me">
            <p>
              Logged in as: <strong>{{ me.displayName }}</strong> (@{{ me.name }})
            </p>
            <button @click="logout" class="test-button">Logout</button>
          </div>
          <div v-else>
            <p>You are not logged in.</p>
            <button @click="toLogin" class="test-button">Login with traQ</button>
          </div>
          <p v-if="error" class="error">{{ error }}</p>
        </div>
      </div>

      <div class="test-section">
        <h2>API Endpoints</h2>
        <div class="endpoint-info">
          <h3>Authentication Endpoints:</h3>
          <ul>
            <li><code>GET /api/auth/me</code> - Get current user information</li>
            <li><code>POST /api/auth/logout</code> - Logout current user</li>
            <li><code>GET /api/auth/request</code> - Initiate OAuth login flow</li>
          </ul>
          <h3>Usage:</h3>
          <ul>
            <li>Login button redirects to traQ OAuth authorization</li>
            <li>Logout clears session and redirects to login</li>
            <li>User information is stored in cookies</li>
          </ul>
        </div>
      </div>
    </div>
  </div>
</template>

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

<style scoped>
.oauth-debug-test {
  min-height: 100vh;
  background: #f8f9fa;
  padding: 20px;
}

.container {
  max-width: 1200px;
  margin: 0 auto;
}

h1 {
  text-align: center;
  color: #333;
  margin-bottom: 40px;
  font-size: 32px;
}

.test-section {
  margin-bottom: 50px;
}

.test-section h2 {
  color: #555;
  margin-bottom: 20px;
  font-size: 24px;
  border-bottom: 2px solid #007bff;
  padding-bottom: 10px;
}

.test-button {
  padding: 8px 12px;
  background: #007bff;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

.test-button:hover {
  background: #0056b3;
}

.endpoint-info {
  background: #fff;
  padding: 20px;
  border-radius: 8px;
  border: 1px solid #ddd;
}

.endpoint-info h3 {
  color: #333;
  margin-top: 0;
  margin-bottom: 10px;
}

.endpoint-info ul {
  margin: 0;
  padding-left: 20px;
}

.endpoint-info li {
  margin-bottom: 8px;
}

.endpoint-info code {
  background: #f4f4f4;
  padding: 2px 4px;
  border-radius: 3px;
  font-family: monospace;
}

.error {
  color: #c00;
  margin-top: 10px;
}

@media (max-width: 768px) {
  .container {
    padding: 0 10px;
  }
}
</style>
