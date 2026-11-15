<template>
  <div class="relative min-h-screen w-full flex flex-col justify-between overflow-x-hidden bg-[#FAFBFC] dark:bg-[#0B111E] transition-colors duration-300 select-none">
    
    <!-- Top-Right Ambient Purple/Violet Gradient Glow (matching user reference image) -->
    <div class="absolute -top-32 -right-32 w-[540px] h-[540px] bg-gradient-to-br from-purple-400/25 via-indigo-300/15 to-transparent rounded-full blur-[110px] pointer-events-none"></div>
    <div class="absolute -bottom-40 -left-40 w-[460px] h-[460px] bg-gradient-to-tr from-blue-400/12 via-sky-300/8 to-transparent rounded-full blur-[100px] pointer-events-none"></div>

    <!-- Top Floating Utility Bar (Non-intrusive, minimal) -->
    <header class="w-full max-w-[1280px] mx-auto px-4 sm:px-8 pt-4 sm:pt-6 flex items-center justify-between z-20">
      <!-- Active Outlet Selector Pill -->
      <button 
        type="button"
        @click="showStoreModal = true"
        class="group inline-flex items-center gap-2.5 px-3.5 py-1.5 rounded-full bg-white/80 dark:bg-slate-900/80 backdrop-blur-md border border-slate-200/80 dark:border-slate-800 shadow-xs hover:border-blue-400 dark:hover:border-blue-600 transition-all cursor-pointer text-xs font-semibold text-slate-700 dark:text-slate-200"
      >
        <span class="w-2 h-2 rounded-full bg-emerald-500 animate-pulse"></span>
        <Store class="w-3.5 h-3.5 text-[#4361EE] dark:text-blue-400" />
        <span class="truncate max-w-[170px] sm:max-w-none">{{ activeStore.name }}</span>
        <span class="text-[10px] text-slate-400 dark:text-slate-500">• {{ activeStore.latency }}</span>
        <ChevronDown class="w-3 h-3 text-slate-400 group-hover:translate-y-0.5 transition-transform" />
      </button>

      <!-- Right Controls: UI Version Switcher, Theme & IT Support -->
      <div class="flex items-center gap-2 sm:gap-3">
        <!-- UI Version Switcher (Instant Rollback to V1) -->
        <button
          type="button"
          @click="toggleUiVersion"
          class="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-full text-xs font-bold transition-all cursor-pointer shadow-xs border"
          :class="uiVersionStore.isV2 
            ? 'bg-blue-50 dark:bg-blue-950/60 text-[#4361EE] dark:text-blue-400 border-blue-200 dark:border-blue-800 hover:bg-blue-100' 
            : 'bg-slate-100 dark:bg-slate-800 text-slate-700 dark:text-slate-300 border-slate-300 dark:border-slate-700'"
          title="Beralih versi UI (V1 Classic vs V2 Modern)"
        >
          <Sparkles class="w-3.5 h-3.5 text-amber-500" />
          <span>{{ uiVersionStore.isV2 ? 'UI V2 (Modern)' : 'UI V1 (Classic)' }}</span>
          <span class="text-[10px] opacity-70 underline ml-0.5">Switch</span>
        </button>

        <!-- Theme Toggle -->
        <button 
          type="button" 
          @click="toggleTheme" 
          class="w-8 h-8 rounded-full bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 flex items-center justify-center text-slate-600 dark:text-slate-300 hover:text-[#4361EE] shadow-xs transition-transform active:scale-90 cursor-pointer"
          :title="appStore.theme === 'dark' ? 'Mode Terang' : 'Mode Gelap'"
        >
          <Sun v-if="appStore.theme === 'dark'" class="w-4 h-4 text-amber-400" />
          <Moon v-else class="w-4 h-4 text-[#4361EE]" />
        </button>

        <!-- IT Support -->
        <button
          type="button"
          @click="showForgotModal = true"
          class="hidden sm:inline-flex items-center gap-1 px-3 py-1.5 rounded-full text-xs font-semibold text-slate-500 hover:text-slate-800 dark:text-slate-400 dark:hover:text-white bg-white/60 dark:bg-slate-900/60 border border-slate-200/60 dark:border-slate-800 cursor-pointer transition-colors"
        >
          <Headphones class="w-3.5 h-3.5" />
          <span>Bantuan</span>
        </button>
      </div>
    </header>

    <!-- Main Container: 2-Column Split matching the uploaded reference image 1:1 -->
    <main class="w-full max-w-[1240px] mx-auto px-4 sm:px-6 py-6 sm:py-8 flex-1 flex items-center justify-center z-10">
      <div class="w-full grid grid-cols-1 lg:grid-cols-12 gap-8 lg:gap-12 items-center">
        
        <!-- ========================================================
             LEFT PANEL: Hero Card with Vibrant Blue & 3D Character
             100% faithful to the reference image
             ======================================================== -->
        <div class="lg:col-span-6 xl:col-span-6 order-2 lg:order-1">
          <div 
            class="relative rounded-[36px] sm:rounded-[42px] bg-gradient-to-b from-[#4F70FE] via-[#4361EE] to-[#3A53DC] text-white p-8 sm:p-11 lg:p-12 shadow-2xl shadow-blue-500/25 overflow-hidden flex flex-col justify-between min-h-[540px] lg:min-h-[640px] border border-white/20 select-none"
          >
            <!-- Background Radial Glows inside Card -->
            <div class="absolute -top-20 -right-20 w-72 h-72 rounded-full bg-white/15 blur-2xl pointer-events-none"></div>
            <div class="absolute -bottom-20 -left-20 w-80 h-80 rounded-full bg-blue-900/30 blur-2xl pointer-events-none"></div>

            <!-- Ambient Floating Micro Dots -->
            <div class="absolute top-1/4 right-8 w-2 h-2 rounded-full bg-white/40 animate-pulse pointer-events-none"></div>
            <div class="absolute top-1/2 left-8 w-1.5 h-1.5 rounded-full bg-white/30 animate-pulse pointer-events-none" style="animation-delay: 1.2s;"></div>
            
            <!-- Card Top Header / Slogan (Exact copy from the reference image) -->
            <div class="relative z-10">
              <h1 class="text-3xl sm:text-[36px] lg:text-[40px] font-black tracking-tight text-white leading-[1.15]">
                Simplify management with our dashboard
              </h1>
              <p class="text-sm sm:text-base text-blue-100/90 font-normal leading-relaxed mt-3.5 max-w-md">
                Simplify your cafe & resto management with our user friendly admin dashboard
              </p>
            </div>

            <!-- ===================================================
                 CENTER/BOTTOM: 3D Clay Character & Desk Illustration
                 Matching character with peace sign ✌️, red laptop,
                 dark monitor with yellow note, and sleek animations!
                 =================================================== -->
            <div class="relative z-10 mt-6 mb-2 flex items-center justify-center">
              <div class="relative w-full max-w-[430px] aspect-[4/3] flex items-center justify-center">
                
                <svg viewBox="0 0 500 400" class="w-full h-full drop-shadow-2xl overflow-visible select-none" xmlns="http://www.w3.org/2000/svg">
                  <defs>
                    <!-- Shading Gradients for 3D Clay Appearance -->
                    <linearGradient id="deskTopGrad" x1="0%" y1="0%" x2="0%" y2="100%">
                      <stop offset="0%" stop-color="#24304A" />
                      <stop offset="100%" stop-color="#141C2E" />
                    </linearGradient>
                    <linearGradient id="monitorBodyGrad" x1="0%" y1="0%" x2="100%" y2="100%">
                      <stop offset="0%" stop-color="#2A374F" />
                      <stop offset="100%" stop-color="#182030" />
                    </linearGradient>
                    <linearGradient id="laptopGrad" x1="0%" y1="0%" x2="100%" y2="100%">
                      <stop offset="0%" stop-color="#EF4444" />
                      <stop offset="100%" stop-color="#B91C1C" />
                    </linearGradient>
                    <linearGradient id="faceGrad" x1="0%" y1="0%" x2="0%" y2="100%">
                      <stop offset="0%" stop-color="#FFE0C8" />
                      <stop offset="100%" stop-color="#F3C3A2" />
                    </linearGradient>
                    <linearGradient id="sweaterGrad" x1="0%" y1="0%" x2="100%" y2="100%">
                      <stop offset="0%" stop-color="#FFFFFF" />
                      <stop offset="100%" stop-color="#E2E8F0" />
                    </linearGradient>
                    <linearGradient id="pantsGrad" x1="0%" y1="0%" x2="0%" y2="100%">
                      <stop offset="0%" stop-color="#2F66EE" />
                      <stop offset="100%" stop-color="#1E4CB8" />
                    </linearGradient>
                    <linearGradient id="hatGrad" x1="0%" y1="0%" x2="100%" y2="100%">
                      <stop offset="0%" stop-color="#C68953" />
                      <stop offset="100%" stop-color="#9A6032" />
                    </linearGradient>
                    <linearGradient id="chartGrad" x1="0%" y1="0%" x2="100%" y2="0%">
                      <stop offset="0%" stop-color="#38BDF8" />
                      <stop offset="50%" stop-color="#818CF8" />
                      <stop offset="100%" stop-color="#C084FC" />
                    </linearGradient>
                  </defs>

                  <!-- 1. Workstation Desk & Soft Ground Shadow -->
                  <g class="desk-group">
                    <!-- Ground Cast Shadow -->
                    <ellipse cx="250" cy="370" rx="195" ry="18" fill="#0E172F" opacity="0.45" />
                    
                    <!-- Desk Legs (Minimalist cylindrical legs) -->
                    <rect x="78" y="275" width="12" height="92" rx="6" fill="#1E2A44" />
                    <rect x="410" y="275" width="12" height="92" rx="6" fill="#1E2A44" />
                    <rect x="85" y="325" width="330" height="5" rx="2.5" fill="#2E3C59" opacity="0.5" />
                    
                    <!-- Main Desk Top Slab -->
                    <rect x="52" y="262" width="396" height="20" rx="10" fill="url(#deskTopGrad)" />
                    <!-- Subtle Desk Edge Bevel Highlight -->
                    <rect x="54" y="263" width="392" height="3" rx="1.5" fill="#38BDF8" opacity="0.35" />
                  </g>

                  <!-- 2. Animated Sitting 3D Character (Bobbing gently) -->
                  <g class="character-group animate-character-bob">
                    
                    <!-- Legs & Blue Jeans -->
                    <rect x="214" y="286" width="30" height="66" rx="15" fill="url(#pantsGrad)" />
                    <rect x="256" y="286" width="30" height="66" rx="15" fill="url(#pantsGrad)" />
                    <!-- Cute Clean Shoes -->
                    <rect x="210" y="342" width="34" height="17" rx="8.5" fill="#FFFFFF" />
                    <rect x="256" y="342" width="34" height="17" rx="8.5" fill="#FFFFFF" />

                    <!-- Torso / White Turtleneck Sweater -->
                    <ellipse cx="250" cy="226" rx="48" ry="54" fill="url(#sweaterGrad)" />
                    <rect x="237" y="166" width="26" height="22" rx="11" fill="#E2E8F0" />

                    <!-- Left Arm resting casually on desk -->
                    <path d="M 208 216 Q 184 242 178 264" stroke="#FFFFFF" stroke-width="22" stroke-linecap="round" fill="none" />
                    <!-- Left Hand resting on desk surface -->
                    <circle cx="178" cy="264" r="11" fill="url(#faceGrad)" />

                    <!-- Right Arm Waving Peace Sign ✌️ (Animated wave & energy rays) -->
                    <g class="peace-arm-group animate-hand-wave">
                      <path d="M 288 212 Q 324 180 342 138" stroke="#FFFFFF" stroke-width="22" stroke-linecap="round" fill="none" />
                      <!-- Right Hand Palm -->
                      <circle cx="342" cy="136" r="14" fill="url(#faceGrad)" />
                      <!-- Peace Sign Fingers (Index & Middle) -->
                      <rect x="334" y="104" width="7" height="26" rx="3.5" fill="url(#faceGrad)" transform="rotate(-15 334 104)" />
                      <rect x="347" y="108" width="7" height="26" rx="3.5" fill="url(#faceGrad)" transform="rotate(14 347 108)" />
                      
                      <!-- Radiant Motion Rays \ | / (Matching reference image exactly) -->
                      <g class="animate-rays-pulse">
                        <line x1="322" y1="92" x2="315" y2="80" stroke="#FFFFFF" stroke-width="3.5" stroke-linecap="round" />
                        <line x1="344" y1="86" x2="344" y2="72" stroke="#FFFFFF" stroke-width="3.5" stroke-linecap="round" />
                        <line x1="366" y1="94" x2="375" y2="84" stroke="#FFFFFF" stroke-width="3.5" stroke-linecap="round" />
                      </g>
                    </g>

                    <!-- Character Head & Friendly Clay Face -->
                    <circle cx="250" cy="148" r="38" fill="url(#faceGrad)" />
                    <!-- Hair underneath Hat -->
                    <path d="M 214 138 C 218 114 282 114 286 138 C 276 142 265 131 250 131 C 235 131 224 142 214 138 Z" fill="#221815" />
                    <!-- Eyes with Specular Sparkle -->
                    <circle cx="237" cy="148" r="4.5" fill="#1A2234" />
                    <circle cx="239" cy="146" r="1.5" fill="#FFFFFF" />
                    <circle cx="263" cy="148" r="4.5" fill="#1A2234" />
                    <circle cx="265" cy="146" r="1.5" fill="#FFFFFF" />
                    <!-- Rosy Pink Cheeks -->
                    <circle cx="228" cy="157" r="7" fill="#F472B6" opacity="0.5" />
                    <circle cx="272" cy="157" r="7" fill="#F472B6" opacity="0.5" />
                    <!-- Big Happy Smile -->
                    <path d="M 242 161 Q 250 171 258 161" stroke="#831843" stroke-width="3.2" stroke-linecap="round" fill="none" />

                    <!-- Stylized Hat / Beanie with Brim (Clay 3D styling) -->
                    <ellipse cx="250" cy="118" rx="38" ry="13" fill="url(#hatGrad)" />
                    <path d="M 216 118 C 220 88 280 88 284 118 Z" fill="url(#hatGrad)" />
                    <circle cx="250" cy="88" r="5" fill="#F1F5F9" />
                  </g>

                  <!-- 3. Ultrawide Dark Monitor on Desk (Left Side) -->
                  <g class="monitor-group">
                    <!-- Monitor Stand -->
                    <rect x="145" y="238" width="10" height="26" fill="#64748B" rx="2" />
                    <ellipse cx="150" cy="264" rx="24" ry="5" fill="#334155" />
                    
                    <!-- Monitor Frame -->
                    <rect x="84" y="142" width="132" height="100" rx="10" fill="url(#monitorBodyGrad)" stroke="#475569" stroke-width="2" />
                    <!-- Monitor Screen -->
                    <rect x="90" y="148" width="120" height="88" rx="6" fill="#0A0F1D" />
                    
                    <!-- Dashboard UI on screen: Top Bar & Metric -->
                    <rect x="96" y="155" width="32" height="4" rx="2" fill="#38BDF8" />
                    <circle cx="198" cy="157" r="2.5" fill="#10B981" />
                    <circle cx="191" cy="157" r="2.5" fill="#F59E0B" />
                    
                    <!-- Line Chart Wave -->
                    <path d="M 98 202 Q 115 174 135 188 T 172 168 T 202 182" fill="none" stroke="url(#chartGrad)" stroke-width="3" stroke-linecap="round" />
                    
                    <!-- Mini Bar Charts -->
                    <rect x="100" y="214" width="8" height="15" rx="2" fill="#38BDF8" opacity="0.8" />
                    <rect x="114" y="208" width="8" height="21" rx="2" fill="#818CF8" opacity="0.8" />
                    <rect x="128" y="216" width="8" height="13" rx="2" fill="#38BDF8" opacity="0.8" />
                    <rect x="142" y="204" width="8" height="25" rx="2" fill="#C084FC" opacity="0.8" />
                    
                    <!-- Yellow Sticky Note on Monitor Corner (Matches Reference Image) -->
                    <rect x="194" y="166" width="18" height="20" rx="2" fill="#FEF08A" transform="rotate(9 194 166)" />
                    <line x1="197" y1="172" x2="207" y2="173" stroke="#CA8A04" stroke-width="1.2" stroke-linecap="round" transform="rotate(9 194 166)" />
                    <line x1="197" y1="176" x2="205" y2="177" stroke="#CA8A04" stroke-width="1.2" stroke-linecap="round" transform="rotate(9 194 166)" />
                  </g>

                  <!-- 4. Coral/Red Open Laptop on Desk (Right Side, Matches Reference Image) -->
                  <g class="laptop-group">
                    <!-- Laptop Screen Outer Lid (Red) -->
                    <rect x="306" y="196" width="76" height="54" rx="6" fill="url(#laptopGrad)" stroke="#DC2626" stroke-width="1.5" transform="rotate(-6 306 196)" />
                    <!-- Laptop Screen Inner Display -->
                    <rect x="311" y="201" width="66" height="44" rx="4" fill="#0F172A" transform="rotate(-6 306 196)" />
                    <!-- Laptop Display UI -->
                    <rect x="318" y="208" width="32" height="3" rx="1.5" fill="#F87171" opacity="0.9" transform="rotate(-6 306 196)" />
                    <rect x="318" y="216" width="48" height="2.5" rx="1.2" fill="#94A3B8" opacity="0.8" transform="rotate(-6 306 196)" />
                    <rect x="318" y="222" width="40" height="2.5" rx="1.2" fill="#94A3B8" opacity="0.8" transform="rotate(-6 306 196)" />
                    <!-- Laptop Keyboard Base Plate on Desk -->
                    <polygon points="296,259 388,251 392,261 294,265" fill="#B91C1C" />
                  </g>

                  <!-- 5. Desk Accessories: Green Ball & Steaming Coffee Mug -->
                  <!-- Cute Green Ball/Fruit on Desk (Left, matching reference) -->
                  <circle cx="95" cy="260" r="7.5" fill="#10B981" />
                  <ellipse cx="94" cy="258" rx="2" ry="1.5" fill="#A7F3D0" />

                  <!-- Steaming Coffee Mug -->
                  <g class="coffee-cup-group">
                    <rect x="62" y="246" width="15" height="17" rx="3" fill="#FFFFFF" />
                    <!-- Mug Handle -->
                    <path d="M 62 250 C 56 250 56 258 62 258" stroke="#FFFFFF" stroke-width="2" fill="none" />
                    <!-- Hot Coffee Inside -->
                    <ellipse cx="69.5" cy="247" rx="6.5" ry="2" fill="#583115" />
                    <!-- Animated Rising Steam -->
                    <path class="animate-steam" d="M 67 242 Q 69 234 67 226" stroke="#FFFFFF" stroke-width="1.8" stroke-linecap="round" fill="none" opacity="0.75" />
                    <path class="animate-steam-delayed" d="M 72 242 Q 74 232 72 224" stroke="#FFFFFF" stroke-width="1.8" stroke-linecap="round" fill="none" opacity="0.75" />
                  </g>
                </svg>

              </div>
            </div>

            <!-- Card Bottom Bar: Subtle System Status -->
            <div class="relative z-10 pt-4 border-t border-white/15 flex items-center justify-between text-xs text-blue-100/80 font-medium">
              <div class="flex items-center gap-2">
                <span class="w-2 h-2 rounded-full bg-emerald-400 animate-pulse"></span>
                <span>Active Branch: Senopati Flagship</span>
              </div>
              <span>KopiERP Enterprise v2.4</span>
            </div>

          </div>
        </div>

        <!-- ========================================================
             RIGHT PANEL: Clean, Modern & Friendly Login Form
             100% matched to reference image layout & typography
             ======================================================== -->
        <div class="lg:col-span-6 xl:col-span-6 order-1 lg:order-2 flex flex-col justify-center">
          <div class="w-full max-w-[420px] mx-auto py-2 sm:py-6">
            
            <!-- Brand Logo: Centered with Iconic Droplet Ring matching Logoipsum -->
            <div class="flex items-center justify-center gap-2.5">
              <!-- Iconic Logoipsum Droplet SVG Glyph -->
              <svg viewBox="0 0 32 32" class="w-8 h-8 text-[#4361EE] fill-current" xmlns="http://www.w3.org/2000/svg">
                <circle cx="16" cy="16" r="13" stroke="currentColor" stroke-width="4" fill="none" />
                <circle cx="20" cy="12" r="4.5" fill="currentColor" />
              </svg>
              <span class="text-2xl font-black tracking-tight text-slate-900 dark:text-white">
                Logoipsum
              </span>
            </div>

            <!-- Header Titles: Centered matching screenshot -->
            <div class="mt-6 mb-7 text-center">
              <h2 class="text-3xl sm:text-[34px] font-black text-slate-900 dark:text-white tracking-tight">
                Welcome Back
              </h2>
              <p class="text-sm text-slate-400 dark:text-slate-400 mt-1.5 font-normal">
                Please Login to your account
              </p>
            </div>

            <!-- Lockout or Error Alert with Smooth Shake -->
            <transition name="fade">
              <div 
                v-if="isLockedOut" 
                class="p-3.5 mb-4 rounded-2xl border flex items-center gap-2.5 bg-amber-50 dark:bg-amber-950/40 border-amber-200 dark:border-amber-900 text-amber-800 dark:text-amber-300 text-xs font-semibold shadow-xs"
              >
                <AlertTriangle class="w-4 h-4 text-amber-600 dark:text-amber-400 shrink-0" />
                <span class="flex-1">Terminal dikunci sementara. Silakan tunggu <strong>{{ lockoutSeconds }}s</strong>.</span>
              </div>
              <div 
                v-else-if="errorMessage" 
                class="p-3.5 mb-4 rounded-2xl border flex items-center gap-2.5 bg-rose-50 dark:bg-rose-950/40 border-rose-200 dark:border-rose-900 text-rose-700 dark:text-rose-300 text-xs font-semibold shadow-xs shake-error"
              >
                <AlertCircle class="w-4 h-4 text-rose-500 shrink-0" />
                <span class="flex-1 leading-snug">{{ errorMessage }}</span>
              </div>
            </transition>

            <!-- Mode Switcher Pill (Sleek, Minimal) -->
            <div class="flex items-center justify-center mb-5">
              <div class="inline-flex bg-slate-100 dark:bg-slate-800/80 p-1 rounded-xl text-xs font-semibold border border-slate-200/50 dark:border-slate-700/50">
                <button
                  type="button"
                  @click="switchMode('email')"
                  class="px-3.5 py-1.5 rounded-lg transition-all cursor-pointer"
                  :class="authMode === 'email' 
                    ? 'bg-white dark:bg-slate-900 text-[#4361EE] dark:text-blue-400 font-bold shadow-xs' 
                    : 'text-slate-500 hover:text-slate-800 dark:hover:text-slate-200'"
                >
                  Email & Password
                </button>
                <button
                  type="button"
                  @click="switchMode('pin')"
                  class="px-3.5 py-1.5 rounded-lg transition-all cursor-pointer"
                  :class="authMode === 'pin' 
                    ? 'bg-white dark:bg-slate-900 text-[#4361EE] dark:text-blue-400 font-bold shadow-xs' 
                    : 'text-slate-500 hover:text-slate-800 dark:hover:text-slate-200'"
                >
                  PIN Kasir POS
                </button>
              </div>
            </div>

            <!-- ==========================================
                 VIEW A: Clean Email & Password Form
                 Matching user reference image 1:1
                 ========================================== -->
            <form v-if="authMode === 'email'" @submit.prevent="handleEmailLogin" class="space-y-4" novalidate>
              <!-- Email Input (Clean, Soft Rounded Pill) -->
              <div class="relative">
                <input
                  id="user-email"
                  v-model="email"
                  type="email"
                  autocomplete="email"
                  required
                  placeholder="Email"
                  class="w-full bg-[#F4F6F9] dark:bg-slate-800/80 border border-transparent focus:border-[#4361EE] focus:bg-white dark:focus:bg-slate-900 focus:ring-4 focus:ring-blue-500/10 rounded-2xl px-4 py-3.5 text-sm sm:text-base font-medium outline-none text-slate-900 dark:text-white placeholder:text-slate-400 dark:placeholder:text-slate-500 transition-all duration-200"
                />
              </div>

              <!-- Password Input with Interactive Eye Toggle (Clean, Soft Rounded Pill) -->
              <div class="relative flex items-center">
                <input
                  id="user-password"
                  v-model="password"
                  :type="showPassword ? 'text' : 'password'"
                  autocomplete="current-password"
                  required
                  placeholder="Password"
                  class="w-full bg-[#F4F6F9] dark:bg-slate-800/80 border border-transparent focus:border-[#4361EE] focus:bg-white dark:focus:bg-slate-900 focus:ring-4 focus:ring-blue-500/10 rounded-2xl px-4 py-3.5 pr-11 text-sm sm:text-base font-medium outline-none text-slate-900 dark:text-white placeholder:text-slate-400 dark:placeholder:text-slate-500 transition-all duration-200"
                />
                <button
                  type="button"
                  @click="showPassword = !showPassword"
                  class="absolute right-3.5 p-1 rounded-lg text-slate-400 hover:text-slate-700 dark:hover:text-slate-200 transition-colors cursor-pointer"
                  :title="showPassword ? 'Sembunyikan sandi' : 'Tampilkan sandi'"
                >
                  <EyeOff v-if="showPassword" class="w-5 h-5" />
                  <Eye v-else class="w-5 h-5" />
                </button>
              </div>

              <!-- Forgot Password Link (Right-aligned, matching screenshot) -->
              <div class="flex justify-end pt-0.5">
                <button 
                  type="button" 
                  @click="showForgotModal = true" 
                  class="font-normal text-slate-400 hover:text-[#4361EE] transition-colors cursor-pointer text-xs"
                >
                  Forgot Password
                </button>
              </div>

              <!-- Main Solid Blue CTA Button (Matching screenshot 1:1) -->
              <button 
                type="submit" 
                :disabled="loading || isLockedOut" 
                class="btn-shimmer w-full bg-[#4361EE] hover:bg-[#3752DE] active:bg-[#2F47D0] text-white font-bold py-3.5 sm:py-4 px-6 rounded-2xl text-sm sm:text-base flex items-center justify-center gap-2 shadow-lg shadow-blue-500/25 hover:shadow-xl hover:shadow-blue-500/35 hover:-translate-y-0.5 active:translate-y-0 active:scale-[0.99] transition-all duration-200 cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed mt-2"
              >
                <span v-if="loading" class="animate-spin w-5 h-5 border-2 border-white border-t-transparent rounded-full"></span>
                <span v-else>Login</span>
              </button>
            </form>

            <!-- ==========================================
                 VIEW B: Cashier 4-Digit PIN & Numeric Keypad
                 ========================================== -->
            <div v-else class="space-y-4">
              <!-- PIN Dots Indicator with Smooth Pulsing -->
              <div class="flex flex-col items-center justify-center py-2">
                <div class="text-xs font-semibold text-slate-500 dark:text-slate-400 mb-3">
                  Ketik 4-digit PIN Kasir / Barista
                </div>
                <div class="flex items-center gap-3.5">
                  <div 
                    v-for="i in 4" 
                    :key="i" 
                    class="w-4 h-4 rounded-full transition-all duration-200"
                    :class="pin.length >= i 
                      ? 'bg-[#4361EE] ring-4 ring-blue-100 dark:ring-blue-900/60 scale-125' 
                      : 'bg-slate-200 dark:bg-slate-700'"
                  ></div>
                </div>
              </div>

              <!-- 3x4 Touch Ergonomic Keypad -->
              <div class="grid grid-cols-3 gap-2.5">
                <button 
                  v-for="digit in ['1', '2', '3', '4', '5', '6', '7', '8', '9']" 
                  :key="digit" 
                  type="button" 
                  @click="appendDigit(digit)" 
                  class="h-12 rounded-2xl bg-white dark:bg-slate-800 hover:bg-slate-50 dark:hover:bg-slate-700 text-lg font-bold text-slate-800 dark:text-slate-100 flex items-center justify-center cursor-pointer transition active:scale-95 border border-slate-200/80 dark:border-slate-700 shadow-2xs"
                >
                  {{ digit }}
                </button>
                
                <!-- Biometric 1-Click Fast Touch -->
                <button 
                  type="button" 
                  @click="handleFingerprint" 
                  title="Simulasi Otorisasi Sidik Jari Kasir (1-Klik Cepat)" 
                  class="h-12 rounded-2xl bg-blue-50 dark:bg-blue-950/60 hover:bg-blue-100 dark:hover:bg-blue-900/70 text-[#4361EE] dark:text-blue-400 flex items-center justify-center cursor-pointer transition active:scale-95 border border-blue-200 dark:border-blue-900 shadow-2xs"
                >
                  <Fingerprint class="w-6 h-6" />
                </button>

                <!-- '0' Key -->
                <button 
                  type="button" 
                  @click="appendDigit('0')" 
                  class="h-12 rounded-2xl bg-white dark:bg-slate-800 hover:bg-slate-50 dark:hover:bg-slate-700 text-lg font-bold text-slate-800 dark:text-slate-100 flex items-center justify-center cursor-pointer transition active:scale-95 border border-slate-200/80 dark:border-slate-700 shadow-2xs"
                >
                  0
                </button>

                <!-- Delete Backspace Key -->
                <button 
                  type="button" 
                  @click="deleteDigit" 
                  title="Hapus Digit" 
                  class="h-12 rounded-2xl bg-slate-100 dark:bg-slate-800 hover:bg-slate-200 dark:hover:bg-slate-700 text-slate-600 dark:text-slate-300 flex items-center justify-center cursor-pointer transition active:scale-95 border border-slate-200 dark:border-slate-700 shadow-2xs"
                >
                  <Delete class="w-5 h-5" />
                </button>
              </div>

              <!-- Submit Button for PIN -->
              <button 
                type="button" 
                @click="submitPinLogin" 
                :disabled="loading || isLockedOut" 
                class="btn-shimmer w-full bg-[#4361EE] hover:bg-[#3752DE] text-white font-bold py-3 px-4 rounded-2xl text-xs sm:text-sm flex items-center justify-center gap-2 shadow-md shadow-blue-500/25 transition cursor-pointer disabled:opacity-50 mt-1"
              >
                <span v-if="loading" class="animate-spin w-4 h-4 border-2 border-white border-t-transparent rounded-full"></span>
                <Unlock v-else class="w-4 h-4" />
                <span>{{ loading ? 'Mengotentikasi...' : isLockedOut ? `Terkunci (${lockoutSeconds}s)` : 'Buka Shift & Masuk POS Kasir' }}</span>
              </button>
            </div>

            <!-- ========================================================
                 QUICK 1-CLICK DEMO ACCOUNTS BAR
                 Clean, subtle chips below the login form
                 ======================================================== -->
            <div class="mt-7 pt-4 border-t border-slate-200/70 dark:border-slate-800 text-center">
              <div class="text-[11px] font-bold text-slate-400 dark:text-slate-500 uppercase tracking-wider mb-2.5">
                Quick Demo Accounts (1-Klik Isi)
              </div>
              <div class="grid grid-cols-3 gap-1.5">
                <button
                  v-for="role in demoRoles"
                  :key="role.email"
                  type="button"
                  @click="selectDemoAccount(role.email, role.roleName)"
                  class="py-1.5 px-2 rounded-xl border flex items-center justify-center gap-1.5 text-xs font-semibold transition-all cursor-pointer hover:border-blue-400 dark:hover:border-blue-600"
                  :class="email === role.email && authMode === 'email'
                    ? 'bg-blue-50 dark:bg-blue-950/60 border-blue-400 text-[#4361EE] dark:text-blue-400 font-bold shadow-2xs' 
                    : 'text-slate-600 dark:text-slate-400 bg-white dark:bg-slate-900 border-slate-200 dark:border-slate-800'"
                >
                  <span class="text-xs">{{ role.icon }}</span>
                  <span class="truncate">{{ role.shortName }}</span>
                </button>
              </div>
            </div>

            <!-- Bottom Offline SQLite Mode Link -->
            <div class="mt-4 flex items-center justify-between text-xs text-slate-400 dark:text-slate-500">
              <button 
                type="button" 
                @click="showOfflineModal = true" 
                class="hover:text-[#4361EE] dark:hover:text-blue-400 transition-colors flex items-center gap-1 cursor-pointer"
              >
                <WifiOff class="w-3.5 h-3.5 text-slate-400" />
                <span>Mode Darurat Offline (SQLite)</span>
              </button>
              <button 
                type="button" 
                @click="showForgotModal = true" 
                class="hover:text-slate-700 dark:hover:text-slate-300 transition-colors cursor-pointer"
              >
                Bantuan IT
              </button>
            </div>

          </div>
        </div>

      </div>
    </main>

    <!-- Footer Copyright & Info -->
    <footer class="w-full max-w-[1240px] mx-auto px-4 sm:px-8 py-4 flex flex-col sm:flex-row items-center justify-between text-xs text-slate-400 dark:text-slate-500 gap-2 z-10">
      <div>&copy; 2026 Logoipsum / KopiERP Enterprise. All rights reserved.</div>
      <div class="flex items-center gap-3">
        <button type="button" @click="toggleLang" class="hover:text-[#4361EE] dark:hover:text-blue-400 font-semibold cursor-pointer">
          {{ lang === 'id' ? 'ID / Bahasa' : 'EN / English' }}
        </button>
        <span>•</span>
        <span>Zero-Downtime Multi-Tenant Cluster</span>
      </div>
    </footer>

    <!-- ========================================================
         MODALS (Outlet Selector, Offline Info, IT Help)
         ======================================================== -->

    <!-- Modal 1: Ganti Gerai / Switch Outlet -->
    <Teleport to="body">
      <transition name="fade">
        <div 
          v-if="showStoreModal" 
          class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-950/60 backdrop-blur-xs"
          @click.self="showStoreModal = false"
        >
          <div 
            class="w-full max-w-sm rounded-3xl p-6 shadow-2xl border bg-white dark:bg-slate-900 border-slate-200 dark:border-slate-800 animate-in zoom-in-95 duration-200"
            role="dialog"
            aria-modal="true"
          >
            <div class="flex items-center justify-between pb-3 border-b border-slate-100 dark:border-slate-800">
              <div class="flex items-center gap-2.5">
                <div class="w-8 h-8 rounded-xl bg-blue-50 dark:bg-blue-950 text-[#4361EE] dark:text-blue-400 flex items-center justify-center">
                  <Store class="w-4 h-4" />
                </div>
                <div>
                  <h3 class="font-bold text-sm text-slate-900 dark:text-white">Pilih Gerai Operasional</h3>
                  <p class="text-[11px] text-slate-400">Sinkronisasi terminal POS ke server outlet</p>
                </div>
              </div>
              <button @click="showStoreModal = false" class="p-1 rounded-lg text-slate-400 hover:text-slate-700 dark:hover:text-slate-200 cursor-pointer text-xs">
                ✕
              </button>
            </div>
            <div class="py-3.5 space-y-2">
              <button
                v-for="store in availableStores"
                :key="store.name"
                type="button"
                @click="selectStore(store)"
                class="w-full p-3 rounded-2xl border text-left flex items-center justify-between transition-all cursor-pointer"
                :class="activeStore.name === store.name 
                  ? 'border-blue-500 bg-blue-50/70 dark:bg-blue-950/40 text-[#4361EE] dark:text-blue-400' 
                  : 'border-slate-200 dark:border-slate-800 hover:bg-slate-50 dark:hover:bg-slate-800/60'"
              >
                <div>
                  <div class="text-xs font-bold text-slate-900 dark:text-white">{{ store.name }}</div>
                  <div class="text-[10px] text-slate-400 mt-0.5">Latensi {{ store.latency }} • Cluster {{ store.dbNode }}</div>
                </div>
                <span v-if="activeStore.name === store.name" class="text-xs font-extrabold text-[#4361EE] dark:text-blue-400 bg-white dark:bg-slate-800 px-2 py-0.5 rounded-full shadow-2xs">Aktif</span>
              </button>
            </div>
          </div>
        </div>
      </transition>
    </Teleport>

    <!-- Modal 2: Mode Darurat Offline SQLite Info -->
    <Teleport to="body">
      <transition name="fade">
        <div 
          v-if="showOfflineModal" 
          class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-950/60 backdrop-blur-xs"
          @click.self="showOfflineModal = false"
        >
          <div 
            class="w-full max-w-sm rounded-3xl p-6 shadow-2xl border bg-white dark:bg-slate-900 border-slate-200 dark:border-slate-800 animate-in zoom-in-95 duration-200"
            role="dialog"
            aria-modal="true"
          >
            <div class="flex items-center justify-between pb-3 border-b border-slate-100 dark:border-slate-800">
              <div class="flex items-center gap-2.5">
                <div class="w-8 h-8 rounded-xl bg-blue-50 dark:bg-blue-950 text-[#4361EE] dark:text-blue-400 flex items-center justify-center">
                  <WifiOff class="w-4 h-4" />
                </div>
                <h3 class="font-bold text-sm text-slate-900 dark:text-white">Sinkronisasi Offline (SQLite)</h3>
              </div>
              <button @click="showOfflineModal = false" class="p-1 rounded-lg text-slate-400 hover:text-slate-700 dark:hover:text-slate-200 cursor-pointer text-xs">
                ✕
              </button>
            </div>
            <div class="py-4 text-xs text-slate-600 dark:text-slate-300 leading-relaxed space-y-2.5">
              <p>
                Terminal POS dilengkapi dengan replikasi lokal <strong>SQLite Cache</strong>. Jika koneksi cloud terputus:
              </p>
              <ul class="list-disc list-inside space-y-1 text-[11px] text-slate-500 dark:text-slate-400">
                <li>Transaksi kasir tetap dapat berjalan tanpa jeda.</li>
                <li>Receipt thermal tetap dapat dicetak normal.</li>
                <li>Data transaksi otomatis tersinkronisasi kembali saat jaringan pulih.</li>
              </ul>
              <div class="p-2.5 rounded-xl bg-emerald-50 dark:bg-emerald-950/40 text-emerald-800 dark:text-emerald-300 text-[11px] font-semibold border border-emerald-200 dark:border-emerald-800 flex items-center gap-2">
                <CheckCircle2 class="w-4 h-4 text-emerald-600 dark:text-emerald-400 shrink-0" />
                <span>Replikasi lokal terkini: Aktif & Siap digunakan</span>
              </div>
            </div>
            <div class="pt-3 border-t border-slate-100 dark:border-slate-800 flex justify-end">
              <button 
                type="button" 
                @click="showOfflineModal = false" 
                class="px-4 py-2 rounded-xl text-xs font-bold text-white bg-[#4361EE] hover:bg-[#3752DE] cursor-pointer"
              >
                Tutup
              </button>
            </div>
          </div>
        </div>
      </transition>
    </Teleport>

    <!-- Modal 3: Bantuan IT / Reset Sandi -->
    <Teleport to="body">
      <transition name="fade">
        <div 
          v-if="showForgotModal" 
          class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-950/60 backdrop-blur-xs"
          @click.self="showForgotModal = false"
        >
          <div 
            class="w-full max-w-sm rounded-3xl p-6 shadow-2xl border bg-white dark:bg-slate-900 border-slate-200 dark:border-slate-800 animate-in zoom-in-95 duration-200"
            role="dialog"
            aria-modal="true"
          >
            <div class="flex items-center justify-between pb-3 border-b border-slate-100 dark:border-slate-800">
              <div class="flex items-center gap-2.5">
                <div class="w-8 h-8 rounded-xl bg-blue-50 dark:bg-blue-950 text-[#4361EE] dark:text-blue-400 flex items-center justify-center">
                  <Headphones class="w-4 h-4" />
                </div>
                <h3 class="font-bold text-sm text-slate-900 dark:text-white">Bantuan IT & Otorisasi</h3>
              </div>
              <button @click="showForgotModal = false" class="p-1 rounded-lg text-slate-400 hover:text-slate-700 dark:hover:text-slate-200 cursor-pointer text-xs">
                ✕
              </button>
            </div>
            <div class="py-4 text-xs text-slate-600 dark:text-slate-300 leading-relaxed space-y-2.5">
              <p>
                PIN atau kata sandi dapat di-reset melalui Super Administrator outlet Anda atau tim IT Senopati HQ.
              </p>
              <div class="p-3 rounded-2xl border bg-slate-50 dark:bg-slate-800 border-slate-200 dark:border-slate-700 text-[11px]">
                <div class="font-bold text-slate-900 dark:text-white mb-0.5">Hotline Dukungan POS Resto:</div>
                <div class="text-[#4361EE] dark:text-blue-400 font-extrabold text-xs">+62 812-9988-7766 (Senopati HQ)</div>
                <div class="text-slate-400 mt-1">Layanan teknis 24/7 untuk kasir, barista & manager cabang.</div>
              </div>
            </div>
            <div class="pt-3 border-t border-slate-100 dark:border-slate-800 flex justify-end">
              <button 
                type="button" 
                @click="showForgotModal = false" 
                class="px-4 py-2 rounded-xl text-xs font-bold text-white bg-[#4361EE] hover:bg-[#3752DE] cursor-pointer"
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
  Store,
  Fingerprint,
  Delete,
  Unlock,
  Headphones,
  WifiOff,
  ChevronDown,
  Eye,
  EyeOff,
  Sun,
  Moon,
  AlertCircle,
  AlertTriangle,
  Sparkles,
  CheckCircle2
} from 'lucide-vue-next'
import { useAuthStore } from '@/stores/auth.store'
import { useAppStore } from '@/stores/app.store'
import { useNotificationStore } from '@/stores/notification.store'
import { useUiVersionStore } from '@/stores/uiVersion.store'
import axios from 'axios'

