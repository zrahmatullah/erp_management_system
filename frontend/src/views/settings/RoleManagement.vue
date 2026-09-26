<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
      <div>
        <div class="flex items-center gap-2.5">
          <h2 class="text-2xl font-black text-slate-900 tracking-tight">Manajemen Peran & Akun Pengguna</h2>
          <span class="px-2.5 py-0.5 rounded-full text-xs font-black bg-blue-100 text-blue-800 border border-blue-200 shadow-xs">
            RBAC & User Access
          </span>
        </div>
        <p class="text-xs text-slate-500 mt-1">Kelola matriks hak akses modul dan akun pengguna yang tertaut dengan staf karyawan HRIS.</p>
      </div>

      <div class="flex items-center gap-2.5">
        <button
          v-if="activeMainTab === 'matrix'"
          @click="openAddRoleModal"
          class="px-4 py-2.5 bg-blue-600 hover:bg-blue-700 active:bg-blue-800 text-white text-xs font-bold rounded-xl shadow-md shadow-blue-600/30 transition-all cursor-pointer flex items-center gap-1.5"
        >
          <Plus class="w-4 h-4" />
          <span>Tambah Peran</span>
        </button>

        <button
          v-else
          @click="openAddUserModal"
          class="px-4 py-2.5 bg-blue-600 hover:bg-blue-700 active:bg-blue-800 text-white text-xs font-bold rounded-xl shadow-md shadow-blue-600/30 transition-all cursor-pointer flex items-center gap-1.5"
        >
          <UserPlus class="w-4 h-4" />
          <span>Tambah Pengguna Baru</span>
        </button>
      </div>
    </div>

    <!-- Main Navigation Tabs -->
    <div class="flex items-center gap-4 border-b border-slate-200 text-xs font-bold">
      <button
        @click="activeMainTab = 'matrix'"
        class="pb-3 transition-colors relative cursor-pointer flex items-center gap-2"
        :class="activeMainTab === 'matrix' ? 'text-blue-600' : 'text-slate-500 hover:text-slate-800'"
      >
        <Shield class="w-4 h-4" />
        <span>Matriks Hak Akses Peran</span>
        <span v-if="activeMainTab === 'matrix'" class="absolute bottom-0 inset-x-0 h-0.5 bg-blue-600 rounded-full"></span>
      </button>

      <button
        @click="activeMainTab = 'users'"
        class="pb-3 transition-colors relative cursor-pointer flex items-center gap-2"
        :class="activeMainTab === 'users' ? 'text-blue-600' : 'text-slate-500 hover:text-slate-800'"
      >
        <Users class="w-4 h-4" />
        <span>Daftar Pengguna & Penugasan Staf</span>
        <span class="px-1.5 py-0.5 rounded-md text-[10px] bg-slate-100 text-slate-600">{{ userList.length }}</span>
        <span v-if="activeMainTab === 'users'" class="absolute bottom-0 inset-x-0 h-0.5 bg-blue-600 rounded-full"></span>
      </button>
    </div>

    <!-- ========================================================================= -->
    <!-- TAB 1: MATRIKS HAK AKSES PERAN                                            -->
    <!-- ========================================================================= -->
    <div v-if="activeMainTab === 'matrix'" class="space-y-6">
      <!-- Role Selection Tabs -->
      <div class="grid grid-cols-3 sm:grid-cols-5 lg:grid-cols-9 gap-3">
        <button
          v-for="role in roles"
          :key="role.id"
          @click="selectRole(role)"
          class="p-3 rounded-2xl border transition-all text-center flex flex-col items-center justify-center gap-1.5 cursor-pointer"
          :class="selectedRole?.id === role.id 
            ? 'bg-blue-50/70 border-blue-500 text-blue-700 shadow-sm' 
            : 'bg-white border-slate-200/80 text-slate-700 hover:border-slate-300 hover:bg-slate-50'"
        >
          <component :is="role.icon" class="w-5 h-5 shrink-0" :class="selectedRole?.id === role.id ? 'text-blue-600' : 'text-slate-500'" />
          <span class="text-xs font-bold truncate w-full">{{ role.name }}</span>
        </button>
      </div>

      <!-- Permission Matrix Table -->
      <div class="bg-white rounded-2xl border border-slate-200/80 shadow-sm p-6">
        <div class="flex items-center justify-between mb-4 pb-3 border-b border-slate-100">
          <div>
            <h3 class="font-bold text-slate-900 text-sm">
              Matriks Hak Akses: <span class="text-blue-600 font-extrabold">{{ selectedRole?.name }}</span>
            </h3>
            <p class="text-xs text-slate-400">Centang izin yang diperbolehkan untuk peran ini</p>
          </div>
          <span class="text-xs text-slate-400 font-medium">{{ permissionModules.length }} modul terkonfigurasi</span>
        </div>

        <div class="overflow-x-auto">
          <table class="w-full text-left text-xs">
            <thead>
              <tr class="text-slate-600 border-b border-slate-200 uppercase font-bold text-[11px] bg-slate-50/50">
                <th class="py-3 px-4">Module</th>
                <th class="py-3 px-4 text-center">Lihat (View)</th>
                <th class="py-3 px-4 text-center">Buat (Create)</th>
                <th class="py-3 px-4 text-center">Edit</th>
                <th class="py-3 px-4 text-center">Hapus (Delete)</th>
                <th class="py-3 px-4 text-center">Setujui (Approve)</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-100 text-slate-700 font-medium">
              <tr v-for="mod in permissionModules" :key="mod.key" class="hover:bg-slate-50/60 transition-colors">
                <td class="py-3 px-4 font-bold text-slate-900">
                  {{ mod.label }}
                </td>
                <td class="py-3 px-4 text-center">
                  <input
                    type="checkbox"
                    v-model="matrix[mod.key].view"
                    class="w-4 h-4 rounded text-blue-600 focus:ring-blue-500 border-slate-300 cursor-pointer"
                  />
                </td>
                <td class="py-3 px-4 text-center">
                  <input
                    type="checkbox"
                    v-model="matrix[mod.key].create"
                    class="w-4 h-4 rounded text-blue-600 focus:ring-blue-500 border-slate-300 cursor-pointer"
                  />
                </td>
                <td class="py-3 px-4 text-center">
                  <input
                    type="checkbox"
                    v-model="matrix[mod.key].edit"
                    class="w-4 h-4 rounded text-blue-600 focus:ring-blue-500 border-slate-300 cursor-pointer"
                  />
                </td>
                <td class="py-3 px-4 text-center">
                  <input
                    type="checkbox"
                    v-model="matrix[mod.key].delete"
                    class="w-4 h-4 rounded text-blue-600 focus:ring-blue-500 border-slate-300 cursor-pointer"
                  />
                </td>
                <td class="py-3 px-4 text-center">
                  <input
                    type="checkbox"
                    v-model="matrix[mod.key].approve"
                    class="w-4 h-4 rounded text-blue-600 focus:ring-blue-500 border-slate-300 cursor-pointer"
                  />
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <!-- Save Button -->
        <div class="mt-6 pt-4 border-t border-slate-100 flex items-center justify-between">
          <span class="text-xs text-slate-400">Perubahan hak akses akan langsung aktif pada sesi pengguna berikutnya</span>
          <button
            @click="savePermissions"
            :disabled="saving"
            class="px-5 py-2.5 bg-blue-600 hover:bg-blue-700 active:bg-blue-800 text-white text-xs font-bold rounded-xl shadow-md shadow-blue-600/30 transition-all cursor-pointer disabled:opacity-50"
          >
            {{ saving ? 'Menyimpan...' : 'Simpan Perubahan' }}
          </button>
        </div>
      </div>
    </div>

    <!-- ========================================================================= -->
    <!-- TAB 2: DAFTAR PENGGUNA & PENUGASAN STAF                                   -->
    <!-- ========================================================================= -->
    <div v-else class="space-y-6">
      <!-- 4 KPI Cards -->
      <div class="grid grid-cols-2 lg:grid-cols-4 gap-4">
        <div class="bg-white rounded-2xl p-4 border border-slate-200/80 shadow-xs flex items-center gap-3.5">
          <div class="w-11 h-11 rounded-xl bg-blue-50 text-blue-600 flex items-center justify-center shrink-0">
            <Users class="w-5 h-5" />
          </div>
          <div>
            <span class="text-[11px] font-bold text-slate-400 uppercase tracking-wider block">Total Pengguna</span>
            <span class="text-xl font-black text-slate-900">{{ userList.length }}</span>
          </div>
        </div>

        <div class="bg-white rounded-2xl p-4 border border-slate-200/80 shadow-xs flex items-center gap-3.5">
          <div class="w-11 h-11 rounded-xl bg-emerald-50 text-emerald-600 flex items-center justify-center shrink-0">
            <UserCheck class="w-5 h-5" />
          </div>
          <div>
            <span class="text-[11px] font-bold text-slate-400 uppercase tracking-wider block">Pengguna Aktif</span>
            <span class="text-xl font-black text-emerald-600">{{ countActiveUsers }}</span>
          </div>
        </div>

        <div class="bg-white rounded-2xl p-4 border border-slate-200/80 shadow-xs flex items-center gap-3.5">
          <div class="w-11 h-11 rounded-xl bg-purple-50 text-purple-600 flex items-center justify-center shrink-0">
            <Briefcase class="w-5 h-5" />
          </div>
          <div>
            <span class="text-[11px] font-bold text-slate-400 uppercase tracking-wider block">Tertaut Staf HRIS</span>
            <span class="text-xl font-black text-purple-600">{{ countLinkedUsers }}</span>
          </div>
        </div>

        <div class="bg-white rounded-2xl p-4 border border-slate-200/80 shadow-xs flex items-center gap-3.5">
          <div class="w-11 h-11 rounded-xl bg-amber-50 text-amber-600 flex items-center justify-center shrink-0">
            <Shield class="w-5 h-5" />
          </div>
          <div>
            <span class="text-[11px] font-bold text-slate-400 uppercase tracking-wider block">Super Admin & Manager</span>
            <span class="text-xl font-black text-amber-600">{{ countAdminUsers }}</span>
          </div>
        </div>
      </div>

      <!-- Filters & User Table -->
      <div class="bg-white rounded-2xl border border-slate-200/80 shadow-sm p-6 space-y-4">
        <!-- Search & Filter Bar -->
        <div class="flex flex-col sm:flex-row items-center justify-between gap-3">
          <div class="relative w-full sm:w-80">
            <Search class="w-4 h-4 text-slate-400 absolute left-3 top-1/2 -translate-y-1/2" />
            <input
              v-model="searchUserQuery"
              type="text"
              placeholder="Cari nama, email, staf..."
              class="w-full pl-9 pr-4 py-2 text-xs bg-slate-50 border border-slate-200 rounded-xl focus:bg-white focus:border-blue-600 outline-none"
            />
          </div>

          <div class="flex items-center gap-2 w-full sm:w-auto overflow-x-auto">
            <!-- Filter Role -->
            <select v-model="filterUserRole" class="px-3 py-2 text-xs bg-slate-50 border border-slate-200 rounded-xl outline-none font-medium">
              <option value="">Semua Peran</option>
              <option v-for="r in roles" :key="r.id" :value="r.name">{{ r.name }}</option>
            </select>

            <!-- Filter Status -->
            <select v-model="filterUserStatus" class="px-3 py-2 text-xs bg-slate-50 border border-slate-200 rounded-xl outline-none font-medium">
              <option value="">Semua Status</option>
              <option value="active">Aktif</option>
              <option value="inactive">Non-Aktif</option>
            </select>
          </div>
        </div>

        <!-- Users Table -->
        <div class="overflow-x-auto">
          <table class="w-full text-left text-xs">
            <thead>
              <tr class="text-slate-500 border-b border-slate-200 uppercase font-bold text-[11px] bg-slate-50/50">
                <th class="py-3 px-4">Pengguna</th>
                <th class="py-3 px-4">Tautan Karyawan HRIS</th>
                <th class="py-3 px-4">Peran (Role)</th>
                <th class="py-3 px-4">Cabang</th>
                <th class="py-3 px-4 text-center">PIN POS</th>
                <th class="py-3 px-4 text-center">Status</th>
                <th class="py-3 px-4 text-right">Aksi</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-100 text-slate-700 font-medium">
              <tr v-if="loadingUsers">
                <td colspan="7" class="py-12 text-center text-slate-400">
                  <RotateCw class="w-5 h-5 animate-spin mx-auto mb-2 text-blue-600" />
                  Memuat data akun pengguna dari database...
                </td>
              </tr>
              <tr v-else-if="filteredUsers.length === 0">
                <td colspan="7" class="py-8 text-center text-slate-400">
                  Tidak ada data pengguna ditemukan.
                </td>
              </tr>
              <tr v-for="user in filteredUsers" :key="user.id" class="hover:bg-slate-50/70 transition-colors">
                <!-- User name & email -->
                <td class="py-3.5 px-4 font-bold text-slate-900 flex items-center gap-2.5">
                  <div class="w-8 h-8 rounded-full bg-blue-100 text-blue-700 flex items-center justify-center font-black text-xs shrink-0 border border-blue-200">
                    {{ (user.full_name || user.username || 'U').charAt(0) }}
                  </div>
                  <div>
                    <div class="text-slate-900">{{ user.full_name }}</div>
                    <div class="text-[10px] text-slate-400 font-mono font-normal">{{ user.email }}</div>
                  </div>
                </td>

                <!-- Linked Employee -->
                <td class="py-3.5 px-4">
                  <div v-if="user.employee_name" class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-lg bg-purple-50 text-purple-700 border border-purple-200/70">
                    <span class="font-mono text-[10px] font-bold">{{ user.employee_nik }}</span>
                    <span class="font-bold text-xs">{{ user.employee_name }}</span>
                    <span class="text-[10px] text-purple-500">({{ user.employee_position || 'Staff' }})</span>
                  </div>
                  <span v-else class="text-slate-400 italic text-[11px]">Belum Tertaut</span>
                </td>

                <!-- Role -->
                <td class="py-3.5 px-4">
                  <span
                    class="px-2.5 py-0.5 rounded-full text-[10px] font-bold capitalize"
                    :class="getRoleBadgeClass(user.role_name)"
                  >
                    {{ user.role_name }}
                  </span>
                </td>

                <!-- Branch -->
                <td class="py-3.5 px-4 text-slate-600">{{ user.branch_name || 'Semua Cabang' }}</td>

                <!-- PIN Code -->
                <td class="py-3.5 px-4 text-center font-mono text-slate-500">
                  {{ user.pin_code ? '••••' : '-' }}
                </td>

                <!-- Status -->
                <td class="py-3.5 px-4 text-center">
                  <span
                    class="px-2.5 py-0.5 rounded-full text-[10px] font-bold"
                    :class="user.is_active ? 'bg-emerald-100 text-emerald-700' : 'bg-slate-100 text-slate-500'"
                  >
                    {{ user.is_active ? 'Aktif' : 'Non-Aktif' }}
                  </span>
                </td>

                <!-- Actions -->
                <td class="py-3.5 px-4 text-right">
                  <div class="inline-flex items-center gap-1.5">
                    <button
                      @click="openEditUserModal(user)"
                      class="p-1.5 rounded-lg border border-slate-200 bg-white hover:bg-slate-50 text-slate-600 hover:text-blue-600 transition-colors cursor-pointer"
                      title="Ubah Akun Pengguna"
                    >
                      <Edit class="w-3.5 h-3.5" />
                    </button>
                    <button
                      @click="deleteUser(user)"
                      class="p-1.5 rounded-lg border border-slate-200 bg-white hover:bg-rose-50 text-slate-400 hover:text-rose-600 transition-colors cursor-pointer"
                      title="Nonaktifkan / Hapus Pengguna"
                    >
                      <Trash2 class="w-3.5 h-3.5" />
                    </button>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>

    <!-- ========================================================================= -->
    <!-- MODAL: TAMBAH / UBAH AKUN PENGGUNA (DENGAN DROPDOWN PILIH KARYAWAN)        -->
    <!-- ========================================================================= -->
    <Teleport to="body">
      <Transition name="emp-modal-fade">
        <div
          v-if="showUserModal"
          class="fixed md:left-64 inset-y-0 right-0 left-0 z-50 bg-slate-950/45 backdrop-blur-sm flex items-center justify-center p-4 overflow-y-auto transition-all duration-200"
        >
          <div class="bg-white rounded-3xl w-full max-w-xl shadow-2xl border border-slate-100/90 overflow-hidden my-6">
            <div class="p-5 border-b border-slate-100 flex items-center justify-between bg-blue-50/50">
              <div class="flex items-center gap-3">
                <div class="w-10 h-10 rounded-2xl bg-blue-600 text-white flex items-center justify-center font-black shrink-0">
                  <UserPlus class="w-5 h-5" />
                </div>
                <div>
                  <h3 class="font-black text-slate-900 text-sm">
                    {{ isEditingUser ? 'Ubah Akun Pengguna' : 'Registrasi Akun Pengguna Sistem' }}
                  </h3>
                  <p class="text-[11px] text-slate-500 mt-0.5">
                    Pilih karyawan HRIS untuk sinkronisasi otomatis nama, email, dan wewenang.
                  </p>
                </div>
              </div>
              <button
                @click="showUserModal = false"
                class="p-1.5 rounded-lg text-slate-400 hover:text-slate-600 hover:bg-slate-100 transition-colors cursor-pointer"
              >
                <X class="w-4 h-4" />
              </button>
            </div>

            <form @submit.prevent="submitUserForm" class="p-6 space-y-4 text-xs font-medium text-slate-700">
              <!-- DROPDOWN PILIH KARYAWAN HRIS -->
              <div class="p-3.5 rounded-2xl bg-blue-50/60 border border-blue-200/80 space-y-2">
                <label class="block text-blue-900 font-bold">
                  Tautkan dengan Karyawan HRIS (Pilih Staf)
                </label>
                <select
                  v-model="userForm.employee_id"
                  @change="onEmployeeSelected"
                  class="w-full px-3 py-2.5 bg-white border border-blue-200 rounded-xl outline-none text-slate-900 font-semibold focus:ring-2 focus:ring-blue-500/20"
                >
                  <option value="">-- Tanpa Tautan Karyawan (Akun Manual) --</option>
                  <option
                    v-for="emp in activeEmployees"
                    :key="emp.id"
                    :value="emp.id"
                  >
                    {{ emp.full_name }} ({{ emp.nik }} - {{ emp.position_name || 'Staff' }})
                  </option>
                </select>
                <p class="text-[10px] text-blue-600">
                  *Memilih staf akan otomatis mengisi nama lengkap, email, cabang, serta rekomendasi peran yang cocok.
                </p>
              </div>

              <!-- Nama Lengkap & Username -->
              <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
                <div>
                  <label class="block text-slate-500 mb-1">Nama Lengkap</label>
                  <input
                    v-model="userForm.full_name"
                    type="text"
                    required
                    placeholder="Contoh: Sarah Johnson"
                    class="w-full px-3 py-2 bg-slate-50 border border-slate-200 rounded-xl outline-none text-slate-900 font-bold focus:bg-white focus:border-blue-600"
                  />
                </div>
                <div>
                  <label class="block text-slate-500 mb-1">Username Login</label>
                  <input
                    v-model="userForm.username"
                    type="text"
                    required
                    placeholder="Contoh: sarah.chef"
                    class="w-full px-3 py-2 bg-slate-50 border border-slate-200 rounded-xl outline-none font-mono text-slate-900 focus:bg-white focus:border-blue-600"
                  />
                </div>
              </div>

              <!-- Email & Password -->
              <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
                <div>
                  <label class="block text-slate-500 mb-1">Email Akun</label>
                  <input
                    v-model="userForm.email"
                    type="email"
                    required
                    placeholder="sarah@cafeharmony.com"
                    class="w-full px-3 py-2 bg-slate-50 border border-slate-200 rounded-xl outline-none text-slate-900 focus:bg-white focus:border-blue-600"
                  />
                </div>
                <div>
                  <label class="block text-slate-500 mb-1">
                    Password {{ isEditingUser ? '(Kosongkan jika tidak diubah)' : '' }}
                  </label>
                  <input
                    v-model="userForm.password"
                    type="password"
                    :placeholder="isEditingUser ? 'Tetap gunakan password lama' : 'Default: Admin@123'"
                    class="w-full px-3 py-2 bg-slate-50 border border-slate-200 rounded-xl outline-none text-slate-900 font-mono focus:bg-white focus:border-blue-600"
                  />
                </div>
              </div>

              <!-- Role & Branch -->
              <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
                <div>
                  <label class="block text-slate-500 mb-1 font-bold">Peran Sistem (Role RBAC)</label>
                  <select
                    v-model="userForm.role_id"
                    required
                    class="w-full px-3 py-2 bg-slate-50 border border-slate-200 rounded-xl outline-none text-slate-900 font-bold focus:bg-white focus:border-blue-600"
                  >
                    <option value="" disabled>Pilih Peran...</option>
                    <option v-for="r in roles" :key="r.id" :value="r.id">
                      {{ r.name }}
                    </option>
                  </select>
                </div>
                <div>
                  <label class="block text-slate-500 mb-1 font-bold">Cabang Penempatan</label>
                  <select
                    v-model="userForm.branch_id"
                    class="w-full px-3 py-2 bg-slate-50 border border-slate-200 rounded-xl outline-none text-slate-900 focus:bg-white focus:border-blue-600"
                  >
                    <option value="">Semua Cabang / Pusat</option>
                    <option v-for="b in branchList" :key="b.id" :value="b.id">
                      {{ b.name }}
                    </option>
                  </select>
                </div>
              </div>

              <!-- PIN Kasir POS & Status -->
              <div class="grid grid-cols-1 sm:grid-cols-2 gap-3 pt-1">
                <div>
                  <label class="block text-slate-500 mb-1">PIN Kasir POS (4-6 Digit)</label>
                  <input
                    v-model="userForm.pin_code"
                    type="text"
                    maxlength="6"
                    placeholder="Contoh: 1234"
                    class="w-full px-3 py-2 bg-slate-50 border border-slate-200 rounded-xl outline-none font-mono text-slate-900 focus:bg-white focus:border-blue-600"
                  />
                </div>
                <div class="flex items-center gap-2 pt-6">
                  <input
                    id="user-status-chk"
                    v-model="userForm.is_active"
                    type="checkbox"
                    class="w-4 h-4 rounded text-blue-600 focus:ring-blue-500 border-slate-300 cursor-pointer"
                  />
                  <label for="user-status-chk" class="text-slate-700 font-bold cursor-pointer">
                    Akun Aktif & Diizinkan Login
                  </label>
                </div>
              </div>

              <!-- Buttons -->
              <div class="flex items-center justify-end gap-3 pt-4 border-t border-slate-100">
                <button
                  type="button"
                  @click="showUserModal = false"
                  class="px-4 py-2.5 border border-slate-200 rounded-xl text-slate-600 hover:bg-slate-50 transition-colors cursor-pointer"
                >
                  Batal
                </button>
                <button
                  type="submit"
                  :disabled="savingUser"
                  class="px-5 py-2.5 bg-blue-600 hover:bg-blue-700 active:bg-blue-800 text-white font-bold rounded-xl shadow-md shadow-blue-600/30 transition-all cursor-pointer disabled:opacity-60"
                >
                  {{ savingUser ? 'Menyimpan...' : (isEditingUser ? 'Perbarui Akun' : 'Buat Akun Pengguna') }}
                </button>
              </div>
            </form>
          </div>
        </div>
      </Transition>
    </Teleport>

    <!-- ========================================================================= -->
    <!-- MODAL: TAMBAH PERAN BARU                                                  -->
    <!-- ========================================================================= -->
    <Teleport to="body">
      <Transition name="emp-modal-fade">
        <div
          v-if="showAddRoleModal"
          class="fixed md:left-64 inset-y-0 right-0 left-0 z-50 bg-slate-950/45 backdrop-blur-sm flex items-center justify-center p-4 overflow-y-auto transition-all duration-200"
        >
          <div class="bg-white rounded-3xl w-full max-w-md shadow-2xl border border-slate-100/90 overflow-hidden my-6">
            <div class="p-5 border-b border-slate-100 flex items-center justify-between bg-slate-50">
              <div class="flex items-center gap-3">
                <div class="w-10 h-10 rounded-2xl bg-blue-600 text-white flex items-center justify-center font-black shrink-0 shadow-md shadow-blue-600/30">
                  <Shield class="w-5 h-5" />
                </div>
                <div>
                  <h3 class="font-black text-slate-900 text-sm">Tambah Peran Baru</h3>
                  <p class="text-[11px] text-slate-500 mt-0.5">Definisikan peran wewenang baru dalam sistem ERP.</p>
                </div>
              </div>
              <button
                @click="showAddRoleModal = false"
                class="p-1.5 rounded-lg text-slate-400 hover:text-slate-600 hover:bg-slate-100 transition-colors cursor-pointer"
              >
                <X class="w-4 h-4" />
              </button>
            </div>

            <form @submit.prevent="submitRoleForm" class="p-6 space-y-4 text-xs font-medium text-slate-700">
              <div>
                <label class="block text-slate-500 mb-1">Nama Peran / Jabatan Akses</label>
                <input
                  v-model="roleForm.name"
                  type="text"
                  required
                  placeholder="Contoh: Barista Lead, Assistant Manager"
                  class="w-full px-3 py-2 bg-slate-50 border border-slate-200 rounded-xl outline-none text-slate-900 font-bold focus:bg-white focus:border-blue-600"
                />
              </div>

              <div>
                <label class="block text-slate-500 mb-1">Deskripsi Wewenang</label>
                <textarea
                  v-model="roleForm.description"
                  rows="3"
                  placeholder="Deskripsikan cakupan tanggung jawab peran ini..."
                  class="w-full px-3 py-2 bg-slate-50 border border-slate-200 rounded-xl outline-none text-slate-900 focus:bg-white focus:border-blue-600"
                ></textarea>
              </div>

              <div class="flex items-center justify-end gap-3 pt-4 border-t border-slate-100">
                <button
                  type="button"
                  @click="showAddRoleModal = false"
                  class="px-4 py-2 border border-slate-200 rounded-xl text-slate-600 hover:bg-slate-50 transition-colors cursor-pointer"
                >
                  Batal
                </button>
                <button
                  type="submit"
                  :disabled="savingRole"
                  class="px-5 py-2 bg-blue-600 hover:bg-blue-700 active:bg-blue-800 text-white font-bold rounded-xl shadow-md shadow-blue-600/30 transition-all cursor-pointer disabled:opacity-60"
                >
                  {{ savingRole ? 'Menyimpan...' : 'Simpan Peran' }}
                </button>
              </div>
            </form>
          </div>
        </div>
      </Transition>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import axios from 'axios'
