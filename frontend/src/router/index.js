import { createRouter, createWebHistory } from 'vue-router'
import Overview from '../views/Overview.vue'
import Proxy from '../views/Proxy.vue'
import Rules from '../views/Rules.vue'
import Settings from '../views/Settings.vue'
import Logs from '../views/Logs.vue'
import ApiSetup from '../views/ApiSetup.vue' // Import the new view
import { isApiConfigured } from '@/services/apiConfig' // Import the checker

const routes = [
  { path: '/', component: Overview, meta: { title: '总览', requiresAuth: true } },
  { path: '/proxy', component: Proxy, meta: { title: '代理', requiresAuth: true } },
  { path: '/rules', component: Rules, meta: { title: '规则', requiresAuth: true } },
  { path: '/settings', component: Settings, meta: { title: '配置', requiresAuth: true } },
  { path: '/logs', component: Logs, meta: { title: '日志', requiresAuth: true } },
  {
    path: '/api-setup',
    component: ApiSetup,
    meta: { title: 'API Configuration', requiresAuth: false }
  },
  // Later, Dashboard.vue will replace Overview.vue for the '/' path.
]


const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes
})

router.beforeEach((to, from, next) => {
  const requiresAuth = to.matched.some(record => record.meta.requiresAuth);
  const configured = isApiConfigured();

  if (requiresAuth && !configured) {
    next('/api-setup');
  } else if (to.path === '/api-setup' && configured) {
    // If user is configured and tries to go to /api-setup, redirect to home
    next('/');
  }
  else {
    next();
  }
});

export default router
