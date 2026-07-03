import React, { useState, useEffect } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { apiCall, SHORT_URL_BASE } from '../utils/api';
import { Link as LinkIcon, Plus, BarChart2, Copy, CheckCircle2 } from 'lucide-react';

export default function Dashboard() {
  const [links, setLinks] = useState([]);
  const [originalUrl, setOriginalUrl] = useState('');
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [copiedSlug, setCopiedSlug] = useState(null);
  const navigate = useNavigate();

  useEffect(() => {
    fetchLinks();
  }, []);

  const fetchLinks = async () => {
    try {
      const data = await apiCall('/links');
      setLinks(data || []);
      setLoading(false);
    } catch (err) {
      if (err.message.includes('401')) {
        localStorage.removeItem('token');
        navigate('/login');
      }
      setError(err.message);
      setLoading(false);
    }
  };

  const handleShorten = async (e) => {
    e.preventDefault();
    setError('');
    try {
      await apiCall('/shorten', {
        method: 'POST',
        body: JSON.stringify({ original_url: originalUrl }),
      });
      setOriginalUrl('');
      fetchLinks(); // Refresh list
    } catch (err) {
      setError(err.message);
    }
  };

  const copyToClipboard = (slug) => {
    navigator.clipboard.writeText(`${SHORT_URL_BASE}${slug}`);
    setCopiedSlug(slug);
    setTimeout(() => setCopiedSlug(null), 2000);
  };

  return (
    <div className="container page-container">
      <div className="glass-panel card" style={{ marginBottom: '2rem' }}>
        <h2 className="gradient-text" style={{ display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
          <LinkIcon size={28} /> Create New Short Link
        </h2>
        {error && <div className="alert-error">{error}</div>}
        
        <form onSubmit={handleShorten} style={{ display: 'flex', gap: '1rem', marginTop: '1.5rem' }}>
          <div className="form-group" style={{ margin: 0, flex: 1 }}>
            <input 
              type="url" 
              className="form-input" 
              placeholder="Paste your long URL here (e.g., https://example.com/very/long/path)"
              value={originalUrl}
              onChange={(e) => setOriginalUrl(e.target.value)}
              required 
            />
          </div>
          <button type="submit" className="btn btn-primary" style={{ whiteSpace: 'nowrap' }}>
            <Plus size={20} /> Shorten URL
          </button>
        </form>
      </div>

      <h3 style={{ marginTop: '2rem' }}>Your Links</h3>
      
      {loading ? (
        <p>Loading your links...</p>
      ) : links.length === 0 ? (
        <p>You haven't created any links yet.</p>
      ) : (
        <div className="link-grid">
          {links.map((link) => (
            <div key={link.id} className="glass-panel link-card">
              <div className="link-card-header">
                <h4 className="link-title" style={{ color: 'var(--primary-color)' }}>
                  {SHORT_URL_BASE}{link.slug}
                </h4>
                <button 
                  onClick={() => copyToClipboard(link.slug)}
                  className="btn btn-outline"
                  style={{ padding: '0.4rem', borderRadius: '6px', border: 'none' }}
                  title="Copy to clipboard"
                >
                  {copiedSlug === link.slug ? <CheckCircle2 size={18} color="var(--success-color)"/> : <Copy size={18} />}
                </button>
              </div>
              
              <div className="link-original" title={link.original_url}>
                {link.original_url}
              </div>
              
              <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginTop: 'auto', paddingTop: '1rem', borderTop: '1px solid var(--surface-border)' }}>
                <div className="link-stats">
                  <BarChart2 size={16} /> {link.total_clicks} clicks
                </div>
                <Link to={`/analytics/${link.slug}`} className="btn btn-outline" style={{ padding: '0.4rem 0.8rem', fontSize: '0.85rem' }}>
                  View Analytics
                </Link>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