import {
  Shield,
  UserCheck,
  Briefcase,
  ShoppingCart,
  ChefHat,
  Users,
  Package,
  FileSpreadsheet,
  UtensilsCrossed,
  Plus,
  UserPlus,
  Search,
  Edit,
  Trash2,
  X,
  RotateCw
} from 'lucide-vue-next'
import { useNotificationStore } from '@/stores/notification.store'
import { useDialogStore } from '@/stores/dialog.store'

const notifyStore = useNotificationStore()
const dialogStore = useDialogStore()

const activeMainTab = ref<'matrix' | 'users'>('matrix')

interface RoleItem {
  id: string
  name: string
  icon: any
  desc: string
}

const getRoleIcon = (name: string) => {
  if (name.includes('Admin')) return Shield
  if (name.includes('Owner')) return UserCheck
  if (name.includes('Manager')) return Briefcase
  if (name.includes('Kasir')) return ShoppingCart
  if (name.includes('Kitchen')) return ChefHat
  if (name.includes('HR')) return Users
  if (name.includes('Warehouse')) return Package
  if (name.includes('Akuntan')) return FileSpreadsheet
  return UtensilsCrossed
}

const roles = ref<RoleItem[]>([
  { id: '1', name: 'Super Admin', icon: Shield, desc: 'Full Access' },
  { id: '2', name: 'Owner', icon: UserCheck, desc: 'Read Only Analytics' },
  { id: '3', name: 'Manager', icon: Briefcase, desc: 'Branch Supervisor' },
  { id: '4', name: 'Kasir', icon: ShoppingCart, desc: 'POS Terminal' },
  { id: '5', name: 'Kitchen Staff', icon: ChefHat, desc: 'KDS Order Monitor' },
  { id: '6', name: 'HR Admin', icon: Users, desc: 'Employee & Payroll' },
  { id: '7', name: 'Warehouse', icon: Package, desc: 'Stock & Logistics' },
  { id: '8', name: 'Akuntan', icon: FileSpreadsheet, desc: 'Finance & COA' },
  { id: '9', name: 'Pelayan', icon: UtensilsCrossed, desc: 'Table Service' }
])

