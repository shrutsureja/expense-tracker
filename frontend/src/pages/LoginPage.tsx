import { useState, FormEvent } from 'react';
import { useNavigate } from 'react-router-dom';
import { Wallet } from 'lucide-react';
import { useAuth } from '../context/AuthContext';
import { Input } from '../components/ui/Input';
import { Button } from '../components/ui/Button';
import { PinInput } from '../components/auth/PinInput';

export function LoginPage() {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [pin, setPin] = useState('');
  const [step, setStep] = useState<'username' | 'pin' | 'password'>('username');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const { login } = useAuth();
  const navigate = useNavigate();

  // Simple heuristic: if username looks like a "family member" (short), use PIN
  // In practice, attempt login and if it fails, try the other
  const isLikelyPin = (un: string) => un !== 'shrutsureja' && un.length <= 10;

  const handleUsernameSubmit = (e: FormEvent) => {
    e.preventDefault();
    if (!username.trim()) return;
    setError('');
    setStep(isLikelyPin(username) ? 'pin' : 'password');
  };

  const handleLogin = async (cred: string) => {
    setLoading(true);
    setError('');
    try {
      await login(username.trim(), cred);
      navigate('/', { replace: true });
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : 'Login failed';
      setError(msg);
    } finally {
      setLoading(false);
    }
  };

  const handlePasswordSubmit = (e: FormEvent) => {
    e.preventDefault();
    if (!password) return;
    handleLogin(password);
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100 flex flex-col items-center justify-center p-6">
      <div className="bg-white rounded-3xl shadow-xl w-full max-w-sm p-8">
        {/* Logo */}
        <div className="flex flex-col items-center mb-8">
          <div className="bg-blue-600 rounded-2xl p-4 mb-3 shadow-md">
            <Wallet size={36} className="text-white" />
          </div>
          <h1 className="text-2xl font-bold text-gray-900">Family Wallet</h1>
          <p className="text-gray-500 text-sm mt-1">Track family expenses together</p>
        </div>

        {/* Step 1: Username */}
        {step === 'username' && (
          <form onSubmit={handleUsernameSubmit} className="flex flex-col gap-5">
            <Input
              label="Who are you?"
              placeholder="Enter your username"
              value={username}
              onChange={e => setUsername(e.target.value)}
              autoFocus
              autoCapitalize="none"
              autoCorrect="off"
            />
            {error && <p className="text-red-500 text-sm text-center">{error}</p>}
            <Button type="submit" size="lg" fullWidth>
              Continue →
            </Button>
          </form>
        )}

        {/* Step 2a: PIN */}
        {step === 'pin' && (
          <div className="flex flex-col gap-5">
            <p className="text-center text-gray-600 text-base">
              Enter your PIN, <strong>{username}</strong>
            </p>
            <PinInput value={pin} onChange={setPin} maxLength={6} />
            {error && <p className="text-red-500 text-sm text-center">{error}</p>}
            <Button
              size="lg"
              fullWidth
              onClick={() => handleLogin(pin)}
              disabled={pin.length < 4 || loading}
            >
              {loading ? 'Signing in…' : 'Sign In'}
            </Button>
            <button
              type="button"
              onClick={() => { setStep('username'); setPin(''); setError(''); }}
              className="text-sm text-gray-400 hover:text-gray-600 text-center"
            >
              ← Back
            </button>
          </div>
        )}

        {/* Step 2b: Password */}
        {step === 'password' && (
          <form onSubmit={handlePasswordSubmit} className="flex flex-col gap-5">
            <p className="text-center text-gray-600 text-base">
              Enter password for <strong>{username}</strong>
            </p>
            <Input
              label="Password"
              type="password"
              placeholder="Enter password"
              value={password}
              onChange={e => setPassword(e.target.value)}
              autoFocus
            />
            {error && <p className="text-red-500 text-sm text-center">{error}</p>}
            <Button type="submit" size="lg" fullWidth disabled={loading}>
              {loading ? 'Signing in…' : 'Sign In'}
            </Button>
            <button
              type="button"
              onClick={() => { setStep('username'); setPassword(''); setError(''); }}
              className="text-sm text-gray-400 hover:text-gray-600 text-center"
            >
              ← Back
            </button>
          </form>
        )}
      </div>
    </div>
  );
}
