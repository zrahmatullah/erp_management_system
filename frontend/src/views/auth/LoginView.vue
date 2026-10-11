<template>
  <!-- Fullscreen Zero-Scroll Viewport Container -->
  <div 
    class="fixed inset-0 w-full h-full overflow-hidden flex items-center justify-center p-3 sm:p-5 select-none bg-[#EEF2F6] dark:bg-[#070D18]"
    style="color: var(--text-primary);"
  >
    <!-- Contained subtle ambient background glow -->
    <div class="absolute inset-0 overflow-hidden pointer-events-none">
      <div class="absolute -top-32 -left-32 w-96 h-96 bg-blue-400/10 dark:bg-blue-600/15 rounded-full blur-3xl"></div>
      <div class="absolute -bottom-32 -right-32 w-96 h-96 bg-indigo-400/10 dark:bg-indigo-600/15 rounded-full blur-3xl"></div>
    </div>

    <!-- Main 2-Panel Card Container matching uploaded design -->
    <div class="w-full max-w-[1040px] grid grid-cols-1 lg:grid-cols-12 gap-5 z-10 animate-in fade-in zoom-in-95 duration-200">
      
      <!-- ========================================================
           LEFT PANEL: Enterprise Resto Overview & Shift Operations
           ======================================================== -->
      <div 
        class="lg:col-span-7 bg-white dark:bg-slate-900 rounded-3xl border border-slate-200/80 dark:border-slate-800 shadow-xl p-5 sm:p-6 flex flex-col justify-between"
      >
        <!-- Header: Logo, Brand & Cloud Sync Status -->
        <div class="flex items-center justify-between pb-2">
          <div class="flex items-center gap-3">
            <div class="w-10 h-10 rounded-xl bg-slate-950 dark:bg-slate-800 text-white flex items-center justify-center shadow-md shrink-0">
              <Coffee class="w-5 h-5 text-white" />
            </div>
            <div>
              <div class="flex items-center gap-2">
                <span class="text-base font-black tracking-tight text-slate-900 dark:text-white">KopiERP</span>
                <span class="text-[9px] font-extrabold uppercase px-2 py-0.5 rounded-md bg-blue-100 dark:bg-blue-900/60 text-blue-700 dark:text-blue-300 tracking-wider">
                  ENTERPRISE PRO
                </span>
              </div>
              <p class="text-[11px] text-slate-500 dark:text-slate-400">Sistem Manajemen POS, Resto & Rantai Suplai</p>
            </div>
          </div>

          <!-- Cloud Sync Pill -->
          <div class="px-2.5 py-1 rounded-full bg-slate-100 dark:bg-slate-800 border border-slate-200/60 dark:border-slate-700 flex items-center gap-1.5 text-[11px] font-semibold text-slate-600 dark:text-slate-300 shrink-0">
            <span class="w-2 h-2 rounded-full bg-blue-500 animate-pulse"></span>
            <span>Sync Cloud Aktif</span>
          </div>
        </div>

        <!-- Headline & Description -->
        <div class="my-2.5">
          <h1 class="text-xl sm:text-[22px] font-black text-slate-900 dark:text-white tracking-tight leading-tight mb-1.5">
            Sistem ERP Terpadu untuk Rantai Kafe & Resto Modern
          </h1>
          <p class="text-xs text-slate-500 dark:text-slate-400 leading-relaxed">
            Gateway otorisasi presisi tinggi untuk pembukaan shift kasir, kalibrasi terminal kopi, dan rekonsiliasi arus kas harian secara real-time.
          </p>
        </div>

        <!-- Store / Outlet Selector Widget -->
        <div class="bg-blue-50/70 dark:bg-blue-950/30 border border-blue-100 dark:border-blue-900/50 rounded-2xl p-2.5 sm:p-3 flex items-center justify-between my-1.5">
          <div class="flex items-center gap-2.5 min-w-0">
            <div class="w-8 h-8 rounded-xl bg-blue-600 text-white flex items-center justify-center shrink-0 shadow-xs">
              <Store class="w-4 h-4" />
            </div>
            <div class="min-w-0">
              <div class="text-xs font-bold text-slate-900 dark:text-white truncate">
                {{ activeStore.name }}
              </div>
              <div class="text-[10px] text-slate-500 dark:text-slate-400 truncate">
                Online • Latensi {{ activeStore.latency }} • Database Cloud {{ activeStore.dbNode }}
              </div>
            </div>
          </div>
          <button 
            type="button"
            @click="showStoreModal = true" 
            class="text-xs font-semibold text-blue-600 dark:text-blue-400 bg-white dark:bg-slate-800 px-2.5 py-1 rounded-xl border border-blue-200 dark:border-blue-800 shadow-2xs hover:bg-blue-50 shrink-0 flex items-center gap-1 cursor-pointer transition-all"
          >
            <span>Ganti Gerai</span>
            <ChevronDown class="w-3.5 h-3.5" />
          </button>
        </div>

        <!-- Operational Status: IoT Terminal & Active Shift Cards -->
        <div class="grid grid-cols-2 gap-2.5 my-2">
          <!-- Subcard 1: IoT Terminal Connection -->
          <div class="bg-[#0B1528] text-white p-3 rounded-2xl border border-slate-800 flex flex-col justify-between relative overflow-hidden">
            <div class="flex items-center justify-between text-[10px] text-slate-400 mb-1.5">
              <span class="bg-slate-800/80 px-1.5 py-0.5 rounded text-[9px] font-semibold text-slate-300">Terminal IoT</span>
              <Radio class="w-3.5 h-3.5 text-blue-400 animate-pulse" />
            </div>
            <div class="text-center my-1">
              <div class="w-9 h-9 mx-auto rounded-xl bg-slate-800/80 border border-slate-700/60 flex items-center justify-center text-blue-400 shadow-inner mb-1.5">
                <Monitor class="w-5 h-5" />
              </div>
              <div class="text-xs font-bold text-white">Station POS–01 Terhubung</div>
              <div class="text-[10px] text-slate-300">Sinkronisasi Siap Pakai</div>
            </div>
            <div class="mt-1.5 text-[9px] text-slate-400 bg-slate-900/90 py-1 px-1.5 rounded-lg text-center truncate">
              Status KDS: Aktif • Printer: 80mm Siap
            </div>
          </div>

          <!-- Subcard 2: Active Shift Status & Coffee Metrics -->
          <div class="bg-slate-50 dark:bg-slate-800/40 border border-slate-200/70 dark:border-slate-800 p-3 rounded-2xl flex flex-col justify-between">
            <div class="flex items-center justify-between mb-1">
              <span class="text-[10px] font-black uppercase text-blue-600 dark:text-blue-400 tracking-wider">SHIFT BERJALAN</span>
              <span class="bg-blue-100 dark:bg-blue-900/50 text-blue-700 dark:text-blue-300 text-[9px] font-extrabold px-1.5 py-0.5 rounded-full">
                Shift Pagi
              </span>
            </div>
            <div class="grid grid-cols-2 gap-1 text-[10px] my-1">
              <div>
                <div class="text-slate-400 text-[9px]">Jadwal Shift</div>
                <div class="font-bold text-slate-800 dark:text-slate-200">06:30 – 15:00 WIB</div>
              </div>
              <div>
                <div class="text-slate-400 text-[9px]">Kas Modal Awal</div>
                <div class="font-bold text-slate-800 dark:text-slate-200">Rp 500.000</div>
              </div>
            </div>
            <div class="text-[10px] text-slate-600 dark:text-slate-300 flex items-center gap-1 truncate mb-1">
              <User class="w-3.5 h-3.5 text-blue-500 shrink-0" />
              <span class="truncate">Barista Bertugas: <strong class="text-slate-900 dark:text-white">Nadia & Reza (#A1)</strong></span>
            </div>
            <div class="grid grid-cols-2 gap-1.5 pt-1 border-t border-slate-200/60 dark:border-slate-700">
              <div class="bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-lg p-1 flex items-center gap-1 text-[9px] text-slate-700 dark:text-slate-300">
                <Thermometer class="w-3 h-3 text-blue-500 shrink-0" />
                <span class="truncate font-medium">92.5°C Stabil</span>
              </div>
              <div class="bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-lg p-1 flex items-center gap-1 text-[9px] text-slate-700 dark:text-slate-300">
                <Sliders class="w-3 h-3 text-blue-500 shrink-0" />
                <span class="truncate font-medium">Dial #2.4 Fine</span>
              </div>
            </div>
          </div>
        </div>

        <!-- Bottom Real-time Analytics Bar -->
        <div class="bg-[#0B1323] text-slate-200 px-3.5 py-2 rounded-2xl flex items-center justify-between text-[11px] font-medium border border-slate-800 mt-2">
          <div class="flex items-center gap-1.5">
            <Receipt class="w-3.5 h-3.5 text-blue-400 shrink-0" />
            <span class="truncate">542 Transaksi Terverifikasi Hari Ini</span>
          </div>
          <div class="flex items-center gap-1.5">
            <Armchair class="w-3.5 h-3.5 text-amber-400 shrink-0" />
            <span class="truncate">18 / 25 Meja (72%)</span>
          </div>
          <div class="flex items-center gap-1.5">
            <Droplets class="w-3.5 h-3.5 text-cyan-400 shrink-0" />
            <span class="truncate">Fresh Milk: 48L</span>
          </div>
        </div>
      </div>

      <!-- ========================================================
           RIGHT PANEL: Cashier PIN Keypad & Manager Authentication
           ======================================================== -->
      <div 
        class="lg:col-span-5 bg-white dark:bg-slate-900 rounded-3xl border border-slate-200/80 dark:border-slate-800 shadow-xl p-5 sm:p-6 flex flex-col justify-between relative overflow-hidden"
      >
        <!-- Blue Top Accent Line (matching screenshot) -->
        <div class="absolute top-0 left-6 right-6 h-1 bg-blue-600 rounded-b-full"></div>

        <!-- Role Tab Switcher: Barista/Kasir vs Manager/Spv -->
        <div class="bg-slate-100 dark:bg-slate-800/80 p-1 rounded-2xl flex items-center gap-1 mb-2">
          <button 
            type="button"
            @click="switchTab('cashier')" 
            class="flex-1 py-1.5 rounded-xl text-xs font-bold flex items-center justify-center gap-1.5 transition-all cursor-pointer"
            :class="authTab === 'cashier' ? 'bg-white dark:bg-slate-900 text-slate-900 dark:text-white shadow-2xs' : 'text-slate-500 hover:text-slate-900 dark:hover:text-white'"
          >
            <Monitor class="w-3.5 h-3.5" />
            <span>Barista / Kasir</span>
          </button>
          <button 
            type="button"
            @click="switchTab('manager')" 
            class="flex-1 py-1.5 rounded-xl text-xs font-bold flex items-center justify-center gap-1.5 transition-all cursor-pointer"
            :class="authTab === 'manager' ? 'bg-white dark:bg-slate-900 text-slate-900 dark:text-white shadow-2xs' : 'text-slate-500 hover:text-slate-900 dark:hover:text-white'"
          >
            <ShieldCheck class="w-3.5 h-3.5" />
            <span>Manager / Spv</span>
          </button>
        </div>

        <!-- ==========================================
             TAB 1: Cashier 4-Digit PIN & Numeric Keypad
             ========================================== -->
        <div v-if="authTab === 'cashier'" class="flex-1 flex flex-col justify-between">
          <div class="text-center my-1">
            <h2 class="text-base font-bold text-slate-900 dark:text-white">Otorisasi Shift Kasir</h2>
            <p class="text-xs text-slate-500 dark:text-slate-400 mt-0.5">
              Masukkan 4-digit PIN pegawai untuk mengunci sesi register
            </p>
          </div>

          <!-- 4-Digit PIN Dots Indicator -->
          <div class="flex items-center justify-center gap-3 my-2">
            <div 
              v-for="i in 4" 
              :key="i" 
              class="w-3 h-3 rounded-full transition-all duration-200"
              :class="pin.length >= i ? 'bg-blue-600 ring-4 ring-blue-100 dark:ring-blue-900/50 scale-110' : 'bg-slate-200 dark:bg-slate-700'"
            ></div>
          </div>

          <!-- Alert message if lockout or error -->
          <div 
            v-if="isLockedOut" 
            class="p-2 my-1 rounded-xl border flex items-center gap-2 bg-amber-50 dark:bg-amber-950/40 border-amber-200 dark:border-amber-900 text-amber-800 dark:text-amber-300 text-xs font-medium"
          >
            <AlertTriangle class="w-3.5 h-3.5 text-amber-600 dark:text-amber-400 shrink-0" />
            <span class="text-[10px]">Terkunci sementara. Tunggu <strong>{{ lockoutSeconds }}s</strong>.</span>
          </div>
          <div 
            v-else-if="errorMessage" 
            class="p-1.5 my-1 rounded-xl border flex items-center gap-2 bg-rose-50 dark:bg-rose-950/40 border-rose-200 dark:border-rose-900 text-rose-700 dark:text-rose-300 text-xs font-medium shake-error"
          >
            <AlertCircle class="w-3.5 h-3.5 text-rose-500 shrink-0" />
            <span class="text-[10px] leading-tight flex-1">{{ errorMessage }}</span>
          </div>

          <!-- 3x4 Numeric Keypad with Soft Rounded Buttons -->
          <div class="grid grid-cols-3 gap-2 my-1.5">
            <button 
              v-for="digit in ['1', '2', '3', '4', '5', '6', '7', '8', '9']" 
              :key="digit" 
              type="button" 
              @click="appendDigit(digit)" 
              class="h-11 rounded-2xl bg-slate-50 dark:bg-slate-800 hover:bg-slate-100 dark:hover:bg-slate-700/80 text-lg font-bold text-slate-800 dark:text-slate-100 flex items-center justify-center cursor-pointer transition active:scale-95 border border-slate-100 dark:border-slate-700/60 shadow-2xs"
            >
              {{ digit }}
            </button>
            
            <!-- Biometric Fingerprint / RFID Key -->
            <button 
              type="button" 
              @click="handleFingerprint" 
              title="Simulasi Otorisasi Sidik Jari Kasir (1-Klik Cepat)" 
              class="h-11 rounded-2xl bg-blue-50 dark:bg-blue-950/40 hover:bg-blue-100 dark:hover:bg-blue-900/50 text-blue-600 dark:text-blue-400 flex items-center justify-center cursor-pointer transition active:scale-95 border border-blue-200 dark:border-blue-900 shadow-2xs"
            >
              <Fingerprint class="w-5 h-5" />
            </button>

            <!-- '0' Key -->
            <button 
              type="button" 
              @click="appendDigit('0')" 
              class="h-11 rounded-2xl bg-slate-50 dark:bg-slate-800 hover:bg-slate-100 dark:hover:bg-slate-700/80 text-lg font-bold text-slate-800 dark:text-slate-100 flex items-center justify-center cursor-pointer transition active:scale-95 border border-slate-100 dark:border-slate-700/60 shadow-2xs"
            >
              0
            </button>

            <!-- Backspace Key -->
            <button 
              type="button" 
              @click="deleteDigit" 
              title="Hapus Digit" 
              class="h-11 rounded-2xl bg-slate-50 dark:bg-slate-800 hover:bg-slate-100 dark:hover:bg-slate-700/80 text-slate-600 dark:text-slate-300 flex items-center justify-center cursor-pointer transition active:scale-95 border border-slate-100 dark:border-slate-700/60 shadow-2xs"
            >
              <Delete class="w-5 h-5" />
            </button>
          </div>

          <!-- Shift Open Action Button -->
          <button 
            type="button" 
            @click="submitPinLogin" 
            :disabled="loading || isLockedOut" 
            class="btn-press w-full bg-[#0B1323] hover:bg-slate-800 text-white font-bold py-2.5 px-4 rounded-2xl text-xs flex items-center justify-center gap-2 shadow-md transition cursor-pointer disabled:opacity-50 mt-1"
          >
            <span v-if="loading" class="animate-spin w-3.5 h-3.5 border-2 border-white border-t-transparent rounded-full"></span>
            <Unlock v-else class="w-3.5 h-3.5" />
            <span>{{ loading ? 'Mengotentikasi...' : isLockedOut ? `Terkunci (${lockoutSeconds}s)` : 'Buka Shift & Masuk POS Kasir' }}</span>
          </button>
        </div>

        <!-- ==========================================
             TAB 2: Manager & Back-Office Authentication
             ========================================== -->
        <div v-else class="flex-1 flex flex-col justify-between py-1">
          <div class="text-center mb-2">
            <h2 class="text-base font-bold text-slate-900 dark:text-white">Otorisasi Manager & Executive</h2>
            <p class="text-xs text-slate-500 dark:text-slate-400 mt-0.5">
              Akses back-office ERP, manajemen stok, keuangan & HRIS
            </p>
          </div>

          <!-- Quick 1-Click Role Fill Chips -->
          <div class="mb-2">
            <div class="text-[10px] font-bold uppercase tracking-wider mb-1.5 flex items-center justify-between text-slate-500">
              <span>Pilih Peran Demo</span>
              <span class="text-blue-600 dark:text-blue-400 font-semibold lowercase">1-klik isi</span>
            </div>
            <div class="grid grid-cols-3 gap-1">
              <button
                v-for="role in demoRoles"
                :key="role.email"
                type="button"
                @click="selectDemoAccount(role.email, role.roleName)"
                class="py-1 px-1.5 rounded-lg border flex items-center justify-center gap-1 text-[11px] font-semibold transition-all cursor-pointer hover:border-blue-500/50"
                :class="email === role.email ? 'bg-blue-50 dark:bg-blue-900/40 border-blue-400 text-blue-600 dark:text-blue-300 font-bold' : 'text-slate-600 dark:text-slate-400 bg-slate-50 dark:bg-slate-800 border-slate-200 dark:border-slate-700'"
              >
                <span>{{ role.icon }}</span>
                <span class="truncate">{{ role.shortName }}</span>
              </button>
            </div>
          </div>

          <!-- Alert message if any -->
          <div 
            v-if="errorMessage" 
            class="p-2 mb-2 rounded-xl border flex items-center gap-2 bg-rose-50 dark:bg-rose-950/40 border-rose-200 dark:border-rose-900 text-rose-700 dark:text-rose-300 text-xs font-medium shake-error"
          >
            <AlertCircle class="w-3.5 h-3.5 text-rose-500 shrink-0" />
            <span class="text-[10px] leading-tight flex-1">{{ errorMessage }}</span>
          </div>

          <!-- Email & Password Form -->
          <form @submit.prevent="handleManagerLogin" class="space-y-2" novalidate>
            <div>
              <label for="mgr-email" class="block text-[11px] font-bold mb-0.5 text-slate-700 dark:text-slate-300">
                Email Pengguna <span class="text-rose-500">*</span>
              </label>
              <div class="rounded-xl px-3 py-1.5 flex items-center gap-2 border bg-slate-50 dark:bg-slate-800 border-slate-200 dark:border-slate-700 focus-within:border-blue-500">
                <Mail class="w-3.5 h-3.5 text-slate-400 shrink-0" />
                <input
                  id="mgr-email"
                  v-model="email"
                  type="email"
                  autocomplete="email"
                  required
                  placeholder="manager@cafe-erp.com"
                  class="w-full bg-transparent p-0 text-xs font-semibold outline-none border-none ring-0 text-slate-900 dark:text-white"
                />
              </div>
            </div>

            <div>
              <label for="mgr-pass" class="block text-[11px] font-bold mb-0.5 text-slate-700 dark:text-slate-300">
                Kata Sandi <span class="text-rose-500">*</span>
              </label>
              <div class="rounded-xl px-3 py-1.5 flex items-center gap-2 border bg-slate-50 dark:bg-slate-800 border-slate-200 dark:border-slate-700 focus-within:border-blue-500">
                <Lock class="w-3.5 h-3.5 text-slate-400 shrink-0" />
                <input
                  id="mgr-pass"
                  v-model="password"
                  :type="showPassword ? 'text' : 'password'"
                  autocomplete="current-password"
                  required
                  placeholder="••••••••••••"
                  class="w-full bg-transparent p-0 text-xs font-semibold outline-none border-none ring-0 text-slate-900 dark:text-white"
                />
                <button
                  type="button"
                  @click="showPassword = !showPassword"
                  class="text-slate-400 hover:text-slate-600 dark:hover:text-slate-200 cursor-pointer"
                >
                  <EyeOff v-if="showPassword" class="w-3.5 h-3.5" />
                  <Eye v-else class="w-3.5 h-3.5" />
                </button>
              </div>
            </div>

            <div class="flex items-center justify-between text-xs pt-0.5">
              <label class="flex items-center gap-1.5 cursor-pointer text-slate-500">
                <input v-model="rememberMe" type="checkbox" class="w-3.5 h-3.5 rounded text-blue-600" />
                <span class="text-[11px]">Ingat Akun</span>
              </label>
              <button 
                type="button" 
                @click="showForgotModal = true" 
                class="text-blue-600 dark:text-blue-400 font-bold hover:underline text-[11px] cursor-pointer"
              >
                Lupa Sandi?
              </button>
            </div>

            <button 
              type="submit" 
              :disabled="loading" 
              class="btn-press w-full bg-blue-600 hover:bg-blue-700 text-white font-bold py-2 px-3 rounded-xl text-xs flex items-center justify-center gap-2 shadow-md transition cursor-pointer disabled:opacity-50 mt-1"
            >
              <span v-if="loading" class="animate-spin w-3.5 h-3.5 border-2 border-white border-t-transparent rounded-full"></span>
              <LogIn v-else class="w-3.5 h-3.5" />
              <span>Masuk ke Sistem ERP</span>
            </button>
          </form>
        </div>

        <!-- Right Panel Footer Links matching design -->
        <div class="pt-2.5 border-t border-slate-200/70 dark:border-slate-800 mt-2 space-y-1.5 text-xs">
          <div class="flex items-center justify-between">
            <button 
              type="button" 
              @click="showOfflineModal = true" 
              class="text-blue-600 dark:text-blue-400 hover:underline flex items-center gap-1 font-semibold text-[11px] cursor-pointer"
            >
              <WifiOff class="w-3.5 h-3.5" />
              <span>Mode Darurat Offline (SQLite)</span>
            </button>
            <button 
              type="button" 
              @click="showForgotModal = true" 
              class="text-slate-500 dark:text-slate-400 hover:underline flex items-center gap-1 text-[11px] cursor-pointer"
            >
              <Headphones class="w-3.5 h-3.5" />
              <span>Bantuan IT</span>
            </button>
          </div>
          <div class="flex items-center justify-between text-[11px] text-slate-400 dark:text-slate-500">
            <div class="flex items-center gap-2">
              <button type="button" @click="toggleLang" class="font-bold hover:text-slate-600 dark:hover:text-slate-300 cursor-pointer">
                {{ lang === 'id' ? 'ID / EN' : 'EN / ID' }}
              </button>
              <span>•</span>
              <button type="button" @click="toggleTheme" class="hover:text-slate-600 dark:hover:text-slate-300 cursor-pointer flex items-center gap-1" title="Ganti Tema">
                <Sun v-if="appStore.theme === 'dark'" class="w-3 h-3 text-amber-500" />
                <Moon v-else class="w-3 h-3 text-blue-500" />
                <span>{{ appStore.theme === 'dark' ? 'Light' : 'Dark' }}</span>
              </button>
            </div>
            <span>KopiERP Web Enterprise v2.4.8</span>
          </div>
        </div>

      </div>
    </div>

    <!-- Modal: Ganti Gerai / Switch Outlet -->
    <Teleport to="body">
      <transition name="fade">
        <div 
          v-if="showStoreModal" 
          class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-950/60 backdrop-blur-xs"
          @click.self="showStoreModal = false"
        >
          <div 
            class="w-full max-w-sm rounded-2xl p-5 shadow-2xl border bg-white dark:bg-slate-900 border-slate-200 dark:border-slate-800 animate-in zoom-in-95 duration-200"
            role="dialog"
            aria-modal="true"
          >
            <div class="flex items-center justify-between pb-2.5 border-b border-slate-200 dark:border-slate-800">
              <div class="flex items-center gap-2">
                <Store class="w-4 h-4 text-blue-600" />
                <h3 class="font-bold text-sm text-slate-900 dark:text-white">Pilih Gerai Operasional</h3>
              </div>
              <button @click="showStoreModal = false" class="p-1 rounded-lg text-slate-400 hover:text-slate-700 cursor-pointer text-xs">
                ✕
              </button>
            </div>
            <div class="py-3 space-y-2">
              <button
                v-for="store in availableStores"
                :key="store.name"
                type="button"
                @click="selectStore(store)"
                class="w-full p-2.5 rounded-xl border text-left flex items-center justify-between transition-all cursor-pointer"
                :class="activeStore.name === store.name ? 'border-blue-500 bg-blue-50 dark:bg-blue-900/30' : 'border-slate-200 dark:border-slate-800 hover:bg-slate-50 dark:hover:bg-slate-800/60'"
              >
                <div>
                  <div class="text-xs font-bold text-slate-900 dark:text-white">{{ store.name }}</div>
                  <div class="text-[10px] text-slate-500">Latensi {{ store.latency }} • {{ store.dbNode }}</div>
                </div>
                <span v-if="activeStore.name === store.name" class="text-xs font-bold text-blue-600">Aktif</span>
              </button>
            </div>
          </div>
        </div>
      </transition>
    </Teleport>

    <!-- Modal: Offline SQLite Mode Info -->
    <Teleport to="body">
      <transition name="fade">
        <div 
          v-if="showOfflineModal" 
          class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-950/60 backdrop-blur-xs"
          @click.self="showOfflineModal = false"
        >
          <div 
            class="w-full max-w-sm rounded-2xl p-5 shadow-2xl border bg-white dark:bg-slate-900 border-slate-200 dark:border-slate-800 animate-in zoom-in-95 duration-200"
            role="dialog"
            aria-modal="true"
          >
            <div class="flex items-center justify-between pb-2.5 border-b border-slate-200 dark:border-slate-800">
              <div class="flex items-center gap-2">
                <WifiOff class="w-4 h-4 text-blue-600" />
                <h3 class="font-bold text-sm text-slate-900 dark:text-white">Sinkronisasi Offline (SQLite)</h3>
              </div>
              <button @click="showOfflineModal = false" class="p-1 rounded-lg text-slate-400 hover:text-slate-700 cursor-pointer text-xs">
                ✕
              </button>
            </div>
            <div class="py-3 text-xs text-slate-600 dark:text-slate-300 leading-relaxed space-y-2">
              <p>
                Terminal POS dilengkapi dengan replikasi lokal <strong>SQLite Cache</strong>. Jika koneksi cloud terputus:
              </p>
              <ul class="list-disc list-inside space-y-1 text-[11px] text-slate-500">
                <li>Transaksi kasir tetap dapat berjalan tanpa jeda.</li>
                <li>Receipt dapat dicetak normal via thermal printer lokal.</li>
                <li>Data transaksi otomatis tersinkronisasi kembali saat online.</li>
              </ul>
              <div class="p-2 rounded-xl bg-emerald-50 dark:bg-emerald-950/40 text-emerald-800 dark:text-emerald-300 text-[11px] font-semibold border border-emerald-200 dark:border-emerald-800">
                ✓ Replikasi lokal terkini: 2 menit yang lalu (Siap digunakan)
              </div>
            </div>
            <div class="pt-2 border-t border-slate-200 dark:border-slate-800 flex justify-end">
              <button 
                type="button" 
                @click="showOfflineModal = false" 
                class="btn-press px-3.5 py-1.5 rounded-xl text-xs font-bold text-white bg-blue-600 hover:bg-blue-700 cursor-pointer"
              >
                Tutup
              </button>
            </div>
          </div>
        </div>
      </transition>
    </Teleport>

    <!-- Modal: Bantuan IT -->
    <Teleport to="body">
      <transition name="fade">
        <div 
          v-if="showForgotModal" 
          class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-950/60 backdrop-blur-xs"
          @click.self="showForgotModal = false"
        >
          <div 
            class="w-full max-w-sm rounded-2xl p-5 shadow-2xl border bg-white dark:bg-slate-900 border-slate-200 dark:border-slate-800 animate-in zoom-in-95 duration-200"
            role="dialog"
            aria-modal="true"
          >
            <div class="flex items-center justify-between pb-2.5 border-b border-slate-200 dark:border-slate-800">
              <div class="flex items-center gap-2">
                <Headphones class="w-4 h-4 text-blue-600" />
                <h3 class="font-bold text-sm text-slate-900 dark:text-white">Bantuan IT & Otorisasi</h3>
              </div>
              <button @click="showForgotModal = false" class="p-1 rounded-lg text-slate-400 hover:text-slate-700 cursor-pointer text-xs">
                ✕
              </button>
            </div>
            <div class="py-3 text-xs text-slate-600 dark:text-slate-300 leading-relaxed space-y-2">
              <p>
                PIN atau kata sandi dapat di-reset melalui Super Administrator outlet Anda.
              </p>
              <div class="p-2.5 rounded-xl border bg-slate-50 dark:bg-slate-800 border-slate-200 dark:border-slate-700 text-[11px]">
                <div class="font-bold text-slate-900 dark:text-white mb-0.5">Hotline Dukungan POS Resto:</div>
                <div class="text-blue-600 dark:text-blue-400 font-bold">+62 812-9988-7766 (Senopati HQ)</div>
                <div class="text-slate-400 mt-1">Layanan teknis 24/7 untuk kasir & barista outlet.</div>
              </div>
            </div>
            <div class="pt-2 border-t border-slate-200 dark:border-slate-800 flex justify-end">
              <button 
                type="button" 
                @click="showForgotModal = false" 
                class="btn-press px-3.5 py-1.5 rounded-xl text-xs font-bold text-white bg-blue-600 hover:bg-blue-700 cursor-pointer"
              >
                Tutup
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
  Store,
  Monitor,
  Radio,
  Receipt,
  Armchair,
  Droplets,
  Thermometer,
  Sliders,
  User,
  Fingerprint,
  Delete,
  Lock,
  Unlock,
  ShieldCheck,
  Headphones,
  WifiOff,
  ChevronDown,
  Mail,
  Eye,
  EyeOff,
  Sun,
  Moon,
  LogIn,
  AlertCircle,
  AlertTriangle
} from 'lucide-vue-next'
import { useAuthStore } from '@/stores/auth.store'
import { useAppStore } from '@/stores/app.store'
import { useNotificationStore } from '@/stores/notification.store'
import axios from 'axios'

