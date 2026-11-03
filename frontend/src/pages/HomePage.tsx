import { useEffect, useState } from 'react';
import { getHealth } from '../lib/health';
import './HomePage.css';

export function HomePage() {
  const [status, setStatus] = useState<string>('checking...');
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    getHealth()
      .then((data) => {
        setStatus(data.status || 'ok');
        setLoading(false);
      })
      .catch((err: Error) => {
        setStatus('unreachable');
        setError(err.message);
        setLoading(false);
      });
  }, []);

  return (
    <main className="home-container">
      <header className="home-header">
        <h1 className="home-title">TSETMC Options Platform</h1>
        <p className="home-subtitle">Tick Data & Analytics Platform</p>
      </header>

      <section className="card">
        <h2 className="card-title">System Status</h2>
        <div className="status-row">
          <span className="status-label">Backend API Status</span>
          <span
            className={`status-badge ${
              loading ? 'loading' : error ? 'error' : 'ok'
            }`}
          >
            <span className="status-dot" />
            {loading ? 'checking...' : status}
          </span>
        </div>
      </section>
    </main>
  );
}

export default HomePage;
