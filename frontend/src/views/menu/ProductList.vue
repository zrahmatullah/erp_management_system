<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
      <div>
        <div class="flex items-center gap-2.5">
          <h2 class="text-2xl font-black text-slate-900 tracking-tight">Katalog Menu & Resep (BOM)</h2>
          <span class="px-2.5 py-0.5 rounded-full text-xs font-black bg-emerald-100 text-emerald-800 border border-emerald-200 shadow-xs">
            Sinkron Stok Gudang
          </span>
        </div>
        <p class="text-xs text-slate-500 mt-1">
          Sinkronisasi otomatis porsi menu POS dari stok bahan baku gudang, kalkulasi HPP otomatis, & pengajuan Draft PO saat stok menipis.
        </p>
      </div>

      <div class="flex items-center gap-2.5">
        <button
          @click="fetchProducts"
          class="p-2.5 rounded-xl border border-slate-200 bg-white hover:bg-slate-50 text-slate-600 transition-colors shadow-xs cursor-pointer"
          title="Segarkan Data"
        >
          <RotateCw class="w-4 h-4" :class="{ 'animate-spin': loading }" />
        </button>

        <button
          @click="$router.push('/menu/products/create')"
          class="px-4 py-2.5 bg-blue-600 hover:bg-blue-700 active:bg-blue-800 text-white text-xs font-bold rounded-xl shadow-md shadow-blue-600/30 transition-all cursor-pointer flex items-center gap-1.5"
        >
          <Plus class="w-4 h-4" />
          <span>Tambah Menu Baru</span>
        </button>
      </div>
    </div>

    <!-- Summary KPI Cards -->
    <div class="grid grid-cols-2 lg:grid-cols-4 gap-4">
      <!-- Total Menu -->
      <div class="bg-white rounded-2xl p-4 border border-slate-200/80 shadow-xs flex items-center gap-3.5">
        <div class="w-11 h-11 rounded-xl bg-blue-50 text-blue-600 flex items-center justify-center shrink-0">
          <Coffee class="w-5 h-5" />
        </div>
        <div>
          <span class="text-[11px] font-bold text-slate-400 uppercase tracking-wider block">Total Menu POS</span>
          <span class="text-xl font-black text-slate-900">{{ products.length }}</span>
        </div>
      </div>

      <!-- Ready in Stock -->
      <div class="bg-white rounded-2xl p-4 border border-slate-200/80 shadow-xs flex items-center gap-3.5">
        <div class="w-11 h-11 rounded-xl bg-emerald-50 text-emerald-600 flex items-center justify-center shrink-0">
          <CheckCircle2 class="w-5 h-5" />
        </div>
        <div>
          <span class="text-[11px] font-bold text-slate-400 uppercase tracking-wider block">Menu Siap Jual</span>
          <span class="text-xl font-black text-emerald-600">{{ readyCount }}</span>
        </div>
      </div>

      <!-- Out of Stock / Depleted -->
      <div class="bg-white rounded-2xl p-4 border border-slate-200/80 shadow-xs flex items-center gap-3.5">
        <div class="w-11 h-11 rounded-xl bg-rose-50 text-rose-600 flex items-center justify-center shrink-0">
          <AlertTriangle class="w-5 h-5" />
        </div>
        <div>
          <span class="text-[11px] font-bold text-slate-400 uppercase tracking-wider block">Habis / Perlu PO</span>
          <span class="text-xl font-black text-rose-600">{{ outOfStockCount }}</span>
        </div>
      </div>

      <!-- BOM Recipe Mapped -->
      <div class="bg-white rounded-2xl p-4 border border-slate-200/80 shadow-xs flex items-center gap-3.5">
        <div class="w-11 h-11 rounded-xl bg-violet-50 text-violet-600 flex items-center justify-center shrink-0">
          <Layers class="w-5 h-5" />
        </div>
        <div>
          <span class="text-[11px] font-bold text-slate-400 uppercase tracking-wider block">Terpetakan Resep (BOM)</span>
          <span class="text-xl font-black text-violet-600">{{ mappedCount }} / {{ products.length }}</span>
        </div>
      </div>
    </div>

    <!-- Filter & Search Bar -->
    <div class="flex flex-col sm:flex-row items-center justify-between gap-3 bg-white p-3.5 rounded-2xl border border-slate-200/80 shadow-xs">
      <div class="relative w-full sm:w-80">
        <Search class="w-4 h-4 text-slate-400 absolute left-3 top-1/2 -translate-y-1/2" />
        <input
          v-model="searchQuery"
          type="text"
          placeholder="Cari menu, SKU, stasiun..."
          class="w-full pl-9 pr-4 py-2 text-xs bg-slate-50 border border-slate-200 rounded-xl focus:bg-white focus:border-blue-600 focus:ring-2 focus:ring-blue-100 outline-hidden transition-all"
        />
      </div>

      <div class="flex items-center gap-2 w-full sm:w-auto overflow-x-auto pb-1 sm:pb-0 text-xs">
        <button
          v-for="st in statusFilters"
          :key="st.value"
          @click="selectedStatusFilter = st.value"
          class="px-3 py-1.5 rounded-xl font-bold whitespace-nowrap transition-colors cursor-pointer"
          :class="selectedStatusFilter === st.value ? 'bg-blue-600 text-white shadow-xs' : 'bg-slate-100 text-slate-600 hover:bg-slate-200'"
        >
          {{ st.label }}
        </button>
      </div>
    </div>

    <!-- Main Products Table -->
    <div class="bg-white rounded-2xl border border-slate-200/80 shadow-xs overflow-hidden">
      <div class="overflow-x-auto">
        <table class="w-full text-left text-xs text-slate-600">
          <thead class="bg-slate-50 border-b border-slate-200/80 text-[11px] uppercase font-bold text-slate-500 tracking-wider">
            <tr>
              <th class="py-3 px-4">Menu Produk</th>
              <th class="py-3 px-4">Kategori / Station</th>
              <th class="py-3 px-4 text-right">Harga Jual</th>
              <th class="py-3 px-4 text-right">HPP Bahan (COGS)</th>
              <th class="py-3 px-4 text-center">Margin Kotor</th>
              <th class="py-3 px-4 text-center">Stok Gudang (Porsi)</th>
              <th class="py-3 px-4 text-center">Status Resep (BOM)</th>
              <th class="py-3 px-4 text-right">Aksi & Integrasi</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100">
            <tr v-if="loading">
              <td colspan="8" class="py-12 text-center text-slate-400">
                <RotateCw class="w-6 h-6 animate-spin mx-auto mb-2 text-blue-600" />
                Memuat katalog menu & kalkulasi stok gudang...
              </td>
            </tr>
            <tr v-else-if="filteredProducts.length === 0">
              <td colspan="8" class="py-12 text-center text-slate-400">
                Tidak ada menu yang sesuai dengan filter pencarian.
              </td>
            </tr>
            <tr
              v-for="p in filteredProducts"
              :key="p.id"
              class="hover:bg-slate-50/70 transition-colors"
            >
              <!-- Menu & SKU -->
              <td class="py-3 px-4">
                <div class="flex items-center gap-3">
                  <div class="w-11 h-11 rounded-xl overflow-hidden bg-slate-100 shrink-0 border border-slate-200/60">
                    <img
                      :src="p.image_url || 'https://images.unsplash.com/photo-1514432324607-a09d9b4aefdd?w=200'"
                      :alt="p.name"
                      class="w-full h-full object-cover"
                    />
                  </div>
                  <div>
                    <div class="font-bold text-slate-900 text-xs flex items-center gap-1.5">
                      <span>{{ p.name }}</span>
                    </div>
                    <div class="text-[10px] font-mono text-slate-400 mt-0.5">{{ p.sku || 'SKU-NONE' }}</div>
                  </div>
                </div>
              </td>

              <!-- Kategori / Station -->
              <td class="py-3 px-4">
                <span class="inline-block px-2 py-0.5 rounded-lg text-[10px] font-black uppercase tracking-wider bg-slate-100 text-slate-700">
                  {{ p.target_station || p.category_name || 'Kitchen & Bar' }}
                </span>
              </td>

              <!-- Harga Jual -->
              <td class="py-3 px-4 text-right font-bold text-slate-900">
                Rp {{ formatRupiah(p.price || p.base_price) }}
              </td>

              <!-- HPP (COGS) -->
              <td class="py-3 px-4 text-right">
                <span v-if="p.has_recipe" class="font-bold text-slate-800">
                  Rp {{ formatRupiah(p.cogs) }}
                </span>
                <span v-else class="text-[11px] text-slate-400 italic">
                  Belum di-set
                </span>
              </td>

              <!-- Margin Kotor -->
              <td class="py-3 px-4 text-center">
                <div v-if="p.has_recipe" class="flex flex-col items-center">
                  <span
                    class="px-2 py-0.5 rounded-full text-[10px] font-black"
                    :class="getMarginBadgeClass(p.margin_percent)"
                  >
                    {{ Math.round(p.margin_percent) }}%
                  </span>
                  <span class="text-[10px] text-slate-400 mt-0.5">
                    Rp {{ formatRupiah(p.gross_profit) }}
                  </span>
                </div>
                <span v-else class="text-[11px] text-slate-400">-</span>
              </td>

              <!-- Stok Gudang (Porsi) -->
              <td class="py-3 px-4 text-center">
                <div class="flex flex-col items-center">
                  <span
                    class="px-2.5 py-1 rounded-xl text-xs font-black"
                    :class="(p.stock <= 0 || p.is_out_of_stock) 
                      ? 'bg-rose-100 text-rose-700 border border-rose-200' 
                      : (p.stock <= 5 
                        ? 'bg-amber-100 text-amber-700 border border-amber-200' 
                        : 'bg-emerald-50 text-emerald-700 border border-emerald-200')"
                  >
                    {{ p.stock }} Porsi
                  </span>
                  <span
                    v-if="p.bottleneck_ingredient && (p.stock <= 5 || p.is_out_of_stock)"
                    class="text-[9px] font-medium text-rose-500 max-w-[140px] truncate mt-1"
                    :title="p.bottleneck_ingredient"
                  >
                    Habis: {{ p.bottleneck_ingredient }}
                  </span>
                </div>
              </td>

              <!-- Status Resep (BOM) -->
              <td class="py-3 px-4 text-center">
                <button
                  @click="openRecipeModal(p)"
                  class="inline-flex items-center gap-1 px-2.5 py-1 rounded-xl text-[11px] font-bold transition-all cursor-pointer"
                  :class="p.has_recipe 
                    ? 'bg-violet-50 text-violet-700 border border-violet-200 hover:bg-violet-100' 
                    : 'bg-slate-100 text-slate-500 border border-dashed border-slate-300 hover:bg-slate-200'"
                >
                  <Layers class="w-3.5 h-3.5" />
                  <span>{{ p.has_recipe ? `${p.total_ingredients} Bahan Baku` : '+ Set Resep' }}</span>
                </button>
              </td>

              <!-- Aksi & Integrasi Draft PO -->
              <td class="py-3 px-4 text-right">
                <div class="flex items-center justify-end gap-1.5">
                  <!-- Minta PO (Draft) button - highlighted if stock is 0 -->
                  <button
                    v-if="p.has_recipe"
                    @click="openDraftPOModal(p)"
                    class="px-2.5 py-1.5 rounded-xl text-[11px] font-bold transition-all flex items-center gap-1 cursor-pointer shadow-xs"
                    :class="(p.stock <= 0 || p.is_out_of_stock)
                      ? 'bg-amber-500 hover:bg-amber-600 text-white font-black animate-pulse'
                      : 'bg-slate-100 hover:bg-amber-50 text-slate-700 hover:text-amber-700 border border-slate-200'"
                    :title="(p.stock <= 0) ? 'Stok habis! Ajukan PO ke Tim Gudang' : 'Ajukan Draft PO bahan baku'"
                  >
                    <ShoppingBag class="w-3.5 h-3.5" />
                    <span>Minta PO</span>
                  </button>

                  <!-- Resep Modal Button -->
                  <button
                    @click="openRecipeModal(p)"
                    class="p-1.5 rounded-xl border border-slate-200 bg-white hover:bg-slate-50 text-slate-600 hover:text-violet-600 transition-colors cursor-pointer"
                    title="Edit Mapping Resep (BOM)"
                  >
                    <ChefHat class="w-4 h-4" />
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- ========================================================================= -->
    <!-- MODAL 1: MAPPING RESEP & BAHAN BAKU (BOM)                                -->
    <!-- ========================================================================= -->
    <Teleport to="body">
      <div
        v-if="showRecipeModal"
        class="fixed md:left-64 inset-y-0 right-0 left-0 z-50 bg-slate-950/45 backdrop-blur-sm flex items-center justify-center p-4 overflow-y-auto"
      >
        <div class="bg-white rounded-2xl w-full max-w-4xl shadow-2xl border border-slate-200 overflow-hidden my-6">
          <!-- Modal Header -->
          <div class="p-5 border-b border-slate-100 flex items-center justify-between bg-slate-50/50">
            <div class="flex items-center gap-3">
              <div class="w-10 h-10 rounded-xl bg-violet-100 text-violet-700 flex items-center justify-center font-black shrink-0">
                <ChefHat class="w-5 h-5" />
              </div>
              <div>
                <div class="flex items-center gap-2">
                  <h3 class="font-black text-slate-900 text-base">Mapping Resep & Bahan Baku (BOM)</h3>
                  <span class="px-2 py-0.5 rounded text-[10px] font-black bg-blue-100 text-blue-800">
                    {{ selectedProductForRecipe?.name }}
                  </span>
                </div>
                <p class="text-xs text-slate-500 mt-0.5">
                  Tentukan takaran bahan baku per 1 porsi menu. Sistem menghitung HPP & stok porsi otomatis dari gudang.
                </p>
              </div>
            </div>
            <button
              @click="showRecipeModal = false"
              class="p-2 rounded-xl text-slate-400 hover:text-slate-700 hover:bg-slate-100 transition-colors cursor-pointer"
            >
              <X class="w-5 h-5" />
            </button>
          </div>

          <!-- Recipe Financial & Portions Live Indicator Bar -->
          <div class="grid grid-cols-2 sm:grid-cols-4 gap-3 p-4 bg-slate-50 border-b border-slate-200/80 text-xs">
            <div class="p-2.5 rounded-xl bg-white border border-slate-200 shadow-xs">
              <span class="text-[10px] font-bold text-slate-400 uppercase tracking-wider block">Harga Jual Menu</span>
              <span class="text-sm font-extrabold text-blue-600">
                Rp {{ formatRupiah(selectedProductForRecipe?.base_price || selectedProductForRecipe?.price) }}
              </span>
            </div>
            <div class="p-2.5 rounded-xl bg-white border border-slate-200 shadow-xs">
              <span class="text-[10px] font-bold text-slate-400 uppercase tracking-wider block">Total HPP (BOM)</span>
              <span class="text-sm font-extrabold text-slate-900">
                Rp {{ formatRupiah(recipeTotalCOGS) }}
              </span>
            </div>
            <div class="p-2.5 rounded-xl bg-white border border-slate-200 shadow-xs">
              <span class="text-[10px] font-bold text-slate-400 uppercase tracking-wider block">Estimasi Margin</span>
              <span class="text-sm font-extrabold" :class="recipeMarginPercent >= 50 ? 'text-emerald-600' : 'text-amber-600'">
                {{ Math.round(recipeMarginPercent) }}% (Rp {{ formatRupiah(recipeGrossProfit) }})
              </span>
            </div>
            <div class="p-2.5 rounded-xl bg-white border border-slate-200 shadow-xs">
              <span class="text-[10px] font-bold text-slate-400 uppercase tracking-wider block">Stok Porsi Cookable</span>
              <span
                class="text-sm font-black"
                :class="recipeCookablePortions <= 0 ? 'text-rose-600' : 'text-emerald-600'"
              >
                {{ recipeCookablePortions }} Porsi Tersedia
              </span>
            </div>
          </div>

          <!-- Recipe Items Builder Table -->
          <div class="p-5 space-y-4 max-h-[50vh] overflow-y-auto">
            <div class="flex items-center justify-between">
              <h4 class="text-xs font-black text-slate-900 uppercase tracking-wider">
                Daftar Bahan Baku Komposisi
              </h4>
              <button
                type="button"
                @click="addRecipeRow"
                class="px-3 py-1.5 bg-violet-50 hover:bg-violet-100 text-violet-700 text-xs font-bold rounded-xl border border-violet-200 transition-colors flex items-center gap-1 cursor-pointer"
              >
                <Plus class="w-3.5 h-3.5" />
                <span>Tambah Bahan Baku</span>
              </button>
            </div>

            <div v-if="loadingRecipe" class="py-12 text-center text-slate-400 text-xs">
              <RotateCw class="w-5 h-5 animate-spin mx-auto mb-2 text-violet-600" />
              Memuat data resep & katalog inventori...
            </div>

            <div v-else-if="recipeRows.length === 0" class="py-10 text-center text-slate-400 text-xs border border-dashed border-slate-200 rounded-xl">
              Belum ada bahan baku pada resep menu ini. Klik "Tambah Bahan Baku" di atas untuk memulai mapping komposisi.
            </div>

            <div v-else class="space-y-2.5">
              <div
                v-for="(row, idx) in recipeRows"
                :key="idx"
                class="p-3 rounded-xl border border-slate-200 bg-slate-50/40 hover:bg-slate-50 transition-colors flex flex-col md:flex-row md:items-center gap-3 text-xs"
              >
                <!-- Raw Material Select -->
                <div class="flex-1 min-w-[200px]">
                  <label class="block text-[10px] font-bold text-slate-500 mb-1">Bahan Baku Gudang</label>
                  <select
                    v-model="row.inventory_item_id"
                    @change="onIngredientSelect(row)"
                    class="w-full px-2.5 py-1.5 text-xs bg-white border border-slate-200 rounded-lg focus:border-blue-600 font-semibold text-slate-800 outline-none"
                  >
                    <option value="" disabled>Pilih bahan baku dari inventori...</option>
                    <option
                      v-for="inv in inventoryItemList"
                      :key="inv.id"
                      :value="inv.id"
                    >
                      {{ inv.name }} (Stok: {{ inv.current_stock }} {{ inv.uom }} | Rp {{ formatRupiah(inv.average_cost) }}/{{ inv.uom }})
                    </option>
                  </select>
                </div>

                <!-- Quantity Required -->
                <div class="w-28">
                  <label class="block text-[10px] font-bold text-slate-500 mb-1">Takaran / Porsi</label>
                  <input
                    v-model.number="row.quantity_required"
                    type="number"
                    step="0.001"
                    min="0.001"
                    class="w-full px-2.5 py-1.5 text-xs bg-white border border-slate-200 rounded-lg focus:border-blue-600 font-bold text-slate-900 outline-none text-right"
                  />
                </div>

                <!-- UOM -->
                <div class="w-20">
                  <label class="block text-[10px] font-bold text-slate-500 mb-1">Satuan</label>
                  <input
                    v-model="row.uom"
                    type="text"
                    placeholder="ml/gr/pcs"
                    class="w-full px-2 py-1.5 text-xs bg-slate-100 border border-slate-200 rounded-lg text-slate-600 font-mono text-center outline-none"
                  />
                </div>

                <!-- Live Subtotal Cost -->
                <div class="w-28 text-right">
                  <label class="block text-[10px] font-bold text-slate-500 mb-1">Subtotal Biaya</label>
                  <div class="py-1.5 font-extrabold text-slate-800">
                    Rp {{ formatRupiah((row.quantity_required || 0) * (row.unit_cost || 0)) }}
                  </div>
                </div>

                <!-- Cookable Portions Indicator -->
                <div class="w-24 text-center">
                  <label class="block text-[10px] font-bold text-slate-500 mb-1">Stok Porsi</label>
                  <div
                    class="py-1 px-1.5 rounded-md font-black text-[11px]"
                    :class="(row.can_make_portions ?? 999) <= 0 
                      ? 'bg-rose-100 text-rose-700' 
                      : 'bg-emerald-50 text-emerald-700'"
                  >
                    {{ row.can_make_portions ?? Math.floor((row.current_stock || 0) / (row.quantity_required || 1)) }} porsi
                  </div>
                </div>

                <!-- Instructions / Notes -->
                <div class="w-36">
                  <label class="block text-[10px] font-bold text-slate-500 mb-1">Catatan / Takaran</label>
                  <input
                    v-model="row.instructions"
                    type="text"
                    placeholder="Contoh: Double shot"
                    class="w-full px-2 py-1.5 text-[11px] bg-white border border-slate-200 rounded-lg outline-none text-slate-700"
                  />
                </div>

                <!-- Delete Action -->
                <div class="pt-4 flex items-center justify-end">
                  <button
                    type="button"
                    @click="removeRecipeRow(idx)"
                    class="p-1.5 rounded-lg text-rose-500 hover:bg-rose-50 hover:text-rose-700 transition-colors cursor-pointer"
                    title="Hapus baris bahan baku"
                  >
                    <Trash2 class="w-4 h-4" />
                  </button>
                </div>
              </div>
            </div>
          </div>

          <!-- Modal Footer -->
          <div class="p-4 border-t border-slate-100 bg-slate-50/70 flex items-center justify-between">
            <span class="text-[11px] text-slate-500">
              * Perubahan resep langsung memperbarui batas stok menu di kasir POS.
            </span>
            <div class="flex items-center gap-2">
              <button
                type="button"
                @click="showRecipeModal = false"
                class="px-4 py-2 rounded-xl border border-slate-200 bg-white hover:bg-slate-100 text-slate-700 text-xs font-bold transition-colors cursor-pointer"
              >
                Batal
              </button>
              <button
                type="button"
                @click="saveRecipe"
                :disabled="savingRecipe"
                class="px-5 py-2 bg-violet-600 hover:bg-violet-700 active:bg-violet-800 disabled:opacity-50 text-white text-xs font-bold rounded-xl shadow-md shadow-violet-600/30 transition-all cursor-pointer flex items-center gap-1.5"
              >
                <ChefHat class="w-4 h-4" />
                <span>{{ savingRecipe ? 'Menyimpan...' : 'Simpan Resep (BOM)' }}</span>
              </button>
            </div>
          </div>
        </div>
      </div>
    </Teleport>

    <!-- ========================================================================= -->
    <!-- MODAL 2: AJUKAN DRAFT PO DARI MENU (RESTOCK REQUISITION)                  -->
    <!-- ========================================================================= -->
    <Teleport to="body">
      <div
        v-if="showDraftPOModal"
        class="fixed md:left-64 inset-y-0 right-0 left-0 z-50 bg-slate-950/45 backdrop-blur-sm flex items-center justify-center p-4 overflow-y-auto"
      >
        <div class="bg-white rounded-2xl w-full max-w-2xl shadow-2xl border border-slate-200 overflow-hidden my-6">
          <!-- Modal Header -->
          <div class="p-5 border-b border-slate-100 flex items-center justify-between bg-amber-50/50">
            <div class="flex items-center gap-3">
              <div class="w-10 h-10 rounded-xl bg-amber-100 text-amber-700 flex items-center justify-center font-black shrink-0">
                <ShoppingBag class="w-5 h-5" />
              </div>
              <div>
                <h3 class="font-black text-slate-900 text-base">Ajukan Draft PO Bahan Baku</h3>
                <p class="text-xs text-slate-500 mt-0.5">
                  Permintaan pembelian otomatis ke Tim Gudang untuk pengadaan bahan baku menu yang menipis/habis.
                </p>
              </div>
            </div>
            <button
              @click="showDraftPOModal = false"
              class="p-2 rounded-xl text-slate-400 hover:text-slate-700 hover:bg-slate-100 transition-colors cursor-pointer"
            >
              <X class="w-5 h-5" />
            </button>
          </div>

          <div class="p-5 space-y-4">
            <!-- Product Target Card -->
            <div class="p-3.5 rounded-xl border border-slate-200 bg-slate-50/60 flex items-center gap-3">
              <div class="w-12 h-12 rounded-xl overflow-hidden bg-slate-200 shrink-0 border border-slate-200">
                <img
                  :src="selectedProductForPO?.image_url || 'https://images.unsplash.com/photo-1514432324607-a09d9b4aefdd?w=200'"
                  :alt="selectedProductForPO?.name"
                  class="w-full h-full object-cover"
                />
              </div>
              <div class="min-w-0 flex-1">
                <div class="text-xs font-bold text-slate-900">{{ selectedProductForPO?.name }}</div>
                <div class="text-[11px] text-slate-500">
                  Stok Saat Ini: 
                  <span class="font-bold text-rose-600">{{ selectedProductForPO?.stock || 0 }} Porsi</span>
                  <span v-if="selectedProductForPO?.bottleneck_ingredient" class="ml-1 text-slate-400">
                    ({{ selectedProductForPO?.bottleneck_ingredient }})
                  </span>
                </div>
              </div>
              <span class="px-2.5 py-1 rounded-full text-[10px] font-black uppercase tracking-wider bg-rose-100 text-rose-800">
                Stok Habis / Menipis
              </span>
            </div>

            <!-- Target Portions Selection -->
            <div>
              <label class="block text-xs font-bold text-slate-700 mb-1.5">
                Target Porsi yang Akan Disediakan <span class="text-rose-500">*</span>
              </label>
              <div class="flex items-center gap-2 mb-2">
                <button
                  v-for="amt in [25, 50, 100, 200]"
                  :key="amt"
                  type="button"
                  @click="draftPOForm.target_portions = amt"
                  class="px-3 py-1.5 rounded-xl text-xs font-bold transition-all cursor-pointer"
                  :class="draftPOForm.target_portions === amt 
                    ? 'bg-amber-500 text-white shadow-xs' 
                    : 'bg-slate-100 text-slate-700 hover:bg-slate-200'"
                >
                  {{ amt }} Porsi
                </button>
              </div>
              <input
                v-model.number="draftPOForm.target_portions"
                type="number"
                min="1"
                class="w-full px-3.5 py-2 text-xs bg-slate-50 border border-slate-200 rounded-xl focus:bg-white focus:border-amber-500 outline-none font-bold text-slate-900"
                placeholder="Atau ketik jumlah porsi kustom..."
              />
            </div>

            <!-- Department & Notes -->
            <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
              <div>
                <label class="block text-xs font-bold text-slate-700 mb-1">Departemen Pemohon</label>
                <input
                  v-model="draftPOForm.department"
                  type="text"
                  class="w-full px-3 py-2 text-xs bg-slate-50 border border-slate-200 rounded-xl focus:bg-white outline-none font-medium text-slate-800"
                />
              </div>
              <div>
                <label class="block text-xs font-bold text-slate-700 mb-1">Catatan Tambahan</label>
                <input
                  v-model="draftPOForm.notes"
                  type="text"
                  placeholder="Catatan prioritas restock..."
                  class="w-full px-3 py-2 text-xs bg-slate-50 border border-slate-200 rounded-xl focus:bg-white outline-none font-medium text-slate-800"
                />
              </div>
            </div>

            <!-- Explanation Callout -->
            <div class="p-3.5 rounded-xl bg-blue-50/70 border border-blue-200 text-xs text-blue-900 space-y-1">
              <div class="font-bold flex items-center gap-1.5">
                <CheckCircle2 class="w-4 h-4 text-blue-600" />
                <span>Alur Proses Setelah Pengajuan:</span>
              </div>
              <p class="text-[11px] text-blue-800 pl-5">
                1. Pengajuan ini langsung masuk ke <strong>Daftar PR (Purchase Requisition)</strong> di modul Pembelian.<br/>
                2. Tim Gudang akan mereview kebutuhan bahan baku, memilih supplier rekanan, dan menerbitkan <strong>PO Resmi</strong>.<br/>
                3. Begitu barang tiba dan dicatat via GRN, stok menu ini akan otomatis bertambah di kasir POS.
              </p>
            </div>
          </div>

          <!-- Modal Footer -->
          <div class="p-4 border-t border-slate-100 bg-slate-50/70 flex items-center justify-end gap-2">
            <button
              type="button"
              @click="showDraftPOModal = false"
              class="px-4 py-2 rounded-xl border border-slate-200 bg-white hover:bg-slate-100 text-slate-700 text-xs font-bold transition-colors cursor-pointer"
            >
              Batal
            </button>
            <button
              type="button"
              @click="submitDraftPO"
              :disabled="submittingPO"
              class="px-5 py-2 bg-amber-500 hover:bg-amber-600 active:bg-amber-700 disabled:opacity-50 text-white text-xs font-bold rounded-xl shadow-md shadow-amber-500/30 transition-all cursor-pointer flex items-center gap-1.5"
            >
              <ShoppingBag class="w-4 h-4" />
              <span>{{ submittingPO ? 'Mengajukan...' : 'Kirim Draft PR ke Tim Gudang' }}</span>
            </button>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import axios from 'axios'
