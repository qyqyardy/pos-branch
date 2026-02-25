<template>
  <div class="min-h-screen">
    <TopBar />

    <main class="mx-auto max-w-xl px-4 py-10">
      <section
        class="animate-float-in rounded-2xl border border-black/10 bg-[color:var(--paper)] backdrop-blur p-6 shadow-[0_18px_60px_rgba(0,0,0,0.10)]"
      >
        <div class="text-xs tracking-widest uppercase text-[color:var(--muted)]">Akses dibatasi</div>
        <h1 class="mt-2 font-brand text-3xl leading-tight">Fitur tidak tersedia</h1>
        <p class="mt-2 text-sm text-[color:var(--muted)]">
          {{ message }}
        </p>

        <div class="mt-5 rounded-2xl border border-black/10 bg-white/70 p-4">
          <div class="grid grid-cols-1 gap-3 sm:grid-cols-3">
            <div>
              <div class="text-xs tracking-widest uppercase text-[color:var(--muted)]">Plan</div>
              <div class="mt-1 font-semibold">{{ planLabel }}</div>
            </div>
            <div>
              <div class="text-xs tracking-widest uppercase text-[color:var(--muted)]">Status</div>
              <div
                class="mt-1 font-semibold"
                :class="subscriptionActive ? 'text-emerald-700' : 'text-red-700'"
              >
                {{ subscriptionActive ? 'Active' : 'Expired' }}
              </div>
            </div>
            <div>
              <div class="text-xs tracking-widest uppercase text-[color:var(--muted)]">Paid until</div>
              <div class="mt-1 font-mono text-xs">{{ paidUntilLabel }}</div>
            </div>
          </div>
        </div>

        <div v-if="canGoPOS" class="mt-5 flex flex-wrap gap-2">
          <button
            type="button"
            class="rounded-xl bg-[color:var(--accent)] px-4 py-2 text-sm font-semibold text-white shadow-[0_14px_30px_rgba(193,122,59,0.30)] transition hover:brightness-95"
            @click="router.push('/pos')"
          >
            Kembali ke POS
          </button>
        </div>

        <div class="mt-4 text-xs text-[color:var(--muted)]">
          Hubungi admin/vendor untuk upgrade plan atau perpanjang subscription.
        </div>
      </section>
    </main>
  </div>
</template>

<script setup>
import { computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import TopBar from '../components/TopBar.vue'
import { useAuthStore } from '../stores/auth'
import { useSettingsStore } from '../stores/settings'

const router = useRouter()
const auth = useAuthStore()
const settings = useSettingsStore()

onMounted(async () => {
  try {
    await settings.loadStore(auth.token)
  } catch {}
})

const role = computed(() => auth.user?.role || '')
const plan = computed(() => settings.store?.plan || 'premium')
const subscriptionActive = computed(() => settings.store?.subscription_active !== false)

const planLabel = computed(() => (plan.value === 'standard' ? 'Standard' : 'Premium'))
const paidUntilLabel = computed(() => {
  const iso = settings.store?.paid_until
  if (!iso) return '-'
  try {
    return new Date(iso).toLocaleDateString('id-ID', { year: 'numeric', month: 'short', day: '2-digit' })
  } catch {
    return String(iso)
  }
})

const message = computed(() => {
  if (plan.value !== 'premium') return 'Akun kamu membutuhkan plan Premium untuk membuka menu Finance.'
  if (!subscriptionActive.value) return 'Subscription kamu sudah habis. Sistem dalam mode terbatas.'
  return 'Akses dibatasi oleh konfigurasi subscription.'
})

const canGoPOS = computed(() => role.value === 'admin' || role.value === 'cashier')
</script>

