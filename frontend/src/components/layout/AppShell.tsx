import { type ReactNode } from 'react';
import { Header } from './Header';
import { BottomNav } from './BottomNav';

interface AppShellProps {
  children: ReactNode;
  title?: string;
}

export function AppShell({ children, title }: AppShellProps) {
  return (
    <div className="min-h-screen bg-slate-50 flex flex-col">
      <Header title={title} />
      <main className="flex-1 pb-20 max-w-2xl w-full mx-auto">
        {children}
      </main>
      <BottomNav />
    </div>
  );
}
