import type { BoardHabit } from '@/api/board'
import type { CSSProperties } from 'react'
import { BooleanHabitToggle } from './BooleanHabitToggle'
import { NumericHabitInput } from './NumericHabitInput'
import { DurationHabitInput } from './DurationHabitInput'
import { Zap } from 'lucide-react'

type HabitRowProps = {
  habit: BoardHabit
  date: string
  isEditable: boolean
  entryValue?: number | null
}

// Habit color mapping for streak pulse animation
const habitColorMap: Record<string, string> = {
  default: 'rgba(234, 179, 8, 0.4)', // orange-500 with 0.4 alpha
  red: 'rgba(239, 68, 68, 0.4)',
  blue: 'rgba(59, 130, 246, 0.4)',
  green: 'rgba(34, 197, 94, 0.4)',
  purple: 'rgba(147, 51, 234, 0.4)',
  pink: 'rgba(236, 72, 153, 0.4)',
}

export function HabitRow({
  habit,
  date,
  isEditable,
  entryValue,
}: HabitRowProps) {
  const shouldShowPulse = habit.streak >= 7
  const habitColor = habitColorMap.default

  return (
    <div
      style={shouldShowPulse ? { '--habit-color-alpha': habitColor } as CSSProperties : undefined}
    >
      <div
        className="rounded-lg border border-zinc-800 bg-zinc-950/50 p-4 hover:border-zinc-700 transition-colors"
      >
        <div className="flex items-start justify-between gap-4">
          <div className="flex-1 min-w-0">
            <div className="flex items-center gap-2 mb-1">
              <h3 className="font-medium text-zinc-100 truncate">{habit.name}</h3>
              {habit.streak > 0 && (
                <div
                  className={`flex items-center gap-1 text-xs bg-orange-500/10 text-orange-400 px-2 py-1 rounded flex-shrink-0 ${
                    shouldShowPulse ? 'streak-pulse' : ''
                  }`}
                >
                  <Zap className="h-3 w-3" />
                  {habit.streak}
                </div>
              )}
            </div>
            {habit.description && (
              <p className="text-xs text-zinc-500 line-clamp-1">{habit.description}</p>
            )}
          </div>

          <div className="flex-shrink-0">
            {habit.type === 'boolean' && (
              <BooleanHabitToggle habit={habit} date={date} isEditable={isEditable} />
            )}
            {habit.type === 'numeric' && (
              <NumericHabitInput
                habit={habit}
                date={date}
                isEditable={isEditable}
                entryValue={entryValue}
              />
            )}
            {habit.type === 'duration' && (
              <DurationHabitInput
                habit={habit}
                date={date}
                isEditable={isEditable}
                entryValue={entryValue}
              />
            )}
          </div>
        </div>
      </div>
      {/* Separator */}
      <div className="h-px bg-gradient-to-r from-transparent via-zinc-700 to-transparent" />
    </div>
  )
}
