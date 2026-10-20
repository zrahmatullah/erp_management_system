<template>
  <div class="dashboard-v2 space-y-6 animate-in fade-in duration-200">
    <!-- 1. TOP WELCOME HERO & PERSPECTIVE TABS -->
    <div 
      class="p-6 sm:p-7 rounded-2xl border transition-all duration-200 relative overflow-hidden backdrop-blur-md"
      style="background-color: var(--bg-primary); border-color: var(--border-color);"
    >
      <div class="flex flex-col lg:flex-row lg:items-center lg:justify-between gap-6">
        <!-- User Identity & Greeting -->
        <div class="flex items-center gap-4 min-w-0">
          <div 
            class="w-14 h-14 rounded-2xl flex items-center justify-center text-xl font-bold text-white shadow-xs shrink-0"
            :class="isOperationalStaff ? 'bg-gradient-to-tr from-amber-500 to-orange-500' :
                    isFinanceRole ? 'bg-gradient-to-tr from-purple-600 to-indigo-600' :
                    isHRISRole ? 'bg-gradient-to-tr from-rose-500 to-pink-600' :
                    'bg-gradient-to-tr from-blue-600 to-indigo-600'"
          >
            {{ userInitials }}
          </div>

          <div class="min-w-0">
            <div class="flex items-center gap-2 flex-wrap mb-1">
              <span 
                class="text-[11px] font-semibold px-2.5 py-0.5 rounded-full border capitalize"
                :class="isOperationalStaff ? 'bg-amber-50 text-amber-700 border-amber-200/80 dark:bg-amber-950/40 dark:text-amber-300' :
                        isFinanceRole ? 'bg-purple-50 text-purple-700 border-purple-200/80 dark:bg-purple-950/40 dark:text-purple-300' :
                        isHRISRole ? 'bg-rose-50 text-rose-700 border-rose-200/80 dark:bg-rose-950/40 dark:text-rose-300' :
                        'bg-blue-50 text-blue-700 border-blue-200/80 dark:bg-blue-950/40 dark:text-blue-300'"
              >
                ● {{ currentRoleName }}
              </span>
              <span v-if="portalData.employee?.nik" class="text-[11px] font-mono px-2 py-0.5 rounded-md border text-slate-500 dark:text-slate-400"
                    style="border-color: var(--border-color); background-color: var(--bg-content);">
                NIK: {{ portalData.employee.nik }}
              </span>
              <span class="inline-flex items-center gap-1.5 text-[11px] font-semibold text-emerald-600 dark:text-emerald-400 bg-emerald-50 dark:bg-emerald-950/40 border border-emerald-200/80 dark:border-emerald-800/60 px-2 py-0.5 rounded-md">
                <span class="w-1.5 h-1.5 rounded-full bg-emerald-500 animate-pulse"></span>
                Staf Aktif
              </span>
            </div>

            <h1 class="text-xl sm:text-2xl font-bold tracking-tight text-slate-900 dark:text-slate-50 truncate">
              {{ greetingTime }}, {{ portalData.employee?.full_name || userFullName }}!
            </h1>
            <p class="text-xs sm:text-sm text-slate-500 dark:text-slate-400 mt-0.5 truncate">
              <span v-if="isOperationalStaff">Kelola presensi harian, periksa giliran shift, dan tugas operasional outlet.</span>
              <span v-else-if="isFinanceRole">Pantau arus kas, realisasi jurnal umum, dan profitabilitas cafe.</span>
              <span v-else-if="isHRISRole">Monitoring kehadiran karyawan, approval cuti, dan administrasi tim.</span>
              <span v-else>Executive overview performa POS, ketersediaan meja, dan mutasi inventaris.</span>
            </p>
          </div>
        </div>

        <!-- Right: WIB Clock & Fast Actions -->
        <div class="flex items-center gap-3 shrink-0 self-start lg:self-auto flex-wrap">
          <!-- Live WIB Clock Card -->
          <div 
            class="px-3.5 py-2 rounded-xl border text-right hidden sm:block"
            style="border-color: var(--border-color); background-color: var(--bg-content);"
          >
            <div class="flex items-center gap-2 text-xs font-bold text-blue-600 dark:text-blue-400">
              <Clock class="w-3.5 h-3.5" />
              <span class="font-mono text-sm tracking-wide">{{ liveTime }}</span>
              <span class="text-[10px] text-slate-400 font-normal">WIB</span>
            </div>
            <div class="text-[11px] text-slate-500 dark:text-slate-400 mt-0.5">
              {{ liveDate }}
            </div>
          </div>

          <!-- Action Buttons -->
          <button
            v-if="isOperationalStaff || activeViewTab === 'employee'"
            type="button"
            @click="openLeaveModal"
            class="px-3.5 py-2 rounded-xl text-xs font-semibold bg-blue-600 hover:bg-blue-700 text-white transition-all shadow-xs flex items-center gap-1.5 cursor-pointer"
          >
            <Calendar class="w-3.5 h-3.5" />
            <span>Ajukan Cuti</span>
          </button>

          <button
            v-if="isOperationalStaff || activeViewTab === 'employee'"
            type="button"
            @click="openProfileModal"
            class="px-3.5 py-2 rounded-xl text-xs font-semibold border transition-all flex items-center gap-1.5 cursor-pointer text-slate-700 dark:text-slate-200 hover:bg-slate-100 dark:hover:bg-slate-800"
            style="border-color: var(--border-color); background-color: var(--bg-content);"
          >
            <User class="w-3.5 h-3.5" />
            <span>Profil Saya</span>
          </button>

          <button
            type="button"
            @click="refreshCurrentView"
            class="p-2 rounded-xl border transition-all cursor-pointer text-slate-500 hover:text-slate-900 dark:hover:text-slate-100 hover:bg-slate-100 dark:hover:bg-slate-800"
            style="border-color: var(--border-color); background-color: var(--bg-content);"
            title="Muat Ulang Data"
          >
            <RefreshCw class="w-4 h-4" :class="{ 'animate-spin': isRefreshing }" />
          </button>
        </div>
      </div>

      <!-- Segmented Perspective Tabs (For Admin, Owner, Manager) -->
      <div 
        v-if="isManagerialRole" 
        class="mt-6 pt-4 border-t flex items-center gap-2 overflow-x-auto select-none"
        style="border-color: var(--border-color);"
      >
        <button
          type="button"
          v-for="tab in perspectiveTabs"
          :key="tab.key"
          @click="activeViewTab = tab.key"
          class="px-3.5 py-1.5 rounded-xl text-xs font-semibold transition-all cursor-pointer flex items-center gap-2 shrink-0 border"
          :class="activeViewTab === tab.key 
            ? 'bg-blue-600 text-white border-blue-600 shadow-xs' 
            : 'text-slate-600 dark:text-slate-400 hover:text-slate-900 dark:hover:text-slate-200 border-transparent hover:bg-slate-100/80 dark:hover:bg-slate-800/60'"
        >
          <component :is="tab.icon" class="w-3.5 h-3.5" />
          <span>{{ tab.label }}</span>
        </button>
      </div>
    </div>

    <!-- ===================================================================== -->
    <!-- 2. OPERATIONAL STAFF PORTAL (BARISTA, CASHIER, WAREHOUSE)             -->
    <!-- ===================================================================== -->
    <div v-if="activePerspective === 'employee'" class="space-y-6">
      <!-- HERO ATTENDANCE CLOCK IN / OUT CARD -->
      <div 
        class="p-6 rounded-2xl border shadow-xs relative overflow-hidden"
        style="background-color: var(--bg-primary); border-color: var(--border-color);"
      >
        <div class="grid grid-cols-1 lg:grid-cols-12 gap-6 items-center">
          <!-- Status & Context Info (7 cols) -->
          <div class="lg:col-span-7 space-y-3">
            <div class="flex items-center gap-2">
              <span 
                class="px-2.5 py-0.5 rounded-full text-xs font-semibold border inline-flex items-center gap-1.5"
                :class="attendanceState.hasClockedIn && !attendanceState.hasClockedOut 
                  ? 'bg-emerald-50 text-emerald-700 border-emerald-200/80 dark:bg-emerald-950/40 dark:text-emerald-300' 
                  : attendanceState.hasClockedOut 
                  ? 'bg-blue-50 text-blue-700 border-blue-200/80 dark:bg-blue-950/40 dark:text-blue-300'
                  : 'bg-amber-50 text-amber-700 border-amber-200/80 dark:bg-amber-950/40 dark:text-amber-300'"
              >
                <span class="w-1.5 h-1.5 rounded-full" 
                      :class="attendanceState.hasClockedIn && !attendanceState.hasClockedOut ? 'bg-emerald-500 animate-pulse' : 'bg-slate-400'"></span>
                {{ attendanceState.hasClockedIn && !attendanceState.hasClockedOut ? 'Sedang Bertugas' :
                   attendanceState.hasClockedOut ? 'Shift Selesai' : 'Belum Absen Masuk' }}
              </span>
              <span class="text-xs text-slate-400">Jadwal: {{ todayShiftName }}</span>
            </div>

            <h2 class="text-xl sm:text-2xl font-bold tracking-tight text-slate-900 dark:text-slate-50">
              Presensi Kehadiran Kerja Hari Ini
            </h2>

            <p class="text-xs sm:text-sm text-slate-500 dark:text-slate-400 leading-relaxed max-w-xl">
              Catat waktu kedatangan saat memulai shift dan waktu kepulangan setelah menyelesaikan operasional kafe.
            </p>

            <!-- Timestamps Summary -->
            <div class="flex items-center gap-6 pt-2 text-xs">
              <div>
                <span class="text-slate-400 block text-[11px]">Waktu Masuk (Clock In)</span>
                <span class="font-mono font-bold text-sm text-slate-900 dark:text-slate-100">
                  {{ attendanceState.clockInTime || '--:--' }}
                </span>
              </div>
              <div class="h-8 w-px bg-slate-200 dark:bg-slate-800"></div>
              <div>
                <span class="text-slate-400 block text-[11px]">Waktu Pulang (Clock Out)</span>
                <span class="font-mono font-bold text-sm text-slate-900 dark:text-slate-100">
                  {{ attendanceState.clockOutTime || '--:--' }}
                </span>
              </div>
              <div class="h-8 w-px bg-slate-200 dark:bg-slate-800"></div>
              <div>
                <span class="text-slate-400 block text-[11px]">Status Durasi</span>
                <span class="font-semibold text-slate-700 dark:text-slate-300">
                  {{ attendanceState.workedHours }}
                </span>
              </div>
            </div>
          </div>

          <!-- Action Buttons & Notes (5 cols) -->
          <div class="lg:col-span-5 flex flex-col gap-3">
            <input
              v-model="attendanceNoteInput"
              type="text"
              placeholder="Catatan kehadiran (opsional, misal: lembur, tukar shift)..."
              class="w-full px-3.5 py-2 text-xs rounded-xl border outline-none transition-all focus:ring-1 focus:ring-blue-500"
              style="background-color: var(--input-bg); border-color: var(--input-border); color: var(--input-text);"
            />

            <div class="grid grid-cols-2 gap-3">
              <!-- Clock In Button -->
              <button
                type="button"
                @click="handleClockIn"
                :disabled="attendanceState.hasClockedIn || attendanceActionLoading"
                class="py-3 px-4 rounded-xl text-xs font-bold transition-all flex items-center justify-center gap-2 cursor-pointer shadow-xs disabled:opacity-50 disabled:cursor-not-allowed"
                :class="!attendanceState.hasClockedIn 
                  ? 'bg-emerald-600 hover:bg-emerald-700 text-white' 
                  : 'bg-slate-100 dark:bg-slate-800 text-slate-400 border border-slate-200 dark:border-slate-700'"
              >
                <CheckCircle2 class="w-4 h-4" />
                <span>{{ attendanceState.hasClockedIn ? 'Sudah Masuk' : 'Clock In Masuk' }}</span>
              </button>

              <!-- Clock Out Button -->
              <button
                type="button"
                @click="handleClockOut"
                :disabled="!attendanceState.hasClockedIn || attendanceState.hasClockedOut || attendanceActionLoading"
                class="py-3 px-4 rounded-xl text-xs font-bold transition-all flex items-center justify-center gap-2 cursor-pointer shadow-xs disabled:opacity-50 disabled:cursor-not-allowed"
                :class="attendanceState.hasClockedIn && !attendanceState.hasClockedOut 
                  ? 'bg-rose-600 hover:bg-rose-700 text-white' 
                  : 'bg-slate-100 dark:bg-slate-800 text-slate-400 border border-slate-200 dark:border-slate-700'"
              >
                <Clock class="w-4 h-4" />
                <span>{{ attendanceState.hasClockedOut ? 'Selesai Shift' : 'Clock Out Pulang' }}</span>
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- 3 STAT CARDS: CUTI, GAJI, PROFIL -->
      <div class="grid grid-cols-1 sm:grid-cols-3 gap-5">
        <!-- Card 1: Kuota Cuti -->
        <StatCardV2
          title="Sisa Kuota Cuti"
          :value="`${portalData.leave_summary?.remaining_days ?? 12} Hari`"
          subtitle="Dari total jatah tahunan 12 hari"
          color="blue"
          :icon="CalendarClock"
        />

        <!-- Card 2: Slip Gaji -->
        <div 
          @click="openSalarySlipModal"
          class="p-5 rounded-2xl border transition-all duration-200 hover:shadow-xs group cursor-pointer relative"
          style="background-color: var(--bg-primary); border-color: var(--border-color);"
        >
          <div class="flex items-start justify-between">
            <div>
              <p class="text-xs font-semibold text-slate-500 dark:text-slate-400 tracking-wide uppercase">
                Slip Gaji Terakhir
              </p>
              <h3 class="text-2xl font-bold tracking-tight text-slate-900 dark:text-slate-50 mt-2">
                Rp {{ formatNum(portalData.latest_payroll?.net_salary || 4850000) }}
              </h3>
              <p class="text-[11px] text-blue-600 dark:text-blue-400 font-semibold mt-3 group-hover:underline flex items-center gap-1">
                <span>Lihat Rincian Slip Gaji</span>
                <ArrowRight class="w-3 h-3" />
              </p>
            </div>
            <div class="w-11 h-11 rounded-xl bg-purple-50 text-purple-600 dark:bg-purple-950/60 dark:text-purple-400 flex items-center justify-center shrink-0">
              <Briefcase class="w-5 h-5" />
            </div>
          </div>
        </div>

        <!-- Card 3: Profil Pegawai -->
        <div 
          @click="openProfileModal"
          class="p-5 rounded-2xl border transition-all duration-200 hover:shadow-xs group cursor-pointer relative"
          style="background-color: var(--bg-primary); border-color: var(--border-color);"
        >
          <div class="flex items-start justify-between">
            <div class="min-w-0">
              <p class="text-xs font-semibold text-slate-500 dark:text-slate-400 tracking-wide uppercase">
                Identitas Staf
              </p>
              <h3 class="text-lg font-bold tracking-tight text-slate-900 dark:text-slate-50 mt-2 truncate">
                {{ portalData.employee?.full_name || userFullName }}
              </h3>
              <p class="text-[11px] text-slate-400 mt-0.5 truncate">
                {{ portalData.employee?.department || 'Operasional' }} • {{ portalData.employee?.job_title || currentRoleName }}
              </p>
              <p class="text-[11px] text-blue-600 dark:text-blue-400 font-semibold mt-2.5 group-hover:underline flex items-center gap-1">
                <span>Kelola Data Lengkap</span>
                <ArrowRight class="w-3 h-3" />
              </p>
            </div>
            <div class="w-11 h-11 rounded-xl bg-emerald-50 text-emerald-600 dark:bg-emerald-950/60 dark:text-emerald-400 flex items-center justify-center shrink-0">
              <UserCheck class="w-5 h-5" />
            </div>
          </div>
        </div>
      </div>

      <!-- JADWAL ROSTER SHIFT MINGGUAN -->
      <div 
        class="p-6 rounded-2xl border shadow-xs"
        style="background-color: var(--bg-primary); border-color: var(--border-color);"
      >
        <div class="flex items-center justify-between mb-4">
          <div>
            <h3 class="text-base font-bold text-slate-900 dark:text-slate-50">
              Jadwal Shift Kerja Mingguan
            </h3>
            <p class="text-xs text-slate-400 mt-0.5">Roster jam kerja operasional Anda minggu ini</p>
          </div>
          <router-link to="/hris/shifts" class="text-xs font-semibold text-blue-600 dark:text-blue-400 hover:underline">
            Lihat Kalender Lengkap →
          </router-link>
        </div>

        <div class="grid grid-cols-2 sm:grid-cols-4 lg:grid-cols-7 gap-3">
          <div 
            v-for="(sched, idx) in weeklySchedules" 
            :key="idx"
            class="p-3.5 rounded-xl border text-center transition-all relative overflow-hidden"
            :class="sched.is_today 
              ? 'border-blue-500 bg-blue-50/70 dark:bg-blue-950/40 shadow-xs' 
              : 'border-slate-200/80 dark:border-slate-800 bg-slate-50/50 dark:bg-slate-900/40'"
          >
            <div v-if="sched.is_today" class="absolute top-1.5 right-1.5 w-2 h-2 rounded-full bg-blue-600 animate-ping"></div>
            <span class="text-xs font-semibold text-slate-500 dark:text-slate-400 block uppercase">
              {{ formatDayName(sched.day_name) }}
            </span>
            <div class="my-1.5">
              <span 
                class="px-2 py-0.5 rounded-md text-[11px] font-bold"
                :class="sched.shift_name === 'Libur' 
                  ? 'bg-slate-200 text-slate-600 dark:bg-slate-800 dark:text-slate-400' 
                  : 'bg-blue-100 text-blue-700 dark:bg-blue-900 dark:text-blue-200'"
              >
                {{ sched.shift_name }}
              </span>
            </div>
            <div class="text-[11px] font-mono text-slate-500 dark:text-slate-400">
              {{ sched.start_time }} - {{ sched.end_time }}
            </div>
          </div>
        </div>
      </div>

      <!-- ROLE FAST SHORTCUTS -->
      <div class="grid grid-cols-1 sm:grid-cols-3 gap-4">
        <router-link 
          to="/pos"
          class="p-4 rounded-xl border flex items-center justify-between transition-all hover:border-blue-400 hover:shadow-xs group cursor-pointer"
          style="background-color: var(--bg-primary); border-color: var(--border-color);"
        >
          <div class="flex items-center gap-3">
            <div class="w-10 h-10 rounded-lg bg-blue-50 text-blue-600 dark:bg-blue-950/50 dark:text-blue-400 flex items-center justify-center">
              <ShoppingCart class="w-5 h-5" />
            </div>
            <div>
              <div class="text-xs font-bold text-slate-900 dark:text-slate-100 group-hover:text-blue-600 transition-colors">Point of Sale (POS)</div>
              <div class="text-[11px] text-slate-400">Buka kasir & terima pesanan</div>
            </div>
          </div>
          <ArrowRight class="w-4 h-4 text-slate-400 group-hover:text-blue-600 transition-transform group-hover:translate-x-1" />
        </router-link>

        <router-link 
          to="/pos/kds"
          class="p-4 rounded-xl border flex items-center justify-between transition-all hover:border-blue-400 hover:shadow-xs group cursor-pointer"
          style="background-color: var(--bg-primary); border-color: var(--border-color);"
        >
          <div class="flex items-center gap-3">
            <div class="w-10 h-10 rounded-lg bg-amber-50 text-amber-600 dark:bg-amber-950/50 dark:text-amber-400 flex items-center justify-center">
              <MonitorDot class="w-5 h-5" />
            </div>
            <div>
              <div class="text-xs font-bold text-slate-900 dark:text-slate-100 group-hover:text-blue-600 transition-colors">Kitchen Display (KDS)</div>
              <div class="text-[11px] text-slate-400">Antrean tiket pesanan barista</div>
            </div>
          </div>
          <ArrowRight class="w-4 h-4 text-slate-400 group-hover:text-blue-600 transition-transform group-hover:translate-x-1" />
        </router-link>

        <router-link 
          to="/inventory"
          class="p-4 rounded-xl border flex items-center justify-between transition-all hover:border-blue-400 hover:shadow-xs group cursor-pointer"
          style="background-color: var(--bg-primary); border-color: var(--border-color);"
        >
          <div class="flex items-center gap-3">
            <div class="w-10 h-10 rounded-lg bg-emerald-50 text-emerald-600 dark:bg-emerald-950/50 dark:text-emerald-400 flex items-center justify-center">
              <Package class="w-5 h-5" />
            </div>
            <div>
              <div class="text-xs font-bold text-slate-900 dark:text-slate-100 group-hover:text-blue-600 transition-colors">Inventaris & Gudang</div>
              <div class="text-[11px] text-slate-400">Cek stok bahan & mutasi</div>
            </div>
          </div>
          <ArrowRight class="w-4 h-4 text-slate-400 group-hover:text-blue-600 transition-transform group-hover:translate-x-1" />
        </router-link>
      </div>
    </div>

    <!-- ===================================================================== -->
    <!-- 3. FINANCE DASHBOARD (AKUNTAN / KEUANGAN)                              -->
    <!-- ===================================================================== -->
    <div v-else-if="activePerspective === 'finance'" class="space-y-6">
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-5">
        <StatCardV2
          title="Total Pendapatan"
          :value="`Rp ${formatNum(financeOverview.total_revenue)}`"
          trend="up"
          trendValue="+14.2%"
          subtitle="Bulan berjalan"
          color="emerald"
          :icon="DollarSign"
        />
        <StatCardV2
          title="Beban Operasional"
          :value="`Rp ${formatNum(financeOverview.total_expenses)}`"
          trend="neutral"
          trendValue="Normal"
          subtitle="Bahan baku & utilitas"
          color="rose"
          :icon="Receipt"
        />
        <StatCardV2
          title="Laba Bersih"
          :value="`Rp ${formatNum(financeOverview.net_profit)}`"
          trend="up"
          trendValue="+18.5%"
          subtitle="Margin 39.0%"
          color="blue"
          :icon="TrendingUp"
        />
        <StatCardV2
          title="Saldo Kas & Bank"
          :value="`Rp ${formatNum(financeOverview.cash_on_hand)}`"
          subtitle="Likuiditas lancar"
          color="indigo"
          :icon="Landmark"
        />
      </div>

      <!-- Akun Kas & Rekening BCA + Jurnal Terakhir -->
      <div class="grid grid-cols-1 lg:grid-cols-12 gap-6">
        <!-- Kas & Bank Cards (5 cols) -->
        <div 
          class="lg:col-span-5 p-6 rounded-2xl border shadow-xs space-y-4"
          style="background-color: var(--bg-primary); border-color: var(--border-color);"
        >
          <div class="flex items-center justify-between">
            <h3 class="text-base font-bold text-slate-900 dark:text-slate-50">Saldo Rekening & Kas</h3>
            <span class="text-xs font-semibold text-emerald-600">Reconciled</span>
          </div>

          <div class="space-y-3">
            <div class="p-4 rounded-xl border flex items-center justify-between"
                 style="background-color: var(--bg-content); border-color: var(--border-color);">
              <div class="flex items-center gap-3">
                <div class="w-9 h-9 rounded-lg bg-blue-100 text-blue-700 flex items-center justify-center font-bold text-xs">
                  BCA
                </div>
                <div>
                  <div class="text-xs font-bold text-slate-900 dark:text-slate-100">Bank BCA Operasional</div>
                  <div class="text-[11px] font-mono text-slate-400">882-0192-331</div>
                </div>
              </div>
              <div class="text-right">
                <div class="text-xs font-bold font-mono text-slate-900 dark:text-slate-100">Rp 24.850.000</div>
                <div class="text-[10px] text-emerald-600 font-semibold">Aktif</div>
              </div>
            </div>

            <div class="p-4 rounded-xl border flex items-center justify-between"
                 style="background-color: var(--bg-content); border-color: var(--border-color);">
              <div class="flex items-center gap-3">
                <div class="w-9 h-9 rounded-lg bg-amber-100 text-amber-700 flex items-center justify-center font-bold text-xs">
                  CASH
                </div>
                <div>
                  <div class="text-xs font-bold text-slate-900 dark:text-slate-100">Petty Cash Kasir</div>
                  <div class="text-[11px] text-slate-400">Drawer Outlet Senopati</div>
                </div>
              </div>
              <div class="text-right">
                <div class="text-xs font-bold font-mono text-slate-900 dark:text-slate-100">Rp 7.650.000</div>
                <div class="text-[10px] text-emerald-600 font-semibold">Aktif</div>
              </div>
            </div>
          </div>

          <div class="pt-2">
            <router-link to="/finance/journal" class="w-full py-2 px-3 rounded-xl border text-xs font-semibold text-blue-600 dark:text-blue-400 flex items-center justify-center gap-1.5 hover:bg-blue-50 dark:hover:bg-blue-950/40 transition-colors">
              <BookOpen class="w-3.5 h-3.5" />
              <span>Buat Jurnal Manual Baru</span>
            </router-link>
          </div>
        </div>

        <!-- Jurnal Terbaru (7 cols) -->
        <div 
          class="lg:col-span-7 p-6 rounded-2xl border shadow-xs"
          style="background-color: var(--bg-primary); border-color: var(--border-color);"
        >
          <div class="flex items-center justify-between mb-4">
            <div>
              <h3 class="text-base font-bold text-slate-900 dark:text-slate-50">Jurnal Transaksi Terbaru</h3>
              <p class="text-xs text-slate-400">Entri debit dan kredit sistem akuntansi</p>
            </div>
            <router-link to="/finance" class="text-xs font-semibold text-blue-600 hover:underline">
              Buku Besar →
            </router-link>
          </div>

          <div class="space-y-2.5">
            <div 
              v-for="j in dummyJournals" 
              :key="j.ref"
              class="p-3.5 rounded-xl border flex items-center justify-between text-xs"
              style="background-color: var(--bg-content); border-color: var(--border-color);"
            >
              <div class="min-w-0 pr-3">
                <div class="flex items-center gap-2 mb-0.5">
                  <span class="font-mono text-[11px] font-bold text-blue-600">{{ j.ref }}</span>
                  <span class="text-[10px] text-slate-400">{{ j.date }}</span>
                </div>
                <div class="font-medium text-slate-700 dark:text-slate-300 truncate">{{ j.description }}</div>
              </div>
              <div class="text-right shrink-0">
                <div class="font-mono font-bold text-slate-900 dark:text-slate-100">
                  Rp {{ formatNum(j.debit) }}
                </div>
                <span class="text-[10px] text-emerald-600 font-semibold uppercase">Posted</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- ===================================================================== -->
    <!-- 4. HRIS DASHBOARD (SDM & KEPEGAWAIAN)                                  -->
    <!-- ===================================================================== -->
    <div v-else-if="activePerspective === 'hris'" class="space-y-6">
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-5">
        <StatCardV2
          title="Total Staf Aktif"
          :value="hrisTodayStats.total_employees"
          subtitle="Outlet Senopati"
          color="blue"
          :icon="Users"
        />
        <StatCardV2
          title="Hadir Hari Ini"
          :value="hrisTodayStats.total_present"
          trend="up"
          trendValue="100%"
          subtitle="Shift Pagi Operasional"
          color="emerald"
          :icon="UserCheck"
        />
        <StatCardV2
          title="Tepat Waktu"
          :value="hrisTodayStats.on_time"
          subtitle="Nol keterlambatan"
          color="indigo"
          :icon="CheckCircle"
        />
        <StatCardV2
          title="Permohonan Cuti"
          :value="`${pendingLeavesList.length} Pending`"
          subtitle="Menunggu persetujuan"
          color="amber"
          :icon="CalendarClock"
        />
      </div>

      <div class="grid grid-cols-1 lg:grid-cols-12 gap-6">
        <!-- Live Attendance Staff Hari Ini (7 cols) -->
        <div 
          class="lg:col-span-7 p-6 rounded-2xl border shadow-xs"
          style="background-color: var(--bg-primary); border-color: var(--border-color);"
        >
          <div class="flex items-center justify-between mb-4">
            <div>
              <h3 class="text-base font-bold text-slate-900 dark:text-slate-50">Presensi Tim Live Hari Ini</h3>
              <p class="text-xs text-slate-400">Monitoring kehadiran staf kafe real-time</p>
            </div>
            <router-link to="/hris/attendance" class="text-xs font-semibold text-blue-600 hover:underline">
              Rekap Bulanan →
            </router-link>
          </div>

          <div class="space-y-2.5">
            <div 
              v-for="emp in staffAttendanceList" 
              :key="emp.id"
              class="p-3.5 rounded-xl border flex items-center justify-between text-xs"
              style="background-color: var(--bg-content); border-color: var(--border-color);"
            >
              <div class="flex items-center gap-3">
                <div class="w-9 h-9 rounded-xl bg-blue-100 text-blue-700 font-bold flex items-center justify-center text-xs">
                  {{ emp.name[0] }}
                </div>
                <div>
                  <div class="font-bold text-slate-900 dark:text-slate-100">{{ emp.name }}</div>
                  <div class="text-[11px] text-slate-400">{{ emp.role }} • {{ emp.shift }}</div>
                </div>
              </div>
              <div class="text-right">
                <StatusBadgeV2 :status="emp.status === 'Hadir' ? 'active' : 'pending'" :label="emp.status" />
                <div class="text-[10px] text-slate-400 mt-1">Clock In: {{ emp.clockIn }}</div>
              </div>
            </div>
          </div>
        </div>

        <!-- Pending Leave Approvals (5 cols) -->
        <div 
          class="lg:col-span-5 p-6 rounded-2xl border shadow-xs flex flex-col justify-between"
          style="background-color: var(--bg-primary); border-color: var(--border-color);"
        >
          <div>
            <div class="flex items-center justify-between mb-3">
              <h3 class="text-base font-bold text-slate-900 dark:text-slate-50">Persetujuan Cuti</h3>
              <span class="text-xs font-bold text-amber-600 bg-amber-50 dark:bg-amber-950/40 px-2 py-0.5 rounded-md border border-amber-200">
                {{ pendingLeavesList.length }} Pending
              </span>
            </div>
            <p class="text-xs text-slate-400 mb-4">Verifikasi permohonan istirahat dan izin staf</p>

            <div v-if="pendingLeavesList.length" class="space-y-3">
              <div 
                v-for="pl in pendingLeavesList" 
                :key="pl.id"
                class="p-3.5 rounded-xl border text-xs space-y-2"
                style="background-color: var(--bg-content); border-color: var(--border-color);"
              >
                <div class="flex items-center justify-between font-bold">
                  <span class="text-slate-900 dark:text-slate-100">{{ pl.name }}</span>
                  <span class="text-blue-600">{{ pl.leave_type }}</span>
                </div>
                <div class="text-[11px] text-slate-500">
                  {{ pl.start_date }} s/d {{ pl.end_date }} ({{ pl.total_days }} Hari)
                </div>
                <div class="text-[11px] italic text-slate-400">
                  "{{ pl.reason }}"
                </div>
                <div class="pt-2 border-t flex items-center justify-end gap-2" style="border-color: var(--border-color);">
                  <button 
                    type="button"
                    @click="approveLeave(pl.id)"
                    class="px-3 py-1 rounded-lg text-xs font-bold bg-emerald-600 hover:bg-emerald-700 text-white cursor-pointer"
                  >
                    Setujui
                  </button>
                  <button 
                    type="button"
                    @click="rejectLeave(pl.id)"
                    class="px-3 py-1 rounded-lg text-xs font-bold bg-rose-600 hover:bg-rose-700 text-white cursor-pointer"
                  >
                    Tolak
                  </button>
                </div>
              </div>
            </div>

            <div v-else class="p-8 text-center text-xs text-slate-400 border border-dashed rounded-xl" style="border-color: var(--border-color);">
              ✨ Tidak ada permohonan cuti yang menunggu saat ini.
            </div>
          </div>

          <div class="mt-4 pt-3 border-t text-right" style="border-color: var(--border-color);">
            <router-link to="/hris/leaves" class="text-xs font-semibold text-blue-600 hover:underline">
              Kelola Pengajuan Cuti Seluruh Staf →
            </router-link>
          </div>
        </div>
      </div>
    </div>

    <!-- ===================================================================== -->
    <!-- 5. EXECUTIVE OVERVIEW (SUPER ADMIN / OWNER / STORE MANAGER)           -->
    <!-- ===================================================================== -->
    <div v-else class="space-y-6">
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-5">
        <StatCardV2
          title="Penjualan Hari Ini"
          :value="`Rp ${formatNum(stats.today_sales)}`"
          trend="up"
          trendValue="+12.4%"
          subtitle="48 pesanan kasir"
          color="blue"
          :icon="DollarSign"
        />
        <StatCardV2
          title="Total Transaksi"
          :value="stats.total_orders"
          subtitle="Rata-rata Rp 309.375"
          color="emerald"
          :icon="Receipt"
        />
        <StatCardV2
          title="Kapasitas Meja Terisi"
          :value="`${stats.active_tables} / 8 Meja`"
          trend="neutral"
          trendValue="50% Okupansi"
          subtitle="Ruang Indoor & Outdoor"
          color="amber"
          :icon="Armchair"
        />
        <StatCardV2
          title="Bahan Baku Kritis"
          :value="`${stats.low_stock_count} Item`"
          trend="down"
          trendValue="Perlu Restok"
          subtitle="Sirup, Susu Fresh, Mint"
          color="rose"
          :icon="Package"
        />
      </div>

      <!-- Denah Meja & POS Quick Status -->
      <div 
        class="p-6 rounded-2xl border shadow-xs"
        style="background-color: var(--bg-primary); border-color: var(--border-color);"
      >
        <div class="flex items-center justify-between mb-4">
          <div>
            <h3 class="text-base font-bold text-slate-900 dark:text-slate-50">Status Meja Kasir POS</h3>
            <p class="text-xs text-slate-400">Monitoring meja tamu yang sedang aktif di area dine-in</p>
          </div>
          <router-link to="/pos/tables" class="text-xs font-semibold text-blue-600 hover:underline">
            Denah Meja Lengkap →
          </router-link>
        </div>

        <div class="grid grid-cols-2 sm:grid-cols-4 lg:grid-cols-8 gap-3">
          <div 
            v-for="tbl in defaultTables" 
            :key="tbl.id"
            class="p-3 rounded-xl border text-center transition-all cursor-pointer hover:border-blue-400 hover:shadow-xs"
            :class="tbl.status === 'occupied' 
              ? 'border-amber-300 bg-amber-50/70 dark:bg-amber-950/30' 
              : tbl.status === 'reserved'
              ? 'border-blue-300 bg-blue-50/70 dark:bg-blue-950/30'
              : 'border-slate-200 dark:border-slate-800 bg-slate-50/50 dark:bg-slate-900/40'"
          >
            <div class="font-bold text-sm text-slate-900 dark:text-slate-100">{{ tbl.table_number }}</div>
            <div class="text-[10px] text-slate-400 mt-0.5">{{ tbl.capacity }} Kursi</div>
            <div class="mt-2">
              <span 
                class="px-2 py-0.5 rounded text-[10px] font-bold uppercase tracking-wider"
                :class="tbl.status === 'occupied' ? 'bg-amber-100 text-amber-700' :
                        tbl.status === 'reserved' ? 'bg-blue-100 text-blue-700' :
                        'bg-emerald-100 text-emerald-700'"
              >
                {{ tbl.status === 'occupied' ? 'Terisi' : tbl.status === 'reserved' ? 'Booking' : 'Kosong' }}
              </span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- ===================================================================== -->
    <!-- 6. MODALS (100% FEATURE PARITY DENGAN V1)                             -->
    <!-- ===================================================================== -->
    
    <!-- Modal Ajukan Cuti -->
    <Teleport to="body">
      <div v-if="showLeaveModal" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-950/50 backdrop-blur-sm animate-in fade-in">
        <div 
          class="w-full max-w-lg p-6 rounded-3xl border shadow-xl space-y-4 animate-in zoom-in-95 duration-150"
          style="background-color: var(--bg-primary); border-color: var(--border-color);"
        >
          <div class="flex items-center justify-between pb-3 border-b" style="border-color: var(--border-color);">
            <div>
              <h3 class="text-base font-bold text-slate-900 dark:text-slate-50">Pengajuan Cuti Karyawan</h3>
              <p class="text-xs text-slate-400">Sisa kuota tahunan: {{ portalData.leave_summary?.remaining_days ?? 12 }} hari</p>
            </div>
            <button @click="showLeaveModal = false" class="p-1.5 rounded-lg text-slate-400 hover:text-slate-600 cursor-pointer">
              <X class="w-4 h-4" />
            </button>
          </div>

          <form @submit.prevent="submitLeaveForm" class="space-y-3.5 text-xs">
            <div>
              <label class="block font-semibold mb-1 text-slate-700 dark:text-slate-300">Jenis Cuti</label>
              <select 
                v-model="leaveForm.leave_type"
                class="w-full px-3 py-2 rounded-xl border outline-none"
                style="background-color: var(--input-bg); border-color: var(--input-border); color: var(--input-text);"
              >
                <option value="Cuti Tahunan">Cuti Tahunan</option>
                <option value="Cuti Sakit">Cuti Sakit</option>
                <option value="Cuti Khusus/Penting">Cuti Khusus / Keperluan Mendesak</option>
              </select>
            </div>

            <div class="grid grid-cols-2 gap-3">
              <div>
                <label class="block font-semibold mb-1 text-slate-700 dark:text-slate-300">Tanggal Mulai</label>
                <input 
                  type="date" 
                  v-model="leaveForm.start_date" 
                  required
                  class="w-full px-3 py-2 rounded-xl border outline-none"
                  style="background-color: var(--input-bg); border-color: var(--input-border); color: var(--input-text);"
                />
              </div>
              <div>
                <label class="block font-semibold mb-1 text-slate-700 dark:text-slate-300">Tanggal Selesai</label>
                <input 
                  type="date" 
                  v-model="leaveForm.end_date" 
                  required
                  class="w-full px-3 py-2 rounded-xl border outline-none"
                  style="background-color: var(--input-bg); border-color: var(--input-border); color: var(--input-text);"
                />
              </div>
            </div>

            <div>
              <label class="block font-semibold mb-1 text-slate-700 dark:text-slate-300">Alasan Cuti</label>
              <textarea 
                v-model="leaveForm.reason" 
                rows="3" 
                required
                placeholder="Jelaskan keperluan pengajuan cuti Anda..."
                class="w-full px-3 py-2 rounded-xl border outline-none"
                style="background-color: var(--input-bg); border-color: var(--input-border); color: var(--input-text);"
              ></textarea>
            </div>

            <div class="pt-3 border-t flex justify-end gap-2.5" style="border-color: var(--border-color);">
              <button 
                type="button" 
                @click="showLeaveModal = false" 
                class="px-4 py-2 rounded-xl border text-xs font-semibold hover:bg-slate-100 dark:hover:bg-slate-800 cursor-pointer"
                style="border-color: var(--border-color);"
              >
                Batal
              </button>
              <button 
                type="submit" 
                :disabled="leaveSubmitting"
                class="px-5 py-2 rounded-xl text-xs font-semibold bg-blue-600 hover:bg-blue-700 text-white cursor-pointer shadow-xs disabled:opacity-50"
              >
                {{ leaveSubmitting ? 'Mengirim...' : 'Kirim Pengajuan Cuti' }}
              </button>
            </div>
          </form>
        </div>
      </div>
    </Teleport>

    <!-- Modal Slip Gaji Elektronik -->
    <Teleport to="body">
      <div v-if="showSalarySlipModal" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-950/50 backdrop-blur-sm animate-in fade-in">
        <div 
          class="w-full max-w-md p-6 rounded-3xl border shadow-xl space-y-4 animate-in zoom-in-95 duration-150"
          style="background-color: var(--bg-primary); border-color: var(--border-color);"
        >
          <div class="flex items-center justify-between pb-3 border-b" style="border-color: var(--border-color);">
            <div>
              <h3 class="text-base font-bold text-slate-900 dark:text-slate-50">Slip Gaji Elektronik</h3>
              <p class="text-xs text-slate-400">Periode: September 2026</p>
            </div>
            <button @click="showSalarySlipModal = false" class="p-1.5 rounded-lg text-slate-400 hover:text-slate-600 cursor-pointer">
              <X class="w-4 h-4" />
            </button>
          </div>

          <div class="space-y-3 text-xs">
            <div class="p-3.5 rounded-xl border flex items-center justify-between"
                 style="background-color: var(--bg-content); border-color: var(--border-color);">
              <div>
                <div class="font-bold text-slate-900 dark:text-slate-100">{{ portalData.employee?.full_name || userFullName }}</div>
                <div class="text-[11px] text-slate-400">NIK: {{ portalData.employee?.nik || 'EMP-001' }}</div>
              </div>
              <div class="text-right">
                <span class="text-[10px] font-bold text-emerald-600 bg-emerald-50 px-2 py-0.5 rounded border border-emerald-200">
                  LUNAS / PAID
                </span>
              </div>
            </div>

            <div class="space-y-2 border-b pb-3" style="border-color: var(--border-color);">
              <div class="flex justify-between font-semibold text-slate-700 dark:text-slate-300">
                <span>Gaji Pokok:</span>
                <span class="font-mono">Rp 4.500.000</span>
              </div>
              <div class="flex justify-between text-slate-600 dark:text-slate-400">
                <span>Tunjangan Kehadiran:</span>
                <span class="font-mono">Rp 350.000</span>
              </div>
              <div class="flex justify-between text-slate-600 dark:text-slate-400">
                <span>Uang Makan & Lembur:</span>
                <span class="font-mono">Rp 200.000</span>
              </div>
              <div class="flex justify-between text-rose-600 font-medium">
                <span>Potongan (BPJS & PPh21):</span>
                <span class="font-mono">- Rp 200.000</span>
              </div>
            </div>

            <div class="p-4 rounded-xl bg-blue-50 dark:bg-blue-950/40 border border-blue-200/80 dark:border-blue-800/60 flex items-center justify-between">
              <div>
                <span class="text-xs font-bold text-blue-900 dark:text-blue-200 block">Take Home Pay</span>
                <span class="text-[10px] text-blue-700 dark:text-blue-400">Transfer ke Rekening Mandiri</span>
              </div>
              <div class="text-lg font-black font-mono text-blue-700 dark:text-blue-300">
                Rp {{ formatNum(portalData.latest_payroll?.net_salary || 4850000) }}
              </div>
            </div>
          </div>

          <div class="pt-3 border-t text-right" style="border-color: var(--border-color);">
            <button 
              type="button" 
              @click="showSalarySlipModal = false" 
              class="px-5 py-2 rounded-xl text-xs font-semibold bg-blue-600 hover:bg-blue-700 text-white cursor-pointer shadow-xs"
            >
              Tutup Slip Gaji
            </button>
          </div>
        </div>
      </div>
    </Teleport>

    <!-- Modal Profil Staf -->
    <Teleport to="body">
      <div v-if="showProfileModal" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-950/50 backdrop-blur-sm animate-in fade-in">
        <div 
          class="w-full max-w-md p-6 rounded-3xl border shadow-xl space-y-4 animate-in zoom-in-95 duration-150"
          style="background-color: var(--bg-primary); border-color: var(--border-color);"
        >
          <div class="flex items-center justify-between pb-3 border-b" style="border-color: var(--border-color);">
            <h3 class="text-base font-bold text-slate-900 dark:text-slate-50">Profil Karyawan</h3>
            <button @click="showProfileModal = false" class="p-1.5 rounded-lg text-slate-400 hover:text-slate-600 cursor-pointer">
              <X class="w-4 h-4" />
            </button>
          </div>

          <div class="space-y-3 text-xs">
            <div class="flex items-center gap-3 p-3.5 rounded-xl border"
                 style="background-color: var(--bg-content); border-color: var(--border-color);">
              <div class="w-12 h-12 rounded-xl bg-blue-600 text-white font-bold text-lg flex items-center justify-center">
                {{ userInitials }}
              </div>
              <div>
                <div class="text-sm font-bold text-slate-900 dark:text-slate-100">{{ portalData.employee?.full_name || userFullName }}</div>
                <div class="text-xs text-blue-600 font-semibold">{{ currentRoleName }}</div>
                <div class="text-[11px] text-slate-400">Senopati Central Branch</div>
              </div>
            </div>

            <div class="space-y-2 border-t pt-3" style="border-color: var(--border-color);">
              <div class="flex justify-between py-1 border-b border-dashed" style="border-color: var(--border-color);">
                <span class="text-slate-400">Nomor Induk Karyawan:</span>
                <span class="font-mono font-bold">{{ portalData.employee?.nik || 'EMP-001' }}</span>
              </div>
              <div class="flex justify-between py-1 border-b border-dashed" style="border-color: var(--border-color);">
                <span class="text-slate-400">Status Kepegawaian:</span>
                <span class="font-semibold text-emerald-600">Tetap (Permanent)</span>
              </div>
              <div class="flex justify-between py-1 border-b border-dashed" style="border-color: var(--border-color);">
                <span class="text-slate-400">Email Kerja:</span>
                <span class="font-medium">{{ authStore.user?.email || 'karyawan@cafe-erp.local' }}</span>
              </div>
              <div class="flex justify-between py-1">
                <span class="text-slate-400">Sisa Cuti Tahunan:</span>
                <span class="font-bold text-blue-600">{{ portalData.leave_summary?.remaining_days ?? 12 }} Hari</span>
              </div>
            </div>
          </div>

          <div class="pt-3 border-t text-right" style="border-color: var(--border-color);">
            <button 
              type="button" 
              @click="showProfileModal = false" 
              class="px-5 py-2 rounded-xl text-xs font-semibold bg-blue-600 hover:bg-blue-700 text-white cursor-pointer shadow-xs"
            >
              Tutup
            </button>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useAuthStore } from '@/stores/auth.store'