const router = useRouter()
const authStore = useAuthStore()
const appStore = useAppStore()
const notifyStore = useNotificationStore()

// Tabs: 'cashier' (PIN keypad) | 'manager' (email/password)
const authTab = ref<'cashier' | 'manager'>('cashier')

// Cashier PIN state
const pin = ref('')

// Manager Form state
const email = ref('admin@cafe-erp.com')
const password = ref('Admin@123')
const rememberMe = ref(true)
const showPassword = ref(false)

// Shared auth state
const loading = ref(false)
const errorMessage = ref('')
const lang = ref<'id' | 'en'>('id')

// Modals
const showStoreModal = ref(false)
const showOfflineModal = ref(false)
const showForgotModal = ref(false)

// Rate limiting state
const failedAttempts = ref(0)
const isLockedOut = ref(false)
const lockoutSeconds = ref(0)
let lockoutTimer: any = null

const REMEMBER_KEY = 'cafe_erp_remembered_email'

// Available store outlets
const availableStores = [
  { name: 'Sudirman Flagship Store #01', latency: '24ms', dbNode: 'ID-JKT-02' },
  { name: 'Senopati Branch #02', latency: '28ms', dbNode: 'ID-JKT-01' },
  { name: 'Grand Indonesia Central #03', latency: '19ms', dbNode: 'ID-JKT-03' },
  { name: 'Dharmawangsa Roastery #04', latency: '32ms', dbNode: 'ID-JKT-02' }
]
const activeStore = ref(availableStores[0])

