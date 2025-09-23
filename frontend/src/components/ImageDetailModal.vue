<template>
  <div v-if="visible" class="modal-overlay" @click="onOverlay">
    <div class="modal" @click.stop>
      <header class="modal-header">
        <h3>投稿の詳細</h3>
        <button class="close" @click="$emit('close')">×</button>
      </header>
      <section class="modal-body">
        <div v-if="loading" class="loading">読み込み中...</div>
        <div v-else-if="error" class="error">{{ error }}</div>
        <div v-else-if="message">
          <div class="meta">
            <div class="meta-row">
              <span class="label">ユーザーID</span>
              <span class="value">{{ message.userId }}</span>
            </div>
            <div class="meta-row">
              <span class="label">チャンネル</span>
              <span class="value">{{ message.channelId }}</span>
            </div>
            <div class="meta-row">
              <span class="label">日時</span>
              <span class="value">{{ formatDate(message.createdAt) }}</span>
            </div>
          </div>
          <div class="content" v-html="formattedContent"></div>
        </div>
      </section>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, watch, ref } from 'vue'

interface Props {
  visible: boolean
  imageId: string | null
}

const props = defineProps<Props>()
const emit = defineEmits<{ close: [] }>()

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

watch(
  () => props.imageId,
  async (id) => {
    if (!props.visible || !id) return
    await fetchDetail(id)
  },
  { immediate: true },
)

watch(
  () => props.visible,
  async (v) => {
    if (v && props.imageId) {
      await fetchDetail(props.imageId)
    }
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
    if (!res.ok) {
      throw new Error(await res.text())
    }
    message.value = (await res.json()) as Message
  } catch (e: unknown) {
    error.value = e instanceof Error ? e.message : String(e)
  } finally {
    loading.value = false
  }
}

function onOverlay(e: MouseEvent) {
  if (e.target === e.currentTarget) emit('close')
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
  // 画像URLの除去（replaceAll互換のため split/join を使用）
  const url = `https://q.trap.jp/files/${props.imageId ?? ''}`
  if (url) {
    text = text.split(url).join('')
  }
  // 末尾の空行などを削除
  text = text.replace(/\n+$/g, '')
  // 改行を <br> に
  return text.replace(/\n/g, '<br>')
})
</script>

<style scoped>
.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.4);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}
.modal {
  width: min(720px, 92vw);
  max-height: 80vh;
  background: #fff;
  border-radius: 10px;
  overflow: hidden;
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.2);
  display: flex;
  flex-direction: column;
}
.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  border-bottom: 1px solid #eee;
}
.close {
  background: none;
  border: none;
  font-size: 20px;
  cursor: pointer;
}
.modal-body {
  padding: 16px;
  overflow: auto;
}
.meta {
  display: grid;
  grid-template-columns: 120px 1fr;
  gap: 6px 12px;
  font-size: 14px;
  margin-bottom: 12px;
  color: #333;
}
.label {
  color: #666;
}
.content {
  padding: 12px;
  background: #fafafa;
  border: 1px solid #eee;
  border-radius: 6px;
  white-space: pre-wrap;
}
.loading {
  color: #666;
}
.error {
  color: #b00020;
}
</style>
