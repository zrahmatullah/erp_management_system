import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router'
import AppLayout from '@/components/layout/AppLayout.vue'
import { createVersionedRoute } from '@/utils/versionedRoute'

const routes: Array<RouteRecordRaw> = [
  {
    path: '/login',
    name: 'Login',
    component: createVersionedRoute(
      () => import('@/views/auth/LoginView.vue'),
      () => import('@/views/v2/auth/LoginViewV2.vue')
    ),
    meta: { requiresAuth: false }
  },
  {
    path: '/',
    component: AppLayout,
    meta: { requiresAuth: true },
    children: [
      {
        path: '',
        redirect: '/dashboard'
      },
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: createVersionedRoute(
          () => import('@/views/dashboard/DashboardView.vue'),
          () => import('@/views/v2/dashboard/DashboardViewV2.vue')
        )
      },
      // POS & Restaurant Service
      {
        path: 'pos',
        name: 'POS',
        component: createVersionedRoute(
          () => import('@/views/pos/POSView.vue'),
          () => import('@/views/v2/pos/POSViewV2.vue')
        )
      },
      {
        path: 'pos/tables',
        name: 'TableManagement',
        component: () => import('@/views/pos/TableManagementView.vue')
      },
      {
        path: 'pos/kds',
        name: 'KitchenDisplay',
        component: createVersionedRoute(
          () => import('@/views/pos/KitchenDisplayView.vue'),
          () => import('@/views/v2/pos/KitchenDisplayViewV2.vue')
        )
      },
      {
        path: 'pos/transactions',
        name: 'POSTransactions',
        component: () => import('@/views/pos/TransactionHistoryView.vue')
      },
      {
        path: 'pos/shifts',
        name: 'CashierShifts',
        component: () => import('@/views/pos/CashierShiftView.vue')
      },
      // Inventory & Supply Chain
      {
        path: 'inventory',
        name: 'InventoryStock',
        component: createVersionedRoute(
          () => import('@/views/inventory/StockList.vue'),
          () => import('@/views/v2/inventory/StockListV2.vue')
        )
      },
      {
        path: 'inventory/stock-card',
        name: 'StockCard',
        component: () => import('@/views/inventory/StockCardView.vue')
      },
      {
        path: 'inventory/opname',
        name: 'StockOpname',
        component: () => import('@/views/inventory/StockOpnameView.vue')
      },
      {
        path: 'inventory/po',
        name: 'PurchaseOrders',
        component: () => import('@/views/inventory/PurchaseOrderList.vue')
      },
      // HRIS & Payroll
      {
        path: 'hris',
        name: 'EmployeeList',
        component: createVersionedRoute(
          () => import('@/views/hris/EmployeeList.vue'),
          () => import('@/views/v2/hris/EmployeeListV2.vue')
        )
      },
      {
        path: 'hris/attendance',
        name: 'Attendance',
        component: () => import('@/views/hris/AttendanceView.vue')
      },
      {
        path: 'hris/shifts',
        name: 'ShiftSchedule',
        component: () => import('@/views/hris/ShiftScheduleView.vue')
      },
      {
        path: 'hris/leaves',
        name: 'LeaveManagement',
        component: () => import('@/views/hris/LeaveManagementView.vue')
      },
      {
        path: 'hris/payroll',
        name: 'PayrollRun',
        component: () => import('@/views/hris/PayrollRun.vue')
      },
      // Finance & Accounting
      {
        path: 'finance',
        name: 'FinanceOverview',
        component: createVersionedRoute(
          () => import('@/views/finance/FinanceOverview.vue'),
          () => import('@/views/v2/finance/FinanceOverviewV2.vue')
        )
      },
      {
        path: 'finance/journal',
        name: 'JournalEntry',
        component: () => import('@/views/finance/JournalEntryForm.vue')
      },
      // Reports & Analytics
      {
        path: 'reports',
        name: 'Reports',
        component: createVersionedRoute(
          () => import('@/views/reports/ReportsView.vue'),
          () => import('@/views/v2/reports/ReportsViewV2.vue')
        )
      },
      // Super Admin Settings & Master Hub
      {
        path: 'settings/master',
        name: 'MasterDataHub',
        component: () => import('@/views/admin/MasterDataHub.vue')
      },
      {
        path: 'settings/roles',
        name: 'RoleManagement',
        component: () => import('@/views/settings/RoleManagement.vue')
      },
      {
        path: 'settings/company',
        name: 'CompanyProfile',
        component: () => import('@/views/settings/CompanyProfile.vue')
      },
      // Menu Management
      {
        path: 'menu/products',
        name: 'ProductList',
        component: () => import('@/views/menu/ProductList.vue')
      },
      {
        path: 'menu/products/create',
        name: 'ProductCreate',
        component: () => import('@/views/menu/ProductForm.vue')
      }
    ]
  },
  {
    path: '/404',
    name: 'NotFound',
    component: () => import('@/views/error/NotFoundView.vue'),
    meta: { requiresAuth: false }
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'CatchAll',
    component: () => import('@/views/error/NotFoundView.vue'),
    meta: { requiresAuth: false }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

import { useLoadingStore } from '@/stores/loading.store'

// Navigation Guard
router.beforeEach((to, _from, next) => {
  const loadingStore = useLoadingStore()
  loadingStore.start()

  const token = localStorage.getItem('token')
  if (to.meta.requiresAuth && !token) {
    next('/login')
  } else if (to.path === '/login' && token) {
    next('/dashboard')
  } else {
    next()
  }
})

router.afterEach(() => {
  const loadingStore = useLoadingStore()
  loadingStore.finish()
})

export default router
