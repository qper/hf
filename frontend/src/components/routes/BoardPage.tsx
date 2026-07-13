import { getBoard } from '@/api/board'
import { DateNavBar } from './DateNavBar'
import { HabitRow } from './HabitRow'
import { useQuery } from '@tanstack/react-query'
import { useNavigate, useParams } from '@tanstack/react-router'
import { Skeleton } from '@/components/ui/skeleton'
import { useCallback } from 'react'
import { useTranslation } from 'react-i18next'

type BoardPageProps = {
  date?: string
}

function BoardPageComponent({ date }: BoardPageProps) {
  const { t } = useTranslation()
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
        <h2 className="text-3xl font-semibold">{t('board.title')}</h2>
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
              {t('board.noHabits')}
            </p>
          ) : (
            board.habits.map((habit) => (
              <HabitRow
                key={habit.id}
                habit={habit}
                date={currentDate}
                isEditable={board.is_editable}
              />
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