import { useNotificationStore } from '@/stores/notification.store'
import api from '@/plugins/axios'
import StatCardV2 from '@/components/v2/common/StatCardV2.vue'
import StatusBadgeV2 from '@/components/v2/common/StatusBadgeV2.vue'
import {
  Clock,
  Calendar,
  CalendarClock,
  UserCheck,
  CheckCircle,
  CheckCircle2,
  Users,
  Briefcase,
  RefreshCw,
  ArrowRight,
  User,
  X,
  ShoppingCart,
  MonitorDot,
  Package,
  DollarSign,
  Receipt,
  TrendingUp,
  Landmark,
  BookOpen,
  Armchair
} from 'lucide-vue-next'

const authStore = useAuthStore()
const notifyStore = useNotificationStore()

// State
const activeViewTab = ref<'employee' | 'finance' | 'hris' | 'executive'>('executive')
const isRefreshing = ref(false)

// Digital Clock
const liveTime = ref('')
const liveDate = ref('')
let clockTimer: any = null

// Employee Portal Data
const portalData = ref<any>({
  employee: null,
  today_attendance: null,
  leave_summary: null,
  latest_payroll: null,
  weekly_schedule: []
})

// Attendance Action State
const attendanceActionLoading = ref(false)
const attendanceNoteInput = ref('')

