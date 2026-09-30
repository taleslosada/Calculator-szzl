import type { FormEvent } from 'react'
import { useCalculator, type CalculationResult } from '../hooks/useCalculator'
import { getOperation, OPERATIONS } from '../lib/operations'
import { formatResult } from '../lib/validation'

export function Calculator() {
  const calc = useCalculator()

  const handleSubmit = (e: FormEvent) => {
    e.preventDefault()
    void calc.submit()
  }

  return (
    <form className="calculator" onSubmit={handleSubmit} noValidate>
      <h1>Calculator</h1>

      <div className="display" aria-live="polite" data-testid="display">
        {calc.error ? (
          <p className="display__error" role="alert">{calc.error}</p>
        ) : calc.result ? (
          <>
            <p className="display__expression">{describe(calc.result)}</p>
            <p className="display__value" data-testid="result">
              {formatResult(calc.result.value)}
            </p>
          </>
        ) : (
          <p className="display__placeholder">Enter values and pick an operation</p>
        )}
      </div>

      <div className="fields">
        <Field
          id="a"
          label={calc.isUnary ? 'Value' : 'A'}
          value={calc.a}
          onChange={calc.setA}
          error={calc.fieldErrors.a}
        />
        {!calc.isUnary && (
          <Field id="b" label="B" value={calc.b} onChange={calc.setB} error={calc.fieldErrors.b} />
        )}
      </div>

      <fieldset className="operations">
        <legend>Operation</legend>
        {OPERATIONS.map((op) => (
          <button
            key={op.id}
            type="button"
            className="op"
            aria-pressed={calc.operation === op.id}
            aria-label={op.label}
            title={op.label}
            onClick={() => calc.setOperation(op.id)}
          >
            {op.symbol}
          </button>
        ))}
      </fieldset>

      <div className="actions">
        <button type="button" className="btn btn--secondary" onClick={calc.reset}>
          Clear
        </button>
        <button type="submit" className="btn btn--primary" disabled={calc.loading}>
          {calc.loading ? 'Calculating…' : '='}
        </button>
      </div>
    </form>
  )
}

interface FieldProps {
  id: string
  label: string
  value: string
  onChange: (value: string) => void
  error?: string
}

function Field({ id, label, value, onChange, error }: FieldProps) {
  const errorId = `${id}-error`
  return (
    <div className="field">
      <label htmlFor={id}>{label}</label>
      <input
        id={id}
        type="text"
        inputMode="decimal"
        autoComplete="off"
        value={value}
        onChange={(e) => onChange(e.target.value)}
        aria-invalid={Boolean(error)}
        aria-describedby={error ? errorId : undefined}
      />
      {error && (
        <span id={errorId} className="field__error">
          {error}
        </span>
      )}
    </div>
  )
}

function describe({ operation, operands }: CalculationResult): string {
  const [a, b] = operands.map(formatResult)
  switch (operation) {
    case 'sqrt':
      return `√${a} =`
    case 'percentage':
      return `${a}% of ${b} =`
    case 'power':
      return `${a} ^ ${b} =`
    default:
      return `${a} ${getOperation(operation).symbol} ${b} =`
  }
}
