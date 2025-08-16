import React, { useEffect, useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { useAuthStore } from '@/store/authStore'
import { ShoppingBag, User, LogOut, Menu, X, MessageCircle, Home, Package } from 'lucide-react'
import ThemeToggle from '@/components/ui/ThemeToggle'
import { Category } from '@/types'
import { apiService } from '@/services/api'

const Header: React.FC = () => {
  const { user, isAuthenticated, logout } = useAuthStore()
  const navigate = useNavigate()
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false)
  const [categories, setCategories] = useState<Category[]>([])
  const [isLoadingCategories, setIsLoadingCategories] = useState(false)

  const handleLogout = () => {
    logout()
    navigate('/login')
  }

  const toggleMobileMenu = () => {
    setIsMobileMenuOpen(!isMobileMenuOpen)
  }

  useEffect(() => {
    if (!isMobileMenuOpen) return
    let isMounted = true
    const fetchCategories = async () => {
      try {
        setIsLoadingCategories(true)
        const data = await apiService.getCategories()
        if (isMounted) setCategories(data)
      } catch (e) {
        // ignore
      } finally {
        if (isMounted) setIsLoadingCategories(false)
      }
    }
    fetchCategories()
    return () => { isMounted = false }
  }, [isMobileMenuOpen])

  return (
    <header className="bg-white shadow-sm border-b border-gray-200 dark:bg-gray-950 dark:border-gray-900">
      <div className="max-w-7xl mx-auto px-3 sm:px-6 lg:px-8">
        <div className="flex justify-between items-center h-14 sm:h-16">
          {/* Logo */}
          <Link to="/" className="flex items-center space-x-2">
            <div className="w-8 h-8 bg-gradient-to-r from-blue-600 to-purple-600 rounded-lg flex items-center justify-center">
              <span className="text-white font-bold text-sm">MS</span>
            </div>
            <span className="text-xl font-bold text-gray-900 dark:text-gray-100">MicroShopik</span>
          </Link>

          {/* Desktop Navigation */}
          <nav className="hidden md:flex space-x-8">
            <Link
              to="/"
              className="text-gray-700 hover:text-blue-600 px-3 py-2 rounded-md text-sm font-medium transition-colors dark:text-gray-300 dark:hover:text-blue-400"
            >
              Home
            </Link>
            <Link
              to="/products"
              className="text-gray-700 hover:text-blue-600 px-3 py-2 rounded-md text-sm font-medium transition-colors dark:text-gray-300 dark:hover:text-blue-400"
            >
              Digital Goods
            </Link>
            {user?.roles.some(role => role.name === 'admin') && (
              <Link
                to="/api-test"
                className="text-gray-700 hover:text-blue-600 px-3 py-2 rounded-md text-sm font-medium transition-colors dark:text-gray-300 dark:hover:text-blue-400"
              >
                API Test
              </Link>
            )}
            {isAuthenticated && (
              <>
                <Link
                  to="/orders"
                  className="text-gray-700 hover:text-blue-600 px-3 py-2 rounded-md text-sm font-medium transition-colors dark:text-gray-300 dark:hover:text-blue-400"
                >
                  My Orders
                </Link>
                <Link
                  to="/conversations"
                  className="text-gray-700 hover:text-blue-600 px-3 py-2 rounded-md text-sm font-medium transition-colors dark:text-gray-300 dark:hover:text-blue-400"
                >
                  Messages
                </Link>
                {user?.roles.some(role => role.name === 'seller') && (
                  <>
                    <Link
                      to="/seller-dashboard"
                      className="text-gray-700 hover:text-blue-600 px-3 py-2 rounded-md text-sm font-medium transition-colors dark:text-gray-300 dark:hover:text-blue-400"
                    >
                      Seller Dashboard
                    </Link>
                  </>
                )}
                {user?.roles.some(role => role.name === 'admin') && (
                  <Link
                    to="/admin-dashboard"
                    className="text-gray-700 hover:text-blue-600 px-3 py-2 rounded-md text-sm font-medium transition-colors dark:text-gray-300 dark:hover:text-blue-400"
                  >
                    Admin Dashboard
                  </Link>
                )}
              </>
            )}
          </nav>

          {/* User Menu */}
          <div className="flex items-center space-x-4">
            <div className="hidden md:block"><ThemeToggle /></div>
            {isAuthenticated ? (
              <div className="hidden md:flex items-center space-x-3">
                <Link
                  to="/profile"
                  className="flex items-center space-x-2 text-gray-700 hover:text-blue-600 transition-colors dark:text-gray-300 dark:hover:text-blue-400"
                >
                  <User className="w-5 h-5" />
                  <span className="text-sm font-medium">{user?.username}</span>
                </Link>
                <button
                  onClick={handleLogout}
                  className="flex items-center space-x-2 text-gray-700 hover:text-red-600 transition-colors dark:text-gray-300 dark:hover:text-red-400"
                >
                  <LogOut className="w-5 h-5" />
                  <span className="text-sm font-medium">Logout</span>
                </button>
              </div>
            ) : (
              <div className="hidden md:flex items-center space-x-4">
                <Link
                  to="/login"
                  className="text-gray-700 hover:text-blue-600 px-3 py-2 rounded-md text-sm font-medium transition-colors dark:text-gray-300 dark:hover:text-blue-400"
                >
                  Login
                </Link>
                <Link
                  to="/register"
                  className="bg-blue-600 text-white px-4 py-2 rounded-md text-sm font-medium hover:bg-blue-700 transition-colors dark:bg-blue-500 dark:hover:bg-blue-600"
                >
                  Register
                </Link>
              </div>
            )}

            {/* Mobile menu button */}
            <button
              onClick={toggleMobileMenu}
              className="md:hidden p-2 rounded-md text-gray-700 hover:text-blue-600 hover:bg-gray-100 transition-colors dark:text-gray-300 dark:hover:text-blue-400 dark:hover:bg-gray-900"
            >
              {isMobileMenuOpen ? (
                <X className="w-6 h-6" />
              ) : (
                <Menu className="w-6 h-6" />
              )}
            </button>
          </div>
        </div>

        {/* Mobile menu */}
        {isMobileMenuOpen && (
          <div className="md:hidden border-t border-gray-200 py-4 dark:border-gray-900">
            <div className="space-y-2">
              <div className="px-3 pb-2"><ThemeToggle /></div>
              {/* Main navigation */}
              <div className="px-3 pt-2">
                <div className="text-xs font-semibold text-gray-500 uppercase tracking-wider mb-2 dark:text-gray-400">Main</div>
                <Link
                  to="/"
                  className="block px-3 py-2 text-gray-700 hover:text-blue-600 hover:bg-gray-50 rounded-md text-sm font-medium transition-colors dark:text-gray-300 dark:hover:text-blue-400 dark:hover:bg-gray-900"
                  onClick={() => setIsMobileMenuOpen(false)}
                >
                  <div className="flex items-center space-x-2">
                    <Home className="w-5 h-5" />
                    <span>Home</span>
                  </div>
                </Link>
                <Link
                  to="/products"
                  className="mt-1 block px-3 py-2 text-gray-700 hover:text-blue-600 hover:bg-gray-50 rounded-md text-sm font-medium transition-colors dark:text-gray-300 dark:hover:text-blue-400 dark:hover:bg-gray-900"
                  onClick={() => setIsMobileMenuOpen(false)}
                >
                  <div className="flex items-center space-x-2">
                    <Package className="w-5 h-5" />
                    <span>Digital Goods</span>
                  </div>
                </Link>
              </div>
              {/* Categories */}
              <div className="px-3 pt-4">
                <div className="text-xs font-semibold text-gray-500 uppercase tracking-wider mb-2 dark:text-gray-400">Categories</div>
                {isLoadingCategories ? (
                  <div className="space-y-2">
                    {[...Array(4)].map((_, i) => (
                      <div key={i} className="h-9 bg-gray-100 rounded-md animate-pulse dark:bg-gray-800" />
                    ))}
                  </div>
                ) : (
                  <div className="space-y-1">
                    {categories.map((cat) => (
                      <Link
                        key={cat.id}
                        to={`/products?category_id=${cat.id}`}
                        className="block px-3 py-2 text-gray-700 hover:text-blue-600 hover:bg-gray-50 rounded-md text-sm font-medium transition-colors dark:text-gray-300 dark:hover:text-blue-400 dark:hover:bg-gray-900"
                        onClick={() => setIsMobileMenuOpen(false)}
                      >
                        {cat.name}
                      </Link>
                    ))}
                  </div>
                )}
              </div>
              {isAuthenticated ? (
                <>
                  <Link
                    to="/profile"
                    className="block px-3 py-2 text-gray-700 hover:text-blue-600 hover:bg-gray-50 rounded-md text-sm font-medium transition-colors dark:text-gray-300 dark:hover:text-blue-400 dark:hover:bg-gray-900"
                    onClick={() => setIsMobileMenuOpen(false)}
                  >
                    <div className="flex items-center space-x-2">
                      <User className="w-5 h-5" />
                      <span>Profile</span>
                    </div>
                  </Link>
                  <Link
                    to="/orders"
                    className="block px-3 py-2 text-gray-700 hover:text-blue-600 hover:bg-gray-50 rounded-md text-sm font-medium transition-colors dark:text-gray-300 dark:hover:text-blue-400 dark:hover:bg-gray-900"
                    onClick={() => setIsMobileMenuOpen(false)}
                  >
                    <div className="flex items-center space-x-2">
                      <ShoppingBag className="h-5 w-5" />
                      <span>My Orders</span>
                    </div>
                  </Link>
                  <Link
                    to="/conversations"
                    className="block px-3 py-2 text-gray-700 hover:text-blue-600 hover:bg-gray-50 rounded-md text-sm font-medium transition-colors dark:text-gray-300 dark:hover:text-blue-400 dark:hover:bg-gray-900"
                    onClick={() => setIsMobileMenuOpen(false)}
                  >
                    <div className="flex items-center space-x-2">
                      <MessageCircle className="h-5 w-5" />
                      <span>Messages</span>
                    </div>
                  </Link>
                  {user?.roles.some(role => role.name === 'seller') && (
                    <Link
                      to="/seller-dashboard"
                      className="block px-3 py-2 text-gray-700 hover:text-blue-600 hover:bg-gray-50 rounded-md text-sm font-medium transition-colors dark:text-gray-300 dark:hover:text-blue-400 dark:hover:bg-gray-900"
                      onClick={() => setIsMobileMenuOpen(false)}
                    >
                      Seller Dashboard
                    </Link>
                  )}
                  {user?.roles.some(role => role.name === 'seller') && (
                    <Link
                      to="/seller/products"
                      className="block px-3 py-2 text-gray-700 hover:text-blue-600 hover:bg-gray-50 rounded-md text-sm font-medium transition-colors dark:text-gray-300 dark:hover:text-blue-400 dark:hover:bg-gray-900"
                      onClick={() => setIsMobileMenuOpen(false)}
                    >
                      My Products
                    </Link>
                  )}
                  {user?.roles.some(role => role.name === 'admin') && (
                    <Link
                      to="/admin-dashboard"
                      className="block px-3 py-2 text-gray-700 hover:text-blue-600 hover:bg-gray-50 rounded-md text-sm font-medium transition-colors dark:text-gray-300 dark:hover:text-blue-400 dark:hover:bg-gray-900"
                      onClick={() => setIsMobileMenuOpen(false)}
                    >
                      Admin Dashboard
                    </Link>
                  )}
                  {user?.roles.some(role => role.name === 'admin') && (
                    <Link
                      to="/api-test"
                      className="block px-3 py-2 text-gray-700 hover:text-blue-600 hover:bg-gray-50 rounded-md text-sm font-medium transition-colors dark:text-gray-300 dark:hover:text-blue-400 dark:hover:bg-gray-900"
                      onClick={() => setIsMobileMenuOpen(false)}
                    >
                      API Test
                    </Link>
                  )}
                  <button
                    onClick={() => {
                      handleLogout()
                      setIsMobileMenuOpen(false)
                    }}
                    className="w-full text-left px-3 py-2 text-gray-700 hover:text-red-600 hover:bg-gray-50 rounded-md text-sm font-medium transition-colors dark:text-gray-300 dark:hover:text-red-400 dark:hover:bg-gray-900"
                  >
                    <div className="flex items-center space-x-2">
                      <LogOut className="w-5 h-5" />
                      <span>Logout</span>
                    </div>
                  </button>
                </>
              ) : (
                <>
                  <Link
                    to="/login"
                    className="block px-3 py-2 text-gray-700 hover:text-blue-600 hover:bg-gray-50 rounded-md text-sm font-medium transition-colors dark:text-gray-300 dark:hover:text-blue-400 dark:hover:bg-gray-900"
                    onClick={() => setIsMobileMenuOpen(false)}
                  >
                    Login
                  </Link>
                  <Link
                    to="/register"
                    className="block px-3 py-2 bg-blue-600 text-white rounded-md text-sm font-medium hover:bg-blue-700 transition-colors dark:bg-blue-500 dark:hover:bg-blue-600"
                    onClick={() => setIsMobileMenuOpen(false)}
                  >
                    Register
                  </Link>
                </>
              )}
            </div>
          </div>
        )}
      </div>
    </header>
  )
}

export default Header

