import { Button } from '@/components/ui/button'
import { ChevronLeft, ChevronRight, Lock } from 'lucide-react'
import { useCallback } from 'react'
import { useTranslation } from 'react-i18next'

type DateNavBarProps = {
  date: string
  onDateChange: (date: string) => void
  progress?: { done: number; total: number }
  isEditable?: boolean
}

export function DateNavBar({ date, onDateChange, progress, isEditable = true }: DateNavBarProps) {
  const { i18n, t } = useTranslation()
  const today = new Date().toISOString().split('T')[0]
  const isToday = date === today

  const getLocalizedDate = useCallback((): string => {
    const dateObj = new Date(date + 'T00:00:00')
    const locale = i18n.language === 'ru' ? 'ru-RU' : 'en-US'

    // Get day name (Mon, Tue, etc) but we want abbreviated
    const dayName = dateObj.toLocaleDateString(locale, { weekday: 'short' })
    // Get month name (January, February, etc)
    const monthName = dateObj.toLocaleDateString(locale, { month: 'long' })
    const dayNum = dateObj.getDate()

    return `${dayName}, ${dayNum} ${monthName}`
  }, [date, i18n.language])

  const handlePrevDay = useCallback(() => {
    const d = new Date(date)
    d.setDate(d.getDate() - 1)
    onDateChange(d.toISOString().split('T')[0])
  }, [date, onDateChange])

  const handleNextDay = useCallback(() => {
    const d = new Date(date)
    d.setDate(d.getDate() + 1)
    onDateChange(d.toISOString().split('T')[0])
  }, [date, onDateChange])

  const handleToday = useCallback(() => {
    onDateChange(today)
  }, [today, onDateChange])

  const canGoForward = date < today

  const progressPercent =
    progress && progress.total > 0
      ? Math.round((progress.done / progress.total) * 100)
      : 0

  return (
    <div className="space-y-4">
      {/* Date Navigation */}
      <div className="flex items-center justify-between">
        <Button
          variant="ghost"
          size="sm"
          onClick={handlePrevDay}
          disabled={!isEditable}
          className="text-zinc-400 hover:text-zinc-100 disabled:opacity-50 disabled:cursor-not-allowed"
          title={!isEditable ? 'Вне окна редактирования' : ''}
        >
          <ChevronLeft className="h-4 w-4" />
        </Button>

        <div className="text-center">
          <div className="flex items-center justify-center gap-1">
            <p className="text-sm text-zinc-400">{getLocalizedDate()}</p>
            {!isEditable && (
              <div title="Вне окна редактирования" className="cursor-help">
                <Lock className="h-4 w-4 text-zinc-500" />
              </div>
            )}
          </div>
        </div>

        <Button
          variant="ghost"
          size="sm"
          onClick={handleNextDay}
          disabled={!isEditable || !canGoForward}
          className="text-zinc-400 hover:text-zinc-100 disabled:opacity-50 disabled:cursor-not-allowed"
          title={!isEditable ? 'Вне окна редактирования' : ''}
        >
          <ChevronRight className="h-4 w-4" />
        </Button>
      </div>

      {/* Today Button */}
      {!isToday && (
        <div className="flex justify-center">
          <Button
            variant="outline"
            size="sm"
            onClick={handleToday}
            disabled={!isEditable}
            className="text-xs disabled:opacity-50 disabled:cursor-not-allowed"
            title={!isEditable ? 'Вне окна редактирования' : ''}
          >
            {t('common.today')}
          </Button>
        </div>
      )}

      {/* Progress Bar */}
      {progress && (
        <div className="space-y-2">
          <div className="h-2 bg-zinc-800 rounded-full overflow-hidden">
            <div
              className="h-full bg-gradient-to-r from-cyan-500 to-blue-500 transition-all"
              style={{ width: `${progressPercent}%` }}
            />
          </div>
          <p className="text-xs text-center text-zinc-400">
            {progress.done} / {progress.total}
          </p>
        </div>
      )}
    </div>
  )
}

