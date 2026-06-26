import React, { useState } from 'react';
import { Search } from 'lucide-react';

const AgentGrid = ({ agents }) => {
  const [searchTerm, setSearchTerm] = useState('');

  if (agents.length === 0) {
    return <div className="no-data">Waiting for agent telemetry...</div>;
  }

  const filtered = agents.filter(a => a.agent_id.toLowerCase().includes(searchTerm.toLowerCase()));

  return (
    <div className="dashboard-section">
      <div className="section-header">
        <h2>Agent Nodes</h2>
        <div className="search-box">
          <Search size={18} />
          <input 
            type="text" 
            placeholder="Search nodes..." 
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
          />
        </div>
      </div>
      <div className="table-container">
        <table>
          <thead>
            <tr>
              <th>Agent ID</th>
              <th>Status</th>
              <th>Uptime</th>
              <th>CPU (%)</th>
              <th>Memory (%)</th>
              <th>Disk I/O (MB/s)</th>
              <th>Network (MB/s)</th>
            </tr>
          </thead>
          <tbody>
            {filtered.map(agent => (
              <tr key={agent.agent_id}>
                <td><strong>{agent.agent_id}</strong></td>
                <td>
                  <div className="uptime-badge">
                    <span className={`status-dot ${agent.status.toLowerCase() === 'healthy' ? 'up' : 'down'}`}></span>
                    {agent.status}
                  </div>
                </td>
                <td>{Math.floor(agent.uptime / 3600)}h {Math.floor((agent.uptime % 3600) / 60)}m</td>
                <td>{agent.cpu_utilization.toFixed(1)}%</td>
                <td>{agent.memory_utilization.toFixed(1)}%</td>
                <td>{agent.disk_io.toFixed(1)}</td>
                <td>{(agent.network_ingress + agent.network_egress).toFixed(1)}</td>
              </tr>
            ))}
            {filtered.length === 0 && (
              <tr><td colSpan="7" className="text-center">No nodes found.</td></tr>
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
};

export default AgentGrid;
