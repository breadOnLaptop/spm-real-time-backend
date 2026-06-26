import React, { useState } from 'react';
import { useTelemetry } from './hooks/useTelemetry';
import Sidebar from './components/Sidebar';
import Overview from './components/Dashboard/Overview';
import AgentGrid from './components/Dashboard/AgentGrid';
import ProcessList from './components/Dashboard/ProcessList';
import DatabaseView from './components/Dashboard/DatabaseView';
import { Menu } from 'lucide-react';
import './styles/App.css';

function App() {
  const agents = useTelemetry();
  const [activeTab, setActiveTab] = useState('overview');
  const [sidebarOpen, setSidebarOpen] = useState(true);

  const renderContent = () => {
    switch (activeTab) {
      case 'overview': return <Overview agents={agents} />;
      case 'nodes': return <AgentGrid agents={agents} />;
      case 'processes': return <ProcessList agents={agents} />;
      case 'database': return <DatabaseView />;
      default: return <Overview agents={agents} />;
    }
  };

  return (
    <div className="app-layout">
      <div className={`sidebar ${!sidebarOpen ? 'hidden' : ''}`}>
        <Sidebar activeTab={activeTab} setActiveTab={setActiveTab} />
      </div>
      <main className="main-content">
        <header className="topbar" style={{ display: 'flex', alignItems: 'center', gap: '15px' }}>
          <button onClick={() => setSidebarOpen(!sidebarOpen)} style={{ background: 'none', border: 'none', cursor: 'pointer', padding: 0 }}>
            <Menu size={20} color="#232f3e" />
          </button>
          <h1>{activeTab.charAt(0).toUpperCase() + activeTab.slice(1)}</h1>
        </header>
        <div className="content-area" style={{ padding: '0', margin: '0' }}>
          <div style={{ padding: '1.5rem 2rem', height: '100%', overflowY: 'auto' }}>
            {renderContent()}
          </div>
        </div>
      </main>
    </div>
  );
}

export default App;
