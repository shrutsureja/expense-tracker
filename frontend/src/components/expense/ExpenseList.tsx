import { formatDate } from '../../utils/formatters';
import { ExpenseCard } from './ExpenseCard';
import type { Expense } from '../../types';

interface ExpenseListProps {
  expenses: Expense[];
  onDelete?: (id: number) => void;
  onEdit?: (expense: Expense) => void;
  showPerson?: boolean;
}

export function ExpenseList({ expenses, onDelete, onEdit, showPerson = true }: ExpenseListProps) {
  if (expenses.length === 0) return null;

  // Group by date
  const grouped = expenses.reduce((acc, expense) => {
    const date = expense.expense_date;
    if (!acc[date]) acc[date] = [];
    acc[date].push(expense);
    return acc;
  }, {} as Record<string, Expense[]>);

  const sortedDates = Object.keys(grouped).sort((a, b) => b.localeCompare(a));

  return (
    <div className="flex flex-col">
      {sortedDates.map(date => (
        <div key={date}>
          <div className="px-4 py-2 bg-slate-100 sticky top-14 z-10">
            <span className="text-xs font-bold text-gray-500 uppercase tracking-wider">
              {formatDate(date)}
            </span>
          </div>
          <div className="divide-y divide-gray-100">
            {grouped[date].map(expense => (
              <ExpenseCard
                key={expense.id}
                expense={expense}
                onDelete={onDelete}
                onEdit={onEdit}
                showPerson={showPerson}
              />
            ))}
          </div>
        </div>
      ))}
    </div>
  );
}
