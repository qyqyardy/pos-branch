<template>
  <div class="min-h-screen">
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
          <div class="min-w-0">
            <div class="font-brand text-2xl leading-none">{{ settings.store?.name || 'Warkop' }}</div>
            <div class="mt-1 text-xs tracking-widest uppercase text-[color:var(--muted)]">
              {{ settings.store?.tagline || 'Point of Sale' }}
            </div>
          </div>
        </div>

        <div class="flex-1">
          <input
            v-model="search"
            type="search"
            placeholder="Cari menu..."
            class="w-full rounded-2xl border border-black/10 bg-white/70 px-4 py-3 text-sm outline-none focus:border-black/20 focus:shadow-[0_0_0_4px_rgba(193,122,59,0.15)]"
          />
        </div>

        <button
          v-if="auth.user?.role === 'admin' || auth.user?.role === 'finance'"
          type="button"
          class="rounded-xl border border-black/10 bg-white/70 px-3 py-2 text-sm font-semibold hover:bg-white"
          @click="router.push('/finance')"
        >
          Finance
        </button>

        <button
          v-if="auth.user?.role === 'admin'"
          type="button"
          class="rounded-xl border border-black/10 bg-white/70 px-3 py-2 text-sm font-semibold hover:bg-white"
          @click="router.push('/settings')"
        >
          Setting
        </button>

        <button
          type="button"
          class="rounded-xl border border-black/10 bg-white/70 px-3 py-2 text-sm font-semibold hover:bg-white"
          @click="handleLogout"
        >
          Keluar
        </button>
      </div>
    </header>

    <main class="mx-auto grid max-w-7xl gap-6 px-4 py-6 lg:grid-cols-[1fr,380px]">
      <section class="min-w-0">
        <div class="mb-4 flex items-end justify-between gap-4">
          <div>
            <h2 class="font-brand text-xl">Menu</h2>
            <p class="mt-1 text-sm text-[color:var(--muted)]">
              Klik item untuk menambah ke keranjang.
            </p>
          </div>
          <div class="text-right text-sm text-[color:var(--muted)]">
            <div>{{ filteredProducts.length }} item</div>
          </div>
        </div>

        <div
          v-if="loadingProducts"
          class="grid grid-cols-2 gap-3 sm:grid-cols-3 lg:grid-cols-4"
        >
          <div
            v-for="n in 12"
            :key="n"
            class="h-28 rounded-2xl border border-black/10 bg-white/40 animate-pulse-soft"
          />
        </div>

        <div
          v-else-if="productsError"
          class="rounded-2xl border border-black/10 bg-white/60 px-5 py-5 shadow-[0_10px_30px_rgba(0,0,0,0.08)]"
        >
          <div class="font-semibold">Gagal memuat produk</div>
          <div class="mt-1 text-sm text-[color:var(--muted)]">
            {{ productsError }}
          </div>
          <button
            type="button"
            class="mt-4 rounded-xl border border-black/10 bg-white/70 px-4 py-2 text-sm font-semibold hover:bg-white"
            @click="loadProducts"
          >
            Coba lagi
          </button>
        </div>

	        <div
	          v-else
	          class="grid grid-cols-2 gap-3 sm:grid-cols-3 lg:grid-cols-4"
	        >
	          <button
	            v-for="p in filteredProducts"
	            :key="p.id"
	            type="button"
	            @click="addToCart(p)"
	            class="group relative overflow-hidden rounded-2xl border border-black/10 bg-white/60 px-4 py-4 text-left shadow-[0_10px_30px_rgba(0,0,0,0.08)] transition hover:-translate-y-0.5 hover:bg-white/85 hover:shadow-[0_18px_50px_rgba(0,0,0,0.12)]"
	          >
	            <div class="absolute right-3 top-3 rounded-full bg-[color:rgba(193,122,59,0.12)] px-2 py-1 text-xs font-semibold text-[color:var(--accent)]">
	              + Tambah
	            </div>

	            <div class="mb-3 overflow-hidden rounded-xl border border-black/10 bg-white/50">
	              <img
	                :src="menuImageFor(p.name)"
	                alt=""
	                class="h-24 w-full object-cover"
	                loading="lazy"
	                @error="onMenuImgError"
	              />
	            </div>

	            <div class="min-h-[2.75rem] font-semibold leading-snug">
	              {{ p.name }}
	            </div>

	            <div class="mt-3 flex items-baseline justify-between">
              <div class="text-xs tracking-widest uppercase text-[color:var(--muted)]">
                Harga
              </div>
              <div class="font-semibold">
                {{ formatIDR(p.price) }}
              </div>
            </div>
          </button>
        </div>
      </section>

      <aside class="lg:sticky lg:top-24">
        <div class="rounded-2xl border border-black/10 bg-[color:var(--paper)] backdrop-blur shadow-[0_18px_60px_rgba(0,0,0,0.10)]">
          <div class="flex items-center justify-between px-5 py-4">
            <div>
              <div class="font-brand text-lg">Keranjang</div>
              <div class="mt-1 text-xs tracking-widest uppercase text-[color:var(--muted)]">
                {{ itemsCount }} item
              </div>
            </div>
            <button
              type="button"
              class="rounded-xl border border-black/10 bg-white/70 px-3 py-2 text-xs font-semibold hover:bg-white disabled:opacity-50"
              :disabled="cart.length === 0"
              @click="clearCart"
            >
              Bersihkan
            </button>
          </div>

          <div class="max-h-[55vh] overflow-auto px-5 pb-2">
            <div
              v-if="cart.length === 0"
              class="rounded-xl border border-dashed border-black/15 bg-white/40 px-4 py-6 text-center"
            >
              <div class="font-semibold">Keranjang kosong</div>
              <div class="mt-1 text-sm text-[color:var(--muted)]">
                Tambah menu dari daftar di kiri.
              </div>
            </div>

            <div v-else class="space-y-3">
              <div
                v-for="item in cart"
                :key="item.id"
                class="rounded-xl border border-black/10 bg-white/60 px-4 py-3"
              >
                <div class="flex items-start justify-between gap-3">
                  <div class="min-w-0">
                    <div class="truncate font-semibold">{{ item.name }}</div>
                    <div class="mt-1 text-xs text-[color:var(--muted)]">
                      {{ formatIDR(item.price) }} / item
                    </div>
                  </div>

                  <button
                    type="button"
                    class="rounded-lg px-2 py-1 text-sm font-semibold text-[color:var(--danger)] hover:bg-red-50"
                    @click="remove(item.id)"
                  >
                    Hapus
                  </button>
                </div>

                <div class="mt-3 flex items-center justify-between gap-3">
                  <div class="flex items-center gap-2">
                    <button
                      type="button"
                      class="h-9 w-9 rounded-xl border border-black/10 bg-white/80 text-lg font-semibold hover:bg-white"
                      @click="decrease(item)"
                    >
                      -
                    </button>
                    <div class="min-w-10 text-center font-semibold">
                      {{ item.qty }}
                    </div>
                    <button
                      type="button"
                      class="h-9 w-9 rounded-xl border border-black/10 bg-white/80 text-lg font-semibold hover:bg-white"
                      @click="increase(item)"
                    >
                      +
                    </button>
                  </div>

                  <div class="text-right">
                    <div class="text-xs tracking-widest uppercase text-[color:var(--muted)]">
                      Subtotal
                    </div>
                    <div class="font-semibold">
                      {{ formatIDR(item.qty * item.price) }}
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <div class="border-t border-black/10 px-5 py-4">
            <div class="flex items-center justify-between">
              <div class="text-sm text-[color:var(--muted)]">Total</div>
              <div class="font-brand text-2xl">
                {{ formatIDR(total) }}
              </div>
            </div>

            <button
              type="button"
              class="mt-4 w-full rounded-xl bg-[color:var(--accent)] px-4 py-3 text-base font-semibold text-white shadow-[0_14px_30px_rgba(193,122,59,0.35)] transition hover:brightness-95 disabled:cursor-not-allowed disabled:opacity-60"
              :disabled="cart.length === 0"
              @click="openPay"
            >
              Bayar
            </button>

            <div class="mt-2 text-center text-xs text-[color:var(--muted)]">
              Order tersimpan setelah konfirmasi pembayaran.
            </div>
          </div>
        </div>
      </aside>
    </main>

    <div
      v-if="isPayOpen"
      class="fixed inset-0 z-50 flex items-end justify-center p-4 sm:items-center"
    >
      <div class="absolute inset-0 bg-black/35" @click="closePay" />
      <div
        class="relative w-full max-w-lg animate-float-in rounded-2xl border border-black/10 bg-[color:var(--paper-strong)] backdrop-blur p-5 shadow-[0_24px_80px_rgba(0,0,0,0.18)]"
      >
        <div class="flex items-start justify-between gap-3">
          <div>
            <div class="text-xs tracking-widest uppercase text-[color:var(--muted)]">
              Pembayaran
            </div>
            <div class="mt-1 font-brand text-2xl">{{ formatIDR(total) }}</div>
          </div>
          <button
            type="button"
            class="rounded-xl border border-black/10 bg-white/70 px-3 py-2 text-sm font-semibold hover:bg-white"
            @click="closePay"
          >
            Tutup
          </button>
        </div>

	        <div class="mt-4 grid grid-cols-2 gap-2">
	          <button
	            type="button"
	            class="rounded-xl border px-4 py-2 text-sm font-semibold transition"
            :class="
              payMethod === 'cash'
                ? 'border-black/15 bg-white shadow-[0_10px_25px_rgba(0,0,0,0.10)]'
                : 'border-black/10 bg-white/60 hover:bg-white/80'
            "
            @click="payMethod = 'cash'"
          >
            Cash
          </button>
          <button
            type="button"
            class="rounded-xl border px-4 py-2 text-sm font-semibold transition"
            :class="
              payMethod === 'qris'
                ? 'border-black/15 bg-white shadow-[0_10px_25px_rgba(0,0,0,0.10)]'
                : 'border-black/10 bg-white/60 hover:bg-white/80'
            "
            @click="payMethod = 'qris'"
          >
            QRIS
	          </button>
	        </div>

	        <div class="mt-4 rounded-2xl border border-black/10 bg-white/70 p-4">
	          <div class="text-sm font-semibold">Detail pesanan</div>

	          <div class="mt-3">
	            <div class="text-xs font-medium text-[color:var(--muted)]">Tipe</div>
	            <div class="mt-1 grid grid-cols-2 gap-2">
	              <button
	                type="button"
	                class="rounded-xl border px-4 py-2 text-sm font-semibold transition"
	                :class="
	                  orderType === 'dine_in'
	                    ? 'border-black/15 bg-white shadow-[0_10px_25px_rgba(0,0,0,0.10)]'
	                    : 'border-black/10 bg-white/60 hover:bg-white/80'
	                "
	                @click="setOrderType('dine_in')"
	              >
	                Dine In
	              </button>
	              <button
	                type="button"
	                class="rounded-xl border px-4 py-2 text-sm font-semibold transition"
	                :class="
	                  orderType === 'take_away'
	                    ? 'border-black/15 bg-white shadow-[0_10px_25px_rgba(0,0,0,0.10)]'
	                    : 'border-black/10 bg-white/60 hover:bg-white/80'
	                "
	                @click="setOrderType('take_away')"
	              >
	                Take Away
	              </button>
	            </div>
	          </div>

	          <div v-if="orderType === 'dine_in'" class="mt-3 grid grid-cols-2 gap-2">
	            <label class="block">
	              <span class="text-xs font-medium text-[color:var(--muted)]">Meja</span>
	              <input
	                v-model.trim="tableNo"
	                inputmode="numeric"
	                placeholder="01"
	                class="mt-1 w-full rounded-xl border border-black/10 bg-white px-3 py-2 text-sm outline-none focus:border-black/20 focus:shadow-[0_0_0_4px_rgba(193,122,59,0.15)]"
	              />
	            </label>
	            <label class="block">
	              <span class="text-xs font-medium text-[color:var(--muted)]">Tamu</span>
	              <input
	                v-model.trim="guestCount"
	                inputmode="numeric"
	                placeholder="3"
	                class="mt-1 w-full rounded-xl border border-black/10 bg-white px-3 py-2 text-sm outline-none focus:border-black/20 focus:shadow-[0_0_0_4px_rgba(193,122,59,0.15)]"
	              />
	            </label>
	          </div>
	          <div v-else class="mt-3 text-sm text-[color:var(--muted)]">
	            Take away (tanpa meja).
	          </div>

	          <label class="mt-2 block">
	            <span class="text-xs font-medium text-[color:var(--muted)]">Nama pemesan</span>
	            <input
	              v-model.trim="customerName"
	              placeholder="Budi"
	              class="mt-1 w-full rounded-xl border border-black/10 bg-white px-3 py-2 text-sm outline-none focus:border-black/20 focus:shadow-[0_0_0_4px_rgba(193,122,59,0.15)]"
	            />
	          </label>
	        </div>

	        <div
	          v-if="payMethod === 'cash'"
	          class="mt-4 rounded-2xl border border-black/10 bg-white/70 p-4"
        >
          <label class="block">
            <span class="text-sm font-medium">Uang diterima</span>
            <input
              v-model="cashReceived"
              inputmode="numeric"
              placeholder="Contoh: 50000"
              class="mt-1 w-full rounded-xl border border-black/10 bg-white px-4 py-3 text-sm outline-none focus:border-black/20 focus:shadow-[0_0_0_4px_rgba(193,122,59,0.15)]"
            />
          </label>
          <div class="mt-3 flex items-center justify-between text-sm">
            <div class="text-[color:var(--muted)]">Kembalian</div>
            <div
              class="font-semibold"
              :class="change != null && change < 0 ? 'text-[color:var(--danger)]' : ''"
            >
              {{ change == null ? '-' : formatIDR(change) }}
            </div>
          </div>
        </div>

        <div
          v-else
          class="mt-4 rounded-2xl border border-black/10 bg-white/70 p-4"
        >
          <div class="text-sm font-semibold">QRIS</div>
          <div class="mt-1 text-sm text-[color:var(--muted)]">
            Scan QR di perangkat kasir, lalu klik konfirmasi.
          </div>
        </div>

        <button
          type="button"
          class="mt-4 w-full rounded-xl bg-[color:var(--accent)] px-4 py-3 font-semibold text-white shadow-[0_14px_30px_rgba(193,122,59,0.35)] transition hover:brightness-95 disabled:cursor-not-allowed disabled:opacity-60"
          :disabled="submitting || (payMethod === 'cash' && !canConfirmCash)"
          @click="confirmPay"
        >
          <span v-if="!submitting">Konfirmasi pembayaran</span>
          <span v-else>Memproses...</span>
	        </button>
	      </div>
	    </div>

	    <div
	      v-if="isDoneOpen && lastOrder"
	      class="fixed inset-0 z-50 flex items-end justify-center p-4 sm:items-center"
	    >
	      <div class="absolute inset-0 bg-black/35" @click="closeDone" />
	      <div
	        class="relative w-full max-w-lg animate-float-in rounded-2xl border border-black/10 bg-[color:var(--paper-strong)] backdrop-blur p-5 shadow-[0_24px_80px_rgba(0,0,0,0.18)]"
	      >
	        <div class="flex items-start justify-between gap-3">
	          <div>
	            <div class="text-xs tracking-widest uppercase text-[color:var(--muted)]">
	              Transaksi berhasil
	            </div>
		            <div class="mt-1 font-brand text-2xl">{{ formatIDR(lastOrder.total) }}</div>
		            <div class="mt-1 text-sm text-[color:var(--muted)]">
		              Order #{{ lastOrder.orderNo || shortOrderId(lastOrder.orderId) }}
		            </div>
		            <div class="mt-1 text-sm text-[color:var(--muted)]">
		              {{ lastOrder.orderType === 'take_away' ? 'Take Away' : 'Dine In' }}
		            </div>
		            <div
		              v-if="lastOrder.tableNo || lastOrder.guestCount || lastOrder.customerName"
		              class="mt-2 text-sm text-[color:var(--muted)]"
		            >
		              <div v-if="lastOrder.tableNo || lastOrder.guestCount">
		                <span v-if="lastOrder.tableNo">Meja {{ lastOrder.tableNo }}</span>
		                <span v-if="lastOrder.guestCount">, {{ lastOrder.guestCount }} Tamu</span>
		              </div>
		              <div v-if="lastOrder.customerName">Nama: {{ lastOrder.customerName }}</div>
		            </div>
		          </div>
	          <button
	            type="button"
	            class="rounded-xl border border-black/10 bg-white/70 px-3 py-2 text-sm font-semibold hover:bg-white"
	            @click="closeDone"
	          >
	            Tutup
	          </button>
	        </div>

	        <div class="mt-4 rounded-2xl border border-black/10 bg-white/70 p-4 text-sm">
	          <div class="flex items-center justify-between">
	            <div class="text-[color:var(--muted)]">Metode</div>
	            <div class="font-semibold">
	              {{ lastOrder.payMethod === 'cash' ? 'Cash' : 'QRIS' }}
	            </div>
	          </div>
	          <div v-if="lastOrder.payMethod === 'cash'" class="mt-2 space-y-1">
	            <div class="flex items-center justify-between">
	              <div class="text-[color:var(--muted)]">Diterima</div>
	              <div class="font-semibold">{{ formatIDR(lastOrder.received) }}</div>
	            </div>
	            <div class="flex items-center justify-between">
	              <div class="text-[color:var(--muted)]">Kembalian</div>
	              <div class="font-semibold">{{ formatIDR(lastOrder.change) }}</div>
	            </div>
	          </div>
	          <div class="mt-2 flex items-center justify-between">
	            <div class="text-[color:var(--muted)]">Waktu</div>
	            <div class="font-semibold">{{ formatTime(lastOrder.paidAt) }}</div>
	          </div>
	        </div>

	        <div class="mt-4 grid grid-cols-1 gap-2 sm:grid-cols-3">
	          <button
	            type="button"
	            class="rounded-xl bg-[color:var(--accent)] px-4 py-3 text-sm font-semibold text-white shadow-[0_14px_30px_rgba(193,122,59,0.30)] transition hover:brightness-95 disabled:cursor-not-allowed disabled:opacity-60"
	            :disabled="printing"
	            @click="printReceipt"
	          >
	            Cetak struk
	          </button>
	          <button
	            type="button"
	            class="rounded-xl border border-black/10 bg-white/70 px-4 py-3 text-sm font-semibold transition hover:bg-white disabled:cursor-not-allowed disabled:opacity-60"
	            :disabled="printing"
	            @click="printKitchen"
	          >
	            Cetak kitchen
	          </button>
	          <button
	            type="button"
	            class="rounded-xl border border-black/10 bg-white/70 px-4 py-3 text-sm font-semibold transition hover:bg-white disabled:cursor-not-allowed disabled:opacity-60"
	            :disabled="printing"
	            @click="printBoth"
	          >
	            Cetak 2x
	          </button>
	        </div>

	        <button
	          type="button"
	          class="mt-3 w-full rounded-xl border border-black/10 bg-white/70 px-4 py-3 text-sm font-semibold transition hover:bg-white disabled:cursor-not-allowed disabled:opacity-60"
	          :disabled="printing"
	          @click="newSale"
	        >
	          Transaksi baru
	        </button>
	      </div>
	    </div>

	    <div
	      v-if="toast"
	      class="fixed bottom-4 left-1/2 z-50 w-[min(560px,calc(100%-2rem))] -translate-x-1/2"
	    >
      <div
        class="rounded-2xl border border-black/10 bg-[color:var(--paper-strong)] px-4 py-3 shadow-[0_18px_60px_rgba(0,0,0,0.14)]"
        :class="
          toast.type === 'error'
            ? 'border-red-200'
            : toast.type === 'success'
              ? 'border-emerald-200'
              : ''
        "
      >
        <div class="flex items-start justify-between gap-3">
          <div class="text-sm font-medium">{{ toast.message }}</div>
          <button
            type="button"
            class="rounded-lg px-2 py-1 text-sm font-semibold text-[color:var(--muted)] hover:bg-black/5"
            @click="toast = null"
          >
            Tutup
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
	import { computed, onMounted, ref } from 'vue'
	import { useRouter } from 'vue-router'
	import { useAuthStore } from '../stores/auth'
	import { useSettingsStore } from '../stores/settings'
	import { ApiError, createOrder, getProducts } from '../api/api'
	import { receiptConfig } from '../config/receipt'

