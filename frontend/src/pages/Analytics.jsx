import React, { useState, useEffect } from 'react';
import { useParams, Link } from 'react-router-dom';
import { apiCall, SHORT_URL_BASE } from '../utils/api';
import { ArrowLeft, Sparkles, TrendingUp, Globe, Smartphone } from 'lucide-react';
import { 
  LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer,
  BarChart, Bar, PieChart, Pie, Cell
} from 'recharts';

const COLORS = ['#66fcf1', '#45a29e', '#c5c6c7', '#ffffff'];

export default function Analytics() {
  const { slug } = useParams();
  const [data, setData] = useState(null);
  const [summary, setSummary] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    const fetchData = async () => {
      try {
        const analyticsData = await apiCall(`/links/${slug}/analytics`);
        setData(analyticsData);
        
        try {
          const summaryData = await apiCall(`/links/${slug}/summary`);
          setSummary(summaryData);
        } catch (sumErr) {
          console.log("Summary not available yet", sumErr);
        }
        
        setLoading(false);
      } catch (err) {
        setError(err.message);
        setLoading(false);
      }
    };
    fetchData();
  }, [slug]);

  if (loading) return <div className="container page-container"><p>Loading analytics...</p></div>;
  if (error) return <div className="container page-container"><div className="alert-error">{error}</div></div>;

  return (
    <div className="container page-container">
      <div style={{ marginBottom: '2rem' }}>
        <Link to="/dashboard" className="nav-link" style={{ display: 'inline-flex', alignItems: 'center', gap: '0.5rem', marginBottom: '1rem' }}>
          <ArrowLeft size={16} /> Back to Dashboard
        </Link>
        <h2 className="gradient-text">Analytics for {SHORT_URL_BASE}{slug}</h2>
      </div>

      {summary && (
        <div className="glass-panel ai-summary-card">
          <h3 style={{ display: 'flex', alignItems: 'center', gap: '0.5rem', color: 'var(--primary-color)' }}>
            <Sparkles size={20} /> AI Insight
          </h3>
          <p style={{ color: '#fff', fontSize: '1.1rem', margin: 0 }}>
            {summary.summary}
          </p>
          <div style={{ fontSize: '0.8rem', color: 'var(--text-secondary)', marginTop: '1rem' }}>
            Generated {new Date(summary.generated_at).toLocaleString()}
          </div>
        </div>
      )}

      <div className="link-grid" style={{ marginTop: '0' }}>
        
        {/* Clicks by Date */}
        <div className="glass-panel card" style={{ gridColumn: '1 / -1' }}>
          <h3 style={{ display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
            <TrendingUp size={20} /> Traffic Over Time
          </h3>
          <div style={{ height: '300px', marginTop: '1.5rem' }}>
            {data.clicks_by_date && data.clicks_by_date.length > 0 ? (
              <ResponsiveContainer width="100%" height="100%">
                <LineChart data={data.clicks_by_date}>
                  <CartesianGrid strokeDasharray="3 3" stroke="rgba(255,255,255,0.1)" />
                  <XAxis dataKey="date" tickFormatter={(tick) => new Date(tick).toLocaleDateString()} stroke="#c5c6c7" />
                  <YAxis stroke="#c5c6c7" />
                  <Tooltip 
                    contentStyle={{ backgroundColor: 'rgba(11,12,16,0.9)', border: '1px solid rgba(255,255,255,0.1)' }}
                    labelFormatter={(label) => new Date(label).toLocaleDateString()}
                  />
                  <Line type="monotone" dataKey="count" stroke="#66fcf1" strokeWidth={3} dot={{ fill: '#45a29e', r: 4 }} activeDot={{ r: 6 }} />
                </LineChart>
              </ResponsiveContainer>
            ) : <p>No click data available yet.</p>}
          </div>
        </div>

        {/* Clicks by Country */}
        <div className="glass-panel card">
          <h3 style={{ display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
            <Globe size={20} /> Top Countries
          </h3>
          <div style={{ height: '250px', marginTop: '1.5rem' }}>
            {data.clicks_by_country && data.clicks_by_country.length > 0 ? (
              <ResponsiveContainer width="100%" height="100%">
                <BarChart data={data.clicks_by_country} layout="vertical" margin={{ top: 0, right: 20, left: 60, bottom: 0 }}>
                  <CartesianGrid strokeDasharray="3 3" stroke="rgba(255,255,255,0.1)" horizontal={true} vertical={false} />
                  <XAxis type="number" stroke="#c5c6c7" />
                  <YAxis dataKey="country" type="category" stroke="#c5c6c7" width={80} />
                  <Tooltip contentStyle={{ backgroundColor: 'rgba(11,12,16,0.9)', border: '1px solid rgba(255,255,255,0.1)' }} />
                  <Bar dataKey="count" fill="#45a29e" radius={[0, 4, 4, 0]} />
                </BarChart>
              </ResponsiveContainer>
            ) : <p>No country data available.</p>}
          </div>
        </div>

        {/* Clicks by Device */}
        <div className="glass-panel card">
          <h3 style={{ display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
            <Smartphone size={20} /> Device Types
          </h3>
          <div style={{ height: '250px', marginTop: '1.5rem' }}>
            {data.clicks_by_device && data.clicks_by_device.length > 0 ? (
              <ResponsiveContainer width="100%" height="100%">
                <PieChart>
                  <Pie
                    data={data.clicks_by_device}
                    cx="50%"
                    cy="50%"
                    innerRadius={60}
                    outerRadius={80}
                    paddingAngle={5}
                    dataKey="count"
                    nameKey="device_type"
                    label={({ name, percent }) => `${name} ${(percent * 100).toFixed(0)}%`}
                  >
                    {data.clicks_by_device.map((entry, index) => (
                      <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                    ))}
                  </Pie>
                  <Tooltip contentStyle={{ backgroundColor: 'rgba(11,12,16,0.9)', border: '1px solid rgba(255,255,255,0.1)' }} />
                </PieChart>
              </ResponsiveContainer>
            ) : <p>No device data available.</p>}
          </div>
        </div>

      </div>
    </div>
  );
}
