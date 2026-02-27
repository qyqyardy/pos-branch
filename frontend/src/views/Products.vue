<template>
  <div class="min-h-screen">
    <TopBar />

    <main class="mx-auto max-w-7xl px-4 py-6 space-y-6">
      <section class="rounded-2xl border border-black/10 bg-[color:var(--paper)] backdrop-blur shadow-[0_18px_60px_rgba(0,0,0,0.10)]">
        <div class="flex items-start justify-between gap-4 px-5 py-4">
          <div>
            <div class="font-brand text-xl">Produk (Menu)</div>
            <div class="mt-1 text-sm text-[color:var(--muted)]">
              Atur daftar makanan dan minuman.
            </div>
          </div>
          <button
            type="button"
            class="rounded-xl border border-black/10 bg-white/70 px-4 py-2 text-sm font-semibold hover:bg-white"
            @click="openCreateProduct"
          >
            Tambah produk
          </button>
        </div>

        <div class="border-t border-black/10 px-5 py-5">
          <div v-if="productsLoading" class="text-sm text-[color:var(--muted)]">
            Memuat produk...
          </div>

          <div v-else-if="productsError" class="rounded-xl border border-black/10 bg-white/70 p-4">
            <div class="font-semibold">Gagal memuat produk</div>
            <div class="mt-1 text-sm text-[color:var(--muted)]">{{ productsError }}</div>
            <button
              type="button"
              class="mt-4 rounded-xl border border-black/10 bg-white/70 px-4 py-2 text-sm font-semibold hover:bg-white"
              @click="loadProducts"
            >
              Coba lagi
            </button>
          </div>

          <div v-else class="overflow-auto">
            <table class="w-full min-w-[720px] text-left text-sm">
              <thead class="text-xs tracking-widest uppercase text-[color:var(--muted)]">
                <tr>
                  <th class="py-2" width="80">Gambar</th>
                  <th class="py-2">Nama</th>
                   <th class="py-2">Harga</th>
                  <th class="py-2 text-center" width="100">Stok</th>
                  <th class="py-2 text-center">Status</th>
                  <th class="py-2 text-right">Aksi</th>
                </tr>
              </thead>
              <tbody>
                <tr
                  v-for="p in products"
                  :key="p.id"
                  class="border-t border-black/10"
                >
                  <td class="py-3">
                    <div class="h-12 w-12 overflow-hidden rounded-xl border border-black/10 bg-white">
                      <img
                        :src="p.image_data_url || '/menu/_default.svg'"
                        alt=""
                        class="h-full w-full object-cover"
                      />
                    </div>
                  </td>
                  <td class="py-3">
                    <div class="font-semibold">{{ p.name }}</div>
                  </td>
                   <td class="py-3 font-mono text-xs">{{ formatIDR(p.price) }}</td>
                  <td class="py-3 text-center">
                    <span class="font-mono text-xs" :class="p.stock <= 5 ? 'text-red-600 font-bold' : ''">
                      {{ p.stock }}
                    </span>
                  </td>
                  <td class="py-3 text-center">
                    <span
                      class="rounded-full px-2 py-1 text-[10px] font-bold uppercase tracking-wider"
                      :class="p.is_active ? 'bg-emerald-100 text-emerald-700' : 'bg-red-100 text-red-700'"
                    >
                      {{ p.is_active ? 'Tersedia' : 'Kosong' }}
                    </span>
                  </td>
                  <td class="py-3 text-right">
                    <div class="flex justify-end gap-2">
                      <button
                        type="button"
                        class="rounded-xl border border-black/10 bg-white/70 px-3 py-2 text-xs font-semibold hover:bg-white"
                        @click="openEditProduct(p)"
                      >
                        Edit
                      </button>
                      <button
                        type="button"
                        class="rounded-xl px-3 py-2 text-xs font-semibold text-[color:var(--danger)] hover:bg-red-50"
                        @click="confirmDeleteProduct(p)"
                      >
                        Hapus
                      </button>
                    </div>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </section>
    </main>

    <div
      v-if="productModal.open"
      class="fixed inset-0 z-50 flex items-end justify-center p-4 sm:items-center"
    >
      <div class="absolute inset-0 bg-black/35" @click="closeProductModal" />
      <div
        class="relative w-full max-w-lg animate-float-in rounded-2xl border border-black/10 bg-[color:var(--paper-strong)] backdrop-blur p-5 shadow-[0_24px_80px_rgba(0,0,0,0.18)]"
      >
        <div class="flex items-start justify-between gap-3">
          <div>
            <div class="text-xs tracking-widest uppercase text-[color:var(--muted)]">
              {{ productModal.mode === 'create' ? 'Tambah produk' : 'Edit produk' }}
            </div>
            <div class="mt-1 font-brand text-2xl">
              {{ productModal.mode === 'create' ? 'Produk baru' : 'Perbarui produk' }}
            </div>
          </div>
          <button
            type="button"
            class="rounded-xl border border-black/10 bg-white/70 px-3 py-2 text-sm font-semibold hover:bg-white"
            @click="closeProductModal"
          >
            Tutup
          </button>
        </div>

        <div class="mt-4 grid grid-cols-1 gap-3">
          <div class="rounded-2xl border border-black/10 bg-white/70 p-4">
            <div class="text-sm font-semibold">Gambar Produk</div>
            <div class="mt-3 flex items-center gap-4">
              <div class="h-20 w-20 overflow-hidden rounded-2xl border border-black/10 bg-white flex items-center justify-center">
                <img
                  v-if="productModal.imgPreview"
                  :src="productModal.imgPreview"
                  alt=""
                  class="h-full w-full object-cover"
                />
                <div v-else class="text-xs text-[color:var(--muted)]">No image</div>
              </div>
              <div class="flex-1">
                <input
                  type="file"
                  accept="image/*"
                  class="block w-full text-xs"
                  @change="onPickProductImg"
                />
                <button
                  v-if="productModal.imgPreview"
                  type="button"
                  class="mt-2 text-xs font-semibold text-[color:var(--danger)]"
                  @click="productModal.imgPreview = ''"
                >
                  Hapus gambar
                </button>
              </div>
            </div>
          </div>

          <label class="block">
            <span class="text-sm font-medium">Nama Produk</span>
            <input
              v-model.trim="productModal.form.name"
              class="mt-1 w-full rounded-xl border border-black/10 bg-white px-4 py-3 text-sm outline-none focus:border-black/20 focus:shadow-[0_0_0_4px_rgba(193,122,59,0.15)]"
              placeholder="Contoh: Kopi Susu"
            />
          </label>

           <label class="block">
            <span class="text-sm font-medium">Harga (IDR)</span>
            <input
              v-model.number="productModal.form.price"
              type="number"
              class="mt-1 w-full rounded-xl border border-black/10 bg-white px-4 py-3 text-sm outline-none focus:border-black/20 focus:shadow-[0_0_0_4px_rgba(193,122,59,0.15)]"
              placeholder="15000"
            />
          </label>

          <label class="block">
            <span class="text-sm font-medium">Stok (pcs/porsi)</span>
            <input
              v-model.number="productModal.form.stock"
              type="number"
              class="mt-1 w-full rounded-xl border border-black/10 bg-white px-4 py-3 text-sm outline-none focus:border-black/20 focus:shadow-[0_0_0_4px_rgba(193,122,59,0.15)]"
              placeholder="100"
            />
          </label>

          <label class="flex items-center gap-3 rounded-2xl border border-black/10 bg-white/70 p-4 cursor-pointer hover:bg-white transition">
            <input
              type="checkbox"
              v-model="productModal.form.is_active"
              class="h-5 w-5 rounded-lg border-black/10 text-[color:var(--accent)] focus:ring-[color:var(--accent)]"
            />
            <div>
              <div class="text-sm font-semibold">Tersedia di Menu</div>
              <div class="text-xs text-[color:var(--muted)]">Matikan jika stok sedang kosong agar tidak muncul di POS.</div>
            </div>
          </label>
        </div>

        <p v-if="productModal.error" class="mt-3 text-sm text-red-700">
          {{ productModal.error }}
        </p>

        <button
          type="button"
          class="mt-4 w-full rounded-xl bg-[color:var(--accent)] px-4 py-3 font-semibold text-white shadow-[0_14px_30px_rgba(193,122,59,0.35)] transition hover:brightness-95 disabled:cursor-not-allowed disabled:opacity-60"
          :disabled="productModal.submitting"
          @click="submitProductModal"
        >
          <span v-if="!productModal.submitting">Simpan</span>
          <span v-else>Memproses...</span>
        </button>
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
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import TopBar from '../components/TopBar.vue'
import { useAuthStore } from '../stores/auth'
import { ApiError, createProduct, deleteProduct, getProducts, updateProduct } from '../api/api'

