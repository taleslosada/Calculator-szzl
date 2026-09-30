export type ParseResult =
  | { ok: true; value: number }
  | { ok: false; error: string }

// Accepts optional sign, digits, optional decimal part and optional exponent.
// Rejects things Number() would happily accept, like "", "  ", "0x1F" or "Infinity".
const NUMBER_PATTERN = /^[+-]?(\d+\.?\d*|\.\d+)([eE][+-]?\d+)?$/

/**
 * Parses user input into a finite number.
 * Client-side validation gives instant feedback; the backend still
 * validates everything because the client can never be trusted.
 */
export function parseOperand(raw: string, name = 'Value'): ParseResult {
  const input = raw.trim()
  if (input === '') {
    return { ok: false, error: `${name} is required` }
  }
  if (!NUMBER_PATTERN.test(input)) {
    return { ok: false, error: `${name} must be a valid number` }
  }
  const value = Number(input)
  if (!Number.isFinite(value)) {
    return { ok: false, error: `${name} is too large` }
  }
  return { ok: true, value }
}

/**
 * Formats a result for display, trimming floating point noise
 * (e.g. 0.1 + 0.2 = 0.30000000000000004 → "0.3").
 */
export function formatResult(value: number): string {
  if (Number.isInteger(value) && Math.abs(value) < 1e15) {
    return value.toString()
  }
  const precise = Number.parseFloat(value.toPrecision(12))
  return Math.abs(precise) >= 1e15 || (precise !== 0 && Math.abs(precise) < 1e-6)
    ? precise.toExponential()
    : precise.toString()
}
