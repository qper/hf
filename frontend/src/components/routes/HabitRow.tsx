import { BoardHabit } from '@/api/board'
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

export function HabitRow({
  habit,
  date,
  isEditable,
  entryValue,
}: HabitRowProps) {
  return (
    <div className="rounded-lg border border-zinc-800 bg-zinc-950/50 p-4 hover:border-zinc-700 transition-colors">
      <div className="flex items-start justify-between gap-4">
        <div className="flex-1 min-w-0">
          <div className="flex items-center gap-2 mb-1">
            <h3 className="font-medium text-zinc-100 truncate">{habit.name}</h3>
            {habit.streak > 0 && (
              <div className="flex items-center gap-1 text-xs bg-orange-500/10 text-orange-400 px-2 py-1 rounded flex-shrink-0">
                <Zap className="h-3 w-3" />
                {habit.streak}
              </div>
            )}
          </div>
          {habit.description && (
            <p className="text-xs text-zinc-500">{habit.description}</p>
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
  )
}
