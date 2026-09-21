<template>
  <div class="min-h-screen w-full flex items-center justify-center bg-gradient-to-br from-slate-950 via-slate-900 to-indigo-950 p-4 relative overflow-hidden select-none">
    <!-- Ambient glowing backgrounds -->
    <div class="absolute -top-40 -left-40 w-96 h-96 bg-blue-600/20 rounded-full blur-3xl pointer-events-none"></div>
    <div class="absolute -bottom-40 -right-40 w-96 h-96 bg-indigo-600/20 rounded-full blur-3xl pointer-events-none"></div>
    <div class="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[600px] h-[600px] bg-blue-500/5 rounded-full blur-3xl pointer-events-none"></div>

    <!-- Decorative Coffee Line Art on right -->
    <div class="hidden lg:block absolute right-12 bottom-0 opacity-15 pointer-events-none max-w-lg select-none">
      <svg viewBox="0 0 400 400" fill="none" stroke="white" stroke-width="1.5" class="w-full h-auto">
        <path d="M120 280 C120 330, 260 330, 260 280 L250 210 L130 210 Z" />
        <path d="M250 230 C280 230, 290 260, 260 270" />
        <path d="M160 180 C160 150, 180 150, 180 120" stroke-dasharray="3 3" />
        <path d="M200 180 C200 150, 220 150, 220 120" stroke-dasharray="3 3" />
        <ellipse cx="190" cy="210" rx="60" ry="12" />
        <ellipse cx="190" cy="330" rx="90" ry="16" />
        <ellipse cx="90" cy="300" rx="18" ry="12" transform="rotate(-30 90 300)" />
        <path d="M82 292 Q90 300 98 308" />
        <ellipse cx="290" cy="310" rx="20" ry="14" transform="rotate(25 290 310)" />
        <path d="M280 305 Q290 310 300 315" />
      </svg>
    </div>

    <!-- Centered Glass Card -->
    <div class="w-full max-w-[420px] bg-white rounded-3xl shadow-2xl p-7 sm:p-9 border border-slate-100/90 relative z-10 transition-all duration-300">
      
      <!-- Top Brand -->
      <div class="flex items-center justify-between mb-6 pb-4 border-b border-slate-100">
        <div class="flex items-center gap-3">
          <div class="w-10 h-10 rounded-2xl bg-blue-600 flex items-center justify-center text-white shadow-md shadow-blue-500/25">
            <Coffee class="w-5 h-5" />
          </div>
          <div>
            <h2 class="text-lg font-black tracking-tight text-slate-900 leading-tight">Cafe ERP</h2>
            <p class="text-[10px] font-bold text-blue-600 uppercase tracking-wider">Enterprise Suite</p>
          </div>
        </div>
        <span class="px-2.5 py-1 rounded-full text-[10px] font-extrabold bg-emerald-50 text-emerald-700 border border-emerald-200/60">
          Online
        </span>
      </div>

      <!-- Heading -->
      <div class="mb-5">
        <h1 class="text-2xl font-black text-slate-900 tracking-tight">Selamat Datang</h1>
        <p class="text-xs text-slate-500 mt-1 font-medium">Masuk ke akun operasional restoran & cafe Anda</p>
      </div>

      <!-- Login Form -->
      <form @submit.prevent="handleLogin" class="space-y-4">
        <!-- Email Field -->
        <div>
          <label class="block text-xs font-bold text-slate-700 mb-1.5">Alamat Email / Akun</label>
          <div class="relative">
            <Mail class="w-4 h-4 text-slate-400 absolute left-3.5 top-1/2 -translate-y-1/2 pointer-events-none" />
            <input
              v-model="email"
              type="email"
              required
              placeholder="nama@cafe-erp.com"
              class="w-full pl-10 pr-4 py-2.5 text-xs sm:text-sm bg-slate-50 border border-slate-200 rounded-xl focus:bg-white focus:border-blue-600 focus:ring-3 focus:ring-blue-100 outline-hidden transition-all text-slate-800 font-medium placeholder:text-slate-400"
            />
          </div>
        </div>

        <!-- Password Field -->
        <div>
          <label class="block text-xs font-bold text-slate-700 mb-1.5">Kata Sandi (Password)</label>
          <div class="relative">
            <Lock class="w-4 h-4 text-slate-400 absolute left-3.5 top-1/2 -translate-y-1/2 pointer-events-none" />
            <input
              v-model="password"
              :type="showPassword ? 'text' : 'password'"
              required
              placeholder="Masukkan kata sandi"
              class="w-full pl-10 pr-10 py-2.5 text-xs sm:text-sm bg-slate-50 border border-slate-200 rounded-xl focus:bg-white focus:border-blue-600 focus:ring-3 focus:ring-blue-100 outline-hidden transition-all text-slate-800 font-medium placeholder:text-slate-400"
            />
            <button
              type="button"
              @click="showPassword = !showPassword"
              class="absolute right-3.5 top-1/2 -translate-y-1/2 text-slate-400 hover:text-slate-600 p-1 cursor-pointer transition-colors"
              tabindex="-1"
            >
              <EyeOff v-if="showPassword" class="w-4 h-4" />
              <Eye v-else class="w-4 h-4" />
            </button>
          </div>
        </div>

        <!-- Remember me & Forgot Password -->
        <div class="flex items-center justify-between text-xs pt-0.5">
          <label class="flex items-center gap-2 cursor-pointer text-slate-600 select-none">
            <input 
              v-model="rememberMe" 
              type="checkbox" 
              class="w-3.5 h-3.5 rounded border-slate-300 text-blue-600 focus:ring-blue-500 cursor-pointer" 
            />
            <span class="font-medium text-[11px]">Ingat saya</span>
          </label>
          <a href="#" @click.prevent="handleForgotPassword" class="text-blue-600 font-bold hover:underline text-[11px]">
            Lupa Password?
          </a>
        </div>

        <!-- Error Alert -->
        <div v-if="errorMessage" class="p-3 bg-red-50 border border-red-200 rounded-xl text-red-600 text-xs font-medium flex items-center gap-2">
          <span class="w-1.5 h-1.5 rounded-full bg-red-500 shrink-0"></span>
          <span>{{ errorMessage }}</span>
        </div>

        <!-- Sign In Button -->
        <button
          type="submit"
          :disabled="loading"
          class="w-full py-2.5 sm:py-3 bg-blue-600 hover:bg-blue-700 active:bg-blue-800 text-white font-bold rounded-xl shadow-md shadow-blue-600/30 transition-all flex items-center justify-center gap-2 disabled:opacity-60 cursor-pointer text-xs sm:text-sm tracking-wide mt-2"
        >
          <span v-if="loading" class="animate-spin w-4 h-4 border-2 border-white border-t-transparent rounded-full"></span>
          <span>{{ loading ? 'Memverifikasi...' : 'Masuk ke Sistem' }}</span>
        </button>
      </form>

      <!-- Quick Demo Login Switcher -->
      <div class="mt-6 pt-5 border-t border-slate-100">
        <div class="flex items-center justify-between mb-3">
          <span class="text-[11px] font-bold text-slate-400 uppercase tracking-wider flex items-center gap-1.5">
            <Sparkles class="w-3.5 h-3.5 text-blue-500" />
            Pilih Peran Demo Akun:
          </span>
          <span class="text-[10px] text-slate-400">Klik untuk isi otomatis</span>
        </div>

        <div class="grid grid-cols-2 sm:grid-cols-3 gap-2">
          <button
            type="button"
            @click="quickLogin('admin@cafe-erp.com', 'Admin@123')"
            class="flex items-center gap-2 p-2 rounded-xl border border-slate-200 bg-slate-50/70 hover:bg-blue-50 hover:border-blue-200 text-left transition-all cursor-pointer group"
          >
            <div class="w-7 h-7 rounded-lg bg-blue-100 text-blue-600 flex items-center justify-center shrink-0 group-hover:bg-blue-600 group-hover:text-white transition-colors">
              <Shield class="w-3.5 h-3.5" />
            </div>
            <div class="min-w-0">
              <div class="text-xs font-bold text-slate-800 group-hover:text-blue-600 truncate">Super Admin</div>
              <div class="text-[10px] text-slate-400 truncate">Akses Penuh ERP</div>
            </div>
          </button>

          <button
            type="button"
            @click="quickLogin('manager@cafe-erp.com', 'Password@123')"
            class="flex items-center gap-2 p-2 rounded-xl border border-slate-200 bg-slate-50/70 hover:bg-indigo-50 hover:border-indigo-200 text-left transition-all cursor-pointer group"
          >
            <div class="w-7 h-7 rounded-lg bg-indigo-100 text-indigo-600 flex items-center justify-center shrink-0 group-hover:bg-indigo-600 group-hover:text-white transition-colors">
              <UserCheck class="w-3.5 h-3.5" />
            </div>
            <div class="min-w-0">
              <div class="text-xs font-bold text-slate-800 group-hover:text-indigo-600 truncate">Store Manager</div>
              <div class="text-[10px] text-indigo-600 font-semibold truncate">Approval Cuti/Shift</div>
            </div>
          </button>

          <button
            type="button"
            @click="quickLogin('kasir@cafe-erp.com', 'Password@123')"
            class="flex items-center gap-2 p-2 rounded-xl border border-slate-200 bg-slate-50/70 hover:bg-emerald-50 hover:border-emerald-200 text-left transition-all cursor-pointer group"
          >
            <div class="w-7 h-7 rounded-lg bg-emerald-100 text-emerald-600 flex items-center justify-center shrink-0 group-hover:bg-emerald-600 group-hover:text-white transition-colors">
              <ShoppingBag class="w-3.5 h-3.5" />
            </div>
            <div class="min-w-0">
              <div class="text-xs font-bold text-slate-800 group-hover:text-emerald-600 truncate">Kasir POS</div>
              <div class="text-[10px] text-slate-400 truncate">Transaksi & Struk</div>
            </div>
          </button>

          <button
            type="button"
            @click="quickLogin('inventory@cafe-erp.com', 'Password@123')"
            class="flex items-center gap-2 p-2 rounded-xl border border-slate-200 bg-slate-50/70 hover:bg-amber-50 hover:border-amber-200 text-left transition-all cursor-pointer group"
          >
            <div class="w-7 h-7 rounded-lg bg-amber-100 text-amber-600 flex items-center justify-center shrink-0 group-hover:bg-amber-600 group-hover:text-white transition-colors">
              <Package class="w-3.5 h-3.5" />
            </div>
            <div class="min-w-0">
              <div class="text-xs font-bold text-slate-800 group-hover:text-amber-600 truncate">Staf Inventory</div>
              <div class="text-[10px] text-slate-400 truncate">Stok & P2P Gudang</div>
            </div>
          </button>

          <button
            type="button"
            @click="quickLogin('finance@cafe-erp.com', 'Password@123')"
            class="flex items-center gap-2 p-2 rounded-xl border border-slate-200 bg-slate-50/70 hover:bg-purple-50 hover:border-purple-200 text-left transition-all cursor-pointer group"
          >
            <div class="w-7 h-7 rounded-lg bg-purple-100 text-purple-600 flex items-center justify-center shrink-0 group-hover:bg-purple-600 group-hover:text-white transition-colors">
              <DollarSign class="w-3.5 h-3.5" />
            </div>
            <div class="min-w-0">
              <div class="text-xs font-bold text-slate-800 group-hover:text-purple-600 truncate">Finance</div>
              <div class="text-[10px] text-slate-400 truncate">Jurnal & Akuntansi</div>
            </div>
          </button>
        </div>
      </div>

    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import {
  Coffee,
  Mail,
  Lock,
  Eye,
  EyeOff,
  Sparkles,
  Shield,
  UserCheck,
  ShoppingBag,
  Package,
  DollarSign
} from 'lucide-vue-next'
import { useAuthStore } from '@/stores/auth.store'
import { useNotificationStore } from '@/stores/notification.store'
import axios from 'axios'

const router = useRouter()
const authStore = useAuthStore()
const notifyStore = useNotificationStore()

const email = ref('')
const password = ref('')
const rememberMe = ref(false)
const showPassword = ref(false)
const loading = ref(false)
const errorMessage = ref('')

const handleLogin = async () => {
  if (!email.value || !password.value) {
    errorMessage.value = 'Silakan isi email dan kata sandi.'
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

      notifyStore.success(`Selamat datang, ${data.user.full_name || 'Pengguna'}!`, `Peran: ${userObj.roleId}`)
      await router.push('/dashboard')
    } else {
      throw new Error('Respons otentikasi tidak valid dari server')
    }
  } catch (err: any) {
    const msg = err.response?.data?.message || err.message || 'Email atau kata sandi tidak valid.'
    errorMessage.value = msg
    notifyStore.error(msg, 'Gagal Masuk')
  } finally {
    loading.value = false
  }
}

const quickLogin = (e: string, p: string) => {
  email.value = e
  password.value = p
  handleLogin()
}

const handleForgotPassword = () => {
  notifyStore.info('Silakan hubungi administrator IT untuk mereset kata sandi akun Anda.', 'Bantuan Akun')
}
</script>