// Finance Data
const financeOverview = ref<any>({
  total_revenue: 128450000,
  total_expenses: 78320000,
  net_profit: 50130000,
  cash_on_hand: 32500000,
  recent_journals: []
})

// HRIS Data
const hrisTodayStats = ref<any>({
  total_employees: 7,
  total_present: 3,
  on_time: 3,
  late: 0,
  on_leave: 0
})
const pendingLeavesList = ref<any[]>([])

// Executive Stats
const stats = ref({
  today_sales: 14850000,
  total_orders: 48,
  low_stock_count: 3,
  active_tables: 4,
  avg_order_value: 0
})

// Modals
const showLeaveModal = ref(false)
const showSalarySlipModal = ref(false)
const showProfileModal = ref(false)
const leaveSubmitting = ref(false)
const leaveForm = ref({
  leave_type: 'Cuti Tahunan',
  start_date: '',
  end_date: '',
  reason: ''
})

const defaultTables = [
  { id: '1', table_number: 'T-01', capacity: 2, status: 'available' },
  { id: '2', table_number: 'T-02', capacity: 2, status: 'available' },
  { id: '3', table_number: 'T-03', capacity: 2, status: 'available' },
  { id: '4', table_number: 'T-04', capacity: 4, status: 'occupied' },
  { id: '5', table_number: 'T-05', capacity: 4, status: 'occupied' },
  { id: '6', table_number: 'T-06', capacity: 6, status: 'available' },
  { id: '7', table_number: 'T-07', capacity: 8, status: 'reserved' },
  { id: '8', table_number: 'T-08', capacity: 6, status: 'available' }
]

