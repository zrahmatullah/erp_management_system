<template>
  <div class="min-h-screen w-full flex items-center justify-center p-4 sm:p-6 lg:p-8 transition-colors duration-250 relative select-none" style="background-color: var(--bg-content); color: var(--text-primary);">
    <!-- Ambient subtle background decorative glow -->
    <div class="absolute -top-32 -left-32 w-96 h-96 bg-blue-500/10 dark:bg-blue-600/15 rounded-full blur-3xl pointer-events-none"></div>
    <div class="absolute -bottom-32 -right-32 w-96 h-96 bg-indigo-500/10 dark:bg-indigo-600/15 rounded-full blur-3xl pointer-events-none"></div>

    <!-- Floating Theme Switcher (Top Right) -->
    <div class="absolute top-4 right-4 sm:top-6 sm:right-6 z-30">
      <div 
        role="radiogroup" 
        aria-label="Pilih tema tampilan login"
        class="p-1 rounded-xl border flex items-center text-xs font-semibold gap-1 backdrop-blur-md shadow-sm transition-all"
        style="background-color: var(--bg-primary); border-color: var(--border-color);"
      >
        <button
          type="button"
          role="radio"
          :aria-checked="appStore.themeMode === 'light'"
          aria-label="Mode Terang"
          @click="appStore.setTheme('light')"
          class="py-1 px-2 rounded-lg flex items-center justify-center gap-1.5 transition-all cursor-pointer select-none focus-visible:ring-2 focus-visible:ring-blue-500"
          :class="appStore.themeMode === 'light' ? 'bg-blue-50 text-blue-600 dark:bg-blue-900/40 dark:text-blue-400 font-bold shadow-xs' : 'text-slate-500 hover:text-slate-900 dark:text-slate-400 dark:hover:text-white'"
          title="Mode Terang (Putih & Biru)"
        >
          <Sun class="w-3.5 h-3.5 text-amber-500 shrink-0" />
          <span class="text-[11px] hidden sm:inline">Light</span>
        </button>

        <button
          type="button"
          role="radio"
          :aria-checked="appStore.themeMode === 'dark'"
          aria-label="Mode Gelap"
          @click="appStore.setTheme('dark')"
          class="py-1 px-2 rounded-lg flex items-center justify-center gap-1.5 transition-all cursor-pointer select-none focus-visible:ring-2 focus-visible:ring-blue-500"
          :class="appStore.themeMode === 'dark' ? 'bg-slate-800 text-blue-400 font-bold shadow-xs' : 'text-slate-500 hover:text-slate-900 dark:text-slate-400 dark:hover:text-white'"
          title="Mode Gelap (Hitam & Biru)"
        >
          <Moon class="w-3.5 h-3.5 text-blue-400 shrink-0" />
          <span class="text-[11px] hidden sm:inline">Dark</span>
        </button>

        <button
          type="button"
          role="radio"
          :aria-checked="appStore.themeMode === 'system'"
          aria-label="Ikuti Tema Sistem OS"
          @click="appStore.setTheme('system')"
          class="py-1 px-2 rounded-lg flex items-center justify-center gap-1.5 transition-all cursor-pointer select-none focus-visible:ring-2 focus-visible:ring-blue-500"
          :class="appStore.themeMode === 'system' ? 'bg-blue-600 text-white font-bold shadow-xs' : 'text-slate-500 hover:text-slate-900 dark:text-slate-400 dark:hover:text-white'"
          title="Otomatis mengikuti preferensi OS"
        >
          <Monitor class="w-3.5 h-3.5 shrink-0" :class="appStore.themeMode === 'system' ? 'text-white' : 'text-slate-400'" />
          <span class="text-[11px] hidden sm:inline">Auto</span>
        </button>
      </div>
    </div>

    <!-- Main Split Login Container Card -->
    <div 
      class="w-full max-w-5xl rounded-3xl shadow-2xl border grid grid-cols-1 lg:grid-cols-12 overflow-hidden transition-all duration-300 z-10 animate-in fade-in zoom-in-95 duration-200"
      style="background-color: var(--bg-primary); border-color: var(--border-color);"
    >
      
      <!-- Left Column: Enterprise Branding & Showcase (5 cols, hidden on small screens) -->
      <div 
        class="hidden lg:flex lg:col-span-5 p-8 sm:p-10 flex-col justify-between relative border-r overflow-hidden"
        style="background: linear-gradient(145deg, rgba(37,99,235,0.06) 0%, rgba(99,102,241,0.03) 100%); border-color: var(--border-color);"
      >
        <!-- Top Branding -->
        <div>
          <div class="flex items-center gap-3 mb-6">
            <div class="w-10 h-10 rounded-2xl bg-gradient-to-tr from-blue-600 via-blue-500 to-indigo-600 flex items-center justify-center text-white shadow-lg shadow-blue-500/25 shrink-0">
              <Coffee class="w-5 h-5" />
            </div>
            <div>
              <div class="flex items-center gap-2">
                <span class="font-black text-lg tracking-tight" style="color: var(--text-primary);">Cafe ERP</span>
                <span class="text-[10px] font-extrabold uppercase px-1.5 py-0.5 rounded-full bg-blue-100 dark:bg-blue-900/60 text-blue-700 dark:text-blue-300 border border-blue-200 dark:border-blue-800">
                  Pro v2.4
                </span>
              </div>
              <p class="text-xs font-medium" style="color: var(--text-muted);">Sistem ERP & POS Resto Terpadu</p>
            </div>
          </div>

          <!-- Headline & Pitch -->
          <h2 class="text-xl font-bold tracking-tight mb-2 leading-snug" style="color: var(--text-primary);">
            Kendali Penuh Seluruh Operasional Cafe & Resto Anda.
          </h2>
          <p class="text-xs leading-relaxed mb-6" style="color: var(--text-muted);">
            Solusi cloud terintegrasi untuk kasir POS, manajemen stok & resep otomatis, presensi karyawan, hingga pelaporan laba rugi real-time.
          </p>

          <!-- 4 Core Enterprise Features Highlight -->
          <div class="space-y-3">
            <div class="p-3 rounded-xl border flex items-start gap-3 transition-colors" style="background-color: var(--bg-content); border-color: var(--border-color);">
              <div class="w-7 h-7 rounded-lg bg-blue-100 dark:bg-blue-900/50 text-blue-600 dark:text-blue-400 flex items-center justify-center shrink-0 mt-0.5">
                <CreditCard class="w-4 h-4" />
              </div>
              <div class="min-w-0">
                <h4 class="text-xs font-bold leading-tight" style="color: var(--text-primary);">POS Kasir & KDS Real-Time</h4>
                <p class="text-[11px] leading-tight mt-0.5" style="color: var(--text-muted);">Transaksi kilat, split-bill, denah meja, & pesanan langsung terkirim ke dapur.</p>
              </div>
            </div>

            <div class="p-3 rounded-xl border flex items-start gap-3 transition-colors" style="background-color: var(--bg-content); border-color: var(--border-color);">
              <div class="w-7 h-7 rounded-lg bg-emerald-100 dark:bg-emerald-900/50 text-emerald-600 dark:text-emerald-400 flex items-center justify-center shrink-0 mt-0.5">
                <Package class="w-4 h-4" />
              </div>
              <div class="min-w-0">
                <h4 class="text-xs font-bold leading-tight" style="color: var(--text-primary);">Inventori & Resep Otomatis</h4>
                <p class="text-[11px] leading-tight mt-0.5" style="color: var(--text-muted);">Potong stok bahan baku presisi berdasarkan porsi menu yang terjual.</p>
              </div>
            </div>

            <div class="p-3 rounded-xl border flex items-start gap-3 transition-colors" style="background-color: var(--bg-content); border-color: var(--border-color);">
              <div class="w-7 h-7 rounded-lg bg-purple-100 dark:bg-purple-900/50 text-purple-600 dark:text-purple-400 flex items-center justify-center shrink-0 mt-0.5">
                <TrendingUp class="w-4 h-4" />
              </div>
              <div class="min-w-0">
                <h4 class="text-xs font-bold leading-tight" style="color: var(--text-primary);">Laporan Keuangan & HPP</h4>
                <p class="text-[11px] leading-tight mt-0.5" style="color: var(--text-muted);">Kalkulasi margin profit, jurnal akuntansi, dan ringkasan omset harian.</p>
              </div>
            </div>

            <div class="p-3 rounded-xl border flex items-start gap-3 transition-colors" style="background-color: var(--bg-content); border-color: var(--border-color);">
              <div class="w-7 h-7 rounded-lg bg-amber-100 dark:bg-amber-900/50 text-amber-600 dark:text-amber-400 flex items-center justify-center shrink-0 mt-0.5">
                <Users class="w-4 h-4" />
              </div>
              <div class="min-w-0">
                <h4 class="text-xs font-bold leading-tight" style="color: var(--text-primary);">HRIS & Shift Staf</h4>
                <p class="text-[11px] leading-tight mt-0.5" style="color: var(--text-muted);">Jadwal kerja barista/waiter, absensi live, dan penggajian otomatis.</p>
              </div>
            </div>
          </div>
        </div>

        <!-- Trust Badges Footer -->
        <div class="pt-6 border-t mt-6 flex items-center justify-between text-[11px]" style="border-color: var(--border-color); color: var(--text-muted);">
          <div class="flex items-center gap-1.5 font-semibold">
            <ShieldCheck class="w-4 h-4 text-emerald-500" />
            <span>256-Bit SSL Secured</span>
          </div>
          <div class="flex items-center gap-1.5 font-semibold">
            <CheckCircle2 class="w-4 h-4 text-blue-500" />
            <span>Multi-Outlet Sync</span>
          </div>
        </div>
      </div>

      <!-- Right Column: Login Form & Demo Quick Login (7 cols on lg) -->
      <div class="col-span-1 lg:col-span-7 p-6 sm:p-10 lg:p-12 flex flex-col justify-center">
        <!-- Mobile Header (Visible on small screens) -->
        <div class="lg:hidden flex items-center gap-2.5 mb-6">
          <div class="w-9 h-9 rounded-xl bg-gradient-to-tr from-blue-600 to-indigo-600 flex items-center justify-center text-white shadow-md shadow-blue-500/25 shrink-0">
            <Coffee class="w-4.5 h-4.5" />
          </div>
          <div>
            <h1 class="font-black text-lg tracking-tight leading-none" style="color: var(--text-primary);">Cafe ERP Pro</h1>
            <p class="text-xs font-medium mt-0.5" style="color: var(--text-muted);">Enterprise Management System</p>
          </div>
        </div>

        <!-- Form Intro Title -->
        <div class="mb-5">
          <h2 class="text-2xl font-black tracking-tight" style="color: var(--text-primary);">
            Selamat Datang Kembali
          </h2>
          <p class="text-xs mt-1" style="color: var(--text-muted);">
            Silakan masukkan kredensial akun Anda untuk mengakses sistem manajemen cafe.
          </p>
        </div>

        <!-- Quick Demo Profiles (1-Click Fill) -->
        <div class="mb-5">
          <div class="text-[11px] font-bold uppercase tracking-wider mb-2 flex items-center justify-between" style="color: var(--text-muted);">
            <span>Pilih Akun Demo Cepat</span>
            <span class="text-[10px] font-medium text-blue-500 dark:text-blue-400">1-Klik Langsung Isi</span>
          </div>

          <div class="grid grid-cols-2 sm:grid-cols-4 gap-2">
            <button
              type="button"
              @click="quickLogin('admin@cafe-erp.com', 'Admin@123', 'Super Admin')"
              class="btn-press p-2 rounded-xl border flex flex-col items-center justify-center text-center transition-all cursor-pointer group"
              style="background-color: var(--bg-content); border-color: var(--border-color);"
              title="Masuk sebagai Super Administrator"
            >
              <div class="w-6 h-6 rounded-lg bg-blue-100 dark:bg-blue-900/60 text-blue-600 dark:text-blue-300 flex items-center justify-center text-xs font-black mb-1 group-hover:scale-105 transition-transform">
                SA
              </div>
              <span class="text-[11px] font-bold truncate max-w-full" style="color: var(--text-primary);">Super Admin</span>
              <span class="text-[9.5px] truncate max-w-full" style="color: var(--text-muted);">Semua Akses</span>
            </button>

            <button
              type="button"
              @click="quickLogin('kasir@cafe-erp.com', 'Password@123', 'Kasir POS')"
              class="btn-press p-2 rounded-xl border flex flex-col items-center justify-center text-center transition-all cursor-pointer group"
              style="background-color: var(--bg-content); border-color: var(--border-color);"
              title="Masuk sebagai Kasir POS"
            >
              <div class="w-6 h-6 rounded-lg bg-emerald-100 dark:bg-emerald-900/60 text-emerald-600 dark:text-emerald-300 flex items-center justify-center text-xs font-black mb-1 group-hover:scale-105 transition-transform">
                KS
              </div>
              <span class="text-[11px] font-bold truncate max-w-full" style="color: var(--text-primary);">Kasir POS</span>
              <span class="text-[9.5px] truncate max-w-full" style="color: var(--text-muted);">Ke Halaman POS</span>
            </button>

            <button
              type="button"
              @click="quickLogin('manager@cafe-erp.com', 'Password@123', 'Store Manager')"
              class="btn-press p-2 rounded-xl border flex flex-col items-center justify-center text-center transition-all cursor-pointer group"
              style="background-color: var(--bg-content); border-color: var(--border-color);"
              title="Masuk sebagai Store Manager"
            >
              <div class="w-6 h-6 rounded-lg bg-amber-100 dark:bg-amber-900/60 text-amber-600 dark:text-amber-300 flex items-center justify-center text-xs font-black mb-1 group-hover:scale-105 transition-transform">
                MN
              </div>
              <span class="text-[11px] font-bold truncate max-w-full" style="color: var(--text-primary);">Manager</span>
              <span class="text-[9.5px] truncate max-w-full" style="color: var(--text-muted);">Stok & Laporan</span>
            </button>

            <button
              type="button"
              @click="quickLogin('finance@cafe-erp.com', 'Password@123', 'Finance Staff')"
              class="btn-press p-2 rounded-xl border flex flex-col items-center justify-center text-center transition-all cursor-pointer group"
              style="background-color: var(--bg-content); border-color: var(--border-color);"
              title="Masuk sebagai Staf Keuangan"
            >
              <div class="w-6 h-6 rounded-lg bg-purple-100 dark:bg-purple-900/60 text-purple-600 dark:text-purple-300 flex items-center justify-center text-xs font-black mb-1 group-hover:scale-105 transition-transform">
                FN
              </div>
              <span class="text-[11px] font-bold truncate max-w-full" style="color: var(--text-primary);">Finance</span>
              <span class="text-[9.5px] truncate max-w-full" style="color: var(--text-muted);">Jurnal & Kas</span>
            </button>
          </div>
        </div>

        <!-- Divider -->
        <div class="relative my-3 text-center">
          <div class="absolute inset-0 flex items-center">
            <div class="w-full border-t" style="border-color: var(--border-color);"></div>
          </div>
          <span class="relative px-3 text-[11px] font-bold uppercase tracking-widest" style="background-color: var(--bg-primary); color: var(--text-muted);">
            Atau Gunakan Akun Terdaftar
          </span>
        </div>

        <!-- Lockout Warning Alert if Rate Limited -->
        <div 
          v-if="isLockedOut" 
          class="p-3.5 mb-4 rounded-xl border flex items-center gap-3 bg-amber-50 dark:bg-amber-950/40 border-amber-200 dark:border-amber-900 text-amber-800 dark:text-amber-300 text-xs font-medium shake-error"
        >
          <AlertTriangle class="w-5 h-5 text-amber-600 dark:text-amber-400 shrink-0" />
          <div class="min-w-0">
            <div class="font-bold">Akses Terkunci Sementara</div>
            <div class="text-[11px] mt-0.5">
              Terlalu banyak percobaan login gagal berturut-turut. Demi keamanan, silakan tunggu 
              <span class="font-black underline">{{ lockoutSeconds }} detik</span> sebelum mencoba kembali.
            </div>
          </div>
        </div>

        <!-- Error Alert (Generic Message for Security) -->
        <div 
          v-else-if="errorMessage" 
          class="p-3 mb-4 rounded-xl border flex items-center gap-3 bg-rose-50 dark:bg-rose-950/40 border-rose-200 dark:border-rose-900 text-rose-700 dark:text-rose-300 text-xs font-medium shake-error"
        >
          <AlertCircle class="w-4 h-4 text-rose-500 shrink-0" />
          <span class="min-w-0 flex-1">{{ errorMessage }}</span>
        </div>

        <!-- Main Form -->
        <form @submit.prevent="handleLogin" class="space-y-4" novalidate>
          <!-- Email Field -->
          <div>
            <label for="login-email" class="block text-xs font-bold mb-1.5" style="color: var(--text-secondary);">
              Email Pengguna <span class="text-rose-500">*</span>
            </label>
            <div 
              class="rounded-xl px-3.5 py-2.5 flex items-center gap-3 border transition-all"
              :class="emailError ? 'border-rose-500 ring-1 ring-rose-500/20' : 'focus-within:border-blue-500 focus-within:ring-2 focus-within:ring-blue-500/20'"
              style="background-color: var(--input-bg); border-color: emailError ? '#ef4444' : var(--input-border);"
            >
              <Mail class="w-4 h-4 shrink-0" :class="emailError ? 'text-rose-500' : 'text-slate-400 dark:text-slate-500'" />
              <input
                id="login-email"
                v-model="email"
                type="email"
                autocomplete="email"
                required
                :disabled="loading || isLockedOut"
                placeholder="nama@cafe-erp.com"
                @input="clearEmailError"
                class="w-full bg-transparent p-0 text-xs font-semibold outline-none border-none ring-0 focus:ring-0"
                style="color: var(--input-text);"
              />
            </div>
            <p v-if="emailError" class="text-[11px] text-rose-500 font-semibold mt-1 flex items-center gap-1">
              <span>•</span> {{ emailError }}
            </p>
          </div>

          <!-- Password Field -->
          <div>
            <label for="login-password" class="block text-xs font-bold mb-1.5" style="color: var(--text-secondary);">
              Kata Sandi <span class="text-rose-500">*</span>
            </label>
            <div 
              class="rounded-xl px-3.5 py-2.5 flex items-center gap-3 border transition-all"
              :class="passwordError ? 'border-rose-500 ring-1 ring-rose-500/20' : 'focus-within:border-blue-500 focus-within:ring-2 focus-within:ring-blue-500/20'"
              style="background-color: var(--input-bg); border-color: passwordError ? '#ef4444' : var(--input-border);"
            >
              <Lock class="w-4 h-4 shrink-0" :class="passwordError ? 'text-rose-500' : 'text-slate-400 dark:text-slate-500'" />
              <input
                id="login-password"
                v-model="password"
                :type="showPassword ? 'text' : 'password'"
                autocomplete="current-password"
                required
                :disabled="loading || isLockedOut"
                placeholder="••••••••••••"
                @input="clearPasswordError"
                class="w-full bg-transparent p-0 text-xs font-semibold outline-none border-none ring-0 focus:ring-0"
                style="color: var(--input-text);"
              />
              <button
                type="button"
                @click="showPassword = !showPassword"
                class="text-slate-400 hover:text-slate-700 dark:hover:text-slate-200 p-0.5 cursor-pointer shrink-0 transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500 rounded"
                :aria-label="showPassword ? 'Sembunyikan kata sandi' : 'Tampilkan kata sandi'"
                :title="showPassword ? 'Sembunyikan kata sandi' : 'Tampilkan kata sandi'"
              >
                <EyeOff v-if="showPassword" class="w-4 h-4" />
                <Eye v-else class="w-4 h-4" />
              </button>
            </div>
            <p v-if="passwordError" class="text-[11px] text-rose-500 font-semibold mt-1 flex items-center gap-1">
              <span>•</span> {{ passwordError }}
            </p>
          </div>

          <!-- Remember Me & Forgot Password -->
          <div class="flex items-center justify-between text-xs pt-0.5">
            <label class="flex items-center gap-2 cursor-pointer select-none" style="color: var(--text-muted);">
              <input 
                v-model="rememberMe" 
                type="checkbox" 
                :disabled="loading || isLockedOut"
                class="w-4 h-4 rounded border-slate-300 dark:border-slate-700 text-blue-600 focus:ring-blue-500 cursor-pointer" 
              />
              <span class="text-xs font-medium">Ingat Akun Saya</span>
            </label>
            <button 
              type="button" 
              @click="showForgotModal = true" 
              class="text-blue-600 dark:text-blue-400 font-bold hover:underline text-xs cursor-pointer focus-visible:outline-none"
            >
              Lupa Kata Sandi?
            </button>
          </div>

          <!-- Submit Button -->
          <button
            type="submit"
            :disabled="loading || isLockedOut"
            class="btn-press w-full py-3 px-4 bg-blue-600 hover:bg-blue-700 active:bg-blue-800 text-white font-bold rounded-xl shadow-lg shadow-blue-600/25 transition-all flex items-center justify-center gap-2 disabled:opacity-50 disabled:cursor-not-allowed cursor-pointer text-xs sm:text-sm tracking-wide mt-2"
          >
            <span v-if="loading" class="animate-spin w-4 h-4 border-2 border-white border-t-transparent rounded-full"></span>
            <LogIn v-else class="w-4 h-4" />
            <span>{{ loading ? 'Mengotentikasi Sesi...' : isLockedOut ? `Terkunci (${lockoutSeconds}s)` : 'Masuk ke Sistem ERP' }}</span>
          </button>
        </form>

        <!-- Footer Notice & Copyright -->
        <div class="mt-8 pt-4 border-t flex flex-col sm:flex-row items-center justify-between text-[11px] gap-2" style="border-color: var(--border-color); color: var(--text-muted);">
          <span>&copy; {{ new Date().getFullYear() }} Cafe Monitoring & ERP System. All rights reserved.</span>
          <div class="flex items-center gap-3">
            <button type="button" @click="handleSecurityPolicy" class="hover:underline cursor-pointer">Kebijakan Keamanan</button>
            <span>•</span>
            <button type="button" @click="showForgotModal = true" class="hover:underline cursor-pointer">Bantuan IT</button>
          </div>
        </div>

      </div>
    </div>

    <!-- Modal: Bantuan Lupa Kata Sandi (A11y Friendly) -->
    <Teleport to="body">
      <transition name="fade">
        <div 
          v-if="showForgotModal" 
          class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-950/60 backdrop-blur-xs"
          @click.self="showForgotModal = false"
        >
          <div 
            class="w-full max-w-md rounded-2xl p-6 shadow-2xl border animate-in zoom-in-95 duration-200"
            style="background-color: var(--bg-primary); border-color: var(--border-color); color: var(--text-primary);"
            role="dialog"
            aria-modal="true"
            aria-labelledby="forgot-title"
          >
            <div class="flex items-center justify-between pb-3 border-b" style="border-color: var(--border-color);">
              <div class="flex items-center gap-2.5">
                <div class="w-8 h-8 rounded-lg bg-blue-100 dark:bg-blue-900/50 text-blue-600 dark:text-blue-400 flex items-center justify-center">
                  <ShieldCheck class="w-4 h-4" />
                </div>
                <h3 id="forgot-title" class="font-bold text-base">Pemulihan Akun & Akses</h3>
              </div>
              <button 
                @click="showForgotModal = false" 
                class="p-1 rounded-lg text-slate-400 hover:text-slate-700 dark:hover:text-slate-200 cursor-pointer"
                aria-label="Tutup dialog"
              >
                ✕
              </button>
            </div>

            <div class="py-4 space-y-3 text-xs leading-relaxed" style="color: var(--text-secondary);">
              <p>
                Demi standar keamanan data audit perbankan & ERP, reset kata sandi tidak dikirimkan melalui email publik otomatis.
              </p>
              <div class="p-3 rounded-xl border" style="background-color: var(--bg-content); border-color: var(--border-color);">
                <div class="font-bold mb-1" style="color: var(--text-primary);">Langkah Pemulihan:</div>
                <ul class="list-disc list-inside space-y-1" style="color: var(--text-muted);">
                  <li>Hubungi <strong>Super Administrator</strong> atau <strong>IT Security Dept</strong> outlet Anda.</li>
                  <li>Sebutkan NIK Karyawan atau email terdaftar untuk verifikasi dua arah.</li>
                  <li>Admin akan menerbitkan token reset kredensial sementara melalui menu <em>Master Data &gt; Pengguna</em>.</li>
                </ul>
              </div>
              <div class="text-[11px]" style="color: var(--text-muted);">
                Hotline Bantuan IT Outlet: <span class="font-bold text-blue-600 dark:text-blue-400">+62 812-9988-7766</span> (Jam Operasional Cafe)
              </div>
            </div>

            <div class="pt-3 border-t flex justify-end" style="border-color: var(--border-color);">
              <button 
                type="button" 
                @click="showForgotModal = false" 
                class="btn-press px-4 py-2 rounded-xl text-xs font-bold text-white bg-blue-600 hover:bg-blue-700 transition-colors cursor-pointer"
              >
                Saya Mengerti
              </button>
            </div>
          </div>
        </div>
      </transition>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import {
  Coffee,
  Mail,
  Lock,
  Eye,
  EyeOff,
  Sun,
  Moon,
  Monitor,
  LogIn,
  AlertCircle,
  AlertTriangle,
  ShieldCheck,
  CheckCircle2,
  CreditCard,
  Package,
  TrendingUp,
  Users
} from 'lucide-vue-next'
import { useAuthStore } from '@/stores/auth.store'
import { useAppStore } from '@/stores/app.store'
import { useNotificationStore } from '@/stores/notification.store'
import axios from 'axios'