import {
  RotateCw,
  Plus,
  Coffee,
  CheckCircle2,
  AlertTriangle,
  Layers,
  Search,
  ChefHat,
  Trash2,
  ShoppingBag,
  X
} from 'lucide-vue-next'
import { useNotificationStore } from '@/stores/notification.store'

const notifyStore = useNotificationStore()

// Main state
const loading = ref(false)
const products = ref<any[]>([])
const searchQuery = ref('')
const selectedStatusFilter = ref('all')

const statusFilters = [
  { label: 'Semua Menu', value: 'all' },
  { label: 'Siap Jual', value: 'ready' },
  { label: 'Habis / Perlu PO', value: 'out_of_stock' },
  { label: 'Belum Ada Resep', value: 'no_recipe' }
]

// Recipe Modal state
const showRecipeModal = ref(false)
const loadingRecipe = ref(false)
const savingRecipe = ref(false)
const selectedProductForRecipe = ref<any>(null)
const recipeRows = ref<any[]>([])
const inventoryItemList = ref<any[]>([])

// Draft PO Modal state
const showDraftPOModal = ref(false)
const submittingPO = ref(false)
const selectedProductForPO = ref<any>(null)
const draftPOForm = ref({
  target_portions: 50,
  department: 'Kitchen & Bar',
  notes: ''
})

