import { getBoard } from '@/api/board'
import { DateNavBar } from './DateNavBar'
import { useQuery } from '@tanstack/react-query'
import { useNavigate, useParams } from '@tanstack/react-router'
import { Skeleton } from '@/components/ui/skeleton'
import { CheckCircle2, Circle, Zap } from 'lucide-react'
import { useCallback } from 'react'

type BoardPageProps = {
  date?: string
}

function BoardPageComponent({ date }: BoardPageProps) {
  const navigate = useNavigate()
  const today = new Date().toISOString().split('T')[0]
  const currentDate = date || today

  const { data: board, isPending, error } = useQuery({
    queryKey: ['board', currentDate],
    queryFn: () => getBoard(currentDate),
  })

  const handleDateChange = useCallback(
    (newDate: string) => {
      navigate({
        to: '/board/$date',
        params: { date: newDate },
      })
    },
    [navigate],
  )

  if (error) {
    return (
      <div className="space-y-4">
        <p className="text-sm text-red-400">Error loading board: {error.message}</p>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      <div>
        <p className="text-sm uppercase tracking-[0.3em] text-cyan-400">
          Board
        </p>
        <h2 className="text-3xl font-semibold">Daily Tracker</h2>
      </div>

      <DateNavBar
        date={currentDate}
        onDateChange={handleDateChange}
        progress={board?.progress}
      />

      {isPending ? (
        <div className="space-y-3">
          {[...Array(3)].map((_, i) => (
            <div key={i} className="rounded-lg border border-zinc-800 p-4">
              <div className="space-y-2">
                <Skeleton className="h-4 w-32" />
                <Skeleton className="h-3 w-24" />
              </div>
            </div>
          ))}
        </div>
      ) : board ? (
        <div className="space-y-3">
          {board.habits.length === 0 ? (
            <p className="text-sm text-zinc-400 text-center py-8">
              No habits yet. Create one to get started!
            </p>
          ) : (
            board.habits.map((habit) => (
              <div
                key={habit.id}
                className="rounded-lg border border-zinc-800 bg-zinc-950/50 p-4 hover:border-zinc-700 transition-colors"
              >
                <div className="flex items-start justify-between gap-4">
                  <div className="flex-1 min-w-0">
                    <div className="flex items-center gap-2 mb-1">
                      <h3 className="font-medium text-zinc-100 truncate">
                        {habit.name}
                      </h3>
                      {habit.streak > 0 && (
                        <div className="flex items-center gap-1 text-xs bg-orange-500/10 text-orange-400 px-2 py-1 rounded">
                          <Zap className="h-3 w-3" />
                          {habit.streak}
                        </div>
                      )}
                    </div>
                    {habit.description && (
                      <p className="text-xs text-zinc-500">{habit.description}</p>
                    )}
                  </div>
                  <button
                    className="flex-shrink-0 p-1"
                    onClick={() => {
                      // TODO: Handle habit completion toggle
                    }}
                  >
                    {habit.is_completed ? (
                      <CheckCircle2 className="h-6 w-6 text-green-500" />
                    ) : (
                      <Circle className="h-6 w-6 text-zinc-500 hover:text-zinc-400" />
                    )}
                  </button>
                </div>
              </div>
            ))
          )}
        </div>
      ) : null}
    </div>
  )
}

export function BoardPage({ date }: BoardPageProps) {
  return <BoardPageComponent date={date} />
}

export function BoardPageWithRoute() {
  const { date } = useParams({ from: '/board/$date' })
  return <BoardPageComponent date={date} />
}

