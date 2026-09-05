import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router'
import AppLayout from '@/components/layout/AppLayout.vue'
import { requireAuth } from './guards'

const routes: Array<RouteRecordRaw> = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/auth/LoginView.vue'),
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
        component: () => import('@/views/dashboard/DashboardView.vue')
      },
      // POS & Restaurant Service
      {
        path: 'pos',
        name: 'POS',
        component: () => import('@/views/pos/POSView.vue')
      },
      {
        path: 'pos/tables',
        name: 'TableManagement',
        component: () => import('@/views/pos/TableManagementView.vue')
      },
      {
        path: 'pos/kds',
        name: 'KitchenDisplay',
        component: () => import('@/views/pos/KitchenDisplayView.vue')
      },
      // Inventory & Supply Chain
      {
        path: 'inventory',
        name: 'InventoryStock',
        component: () => import('@/views/inventory/StockList.vue')
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
        component: () => import('@/views/hris/EmployeeList.vue')
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
        component: () => import('@/views/finance/FinanceOverview.vue')
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
        component: () => import('@/views/reports/SalesReport.vue')
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
      }
    ]
  },
  {
    path: '/:pathMatch(.*)*',
    redirect: '/dashboard'
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