const selectStore = (store: typeof availableStores[0]) => {
  activeStore.value = store
  showStoreModal.value = false
  notifyStore.info(`Gerai aktif dialihkan ke ${store.name}`, 'Gerai Terpilih')
}

// Demo quick roles
const demoRoles = [
  { icon: '👑', shortName: 'Admin', roleName: 'Super Admin', email: 'admin@cafe-erp.com' },
  { icon: '🏷️', shortName: 'Kasir', roleName: 'Kasir POS', email: 'kasir@cafe-erp.com' },
  { icon: '☕', shortName: 'Barista', roleName: 'Barista', email: 'barista@cafe-erp.com' },
  { icon: '📦', shortName: 'Gudang', roleName: 'Inventory', email: 'inventory@cafe-erp.com' },
  { icon: '💰', shortName: 'Finance', roleName: 'Finance', email: 'finance@cafe-erp.com' },
  { icon: '👥', shortName: 'HRIS', roleName: 'HR Admin', email: 'hris@cafe-erp.com' }
]

// Keyboard listener for keypad
const handleKeyDown = (e: KeyboardEvent) => {
  if (authTab.value !== 'cashier') return
  if (e.key >= '0' && e.key <= '9') {
    appendDigit(e.key)
  } else if (e.key === 'Backspace') {
    deleteDigit()
  } else if (e.key === 'Enter') {
    submitPinLogin()
  }
}

