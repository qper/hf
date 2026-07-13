import { useQuery } from '@tanstack/react-query'
import { useTranslation } from 'react-i18next'

const fetchGreeting = async (greeting: string) => {
  await new Promise((resolve) => window.setTimeout(resolve, 100))
  return greeting
}

export function HomePage() {
  const { t } = useTranslation()
  const { data, isPending } = useQuery({
    queryKey: ['home-greeting', t('common.homeGreeting')],
    queryFn: () => fetchGreeting(t('common.homeGreeting')),
    staleTime: 30_000,
    retry: 2,
  })

  return (
    <div className="space-y-4">
      <p className="text-sm uppercase tracking-[0.3em] text-cyan-400">
        {t('common.dashboard')}
      </p>
      <h2 className="text-3xl font-semibold">{t('common.homeRoute')}</h2>
      <p className="max-w-2xl text-zinc-400">
        {t('common.loadingGreeting')}
      </p>
      <div className="rounded-xl border border-zinc-800 bg-zinc-950/70 p-4">
        {isPending ? (
          <p className="text-sm text-zinc-400">{t('common.loadingGreeting')}</p>
        ) : (
          <p className="text-lg text-zinc-100">{data}</p>
        )}
      </div>
    </div>
  )
}
