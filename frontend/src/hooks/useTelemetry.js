import { useState, useEffect } from 'react';
import { createWebSocketConnection } from '../services/websocket';

export const useTelemetry = () => {
  const [agents, setAgents] = useState({});

  useEffect(() => {
    const ws = createWebSocketConnection((data) => {
      setAgents(prev => ({
        ...prev,
        [data.agent_id]: data
      }));
    });

    return () => {
      ws.close();
    };
  }, []);

  return Object.values(agents);
};
