import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import AlbumsView from '@/views/AlbumsView.vue'
import SavedView from '@/views/SavedView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView,
    },
    {
      path: '/about',
      name: 'about',
      // route level code-splitting
      // this generates a separate chunk (About.[hash].js) for this route
      // which is lazy-loaded when the route is visited.
      component: () => import('../views/AboutView.vue'),
    },
    {
      path: '/albums',
      name: 'albums',
      component: AlbumsView,
    },
    {
      path: '/saved',
      name: 'saved',
      component: SavedView,
    },
  ],
})

export default router