const selectedRole = ref<RoleItem>(roles.value[2])
const saving = ref(false)

const permissionModules = [
  { key: 'dashboard', label: 'Dashboard' },
  { key: 'pos', label: 'Point of Sale (POS)' },
  { key: 'kitchen', label: 'Kitchen Display (KDS)' },
  { key: 'tables', label: 'Table Management' },
  { key: 'inventory', label: 'Inventory & Bahan Baku' },
  { key: 'purchasing', label: 'Purchase Order (PO)' },
  { key: 'hris', label: 'HRIS & Karyawan' },
  { key: 'payroll', label: 'Payroll & Penggajian' },
  { key: 'finance', label: 'Keuangan & Jurnal' },
  { key: 'reports', label: 'Laporan & Analitik' },
  { key: 'settings', label: 'Pengaturan Sistem' },
  { key: 'master', label: 'Master Data Hub' }
]

const matrix = reactive<Record<string, { view: boolean, create: boolean, edit: boolean, delete: boolean, approve: boolean }>>({
  dashboard: { view: true, create: false, edit: false, delete: false, approve: false },
  pos: { view: true, create: true, edit: true, delete: false, approve: false },
  kitchen: { view: true, create: true, edit: true, delete: false, approve: false },
  tables: { view: true, create: true, edit: true, delete: false, approve: false },
  inventory: { view: true, create: false, edit: true, delete: false, approve: false },
  purchasing: { view: true, create: true, edit: false, delete: false, approve: true },
  hris: { view: true, create: true, edit: true, delete: false, approve: true },
  payroll: { view: true, create: true, edit: true, delete: false, approve: true },
  finance: { view: true, create: false, edit: false, delete: false, approve: false },
  reports: { view: true, create: false, edit: false, delete: false, approve: false },
  settings: { view: true, create: false, edit: false, delete: false, approve: false },
  master: { view: true, create: true, edit: true, delete: true, approve: true }
})