// Metrics
const readyCount = computed(() => products.value.filter(p => p.stock > 0 && !p.is_out_of_stock).length)
const outOfStockCount = computed(() => products.value.filter(p => p.stock <= 0 || p.is_out_of_stock).length)
const mappedCount = computed(() => products.value.filter(p => p.has_recipe).length)

// Filtered Products
const filteredProducts = computed(() => {
  let list = products.value
  if (searchQuery.value) {
    const q = searchQuery.value.toLowerCase()
    list = list.filter(p =>
      p.name?.toLowerCase().includes(q) ||
      p.sku?.toLowerCase().includes(q) ||
      p.target_station?.toLowerCase().includes(q) ||
      p.category_name?.toLowerCase().includes(q)
    )
  }
  if (selectedStatusFilter.value === 'ready') {
    list = list.filter(p => p.stock > 0 && !p.is_out_of_stock)
  } else if (selectedStatusFilter.value === 'out_of_stock') {
    list = list.filter(p => p.stock <= 0 || p.is_out_of_stock)
  } else if (selectedStatusFilter.value === 'no_recipe') {
    list = list.filter(p => !p.has_recipe)
  }
  return list
})

// Helpers
const formatRupiah = (val: any) => {
  const num = Number(val || 0)
  return num.toLocaleString('id-ID')
}

