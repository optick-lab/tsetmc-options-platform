import { useEffect, useState } from 'react';
import { getHealth } from '../lib/health';
import './HomePage.css';

export function HomePage() {
  const [status, setStatus] = useState<string>('در حال بررسی...');
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    getHealth()
      .then((data) => {
        setStatus(data.status === 'ok' ? 'فعال و آماده‌به‌کار' : data.status || 'فعال');
        setLoading(false);
      })
      .catch((err: Error) => {
        setStatus('عدم برقراری ارتباط');
        setError(err.message);
        setLoading(false);
      });
  }, []);

  return (
    <main className="home-container" dir="rtl">
      <header className="home-header">
        <h1 className="home-title">سامانه اختیار معامله TSETMC</h1>
        <p className="home-subtitle">سکوی جمع‌آوری و پردازش داده‌های تیک بازار ابزارهای مشتقه</p>
      </header>

      <section className="card">
        <h2 className="card-title">وضعیت سیستم</h2>
        <div className="status-row">
          <span className="status-label">وضعیت سرویس API بک‌اند</span>
          <span
            className={`status-badge ${
              loading ? 'loading' : error ? 'error' : 'ok'
            }`}
          >
            <span className="status-dot" />
            {loading ? 'در حال بررسی...' : status}
          </span>
        </div>
      </section>
    </main>
  );
}

export default HomePage;
