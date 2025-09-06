<template>
  <div :class="$style.userProfile">
    <div v-if="isLoggedIn" :class="$style.loggedInSection">
      <!-- ユーザーアバター（アイコン） -->
      <div :class="$style.avatarContainer">
        <div :class="$style.avatar">
          <img
            v-if="userStore.me"
            :src="`https://q.trap.jp/api/v3/public/icon/${userStore.me.name}`"
            :alt="userStore.me.displayName || userStore.me.name"
            :class="$style.avatarImg"
            @error="avatarError = true"
            @load="avatarError = false"
            v-show="!avatarError"
          />
          <span v-if="avatarError || !userStore.me" :class="$style.avatarText">
            {{ userInitials }}
          </span>
        </div>
      </div>

      <!-- ユーザー情報 -->
      <div :class="$style.userInfo">
        <div :class="$style.displayName">{{ userStore.me?.displayName || userStore.me?.name }}</div>
        <div :class="$style.username">@{{ userStore.me?.name }}</div>
      </div>

      <!-- ログアウトボタン -->
      <button :class="$style.logoutButton" @click="logout" title="ログアウト">
        <svg
          width="16"
          height="16"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
        >
          <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"></path>
          <polyline points="16,17 21,12 16,7"></polyline>
          <line x1="21" y1="12" x2="9" y2="12"></line>
        </svg>
      </button>
    </div>

    <div v-else :class="$style.loggedOutSection" @click="login">
      <!-- ログインアイコン -->
      <div :class="$style.loginIcon">
        <svg
          width="20"
          height="20"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
        >
          <path d="M15 3h4a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2h-4"></path>
          <polyline points="10,17 15,12 10,7"></polyline>
          <line x1="15" y1="12" x2="3" y2="12"></line>
        </svg>
      </div>

      <!-- ログインテキスト -->
      <div :class="$style.loginText">
        <div :class="$style.loginMainText">ログイン</div>
        <div :class="$style.loginSubText">traQでログイン</div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { useUserStore } from '@/stores/userStore'

const userStore = useUserStore()

const avatarError = ref(false)

const isLoggedIn = computed(() => userStore.me !== null)

const userInitials = computed(() => {
  if (!userStore.me) return ''
  const displayName = userStore.me.displayName || userStore.me.name
  return displayName.slice(0, 2).toUpperCase()
})

const login = () => {
  const cb = location.pathname + location.search + location.hash
  location.href = `/api/auth/request?callback=${encodeURIComponent(cb)}`
}

const logout = async () => {
  await fetch('/api/auth/logout', { method: 'POST', credentials: 'include' })
  await userStore.fetchMe()
}

onMounted(async () => {
  if (userStore.me === null) {
    await userStore.fetchMe()
  }
})

// userStore.me が変わったら avatarError をリセット
watch(
  () => userStore.me,
  () => {
    avatarError.value = false
  },
)
</script>

<style lang="scss" module>
.userProfile {
  border-radius: 10px;
  overflow: hidden;
  background: #ffffff;
  color: #222;
  border: 1px solid rgba(16, 24, 40, 0.06); /* subtle card border */
}

.loggedInSection {
  display: flex;
  align-items: center;
  padding: 10px 12px;
  gap: 12px;
}

.avatarContainer {
  flex-shrink: 0;
}

.avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: #f3f4f6; /* neutral gray */
  display: flex;
  align-items: center;
  justify-content: center;
  border: 1px solid rgba(16, 24, 40, 0.06);
  overflow: hidden;
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
  color: #111827; /* dark text */
}

.userInfo {
  flex: 1;
  min-width: 0;
}

.displayName {
  font-size: 13px;
  font-weight: 600;
  color: #0f172a;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  line-height: 1.2;
}

.username {
  font-size: 12px;
  color: #64748b; /* cool gray */
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  line-height: 1.2;
  margin-top: 2px;
}

.logoutButton {
  flex-shrink: 0;
  background: transparent;
  border: 1px solid rgba(16, 24, 40, 0.06);
  border-radius: 8px;
  width: 34px;
  height: 34px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #0f172a;
  cursor: pointer;
  transition: all 0.12s ease;

  &:hover {
    background: #f8fafc;
    transform: translateY(-1px);
  }

  &:active {
    transform: translateY(0);
  }
}

.loggedOutSection {
  display: flex;
  align-items: center;
  padding: 12px;
  gap: 12px;
  cursor: pointer;
  transition:
    background 0.12s ease,
    transform 0.12s ease;

  &:hover {
    background: #f8fafc;
    transform: translateY(-1px);
  }

  &:active {
    transform: translateY(0);
  }
}

.loginIcon {
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: #eef2ff; /* soft accent */
  border: 1px solid rgba(16, 24, 40, 0.04);
}

.loginText {
  flex: 1;
}

.loginMainText {
  font-size: 13px;
  font-weight: 600;
  color: #0f172a;
  line-height: 1.2;
}

.loginSubText {
  font-size: 11px;
  color: #64748b;
  line-height: 1.2;
  margin-top: 2px;
}
</style>
