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
    <div className="status-container">
      <section className="status-card">
        <div className="card-header">
          <h2 className="card-title">پایش سرویس‌های سامانه</h2>
          <span className="card-badge">سیستم برخط</span>
        </div>

        <div className="status-list">
          <div className="status-item">
            <div className="item-meta">
              <span className="item-title">سرویس مرکزی API</span>
              <span className="item-desc">پاسخگویی به درخواست‌های وب و کلاینت</span>
            </div>
            <span
              className={`status-pill ${
                loading ? 'loading' : error ? 'error' : 'ok'
              }`}
            >
              <span className="status-dot" />
              {loading ? 'در حال بررسی...' : status}
            </span>
          </div>

          <div className="status-item">
            <div className="item-meta">
              <span className="item-title">پایگاه‌داده تیک دیتا (PostgreSQL)</span>
              <span className="item-desc">ذخیره‌سازی اطلاعات دارایی‌های پایه و مظنه‌ها</span>
            </div>
            <span
              className={`status-pill ${
                loading ? 'loading' : error ? 'error' : 'ok'
              }`}
            >
              <span className="status-dot" />
              {loading ? 'در حال بررسی...' : error ? 'عدم برقراری ارتباط' : 'متصل و پایدار'}
            </span>
          </div>
        </div>

        {error && (
          <div className="error-note">
            <span className="error-icon">!</span>
            <span>جزئیات خطا: {error}</span>
          </div>
        )}
      </section>
    </div>
  );
}

export default HomePage;
