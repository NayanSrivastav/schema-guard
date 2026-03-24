import { useEffect, useState } from 'react';
import { AreaChart, Area, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from 'recharts';
import { ShieldCheck, ArrowUpRight, ArrowDownRight, Activity } from 'lucide-react';
import './index.css';

export default function App() {
  const [data, setData] = useState<any>(null);

  useEffect(() => {
    // In production, this hits the real Go API endpoints defined in `api/server.go`.
    // Simulating REST API latency fetching for local rendering MVP
    setTimeout(() => {
        setData({
            kpi: {
                pass_rate: 92.4,
                validations_month: 42104,
                cost_saved: 842.10
            },
            heatmap: [
                { name: "user.address.zipcode", failures: 412 },
                { name: "items[].price", failures: 234 },
                { name: "metadata.tags", failures: 182 },
                { name: "config.flags", failures: 89 },
                { name: "id", failures: 45 },
            ],
            timeseries: [
                { day: "Mon", success: 4000, errors: 300 },
                { day: "Tue", success: 5200, errors: 420 },
                { day: "Wed", success: 6120, errors: 390 },
                { day: "Thu", success: 5900, errors: 800 },
                { day: "Fri", success: 7200, errors: 210 },
                { day: "Sat", success: 6800, errors: 190 },
                { day: "Sun", success: 6900, errors: 150 },
            ]
        });
    }, 500);
  }, []);

  if (!data) return <div style={{ color: "white", padding: "2rem", display: "flex", alignItems: "center", gap: "10px" }}><Activity className="lucide-spin" /> Initializing Telemetry Engine...</div>;

  return (
    <div className="dashboard-container">
      <header>
        <div className="title">
          <ShieldCheck size={28} color="var(--accent)" />
          <span>SchemaGuard Validation Telemetry</span>
        </div>
        <div style={{ display: "flex", gap: "1rem" }}>
            <span style={{ color: "var(--success)", display: "flex", alignItems: "center", gap: "0.4rem", fontSize: "0.9rem", fontWeight: 600 }}>
                <Activity size={16} /> Live Traffic Stream
            </span>
        </div>
      </header>

      <div className="kpi-grid">
        <div className="glass-panel kpi-card">
          <h3>Overall Pass Rate (7d)</h3>
          <p className="value">{data.kpi.pass_rate}%</p>
          <div className="trend up"><ArrowUpRight size={16} /> 2.1% improvement directly from Prompt Coercion fixes</div>
        </div>
        <div className="glass-panel kpi-card">
          <h3>Total Executions (30d)</h3>
          <p className="value">{data.kpi.validations_month.toLocaleString()}</p>
          <div className="trend up"><ArrowUpRight size={16} /> 12.4% traffic scaling vs previous</div>
        </div>
        <div className="glass-panel kpi-card">
          <h3>Token Overhead Cost Saved</h3>
          <p className="value">${data.kpi.cost_saved.toFixed(2)}</p>
          <div className="trend down" style={{ color: "var(--text-secondary)" }}>
            <ArrowDownRight size={16} /> Circuit breaker successfully prevented runaway cascade expenses
          </div>
        </div>
      </div>

      <div className="charts-grid">
        <div className="glass-panel">
          <div className="chart-header">Validation Throughput & Reliability Metrics</div>
          <div style={{ width: '100%', height: 320 }}>
            <ResponsiveContainer>
              <AreaChart data={data.timeseries} margin={{ top: 10, right: 30, left: 0, bottom: 0 }}>
                <defs>
                  <linearGradient id="colorSuccess" x1="0" y1="0" x2="0" y2="1">
                    <stop offset="5%" stopColor="var(--accent)" stopOpacity={0.8}/>
                    <stop offset="95%" stopColor="var(--accent)" stopOpacity={0}/>
                  </linearGradient>
                  <linearGradient id="colorError" x1="0" y1="0" x2="0" y2="1">
                    <stop offset="5%" stopColor="var(--danger)" stopOpacity={0.6}/>
                    <stop offset="95%" stopColor="var(--danger)" stopOpacity={0}/>
                  </linearGradient>
                </defs>
                <CartesianGrid strokeDasharray="3 3" stroke="rgba(255,255,255,0.06)" vertical={false} />
                <XAxis dataKey="day" stroke="var(--text-secondary)" tick={{fill: 'var(--text-secondary)', fontSize: 12}} dy={10} axisLine={false} tickLine={false} />
                <YAxis stroke="var(--text-secondary)" tick={{fill: 'var(--text-secondary)', fontSize: 12}} dx={-10} axisLine={false} tickLine={false} />
                <Tooltip 
                  contentStyle={{ backgroundColor: 'rgba(22, 27, 34, 0.95)', borderColor: 'var(--border)', borderRadius: '8px', backdropFilter: 'blur(10px)' }}
                  itemStyle={{ color: 'var(--text-primary)', fontWeight: 600 }}
                  labelStyle={{ color: 'var(--text-secondary)', marginBottom: '4px' }}
                />
                <Area type="monotone" dataKey="success" stroke="var(--accent)" strokeWidth={3} fillOpacity={1} fill="url(#colorSuccess)" name="Enforced Passed" />
                <Area type="monotone" dataKey="errors" stroke="var(--danger)" strokeWidth={3} fillOpacity={1} fill="url(#colorError)" name="Rejections" />
              </AreaChart>
            </ResponsiveContainer>
          </div>
        </div>

        <div className="glass-panel">
          <div className="chart-header">Real-Time Failure Heatmap</div>
          <p style={{ fontSize: "0.85rem", color: "var(--text-secondary)", marginBottom: "1.2rem", lineHeight: "1.4" }}>
            These are the top deeply-nested schema properties failing LLM outputs causing internal retry spins.
          </p>
          <div className="heatmap-list">
            {data.heatmap.map((item: any, idx: number) => (
              <div key={idx} className="heatmap-item">
                <span className="field">{item.name}</span>
                <span className="count">{item.failures}</span>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
}
