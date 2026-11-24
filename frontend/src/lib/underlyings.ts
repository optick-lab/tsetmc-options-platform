import { apiFetch } from './api';

export interface UnderlyingItem {
  id: number;
  symbol: string;
  ins_code: string | null;
  asset_type: 'stock' | 'etf' | 'fund' | string;
}

export interface PaginatedUnderlyings {
  items: UnderlyingItem[];
  page: number;
  page_size: number;
  total: number;
}

export async function getUnderlyings(
  page: number = 1,
  pageSize: number = 10
): Promise<PaginatedUnderlyings> {
  return apiFetch<PaginatedUnderlyings>(
    `/underlyings?page=${page}&page_size=${pageSize}`
  );
}
