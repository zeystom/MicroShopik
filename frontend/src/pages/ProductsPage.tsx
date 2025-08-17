import React, { useState, useEffect } from 'react'
import { Link, useSearchParams } from 'react-router-dom'
import { 
  Search, 
  Monitor,
  SortAsc,
  SortDesc,
  Loader2
} from 'lucide-react'
import { apiService } from '@/services/api'
import { Product, Category } from '@/types'
import { toast } from 'react-hot-toast'
import { useAuthStore } from '@/store/authStore'
import ProductStatusBadge from '@/components/ui/ProductStatusBadge'

const ProductsPage: React.FC = () => {
  const [searchParams, setSearchParams] = useSearchParams()
  const [searchTerm, setSearchTerm] = useState('')
  const [selectedCategory, setSelectedCategory] = useState('all')
  const [sortBy, setSortBy] = useState('relevance')
  const [sortOrder, setSortOrder] = useState<'asc' | 'desc'>('desc')
  const [products, setProducts] = useState<Product[]>([])
  const [categories, setCategories] = useState<Category[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState('')
  const { user } = useAuthStore()
  const isAdmin = user?.roles?.some(r => r.name === 'admin')

  // Initialize filters from URL params
  useEffect(() => {
    const categoryId = searchParams.get('category_id')
    const search = searchParams.get('search')
    
    if (categoryId) {
      setSelectedCategory(categoryId)
    }
    if (search) {
      setSearchTerm(search)
    }
  }, [searchParams])

  // Fetch data on component mount
  useEffect(() => {
    const fetchData = async () => {
      try {
        setIsLoading(true)
        // Админы видят все продукты, обычные пользователи только активные
        const productsPromise = isAdmin 
          ? apiService.getAllProducts()
          : apiService.getProducts()
        
        const [productsData, categoriesData] = await Promise.all([
          productsPromise,
          apiService.getCategories()
        ])
        setProducts(productsData)
        setCategories(categoriesData)
      } catch (err: unknown) {
        const message =
          typeof err === 'object' && err !== null && 'response' in err
            ? (err as { response?: { data?: { error?: string } } }).response?.data?.error || 'Failed to load products'
            : 'Failed to load products'
        setError(message)
        toast.error('Failed to load products')
      } finally {
        setIsLoading(false)
      }
    }

    fetchData()
  }, [isAdmin])

  // Update URL when filters change
  const updateURL = (category: string, search: string) => {
    const params = new URLSearchParams()
    if (category !== 'all') {
      params.set('category_id', category)
    }
    if (search) {
      params.set('search', search)
    }
    setSearchParams(params)
  }

  const handleCategoryChange = (category: string) => {
    setSelectedCategory(category)
    updateURL(category, searchTerm)
  }

  const handleSearchChange = (search: string) => {
    setSearchTerm(search)
    updateURL(selectedCategory, search)
  }

  const categoryOptions = [
    { id: 'all', name: 'All Categories', icon: Monitor },
    ...categories.map(cat => ({
      id: cat.id.toString(),
      name: cat.name,
      icon: Monitor // Default icon for all categories
    }))
  ]

  const sortOptions = [
    { value: 'relevance', label: 'Relevance' },
    { value: 'price', label: 'Price' },
    { value: 'rating', label: 'Rating' },
    { value: 'sales', label: 'Sales' },
    { value: 'newest', label: 'Newest' }
  ]

  const filteredProducts = products
    .filter(product => {
      // Обычные пользователи видят только активные продукты
      if (!isAdmin && !product.is_active) {
        return false
      }
      
      return (selectedCategory === 'all' || product.category_id.toString() === selectedCategory) &&
        (searchTerm === '' || 
          product.title.toLowerCase().includes(searchTerm.toLowerCase()) ||
          product.description.toLowerCase().includes(searchTerm.toLowerCase()))
    })
    .sort((a, b) => {
      switch (sortBy) {
        case 'price':
          return sortOrder === 'asc' ? a.price - b.price : b.price - a.price
        case 'sales':
          return sortOrder === 'asc' ? a.sold_count - b.sold_count : b.sold_count - a.sold_count
        case 'newest':
          return sortOrder === 'asc' ? 
            new Date(a.created_at).getTime() - new Date(b.created_at).getTime() : 
            new Date(b.created_at).getTime() - new Date(a.created_at).getTime()
        default:
          return 0
      }
    })

  const toggleSortOrder = () => {
    setSortOrder(sortOrder === 'asc' ? 'desc' : 'asc')
  }

  return (
    <div className="max-w-7xl mx-auto">
      {/* Header */}
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-gray-900 mb-2 dark:text-gray-100">Digital Goods</h1>
        <p className="text-gray-600 dark:text-gray-300">Discover premium digital products from verified sellers</p>
      </div>

      {/* Search and Filters */}
      <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6 mb-8 dark:bg-gray-900 dark:border-gray-800">
        <div className="flex flex-col lg:flex-row gap-4">
          {/* Search */}
          <div className="flex-1">
            <div className="relative">
              <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 w-5 h-5" />
              <input
                type="text"
                placeholder="Search digital goods..."
                value={searchTerm}
                onChange={(e) => handleSearchChange(e.target.value)}
                className="w-full pl-10 pr-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent dark:bg-gray-800 dark:border-gray-700 dark:text-gray-100 placeholder:text-gray-400"
              />
            </div>
          </div>

          {/* Category Filter */}
          <div className="lg:w-48">
            <select
              value={selectedCategory}
              onChange={(e) => handleCategoryChange(e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent dark:bg-gray-800 dark:border-gray-700 dark:text-gray-100"
            >
              {categoryOptions.map((category) => (
                <option key={category.id} value={category.id}>
                  {category.name}
                </option>
              ))}
            </select>
          </div>

          {/* Sort */}
          <div className="lg:w-48">
            <select
              value={sortBy}
              onChange={(e) => setSortBy(e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent dark:bg-gray-800 dark:border-gray-700 dark:text-gray-100"
            >
              {sortOptions.map((option) => (
                <option key={option.value} value={option.value}>
                  {option.label}
                </option>
              ))}
            </select>
          </div>

          {/* Sort Order Toggle */}
          <button
            onClick={toggleSortOrder}
            className="px-4 py-2 border border-gray-300 rounded-lg hover:bg-gray-50 transition-colors dark:border-gray-700 dark:hover:bg-gray-800"
            title={sortOrder === 'asc' ? 'Sort Ascending' : 'Sort Descending'}
          >
            {sortOrder === 'asc' ? <SortAsc className="w-5 h-5" /> : <SortDesc className="w-5 h-5" />}
          </button>
        </div>
      </div>

      {/* Loading State */}
      {isLoading && (
        <div className="flex items-center justify-center py-12">
          <Loader2 className="w-8 h-8 animate-spin text-blue-600" />
          <span className="ml-2 text-gray-600">Loading products...</span>
        </div>
      )}

      {/* Error State */}
      {error && !isLoading && (
        <div className="text-center py-12">
          <div className="w-16 h-16 bg-red-100 rounded-full flex items-center justify-center mx-auto mb-4">
            <Search className="w-8 h-8 text-red-400" />
          </div>
          <h3 className="text-lg font-medium text-gray-900 mb-2">Error loading products</h3>
          <p className="text-gray-600">{error}</p>
        </div>
      )}

      {/* Results Count */}
      {!isLoading && !error && (
        <div className="mb-6">
          <p className="text-gray-600 dark:text-gray-300">
            Showing {filteredProducts.length} of {products.length} digital goods
            {selectedCategory !== 'all' && (
              <span className="ml-2">
                in <span className="font-medium">{categoryOptions.find(cat => cat.id === selectedCategory)?.name}</span>
              </span>
            )}
          </p>
        </div>
      )}

      {/* Products Grid */}
      {!isLoading && !error && (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
          {filteredProducts.map((product) => {
            const category = categories.find(cat => cat.id === product.category_id)
            const sellerName = (product as unknown as { seller?: { username?: string } })?.seller?.username || ''
            return (
              <div key={product.id} className="bg-white rounded-lg shadow-sm border border-gray-200 hover:shadow-md transition-shadow dark:bg-gray-900 dark:border-gray-800">
                <div className="p-6">
                  {/* Header */}
                  <div className="flex items-center justify-between mb-4">
                    <Monitor className="w-8 h-8 text-blue-600" />
                    <span className="text-xs text-gray-500 bg-gray-100 px-2 py-1 rounded-full capitalize dark:bg-gray-800 dark:text-gray-300">
                      {category?.name || 'Digital Good'}
                    </span>
                  </div>

                  {/* Title and Description */}
                  <h3 className="font-semibold text-gray-900 mb-2 line-clamp-2 dark:text-gray-100">
                    {product.title}
                  </h3>
                  <p className="text-gray-600 text-sm mb-4 line-clamp-3 dark:text-gray-300">
                    {product.description}
                  </p>

                  {/* Sales Count */}
                  <div className="flex items-center justify-end mb-4">
                    <span className="text-sm text-gray-500 dark:text-gray-400">{product.sold_count} sales</span>
                  </div>

                  {/* Admin-only meta */}
                  {isAdmin && (
                    <div className="text-sm text-gray-500 mb-4 dark:text-gray-400">
                      <p>Status: {product.is_active ? 'Active' : 'Inactive'}</p>
                      <p>Seller ID: {product.seller_id}</p>
                    </div>
                  )}

                  {/* Status indicator for inactive products */}
                  {!product.is_active && (
                    <div className="mb-4">
                      <ProductStatusBadge isActive={product.is_active} />
                    </div>
                  )}

                  {/* Price and Action */}
                  <div className="flex items-center justify-between">
                    <span className="text-2xl font-bold text-blue-600">${(product.price / 100).toFixed(2)}</span>
                    <span className="text-sm text-gray-500 mr-2 dark:text-gray-400">{sellerName ? `by ${sellerName}` : ''}</span>
                    {product.is_active ? (
                      <Link
                        to={`/products/${product.id}`}
                        className="bg-blue-600 text-white px-4 py-2 rounded-lg text-sm font-medium hover:bg-blue-700 transition-colors dark:bg-blue-500 dark:hover:bg-blue-600"
                      >
                        View Details
                      </Link>
                    ) : (
                      <button
                        disabled
                        className="bg-gray-400 text-white px-4 py-2 rounded-lg text-sm font-medium cursor-not-allowed dark:bg-gray-600"
                        title="This product is currently inactive"
                      >
                        Unavailable
                      </button>
                    )}
                  </div>
                </div>
              </div>
            )
          })}
        </div>
      )}

      {/* No Results */}
      {!isLoading && !error && filteredProducts.length === 0 && (
        <div className="text-center py-12">
          <div className="w-16 h-16 bg-gray-100 rounded-full flex items-center justify-center mx-auto mb-4 dark:bg-gray-800">
            <Search className="w-8 h-8 text-gray-400 dark:text-gray-500" />
          </div>
          <h3 className="text-lg font-medium text-gray-900 mb-2 dark:text-gray-100">No digital goods found</h3>
          <p className="text-gray-600 dark:text-gray-400">Try adjusting your search or filter criteria</p>
        </div>
      )}
    </div>
  )
}

export default ProductsPage