const router = useRouter()
const authStore = useAuthStore()
const appStore = useAppStore()
const notifyStore = useNotificationStore()

// State
const email = ref('')
const password = ref('')
const rememberMe = ref(false)
const showPassword = ref(false)
const loading = ref(false)
const errorMessage = ref('')
const emailError = ref('')
const passwordError = ref('')
const showForgotModal = ref(false)

// Rate Limiting & Brute Force Protection (Client-side)
const failedAttempts = ref(0)
const isLockedOut = ref(false)
const lockoutSeconds = ref(0)
let lockoutTimer: any = null

const REMEMBER_KEY = 'cafe_erp_remembered_email'

onMounted(() => {
  // Check remembered email
  const savedEmail = localStorage.getItem(REMEMBER_KEY)
  if (savedEmail) {
    email.value = savedEmail
    rememberMe.value = true
  }
})

onUnmounted(() => {
  if (lockoutTimer) clearInterval(lockoutTimer)
})

const clearEmailError = () => {
  emailError.value = ''
  errorMessage.value = ''
}

const clearPasswordError = () => {
  passwordError.value = ''
  errorMessage.value = ''
}

const validateForm = (): boolean => {
  let valid = true
  emailError.value = ''
  passwordError.value = ''
  errorMessage.value = ''

  const trimmedEmail = email.value.trim()
  if (!trimmedEmail) {
    emailError.value = 'Email wajib diisi.'
    valid = false
  } else {
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
    if (!emailRegex.test(trimmedEmail)) {
      emailError.value = 'Format alamat email tidak valid (contoh: nama@cafe-erp.com).'
      valid = false
    }
  }

  if (!password.value) {
    passwordError.value = 'Kata sandi wajib diisi.'
    valid = false
  } else if (password.value.length < 6) {
    passwordError.value = 'Kata sandi minimal 6 karakter.'
    valid = false
  }

  return valid
}

