import { Link, Outlet } from '@tanstack/react-router'
import { useTranslation } from 'react-i18next'
import { useUIStore } from '@/stores/uiStore'

export function RootLayout() {
  const { t } = useTranslation()
  const { isSidebarOpen, toggleSidebar } = useUIStore()

  return (
    <div className="min-h-screen bg-zinc-950 text-zinc-100">
      <header className="border-b border-zinc-800 bg-zinc-900/80 px-6 py-4 backdrop-blur">
        <div className="mx-auto flex max-w-5xl items-center justify-between gap-4">
          <div>
            <p className="text-sm font-semibold uppercase tracking-[0.3em] text-cyan-400">
              HabitFlow
            </p>
            <h1 className="text-xl font-semibold">Modern frontend shell</h1>
          </div>
          <nav className="flex flex-wrap items-center gap-3">
            <Link
              to="/"
              className="rounded-full px-3 py-2 text-sm text-zinc-300 transition hover:bg-zinc-800 hover:text-white"
            >
              {t('nav.home')}
            </Link>
            <Link
              to="/login"
              className="rounded-full px-3 py-2 text-sm text-zinc-300 transition hover:bg-zinc-800 hover:text-white"
            >
              {t('nav.login')}
            </Link>
            <Link
              to="/register"
              className="rounded-full px-3 py-2 text-sm text-zinc-300 transition hover:bg-zinc-800 hover:text-white"
            >
              {t('nav.register')}
            </Link>
            <button
              type="button"
              onClick={toggleSidebar}
              className="rounded-full border border-zinc-700 px-3 py-2 text-sm text-zinc-300 transition hover:bg-zinc-800"
            >
              {isSidebarOpen ? t('nav.hideSidebar') : t('nav.showSidebar')}
            </button>
          </nav>
        </div>
      </header>

      <main className="mx-auto flex max-w-5xl gap-6 px-6 py-8">
        <aside
          className={`w-64 rounded-2xl border border-zinc-800 bg-zinc-900/80 p-4 ${isSidebarOpen ? 'block' : 'hidden'}`}
        >
          <p className="text-sm font-semibold uppercase tracking-[0.25em] text-zinc-400">
            {t('nav.sessionState')}
          </p>
          <p className="mt-2 text-sm text-zinc-300">
            {t('nav.sessionHint')}
          </p>
        </aside>

        <section className="flex-1 rounded-2xl border border-zinc-800 bg-zinc-900/70 p-6 shadow-xl">
          <Outlet />
        </section>
      </main>
    </div>
  )
}
