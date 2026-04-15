import { useState, useEffect, useRef } from 'react';
import { Plus } from 'lucide-react';
import { tagsApi } from '../../api/tags';
import { expensesApi } from '../../api/expenses';
import { useAuth } from '../../context/AuthContext';
import { useToast } from '../ui/Toast';
import { Modal } from '../ui/Modal';
import { Button } from '../ui/Button';
import { PAYMENT_METHODS } from '../../utils/constants';
import { todayStr } from '../../utils/formatters';
import type { Tag, PaymentMethod, AddExpenseRequest } from '../../types';

interface QuickAddProps {
  onAdded: () => void;
}

const ICON_SUGGESTIONS = ['💰', '🏷️', '🛍️', '🎯', '📌', '⚡', '🌟', '🔑', '🎪', '🏪'];

export function QuickAdd({ onAdded }: QuickAddProps) {
  const { user } = useAuth();
  const isOwner = user?.role === 'family_owner';

  const [tags, setTags] = useState<Tag[]>([]);
  const [selectedTag, setSelectedTag] = useState<Tag | null>(null);
  const [amount, setAmount] = useState('');
  const [paymentMethod, setPaymentMethod] = useState<PaymentMethod>('upi');
  const [note, setNote] = useState('');
  const [loading, setLoading] = useState(false);

  // Create tag state
  const [showCreateTag, setShowCreateTag] = useState(false);
  const [newTagName, setNewTagName] = useState('');
  const [newTagIcon, setNewTagIcon] = useState('💰');
  const [tagCreateLoading, setTagCreateLoading] = useState(false);

  const toast = useToast();
  const amountRef = useRef<HTMLInputElement>(null);

  const loadTags = () =>
    tagsApi.suggested()
      .then(setTags)
      .catch(() => tagsApi.list().then(t => setTags(t.slice(0, 12))));

  useEffect(() => { loadTags(); }, []);

  useEffect(() => {
    if (selectedTag && amountRef.current) {
      setTimeout(() => amountRef.current?.focus(), 100);
    }
  }, [selectedTag]);

  const handleSubmit = async () => {
    if (!selectedTag || !amount || isNaN(parseFloat(amount))) return;
    setLoading(true);
    try {
      const data: AddExpenseRequest = {
        amount: parseFloat(amount),
        tag_id: selectedTag.id,
        payment_method: paymentMethod,
        note: note.trim() || undefined,
        expense_date: todayStr(),
      };
      await expensesApi.add(data);
      toast(`Added ₹${amount} for ${selectedTag.name}`, 'success');
      setSelectedTag(null);
      setAmount('');
      setNote('');
      setPaymentMethod('upi');
      onAdded();
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : 'Failed to add expense';
      toast(msg, 'error');
    } finally {
      setLoading(false);
    }
  };

  const handleCreateTag = async () => {
    const name = newTagName.trim();
    if (!name) return;
    setTagCreateLoading(true);
    try {
      const created = await tagsApi.create(name, newTagIcon);
      await loadTags();
      setShowCreateTag(false);
      setNewTagName('');
      setNewTagIcon('💰');
      setSelectedTag(created);
      toast(`Tag "${name}" created`, 'success');
    } catch (err: unknown) {
      toast(err instanceof Error ? err.message : 'Failed to create tag', 'error');
    } finally {
      setTagCreateLoading(false);
    }
  };

  return (
    <div>
      {/* Tag chips */}
      <div className="px-4 pb-2">
        <p className="text-xs font-semibold text-gray-400 uppercase tracking-wider mb-3">Quick Add</p>
        <div className="flex flex-wrap gap-2">
          {tags.map(tag => (
            <button
              key={tag.id}
              onClick={() => setSelectedTag(tag)}
              className="flex items-center gap-1.5 px-3 py-2 rounded-2xl bg-white border-2 border-gray-200
                         hover:border-blue-400 hover:bg-blue-50 active:scale-95 transition-all shadow-sm
                         text-sm font-medium text-gray-700"
            >
              <span>{tag.icon}</span>
              <span>{tag.name}</span>
            </button>
          ))}
          {isOwner && (
            <button
              onClick={() => setShowCreateTag(true)}
              className="flex items-center gap-1.5 px-3 py-2 rounded-2xl bg-white border-2 border-dashed border-gray-300
                         hover:border-green-400 hover:bg-green-50 active:scale-95 transition-all shadow-sm
                         text-sm font-medium text-gray-500"
            >
              <Plus size={14} />
              New tag
            </button>
          )}
        </div>
      </div>

      {/* Quick-add expense modal */}
      <Modal
        open={!!selectedTag}
        onClose={() => { setSelectedTag(null); setAmount(''); setNote(''); }}
        title={selectedTag ? `${selectedTag.icon} ${selectedTag.name}` : ''}
      >
        <div className="flex flex-col gap-5">
          {/* Amount input */}
          <div>
            <label className="text-sm font-semibold text-gray-600 block mb-1">Amount (₹)</label>
            <div className="relative">
              <span className="absolute left-4 top-1/2 -translate-y-1/2 text-2xl font-bold text-gray-400">₹</span>
              <input
                ref={amountRef}
                type="number"
                inputMode="decimal"
                placeholder="0"
                value={amount}
                onChange={e => setAmount(e.target.value)}
                className="w-full pl-10 pr-4 py-4 text-3xl font-bold rounded-2xl border-2 border-gray-200
                           focus:border-blue-500 focus:outline-none text-gray-900 bg-gray-50"
              />
            </div>
          </div>

          {/* Payment method */}
          <div>
            <label className="text-sm font-semibold text-gray-600 block mb-2">How did you pay?</label>
            <div className="grid grid-cols-4 gap-2">
              {PAYMENT_METHODS.map(pm => (
                <button
                  key={pm.value}
                  type="button"
                  onClick={() => setPaymentMethod(pm.value)}
                  className={`
                    flex flex-col items-center gap-1 py-2.5 rounded-xl border-2 transition-all text-xs font-semibold
                    ${paymentMethod === pm.value
                      ? 'border-blue-500 bg-blue-50 text-blue-700'
                      : 'border-gray-200 bg-white text-gray-600 hover:border-gray-300'
                    }
                  `}
                >
                  <span className="text-lg">{pm.icon}</span>
                  {pm.label}
                </button>
              ))}
            </div>
          </div>

          {/* Optional note */}
          <div>
            <label className="text-sm font-semibold text-gray-600 block mb-1">Note (optional)</label>
            <input
              type="text"
              placeholder="e.g. Amul milk 2 packets"
              value={note}
              onChange={e => setNote(e.target.value)}
              className="w-full px-4 py-3 rounded-xl border-2 border-gray-200 focus:border-blue-400 focus:outline-none text-base bg-white"
            />
          </div>

          <Button
            size="lg"
            fullWidth
            onClick={handleSubmit}
            disabled={!amount || isNaN(parseFloat(amount)) || parseFloat(amount) <= 0 || loading}
          >
            {loading ? 'Adding…' : `Add ₹${amount || '0'}`}
          </Button>
        </div>
      </Modal>

      {/* Create tag modal (family owners only) */}
      <Modal
        open={showCreateTag}
        onClose={() => { setShowCreateTag(false); setNewTagName(''); setNewTagIcon('💰'); }}
        title="Create New Tag"
      >
        <div className="flex flex-col gap-4">
          <div>
            <label className="text-sm font-semibold text-gray-600 block mb-2">Pick an icon</label>
            <div className="flex flex-wrap gap-2">
              {ICON_SUGGESTIONS.map(icon => (
                <button
                  key={icon}
                  type="button"
                  onClick={() => setNewTagIcon(icon)}
                  className={`w-10 h-10 rounded-xl text-xl flex items-center justify-center transition-all
                    ${newTagIcon === icon ? 'bg-blue-100 border-2 border-blue-400 scale-110' : 'bg-gray-100 hover:bg-gray-200'}`}
                >
                  {icon}
                </button>
              ))}
            </div>
            <p className="text-xs text-gray-400 mt-1">Or type any emoji below in the name field.</p>
          </div>
          <div>
            <label className="text-sm font-semibold text-gray-600 block mb-1">Tag Name</label>
            <input
              type="text"
              placeholder="e.g. Gym, Parking, Party…"
              value={newTagName}
              onChange={e => setNewTagName(e.target.value)}
              className="w-full px-3 py-2.5 rounded-xl border-2 border-gray-200 focus:border-blue-400 focus:outline-none text-sm"
              autoFocus
            />
          </div>
          <Button
            size="lg"
            fullWidth
            onClick={handleCreateTag}
            disabled={!newTagName.trim() || tagCreateLoading}
          >
            {tagCreateLoading ? 'Creating…' : `Create "${newTagName.trim() || '…'}" Tag`}
          </Button>
        </div>
      </Modal>
    </div>
  );
}
