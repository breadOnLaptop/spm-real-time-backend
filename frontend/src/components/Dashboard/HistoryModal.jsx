import React from 'react';

const HistoryModal = ({ isOpen, onClose, agentId, history }) => {
  if (!isOpen) return null;

  return (
    <div className="modal-overlay">
      <div className="modal-content">
        <div className="modal-header">
          <h2>History for {agentId}</h2>
          <button onClick={onClose} className="btn-close">&times;</button>
        </div>
        <div className="modal-body">
          {history.length === 0 ? (
            <p>No history found.</p>
          ) : (
            <table className="history-table">
              <thead>
                <tr>
                  <th>Timestamp</th>
                  <th>CPU (%)</th>
                  <th>RAM (%)</th>
                </tr>
              </thead>
              <tbody>
                {history.map((record, idx) => (
                  <tr key={idx}>
                    <td>{new Date(record.timestamp).toLocaleTimeString()}</td>
                    <td>{record.cpu_utilization.toFixed(1)}</td>
                    <td>{record.memory_utilization.toFixed(1)}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}
        </div>
      </div>
    </div>
  );
};

export default HistoryModal;
