<template>
  <div class="min-h-screen w-full flex items-center justify-center bg-slate-100/90 dark:bg-slate-950 p-4 sm:p-6 lg:p-8 transition-colors duration-200 relative select-none">
    <!-- Ambient subtle background glow for dark mode -->
    <div class="absolute -top-40 -left-40 w-96 h-96 bg-indigo-500/10 dark:bg-indigo-500/15 rounded-full blur-3xl pointer-events-none"></div>
    <div class="absolute -bottom-40 -right-40 w-96 h-96 bg-purple-500/10 dark:bg-purple-500/15 rounded-full blur-3xl pointer-events-none"></div>

    <!-- Floating Theme Switcher Button (Top Right) -->
    <div class="absolute top-4 right-4 sm:top-6 sm:right-6 z-30">
      <button 
        @click="appStore.toggleTheme"
        class="flex items-center gap-2 px-3 py-1.5 rounded-full border border-slate-200 dark:border-slate-800 bg-white/90 dark:bg-slate-900/90 backdrop-blur-xs text-xs font-bold text-slate-700 dark:text-slate-200 hover:border-indigo-500/50 shadow-xs transition-all cursor-pointer"
        :title="appStore.theme === 'dark' ? 'Ganti ke Mode Terang' : 'Ganti ke Mode Gelap'"
      >
        <Sun v-if="appStore.theme === 'dark'" class="w-4 h-4 text-amber-400" />
        <Moon v-else class="w-4 h-4 text-indigo-600" />
        <span class="text-[11px] hidden sm:inline">{{ appStore.theme === 'dark' ? 'Light' : 'Dark' }}</span>
      </button>
    </div>

    <!-- Main Split Login Container Card -->
    <div class="w-full max-w-4xl bg-white dark:bg-slate-900 rounded-3xl sm:rounded-[2.25rem] shadow-2xl shadow-indigo-950/10 dark:shadow-slate-950/70 border border-slate-200/80 dark:border-slate-800 grid grid-cols-1 md:grid-cols-2 overflow-hidden transition-all duration-300 z-10">
      
      <!-- Left Column: Vector Illustration (Matching Mockup) -->
      <div class="bg-gradient-to-br from-slate-50 via-indigo-50/30 to-purple-50/20 dark:from-slate-900 dark:via-slate-850 dark:to-indigo-950/30 p-8 sm:p-10 flex flex-col items-center justify-center relative border-b md:border-b-0 md:border-r border-slate-100 dark:border-slate-800/80">
        
        <!-- Illustration SVG Graphic -->
        <div class="w-full max-w-[320px] aspect-square relative flex items-center justify-center">
          <svg viewBox="0 0 400 400" class="w-full h-full drop-shadow-sm select-none" fill="none" xmlns="http://www.w3.org/2000/svg">
            <defs>
              <linearGradient id="phoneGrad" x1="0%" y1="0%" x2="100%" y2="100%">
                <stop offset="0%" stop-color="#334155" />
                <stop offset="100%" stop-color="#0F172A" />
              </linearGradient>
              <linearGradient id="cardGrad" x1="0%" y1="0%" x2="100%" y2="100%">
                <stop offset="0%" stop-color="#FBBF24" />
                <stop offset="100%" stop-color="#F59E0B" />
              </linearGradient>
              <linearGradient id="badgeGrad" x1="0%" y1="0%" x2="100%" y2="100%">
                <stop offset="0%" stop-color="#FB7185" />
                <stop offset="100%" stop-color="#F43F5E" />
              </linearGradient>
              <linearGradient id="leafGrad" x1="0%" y1="0%" x2="100%" y2="100%">
                <stop offset="0%" stop-color="#E879F9" />
                <stop offset="100%" stop-color="#C084FC" />
              </linearGradient>
              <linearGradient id="charShirt" x1="0%" y1="0%" x2="100%" y2="100%">
                <stop offset="0%" stop-color="#818CF8" />
                <stop offset="100%" stop-color="#6366F1" />
              </linearGradient>
            </defs>

            <!-- Decorative gear background in top center -->
            <g opacity="0.45" stroke="#CBD5E1" stroke-width="3" fill="none">
              <circle cx="210" cy="110" r="32" stroke-dasharray="8 6" />
              <circle cx="210" cy="110" r="14" />
            </g>

            <!-- Floating Padlock Badge (Upper Left) -->
            <g filter="drop-shadow(0 6px 12px rgba(244,63,94,0.25))">
              <circle cx="110" cy="130" r="32" fill="url(#badgeGrad)" />
              <!-- Lock Icon -->
              <rect x="99" y="126" width="22" height="18" rx="4" fill="white" />
              <path d="M104 126 V119 C104 114.5 107.5 111 112 111 C116.5 111 120 114.5 120 119 V126" stroke="white" stroke-width="3" stroke-linecap="round" fill="none" />
              <circle cx="110" cy="134" r="2" fill="#F43F5E" />
            </g>

            <!-- Background Decorative Leaf (Right) -->
            <path d="M260 280 C290 240, 320 170, 270 120 C270 170, 320 220, 290 280 Z" fill="url(#leafGrad)" opacity="0.75" />
            <path d="M280 320 C320 280, 340 210, 310 160 C300 210, 330 260, 300 320 Z" fill="#F472B6" opacity="0.6" />

            <!-- Giant Smartphone Frame (Tilted slightly) -->
            <g transform="rotate(4 230 240)">
              <!-- Phone body -->
              <rect x="160" y="90" width="145" height="235" rx="22" fill="#FFFFFF" stroke="#E2E8F0" stroke-width="4" filter="drop-shadow(0 15px 25px rgba(0,0,0,0.08))" />
              <!-- Screen area -->
              <rect x="168" y="100" width="129" height="215" rx="16" fill="#F8FAFC" />
              <!-- Camera notch -->
              <rect x="215" y="105" width="35" height="5" rx="2.5" fill="#CBD5E1" />

              <!-- Yellow/Amber Card Overlay on screen -->
              <g transform="translate(-10, 20)">
                <rect x="168" y="115" width="130" height="155" rx="14" fill="url(#cardGrad)" filter="drop-shadow(0 10px 18px rgba(245,158,11,0.3))" />
                
                <!-- User Avatar Circle on Card -->
                <circle cx="233" cy="155" r="18" fill="white" />
                <circle cx="233" cy="151" r="7" fill="#F59E0B" />
                <path d="M222 167 C222 161, 227 159, 233 159 C239 159, 244 161, 244 167" fill="#F59E0B" />

                <!-- Input bar 1 -->
                <rect x="183" y="185" width="100" height="14" rx="4" fill="white" />
                <!-- Input bar 2 -->
                <rect x="183" y="206" width="100" height="14" rx="4" fill="white" />
                <!-- Submit mini bar -->
                <rect x="198" y="232" width="70" height="16" rx="5" fill="#D97706" opacity="0.75" />
              </g>
            </g>

            <!-- Character Standing on Left (Matching Mockup) -->
            <!-- Head -->
            <circle cx="108" cy="205" r="11" fill="#FBBF24" />
            <!-- Hair -->
            <path d="M102 201 C102 195, 114 195, 116 198 C118 202, 116 205, 116 205" stroke="#1E293B" stroke-width="4" stroke-linecap="round" fill="none" />
            <!-- Body / Purple Shirt -->
            <path d="M96 220 C96 216, 120 216, 120 220 L124 268 L92 268 Z" fill="url(#charShirt)" />
            <!-- Left Arm bent to hip -->
            <path d="M96 222 L82 242 L94 250" stroke="url(#charShirt)" stroke-width="7" stroke-linecap="round" stroke-linejoin="round" fill="none" />
            <!-- Right Arm extended towards phone -->
            <path d="M120 222 L138 238 L146 235" stroke="url(#charShirt)" stroke-width="7" stroke-linecap="round" stroke-linejoin="round" fill="none" />
            <!-- Belt -->
            <rect x="92" y="266" width="32" height="5" fill="#1E1B4B" />
            <!-- Trousers / Legs -->
            <path d="M94 271 L92 340 L104 340 L107 288 L111 288 L114 340 L126 340 L122 271 Z" fill="#312E81" />
            <!-- Shoes -->
            <ellipse cx="96" cy="342" rx="10" ry="3.5" fill="#0F172A" />
            <ellipse cx="122" cy="342" rx="10" ry="3.5" fill="#0F172A" />

            <!-- Subtle ground shadow -->
            <ellipse cx="200" cy="360" rx="120" ry="10" fill="#E2E8F0" opacity="0.5" />
          </svg>
        </div>

        <!-- Tagline under illustration -->
        <div class="mt-4 text-center">
          <div class="flex items-center justify-center gap-2 mb-1">
            <span class="w-6 h-0.5 rounded-full bg-indigo-500"></span>
            <span class="text-xs font-bold text-slate-800 dark:text-slate-200">Sistem ERP & POS Cafe</span>
            <span class="w-6 h-0.5 rounded-full bg-indigo-500"></span>
          </div>
          <p class="text-[11px] text-slate-400 dark:text-slate-500 max-w-xs leading-relaxed">
            Kelola pesanan kasir, inventori stok gudang, dan staf terintegrasi dalam satu platform.
          </p>
        </div>
      </div>

      <!-- Right Column: Login Form (Matching Mockup) -->
      <div class="p-7 sm:p-10 lg:p-11 flex flex-col justify-center bg-white dark:bg-slate-900">
        <!-- Greetings -->
        <div class="mb-5">
          <span class="text-xs font-medium text-slate-500 dark:text-slate-400">Welcome to</span>
          <h1 class="text-2xl sm:text-3xl font-black text-[#5B50E5] dark:text-indigo-400 tracking-tight mt-0.5">
            Cafe ERP
          </h1>
        </div>

        <!-- Quick Demo Profiles (Matching Social Login Pill Style) -->
        <div class="space-y-2 mb-4">
          <!-- Google Style Pill -> Super Admin -->
          <button
            type="button"
            @click="quickLogin('admin@cafe-erp.com', 'Admin@123')"
            class="w-full py-2.5 px-4 rounded-xl border border-slate-200/90 dark:border-slate-800 bg-white dark:bg-slate-850 hover:bg-slate-50 dark:hover:bg-slate-800 text-slate-700 dark:text-slate-200 text-xs font-bold transition-all shadow-2xs flex items-center justify-center gap-2.5 cursor-pointer group"
          >
            <div class="w-5 h-5 rounded-full bg-red-50 dark:bg-red-950/60 text-red-500 flex items-center justify-center font-bold text-[11px] shrink-0">
              G
            </div>
            <span>Masuk sebagai Super Admin</span>
          </button>

          <!-- Facebook Style Pill -> Kasir POS -->
          <button
            type="button"
            @click="quickLogin('kasir@cafe-erp.com', 'Password@123')"
            class="w-full py-2.5 px-4 rounded-xl border border-slate-200/90 dark:border-slate-800 bg-white dark:bg-slate-850 hover:bg-slate-50 dark:hover:bg-slate-800 text-slate-700 dark:text-slate-200 text-xs font-bold transition-all shadow-2xs flex items-center justify-center gap-2.5 cursor-pointer group"
          >
            <div class="w-5 h-5 rounded-full bg-blue-50 dark:bg-blue-950/60 text-blue-600 flex items-center justify-center font-bold text-[11px] shrink-0">
              f
            </div>
            <span>Masuk sebagai Kasir POS</span>
          </button>

          <!-- Toggle Other Roles Pill (Manager & Finance) -->
          <div class="grid grid-cols-2 gap-2 pt-0.5">
            <button
              type="button"
              @click="quickLogin('manager@cafe-erp.com', 'Password@123')"
              class="py-1.5 px-2.5 rounded-lg border border-slate-100 dark:border-slate-800 bg-slate-50 dark:bg-slate-800/60 hover:bg-indigo-50 dark:hover:bg-indigo-950/40 text-slate-600 dark:text-slate-300 text-[10.5px] font-semibold transition-all cursor-pointer truncate"
            >
              Demo: Store Manager
            </button>
            <button
              type="button"
              @click="quickLogin('finance@cafe-erp.com', 'Password@123')"
              class="py-1.5 px-2.5 rounded-lg border border-slate-100 dark:border-slate-800 bg-slate-50 dark:bg-slate-800/60 hover:bg-purple-50 dark:hover:bg-purple-950/40 text-slate-600 dark:text-slate-300 text-[10.5px] font-semibold transition-all cursor-pointer truncate"
            >
              Demo: Finance Staff
            </button>
          </div>
        </div>

        <!-- OR Divider (Matching Mockup) -->
        <div class="relative my-3 text-center">
          <div class="absolute inset-0 flex items-center">
            <div class="w-full border-t border-slate-200 dark:border-slate-800"></div>
          </div>
          <span class="relative px-3 bg-white dark:bg-slate-900 text-[11px] font-semibold text-slate-400 dark:text-slate-500 uppercase tracking-widest">
            Or
          </span>
        </div>

        <!-- Main Form (Matching Mockup Inputs) -->
        <form @submit.prevent="handleLogin" class="space-y-3.5">
          <!-- Email Field with envelope icon inside shaded pill container -->
          <div>
            <div class="bg-slate-100/90 dark:bg-slate-800 rounded-xl px-3.5 py-2.5 flex items-center gap-3 border border-transparent focus-within:border-[#5B50E5] focus-within:bg-white dark:focus-within:bg-slate-850 transition-all">
              <Mail class="w-4 h-4 text-slate-400 shrink-0" />
              <div class="flex-1 min-w-0">
                <span class="block text-[9.5px] font-semibold text-slate-400 dark:text-slate-500 leading-none mb-0.5">Email</span>
                <input
                  v-model="email"
                  type="email"
                  required
                  placeholder="example@gmail.com"
                  class="w-full bg-transparent p-0 text-xs text-slate-800 dark:text-slate-100 font-medium outline-none placeholder:text-slate-400 border-none ring-0 focus:ring-0"
                />
              </div>
            </div>
          </div>

          <!-- Password Field with key icon inside shaded pill container -->
          <div>
            <div class="bg-slate-100/90 dark:bg-slate-800 rounded-xl px-3.5 py-2.5 flex items-center gap-3 border border-transparent focus-within:border-[#5B50E5] focus-within:bg-white dark:focus-within:bg-slate-850 transition-all">
              <Key class="w-4 h-4 text-slate-400 shrink-0" />
              <div class="flex-1 min-w-0">
                <span class="block text-[9.5px] font-semibold text-slate-400 dark:text-slate-500 leading-none mb-0.5">Password</span>
                <input
                  v-model="password"
                  :type="showPassword ? 'text' : 'password'"
                  required
                  placeholder="••••••••••••"
                  class="w-full bg-transparent p-0 text-xs text-slate-800 dark:text-slate-100 font-medium outline-none placeholder:text-slate-400 border-none ring-0 focus:ring-0"
                />
              </div>
              <button
                type="button"
                @click="showPassword = !showPassword"
                class="text-slate-400 hover:text-slate-600 dark:hover:text-slate-300 p-0.5 cursor-pointer shrink-0 transition-colors"
                tabindex="-1"
              >
                <EyeOff v-if="showPassword" class="w-4 h-4" />
                <Eye v-else class="w-4 h-4" />
              </button>
            </div>
          </div>

          <!-- Remember me & Forgot Password -->
          <div class="flex items-center justify-between text-xs pt-0.5">
            <label class="flex items-center gap-2 cursor-pointer text-slate-500 dark:text-slate-400 select-none">
              <input 
                v-model="rememberMe" 
                type="checkbox" 
                class="w-3.5 h-3.5 rounded border-slate-300 dark:border-slate-700 text-[#5B50E5] focus:ring-[#5B50E5] cursor-pointer" 
              />
              <span class="text-[11px] font-medium">Remember me</span>
            </label>
            <a href="#" @click.prevent="handleForgotPassword" class="text-[#5B50E5] dark:text-indigo-400 font-semibold hover:underline text-[11px]">
              Forgot Password?
            </a>
          </div>

          <!-- Error Alert -->
          <div v-if="errorMessage" class="p-3 bg-red-50 dark:bg-red-950/40 border border-red-200 dark:border-red-900/50 rounded-xl text-red-600 dark:text-red-400 text-xs font-medium flex items-center gap-2">
            <span class="w-1.5 h-1.5 rounded-full bg-red-500 shrink-0"></span>
            <span>{{ errorMessage }}</span>
          </div>

          <!-- Submit Button (Vibrant Purple Matching Mockup) -->
          <button
            type="submit"
            :disabled="loading"
            class="w-full py-3 bg-[#5B50E5] hover:bg-[#4C40D4] active:bg-[#3F33C4] text-white font-bold rounded-xl shadow-lg shadow-indigo-600/30 transition-all flex items-center justify-center gap-2 disabled:opacity-60 cursor-pointer text-xs sm:text-sm tracking-wide mt-3"
          >
            <span v-if="loading" class="animate-spin w-4 h-4 border-2 border-white border-t-transparent rounded-full"></span>
            <span>{{ loading ? 'Signing in...' : 'Login' }}</span>
          </button>
        </form>

        <!-- Footer Register Link (Matching Mockup) -->
        <div class="mt-5 text-center text-xs text-slate-500 dark:text-slate-400">
          Don't have an account? 
          <a href="#" @click.prevent="handleRegisterClick" class="text-[#5B50E5] dark:text-indigo-400 font-bold hover:underline ml-0.5">
            Register
          </a>
        </div>

      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import {
  Mail,
  Key,
  Eye,
  EyeOff,
  Sun,
  Moon
} from 'lucide-vue-next'
import { useAuthStore } from '@/stores/auth.store'
import { useAppStore } from '@/stores/app.store'
import { useNotificationStore } from '@/stores/notification.store'
import axios from 'axios'

const router = useRouter()
const authStore = useAuthStore()
const appStore = useAppStore()
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

const handleRegisterClick = () => {
  notifyStore.info('Registrasi akun pengguna baru dikelola langsung oleh Super Admin melalui menu Master Data > Pengguna.', 'Pendaftaran Pengguna')
}
</script>
