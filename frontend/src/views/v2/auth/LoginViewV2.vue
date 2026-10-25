<template>
  <div class="relative min-h-screen w-full flex flex-col justify-between overflow-x-hidden bg-[#FAFBFC] dark:bg-[#0B111E] transition-colors duration-300 select-none">
    
    <!-- Top-Right Ambient Gradient Glow (matching user reference image) -->
    <div class="absolute -top-32 -right-32 w-[520px] h-[520px] bg-gradient-to-br from-violet-500/18 via-indigo-500/12 to-transparent rounded-full blur-[100px] pointer-events-none"></div>
    <div class="absolute -bottom-40 -left-40 w-[450px] h-[450px] bg-gradient-to-tr from-blue-500/10 via-sky-400/8 to-transparent rounded-full blur-[90px] pointer-events-none"></div>

    <!-- Top Floating Navigation / Utility Bar -->
    <header class="w-full max-w-[1280px] mx-auto px-4 sm:px-6 pt-4 sm:pt-6 flex items-center justify-between z-20">
      <!-- Active Outlet Selector Pill -->
      <button 
        type="button"
        @click="showStoreModal = true"
        class="group inline-flex items-center gap-2.5 px-3.5 py-1.5 rounded-full bg-white/80 dark:bg-slate-900/80 backdrop-blur-md border border-slate-200/80 dark:border-slate-800 shadow-xs hover:border-blue-400 dark:hover:border-blue-600 transition-all cursor-pointer text-xs font-semibold text-slate-700 dark:text-slate-200"
      >
        <span class="w-2 h-2 rounded-full bg-emerald-500 animate-pulse"></span>
        <Store class="w-3.5 h-3.5 text-blue-600 dark:text-blue-400" />
        <span class="truncate max-w-[180px] sm:max-w-none">{{ activeStore.name }}</span>
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
            ? 'bg-blue-50 dark:bg-blue-950/60 text-blue-600 dark:text-blue-400 border-blue-200 dark:border-blue-800 hover:bg-blue-100' 
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
          class="w-8 h-8 rounded-full bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 flex items-center justify-center text-slate-600 dark:text-slate-300 hover:text-blue-600 shadow-xs transition-transform active:scale-90 cursor-pointer"
          :title="appStore.theme === 'dark' ? 'Mode Terang' : 'Mode Gelap'"
        >
          <Sun v-if="appStore.theme === 'dark'" class="w-4 h-4 text-amber-400" />
          <Moon v-else class="w-4 h-4 text-blue-500" />
        </button>

        <!-- IT Support / Offline -->
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

    <!-- Main Container: 2-Column Split matching the uploaded reference image -->
    <main class="w-full max-w-[1240px] mx-auto px-4 sm:px-6 py-6 sm:py-8 flex-1 flex items-center justify-center z-10">
      <div class="w-full grid grid-cols-1 lg:grid-cols-12 gap-8 xl:gap-12 items-center">
        
        <!-- ========================================================
             LEFT PANEL: Hero Card with Vibrant Blue Gradient & 3D Scene
             ======================================================== -->
        <div class="lg:col-span-6 xl:col-span-6 order-2 lg:order-1">
          <div 
            class="relative rounded-[32px] sm:rounded-[36px] bg-gradient-to-b from-[#4F6FF8] via-[#435FE9] to-[#364ED9] text-white p-7 sm:p-10 lg:p-12 shadow-2xl shadow-blue-500/20 overflow-hidden flex flex-col justify-between min-h-[520px] lg:min-h-[640px] border border-white/20"
          >
            <!-- Background Radial Glows & Rings inside Card -->
            <div class="absolute -top-24 -right-24 w-80 h-80 rounded-full bg-white/10 blur-2xl pointer-events-none"></div>
            <div class="absolute bottom-0 left-0 w-full h-48 bg-gradient-to-t from-[#273bb3]/40 to-transparent pointer-events-none"></div>
            
            <!-- Card Top Header / Slogan matching user image -->
            <div class="relative z-10">
              <h1 class="text-3xl sm:text-[34px] lg:text-[38px] font-extrabold tracking-tight text-white leading-[1.18]">
                Simplify management with our dashboard
              </h1>
              <p class="text-sm sm:text-base text-blue-100/90 font-medium leading-relaxed mt-3.5 max-w-md">
                Simplify your cafe & resto management with our user friendly admin dashboard and real-time cloud operations.
              </p>

              <!-- Floating Micro Stats / Capabilities -->
              <div class="flex flex-wrap items-center gap-2 mt-4 pt-1">
                <span class="inline-flex items-center gap-1.5 px-3 py-1 rounded-full bg-white/15 backdrop-blur-md border border-white/20 text-xs font-semibold text-white shadow-2xs">
                  <span class="w-1.5 h-1.5 rounded-full bg-emerald-400 animate-ping"></span>
                  99.9% Uptime
                </span>
                <span class="inline-flex items-center gap-1.5 px-3 py-1 rounded-full bg-white/15 backdrop-blur-md border border-white/20 text-xs font-semibold text-white shadow-2xs">
                  <Coffee class="w-3.5 h-3.5 text-amber-300" />
                  KDS & POS Realtime
                </span>
                <span class="inline-flex items-center gap-1.5 px-3 py-1 rounded-full bg-white/15 backdrop-blur-md border border-white/20 text-xs font-semibold text-white shadow-2xs">
                  <ShieldCheck class="w-3.5 h-3.5 text-sky-300" />
                  OWASP ASVS Ready
                </span>
              </div>
            </div>

            <!-- ===================================================
                 CENTER/BOTTOM: Exquisite 3D Barista & Desk Illustration
                 Constructed with high-fidelity vector & smooth CSS animation
                 =================================================== -->
            <div class="relative z-10 mt-8 mb-2 flex items-center justify-center">
              <div class="relative w-full max-w-[420px] aspect-[4/3] flex items-center justify-center illustration-container">
                
                <!-- Floating Floating Metric Badge 1 (Top Left) -->
                <div class="absolute top-2 left-2 sm:-left-2 z-20 animate-float-badge bg-white/95 dark:bg-slate-900/90 backdrop-blur-md rounded-2xl p-2.5 shadow-xl border border-white/40 flex items-center gap-2.5">
                  <div class="w-8 h-8 rounded-xl bg-emerald-500/15 text-emerald-600 dark:text-emerald-400 flex items-center justify-center shrink-0">
                    <TrendingUp class="w-4 h-4" />
                  </div>
                  <div>
                    <div class="text-[10px] text-slate-500 dark:text-slate-400 font-medium">Omset Hari Ini</div>
                    <div class="text-xs font-extrabold text-slate-900 dark:text-white">+28.4% Naik 🚀</div>
                  </div>
                </div>

                <!-- Floating Metric Badge 2 (Top Right) -->
                <div class="absolute top-6 right-2 sm:-right-2 z-20 animate-float-badge-delayed bg-white/95 dark:bg-slate-900/90 backdrop-blur-md rounded-2xl p-2.5 shadow-xl border border-white/40 flex items-center gap-2.5">
                  <div class="w-8 h-8 rounded-xl bg-blue-500/15 text-blue-600 dark:text-blue-400 flex items-center justify-center shrink-0">
                    <Receipt class="w-4 h-4" />
                  </div>
                  <div>
                    <div class="text-[10px] text-slate-500 dark:text-slate-400 font-medium">Order Selesai</div>
                    <div class="text-xs font-extrabold text-slate-900 dark:text-white">542 Transaksi ☕</div>
                  </div>
                </div>

                <!-- 3D Desk & Character SVG Illustration -->
                <svg viewBox="0 0 500 400" class="w-full h-full drop-shadow-2xl overflow-visible select-none" xmlns="http://www.w3.org/2000/svg">
                  <defs>
                    <!-- Linear and Radial Gradients for 3D Shading -->
                    <linearGradient id="deskGrad" x1="0%" y1="0%" x2="0%" y2="100%">
                      <stop offset="0%" stop-color="#1E2A4A" />
                      <stop offset="100%" stop-color="#0F172A" />
                    </linearGradient>
                    <linearGradient id="monitorGrad" x1="0%" y1="0%" x2="100%" y2="100%">
                      <stop offset="0%" stop-color="#1E293B" />
                      <stop offset="100%" stop-color="#0F172A" />
                    </linearGradient>
                    <linearGradient id="screenChartGrad" x1="0%" y1="0%" x2="100%" y2="0%">
                      <stop offset="0%" stop-color="#38BDF8" />
                      <stop offset="50%" stop-color="#818CF8" />
                      <stop offset="100%" stop-color="#C084FC" />
                    </linearGradient>
                    <linearGradient id="faceGrad" x1="0%" y1="0%" x2="0%" y2="100%">
                      <stop offset="0%" stop-color="#FFDFC4" />
                      <stop offset="100%" stop-color="#F3C3A2" />
                    </linearGradient>
                    <linearGradient id="sweaterGrad" x1="0%" y1="0%" x2="100%" y2="100%">
                      <stop offset="0%" stop-color="#FFFFFF" />
                      <stop offset="100%" stop-color="#E2E8F0" />
                    </linearGradient>
                    <linearGradient id="pantsGrad" x1="0%" y1="0%" x2="0%" y2="100%">
                      <stop offset="0%" stop-color="#2563EB" />
                      <stop offset="100%" stop-color="#1D4ED8" />
                    </linearGradient>
                    <linearGradient id="hatGrad" x1="0%" y1="0%" x2="100%" y2="100%">
                      <stop offset="0%" stop-color="#C28E61" />
                      <stop offset="100%" stop-color="#9C6B43" />
                    </linearGradient>
                    <filter id="glowFilter" x="-20%" y="-20%" width="140%" height="140%">
                      <feGaussianBlur stdDeviation="6" result="blur" />
                      <feComposite in="SourceGraphic" in2="blur" operator="over" />
                    </filter>
                  </defs>

                  <!-- Modern Workstation Desk Surface -->
                  <g class="desk-group">
                    <!-- Desk Shadow -->
                    <ellipse cx="250" cy="365" rx="200" ry="20" fill="#0D1630" opacity="0.4" />
                    <!-- Desk Legs -->
                    <rect x="75" y="270" width="12" height="95" rx="6" fill="#1E293B" />
                    <rect x="415" y="270" width="12" height="95" rx="6" fill="#1E293B" />
                    <rect x="80" y="320" width="340" height="6" rx="3" fill="#334155" opacity="0.6" />
                    <!-- Main Desk Top Plate -->
                    <rect x="50" y="260" width="400" height="20" rx="10" fill="url(#deskGrad)" />
                    <rect x="52" y="261" width="396" height="4" rx="2" fill="#38BDF8" opacity="0.3" />
                  </g>

                  <!-- Center 3D Character Sitting at the Desk (with gentle animated bobbing) -->
                  <g class="character-group animate-character-bob">
                    <!-- Legs & Pants -->
                    <rect x="215" y="285" width="28" height="65" rx="14" fill="url(#pantsGrad)" />
                    <rect x="257" y="285" width="28" height="65" rx="14" fill="url(#pantsGrad)" />
                    <!-- Shoes -->
                    <rect x="210" y="340" width="34" height="18" rx="9" fill="#FFFFFF" />
                    <rect x="256" y="340" width="34" height="18" rx="9" fill="#FFFFFF" />

                    <!-- Torso / White Turtleneck Sweater -->
                    <ellipse cx="250" cy="225" rx="46" ry="52" fill="url(#sweaterGrad)" />
                    <rect x="238" y="165" width="24" height="20" rx="10" fill="#E2E8F0" />

                    <!-- Left Arm on Keyboard -->
                    <path d="M 210 215 Q 185 240 180 262" stroke="#FFFFFF" stroke-width="22" stroke-linecap="round" fill="none" />
                    <!-- Left Hand resting on desk -->
                    <circle cx="180" cy="262" r="10" fill="url(#faceGrad)" />

                    <!-- Right Arm Waving Peace Sign ✌️ (matches user reference image) -->
                    <g class="animate-hand-wave">
                      <path d="M 285 210 Q 320 180 340 140" stroke="#FFFFFF" stroke-width="22" stroke-linecap="round" fill="none" />
                      <!-- Right Hand with Peace Sign ✌️ -->
                      <circle cx="340" cy="138" r="13" fill="url(#faceGrad)" />
                      <!-- Peace Fingers -->
                      <rect x="333" y="108" width="6" height="24" rx="3" fill="url(#faceGrad)" transform="rotate(-15 333 108)" />
                      <rect x="345" y="112" width="6" height="24" rx="3" fill="url(#faceGrad)" transform="rotate(12 345 112)" />
                      <!-- Fun Sparkle lines around the peace hand -->
                      <line x1="322" y1="95" x2="318" y2="85" stroke="#FFFFFF" stroke-width="3" stroke-linecap="round" />
                      <line x1="344" y1="90" x2="344" y2="78" stroke="#FFFFFF" stroke-width="3" stroke-linecap="round" />
                      <line x1="365" y1="98" x2="373" y2="90" stroke="#FFFFFF" stroke-width="3" stroke-linecap="round" />
                    </g>

                    <!-- Head & Smiling Face -->
                    <circle cx="250" cy="148" r="36" fill="url(#faceGrad)" />
                    <!-- Hair -->
                    <path d="M 216 138 C 220 115 280 115 284 138 C 275 142 265 132 250 132 C 235 132 225 142 216 138 Z" fill="#2D201A" />
                    <!-- Eyes with friendly gaze -->
                    <circle cx="238" cy="148" r="4.5" fill="#1E293B" />
                    <circle cx="240" cy="146" r="1.5" fill="#FFFFFF" />
                    <circle cx="262" cy="148" r="4.5" fill="#1E293B" />
                    <circle cx="264" cy="146" r="1.5" fill="#FFFFFF" />
                    <!-- Pink Cheeks -->
                    <circle cx="230" cy="156" r="6" fill="#F472B6" opacity="0.45" />
                    <circle cx="270" cy="156" r="6" fill="#F472B6" opacity="0.45" />
                    <!-- Big Warm Smile -->
                    <path d="M 242 160 Q 250 169 258 160" stroke="#831843" stroke-width="3" stroke-linecap="round" fill="none" />

                    <!-- Barista / Modern Hat / Beanie -->
                    <ellipse cx="250" cy="120" rx="36" ry="12" fill="url(#hatGrad)" />
                    <path d="M 218 120 C 222 92 278 92 282 120 Z" fill="url(#hatGrad)" />
                    <circle cx="250" cy="92" r="5" fill="#E2E8F0" />
                  </g>

                  <!-- Main Ultrawide Monitor on Desk (Left Side) -->
                  <g class="main-monitor">
                    <!-- Monitor Stand -->
                    <rect x="145" y="240" width="10" height="24" fill="#64748B" rx="2" />
                    <ellipse cx="150" cy="263" rx="24" ry="5" fill="#334155" />
                    <!-- Monitor Frame -->
                    <rect x="85" y="145" width="130" height="98" rx="10" fill="url(#monitorGrad)" stroke="#475569" stroke-width="2" />
                    <!-- Screen Display Area -->
                    <rect x="91" y="151" width="118" height="86" rx="6" fill="#0A0F1D" />
                    
                    <!-- Dashboard UI on screen: Top Bar -->
                    <rect x="96" y="157" width="30" height="4" rx="2" fill="#38BDF8" />
                    <circle cx="198" cy="159" r="2" fill="#10B981" />
                    <circle cx="192" cy="159" r="2" fill="#F59E0B" />
                    <!-- Revenue Analytics Wave Graph on Screen -->
                    <path d="M 98 205 Q 115 178 135 192 T 170 172 T 202 185" fill="none" stroke="url(#screenChartGrad)" stroke-width="3" stroke-linecap="round" />
                    <!-- Bar Chart on Screen -->
                    <rect x="100" y="215" width="8" height="14" rx="2" fill="#38BDF8" opacity="0.8" />
                    <rect x="114" y="210" width="8" height="19" rx="2" fill="#818CF8" opacity="0.8" />
                    <rect x="128" y="218" width="8" height="11" rx="2" fill="#38BDF8" opacity="0.8" />
                    <rect x="142" y="206" width="8" height="23" rx="2" fill="#C084FC" opacity="0.8" />
                    <!-- Sticky Note on Monitor Corner (matching reference image) -->
                    <rect x="195" y="170" width="16" height="18" rx="2" fill="#FEF08A" transform="rotate(8 195 170)" />
                  </g>

                  <!-- Secondary Screen: Laptop / Tablet (Right Side) -->
                  <g class="laptop-group">
                    <rect x="305" y="200" width="75" height="52" rx="6" fill="#1E293B" stroke="#475569" stroke-width="1.5" transform="rotate(-6 305 200)" />
                    <!-- Laptop Screen -->
                    <rect x="310" y="205" width="65" height="42" rx="4" fill="#0F172A" transform="rotate(-6 305 200)" />
                    <!-- Laptop Content -->
                    <rect x="316" y="212" width="30" height="3" rx="1.5" fill="#38BDF8" opacity="0.9" transform="rotate(-6 305 200)" />
                    <rect x="316" y="220" width="45" height="2.5" rx="1.2" fill="#64748B" transform="rotate(-6 305 200)" />
                    <rect x="316" y="226" width="38" height="2.5" rx="1.2" fill="#64748B" transform="rotate(-6 305 200)" />
                    <!-- Laptop Base on desk -->
                    <polygon points="295,258 385,252 388,260 292,264" fill="#334155" />
                  </g>

                  <!-- Hot Coffee Mug with Steaming Vapor -->
                  <g class="coffee-cup-group">
                    <rect x="62" y="244" width="16" height="18" rx="3" fill="#FFFFFF" />
                    <!-- Mug Handle -->
                    <path d="M 62 248 C 55 248 55 258 62 258" stroke="#FFFFFF" stroke-width="2.5" fill="none" />
                    <!-- Coffee Liquid inside -->
                    <ellipse cx="70" cy="245" rx="7" ry="2" fill="#583115" />
                    <!-- Rising Steam Waves (animated with CSS) -->
                    <path class="animate-steam" d="M 67 240 Q 69 232 67 225" stroke="#FFFFFF" stroke-width="2" stroke-linecap="round" fill="none" opacity="0.7" />
                    <path class="animate-steam-delayed" d="M 73 240 Q 75 230 73 222" stroke="#FFFFFF" stroke-width="2" stroke-linecap="round" fill="none" opacity="0.7" />
                  </g>
                </svg>

              </div>
            </div>

            <!-- Card Bottom Bar: Cloud Sync Info -->
            <div class="relative z-10 pt-4 border-t border-white/15 flex items-center justify-between text-xs text-blue-100/90 font-medium">
              <div class="flex items-center gap-2">
                <span class="w-2 h-2 rounded-full bg-emerald-400"></span>
                <span>Active Branch: Senopati Flagship</span>
              </div>
              <span>KopiERP Enterprise v2.4</span>
            </div>

          </div>
        </div>

        <!-- ========================================================
             RIGHT PANEL: Ultra Clean, Modern & Friendly Login Form
             Matching the uploaded reference layout & styling 1:1
             ======================================================== -->
        <div class="lg:col-span-6 xl:col-span-6 order-1 lg:order-2 flex flex-col justify-center">
          <div class="w-full max-w-[430px] mx-auto py-2 sm:py-6">
            
            <!-- Brand Logo with Modern Droplet Glyph matching Logoipsum -->
            <div class="flex items-center gap-3">
              <div class="w-10 h-10 rounded-2xl bg-gradient-to-tr from-[#3B54DE] to-[#5572FA] text-white flex items-center justify-center shadow-md shadow-blue-500/25">
                <svg viewBox="0 0 24 24" class="w-6 h-6 fill-current" xmlns="http://www.w3.org/2000/svg">
                  <path d="M12 2C6.477 2 2 6.477 2 12c0 4.991 3.657 9.128 8.438 9.878v-6.987h-2.54V12h2.54V9.797c0-2.506 1.492-3.89 3.777-3.89 1.094 0 2.238.195 2.238.195v2.46h-1.26c-1.243 0-1.63.771-1.63 1.562V12h2.773l-.443 2.891h-2.33v6.988C18.343 21.128 22 16.991 22 12c0-5.523-4.477-10-10-10z" class="hidden" />
                  <!-- Modern Logoipsum Style Ring/Droplet Glyph -->
                  <circle cx="12" cy="12" r="9" stroke="white" stroke-width="3" fill="none" opacity="0.9" />
                  <circle cx="15" cy="9" r="3.5" fill="white" />
                </svg>
              </div>
              <div class="flex items-center gap-1.5">
                <span class="text-2xl font-black tracking-tight text-slate-900 dark:text-white">KopiERP</span>
                <span class="text-[10px] font-extrabold uppercase px-2 py-0.5 rounded-md bg-blue-100 dark:bg-blue-900/60 text-blue-700 dark:text-blue-300 tracking-wider">
                  V2 PRO
                </span>
              </div>
            </div>

            <!-- Page Title & Subtitle matching the reference screenshot -->
            <div class="mt-6 mb-6">
              <h2 class="text-3xl sm:text-4xl font-black text-slate-900 dark:text-white tracking-tight">
                Welcome Back
              </h2>
              <p class="text-sm sm:text-base text-slate-400 dark:text-slate-400 mt-1.5 font-medium">
                Please Login to your account
              </p>
            </div>

            <!-- Authentication Mode Segmented Switch (Email vs POS Cashier PIN) -->
            <div class="bg-slate-100 dark:bg-slate-800/90 p-1 rounded-2xl flex items-center gap-1 mb-5 border border-slate-200/60 dark:border-slate-700/60 shadow-2xs">
              <button 
                type="button"
                @click="switchMode('email')" 
                class="flex-1 py-2 rounded-xl text-xs font-bold flex items-center justify-center gap-2 transition-all cursor-pointer"
                :class="authMode === 'email' 
                  ? 'bg-white dark:bg-slate-900 text-slate-900 dark:text-white shadow-xs' 
                  : 'text-slate-500 hover:text-slate-800 dark:hover:text-slate-200'"
              >
                <Mail class="w-3.5 h-3.5 text-blue-600 dark:text-blue-400" />
                <span>Email & Password</span>
              </button>
              <button 
                type="button"
                @click="switchMode('pin')" 
                class="flex-1 py-2 rounded-xl text-xs font-bold flex items-center justify-center gap-2 transition-all cursor-pointer"
                :class="authMode === 'pin' 
                  ? 'bg-white dark:bg-slate-900 text-slate-900 dark:text-white shadow-xs' 
                  : 'text-slate-500 hover:text-slate-800 dark:hover:text-slate-200'"
              >
                <Monitor class="w-3.5 h-3.5 text-blue-600 dark:text-blue-400" />
                <span>PIN Kasir (Fast POS)</span>
              </button>
            </div>

            <!-- Lockout or Error Alert Message with Shake Animation -->
            <transition name="fade">
              <div 
                v-if="isLockedOut" 
                class="p-3 mb-4 rounded-2xl border flex items-center gap-2.5 bg-amber-50 dark:bg-amber-950/40 border-amber-200 dark:border-amber-900 text-amber-800 dark:text-amber-300 text-xs font-semibold shadow-xs"
              >
                <AlertTriangle class="w-4 h-4 text-amber-600 dark:text-amber-400 shrink-0" />
                <span class="flex-1">Terminal dikunci sementara. Silakan tunggu <strong>{{ lockoutSeconds }}s</strong>.</span>
              </div>
              <div 
                v-else-if="errorMessage" 
                class="p-3 mb-4 rounded-2xl border flex items-center gap-2.5 bg-rose-50 dark:bg-rose-950/40 border-rose-200 dark:border-rose-900 text-rose-700 dark:text-rose-300 text-xs font-semibold shadow-xs shake-error"
              >
                <AlertCircle class="w-4 h-4 text-rose-500 shrink-0" />
                <span class="flex-1 leading-snug">{{ errorMessage }}</span>
              </div>
            </transition>

            <!-- ==========================================
                 VIEW A: Modern Email & Password Login Form
                 ========================================== -->
            <form v-if="authMode === 'email'" @submit.prevent="handleEmailLogin" class="space-y-4" novalidate>
              <!-- Email Input (Clean, Soft Rounded Container) -->
              <div class="space-y-1">
                <div 
                  class="group relative rounded-2xl bg-[#F3F5F9] dark:bg-slate-800/80 border border-transparent focus-within:border-[#435FE9] focus-within:bg-white dark:focus-within:bg-slate-900 focus-within:ring-4 focus-within:ring-blue-500/10 transition-all duration-200 px-4 py-3.5 flex items-center gap-3"
                >
                  <Mail class="w-5 h-5 text-slate-400 group-focus-within:text-[#435FE9] transition-colors shrink-0" />
                  <input
                    id="user-email"
                    v-model="email"
                    type="email"
                    autocomplete="email"
                    required
                    placeholder="Email"
                    class="w-full bg-transparent text-sm sm:text-base font-medium outline-none border-none ring-0 text-slate-900 dark:text-white placeholder:text-slate-400 dark:placeholder:text-slate-500"
                  />
                </div>
              </div>

              <!-- Password Input with Interactive Eye Toggle -->
              <div class="space-y-1">
                <div 
                  class="group relative rounded-2xl bg-[#F3F5F9] dark:bg-slate-800/80 border border-transparent focus-within:border-[#435FE9] focus-within:bg-white dark:focus-within:bg-slate-900 focus-within:ring-4 focus-within:ring-blue-500/10 transition-all duration-200 px-4 py-3.5 flex items-center gap-3"
                >
                  <Lock class="w-5 h-5 text-slate-400 group-focus-within:text-[#435FE9] transition-colors shrink-0" />
                  <input
                    id="user-password"
                    v-model="password"
                    :type="showPassword ? 'text' : 'password'"
                    autocomplete="current-password"
                    required
                    placeholder="Password"
                    class="w-full bg-transparent text-sm sm:text-base font-medium outline-none border-none ring-0 text-slate-900 dark:text-white placeholder:text-slate-400 dark:placeholder:text-slate-500"
                  />
                  <button
                    type="button"
                    @click="showPassword = !showPassword"
                    class="p-1 rounded-lg text-slate-400 hover:text-slate-700 dark:hover:text-slate-200 transition-colors cursor-pointer"
                    :title="showPassword ? 'Sembunyikan sandi' : 'Tampilkan sandi'"
                  >
                    <EyeOff v-if="showPassword" class="w-5 h-5" />
                    <Eye v-else class="w-5 h-5" />
                  </button>
                </div>
              </div>

              <!-- Remember Me & Forgot Password Link -->
              <div class="flex items-center justify-between text-xs sm:text-sm pt-1">
                <label class="flex items-center gap-2 cursor-pointer text-slate-500 dark:text-slate-400 hover:text-slate-800 dark:hover:text-slate-200">
                  <input 
                    v-model="rememberMe" 
                    type="checkbox" 
                    class="w-4 h-4 rounded text-[#435FE9] border-slate-300 dark:border-slate-700 focus:ring-[#435FE9]" 
                  />
                  <span class="font-medium text-xs">Remember Me</span>
                </label>
                <button 
                  type="button" 
                  @click="showForgotModal = true" 
                  class="font-semibold text-slate-400 hover:text-[#435FE9] transition-colors cursor-pointer text-xs"
                >
                  Forgot Password
                </button>
              </div>

              <!-- Main Solid Blue CTA Login Button matching reference -->
              <button 
                type="submit" 
                :disabled="loading || isLockedOut" 
                class="btn-shimmer w-full bg-[#435FE9] hover:bg-[#3752DE] active:bg-[#2F47D0] text-white font-bold py-3.5 sm:py-4 px-6 rounded-2xl text-sm sm:text-base flex items-center justify-center gap-2 shadow-lg shadow-blue-500/25 hover:shadow-xl hover:shadow-blue-500/40 hover:-translate-y-0.5 active:translate-y-0 active:scale-[0.99] transition-all duration-200 cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed mt-2"
              >
                <span v-if="loading" class="animate-spin w-5 h-5 border-2 border-white border-t-transparent rounded-full"></span>
                <span v-else class="flex items-center gap-2">
                  <span>Login</span>
                  <ArrowRight class="w-4 h-4 group-hover:translate-x-1 transition-transform" />
                </span>
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
                      ? 'bg-[#435FE9] ring-4 ring-blue-100 dark:ring-blue-900/60 scale-125' 
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
                  class="h-12 rounded-2xl bg-blue-50 dark:bg-blue-950/60 hover:bg-blue-100 dark:hover:bg-blue-900/70 text-blue-600 dark:text-blue-400 flex items-center justify-center cursor-pointer transition active:scale-95 border border-blue-200 dark:border-blue-900 shadow-2xs"
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
                class="btn-shimmer w-full bg-[#435FE9] hover:bg-[#3752DE] text-white font-bold py-3 px-4 rounded-2xl text-xs sm:text-sm flex items-center justify-center gap-2 shadow-md shadow-blue-500/25 transition cursor-pointer disabled:opacity-50 mt-1"
              >
                <span v-if="loading" class="animate-spin w-4 h-4 border-2 border-white border-t-transparent rounded-full"></span>
                <Unlock v-else class="w-4 h-4" />
                <span>{{ loading ? 'Mengotentikasi...' : isLockedOut ? `Terkunci (${lockoutSeconds}s)` : 'Buka Shift & Masuk POS Kasir' }}</span>
              </button>
            </div>

            <!-- ========================================================
                 QUICK 1-CLICK DEMO ACCOUNTS BAR (Testing & Reviewer Friendly)
                 ======================================================== -->
            <div class="mt-8 pt-5 border-t border-slate-200/80 dark:border-slate-800">
              <div class="flex items-center justify-between text-[11px] font-bold text-slate-400 dark:text-slate-500 uppercase tracking-wider mb-2.5">
                <span class="flex items-center gap-1.5">
                  <Zap class="w-3.5 h-3.5 text-amber-500" />
                  Quick Demo Accounts
                </span>
                <span class="text-blue-600 dark:text-blue-400 lowercase font-medium">1-klik isi akun</span>
              </div>
              <div class="grid grid-cols-3 gap-1.5">
                <button
                  v-for="role in demoRoles"
                  :key="role.email"
                  type="button"
                  @click="selectDemoAccount(role.email, role.roleName)"
                  class="py-1.5 px-2 rounded-xl border flex items-center justify-center gap-1.5 text-xs font-semibold transition-all cursor-pointer hover:border-blue-400 dark:hover:border-blue-600"
                  :class="email === role.email && authMode === 'email'
                    ? 'bg-blue-50 dark:bg-blue-950/60 border-blue-400 text-blue-600 dark:text-blue-400 font-bold shadow-2xs' 
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
                class="hover:text-blue-600 dark:hover:text-blue-400 transition-colors flex items-center gap-1 cursor-pointer"
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

    <!-- Footer Copyright & Version -->
    <footer class="w-full max-w-[1240px] mx-auto px-4 sm:px-6 py-4 flex flex-col sm:flex-row items-center justify-between text-xs text-slate-400 dark:text-slate-500 gap-2 z-10">
      <div>&copy; 2026 KopiERP Enterprise. All rights reserved.</div>
      <div class="flex items-center gap-3">
        <button type="button" @click="toggleLang" class="hover:text-blue-600 dark:hover:text-blue-400 font-semibold cursor-pointer">
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
                <div class="w-8 h-8 rounded-xl bg-blue-50 dark:bg-blue-950 text-blue-600 dark:text-blue-400 flex items-center justify-center">
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
                  ? 'border-blue-500 bg-blue-50/70 dark:bg-blue-950/40 text-blue-600 dark:text-blue-400' 
                  : 'border-slate-200 dark:border-slate-800 hover:bg-slate-50 dark:hover:bg-slate-800/60'"
              >
                <div>
                  <div class="text-xs font-bold text-slate-900 dark:text-white">{{ store.name }}</div>
                  <div class="text-[10px] text-slate-400 mt-0.5">Latensi {{ store.latency }} • Cluster {{ store.dbNode }}</div>
                </div>
                <span v-if="activeStore.name === store.name" class="text-xs font-extrabold text-blue-600 dark:text-blue-400 bg-white dark:bg-slate-800 px-2 py-0.5 rounded-full shadow-2xs">Aktif</span>
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
                <div class="w-8 h-8 rounded-xl bg-blue-50 dark:bg-blue-950 text-blue-600 dark:text-blue-400 flex items-center justify-center">
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
                class="px-4 py-2 rounded-xl text-xs font-bold text-white bg-[#435FE9] hover:bg-[#3752DE] cursor-pointer"
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
                <div class="w-8 h-8 rounded-xl bg-blue-50 dark:bg-blue-950 text-blue-600 dark:text-blue-400 flex items-center justify-center">
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
                <div class="text-[#435FE9] dark:text-blue-400 font-extrabold text-xs">+62 812-9988-7766 (Senopati HQ)</div>
                <div class="text-slate-400 mt-1">Layanan teknis 24/7 untuk kasir, barista & manager cabang.</div>
              </div>
            </div>
            <div class="pt-3 border-t border-slate-100 dark:border-slate-800 flex justify-end">
              <button 
                type="button" 
                @click="showForgotModal = false" 
                class="px-4 py-2 rounded-xl text-xs font-bold text-white bg-[#435FE9] hover:bg-[#3752DE] cursor-pointer"
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
  Receipt,
  TrendingUp,
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
  ArrowRight,
  AlertCircle,
  AlertTriangle,
  Sparkles,
  Zap,
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
@keyframes floatBadge {
  0%, 100% {
    transform: translateY(0px) rotate(0deg);
  }
  50% {
    transform: translateY(-8px) rotate(1deg);
  }
}

