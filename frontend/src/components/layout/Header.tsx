import { LogOut, Wallet } from 'lucide-react';
import { useAuth } from '../../context/AuthContext';
import { useNavigate } from 'react-router-dom';

interface HeaderProps {
  title?: string;
}

export function Header({ title }: HeaderProps) {
  const { user, logout } = useAuth();
  const navigate = useNavigate();

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  return (
    <header className="sticky top-0 bg-white border-b border-gray-200 z-30 shadow-sm">
      <div className="flex items-center justify-between px-4 h-14">
        <div className="flex items-center gap-2">
          <Wallet size={22} className="text-blue-600" />
          <span className="text-lg font-bold text-gray-900">
            {title || 'Family Wallet'}
          </span>
        </div>
        <div className="flex items-center gap-3">
          <span className="text-sm text-gray-500 hidden sm:block">
            Hi, <strong className="text-gray-700">{user?.display_name}</strong>
          </span>
          <button
            onClick={handleLogout}
            className="p-2 rounded-xl hover:bg-gray-100 text-gray-500 transition-colors"
            title="Logout"
          >
            <LogOut size={20} />
          </button>
        </div>
      </div>
    </header>
  );
}
