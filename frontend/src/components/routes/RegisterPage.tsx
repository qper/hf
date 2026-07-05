import { useState } from 'react'
import { useForm } from 'react-hook-form'
import { z } from 'zod'
import { zodResolver } from '@hookform/resolvers/zod'
import * as Dialog from '@radix-ui/react-dialog'
import { Download, FilePlus, Shield } from 'lucide-react'

import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import {
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'

const registerSchema = z
  .object({
    username: z.string().min(3, 'Минимум 3 символа'),
    password: z.string().min(8, 'Минимум 8 символов'),
    confirmPassword: z.string().min(1, 'Подтвердите пароль'),
  })
  .refine((data) => data.password === data.confirmPassword, {
    path: ['confirmPassword'],
    message: 'Пароли не совпадают',
  })

type RegisterFormValues = z.infer<typeof registerSchema>

const initialRecoveryCodes = Array.from({ length: 8 }, (_, index) => `CODE-${index + 1}`)

export function RegisterPage() {
  const [isDialogOpen, setIsDialogOpen] = useState(false)
  const [recoveryCodes, setRecoveryCodes] = useState<string[]>([])
  const [accepted, setAccepted] = useState(false)

  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm<RegisterFormValues>({
    resolver: zodResolver(registerSchema),
  })

  const onSubmit = async (values: RegisterFormValues) => {
    setAccepted(false)
    try {
      const response = await fetch('/api/v1/auth/register', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          username: values.username,
          password: values.password,
          email: `${values.username}@example.com`,
        }),
      })

      if (!response.ok) {
        return
      }

      const data = await response.json()
      setRecoveryCodes(data.recovery_codes || initialRecoveryCodes)
      setIsDialogOpen(true)
    } catch {
      setRecoveryCodes(initialRecoveryCodes)
      setIsDialogOpen(true)
    }
  }

  const copyAll = async () => {
    await navigator.clipboard.writeText(recoveryCodes.join('\n'))
  }

  const downloadCodes = () => {
    const blob = new Blob([recoveryCodes.join('\n')], { type: 'text/plain;charset=utf-8' })
    const link = document.createElement('a')
    link.href = URL.createObjectURL(blob)
    link.download = 'recovery-codes.txt'
    link.click()
    URL.revokeObjectURL(link.href)
  }

  return (
    <div className="mx-auto max-w-md px-4 py-8 sm:px-6">
      <div className="rounded-[2rem] border border-zinc-800 bg-zinc-950/90 p-8 shadow-2xl shadow-black/30">
        <div className="mb-8 space-y-3">
          <p className="text-sm uppercase tracking-[0.3em] text-cyan-400">Access</p>
          <h2 className="text-3xl font-semibold">Register</h2>
          <p className="max-w-xl text-sm text-zinc-400">
            Создайте аккаунт и сохраните recovery codes для безопасного доступа.
          </p>
        </div>

        <form className="space-y-6" onSubmit={handleSubmit(onSubmit)}>
          <div className="space-y-2">
            <label htmlFor="username" className="block text-sm font-medium text-zinc-200">
              Username
            </label>
            <Input id="username" {...register('username')} autoComplete="username" />
            {errors.username ? (
              <p className="text-sm text-rose-400">{errors.username.message}</p>
            ) : null}
          </div>

          <div className="space-y-2">
            <label htmlFor="password" className="block text-sm font-medium text-zinc-200">
              Password
            </label>
            <Input
              id="password"
              type="password"
              {...register('password')}
              autoComplete="new-password"
            />
            {errors.password ? (
              <p className="text-sm text-rose-400">{errors.password.message}</p>
            ) : null}
          </div>

          <div className="space-y-2">
            <label htmlFor="confirmPassword" className="block text-sm font-medium text-zinc-200">
              Confirm password
            </label>
            <Input
              id="confirmPassword"
              type="password"
              {...register('confirmPassword')}
              autoComplete="new-password"
            />
            {errors.confirmPassword ? (
              <p className="text-sm text-rose-400">{errors.confirmPassword.message}</p>
            ) : null}
          </div>

          <Button type="submit" className="w-full" disabled={isSubmitting}>
            {isSubmitting ? 'Создание...' : 'Зарегистрироваться'}
          </Button>
        </form>
      </div>

      <Dialog.Root open={isDialogOpen} onOpenChange={setIsDialogOpen}>
        <Dialog.Portal>
          <DialogContent className="max-w-lg">
            <DialogHeader>
              <DialogTitle>Сохраните recovery codes</DialogTitle>
              <DialogDescription>
                Они нужны для восстановления доступа. Скопируйте или скачайте кодовый файл.
              </DialogDescription>
            </DialogHeader>

            <div className="space-y-4">
              <div className="rounded-2xl border border-zinc-800 bg-zinc-950/90 p-4 text-sm text-zinc-100">
                <div className="mb-4 flex items-center gap-2 text-cyan-300">
                  <FilePlus className="h-4 w-4" />
                  <span>Восстановите доступ с помощью одного из этих кодов.</span>
                </div>
                <div className="grid gap-2 sm:grid-cols-2">
                  {recoveryCodes.map((code) => (
                    <div
                      key={code}
                      className="rounded-2xl border border-zinc-800 bg-zinc-900/80 p-3 text-xs uppercase tracking-[0.2em] text-zinc-100"
                    >
                      {code}
                    </div>
                  ))}
                </div>
              </div>

              <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
                <Button type="button" variant="outline" onClick={copyAll} className="w-full sm:w-auto">
                  <Shield className="mr-2 h-4 w-4" /> Копировать все
                </Button>
                <Button type="button" variant="outline" onClick={downloadCodes} className="w-full sm:w-auto">
                  <Download className="mr-2 h-4 w-4" /> Скачать .txt
                </Button>
              </div>

              <label className="flex items-center gap-3 rounded-2xl border border-zinc-800 bg-zinc-950/80 p-4 text-sm text-zinc-200">
                <input
                  type="checkbox"
                  checked={accepted}
                  onChange={(event) => setAccepted(event.target.checked)}
                  className="h-5 w-5 rounded border-zinc-600 bg-zinc-900 text-cyan-400 focus:ring-cyan-400"
                />
                <span>Я сохранил коды и понимаю, что они нужны для восстановления.</span>
              </label>

              <div className="flex flex-col gap-3 sm:flex-row sm:justify-end">
                <DialogClose asChild>
                  <Button type="button" variant="outline" className="w-full sm:w-auto">
                    Закрыть
                  </Button>
                </DialogClose>
                <Button type="button" disabled={!accepted} className="w-full sm:w-auto">
                  Продолжить
                </Button>
              </div>
            </div>
          </DialogContent>
        </Dialog.Portal>
      </Dialog.Root>
    </div>
  )
}
