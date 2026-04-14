import { useState, useEffect, useCallback } from 'react';
import { expensesApi } from '../api/expenses';
import { analyticsApi } from '../api/analytics';
import { AppShell } from '../components/layout/AppShell';
import { QuickAdd } from '../components/expense/QuickAdd';
import { ExpenseList } from '../components/expense/ExpenseList';
import { ExpenseForm } from '../components/expense/ExpenseForm';
import { Modal } from '../components/ui/Modal';
import { EmptyState } from '../components/ui/EmptyState';
import { Spinner } from '../components/ui/Spinner';
import { useToast } from '../components/ui/Toast';
import { formatCurrency } from '../utils/formatters';
import type { Expense } from '../types';

export function HomePage() {
  const [expenses, setExpenses] = useState<Expense[]>([]);
  const [todayTotal, setTodayTotal] = useState(0);
  const [loading, setLoading] = useState(true);
  const [editExpense, setEditExpense] = useState<Expense | null>(null);
  const [editLoading, setEditLoading] = useState(false);
  const toast = useToast();

  const load = useCallback(async () => {
    try {
      const [recent, summary] = await Promise.all([
        expensesApi.recent(),
        analyticsApi.summary('today'),
      ]);
      setExpenses(recent);
      setTodayTotal(summary.total_amount);
    } catch {
      // silent
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => { load(); }, [load]);

  const handleDelete = async (id: number) => {
    try {
      await expensesApi.delete(id);
      toast('Expense deleted', 'success');
      load();
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : 'Failed to delete';
      toast(msg, 'error');
    }
  };

  const handleEditSubmit = async (data: Parameters<typeof expensesApi.update>[1]) => {
    if (!editExpense) return;
    setEditLoading(true);
    try {
      await expensesApi.update(editExpense.id, data);
      toast('Expense updated', 'success');
      setEditExpense(null);
      load();
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : 'Failed to update';
      toast(msg, 'error');
    } finally {
      setEditLoading(false);
    }
  };

  return (
    <AppShell>
      {/* Today's total */}
      <div className="bg-gradient-to-r from-blue-600 to-blue-700 px-4 py-6 text-white">
        <p className="text-blue-200 text-sm font-medium mb-1">Spent Today</p>
        <p className="text-4xl font-bold">{formatCurrency(todayTotal)}</p>
      </div>

      {/* Quick add */}
      <div className="bg-white border-b border-gray-100 py-4">
        <QuickAdd onAdded={load} />
      </div>

      {/* Recent expenses */}
      <div className="mt-2">
        <div className="px-4 py-3">
          <h2 className="text-base font-bold text-gray-700">Recent Expenses</h2>
        </div>
        {loading ? (
          <div className="flex justify-center py-12"><Spinner className="text-blue-600" /></div>
        ) : expenses.length === 0 ? (
          <EmptyState
            icon="💸"
            title="No expenses yet"
            description="Tap a category above to quickly add your first expense!"
          />
        ) : (
          <div className="bg-white rounded-2xl mx-2 shadow-sm overflow-hidden">
            <ExpenseList
              expenses={expenses}
              onDelete={handleDelete}
              onEdit={setEditExpense}
            />
          </div>
        )}
      </div>

      {/* Edit modal */}
      <Modal
        open={!!editExpense}
        onClose={() => setEditExpense(null)}
        title="Edit Expense"
      >
        {editExpense && (
          <ExpenseForm
            initial={editExpense}
            onSubmit={handleEditSubmit}
            loading={editLoading}
          />
        )}
      </Modal>
    </AppShell>
  );
}
