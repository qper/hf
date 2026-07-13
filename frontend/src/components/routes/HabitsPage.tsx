import { useEffect, useMemo, useState } from 'react'
import { MoreHorizontal, GripVertical, Archive, Trash2, Plus } from 'lucide-react'
import {
  closestCenter,
  DndContext,
  KeyboardSensor,
  PointerSensor,
  TouchSensor,
  type DragEndEvent,
  useSensor,
  useSensors,
} from '@dnd-kit/core'
import {
  SortableContext,
  sortableKeyboardCoordinates,
  useSortable,
  verticalListSortingStrategy,
} from '@dnd-kit/sortable'
import { CSS } from '@dnd-kit/utilities'
import { HabitFormDialog, type HabitFormValues } from './HabitFormDialog'
import { reorderHabits } from './habitReorder'

type Habit = {
  id: string
  name: string
  icon: string
  streak: number
  archived: boolean
  category: string
  type: 'boolean' | 'numeric' | 'duration'
  color: string
  targetValue?: number
  unit?: string
  frequency: 'daily' | 'custom' | 'weekly'
  selectedDays?: string[]
  weeklyCount?: string
}

const STORAGE_KEY = 'hf-habits-order'

const initialHabits: Habit[] = [
  {
    id: '1',
    name: 'Read 20 min',
    icon: '📖',
    streak: 7,
    archived: false,
    category: 'All',
    type: 'boolean',
    color: '#14b8a6',
    frequency: 'daily',
  },
  {
    id: '2',
    name: 'Morning run',
    icon: '🏃',
    streak: 4,
    archived: false,
    category: 'Health',
    type: 'numeric',
    color: '#f59e0b',
    targetValue: 5,
    unit: 'км',
    frequency: 'weekly',
  },
  {
    id: '3',
    name: 'Meditation',
    icon: '🧘',
    streak: 2,
    archived: false,
    category: 'All',
    type: 'boolean',
    color: '#8b5cf6',
    frequency: 'daily',
  },
  {
    id: '4',
    name: 'Deep work',
    icon: '💼',
    streak: 6,
    archived: true,
    category: 'Work',
    type: 'boolean',
    color: '#ef4444',
    frequency: 'daily',
  },
]

const categories = ['All', 'Health', 'Work', 'Learning']

function getStoredHabits(baseHabits: Habit[]): Habit[] {
  if (typeof window === 'undefined') {
    return baseHabits
  }

  try {
    const savedOrder = window.localStorage.getItem(STORAGE_KEY)
    if (!savedOrder) {
      return baseHabits
    }

    const orderedIds = JSON.parse(savedOrder) as string[]
    const orderedHabits = orderedIds
      .map((id) => baseHabits.find((habit) => habit.id === id))
      .filter((habit): habit is Habit => Boolean(habit))
    const remainingHabits = baseHabits.filter((habit) => !orderedIds.includes(habit.id))

    return [...orderedHabits, ...remainingHabits]
  } catch {
    return baseHabits
  }
}

type SortableHabitRowProps = {
  habit: Habit
  menuOpenId: string | null
  setMenuOpenId: (value: string | null) => void
  onArchive: (habitId: string) => void
  onDelete: (habitId: string) => void
  onEdit: (habit: Habit) => void
}

function SortableHabitRow({
  habit,
  menuOpenId,
  setMenuOpenId,
  onArchive,
  onDelete,
  onEdit,
}: SortableHabitRowProps) {
  const { attributes, listeners, setNodeRef, transform, transition, isDragging } = useSortable({ id: habit.id })
  const style = {
    transform: CSS.Transform.toString(transform),
    transition,
  }

  return (
    <div
      ref={setNodeRef}
      style={style}
      data-testid={`habit-row-${habit.id}`}
      className={`flex items-center gap-3 rounded-2xl border border-zinc-800 bg-zinc-950/70 p-4 ${isDragging ? 'border-cyan-500/70 shadow-lg shadow-cyan-500/10' : ''}`}
    >
      <button
        type="button"
        className={`text-zinc-500 ${isDragging ? 'cursor-grabbing' : 'cursor-grab'}`}
        aria-label={`drag-${habit.id}`}
        {...attributes}
        {...listeners}
      >
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
          onClick={() => onArchive(habit.id)}
          className="rounded-full p-2 text-zinc-400 hover:bg-zinc-800 hover:text-white"
        >
          <Archive size={16} />
        </button>
        <div className="relative">
          <button
            type="button"
            aria-label="more"
            onClick={() => setMenuOpenId((current) => (current === habit.id ? null : habit.id))}
            className="rounded-full p-2 text-zinc-400 hover:bg-zinc-800 hover:text-white"
          >
            <MoreHorizontal size={16} />
          </button>
          <div className={`absolute right-0 z-10 mt-2 min-w-[140px] rounded-xl border border-zinc-800 bg-zinc-900 p-2 shadow-lg ${menuOpenId === habit.id ? 'block' : 'hidden'}`}>
            <button
              type="button"
              role="menuitem"
              onClick={() => {
                setMenuOpenId(null)
                onEdit(habit)
              }}
              className="flex w-full items-center gap-2 rounded-lg px-3 py-2 text-sm text-zinc-200 hover:bg-zinc-800"
            >
              Edit
            </button>
            <button
              type="button"
              role="menuitem"
              aria-label="delete"
              onClick={() => onDelete(habit.id)}
              className="flex w-full items-center gap-2 rounded-lg px-3 py-2 text-sm text-rose-300 hover:bg-zinc-800"
            >
              <Trash2 size={14} />
              Delete
            </button>
          </div>
        </div>
      </div>
    </div>
  )
}

