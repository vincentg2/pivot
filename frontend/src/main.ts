import { createPinia } from 'pinia'
import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import { i18n } from './i18n'
import './style.css'

createApp(App).use(createPinia()).use(i18n).use(router).mount('#app')
