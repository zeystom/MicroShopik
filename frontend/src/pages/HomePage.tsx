import React, { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { 
  Gamepad2, 
  Monitor, 
  Smartphone, 
  Tv, 
  Headphones, 
  ArrowRight,
  Shield,
  Zap,
  Users,
  Loader2
} from 'lucide-react'
import { apiService } from '@/services/api'
import { Product, Category } from '@/types'
import { toast } from 'react-hot-toast'

const HomePage: React.FC = () => {
  const [isLoading, setIsLoading] = useState(true)
  const [featuredProducts, setFeaturedProducts] = useState<Product[]>([])
  const [categories, setCategories] = useState<Category[]>([])

  useEffect(() => {
    const load = async () => {
      try {
        setIsLoading(true)
        const [productsData, categoriesData] = await Promise.all([
          apiService.getProducts({ limit: 4 }),
          apiService.getCategories(),
        ])
        setFeaturedProducts(productsData)
        setCategories(categoriesData)
      } catch (err: unknown) {
        const message =
          typeof err === 'object' && err !== null && 'response' in err
            ? (err as { response?: { data?: { error?: string } } }).response?.data?.error || 'Failed to load homepage data'
            : 'Failed to load homepage data'
        toast.error(message)
        setFeaturedProducts([])
        setCategories([])
      } finally {
        setIsLoading(false)
      }
    }
    load()
  }, [])

  const categoryIcon = (idx: number) => {
    const icons = [Gamepad2, Monitor, Smartphone, Tv, Headphones]
    return icons[idx % icons.length]
  }

  const features = [
    {
      icon: Shield,
      title: 'Secure Transactions',
      description: 'End-to-end encryption and secure payment processing'
    },
    {
      icon: Zap,
      title: 'Instant Delivery',
      description: 'Get your digital goods immediately after purchase'
    },
    {
      icon: Users,
      title: 'Verified Sellers',
      description: 'All sellers are verified and rated by our community'
    }
  ]

  return (
    <div className="min-h-screen bg-gray-50 dark:bg-gray-950 dark:text-gray-100">
      {/* Hero Section */}
      <section className="bg-gradient-to-r from-blue-600 via-purple-600 to-blue-800 text-white">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-24">
          <div className="text-center">
            <h1 className="text-4xl md:text-6xl font-bold mb-6">
              Premium Digital Goods
            </h1>
            <p className="text-xl md:text-2xl mb-8 text-blue-100">
              Discover gaming accounts, streaming services, software licenses, and more
            </p>
            <div className="flex flex-col sm:flex-row gap-4 justify-center">
              <Link
                to="/products"
                className="bg-white text-blue-600 px-8 py-3 rounded-lg font-semibold hover:bg-gray-100 transition-colors inline-flex items-center"
              >
                Browse Digital Goods
                <ArrowRight className="ml-2 w-5 h-5" />
              </Link>
            </div>
          </div>
        </div>
      </section>

      {/* Featured Digital Goods */}
      <section className="py-16">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center mb-12">
            <h2 className="text-3xl font-bold text-gray-900 mb-4 dark:text-gray-100">
              Featured Digital Goods
            </h2>
            <p className="text-lg text-gray-600 dark:text-gray-300">
              Handpicked premium digital products from verified sellers
            </p>
          </div>
          
          {isLoading ? (
            <div className="flex items-center justify-center py-12">
              <Loader2 className="w-8 h-8 animate-spin text-blue-600" />
              <span className="ml-2 text-gray-600">Loading featured products...</span>
            </div>
          ) : (
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
              {featuredProducts.slice(0, 4).map((product) => {
                const sellerName = (product as any)?.seller?.username || ''
                return (
                <div key={product.id} className="h-full flex flex-col bg-white rounded-lg shadow-md hover:shadow-lg transition-shadow p-6 dark:bg-gray-900 dark:border-gray-800">
                  <div className="flex items-center justify-start mb-4">
                    <Monitor className="w-8 h-8 text-blue-600" />
                  </div>
                  <h3 className="font-semibold text-gray-900 mb-2 line-clamp-2 dark:text-gray-100">{product.title}</h3>
                  <p className="text-gray-600 text-sm mb-4 line-clamp-3 dark:text-gray-300">{product.description}</p>
                  <div className="flex items-center justify-end mb-4">
                    <span className="text-sm text-gray-500 dark:text-gray-400">{product.sold_count} sales</span>
                  </div>
                  <div className="mt-auto flex items-center justify-between">
                    <span className="text-2xl font-bold text-blue-600">${(product.price / 100).toFixed(2)}</span>
                    <span className="text-sm text-gray-500 mr-2 dark:text-gray-400">{sellerName ? `by ${sellerName}` : ''}</span>
                    <Link
                      to={`/products/${product.id}`}
                      className="bg-blue-600 text-white px-4 py-2 rounded-lg text-sm font-medium hover:bg-blue-700 transition-colors dark:bg-blue-500 dark:hover:bg-blue-600"
                    >
                      View Details
                    </Link>
                  </div>
                </div>
              )})}
            </div>
          )}
        </div>
      </section>

      {/* Categories Section */}
      <section className="py-16 bg-white dark:bg-gray-900">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center mb-12">
            <h2 className="text-3xl font-bold text-gray-900 mb-4 dark:text-gray-100">
              Explore Categories
            </h2>
            <p className="text-lg text-gray-600 dark:text-gray-300">
              Find exactly what you're looking for in our organized digital marketplace
            </p>
          </div>
          
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {categories.map((category, index) => {
              const Icon = categoryIcon(index)
              return (
                <Link
                  key={category.id}
                  to={`/products?category_id=${category.id}`}
                  className="group bg-gray-50 rounded-lg p-6 hover:bg-blue-50 transition-colors dark:bg-gray-800 dark:hover:bg-gray-700"
                >
                  <div className="flex items-center mb-4">
                    <Icon className="w-10 h-10 text-blue-600 group-hover:text-blue-700" />
                    <div className="ml-4">
                      <h3 className="text-lg font-semibold text-gray-900 group-hover:text-blue-700 dark:text-gray-100">{category.name}</h3>
                      <p className="text-sm text-gray-500 dark:text-gray-400">Browse items</p>
                    </div>
                  </div>
                  <p className="text-gray-600 text-sm dark:text-gray-300">Find digital goods in the {category.name} category</p>
                </Link>
              )
            })}
          </div>
        </div>
      </section>

      {/* Features Section */}
      <section className="py-16 bg-gray-50 dark:bg-gray-900">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center mb-12">
            <h2 className="text-3xl font-bold text-gray-900 mb-4 dark:text-gray-100">
              Why Choose MicroShopik?
            </h2>
            <p className="text-lg text-gray-600 dark:text-gray-300">
              We provide the safest and most convenient way to buy digital products
            </p>
          </div>
          
          <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
            {features.map((feature) => {
              const Icon = feature.icon
              return (
                <div key={feature.title} className="text-center">
                  <div className="w-16 h-16 bg-blue-100 rounded-full flex items-center justify-center mx-auto mb-4 dark:bg-blue-900">
                    <Icon className="w-8 h-8 text-blue-600" />
                  </div>
                  <h3 className="text-xl font-semibold text-gray-900 mb-2 dark:text-gray-100">
                    {feature.title}
                  </h3>
                  <p className="text-gray-600 dark:text-gray-300">{feature.description}</p>
                </div>
              )
            })}
          </div>
        </div>
      </section>

      {/* CTA Section */}
      <section className="py-16 bg-blue-600">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 text-center">
          <h2 className="text-3xl font-bold text-white mb-4">
            Ready to Start?
          </h2>
          <p className="text-xl text-blue-100 mb-8">
            Join thousands of users buying and selling digital goods
          </p>
          <div className="flex flex-col sm:flex-row gap-4 justify-center">
            <Link
              to="/products"
              className="bg-white text-blue-600 px-8 py-3 rounded-lg font-semibold hover:bg-gray-100 transition-colors"
            >
              Browse Now
            </Link>
          </div>
        </div>
      </section>
    </div>
  )
}

export default HomePage

