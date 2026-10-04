import api from '@/plugins/axios';

export interface OpenShiftPayload {
  branch_id?: string;
  shift_name: string;
  opening_cash_float: number;
  notes?: string;
}

export interface CashMovementPayload {
  movement_type: 'cash_drop' | 'paid_out' | 'cash_in';
  amount: number;
  reason: string;
  authorized_by?: string;
}

export interface CloseShiftPayload {
  actual_cash_counted: number;
  supervisor_id?: string;
  notes?: string;
}

export const getCurrentShift = async (branchId?: string) =>
  api.get('/pos/shifts/current', { params: { branch_id: branchId } });

export const openShift = async (data: OpenShiftPayload) =>
  api.post('/pos/shifts/open', data);

export const recordCashMovement = async (shiftId: string, data: CashMovementPayload) =>
  api.post(`/pos/shifts/${shiftId}/cash-movement`, data);

export const closeShift = async (shiftId: string, data: CloseShiftPayload) =>
  api.post(`/pos/shifts/${shiftId}/close`, data);

export const getShiftSummary = async (shiftId: string) =>
  api.get(`/pos/shifts/${shiftId}/summary`);

export const getShifts = async (params?: any) =>
  api.get('/pos/shifts', { params });
