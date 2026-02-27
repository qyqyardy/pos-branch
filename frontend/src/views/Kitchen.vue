<template>
  <div class="min-h-screen">
    <TopBar />

    <main class="mx-auto max-w-7xl px-4 py-6">
      <div class="flex flex-wrap items-center justify-between gap-4 mb-6">
        <div>
          <h1 class="font-brand text-3xl">Layar Dapur</h1>
          <p class="text-[color:var(--muted)] text-sm">Monitor pesanan aktif secara real-time.</p>
        </div>
        <div class="flex items-center gap-3">
          <div class="flex items-center gap-2 px-3 py-1.5 rounded-full bg-emerald-50 text-emerald-700 text-xs font-bold border border-emerald-100">
            <span class="relative flex h-2 w-2">
              <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-emerald-400 opacity-75"></span>
              <span class="relative inline-flex rounded-full h-2 w-2 bg-emerald-500"></span>
            </span>
            KDS LIVE
          </div>
          <button @click="loadOrders" class="p-2 rounded-xl border border-black/10 bg-white hover:bg-black/5 transition">
            <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M3 12a9 9 0 0 1 9-9 9.75 9.75 0 0 1 6.74 2.74L21 8"/><path d="M21 3v5h-5"/><path d="M21 12a9 9 0 0 1-9 9 9.75 9.75 0 0 1-6.74-2.74L3 16"/><path d="M8 16H3v5"/></svg>
          </button>
        </div>
      </div>

      <div v-if="loading && !orders.length" class="flex flex-col items-center justify-center py-20 opacity-50">
        <div class="h-10 w-10 animate-spin rounded-full border-4 border-[color:var(--accent)] border-t-transparent"></div>
        <p class="mt-4 font-medium">Memuat pesanan...</p>
      </div>

      <div v-else-if="!orders.length" class="flex flex-col items-center justify-center py-20 grayscale opacity-40">
        <div class="mb-4 text-6xl">🍳</div>
        <p class="text-xl font-brand">Dapur Bersih!</p>
        <p class="text-sm">Belum ada pesanan yang perlu dimasak.</p>
      </div>

      <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6 items-start">
        <div 
          v-for="order in sortedOrders" 
          :key="order.id"
          class="group rounded-3xl border border-black/10 transition-all duration-300 overflow-hidden flex flex-col h-full bg-[color:var(--paper)] shadow-[0_8px_30px_rgba(0,0,0,0.04)] hover:shadow-[0_20px_50px_rgba(0,0,0,0.1)] hover:-translate-y-1"
          :class="{
            'border-orange-200 bg-orange-50/30': order.kitchen_status === 'preparing',
            'border-blue-200 bg-blue-50/30': order.kitchen_status === 'ready'
          }"
        >
          <!-- Order Header -->
          <div class="px-5 py-4 border-b border-black/5 flex items-start justify-between bg-white/50">
            <div>
              <div class="text-[10px] tracking-widest uppercase font-bold text-[color:var(--muted)] mb-1">
                Order #{{ orderNoFromId(order.id) }}
              </div>
              <div class="font-brand text-xl">
                {{ order.table_no ? 'Meja ' + order.table_no : (order.order_type === 'take_away' ? '📦 Take Away' : 'Dine In') }}
              </div>
            </div>
            <div class="text-right">
              <div class="text-xs font-mono font-bold">{{ formatTime(order.created_at) }}</div>
              <div class="text-[10px] text-[color:var(--muted)]">{{ timeAgo(order.created_at) }}</div>
            </div>
          </div>

          <!-- Items List -->
          <div class="flex-1 px-5 py-4 space-y-3">
            <div 
              v-for="(item, idx) in order.details?.items" 
              :key="idx"
              class="flex items-start gap-3"
            >
              <div class="flex-shrink-0 w-8 h-8 rounded-lg bg-black/5 flex items-center justify-center font-bold text-sm">
                {{ item.qty }}x
              </div>
              <div class="flex-1">
                <div class="font-bold text-sm leading-tight">{{ item.name }}</div>
                <div v-if="item.notes" class="text-[10px] text-orange-600 font-medium italic mt-0.5">
                  "{{ item.notes }}"
                </div>
              </div>
            </div>
          </div>

          <!-- Order Footer Actions -->
          <div class="px-5 py-4 bg-white/40 border-t border-black/5">
             <div class="grid grid-cols-1 gap-2">
                <button 
                  v-if="order.kitchen_status === 'pending'"
                  @click="updateStatus(order.id, 'preparing')"
                  class="w-full py-3 rounded-2xl bg-orange-500 hover:bg-orange-400 text-white font-bold text-sm shadow-[0_10px_20px_rgba(249,115,22,0.3)] transition-all active:scale-95 flex items-center justify-center gap-2"
                >
                  <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><path d="M3 14h18"/><path d="M7 14v-2a5 5 0 0 1 5-5v0a5 5 0 0 1 5 5v2"/><path d="M5 14v4a2 2 0 0 0 2 2h10a2 2 0 0 0 2-2v-4"/><path d="M10 9a3 3 0 1 1 4 0"/></svg>
                  Proses Memasak
                </button>

                <button 
                  v-if="order.kitchen_status === 'preparing'"
                  @click="updateStatus(order.id, 'ready')"
                  class="w-full py-3 rounded-2xl bg-blue-600 hover:bg-blue-500 text-white font-bold text-sm shadow-[0_10px_20px_rgba(37,99,235,0.3)] transition-all active:scale-95 flex items-center justify-center gap-2"
                >
                  <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/><polyline points="22 4 12 14.01 9 11.01"/></svg>
                  Siap Disajikan
                </button>

                <button 
                  v-if="order.kitchen_status === 'ready'"
                  @click="updateStatus(order.id, 'done')"
                  class="w-full py-3 rounded-2xl bg-emerald-600 hover:bg-emerald-500 text-white font-bold text-sm shadow-[0_10px_20px_rgba(16,185,129,0.3)] transition-all active:scale-95 flex items-center justify-center gap-2"
                >
                  <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><path d="m5 12 5 5L20 7"/></svg>
                  Sudah Diambil
                </button>

                <div v-if="order.updating" class="text-[10px] text-center font-bold text-black/40 animate-pulse">MEMPROSES...</div>
             </div>
          </div>
        </div>
      </div>
    </main>

    <!-- Simple Toast -->
    <div v-if="toast" class="fixed bottom-6 left-1/2 -translate-x-1/2 z-[100] animate-bounce-in">
      <div class="px-6 py-3 rounded-2xl bg-black text-white text-sm font-bold shadow-2xl flex items-center gap-3">
        <span>{{ toast }}</span>
        <button @click="toast = ''" class="opacity-50 hover:opacity-100">×</button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted, onUnmounted, ref } from 'vue'
