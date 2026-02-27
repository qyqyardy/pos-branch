<template>
  <div class="min-h-screen">
    <TopBar />

    <main class="mx-auto max-w-7xl px-4 py-6 space-y-6">
      <section class="rounded-2xl border border-black/10 bg-[color:var(--paper)] backdrop-blur shadow-[0_18px_60px_rgba(0,0,0,0.10)]">
        <div class="flex items-start justify-between gap-4 px-5 py-4">
          <div>
            <div class="font-brand text-xl">Setting Toko</div>
            <div class="mt-1 text-sm text-[color:var(--muted)]">
              Nama, alamat, dan nomor telepon yang dipakai di struk.
            </div>
          </div>
          <button
            type="button"
            class="rounded-xl bg-[color:var(--accent)] px-4 py-2 text-sm font-semibold text-white shadow-[0_14px_30px_rgba(193,122,59,0.30)] transition hover:brightness-95 disabled:cursor-not-allowed disabled:opacity-60"
            :disabled="savingStore"
            @click="saveStore"
          >
            <span v-if="!savingStore">Simpan</span>
            <span v-else>Menyimpan...</span>
          </button>
        </div>

        <div class="border-t border-black/10 px-5 py-5">
          <div class="mb-5 rounded-2xl border border-black/10 bg-white/70 p-4">
            <div class="text-sm font-semibold">Subscription</div>
            <div class="mt-3 grid grid-cols-1 gap-3 sm:grid-cols-3">
              <div>
                <div class="text-xs tracking-widest uppercase text-[color:var(--muted)]">Plan</div>
                <div class="mt-1 font-semibold">
                  {{ planLabel(settings.store?.plan) }}
                </div>
              </div>
              <div>
                <div class="text-xs tracking-widest uppercase text-[color:var(--muted)]">Berlaku sampai</div>
                <div class="mt-1 font-mono text-xs">
                  {{ formatPaidUntil(settings.store?.paid_until) }}
                </div>
              </div>
              <div>
                <div class="text-xs tracking-widest uppercase text-[color:var(--muted)]">Status</div>
                <div
                  class="mt-1 font-semibold"
                  :class="settings.store?.subscription_active === false ? 'text-red-700' : 'text-emerald-700'"
                >
                  {{ settings.store?.subscription_active === false ? 'Expired' : 'Active' }}
                </div>
              </div>
            </div>
            <div class="mt-3 text-xs text-[color:var(--muted)]">
              Plan dan masa aktif diatur oleh vendor.
            </div>
          </div>

          <div class="mb-5 rounded-2xl border border-black/10 bg-white/70 p-4">
            <div class="text-sm font-semibold">Logo</div>
            <div class="mt-3 flex flex-wrap items-center gap-4">
              <div class="h-20 w-20 overflow-hidden rounded-3xl border border-black/10 bg-white flex items-center justify-center">
                <img
                  v-if="logoPreview"
                  :src="logoPreview"
                  alt=""
                  class="h-full w-full object-contain"
                />
                <div v-else class="text-xs text-[color:var(--muted)]">No logo</div>
              </div>

              <div class="min-w-[240px] flex-1">
                <input
                  type="file"
                  accept="image/*"
                  class="block w-full text-sm"
                  @change="onPickLogo"
                />
                <div class="mt-2 flex items-center gap-2">
                  <button
                    v-if="logoPreview"
                    type="button"
                    class="rounded-xl border border-black/10 bg-white/70 px-3 py-2 text-xs font-semibold hover:bg-white"
                    @click="removeLogo"
                  >
                    Hapus logo
                  </button>
                  <div class="text-xs text-[color:var(--muted)]">
                    Rekomendasi: PNG/JPG/SVG, max 250KB.
                  </div>
                </div>

                <p v-if="logoError" class="mt-2 text-sm text-red-700">
                  {{ logoError }}
                </p>
              </div>
            </div>
          </div>

          <div class="grid grid-cols-1 gap-4 md:grid-cols-2">
            <label class="block">
              <span class="text-sm font-medium">Nama Warkop/Cafe</span>
              <input
                v-model.trim="storeForm.name"
                class="mt-1 w-full rounded-xl border border-black/10 bg-white px-4 py-3 text-sm outline-none focus:border-black/20 focus:shadow-[0_0_0_4px_rgba(193,122,59,0.15)]"
                placeholder="Warkop Kamu"
              />
            </label>

            <label class="block">
              <span class="text-sm font-medium">Tagline</span>
              <input
                v-model.trim="storeForm.tagline"
                class="mt-1 w-full rounded-xl border border-black/10 bg-white px-4 py-3 text-sm outline-none focus:border-black/20 focus:shadow-[0_0_0_4px_rgba(193,122,59,0.15)]"
                placeholder="Point of Sale"
              />
            </label>

            <label class="block">
              <span class="text-sm font-medium">Alamat 1</span>
              <input
                v-model.trim="storeForm.address1"
                class="mt-1 w-full rounded-xl border border-black/10 bg-white px-4 py-3 text-sm outline-none focus:border-black/20 focus:shadow-[0_0_0_4px_rgba(193,122,59,0.15)]"
                placeholder="Jl. ..."
              />
            </label>

            <label class="block">
              <span class="text-sm font-medium">Alamat 2</span>
              <input
                v-model.trim="storeForm.address2"
                class="mt-1 w-full rounded-xl border border-black/10 bg-white px-4 py-3 text-sm outline-none focus:border-black/20 focus:shadow-[0_0_0_4px_rgba(193,122,59,0.15)]"
                placeholder="Kota/RT RW (opsional)"
              />
            </label>

            <label class="block md:col-span-2">
              <span class="text-sm font-medium">Nomor Telepon</span>
              <input
                v-model.trim="storeForm.phone"
                class="mt-1 w-full rounded-xl border border-black/10 bg-white px-4 py-3 text-sm outline-none focus:border-black/20 focus:shadow-[0_0_0_4px_rgba(193,122,59,0.15)]"
                placeholder="08xx-xxxx-xxxx"
              />
            </label>
          </div>

          <p v-if="settings.error" class="mt-4 text-sm text-red-700">
            {{ settings.error }}
          </p>
        </div>
      </section>


      <section class="rounded-2xl border border-black/10 bg-[color:var(--paper)] backdrop-blur shadow-[0_18px_60px_rgba(0,0,0,0.10)]">
        <div class="flex items-start justify-between gap-4 px-5 py-4">
          <div>
            <div class="font-brand text-xl">User</div>
            <div class="mt-1 text-sm text-[color:var(--muted)]">
              Buat user kasir, finance, atau superadmin.
            </div>
          </div>
          <button
            type="button"
            class="rounded-xl border border-black/10 bg-white/70 px-4 py-2 text-sm font-semibold hover:bg-white"
            @click="openCreateUser"
          >
            Tambah user
          </button>
        </div>

        <div class="border-t border-black/10 px-5 py-5">
          <div v-if="usersLoading" class="text-sm text-[color:var(--muted)]">
            Memuat user...
          </div>

          <div v-else-if="usersError" class="rounded-xl border border-black/10 bg-white/70 p-4">
            <div class="font-semibold">Gagal memuat user</div>
            <div class="mt-1 text-sm text-[color:var(--muted)]">{{ usersError }}</div>
            <button
              type="button"
              class="mt-4 rounded-xl border border-black/10 bg-white/70 px-4 py-2 text-sm font-semibold hover:bg-white"
              @click="loadUsers"
            >
              Coba lagi
            </button>
          </div>

          <div v-else class="overflow-auto">
            <table class="w-full min-w-[720px] text-left text-sm">
              <thead class="text-xs tracking-widest uppercase text-[color:var(--muted)]">
                <tr>
                  <th class="py-2">Nama</th>
                  <th class="py-2">Email</th>
                  <th class="py-2">Role</th>
                  <th class="py-2">Dibuat</th>
                  <th class="py-2 text-right">Aksi</th>
                </tr>
              </thead>
              <tbody>
                <tr
                  v-for="u in users"
                  :key="u.id"
                  class="border-t border-black/10"
                >
                  <td class="py-3">
                    <div class="font-semibold">
                      {{ u.name }}
                      <span v-if="u.id === auth.user?.id" class="text-xs text-[color:var(--muted)]">(kamu)</span>
                    </div>
                  </td>
                  <td class="py-3 font-mono text-xs">{{ u.email }}</td>
                  <td class="py-3">
                    <span class="rounded-full border border-black/10 bg-white/70 px-2 py-1 text-xs font-semibold">
                      {{ roleLabel(u.role) }}
                    </span>
                  </td>
                  <td class="py-3 text-xs text-[color:var(--muted)]">
                    {{ formatDate(u.created_at) }}
                  </td>
                  <td class="py-3 text-right">
                    <div class="flex justify-end gap-2">
                      <button
                        type="button"
                        class="rounded-xl border border-black/10 bg-white/70 px-3 py-2 text-xs font-semibold hover:bg-white"
                        @click="openEditUser(u)"
                      >
                        Edit
                      </button>
                      <button
                        type="button"
                        class="rounded-xl px-3 py-2 text-xs font-semibold text-[color:var(--danger)] hover:bg-red-50"
                        :disabled="u.id === auth.user?.id"
                        @click="confirmDeleteUser(u)"
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
      v-if="userModal.open"

      class="fixed inset-0 z-50 flex items-end justify-center p-4 sm:items-center"
    >
      <div class="absolute inset-0 bg-black/35" @click="closeUserModal" />
      <div
        class="relative w-full max-w-lg animate-float-in rounded-2xl border border-black/10 bg-[color:var(--paper-strong)] backdrop-blur p-5 shadow-[0_24px_80px_rgba(0,0,0,0.18)]"
      >
        <div class="flex items-start justify-between gap-3">
          <div>
            <div class="text-xs tracking-widest uppercase text-[color:var(--muted)]">
              {{ userModal.mode === 'create' ? 'Tambah user' : 'Edit user' }}
            </div>
            <div class="mt-1 font-brand text-2xl">
              {{ userModal.mode === 'create' ? 'User baru' : 'Perbarui user' }}
            </div>
          </div>
          <button
            type="button"
            class="rounded-xl border border-black/10 bg-white/70 px-3 py-2 text-sm font-semibold hover:bg-white"
            @click="closeUserModal"
          >
            Tutup
          </button>
        </div>

        <div class="mt-4 grid grid-cols-1 gap-3">
          <label class="block">
            <span class="text-sm font-medium">Nama</span>
            <input
              v-model.trim="userModal.form.name"
              class="mt-1 w-full rounded-xl border border-black/10 bg-white px-4 py-3 text-sm outline-none focus:border-black/20 focus:shadow-[0_0_0_4px_rgba(193,122,59,0.15)]"
              placeholder="Nama kasir"
            />
          </label>

          <label class="block">
            <span class="text-sm font-medium">Email</span>
            <input
              v-model.trim="userModal.form.email"
              type="email"
              autocomplete="off"
              class="mt-1 w-full rounded-xl border border-black/10 bg-white px-4 py-3 text-sm outline-none focus:border-black/20 focus:shadow-[0_0_0_4px_rgba(193,122,59,0.15)]"
              placeholder="kasir@warkop.com"
            />
          </label>

          <label class="block">
            <span class="text-sm font-medium">Role</span>
            <select
              v-model="userModal.form.role"
              class="mt-1 w-full rounded-xl border border-black/10 bg-white px-4 py-3 text-sm outline-none focus:border-black/20 focus:shadow-[0_0_0_4px_rgba(193,122,59,0.15)]"
            >
              <option value="cashier">Kasir</option>
              <option value="finance" :disabled="!isPremium">Finance (Premium)</option>
              <option value="admin">Superadmin</option>
            </select>
            <div v-if="!isPremium" class="mt-2 text-xs text-[color:var(--muted)]">
              Role Finance hanya tersedia untuk plan Premium.
            </div>
          </label>

          <label class="block">
            <span class="text-sm font-medium">
              Password
              <span v-if="userModal.mode === 'edit'" class="text-xs text-[color:var(--muted)]">(kosongkan jika tidak diganti)</span>
            </span>
            <input
              v-model="userModal.form.password"
              type="password"
              autocomplete="new-password"
              class="mt-1 w-full rounded-xl border border-black/10 bg-white px-4 py-3 text-sm outline-none focus:border-black/20 focus:shadow-[0_0_0_4px_rgba(193,122,59,0.15)]"
              placeholder="••••••••"
            />
          </label>
        </div>

        <p v-if="userModal.error" class="mt-3 text-sm text-red-700">
          {{ userModal.error }}
        </p>

        <button
          type="button"
          class="mt-4 w-full rounded-xl bg-[color:var(--accent)] px-4 py-3 font-semibold text-white shadow-[0_14px_30px_rgba(193,122,59,0.35)] transition hover:brightness-95 disabled:cursor-not-allowed disabled:opacity-60"
          :disabled="userModal.submitting"
          @click="submitUserModal"
        >
          <span v-if="!userModal.submitting">Simpan</span>
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
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import TopBar from '../components/TopBar.vue'
import { useAuthStore } from '../stores/auth'
import { useSettingsStore } from '../stores/settings'
import { ApiError, createUser, deleteUser, listUsers, updateUser } from '../api/api'


