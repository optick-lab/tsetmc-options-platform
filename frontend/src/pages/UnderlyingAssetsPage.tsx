import { useEffect, useState } from 'react';
import { useSearchParams } from 'react-router-dom';
import { formatAssetType, toPersianDigits } from '../lib/format';
import { getUnderlyings, PaginatedUnderlyings } from '../lib/underlyings';
import './UnderlyingAssetsPage.css';

export function UnderlyingAssetsPage() {
  const [searchParams, setSearchParams] = useSearchParams();
  const pageParam = parseInt(searchParams.get('page') || '1', 10);
  const initialPage = isNaN(pageParam) || pageParam < 1 ? 1 : pageParam;

  const [data, setData] = useState<PaginatedUnderlyings | null>(null);
  const [page, setPage] = useState<number>(initialPage);
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

  const handlePageChange = (newPage: number) => {
    setPage(newPage);
    setSearchParams({ page: String(newPage) });
  };

  const totalPages = data ? Math.max(1, Math.ceil(data.total / data.page_size)) : 1;

  return (
    <div className="table-card">
      {error && !loading && (
        <div className="state-container">
          <div className="error-card">
            <span className="error-title">خطا در بارگذاری اطلاعات</span>
            <p className="error-message">{error}</p>
            <button
              type="button"
              className="btn-primary"
              onClick={() => fetchData(page)}
            >
              تلاش مجدد
            </button>
          </div>
        </div>
      )}

      {loading && (
        <div className="table-wrapper">
          <table className="data-table">
            <thead>
              <tr>
                <th>نماد</th>
                <th>نوع دارایی</th>
                <th>شناسه ابزار</th>
              </tr>
            </thead>
            <tbody>
              {Array.from({ length: 8 }).map((_, index) => (
                <tr key={index} className="skeleton-row">
                  <td>
                    <div className="skeleton-line skeleton-w-sm" />
                  </td>
                  <td>
                    <div className="skeleton-pill" />
                  </td>
                  <td>
                    <div className="skeleton-line skeleton-w-md" />
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      {!loading && !error && data && data.items.length === 0 && (
        <div className="empty-state">
          <p>هیچ نمادی یافت نشد.</p>
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
                  <th>شناسه ابزار</th>
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
              صفحه {toPersianDigits(data.page)} از {toPersianDigits(totalPages)} (مجموع {toPersianDigits(data.total)} نماد)
            </span>

            <div className="pagination-actions">
              <button
                type="button"
                className="btn-secondary"
                disabled={data.page <= 1 || loading}
                onClick={() => handlePageChange(Math.max(1, page - 1))}
              >
                قبلی
              </button>
              <button
                type="button"
                className="btn-secondary"
                disabled={data.page >= totalPages || loading}
                onClick={() => handlePageChange(page + 1)}
              >
                بعدی
              </button>
            </div>
          </div>
        </>
      )}
    </div>
  );
}

export default UnderlyingAssetsPage;
