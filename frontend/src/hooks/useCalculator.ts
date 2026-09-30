import { useCallback, useRef, useState } from 'react'
import { ApiError, calculate } from '../api/client'
import { getOperation, type Operation } from '../lib/operations'
import { parseOperand } from '../lib/validation'

export interface FieldErrors {
  a?: string
  b?: string
}

export interface CalculationResult {
  operation: Operation
  operands: number[]
  value: number
}

// User-facing copy for known backend error codes.
const FRIENDLY_ERRORS: Record<string, string> = {
  DIVISION_BY_ZERO: 'Cannot divide by zero.',
  NEGATIVE_SQRT: 'Cannot take the square root of a negative number.',
  RESULT_NOT_FINITE: 'The result is too large or undefined.',
  NETWORK_ERROR: 'Could not reach the server. Is the backend running?',
}

export function useCalculator() {
  const [a, setA] = useState('')
  const [b, setB] = useState('')
  const [operation, setOperation] = useState<Operation>('add')
  const [result, setResult] = useState<CalculationResult | null>(null)
  const [error, setError] = useState<string | null>(null)
  const [fieldErrors, setFieldErrors] = useState<FieldErrors>({})
  const [loading, setLoading] = useState(false)

  // Guards against out-of-order responses when the user submits repeatedly.
  const latestRequest = useRef(0)

  const isUnary = getOperation(operation).arity === 1

  const submit = useCallback(async () => {
    const parsedA = parseOperand(a, 'A')
    const parsedB = isUnary ? null : parseOperand(b, 'B')

    const errors: FieldErrors = {}
    if (!parsedA.ok) errors.a = parsedA.error
    if (parsedB && !parsedB.ok) errors.b = parsedB.error
    setFieldErrors(errors)
    setError(null)

    if (!parsedA.ok || (parsedB && !parsedB.ok)) {
      setResult(null)
      return
    }

    const requestId = ++latestRequest.current
    setLoading(true)
    try {
      const res = await calculate({
        operation,
        a: parsedA.value,
        ...(parsedB ? { b: parsedB.value } : {}),
      })
      if (requestId !== latestRequest.current) return
      setResult({ operation: res.operation, operands: res.operands, value: res.result })
    } catch (err) {
      if (requestId !== latestRequest.current) return
      setResult(null)
      setError(
        err instanceof ApiError
          ? (FRIENDLY_ERRORS[err.code] ?? err.message)
          : 'Something went wrong. Please try again.',
      )
    } finally {
      if (requestId === latestRequest.current) setLoading(false)
    }
  }, [a, b, operation, isUnary])

  const reset = useCallback(() => {
    latestRequest.current++
    setA('')
    setB('')
    setResult(null)
    setError(null)
    setFieldErrors({})
    setLoading(false)
  }, [])

  return {
    a, setA,
    b, setB,
    operation, setOperation,
    isUnary,
    result,
    error,
    fieldErrors,
    loading,
    submit,
    reset,
  }
}
