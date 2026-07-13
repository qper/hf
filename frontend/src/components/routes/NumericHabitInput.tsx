import { useMutation, useQueryClient } from '@tanstack/react-query'
import { Minus, Plus } from 'lucide-react'
import { useState } from 'react'
import { createEntry, Board, BoardHabit } from '@/api/board'

type NumericHabitInputProps = {
  habit: BoardHabit
  date: string
  isEditable: boolean
  entryValue?: number | null
}

export function NumericHabitInput({
  habit,
  date,
  isEditable,
  entryValue,
}: NumericHabitInputProps) {
  const queryClient = useQueryClient()
  const [isEditing, setIsEditing] = useState(false)
  const [inputValue, setInputValue] = useState(entryValue?.toString() ?? '')

  const createEntryMutation = useMutation({
    mutationFn: async (value: number) => {
      return createEntry({
        habit_id: habit.id,
        date,
        completed: habit.target_value ? value >= habit.target_value : false,
        value,
      })
    },
    onMutate: async (value: number) => {
      await queryClient.cancelQueries({ queryKey: ['board', date] })
      const previousData = queryClient.getQueryData<Board>(['board', date])

      queryClient.setQueryData(['board', date], (old: Board | undefined) => {
        if (!old) return old
        const wasCompleted = habit.is_completed
        const isNowCompleted = habit.target_value ? value >= habit.target_value : false

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

  const handleIncrement = () => {
    if (!isEditable) return
    const current = entryValue ?? 0
    createEntryMutation.mutate(current + 1)
  }

  const handleDecrement = () => {
    if (!isEditable) return
    const current = entryValue ?? 0
    createEntryMutation.mutate(Math.max(0, current - 1))
  }

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setInputValue(e.target.value)
  }

  const handleInputBlur = () => {
    const value = parseFloat(inputValue)
    if (!isNaN(value)) {
      createEntryMutation.mutate(value)
    }
    setIsEditing(false)
  }

  const handleInputKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter') {
      handleInputBlur()
    } else if (e.key === 'Escape') {
      setIsEditing(false)
      setInputValue(entryValue?.toString() ?? '')
    }
  }

  return (
    <div
      className="flex items-center gap-2"
      style={{
        minHeight: '44px',
        display: 'flex',
        alignItems: 'center',
      }}
    >
      <button
        onClick={handleDecrement}
        disabled={!isEditable || createEntryMutation.isPending}
        className="p-1 rounded hover:bg-zinc-800 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
        style={{
          minWidth: '44px',
          minHeight: '44px',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
        }}
      >
        <Minus className="h-4 w-4 text-zinc-400" />
      </button>

      {isEditing ? (
        <input
          autoFocus
          type="number"
          inputMode="decimal"
          value={inputValue}
          onChange={handleInputChange}
          onBlur={handleInputBlur}
          onKeyDown={handleInputKeyDown}
          className="w-16 text-center text-sm bg-zinc-800 text-zinc-100 rounded px-2 py-1 border border-zinc-700 focus:border-cyan-500 outline-none"
        />
      ) : (
        <button
          onClick={() => isEditable && setIsEditing(true)}
          disabled={!isEditable}
          className="flex-1 text-center text-sm text-zinc-300 px-2 py-1 rounded hover:bg-zinc-800/50 disabled:cursor-not-allowed"
          style={{
            minHeight: '44px',
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
          }}
        >
          {entryValue ?? 0} / {habit.target_value ?? '—'}
        </button>
      )}

      <button
        onClick={handleIncrement}
        disabled={!isEditable || createEntryMutation.isPending}
        className="p-1 rounded hover:bg-zinc-800 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
        style={{
          minWidth: '44px',
          minHeight: '44px',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
        }}
      >
        <Plus className="h-4 w-4 text-zinc-400" />
      </button>
    </div>
  )
}
