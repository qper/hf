import { useMutation, useQueryClient } from '@tanstack/react-query'
import { useState } from 'react'
import type { ChangeEvent, KeyboardEvent } from 'react'
import { createEntry } from '@/api/board'
import type { Board, BoardHabit } from '@/api/board'

type DurationHabitInputProps = {
  habit: BoardHabit
  date: string
  isEditable: boolean
  entryValue?: number | null
}

export function DurationHabitInput({
  habit,
  date,
  isEditable,
  entryValue,
}: DurationHabitInputProps) {
  const queryClient = useQueryClient()
  const [inputValue, setInputValue] = useState(entryValue?.toString() ?? '')

  const createEntryMutation = useMutation({
    mutationFn: async (value: number) => {
      return createEntry({
        habit_id: habit.id,
        date,
        completed: value > 0,
        value,
      })
    },
    onMutate: async (value: number) => {
      await queryClient.cancelQueries({ queryKey: ['board', date] })
      const previousData = queryClient.getQueryData<Board>(['board', date])

      queryClient.setQueryData(['board', date], (old: Board | undefined) => {
        if (!old) return old
        const wasCompleted = habit.is_completed
        const isNowCompleted = value > 0

        return {
          ...old,
          progress: {
            ...old.progress,
            done:
              wasCompleted && !isNowCompleted
                ? old.progress.done - 1
                : !wasCompleted && isNowCompleted
                  ? old.progress.done + 1
                  : old.progress.done,
          },
          habits: old.habits.map((h) =>
            h.id === habit.id ? { ...h, is_completed: isNowCompleted } : h,
          ),
        }
      })

      return { previousData }
    },
    onError: (_err, _value, context) => {
      if (context?.previousData) {
        queryClient.setQueryData(['board', date], context.previousData)
      }
    },
  })

  const handleInputChange = (e: ChangeEvent<HTMLInputElement>) => {
    setInputValue(e.target.value)
  }

  const handleInputBlur = () => {
    const value = parseFloat(inputValue)
    if (!isNaN(value) && value >= 0) {
      createEntryMutation.mutate(value)
    }
  }

  const handleInputKeyDown = (e: KeyboardEvent<HTMLInputElement>) => {
    if (e.key === 'Enter') {
      handleInputBlur()
    }
  }

  return (
    <div
      className="flex items-center"
      style={{
        minHeight: '44px',
        display: 'flex',
        alignItems: 'center',
      }}
    >
      <input
        type="number"
        inputMode="numeric"
        value={inputValue}
        onChange={handleInputChange}
        onBlur={handleInputBlur}
        onKeyDown={handleInputKeyDown}
        disabled={!isEditable || createEntryMutation.isPending}
        placeholder="0 мин"
        className="w-24 text-sm bg-zinc-800 text-zinc-100 rounded px-3 py-2 border border-zinc-700 focus:border-cyan-500 outline-none placeholder-zinc-600 disabled:opacity-50 disabled:cursor-not-allowed"
        min="0"
      />
      <span className="ml-2 text-xs text-zinc-500">мин</span>
    </div>
  )
}
