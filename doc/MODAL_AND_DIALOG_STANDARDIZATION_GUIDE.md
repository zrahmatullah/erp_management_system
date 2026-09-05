# Panduan Standarisasi UI Popup & Modal Dialog

Dokumen ini menjelaskan standarisasi sistem antarmuka untuk seluruh popup dialog, modal form, dan konfirmasi aksi di Cafe ERP System, menggantikan 100% dialog bawaan peramban (*browser native dialog* seperti `prompt()`, `confirm()`, dan `alert()`).

---

## 1. Latar Belakang & Masalah

Sebelum standarisasi, beberapa aksi masih menggunakan fungsi bawaan browser (`window.prompt()` pada tambah meja, `window.confirm()` pada hapus data master, dan `window.alert()` pada notifikasi form). Dialog native browser memiliki banyak kelemahan:
* **Tampilan Kaku**: Menggunakan dialog abu-abu/hitam bawaan OS yang merusak estetika antarmuka modern Cafe ERP.
* **Tidak Mendukung Styling**: Tidak dapat menggunakan warna merek, font Inter/Plus Jakarta Sans, maupun icon vektor `lucide-vue-next`.
* **Pengalaman Pengguna (UX) Rendah**: Membekukan thread eksekusi browser dan tidak mendukung validasi input interaktif (dropdown, radio, switch).

---

## 2. Arsitektur Komponen Modal Terstandarisasi

Sistem popup Cafe ERP kini distandarisasi menjadi 3 layer komponen visual:

```
                          [ Aksi Pengguna ]
                                  │
      ┌───────────────────────────┼───────────────────────────┐
      ▼                           ▼                           ▼
[ Form Kompleks ]      [ Dialog Konfirmasi/Prompt ]    [ Notifikasi Aksi ]
  (Tambah Meja,            (Hapus Data, Tanya,           (Sukses Simpan,
   Detail, Master)          Persetujuan Aksi)             Peringatan Jurnal)
      │                           │                           │
      ▼                           ▼                           ▼
[ AppModal.vue ]      [ GlobalDialogContainer.vue ]   [ ToastContainer.vue ]
(Modal Form Card)     (dialogStore.confirm/prompt)    (notifyStore.success/error)
```

---

## 3. Komponen Modal yang Diterapkan

### A. Modal Form Terstandarisasi (`AppModal.vue`)
* **Lokasi**: [`frontend/src/components/common/AppModal.vue`](file:///d:/Project/cafe-erp-system/frontend/src/components/common/AppModal.vue)
* **Fitur**:
  - `Teleport to="body"` dengan backdrop gelap dan blur halus (`bg-slate-900/60 backdrop-blur-xs`).
  - Animasi transisi masuk halus (*smooth fade-in & zoom-in-95*).
  - Header dengan icon badge Lucide, judul tebal, subjudul penjelas, dan tombol silang `X`.
  - Area konten *scrollable* (`max-h-[75vh]`) agar responsif di layar tablet maupun POS terminal.
  - Footer tombol aksi terstandarisasi: tombol **Batal** (abu-abu Slate 100) dan tombol **Submit** (biru Blue 600 dengan efek *glow shadow*).

#### Contoh Implementasi: Modal Tambah Meja (`TableManagementView.vue`)
Menggantikan input `prompt()` satu baris dengan form terstruktur:
* **Nomor / Kode Meja**: Input teks dengan autofocus (contoh: `T-13`, `VIP-02`).
* **Kapasitas Kursi**: Pilihan 2 Kursi, 4 Kursi, 6 Kursi, 8 Kursi, hingga 10 Kursi.
* **Zona Cafe**: Dropdown pilihan zona (Lantai 1, Lantai 2, Outdoor).
* **Status Awal**: Tersedia (Available) atau Reservasi (Reserved).

---

### B. Global Dialog Container (`GlobalDialogContainer.vue`)
* **Lokasi**: [`frontend/src/components/common/GlobalDialogContainer.vue`](file:///d:/Project/cafe-erp-system/frontend/src/components/common/GlobalDialogContainer.vue)
* **Store Pengendali**: [`frontend/src/stores/dialog.store.ts`](file:///d:/Project/cafe-erp-system/frontend/src/stores/dialog.store.ts)
* **Pemasangan**: Terpasang di root aplikasi [`App.vue`](file:///d:/Project/cafe-erp-system/frontend/src/App.vue).

#### 1. Konfirmasi Hapus / Aksi Kritis (`dialogStore.confirm`)
Menggantikan `window.confirm()` dengan dialog visual elegan dan tipe warna dinamis (`danger`, `warning`, `info`):
```typescript
import { useDialogStore } from '@/stores/dialog.store'
const dialog = useDialogStore()

// Contoh Konfirmasi Hapus di MasterDataHub.vue:
const confirmed = await dialog.confirm({
  title: 'Hapus Data Cabang',
  message: 'Apakah Anda yakin ingin menghapus data cabang ini? Tindakan ini tidak dapat dibatalkan.',
  type: 'danger',
  confirmText: 'Ya, Hapus Data',
  cancelText: 'Batal'
})
if (confirmed) {
  // Lakukan panggilan API hapus
}
```

#### 2. Input Prompt Cepat (`dialogStore.prompt`)
Menggantikan `window.prompt()` dengan modal berinput field ber-focus otomatis:
```typescript
const result = await dialog.prompt({
  title: 'Alasan Penolakan PO',
  message: 'Masukkan alasan pembatalan atau revisi Purchase Order:',
  placeholder: 'e.g. Harga bahan baku tidak sesuai',
  confirmText: 'Kirim Alasan'
})
if (result) {
  // Proses penolakan
}
```

---

### C. Toast Notifikasi Global (`ToastContainer.vue`)
Menggantikan seluruh `window.alert()` untuk pesan-pesan status aplikasi:
* `notifyStore.success('Pesan sukses', 'Judul')` (Hijau)
* `notifyStore.warning('Pesan peringatan', 'Judul')` (Kuning/Amber)
* `notifyStore.error('Pesan kesalahan', 'Judul')` (Merah/Rose)
* `notifyStore.info('Pesan informasi', 'Judul')` (Biru)

---

## 4. Hasil Verifikasi

1. **Pencarian Kode Sumber**: Pencarian global ripgrep untuk fungsi `prompt(`, `confirm(`, dan `alert(` di seluruh direktori `frontend/src` menghasilkan **0 pemanggilan native**.
2. **Kompilasi TypeScript & Vite**: Perintah `npm run build` berhasil dijalankan dengan **Exit code 0**.

