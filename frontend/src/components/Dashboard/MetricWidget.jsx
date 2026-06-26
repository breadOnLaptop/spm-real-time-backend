import React from 'react';

const MetricWidget = ({ label, value }) => {
  return (
    <div className="metric">
      <span className="label">{label}</span>
      <span className="value">{value}</span>
    </div>
  );
};

export default MetricWidget;
