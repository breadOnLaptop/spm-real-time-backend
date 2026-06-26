import React, { useState } from 'react';
import MetricWidget from './MetricWidget';
import HistoryModal from './HistoryModal';
import { Activity, Clock } from 'lucide-react';

const AgentCard = ({ agent }) => {
  const [historyOpen, setHistoryOpen] = useState(false);
  const [historyData, setHistoryData] = useState([]);

  const formatUptime = (seconds) => {
    const h = Math.floor(seconds / 3600);
    const m = Math.floor((seconds % 3600) / 60);
    return `${h}h ${m}m`;
  };

  const loadHistory = async () => {
    try {
      const host = window.location.hostname || 'localhost';
      const res = await fetch(`http://${host}:8000/api/history?agent_id=${agent.agent_id}`);
      const data = await res.json();
      setHistoryData(data || []);
      setHistoryOpen(true);
    } catch (e) {
      alert("Failed to load history");
    }
  };

  return (
    <div className={`agent-card ${agent.status === 'Warning' ? 'border-warning' : ''}`}>
      <div className="agent-header">
        <h2>
          <span className={`status-dot ${agent.status === 'Warning' ? 'down' : 'up'}`}></span>
          {agent.agent_id}
        </h2>
        <span className="uptime-badge"><Clock size={12}/> {formatUptime(agent.uptime)}</span>
      </div>
      
      <div className="metrics">
        <MetricWidget label="CPU" value={`${agent.cpu_utilization.toFixed(1)}%`} />
        <MetricWidget label="RAM" value={`${agent.memory_utilization.toFixed(1)}%`} />
        <MetricWidget label="Temp" value={`${agent.temperature?.toFixed(1)}°C`} />
      </div>

      <div className="actions">
        <button className="btn-secondary" onClick={loadHistory}><Activity size={16}/> View History</button>
      </div>

      <HistoryModal isOpen={historyOpen} onClose={() => setHistoryOpen(false)} agentId={agent.agent_id} history={historyData} />
    </div>
  );
};

export default AgentCard;