const dummyJournals = [
  { ref: 'JV-EXP-002', date: '04 Okt 2026', description: 'Pembelian es batu kristal & mint lokal', debit: 750000, credit: 750000, status: 'posted' },
  { ref: 'JV-EXP-003', date: '03 Okt 2026', description: 'Servis rutin mesin espresso La Marzocco', debit: 1250000, credit: 1250000, status: 'posted' },
  { ref: 'JV-SET-001', date: '02 Okt 2026', description: 'Setoran Kasir POS EDC BCA Weekend', debit: 5120000, credit: 5120000, status: 'posted' },
  { ref: 'JV-PAY-008', date: '28 Sep 2026', description: 'Pencairan Payroll Gaji Staf Operasional', debit: 28260969, credit: 28260969, status: 'posted' }
]

const staffAttendanceList = [
  { id: '1', name: 'Budi Pratama', role: 'Head Barista', shift: 'Pagi (07-15)', status: 'Hadir', clockIn: '06:52' },
  { id: '2', name: 'Jane Doe', role: 'Kasir POS', shift: 'Pagi (07-15)', status: 'Hadir', clockIn: '06:58' },
  { id: '3', name: 'Ahmad Gudang', role: 'Warehouse Staff', shift: 'Pagi (07-15)', status: 'Hadir', clockIn: '07:05' },
  { id: '4', name: 'John Doe', role: 'Head Chef', shift: 'Siang (15-23)', status: 'Jadwal Siang', clockIn: '-' },
  { id: '5', name: 'Sarah Manager', role: 'Store Manager', shift: 'Pagi (07-15)', status: 'Hadir', clockIn: '06:45' }
]

