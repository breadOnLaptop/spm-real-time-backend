import React, { useState, useEffect } from 'react';
import { Database, Download, CheckCircle, XCircle } from 'lucide-react';

const DatabaseView = () => {
  const [health, setHealth] = useState({ database: 'checking...', redis: 'checking...' });

  useEffect(() => {
    const checkHealth = async () => {
      try {
        const host = window.location.hostname || 'localhost';
        const defaultApi = window.location.protocol === 'https:' ? `https://${host}/api` : `http://${host}:8000/api`;
        const apiUrl = import.meta.env.VITE_API_URL || defaultApi;
        
        const res = await fetch(`${apiUrl}/health`);
        const data = await res.json();
        setHealth(data);
      } catch (e) {
        setHealth({ database: 'down', redis: 'down' });
      }
    };
    checkHealth();
    const interval = setInterval(checkHealth, 10000);
    return () => clearInterval(interval);
  }, []);

  const handleBackup = () => {
    const host = window.location.hostname || 'localhost';
    const defaultApi = window.location.protocol === 'https:' ? `https://${host}/api` : `http://${host}:8000/api`;
    const apiUrl = import.meta.env.VITE_API_URL || defaultApi;
    window.location.href = `${apiUrl}/backup`;
  };

  return (
    <div className="database-view">
      <div className="stats-grid">
        <div className="stat-card">
          <div className="stat-info">
            <h4>PostgreSQL Status</h4>
            <h2 className="flex-align">
              {health.database === 'up' ? <CheckCircle color="var(--success)" /> : <XCircle color="var(--danger)" />}
              <span className="ml-2">{health.database.toUpperCase()}</span>
            </h2>
          </div>
        </div>
        <div className="stat-card">
          <div className="stat-info">
            <h4>Redis Cache Status</h4>
            <h2 className="flex-align">
              {health.redis === 'up' ? <CheckCircle color="var(--success)" /> : <XCircle color="var(--danger)" />}
              <span className="ml-2">{health.redis.toUpperCase()}</span>
            </h2>
          </div>
        </div>
      </div>

      <div className="dashboard-section mt-4">
        <div className="section-header">
          <h2>Storage Management</h2>
        </div>
        <div className="backup-panel">
          <Database size={48} color="var(--primary)" />
          <div>
            <h3>Export Telemetry History</h3>
            <p>Download a complete CSV backup of all agent telemetry stored in the relational database.</p>
          </div>
          <button className="btn-primary flex-align" onClick={handleBackup}>
            <Download size={18} className="mr-2" /> Download Backup (CSV)
          </button>
        </div>
      </div>
    </div>
  );
};

export default DatabaseView;