@keyframes floatBadgeDelayed {
  0%, 100% {
    transform: translateY(0px) rotate(0deg);
  }
  50% {
    transform: translateY(-9px) rotate(-1deg);
  }
}

@keyframes characterBob {
  0%, 100% {
    transform: translateY(0px);
  }
  50% {
    transform: translateY(-3px);
  }
}

@keyframes handWave {
  0%, 100% {
    transform: rotate(0deg);
    transform-origin: 320px 170px;
  }
  50% {
    transform: rotate(4deg);
    transform-origin: 320px 170px;
  }
}

@keyframes steamWave {
  0% {
    transform: translateY(0px) scaleX(1);
    opacity: 0.8;
  }
  50% {
    transform: translateY(-6px) scaleX(1.3);
    opacity: 0.4;
  }
  100% {
    transform: translateY(-12px) scaleX(1.6);
    opacity: 0;
  }
}

@keyframes shake {
  0%, 100% { transform: translateX(0); }
  20%, 60% { transform: translateX(-5px); }
  40%, 80% { transform: translateX(5px); }
}

.animate-float-badge {
  animation: floatBadge 4.5s ease-in-out infinite;
}

.animate-float-badge-delayed {
  animation: floatBadgeDelayed 5s ease-in-out infinite 0.7s;
}

.animate-character-bob {
  animation: characterBob 4s ease-in-out infinite;
}

.animate-hand-wave {
  animation: handWave 2.2s ease-in-out infinite;
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
    rgba(255, 255, 255, 0.25),
    transparent
  );
  transform: rotate(25deg);
  transition: all 0.6s ease;
}

.btn-shimmer:hover::after {
  left: 120%;
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