const auth = useAuthStore()
const settings = useSettingsStore()
const router = useRouter()

const products = ref([])
const loadingProducts = ref(true)
const productsError = ref('')

	const cart = ref([])
	const search = ref('')

	const customerName = ref('')
	const tableNo = ref('')
	const guestCount = ref('')
	const orderType = ref('dine_in')

		const isPayOpen = ref(false)
		const payMethod = ref('cash')
		const cashReceived = ref('')
		const submitting = ref(false)
	const isDoneOpen = ref(false)
	const printing = ref(false)
	const lastOrder = ref(null)

	const toast = ref(null)
	let toastTimer = null

function showToast(type, message) {
  toast.value = { type, message }
  if (toastTimer) clearTimeout(toastTimer)
  toastTimer = setTimeout(() => {
    toast.value = null
  }, 3200)
}

	function formatIDR(num) {
	  const amount = Number(num || 0)
	  return new Intl.NumberFormat('id-ID', {
	    style: 'currency',
    currency: 'IDR',
    maximumFractionDigits: 0
	  }).format(amount)
	}

	function formatTime(isoString) {
	  const d = isoString ? new Date(isoString) : new Date()
	  return d.toLocaleString('id-ID')
	}

	function formatRpShort(num) {
	  return new Intl.NumberFormat('id-ID').format(Number(num || 0))
	}

	function escapeHtml(s) {
	  return String(s || '')
	    .replace(/&/g, '&amp;')
	    .replace(/</g, '&lt;')
	    .replace(/>/g, '&gt;')
	    .replace(/"/g, '&quot;')
	    .replace(/'/g, '&#39;')
	}

		function shortOrderId(orderId) {
		  if (!orderId) return '-'
		  return String(orderId).split('-')[0].toUpperCase()
		}

		function orderNoFromId(orderId) {
		  if (!orderId) return '-----'
		  const hex = String(orderId).replace(/-/g, '').slice(0, 8)
		  const n = Number.parseInt(hex, 16)
		  if (!Number.isFinite(n)) return shortOrderId(orderId)
		  return String(n % 100000).padStart(5, '0')
		}

		function normalizeTableNo(value) {
		  const raw = String(value || '').trim()
		  if (!raw) return null
		  const digits = raw.replace(/[^\d]/g, '')
		  if (!digits) return raw.toUpperCase()
		  return digits.padStart(2, '0')
		}

		function setOrderType(type) {
		  orderType.value = type
		  if (type === 'take_away') {
		    tableNo.value = ''
		    guestCount.value = ''
		  }
		}

		function printDoc(title, html, printCfg) {
		  return new Promise(resolve => {
		    if (!document?.body) {
		      showToast('error', 'Gagal memulai cetak')
		      resolve()
		      return
		    }

		    const cfg = printCfg || receiptConfig.printer.receipt
		    const paperWidthMm = Number(cfg.paperWidthMm || 58)
		    const contentWidthMm = Number(cfg.contentWidthMm || paperWidthMm)
		    const paddingMm = Number(cfg.paddingMm || 0)
		    const baseFontPx = Number(cfg.baseFontPx || 11)
		    const titleFontPx = Number(cfg.titleFontPx || 16)

		    const iframe = document.createElement('iframe')
		    iframe.setAttribute('aria-hidden', 'true')
		    iframe.tabIndex = -1
		    iframe.style.position = 'fixed'
		    iframe.style.right = '0'
		    iframe.style.bottom = '0'
		    iframe.style.width = '0'
		    iframe.style.height = '0'
		    iframe.style.border = '0'
		    iframe.style.opacity = '0'
		    iframe.style.pointerEvents = 'none'

		    let cleaned = false
		    const cleanup = () => {
		      if (cleaned) return
		      cleaned = true
		      try {
		        const w = iframe.contentWindow
		        if (w) w.onafterprint = null
		      } catch {}
		      try {
		        iframe.remove()
		      } catch {}
		      resolve()
		    }

			    const docHtml = `<!doctype html>
	<html>
	  <head>
	    <meta charset="utf-8" />
	    <meta name="viewport" content="width=device-width,initial-scale=1" />
	    <title>${escapeHtml(title)}</title>
	    <style id="page-style">
	      @page { size: ${paperWidthMm}mm 200mm; margin: 0; }
	      * { box-sizing: border-box; }
	      html, body { height: 100%; }
	      body {
	        margin: 0;
	        padding: 0;
	        font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, "Liberation Mono", "Courier New", monospace;
	        color: #111;
	        background: #fff;
	        -webkit-print-color-adjust: exact;
	        print-color-adjust: exact;
	        font-size: ${baseFontPx}px;
	      }
	      .paper { width: ${contentWidthMm}mm; margin: 0 auto; padding: ${paddingMm}mm; }
	      img { max-width: 100%; height: auto; }
	      .logo { display: block; margin: 0 auto; max-width: 80%; max-height: 72px; object-fit: contain; }
	      .row { display: flex; justify-content: space-between; gap: 8px; }
	      .muted { color: rgba(0,0,0,0.62); }
	      .hr { border-top: 1px dashed rgba(0,0,0,0.35); margin: 10px 0; }
	      .title { font-size: ${titleFontPx}px; font-weight: 800; letter-spacing: 1px; text-align: center; }
	      .small { font-size: 11px; line-height: 1.35; }
      .items { margin-top: 8px; }
      .item { display: flex; justify-content: space-between; gap: 10px; margin: 6px 0; }
      .name { flex: 1; word-break: break-word; }
      .qty { width: 34px; text-align: right; }
      .money { text-align: right; white-space: nowrap; }
      .stamp { display:inline-block; border: 2px solid rgba(0,0,0,0.55); padding: 2px 8px; border-radius: 999px; font-weight: 800; letter-spacing: 2px; }
    </style>
  </head>
  <body>
    <div class="paper">${html}</div>
  </body>
</html>`

		    iframe.onload = () => {
		      const win = iframe.contentWindow
		      if (!win) {
		        showToast('error', 'Gagal membuka printer')
		        cleanup()
		        return
		      }

		      win.onafterprint = cleanup

		      const doc = iframe.contentDocument
		      const paper = doc?.querySelector?.('.paper')
		      const styleEl = doc?.getElementById?.('page-style')

		      // Auto-fit page height to content to avoid wasting blank paper on thermal printers.
		      // CSS pixel is 1/96 inch.
		      const pxToMm = px => (px * 25.4) / 96

		      const waitForImages = () => {
		        const imgs = Array.from(doc?.images || [])
		        if (!imgs.length) return Promise.resolve()

		        const pending = imgs
		          .filter(img => !img.complete)
		          .map(
		            img =>
		              new Promise(resolveImg => {
		                const done = () => resolveImg()
		                img.addEventListener('load', done, { once: true })
		                img.addEventListener('error', done, { once: true })
		              })
		          )

		        if (!pending.length) return Promise.resolve()

		        const timeoutMs = 1200
		        return Promise.race([
		          Promise.all(pending),
		          new Promise(res => setTimeout(res, timeoutMs))
		        ])
		      }

		      setTimeout(() => {
		        waitForImages().then(() => {
		          try {
		            if (paper && styleEl) {
		              const heightPx = Math.ceil(paper.scrollHeight || 0)
		              const heightMm = Math.max(60, Math.ceil(pxToMm(heightPx)) + paddingMm * 2)
		              styleEl.textContent = styleEl.textContent.replace(
		                /@page\\s*\\{\\s*size:[^;]+;/,
		                `@page { size: ${paperWidthMm}mm ${heightMm}mm;`
		              )
		            }

		            win.focus()
		            win.print()
		          } catch (e) {
		            showToast('error', 'Cetak gagal')
		            cleanup()
		            return
		          }

		          setTimeout(cleanup, 120000)
		        })
		      }, 180)
		    }

		    iframe.srcdoc = docHtml
		    document.body.appendChild(iframe)
		  })
		}

			function receiptHtml(order) {
			  const store = settings.store || receiptConfig.store || {}
			  const storeName = escapeHtml(store.name || 'WARKOP')
			  const storeTagline = escapeHtml(store.tagline || '')
			  const addressLines = (store.address_lines || store.addressLines || []).map(escapeHtml)
			  const phone = escapeHtml(store.phone || '')
			  const logoDataUrl = String(store.logo_data_url || '').trim()

			  const orderNo = escapeHtml(order.orderNo || orderNoFromId(order.orderId))
			  const time = escapeHtml(formatTime(order.paidAt))
			  const cashier = escapeHtml(order.cashierName || auth.user?.name || auth.user?.email || auth.user?.role || 'Kasir')
			  const orderType = order.orderType === 'take_away' ? 'take_away' : 'dine_in'
			  const serviceLabel = orderType === 'take_away' ? 'TAKE AWAY' : 'DINE IN'

			  const table = normalizeTableNo(order.tableNo)
			  const guests =
			    order.guestCount == null || order.guestCount === ''
			      ? null
			      : Number(order.guestCount)
			  const customer = String(order.customerName || '').trim()

			  const tableParts = []
			  if (table) tableParts.push(`Meja ${escapeHtml(table)}`)
			  if (Number.isFinite(guests) && guests > 0) tableParts.push(`${guests} Tamu`)
			  const tableLine = tableParts.join(', ')

			  const addressBlock =
			    receiptConfig.receipt?.showAddress && addressLines.length
			      ? addressLines
			          .map(l => `<div class="small muted" style="text-align:center;">${l}</div>`)
			          .join('')
			      : ''

			  const phoneBlock =
			    receiptConfig.receipt?.showPhone && phone
			      ? `<div class="small muted" style="text-align:center;">Telp: ${phone}</div>`
			      : ''

			  const logoBlock =
			    receiptConfig.receipt?.showLogo && logoDataUrl
			      ? `<div style="text-align:center;margin-bottom:6px;">
  <img class="logo" alt="" src="${escapeHtml(logoDataUrl)}" />
</div>`
			      : ''

			  const footerLines = (receiptConfig.receipt?.footerLines || []).map(escapeHtml)
			  const footerBlock = footerLines.length
			    ? footerLines
			        .map(l => `<div class="small muted" style="text-align:center;">${l}</div>`)
			        .join('')
			    : ''

			  const qtyTotal = order.items.reduce(
			    (sum, i) => sum + (Number(i.qty) || 0),
			    0
			  )
			  const subtotal = order.items.reduce(
			    (sum, i) => sum + (Number(i.qty) || 0) * (Number(i.price) || 0),
			    0
			  )

			  const taxCfg = receiptConfig.receipt?.tax || {}
			  const taxEnabled = !!taxCfg.enabled
			  const taxRate = Number(taxCfg.rate || 0)
			  const taxValue = taxEnabled ? Math.round(subtotal * taxRate) : 0
			  const grandTotal = taxEnabled ? subtotal + taxValue : subtotal

			  const totalToPrint =
			    taxEnabled ? grandTotal : Number(order.total || subtotal)

			  const received =
			    order.payMethod === 'cash' ? `Rp ${formatRpShort(order.received)}` : '-'
			  const changeValue =
			    order.payMethod === 'cash' ? `Rp ${formatRpShort(order.change)}` : '-'

			  const paidStamp =
			    order.payMethod === 'qris'
			      ? `<div style="text-align:center;margin-top:10px;"><span class="stamp">PAID</span></div>`
			      : ''

			  const itemsRows = order.items
			    .map(i => {
			      const name = escapeHtml(i.name)
			      const qty = Number(i.qty) || 0
			      const lineTotal = qty * (Number(i.price) || 0)
			      const unit = `Rp ${formatRpShort(i.price)}`
			      const unitLine = receiptConfig.receipt?.showItemUnitPrice
			        ? `<div class="row small muted"><div>${unit} / item</div><div></div></div>`
			        : ''

			      return `<div class="row small" style="margin:6px 0;">
  <div class="name">${name}</div>
  <div class="qty">${qty}</div>
  <div class="money" style="width:72px;">Rp ${formatRpShort(lineTotal)}</div>
</div>${unitLine}`
			    })
			    .join('')

			  const taxLine = taxEnabled
			    ? `<div class="row small"><div class="muted">${escapeHtml(taxCfg.label || 'Pajak')}</div><div>Rp ${formatRpShort(taxValue)}</div></div>`
			    : ''

			  const paymentLineLabel = order.payMethod === 'cash' ? 'Tunai' : 'Non-tunai'
			  const paymentLineValue = order.payMethod === 'cash' ? received : '-'
			  const changeLine =
			    order.payMethod === 'cash'
			      ? `<div class="row small"><div class="muted">Kembalian</div><div>${changeValue}</div></div>`
			      : ''

				  const tableBlock = tableLine
				    ? `<div class="small" style="text-align:center;font-weight:800;">${tableLine}</div>`
				    : ''
				  const serviceBlock = `<div class="small" style="text-align:center;font-weight:900;letter-spacing:2px;">${escapeHtml(serviceLabel)}</div>`
				  const customerBlock = customer
				    ? `<div class="small muted" style="text-align:center;">Nama: ${escapeHtml(customer)}</div>`
				    : ''

				  return `
	    ${logoBlock}
	    <div class="title">${storeName}</div>
	    ${storeTagline ? `<div class="small muted" style="text-align:center;margin-top:4px;">${storeTagline}</div>` : ''}
	    ${addressBlock}
	    ${phoneBlock}
	    <div class="hr"></div>
	    ${serviceBlock}
	    ${tableBlock}
	    ${customerBlock}
	    <div class="hr"></div>
	    <div class="row small"><div class="muted">Invoice</div><div>${orderNo}</div></div>
	    <div class="row small"><div class="muted">Pesan</div><div>${orderNo}</div></div>
	    <div class="row small"><div class="muted">Pegawai</div><div>${cashier}</div></div>
	    <div class="row small"><div class="muted">Jam</div><div>${time}</div></div>
    <div class="hr"></div>
    <div class="row small muted">
      <div class="name">Barang</div>
      <div class="qty">Qty</div>
      <div class="money" style="width:72px;">Harga</div>
    </div>
    <div class="items">${itemsRows}</div>
    <div class="hr"></div>
    <div class="row small"><div class="muted">Qty</div><div>${qtyTotal}</div></div>
    <div class="row small"><div class="muted">Subtotal</div><div>Rp ${formatRpShort(subtotal)}</div></div>
    ${taxLine}
    <div class="row small"><div class="muted">Total</div><div><b>Rp ${formatRpShort(totalToPrint)}</b></div></div>
    <div class="row small"><div class="muted">${paymentLineLabel}</div><div>${paymentLineValue}</div></div>
    ${changeLine}
    <div class="hr"></div>
    ${paidStamp}
    ${footerBlock}
  `
			}

				function kitchenHtml(order) {
				  const station = receiptConfig.kitchen?.stationName || 'Dapur 1'
				  const orderNo = escapeHtml(order.orderNo || orderNoFromId(order.orderId))
				  const time = escapeHtml(formatTime(order.paidAt))
				  const cashier = escapeHtml(order.cashierName || auth.user?.name || auth.user?.email || auth.user?.role || 'Kasir')
				  const orderType = order.orderType === 'take_away' ? 'take_away' : 'dine_in'
				  const serviceLabel = orderType === 'take_away' ? 'TAKE AWAY' : 'DINE IN'

				  const table = normalizeTableNo(order.tableNo)
				  const guests =
				    order.guestCount == null || order.guestCount === ''
			      ? null
			      : Number(order.guestCount)
			  const customer = String(order.customerName || '').trim()

			  const groups = { MAKANAN: [], MINUMAN: [] }
			  for (const item of order.items) {
			    const n = String(item?.name || '').toLowerCase()
			    const isDrink =
			      n.includes('kopi') ||
			      n.includes('teh') ||
			      n.includes('jeruk') ||
			      n.includes('air ') ||
			      n === 'air' ||
			      n.includes('soda') ||
			      n.includes('es ')
			    groups[isDrink ? 'MINUMAN' : 'MAKANAN'].push(item)
			  }

			  const sectionHtml = (title, items) => {
			    if (!items.length) return ''
			    const lines = items
			      .map(i => {
			        const qty = Number(i.qty) || 0
			        const name = escapeHtml(i.name)
			        return `<div style="margin:8px 0; font-size:16px; font-weight:800;">${qty}X ${name}</div>`
			      })
			      .join('')
			    return `
      <div style="text-align:center;font-weight:900;letter-spacing:2px;margin-top:8px;">--${escapeHtml(title)}--</div>
      ${lines}
    `
			  }

			  const footerLines = (receiptConfig.kitchen?.footerLines || []).map(escapeHtml)
				  const footerBlock = footerLines.length
				    ? footerLines
				        .map(l => `<div class="small muted" style="text-align:center;">${l}</div>`)
				        .join('')
				    : ''

				  const headerLines = []
				  headerLines.push(`<div style="text-align:center;font-size:16px;font-weight:900;letter-spacing:2px;margin-top:4px;">${escapeHtml(serviceLabel)}</div>`)
				  if (table) headerLines.push(`<div style="text-align:center;font-size:20px;font-weight:900;margin-top:6px;">Meja ${escapeHtml(table)}</div>`)
				  if (Number.isFinite(guests) && guests > 0) headerLines.push(`<div style="text-align:center;font-size:16px;font-weight:800;">${guests} Tamu</div>`)
				  if (customer) headerLines.push(`<div class="small muted" style="text-align:center;margin-top:2px;">Nama: ${escapeHtml(customer)}</div>`)

			  return `
    <div class="title">*** ${escapeHtml(station)} ***</div>
    ${headerLines.join('')}
    <div style="text-align:center;font-size:34px;font-weight:900;letter-spacing:6px;margin-top:8px;">${orderNo}</div>
    <div class="hr"></div>
    <div class="row small"><div class="muted">Pegawai</div><div>${cashier}</div></div>
    <div class="row small"><div class="muted">Jam</div><div>${time}</div></div>
    <div class="hr"></div>
    ${sectionHtml('MAKANAN', groups.MAKANAN)}
    ${sectionHtml('MINUMAN', groups.MINUMAN)}
    <div class="hr"></div>
    ${footerBlock}
  `
			}

		async function printReceipt() {
		  if (!lastOrder.value) return
		  printing.value = true
		  try {
		    await printDoc('Struk', receiptHtml(lastOrder.value), receiptConfig.printer.receipt)
		  } finally {
		    printing.value = false
		  }
		}

		async function printKitchen() {
		  if (!lastOrder.value) return
		  printing.value = true
		  try {
		    await printDoc('Kitchen', kitchenHtml(lastOrder.value), receiptConfig.printer.kitchen)
		  } finally {
		    printing.value = false
		  }
		}

		async function printBoth() {
		  if (!lastOrder.value) return
		  printing.value = true
		  try {
		    await printDoc('Kitchen', kitchenHtml(lastOrder.value), receiptConfig.printer.kitchen)
		    await printDoc('Struk', receiptHtml(lastOrder.value), receiptConfig.printer.receipt)
		  } finally {
		    printing.value = false
		  }
		}

	function closeDone() {
	  if (printing.value) return
	  isDoneOpen.value = false
	}

		function newSale() {
		  closeDone()
		  lastOrder.value = null
		  customerName.value = ''
		  tableNo.value = ''
		  guestCount.value = ''
		  orderType.value = 'dine_in'
		  payMethod.value = 'cash'
		  cashReceived.value = ''
		  showToast('success', 'Siap transaksi baru')
		}

	function menuImageFor(name) {
	  const n = String(name || '').toLowerCase()
	  if (n.includes('kopi')) return '/menu/coffee.svg'
	  if (n.includes('teh')) return '/menu/tea.svg'
	  if (n.includes('indomie') || n.includes('mie')) return '/menu/noodle.svg'
	  if (n.includes('roti')) return '/menu/toast.svg'
	  if (n.includes('nasi')) return '/menu/rice.svg'
	  if (n.includes('air') || n.includes('jeruk') || n.includes('soda')) return '/menu/drink.svg'
	  if (n.includes('pisang') || n.includes('tahu') || n.includes('tempe') || n.includes('sate')) return '/menu/snack.svg'
	  return '/menu/_default.svg'
	}

	function onMenuImgError(e) {
	  const img = e?.target
	  if (img && img.src && !img.src.endsWith('/menu/_default.svg')) {
	    img.src = '/menu/_default.svg'
	  }
	}

	async function loadProducts() {
	  loadingProducts.value = true
	  productsError.value = ''
	  try {
    const list = await getProducts(auth.token)
    products.value = Array.isArray(list) ? list : []
  } catch (e) {
    if (e instanceof ApiError && e.status === 401) {
      auth.logout()
      router.push('/')
      return
    }
    productsError.value = e?.message || 'Unknown error'
  } finally {
    loadingProducts.value = false
  }
}

	onMounted(async () => {
	  try {
	    if (auth.token && !auth.user) await auth.loadMe()
	  } catch {
	    router.push('/')
	    return
	  }

	  await settings.loadStore(auth.token)
	  await loadProducts()
	})

const filteredProducts = computed(() => {
  const q = search.value.trim().toLowerCase()
  if (!q) return products.value
  return products.value.filter(p => String(p?.name || '').toLowerCase().includes(q))
})

function addToCart(product) {
  if (!product?.id) return

  const existing = cart.value.find(i => i.id === product.id)

  if (existing) {
    existing.qty++
  } else {
    cart.value.push({
      id: product.id,
      name: product.name,
      price: Number(product.price || 0),
      qty: 1
    })
  }
}

function increase(item) {
  item.qty++
}

function decrease(item) {
  if (item.qty > 1) {
    item.qty--
  } else {
    remove(item.id)
  }
}

function remove(id) {
  cart.value = cart.value.filter(i => i.id !== id)
}

function clearCart() {
  cart.value = []
}

const itemsCount = computed(() =>
  cart.value.reduce((sum, i) => sum + (Number(i.qty) || 0), 0)
)

const total = computed(() =>
  cart.value.reduce(
    (sum, i) => sum + (Number(i.price) || 0) * (Number(i.qty) || 0),
    0
  )
)

const receivedCashNumber = computed(() => {
  const raw = String(cashReceived.value || '').replace(/[^\d]/g, '')
  if (raw === '') return null
  return Number(raw)
})

const change = computed(() => {
  if (payMethod.value !== 'cash') return null
  const received = receivedCashNumber.value
  if (received == null) return null
  return received - total.value
})

const canConfirmCash = computed(() => {
  if (payMethod.value !== 'cash') return true
  const received = receivedCashNumber.value
  return received != null && received >= total.value && total.value > 0
})

function openPay() {
  if (cart.value.length === 0) {
    showToast('error', 'Keranjang masih kosong')
    return
  }
  isPayOpen.value = true
  payMethod.value = 'cash'
  cashReceived.value = ''
}

function closePay() {
  if (submitting.value) return
  isPayOpen.value = false
}

function handleLogout() {
  auth.logout()
  router.push('/')
}

	async function confirmPay() {
	  if (cart.value.length === 0) {
	    showToast('error', 'Keranjang masih kosong')
	    return
	  }

		  const payload = {
		    items: cart.value.map(i => ({
		      product_id: i.id,
		      qty: i.qty,
		      price: i.price
		    })),
		    order_type: orderType.value,
		    table_no: orderType.value === 'dine_in' ? normalizeTableNo(tableNo.value) : null,
		    guest_count:
		      orderType.value === 'dine_in'
		        ? (String(guestCount.value || '').replace(/[^\d]/g, '') ? Number(String(guestCount.value || '').replace(/[^\d]/g, '')) : null)
		        : null,
		    customer_name: customerName.value.trim() || null,
		    payment_method: payMethod.value,
		    received: payMethod.value === 'cash' ? receivedCashNumber.value : null
		  }

	  submitting.value = true
	  try {
	    const snapshotItems = cart.value.map(i => ({
	      id: i.id,
	      name: i.name,
	      qty: Number(i.qty) || 0,
	      price: Number(i.price) || 0
	    }))
	    const snapshotTime = new Date().toISOString()

	    const snapshotReceived =
	      payMethod.value === 'cash' ? receivedCashNumber.value : null

		    const snapshotCustomer = customerName.value.trim() || null
		    const snapshotTable = normalizeTableNo(tableNo.value)
		    const guestsRaw = String(guestCount.value || '').replace(/[^\d]/g, '')
		    const snapshotGuests = guestsRaw ? Number(guestsRaw) : null
		    const snapshotOrderType = orderType.value === 'take_away' ? 'take_away' : 'dine_in'

		    const snapshotCashier =
		      auth.user?.name || auth.user?.email || auth.user?.role || 'Kasir'

	    const res = await createOrder(auth.token, payload)
	    const orderId = res?.order_id || null
	    const orderNo = orderNoFromId(orderId)

	    const apiTotal = Number(res?.total)
	    const finalTotal = Number.isFinite(apiTotal) && apiTotal > 0 ? apiTotal : total.value
	    const finalChange =
	      payMethod.value === 'cash' && snapshotReceived != null
	        ? snapshotReceived - finalTotal
	        : null

	    showToast('success', 'Transaksi tersimpan')
	    lastOrder.value = {
	      orderId,
	      orderNo,
	      paidAt: snapshotTime,
	      items: snapshotItems,
	      total: finalTotal,
	      payMethod: payMethod.value,
	      received: snapshotReceived,
	      change: finalChange,
		      customerName: snapshotCustomer,
		      tableNo: snapshotTable,
		      guestCount: snapshotGuests,
		      orderType: snapshotOrderType,
		      cashierName: snapshotCashier
		    }

	    cart.value = []
	    isPayOpen.value = false
	    cashReceived.value = ''
	    isDoneOpen.value = true
	  } catch (e) {
	    if (e instanceof ApiError && e.status === 401) {
	      auth.logout()
	      router.push('/')
	      return
	    }
	    showToast('error', e?.message || 'Checkout gagal')
	  } finally {
	    submitting.value = false
	  }
	}
</script>
