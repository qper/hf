import { render, screen } from '@testing-library/react'
import { useTranslation } from 'react-i18next'
import { describe, expect, it } from 'vitest'
import './i18n'
import i18n from './i18n'

function LocaleProbe() {
  const { t } = useTranslation()
  return <div>{t('nav.home')}</div>
}

describe('i18n', () => {
  it('uses russian strings by default and switches languages without reload', async () => {
    render(<LocaleProbe />)

    expect(screen.getByText('Главная')).toBeTruthy()

    await i18n.changeLanguage('en')

    expect(screen.getByText('Home')).toBeTruthy()
  })
})
