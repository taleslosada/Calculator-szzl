import { describe, expect, it, vi } from 'vitest'
import { ApiError, calculate } from './client'

function mockFetch(status: number, body: unknown) {
  return vi.fn().mockResolvedValue(
    new Response(typeof body === 'string' ? body : JSON.stringify(body), {
      status,
      headers: { 'Content-Type': 'application/json' },
    }),
  )
}

describe('calculate', () => {
  it('posts the request and returns the parsed response', async () => {
    const fetchFn = mockFetch(200, { operation: 'add', operands: [2, 3], result: 5 })

    const res = await calculate({ operation: 'add', a: 2, b: 3 }, fetchFn)

    expect(res.result).toBe(5)
    expect(fetchFn).toHaveBeenCalledWith('/api/v1/calculate', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ operation: 'add', a: 2, b: 3 }),
    })
  })

  it('throws ApiError with the backend code on error responses', async () => {
    const fetchFn = mockFetch(422, {
      error: { code: 'DIVISION_BY_ZERO', message: 'division by zero' },
    })

    const err = await calculate({ operation: 'divide', a: 1, b: 0 }, fetchFn).catch((e) => e)

    expect(err).toBeInstanceOf(ApiError)
    expect(err).toMatchObject({ code: 'DIVISION_BY_ZERO', status: 422, message: 'division by zero' })
  })

  it('throws UNKNOWN_ERROR when the error body is not the expected shape', async () => {
    const fetchFn = mockFetch(502, '<html>Bad Gateway</html>')

    await expect(calculate({ operation: 'add', a: 1, b: 1 }, fetchFn)).rejects.toMatchObject({
      code: 'UNKNOWN_ERROR',
      status: 502,
    })
  })

  it('throws NETWORK_ERROR when the request cannot be sent', async () => {
    const fetchFn = vi.fn().mockRejectedValue(new TypeError('Failed to fetch'))

    await expect(calculate({ operation: 'add', a: 1, b: 1 }, fetchFn)).rejects.toMatchObject({
      code: 'NETWORK_ERROR',
      status: 0,
    })
  })
})
