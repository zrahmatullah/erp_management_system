export interface ApiResponse<T> {
  data: T;
  message?: string;
  success: boolean;
}

export interface PaginationMeta {
  page: number;
  pageSize: number;
  total: number;
  totalPages: number;
}

export interface PaginatedResponse<T> {
  data: T[];
  meta: PaginationMeta;
}

export interface SelectOption {
  label: string;
  value: string | number;
}

export type StatusType = 'success' | 'warning' | 'error' | 'info' | 'default';

export interface FilterParams {
  [key: string]: any;
}