import TopBar from '../components/TopBar.vue'
import { useAuthStore } from '../stores/auth'
import { listOrders, getOrder, updateOrderStatus, ApiError } from '../api/api'

const auth = useAuthStore()
const orders = ref([])
const loading = ref(false)
const toast = ref('')
let pollTimer = null

const sortedOrders = computed(() => {
  // Sort by created_at ascending (oldest first for FIFO)
  return [...orders.value].sort((a, b) => new Date(a.created_at) - new Date(b.created_at))
})

async function loadOrders() {
  if (loading.value && orders.value.length === 0) return
  loading.value = true
  try {
    const today = new Date().toISOString().split('T')[0]
    // Fetch orders for today
    const list = await listOrders(auth.token, { date: today })
    
    // Filter only non-finished orders for KDS
    const activeOrders = list.filter(o => o.kitchen_status !== 'done')
    
    // For each order, we need items, so we'll fetch details if not already present
    // To keep it simple, we fetch details for all active orders
    const enriched = await Promise.all(activeOrders.map(async (o) => {
      const existing = orders.value.find(prev => prev.id === o.id)
      if (existing && existing.details) {
        return { ...o, details: existing.details }
      }
      try {
        const details = await getOrder(auth.token, o.id)
        return { ...o, details }
      } catch {
        return o
      }
    }))

    orders.value = enriched
  } catch (e) {
    console.error('KDS load error:', e)
  } finally {
    loading.value = false
  }
}

async function updateStatus(id, newStatus) {
  const order = orders.value.find(o => o.id === id)
  if (!order) return
  
  order.updating = true
  try {
    await updateOrderStatus(auth.token, id, newStatus)
    showToast(`Order #${orderNoFromId(id)} updated to ${newStatus}`)
    await loadOrders()
  } catch (e) {
    showToast(e.message || 'Gagal update status')
  } finally {
    if (order) order.updating = false
  }
}

function showToast(msg) {
  toast.value = msg
  setTimeout(() => toast.value = '', 3000)
}

function orderNoFromId(orderId) {
  if (!orderId) return '-----'
  const hex = String(orderId).replace(/-/g, '').slice(0, 8)
  const n = Number.parseInt(hex, 16)
  return String(n % 100000).padStart(5, '0')
}

function formatTime(d) {
  const date = new Date(d)
  return date.toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit' })
}

function timeAgo(d) {
  const seconds = Math.floor((new Date() - new Date(d)) / 1000)
  if (seconds < 60) return 'Baru saja'
  const mins = Math.floor(seconds / 60)
  return `${mins} menit lalu`
}

onMounted(() => {
  loadOrders()
  pollTimer = setInterval(loadOrders, 5000) // Poll every 5 seconds
})

onUnmounted(() => {
  if (pollTimer) clearInterval(pollTimer)
})
</script>

<style scoped>
@keyframes bounce-in {
  0% { transform: translate(-50%, 20px) scale(0.9); opacity: 0; }
  70% { transform: translate(-50%, -5px) scale(1.05); }
  100% { transform: translate(-50%, 0) scale(1); opacity: 1; }
}
.animate-bounce-in {
  animation: bounce-in 0.4s cubic-bezier(0.175, 0.885, 0.32, 1.275) forwards;
}
</style>
