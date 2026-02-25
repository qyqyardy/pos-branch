<template>
  <header class="sticky top-0 z-20 border-b border-black/10 bg-[color:var(--paper)] backdrop-blur">
    <div class="mx-auto flex max-w-7xl items-center gap-4 px-4 py-4">
      <div class="flex min-w-0 items-center gap-3">
        <div
          v-if="settings.store?.logo_data_url"
          class="h-11 w-11 overflow-hidden rounded-2xl border border-black/10 bg-white/70"
        >
          <img
            :src="settings.store.logo_data_url"
            alt=""
            class="h-full w-full object-contain"
            loading="lazy"
          />
        </div>
        <div class="font-brand text-2xl leading-none">
          {{ settings.store?.name || 'Warkop' }}
        </div>
        <div class="hidden sm:block">
          <div class="mt-1 text-xs tracking-widest uppercase text-[color:var(--muted)]">
            {{ settings.store?.tagline || 'Point of Sale' }}
          </div>
        </div>
      </div>

      <div class="flex-1" />

      <nav class="flex items-center gap-2">
        <button
          v-if="showPOS"
          type="button"
          class="rounded-xl border border-black/10 bg-white/70 px-3 py-2 text-sm font-semibold hover:bg-white"
          @click="go('/pos')"
        >
          POS
        </button>

        <button
          v-if="canFinance"
          type="button"
          class="rounded-xl border border-black/10 bg-white/70 px-3 py-2 text-sm font-semibold hover:bg-white"
          @click="go('/finance')"
        >
          Finance
        </button>

        <button
          v-if="canSettings"
          type="button"
          class="rounded-xl border border-black/10 bg-white/70 px-3 py-2 text-sm font-semibold hover:bg-white"
          @click="go('/settings')"
        >
          Setting
        </button>
      </nav>

      <button
        type="button"
        class="rounded-xl border border-black/10 bg-white/70 px-3 py-2 text-sm font-semibold hover:bg-white"
        @click="logout"
      >
        Keluar
      </button>
    </div>
  </header>
</template>

<script setup>
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { useSettingsStore } from '../stores/settings'

const router = useRouter()
const auth = useAuthStore()
const settings = useSettingsStore()

const role = computed(() => auth.user?.role || '')
const plan = computed(() => settings.store?.plan || 'premium')
const canSettings = computed(() => role.value === 'admin')
const canFinance = computed(
  () => (role.value === 'admin' || role.value === 'finance') && plan.value === 'premium'
)
const showPOS = computed(() => role.value === 'admin' || role.value === 'cashier')

function go(path) {
  router.push(path)
}

function logout() {
  auth.logout()
  router.push('/')
}
</script>
