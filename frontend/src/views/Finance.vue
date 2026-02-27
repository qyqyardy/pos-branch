<template>
  <div class="min-h-screen">
    <TopBar />

    <main class="mx-auto max-w-7xl px-4 py-6 space-y-6">
	      <section class="rounded-2xl border border-black/10 bg-[color:var(--paper)] backdrop-blur shadow-[0_18px_60px_rgba(0,0,0,0.10)]">
	        <div class="flex flex-wrap items-end justify-between gap-4 px-5 py-4">
	          <div>
	            <div class="font-brand text-xl">Finance</div>
            <div class="mt-1 text-sm text-[color:var(--muted)]">
              Trace transaksi harian dan pembukuan sederhana.
            </div>
          </div>

          <div class="flex items-end gap-2">
            <label class="block">
              <span class="text-xs font-medium text-[color:var(--muted)]">Tanggal</span>
              <input
                v-model="selectedDate"
                type="date"
                class="mt-1 rounded-xl border border-black/10 bg-white px-3 py-2 text-sm outline-none focus:border-black/20 focus:shadow-[0_0_0_4px_rgba(193,122,59,0.15)]"
                @change="loadAll"
              />
            </label>
            <button
              type="button"
              class="rounded-xl border border-black/10 bg-white/70 px-4 py-2 text-sm font-semibold hover:bg-white disabled:opacity-60"
              :disabled="loading"
              @click="loadAll"
            >
              Refresh
            </button>
          </div>
        </div>

        <div class="border-t border-black/10 px-5 py-5">
          <div v-if="loading" class="text-sm text-[color:var(--muted)]">
            Memuat data...
          </div>

          <div v-else-if="error" class="rounded-xl border border-black/10 bg-white/70 p-4">
            <div class="font-semibold">Gagal memuat</div>
            <div class="mt-1 text-sm text-[color:var(--muted)]">{{ error }}</div>
          </div>

	          <div v-else class="grid grid-cols-2 gap-3 md:grid-cols-6">
            <div class="rounded-2xl border border-black/10 bg-white/70 p-4">
              <div class="text-xs tracking-widest uppercase text-[color:var(--muted)]">Sales</div>
              <div class="mt-2 font-brand text-xl">{{ formatIDR(summary.salesTotal) }}</div>
            </div>
            <div class="rounded-2xl border border-black/10 bg-white/70 p-4">
              <div class="text-xs tracking-widest uppercase text-[color:var(--muted)]">Orders</div>
              <div class="mt-2 font-brand text-xl">{{ summary.ordersCount }}</div>
            </div>
            <div class="rounded-2xl border border-black/10 bg-white/70 p-4">
              <div class="text-xs tracking-widest uppercase text-[color:var(--muted)]">Cash</div>
              <div class="mt-2 font-brand text-xl">{{ formatIDR(summary.cashTotal) }}</div>
            </div>
            <div class="rounded-2xl border border-black/10 bg-white/70 p-4">
              <div class="text-xs tracking-widest uppercase text-[color:var(--muted)]">QRIS</div>
              <div class="mt-2 font-brand text-xl">{{ formatIDR(summary.qrisTotal) }}</div>
            </div>
            <div class="rounded-2xl border border-black/10 bg-white/70 p-4">
              <div class="text-xs tracking-widest uppercase text-[color:var(--muted)]">Expense</div>
              <div class="mt-2 font-brand text-xl">{{ formatIDR(summary.expenseTotal) }}</div>
            </div>
            <div class="rounded-2xl border border-black/10 bg-white/70 p-4">
              <div class="text-xs tracking-widest uppercase text-[color:var(--muted)]">Net</div>
              <div class="mt-2 font-brand text-xl">{{ formatIDR(summary.netTotal) }}</div>
            </div>
          </div>

          <div v-if="!loading && !error" class="mt-8 grid grid-cols-1 gap-6 lg:grid-cols-2">
            <!-- Revenue Trend Chart -->
            <div class="rounded-2xl border border-black/10 bg-white/70 p-5">
              <div class="mb-4 flex items-center justify-between">
                <div class="font-brand text-lg">Tren Omzet (7 Hari)</div>
              </div>
              <div class="h-[280px]">
                <Line
                  v-if="chartData.sales.labels.length"
                  :data="chartData.sales"
                  :options="chartOptions.sales"
                />
                <div v-else class="flex h-full items-center justify-center text-sm text-[color:var(--muted)]">
                  Memuat grafik...
                </div>
              </div>
            </div>

            <!-- Top Products Chart -->
            <div class="rounded-2xl border border-black/10 bg-white/70 p-5">
              <div class="mb-4 flex items-center justify-between">
                <div class="font-brand text-lg">5 Produk Terlaris (30 Hari)</div>
              </div>
              <div class="h-[280px]">
                <Bar
                  v-if="chartData.topProducts.labels.length"
                  :data="chartData.topProducts"
                  :options="chartOptions.topProducts"
                />
                <div v-else class="flex h-full items-center justify-center text-sm text-[color:var(--muted)]">
                  Memuat data produk...
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>

	      <section class="rounded-2xl border border-black/10 bg-[color:var(--paper)] backdrop-blur shadow-[0_18px_60px_rgba(0,0,0,0.10)]">
	        <div class="px-5 py-4">
	          <div class="font-brand text-xl">Neraca Harian (Sederhana)</div>
	          <div class="mt-1 text-sm text-[color:var(--muted)]">
	            Fokus ke saldo kas dan bank/QRIS di tanggal terpilih.
	          </div>
	        </div>

	        <div class="border-t border-black/10 px-5 py-5 space-y-5">
	          <div class="grid grid-cols-1 gap-3 md:grid-cols-4">
	            <label class="block">
	              <span class="text-xs font-medium text-[color:var(--muted)]">Saldo awal kas</span>
	              <input
	                v-model="opening.cash"
	                inputmode="numeric"
	                placeholder="0"
	                class="mt-1 w-full rounded-xl border border-black/10 bg-white px-3 py-2 text-sm outline-none focus:border-black/20 focus:shadow-[0_0_0_4px_rgba(193,122,59,0.15)]"
	              />
	            </label>
	            <label class="block">
	              <span class="text-xs font-medium text-[color:var(--muted)]">Saldo awal bank/QRIS</span>
	              <input
	                v-model="opening.bank"
	                inputmode="numeric"
	                placeholder="0"
	                class="mt-1 w-full rounded-xl border border-black/10 bg-white px-3 py-2 text-sm outline-none focus:border-black/20 focus:shadow-[0_0_0_4px_rgba(193,122,59,0.15)]"
	              />
	            </label>
	            <div class="rounded-xl border border-black/10 bg-white/70 p-3 md:col-span-2">
	              <div class="text-xs tracking-widest uppercase text-[color:var(--muted)]">Check</div>
	              <div class="mt-2 flex flex-wrap items-baseline justify-between gap-2">
	                <div class="text-sm text-[color:var(--muted)]">Aktiva - Modal</div>
	                <div
	                  class="font-semibold"
	                  :class="balance.diff === 0 ? 'text-emerald-700' : 'text-[color:var(--danger)]'"
	                >
	                  {{ formatIDR(balance.diff) }}
	                </div>
	              </div>
	              <div class="mt-1 text-xs text-[color:var(--muted)]">
	                Idealnya 0. Kalau tidak 0, cek saldo awal atau catatan income/expense.
	              </div>
	            </div>
	          </div>

	          <div class="grid grid-cols-1 gap-4 md:grid-cols-2">
	            <div class="rounded-2xl border border-black/10 bg-white/70 p-4">
	              <div class="text-xs tracking-widest uppercase text-[color:var(--muted)]">Aktiva</div>
	              <div class="mt-3 space-y-2 text-sm">
	                <div class="flex items-center justify-between">
	                  <div class="text-[color:var(--muted)]">Kas (akhir)</div>
	                  <div class="font-semibold">{{ formatIDR(balance.cashClosing) }}</div>
	                </div>
	                <div class="flex items-center justify-between">
	                  <div class="text-[color:var(--muted)]">Bank/QRIS (akhir)</div>
	                  <div class="font-semibold">{{ formatIDR(balance.bankClosing) }}</div>
	                </div>
	                <div class="border-t border-black/10 pt-2 flex items-center justify-between">
	                  <div class="font-semibold">Total Aktiva</div>
	                  <div class="font-semibold">{{ formatIDR(balance.assetsTotal) }}</div>
	                </div>
	              </div>
	            </div>

	            <div class="rounded-2xl border border-black/10 bg-white/70 p-4">
	              <div class="text-xs tracking-widest uppercase text-[color:var(--muted)]">Modal</div>
	              <div class="mt-3 space-y-2 text-sm">
	                <div class="flex items-center justify-between">
	                  <div class="text-[color:var(--muted)]">Modal awal</div>
	                  <div class="font-semibold">{{ formatIDR(balance.openingCapital) }}</div>
	                </div>
	                <div class="flex items-center justify-between">
	                  <div class="text-[color:var(--muted)]">Laba bersih (hari ini)</div>
	                  <div class="font-semibold">{{ formatIDR(summary.netTotal) }}</div>
	                </div>
	                <div class="border-t border-black/10 pt-2 flex items-center justify-between">
	                  <div class="font-semibold">Total Modal</div>
	                  <div class="font-semibold">{{ formatIDR(balance.equityTotal) }}</div>
	                </div>
	              </div>
	            </div>
	          </div>
	        </div>
	      </section>

      <section class="rounded-2xl border border-black/10 bg-[color:var(--paper)] backdrop-blur shadow-[0_18px_60px_rgba(0,0,0,0.10)]">
        <div class="flex flex-wrap items-center justify-between gap-4 px-5 py-4">
          <div>
            <div class="font-brand text-xl">Transaksi</div>
            <div class="mt-1 text-sm text-[color:var(--muted)]">
              Daftar order di tanggal terpilih.
            </div>
          </div>
          <div class="flex items-center gap-2">
            <div class="relative">
              <input
                v-model="orderSearch"
                type="search"
                placeholder="Cari Order #..."
                class="w-48 rounded-xl border border-black/10 bg-white px-3 py-2 text-xs outline-none focus:border-black/20 focus:shadow-[0_0_0_4px_rgba(193,122,59,0.15)]"
              />
            </div>
          </div>
        </div>

        <div class="border-t border-black/10 px-5 py-5">
          <div v-if="orders.length === 0" class="text-sm text-[color:var(--muted)]">
            Belum ada transaksi.
          </div>

          <div v-else class="overflow-auto">
            <table class="w-full min-w-[860px] text-left text-sm">
              <thead class="text-xs tracking-widest uppercase text-[color:var(--muted)]">
                <tr>
                  <th class="py-2">Jam</th>
                  <th class="py-2">Order</th>
                  <th class="py-2">Tipe</th>
                  <th class="py-2">Meja</th>
                  <th class="py-2">Kasir</th>
                  <th class="py-2">Bayar</th>
                  <th class="py-2 text-right">Total</th>
                  <th class="py-2 text-right">Aksi</th>
                </tr>
              </thead>
              <tbody>
                <tr
                  v-for="o in filteredOrders"
                  :key="o.id"
                  class="border-t border-black/10"
                >
                  <td class="py-3 text-xs text-[color:var(--muted)]">{{ formatTime(o.created_at) }}</td>
                  <td class="py-3 font-mono text-xs">#{{ orderNoFromId(o.id) }}</td>
                  <td class="py-3">
                    <span class="rounded-full border border-black/10 bg-white/70 px-2 py-1 text-xs font-semibold">
                      {{ o.order_type === 'take_away' ? 'Take Away' : 'Dine In' }}
                    </span>
                  </td>
                  <td class="py-3 text-xs">
                    <span v-if="o.table_no" class="font-semibold">Meja {{ o.table_no }}</span>
                    <span v-else class="text-[color:var(--muted)]">-</span>
                  </td>
                  <td class="py-3 text-xs">
                    <div class="font-semibold">{{ o.cashier?.name || '-' }}</div>
                    <div class="text-[color:var(--muted)]">{{ o.cashier?.email || '' }}</div>
                  </td>
                  <td class="py-3">
                    <span class="rounded-full border border-black/10 bg-white/70 px-2 py-1 text-xs font-semibold">
                      {{ (o.payment_method || 'cash').toUpperCase() }}
                    </span>
                  </td>
                  <td class="py-3 text-right font-semibold">{{ formatIDR(o.total) }}</td>
                  <td class="py-3 text-right">
                    <button
                      type="button"
                      class="rounded-xl border border-black/10 bg-white/70 px-3 py-2 text-xs font-semibold hover:bg-white"
                      @click="openOrderDetail(o.id)"
                    >
                      Detail
                    </button>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </section>

	      <section class="rounded-2xl border border-black/10 bg-[color:var(--paper)] backdrop-blur shadow-[0_18px_60px_rgba(0,0,0,0.10)]">
	        <div class="flex flex-wrap items-start justify-between gap-4 px-5 py-4">
	          <div>
	            <div class="font-brand text-xl">Pembukuan</div>
            <div class="mt-1 text-sm text-[color:var(--muted)]">
              Catat pemasukan/pengeluaran manual (contoh: gas, gula, susu).
            </div>
          </div>
        </div>

	        <div class="border-t border-black/10 px-5 py-5">
          <div class="grid grid-cols-1 gap-3 md:grid-cols-6">
            <label class="block">
              <span class="text-xs font-medium text-[color:var(--muted)]">Tipe</span>
              <select
                v-model="ledgerForm.type"
                class="mt-1 w-full rounded-xl border border-black/10 bg-white px-3 py-2 text-sm outline-none focus:border-black/20 focus:shadow-[0_0_0_4px_rgba(193,122,59,0.15)]"
              >
                <option value="expense">Expense</option>
                <option value="income">Income</option>
              </select>
            </label>
            <label class="block">
              <span class="text-xs font-medium text-[color:var(--muted)]">Metode</span>
              <select
                v-model="ledgerForm.payment_method"
                class="mt-1 w-full rounded-xl border border-black/10 bg-white px-3 py-2 text-sm outline-none focus:border-black/20 focus:shadow-[0_0_0_4px_rgba(193,122,59,0.15)]"
              >
                <option value="cash">Cash</option>
                <option value="bank">Bank/QRIS</option>
              </select>
            </label>
            <label class="block">
              <span class="text-xs font-medium text-[color:var(--muted)]">Kategori</span>
              <input
                v-model.trim="ledgerForm.category"
                class="mt-1 w-full rounded-xl border border-black/10 bg-white px-3 py-2 text-sm outline-none focus:border-black/20 focus:shadow-[0_0_0_4px_rgba(193,122,59,0.15)]"
                placeholder="general / bahan / operasional"
              />
            </label>
            <label class="block">
              <span class="text-xs font-medium text-[color:var(--muted)]">Keterangan</span>
              <input
                v-model.trim="ledgerForm.description"
                class="mt-1 w-full rounded-xl border border-black/10 bg-white px-3 py-2 text-sm outline-none focus:border-black/20 focus:shadow-[0_0_0_4px_rgba(193,122,59,0.15)]"
                placeholder="Contoh: beli gula 2kg"
              />
            </label>
            <label class="block">
              <span class="text-xs font-medium text-[color:var(--muted)]">Nominal</span>
              <input
                v-model="ledgerForm.amount"
                inputmode="numeric"
                class="mt-1 w-full rounded-xl border border-black/10 bg-white px-3 py-2 text-sm outline-none focus:border-black/20 focus:shadow-[0_0_0_4px_rgba(193,122,59,0.15)]"
                placeholder="50000"
              />
            </label>
            <button
              type="button"
              class="md:mt-5 rounded-xl bg-[color:var(--accent)] px-4 py-2 text-sm font-semibold text-white shadow-[0_14px_30px_rgba(193,122,59,0.30)] transition hover:brightness-95 disabled:cursor-not-allowed disabled:opacity-60"
              :disabled="ledgerSubmitting"
              @click="addLedger"
            >
              <span v-if="!ledgerSubmitting">Tambah</span>
              <span v-else>...</span>
            </button>
          </div>

          <div class="mt-5">
            <div v-if="ledger.length === 0" class="text-sm text-[color:var(--muted)]">
              Belum ada catatan pembukuan.
            </div>

	            <div v-else class="overflow-auto">
	              <table class="w-full min-w-[860px] text-left text-sm">
	                <thead class="text-xs tracking-widest uppercase text-[color:var(--muted)]">
	                  <tr>
	                    <th class="py-2">Jam</th>
	                    <th class="py-2">Tipe</th>
	                    <th class="py-2">Metode</th>
	                    <th class="py-2">Kategori</th>
	                    <th class="py-2">Keterangan</th>
	                    <th class="py-2 text-right">Nominal</th>
	                    <th class="py-2">User</th>
	                    <th class="py-2 text-right">Aksi</th>
                  </tr>
                </thead>
                <tbody>
	                  <tr v-for="e in ledger" :key="e.id" class="border-t border-black/10">
	                    <td class="py-3 text-xs text-[color:var(--muted)]">{{ formatTime(e.created_at) }}</td>
	                    <td class="py-3">
	                      <span
	                        class="rounded-full border border-black/10 px-2 py-1 text-xs font-semibold"
	                        :class="e.type === 'expense' ? 'bg-red-50 text-red-700' : 'bg-emerald-50 text-emerald-700'"
	                      >
	                        {{ e.type.toUpperCase() }}
	                      </span>
	                    </td>
	                    <td class="py-3">
	                      <span class="rounded-full border border-black/10 bg-white/70 px-2 py-1 text-xs font-semibold">
	                        {{ (e.payment_method || 'cash').toUpperCase() }}
	                      </span>
	                    </td>
	                    <td class="py-3 text-xs">
	                      <span class="rounded-full border border-black/10 bg-white/70 px-2 py-1 text-xs font-semibold">
	                        {{ e.category || 'general' }}
	                      </span>
	                    </td>
	                    <td class="py-3">
	                      <div class="font-semibold">{{ e.description || '-' }}</div>
	                    </td>
                    <td class="py-3 text-right font-semibold">{{ formatIDR(e.amount) }}</td>
                    <td class="py-3 text-xs">
                      <div class="font-semibold">{{ e.created_by?.name || '-' }}</div>
                      <div class="text-[color:var(--muted)]">{{ e.created_by?.email || '' }}</div>
                    </td>
                    <td class="py-3 text-right">
                      <button
                        type="button"
                        class="rounded-xl px-3 py-2 text-xs font-semibold text-[color:var(--danger)] hover:bg-red-50 disabled:cursor-not-allowed disabled:opacity-60 disabled:hover:bg-transparent"
                        :disabled="false"
                        @click="removeLedger(e)"
                      >
                        Hapus
                      </button>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
        </div>
      </section>
    </main>

    <div
      v-if="detail.open"
      class="fixed inset-0 z-50 flex items-end justify-center p-4 sm:items-center"
    >
      <div class="absolute inset-0 bg-black/35" @click="closeDetail" />
      <div
        class="relative w-full max-w-lg animate-float-in rounded-2xl border border-black/10 bg-[color:var(--paper-strong)] backdrop-blur p-5 shadow-[0_24px_80px_rgba(0,0,0,0.18)]"
      >
        <div class="flex items-start justify-between gap-3">
          <div>
            <div class="text-xs tracking-widest uppercase text-[color:var(--muted)]">Order</div>
            <div class="mt-1 font-brand text-2xl font-mono">
              #{{ detail.data ? orderNoFromId(detail.data.id) : '-----' }}
            </div>
            <div class="mt-1 text-sm text-[color:var(--muted)]">
              {{ detail.data?.order_type === 'take_away' ? 'Take Away' : 'Dine In' }}
              <span v-if="detail.data?.table_no"> • Meja {{ detail.data.table_no }}</span>
              <span v-if="detail.data?.guest_count"> • {{ detail.data.guest_count }} Tamu</span>
              <span v-if="detail.data?.customer_name"> • Nama: {{ detail.data.customer_name }}</span>
            </div>
          </div>
          <button
            type="button"
            class="rounded-xl border border-black/10 bg-white/70 px-3 py-2 text-sm font-semibold hover:bg-white"
            @click="closeDetail"
          >
            Tutup
          </button>
        </div>

        <div v-if="detail.loading" class="mt-4 text-sm text-[color:var(--muted)]">
          Memuat detail...
        </div>
        <div v-else-if="detail.error" class="mt-4 text-sm text-red-700">
          {{ detail.error }}
        </div>
        <div v-else-if="detail.data" class="mt-4 rounded-2xl border border-black/10 bg-white/70 p-4">
          <div class="flex items-center justify-between">
            <div class="text-sm font-semibold">Item</div>
            <div class="text-sm font-semibold">{{ formatIDR(detail.data.total) }}</div>
          </div>
          <div class="mt-3 space-y-2">
            <div
              v-for="(it, idx) in detail.data.items"
              :key="idx"
              class="flex items-start justify-between gap-4 text-sm"
            >
              <div class="min-w-0">
                <div class="font-semibold">{{ it.qty }}x {{ it.name }}</div>
                <div class="text-xs text-[color:var(--muted)]">{{ formatIDR(it.price) }} / item</div>
              </div>
              <div class="font-semibold">{{ formatIDR(it.line_total) }}</div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div
      v-if="toast"
      class="fixed bottom-4 left-1/2 z-50 w-[min(560px,calc(100%-2rem))] -translate-x-1/2"
    >
      <div class="rounded-2xl border border-black/10 bg-[color:var(--paper-strong)] px-4 py-3 shadow-[0_18px_60px_rgba(0,0,0,0.14)]">
        <div class="flex items-start justify-between gap-3">
          <div class="text-sm font-medium">{{ toast }}</div>
          <button
            type="button"
            class="rounded-lg px-2 py-1 text-sm font-semibold text-[color:var(--muted)] hover:bg-black/5"
            @click="toast = ''"
          >
            Tutup
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import TopBar from '../components/TopBar.vue'
import { useAuthStore } from '../stores/auth'
import { useSettingsStore } from '../stores/settings'
import {
  ApiError,
  createLedgerEntry,
  deleteLedgerEntry,
  getOrder,
  listLedger,
  listOrders,
  getSalesAnalytics,
  getTopProducts
} from '../api/api'
import {
  Chart as ChartJS,
  Title,
  Tooltip,
  Legend,
  LineElement,
  LinearScale,
  PointElement,
  CategoryScale,
  BarElement
} from 'chart.js'
import { Line, Bar } from 'vue-chartjs'

