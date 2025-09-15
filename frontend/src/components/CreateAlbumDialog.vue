<template>
  <div v-if="isVisible" class="dialog-overlay" @click="handleOverlayClick">
    <div class="dialog" @click.stop>
      <div class="dialog-header">
        <h3>新しいアルバムを作成</h3>
        <button type="button" class="close-button" @click="$emit('close')">×</button>
      </div>

      <form @submit.prevent="handleSubmit" class="dialog-content">
        <div class="form-group">
          <label for="album-title">アルバム名</label>
          <input
            id="album-title"
            v-model="albumTitle"
            type="text"
            placeholder="アルバム名を入力"
            required
            :disabled="isCreating"
          />
        </div>

        <div class="form-group">
          <label for="album-description">説明（任意）</label>
          <textarea
            id="album-description"
            v-model="albumDescription"
            placeholder="アルバムの説明を入力"
            :disabled="isCreating"
            rows="3"
          ></textarea>
        </div>

        <div class="selected-images-info">
          <p>選択された画像: {{ selectedImageCount }}枚</p>
        </div>

        <!-- エラーメッセージ -->
        <div v-if="errorMessage" class="error-message">
          {{ errorMessage }}
        </div>

        <div class="dialog-actions">
          <button
            type="button"
            @click="$emit('close')"
            :disabled="isCreating"
            class="cancel-button"
          >
            キャンセル
          </button>
          <button type="submit" :disabled="!albumTitle.trim() || isCreating" class="create-button">
            {{ isCreating ? '作成中...' : '作成' }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'

interface Props {
  isVisible: boolean
  selectedImageCount: number
}

interface Emits {
  (e: 'close'): void
  (e: 'create', data: { title: string; description: string }): Promise<void>
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const albumTitle = ref('')
const albumDescription = ref('')
const isCreating = ref(false)
const errorMessage = ref('')

const handleOverlayClick = () => {
  if (!isCreating.value) {
    emit('close')
  }
}

const handleSubmit = async () => {
  if (!albumTitle.value.trim() || isCreating.value) return

  isCreating.value = true
  errorMessage.value = ''

  try {
    await emit('create', {
      title: albumTitle.value.trim(),
      description: albumDescription.value.trim(),
    })
    // 成功時は親コンポーネントがダイアログを閉じる
  } catch (error) {
    errorMessage.value = error instanceof Error ? error.message : 'アルバムの作成に失敗しました'
  } finally {
    isCreating.value = false
  }
}

// ダイアログが閉じられたときにフォームをリセット
const resetForm = () => {
  albumTitle.value = ''
  albumDescription.value = ''
  isCreating.value = false
  errorMessage.value = ''
}

// props.isVisibleが変更されたときにフォームをリセット
watch(
  () => props.isVisible,
  (newValue) => {
    if (!newValue) {
      resetForm()
    }
  },
)
</script>

<style scoped lang="scss">
.dialog-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.dialog {
  background: white;
  border-radius: 8px;
  width: 90%;
  max-width: 500px;
  max-height: 90vh;
  overflow: hidden;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.15);
}

.dialog-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1.5rem;
  border-bottom: 1px solid #e0e0e0;

  h3 {
    margin: 0;
    color: #333;
    font-size: 1.25rem;
  }
}

.close-button {
  background: none;
  border: none;
  font-size: 1.5rem;
  cursor: pointer;
  color: #666;
  padding: 0;
  width: 30px;
  height: 30px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 4px;

  &:hover {
    background-color: #f0f0f0;
    color: #333;
  }
}

.dialog-content {
  padding: 1.5rem;
}

.form-group {
  margin-bottom: 1.5rem;

  label {
    display: block;
    margin-bottom: 0.5rem;
    font-weight: 600;
    color: #333;
  }

  input,
  textarea {
    width: 100%;
    padding: 0.75rem;
    border: 1px solid #ddd;
    border-radius: 4px;
    font-size: 1rem;
    transition: border-color 0.2s;

    &:focus {
      outline: none;
      border-color: #007bff;
      box-shadow: 0 0 0 3px rgba(0, 123, 255, 0.1);
    }

    &:disabled {
      background-color: #f8f9fa;
      cursor: not-allowed;
    }
  }

  textarea {
    resize: vertical;
    min-height: 80px;
  }
}

.selected-images-info {
  background-color: #f8f9fa;
  padding: 1rem;
  border-radius: 4px;
  margin-bottom: 1.5rem;

  p {
    margin: 0;
    color: #666;
    font-size: 0.9rem;
  }
}

.error-message {
  background-color: #f8d7da;
  color: #721c24;
  padding: 1rem;
  border-radius: 4px;
  margin-bottom: 1.5rem;
  font-size: 0.9rem;
  border: 1px solid #f5c6cb;
}
.dialog-actions {
  display: flex;
  gap: 1rem;
  justify-content: flex-end;

  button {
    padding: 0.75rem 1.5rem;
    border: none;
    border-radius: 4px;
    font-size: 1rem;
    cursor: pointer;
    transition: all 0.2s;

    &:disabled {
      cursor: not-allowed;
      opacity: 0.6;
    }
  }
}

.cancel-button {
  background-color: #6c757d;
  color: white;

  &:hover:not(:disabled) {
    background-color: #5a6268;
  }
}

.create-button {
  background-color: #007bff;
  color: white;

  &:hover:not(:disabled) {
    background-color: #0056b3;
  }
}
</style>