const router = useRouter()
const authStore = useAuthStore()
const appStore = useAppStore()
const notifyStore = useNotificationStore()
const uiVersionStore = useUiVersionStore()

// Auth Modes: 'email' (Email & Password) | 'pin' (Cashier PIN Keypad)
const authMode = ref<'email' | 'pin'>('email')

// Email & Password state
const email = ref('admin@cafe-erp.com')
const password = ref('Admin@123')
const rememberMe = ref(true)
const showPassword = ref(false)

// Cashier PIN state
const pin = ref('')

// General UI & Submission state
const loading = ref(false)
const errorMessage = ref('')
const lang = ref<'id' | 'en'>('id')

// Modals
const showStoreModal = ref(false)
const showOfflineModal = ref(false)
const showForgotModal = ref(false)

// Rate limiting & lockout state
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

// 1-Click Demo Accounts for testing & evaluation
const demoRoles = [
  { icon: '👑', shortName: 'Admin', roleName: 'Super Admin', email: 'admin@cafe-erp.com' },
  { icon: '🏷️', shortName: 'Kasir', roleName: 'Kasir POS', email: 'kasir@cafe-erp.com' },
  { icon: '☕', shortName: 'Barista', roleName: 'Barista', email: 'barista@cafe-erp.com' },
  { icon: '📦', shortName: 'Gudang', roleName: 'Inventory', email: 'inventory@cafe-erp.com' },
  { icon: '💰', shortName: 'Finance', roleName: 'Finance', email: 'finance@cafe-erp.com' },
  { icon: '👥', shortName: 'HRIS', roleName: 'HR Admin', email: 'hris@cafe-erp.com' }
]

