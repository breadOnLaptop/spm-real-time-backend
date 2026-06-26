export const createWebSocketConnection = (onMessage) => {
  const host = window.location.hostname || 'localhost';
  const ws = new WebSocket(`ws://${host}:8000/api/ws`); // Update to match API router prefix if any, wait, it's /ws, but router is /api? 
  // Let me check my api router: api_router.include_router(websockets.router, tags=["websockets"]) 
  // without prefix, so it is just /ws.
  
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
