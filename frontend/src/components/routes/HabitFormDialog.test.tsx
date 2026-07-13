import { fireEvent, render, screen } from '@testing-library/react'
import { beforeEach, describe, expect, it } from 'vitest'
import '@/i18n'
import i18n from '@/i18n'
import { HabitFormDialog } from './HabitFormDialog'

beforeEach(async () => {
  await i18n.changeLanguage('ru')
})

describe('HabitFormDialog', () => {
  it('requires target value for numeric habits', async () => {
    render(<HabitFormDialog open onOpenChange={() => {}} />)

    fireEvent.click(screen.getByLabelText(/numeric/i))
    fireEvent.click(screen.getByRole('button', { name: /сохранить/i }))

    expect(await screen.findByText(/цель обязательна/i)).toBeTruthy()
  })

  it('disables type editing in edit mode', () => {
    render(
      <HabitFormDialog
        open
        onOpenChange={() => {}}
        mode="edit"
        defaultValues={{ type: 'numeric' }}
      />,
    )

    const numericRadio = document.body.querySelector('input[value="numeric"][disabled]') as HTMLInputElement | null
    expect(numericRadio).not.toBeNull()
    expect(numericRadio?.disabled).toBe(true)
  })
})
