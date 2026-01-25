import { createFileRoute, Outlet, useParams } from '@tanstack/react-router'
import { useEffect } from 'react'
import { useTranslation } from 'react-i18next'
import i18n from '@/lib/i18n'

export const Route = createFileRoute('/$locale')({
  component: LocaleLayout,
  beforeLoad: ({ params }) => {
    
    if ((params as any).locale && i18n.language !== (params as any).locale) {
      i18n.changeLanguage((params as any).locale)
    }
  }
})

function LocaleLayout() {
  const { locale } = useParams({ from: '/$locale' })
  const { i18n } = useTranslation()

  useEffect(() => {
    if (locale && i18n.language !== locale) {
      i18n.changeLanguage(locale)
    }
  }, [locale, i18n])

  return <Outlet />
}