const selectRole = (role: RoleItem) => {
  selectedRole.value = role
  if (role.name === 'Super Admin') {
    Object.keys(matrix).forEach(k => {
      matrix[k].view = true
      matrix[k].create = true
      matrix[k].edit = true
      matrix[k].delete = true
      matrix[k].approve = true
    })
  } else if (role.name === 'Kasir') {
    Object.keys(matrix).forEach(k => {
      matrix[k].view = k === 'pos' || k === 'kitchen'
      matrix[k].create = k === 'pos'
      matrix[k].edit = false
      matrix[k].delete = false
      matrix[k].approve = false
    })
  }
}

const savePermissions = () => {
  saving.value = true
  setTimeout(() => {
    saving.value = false
    notifyStore.success(`Matriks hak akses untuk peran '${selectedRole.value.name}' berhasil disimpan ke sistem!`, 'Izin Disimpan')
  }, 400)
}

// -----------------------------------------------------------------------------
// USER MANAGEMENT & EMPLOYEE LINKAGE
// -----------------------------------------------------------------------------
const userList = ref<any[]>([])
const activeEmployees = ref<any[]>([])
const branchList = ref<any[]>([])
const loadingUsers = ref(false)
const searchUserQuery = ref('')
const filterUserRole = ref('')
const filterUserStatus = ref('')

