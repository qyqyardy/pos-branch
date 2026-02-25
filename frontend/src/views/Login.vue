<template>
  <div class="min-h-screen px-6 py-10 flex items-center justify-center">
    <div class="w-full max-w-md">
      <div class="mb-6 text-center">
        <div v-if="settings.store?.logo_data_url" class="mx-auto mb-3 h-16 w-16 overflow-hidden rounded-3xl border border-black/10 bg-white/70">
          <img :src="settings.store.logo_data_url" alt="" class="h-full w-full object-contain" loading="lazy" />
        </div>
        <div class="text-sm tracking-widest uppercase text-[color:var(--muted)]">
          {{ settings.store?.name || 'Warkop' }}
        </div>
        <h1 class="mt-2 font-brand text-3xl leading-tight">
          POS Kasir
        </h1>
        <p class="mt-2 text-sm text-[color:var(--muted)]">
          Masuk untuk mulai transaksi.
        </p>
      </div>

      <div class="rounded-2xl border border-black/10 bg-white/70 backdrop-blur px-5 py-6 shadow-[0_18px_60px_rgba(0,0,0,0.12)]">
        <div class="space-y-4">
          <label class="block">
            <span class="text-sm font-medium">Email</span>
            <input
              v-model.trim="email"
              type="email"
              autocomplete="username"
              placeholder="admin@pos.com"
              class="mt-1 w-full rounded-xl border border-black/10 bg-white px-4 py-3 outline-none ring-0 focus:border-black/20 focus:shadow-[0_0_0_4px_rgba(193,122,59,0.15)]"
              @keydown.enter="handleLogin"
            />
          </label>

          <label class="block">
            <span class="text-sm font-medium">Password</span>
            <input
              v-model="password"
              type="password"
              autocomplete="current-password"
              placeholder="••••••••"
              class="mt-1 w-full rounded-xl border border-black/10 bg-white px-4 py-3 outline-none ring-0 focus:border-black/20 focus:shadow-[0_0_0_4px_rgba(193,122,59,0.15)]"
              @keydown.enter="handleLogin"
            />
          </label>

          <p v-if="error" class="text-sm text-red-700">
            {{ error }}
          </p>

          <button
            class="w-full rounded-xl bg-[color:var(--accent)] px-4 py-3 font-semibold text-white shadow-[0_14px_30px_rgba(193,122,59,0.35)] transition hover:brightness-95 disabled:cursor-not-allowed disabled:opacity-60"
            :disabled="loading"
            @click="handleLogin"
          >
            <span v-if="!loading">Masuk</span>
            <span v-else>Memproses...</span>
          </button>
        </div>

        <div class="mt-5 text-xs text-[color:var(--muted)]">
          <span class="font-medium">Tips:</span> akun seed: <span class="font-mono">admin@pos.com</span> (password: <span class="font-mono">123456</span>)
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useAuthStore } from '../stores/auth'
import { useSettingsStore } from '../stores/settings'
import { useRouter } from 'vue-router'

const email = ref('')
const password = ref('')
const auth = useAuthStore()
const settings = useSettingsStore()
const router = useRouter()
const loading = ref(false)
const error = ref('')

async function handleLogin() {
  error.value = ''
  loading.value = true
  try {
    await auth.doLogin(email.value, password.value)
    if (auth.user?.role === 'finance') router.push('/finance')
    else router.push('/pos')
  } catch (e) {
    error.value = e?.message || 'Login gagal'
  } finally {
    loading.value = false
  }
}
</script>