const perspectiveTabs: { key: 'employee' | 'finance' | 'hris' | 'executive'; label: string; icon: any }[] = [
  { key: 'executive', label: 'Executive & POS', icon: Armchair },
  { key: 'employee', label: 'Portal Karyawan', icon: UserCheck },
  { key: 'finance', label: 'Finance & Kas', icon: DollarSign },
  { key: 'hris', label: 'SDM & Presensi Tim', icon: Users }
]

// Computed Helpers
const userFullName = computed(() => authStore.fullName || authStore.user?.firstName || 'Pengguna')
const currentRoleName = computed(() => authStore.currentRole || 'Karyawan')
const roleLower = computed(() => String(authStore.currentRole || authStore.user?.roleId || '').toLowerCase())

const userInitials = computed(() => {
  const name = userFullName.value.trim()
  const parts = name.split(' ')
  if (parts.length >= 2) return `${parts[0][0]}${parts[1][0]}`.toUpperCase()
  return name.slice(0, 2).toUpperCase() || 'U'
})

const isOperationalStaff = computed(() => {
  const r = roleLower.value
  return r.includes('barista') || r.includes('kasir') || r.includes('cashier') ||
         r.includes('warehouse') || r.includes('gudang') || r.includes('inventory') ||
         r.includes('pelayan') || r.includes('waiter') || r.includes('kitchen')
})

