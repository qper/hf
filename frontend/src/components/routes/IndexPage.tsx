import { useRouter } from '@tanstack/react-router'
import { useEffect } from 'react'

export function IndexPage() {
  const router = useRouter()

  useEffect(() => {
    const today = new Date().toISOString().split('T')[0]
    router.navigate({
      to: '/board/$date',
      params: { date: today },
    }).catch(() => {
      // Silently catch any errors during redirect
    })
  }, [router])

  return null
}
