import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useLoadingStore = defineStore('loading', () => {
  const isLoading = ref(false)
  const progress = ref(0)
  let timer: any = null

  const start = () => {
    isLoading.value = true
    progress.value = 15
    if (timer) clearInterval(timer)
    timer = setInterval(() => {
      if (progress.value < 85) {
        progress.value += Math.floor(Math.random() * 12) + 5
      }
    }, 120)
  }

  const finish = () => {
    progress.value = 100
    if (timer) clearInterval(timer)
    setTimeout(() => {
      isLoading.value = false
      setTimeout(() => {
        progress.value = 0
      }, 200)
    }, 300)
  }

  return {
    isLoading,
    progress,
    start,
    finish
  }
})
