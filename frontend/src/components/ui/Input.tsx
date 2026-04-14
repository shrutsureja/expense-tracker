import { InputHTMLAttributes, forwardRef } from 'react';

interface InputProps extends InputHTMLAttributes<HTMLInputElement> {
  label?: string;
  error?: string;
}

export const Input = forwardRef<HTMLInputElement, InputProps>(
  ({ label, error, className = '', ...props }, ref) => (
    <div className="flex flex-col gap-1 w-full">
      {label && (
        <label className="text-sm font-semibold text-gray-600">{label}</label>
      )}
      <input
        ref={ref}
        className={`
          w-full px-4 py-3 text-lg rounded-xl border-2
          min-h-[52px] bg-white text-gray-900
          border-gray-200 focus:border-blue-500 focus:outline-none
          transition-colors
          disabled:bg-gray-50 disabled:text-gray-500
          ${error ? 'border-red-400 focus:border-red-500' : ''}
          ${className}
        `}
        {...props}
      />
      {error && <span className="text-sm text-red-500">{error}</span>}
    </div>
  )
);

Input.displayName = 'Input';
