import { createRouter, createWebHistory } from 'vue-router'
import Login from '../views/Login.vue'
import POS from '../views/POS.vue'
import Products from '../views/Products.vue'
import Settings from '../views/Settings.vue'
import Finance from '../views/Finance.vue'
import { useAuthStore } from '../stores/auth'
import { useSettingsStore } from '../stores/settings'
import Blocked from '../views/Blocked.vue'

const routes = [
  { path: '/', component: Login, meta: { guestOnly: true } },
  { path: '/pos', component: POS, meta: { requiresAuth: true, roles: ['admin', 'cashier'] } },
  { path: '/products', component: Products, meta: { requiresAuth: true, roles: ['admin'] } },
  { path: '/settings', component: Settings, meta: { requiresAuth: true, roles: ['admin'] } },
  { path: '/finance', component: Finance, meta: { requiresAuth: true, roles: ['admin', 'finance'], plan: 'premium' } },
  { path: '/blocked', component: Blocked, meta: { requiresAuth: true } },
  { path: '/:pathMatch(.*)*', redirect: '/' }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach(async (to, from, next) => {
  const auth = useAuthStore()
  const settings = useSettingsStore()

  if (to.meta?.requiresAuth && !auth.token) return next('/')

  if (auth.token && !auth.user) {
    try {
      await auth.loadMe()
    } catch {
      return next('/')
    }
  }

  const role = auth.user?.role || ''

  // Load store settings when plan gating matters (finance role or premium routes).
  if (auth.token && (role === 'finance' || to.meta?.plan)) {
    try {
      await settings.loadStore(auth.token)
    } catch { }
  }

  const plan = settings.store?.plan || 'premium'
  const homeByRole = () => {
    if (role === 'finance') return plan === 'premium' ? '/finance' : '/blocked'
    if (role === 'admin' || role === 'cashier') return '/pos'
    return '/'
  }

  if (to.meta?.guestOnly && auth.token) return next(homeByRole())

  if (to.meta?.plan === 'premium' && plan !== 'premium') {
    if (role === 'finance') return next('/blocked')
    return next(homeByRole())
  }

  const roles = to.meta?.roles
  if (Array.isArray(roles) && roles.length) {
    if (!role || !roles.includes(role)) return next(homeByRole())
  }
  return next()
})

export default router