const router = useRouter()
const auth = useAuthStore()
const settings = useSettingsStore()
const isPremium = computed(() => (settings.store?.plan || 'premium') === 'premium')

const storeForm = ref({
  name: '',
  tagline: '',
  address1: '',
  address2: '',
  phone: ''
})
const savingStore = ref(false)

const logoPreview = ref('')
const logoDirty = ref(false)
const logoError = ref('')

const users = ref([])
const usersLoading = ref(false)
const usersError = ref('')

function formatIDR(val) {
  return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(val || 0)
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
    price: 0
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
    form: { name: '', price: 0 }
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
    form: { name: p.name || '', price: p.price || 0 }
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
    image_data_url: productModal.value.imgPreview
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

const toast = ref('')

let toastTimer = null

function showToast(msg) {
  toast.value = msg
  if (toastTimer) clearTimeout(toastTimer)
  toastTimer = setTimeout(() => (toast.value = ''), 3200)
}

function roleLabel(role) {
  if (role === 'admin') return 'Superadmin'
  if (role === 'finance') return 'Finance'
  return 'Kasir'
}

function planLabel(plan) {
  return String(plan || '').toLowerCase() === 'standard' ? 'Standard' : 'Premium'
}

function formatPaidUntil(iso) {
  if (!iso) return '-'
  try {
    const d = new Date(iso)
    return d.toLocaleDateString('id-ID', { year: 'numeric', month: 'short', day: '2-digit' })
  } catch {
    return String(iso)
  }
}

function formatDate(d) {
  try {
    const dt = typeof d === 'string' ? new Date(d) : d
    return dt?.toLocaleString?.('id-ID') || '-'
  } catch {
    return '-'
  }
}

function syncStoreForm() {
  const s = settings.store || {}
  storeForm.value = {
    name: s.name || '',
    tagline: s.tagline || '',
    address1: (s.address_lines && s.address_lines[0]) || '',
    address2: (s.address_lines && s.address_lines[1]) || '',
    phone: s.phone || ''
  }
  logoPreview.value = s.logo_data_url || ''
  logoDirty.value = false
  logoError.value = ''
}

function onPickLogo(e) {
  logoError.value = ''
  const file = e?.target?.files?.[0]
  if (!file) return

  if (!String(file.type || '').startsWith('image/')) {
    logoError.value = 'File harus gambar'
    return
  }

  const maxBytes = 250 * 1024
  if (file.size > maxBytes) {
    logoError.value = 'Logo terlalu besar (max 250KB)'
    return
  }

  const reader = new FileReader()
  reader.onload = () => {
    logoPreview.value = String(reader.result || '')
    logoDirty.value = true
  }
  reader.onerror = () => {
    logoError.value = 'Gagal membaca file'
  }
  reader.readAsDataURL(file)
}

function removeLogo() {
  logoPreview.value = ''
  logoDirty.value = true
}

async function saveStore() {
  savingStore.value = true
  try {
    const payload = {
      name: storeForm.value.name,
      tagline: storeForm.value.tagline,
      address_lines: [storeForm.value.address1, storeForm.value.address2],
      phone: storeForm.value.phone
    }
    if (logoDirty.value) payload.logo_data_url = logoPreview.value || ''

    await settings.saveStore(auth.token, payload)
    logoDirty.value = false
    showToast('Setting toko tersimpan')
  } catch (e) {
    if (e instanceof ApiError && e.status === 401) {
      auth.logout()
      router.push('/')
      return
    }
    showToast(e?.message || 'Gagal menyimpan setting')
  } finally {
    savingStore.value = false
  }
}

async function loadUsers() {
  usersLoading.value = true
  usersError.value = ''
  try {
    const list = await listUsers(auth.token)
    users.value = Array.isArray(list) ? list : []
  } catch (e) {
    if (e instanceof ApiError && e.status === 401) {
      auth.logout()
      router.push('/')
      return
    }
    usersError.value = e?.message || 'Unknown error'
  } finally {
    usersLoading.value = false
  }
}

const userModal = ref({
  open: false,
  mode: 'create',
  submitting: false,
  error: '',
  userId: null,
  form: {
    name: '',
    email: '',
    role: 'cashier',
    password: ''
  }
})

function openCreateUser() {
  userModal.value = {
    open: true,
    mode: 'create',
    submitting: false,
    error: '',
    userId: null,
    form: { name: '', email: '', role: 'cashier', password: '' }
  }
}

function openEditUser(u) {
  userModal.value = {
    open: true,
    mode: 'edit',
    submitting: false,
    error: '',
    userId: u.id,
    form: { name: u.name || '', email: u.email || '', role: u.role || 'cashier', password: '' }
  }
}

function closeUserModal() {
  if (userModal.value.submitting) return
  userModal.value.open = false
}

async function submitUserModal() {
  userModal.value.error = ''
  userModal.value.submitting = true

  const payload = {
    name: userModal.value.form.name,
    email: userModal.value.form.email,
    role: userModal.value.form.role
  }
  if (userModal.value.form.password?.trim()) payload.password = userModal.value.form.password

  try {
    if (userModal.value.mode === 'create') {
      if (!payload.password) throw new Error('Password wajib diisi')
      await createUser(auth.token, payload)
      showToast('User berhasil dibuat')
    } else {
      await updateUser(auth.token, userModal.value.userId, payload)
      showToast('User berhasil diperbarui')
    }
    closeUserModal()
    await loadUsers()
  } catch (e) {
    if (e instanceof ApiError && e.status === 401) {
      auth.logout()
      router.push('/')
      return
    }
    userModal.value.error = e?.message || 'Gagal menyimpan user'
  } finally {
    userModal.value.submitting = false
  }
}

async function confirmDeleteUser(u) {
  if (!u?.id) return
  if (u.id === auth.user?.id) return
  // eslint-disable-next-line no-alert
  const ok = window.confirm(`Hapus user "${u.name}"?`)
  if (!ok) return

  try {
    await deleteUser(auth.token, u.id)
    showToast('User dihapus')
    await loadUsers()
  } catch (e) {
    if (e instanceof ApiError && e.status === 401) {
      auth.logout()
      router.push('/')
      return
    }
    showToast(e?.message || 'Gagal menghapus user')
  }
}

onMounted(async () => {
  await settings.loadStore(auth.token)
  syncStoreForm()
  await loadUsers()
})

</script>
