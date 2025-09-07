import { createRouter, createWebHistory } from 'vue-router'
import { env } from '@/config/env'
import { useUserStore } from '@/stores/userStore'
import HomeView from '@/views/HomeView.vue'

const routes = [
  {
    path: '/',
    name: 'home',
    component: HomeView,
  },
  {
    path: '/about',
    name: 'about',
    component: () => import('../views/AboutView.vue'),
  },
  {
    path: '/albums',
    name: 'albumList',
    // this generates a separate chunk (Albums.[hash].js) for this route
    // which is lazy-loaded when the route is visited.
    component: () => import('@/views/AlbumListView.vue'),
  },
  {
    path: '/albums/:id',
    name: 'album',
    component: () => import('@/views/AlbumView.vue'),
    props: true,
    // TODO: beforeEnterの処理
  },
  {
    path: '/saved',
    name: 'saved',
    component: () => import('@/views/SavedView.vue'),
  },
  // テスト用ルート（一時的）
  {
    path: '/test',
    name: 'test',
    component: () => import('@/views/TestView.vue'),
  },
  {
    path: '/test/album-card',
    name: 'albumCardTest',
    component: () => import('@/views/AlbumCardTestView.vue'),
  },
  {
    path: '/test/image-card',
    name: 'imageCardTest',
    component: () => import('@/views/ImageCardTestView.vue'),
  },
  {
    path: '/test/traq-file',
    name: 'traqFileTest',
    component: () => import('@/views/TraqFileTestView.vue'),
  },
  {
    path: '/test/traq-image-search',
    name: 'traqImageSearchTest',
    component: () => import('@/views/TraqMessageSearchTestView.vue'),
  },
  {
    path: '/test/image-detail',
    name: 'imagePostDetailTest',
    component: () => import('@/views/ImageDetailTest.vue'),
  },
  {
    path: '/test/oauth-debug',
    name: 'oauthDebugTest',
    component: () => import('@/views/OAuthDebugTestView.vue'),
  },
  {
    path: '/login',
    name: 'login',
    component: () => import('@/views/LoginRedirect.vue'),
  },
]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes,
})

// ログイン必須フラグが有効な場合のみ、未ログインを /login に誘導
router.beforeEach(async (to) => {
  if (!env.VITE_REQUIRE_LOGIN) return true

  // 免除ルート: ログイン導線自体と、開発用のOAuthデバッグ等
  const publicPaths = new Set<string>([
    '/login',
    '/test/oauth-debug',
  ])

  // 任意: /about など完全公開にしたい場合は追加
  // publicPaths.add('/about')

  // /login 自体ではガードせずに通す（LoginRedirect.vue が自前で /api/auth/me を確認し、次へ進める）
  if (to.path === '/login') return true
  if (publicPaths.has(to.path)) return true

  const user = useUserStore()

  // me が未取得の場合は、読み込み中なら完了を待機、そうでなければ取得
  if (user.me === null) {
    if (user.loading) {
      // 読み込み完了を短い間待つ（最大~1秒）
      const sleep = (ms: number) => new Promise((r) => setTimeout(r, ms))
      for (let i = 0; i < 20 && user.loading; i++) {
        await sleep(50)
      }
    }
    if (user.me === null) {
      await user.fetchMe()
    }
  }

  if (user.me === null) {
    return { path: '/login', query: { next: to.fullPath } }
  }

  return true
})

export default router
