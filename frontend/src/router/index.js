import { createRouter, createWebHistory } from 'vue-router'
import Overview from '../views/Overview.vue'
import Proxy from '../views/Proxy.vue'
import Rules from '../views/Rules.vue'
import Settings from '../views/Settings.vue'
import Logs from '../views/Logs.vue'

const routes = [
  { path: '/', component: Overview, meta: { title: '总览' } },
  { path: '/proxy', component: Proxy, meta: { title: '代理' } },
  { path: '/rules', component: Rules, meta: { title: '规则' } },
  { path: '/settings', component: Settings, meta: { title: '配置' } },
  { path: '/logs', component: Logs, meta: { title: '日志' } },
]


const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes
})

export default router
