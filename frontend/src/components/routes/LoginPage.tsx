import { useState } from 'react'
import { useForm } from 'react-hook-form'
import { z } from 'zod'
import { zodResolver } from '@hookform/resolvers/zod'
import { Link, useNavigate } from '@tanstack/react-router'
import { AlertTriangle } from 'lucide-react'
import { useTranslation } from 'react-i18next'

import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import * as auth from '@/api/auth'

const loginSchema = (t: (key: string) => string) =>
  z.object({
    username: z.string().min(1, t('auth.usernameRequired')),
    password: z.string().min(1, t('auth.passwordRequired')),
  })

type LoginFormValues = z.infer<ReturnType<typeof loginSchema>>

export function LoginPage() {
  const navigate = useNavigate()
  const { t } = useTranslation()
  const [errorMessage, setErrorMessage] = useState<string | null>(null)
  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm<LoginFormValues>({
    resolver: zodResolver(loginSchema(t)),
  })

  const onSubmit = async (values: LoginFormValues) => {
    setErrorMessage(null)
    try {
      await auth.login(values.username, values.password)
      const today = new Date().toISOString().split('T')[0]
      navigate({ to: '/board/$date', params: { date: today } })
    } catch (error) {
      if (error instanceof Error && error.message === 'session_expired') {
        setErrorMessage(t('errors.sessionExpired'))
        return
      }

      if (error instanceof Response) {
        const data = await error.json().catch(() => ({}))
        setErrorMessage(auth.getAuthErrorMessage(error.status, data, t))
        return
      }

      setErrorMessage(t('errors.network'))
    }
  }

  return (
    <div className="mx-auto max-w-md px-4 py-8 sm:px-6">
      <div className="rounded-[2rem] border border-zinc-800 bg-zinc-950/90 p-8 shadow-2xl shadow-black/30">
        <div className="mb-8 space-y-3">
          <p className="text-sm uppercase tracking-[0.3em] text-cyan-400">
            {t('common.access')}
          </p>
          <h2 className="text-3xl font-semibold">{t('common.signIn')}</h2>
          <p className="max-w-xl text-sm text-zinc-400">
            {t('common.loginSubtitle')}
          </p>
        </div>

        <form className="space-y-6" onSubmit={handleSubmit(onSubmit)}>
          <div className="space-y-2">
            <label
              htmlFor="username"
              className="block text-sm font-medium text-zinc-200"
            >
              {t('common.username')}
            </label>
            <Input
              id="username"
              {...register('username')}
              autoComplete="username"
            />
            {errors.username ? (
              <p className="text-sm text-rose-400">{errors.username.message}</p>
            ) : null}
          </div>

          <div className="space-y-2">
            <label
              htmlFor="password"
              className="block text-sm font-medium text-zinc-200"
            >
              {t('common.password')}
            </label>
            <Input
              id="password"
              type="password"
              {...register('password')}
              autoComplete="current-password"
            />
            {errors.password ? (
              <p className="text-sm text-rose-400">{errors.password.message}</p>
            ) : null}
          </div>

          {errorMessage ? (
            <div className="rounded-2xl border border-rose-500/20 bg-rose-500/10 p-4 text-sm text-rose-100">
              <div className="flex items-center gap-2">
                <AlertTriangle className="h-4 w-4" />
                <span>{errorMessage}</span>
              </div>
            </div>
          ) : null}

          <div className="space-y-4">
            <Button type="submit" className="w-full" disabled={isSubmitting}>
              {isSubmitting ? t('common.signingIn') : t('common.signIn')}
            </Button>
            <div className="flex flex-col gap-2 text-center text-sm text-zinc-400">
              <Link
                to="/register"
                className="text-cyan-300 underline-offset-4 hover:underline"
              >
                {t('common.loginLink')}
              </Link>
              <Link
                to="/register"
                className="text-cyan-300 underline-offset-4 hover:underline"
              >
                {t('common.recoveryLink')}
              </Link>
            </div>
          </div>
        </form>
      </div>
    </div>
  )
}
