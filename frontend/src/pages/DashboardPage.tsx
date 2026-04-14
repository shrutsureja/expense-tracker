import { useState, useEffect } from 'react';
import {
  PieChart, Pie, Cell, Tooltip, ResponsiveContainer,
  LineChart, Line, XAxis, YAxis, CartesianGrid,
  BarChart, Bar, Legend,
} from 'recharts';
import { analyticsApi } from '../api/analytics';
import { AppShell } from '../components/layout/AppShell';
import { Spinner } from '../components/ui/Spinner';
import { formatCurrency, currentMonthStr, formatShortDate } from '../utils/formatters';
import { DATE_RANGE_LABELS, PAYMENT_METHODS } from '../utils/constants';
import type {
  DateRange, AnalyticsSummary, CategoryBreakdown,
  PersonBreakdown, DailyTotal, MonthlyComparison,
} from '../types';

const COLORS = ['#3b82f6', '#10b981', '#f59e0b', '#ef4444', '#8b5cf6', '#06b6d4', '#ec4899', '#84cc16'];

export function DashboardPage() {
  const [range, setRange] = useState<DateRange>('month');
  const [summary, setSummary] = useState<AnalyticsSummary | null>(null);
  const [byCategory, setByCategory] = useState<CategoryBreakdown[]>([]);
  const [byPerson, setByPerson] = useState<PersonBreakdown[]>([]);
  const [dailyTrend, setDailyTrend] = useState<DailyTotal[]>([]);
  const [monthly, setMonthly] = useState<MonthlyComparison | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    setLoading(true);
    Promise.all([
      analyticsApi.summary(range),
      analyticsApi.byCategory(range),
      analyticsApi.byPerson(range),
      analyticsApi.dailyTrend(currentMonthStr()),
      analyticsApi.monthlyComparison(),
    ]).then(([s, cat, per, daily, mon]) => {
      setSummary(s);
      setByCategory(cat);
      setByPerson(per);
      setDailyTrend(daily);
      setMonthly(mon);
    }).finally(() => setLoading(false));
  }, [range]);

  const pieData = byCategory.slice(0, 8).map(c => ({
    name: `${c.tag_icon} ${c.tag_name}`,
    value: c.total,
  }));

  const trendData = dailyTrend.map(d => ({
    date: formatShortDate(d.date),
    amount: d.total,
  }));

  return (
    <AppShell title="Dashboard">
      {/* Date range tabs */}
      <div className="flex gap-1 p-4 bg-white border-b border-gray-100 overflow-x-auto">
        {(Object.keys(DATE_RANGE_LABELS) as DateRange[]).filter(r => r !== 'custom').map(r => (
          <button
            key={r}
            onClick={() => setRange(r)}
            className={`
              px-4 py-2 rounded-xl text-sm font-semibold whitespace-nowrap transition-all
              ${range === r ? 'bg-blue-600 text-white' : 'bg-gray-100 text-gray-600 hover:bg-gray-200'}
            `}
          >
            {DATE_RANGE_LABELS[r]}
          </button>
        ))}
      </div>

      {loading ? (
        <div className="flex justify-center py-16"><Spinner className="text-blue-600 w-8 h-8 border-4" /></div>
      ) : (
        <div className="p-4 flex flex-col gap-5">
          {/* Summary cards */}
          {summary && (
            <div className="grid grid-cols-3 gap-3">
              <SummaryCard label="Total Spent" value={formatCurrency(summary.total_amount)} color="blue" />
              <SummaryCard label="Expenses" value={String(summary.count)} color="green" />
              <SummaryCard label="Avg/Day" value={formatCurrency(summary.avg_per_day)} color="purple" />
            </div>
          )}

          {/* Monthly comparison */}
          {monthly && (
            <div className="bg-white rounded-2xl p-4 shadow-sm">
              <h3 className="font-bold text-gray-700 mb-3">Month Comparison</h3>
              <div className="flex gap-4">
                <div className="flex-1 text-center p-3 bg-blue-50 rounded-xl">
                  <p className="text-xs text-gray-500 mb-1">This Month</p>
                  <p className="text-xl font-bold text-blue-700">{formatCurrency(monthly.this_month)}</p>
                </div>
                <div className="flex-1 text-center p-3 bg-gray-50 rounded-xl">
                  <p className="text-xs text-gray-500 mb-1">Last Month</p>
                  <p className="text-xl font-bold text-gray-700">{formatCurrency(monthly.last_month)}</p>
                </div>
              </div>
              {monthly.change_percent !== 0 && (
                <p className={`text-center text-sm font-medium mt-2 ${monthly.change_percent > 0 ? 'text-red-500' : 'text-green-600'}`}>
                  {monthly.change_percent > 0 ? '↑' : '↓'} {Math.abs(monthly.change_percent).toFixed(1)}% vs last month
                </p>
              )}
            </div>
          )}

          {/* Category pie chart */}
          {pieData.length > 0 && (
            <div className="bg-white rounded-2xl p-4 shadow-sm">
              <h3 className="font-bold text-gray-700 mb-3">By Category</h3>
              <ResponsiveContainer width="100%" height={200}>
                <PieChart>
                  <Pie
                    data={pieData}
                    cx="50%"
                    cy="50%"
                    innerRadius={50}
                    outerRadius={80}
                    dataKey="value"
                  >
                    {pieData.map((_, idx) => (
                      <Cell key={idx} fill={COLORS[idx % COLORS.length]} />
                    ))}
                  </Pie>
                  <Tooltip formatter={(v: number) => formatCurrency(v)} />
                </PieChart>
              </ResponsiveContainer>
              <div className="flex flex-col gap-1.5 mt-2">
                {byCategory.slice(0, 8).map((c, idx) => (
                  <div key={c.tag_id} className="flex items-center gap-2">
                    <div className="w-3 h-3 rounded-full flex-shrink-0" style={{ background: COLORS[idx % COLORS.length] }} />
                    <span className="text-sm text-gray-600 flex-1">{c.tag_icon} {c.tag_name}</span>
                    <span className="text-sm font-semibold text-gray-800">{formatCurrency(c.total)}</span>
                    <span className="text-xs text-gray-400">{c.percentage.toFixed(0)}%</span>
                  </div>
                ))}
              </div>
            </div>
          )}

          {/* Person breakdown */}
          {byPerson.length > 0 && (
            <div className="bg-white rounded-2xl p-4 shadow-sm">
              <h3 className="font-bold text-gray-700 mb-3">By Person</h3>
              <div className="flex flex-col gap-3">
                {byPerson.map(p => (
                  <div key={p.user_id}>
                    <div className="flex justify-between mb-1">
                      <span className="text-sm font-medium text-gray-700">{p.display_name}</span>
                      <span className="text-sm font-bold text-gray-900">{formatCurrency(p.total)}</span>
                    </div>
                    <div className="h-2.5 bg-gray-100 rounded-full overflow-hidden">
                      <div
                        className="h-full bg-blue-500 rounded-full transition-all"
                        style={{ width: `${p.percentage}%` }}
                      />
                    </div>
                  </div>
                ))}
              </div>
            </div>
          )}

          {/* Daily trend */}
          {trendData.length > 0 && (
            <div className="bg-white rounded-2xl p-4 shadow-sm">
              <h3 className="font-bold text-gray-700 mb-3">Daily Trend (This Month)</h3>
              <ResponsiveContainer width="100%" height={160}>
                <LineChart data={trendData}>
                  <CartesianGrid strokeDasharray="3 3" stroke="#f1f5f9" />
                  <XAxis dataKey="date" tick={{ fontSize: 10 }} interval="preserveStartEnd" />
                  <YAxis tick={{ fontSize: 10 }} tickFormatter={v => `₹${v}`} width={50} />
                  <Tooltip formatter={(v: number) => formatCurrency(v)} />
                  <Line type="monotone" dataKey="amount" stroke="#3b82f6" strokeWidth={2} dot={false} />
                </LineChart>
              </ResponsiveContainer>
            </div>
          )}

          {/* Payment method breakdown */}
          <PaymentMethodChart range={range} />
        </div>
      )}
    </AppShell>
  );
}

