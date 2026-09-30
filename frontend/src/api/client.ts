import type { Operation } from '../lib/operations'

export interface CalculateRequest {
  operation: Operation
  a: number
  b?: number
}

export interface CalculateResponse {
  operation: Operation
  operands: number[]
  result: number
}

interface ErrorResponse {
  error: { code: string; message: string }
}

/** Error raised for any failed API call, carrying the backend's error code. */
export class ApiError extends Error {
  readonly code: string
  readonly status: number

  constructor(code: string, message: string, status: number) {
    super(message)
    this.name = 'ApiError'
    this.code = code
    this.status = status
  }
}

// Relative by default: same origin in Docker, proxied by Vite in development.
const BASE_URL = import.meta.env.VITE_API_URL ?? ''

export async function calculate(
  req: CalculateRequest,
  fetchFn: typeof fetch = fetch,
): Promise<CalculateResponse> {
  let res: Response
  try {
    res = await fetchFn(`${BASE_URL}/api/v1/calculate`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(req),
    })
  } catch {
    throw new ApiError('NETWORK_ERROR', 'Could not reach the server', 0)
  }

  const body: unknown = await res.json().catch(() => null)

  if (!res.ok) {
    if (isErrorResponse(body)) {
      throw new ApiError(body.error.code, body.error.message, res.status)
    }
    throw new ApiError('UNKNOWN_ERROR', `Unexpected error (HTTP ${res.status})`, res.status)
  }

  return body as CalculateResponse
}

function isErrorResponse(body: unknown): body is ErrorResponse {
  return (
    typeof body === 'object' &&
    body !== null &&
    'error' in body &&
    typeof (body as ErrorResponse).error?.code === 'string'
  )
}
