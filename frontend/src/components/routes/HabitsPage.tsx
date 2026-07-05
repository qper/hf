import { useMemo, useState } from 'react'
import { MoreHorizontal, GripVertical, Archive, Trash2, Plus } from 'lucide-react'

type Habit = {
  id: string
  name: string
  icon: string
  streak: number
  archived: boolean
  category: string
}

const initialHabits: Habit[] = [
  {
    id: '1',
    name: 'Read 20 min',
    icon: '📖',
    streak: 7,
    archived: false,
    category: 'All',
  },
  {
    id: '2',
    name: 'Morning run',
    icon: '🏃',
    streak: 4,
    archived: true,
    category: 'Health',
  },
]

const categories = ['All', 'Health', 'Work', 'Learning']

export function HabitsPage() {
  const [activeCategory, setActiveCategory] = useState('All')
  const [habits, setHabits] = useState(initialHabits)
  const [archiveOpen, setArchiveOpen] = useState(false)

  const visibleHabits = useMemo(() => {
    return habits.filter((habit) => {
      const categoryMatch =
        activeCategory === 'All' || habit.category === activeCategory
      return categoryMatch && !habit.archived
    })
  }, [activeCategory, habits])

  const archivedHabits = useMemo(() => {
    return habits.filter((habit) => habit.archived)
  }, [habits])

  const handleArchive = (habitId: string) => {
    setHabits((current) =>
      current.map((habit) =>
        habit.id === habitId ? { ...habit, archived: !habit.archived } : habit,
      ),
    )
  }

  const handleDelete = (habitId: string) => {
    const confirmed = window.confirm('Delete this habit?')
    if (!confirmed) return
    setHabits((current) => current.filter((habit) => habit.id !== habitId))
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-3xl font-semibold text-white">Привычки</h2>
          <p className="text-sm text-zinc-400">Твои цели и привычки</p>
        </div>
        <button className="flex items-center gap-2 rounded-full bg-cyan-500 px-4 py-2 text-sm font-medium text-zinc-950 transition hover:bg-cyan-400">
          <Plus size={16} />
          Новая
        </button>
      </div>

      <div className="flex gap-2 overflow-x-auto pb-2">
        {categories.map((category) => (
          <button
            key={category}
            type="button"
            onClick={() => setActiveCategory(category)}
            className={`rounded-full px-4 py-2 text-sm whitespace-nowrap transition ${
              activeCategory === category
                ? 'bg-cyan-500 text-zinc-950'
                : 'bg-zinc-800 text-zinc-300 hover:bg-zinc-700'
            }`}
          >
            {category}
          </button>
        ))}
      </div>

      <div className="space-y-3">
        {visibleHabits.map((habit) => (
          <div
            key={habit.id}
            data-testid={`habit-row-${habit.id}`}
            className="flex items-center gap-3 rounded-2xl border border-zinc-800 bg-zinc-950/70 p-4"
          >
            <button type="button" className="text-zinc-500" aria-label="drag">
              <GripVertical size={18} />
            </button>
            <div className="flex h-10 w-10 items-center justify-center rounded-full bg-zinc-800 text-lg">
              {habit.icon}
            </div>
            <div className="min-w-0 flex-1">
              <p className="font-medium text-white">{habit.name}</p>
              <p className="text-sm text-zinc-400">{habit.category}</p>
            </div>
            <div className="hidden items-center gap-1 text-sm text-orange-400 sm:flex">
              <span>🔥</span>
              <span>{habit.streak}</span>
            </div>
            <div className="flex items-center gap-2">
              <button
                type="button"
                aria-label="archive"
                onClick={() => handleArchive(habit.id)}
                className="rounded-full p-2 text-zinc-400 hover:bg-zinc-800 hover:text-white"
              >
                <Archive size={16} />
              </button>
              <div className="relative">
                <button
                  type="button"
                  aria-label="more"
                  className="rounded-full p-2 text-zinc-400 hover:bg-zinc-800 hover:text-white"
                >
                  <MoreHorizontal size={16} />
                </button>
                <div className="absolute right-0 z-10 mt-2 hidden min-w-[140px] rounded-xl border border-zinc-800 bg-zinc-900 p-2 shadow-lg">
                  <button type="button" role="menuitem" className="flex w-full items-center gap-2 rounded-lg px-3 py-2 text-sm text-zinc-200 hover:bg-zinc-800">
                    Edit
                  </button>
                  <button
                    type="button"
                    role="menuitem"
                    aria-label="delete"
                    onClick={() => handleDelete(habit.id)}
                    className="flex w-full items-center gap-2 rounded-lg px-3 py-2 text-sm text-rose-300 hover:bg-zinc-800"
                  >
                    <Trash2 size={14} />
                    Delete
                  </button>
                </div>
              </div>
            </div>
          </div>
        ))}
      </div>

      <div className="rounded-2xl border border-zinc-800 bg-zinc-950/60 p-4">
        <button
          type="button"
          onClick={() => setArchiveOpen((current) => !current)}
          className="flex w-full items-center justify-between text-left"
        >
          <span className="font-medium text-white">Архив ({archivedHabits.length})</span>
          <span className="text-sm text-zinc-400">{archiveOpen ? '▾' : '▸'}</span>
        </button>
        {archiveOpen ? (
          <div className="mt-4 space-y-2">
            {archivedHabits.map((habit) => (
              <div key={habit.id} className="flex items-center justify-between rounded-xl border border-zinc-800 bg-zinc-900/70 px-3 py-2 text-sm text-zinc-300">
                <span>{habit.name}</span>
                <button
                  type="button"
                  onClick={() => handleArchive(habit.id)}
                  className="text-cyan-400"
                >
                  Restore
                </button>
              </div>
            ))}
          </div>
        ) : null}
      </div>
    </div>
  )
}
