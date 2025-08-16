import React, { useState, useEffect } from 'react'
import { Link, useLocation } from 'react-router-dom'
import { useAuthStore } from '@/store/authStore'
import { apiService } from '@/services/api'
import { Category } from '@/types'
import { 
  Home, 
  Package, 
  ShoppingBag, 
  User, 
  Settings, 
  ChevronLeft, 
  ChevronRight,
  Gamepad2,
  Monitor,
  Smartphone,
  Tv,
  Headphones,
  MessageSquare
} from 'lucide-react'

interface SidebarProps {
  onClose?: () => void;
}

const Sidebar: React.FC<SidebarProps> = ({ onClose }) => {
  const [isCollapsed, setIsCollapsed] = useState(false)
  const [categories, setCategories] = useState<Category[]>([])
  const [isLoadingCategories, setIsLoadingCategories] = useState(true)
  const { user, isAuthenticated } = useAuthStore()
  const location = useLocation()

  const isActive = (path: string) => location.pathname === path

  // Fetch categories from API
  useEffect(() => {
    const fetchCategories = async () => {
      try {
        const categoriesData = await apiService.getCategories()
        setCategories(categoriesData)
      } catch (error) {
        console.error('Failed to fetch categories:', error)
        // Fallback to default categories if API fails
        setCategories([
          { id: 1, name: 'Gaming', description: 'Gaming accounts and items' },
          { id: 2, name: 'Streaming', description: 'Streaming services' },
          { id: 3, name: 'Mobile', description: 'Mobile apps and services' },
          { id: 4, name: 'Software', description: 'Software licenses' },
          { id: 5, name: 'Music', description: 'Digital music and audio' }
        ])
      } finally {
        setIsLoadingCategories(false)
      }
    }

    fetchCategories()
  }, [])

  const handleLinkClick = () => {
    if (onClose) {
      onClose()
    }
  }

  const mainNav = [
    { name: 'Home', href: '/', icon: Home },
    { name: 'Digital Goods', href: '/products', icon: Package },
  ]

  const authenticatedNav = [
    { name: 'My Orders', href: '/orders', icon: ShoppingBag },
    { name: 'Conversations', href: '/conversations', icon: MessageSquare },
    { name: 'Profile', href: '/profile', icon: User },
  ]

  const sellerNav = [
    { name: 'Seller Dashboard', href: '/seller-dashboard', icon: Settings },
    { name: 'My Products', href: '/seller/products', icon: Package },
  ]

  const adminNav = [
    { name: 'Admin Dashboard', href: '/admin-dashboard', icon: Settings },
  ]

  const getCategoryIcon = (categoryName: string) => {
    const name = categoryName.toLowerCase()
    if (name.includes('gaming') || name.includes('game')) return Gamepad2
    if (name.includes('streaming') || name.includes('video')) return Monitor
    if (name.includes('mobile') || name.includes('app')) return Smartphone
    if (name.includes('software') || name.includes('license')) return Tv
    if (name.includes('music') || name.includes('audio')) return Headphones
    return Package
  }

  return (
    <div className={`bg-white border-r border-gray-200 transition-all duration-300 dark:bg-gray-900 dark:border-gray-800 ${
      isCollapsed ? 'w-16' : 'w-64'
    }`}>
      <div className="flex flex-col h-full">
        {/* Toggle Button - Only show on desktop */}
        <div className="hidden md:flex justify-end p-4 border-b border-gray-200 dark:border-gray-800">
          <button
            onClick={() => setIsCollapsed(!isCollapsed)}
            className="p-2 rounded-lg hover:bg-gray-100 transition-colors dark:hover:bg-gray-800"
          >
            {isCollapsed ? <ChevronRight className="w-4 h-4" /> : <ChevronLeft className="w-4 h-4" />}
          </button>
        </div>

        {/* Navigation */}
        <nav className="flex-1 p-4 space-y-6 overflow-y-auto sidebar-scroll">
          {/* Main Navigation */}
          <div>
            <h3 className={`text-xs font-semibold text-gray-500 uppercase tracking-wider mb-3 dark:text-gray-400 ${
              isCollapsed ? 'sr-only' : ''
            }`}>
              Main
            </h3>
            <ul className="space-y-2">
              {mainNav.map((item) => {
                const Icon = item.icon
                return (
                  <li key={item.name}>
                    <Link
                      to={item.href}
                      onClick={handleLinkClick}
                      className={`flex items-center px-3 py-2 rounded-lg text-sm font-medium transition-colors ${
                        isActive(item.href)
                          ? 'bg-blue-50 text-blue-700 dark:bg-blue-950 dark:text-blue-300'
                          : 'text-gray-700 hover:bg-gray-50 hover:text-gray-900 dark:text-gray-300 dark:hover:bg-gray-800 dark:hover:text-gray-100'
                      }`}
                    >
                      <Icon className="w-5 h-5 flex-shrink-0" />
                      {!isCollapsed && <span className="ml-3">{item.name}</span>}
                    </Link>
                  </li>
                )
              })}
            </ul>
          </div>

          {/* Digital Goods Categories */}
          <div>
            <h3 className={`text-xs font-semibold text-gray-500 uppercase tracking-wider mb-3 dark:text-gray-400 ${
              isCollapsed ? 'sr-only' : ''
            }`}>
              Categories
            </h3>
            {isLoadingCategories ? (
              <div className="space-y-2">
                {[...Array(5)].map((_, i) => (
                  <div key={i} className="animate-pulse">
                    <div className="h-10 bg-gray-200 rounded-lg"></div>
                  </div>
                ))}
              </div>
            ) : (
              <ul className="space-y-2">
                {categories.map((category) => {
                  const Icon = getCategoryIcon(category.name)
                  return (
                    <li key={category.id}>
                      <Link
                        to={`/products?category_id=${category.id}`}
                        onClick={handleLinkClick}
                        className="flex items-center px-3 py-2 rounded-lg text-sm font-medium text-gray-700 hover:bg-gray-50 hover:text-gray-900 transition-colors dark:text-gray-300 dark:hover:bg-gray-800 dark:hover:text-gray-100"
                      >
                        <Icon className="w-5 h-5 flex-shrink-0" />
                        {!isCollapsed && (
                          <>
                            <span className="ml-3 flex-1">{category.name}</span>
                          </>
                        )}
                      </Link>
                    </li>
                  )
                })}
              </ul>
            )}
          </div>

          {/* Authenticated Navigation */}
          {isAuthenticated && (
            <div>
              <h3 className={`text-xs font-semibold text-gray-500 uppercase tracking-wider mb-3 dark:text-gray-400 ${
                isCollapsed ? 'sr-only' : ''
              }`}>
                Account
              </h3>
              <ul className="space-y-2">
                {authenticatedNav.map((item) => {
                  const Icon = item.icon
                  return (
                    <li key={item.name}>
                      <Link
                        to={item.href}
                        onClick={handleLinkClick}
                        className={`flex items-center px-3 py-2 rounded-lg text-sm font-medium transition-colors ${
                          isActive(item.href)
                            ? 'bg-blue-50 text-blue-700 dark:bg-blue-950 dark:text-blue-300'
                            : 'text-gray-700 hover:bg-gray-50 hover:text-gray-900 dark:text-gray-300 dark:hover:bg-gray-800 dark:hover:text-gray-100'
                        }`}
                      >
                        <Icon className="w-5 h-5 flex-shrink-0" />
                        {!isCollapsed && <span className="ml-3">{item.name}</span>}
                      </Link>
                    </li>
                  )
                })}
              </ul>
            </div>
          )}

          {/* Seller Navigation */}
          {isAuthenticated && user?.roles.some(role => role.name === 'seller') && (
            <div>
              <h3 className={`text-xs font-semibold text-gray-500 uppercase tracking-wider mb-3 dark:text-gray-400 ${
                isCollapsed ? 'sr-only' : ''
              }`}>
                Seller
              </h3>
              <ul className="space-y-2">
                {sellerNav.map((item) => {
                  const Icon = item.icon
                  return (
                    <li key={item.name}>
                      <Link
                        to={item.href}
                        onClick={handleLinkClick}
                        className={`flex items-center px-3 py-2 rounded-lg text-sm font-medium transition-colors ${
                          isActive(item.href)
                            ? 'bg-blue-50 text-blue-700 dark:bg-blue-950 dark:text-blue-300'
                            : 'text-gray-700 hover:bg-gray-50 hover:text-gray-900 dark:text-gray-300 dark:hover:bg-gray-800 dark:hover:text-gray-100'
                        }`}
                      >
                        <Icon className="w-5 h-5 flex-shrink-0" />
                        {!isCollapsed && <span className="ml-3">{item.name}</span>}
                      </Link>
                    </li>
                  )
                })}
              </ul>
            </div>
          )}

          {/* Admin Navigation */}
          {isAuthenticated && user?.roles.some(role => role.name === 'admin') && (
            <div>
              <h3 className={`text-xs font-semibold text-gray-500 uppercase tracking-wider mb-3 dark:text-gray-400 ${
                isCollapsed ? 'sr-only' : ''
              }`}>
                Admin
              </h3>
              <ul className="space-y-2">
                {adminNav.map((item) => {
                  const Icon = item.icon
                  return (
                    <li key={item.name}>
                      <Link
                        to={item.href}
                        onClick={handleLinkClick}
                        className={`flex items-center px-3 py-2 rounded-lg text-sm font-medium transition-colors ${
                          isActive(item.href)
                            ? 'bg-blue-50 text-blue-700 dark:bg-blue-950 dark:text-blue-300'
                            : 'text-gray-700 hover:bg-gray-50 hover:text-gray-900 dark:text-gray-300 dark:hover:bg-gray-800 dark:hover:text-gray-100'
                        }`}
                      >
                        <Icon className="w-5 h-5 flex-shrink-0" />
                        {!isCollapsed && <span className="ml-3">{item.name}</span>}
                      </Link>
                    </li>
                  )
                })}
              </ul>
            </div>
          )}
        </nav>
      </div>
    </div>
  )
}

export default Sidebar

