export interface ApiResponse<T = unknown> {
  is_successful: boolean;
  request_uuid: string;
  status_code: number;
  message: string;
  data: T;
  errors: unknown;
  error_code: number;
}

export class ApiError extends Error {
  statusCode: number;
  errorCode: number;
  errors: unknown;
  requestUuid?: string;

  constructor(
    message: string,
    statusCode: number,
    errorCode: number,
    errors: unknown = null,
    requestUuid?: string
  ) {
    super(message);
    this.name = 'ApiError';
    this.statusCode = statusCode;
    this.errorCode = errorCode;
    this.errors = errors;
    this.requestUuid = requestUuid;
  }
}

const BASE_URL = import.meta.env.VITE_API_URL || '/api/v1';

export async function apiFetch<T = unknown>(
  endpoint: string,
  options?: RequestInit
): Promise<T> {
  const normalizedEndpoint = endpoint.startsWith('/') ? endpoint : `/${endpoint}`;
  const url = `${BASE_URL}${normalizedEndpoint}`;

  const response = await fetch(url, {
    headers: {
      'Content-Type': 'application/json',
      ...options?.headers,
    },
    ...options,
  });

  let envelope: ApiResponse<T>;
  try {
    envelope = await response.json();
  } catch (err) {
    throw new ApiError(
      response.statusText || 'Failed to parse JSON response',
      response.status,
      response.status * 100,
      err
    );
  }

  if (!response.ok || !envelope.is_successful) {
    throw new ApiError(
      envelope.message || 'API request failed',
      envelope.status_code || response.status,
      envelope.error_code || 0,
      envelope.errors,
      envelope.request_uuid
    );
  }

  return envelope.data;
}
