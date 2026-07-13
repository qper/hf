import { useEffect, useMemo } from 'react'
import { useForm, useWatch } from 'react-hook-form'
import { z } from 'zod'
import { zodResolver } from '@hookform/resolvers/zod'
import { Check, Circle } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'

type HabitType = 'boolean' | 'numeric' | 'duration'
type HabitFrequency = 'daily' | 'custom' | 'weekly'

export type HabitFormValues = {
  name: string
  type: HabitType
  category: string
  color: string
  icon: string
  targetValue: string
  unit: string
  frequency: HabitFrequency
  selectedDays: string[]
  weeklyCount: string
}

const colorOptions = ['#14b8a6', '#f59e0b', '#8b5cf6', '#ef4444', '#22c55e', '#3b82f6', '#f97316', '#ec4899', '#64748b', '#e2e8f0']
const iconOptions = ['Sparkles', 'BookOpen', 'Dumbbell', 'Moon', 'Sun', 'Footprints', 'TreePine', 'Coffee', 'Target', 'Heart']
const dayOptions = ['Пн', 'Вт', 'Ср', 'Чт', 'Пт', 'Сб', 'Вс']

const habitSchema = (t: (key: string) => string) => z.object({
  name: z.string().trim().min(1, t('habits.requiredName')),
  type: z.enum(['boolean', 'numeric', 'duration'], {
    errorMap: () => ({ message: t('habits.requiredType') }),
  }),
  category: z.string().optional(),
  color: z.string().min(1, t('habits.requiredColor')),
  icon: z.string().min(1, t('habits.requiredIcon')),
  targetValue: z.string().optional(),
  unit: z.string().optional(),
  frequency: z.enum(['daily', 'custom', 'weekly']),
  selectedDays: z.array(z.string()).optional(),
  weeklyCount: z.string().optional(),
}).superRefine((data, ctx) => {
  if (data.type === 'numeric' && (!data.targetValue || Number(data.targetValue) <= 0)) {
    ctx.addIssue({
      code: z.ZodIssueCode.custom,
      path: ['targetValue'],
      message: t('habits.requiredTarget'),
    })
  }
  if (data.frequency === 'custom' && (!data.selectedDays || data.selectedDays.length === 0)) {
    ctx.addIssue({
      code: z.ZodIssueCode.custom,
      path: ['selectedDays'],
      message: t('habits.requiredDays'),
    })
  }
  if (data.frequency === 'weekly' && (!data.weeklyCount || Number(data.weeklyCount) <= 0)) {
    ctx.addIssue({
      code: z.ZodIssueCode.custom,
      path: ['weeklyCount'],
      message: t('habits.requiredWeeklyCount'),
    })
  }
})

type HabitFormDialogProps = {
  open: boolean
  onOpenChange: (open: boolean) => void
  defaultValues?: Partial<HabitFormValues>
  mode?: 'create' | 'edit'
  onSubmit?: (values: HabitFormValues) => void
}

