import { defineStore } from 'pinia';
import { ref, computed } from 'vue';

export interface Toast {
  id: string;
  type: 'success' | 'error' | 'warning' | 'info';
  title?: string;
  message: string;
  duration?: number;
}

interface NotificationItem {
  id: string;
  title: string;
  message: string;
  isRead: boolean;
  createdAt: string;
}

export const useNotificationStore = defineStore('notification', () => {
  const notifications = ref<NotificationItem[]>([]);
  const toasts = ref<Toast[]>([]);
  const unreadCount = computed(() => notifications.value.filter(n => !n.isRead).length);

  const addToast = (toast: Omit<Toast, 'id'>) => {
    const id = 'toast-' + Math.random().toString(36).substr(2, 9);
    const newToast: Toast = {
      id,
      duration: toast.duration ?? 4000,
      ...toast
    };
    toasts.value.push(newToast);

    if (newToast.duration && newToast.duration > 0) {
      setTimeout(() => {
        removeToast(id);
      }, newToast.duration);
    }
    return id;
  };

  const success = (message: string, title: string = 'Berhasil') => {
    return addToast({ type: 'success', title, message });
  };

  const error = (message: string, title: string = 'Terjadi Kesalahan') => {
    return addToast({ type: 'error', title, message, duration: 5000 });
  };

  const warning = (message: string, title: string = 'Peringatan') => {
    return addToast({ type: 'warning', title, message });
  };

  const info = (message: string, title: string = 'Informasi') => {
    return addToast({ type: 'info', title, message });
  };

  const removeToast = (id: string) => {
    toasts.value = toasts.value.filter(t => t.id !== id);
  };

  const fetchNotifications = async () => {
    // API call placeholder
  };

  const markAsRead = async (id: string) => {
    const notif = notifications.value.find(n => n.id === id);
    if (notif) notif.isRead = true;
  };

  const markAllRead = async () => {
    notifications.value.forEach(n => n.isRead = true);
  };

  const deleteNotification = async (id: string) => {
    notifications.value = notifications.value.filter(n => n.id !== id);
  };

  return {
    toasts,
    addToast,
    success,
    error,
    warning,
    info,
    removeToast,
    notifications,
    unreadCount,
    fetchNotifications,
    markAsRead,
    markAllRead,
    deleteNotification
  };
});