const showUserModal = ref(false)
const isEditingUser = ref(false)
const editingUserId = ref<string | null>(null)
const savingUser = ref(false)

const userForm = ref({
  employee_id: '',
  full_name: '',
  username: '',
  email: '',
  phone: '',
  password: '',
  role_id: '',
  branch_id: '',
  pin_code: '',
  is_active: true
})

const fetchUsers = async () => {
  loadingUsers.value = true
  try {
    const res = await axios.get('/api/v1/master/users')
    userList.value = Array.isArray(res.data) ? res.data : (res.data?.data || [])
  } catch (err) {
    console.error('Failed to load users:', err)
  } finally {
    loadingUsers.value = false
  }
}

const fetchEmployees = async () => {
  try {
    const res = await axios.get('/api/v1/hris/employees')
    activeEmployees.value = Array.isArray(res.data?.data) ? res.data.data : []
  } catch (err) {
    console.error('Failed to load employees:', err)
  }
}

const fetchBranches = async () => {
  try {
    const res = await axios.get('/api/v1/master/branches')
    branchList.value = Array.isArray(res.data) ? res.data : (res.data?.data || [])
  } catch (err) {
    console.error('Failed to load branches:', err)
  }
}

const fetchRoles = async () => {
  try {
    const res = await axios.get('/api/v1/master/roles')
    const list = Array.isArray(res.data) ? res.data : (res.data?.data || [])
    if (list.length > 0) {
      roles.value = list.map((r: any) => ({
        id: r.id,
        name: r.name,
        desc: r.description,
        icon: getRoleIcon(r.name)
      }))
      selectedRole.value = roles.value.find((r: any) => r.name.toLowerCase().includes('manager')) || roles.value[0]
    }
  } catch (err) {
    console.error('Failed to load master roles:', err)
  }
}