const getMarginBadgeClass = (pct: number) => {
  if (pct >= 65) return 'bg-emerald-100 text-emerald-800'
  if (pct >= 40) return 'bg-blue-100 text-blue-800'
  if (pct > 0) return 'bg-amber-100 text-amber-800'
  return 'bg-rose-100 text-rose-800'
}

// Fetch products from master endpoint
const fetchProducts = async () => {
  loading.value = true
  try {
    const res = await axios.get('/api/v1/master/products')
    const list = Array.isArray(res.data) ? res.data : (res.data?.data || [])
    products.value = list
  } catch (err: any) {
    notifyStore.error(err.response?.data?.message || 'Gagal memuat katalog produk', 'Error')
  } finally {
    loading.value = false
  }
}

// Load inventory raw materials for recipe builder
const fetchInventoryItems = async () => {
  try {
    const res = await axios.get('/api/v1/master/inventory-items')
    inventoryItemList.value = Array.isArray(res.data) ? res.data : (res.data?.data || [])
  } catch (err) {
    console.error('Failed to load inventory items:', err)
  }
}

// -----------------------------------------------------------------------------
// RECIPE (BOM) LOGIC
// -----------------------------------------------------------------------------
const openRecipeModal = async (product: any) => {
  selectedProductForRecipe.value = product
  showRecipeModal.value = true
  loadingRecipe.value = true
  recipeRows.value = []

  try {
    if (inventoryItemList.value.length === 0) {
      await fetchInventoryItems()
    }
    const res = await axios.get(`/api/v1/master/products/${product.id}/recipe`)
    const recipeData = res.data
    if (recipeData?.items && recipeData.items.length > 0) {
      recipeRows.value = recipeData.items.map((it: any) => ({
        inventory_item_id: it.inventory_item_id,
        item_name: it.item_name,
        quantity_required: it.quantity_required,
        uom: it.uom,
        unit_cost: it.unit_cost,
        current_stock: it.current_stock,
        can_make_portions: it.can_make_portions,
        instructions: it.instructions || ''
      }))
    } else {
      // Default empty row
      addRecipeRow()
    }
  } catch (err) {
    console.error('Failed to load recipe:', err)
    addRecipeRow()
  } finally {
    loadingRecipe.value = false
  }
}

