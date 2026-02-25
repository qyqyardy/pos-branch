export const receiptConfig = {
  store: {
    name: 'WARKOP',
    tagline: 'Point of Sale',
    addressLines: [
      'Jl. Contoh No. 1',
      'Kota Kamu'
    ],
    phone: '08xx-xxxx-xxxx'
  },

  // Thermal printer notes:
  // - Paper 58mm biasanya printable width aman ~48mm.
  // - Paper 80mm biasanya printable width aman ~72mm.
  printer: {
    receipt: {
      paperWidthMm: 58,
      contentWidthMm: 48,
      paddingMm: 3,
      baseFontPx: 11,
      titleFontPx: 16
    },
    kitchen: {
      paperWidthMm: 58,
      contentWidthMm: 48,
      paddingMm: 3,
      baseFontPx: 12,
      titleFontPx: 18
    }
  },

  receipt: {
    showAddress: true,
    showPhone: true,
    showLogo: true,
    tax: {
      enabled: false,
      rate: 0.1,
      label: 'Pajak'
    },
    // Set true kalau mau tampilkan "harga/item" di bawah tiap item (struk jadi lebih panjang).
    showItemUnitPrice: false,
    footerLines: [
      'Terima kasih - Silahkan datang lagi.'
    ]
  },

  kitchen: {
    stationName: 'Dapur 1',
    footerLines: [
      'Tandai bila selesai.'
    ]
  }
}