const selectDemoAccount = (selectedEmail: string, roleName: string) => {
  authMode.value = 'email'
  email.value = selectedEmail
  password.value = 'Admin@123'
  errorMessage.value = ''
  notifyStore.info(`Akun terpasang: ${roleName}`, 'Demo Account')
}

// Switch between Email and PIN modes
const switchMode = (mode: 'email' | 'pin') => {
  authMode.value = mode
  errorMessage.value = ''
  pin.value = ''
}

// Toggle UI Version between V1 (Classic) and V2 (Modern)
const toggleUiVersion = () => {
  uiVersionStore.toggleVersion()
  notifyStore.info(
    `Beralih ke tampilan ${uiVersionStore.isV2 ? 'UI V2 (Modern)' : 'UI V1 (Classic)'}`,
    'UI Switcher'
  )
}

// Theme Toggle
const toggleTheme = () => {
  const nextTheme = appStore.theme === 'dark' ? 'light' : 'dark'
  appStore.setTheme(nextTheme)
}

// Language Toggle
const toggleLang = () => {
  lang.value = lang.value === 'id' ? 'en' : 'id'
  notifyStore.info(`Bahasa diubah ke ${lang.value.toUpperCase()}`, 'Pengaturan Bahasa')
}

// Keyboard input for PIN keypad
const handleKeyDown = (e: KeyboardEvent) => {
  if (authMode.value !== 'pin') return
  if (e.key >= '0' && e.key <= '9') {
    appendDigit(e.key)
  } else if (e.key === 'Backspace') {
    deleteDigit()
  } else if (e.key === 'Enter') {
    submitPinLogin()
  }
}

