import { createApp } from 'vue'
import { createPinia } from 'pinia'
import naive from 'naive-ui'
import axios from 'axios'
import App from './App.vue'
import router from './router'
import './style.css'

// Global OWASP ASVS V3 & V4: Automatic Bearer Authorization Header & Session Handler
axios.interceptors.request.use((config) => {
  const token = localStorage.getItem('token')
  if (token && config.headers) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

axios.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401 && window.location.pathname !== '/login') {
      localStorage.removeItem('token')
      localStorage.removeItem('refreshToken')
      window.location.href = '/login'
    }
    return Promise.reject(error)
  }
)

const app = createApp(App)

app.use(createPinia())
app.use(router)
app.use(naive)

app.mount('#app')
