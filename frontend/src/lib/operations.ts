// Mirrors the operations exposed by the backend. Kept static on purpose:
// the UI needs labels and symbols the API does not (and should not) know about.

export type Operation =
  | 'add'
  | 'subtract'
  | 'multiply'
  | 'divide'
  | 'power'
  | 'sqrt'
  | 'percentage'

export interface OperationMeta {
  id: Operation
  symbol: string
  label: string
  arity: 1 | 2
}

export const OPERATIONS: readonly OperationMeta[] = [
  { id: 'add', symbol: '+', label: 'Add', arity: 2 },
  { id: 'subtract', symbol: '−', label: 'Subtract', arity: 2 },
  { id: 'multiply', symbol: '×', label: 'Multiply', arity: 2 },
  { id: 'divide', symbol: '÷', label: 'Divide', arity: 2 },
  { id: 'power', symbol: 'xʸ', label: 'Power', arity: 2 },
  { id: 'sqrt', symbol: '√', label: 'Square root', arity: 1 },
  { id: 'percentage', symbol: '%', label: 'Percentage (A% of B)', arity: 2 },
]

export function getOperation(id: Operation): OperationMeta {
  const op = OPERATIONS.find((o) => o.id === id)
  if (!op) throw new Error(`Unknown operation: ${id}`)
  return op
}
