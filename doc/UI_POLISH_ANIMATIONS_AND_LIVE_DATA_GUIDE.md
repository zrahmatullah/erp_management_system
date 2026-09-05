# Panduan Pembersihan Data Dummy, Animasi Loading, dan Standarisasi Modal

Dokumen ini merangkum penyempurnaan menyeluruh pada antarmuka sistem Cafe ERP, mencakup pembersihan total data tiruan (dummy data) di frontend, penambahan animasi transisi halaman & status autentikasi, serta peremajaan visual popup modal dengan efek blur modern (frosted glass) dan transisi cubic-bezier yang halus.

---

## 1. Pembersihan Total Data Dummy Frontend (Live Database Integration)

Seluruh komponen antarmuka yang sebelumnya memiliki array statis, data percontohan lokal, atau rumus kalkulasi tiruan kini telah dihubungkan 100% ke REST API backend Go dan basis data PostgreSQL:

| Komponen / Halaman | Perubahan yang Diterapkan | Endpoint API Sumber |
| :--- | :--- | :--- |
| **Penggajian (`PayrollRun.vue`)** | Menghapus 8 baris karyawan hardcoded dan ringkasan statis. Seluruh data staf, gaji pokok, lembur, BPJS, PPh21, dan gaji bersih dimuat langsung dari database. Ringkasan total dihitung secara dinamis via computed property. Mengganti `window.alert()` dengan toast notifikasi modern. | `GET /api/v1/hris/payrolls` |
| **Laporan Penjualan (`SalesReport.vue`)** | Menghapus string hardcoded pada 4 kartu KPI (`Rp 89.750.000`, `Rp 68.500`, `1,310`). Metrik omset, nilai pesanan rata-rata, dan jumlah transaksi kini tersambung ke live stats. Opsi dropdown cabang dimuat dinamis dari data master cabang. Menghapus rumus dummy `120 - (idx * 15)` pada tabel top menu terlaris. | `GET /api/v1/dashboard/stats`<br>`GET /api/v1/master/branches`<br>`GET /api/v1/pos/products`<br>`GET /api/v1/pos/orders` |
| **Dashboard Utama (`DashboardView.vue`)** | Menghapus angka fallback tiruan (`45890000` dan `127`). Nilai KPI omset harian dan pesanan harian kini murni menampilkan nilai riil dari database. | `GET /api/v1/dashboard/stats` |
| **Jadwal Shift (`ShiftScheduleView.vue`)** | Panel "Definisi Shift" yang sebelumnya memiliki card statis (Pagi, Siang, Malam) kini dimuat dinamis dari database master shift beserta jam operasional dan warna badge-nya. | `GET /api/v1/master/shifts` |
| **Manajemen Peran (`RoleManagement.vue`)** | Daftar peran staf dan matriks permission dimuat dinamis dari database master peran. | `GET /api/v1/master/roles` |
| **Input Jurnal Keuangan (`JournalEntryForm.vue`)** | Mengganti dialog native `window.alert()` saat validasi ketidakseimbangan debit-kredit dan posting jurnal dengan `useNotificationStore()`. | `/api/v1/finance/journals` |
| **Pusat Data Master Super Admin (`MasterDataHub.vue`)** | Mengganti dialog konfirmasi native `window.confirm()` pada penghapusan item master dengan custom modal dialog `dialogStore.confirm({ type: 'danger' })`. | `/api/v1/master/*` |

---

## 2. Animasi Loading & Transisi Halaman (Page & Auth Animations)

Sesuai permintaan pengguna, antarmuka kini dilengkapi alur transisi dan animasi visual yang halus saat login, logout, maupun perpindahan antar-halaman / tab:

