import { createRouter, createWebHistory } from 'vue-router'
import DashboardView from '@/views/DashboardView.vue'
import LoginView from '@/views/LoginView.vue'
import ProfileView from '@/views/ProfileView.vue'
import RegisterView from '@/views/RegisterView.vue'
import AdminView from '@/views/AdminView.vue'
import CatalogView from '@/views/CatalogView.vue'
import ClubView from '@/views/ClubView.vue'
import MatchesView from '@/views/MatchesView.vue'
import StandingsView from '@/views/StandingsView.vue'
import BroadcastsView from '@/views/BroadcastsView.vue'
import AdminBroadcastsView from '@/views/AdminBroadcastsView.vue'
import NewsView from '@/views/NewsView.vue'
import AdminNewsView from '@/views/AdminNewsView.vue'
import SetupView from '@/views/SetupView.vue'
import ResetPasswordView from '@/views/ResetPasswordView.vue'
import { setupRequired } from '@/lib/installation'
import { useSessionStore } from '@/stores/session'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', name: 'dashboard', component: DashboardView, meta: { auth: true } },
    { path: '/login', name: 'login', component: LoginView, meta: { guest: true } },
    { path: '/register', name: 'register', component: RegisterView, meta: { guest: true } },
    { path: '/setup', name: 'setup', component: SetupView, meta: { guest: true } },
    {
      path: '/reset-password',
      name: 'reset-password',
      component: ResetPasswordView,
      meta: { authLayout: true },
    },
    { path: '/profile', name: 'profile', component: ProfileView, meta: { auth: true } },
    { path: '/clubs', name: 'clubs', component: CatalogView, meta: { auth: true } },
    { path: '/clubs/:id', name: 'club', component: ClubView, meta: { auth: true } },
    { path: '/matches', name: 'matches', component: MatchesView, meta: { auth: true } },
    { path: '/standings', name: 'standings', component: StandingsView, meta: { auth: true } },
    { path: '/tv', name: 'tv', component: BroadcastsView, meta: { auth: true } },
    { path: '/news', name: 'news', component: NewsView, meta: { auth: true } },
    { path: '/admin', name: 'admin', component: AdminView, meta: { auth: true, admin: true } },
    {
      path: '/admin/tv',
      name: 'admin-tv',
      component: AdminBroadcastsView,
      meta: { auth: true, admin: true },
    },
    {
      path: '/admin/news',
      name: 'admin-news',
      component: AdminNewsView,
      meta: { auth: true, admin: true },
    },
  ],
})

router.beforeEach(async (to) => {
  const needsSetup = await setupRequired()
  if (needsSetup && to.name !== 'setup') return { name: 'setup' }
  if (!needsSetup && to.name === 'setup') return { name: 'login' }
  const session = useSessionStore()
  if (!session.ready) await session.restore()
  if (to.meta.auth && !session.user) return { name: 'login', query: { next: to.fullPath } }
  if (to.meta.admin && session.user?.role !== 'admin') return { name: 'dashboard' }
  if (to.meta.guest && session.user) return { name: 'dashboard' }
})

export default router
