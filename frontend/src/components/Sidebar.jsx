import React from 'react';
import { LayoutDashboard, Server, Activity, Database, PlusCircle } from 'lucide-react';

const Sidebar = ({ activeTab, setActiveTab }) => {
  return (
    <aside className="sidebar">
      <div className="logo">
        <Activity color="#3b82f6" size={28} />
        <h2>SPM</h2>
      </div>
      <nav>
        <button className={activeTab === 'overview' ? 'active' : ''} onClick={() => setActiveTab('overview')}>
          <LayoutDashboard size={20} /> Overview
        </button>
        <button className={activeTab === 'nodes' ? 'active' : ''} onClick={() => setActiveTab('nodes')}>
          <Server size={20} /> Nodes
        </button>
        <button className={activeTab === 'processes' ? 'active' : ''} onClick={() => setActiveTab('processes')}>
          <Activity size={20} /> Processes
        </button>
        <button className={activeTab === 'database' ? 'active' : ''} onClick={() => setActiveTab('database')}>
          <Database size={20} /> Storage & Backup
        </button>
      </nav>
      <div className="sidebar-bottom">
        <button className="btn-add-node" onClick={() => alert("Simulating adding a new node... Run install.sh on target machine.")}>
          <PlusCircle size={20} /> Add Node
        </button>
      </div>
    </aside>
  );
};

export default Sidebar;