// Append PIN digit
const appendDigit = (digit: string) => {
  if (isLockedOut.value) return
  errorMessage.value = ''
  if (pin.value.length < 4) {
    pin.value += digit
    if (pin.value.length === 4) {
      setTimeout(() => {
        submitPinLogin()
      }, 150)
    }
  }
}

// Delete PIN digit
const deleteDigit = () => {
  if (pin.value.length > 0) {
    pin.value = pin.value.slice(0, -1)
    errorMessage.value = ''
  }
}

// Biometric simulation
const handleFingerprint = () => {
  if (isLockedOut.value) return
  pin.value = '1234'
  notifyStore.success('Sidik Jari Kasir Terverifikasi: Jane Doe (#A1)', 'Biometrik Terhubung')
  setTimeout(() => {
    submitPinLogin()
  }, 250)
}

// Core API Authentication method (DRY with V1)
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
      notifyStore.success(`Selamat datang kembali, ${userObj.firstName}!`, `Sesi Terhubung`)
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
    const genericMsg = 'Email atau kata sandi tidak valid. Periksa kembali kredensial Anda.'
    errorMessage.value = genericMsg
    notifyStore.error(genericMsg, 'Otorisasi Gagal')
  } finally {
    loading.value = false
  }
}

// Handle Email & Password Submission
const handleEmailLogin = async () => {
  if (isLockedOut.value) return
  if (!email.value.trim() || !password.value) {
    errorMessage.value = 'Email dan kata sandi wajib diisi.'
    return
  }
  await authenticateWithCredentials(email.value.trim(), password.value, '/dashboard')
}

