import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import { i18n, initTheme } from './i18n'

initTheme()

const app = createApp(App)
app.use(createPinia())
app.use(i18n)
app.mount('#app')