export function HabitsPage() {
  const [activeCategory, setActiveCategory] = useState('All')
  const [habits, setHabits] = useState<Habit[]>(() => getStoredHabits(initialHabits))
  const [archiveOpen, setArchiveOpen] = useState(false)
  const [isDialogOpen, setIsDialogOpen] = useState(false)
  const [dialogMode, setDialogMode] = useState<'create' | 'edit'>('create')
  const [editingHabit, setEditingHabit] = useState<Habit | null>(null)
  const [menuOpenId, setMenuOpenId] = useState<string | null>(null)

  const sensors = useSensors(
    useSensor(PointerSensor, { activationConstraint: { distance: 8 } }),
    useSensor(TouchSensor, { activationConstraint: { delay: 100, tolerance: 8 } }),
    useSensor(KeyboardSensor, { coordinateGetter: sortableKeyboardCoordinates }),
  )

  const visibleHabits = useMemo(() => {
    return habits.filter((habit) => {
      const categoryMatch = activeCategory === 'All' || habit.category === activeCategory
      return categoryMatch && !habit.archived
    })
  }, [activeCategory, habits])

  const archivedHabits = useMemo(() => {
    return habits.filter((habit) => habit.archived)
  }, [habits])

  useEffect(() => {
    if (typeof window === 'undefined' || !window.localStorage) {
      return
    }

    window.localStorage.setItem(STORAGE_KEY, JSON.stringify(habits.filter((habit) => !habit.archived).map((habit) => habit.id)))
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

  const openCreateDialog = () => {
    setEditingHabit(null)
    setDialogMode('create')
    setIsDialogOpen(true)
  }

  const openEditDialog = (habit: Habit) => {
    setEditingHabit(habit)
    setDialogMode('edit')
    setIsDialogOpen(true)
  }

  const handleSaveHabit = (values: HabitFormValues) => {
    const payload: Habit = {
      id: editingHabit?.id ?? `${Date.now()}`,
      name: values.name,
      icon: values.icon,
      streak: editingHabit?.streak ?? 0,
      archived: editingHabit?.archived ?? false,
      category: values.category || 'All',
      type: values.type,
      color: values.color,
      targetValue: values.targetValue ? Number(values.targetValue) : undefined,
      unit: values.type === 'duration' ? 'мин' : values.unit || undefined,
      frequency: values.frequency,
      selectedDays: values.selectedDays,
      weeklyCount: values.weeklyCount,
    }

    setHabits((current) => {
      if (editingHabit) {
        return current.map((habit) => (habit.id === editingHabit.id ? payload : habit))
      }
      return [...current, payload]
    })
    setIsDialogOpen(false)
  }

  const handleReorder = async (activeId: string, overId: string) => {
    const previousHabits = habits

    if (activeId === overId) {
      return
    }

    const nextHabits = reorderHabits(habits, activeId, overId)
    setHabits(nextHabits)

    try {
      const response = await fetch('/api/v1/habits/reorder', {
        method: 'PATCH',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ order: nextHabits.filter((habit) => !habit.archived).map((habit) => habit.id) }),
      })

      if (!response.ok) {
        throw new Error('Failed to reorder habits')
      }
    } catch {
      setHabits(previousHabits)
    }
  }

  const handleDragEnd = (event: DragEndEvent) => {
    const { active, over } = event

    if (over && active.id !== over.id) {
      void handleReorder(String(active.id), String(over.id))
    }
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-3xl font-semibold text-white">Привычки</h2>
          <p className="text-sm text-zinc-400">Твои цели и привычки</p>
        </div>
        <button
          type="button"
          onClick={openCreateDialog}
          className="flex items-center gap-2 rounded-full bg-cyan-500 px-4 py-2 text-sm font-medium text-zinc-950 transition hover:bg-cyan-400"
        >
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

      <DndContext sensors={sensors} collisionDetection={closestCenter} onDragEnd={handleDragEnd}>
        <SortableContext items={visibleHabits.map((habit) => habit.id)} strategy={verticalListSortingStrategy}>
          <div className="space-y-3">
            {visibleHabits.map((habit) => (
              <SortableHabitRow
                key={habit.id}
                habit={habit}
                menuOpenId={menuOpenId}
                setMenuOpenId={setMenuOpenId}
                onArchive={handleArchive}
                onDelete={handleDelete}
                onEdit={openEditDialog}
              />
            ))}
          </div>
        </SortableContext>
      </DndContext>

      <HabitFormDialog
        open={isDialogOpen}
        onOpenChange={setIsDialogOpen}
        mode={dialogMode}
        onSubmit={handleSaveHabit}
        defaultValues={editingHabit ? {
          name: editingHabit.name,
          type: editingHabit.type,
          category: editingHabit.category,
          color: editingHabit.color,
          icon: editingHabit.icon,
          targetValue: editingHabit.targetValue?.toString() ?? '',
          unit: editingHabit.unit ?? '',
          frequency: editingHabit.frequency,
          selectedDays: editingHabit.selectedDays ?? [],
          weeklyCount: editingHabit.weeklyCount ?? '1',
        } : undefined}
      />

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
