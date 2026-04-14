import { Delete } from 'lucide-react';

interface PinInputProps {
  value: string;
  onChange: (val: string) => void;
  maxLength?: number;
}

export function PinInput({ value, onChange, maxLength = 6 }: PinInputProps) {
  const handleKey = (digit: string) => {
    if (value.length < maxLength) onChange(value + digit);
  };

  const handleDelete = () => onChange(value.slice(0, -1));

  const keys = ['1', '2', '3', '4', '5', '6', '7', '8', '9', '', '0', 'del'];

  return (
    <div className="flex flex-col items-center gap-6">
      {/* PIN dots display */}
      <div className="flex gap-3">
        {Array.from({ length: maxLength }).map((_, i) => (
          <div
            key={i}
            className={`w-4 h-4 rounded-full transition-all ${
              i < value.length ? 'bg-blue-600 scale-110' : 'bg-gray-300'
            }`}
          />
        ))}
      </div>

      {/* Keypad */}
      <div className="grid grid-cols-3 gap-3 w-full max-w-[280px]">
        {keys.map((key, i) => {
          if (key === '') return <div key={i} />;
          if (key === 'del') {
            return (
              <button
                key={i}
                type="button"
                onClick={handleDelete}
                className="flex items-center justify-center h-16 rounded-2xl bg-gray-100 hover:bg-gray-200 active:scale-95 transition-all text-gray-600"
              >
                <Delete size={22} />
              </button>
            );
          }
          return (
            <button
              key={i}
              type="button"
              onClick={() => handleKey(key)}
              className="flex items-center justify-center h-16 rounded-2xl bg-white border-2 border-gray-200 hover:bg-blue-50 hover:border-blue-300 active:scale-95 transition-all text-2xl font-bold text-gray-800 shadow-sm"
            >
              {key}
            </button>
          );
        })}
      </div>
    </div>
  );
}
