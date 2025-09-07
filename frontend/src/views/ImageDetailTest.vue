<template>
  <div class="image-detail-test">
    <div class="container">
      <h1>Image Post Detail Test</h1>

      <div class="test-section">
        <h2>Authentication Status</h2>
        <p v-if="isLoggedIn" class="status-logged-in">
          ✅ Logged in as: {{ userInfo?.name || 'Unknown' }}
        </p>
        <p v-else class="status-not-logged-in">
          ❌ Not logged in. Please <router-link to="/account">login</router-link> first.
        </p>
      </div>

      <div v-if="isLoggedIn" class="test-section">
        <h2>Get Post Detail by File UUID</h2>
        <p>
          <code>GET /api/v1/images/{uuid}</code> により、
          <code>https://q.trap.jp/files/{uuid}</code> を語句検索して、
          投稿者・本文・スタンプなどを含む投稿詳細（最も古い1件）を返します。
        </p>

        <div class="input-section">
          <input
            v-model="uuid"
            type="text"
            placeholder="Enter image/file UUID"
            class="uuid-input"
          />
          <button @click="run" :disabled="!uuid || loading" class="test-button">Run</button>
          <button @click="setTemplate" :disabled="loading" class="test-button">Use Template</button>
        </div>

        <div class="result-section">
          <h3>Post Detail (Oldest)</h3>
          <div v-if="loading">Loading...</div>
          <div v-else-if="error" class="error">Error: {{ error }}</div>
          <pre v-else class="json">{{ prettyJson }}</pre>
        </div>
      </div>

      <div class="test-section">
        <h2>API Endpoint</h2>
        <div class="endpoint-info">
          <ul>
            <li>
              <code>GET /api/v1/images/{uuid}</code> - ファイルUUIDから投稿詳細（最も古い1件）を返す
            </li>
          </ul>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useUserStore } from '@/stores/userStore'

const userStore = useUserStore()

const TEMPLATE_UUID = '01991edf-51bb-72ac-b380-a607a9d2b474'
const uuid = ref(TEMPLATE_UUID)
const loading = ref(false)
const error = ref<string | null>(null)
const result = ref<unknown>(null)

const isLoggedIn = computed(() => userStore.me !== null)
const userInfo = computed(() => userStore.me)
const prettyJson = computed(() => JSON.stringify(result.value, null, 2))

function setTemplate() {
  uuid.value = TEMPLATE_UUID
}

async function run() {
  loading.value = true
  error.value = null
  result.value = null
  try {
    const res = await fetch(`/api/v1/images/${encodeURIComponent(uuid.value)}`, {
      credentials: 'same-origin',
    })
    if (!res.ok) {
      const t = await res.text()
      throw new Error(t || `HTTP ${res.status}`)
    }
    result.value = await res.json()
  } catch (e: unknown) {
    error.value = e instanceof Error ? e.message : String(e)
  } finally {
    loading.value = false
  }
}

onMounted(async () => {
  if (userStore.me === null) {
    await userStore.fetchMe()
  }
})
</script>

<style scoped>
.image-detail-test {
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
.input-section {
  display: flex;
  gap: 8px;
  margin: 8px 0;
}
.uuid-input {
  flex: 1;
  padding: 8px;
  border: 1px solid #ddd;
  border-radius: 4px;
}
.test-button,
.copy-button {
  padding: 8px 12px;
  background: #007bff;
  color: white;
  border: none;
  border-radius: 4px;
}
.json {
  background: #fff;
  padding: 12px;
  border-radius: 6px;
  overflow: auto;
}
.error {
  color: #b00020;
}
.sample-list {
  padding: 0;
  list-style: none;
}
.sample-list li {
  padding: 6px 0;
}
@media (max-width: 768px) {
  .container {
    padding: 0 10px;
  }
}
</style>
