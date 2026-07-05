/* eslint-disable react-refresh/only-export-components */
import { createFileRoute } from '@tanstack/react-router'
import { useState } from 'react'
import * as auth from '@/api/auth'

function TestAuthPage() {
  const [status, setStatus] = useState<string>('')

  const doLogin = async () => {
    try {
      await auth.login('alice', 'password')
      setStatus('logged-in')
    } catch {
      setStatus('login-failed')
    }
  }

  const callProtected = async () => {
    try {
      const res = await auth.authFetch('/api/v1/protected')
      setStatus(String(res.status))
    } catch {
      setStatus('error')
    }
  }

  return (
    <div>
      <button id="login-btn" onClick={doLogin}>Login</button>
      <button id="call-btn" onClick={callProtected}>Call protected</button>
      <div id="status">{status}</div>
    </div>
  )
}

export const Route = createFileRoute('/test-auth')({
  component: TestAuthPage,
})
