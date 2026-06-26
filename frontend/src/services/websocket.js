export const createWebSocketConnection = (onMessage) => {
  const host = window.location.hostname || 'localhost';
  const defaultWs = window.location.protocol === 'https:' ? `wss://${host}/api/ws` : `ws://${host}:8000/api/ws`;
  const wsUrl = import.meta.env.VITE_WS_URL || defaultWs;
  const ws = new WebSocket(wsUrl);
  
  ws.onmessage = (event) => {
    try {
      const data = JSON.parse(event.data);
      onMessage(data);
    } catch (e) {
      console.error("Failed to parse websocket message", e);
    }
  };

  return ws;
};