onMounted(() => {
  // Lock document to eliminate any scrollbars
  document.documentElement.classList.add('login-page')
  document.body.classList.add('login-page')
  document.body.style.overflow = 'hidden'

  window.addEventListener('keydown', handleKeyDown)

  const savedEmail = localStorage.getItem(REMEMBER_KEY)
  if (savedEmail) {
    email.value = savedEmail
    rememberMe.value = true
  }
})

onUnmounted(() => {
  document.documentElement.classList.remove('login-page')
  document.body.classList.remove('login-page')
  document.body.style.overflow = ''

  window.removeEventListener('keydown', handleKeyDown)
  if (lockoutTimer) clearInterval(lockoutTimer)
})

const switchTab = (tab: 'cashier' | 'manager') => {
  authTab.value = tab
  errorMessage.value = ''
  pin.value = ''
}

const appendDigit = (digit: string) => {
  if (isLockedOut.value) return
  errorMessage.value = ''
  if (pin.value.length < 4) {
    pin.value += digit
    if (pin.value.length === 4) {
      // Auto-submit on 4th digit for instant cashier ergonomics
      setTimeout(() => {
        submitPinLogin()
      }, 150)
    }
  }
}

const deleteDigit = () => {
  if (pin.value.length > 0) {
    pin.value = pin.value.slice(0, -1)
    errorMessage.value = ''
  }
}

