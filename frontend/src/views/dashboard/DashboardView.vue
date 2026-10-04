<template>
  <div class="space-y-6 select-none animate-in fade-in duration-300">
    <!-- TOP WELCOME BANNER & ROLE SWITCHER -->
    <div 
      class="rounded-3xl p-6 sm:p-8 border shadow-sm relative overflow-hidden transition-all duration-300"
      :style="{ backgroundColor: 'var(--bg-primary)', borderColor: 'var(--border-color)' }"
    >
      <!-- Background Ambient Aura -->
      <div class="absolute -top-24 -right-24 w-80 h-80 rounded-full blur-3xl pointer-events-none opacity-40"
           :class="isOperationalStaff ? 'bg-amber-400/20' : isFinanceRole ? 'bg-purple-500/20' : isHRISRole ? 'bg-rose-500/20' : 'bg-blue-500/20'">
      </div>

      <div class="relative z-10 flex flex-col md:flex-row md:items-center justify-between gap-6">
        <!-- User Greeting & Profile Snippet -->
        <div class="flex items-center gap-4">
          <div 
            class="w-14 h-14 sm:w-16 sm:h-16 rounded-2xl flex items-center justify-center text-2xl sm:text-3xl font-black shadow-md shrink-0 border"
            :class="isOperationalStaff ? 'bg-gradient-to-tr from-amber-500 to-orange-400 text-white border-amber-300' :
                    isFinanceRole ? 'bg-gradient-to-tr from-purple-600 to-indigo-500 text-white border-purple-400' :
                    isHRISRole ? 'bg-gradient-to-tr from-rose-600 to-pink-500 text-white border-rose-400' :
                    'bg-gradient-to-tr from-blue-600 to-cyan-500 text-white border-blue-400'"
          >
            <span v-if="portalData.employee?.first_name">{{ portalData.employee.first_name[0] }}</span>
            <span v-else>{{ userFullName[0] || 'U' }}</span>
          </div>

          <div class="min-w-0">
            <div class="flex flex-wrap items-center gap-2 mb-1">
              <span class="text-xs font-semibold px-2.5 py-0.5 rounded-full border"
                    :class="isOperationalStaff ? 'bg-amber-50 dark:bg-amber-950/40 text-amber-700 dark:text-amber-300 border-amber-200' :
                            isFinanceRole ? 'bg-purple-50 dark:bg-purple-950/40 text-purple-700 dark:text-purple-300 border-purple-200' :
                            isHRISRole ? 'bg-rose-50 dark:bg-rose-950/40 text-rose-700 dark:text-rose-300 border-rose-200' :
                            'bg-blue-50 dark:bg-blue-950/40 text-blue-700 dark:text-blue-300 border-blue-200'">
                ● {{ currentRoleName }}
              </span>
              <span v-if="portalData.employee?.nik" class="text-xs font-mono px-2 py-0.5 rounded-md border text-slate-500 dark:text-slate-400"
                    :style="{ backgroundColor: 'var(--bg-content)', borderColor: 'var(--border-color)' }">
                NIK: {{ portalData.employee.nik }}
              </span>
              <span class="text-xs px-2 py-0.5 rounded-md border text-emerald-600 bg-emerald-50 dark:bg-emerald-950/40 border-emerald-200 flex items-center gap-1 font-semibold">
                <span class="w-1.5 h-1.5 rounded-full bg-emerald-500 animate-pulse"></span> Staf Aktif
              </span>
            </div>

            <h1 class="text-xl sm:text-2xl font-black tracking-tight" :style="{ color: 'var(--text-primary)' }">
              {{ greetingTime }}, {{ portalData.employee?.full_name || userFullName }}! 👋
            </h1>
            <p class="text-xs sm:text-sm mt-0.5" :style="{ color: 'var(--text-muted)' }">
              <span v-if="isOperationalStaff">Selamat bertugas! Silakan catat presensi masuk/keluar dan periksa jadwal shift Anda.</span>
              <span v-else-if="isFinanceRole">Selamat datang di modul keuangan. Pantau arus kas, realisasi jurnal, dan ringkasan laba rugi.</span>
              <span v-else-if="isHRISRole">Pusat kendali SDM & HRIS. Pantau kehadiran tim kafe hari ini dan persetujuan cuti karyawan.</span>
              <span v-else>Ikhtisar performa penjualan kasir POS, mutasi bahan baku, dan monitoring outlet terpadu.</span>
            </p>
          </div>
        </div>

        <!-- Live Clock & Quick Actions -->
        <div class="flex flex-col sm:flex-row md:flex-col items-start md:items-end justify-between gap-3 shrink-0">
          <div class="p-3 rounded-2xl border text-right backdrop-blur-md"
               :style="{ backgroundColor: 'var(--bg-content)', borderColor: 'var(--border-color)' }">
            <div class="flex items-center gap-2 text-xs font-bold text-blue-600 dark:text-blue-400">
              <Clock class="w-3.5 h-3.5 animate-spin-slow" />
              <span class="font-mono text-sm tracking-wider">{{ liveTime }}</span>
              <span class="text-[10px] text-slate-400">WIB</span>
            </div>
            <div class="text-[11px] font-medium mt-0.5" :style="{ color: 'var(--text-muted)' }">
              {{ liveDate }}
            </div>
          </div>

          <!-- Top Role Action Buttons -->
          <div class="flex items-center gap-2">
            <button
              v-if="isOperationalStaff || activeViewTab === 'employee'"
              @click="openLeaveModal"
              class="px-3 py-1.5 rounded-xl text-xs font-bold bg-blue-600 hover:bg-blue-700 text-white transition-all shadow-sm flex items-center gap-1.5 cursor-pointer"
            >
              <Calendar class="w-3.5 h-3.5" />
              <span>Ajukan Cuti</span>
            </button>
            <button
              v-if="isOperationalStaff || activeViewTab === 'employee'"
              @click="openProfileModal"
              class="px-3 py-1.5 rounded-xl text-xs font-bold border hover:border-slate-400 transition-all flex items-center gap-1.5 cursor-pointer"
              :style="{ backgroundColor: 'var(--bg-content)', borderColor: 'var(--border-color)', color: 'var(--text-secondary)' }"
            >
              <User class="w-3.5 h-3.5" />
              <span>Profil Saya</span>
            </button>
            <button
              @click="refreshCurrentView"
              class="p-2 rounded-xl border hover:bg-slate-100 dark:hover:bg-slate-800 transition-all cursor-pointer"
              :style="{ backgroundColor: 'var(--bg-content)', borderColor: 'var(--border-color)', color: 'var(--text-muted)' }"
              title="Perbarui Data"
            >
              <RefreshCw class="w-3.5 h-3.5" :class="{ 'animate-spin': isRefreshing }" />
            </button>
          </div>
        </div>
      </div>

      <!-- ROLE PREVIEW SWITCHER BAR (TABS) -->
      <!-- Visible for Super Admin, Owner, Manager, or as role switcher pill -->
      <div 
        v-if="isManagerialRole || availableRolesCount > 1" 
        class="mt-6 pt-5 border-t flex flex-wrap items-center gap-2"
        :style="{ borderColor: 'var(--border-color)' }"
      >
        <span class="text-xs font-bold mr-1 flex items-center gap-1" :style="{ color: 'var(--text-muted)' }">
          <Sparkles class="w-3.5 h-3.5 text-amber-500" />
          Pilih Perspektif Dashboard:
        </span>

        <button
          @click="activeViewTab = 'employee'"
          class="px-3.5 py-1.5 rounded-xl text-xs font-bold border transition-all flex items-center gap-2 cursor-pointer"
          :class="activeViewTab === 'employee' ? 'bg-amber-500 text-white border-amber-600 shadow-sm' : 'hover:bg-slate-100 dark:hover:bg-slate-800 text-slate-600 dark:text-slate-300'"
          :style="activeViewTab !== 'employee' ? { backgroundColor: 'var(--bg-content)', borderColor: 'var(--border-color)' } : {}"
        >
          <span>☕</span>
          <span>Portal Karyawan (Barista/Kasir/Gudang)</span>
        </button>

        <button
          @click="activeViewTab = 'finance'"
          class="px-3.5 py-1.5 rounded-xl text-xs font-bold border transition-all flex items-center gap-2 cursor-pointer"
          :class="activeViewTab === 'finance' ? 'bg-purple-600 text-white border-purple-700 shadow-sm' : 'hover:bg-slate-100 dark:hover:bg-slate-800 text-slate-600 dark:text-slate-300'"
          :style="activeViewTab !== 'finance' ? { backgroundColor: 'var(--bg-content)', borderColor: 'var(--border-color)' } : {}"
        >
          <span>💰</span>
          <span>Dashboard Finance</span>
        </button>

        <button
          @click="activeViewTab = 'hris'"
          class="px-3.5 py-1.5 rounded-xl text-xs font-bold border transition-all flex items-center gap-2 cursor-pointer"
          :class="activeViewTab === 'hris' ? 'bg-rose-600 text-white border-rose-700 shadow-sm' : 'hover:bg-slate-100 dark:hover:bg-slate-800 text-slate-600 dark:text-slate-300'"
          :style="activeViewTab !== 'hris' ? { backgroundColor: 'var(--bg-content)', borderColor: 'var(--border-color)' } : {}"
        >
          <span>👥</span>
          <span>Dashboard HRIS</span>
        </button>

        <button
          @click="activeViewTab = 'executive'"
          class="px-3.5 py-1.5 rounded-xl text-xs font-bold border transition-all flex items-center gap-2 cursor-pointer"
          :class="activeViewTab === 'executive' ? 'bg-blue-600 text-white border-blue-700 shadow-sm' : 'hover:bg-slate-100 dark:hover:bg-slate-800 text-slate-600 dark:text-slate-300'"
          :style="activeViewTab !== 'executive' ? { backgroundColor: 'var(--bg-content)', borderColor: 'var(--border-color)' } : {}"
        >
          <span>🏢</span>
          <span>Executive Overview & POS</span>
        </button>
      </div>
    </div>

    <!-- ===================================================================== -->
    <!-- 1. VIEW: EMPLOYEE SELF-SERVICE PORTAL (BARISTA, KASIR, INVENTORY)     -->
    <!-- ===================================================================== -->
    <div v-if="activeViewTab === 'employee'" class="space-y-6">
      
      <!-- HERO CARD: LIVE CLOCK-IN & CLOCK-OUT HERO -->
      <div 
        class="rounded-3xl p-6 sm:p-7 border shadow-sm relative overflow-hidden transition-all"
        :style="{ backgroundColor: 'var(--bg-primary)', borderColor: 'var(--border-color)' }"
      >
        <div class="grid grid-cols-1 lg:grid-cols-12 gap-6 items-center">
          
          <!-- Attendance Status & Live Indicator (7 cols) -->
          <div class="lg:col-span-7 space-y-4">
            <div class="flex items-center gap-2.5">
              <div class="p-2 rounded-xl"
                   :class="attendanceState.hasClockedOut ? 'bg-slate-100 dark:bg-slate-800 text-slate-600' :
                           attendanceState.hasClockedIn ? 'bg-emerald-100 dark:bg-emerald-950/50 text-emerald-600' :
                           'bg-amber-100 dark:bg-amber-950/50 text-amber-600'">
                <CalendarClock class="w-6 h-6" />
              </div>
              <div>
                <div class="text-xs font-bold uppercase tracking-wider text-slate-400">Status Presensi Hari Ini</div>
                <div class="text-lg sm:text-xl font-black flex items-center gap-2" :style="{ color: 'var(--text-primary)' }">
                  <span v-if="attendanceState.hasClockedOut" class="text-slate-600 dark:text-slate-400">
                    Selesai Shift Kerja (Clocked Out)
                  </span>
                  <span v-else-if="attendanceState.hasClockedIn" class="text-emerald-600 dark:text-emerald-400 flex items-center gap-2">
                    <span class="w-2.5 h-2.5 rounded-full bg-emerald-500 animate-ping"></span>
                    Sedang Bertugas (On Duty)
                  </span>
                  <span v-else class="text-amber-600 dark:text-amber-400 flex items-center gap-2">
                    <span class="w-2.5 h-2.5 rounded-full bg-amber-500 animate-pulse"></span>
                    Belum Absen Masuk (Siap Mulai Shift)
                  </span>
                </div>
              </div>
            </div>

            <!-- Shift info pills -->
            <div class="flex flex-wrap items-center gap-3 text-xs" :style="{ color: 'var(--text-secondary)' }">
              <div class="px-3 py-1 rounded-xl border flex items-center gap-2"
                   :style="{ backgroundColor: 'var(--bg-content)', borderColor: 'var(--border-color)' }">
                <Clock class="w-3.5 h-3.5 text-blue-500" />
                <span>Shift Terjadwal: <strong>{{ todayShiftName }}</strong></span>
              </div>
              <div v-if="attendanceState.clockInTime" class="px-3 py-1 rounded-xl border flex items-center gap-2 bg-emerald-50 dark:bg-emerald-950/30 border-emerald-200 text-emerald-800 dark:text-emerald-300">
                <CheckCircle2 class="w-3.5 h-3.5 text-emerald-500" />
                <span>Masuk: <strong>{{ attendanceState.clockInTime }}</strong></span>
              </div>
              <div v-if="attendanceState.clockOutTime" class="px-3 py-1 rounded-xl border flex items-center gap-2 bg-blue-50 dark:bg-blue-950/30 border-blue-200 text-blue-800 dark:text-blue-300">
                <CheckCircle2 class="w-3.5 h-3.5 text-blue-500" />
                <span>Keluar: <strong>{{ attendanceState.clockOutTime }}</strong></span>
              </div>
              <div v-if="attendanceState.workedHours" class="px-3 py-1 rounded-xl border flex items-center gap-2 font-semibold"
                   :style="{ backgroundColor: 'var(--bg-content)', borderColor: 'var(--border-color)' }">
                <span>⏱️ Durasi: {{ attendanceState.workedHours }}</span>
              </div>
            </div>

            <!-- Notes preview if available -->
            <p v-if="attendanceState.notes" class="text-xs italic text-slate-500 bg-slate-50 dark:bg-slate-900/40 p-2.5 rounded-xl border border-slate-200 dark:border-slate-800">
              💬 Catatan Harian: "{{ attendanceState.notes }}"
            </p>
          </div>

          <!-- Action Buttons (5 cols) -->
          <div class="lg:col-span-5 flex flex-col sm:flex-row lg:flex-col gap-3 justify-center">
            <!-- Clock In Button -->
            <button
              v-if="!attendanceState.hasClockedIn"
              @click="handleClockIn"
              :disabled="attendanceActionLoading"
              class="w-full py-4 px-6 rounded-2xl font-black text-sm text-white bg-gradient-to-r from-emerald-600 to-teal-600 hover:from-emerald-700 hover:to-teal-700 active:scale-98 shadow-lg shadow-emerald-600/25 transition-all flex items-center justify-center gap-2 cursor-pointer disabled:opacity-50"
            >
              <span v-if="attendanceActionLoading" class="w-5 h-5 border-2 border-white border-t-transparent rounded-full animate-spin"></span>
              <CheckCircle class="w-5 h-5" v-else />
              <span>CLOCK IN SEKARANG (MASUK SHIFT)</span>
            </button>

            <!-- Clock Out Button -->
            <button
              v-else-if="!attendanceState.hasClockedOut"
              @click="handleClockOut"
              :disabled="attendanceActionLoading"
              class="w-full py-4 px-6 rounded-2xl font-black text-sm text-white bg-gradient-to-r from-rose-600 to-red-600 hover:from-rose-700 hover:to-red-700 active:scale-98 shadow-lg shadow-rose-600/25 transition-all flex items-center justify-center gap-2 cursor-pointer disabled:opacity-50"
            >
              <span v-if="attendanceActionLoading" class="w-5 h-5 border-2 border-white border-t-transparent rounded-full animate-spin"></span>
              <XCircle class="w-5 h-5" v-else />
              <span>CLOCK OUT (SELESAI SHIFT)</span>
            </button>

            <!-- Completed state info banner -->
            <div 
              v-else 
              class="w-full p-4 rounded-2xl border text-center text-xs font-semibold bg-slate-50 dark:bg-slate-900/40 border-slate-200 dark:border-slate-800 text-slate-600 dark:text-slate-400"
            >
              🎉 Presensi shift Anda hari ini telah lengkap. Terima kasih atas kerja keras Anda!
            </div>

            <!-- Optional Quick Note Input -->
            <div v-if="!attendanceState.hasClockedOut" class="relative">
              <input
                v-model="attendanceNoteInput"
                type="text"
                placeholder="Tambah catatan shift (opsional)..."
                class="w-full text-xs px-3.5 py-2 rounded-xl border outline-none transition-all focus:border-blue-500"
                :style="{ backgroundColor: 'var(--bg-content)', borderColor: 'var(--border-color)', color: 'var(--text-primary)' }"
              />
            </div>
          </div>

        </div>
      </div>

      <!-- 3 CORE SELF-SERVICE KPI CARDS -->
      <div class="grid grid-cols-1 md:grid-cols-3 gap-5">
        
        <!-- Card 1: Leave Quota (Jumlah Cuti) -->
        <div 
          class="rounded-3xl p-6 border shadow-sm flex flex-col justify-between hover:shadow-md transition-all relative overflow-hidden"
          :style="{ backgroundColor: 'var(--bg-primary)', borderColor: 'var(--border-color)' }"
        >
          <div class="flex items-start justify-between">
            <div>
              <div class="text-xs font-bold uppercase tracking-wider text-slate-400 mb-1">Kuota Cuti Tahunan</div>
              <div class="text-3xl font-black tracking-tight text-blue-600 dark:text-blue-400">
                {{ portalData.leave_summary?.remaining_leave ?? 12 }} <span class="text-sm font-semibold text-slate-400">Hari</span>
              </div>
              <div class="text-xs text-slate-500 mt-1">
                Dari total alokasi 12 hari/tahun
              </div>
            </div>
            <div class="w-12 h-12 rounded-2xl bg-blue-50 dark:bg-blue-950/50 text-blue-600 dark:text-blue-400 flex items-center justify-center shrink-0">
              <Calendar class="w-6 h-6" />
            </div>
          </div>

          <div class="mt-4 pt-4 border-t flex items-center justify-between" :style="{ borderColor: 'var(--border-color)' }">
            <span class="text-xs text-slate-500">
              {{ portalData.leave_summary?.pending_count || 0 }} Pengajuan Menunggu
            </span>
            <button 
              @click="openLeaveModal"
              class="text-xs font-bold text-blue-600 hover:text-blue-700 dark:text-blue-400 flex items-center gap-1 cursor-pointer"
            >
              + Ajukan Cuti →
            </button>
          </div>
        </div>

        <!-- Card 2: Take Home Pay / Gaji Terakhir -->
        <div 
          class="rounded-3xl p-6 border shadow-sm flex flex-col justify-between hover:shadow-md transition-all relative overflow-hidden"
          :style="{ backgroundColor: 'var(--bg-primary)', borderColor: 'var(--border-color)' }"
        >
          <div class="flex items-start justify-between">
            <div>
              <div class="text-xs font-bold uppercase tracking-wider text-slate-400 mb-1">Take Home Pay Terakhir</div>
              <div class="text-2xl sm:text-3xl font-black tracking-tight text-emerald-600 dark:text-emerald-400">
                Rp {{ formatNum(portalData.latest_payroll?.net_salary || portalData.employee?.basic_salary || 5200000) }}
              </div>
              <div class="text-xs text-slate-500 mt-1 flex items-center gap-1">
                <span>Periode: {{ portalData.latest_payroll?.period_label || 'Bulan Berjalan' }}</span>
              </div>
            </div>
            <div class="w-12 h-12 rounded-2xl bg-emerald-50 dark:bg-emerald-950/50 text-emerald-600 dark:text-emerald-400 flex items-center justify-center shrink-0">
              <Receipt class="w-6 h-6" />
            </div>
          </div>

          <div class="mt-4 pt-4 border-t flex items-center justify-between" :style="{ borderColor: 'var(--border-color)' }">
            <span class="text-xs px-2 py-0.5 rounded-full font-semibold"
                  :class="portalData.latest_payroll?.is_paid ? 'bg-emerald-100 text-emerald-700' : 'bg-amber-100 text-amber-700'">
              {{ portalData.latest_payroll?.is_paid ? '✅ Sudah Ditransfer' : '⏳ Diproses HR' }}
            </span>
            <button 
              @click="openSalarySlipModal"
              class="text-xs font-bold text-emerald-600 hover:text-emerald-700 dark:text-emerald-400 flex items-center gap-1 cursor-pointer"
            >
              Rincian Slip Gaji →
            </button>
          </div>
        </div>

        <!-- Card 3: Profil Kepegawaian & Posisi -->
        <div 
          class="rounded-3xl p-6 border shadow-sm flex flex-col justify-between hover:shadow-md transition-all relative overflow-hidden"
          :style="{ backgroundColor: 'var(--bg-primary)', borderColor: 'var(--border-color)' }"
        >
          <div class="flex items-start justify-between">
            <div>
              <div class="text-xs font-bold uppercase tracking-wider text-slate-400 mb-1">Posisi & Departemen</div>
              <div class="text-xl font-black tracking-tight" :style="{ color: 'var(--text-primary)' }">
                {{ portalData.employee?.job_title || 'Head Barista' }}
              </div>
              <div class="text-xs text-slate-500 mt-1">
                {{ portalData.employee?.department || 'Front of House' }} • {{ portalData.employee?.position || 'Staff' }}
              </div>
            </div>
            <div class="w-12 h-12 rounded-2xl bg-amber-50 dark:bg-amber-950/50 text-amber-600 dark:text-amber-400 flex items-center justify-center shrink-0">
              <Briefcase class="w-6 h-6" />
            </div>
          </div>

          <div class="mt-4 pt-4 border-t flex items-center justify-between" :style="{ borderColor: 'var(--border-color)' }">
            <span class="text-xs text-slate-500">
              Sejak: {{ portalData.employee?.join_date || '2023-03-01' }}
            </span>
            <button 
              @click="openProfileModal"
              class="text-xs font-bold text-amber-600 hover:text-amber-700 dark:text-amber-400 flex items-center gap-1 cursor-pointer"
            >
              Lihat Detail Profil →
            </button>
          </div>
        </div>

      </div>

      <!-- JADWAL SHIFT ROSTER MINGGU INI (7 HARI) -->
      <div 
        class="rounded-3xl p-6 sm:p-7 border shadow-sm transition-all"
        :style="{ backgroundColor: 'var(--bg-primary)', borderColor: 'var(--border-color)' }"
      >
        <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3 mb-5">
          <div>
            <h3 class="text-base font-bold tracking-tight" :style="{ color: 'var(--text-primary)' }">
              📅 Jadwal Kerja & Shift Mingguan Anda
            </h3>
            <p class="text-xs text-slate-400">Roster jadwal shift kerja operasional minggu ini</p>
          </div>
          <div class="flex items-center gap-3 text-xs font-medium">
            <span class="flex items-center gap-1.5"><span class="w-2.5 h-2.5 rounded-full bg-blue-600"></span> Shift Pagi (07-15)</span>
            <span class="flex items-center gap-1.5"><span class="w-2.5 h-2.5 rounded-full bg-emerald-600"></span> Shift Siang (15-23)</span>
            <span class="flex items-center gap-1.5"><span class="w-2.5 h-2.5 rounded-full bg-slate-400"></span> Libur (OFF)</span>
          </div>
        </div>

        <!-- 7-Day Grid Cards -->
        <div class="grid grid-cols-2 sm:grid-cols-4 lg:grid-cols-7 gap-3">
          <div 
            v-for="s in weeklySchedules" 
            :key="s.date"
            class="p-3.5 rounded-2xl border text-center transition-all relative overflow-hidden"
            :class="s.is_today ? 'ring-2 ring-blue-500 shadow-md bg-blue-50/50 dark:bg-blue-950/30' : 'hover:border-slate-400'"
            :style="{ backgroundColor: 'var(--bg-content)', borderColor: 'var(--border-color)' }"
          >
            <!-- Badge Hari Ini -->
            <div v-if="s.is_today" class="text-[9px] font-black uppercase tracking-wider text-blue-600 dark:text-blue-400 mb-1">
              ⭐ HARI INI
            </div>
            <div class="text-xs font-bold" :style="{ color: 'var(--text-primary)' }">
              {{ formatDayName(s.day_name) }}
            </div>
            <div class="text-[11px] font-mono text-slate-400 mb-2">
              {{ s.date ? s.date.slice(5) : '' }}
            </div>

            <!-- Shift Pill -->
            <div 
              class="py-1 px-2 rounded-xl text-xs font-bold text-white shadow-xs"
              :style="{ backgroundColor: s.color || '#2563eb' }"
            >
              {{ s.shift_name }}
            </div>

            <div class="text-[10px] font-semibold mt-1.5 text-slate-500">
              {{ s.shift_name === 'Libur' ? 'Off Day' : `${s.start_time} - ${s.end_time}` }}
            </div>
          </div>
        </div>
      </div>

      <!-- RIWAYAT PENGAJUAN CUTI & AKSI CEPAT ROLE -->
      <div class="grid grid-cols-1 lg:grid-cols-12 gap-6">
        
        <!-- Riwayat Cuti Terkini (7 cols) -->
        <div 
          class="lg:col-span-7 rounded-3xl p-6 border shadow-sm flex flex-col justify-between"
          :style="{ backgroundColor: 'var(--bg-primary)', borderColor: 'var(--border-color)' }"
        >
          <div>
            <div class="flex items-center justify-between mb-3">
              <h3 class="text-base font-bold" :style="{ color: 'var(--text-primary)' }">
                Riwayat Pengajuan Cuti Anda
              </h3>
              <button @click="openLeaveModal" class="text-xs text-blue-600 dark:text-blue-400 font-bold hover:underline cursor-pointer">
                + Tambah Pengajuan
              </button>
            </div>
            <p class="text-xs text-slate-400 mb-4">Status persetujuan permohonan izin cuti terakhir</p>

            <div v-if="portalData.leave_summary?.history?.length" class="space-y-2.5">
              <div 
                v-for="l in portalData.leave_summary.history" 
                :key="l.id"
                class="p-3 rounded-2xl border flex items-center justify-between text-xs transition-colors"
                :style="{ backgroundColor: 'var(--bg-content)', borderColor: 'var(--border-color)' }"
              >
                <div>
                  <div class="font-bold" :style="{ color: 'var(--text-primary)' }">
                    {{ l.leave_type }} ({{ l.total_days }} Hari)
                  </div>
                  <div class="text-[11px] text-slate-500 mt-0.5">
                    {{ l.start_date }} s/d {{ l.end_date }} • {{ l.reason }}
                  </div>
                </div>
                <div class="text-right">
                  <span 
                    class="px-2.5 py-0.5 rounded-full text-[10.5px] font-bold uppercase tracking-wider"
                    :class="l.status === 'approved' ? 'bg-emerald-100 text-emerald-700 dark:bg-emerald-950/60 dark:text-emerald-300' :
                            l.status === 'rejected' ? 'bg-rose-100 text-rose-700 dark:bg-rose-950/60 dark:text-rose-300' :
                            'bg-amber-100 text-amber-700 dark:bg-amber-950/60 dark:text-amber-300'"
                  >
                    {{ l.status }}
                  </span>
                  <div class="text-[10px] text-slate-400 mt-0.5">{{ l.created_at }}</div>
                </div>
              </div>
            </div>

            <div v-else class="p-6 text-center text-xs text-slate-400 border border-dashed rounded-2xl" :style="{ borderColor: 'var(--border-color)' }">
              Belum ada riwayat pengajuan cuti yang tercatat.
            </div>
          </div>
        </div>

        <!-- Quick Operational Module Links by Role (5 cols) -->
        <div 
          class="lg:col-span-5 rounded-3xl p-6 border shadow-sm flex flex-col justify-between"
          :style="{ backgroundColor: 'var(--bg-primary)', borderColor: 'var(--border-color)' }"
        >
          <div>
            <h3 class="text-base font-bold mb-1" :style="{ color: 'var(--text-primary)' }">
              ⚡ Pintasan Menu Cepat Anda
            </h3>
            <p class="text-xs text-slate-400 mb-4">Navigasi langsung ke modul tugas harian</p>

            <div class="space-y-2.5">
              <!-- If Barista -->
              <router-link 
                to="/kitchen" 
                class="p-3.5 rounded-2xl border flex items-center justify-between hover:border-blue-500 transition-all group"
                :style="{ backgroundColor: 'var(--bg-content)', borderColor: 'var(--border-color)' }"
              >
                <div class="flex items-center gap-3">
                  <div class="w-8 h-8 rounded-xl bg-amber-100 dark:bg-amber-950/50 text-amber-600 flex items-center justify-center font-bold">
                    ☕
                  </div>
                  <div>
                    <div class="text-xs font-bold" :style="{ color: 'var(--text-primary)' }">Kitchen Display (KDS)</div>
                    <div class="text-[11px] text-slate-400">Antrian pesanan kopi & minuman</div>
                  </div>
                </div>
                <ArrowRight class="w-4 h-4 text-slate-400 group-hover:translate-x-1 transition-transform" />
              </router-link>

              <!-- If Kasir -->
              <router-link 
                to="/pos" 
                class="p-3.5 rounded-2xl border flex items-center justify-between hover:border-blue-500 transition-all group"
                :style="{ backgroundColor: 'var(--bg-content)', borderColor: 'var(--border-color)' }"
              >
                <div class="flex items-center gap-3">
                  <div class="w-8 h-8 rounded-xl bg-emerald-100 dark:bg-emerald-950/50 text-emerald-600 flex items-center justify-center font-bold">
                    🏷️
                  </div>
                  <div>
                    <div class="text-xs font-bold" :style="{ color: 'var(--text-primary)' }">Buka Kasir POS</div>
                    <div class="text-[11px] text-slate-400">Transaksi penjualan & pembayaran</div>
                  </div>
                </div>
                <ArrowRight class="w-4 h-4 text-slate-400 group-hover:translate-x-1 transition-transform" />
              </router-link>

              <!-- If Inventory / Gudang -->
              <router-link 
                to="/inventory" 
                class="p-3.5 rounded-2xl border flex items-center justify-between hover:border-blue-500 transition-all group"
                :style="{ backgroundColor: 'var(--bg-content)', borderColor: 'var(--border-color)' }"
              >
                <div class="flex items-center gap-3">
                  <div class="w-8 h-8 rounded-xl bg-cyan-100 dark:bg-cyan-950/50 text-cyan-600 flex items-center justify-center font-bold">
                    📦
                  </div>
                  <div>
                    <div class="text-xs font-bold" :style="{ color: 'var(--text-primary)' }">Manajemen Bahan & Stok</div>
                    <div class="text-[11px] text-slate-400">Penerimaan barang & stock opname</div>
                  </div>
                </div>
                <ArrowRight class="w-4 h-4 text-slate-400 group-hover:translate-x-1 transition-transform" />
              </router-link>
            </div>
          </div>

          <div class="mt-4 pt-3 border-t text-[11px] text-slate-400 flex items-center justify-between" :style="{ borderColor: 'var(--border-color)' }">
            <span>Perlu bantuan shift atau kendala pos?</span>
            <span class="font-bold text-blue-600 dark:text-blue-400">Hubungi Manager</span>
          </div>
        </div>

      </div>

    </div>

    <!-- ===================================================================== -->
    <!-- 2. VIEW: FINANCE DASHBOARD (AKUNTAN / KEUANGAN)                       -->
    <!-- ===================================================================== -->
    <div v-else-if="activeViewTab === 'finance'" class="space-y-6">
      
      <!-- 4 Finance KPI Cards -->
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-5">
        <div class="rounded-3xl p-5 border shadow-sm flex items-start justify-between"
             :style="{ backgroundColor: 'var(--bg-primary)', borderColor: 'var(--border-color)' }">
          <div>
            <div class="text-xs font-semibold text-slate-400 mb-1">Total Pendapatan (Revenue)</div>
            <div class="text-2xl font-black tracking-tight text-emerald-600 dark:text-emerald-400">
              Rp {{ formatNum(financeOverview.total_revenue || 128450000) }}
            </div>
            <div class="text-[11px] text-emerald-600 mt-1 font-semibold flex items-center gap-1">
              <TrendingUp class="w-3.5 h-3.5" /> +14.2% MoM
            </div>
          </div>
          <div class="w-11 h-11 rounded-2xl bg-emerald-50 dark:bg-emerald-950/50 text-emerald-600 flex items-center justify-center">
            <DollarSign class="w-6 h-6" />
          </div>
        </div>

        <div class="rounded-3xl p-5 border shadow-sm flex items-start justify-between"
             :style="{ backgroundColor: 'var(--bg-primary)', borderColor: 'var(--border-color)' }">
          <div>
            <div class="text-xs font-semibold text-slate-400 mb-1">Beban & Pengeluaran</div>
            <div class="text-2xl font-black tracking-tight text-rose-600 dark:text-rose-400">
              Rp {{ formatNum(financeOverview.total_expenses || 78320000) }}
            </div>
            <div class="text-[11px] text-slate-400 mt-1">Bahan baku & operasional</div>
          </div>
          <div class="w-11 h-11 rounded-2xl bg-rose-50 dark:bg-rose-950/50 text-rose-600 flex items-center justify-center">
            <CreditCard class="w-6 h-6" />
          </div>
        </div>

        <div class="rounded-3xl p-5 border shadow-sm flex items-start justify-between"
             :style="{ backgroundColor: 'var(--bg-primary)', borderColor: 'var(--border-color)' }">
          <div>
            <div class="text-xs font-semibold text-slate-400 mb-1">Laba Bersih Operasional</div>
            <div class="text-2xl font-black tracking-tight text-blue-600 dark:text-blue-400">
              Rp {{ formatNum(financeOverview.net_profit || 50130000) }}
            </div>
            <div class="text-[11px] text-blue-600 mt-1 font-semibold">Margin Laba: 39.0%</div>
          </div>
          <div class="w-11 h-11 rounded-2xl bg-blue-50 dark:bg-blue-950/50 text-blue-600 flex items-center justify-center">
            <TrendingUp class="w-6 h-6" />
          </div>
        </div>

        <div class="rounded-3xl p-5 border shadow-sm flex items-start justify-between"
             :style="{ backgroundColor: 'var(--bg-primary)', borderColor: 'var(--border-color)' }">
          <div>
            <div class="text-xs font-semibold text-slate-400 mb-1">Kas & Rekening Bank</div>
            <div class="text-2xl font-black tracking-tight text-purple-600 dark:text-purple-400">
              Rp {{ formatNum((financeOverview.cash_on_hand || 32500000) + 51650000) }}
            </div>
            <div class="text-[11px] text-purple-600 mt-1 font-semibold">Cash Drawer + Bank BCA</div>
          </div>
          <div class="w-11 h-11 rounded-2xl bg-purple-50 dark:bg-purple-950/50 text-purple-600 flex items-center justify-center">
            <Landmark class="w-6 h-6" />
          </div>
        </div>
      </div>

      <!-- Detail Finance: Rekening Kas/Bank & Jurnal Terbaru -->
      <div class="grid grid-cols-1 lg:grid-cols-12 gap-6">
        
        <!-- Jurnal Transaksi Terkini (7 cols) -->
        <div 
          class="lg:col-span-7 rounded-3xl p-6 border shadow-sm"
          :style="{ backgroundColor: 'var(--bg-primary)', borderColor: 'var(--border-color)' }"
        >
          <div class="flex items-center justify-between mb-4">
            <div>
              <h3 class="text-base font-bold" :style="{ color: 'var(--text-primary)' }">
                Jurnal Akuntansi Terbaru
              </h3>
              <p class="text-xs text-slate-400">Buku jurnal transaksi otomatis dan manual</p>
            </div>
            <router-link to="/finance" class="text-xs text-blue-600 font-bold hover:underline">
              Buka Semua Jurnal →
            </router-link>
          </div>

          <div class="overflow-x-auto">
            <table class="w-full text-xs">
              <thead>
                <tr class="text-left text-slate-400 border-b" :style="{ borderColor: 'var(--border-color)' }">
                  <th class="pb-2">Tanggal / Ref</th>
                  <th class="pb-2">Keterangan</th>
                  <th class="pb-2 text-right">Debit</th>
                  <th class="pb-2 text-right">Kredit</th>
                  <th class="pb-2 text-right">Status</th>
                </tr>
              </thead>
              <tbody class="divide-y" :style="{ borderColor: 'var(--border-color)' }">
                <tr v-for="(j, idx) in (financeOverview.recent_journals || dummyJournals)" :key="idx" class="hover:bg-slate-50/50">
                  <td class="py-2.5 font-mono text-[11px]">
                    <div class="font-bold" :style="{ color: 'var(--text-primary)' }">{{ j.ref }}</div>
                    <div class="text-slate-400 text-[10px]">{{ j.date }}</div>
                  </td>
                  <td class="py-2.5 max-w-xs truncate" :title="j.description" :style="{ color: 'var(--text-secondary)' }">
                    {{ j.description }}
                  </td>
                  <td class="py-2.5 text-right font-mono font-bold text-slate-700 dark:text-slate-300">
                    Rp {{ formatNum(j.debit) }}
                  </td>
                  <td class="py-2.5 text-right font-mono font-bold text-slate-700 dark:text-slate-300">
                    Rp {{ formatNum(j.credit) }}
                  </td>
                  <td class="py-2.5 text-right">
                    <span class="px-2 py-0.5 rounded-full text-[10px] font-bold bg-emerald-100 text-emerald-800">
                      {{ j.status || 'posted' }}
                    </span>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        <!-- Rekening Bank & Pintasan Modul (5 cols) -->
        <div 
          class="lg:col-span-5 rounded-3xl p-6 border shadow-sm space-y-5"
          :style="{ backgroundColor: 'var(--bg-primary)', borderColor: 'var(--border-color)' }"
        >
          <div>
            <h3 class="text-base font-bold mb-1" :style="{ color: 'var(--text-primary)' }">
              Rekening & Posisi Saldo
            </h3>
            <p class="text-xs text-slate-400 mb-3">Saldo likuiditas tunai dan perbankan</p>

            <div class="space-y-2.5">
              <div class="p-3.5 rounded-2xl border flex items-center justify-between"
                   :style="{ backgroundColor: 'var(--bg-content)', borderColor: 'var(--border-color)' }">
                <div class="flex items-center gap-3">
                  <div class="w-8 h-8 rounded-xl bg-purple-100 text-purple-600 flex items-center justify-center font-bold">
                    🏦
                  </div>
                  <div>
                    <div class="text-xs font-bold" :style="{ color: 'var(--text-primary)' }">Bank BCA Operasional (1102)</div>
                    <div class="text-[11px] text-slate-400">Rekonsiliasi matched 9/11</div>
                  </div>
                </div>
                <div class="text-right font-mono font-black text-sm text-purple-600">
                  Rp 51.650.000
                </div>
              </div>

              <div class="p-3.5 rounded-2xl border flex items-center justify-between"
                   :style="{ backgroundColor: 'var(--bg-content)', borderColor: 'var(--border-color)' }">
                <div class="flex items-center gap-3">
                  <div class="w-8 h-8 rounded-xl bg-emerald-100 text-emerald-600 flex items-center justify-center font-bold">
                    💵
                  </div>
                  <div>
                    <div class="text-xs font-bold" :style="{ color: 'var(--text-primary)' }">Kas Kasir (Cash Drawer) (1101)</div>
                    <div class="text-[11px] text-slate-400">Setoran kas harian POS</div>
                  </div>
                </div>
                <div class="text-right font-mono font-black text-sm text-emerald-600">
                  Rp {{ formatNum(financeOverview.cash_on_hand || 32500000) }}
                </div>
              </div>
            </div>
          </div>

          <!-- Quick Navigation to Finance Modules -->
          <div class="pt-4 border-t" :style="{ borderColor: 'var(--border-color)' }">
            <div class="text-xs font-bold mb-2.5" :style="{ color: 'var(--text-muted)' }">Pintasan Fitur Keuangan:</div>
            <div class="grid grid-cols-2 gap-2 text-xs">
              <router-link to="/finance" class="p-2.5 rounded-xl border text-center font-bold hover:border-blue-500 transition-colors"
                           :style="{ backgroundColor: 'var(--bg-content)', borderColor: 'var(--border-color)', color: 'var(--text-primary)' }">
                📒 Entri Jurnal
              </router-link>
              <router-link to="/finance" class="p-2.5 rounded-xl border text-center font-bold hover:border-blue-500 transition-colors"
                           :style="{ backgroundColor: 'var(--bg-content)', borderColor: 'var(--border-color)', color: 'var(--text-primary)' }">
                💳 Rekonsiliasi Bank
              </router-link>
              <router-link to="/reports" class="p-2.5 rounded-xl border text-center font-bold hover:border-blue-500 transition-colors"
                           :style="{ backgroundColor: 'var(--bg-content)', borderColor: 'var(--border-color)', color: 'var(--text-primary)' }">
                📊 Laporan Keuangan
              </router-link>
              <router-link to="/finance" class="p-2.5 rounded-xl border text-center font-bold hover:border-blue-500 transition-colors"
                           :style="{ backgroundColor: 'var(--bg-content)', borderColor: 'var(--border-color)', color: 'var(--text-primary)' }">
                🏷️ Master Akun (COA)
              </router-link>
            </div>
          </div>
        </div>

      </div>

    </div>

    <!-- ===================================================================== -->
    <!-- 3. VIEW: HRIS DASHBOARD (HR ADMIN / PEOPLE OPERATIONS)               -->
    <!-- ===================================================================== -->
    <div v-else-if="activeViewTab === 'hris'" class="space-y-6">
      
      <!-- 4 HRIS KPI Cards -->
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-5">
        <div class="rounded-3xl p-5 border shadow-sm flex items-start justify-between"
             :style="{ backgroundColor: 'var(--bg-primary)', borderColor: 'var(--border-color)' }">
          <div>
            <div class="text-xs font-semibold text-slate-400 mb-1">Total Karyawan Aktif</div>
            <div class="text-2xl font-black tracking-tight" :style="{ color: 'var(--text-primary)' }">
              {{ hrisTodayStats.total_employees || 7 }} Staf
            </div>
            <div class="text-[11px] text-emerald-600 mt-1 font-semibold">100% Kontrak Aktif</div>
          </div>
          <div class="w-11 h-11 rounded-2xl bg-blue-50 dark:bg-blue-950/50 text-blue-600 flex items-center justify-center">
            <Users class="w-6 h-6" />
          </div>
        </div>

        <div class="rounded-3xl p-5 border shadow-sm flex items-start justify-between"
             :style="{ backgroundColor: 'var(--bg-primary)', borderColor: 'var(--border-color)' }">
          <div>
            <div class="text-xs font-semibold text-slate-400 mb-1">Presensi Hari Ini</div>
            <div class="text-2xl font-black tracking-tight text-emerald-600">
              {{ hrisTodayStats.total_present || 3 }} Hadir
            </div>
            <div class="text-[11px] text-slate-400 mt-1">
              {{ hrisTodayStats.on_time || 3 }} Tepat Waktu • {{ hrisTodayStats.late || 0 }} Terlambat
            </div>
          </div>
          <div class="w-11 h-11 rounded-2xl bg-emerald-50 dark:bg-emerald-950/50 text-emerald-600 flex items-center justify-center">
            <UserCheck class="w-6 h-6" />
          </div>
        </div>

        <div class="rounded-3xl p-5 border shadow-sm flex items-start justify-between"
             :style="{ backgroundColor: 'var(--bg-primary)', borderColor: 'var(--border-color)' }">
          <div>
            <div class="text-xs font-semibold text-slate-400 mb-1">Pending Approval Cuti</div>
            <div class="text-2xl font-black tracking-tight text-amber-600">
              {{ pendingLeavesList.length }} Pengajuan
            </div>
            <div class="text-[11px] text-amber-600 mt-1 font-semibold">Perlu Verifikasi HR</div>
          </div>
          <div class="w-11 h-11 rounded-2xl bg-amber-50 dark:bg-amber-950/50 text-amber-600 flex items-center justify-center">
            <CalendarClock class="w-6 h-6" />
          </div>
        </div>

        <div class="rounded-3xl p-5 border shadow-sm flex items-start justify-between"
             :style="{ backgroundColor: 'var(--bg-primary)', borderColor: 'var(--border-color)' }">
          <div>
            <div class="text-xs font-semibold text-slate-400 mb-1">Estimasi Payroll Periode Ini</div>
            <div class="text-2xl font-black tracking-tight text-rose-600">
              Rp 44.500.000
            </div>
            <div class="text-[11px] text-slate-400 mt-1">Batch penggajian 7 staf</div>
          </div>
          <div class="w-11 h-11 rounded-2xl bg-rose-50 dark:bg-rose-950/50 text-rose-600 flex items-center justify-center">
            <Receipt class="w-6 h-6" />
          </div>
        </div>
      </div>

      <!-- Presensi Live & Approval Cuti Staf -->
      <div class="grid grid-cols-1 lg:grid-cols-12 gap-6">
        
        <!-- Live Attendance Staff Hari Ini (7 cols) -->
        <div 
          class="lg:col-span-7 rounded-3xl p-6 border shadow-sm"
          :style="{ backgroundColor: 'var(--bg-primary)', borderColor: 'var(--border-color)' }"
        >
          <div class="flex items-center justify-between mb-4">
            <div>
              <h3 class="text-base font-bold" :style="{ color: 'var(--text-primary)' }">
                Presensi Staf Live Hari Ini
              </h3>
              <p class="text-xs text-slate-400">Monitoring absensi real-time seluruh staf kafe</p>
            </div>
            <router-link to="/hris" class="text-xs text-blue-600 font-bold hover:underline">
              Kelola Staf HR →
            </router-link>
          </div>

          <div class="space-y-2.5">
            <div 
              v-for="emp in staffAttendanceList" 
              :key="emp.id"
              class="p-3.5 rounded-2xl border flex items-center justify-between text-xs"
              :style="{ backgroundColor: 'var(--bg-content)', borderColor: 'var(--border-color)' }"
            >
              <div class="flex items-center gap-3">
                <div class="w-9 h-9 rounded-xl bg-blue-100 text-blue-700 font-bold flex items-center justify-center text-xs">
                  {{ emp.name[0] }}
                </div>
                <div>
                  <div class="font-bold text-sm" :style="{ color: 'var(--text-primary)' }">{{ emp.name }}</div>
                  <div class="text-[11px] text-slate-400">{{ emp.role }} • Shift {{ emp.shift }}</div>
                </div>
              </div>
              <div class="text-right">
                <span 
                  class="px-2.5 py-0.5 rounded-full text-[10.5px] font-bold"
                  :class="emp.status === 'Hadir' ? 'bg-emerald-100 text-emerald-700' :
                          emp.status === 'Terlambat' ? 'bg-amber-100 text-amber-700' :
                          'bg-slate-100 text-slate-500'"
                >
                  {{ emp.status }}
                </span>
                <div class="text-[10px] text-slate-400 mt-0.5">Clock In: {{ emp.clockIn }}</div>
              </div>
            </div>
          </div>
        </div>

        <!-- Pending Leave Approvals (5 cols) -->
        <div 
          class="lg:col-span-5 rounded-3xl p-6 border shadow-sm flex flex-col justify-between"
          :style="{ backgroundColor: 'var(--bg-primary)', borderColor: 'var(--border-color)' }"
        >
          <div>
            <div class="flex items-center justify-between mb-3">
              <h3 class="text-base font-bold" :style="{ color: 'var(--text-primary)' }">
                Daftar Permohonan Cuti
              </h3>
              <span class="text-xs font-bold text-amber-600">{{ pendingLeavesList.length }} Pending</span>
            </div>
            <p class="text-xs text-slate-400 mb-4">Verifikasi dan persetujuan pengajuan cuti staf</p>

            <div v-if="pendingLeavesList.length" class="space-y-3">
              <div 
                v-for="pl in pendingLeavesList" 
                :key="pl.id"
                class="p-3.5 rounded-2xl border text-xs space-y-2"
                :style="{ backgroundColor: 'var(--bg-content)', borderColor: 'var(--border-color)' }"
              >
                <div class="flex items-center justify-between font-bold">
                  <span :style="{ color: 'var(--text-primary)' }">{{ pl.name }}</span>
                  <span class="text-blue-600">{{ pl.leave_type }}</span>
                </div>
                <div class="text-[11px] text-slate-500">
                  {{ pl.start_date }} s/d {{ pl.end_date }} ({{ pl.total_days }} Hari)
                </div>
                <div class="text-[11px] italic text-slate-400">
                  "{{ pl.reason }}"
                </div>
                <div class="pt-2 border-t flex items-center justify-end gap-2" :style="{ borderColor: 'var(--border-color)' }">
                  <button 
                    @click="approveLeave(pl.id)"
                    class="px-3 py-1 rounded-xl text-xs font-bold bg-emerald-600 hover:bg-emerald-700 text-white cursor-pointer"
                  >
                    Setujui
                  </button>
                  <button 
                    @click="rejectLeave(pl.id)"
                    class="px-3 py-1 rounded-xl text-xs font-bold bg-rose-600 hover:bg-rose-700 text-white cursor-pointer"
                  >
                    Tolak
                  </button>
                </div>
              </div>
            </div>

            <div v-else class="p-8 text-center text-xs text-slate-400 border border-dashed rounded-2xl" :style="{ borderColor: 'var(--border-color)' }">
              ✨ Tidak ada permohonan cuti yang menunggu persetujuan saat ini.
            </div>
          </div>

          <div class="mt-4 pt-3 border-t text-right" :style="{ borderColor: 'var(--border-color)' }">
            <router-link to="/hris" class="text-xs text-blue-600 font-bold hover:underline">
              Buka Modul HRIS Lengkap →
            </router-link>
          </div>
        </div>

      </div>

    </div>

    <!-- ===================================================================== -->
    <!-- 4. VIEW: EXECUTIVE OVERVIEW & POS (SUPER ADMIN / OWNER / MANAGER)     -->
    <!-- ===================================================================== -->
    <div v-else class="space-y-6">
      
      <!-- Stat Cards 4 Kolom -->
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-5">
        <!-- Card 1: Revenue -->
        <div class="card-hover-lift rounded-2xl p-5 border shadow-sm flex items-start justify-between cursor-pointer"
             :style="{ backgroundColor: 'var(--bg-primary)', borderColor: 'var(--border-color)' }">
          <div>
            <div class="text-xs font-semibold text-slate-500 mb-1">1. Total Revenue Hari Ini</div>
            <div class="text-2xl font-black tracking-tight" :style="{ color: 'var(--text-primary)' }">
              Rp {{ formatNum(stats.today_sales || 14850000) }}
            </div>
            <div class="inline-flex items-center gap-1 text-xs font-semibold text-emerald-600 mt-2 bg-emerald-50 px-2 py-0.5 rounded-full">
              <TrendingUp class="w-3.5 h-3.5" /> 12.5% vs kemarin
            </div>
          </div>
          <div class="w-11 h-11 rounded-xl bg-blue-50 text-blue-600 flex items-center justify-center">
            <DollarSign class="w-6 h-6" />
          </div>
        </div>

        <!-- Card 2: Orders Today -->
        <div class="card-hover-lift rounded-2xl p-5 border shadow-sm flex items-start justify-between cursor-pointer"
             :style="{ backgroundColor: 'var(--bg-primary)', borderColor: 'var(--border-color)' }">
          <div>
            <div class="text-xs font-semibold text-slate-500 mb-1">2. Pesanan Masuk (POS)</div>
            <div class="text-2xl font-black tracking-tight" :style="{ color: 'var(--text-primary)' }">
              {{ stats.total_orders || 48 }}
            </div>
            <div class="inline-flex items-center gap-1 text-xs font-semibold text-emerald-600 mt-2 bg-emerald-50 px-2 py-0.5 rounded-full">
              <TrendingUp class="w-3.5 h-3.5" /> 8.2%
            </div>
          </div>
          <div class="w-11 h-11 rounded-xl bg-blue-50 text-blue-600 flex items-center justify-center">
            <ClipboardList class="w-6 h-6" />
          </div>
        </div>

        <!-- Card 3: Low Stock Warnings -->
        <div class="card-hover-lift rounded-2xl p-5 border shadow-sm flex items-start justify-between cursor-pointer"
             :style="{ backgroundColor: 'var(--bg-primary)', borderColor: 'var(--border-color)' }">
          <div>
            <div class="text-xs font-semibold text-slate-500 mb-1">3. Stok Bahan Menipis</div>
            <div class="text-2xl font-black tracking-tight text-amber-600">
              {{ stats.low_stock_count || 3 }} items
            </div>
            <div class="inline-flex items-center gap-1 text-xs font-semibold text-amber-600 mt-2 bg-amber-50 px-2 py-0.5 rounded-full">
              <AlertTriangle class="w-3.5 h-3.5" /> Perlu Restock PO
            </div>
          </div>
          <div class="w-11 h-11 rounded-xl bg-amber-50 text-amber-600 flex items-center justify-center">
            <Package class="w-6 h-6" />
          </div>
        </div>

        <!-- Card 4: Active Occupied Tables -->
        <div class="card-hover-lift rounded-2xl p-5 border shadow-sm flex items-start justify-between cursor-pointer"
             :style="{ backgroundColor: 'var(--bg-primary)', borderColor: 'var(--border-color)' }">
          <div>
            <div class="text-xs font-semibold text-slate-500 mb-1">4. Meja Terisi Aktif</div>
            <div class="text-2xl font-black tracking-tight" :style="{ color: 'var(--text-primary)' }">
              {{ stats.active_tables || 4 }} Meja
            </div>
            <div class="inline-flex items-center gap-1 text-xs font-semibold text-emerald-600 mt-2 bg-emerald-50 px-2 py-0.5 rounded-full">
              <TrendingUp class="w-3.5 h-3.5" /> Dinamis POS
            </div>
          </div>
          <div class="w-11 h-11 rounded-xl bg-blue-50 text-blue-600 flex items-center justify-center">
            <Armchair class="w-6 h-6" />
          </div>
        </div>
      </div>

      <!-- Middle Section: Revenue Trend Curve & Table Layout -->
      <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <div class="lg:col-span-2 rounded-2xl p-6 border shadow-sm"
             :style="{ backgroundColor: 'var(--bg-primary)', borderColor: 'var(--border-color)' }">
          <div class="flex items-center justify-between mb-4">
            <div>
              <h3 class="font-bold text-base" :style="{ color: 'var(--text-primary)' }">Trend Omset Mingguan</h3>
              <p class="text-xs text-slate-400">Pola pergerakan omset kasir seluruh outlet</p>
            </div>
            <div class="flex items-center gap-2 text-xs font-medium">
              <span class="inline-flex items-center gap-1 text-blue-600">
                <span class="w-2.5 h-2.5 rounded-full bg-blue-600"></span> Minggu Ini
              </span>
              <span class="inline-flex items-center gap-1 text-slate-400">
                <span class="w-2.5 h-2.5 rounded-full bg-slate-300"></span> Minggu Lalu
              </span>
            </div>
          </div>

          <div class="h-64 flex flex-col justify-between">
            <div class="relative h-52 w-full border-b border-slate-100">
              <svg class="w-full h-full" viewBox="0 0 600 200" preserveAspectRatio="none">
                <defs>
                  <linearGradient id="gradientArea" x1="0%" y1="0%" x2="0%" y2="100%">
                    <stop offset="0%" stop-color="#2563EB" stop-opacity="0.25" />
                    <stop offset="100%" stop-color="#2563EB" stop-opacity="0.0" />
                  </linearGradient>
                </defs>
                <path d="M 0 160 Q 100 120, 200 130 T 400 70 T 600 40 L 600 200 L 0 200 Z" fill="url(#gradientArea)" />
                <path d="M 0 160 Q 100 120, 200 130 T 400 70 T 600 40" fill="none" stroke="#2563eb" stroke-width="3" />
              </svg>
            </div>
            <div class="flex justify-between text-xs text-slate-400 pt-2 font-medium">
              <span>Sen</span><span>Sel</span><span>Rab</span><span>Kam</span><span>Jum</span><span>Sab</span><span>Min</span>
            </div>
          </div>
        </div>

        <div class="rounded-2xl p-6 border shadow-sm flex flex-col justify-between"
             :style="{ backgroundColor: 'var(--bg-primary)', borderColor: 'var(--border-color)' }">
          <div>
            <div class="flex items-center justify-between mb-2">
              <h3 class="font-bold text-base" :style="{ color: 'var(--text-primary)' }">Denah Meja Outlet</h3>
              <router-link to="/pos/tables" class="text-xs text-blue-600 font-semibold hover:underline">
                Buka Denah →
              </router-link>
            </div>
            <p class="text-xs text-slate-400">Status keterisian meja makan realtime</p>
          </div>

          <div class="rounded-xl p-4 border my-3 grid grid-cols-4 gap-2.5"
               :style="{ backgroundColor: 'var(--bg-content)', borderColor: 'var(--border-color)' }">
            <div 
              v-for="t in (floorTables.length ? floorTables.slice(0, 8) : defaultTables)" 
              :key="t.id"
              class="p-2.5 rounded-lg border text-center text-xs font-bold transition-all shadow-xs"
              :class="t.status === 'occupied' 
                ? 'bg-red-500 text-white border-red-600' 
                : t.status === 'reserved' 
                  ? 'bg-amber-400 text-slate-900 border-amber-500' 
                  : 'bg-white border-emerald-500 text-emerald-700'"
            >
              <div>{{ t.table_number }}</div>
              <div class="text-[10px] font-normal" :class="t.status === 'occupied' ? 'text-red-100' : 'text-slate-400'">
                {{ t.status === 'occupied' ? 'Terisi' : `${t.capacity} Seat` }}
              </div>
            </div>
          </div>

          <div class="flex justify-around text-[11px] text-slate-500 pt-2 border-t" :style="{ borderColor: 'var(--border-color)' }">
            <div class="flex items-center gap-1.5"><span class="w-2.5 h-2.5 rounded bg-emerald-500"></span> Tersedia</div>
            <div class="flex items-center gap-1.5"><span class="w-2.5 h-2.5 rounded bg-red-500"></span> Terisi</div>
            <div class="flex items-center gap-1.5"><span class="w-2.5 h-2.5 rounded bg-amber-400"></span> Reservasi</div>
          </div>
        </div>
      </div>

    </div>

    <!-- ===================================================================== -->
    <!-- MODAL 1: AJUKAN CUTI KARYAWAN                                         -->
    <!-- ===================================================================== -->
    <Teleport to="body">
      <div v-if="showLeaveModal" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-950/60 backdrop-blur-xs">
        <div class="w-full max-w-md rounded-3xl p-6 sm:p-7 border shadow-2xl animate-in zoom-in-95 duration-200"
             :style="{ backgroundColor: 'var(--bg-primary)', borderColor: 'var(--border-color)', color: 'var(--text-primary)' }">
          
          <div class="flex items-center justify-between pb-3 border-b" :style="{ borderColor: 'var(--border-color)' }">
            <div class="flex items-center gap-2.5">
              <div class="w-9 h-9 rounded-xl bg-blue-100 dark:bg-blue-950/50 text-blue-600 flex items-center justify-center">
                <Calendar class="w-5 h-5" />
              </div>
              <div>
                <h3 class="font-bold text-base">Formulir Pengajuan Cuti</h3>
                <p class="text-xs text-slate-400">Sisa kuota: {{ portalData.leave_summary?.remaining_leave ?? 10 }} hari</p>
              </div>
            </div>
            <button @click="showLeaveModal = false" class="p-1.5 rounded-lg text-slate-400 hover:text-slate-700 cursor-pointer">✕</button>
          </div>

          <form @submit.prevent="submitLeaveForm" class="py-4 space-y-3.5 text-xs">
            <div>
              <label class="block font-bold mb-1" :style="{ color: 'var(--text-secondary)' }">Jenis Cuti</label>
              <select v-model="leaveForm.leave_type" class="w-full px-3 py-2.5 rounded-xl border outline-none"
                      :style="{ backgroundColor: 'var(--input-bg)', borderColor: 'var(--input-border)', color: 'var(--input-text)' }">
                <option value="Cuti Tahunan">Cuti Tahunan (Tahunan Reguler)</option>
                <option value="Cuti Sakit">Cuti Sakit (Surat Dokter)</option>
                <option value="Cuti Alasan Penting">Cuti Alasan Penting (Keluarga)</option>
              </select>
            </div>

            <div class="grid grid-cols-2 gap-3">
              <div>
                <label class="block font-bold mb-1" :style="{ color: 'var(--text-secondary)' }">Tanggal Mulai</label>
                <input v-model="leaveForm.start_date" type="date" required class="w-full px-3 py-2 rounded-xl border outline-none"
                       :style="{ backgroundColor: 'var(--input-bg)', borderColor: 'var(--input-border)', color: 'var(--input-text)' }" />
              </div>
              <div>
                <label class="block font-bold mb-1" :style="{ color: 'var(--text-secondary)' }">Tanggal Selesai</label>
                <input v-model="leaveForm.end_date" type="date" required class="w-full px-3 py-2 rounded-xl border outline-none"
                       :style="{ backgroundColor: 'var(--input-bg)', borderColor: 'var(--input-border)', color: 'var(--input-text)' }" />
              </div>
            </div>

            <div>
              <label class="block font-bold mb-1" :style="{ color: 'var(--text-secondary)' }">Alasan Pengajuan</label>
              <textarea v-model="leaveForm.reason" rows="3" required placeholder="Tuliskan keterangan keperluan cuti..."
                        class="w-full px-3 py-2 rounded-xl border outline-none"
                        :style="{ backgroundColor: 'var(--input-bg)', borderColor: 'var(--input-border)', color: 'var(--input-text)' }"></textarea>
            </div>

            <div class="pt-3 border-t flex justify-end gap-2" :style="{ borderColor: 'var(--border-color)' }">
              <button type="button" @click="showLeaveModal = false" class="px-4 py-2 rounded-xl border font-bold text-slate-500 cursor-pointer">
                Batal
              </button>
              <button type="submit" :disabled="leaveSubmitting" class="px-5 py-2 rounded-xl font-bold text-white bg-blue-600 hover:bg-blue-700 cursor-pointer disabled:opacity-50">
                {{ leaveSubmitting ? 'Mengirimkan...' : 'Kirim Pengajuan Cuti' }}
              </button>
            </div>
          </form>

        </div>
      </div>
    </Teleport>

    <!-- ===================================================================== -->
    <!-- MODAL 2: RINCIAN SLIP GAJI RESMI                                      -->
    <!-- ===================================================================== -->
    <Teleport to="body">
      <div v-if="showSalarySlipModal" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-950/60 backdrop-blur-xs">
        <div class="w-full max-w-lg rounded-3xl p-6 sm:p-8 border shadow-2xl animate-in zoom-in-95 duration-200"
             :style="{ backgroundColor: 'var(--bg-primary)', borderColor: 'var(--border-color)', color: 'var(--text-primary)' }">
          
          <div class="flex items-center justify-between pb-4 border-b" :style="{ borderColor: 'var(--border-color)' }">
            <div class="flex items-center gap-3">
              <div class="w-10 h-10 rounded-2xl bg-emerald-100 dark:bg-emerald-950/50 text-emerald-600 flex items-center justify-center font-bold">
                ☕
              </div>
              <div>
                <h3 class="font-black text-base">SLIP GAJI ELEKTRONIK</h3>
                <p class="text-xs text-slate-400">Cafe ERP System • Bukti Penghasilan Resmi</p>
              </div>
            </div>
            <button @click="showSalarySlipModal = false" class="p-1.5 rounded-lg text-slate-400 hover:text-slate-700 cursor-pointer">✕</button>
          </div>

          <div class="py-4 space-y-4 text-xs">
            <!-- Header Identitas Pegawai -->
            <div class="p-3.5 rounded-2xl border grid grid-cols-2 gap-2"
                 :style="{ backgroundColor: 'var(--bg-content)', borderColor: 'var(--border-color)' }">
              <div>
                <span class="text-slate-400 text-[10px] block">NAMA KARYAWAN:</span>
                <span class="font-bold text-sm">{{ portalData.employee?.full_name || userFullName }}</span>
              </div>
              <div>
                <span class="text-slate-400 text-[10px] block">NIK / POSISI:</span>
                <span class="font-bold">{{ portalData.employee?.nik || 'EMP-003' }} ({{ portalData.employee?.job_title || 'Head Barista' }})</span>
              </div>
              <div>
                <span class="text-slate-400 text-[10px] block">PERIODE PENGGAJIAN:</span>
                <span class="font-semibold">{{ portalData.latest_payroll?.period_label || 'Agustus 2026' }}</span>
              </div>
              <div>
                <span class="text-slate-400 text-[10px] block">STATUS PEMBAYARAN:</span>
                <span class="font-bold text-emerald-600">LUNAS / DITRANSFER</span>
              </div>
            </div>

            <!-- Rincian Penghasilan (Earnings) -->
            <div>
              <div class="font-bold uppercase tracking-wider text-[11px] text-slate-400 mb-1.5">Penghasilan (Earnings)</div>
              <div class="space-y-1.5">
                <div class="flex justify-between py-1 border-b" :style="{ borderColor: 'var(--border-color)' }">
                  <span>Gaji Pokok</span>
                  <span class="font-mono font-bold">Rp {{ formatNum(portalData.latest_payroll?.basic_salary || 5000000) }}</span>
                </div>
                <div class="flex justify-between py-1 border-b" :style="{ borderColor: 'var(--border-color)' }">
                  <span>Tunjangan Posisi & Kehadiran</span>
                  <span class="font-mono font-bold">Rp {{ formatNum(portalData.latest_payroll?.allowances || 800000) }}</span>
                </div>
                <div class="flex justify-between py-1 border-b" :style="{ borderColor: 'var(--border-color)' }">
                  <span>Uang Lembur (Overtime)</span>
                  <span class="font-mono font-bold">Rp {{ formatNum(portalData.latest_payroll?.overtime_pay || 300000) }}</span>
                </div>
                <div class="flex justify-between py-1.5 font-bold text-slate-800 dark:text-slate-200">
                  <span>Total Penghasilan Kotor</span>
                  <span class="font-mono text-emerald-600">Rp {{ formatNum(portalData.latest_payroll?.gross_salary || 6100000) }}</span>
                </div>
              </div>
            </div>

            <!-- Rincian Potongan (Deductions) -->
            <div>
              <div class="font-bold uppercase tracking-wider text-[11px] text-slate-400 mb-1.5">Potongan (Deductions)</div>
              <div class="space-y-1.5">
                <div class="flex justify-between py-1 border-b" :style="{ borderColor: 'var(--border-color)' }">
                  <span>Iuran BPJS Ketenagakerjaan & Kesehatan</span>
                  <span class="font-mono text-rose-500 font-bold">- Rp 200.000</span>
                </div>
                <div class="flex justify-between py-1 border-b" :style="{ borderColor: 'var(--border-color)' }">
                  <span>Potongan Pajak PPh 21</span>
                  <span class="font-mono text-rose-500 font-bold">- Rp 150.000</span>
                </div>
                <div class="flex justify-between py-1.5 font-bold text-slate-800 dark:text-slate-200">
                  <span>Total Potongan</span>
                  <span class="font-mono text-rose-600">- Rp {{ formatNum(portalData.latest_payroll?.total_deductions || 350000) }}</span>
                </div>
              </div>
            </div>

            <!-- Net Total Take Home Pay -->
            <div class="p-4 rounded-2xl bg-emerald-50 dark:bg-emerald-950/40 border border-emerald-200 dark:border-emerald-800 flex items-center justify-between">
              <div>
                <div class="text-[10px] font-black uppercase text-emerald-800 dark:text-emerald-300">TOTAL TAKE HOME PAY</div>
                <div class="text-xs text-emerald-700 dark:text-emerald-400 font-semibold">Gaji Bersih Masuk Rekening</div>
              </div>
              <div class="text-2xl font-black font-mono text-emerald-600 dark:text-emerald-400">
                Rp {{ formatNum(portalData.latest_payroll?.net_salary || 5750000) }}
              </div>
            </div>
          </div>

          <div class="pt-4 border-t flex justify-end gap-2" :style="{ borderColor: 'var(--border-color)' }">
            <button @click="showSalarySlipModal = false" class="px-5 py-2 rounded-xl font-bold bg-blue-600 hover:bg-blue-700 text-white cursor-pointer">
              Tutup Slip Gaji
            </button>
          </div>

        </div>
      </div>
    </Teleport>

    <!-- ===================================================================== -->
    <!-- MODAL 3: PROFIL KEPEGAWAIAN LENGKAP                                   -->
    <!-- ===================================================================== -->
    <Teleport to="body">
      <div v-if="showProfileModal" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-950/60 backdrop-blur-xs">
        <div class="w-full max-w-md rounded-3xl p-6 sm:p-7 border shadow-2xl animate-in zoom-in-95 duration-200"
             :style="{ backgroundColor: 'var(--bg-primary)', borderColor: 'var(--border-color)', color: 'var(--text-primary)' }">
          
          <div class="flex items-center justify-between pb-3 border-b" :style="{ borderColor: 'var(--border-color)' }">
            <div class="flex items-center gap-2.5">
              <div class="w-9 h-9 rounded-xl bg-amber-100 text-amber-600 flex items-center justify-center font-bold">
                👤
              </div>
              <h3 class="font-bold text-base">Profil Lengkap Karyawan</h3>
            </div>
            <button @click="showProfileModal = false" class="p-1.5 rounded-lg text-slate-400 hover:text-slate-700 cursor-pointer">✕</button>
          </div>

          <div class="py-4 space-y-3 text-xs">
            <div class="p-3 rounded-xl border flex items-center justify-between"
                 :style="{ backgroundColor: 'var(--bg-content)', borderColor: 'var(--border-color)' }">
              <span class="text-slate-400">Nama Lengkap</span>
              <span class="font-bold text-sm">{{ portalData.employee?.full_name || userFullName }}</span>
            </div>
            <div class="p-3 rounded-xl border flex items-center justify-between"
                 :style="{ backgroundColor: 'var(--bg-content)', borderColor: 'var(--border-color)' }">
              <span class="text-slate-400">NIK Staf</span>
              <span class="font-mono font-bold">{{ portalData.employee?.nik || 'EMP-003' }}</span>
            </div>
            <div class="p-3 rounded-xl border flex items-center justify-between"
                 :style="{ backgroundColor: 'var(--bg-content)', borderColor: 'var(--border-color)' }">
              <span class="text-slate-400">Jabatan & Posisi</span>
              <span class="font-bold text-blue-600">{{ portalData.employee?.job_title || 'Head Barista' }}</span>
            </div>
            <div class="p-3 rounded-xl border flex items-center justify-between"
                 :style="{ backgroundColor: 'var(--bg-content)', borderColor: 'var(--border-color)' }">
              <span class="text-slate-400">Departemen</span>
              <span class="font-semibold">{{ portalData.employee?.department || 'Front of House' }}</span>
            </div>
            <div class="p-3 rounded-xl border flex items-center justify-between"
                 :style="{ backgroundColor: 'var(--bg-content)', borderColor: 'var(--border-color)' }">
              <span class="text-slate-400">Email Terdaftar</span>
              <span class="font-mono">{{ portalData.employee?.email || authStore.user?.email }}</span>
            </div>
            <div class="p-3 rounded-xl border flex items-center justify-between"
                 :style="{ backgroundColor: 'var(--bg-content)', borderColor: 'var(--border-color)' }">
              <span class="text-slate-400">Nomor Telepon</span>
              <span class="font-mono">{{ portalData.employee?.phone || '081234567821' }}</span>
            </div>
            <div class="p-3 rounded-xl border flex items-center justify-between"
                 :style="{ backgroundColor: 'var(--bg-content)', borderColor: 'var(--border-color)' }">
              <span class="text-slate-400">Tanggal Bergabung</span>
              <span class="font-semibold">{{ portalData.employee?.join_date || '2023-03-01' }}</span>
            </div>
            <div class="p-3 rounded-xl border flex items-center justify-between"
                 :style="{ backgroundColor: 'var(--bg-content)', borderColor: 'var(--border-color)' }">
              <span class="text-slate-400">Sisa Kuota Cuti</span>
              <span class="font-bold text-emerald-600">{{ portalData.leave_summary?.remaining_leave ?? 10 }} Hari</span>
            </div>
          </div>

          <div class="pt-3 border-t flex justify-end" :style="{ borderColor: 'var(--border-color)' }">
            <button @click="showProfileModal = false" class="px-5 py-2 rounded-xl font-bold bg-blue-600 text-white cursor-pointer">
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
import {
  DollarSign,
  ClipboardList,
  Package,
  Armchair,
  TrendingUp,
  AlertTriangle,
  Clock,
  Calendar,
  CalendarClock,
  UserCheck,
  CheckCircle,
  CheckCircle2,
  XCircle,
  Users,
  CreditCard,
  Landmark,
  Receipt,
  Briefcase,
  RefreshCw,
  Sparkles,
  ArrowRight,
  User
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
const floorTables = ref<any[]>([])

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

// Fallback Sample Data
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

// Computed Role Helpers
const userFullName = computed(() => {
  return authStore.fullName || authStore.user?.firstName || 'Pengguna'
})

const currentRoleName = computed(() => {
  return authStore.currentRole || 'Karyawan'
})

const roleLower = computed(() => {
  return String(authStore.currentRole || authStore.user?.roleId || '').toLowerCase()
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

const availableRolesCount = computed(() => {
  return isManagerialRole.value ? 4 : 1
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
  // Fallback 7-day view
  const days = ['Senin', 'Selasa', 'Rabu', 'Kamis', 'Jumat', 'Sabtu', 'Minggu']
  return days.map((d, i) => ({
    date: `2026-09-${28 + i}`,
    day_name: d,
    shift_name: i === 6 ? 'Libur' : i % 2 === 0 ? 'Pagi' : 'Siang',
    start_time: i === 6 ? '00:00' : i % 2 === 0 ? '07:00' : '15:00',
    end_time: i === 6 ? '00:00' : i % 2 === 0 ? '15:00' : '23:00',
    color: i === 6 ? '#94a3b8' : i % 2 === 0 ? '#2563eb' : '#16a34a',
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

const formatNum = (val: number | string) => {
  return Number(val || 0).toLocaleString('id-ID')
}

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

// Fetch Portal Data
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

// Fetch Finance Overview
const fetchFinanceData = async () => {
  try {
    const res = await api.get('/finance/overview')
    if (res.data?.data) {
      financeOverview.value = res.data.data
    }
  } catch (err) {
    console.warn('Could not fetch finance overview:', err)
  }
}

// Fetch HRIS Overview
const fetchHRISData = async () => {
  try {
    const [attRes, leaveRes] = await Promise.all([
      api.get('/hris/attendances/today-status'),
      api.get('/hris/leaves')
    ])
    if (attRes.data) {
      hrisTodayStats.value = attRes.data
    }
    if (leaveRes.data?.data) {
      pendingLeavesList.value = leaveRes.data.data.filter((l: any) => l.status === 'pending')
    }
  } catch (err) {
    console.warn('Could not fetch HRIS data:', err)
  }
}

// Fetch Executive Overview
const fetchExecutiveData = async () => {
  try {
    const [sRes, tRes] = await Promise.all([
      api.get('/dashboard/stats'),
      api.get('/pos/tables')
    ])
    if (sRes.data?.data) stats.value = sRes.data.data
    if (tRes.data?.data) floorTables.value = tRes.data.data
  } catch (err) {
    console.warn('Could not fetch executive dashboard stats:', err)
  }
}

// Clock In Action
const handleClockIn = async () => {
  attendanceActionLoading.value = true
  try {
    const res = await api.post('/hris/attendances/clock-in', {
      employee_id: portalData.value.employee?.id || '',
      notes: attendanceNoteInput.value.trim() || 'Hadir siap bertugas'
    })
    notifyStore.success(res.data?.message || 'Clock In berhasil dicatat!', 'Presensi Sukses')
    attendanceNoteInput.value = ''
    await fetchPortalData()
    await fetchHRISData()
  } catch (err: any) {
    const msg = err.response?.data?.message || 'Gagal melakukan Clock In'
    notifyStore.error(msg, 'Presensi')
  } finally {
    attendanceActionLoading.value = false
  }
}

// Clock Out Action
const handleClockOut = async () => {
  attendanceActionLoading.value = true
  try {
    const res = await api.post('/hris/attendances/clock-out', {
      employee_id: portalData.value.employee?.id || '',
      notes: attendanceNoteInput.value.trim() || 'Shift selesai'
    })
    notifyStore.success(res.data?.message || 'Clock Out berhasil dicatat. Terima kasih!', 'Selesai Shift')
    attendanceNoteInput.value = ''
    await fetchPortalData()
    await fetchHRISData()
  } catch (err: any) {
    const msg = err.response?.data?.message || 'Gagal melakukan Clock Out'
    notifyStore.error(msg, 'Presensi')
  } finally {
    attendanceActionLoading.value = false
  }
}

// Leave Form Handlers
const openLeaveModal = () => {
  const tomorrow = new Date()
  tomorrow.setDate(tomorrow.getDate() + 1)
  const tomorrowStr = tomorrow.toISOString().slice(0, 10)
  leaveForm.value = {
    leave_type: 'Cuti Tahunan',
    start_date: tomorrowStr,
    end_date: tomorrowStr,
    reason: ''
  }
  showLeaveModal.value = true
}

const submitLeaveForm = async () => {
  leaveSubmitting.value = true
  try {
    await api.post('/hris/leaves', {
      employee_id: portalData.value.employee?.id || '',
      leave_type: leaveForm.value.leave_type,
      start_date: leaveForm.value.start_date,
      end_date: leaveForm.value.end_date,
      reason: leaveForm.value.reason
    })
    notifyStore.success('Permohonan cuti berhasil diajukan dan menunggu persetujuan HR/Manager.', 'Sukses')
    showLeaveModal.value = false
    await fetchPortalData()
    await fetchHRISData()
  } catch (err: any) {
    const msg = err.response?.data?.message || 'Gagal mengirimkan permohonan cuti'
    notifyStore.error(msg, 'Gagal Mengajukan')
  } finally {
    leaveSubmitting.value = false
  }
}

const approveLeave = async (id: string) => {
  try {
    await api.put(`/hris/leaves/${id}/status`, { status: 'approved' })
    notifyStore.success('Permohonan cuti karyawan telah disetujui.', 'Approval Cuti')
    await fetchHRISData()
    await fetchPortalData()
  } catch (err: any) {
    notifyStore.error(err.response?.data?.message || 'Gagal menyetujui cuti', 'Error')
  }
}

const rejectLeave = async (id: string) => {
  try {
    await api.put(`/hris/leaves/${id}/status`, { status: 'rejected' })
    notifyStore.info('Permohonan cuti karyawan ditolak.', 'Cuti Ditolak')
    await fetchHRISData()
    await fetchPortalData()
  } catch (err: any) {
    notifyStore.error(err.response?.data?.message || 'Gagal menolak cuti', 'Error')
  }
}

const openSalarySlipModal = () => {
  showSalarySlipModal.value = true
}

const openProfileModal = () => {
  showProfileModal.value = true
}

const refreshCurrentView = async () => {
  isRefreshing.value = true
  try {
    await Promise.all([
      fetchPortalData(),
      fetchFinanceData(),
      fetchHRISData(),
      fetchExecutiveData()
    ])
    notifyStore.info('Data dashboard berhasil diperbarui.', 'Update')
  } finally {
    isRefreshing.value = false
  }
}

// Lifecycle Init
onMounted(() => {
  updateClock()
  clockTimer = setInterval(updateClock, 1000)

  // Determine initial view based on role
  if (isOperationalStaff.value) {
    activeViewTab.value = 'employee'
  } else if (isFinanceRole.value) {
    activeViewTab.value = 'finance'
  } else if (isHRISRole.value) {
    activeViewTab.value = 'hris'
  } else {
    activeViewTab.value = 'executive'
  }

  // Load all data
  fetchPortalData()
  fetchFinanceData()
  fetchHRISData()
  fetchExecutiveData()
})

onUnmounted(() => {
  if (clockTimer) clearInterval(clockTimer)
})
</script>

<style scoped>
.animate-spin-slow {
  animation: spin 8s linear infinite;
}
@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}
</style>
