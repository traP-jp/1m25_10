<template>
  <div class="traq-file-test">
    <div class="container">
      <h1>traQ File Proxy Test</h1>

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
        <h2>File Display Test</h2>
        <p>Test the traQ file proxy endpoints by entering a file UUID:</p>

        <div class="input-section">
          <input
            v-model="fileUuid"
            type="text"
            placeholder="Enter traQ file UUID"
            class="uuid-input"
          />
          <button @click="testFile" :disabled="!fileUuid" class="test-button">Test File</button>
          <button @click="testThumbnail" :disabled="!fileUuid" class="test-button">
            Test Thumbnail
          </button>
        </div>

        <div v-if="currentTest" class="result-section">
          <h3>Test Result: {{ currentTest.type }}</h3>
          <div class="result-content">
            <!-- 常に <img> を描画してネットワークリクエストを発生させる -->
            <img
              v-if="currentTest.url"
              :key="currentTest.url"
              :src="currentTest.url"
              :alt="currentTest.type"
              class="test-image"
              @error="handleImageError"
              @load="handleImageLoad"
            />

            <div v-if="currentTest.loading" class="loading">Loading...</div>
            <div v-else-if="currentTest.error" class="error">Error: {{ currentTest.error }}</div>
            <div v-else class="success">
              <p class="image-info">Image loaded successfully! ({{ currentTest.type }})</p>
            </div>
          </div>
        </div>
      </div>

      <div class="test-section">
        <h2>Sample UUIDs for Testing</h2>
        <p>You can use these sample UUIDs to test the functionality:</p>
        <ul class="sample-list">
          <li v-for="sample in sampleUuids" :key="sample.id">
            <code>{{ sample.uuid }}</code> - {{ sample.description }}
            <button @click="fileUuid = sample.uuid" class="copy-button">Use This</button>
          </li>
        </ul>
      </div>

      <div class="test-section">
        <h2>Embed Test with ImageCard</h2>
        <p>
          Below is an embedded ImageCard using the new thumbnail endpoint. Replace the UUID with a
          valid one you have access to.
        </p>
        <div class="embed-grid">
          <ImageCard
            :url="sampleEmbedUrl"
            alt="Thumbnail via /api/v1/traq/files/{uuid}/thumbnail"
          />
        </div>
        <div class="input-section">
          <input
            v-model="embedUuid"
            type="text"
            placeholder="Enter UUID for ImageCard embed"
            class="uuid-input"
          />
        </div>
        <p class="note">
          Current URL: <code>{{ sampleEmbedUrl }}</code>
        </p>
      </div>

      <div class="test-section">
        <h2>API Endpoints</h2>
        <div class="endpoint-info">
          <h3>File Endpoints:</h3>
          <ul>
            <li><code>GET /api/v1/traq/files/{uuid}</code> - Get file content</li>
            <li><code>GET /api/v1/traq/files/{uuid}/thumbnail</code> - Get file thumbnail</li>
          </ul>
          <h3>Usage in Frontend:</h3>
          <pre><code>&lt;img src="/api/v1/traq/files/YOUR_FILE_UUID" /&gt;
&lt;img src="/api/v1/traq/files/YOUR_FILE_UUID/thumbnail" /&gt;</code></pre>
        </div>
      </div>

      <div class="test-section">
        <h2>Features Tested</h2>
        <ul class="feature-list">
          <li>✅ Authentication check (Cookie-based)</li>
          <li>✅ File proxy with Bearer token</li>
          <li>✅ Thumbnail proxy with Bearer token</li>
          <li>✅ Error handling (403, 404, etc.)</li>
          <li>✅ Content-Type header forwarding</li>
          <li>✅ Response streaming</li>
        </ul>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useUserStore } from '@/stores/userStore'
import ImageCard from '@/components/ImageCard.vue'

const userStore = useUserStore()

// Reactive data
const fileUuid = ref('')
const currentTest = ref<{
  type: string
  url: string
  loading: boolean
  error: string | null
} | null>(null)

// Sample UUIDs for testing (these are just examples)
const sampleUuids = [
  {
    id: 1,
    uuid: '0198cffa-4228-7afe-ac39-e0c88a0f0b81',
    description: 'Sample image file (KAOMOJI)',
  },
  {
    id: 2,
    uuid: '0198b7bc-e241-7a7d-b4a2-4f56a756f92e',
    description: 'Sample image file (HENKOU 3219)',
  },
]

// Computed properties
const isLoggedIn = computed(() => userStore.me !== null)
const userInfo = computed(() => userStore.me)

// ImageCard 埋め込み用
const embedUuid = ref('0198cffa-4228-7afe-ac39-e0c88a0f0b81')
const sampleEmbedUrl = computed(() => `/api/v1/traq/files/${embedUuid.value}/thumbnail`)

// Methods
const testFile = () => {
  if (!fileUuid.value) return

  currentTest.value = {
    type: 'File',
    url: `/api/v1/traq/files/${fileUuid.value}`,
    loading: true,
    error: null,
  }
}

const testThumbnail = () => {
  if (!fileUuid.value) return

  currentTest.value = {
    type: 'Thumbnail',
    url: `/api/v1/traq/files/${fileUuid.value}/thumbnail`,
    loading: true,
    error: null,
  }
}

const handleImageLoad = () => {
  if (currentTest.value) {
    currentTest.value.loading = false
  }
}

const handleImageError = () => {
  if (currentTest.value) {
    currentTest.value.loading = false
    currentTest.value.error = 'Failed to load image. Check the UUID and your permissions.'
  }
}

// Initialize user info on mount
onMounted(async () => {
  if (userStore.me === null) {
    await userStore.fetchMe()
  }
})
</script>

<style scoped>
.traq-file-test {
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

.test-button {
  padding: 8px 12px;
  background: #007bff;
  color: white;
  border: none;
  border-radius: 4px;
}

.sample-list {
  padding: 0;
  list-style: none;
}

.sample-list li {
  padding: 6px 0;
}

.test-image {
  max-width: 100%;
  border-radius: 4px;
}

@media (max-width: 768px) {
  .container {
    padding: 0 10px;
  }
}
</style>
