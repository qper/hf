import { useMutation, useQueryClient } from '@tanstack/react-query'
import { CheckCircle2, Circle } from 'lucide-react'
import { useState } from 'react'
import { createEntry, Board, BoardHabit } from '@/api/board'

type BooleanHabitToggleProps = {
  habit: BoardHabit
  date: string
  isEditable: boolean
}

export function BooleanHabitToggle({
  habit,
  date,
  isEditable,
}: BooleanHabitToggleProps) {
  const queryClient = useQueryClient()
  const [isAnimating, setIsAnimating] = useState(false)

  const toggleMutation = useMutation({
    mutationFn: async (completed: boolean) => {
      return createEntry({
        habit_id: habit.id,
        date,
        completed,
      })
    },
    onMutate: async (completed: boolean) => {
      setIsAnimating(true)
      // Cancel outgoing refetches to prevent overwriting optimistic update
      await queryClient.cancelQueries({ queryKey: ['board', date] })

      // Snapshot the previous value
      const previousData = queryClient.getQueryData<Board>(['board', date])

      // Optimistically update to the new value
      queryClient.setQueryData(['board', date], (old: Board | undefined) => {
        if (!old) return old
        return {
          ...old,
          progress: {
            ...old.progress,
            done: completed ? old.progress.done + 1 : old.progress.done - 1,
          },
          habits: old.habits.map((h) =>
            h.id === habit.id ? { ...h, is_completed: completed } : h,
          ),
        }
      })

      // Return context with previous data
      return { previousData }
    },
    onError: (_err, _completed, context) => {
      // Rollback to previous data on error
      if (context?.previousData) {
        queryClient.setQueryData(['board', date], context.previousData)
      }
      setIsAnimating(false)
    },
    onSuccess: () => {
      setTimeout(() => setIsAnimating(false), 50)
    },
  })

  const handleToggle = () => {
    if (!isEditable) return
    toggleMutation.mutate(!habit.is_completed)
  }

  return (
    <button
      onClick={handleToggle}
      disabled={!isEditable || toggleMutation.isPending}
      className={`flex-shrink-0 p-1 transition-transform ${
        isAnimating ? 'scale-95' : 'scale-100'
      } ${!isEditable ? 'opacity-50 cursor-not-allowed' : 'cursor-pointer'}`}
      style={{
        minWidth: '44px',
        minHeight: '44px',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
      }}
    >
      {habit.is_completed ? (
        <CheckCircle2
          className="h-6 w-6 text-green-500 transition-all"
          style={{
            filter: isAnimating ? 'brightness(1.3)' : 'brightness(1)',
          }}
        />
      ) : (
        <Circle className="h-6 w-6 text-zinc-500 hover:text-zinc-400 transition-colors" />
      )}
    </button>
  )
}
