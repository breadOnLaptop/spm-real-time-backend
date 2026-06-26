import React, { useState } from 'react';
import { Search } from 'lucide-react';

const ProcessList = ({ agents }) => {
  const [nodeSearch, setNodeSearch] = useState('');
  const [procSearch, setProcSearch] = useState('');
  
  const filteredAgents = agents.filter(a => a.agent_id.toLowerCase().includes(nodeSearch.toLowerCase()));

  return (
    <div>
      <div className="dashboard-section" style={{ marginBottom: '1rem' }}>
        <div className="section-header" style={{ marginBottom: 0, border: 'none', paddingBottom: 0 }}>
          <h2 style={{ display: 'flex', gap: '15px' }}>
            <div className="search-box">
              <Search size={18} />
              <input 
                type="text" 
                placeholder="Filter Nodes..." 
                value={nodeSearch}
                onChange={(e) => setNodeSearch(e.target.value)}
              />
            </div>
            <div className="search-box">
              <Search size={18} />
              <input 
                type="text" 
                placeholder="Filter Processes..." 
                value={procSearch}
                onChange={(e) => setProcSearch(e.target.value)}
              />
            </div>
          </h2>
        </div>
      </div>

      {filteredAgents.length > 0 ? filteredAgents.map(agent => {
        const procs = (agent.top_processes || []).filter(p => p.executable_name.toLowerCase().includes(procSearch.toLowerCase()));
        
        return (
          <div key={agent.agent_id} className="dashboard-section" style={{ marginBottom: '1rem', padding: '1rem' }}>
            <div style={{ marginBottom: '0.75rem', paddingBottom: '0.5rem', borderBottom: '1px solid #eaeded' }}>
              <strong>{agent.agent_id}</strong> - <span style={{ color: '#5f6b7a', fontSize: '0.85rem' }}>{procs.length} processes matching</span>
            </div>
            {procs.length > 0 ? (
              <div className="table-container">
                <table style={{ margin: 0 }}>
                  <thead>
                    <tr>
                      <th style={{ background: 'transparent' }}>PID</th>
                      <th style={{ background: 'transparent' }}>Executable</th>
                      <th style={{ background: 'transparent' }}>Resource Utilization</th>
                    </tr>
                  </thead>
                  <tbody>
                    {procs.map((p, idx) => (
                      <tr key={idx}>
                        <td style={{ width: '150px' }}>{p.pid}</td>
                        <td style={{ width: '250px' }}><strong>{p.executable_name}</strong></td>
                        <td>
                          <div className="progress-bar">
                            <div className="progress-fill" style={{width: `${Math.min(p.resource_utilization, 100)}%`}}></div>
                          </div>
                          <span className="progress-text">{p.resource_utilization.toFixed(1)}%</span>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            ) : (
              <div style={{ padding: '0.5rem', color: '#5f6b7a', fontSize: '0.85rem' }}>No matching processes on this node.</div>
            )}
          </div>
        );
      }) : (
        <div className="dashboard-section text-center">No nodes match your filter.</div>
      )}
    </div>
  );
};

export default ProcessList;
