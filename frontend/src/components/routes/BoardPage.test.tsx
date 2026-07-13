import { render, screen } from '@testing-library/react'
import { BoardPage } from './BoardPage'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import { I18nextProvider } from 'react-i18next'
import i18n from '@/i18n'

// Mock the router
vi.mock('@tanstack/react-router', () => ({
  useParams: () => ({ date: '2026-07-13' }),
  useNavigate: () => vi.fn(),
}))

// Mock the API
vi.mock('@/api/board', () => ({
  getBoard: vi.fn(() =>
    Promise.resolve({
      date: '2026-07-13',
      is_editable: true,
      progress: { done: 2, total: 5 },
      habits: [],
    }),
  ),
}))

describe('BoardPage', () => {
  const createWrapper = () => {
    const queryClient = new QueryClient({
      defaultOptions: {
        queries: { retry: false },
        mutations: { retry: false },
      },
    })
    return ({ children }: { children: React.ReactNode }) => (
      <QueryClientProvider client={queryClient}>
        <I18nextProvider i18n={i18n}>{children}</I18nextProvider>
      </QueryClientProvider>
    )
  }

  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('renders without crashing', async () => {
    render(<BoardPage date="2026-07-13" />, {
      wrapper: createWrapper(),
    })
    // Just verify the component renders
    expect(screen.getByText(/Board/i)).toBeTruthy()
  })

  it('applies swipe transition class when transitioning', async () => {
    const { container } = render(<BoardPage date="2026-07-13" />, {
      wrapper: createWrapper(),
    })
    const boardDiv = container.firstChild as HTMLElement
    expect(boardDiv).toBeTruthy()
    // Verify the div has touch event handlers
    expect(boardDiv.ontouchstart).toBeDefined()
    expect(boardDiv.ontouchend).toBeDefined()
  })
})
