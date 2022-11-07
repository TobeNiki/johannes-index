import { defineStore } from 'pinia'
import { LocalStorage } from 'quasar'
export const useAuthStore = defineStore('auth', {
  state: () => ({
    token: '',
    expire: '',
    refreshInterval: 30 * 60 * 1000 // 30minute
  }),
  getters: {
    getAuthrizationHeader: () => 'Bearer ' + LocalStorage.getItem('token')
  },
  actions: {
    setToken (token, expire) {
      this.expire = expire
      this.token = token
      LocalStorage.set('token', token)
    }
  }
})
