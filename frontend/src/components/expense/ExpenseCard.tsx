import { Trash2, Edit2 } from 'lucide-react';
import { useState } from 'react';
import { formatCurrency } from '../../utils/formatters';
import { PAYMENT_METHODS } from '../../utils/constants';
import type { Expense } from '../../types';

interface ExpenseCardProps {
  expense: Expense;
  onDelete?: (id: number) => void;
  onEdit?: (expense: Expense) => void;
  showPerson?: boolean;
}

export function ExpenseCard({ expense, onDelete, onEdit, showPerson = true }: ExpenseCardProps) {
  const [confirming, setConfirming] = useState(false);
  const pm = PAYMENT_METHODS.find(p => p.value === expense.payment_method);

  return (
    <div className="flex items-center gap-3 bg-white px-4 py-3 hover:bg-gray-50 transition-colors">
      {/* Tag icon */}
      <div className="w-11 h-11 rounded-2xl bg-blue-50 flex items-center justify-center text-xl flex-shrink-0">
        {expense.tag_icon || '💰'}
      </div>

      {/* Details */}
      <div className="flex-1 min-w-0">
        <div className="flex items-center gap-2">
          <span className="font-semibold text-gray-900 text-base">{expense.tag_name}</span>
          <span className={`text-xs px-2 py-0.5 rounded-full font-medium border ${pm?.color || 'bg-gray-100 text-gray-500 border-gray-200'}`}>
            {pm?.label || expense.payment_method}
          </span>
        </div>
        <div className="flex items-center gap-2 mt-0.5">
          {showPerson && (
            <span className="text-sm text-gray-500">{expense.user_display_name}</span>
          )}
          {expense.note && (
            <span className="text-sm text-gray-400 truncate">· {expense.note}</span>
          )}
        </div>
      </div>

      {/* Amount + actions */}
      <div className="flex flex-col items-end gap-1 flex-shrink-0">
        <span className="text-lg font-bold text-gray-900">{formatCurrency(expense.amount)}</span>
        <div className="flex gap-1">
          {onEdit && (
            <button
              onClick={() => onEdit(expense)}
              className="p-1.5 rounded-lg hover:bg-gray-100 text-gray-400 hover:text-blue-500 transition-colors"
            >
              <Edit2 size={14} />
            </button>
          )}
          {onDelete && !confirming && (
            <button
              onClick={() => setConfirming(true)}
              className="p-1.5 rounded-lg hover:bg-gray-100 text-gray-400 hover:text-red-500 transition-colors"
            >
              <Trash2 size={14} />
            </button>
          )}
          {onDelete && confirming && (
            <div className="flex gap-1">
              <button
                onClick={() => { onDelete(expense.id); setConfirming(false); }}
                className="px-2 py-1 text-xs bg-red-500 text-white rounded-lg font-medium"
              >
                Delete
              </button>
              <button
                onClick={() => setConfirming(false)}
                className="px-2 py-1 text-xs bg-gray-200 text-gray-600 rounded-lg font-medium"
              >
                No
              </button>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