const handleFingerprint = () => {
  if (isLockedOut.value) return
  pin.value = '1234'
  notifyStore.success('Sidik Jari Kasir Terverifikasi: Jane Doe (#A1)', 'Biometrik Terhubung')
  setTimeout(() => {
    submitPinLogin()
  }, 250)
}

const toggleLang = () => {
  lang.value = lang.value === 'id' ? 'en' : 'id'
  notifyStore.info(`Bahasa diubah ke ${lang.value.toUpperCase()}`, 'Pengaturan Bahasa')
}

const toggleTheme = () => {
  const nextTheme = appStore.theme === 'dark' ? 'light' : 'dark'
  appStore.setTheme(nextTheme)
}

const selectDemoAccount = (selectedEmail: string, roleName: string) => {
  email.value = selectedEmail
  password.value = 'Admin@123'
  errorMessage.value = ''
  notifyStore.info(`Akun terpasang: ${roleName}`, 'Demo Account')
}

// Execute authentication against backend API
const authenticateWithCredentials = async (loginEmail: string, loginPass: string, targetRoute: string) => {
  loading.value = true
  errorMessage.value = ''

  try {
    const res = await axios.post('/api/v1/auth/login', {
      email: loginEmail.trim(),
      password: loginPass
    })

    if (res.data?.data) {
      const data = res.data.data
      const userObj = {
        id: data.user.id,
        email: data.user.email,
        firstName: data.user.full_name || 'Pengguna',
        lastName: '',
        roleId: data.role || data.user.role_name || 'Kasir POS'
      }

      authStore.login(
        data.token,
        data.refresh_token || data.token,
        userObj,
        data.permissions || []
      )

      if (rememberMe.value) {
        localStorage.setItem(REMEMBER_KEY, loginEmail.trim())
      } else {
        localStorage.removeItem(REMEMBER_KEY)
      }

      failedAttempts.value = 0
      notifyStore.success(`Shift Dibuka: ${userObj.firstName}`, `Station POS-01 Siap Digunakan`)
      await router.push(targetRoute)
    } else {
      throw new Error('Respons otentikasi tidak valid dari server')
    }
  } catch (err: any) {
    failedAttempts.value += 1
    if (failedAttempts.value >= 5) {
      triggerLockout(30)
      return
    }
    const genericMsg = 'PIN atau kredensial otorisasi tidak valid. Periksa kembali akses Anda.'
    errorMessage.value = genericMsg
    notifyStore.error(genericMsg, 'Otorisasi Gagal')
  } finally {
    loading.value = false
  }
}

