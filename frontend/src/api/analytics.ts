import { api } from './client';
import type {
  AnalyticsSummary,
  CategoryBreakdown,
  PersonBreakdown,
  PaymentMethodBreakdown,
  DailyTotal,
  MonthlyComparison,
  DateRange,
} from '../types';

function buildDateParams(range: DateRange, dateFrom?: string, dateTo?: string): string {
  const params = new URLSearchParams({ range });
  if (range === 'custom') {
    if (dateFrom) params.set('from', dateFrom);
    if (dateTo) params.set('to', dateTo);
  }
  return params.toString();
}

export const analyticsApi = {
  summary: (range: DateRange, dateFrom?: string, dateTo?: string) =>
    api.get<AnalyticsSummary>(`/analytics/summary?${buildDateParams(range, dateFrom, dateTo)}`),

  byPerson: (range: DateRange, dateFrom?: string, dateTo?: string) =>
    api.get<PersonBreakdown[]>(`/analytics/by-person?${buildDateParams(range, dateFrom, dateTo)}`),

  byCategory: (range: DateRange, dateFrom?: string, dateTo?: string) =>
    api.get<CategoryBreakdown[]>(`/analytics/by-category?${buildDateParams(range, dateFrom, dateTo)}`),

  byPaymentMethod: (range: DateRange, dateFrom?: string, dateTo?: string) =>
    api.get<PaymentMethodBreakdown[]>(`/analytics/by-payment-method?${buildDateParams(range, dateFrom, dateTo)}`),

  dailyTrend: (month: string) =>
    api.get<DailyTotal[]>(`/analytics/daily-trend?month=${month}`),

  monthlyComparison: () =>
    api.get<MonthlyComparison>('/analytics/monthly-comparison'),
};