const handleLogin = async () => {
  if (isLockedOut.value) return

  if (!validateForm()) {
    return
  }

  loading.value = true
  errorMessage.value = ''

  try {
    const res = await axios.post('/api/v1/auth/login', {
      email: email.value.trim(),
      password: password.value
    })

    if (res.data?.data) {
      const data = res.data.data
      const userObj = {
        id: data.user.id,
        email: data.user.email,
        firstName: data.user.full_name || 'Pengguna',
        lastName: '',
        roleId: data.role || data.user.role_name || 'Super Admin'
      }

      authStore.login(
        data.token,
        data.refresh_token || data.token,
        userObj,
        data.permissions || []
      )

      // Handle remember me
      if (rememberMe.value) {
        localStorage.setItem(REMEMBER_KEY, email.value.trim())
      } else {
        localStorage.removeItem(REMEMBER_KEY)
      }

      // Reset failed attempts upon success
      failedAttempts.value = 0

      notifyStore.success(`Selamat datang kembali, ${userObj.firstName}!`, `Peran: ${userObj.roleId}`)

      // Role-Based Smart Navigation
      const roleLower = String(userObj.roleId).toLowerCase()
      if (roleLower.includes('kasir') || roleLower.includes('cashier')) {
        await router.push('/pos')
      } else if (roleLower.includes('finance') || roleLower.includes('keuangan')) {
        await router.push('/finance')
      } else if (roleLower.includes('gudang') || roleLower.includes('inventory')) {
        await router.push('/inventory')
      } else {
        await router.push('/dashboard')
      }
    } else {
      throw new Error('Respons otentikasi tidak valid dari server')
    }
  } catch (err: any) {
    failedAttempts.value += 1

    // Anti-Brute Force Lockout trigger
    if (failedAttempts.value >= 5) {
      triggerLockout(30)
      return
    }

    // Generic OWASP-compliant error message to prevent user enumeration
    const genericMsg = 'Kredensial login tidak cocok. Silakan periksa kembali email dan kata sandi Anda.'
    errorMessage.value = genericMsg
    notifyStore.error(genericMsg, 'Gagal Masuk')
  } finally {
    loading.value = false
  }
}

const triggerLockout = (seconds: number) => {
  isLockedOut.value = true
  lockoutSeconds.value = seconds
  errorMessage.value = `Terlalu banyak percobaan gagal. Akses dibatasi sementara selama ${seconds} detik.`
  notifyStore.error(errorMessage.value, 'Keamanan Sistem')

  if (lockoutTimer) clearInterval(lockoutTimer)
  lockoutTimer = setInterval(() => {
    lockoutSeconds.value -= 1
    if (lockoutSeconds.value <= 0) {
      clearInterval(lockoutTimer)
      isLockedOut.value = false
      failedAttempts.value = 0
      errorMessage.value = ''
    }
  }, 1000)
}

const quickLogin = (e: string, p: string, _roleName?: string) => {
  if (isLockedOut.value) return
  email.value = e
  password.value = p
  clearEmailError()
  clearPasswordError()
  handleLogin()
}

const handleSecurityPolicy = () => {
  notifyStore.info('Sistem mematuhi standar enkripsi TLS 1.3, hashing bcrypt salted, dan proteksi sesi OWASP ASVS V3.', 'Kebijakan Keamanan')
}
</script>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
