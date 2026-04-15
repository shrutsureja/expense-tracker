import { useState, useEffect } from 'react';
import { UserPlus, Edit2, UserX } from 'lucide-react';
import { familiesApi } from '../api/families';
import type { AddMemberRequest } from '../api/families';
import { AppShell } from '../components/layout/AppShell';
import { Modal } from '../components/ui/Modal';
import { Button } from '../components/ui/Button';
import { Input } from '../components/ui/Input';
import { PinInput } from '../components/auth/PinInput';
import { Spinner } from '../components/ui/Spinner';
import { EmptyState } from '../components/ui/EmptyState';
import { useToast } from '../components/ui/Toast';
import type { User } from '../types';

export function FamilyPage() {
  const [members, setMembers] = useState<User[]>([]);
  const [loading, setLoading] = useState(true);
  const [showAdd, setShowAdd] = useState(false);
  const [editMember, setEditMember] = useState<User | null>(null);
  const [form, setForm] = useState({ username: '', display_name: '' });
  const [addPin, setAddPin] = useState('');
  const [editPin, setEditPin] = useState('');
  const [formLoading, setFormLoading] = useState(false);
  const toast = useToast();

  const load = () => {
    setLoading(true);
    familiesApi.listMembers()
      .then(setMembers)
      .finally(() => setLoading(false));
  };

  useEffect(() => { load(); }, []);

  const resetForm = () => {
    setForm({ username: '', display_name: '' });
    setAddPin('');
    setEditPin('');
  };

  const handleAdd = async () => {
    if (!form.username || !addPin || !form.display_name) return;
    if (addPin.length !== 6) {
      toast('PIN must be exactly 6 digits', 'error');
      return;
    }
    setFormLoading(true);
    try {
      const req: AddMemberRequest = {
        username: form.username,
        pin: addPin,
        display_name: form.display_name,
      };
      await familiesApi.addMember(req);
      toast(`Added ${form.display_name}`, 'success');
      setShowAdd(false);
      resetForm();
      load();
    } catch (err: unknown) {
      toast(err instanceof Error ? err.message : 'Failed to add member', 'error');
    } finally {
      setFormLoading(false);
    }
  };

  const handleUpdate = async () => {
    if (!editMember) return;
    if (editPin && editPin.length !== 6) {
      toast('PIN must be exactly 6 digits', 'error');
      return;
    }
    setFormLoading(true);
    try {
      const update: Record<string, string> = { display_name: form.display_name };
      if (editPin) update.pin = editPin;
      await familiesApi.updateMember(editMember.id, update);
      toast('Member updated', 'success');
      setEditMember(null);
      resetForm();
      load();
    } catch (err: unknown) {
      toast(err instanceof Error ? err.message : 'Failed to update', 'error');
    } finally {
      setFormLoading(false);
    }
  };

  const handleDeactivate = async (member: User) => {
    if (!confirm(`Deactivate ${member.display_name}?`)) return;
    try {
      await familiesApi.deactivateMember(member.id);
      toast(`${member.display_name} deactivated`, 'success');
      load();
    } catch (err: unknown) {
      toast(err instanceof Error ? err.message : 'Failed', 'error');
    }
  };

  return (
    <AppShell title="Family">
      <div className="p-4">
        <div className="flex justify-between items-center mb-4">
          <h2 className="text-xl font-bold text-gray-800">Family Members</h2>
          <Button size="sm" onClick={() => { resetForm(); setShowAdd(true); }}>
            <UserPlus size={16} /> Add Member
          </Button>
        </div>

        {loading ? (
          <div className="flex justify-center py-12"><Spinner className="text-blue-600" /></div>
        ) : members.length === 0 ? (
          <EmptyState icon="👨‍👩‍👧" title="No members yet" description="Add family members so they can track expenses too." />
        ) : (
          <div className="bg-white rounded-2xl shadow-sm divide-y divide-gray-100">
            {members.map(m => (
              <div key={m.id} className="flex items-center gap-3 p-4">
                <div className="w-11 h-11 rounded-full bg-blue-100 flex items-center justify-center text-blue-700 font-bold text-lg flex-shrink-0">
                  {m.display_name[0].toUpperCase()}
                </div>
                <div className="flex-1 min-w-0">
                  <p className="font-semibold text-gray-900">{m.display_name}</p>
                  <p className="text-sm text-gray-500">@{m.username} · {m.role.replace('_', ' ')}</p>
                </div>
                <div className="flex gap-1">
                  <button
                    onClick={() => {
                      setEditMember(m);
                      setForm({ username: m.username, display_name: m.display_name });
                      setEditPin('');
                    }}
                    className="p-2 rounded-xl hover:bg-gray-100 text-gray-400 hover:text-blue-500 transition-colors"
                  >
                    <Edit2 size={16} />
                  </button>
                  {m.role !== 'family_owner' && (
                    <button
                      onClick={() => handleDeactivate(m)}
                      className="p-2 rounded-xl hover:bg-gray-100 text-gray-400 hover:text-red-500 transition-colors"
                    >
                      <UserX size={16} />
                    </button>
                  )}
                </div>
              </div>
            ))}
          </div>
        )}
      </div>

      {/* Add member modal */}
      <Modal open={showAdd} onClose={() => { setShowAdd(false); resetForm(); }} title="Add Family Member">
        <div className="flex flex-col gap-4">
          <Input
            label="Display Name"
            placeholder="e.g. Mom"
            value={form.display_name}
            onChange={e => setForm(f => ({ ...f, display_name: e.target.value }))}
          />
          <Input
            label="Username"
            placeholder="e.g. mom"
            value={form.username}
            onChange={e => setForm(f => ({ ...f, username: e.target.value }))}
            autoCapitalize="none"
          />
          <div>
            <p className="text-sm font-semibold text-gray-600 mb-3">
              PIN <span className="text-gray-400 font-normal">(6 digits)</span>
            </p>
            <PinInput value={addPin} onChange={setAddPin} maxLength={6} />
          </div>
          <p className="text-sm text-gray-400 text-center">Share the username and PIN with the family member verbally.</p>
          <Button
            size="lg"
            fullWidth
            onClick={handleAdd}
            disabled={formLoading || !form.username || addPin.length !== 6 || !form.display_name}
          >
            {formLoading ? 'Adding…' : 'Add Member'}
          </Button>
        </div>
      </Modal>

      {/* Edit member modal */}
      <Modal open={!!editMember} onClose={() => { setEditMember(null); resetForm(); }} title="Edit Member">
        <div className="flex flex-col gap-4">
          <Input
            label="Display Name"
            value={form.display_name}
            onChange={e => setForm(f => ({ ...f, display_name: e.target.value }))}
          />
          <div>
            <p className="text-sm font-semibold text-gray-600 mb-1">
              New PIN <span className="text-gray-400 font-normal">(leave blank to keep current)</span>
            </p>
            {editPin.length === 0 ? (
              <button
                type="button"
                onClick={() => setEditPin(' ')}
                className="w-full py-3 text-sm text-blue-600 border-2 border-dashed border-blue-200 rounded-xl hover:bg-blue-50 transition-colors"
              >
                Tap to set a new PIN
              </button>
            ) : (
              <PinInput value={editPin.trim()} onChange={v => setEditPin(v)} maxLength={6} />
            )}
          </div>
          <Button
            size="lg"
            fullWidth
            onClick={handleUpdate}
            disabled={formLoading || (editPin.trim().length > 0 && editPin.trim().length !== 6)}
          >
            {formLoading ? 'Saving…' : 'Save Changes'}
          </Button>
        </div>
      </Modal>
    </AppShell>
  );
}
