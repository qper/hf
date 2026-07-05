/* eslint-disable react-refresh/only-export-components */
import React, { createContext, useContext, useEffect, useState } from 'react'
import * as auth from '@/api/auth'

type UserShape = Record<string, unknown>

type AuthContextValue = {
  user: UserShape | null
  isLoading: boolean
  login: (username: string, password: string) => Promise<void>
  logout: () => Promise<void>
}

const AuthContext = createContext<AuthContextValue | undefined>(undefined)

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<UserShape | null>(null)
  const [isLoading, setIsLoading] = useState(true)

  useEffect(() => {
    let mounted = true

    ;(async () => {
      setIsLoading(true)
      const ok = await auth.tryRefresh()
      if (!mounted) return

      if (ok) {
        try {
          const res = await auth.authFetch('/api/v1/me')
          if (res.ok) {
            const data = await res.json()
            setUser(data)
          } else {
            setUser(null)
          }
        } catch {
          setUser(null)
        }
      } else {
        setUser(null)
      }

      if (mounted) setIsLoading(false)
    })()

    return () => {
      mounted = false
    }
  }, [])

  const login = async (username: string, password: string) => {
    setIsLoading(true)
    try {
      await auth.login(username, password)
      const res = await auth.authFetch('/api/v1/me')
      if (res.ok) {
        const data = await res.json()
        setUser(data)
      }
    } finally {
      setIsLoading(false)
    }
  }

  const logout = async () => {
    await auth.logout()
    setUser(null)
  }

  return (
    <AuthContext.Provider value={{ user, isLoading, login, logout }}>
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth() {
  const ctx = useContext(AuthContext)
  if (!ctx) throw new Error('useAuth must be used within AuthProvider')
  return ctx
}
