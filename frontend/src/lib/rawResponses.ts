import { apiFetch } from './api';

export interface RawResponseItem {
  id: number;
  endpoint: string;
  ins_code: string | null;
  http_status: number;
  is_blocked: boolean;
  fetched_at: string;
  body_size: number;
  trade_date: string | null;
}

export interface TradeItem {
  n_tran: number;
  time: string;
  volume: number;
  price: number;
  canceled: boolean;
}

export interface TradeSummary {
  trade_count: number;
  total_volume: number;
  total_value: number;
  min_price: number;
  max_price: number;
  first_price: number;
  last_price: number;
}

export interface TradesData {
  summary: TradeSummary;
  trades: TradeItem[];
}

export async function getRawResponses(): Promise<RawResponseItem[]> {
  return apiFetch<RawResponseItem[]>('/raw-responses');
}

export async function getRawResponseTrades(id: number): Promise<TradesData> {
  return apiFetch<TradesData>(`/raw-responses/${id}/trades`);
}

export function formatTradeDateJalali(tradeDateStr?: string | null): string {
  if (!tradeDateStr || tradeDateStr.length !== 8) return '—';
  const year = parseInt(tradeDateStr.slice(0, 4), 10);
  const month = parseInt(tradeDateStr.slice(4, 6), 10) - 1;
  const day = parseInt(tradeDateStr.slice(6, 8), 10);
  const date = new Date(year, month, day);
  return new Intl.DateTimeFormat('fa-IR', {
    calendar: 'persian',
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
  }).format(date);
}