function SummaryCard({ label, value, color }: { label: string; value: string; color: 'blue' | 'green' | 'purple' }) {
  const colors = {
    blue: 'bg-blue-50 text-blue-700',
    green: 'bg-green-50 text-green-700',
    purple: 'bg-purple-50 text-purple-700',
  };
  return (
    <div className={`rounded-2xl p-3 text-center ${colors[color]}`}>
      <p className="text-xs font-medium opacity-70 mb-1">{label}</p>
      <p className="text-lg font-bold">{value}</p>
    </div>
  );
}

function PaymentMethodChart({ range }: { range: DateRange }) {
  const [data, setData] = useState<{ name: string; amount: number }[]>([]);

  useEffect(() => {
    analyticsApi.byPaymentMethod(range).then(methods => {
      setData(methods.map(m => {
        const pm = PAYMENT_METHODS.find(p => p.value === m.payment_method);
        return { name: `${pm?.icon || ''} ${pm?.label || m.payment_method}`, amount: m.total };
      }));
    });
  }, [range]);

  if (data.length === 0) return null;

  return (
    <div className="bg-white rounded-2xl p-4 shadow-sm">
      <h3 className="font-bold text-gray-700 mb-3">By Payment Method</h3>
      <ResponsiveContainer width="100%" height={140}>
        <BarChart data={data} layout="vertical">
          <XAxis type="number" tick={{ fontSize: 10 }} tickFormatter={v => `₹${v}`} />
          <YAxis type="category" dataKey="name" tick={{ fontSize: 11 }} width={80} />
          <Tooltip formatter={(v: number) => formatCurrency(v)} />
          <Legend />
          <Bar dataKey="amount" fill="#3b82f6" radius={[0, 6, 6, 0]} />
        </BarChart>
      </ResponsiveContainer>
    </div>
  );
}
