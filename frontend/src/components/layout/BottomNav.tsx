import { NavLink } from 'react-router-dom';
import { Home, PlusCircle, BarChart2, Users, ShieldCheck } from 'lucide-react';
import { useAuth } from '../../context/AuthContext';

export function BottomNav() {
  const { user } = useAuth();

  const navItems = [
    { to: '/', icon: Home, label: 'Home' },
    { to: '/add', icon: PlusCircle, label: 'Add', primary: true },
    { to: '/dashboard', icon: BarChart2, label: 'Dashboard' },
    ...(user?.role === 'family_owner' || user?.role === 'super_admin'
      ? [{ to: '/family', icon: Users, label: 'Family' }]
      : []),
    ...(user?.role === 'super_admin'
      ? [{ to: '/admin', icon: ShieldCheck, label: 'Admin' }]
      : []),
  ];

  return (
    <nav className="fixed bottom-0 left-0 right-0 bg-white border-t border-gray-200 safe-area-pb z-40 shadow-lg">
      <div className="flex items-center justify-around px-2 h-16">
        {navItems.map(({ to, icon: Icon, label, primary }) => (
          <NavLink
            key={to}
            to={to}
            end={to === '/'}
            className={({ isActive }) => `
              flex flex-col items-center justify-center gap-0.5 px-3 py-1 rounded-xl
              min-w-[56px] min-h-[56px] transition-all
              ${primary
                ? isActive
                  ? 'text-blue-600'
                  : 'text-gray-400 hover:text-blue-500'
                : isActive
                  ? 'text-blue-600'
                  : 'text-gray-400 hover:text-gray-600'
              }
            `}
          >
            <Icon size={primary ? 28 : 22} strokeWidth={primary ? 2.5 : 2} />
            <span className={`text-xs font-medium ${primary ? 'text-[11px]' : ''}`}>{label}</span>
          </NavLink>
        ))}
      </div>
    </nav>
  );
}