const router = useRouter()
const auth = useAuthStore()

const products = ref([])
const productsLoading = ref(false)
const productsError = ref('')

const toast = ref('')
let toastTimer = null

function showToast(msg) {
  toast.value = msg
  if (toastTimer) clearTimeout(toastTimer)
  toastTimer = setTimeout(() => (toast.value = ''), 3200)
}

function formatIDR(val) {
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    minimumFractionDigits: 0
  }).format(val || 0)
}

async function loadProducts() {
  productsLoading.value = true
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
    productsError.value = e?.message || 'Gagal memuat produk'
  } finally {
    productsLoading.value = false
  }
}

const productModal = ref({
  open: false,
  mode: 'create',
  submitting: false,
  error: '',
  productId: null,
  imgPreview: '',
  form: {
    name: '',
    price: 0,
    is_active: true,
    stock: 0
  }
})

function openCreateProduct() {
  productModal.value = {
    open: true,
    mode: 'create',
    submitting: false,
    error: '',
    productId: null,
    imgPreview: '',
    form: { name: '', price: 0, is_active: true, stock: 0 }
  }
}

function openEditProduct(p) {
  productModal.value = {
    open: true,
    mode: 'edit',
    submitting: false,
    error: '',
    productId: p.id,
    imgPreview: p.image_data_url || '',
    form: {
      name: p.name || '',
      price: p.price || 0,
      is_active: p.is_active !== false,
      stock: p.stock || 0
    }
  }
}

