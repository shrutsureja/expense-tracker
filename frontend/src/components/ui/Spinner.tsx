export function Spinner({ className = '' }: { className?: string }) {
  return (
    <div
      className={`inline-block w-6 h-6 border-2 border-current border-t-transparent rounded-full animate-spin ${className}`}
    />
  );
}

export function FullPageSpinner() {
  return (
    <div className="min-h-screen flex items-center justify-center bg-slate-50">
      <Spinner className="text-blue-600 w-10 h-10 border-4" />
    </div>
  );
}
