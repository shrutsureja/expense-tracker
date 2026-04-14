import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { expensesApi } from '../api/expenses';
import { AppShell } from '../components/layout/AppShell';
import { ExpenseForm } from '../components/expense/ExpenseForm';
import { useToast } from '../components/ui/Toast';
import type { AddExpenseRequest } from '../types';

export function AddExpensePage() {
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();
  const toast = useToast();

  const handleSubmit = async (data: AddExpenseRequest) => {
    setLoading(true);
    try {
      await expensesApi.add(data);
      toast('Expense added!', 'success');
      navigate('/');
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : 'Failed to add expense';
      toast(msg, 'error');
    } finally {
      setLoading(false);
    }
  };

  return (
    <AppShell title="Add Expense">
      <ExpenseForm onSubmit={handleSubmit} loading={loading} />
    </AppShell>
  );
}