const addRecipeRow = () => {
  recipeRows.value.push({
    inventory_item_id: '',
    item_name: '',
    quantity_required: 1,
    uom: 'pcs',
    unit_cost: 0,
    current_stock: 0,
    can_make_portions: 0,
    instructions: ''
  })
}

const removeRecipeRow = (idx: number) => {
  recipeRows.value.splice(idx, 1)
}

const onIngredientSelect = (row: any) => {
  const found = inventoryItemList.value.find(i => i.id === row.inventory_item_id)
  if (found) {
    row.item_name = found.name
    row.uom = found.uom
    row.unit_cost = found.average_cost || 0
    row.current_stock = found.current_stock || 0
    if (row.quantity_required > 0) {
      row.can_make_portions = Math.floor(row.current_stock / row.quantity_required)
    }
  }
}

// Live recipe totals
const recipeTotalCOGS = computed(() => {
  return recipeRows.value.reduce((acc, r) => {
    const qty = Number(r.quantity_required) || 0
    const cost = Number(r.unit_cost) || 0
    return acc + (qty * cost)
  }, 0)
})

const recipeGrossProfit = computed(() => {
  const price = Number(selectedProductForRecipe.value?.base_price || selectedProductForRecipe.value?.price || 0)
  return price - recipeTotalCOGS.value
})

