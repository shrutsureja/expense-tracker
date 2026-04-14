import { useState, useEffect } from 'react';
import { Search } from 'lucide-react';
import { tagsApi } from '../../api/tags';
import { Button } from '../ui/Button';
import { PAYMENT_METHODS } from '../../utils/constants';
import { todayStr, yesterdayStr } from '../../utils/formatters';
import type { Tag, PaymentMethod, AddExpenseRequest, Expense } from '../../types';

interface ExpenseFormProps {
  initial?: Expense;
  onSubmit: (data: AddExpenseRequest) => Promise<void>;
  loading?: boolean;
}

export function ExpenseForm({ initial, onSubmit, loading }: ExpenseFormProps) {
  const [tags, setTags] = useState<Tag[]>([]);
  const [tagSearch, setTagSearch] = useState('');
  const [selectedTag, setSelectedTag] = useState<Tag | null>(null);
  const [amount, setAmount] = useState(initial ? String(initial.amount) : '');
  const [paymentMethod, setPaymentMethod] = useState<PaymentMethod>(initial?.payment_method || 'upi');
  const [note, setNote] = useState(initial?.note || '');
  const [date, setDate] = useState(initial?.expense_date || todayStr());

  useEffect(() => {
    tagsApi.list().then(setTags);
  }, []);

  useEffect(() => {
    if (initial) {
      const tag = tags.find(t => t.id === initial.tag_id);
      if (tag) setSelectedTag(tag);
    }
  }, [initial, tags]);

  const filteredTags = tagSearch
    ? tags.filter(t => t.name.toLowerCase().includes(tagSearch.toLowerCase()))
    : tags;

  const handleSubmit = async () => {
    if (!selectedTag || !amount) return;
    await onSubmit({
      amount: parseFloat(amount),
      tag_id: selectedTag.id,
      payment_method: paymentMethod,
      note: note.trim() || undefined,
      expense_date: date,
    });
  };

  return (
    <div className="flex flex-col gap-5 p-4">
      {/* Amount */}
      <div>
        <label className="text-sm font-semibold text-gray-600 block mb-1">Amount (₹) *</label>
        <div className="relative">
          <span className="absolute left-4 top-1/2 -translate-y-1/2 text-2xl font-bold text-gray-400">₹</span>
          <input
            type="number"
            inputMode="decimal"
            placeholder="0"
            value={amount}
            onChange={e => setAmount(e.target.value)}
            className="w-full pl-10 pr-4 py-4 text-3xl font-bold rounded-2xl border-2 border-gray-200
                       focus:border-blue-500 focus:outline-none text-gray-900 bg-gray-50"
            autoFocus
          />
        </div>
      </div>

      {/* Category */}
      <div>
        <label className="text-sm font-semibold text-gray-600 block mb-2">Category *</label>
        {selectedTag ? (
          <div className="flex items-center gap-2 p-3 bg-blue-50 border-2 border-blue-200 rounded-2xl">
            <span className="text-2xl">{selectedTag.icon}</span>
            <span className="font-semibold text-blue-800">{selectedTag.name}</span>
            <button
              type="button"
              onClick={() => setSelectedTag(null)}
              className="ml-auto text-blue-400 hover:text-blue-600 text-sm"
            >
              Change
            </button>
          </div>
        ) : (
          <div className="flex flex-col gap-2">
            <div className="relative">
              <Search size={16} className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" />
              <input
                type="text"
                placeholder="Search categories…"
                value={tagSearch}
                onChange={e => setTagSearch(e.target.value)}
                className="w-full pl-9 pr-4 py-2.5 rounded-xl border-2 border-gray-200 focus:border-blue-400 focus:outline-none text-sm"
              />
            </div>
            <div className="max-h-48 overflow-y-auto">
              <div className="flex flex-wrap gap-2 p-1">
                {filteredTags.map(tag => (
                  <button
                    key={tag.id}
                    type="button"
                    onClick={() => setSelectedTag(tag)}
                    className="flex items-center gap-1.5 px-3 py-2 rounded-xl bg-white border-2 border-gray-200
                               hover:border-blue-400 hover:bg-blue-50 active:scale-95 transition-all text-sm font-medium text-gray-700"
                  >
                    <span>{tag.icon}</span>
                    <span>{tag.name}</span>
                  </button>
                ))}
              </div>
            </div>
          </div>
        )}
      </div>

      {/* Payment method */}
      <div>
        <label className="text-sm font-semibold text-gray-600 block mb-2">Payment Method *</label>
        <div className="grid grid-cols-4 gap-2">
          {PAYMENT_METHODS.map(pm => (
            <button
              key={pm.value}
              type="button"
              onClick={() => setPaymentMethod(pm.value)}
              className={`
                flex flex-col items-center gap-1 py-3 rounded-xl border-2 transition-all text-xs font-semibold
                ${paymentMethod === pm.value
                  ? 'border-blue-500 bg-blue-50 text-blue-700'
                  : 'border-gray-200 bg-white text-gray-600 hover:border-gray-300'
                }
              `}
            >
              <span className="text-xl">{pm.icon}</span>
              {pm.label}
            </button>
          ))}
        </div>
      </div>

      {/* Date */}
      <div>
        <label className="text-sm font-semibold text-gray-600 block mb-2">Date *</label>
        <div className="flex gap-2 flex-wrap">
          <button
            type="button"
            onClick={() => setDate(todayStr())}
            className={`px-4 py-2 rounded-xl border-2 text-sm font-semibold transition-all
              ${date === todayStr() ? 'border-blue-500 bg-blue-50 text-blue-700' : 'border-gray-200 bg-white text-gray-600'}`}
          >
            Today
          </button>
          <button
            type="button"
            onClick={() => setDate(yesterdayStr())}
            className={`px-4 py-2 rounded-xl border-2 text-sm font-semibold transition-all
              ${date === yesterdayStr() ? 'border-blue-500 bg-blue-50 text-blue-700' : 'border-gray-200 bg-white text-gray-600'}`}
          >
            Yesterday
          </button>
          <input
            type="date"
            value={date}
            onChange={e => setDate(e.target.value)}
            className="flex-1 min-w-[140px] px-3 py-2 rounded-xl border-2 border-gray-200 focus:border-blue-400 focus:outline-none text-sm"
          />
        </div>
      </div>

      {/* Note */}
      <div>
        <label className="text-sm font-semibold text-gray-600 block mb-1">Note (optional)</label>
        <input
          type="text"
          placeholder="e.g. Monthly groceries from Big Bazaar"
          value={note}
          onChange={e => setNote(e.target.value)}
          className="w-full px-4 py-3 rounded-xl border-2 border-gray-200 focus:border-blue-400 focus:outline-none text-base bg-white"
        />
      </div>

      <Button
        size="lg"
        fullWidth
        onClick={handleSubmit}
        disabled={!amount || !selectedTag || isNaN(parseFloat(amount)) || parseFloat(amount) <= 0 || loading}
      >
        {loading ? 'Saving…' : initial ? 'Update Expense' : '✓ Add Expense'}
      </Button>
    </div>
  );
}
