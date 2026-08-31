import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import { initTheme } from './lib/theme'
import { useAuth } from './lib/useAuth'
import './assets/index.css'

initTheme()

const { hydrate } = useAuth()

void hydrate().then(() => {
  createApp(App).use(router).mount('#app')
})