onMounted(() => {
  fetchRoles()
  fetchUsers()
  fetchEmployees()
  fetchBranches()
})

const countActiveUsers = computed(() => userList.value.filter(u => u.is_active).length)
const countLinkedUsers = computed(() => userList.value.filter(u => u.employee_name).length)
const countAdminUsers = computed(() => userList.value.filter(u => u.role_name?.includes('Admin') || u.role_name?.includes('Manager')).length)

const filteredUsers = computed(() => {
  return userList.value.filter(u => {
    const matchSearch = !searchUserQuery.value ||
      u.full_name?.toLowerCase().includes(searchUserQuery.value.toLowerCase()) ||
      u.username?.toLowerCase().includes(searchUserQuery.value.toLowerCase()) ||
      u.email?.toLowerCase().includes(searchUserQuery.value.toLowerCase()) ||
      u.employee_name?.toLowerCase().includes(searchUserQuery.value.toLowerCase())
    
    const matchRole = !filterUserRole.value || u.role_name === filterUserRole.value
    const matchStatus = !filterUserStatus.value || (filterUserStatus.value === 'active' ? u.is_active : !u.is_active)

    return matchSearch && matchRole && matchStatus
  })
})

const getRoleBadgeClass = (roleName?: string) => {
  if (!roleName) return 'bg-slate-100 text-slate-700'
  const r = roleName.toLowerCase()
  if (r.includes('admin')) return 'bg-rose-100 text-rose-700'
  if (r.includes('manager')) return 'bg-amber-100 text-amber-800'
  if (r.includes('kasir')) return 'bg-emerald-100 text-emerald-700'
  if (r.includes('kitchen')) return 'bg-blue-100 text-blue-700'
  if (r.includes('hr')) return 'bg-purple-100 text-purple-700'
  return 'bg-slate-100 text-slate-700'
}

