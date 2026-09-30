import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { ApiError, calculate } from '../api/client'
import { Calculator } from './Calculator'

// Mock only the network call; keep the real ApiError class.
vi.mock('../api/client', async (importOriginal) => {
  const actual = await importOriginal<typeof import('../api/client')>()
  return { ...actual, calculate: vi.fn() }
})

const mockedCalculate = vi.mocked(calculate)

function setup() {
  const user = userEvent.setup()
  render(<Calculator />)
  return {
    user,
    inputA: () => screen.getByLabelText(/^(A|Value)$/),
    inputB: () => screen.queryByLabelText('B'),
    submit: () => user.click(screen.getByRole('button', { name: '=' })),
  }
}

describe('Calculator', () => {
  beforeEach(() => {
    mockedCalculate.mockReset()
  })

  it('shows the result of a successful calculation', async () => {
    mockedCalculate.mockResolvedValue({ operation: 'add', operands: [2, 3], result: 5 })
    const { user, inputA, inputB, submit } = setup()

    await user.type(inputA(), '2')
    await user.type(inputB()!, '3')
    await submit()

    expect(mockedCalculate).toHaveBeenCalledWith({ operation: 'add', a: 2, b: 3 })
    expect(await screen.findByTestId('result')).toHaveTextContent('5')
    expect(screen.getByText('2 + 3 =')).toBeInTheDocument()
  })

  it('sends the selected operation', async () => {
    mockedCalculate.mockResolvedValue({ operation: 'multiply', operands: [4, 5], result: 20 })
    const { user, inputA, inputB, submit } = setup()

    await user.click(screen.getByRole('button', { name: 'Multiply' }))
    await user.type(inputA(), '4')
    await user.type(inputB()!, '5')
    await submit()

    expect(mockedCalculate).toHaveBeenCalledWith({ operation: 'multiply', a: 4, b: 5 })
    expect(screen.getByRole('button', { name: 'Multiply' })).toHaveAttribute('aria-pressed', 'true')
  })

  it('hides the second input for square root and sends a single operand', async () => {
    mockedCalculate.mockResolvedValue({ operation: 'sqrt', operands: [16], result: 4 })
    const { user, inputA, inputB, submit } = setup()

    await user.click(screen.getByRole('button', { name: 'Square root' }))
    expect(inputB()).not.toBeInTheDocument()

    await user.type(inputA(), '16')
    await submit()

    expect(mockedCalculate).toHaveBeenCalledWith({ operation: 'sqrt', a: 16 })
    expect(await screen.findByText('√16 =')).toBeInTheDocument()
  })

  it('shows field errors and does not call the API on invalid input', async () => {
    const { user, inputA, inputB, submit } = setup()

    await user.type(inputA(), 'abc')
    await submit()

    expect(screen.getByText('A must be a valid number')).toBeInTheDocument()
    expect(screen.getByText('B is required')).toBeInTheDocument()
    expect(inputA()).toHaveAttribute('aria-invalid', 'true')
    expect(inputB()).toHaveAttribute('aria-invalid', 'true')
    expect(mockedCalculate).not.toHaveBeenCalled()
  })

  it('shows a friendly message for backend domain errors', async () => {
    mockedCalculate.mockRejectedValue(new ApiError('DIVISION_BY_ZERO', 'division by zero', 422))
    const { user, inputA, inputB, submit } = setup()

    await user.click(screen.getByRole('button', { name: 'Divide' }))
    await user.type(inputA(), '1')
    await user.type(inputB()!, '0')
    await submit()

    expect(await screen.findByRole('alert')).toHaveTextContent('Cannot divide by zero.')
  })

  it('falls back to the backend message for unmapped error codes', async () => {
    mockedCalculate.mockRejectedValue(new ApiError('SOMETHING_NEW', 'custom message', 400))
    const { user, inputA, inputB, submit } = setup()

    await user.type(inputA(), '1')
    await user.type(inputB()!, '1')
    await submit()

    expect(await screen.findByRole('alert')).toHaveTextContent('custom message')
  })

  it('shows a generic message for unexpected errors', async () => {
    mockedCalculate.mockRejectedValue(new Error('boom'))
    const { user, inputA, inputB, submit } = setup()

    await user.type(inputA(), '1')
    await user.type(inputB()!, '1')
    await submit()

    expect(await screen.findByRole('alert')).toHaveTextContent('Something went wrong')
  })

  it('submits with the Enter key', async () => {
    mockedCalculate.mockResolvedValue({ operation: 'add', operands: [1, 1], result: 2 })
    const { user, inputA, inputB } = setup()

    await user.type(inputA(), '1')
    await user.type(inputB()!, '1{Enter}')

    expect(await screen.findByTestId('result')).toHaveTextContent('2')
  })

  it('clears inputs and result', async () => {
    mockedCalculate.mockResolvedValue({ operation: 'add', operands: [1, 1], result: 2 })
    const { user, inputA, inputB, submit } = setup()

    await user.type(inputA(), '1')
    await user.type(inputB()!, '1')
    await submit()
    await screen.findByTestId('result')

    await user.click(screen.getByRole('button', { name: 'Clear' }))

    expect(inputA()).toHaveValue('')
    expect(inputB()).toHaveValue('')
    expect(screen.queryByTestId('result')).not.toBeInTheDocument()
  })

  it.each([
    ['Percentage (A% of B)', 'percentage', '15% of 200 ='],
    ['Power', 'power', '15 ^ 200 ='],
  ] as const)('describes %s results', async (label, operation, expression) => {
    mockedCalculate.mockResolvedValue({ operation, operands: [15, 200], result: 30 })
    const { user, inputA, inputB, submit } = setup()

    await user.click(screen.getByRole('button', { name: label }))
    await user.type(inputA(), '15')
    await user.type(inputB()!, '200')
    await submit()

    expect(await screen.findByText(expression)).toBeInTheDocument()
  })
})
