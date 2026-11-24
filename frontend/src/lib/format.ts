const persianNumberFormatter = new Intl.NumberFormat('fa-IR');

/**
 * Format a number into Persian digits using Intl.NumberFormat('fa-IR').
 */
export function toPersianDigits(value: number | string): string {
  const num = typeof value === 'string' ? Number(value) : value;
  if (isNaN(num)) {
    return String(value);
  }
  return persianNumberFormatter.format(num);
}

/**
 * Persian labels for underlying asset types.
 */
export const ASSET_TYPE_PERSIAN_LABELS: Record<string, string> = {
  stock: 'سهم',
  etf: 'صندوق ETF',
  fund: 'صندوق',
};

/**
 * Returns Persian label for a given asset_type.
 */
export function formatAssetType(assetType: string): string {
  const normalized = (assetType || '').toLowerCase();
  return ASSET_TYPE_PERSIAN_LABELS[normalized] || assetType;
}
