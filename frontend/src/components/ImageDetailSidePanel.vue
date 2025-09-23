<template>
  <aside :class="['side-panel', { open: visible }]">
    <header class="header">
      <h3>投稿の詳細</h3>
      <button class="close" @click="$emit('close')">×</button>
    </header>
    <div class="body">
      <div v-if="!imageId" class="placeholder">画像をクリックすると投稿詳細が表示されます</div>
      <div v-else>
        <div v-if="loading" class="loading">読み込み中...</div>
        <div v-else-if="error" class="error">{{ error }}</div>
        <div v-else-if="message" class="content">
          <div class="author">
            <div class="avatar">
              <img
                v-if="authorIconUrl && !authorIconError"
                :src="authorIconUrl"
                :alt="authorDisplayName || authorName || 'author'"
                class="avatarImg"
                @error="authorIconError = true"
                @load="authorIconError = false"
              />
              <span v-else class="avatarText">{{ authorInitials }}</span>
            </div>
            <div class="info">
              <div class="name">{{ authorDisplayName || authorName || message.userId }}</div>
              <div class="date">{{ formatDate(message.createdAt) }}</div>
            </div>
          </div>
          <div class="text" v-html="formattedContent"></div>
          <div class="channel">#{{ message.channelId }}</div>
          <div class="actions">
            <a :href="traqLink" target="_blank" rel="noopener" class="link">traQで開く</a>
          </div>
        </div>
      </div>
    </div>
  </aside>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'

interface Props {
  visible: boolean
  imageId: string | null
}

const props = defineProps<Props>()
defineEmits<{ close: [] }>()

type Message = {
  id: string
  userId: string
  channelId: string
  content: string
  createdAt: string
}

const loading = ref(false)
const error = ref<string | null>(null)
const message = ref<Message | null>(null)
const authorName = ref<string>('')
const authorDisplayName = ref<string>('')
const authorIconError = ref<boolean>(false)

watch(
  () => [props.visible, props.imageId] as const,
  async ([vis, id], [, prevId]) => {
    if (!vis || !id || id === prevId) return
    await fetchDetail(id)
  },
)

async function fetchDetail(id: string) {
  loading.value = true
  error.value = null
  message.value = null
  try {
    const res = await fetch(`/api/v1/images/${encodeURIComponent(id)}`, {
      credentials: 'same-origin',
    })
    if (!res.ok) throw new Error(await res.text())
    message.value = (await res.json()) as Message
    // 併せて投稿者情報を取得
    await fetchAuthor(message.value.userId)
  } catch (e: unknown) {
    error.value = e instanceof Error ? e.message : String(e)
  } finally {
    loading.value = false
  }
}

async function fetchAuthor(userId: string) {
  authorName.value = ''
  authorDisplayName.value = ''
  authorIconError.value = false
  try {
    const res = await fetch(`/api/v1/traq/users/${encodeURIComponent(userId)}`, {
      credentials: 'same-origin',
    })
    if (!res.ok) return
    const data = (await res.json()) as { name?: string; displayName?: string }
    authorName.value = data?.name || ''
    authorDisplayName.value = data?.displayName || ''
  } catch {
    // noop: フォールバックでイニシャル表示
  }
}

function formatDate(iso: string): string {
  try {
    const d = new Date(iso)
    const y = d.getFullYear()
    const m = String(d.getMonth() + 1).padStart(2, '0')
    const day = String(d.getDate()).padStart(2, '0')
    const hh = String(d.getHours()).padStart(2, '0')
    const mm = String(d.getMinutes()).padStart(2, '0')
    return `${y}/${m}/${day} ${hh}:${mm}`
  } catch {
    return iso
  }
}

const formattedContent = computed(() => {
  if (!message.value) return ''
  let text = message.value.content || ''
  const url = props.imageId ? `https://q.trap.jp/files/${props.imageId}` : ''
  if (url) text = text.split(url).join('')
  text = text.replace(/\n+$/g, '')
  return text.replace(/\n/g, '<br>')
})

const traqLink = computed(() => (props.imageId ? `https://q.trap.jp/files/${props.imageId}` : '#'))

const authorIconUrl = computed(() =>
  authorName.value
    ? `https://q.trap.jp/api/v3/public/icon/${encodeURIComponent(authorName.value)}`
    : '',
)

const authorInitials = computed(() => {
  const base = authorDisplayName.value || authorName.value
  return base ? base.slice(0, 2).toUpperCase() : '👤'
})
</script>

<style scoped>
.side-panel {
  /* in-flow side column by default */
  position: relative;
  width: 100%;
  height: auto;
  background: #fff;
  border-left: 1px solid #e5e7eb;
  box-shadow: none;
  transform: none;
  transition: none;
  display: flex;
  flex-direction: column;
  z-index: 1;
}
.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 14px;
  border-bottom: 1px solid #f0f0f0;
}
.close {
  background: none;
  border: none;
  font-size: 20px;
  cursor: pointer;
}
.body {
  padding: 12px;
  overflow: auto;
  max-height: calc(100vh - 54px);
}
.placeholder {
  color: #777;
}
.loading {
  color: #666;
}
.error {
  color: #b00020;
  white-space: pre-wrap;
}
.author {
  display: flex;
  gap: 12px;
  align-items: center;
  margin-bottom: 10px;
}
.avatar {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  background: #eef2f7;
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
  border: 1px solid rgba(16, 24, 40, 0.06);
}
.avatarImg {
  width: 100%;
  height: 100%;
  object-fit: cover;
  display: block;
}
.avatarText {
  font-size: 14px;
  font-weight: 600;
  color: #111827;
}
.name {
  font-weight: 600;
}
.date {
  color: #666;
  font-size: 12px;
}
.text {
  margin: 8px 0 12px;
  white-space: pre-wrap;
  background: #fafafa;
  border: 1px solid #eee;
  border-radius: 6px;
  padding: 10px;
}
.channel {
  color: #6b7280;
  font-size: 12px;
  margin-top: 8px;
}
.actions {
  margin-top: 12px;
}
.link {
  color: #2563eb;
  text-decoration: none;
}
.link:hover {
  text-decoration: underline;
}

/* Small screens: use overlay behavior to avoid pushing content vertically */
@media (max-width: 768px) {
  .side-panel {
    position: fixed;
    top: 0;
    right: 0;
    width: min(420px, 92vw);
    height: 100vh;
    box-shadow: -8px 0 24px rgba(0, 0, 0, 0.08);
    transform: translateX(100%);
    transition: transform 0.25s ease;
    z-index: 900;
  }
  .side-panel.open {
    transform: translateX(0);
  }
  .body {
    height: calc(100vh - 54px);
  }
}
</style>
