import React, { useState, useEffect } from 'react';

const Header = () => {
  const [health, setHealth] = useState({ database: 'checking...', redis: 'checking...' });

  useEffect(() => {
    const checkHealth = async () => {
      try {
        const host = window.location.hostname || 'localhost';
        const res = await fetch(`http://${host}:8000/api/health`);
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
    window.location.href = `http://${host}:8000/api/backup`;
  };

  return (
    <header className="app-header">
      <div className="header-content">
        <div>
          <h1>Smart Process Manager</h1>
          <p>Real-time edge telemetry and process visibility</p>
        </div>
        <div className="system-status">
          <div className="status-item">
            <span className={`status-dot ${health.database === 'up' ? 'up' : 'down'}`}></span>
            DB: {health.database}
          </div>
          <div className="status-item">
            <span className={`status-dot ${health.redis === 'up' ? 'up' : 'down'}`}></span>
            Cache: {health.redis}
          </div>
          <button className="btn-primary" onClick={handleBackup}>Download DB Backup (CSV)</button>
        </div>
      </div>
    </header>
  );
};

export default Header;
