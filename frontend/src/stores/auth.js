import { defineStore } from 'pinia'
import { getMe, login } from '../api/api'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    token: localStorage.getItem('token') || null,
    user: JSON.parse(localStorage.getItem('user') || 'null')
  }),

  actions: {
    async doLogin(email, password) {
      const res = await login(email, password)
      if (!res?.token) throw new Error('Login gagal')
      this.token = res.token
      localStorage.setItem('token', res.token)
      await this.loadMe()
    },

    async loadMe() {
      if (!this.token) return

      try {
        const me = await getMe(this.token)
        this.user = me || null
        localStorage.setItem('user', JSON.stringify(this.user))
      } catch (e) {
        // Token exists but server rejects it (e.g. DB reset). Force re-login.
        this.logout()
        throw e
      }
    },

    logout() {
      this.token = null
      this.user = null
      localStorage.removeItem('token')
      localStorage.removeItem('user')
    }
  }
})
