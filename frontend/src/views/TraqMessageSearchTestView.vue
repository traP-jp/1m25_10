<template>
  <div class="traq-image-search-test">
    <div class="container">
      <h1>traQ Image Search Test</h1>

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
        <h2>Image Search (via traQ message search)</h2>
        <div class="form-grid">
          <label>
            word
            <input v-model="form.word" placeholder="例: テスト" />
          </label>
          <label>
            after (RFC3339)
            <input v-model="form.after" placeholder="2025-01-01T00:00:00Z" />
          </label>
          <label>
            before (RFC3339)
            <input v-model="form.before" placeholder="2025-12-31T23:59:59Z" />
          </label>
          <label>
            in (channel uuid)
            <input v-model="form.in" />
          </label>
          <label>
            to (comma separated)
            <input v-model="toText" placeholder="u1,u2" />
          </label>
          <label>
            from (comma separated)
            <input v-model="fromText" placeholder="f1,f2" />
          </label>
          <label>
            citation (uuid)
            <input v-model="form.citation" />
          </label>
          <label class="checkbox"><input type="checkbox" v-model="form.bot" /> bot</label>
          <label class="checkbox"><input type="checkbox" v-model="form.hasURL" /> hasURL</label>
          <label class="checkbox"
            ><input type="checkbox" v-model="form.hasAttachments" /> hasAttachments</label
          >
          <label class="checkbox"><input type="checkbox" v-model="form.hasImage" /> hasImage</label>
          <label class="checkbox"><input type="checkbox" v-model="form.hasVideo" /> hasVideo</label>
          <label class="checkbox"><input type="checkbox" v-model="form.hasAudio" /> hasAudio</label>
          <label>
            limit
            <input type="number" v-model.number="form.limit" />
          </label>
          <label>
            offset
            <input type="number" v-model.number="form.offset" />
          </label>
          <label>
            sort
            <select v-model="form.sort">
              <option value="">(none)</option>
              <option value="createdAt">createdAt</option>
              <option value="-createdAt">-createdAt</option>
              <option value="updatedAt">updatedAt</option>
              <option value="-updatedAt">-updatedAt</option>
            </select>
          </label>
        </div>

        <div class="actions">
          <button @click="setExample" :disabled="loading">[Set Example]</button>
          <button @click="runSearch" :disabled="loading">[Message Search]</button>
          <button @click="runExtract" :disabled="loading">[Extract Image UUIDs]</button>
        </div>

        <div class="result">
          <div class="search-result">
            <h3>Search Result:</h3>
            <div v-if="loading">Loading...</div>
            <div v-else-if="error" class="error">Error: {{ error }}</div>
            <pre v-else class="json">{{ prettyJson }}</pre>
          </div>
          <div class="image-extract">
            <h3>Extracted Image UUID Results:</h3>
            <div v-if="imagesLoading">Loading images...</div>
            <div v-else-if="imagesError" class="error">Error: {{ imagesError }}</div>
            <pre v-else class="json">{{ JSON.stringify(imagesResult, null, 2) }}</pre>
          </div>
        </div>
      </div>

      <div class="test-section">
        <h2>API Endpoints</h2>
        <div class="endpoint-info">
          <h3>Relevant Endpoints:</h3>
          <ul>
            <li>
              <code>GET /api/v1/traq/messages</code> - traQ message search proxy (all traQ query
              params supported)
            </li>
            <li>
              <code>GET /api/v1/images</code> - search with hasImage=true and extract image UUIDs
              from message content
            </li>
          </ul>
          <h3>Usage in Frontend:</h3>
          <pre><code>// search
fetch('/api/v1/traq/messages?word=テスト', { credentials: 'same-origin' })
// extract images
fetch('/api/v1/images?word=テスト', { credentials: 'same-origin' })
</code></pre>
        </div>
      </div>

      <div class="test-section">
        <h2>Features Tested</h2>
        <ul class="feature-list">
          <li>✅ Authentication check (Cookie-based)</li>
          <li>✅ Proxying traQ message search</li>
          <li>✅ hasImage=true enforced and image extraction</li>
          <li>✅ Extraction of multiple image URLs per message</li>
          <li>✅ Returns totalHits and UUID array</li>
        </ul>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref, computed, onMounted } from 'vue'
