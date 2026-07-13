export type HabitOrderItem = {
  id: string
}

export function reorderHabits<T extends HabitOrderItem>(items: T[], activeId: string, overId: string): T[] {
  const oldItems = [...items]
  const oldIndex = oldItems.findIndex((item) => item.id === activeId)
  const newIndex = oldItems.findIndex((item) => item.id === overId)

  if (oldIndex === -1 || newIndex === -1 || oldIndex === newIndex) {
    return oldItems
  }

  const [moved] = oldItems.splice(oldIndex, 1)
  oldItems.splice(newIndex, 0, moved)

  return oldItems
}
