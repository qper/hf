import { create } from 'zustand'
import { createJSONStorage, persist } from 'zustand/middleware'

type UIState = {
  isSidebarOpen: boolean
  toggleSidebar: () => void
}

export const useUIStore = create<UIState>()(
  persist(
    (set) => ({
      isSidebarOpen: true,
      toggleSidebar: () =>
        set((state) => ({ isSidebarOpen: !state.isSidebarOpen })),
    }),
    {
      name: 'hf-ui-store',
      storage: createJSONStorage(() => sessionStorage),
    },
  ),
)