const recipeMarginPercent = computed(() => {
  const price = Number(selectedProductForRecipe.value?.base_price || selectedProductForRecipe.value?.price || 0)
  if (price <= 0) return 0
  return (recipeGrossProfit.value / price) * 100
})

const recipeCookablePortions = computed(() => {
  if (recipeRows.value.length === 0) return 0
  let minPortion: number | null = null
  for (const r of recipeRows.value) {
    const qty = Number(r.quantity_required) || 0
    const stock = Number(r.current_stock) || 0
    if (qty > 0) {
      const p = Math.floor(stock / qty)
      if (minPortion === null || p < minPortion) {
        minPortion = p
      }
    }
  }
  return minPortion ?? 0
})

const saveRecipe = async () => {
  if (!selectedProductForRecipe.value) return
  const validItems = recipeRows.value
    .filter(r => r.inventory_item_id && r.quantity_required > 0)
    .map(r => ({
      inventory_item_id: r.inventory_item_id,
      quantity_required: Number(r.quantity_required),
      uom: r.uom || 'pcs',
      instructions: r.instructions || ''
    }))

  if (validItems.length === 0) {
    notifyStore.warning('Masukkan minimal 1 bahan baku valid beserta takarannya!', 'Validasi')
    return
  }

  savingRecipe.value = true
  try {
    const res = await axios.put(`/api/v1/master/products/${selectedProductForRecipe.value.id}/recipe`, {
      items: validItems
    })
    notifyStore.success(res.data?.message || 'Resep BOM berhasil disimpan!', 'Sukses')
    showRecipeModal.value = false
    await fetchProducts()
  } catch (err: any) {
    notifyStore.error(err.response?.data?.message || 'Gagal menyimpan resep', 'Error')
  } finally {
    savingRecipe.value = false
  }
}

