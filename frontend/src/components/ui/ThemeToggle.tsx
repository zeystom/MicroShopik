import React, { useState } from 'react'
import { Sun, Moon } from 'lucide-react'
import { useTheme } from '@/hooks'

const ThemeToggle: React.FC = () => {
  const { theme, toggleTheme } = useTheme()
  const [burstId, setBurstId] = useState(0)

  const handleClick = () => {
    setBurstId(Date.now())
    toggleTheme()
  }

  const isDark = theme === 'dark'

  return (
    <button
      onClick={handleClick}
      aria-label="Toggle theme"
      aria-pressed={isDark}
      className={`theme-toggle inline-flex items-center justify-center rounded-md p-2 transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500 focus-visible:ring-offset-2 hover:bg-gray-100 text-gray-700 dark:hover:bg-gray-800 dark:text-gray-300 ${isDark ? 'is-dark' : ''}`}
    >
      <span className="relative inline-block h-5 w-5">
        <Sun className="icon sun-icon absolute inset-0 h-5 w-5" />
        <Moon className="icon moon-icon absolute inset-0 h-5 w-5" />
      </span>

      {/* Sparkle burst */}
      <span key={burstId} className="relative">
        <span className="sparkle s1" />
        <span className="sparkle s2" />
        <span className="sparkle s3" />
        <span className="sparkle s4" />
        <span className="sparkle s5" />
        <span className="sparkle s6" />
      </span>
    </button>
  )
}

export default ThemeToggle


