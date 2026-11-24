import { useState, useEffect } from 'react';
import {
  getRawResponses,
  getRawResponseTrades,
  formatTradeDateJalali,
  RawResponseItem,
  TradesData,
} from '../lib/rawResponses';
import './RawTradesPage.css';

const faNum = new Intl.NumberFormat('fa-IR');
const faCurrency = new Intl.NumberFormat('fa-IR', {
  maximumFractionDigits: 0,
});

function formatJalaliDateTime(isoString: string): string {
  try {
    const d = new Date(isoString);
    return new Intl.DateTimeFormat('fa-IR', {
      calendar: 'persian',
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit',
      hour12: false,
    }).format(d);
  } catch {
    return isoString;
  }
}

export function RawTradesPage() {
  const [responses, setResponses] = useState<RawResponseItem[]>([]);
  const [selectedId, setSelectedId] = useState<number | null>(null);
  const [tradesData, setTradesData] = useState<TradesData | null>(null);

  const [isLoadingResponses, setIsLoadingResponses] = useState(true);
  const [isLoadingTrades, setIsLoadingTrades] = useState(false);
  const [errorResponses, setErrorResponses] = useState<string | null>(null);
  const [errorTrades, setErrorTrades] = useState<string | null>(null);

  // Initial load: fetch responses list and immediately load trades for the first response
  useEffect(() => {
    let mounted = true;
    setIsLoadingResponses(true);
    setErrorResponses(null);

    getRawResponses()
      .then((items) => {
        if (!mounted) return;
        setResponses(items);
        if (items.length > 0) {
          const firstId = items[0].id;
          setSelectedId(firstId);
          setIsLoadingTrades(true);
          return getRawResponseTrades(firstId).then((data) => {
            if (!mounted) return;
            setTradesData(data);
          });
        }
      })
      .catch((err) => {
        if (!mounted) return;
        setErrorResponses(err instanceof Error ? err.message : 'خطا در بارگذاری پاسخ‌های خام');
      })
      .finally(() => {
        if (mounted) {
          setIsLoadingResponses(false);
          setIsLoadingTrades(false);
        }
      });

    return () => {
      mounted = false;
    };
  }, []);

  const handleSelectResponse = (id: number) => {
    if (selectedId === id) return;
    setSelectedId(id);
    setIsLoadingTrades(true);
    setErrorTrades(null);

    getRawResponseTrades(id)
      .then((data) => {
        setTradesData(data);
      })
      .catch((err) => {
        setErrorTrades(err instanceof Error ? err.message : 'خطا در بارگذاری معاملات');
      })
      .finally(() => {
        setIsLoadingTrades(false);
      });
  };

  return (
    <div className="raw-trades-page">
      {/* Small note banner */}
      <div className="source-note-banner">
        <span className="source-note-dot" />
        <span>داده از پاسخ خام ذخیره‌شده‌ی TSETMC خوانده می‌شود.</span>
      </div>

      {/* Responses Selection */}
      <section>
        <h2 className="section-title">پاسخ‌های خام ذخیره‌شده (TSETMC)</h2>

        {isLoadingResponses && (
          <div className="loading-indicator">
            <span>در حال بارگذاری لیست پاسخ‌های خام...</span>
          </div>
        )}

        {errorResponses && (
          <div className="error-indicator">
            <span>{errorResponses}</span>
          </div>
        )}

        {!isLoadingResponses && !errorResponses && responses.length === 0 && (
          <div className="empty-indicator">
            <span>هیچ پاسخ خامی در پایگاه داده ذخیره نشده است.</span>
          </div>
        )}

        {!isLoadingResponses && responses.length > 0 && (
          <div className="responses-scroll-container">
            <div className="responses-grid">
              {responses.map((item) => (
                <div
                  key={item.id}
                  className={`response-card ${selectedId === item.id ? 'selected' : ''}`}
                  onClick={() => handleSelectResponse(item.id)}
                >
                  <div className="response-card-header">
                    <span className="response-ins-code">
                      {item.ins_code ? `کد: ${item.ins_code}` : 'بدون کد ابزار'}
                    </span>
                    <span
                      className={`response-badge ${
                        item.is_blocked ? 'blocked' : item.http_status === 200 ? 'success' : 'blocked'
                      }`}
                    >
                      {item.is_blocked
                        ? 'مسدود شده'
                        : item.http_status === 200
                        ? '۲۰۰ موفق'
                        : `خطا ${item.http_status}`}
                    </span>
                  </div>

                  <div className="response-trade-date-line">
                    <span className="trade-date-label">تاریخ معامله:</span>
                    <span className="trade-date-value">
                      {formatTradeDateJalali(item.trade_date)}
                    </span>
                    {item.trade_date && (
                      <span className="trade-date-raw">({item.trade_date})</span>
                    )}
                  </div>

                  <div className="response-card-details">
                    <span>دریافت: {formatJalaliDateTime(item.fetched_at)}</span>
                    <span>{faNum.format(item.body_size)} بایت</span>
                  </div>
                </div>
              ))}
            </div>
          </div>
        )}
      </section>

      {/* Trades Details Section */}
      {selectedId !== null && (
        <section className="trades-section">
          {isLoadingTrades && (
            <div className="loading-indicator">
              <span>در حال استخراج و پردازش معاملات از بدنه خام...</span>
            </div>
          )}

          {errorTrades && (
            <div className="error-indicator">
              <span>{errorTrades}</span>
            </div>
          )}

          {!isLoadingTrades && !errorTrades && tradesData && (
            <>
              {/* Summary Cards */}
              <div className="summary-grid">
                <div className="summary-card">
                  <span className="summary-label">تعداد معاملات</span>
                  <span className="summary-value highlight">
                    {faNum.format(tradesData.summary.trade_count)}
                  </span>
                </div>

                <div className="summary-card">
                  <span className="summary-label">حجم کل معاملات</span>
                  <span className="summary-value">
                    {faNum.format(tradesData.summary.total_volume)}
                  </span>
                </div>

                <div className="summary-card">
                  <span className="summary-label">ارزش کل (ریال)</span>
                  <span className="summary-value">
                    {faCurrency.format(tradesData.summary.total_value)}
                  </span>
                </div>

                <div className="summary-card">
                  <span className="summary-label">کمترین / بیشترین قیمت</span>
                  <span className="summary-value">
                    {faNum.format(tradesData.summary.min_price)} /{' '}
                    {faNum.format(tradesData.summary.max_price)}
                  </span>
                </div>

                <div className="summary-card">
                  <span className="summary-label">اولین / آخرین قیمت</span>
                  <span className="summary-value">
                    {faNum.format(tradesData.summary.first_price)} /{' '}
                    {faNum.format(tradesData.summary.last_price)}
                  </span>
                </div>
              </div>

              {/* Trades Table */}
              <div className="trades-header">
                <h3 className="section-title" style={{ margin: 0 }}>
                  ریز معاملات ثبت‌شده
                </h3>
                <span className="trades-count-pill">
                  {faNum.format(tradesData.trades.length)} ردیف
                </span>
              </div>

              <div className="trades-card">
                <table className="trades-table">
                  <thead>
                    <tr>
                      <th>شماره (nTran)</th>
                      <th>زمان معامله</th>
                      <th>حجم</th>
                      <th>قیمت (ریال)</th>
                      <th>وضعیت</th>
                    </tr>
                  </thead>
                  <tbody>
                    {tradesData.trades.map((trade) => (
                      <tr
                        key={trade.n_tran}
                        className={trade.canceled ? 'trade-canceled-row' : ''}
                      >
                        <td>{faNum.format(trade.n_tran)}</td>
                        <td>{trade.time}</td>
                        <td>{faNum.format(trade.volume)}</td>
                        <td>{faCurrency.format(trade.price)}</td>
                        <td>
                          <span
                            className={`status-badge ${
                              trade.canceled ? 'canceled' : 'valid'
                            }`}
                          >
                            {trade.canceled ? 'باطل‌شده' : 'معتبر'}
                          </span>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            </>
          )}
        </section>
      )}
    </div>
  );
}
