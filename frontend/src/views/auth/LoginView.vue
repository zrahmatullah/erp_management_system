<template>
  <div class="min-h-screen w-full flex items-center justify-center bg-gradient-to-br from-[#0F172A] via-[#1E1B4B] to-[#3B0764] p-4 relative overflow-hidden">
    <!-- Ambient glow circles -->
    <div class="absolute -top-32 -left-32 w-96 h-96 bg-blue-600/20 rounded-full blur-3xl pointer-events-none"></div>
    <div class="absolute -bottom-32 -right-32 w-96 h-96 bg-purple-600/25 rounded-full blur-3xl pointer-events-none"></div>

    <!-- Decorative Coffee Line Art on right (matching mockup) -->
    <div class="hidden lg:block absolute right-12 bottom-0 opacity-15 pointer-events-none max-w-lg select-none">
      <svg viewBox="0 0 400 400" fill="none" stroke="white" stroke-width="1.5" class="w-full h-auto">
        <!-- Cup, beans, brewing contours -->
        <path d="M120 280 C120 330, 260 330, 260 280 L250 210 L130 210 Z" />
        <path d="M250 230 C280 230, 290 260, 260 270" />
        <path d="M160 180 C160 150, 180 150, 180 120" stroke-dasharray="3 3" />
        <path d="M200 180 C200 150, 220 150, 220 120" stroke-dasharray="3 3" />
        <ellipse cx="190" cy="210" rx="60" ry="12" />
        <ellipse cx="190" cy="330" rx="90" ry="16" />
        <!-- Coffee beans -->
        <ellipse cx="90" cy="300" rx="18" ry="12" transform="rotate(-30 90 300)" />
        <path d="M82 292 Q90 300 98 308" />
        <ellipse cx="290" cy="310" rx="20" ry="14" transform="rotate(25 290 310)" />
        <path d="M280 305 Q290 310 300 315" />
      </svg>
    </div>

    <!-- Centered Login Card -->
    <div class="w-full max-w-[440px] bg-white rounded-3xl shadow-2xl p-8 sm:p-10 border border-slate-100 relative z-10 animate-in fade-in zoom-in-95 duration-200 overflow-hidden">
      
      <!-- Loading Overlay -->
      <transition name="fade">
        <div 
          v-if="loading" 
          class="absolute inset-0 bg-white/92 backdrop-blur-xs rounded-3xl z-20 flex flex-col items-center justify-center p-6 text-center animate-in fade-in duration-200"
        >
          <div class="relative w-16 h-16 mb-4 flex items-center justify-center">
            <div class="absolute inset-0 rounded-full border-3 border-blue-100 border-t-blue-600 animate-spin"></div>
            <div class="w-10 h-10 rounded-xl bg-blue-50 text-blue-600 flex items-center justify-center animate-pulse">
              <Coffee class="w-5 h-5" />
            </div>
          </div>
          <h3 class="text-base font-bold text-slate-900 tracking-tight">Memproses Masuk...</h3>
          <p class="text-xs text-slate-500 mt-1 max-w-[220px]">Memverifikasi kredensial dan menyiapkan sesi dashboard</p>
        </div>
      </transition>

      <!-- Brand & Header -->
      <div class="flex items-center justify-center gap-3 mb-6">
        <div class="w-11 h-11 rounded-2xl bg-gradient-to-tr from-blue-600 to-indigo-600 flex items-center justify-center text-white shadow-md shadow-blue-500/20">
          <Coffee class="w-6 h-6" />
        </div>
        <span class="text-2xl font-bold tracking-tight text-slate-900">Cafe ERP</span>
      </div>

      <div class="text-center mb-8">
        <h1 class="text-2xl sm:text-3xl font-black text-slate-900 tracking-tight">Welcome Back</h1>
        <p class="text-sm text-slate-500 mt-1.5 font-medium">Sign in to your account</p>
      </div>

      <!-- Login Form -->
      <form @submit.prevent="handleLogin" class="space-y-4">
        <!-- Email Field -->
        <div>
          <label class="block text-xs font-semibold text-slate-700 mb-1.5">Email Address</label>
          <div class="relative">
            <Mail class="w-4 h-4 text-slate-400 absolute left-3.5 top-1/2 -translate-y-1/2 pointer-events-none" />
            <input
              v-model="email"
              type="text"
              required
              placeholder="Enter your email"
              class="w-full pl-10 pr-4 py-2.5 text-sm bg-slate-50 border border-slate-200 rounded-xl focus:bg-white focus:border-blue-600 focus:ring-3 focus:ring-blue-100 outline-hidden transition-all text-slate-800 font-medium placeholder:text-slate-400"
            />
          </div>
        </div>

        <!-- Password Field -->
        <div>
          <label class="block text-xs font-semibold text-slate-700 mb-1.5">Password</label>
          <div class="relative">
            <Lock class="w-4 h-4 text-slate-400 absolute left-3.5 top-1/2 -translate-y-1/2 pointer-events-none" />
            <input
              v-model="password"
              :type="showPassword ? 'text' : 'password'"
              required
              placeholder="Enter your password"
              class="w-full pl-10 pr-10 py-2.5 text-sm bg-slate-50 border border-slate-200 rounded-xl focus:bg-white focus:border-blue-600 focus:ring-3 focus:ring-blue-100 outline-hidden transition-all text-slate-800 font-medium placeholder:text-slate-400"
            />
            <button
              type="button"
              @click="showPassword = !showPassword"
              class="absolute right-3.5 top-1/2 -translate-y-1/2 text-slate-400 hover:text-slate-600 p-0.5"
            >
              <EyeOff v-if="showPassword" class="w-4 h-4" />
              <Eye v-else class="w-4 h-4" />
            </button>
          </div>
        </div>

        <!-- Remember me & Forgot Password -->
        <div class="flex items-center justify-between text-xs pt-1">
          <label class="flex items-center gap-2 cursor-pointer text-slate-600 select-none">
            <input 
              v-model="rememberMe" 
              type="checkbox" 
              class="w-4 h-4 rounded-md border-slate-300 text-blue-600 focus:ring-blue-500 cursor-pointer" 
            />
            <span class="font-medium">Remember me</span>
          </label>
          <a href="#" @click.prevent="handleForgotPassword" class="text-blue-600 font-semibold hover:underline">
            Forgot Password?
          </a>
        </div>

        <!-- Error Alert -->
        <div v-if="errorMessage" class="p-3 bg-red-50 border border-red-200 rounded-xl text-red-600 text-xs font-medium">
          {{ errorMessage }}
        </div>

        <!-- Sign In Button -->
        <button
          type="submit"
          :disabled="loading"
          class="w-full py-3 bg-blue-600 hover:bg-blue-700 active:bg-blue-800 text-white font-semibold rounded-xl shadow-md shadow-blue-600/25 transition-all flex items-center justify-center gap-2 disabled:opacity-60 cursor-pointer text-sm tracking-wide mt-2"
        >
          <span v-if="loading" class="animate-spin w-4 h-4 border-2 border-white border-t-transparent rounded-full"></span>
          <span>{{ loading ? 'Signing In...' : 'Sign In' }}</span>
        </button>
      </form>

      <!-- SSO Divider -->
      <div class="relative my-6 text-center">
        <div class="absolute inset-0 flex items-center">
          <div class="w-full border-t border-slate-200"></div>
        </div>
        <span class="relative bg-white px-3 text-xs text-slate-400 font-medium">Or sign in with SSO</span>
      </div>

      <!-- Quick Demo Login Switcher -->
      <div class="mb-5">
        <div class="text-[11px] font-semibold text-slate-400 uppercase tracking-wider text-center mb-2 flex items-center justify-center gap-1">
          <Sparkles class="w-3 h-3 text-amber-500" />
          1-Click Demo Login
        </div>
        <div class="grid grid-cols-2 gap-2">
          <button
            type="button"
            @click="quickLogin('admin@cafe-erp.com', 'Admin@123', 'Super Admin')"
            class="px-3 py-2 text-xs bg-slate-50 hover:bg-blue-50 hover:text-blue-600 text-slate-700 font-semibold rounded-xl border border-slate-200 transition-colors text-left flex items-center gap-1.5"
          >
            <Shield class="w-3.5 h-3.5 text-blue-600" />
            Super Admin
          </button>
          <button
            type="button"
            @click="quickLogin('cashier@cafe-erp.com', 'Admin@123', 'Kasir Jane')"
            class="px-3 py-2 text-xs bg-slate-50 hover:bg-blue-50 hover:text-blue-600 text-slate-700 font-semibold rounded-xl border border-slate-200 transition-colors text-left flex items-center gap-1.5"
          >
            <Coffee class="w-3.5 h-3.5 text-amber-600" />
            Kasir Jane
          </button>
        </div>
      </div>

      <!-- Language selector at bottom (matching mockup) -->
      <div class="text-center pt-1">
        <button class="inline-flex items-center gap-1.5 text-xs text-slate-500 hover:text-slate-800 bg-slate-50 px-3.5 py-1.5 rounded-lg border border-slate-200/70 font-medium transition-colors">
          <Languages class="w-3.5 h-3.5 text-slate-400" />
          <span>English (US)</span>
          <ChevronDown class="w-3.5 h-3.5 text-slate-400 ml-0.5" />
        </button>
      </div>

    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth.store'
