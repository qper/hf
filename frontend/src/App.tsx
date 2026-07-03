import '@/App.css'
import { formatMessage } from '@/utils/format'

function App() {
  return (
    <main className="app-shell">
      <h1>HabitFlow Frontend</h1>
      <p>{formatMessage('Vite + React + TypeScript is ready.')}</p>
    </main>
  )
}

export default App