### A. Top Progress Bar (Global Route Navigation Loader)
- Dibuat komponen [`frontend/src/components/common/TopProgressBar.vue`](file:///d:/Project/cafe-erp-system/frontend/src/components/common/TopProgressBar.vue) dan state manager [`frontend/src/stores/loading.store.ts`](file:///d:/Project/cafe-erp-system/frontend/src/stores/loading.store.ts).
- Terpasang pada guard `router.beforeEach` dan `router.afterEach`.
- Menampilkan garis progres ramping (3px) dengan gradasi modern (`from-blue-600 via-sky-400 to-indigo-600`) dan efek cahaya glowing tip di sepanjang tepi atas layar saat bernavigasi antar halaman.

### B. Transisi Antar Halaman (`page-fade`)
- Pada [`frontend/src/components/layout/AppLayout.vue`](file:///d:/Project/cafe-erp-system/frontend/src/components/layout/AppLayout.vue), konten router dibungkus transisi Vue:
  ```html
  <router-view v-slot="{ Component }">
    <transition name="page-fade" mode="out-in">
      <component :is="Component" :key="$route.fullPath" />
    </transition>
  </router-view>
  ```
- Memberikan efek lembut fade & subtle glide vertikal (6px) saat berpindah menu/tab sehingga tidak terjadi loncatan tampilan yang kasar.

### C. Animasi Loading Login (`LoginView.vue`)
- Saat pengguna menekan tombol "Sign In" atau tombol "1-Click Demo Login", muncul loading overlay elegan di dalam kartu login:
  - Lingkaran ring gradasi berputar (`animate-spin`) mengelilingi ikon cangkir kopi (`animate-pulse`).
  - Teks status interaktif: *"Memproses Masuk... Memverifikasi kredensial dan menyiapkan sesi dashboard"*.
  - Transisi halus 350ms sebelum berpindah ke dashboard untuk memastikan user merasakan feedback visual yang nyaman.

### D. Animasi Loading Logout (`AppHeader.vue`)
- Saat opsi "Logout" dipilih dari menu profil:
  - Muncul dialog transparan full-screen berlatar frosted glass (`bg-slate-950/40 backdrop-blur-md`).
  - Menampilkan spinner merah lembut dengan ikon logout dan teks *"Mengakhiri Sesi... Menutup akses aman dan mengalihkan ke halaman login"*.
  - Menghapus token autentikasi, memunculkan toast pemberitahuan sesi berakhir, dan mengarahkan ke halaman login.

---

## 3. Standarisasi Backdrop Blur & Animasi Modal Popup

Latar belakang abu-abu pekat/kusam (`bg-slate-900/60 backdrop-blur-xs` atau `bg-gray-600/50`) telah diperbarui di seluruh popup modal sistem menjadi standar desain modern:

1. **Backdrop Frosted Glass Transparan**:
   - Menggunakan kelas `bg-slate-950/40 backdrop-blur-md` (atau `bg-black/35 backdrop-blur-md`).
   - Latar belakang halaman di balik modal terlihat buram halus (*soft blur*) dan tetap terbaca secara elegan tanpa warna abu-abu gelap pekat.

2. **Animasi Masuk & Keluar (Cubic-Bezier Spring Zoom)**:
   - Setiap modal dibungkus dalam `<Transition>` Vue dengan kurva easing `cubic-bezier(0.16, 1, 0.3, 1)`:
     - Saat terbuka: Muncul dari `opacity: 0, scale(0.94), translateY(8px)` menuju posisi normal dengan efek pantul halus.
     - Saat ditutup: Menghilang secara proporsional dan mulus tanpa jeda kaku.

3. **Komponen Modal yang Telah distandarisasi**:
   - `frontend/src/components/common/AppModal.vue` (Modal form serbaguna)
   - `frontend/src/components/common/GlobalDialogContainer.vue` (Dialog konfirmasi & prompt global)
   - `frontend/src/views/pos/PaymentModal.vue` (Modal pembayaran kasir multi-metode)
   - `frontend/src/views/pos/TableDetailModal.vue` (Modal detail & pemesanan meja interaktif)
   - `frontend/src/views/hris/EmployeeDetailModal.vue` (Modal rincian profil staf)
   - `frontend/src/views/hris/PayslipModal.vue` (Modal slip gaji resmi siap cetak)
   - `frontend/src/views/admin/MasterDataHub.vue` (Modal tambah/edit data master)

---

## 4. Verifikasi & Pengujian Kompilasi

Perintah kompilasi frontend telah diuji dan menghasilkan **Exit Code 0** tanpa kesalahan:
```bash
npm run build
```
- Total 4.523 modul berhasil ditransformasikan oleh Vite.
- Seluruh tipe TypeScript (`vue-tsc`) tervalidasi 100% valid.
