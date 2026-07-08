import { createRouter, createWebHashHistory } from 'vue-router'
import Home from '../views/Home.vue'
import DeckManage from '../views/DeckManage.vue'
import Training from '../views/Training.vue'

const routes = [
  { path: '/', name: 'Home', component: Home },
  { path: '/deck/:id', name: 'DeckManage', component: DeckManage },
  { path: '/training', name: 'Training', component: Training }
]

const router = createRouter({
  history: createWebHashHistory(),
  routes
})

export default router