export function HabitFormDialog({
  open,
  onOpenChange,
  defaultValues,
  mode = 'create',
  onSubmit: onSubmitHabit,
}: HabitFormDialogProps) {
  const { t } = useTranslation()
  const {
    register,
    handleSubmit,
    reset,
    setValue,
    control,
    formState: { errors, isSubmitting },
  } = useForm<HabitFormValues>({
    resolver: zodResolver(habitSchema(t)),
    defaultValues: {
      name: '',
      type: 'boolean',
      category: 'All',
      color: colorOptions[0],
      icon: iconOptions[0],
      targetValue: '',
      unit: mode === 'edit' ? 'мин' : '',
      frequency: 'daily',
      selectedDays: [],
      weeklyCount: '1',
      ...defaultValues,
    },
  })

  const selectedType = useWatch({ control, name: 'type' }) as HabitType | undefined
  const selectedFrequency = useWatch({ control, name: 'frequency' }) as HabitFrequency | undefined
  const selectedColor = useWatch({ control, name: 'color' }) as string | undefined
  const selectedIcon = useWatch({ control, name: 'icon' }) as string | undefined
  const selectedDays = (useWatch({ control, name: 'selectedDays' }) as string[] | undefined) ?? []

  const resolvedDefaultValues = useMemo(() => ({
    name: '',
    type: 'boolean' as HabitType,
    category: 'All',
    color: colorOptions[0],
    icon: iconOptions[0],
    targetValue: '',
    unit: '',
    frequency: 'daily' as HabitFrequency,
    selectedDays: [] as string[],
    weeklyCount: '1',
    ...defaultValues,
  }), [defaultValues])

  useEffect(() => {
    reset(resolvedDefaultValues)
  }, [open, reset, resolvedDefaultValues])

  const onSubmit = (values: HabitFormValues) => {
    const payload = {
      ...values,
      selectedDays,
      unit: values.type === 'duration' ? 'мин' : values.unit,
    }
    onSubmitHabit?.(payload)
    onOpenChange(false)
  }

  const toggleDay = (day: string) => {
    const current = selectedDays ?? []
    const next = current.includes(day) ? current.filter((item) => item !== day) : [...current, day]
    setValue('selectedDays', next, { shouldDirty: true, shouldTouch: true })
  }

  const frequencyHint =
    selectedFrequency === 'daily'
      ? t('habits.everyDay')
      : selectedFrequency === 'custom'
        ? t('habits.chooseDays')
        : t('habits.timesPerWeek')

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-2xl">
        <DialogHeader>
          <DialogTitle>{mode === 'edit' ? t('habits.editHabit') : t('habits.newHabit')}</DialogTitle>
          <DialogDescription>
            {t('habits.habitDescription')}
          </DialogDescription>
        </DialogHeader>

        <form onSubmit={handleSubmit(onSubmit)} className="space-y-5">
          <div className="space-y-2">
            <label className="text-sm font-medium text-zinc-200">{t('habits.name')}</label>
            <Input {...register('name')} placeholder={t('habits.namePlaceholder')} />
            {errors.name ? <p className="text-sm text-rose-400">{errors.name.message}</p> : null}
          </div>

          <div className="space-y-2">
            <label className="text-sm font-medium text-zinc-200">{t('habits.type')}</label>
            <div className="flex flex-wrap gap-3">
              {(['boolean', 'numeric', 'duration'] as HabitType[]).map((type) => (
                <label key={type} className="flex items-center gap-2 rounded-full border border-zinc-700 px-3 py-2 text-sm text-zinc-200">
                  <input
                    type="radio"
                    value={type}
                    disabled={mode === 'edit'}
                    {...register('type')}
                  />
                  {type === 'boolean' ? 'boolean' : type === 'numeric' ? 'numeric' : 'duration'}
                </label>
              ))}
            </div>
            {errors.type ? <p className="text-sm text-rose-400">{errors.type.message}</p> : null}
          </div>

          <div className="grid gap-4 md:grid-cols-2">
            <div className="space-y-2">
              <label className="text-sm font-medium text-zinc-200">{t('habits.category')}</label>
              <select className="flex h-10 w-full rounded-md border border-zinc-700 bg-zinc-900 px-3 py-2 text-sm text-zinc-100" {...register('category')}>
                <option value="All">{t('habits.categories.all')}</option>
                <option value="Health">{t('habits.categories.health')}</option>
                <option value="Work">{t('habits.categories.work')}</option>
                <option value="Learning">{t('habits.categories.learning')}</option>
              </select>
            </div>
            <div className="space-y-2">
              <label className="text-sm font-medium text-zinc-200">{t('habits.color')}</label>
              <div className="flex flex-wrap gap-2">
                {colorOptions.map((color) => (
                  <button
                    key={color}
                    type="button"
                    onClick={() => setValue('color', color)}
                    className={`h-8 w-8 rounded-full border-2 ${selectedColor === color ? 'border-white' : 'border-transparent'}`}
                    style={{ backgroundColor: color }}
                    aria-label={`color-${color}`}
                  />
                ))}
              </div>
            </div>
          </div>

          <div className="space-y-2">
            <label className="text-sm font-medium text-zinc-200">{t('habits.icon')}</label>
            <div className="grid grid-cols-5 gap-2">
              {iconOptions.map((icon) => (
                <button
                  key={icon}
                  type="button"
                  onClick={() => setValue('icon', icon)}
                  className={`flex h-10 items-center justify-center rounded-lg border ${selectedIcon === icon ? 'border-cyan-400 bg-zinc-800' : 'border-zinc-700 bg-zinc-900'}`}
                >
                  {icon}
                </button>
              ))}
            </div>
          </div>

          {selectedType === 'numeric' ? (
            <div className="grid gap-4 md:grid-cols-2">
              <div className="space-y-2">
                <label className="text-sm font-medium text-zinc-200">{t('habits.targetValue')}</label>
                <Input type="number" {...register('targetValue')} placeholder="20" />
                {errors.targetValue ? <p className="text-sm text-rose-400">{errors.targetValue.message}</p> : null}
              </div>
              <div className="space-y-2">
                <label className="text-sm font-medium text-zinc-200">{t('habits.unit')}</label>
                <Input {...register('unit')} placeholder="шт" />
              </div>
            </div>
          ) : null}

          <div className="space-y-2">
            <label className="text-sm font-medium text-zinc-200">{t('habits.frequency')}</label>
            <div className="flex flex-wrap gap-2">
              {([
                { value: 'daily', label: t('habits.daily') },
                { value: 'custom', label: t('habits.custom') },
                { value: 'weekly', label: t('habits.weekly') },
              ] as const).map((option) => (
                <label key={option.value} className="flex items-center gap-2 rounded-full border border-zinc-700 px-3 py-2 text-sm text-zinc-200">
                  <input type="radio" value={option.value} {...register('frequency')} />
                  {option.label}
                </label>
              ))}
            </div>
            <p className="text-sm text-zinc-400">{frequencyHint}</p>
          </div>

          {selectedFrequency === 'custom' ? (
            <div className="flex flex-wrap gap-2">
              {dayOptions.map((day) => (
                <button
                  key={day}
                  type="button"
                  onClick={() => toggleDay(day)}
                  className={`flex items-center gap-2 rounded-full border px-3 py-2 text-sm ${selectedDays.includes(day) ? 'border-cyan-400 bg-cyan-500/10 text-cyan-300' : 'border-zinc-700 bg-zinc-900 text-zinc-300'}`}
                >
                  {selectedDays.includes(day) ? <Check size={14} /> : <Circle size={14} />}
                  {day}
                </button>
              ))}
            </div>
          ) : null}

          {selectedFrequency === 'weekly' ? (
            <div className="space-y-2">
              <label className="text-sm font-medium text-zinc-200">{t('habits.weeklyCount')}</label>
              <Input type="number" {...register('weeklyCount')} min="1" />
              {errors.weeklyCount ? <p className="text-sm text-rose-400">{errors.weeklyCount.message}</p> : null}
            </div>
          ) : null}

          <div className="flex justify-end gap-2">
            <DialogClose asChild>
              <Button type="button" variant="outline">
                {t('habits.cancel')}
              </Button>
            </DialogClose>
            <Button type="submit" disabled={isSubmitting}>
              {t('habits.save')}
            </Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  )
}
