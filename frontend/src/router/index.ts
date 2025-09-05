import { createRouter, createWebHistory } from 'vue-router'
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
  {
    path: '/account',
    name: 'account',
    component: () => import('@/views/AccountView.vue'),
  },
  // テスト用ルート（一時的）
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
]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes,
})

export default router
