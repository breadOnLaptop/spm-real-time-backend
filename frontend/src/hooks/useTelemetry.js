import { useState, useEffect } from 'react';
import { createWebSocketConnection } from '../services/websocket';

export const useTelemetry = () => {
  const [agents, setAgents] = useState({});

  useEffect(() => {
    // 1. Fetch initial state from the API cache
    const fetchInitial = async () => {
      try {
        const host = window.location.hostname || 'localhost';
        const defaultApi = window.location.protocol === 'https:' ? `https://${host}/api` : `http://${host}:8000/api`;
        const apiUrl = import.meta.env.VITE_API_URL || defaultApi;
        
        const res = await fetch(`${apiUrl}/agents`);
        const data = await res.json();
        
        if (data && Array.isArray(data)) {
          const initialState = {};
          const now = Date.now();
          data.forEach(agent => {
            // Assume initial load might be stale
            initialState[agent.agent_id] = { ...agent, last_seen: now };
          });
          setAgents(initialState);
        }
      } catch (e) {
        console.error("Failed to fetch initial agents", e);
      }
    };
    fetchInitial();

    // 2. Connect WebSocket for live updates
    const ws = createWebSocketConnection((data) => {
      setAgents(prev => ({
        ...prev,
        [data.agent_id]: { ...data, last_seen: Date.now() }
      }));
    });

    // 3. Heartbeat watchdog to mark disconnected agents as OFFLINE
    const interval = setInterval(() => {
      setAgents(prev => {
        const now = Date.now();
        let changed = false;
        const updated = { ...prev };
        
        for (const id in updated) {
          if (updated[id].status !== 'OFFLINE' && now - updated[id].last_seen > 5000) {
            updated[id] = { 
              ...updated[id], 
              status: 'OFFLINE',
              cpu_utilization: 0,
              memory_utilization: 0,
              disk_io: 0,
              network_ingress: 0,
              network_egress: 0,
              top_processes: []
            };
            changed = true;
          }
        }
        return changed ? updated : prev;
      });
    }, 2000);

    return () => {
      ws.close();
      clearInterval(interval);
    };
  }, []);

  return Object.values(agents);
};