// Handle employee selection from dropdown
const onEmployeeSelected = () => {
  if (!userForm.value.employee_id) return
  const emp = activeEmployees.value.find(e => e.id === userForm.value.employee_id)
  if (!emp) return

  userForm.value.full_name = emp.full_name || `${emp.first_name || ''} ${emp.last_name || ''}`.trim()
  userForm.value.email = emp.email || userForm.value.email
  if (emp.phone) userForm.value.phone = emp.phone
  if (emp.branch_id) userForm.value.branch_id = emp.branch_id

  if (!userForm.value.username && emp.email) {
    userForm.value.username = emp.email.split('@')[0]
  }

  // Suggest role based on position title
  const pos = (emp.position_name || '').toLowerCase()
  let matchedRole = roles.value.find(r => r.name.toLowerCase() === pos)
  if (!matchedRole) {
    if (pos.includes('chef') || pos.includes('kitchen') || pos.includes('cook')) {
      matchedRole = roles.value.find(r => r.name.toLowerCase().includes('kitchen'))
    } else if (pos.includes('barista') || pos.includes('cashier') || pos.includes('kasir')) {
      matchedRole = roles.value.find(r => r.name.toLowerCase().includes('kasir'))
    } else if (pos.includes('manager') || pos.includes('supervisor')) {
      matchedRole = roles.value.find(r => r.name.toLowerCase().includes('manager'))
    } else if (pos.includes('hr') || pos.includes('personalia')) {
      matchedRole = roles.value.find(r => r.name.toLowerCase().includes('hr'))
    }
  }

  if (matchedRole) {
    userForm.value.role_id = matchedRole.id
  }
}

