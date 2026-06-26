import React from 'react';
import { ShieldAlert, Cpu, Thermometer, Clock } from 'lucide-react';

const Overview = ({ agents }) => {
  const totalNodes = agents.length;
  const avgCpu = totalNodes ? agents.reduce((acc, a) => acc + a.cpu_utilization, 0) / totalNodes : 0;
  const avgTemp = totalNodes ? agents.reduce((acc, a) => acc + a.temperature, 0) / totalNodes : 0;
  
  const warnings = agents.filter(a => a.status === 'Warning');

  return (
    <div className="overview-container">
      <div className="stats-grid">
        <div className="stat-card">
          <div className="stat-icon"><Cpu size={24} /></div>
          <div className="stat-info">
            <h4>Total Active Nodes</h4>
            <h2>{totalNodes}</h2>
          </div>
        </div>
        <div className="stat-card">
          <div className="stat-icon"><Clock size={24} /></div>
          <div className="stat-info">
            <h4>Avg Fleet CPU</h4>
            <h2>{avgCpu.toFixed(1)}%</h2>
          </div>
        </div>
        <div className="stat-card">
          <div className="stat-icon"><Thermometer size={24} /></div>
          <div className="stat-info">
            <h4>Avg Fleet Temp</h4>
            <h2>{avgTemp.toFixed(1)}°C</h2>
          </div>
        </div>
        <div className={`stat-card ${warnings.length > 0 ? 'warning-bg' : ''}`}>
          <div className="stat-icon"><ShieldAlert size={24} /></div>
          <div className="stat-info">
            <h4>Active Alerts</h4>
            <h2>{warnings.length}</h2>
          </div>
        </div>
      </div>

      <div className="dashboard-section">
        <h3>System Alerts</h3>
        {warnings.length === 0 ? (
          <div className="no-alerts">All systems operating normally.</div>
        ) : (
          <ul className="alerts-list">
            {warnings.map(w => (
              <li key={w.agent_id}>
                <strong>{w.agent_id}</strong> is running hot (Temp: {w.temperature.toFixed(1)}°C, CPU: {w.cpu_utilization.toFixed(1)}%)
              </li>
            ))}
          </ul>
        )}
      </div>
    </div>
  );
};

export default Overview;
