import { useEffect, useState } from 'react';
import { formatAssetType, toPersianDigits } from '../lib/format';
import { getUnderlyings, PaginatedUnderlyings } from '../lib/underlyings';
import './UnderlyingAssetsPage.css';

export function UnderlyingAssetsPage() {
  const [data, setData] = useState<PaginatedUnderlyings | null>(null);
  const [page, setPage] = useState<number>(1);
  const [pageSize] = useState<number>(10);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  const fetchData = (targetPage: number) => {
    setLoading(true);
    setError(null);
    getUnderlyings(targetPage, pageSize)
      .then((res) => {
        setData(res);
        setPage(res.page);
        setLoading(false);
      })
      .catch((err: Error) => {
        setError(err.message || 'خطا در برقراری ارتباط با سرور');
        setLoading(false);
      });
  };

  useEffect(() => {
    fetchData(page);
  }, [page]);

  const totalPages = data ? Math.max(1, Math.ceil(data.total / data.page_size)) : 1;

  return (
    <main className="underlyings-container" dir="rtl">
      <header className="page-header">
        <h1 className="page-title">دارایی‌های پایه</h1>
        <p className="page-subtitle">
          ابزارهای بازار نقدی (سهام و صندوق‌ها) مشمول قراردادهای اختیار معامله
        </p>
      </header>

      <div className="table-card">
        {loading && (
          <div className="state-container">
            <p className="loading-text">در حال بارگذاری دارایی‌های پایه...</p>
          </div>
        )}

        {error && !loading && (
          <div className="state-container">
            <div className="error-box">
              <p>خطا در دریافت اطلاعات: {error}</p>
              <button
                type="button"
                className="retry-btn"
                onClick={() => fetchData(page)}
              >
                تلاش مجدد
              </button>
            </div>
          </div>
        )}

        {!loading && !error && data && data.items.length === 0 && (
          <div className="state-container">
            <p className="empty-text">هیچ دارایی پایه‌ای یافت نشد.</p>
          </div>
        )}

        {!loading && !error && data && data.items.length > 0 && (
          <>
            <div className="table-wrapper">
              <table className="data-table">
                <thead>
                  <tr>
                    <th>نماد</th>
                    <th>نوع دارایی</th>
                    <th>کد نماد</th>
                  </tr>
                </thead>
                <tbody>
                  {data.items.map((item) => (
                    <tr key={item.id}>
                      <td className="symbol-cell">{item.symbol}</td>
                      <td>
                        <span className={`type-badge ${item.asset_type.toLowerCase()}`}>
                          {formatAssetType(item.asset_type)}
                        </span>
                      </td>
                      <td className="code-cell" dir="ltr">
                        {item.ins_code ? (
                          item.ins_code
                        ) : (
                          <span className="code-null">—</span>
                        )}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>

            <div className="pagination-bar">
              <span className="pagination-info">
                صفحه {toPersianDigits(data.page)} از {toPersianDigits(totalPages)} (مجموع: {toPersianDigits(data.total)} دارایی)
              </span>

              <div className="pagination-actions">
                <button
                  type="button"
                  className="pagination-btn"
                  disabled={data.page <= 1 || loading}
                  onClick={() => setPage((prev) => Math.max(1, prev - 1))}
                >
                  قبلی
                </button>
                <button
                  type="button"
                  className="pagination-btn"
                  disabled={data.page >= totalPages || loading}
                  onClick={() => setPage((prev) => prev + 1)}
                >
                  بعدی
                </button>
              </div>
            </div>
          </>
        )}
      </div>
    </main>
  );
}

export default UnderlyingAssetsPage;
