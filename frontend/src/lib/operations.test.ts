import { describe, expect, it } from 'vitest'
import { getOperation, OPERATIONS, type Operation } from './operations'

describe('operations', () => {
  it('only square root is unary', () => {
    const unary = OPERATIONS.filter((o) => o.arity === 1).map((o) => o.id)
    expect(unary).toEqual(['sqrt'])
  })

  it('getOperation returns metadata by id', () => {
    expect(getOperation('divide').symbol).toBe('÷')
  })

  it('getOperation throws for unknown ids', () => {
    expect(() => getOperation('modulo' as Operation)).toThrow('Unknown operation: modulo')
  })
})