// Handle Cashier PIN Submission
const submitPinLogin = async () => {
  if (isLockedOut.value) return
  if (pin.value.length < 4) {
    errorMessage.value = 'Masukkan lengkap 4-digit PIN pegawai.'
    return
  }

  // Map PIN to role credentials
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

// Lockout timer if repeatedly failed
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

onMounted(() => {
  window.addEventListener('keydown', handleKeyDown)

  const savedEmail = localStorage.getItem(REMEMBER_KEY)
  if (savedEmail) {
    email.value = savedEmail
    rememberMe.value = true
  }
})

onUnmounted(() => {
  window.removeEventListener('keydown', handleKeyDown)
  if (lockoutTimer) clearInterval(lockoutTimer)
})
</script>

<style scoped>
/* Keyframes for Hero 3D Illustration & Micro-interactions */
@keyframes characterBob {
  0%, 100% {
    transform: translateY(0px);
  }
  50% {
    transform: translateY(-3.5px);
  }
}

@keyframes handWave {
  0%, 100% {
    transform: rotate(0deg);
    transform-origin: 324px 175px;
  }
  50% {
    transform: rotate(5deg);
    transform-origin: 324px 175px;
  }
}

@keyframes raysPulse {
  0%, 100% {
    opacity: 0.9;
    transform: scale(1);
    transform-origin: 344px 85px;
  }
  50% {
    opacity: 0.4;
    transform: scale(0.85);
    transform-origin: 344px 85px;
  }
}

@keyframes steamWave {
  0% {
    transform: translateY(0px) scaleX(1);
    opacity: 0.8;
  }
  50% {
    transform: translateY(-5px) scaleX(1.3);
    opacity: 0.4;
  }
  100% {
    transform: translateY(-11px) scaleX(1.6);
    opacity: 0;
  }
}

@keyframes shake {
  0%, 100% { transform: translateX(0); }
  20%, 60% { transform: translateX(-5px); }
  40%, 80% { transform: translateX(5px); }
}

.animate-character-bob {
  animation: characterBob 3.8s ease-in-out infinite;
}

.animate-hand-wave {
  animation: handWave 2.2s ease-in-out infinite;
}

.animate-rays-pulse {
  animation: raysPulse 1.8s ease-in-out infinite;
}

.animate-steam {
  animation: steamWave 2.4s ease-out infinite;
}

.animate-steam-delayed {
  animation: steamWave 2.4s ease-out infinite 1.2s;
}

.shake-error {
  animation: shake 0.4s ease-in-out;
}

/* Shimmer shine effect across primary Login button on hover */
.btn-shimmer {
  position: relative;
  overflow: hidden;
}

.btn-shimmer::after {
  content: '';
  position: absolute;
  top: -50%;
  left: -60%;
  width: 50%;
  height: 200%;
  background: linear-gradient(
    60deg,
    transparent,
    rgba(255, 255, 255, 0.28),
    transparent
  );
  transform: rotate(25deg);
  transition: all 0.65s ease;
}

.btn-shimmer:hover::after {
  left: 125%;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
