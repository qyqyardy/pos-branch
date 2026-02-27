import { defineStore } from 'pinia'
import { getStoreSettings, updateStoreSettings } from '../api/api'

const LS_KEY = 'store_settings'

function defaultStore() {
  return {
    name: 'Warkop',
    tagline: 'Point of Sale',
    address_lines: ['', ''],
    phone: '',
    logo_data_url: ''
  }
}

export const useSettingsStore = defineStore('settings', {
  state: () => ({
    store: (() => {
      try {
        const raw = localStorage.getItem(LS_KEY)
        const parsed = raw ? JSON.parse(raw) : null
        return parsed && typeof parsed === 'object'
          ? { ...defaultStore(), ...parsed }
          : defaultStore()
      } catch {
        return defaultStore()
      }
    })(),
    loading: false,
    error: ''
  }),

  actions: {
    setStore(next) {
      this.store = next || defaultStore()
      try {
        localStorage.setItem(LS_KEY, JSON.stringify(this.store))
      } catch { }
    },

    async loadStore(token) {
      if (!token) return
      this.loading = true
      this.error = ''
      try {
        const data = await getStoreSettings(token)
        this.setStore({
          name: data?.name || 'Warkop',
          tagline: data?.tagline || '',
          address_lines: Array.isArray(data?.address_lines)
            ? [data.address_lines[0] || '', data.address_lines[1] || '']
            : ['', ''],
          phone: data?.phone || '',
          logo_data_url: data?.logo_data_url || ''
        })
      } catch (e) {
        this.error = e?.message || 'Gagal memuat setting'
      } finally {
        this.loading = false
      }
    },

    async saveStore(token, payload) {
      if (!token) throw new Error('Unauthorized')
      this.loading = true
      this.error = ''
      try {
        const data = await updateStoreSettings(token, payload)
        this.setStore({
          name: data?.name || payload?.name || 'Warkop',
          tagline: data?.tagline || payload?.tagline || '',
          address_lines: Array.isArray(data?.address_lines)
            ? [data.address_lines[0] || '', data.address_lines[1] || '']
            : Array.isArray(payload?.address_lines)
              ? [payload.address_lines[0] || '', payload.address_lines[1] || '']
              : ['', ''],
          phone: data?.phone || payload?.phone || '',
          logo_data_url:
            data?.logo_data_url != null
              ? data.logo_data_url
              : payload?.logo_data_url || this.store?.logo_data_url || ''
        })
        return this.store
      } catch (e) {
        this.error = e?.message || 'Gagal menyimpan setting'
        throw e
      } finally {
        this.loading = false
      }
    }
  }
})
