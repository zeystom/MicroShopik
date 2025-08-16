import { useEffect, useState, useCallback } from 'react'

type Theme = 'light' | 'dark'

const STORAGE_KEY = 'microshopik_theme'

export function useTheme() {
  const [theme, setThemeState] = useState<Theme>(() => {
    if (typeof window === 'undefined') return 'light'
    const stored = localStorage.getItem(STORAGE_KEY) as Theme | null
    if (stored === 'light' || stored === 'dark') return stored
    const prefersDark = window.matchMedia?.('(prefers-color-scheme: dark)').matches
    return prefersDark ? 'dark' : 'light'
  })

  const applyThemeClass = useCallback((next: Theme) => {
    const root = document.documentElement
    if (next === 'dark') root.classList.add('dark')
    else root.classList.remove('dark')
  }, [])

  useEffect(() => {
    try {
      localStorage.setItem(STORAGE_KEY, theme)
    } catch {}
    applyThemeClass(theme)
  }, [theme, applyThemeClass])

  const setTheme = useCallback((next: Theme) => setThemeState(next), [])

  const toggleTheme = useCallback(() => {
    // View Transitions API radial color change
    const next: Theme = theme === 'dark' ? 'light' : 'dark'
    // If unsupported, just toggle
    // @ts-ignore
    if (!document.startViewTransition) {
      setTheme(next)
      return
    }
    const x = (window.innerWidth / 2)
    const y = (40) // roughly header area; we keep it simple
    // @ts-ignore
    document.startViewTransition(() => setTheme(next)).ready.then(() => {
      const endRadius = Math.hypot(Math.max(x, window.innerWidth - x), Math.max(y, window.innerHeight - y))
      document.documentElement.animate(
        {
          clipPath: [
            `circle(0px at ${x}px ${y}px)`,
            `circle(${endRadius}px at ${x}px ${y}px)`
          ]
        },
        {
          duration: 400,
          easing: 'ease-out',
          pseudoElement: '::view-transition-new(root)'
        }
      )
    })
  }, [theme, setTheme])

  return { theme, setTheme, toggleTheme }
}