const isFinanceRole = computed(() => {
  const r = roleLower.value
  return r.includes('akuntan') || r.includes('finance') || r.includes('keuangan')
})

const isHRISRole = computed(() => {
  const r = roleLower.value
  return r.includes('hr') || r.includes('human') || r.includes('hris')
})

const isManagerialRole = computed(() => {
  const r = roleLower.value
  return r.includes('admin') || r.includes('owner') || r.includes('manager')
})

const activePerspective = computed(() => {
  if (isManagerialRole.value) {
    return activeViewTab.value
  }
  if (isOperationalStaff.value) return 'employee'
  if (isFinanceRole.value) return 'finance'
  if (isHRISRole.value) return 'hris'
  return 'executive'
})

const greetingTime = computed(() => {
  const hour = new Date().getHours()
  if (hour < 11) return 'Selamat Pagi'
  if (hour < 15) return 'Selamat Siang'
  if (hour < 18) return 'Selamat Sore'
  return 'Selamat Malam'
})

const attendanceState = computed(() => {
  const att = portalData.value.today_attendance
  if (!att) {
    return {
      hasClockedIn: false,
      hasClockedOut: false,
      clockInTime: '',
      clockOutTime: '',
      workedHours: '',
      notes: ''
    }
  }
  return {
    hasClockedIn: !!att.has_clocked_in,
    hasClockedOut: !!att.has_clocked_out,
    clockInTime: att.clock_in ? att.clock_in.slice(0, 5) : '',
    clockOutTime: att.clock_out ? att.clock_out.slice(0, 5) : '',
    workedHours: att.has_clocked_in && !att.has_clocked_out ? 'Sedang Berjalan' : 'Selesai Shift',
    notes: att.notes || ''
  }
})