function closeProductModal() {
  if (productModal.value.submitting) return
  productModal.value.open = false
}

function onPickProductImg(e) {
  productModal.value.error = ''
  const file = e?.target?.files?.[0]
  if (!file) return

  const reader = new FileReader()
  reader.onload = () => {
    productModal.value.imgPreview = String(reader.result || '')
  }
  reader.readAsDataURL(file)
}

async function submitProductModal() {
  productModal.value.error = ''
  productModal.value.submitting = true

  const payload = {
    name: productModal.value.form.name,
    price: parseInt(productModal.value.form.price, 10),
    image_data_url: productModal.value.imgPreview,
    is_active: productModal.value.form.is_active,
    stock: parseInt(productModal.value.form.stock, 10) || 0
  }

  try {
    if (productModal.value.mode === 'create') {
      await createProduct(auth.token, payload)
      showToast('Produk berhasil dibuat')
    } else {
      await updateProduct(auth.token, productModal.value.productId, payload)
      showToast('Produk berhasil diperbarui')
    }
    closeProductModal()
    await loadProducts()
  } catch (e) {
    if (e instanceof ApiError && e.status === 401) {
      auth.logout()
      router.push('/')
      return
    }
    productModal.value.error = e?.message || 'Gagal menyimpan produk'
  } finally {
    productModal.value.submitting = false
  }
}

async function confirmDeleteProduct(p) {
  if (!p?.id) return
  // eslint-disable-next-line no-alert
  const ok = window.confirm(`Hapus produk "${p.name}"?`)
  if (!ok) return

  try {
    await deleteProduct(auth.token, p.id)
    showToast('Produk dihapus')
    await loadProducts()
  } catch (e) {
    if (e instanceof ApiError && e.status === 401) {
      auth.logout()
      router.push('/')
      return
    }
    showToast(e?.message || 'Gagal menghapus produk')
  }
}

onMounted(() => {
  loadProducts()
})
</script>
