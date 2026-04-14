export interface User {
  id: number;
  username: string;
  display_name: string;
  role: 'super_admin' | 'family_owner' | 'family_member';
  family_id: number | null;
  is_active: boolean;
}

export interface Family {
  id: number;
  name: string;
  created_by: number;
  owner_id: number | null;
  created_at: string;
}

export interface FamilyWithOwner {
  family: Family;
  owner: User;
}

export interface Tag {
  id: number;
  name: string;
  icon: string;
  family_id: number | null;
  sort_order: number;
}

export type PaymentMethod = 'cash' | 'online' | 'upi' | 'card';

export interface Expense {
  id: number;
  family_id: number;
  user_id: number;
  amount: number;
  tag_id: number;
  payment_method: PaymentMethod;
  note: string;
  expense_date: string;
  user_display_name: string;
  tag_name: string;
  tag_icon: string;
  created_at: string;
}

export interface ExpenseListResponse {
  expenses: Expense[];
  total: number;
  page: number;
  page_size: number;
}

export type DateRange = 'today' | 'week' | 'month' | 'custom';

export interface AnalyticsSummary {
  total_amount: number;
  count: number;
  avg_per_day: number;
}

export interface CategoryBreakdown {
  tag_id: number;
  tag_name: string;
  tag_icon: string;
  total: number;
  count: number;
  percentage: number;
}

export interface PersonBreakdown {
  user_id: number;
  display_name: string;
  total: number;
  count: number;
  percentage: number;
}

export interface PaymentMethodBreakdown {
  payment_method: string;
  total: number;
  count: number;
  percentage: number;
}

export interface DailyTotal {
  date: string;
  total: number;
}

export interface MonthlyComparison {
  this_month: number;
  last_month: number;
  change_percent: number;
}

export interface AddExpenseRequest {
  amount: number;
  tag_id: number;
  payment_method: PaymentMethod;
  note?: string;
  expense_date: string;
}