// Cashier PIN submission logic
const submitPinLogin = async () => {
  if (isLockedOut.value) return
  if (pin.value.length < 4) {
    errorMessage.value = 'Masukkan lengkap 4-digit PIN pegawai.'
    return
  }

  // Map PIN to role credentials
  // 1234: Kasir | 1111: Barista | 1122: Manager | 0000: Admin
  if (pin.value === '1234') {
    await authenticateWithCredentials('kasir@cafe-erp.com', 'Admin@123', '/pos')
  } else if (pin.value === '1111') {
    await authenticateWithCredentials('barista@cafe-erp.com', 'Admin@123', '/pos')
  } else if (pin.value === '1122') {
    await authenticateWithCredentials('manager@cafe-erp.com', 'Admin@123', '/dashboard')
  } else if (pin.value === '0000' || pin.value === '123456') {
    await authenticateWithCredentials('admin@cafe-erp.com', 'Admin@123', '/dashboard')
  } else {
    failedAttempts.value += 1
    if (failedAttempts.value >= 5) {
      triggerLockout(30)
      return
    }
    pin.value = ''
    errorMessage.value = 'PIN kasir tidak cocok. (Gunakan 1234 untuk Kasir atau 1111 untuk Barista)'
  }
}

// Manager email/password submission logic
const handleManagerLogin = async () => {
  if (isLockedOut.value) return
  if (!email.value.trim() || !password.value) {
    errorMessage.value = 'Email dan kata sandi wajib diisi.'
    return
  }
  await authenticateWithCredentials(email.value.trim(), password.value, '/dashboard')
}

const triggerLockout = (seconds: number) => {
  isLockedOut.value = true
  lockoutSeconds.value = seconds
  errorMessage.value = `Terlalu banyak percobaan gagal. Terminal dikunci sementara selama ${seconds}s.`
  notifyStore.error(errorMessage.value, 'Keamanan Terminal')

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
