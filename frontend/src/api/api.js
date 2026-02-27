const API = import.meta.env.VITE_API_URL || 'http://localhost:8080'

export class ApiError extends Error {
  constructor(status, message, data) {
    super(message)
    this.name = 'ApiError'
    this.status = status
    this.data = data
  }
}

async function readBody(res) {
  const contentType = res.headers.get('content-type') || ''

  // 201 from the backend can be empty; don't fail parsing in that case.
  if (res.status === 204) return null

  if (contentType.includes('application/json')) {
    try {
      return await res.json()
    } catch {
      return null
    }
  }

  try {
    const text = await res.text()
    return text === '' ? null : text
  } catch {
    return null
  }
}

async function request(path, { method = 'GET', token, body } = {}) {
  const headers = {}
  if (body != null) headers['Content-Type'] = 'application/json'
  if (token) headers.Authorization = `Bearer ${token}`

  const res = await fetch(`${API}${path}`, {
    method,
    headers,
    body: body != null ? JSON.stringify(body) : undefined
  })

  const data = await readBody(res)

  if (!res.ok) {
    const message =
      (typeof data === 'string' && data) ||
      data?.message ||
      data?.error ||
      `Request failed (${res.status})`
    throw new ApiError(res.status, message, data)
  }

  return data
}

export async function login(email, password) {
  return request('/login', {
    method: 'POST',
    body: { email, password }
  })
}

export async function getProducts(token) {
  return request('/api/products', { token })
}

export async function createProduct(token, payload) {
  return request('/api/admin/products', {
    method: 'POST',
    token,
    body: payload
  })
}

export async function updateProduct(token, id, payload) {
  return request(`/api/admin/products/${id}`, {
    method: 'PUT',
    token,
    body: payload
  })
}

export async function deleteProduct(token, id) {
  return request(`/api/admin/products/${id}`, {
    method: 'DELETE',
    token
  })
}


export async function createOrder(token, payload) {
  return request('/api/orders', {
    method: 'POST',
    token,
    body: payload
  })
}

export async function getMe(token) {
  return request('/api/me', { token })
}

export async function getStoreSettings(token) {
  return request('/api/settings/store', { token })
}

export async function updateStoreSettings(token, payload) {
  return request('/api/settings/store', {
    method: 'PUT',
    token,
    body: payload
  })
}

export async function listUsers(token) {
  return request('/api/admin/users', { token })
}

export async function createUser(token, payload) {
  return request('/api/admin/users', {
    method: 'POST',
    token,
    body: payload
  })
}

export async function updateUser(token, id, payload) {
  return request(`/api/admin/users/${id}`, {
    method: 'PATCH',
    token,
    body: payload
  })
}

export async function deleteUser(token, id) {
  return request(`/api/admin/users/${id}`, {
    method: 'DELETE',
    token
  })
}

export async function listOrders(token, { date } = {}) {
  const qs = date ? `?date=${encodeURIComponent(date)}` : ''
  return request(`/api/orders${qs}`, { token })
}

export async function getOrder(token, id) {
  return request(`/api/orders/${id}`, { token })
}

export async function listLedger(token, { date } = {}) {
  const qs = date ? `?date=${encodeURIComponent(date)}` : ''
  return request(`/api/ledger${qs}`, { token })
}

export async function createLedgerEntry(token, payload) {
  return request('/api/ledger', {
    method: 'POST',
    token,
    body: payload
  })
}

export async function deleteLedgerEntry(token, id) {
  return request(`/api/ledger/${id}`, {
    method: 'DELETE',
    token
  })
}

