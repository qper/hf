import { useQuery } from '@tanstack/react-query'

const fetchGreeting = async () => {
  await new Promise((resolve) => window.setTimeout(resolve, 100))
  return 'Welcome back to HabitFlow'
}

export function HomePage() {
  const { data, isPending } = useQuery({
    queryKey: ['home-greeting'],
    queryFn: fetchGreeting,
    staleTime: 30_000,
    retry: 2,
  })

  return (
    <div className="space-y-4">
      <p className="text-sm uppercase tracking-[0.3em] text-cyan-400">Dashboard</p>
      <h2 className="text-3xl font-semibold">Home route</h2>
      <p className="max-w-2xl text-zinc-400">
        This page uses React Query to cache data across renders and route changes.
      </p>
      <div className="rounded-xl border border-zinc-800 bg-zinc-950/70 p-4">
        {isPending ? (
          <p className="text-sm text-zinc-400">Loading greeting...</p>
        ) : (
          <p className="text-lg text-zinc-100">{data}</p>
        )}
      </div>
    </div>
  )
}
