import { useEffect, useState } from 'react';
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
        setError(err.message || 'Failed to load underlying assets');
        setLoading(false);
      });
  };

  useEffect(() => {
    fetchData(page);
  }, [page]);

  const totalPages = data ? Math.max(1, Math.ceil(data.total / data.page_size)) : 1;

  return (
    <main className="underlyings-container">
      <header className="page-header">
        <h1 className="page-title">Underlying Assets</h1>
        <p className="page-subtitle">
          Spot market instruments (stocks, ETFs, funds) eligible for options contracts.
        </p>
      </header>

      <div className="table-card">
        {loading && (
          <div className="state-container">
            <p className="loading-text">Loading underlying assets...</p>
          </div>
        )}

        {error && !loading && (
          <div className="state-container">
            <div className="error-box">
              <p>Error: {error}</p>
              <button
                type="button"
                className="retry-btn"
                onClick={() => fetchData(page)}
              >
                Retry
              </button>
            </div>
          </div>
        )}

        {!loading && !error && data && (
          <>
            <div className="table-wrapper">
              <table className="data-table">
                <thead>
                  <tr>
                    <th>Symbol</th>
                    <th>Asset Type</th>
                    <th>Instrument Code</th>
                  </tr>
                </thead>
                <tbody>
                  {data.items.map((item) => (
                    <tr key={item.id}>
                      <td className="symbol-cell" dir="rtl">{item.symbol}</td>
                      <td>
                        <span className={`type-badge ${item.asset_type.toLowerCase()}`}>
                          {item.asset_type}
                        </span>
                      </td>
                      <td className="code-cell">
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
                Page {data.page} of {totalPages} (Total: {data.total} assets)
              </span>

              <div className="pagination-actions">
                <button
                  type="button"
                  className="pagination-btn"
                  disabled={data.page <= 1 || loading}
                  onClick={() => setPage((prev) => Math.max(1, prev - 1))}
                >
                  Previous
                </button>
                <button
                  type="button"
                  className="pagination-btn"
                  disabled={data.page >= totalPages || loading}
                  onClick={() => setPage((prev) => prev + 1)}
                >
                  Next
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
