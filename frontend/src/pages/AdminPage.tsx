import { useState, useEffect } from 'react';
import { Trash2, PlusCircle } from 'lucide-react';
import { familiesApi } from '../api/families';
import type { CreateFamilyRequest } from '../api/families';
import { AppShell } from '../components/layout/AppShell';
import { Modal } from '../components/ui/Modal';
import { Button } from '../components/ui/Button';
import { Input } from '../components/ui/Input';
import { PinInput } from '../components/auth/PinInput';
import { Spinner } from '../components/ui/Spinner';
import { EmptyState } from '../components/ui/EmptyState';
import { useToast } from '../components/ui/Toast';
import type { FamilyWithOwner } from '../types';

export function AdminPage() {
  const [families, setFamilies] = useState<FamilyWithOwner[]>([]);
  const [loading, setLoading] = useState(true);
  const [showCreate, setShowCreate] = useState(false);
  const [form, setForm] = useState<CreateFamilyRequest>({
    name: '', owner_username: '', owner_password: '', owner_display_name: '',
  });
  const [ownerPin, setOwnerPin] = useState('');
  const [formLoading, setFormLoading] = useState(false);
  const toast = useToast();

  const load = () => {
    setLoading(true);
    familiesApi.listAll().then(setFamilies).finally(() => setLoading(false));
  };

  useEffect(() => { load(); }, []);

  const handleCreate = async () => {
    if (!form.name || !form.owner_username || !ownerPin || !form.owner_display_name) return;
    if (ownerPin.length !== 6) {
      toast('Owner PIN must be exactly 6 digits', 'error');
      return;
    }
    setFormLoading(true);
    try {
      await familiesApi.create({ ...form, owner_password: ownerPin });
      toast(`Family "${form.name}" created`, 'success');
      setShowCreate(false);
      setForm({ name: '', owner_username: '', owner_password: '', owner_display_name: '' });
      setOwnerPin('');
      load();
    } catch (err: unknown) {
      toast(err instanceof Error ? err.message : 'Failed to create family', 'error');
    } finally {
      setFormLoading(false);
    }
  };

  const handleDelete = async (id: number, name: string) => {
    if (!confirm(`Delete family "${name}" and all its data? This cannot be undone.`)) return;
    try {
      await familiesApi.deleteFamily(id);
      toast(`Family "${name}" deleted`, 'success');
      load();
    } catch (err: unknown) {
      toast(err instanceof Error ? err.message : 'Failed to delete', 'error');
    }
  };

  return (
    <AppShell title="Admin">
      <div className="p-4">
        <div className="flex justify-between items-center mb-4">
          <h2 className="text-xl font-bold text-gray-800">All Families</h2>
          <Button size="sm" onClick={() => setShowCreate(true)}>
            <PlusCircle size={16} /> New Family
          </Button>
        </div>

        {loading ? (
          <div className="flex justify-center py-12"><Spinner className="text-blue-600" /></div>
        ) : families.length === 0 ? (
          <EmptyState icon="🏠" title="No families yet" description="Create a family to get started." />
        ) : (
          <div className="bg-white rounded-2xl shadow-sm divide-y divide-gray-100">
            {families.map(({ family, owner }) => (
              <div key={family.id} className="flex items-center gap-3 p-4">
                <div className="w-11 h-11 rounded-2xl bg-indigo-100 flex items-center justify-center text-indigo-700 font-bold text-xl flex-shrink-0">
                  🏠
                </div>
                <div className="flex-1 min-w-0">
                  <p className="font-semibold text-gray-900">{family.name}</p>
                  <p className="text-sm text-gray-500">Owner: {owner.display_name} (@{owner.username})</p>
                </div>
                <button
                  onClick={() => handleDelete(family.id, family.name)}
                  className="p-2 rounded-xl hover:bg-gray-100 text-gray-400 hover:text-red-500 transition-colors"
                >
                  <Trash2 size={16} />
                </button>
              </div>
            ))}
          </div>
        )}
      </div>

      <Modal open={showCreate} onClose={() => { setShowCreate(false); setOwnerPin(''); }} title="Create Family">
        <div className="flex flex-col gap-4">
          <Input
            label="Family Name"
            placeholder="e.g. Sureja Family"
            value={form.name}
            onChange={e => setForm(f => ({ ...f, name: e.target.value }))}
          />
          <hr className="border-gray-100" />
          <p className="text-sm font-semibold text-gray-500">Family Owner Account</p>
          <Input
            label="Display Name"
            placeholder="e.g. Shrut"
            value={form.owner_display_name}
            onChange={e => setForm(f => ({ ...f, owner_display_name: e.target.value }))}
          />
          <Input
            label="Username"
            placeholder="e.g. shrut"
            value={form.owner_username}
            onChange={e => setForm(f => ({ ...f, owner_username: e.target.value }))}
            autoCapitalize="none"
          />
          <div>
            <p className="text-sm font-semibold text-gray-600 mb-3">
              Owner PIN <span className="text-gray-400 font-normal">(6 digits)</span>
            </p>
            <PinInput value={ownerPin} onChange={setOwnerPin} maxLength={6} />
            <p className="text-xs text-gray-400 text-center mt-2">
              Share this PIN with the family owner verbally.
            </p>
          </div>
          <Button
            size="lg"
            fullWidth
            onClick={handleCreate}
            disabled={formLoading || !form.name || !form.owner_username || ownerPin.length !== 6 || !form.owner_display_name}
          >
            {formLoading ? 'Creating…' : 'Create Family'}
          </Button>
        </div>
      </Modal>
    </AppShell>
  );
}
