import type { PaymentMethod } from '../types';

export const PAYMENT_METHODS: { value: PaymentMethod; label: string; icon: string; color: string }[] = [
  { value: 'upi', label: 'UPI', icon: '📱', color: 'bg-purple-100 text-purple-700 border-purple-300' },
  { value: 'cash', label: 'Cash', icon: '💵', color: 'bg-green-100 text-green-700 border-green-300' },
  { value: 'online', label: 'Online', icon: '🌐', color: 'bg-blue-100 text-blue-700 border-blue-300' },
  { value: 'card', label: 'Card', icon: '💳', color: 'bg-orange-100 text-orange-700 border-orange-300' },
];

export const DATE_RANGE_LABELS = {
  today: 'Today',
  week: 'This Week',
  month: 'This Month',
  custom: 'Custom',
};