ChartJS.register(
  Title,
  Tooltip,
  Legend,
  LineElement,
  LinearScale,
  PointElement,
  CategoryScale,
  BarElement
)

const router = useRouter()
const auth = useAuthStore()
const settings = useSettingsStore()

function todayLocalISO() {
  const d = new Date()
  const pad = n => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())}`
}

const selectedDate = ref(todayLocalISO())
const loading = ref(false)
const error = ref('')

const orders = ref([])
const orderSearch = ref('')
const ledger = ref([])

const OPENING_LS_KEY = 'finance_opening_by_date_v1'

const opening = ref({
  cash: '',
  bank: ''
})

function parseMoney(raw) {
  const digits = String(raw || '').replace(/[^\d]/g, '')
  if (!digits) return 0
  const n = Number(digits)
  return Number.isFinite(n) ? n : 0
}

function readOpeningMap() {
  try {
    const raw = localStorage.getItem(OPENING_LS_KEY)
    const parsed = raw ? JSON.parse(raw) : {}
    return parsed && typeof parsed === 'object' ? parsed : {}
  } catch {
    return {}
  }
}

function saveOpeningMap(map) {
  try {
    localStorage.setItem(OPENING_LS_KEY, JSON.stringify(map || {}))
  } catch {}
}

function loadOpeningForDate(date) {
  const map = readOpeningMap()
  const cur = map?.[date] || {}
  opening.value = {
    cash: cur?.cash != null ? String(cur.cash) : '',
    bank: cur?.bank != null ? String(cur.bank) : ''
  }
}

function persistOpeningForDate(date) {
  const map = readOpeningMap()
  map[date] = {
    cash: parseMoney(opening.value.cash),
    bank: parseMoney(opening.value.bank)
  }
  saveOpeningMap(map)
}

const ledgerForm = ref({
  type: 'expense',
  payment_method: 'cash',
  category: 'general',
  description: '',
  amount: ''
})
const ledgerSubmitting = ref(false)

const toast = ref('')
let toastTimer = null
function showToast(msg) {
  toast.value = msg
  if (toastTimer) clearTimeout(toastTimer)
  toastTimer = setTimeout(() => (toast.value = ''), 3200)
}

function formatIDR(num) {
  const amount = Number(num || 0)
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    maximumFractionDigits: 0
  }).format(amount)
}

function formatTime(d) {
  try {
    const dt = typeof d === 'string' ? new Date(d) : d
    return dt?.toLocaleTimeString?.('id-ID', { hour: '2-digit', minute: '2-digit' }) || '-'
  } catch {
    return '-'
  }
}

function orderNoFromId(orderId) {
  if (!orderId) return '-----'
  const hex = String(orderId).replace(/-/g, '').slice(0, 8)
  const n = Number.parseInt(hex, 16)
  if (!Number.isFinite(n)) return String(orderId).split('-')[0].toUpperCase()
  return String(n % 100000).padStart(5, '0')
}

async function loadAnalytics() {
  try {
    const [salesRaw, topRaw] = await Promise.all([
      getSalesAnalytics(auth.token),
      getTopProducts(auth.token)
    ])

    // Format Sales Chart
    chartData.value.sales = {
      labels: salesRaw.map(s => {
        const d = new Date(s.date)
        return d.toLocaleDateString('id-ID', { day: 'numeric', month: 'short' })
      }),
      datasets: [{
        label: 'Omzet',
        backgroundColor: '#C17A3B',
        borderColor: '#C17A3B',
        borderWidth: 2,
        tension: 0.3,
        data: salesRaw.map(s => s.total)
      }]
    }

    // Format Top Products Chart
    chartData.value.topProducts = {
      labels: topRaw.map(p => p.name),
      datasets: [{
        label: 'Terjual',
        backgroundColor: 'rgba(193, 122, 59, 0.2)',
        borderColor: '#C17A3B',
        borderWidth: 1,
        borderRadius: 8,
        data: topRaw.map(p => p.qty)
      }]
    }
  } catch (e) {
    console.error('Failed to load analytics', e)
  }
}

async function loadAll() {
  loading.value = true
  error.value = ''
  try {
    const [o, l] = await Promise.all([
      listOrders(auth.token, { date: selectedDate.value }),
      listLedger(auth.token, { date: selectedDate.value }),
      loadAnalytics()
    ])
    orders.value = Array.isArray(o) ? o : []
    ledger.value = Array.isArray(l) ? l : []
  } catch (e) {
    if (e instanceof ApiError && e.status === 401) {
      auth.logout()
      router.push('/')
      return
    }
    if (e instanceof ApiError && e.status === 402) {
      error.value = e?.message || 'Akses dibatasi'
      if (auth.user?.role === 'finance') router.push('/blocked')
      return
    }
    error.value = e?.message || 'Unknown error'
  } finally {
    loading.value = false
  }
}

const filteredOrders = computed(() => {
  const q = orderSearch.value.trim().toLowerCase()
  if (!q) return orders.value
  return orders.value.filter(o => {
    const orderNo = orderNoFromId(o.id).toLowerCase()
    const orderId = String(o.id || '').toLowerCase()
    return orderNo.includes(q) || orderId.includes(q)
  })
})

const summary = computed(() => {
  const salesTotal = orders.value.reduce((s, o) => s + Number(o?.total || 0), 0)
  const ordersCount = orders.value.length
  const cashTotal = orders.value
    .filter(o => (o?.payment_method || 'cash') === 'cash')
    .reduce((s, o) => s + Number(o?.total || 0), 0)
  const qrisTotal = orders.value
    .filter(o => (o?.payment_method || '') === 'qris')
    .reduce((s, o) => s + Number(o?.total || 0), 0)
  const incomeTotal = ledger.value
    .filter(e => e?.type === 'income')
    .reduce((s, e) => s + Number(e?.amount || 0), 0)
  const expenseTotal = ledger.value
    .filter(e => e?.type === 'expense')
    .reduce((s, e) => s + Number(e?.amount || 0), 0)

  const incomeCash = ledger.value
    .filter(e => e?.type === 'income' && (e?.payment_method || 'cash') === 'cash')
    .reduce((s, e) => s + Number(e?.amount || 0), 0)
  const incomeBank = ledger.value
    .filter(e => e?.type === 'income' && (e?.payment_method || 'cash') === 'bank')
    .reduce((s, e) => s + Number(e?.amount || 0), 0)
  const expenseCash = ledger.value
    .filter(e => e?.type === 'expense' && (e?.payment_method || 'cash') === 'cash')
    .reduce((s, e) => s + Number(e?.amount || 0), 0)
  const expenseBank = ledger.value
    .filter(e => e?.type === 'expense' && (e?.payment_method || 'cash') === 'bank')
    .reduce((s, e) => s + Number(e?.amount || 0), 0)

  const netTotal = salesTotal + incomeTotal - expenseTotal
  return {
    salesTotal,
    ordersCount,
    cashTotal,
    qrisTotal,
    incomeTotal,
    expenseTotal,
    netTotal,
    incomeCash,
    incomeBank,
    expenseCash,
    expenseBank
  }
})

const balance = computed(() => {
  const openingCash = parseMoney(opening.value.cash)
  const openingBank = parseMoney(opening.value.bank)
  const cashClosing = openingCash + summary.value.cashTotal + summary.value.incomeCash - summary.value.expenseCash
  const bankClosing = openingBank + summary.value.qrisTotal + summary.value.incomeBank - summary.value.expenseBank
  const assetsTotal = cashClosing + bankClosing
  const openingCapital = openingCash + openingBank
  const equityTotal = openingCapital + summary.value.netTotal
  const diff = assetsTotal - equityTotal
  return { openingCash, openingBank, cashClosing, bankClosing, assetsTotal, openingCapital, equityTotal, diff }
})

const detail = ref({
  open: false,
  loading: false,
  error: '',
  data: null
})

const chartData = ref({
  sales: { labels: [], datasets: [] },
  topProducts: { labels: [], datasets: [] }
})

const chartOptions = {
  sales: {
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
      legend: { display: false },
      tooltip: {
        callbacks: {
          label: (ctx) => `Omzet: ${formatIDR(ctx.parsed.y)}`
        }
      }
    },
    scales: {
      y: {
        beginAtZero: true,
        grid: { color: 'rgba(0,0,0,0.05)' },
        ticks: {
          callback: v => v >= 1000000 ? `${v / 1000000}jt` : v >= 1000 ? `${v / 1000}rb` : v
        }
      },
      x: { grid: { display: false } }
    }
  },
  topProducts: {
    responsive: true,
    maintainAspectRatio: false,
    indexAxis: 'y',
    plugins: { legend: { display: false } },
    scales: {
      x: {
        beginAtZero: true,
        grid: { color: 'rgba(0,0,0,0.05)' }
      },
      y: { grid: { display: false } }
    }
  }
}

async function openOrderDetail(id) {
  detail.value = { open: true, loading: true, error: '', data: null }
  try {
    const data = await getOrder(auth.token, id)
    detail.value = { open: true, loading: false, error: '', data }
  } catch (e) {
    if (e instanceof ApiError && e.status === 401) {
      auth.logout()
      router.push('/')
      return
    }
    if (e instanceof ApiError && e.status === 402) {
      detail.value = { open: true, loading: false, error: e?.message || 'Akses dibatasi', data: null }
      return
    }
    detail.value = { open: true, loading: false, error: e?.message || 'Gagal memuat detail', data: null }
  }
}

function closeDetail() {
  detail.value.open = false
}

async function addLedger() {
  const amount = Number(String(ledgerForm.value.amount || '').replace(/[^\d]/g, ''))
  if (!Number.isFinite(amount) || amount <= 0) {
    showToast('Nominal tidak valid')
    return
  }

  ledgerSubmitting.value = true
  try {
    await createLedgerEntry(auth.token, {
      entry_date: selectedDate.value,
      type: ledgerForm.value.type,
      amount,
      payment_method: ledgerForm.value.payment_method,
      category: ledgerForm.value.category?.trim() || 'general',
      description: ledgerForm.value.description?.trim() || null
    })
    ledgerForm.value.amount = ''
    ledgerForm.value.description = ''
    showToast('Catatan ditambahkan')
    await loadAll()
  } catch (e) {
    if (e instanceof ApiError && e.status === 401) {
      auth.logout()
      router.push('/')
      return
    }
    if (e instanceof ApiError && e.status === 402) {
      showToast(e?.message || 'Akses dibatasi')
      return
    }
    showToast(e?.message || 'Gagal menambah catatan')
  } finally {
    ledgerSubmitting.value = false
  }
}

async function removeLedger(e) {
  if (!e?.id) return
  // eslint-disable-next-line no-alert
  const ok = window.confirm('Hapus catatan ini?')
  if (!ok) return

  try {
    await deleteLedgerEntry(auth.token, e.id)
    showToast('Catatan dihapus')
    await loadAll()
  } catch (err) {
    if (err instanceof ApiError && err.status === 401) {
      auth.logout()
      router.push('/')
      return
    }
    if (err instanceof ApiError && err.status === 402) {
      showToast(err?.message || 'Akses dibatasi')
      return
    }
    showToast(err?.message || 'Gagal menghapus catatan')
  }
}

onMounted(async () => {
  await settings.loadStore(auth.token)
  loadOpeningForDate(selectedDate.value)
  await loadAll()
})

watch(
  () => selectedDate.value,
  date => {
    loadOpeningForDate(date)
  }
)

watch(
  () => [opening.value.cash, opening.value.bank, selectedDate.value],
  () => {
    persistOpeningForDate(selectedDate.value)
  }
)
</script>
