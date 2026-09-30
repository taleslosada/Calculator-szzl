import { describe, expect, it } from 'vitest'
import { formatResult, parseOperand } from './validation'

describe('parseOperand', () => {
  it.each([
    ['42', 42],
    ['-7', -7],
    ['+3', 3],
    ['3.14', 3.14],
    ['.5', 0.5],
    ['5.', 5],
    ['1e3', 1000],
    ['2.5E-2', 0.025],
    ['  12  ', 12],
    ['0', 0],
  ])('accepts %j', (input, expected) => {
    expect(parseOperand(input)).toEqual({ ok: true, value: expected })
  })

  it.each(['abc', '1,5', '1.2.3', '0x1F', 'Infinity', 'NaN', '--1', '1e', '.', '-'])(
    'rejects %j',
    (input) => {
      const result = parseOperand(input, 'A')
      expect(result).toEqual({ ok: false, error: 'A must be a valid number' })
    },
  )

  it('rejects empty input as required', () => {
    expect(parseOperand('   ', 'B')).toEqual({ ok: false, error: 'B is required' })
  })

  it('rejects numbers that overflow to Infinity', () => {
    expect(parseOperand('1e400')).toEqual({ ok: false, error: 'Value is too large' })
  })
})

describe('formatResult', () => {
  it.each([
    [5, '5'],
    [-12, '-12'],
    [0.1 + 0.2, '0.3'],
    [2.5, '2.5'],
    [1 / 3, '0.333333333333'],
    [1e20, '1e+20'],
    [1.5e-9, '1.5e-9'],
    [0, '0'],
  ])('formats %d as %s', (value, expected) => {
    expect(formatResult(value)).toBe(expected)
  })
})