import { useNotificationStore } from '@/stores/notification.store'
import axios from 'axios'
import { 
  Coffee, 
  Mail, 
  Lock, 
  Eye, 
  EyeOff, 
  Shield, 
  Languages, 
  ChevronDown, 
  Sparkles 
} from 'lucide-vue-next'

const router = useRouter()
const authStore = useAuthStore()
const notifyStore = useNotificationStore()

const email = ref('admin@cafe-erp.com')
const password = ref('Admin@123')
const rememberMe = ref(true)
const showPassword = ref(false)
const loading = ref(false)
const errorMessage = ref('')

const handleLogin = async () => {
  loading.value = true
  errorMessage.value = ''
  try {
    const res = await axios.post('/api/v1/auth/login', {
      email: email.value,
      password: password.value
    })

    if (res.data?.data) {
      const data = res.data.data
      authStore.login(
        data.token,
        data.token,
        {
          id: data.user.id,
          email: data.user.email,
          firstName: data.user.full_name || 'Admin',
          lastName: '',
          roleId: 'super_admin'
        },
        []
      )
      notifyStore.success(`Selamat datang kembali, ${data.user.full_name || 'Pengguna'}!`, 'Login Berhasil')
      setTimeout(() => {
        router.push('/dashboard')
      }, 350)
    } else {
      throw new Error('Respons otentikasi tidak valid dari server')
    }
  } catch (err: any) {
    const msg = err.response?.data?.message || err.message || 'Email atau kata sandi tidak valid.'
    errorMessage.value = msg
    notifyStore.error(msg, 'Gagal Masuk')
    loading.value = false
  }
}

const quickLogin = (e: string, p: string, _role: string) => {
  email.value = e
  password.value = p
  handleLogin()
}

const handleForgotPassword = () => {
  notifyStore.info('Silakan hubungi administrator IT untuk mereset kata sandi akun Anda.', 'Bantuan Akun')
}
</script>