const todayShiftName = computed(() => {
  const todayStr = new Date().toISOString().slice(0, 10)
  const todaySched = portalData.value.weekly_schedule?.find((s: any) => s.date === todayStr)
  if (todaySched) {
    return `${todaySched.shift_name} (${todaySched.start_time} - ${todaySched.end_time})`
  }
  return 'Shift Pagi (07:00 - 15:00)'
})

const weeklySchedules = computed(() => {
  if (portalData.value.weekly_schedule?.length) {
    return portalData.value.weekly_schedule
  }
  const days = ['Senin', 'Selasa', 'Rabu', 'Kamis', 'Jumat', 'Sabtu', 'Minggu']
  return days.map((d, i) => ({
    date: `2026-09-${28 + i}`,
    day_name: d,
    shift_name: i === 6 ? 'Libur' : i % 2 === 0 ? 'Pagi' : 'Siang',
    start_time: i === 6 ? '00:00' : i % 2 === 0 ? '07:00' : '15:00',
    end_time: i === 6 ? '00:00' : i % 2 === 0 ? '15:00' : '23:00',
    is_today: i === 6
  }))
})

// Methods
const updateClock = () => {
  const now = new Date()
  liveTime.value = now.toLocaleTimeString('id-ID', { hour12: false })
  liveDate.value = now.toLocaleDateString('id-ID', {
    weekday: 'long',
    day: 'numeric',
    month: 'long',
    year: 'numeric'
  })
}