// -----------------------------------------------------------------------------
// DRAFT PO REQUISITION LOGIC
// -----------------------------------------------------------------------------
const openDraftPOModal = (product: any) => {
  selectedProductForPO.value = product
  draftPOForm.value = {
    target_portions: 50,
    department: 'Kitchen & Bar',
    notes: `Permintaan bahan baku menu ${product.name} habis/menipis.`
  }
  showDraftPOModal.value = true
}

const submitDraftPO = async () => {
  if (!selectedProductForPO.value) return
  if (draftPOForm.value.target_portions <= 0) {
    notifyStore.warning('Masukkan jumlah target porsi yang valid!', 'Validasi')
    return
  }

  submittingPO.value = true
  try {
    const res = await axios.post('/api/v1/inventory/pr/from-menu', {
      product_id: selectedProductForPO.value.id,
      target_portions: Number(draftPOForm.value.target_portions),
      department: draftPOForm.value.department,
      notes: draftPOForm.value.notes
    })
    notifyStore.success(
      res.data?.message || `Draft PR ${res.data?.pr_number || ''} berhasil dibuat & dikirim ke Tim Gudang!`,
      'Draft PO Terkirim'
    )
    showDraftPOModal.value = false
  } catch (err: any) {
    notifyStore.error(err.response?.data?.message || 'Gagal mengajukan draft PO', 'Error')
  } finally {
    submittingPO.value = false
  }
}

onMounted(() => {
  fetchProducts()
  fetchInventoryItems()
})
</script>