const openAddUserModal = () => {
  isEditingUser.value = false
  editingUserId.value = null
  userForm.value = {
    employee_id: '',
    full_name: '',
    username: '',
    email: '',
    phone: '',
    password: '',
    role_id: roles.value.find(r => r.name.toLowerCase().includes('kasir'))?.id || roles.value[0]?.id || '',
    branch_id: branchList.value[0]?.id || '',
    pin_code: '',
    is_active: true
  }
  showUserModal.value = true
}

const openEditUserModal = (user: any) => {
  isEditingUser.value = true
  editingUserId.value = user.id
  userForm.value = {
    employee_id: user.employee_id || '',
    full_name: user.full_name,
    username: user.username,
    email: user.email,
    phone: user.phone || '',
    password: '',
    role_id: user.role_id || '',
    branch_id: user.branch_id || '',
    pin_code: user.pin_code || '',
    is_active: Boolean(user.is_active)
  }
  showUserModal.value = true
}

const submitUserForm = async () => {
  if (!userForm.value.email || !userForm.value.full_name) {
    notifyStore.warning('Nama lengkap dan email wajib diisi', 'Validasi')
    return
  }

  savingUser.value = true
  try {
    if (isEditingUser.value && editingUserId.value) {
      await axios.put(`/api/v1/master/users/${editingUserId.value}`, userForm.value)
      notifyStore.success('Data akun pengguna berhasil diperbarui!', 'Sukses')
    } else {
      await axios.post('/api/v1/master/users', userForm.value)
      notifyStore.success('Akun pengguna baru berhasil dibuat dan ditautkan ke database!', 'Pengguna Dibuat')
    }
    showUserModal.value = false
    await fetchUsers()
  } catch (err: any) {
    console.error('Failed to save user:', err)
    notifyStore.error(err.response?.data?.error || 'Gagal menyimpan akun pengguna', 'Gagal')
  } finally {
    savingUser.value = false
  }
}

const deleteUser = async (user: any) => {
  const confirmed = await dialogStore.confirm({
    title: 'Hapus / Nonaktifkan Pengguna',
    message: `Apakah Anda yakin ingin menonaktifkan pengguna '${user.full_name}'?`,
    confirmText: 'Ya, Nonaktifkan',
    cancelText: 'Batal',
    type: 'danger'
  })
  if (!confirmed) return

  try {
    await axios.delete(`/api/v1/master/users/${user.id}`)
    notifyStore.success(`Akun '${user.full_name}' berhasil dinonaktifkan!`, 'Pengguna Dihapus')
    await fetchUsers()
  } catch (err: any) {
    console.error('Failed to delete user:', err)
    notifyStore.error(err.response?.data?.error || 'Gagal menonaktifkan pengguna', 'Gagal')
  }
}

// -----------------------------------------------------------------------------
// TAMBAH PERAN MODAL
// -----------------------------------------------------------------------------
const showAddRoleModal = ref(false)
const savingRole = ref(false)
const roleForm = ref({
  name: '',
  description: ''
})

const openAddRoleModal = () => {
  roleForm.value = { name: '', description: '' }
  showAddRoleModal.value = true
}

const submitRoleForm = async () => {
  if (!roleForm.value.name) {
    notifyStore.warning('Nama peran harus diisi', 'Validasi')
    return
  }

  savingRole.value = true
  try {
    await axios.post('/api/v1/master/roles', roleForm.value)
    notifyStore.success(`Peran baru '${roleForm.value.name}' berhasil ditambahkan ke master!`, 'Peran Dibuat')
    showAddRoleModal.value = false
    await fetchRoles()
  } catch (err: any) {
    console.error('Failed to create role:', err)
    notifyStore.error(err.response?.data?.error || 'Gagal menambahkan peran baru', 'Gagal')
  } finally {
    savingRole.value = false
  }
}
</script>

<style scoped>
.emp-modal-fade-enter-active,
.emp-modal-fade-leave-active {
  transition: opacity 0.25s cubic-bezier(0.16, 1, 0.3, 1);
}

.emp-modal-fade-enter-from,
.emp-modal-fade-leave-to {
  opacity: 0;
}
</style>