const formatNum = (val: number | string) => Number(val || 0).toLocaleString('id-ID')

const formatDayName = (name: string) => {
  const map: Record<string, string> = {
    Monday: 'Senin',
    Tuesday: 'Selasa',
    Wednesday: 'Rabu',
    Thursday: 'Kamis',
    Friday: 'Jumat',
    Saturday: 'Sabtu',
    Sunday: 'Minggu'
  }
  return map[name] || name
}

const fetchPortalData = async () => {
  try {
    const res = await api.get('/hris/my-portal')
    if (res.data?.data) {
      portalData.value = res.data.data
    }
  } catch (err) {
    console.warn('Could not fetch employee portal data directly:', err)
  }
}

const handleClockIn = async () => {
  attendanceActionLoading.value = true
  try {
    await api.post('/hris/attendance/clock-in', {
      notes: attendanceNoteInput.value || 'Clock in via Welcome Portal V2'
    })
    notifyStore.success('Presensi Masuk Berhasil!', 'Selamat Bertugas')
    attendanceNoteInput.value = ''
    await fetchPortalData()
  } catch (err: any) {
    notifyStore.error(err.response?.data?.message || 'Gagal melakukan clock-in', 'Presensi')
  } finally {
    attendanceActionLoading.value = false
  }
}

const handleClockOut = async () => {
  attendanceActionLoading.value = true
  try {
    await api.post('/hris/attendance/clock-out', {
      notes: attendanceNoteInput.value || 'Clock out via Welcome Portal V2'
    })
    notifyStore.success('Presensi Keluar Berhasil! Terima kasih atas dedikasi Anda hari ini.', 'Shift Selesai')
    attendanceNoteInput.value = ''
    await fetchPortalData()
  } catch (err: any) {
    notifyStore.error(err.response?.data?.message || 'Gagal melakukan clock-out', 'Presensi')
  } finally {
    attendanceActionLoading.value = false
  }
}

const openLeaveModal = () => {
  const today = new Date().toISOString().slice(0, 10)
  leaveForm.value = {
    leave_type: 'Cuti Tahunan',
    start_date: today,
    end_date: today,
    reason: ''
  }
  showLeaveModal.value = true
}

const submitLeaveForm = async () => {
  leaveSubmitting.value = true
  try {
    await api.post('/hris/leaves', leaveForm.value)
    notifyStore.success('Permohonan cuti telah berhasil dikirim ke manajer HR!', 'Pengajuan Cuti')
    showLeaveModal.value = false
    await fetchPortalData()
  } catch (err: any) {
    notifyStore.error(err.response?.data?.message || 'Gagal mengajukan cuti', 'Pengajuan Cuti')
  } finally {
    leaveSubmitting.value = false
  }
}

const openSalarySlipModal = () => {
  showSalarySlipModal.value = true
}

const openProfileModal = () => {
  showProfileModal.value = true
}

const approveLeave = (id: string) => {
  pendingLeavesList.value = pendingLeavesList.value.filter(l => l.id !== id)
  notifyStore.success('Permohonan cuti staf telah disetujui', 'HRIS Approval')
}

const rejectLeave = (id: string) => {
  pendingLeavesList.value = pendingLeavesList.value.filter(l => l.id !== id)
  notifyStore.info('Permohonan cuti staf telah ditolak', 'HRIS Approval')
}

const refreshCurrentView = async () => {
  isRefreshing.value = true
  await fetchPortalData()
  setTimeout(() => {
    isRefreshing.value = false
    notifyStore.success('Data dashboard telah diperbarui', 'Refresh Berhasil')
  }, 400)
}

onMounted(() => {
  updateClock()
  clockTimer = setInterval(updateClock, 1000)
  fetchPortalData()

  // Set default view perspective based on logged-in role
  if (isOperationalStaff.value) {
    activeViewTab.value = 'employee'
  } else if (isFinanceRole.value) {
    activeViewTab.value = 'finance'
  } else if (isHRISRole.value) {
    activeViewTab.value = 'hris'
  } else {
    activeViewTab.value = 'executive'
  }
})

onUnmounted(() => {
  if (clockTimer) clearInterval(clockTimer)
})
</script>
