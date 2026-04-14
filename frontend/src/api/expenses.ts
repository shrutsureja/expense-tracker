import { api } from './client';
import type { Expense, ExpenseListResponse, AddExpenseRequest } from '../types';

export interface ExpenseFilters {
  page?: number;
  page_size?: number;
  user_id?: number;
  tag_id?: number;
  payment_method?: string;
  date_from?: string;
  date_to?: string;
}

export const expensesApi = {
  list: (filters: ExpenseFilters = {}) => {
    const params = new URLSearchParams();
    if (filters.page) params.set('page', String(filters.page));
    if (filters.page_size) params.set('page_size', String(filters.page_size));
    if (filters.user_id) params.set('user_id', String(filters.user_id));
    if (filters.tag_id) params.set('tag_id', String(filters.tag_id));
    if (filters.payment_method) params.set('payment_method', filters.payment_method);
    if (filters.date_from) params.set('date_from', filters.date_from);
    if (filters.date_to) params.set('date_to', filters.date_to);
    const qs = params.toString();
    return api.get<ExpenseListResponse>(`/expenses/${qs ? `?${qs}` : ''}`);
  },

  recent: () => api.get<Expense[]>('/expenses/recent'),

  add: (data: AddExpenseRequest) => api.post<Expense>('/expenses/', data),

  update: (id: number, data: Partial<AddExpenseRequest>) =>
    api.put<Expense>(`/expenses/${id}`, data),

  delete: (id: number) => api.delete<void>(`/expenses/${id}`),
};