import { useUserStore } from '@/stores/userStore'
import { searchTraqMessages, type TraqMessageSearchParams } from '@/services/traqService'

const userStore = useUserStore()

const form = reactive<TraqMessageSearchParams>({
  word: '',
  after: '',
  before: '',
  in: '',
  to: [],
  from: [],
  citation: '',
  bot: undefined,
  hasURL: undefined,
  hasAttachments: undefined,
  hasImage: undefined,
  hasVideo: undefined,
  hasAudio: undefined,
  limit: undefined,
  offset: undefined,
  sort: '',
})

const toText = ref('')
const fromText = ref('')
const loading = ref(false)
const error = ref<string | null>(null)
const result = ref<unknown>(null)

const prettyJson = computed(() => JSON.stringify(result.value, null, 2))

const isLoggedIn = computed(() => userStore.me !== null)
const userInfo = computed(() => userStore.me)

function setExample() {
  form.word = 'デバッグ用の写真。'
}

async function runSearch() {
  loading.value = true
  error.value = null
  result.value = null

  // 文字列→配列
  form.to = toText.value
    .split(',')
    .map((s) => s.trim())
    .filter((s) => s.length > 0)
  form.from = fromText.value
    .split(',')
    .map((s) => s.trim())
    .filter((s) => s.length > 0)

  try {
    result.value = await searchTraqMessages(form)
  } catch (e: unknown) {
    if (typeof e === 'object' && e !== null && 'message' in e) {
      const m = (e as { message?: unknown }).message
      error.value = typeof m === 'string' ? m : JSON.stringify(e)
    } else {
      error.value = String(e)
    }
  } finally {
    loading.value = false
  }
}

const imagesLoading = ref(false)
const imagesError = ref<string | null>(null)
const imagesResult = ref<unknown>(null)

async function runExtract() {
  imagesLoading.value = true
  imagesError.value = null
  imagesResult.value = null

  // build params same as runSearch
  form.to = toText.value
    .split(',')
    .map((s) => s.trim())
    .filter((s) => s.length > 0)
  form.from = fromText.value
    .split(',')
    .map((s) => s.trim())
    .filter((s) => s.length > 0)

  const q: Record<string, unknown> = { ...form }
  // remove undefined values
  for (const k of Object.keys(q)) {
    if (q[k] === undefined || q[k] === null || q[k] === '') delete q[k]
  }

  try {
    const query = new URLSearchParams()
    for (const [k, v] of Object.entries(q)) {
      if (Array.isArray(v)) {
        for (const e of v) query.append(k, String(e))
      } else {
        query.append(k, String(v))
      }
    }
    const url = `/api/v1/images?${query.toString()}`
    const res = await fetch(url, { credentials: 'same-origin' })
    if (!res.ok) {
      const t = await res.text()
      imagesError.value = t || `HTTP ${res.status}`
    } else {
      imagesResult.value = await res.json()
    }
  } catch (e: unknown) {
    imagesError.value = String(e)
  } finally {
    imagesLoading.value = false
  }
}

onMounted(async () => {
  if (userStore.me === null) {
    await userStore.fetchMe()
  }
})
</script>

<style scoped>
.traq-image-search-test {
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

.form-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(220px, 1fr));
  gap: 12px;
}
.checkbox {
  display: flex;
  align-items: center;
  gap: 6px;
}
.actions {
  margin: 16px 0;
  display: flex;
  gap: 8px;
}
.result .json {
  background: #fff;
  padding: 12px;
  border-radius: 6px;
  overflow: auto;
}
.error {
  color: #b00020;
}

/* Styles borrowed from TraqFileTestView to unify the design */
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
  flex-wrap: wrap;
}

.input-section label {
  display: flex;
  flex-direction: column;
  min-width: 220px;
  flex: 1 1 220px;
}

.uuid-input,
input[type='text'],
input[type='number'],
select {
  padding: 8px;
  border: 1px solid #ddd;
  border-radius: 4px;
}

.test-button,
.actions button {
  padding: 8px 12px;
  background: #007bff;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

.sample-list {
  padding: 0;
  list-style: none;
}

.sample-list li {
  padding: 6px 0;
}

.embed-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  gap: 12px;
}

.note {
  margin-top: 8px;
  font-size: 13px;
  color: #666;
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
