import { useEffect, useState } from 'react';
import { AreaChart, Area, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from 'recharts';
import { ShieldCheck, ArrowUpRight, ArrowDownRight, Activity } from 'lucide-react';
import './index.css';

export default function App() {
  const [data, setData] = useState<any>(null);
  const [selectedField, setSelectedField] = useState<string | null>(null);

  useEffect(() => {
    const fetchLiveTelemetry = async () => {
      try {
        const response = await fetch('http://localhost:8085/v1/stats');
        const json = await response.json();
        
        // Failsafe mappings to gracefully handle empty null states explicitly safely 
        if (!json.heatmap || json.heatmap.length === 0) {
            json.heatmap = [{ name: "No errors recorded yet", failures: 0 }];
        }
        
        setData(json);
      } catch (err) {
        console.error("SchemaGuard Node Connection Severed:", err);
      }
    };
    
    fetchLiveTelemetry(); // Initial physical hook
    const interval = setInterval(fetchLiveTelemetry, 2000); // 2s live operational polling loop mapping dynamically 
    
    return () => clearInterval(interval);
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
          <p className="value">{data.kpi.pass_rate.toFixed(1)}%</p>
          <div className="trend up" style={{ color: "var(--text-secondary)" }}><ArrowUpRight size={16} /> Live passing metrics computed directly over internal network traffic</div>
        </div>
        <div className="glass-panel kpi-card">
          <h3>Total Executions (30d)</h3>
          <p className="value">{data.kpi.validations_month.toLocaleString()}</p>
          <div className="trend up" style={{ color: "var(--text-secondary)" }}><ArrowUpRight size={16} /> Absolute network connections strictly hitting the Go verification engine</div>
        </div>
        <div className="glass-panel kpi-card">
          <h3>Token Overhead Cost Saved</h3>
          <p className="value">${data.kpi.cost_saved.toFixed(2)}</p>
          <div className="trend down" style={{ color: "var(--text-secondary)" }}>
            <ArrowDownRight size={16} /> 
            {data.kpi.cost_saved > 0 
              ? "Short-circuited failures safely preventing downstream LLM cascade expenses" 
              : "Monitoring boundary layers seamlessly for structural formatting drifts"}
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
              <div key={idx} className="heatmap-item" onClick={() => setSelectedField(item.name)}>
                <span className="field">{item.name}</span>
                <span className="count">{item.failures}</span>
              </div>
            ))}
          </div>
        </div>
      </div>

      {selectedField && (
        <div className="glass-panel" style={{ position: 'fixed', bottom: '2rem', right: '2rem', zIndex: 1000, borderLeft: '4px solid var(--accent)', animation: 'slideIn 0.3s ease-out' }}>
          <h4 style={{ margin: '0 0 0.5rem 0', display: 'flex', alignItems: 'center', gap: '8px' }}>
            <Activity size={16} color="var(--accent)" /> Deep Trace: <span style={{ fontFamily: 'monospace', color: '#ff7b72' }}>{selectedField}</span>
          </h4>
          <p style={{ margin: 0, fontSize: '0.9rem', color: 'var(--text-secondary)', maxWidth: '300px', lineHeight: '1.4' }}>
            Historical payload forensics and Deep Tracing requires <b>SchemaGuard Enterprise</b>.
            <br />
            <button 
              onClick={() => setSelectedField(null)} 
              style={{ background: 'transparent', border: 'none', color: 'var(--accent)', marginTop: '0.8rem', cursor: 'pointer', padding: 0, fontWeight: 600 }}
            >
              Dismiss
            </button>
          </p>
        </div>
      )}
    </div>
  );
}
