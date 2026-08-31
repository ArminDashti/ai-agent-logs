import { createRouter, createWebHistory } from 'vue-router'
import LoginView from '@/views/LoginView.vue'
import DashboardView from '@/views/DashboardView.vue'
import SessionsView from '@/views/SessionsView.vue'
import SessionDetailView from '@/views/SessionDetailView.vue'
import AboutMeView from '@/views/AboutMeView.vue'
import { getToken } from '@/lib/auth'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    { path: '/', redirect: '/dashboard' },
    { path: '/login', name: 'login', component: LoginView, meta: { guest: true } },
    { path: '/dashboard', name: 'dashboard', component: DashboardView, meta: { requiresAuth: true } },
    { path: '/sessions', name: 'sessions', component: SessionsView, meta: { requiresAuth: true } },
    {
      path: '/sessions/:id',
      name: 'session-detail',
      component: SessionDetailView,
      meta: { requiresAuth: true },
    },
    { path: '/about', name: 'about', component: AboutMeView },
  ],
})

router.beforeEach((to) => {
  const token = getToken()
  if (to.meta.requiresAuth && !token) {
    return { name: 'login', query: { redirect: to.fullPath } }
  }
  if (to.meta.guest && token) {
    return { name: 'dashboard' }
  }
  return true
})

export default router
