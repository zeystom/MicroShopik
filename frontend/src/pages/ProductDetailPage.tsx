import React, { useState, useEffect } from 'react'
import { useParams, Link, useNavigate } from 'react-router-dom'
import { 
  Shield, 
  Clock, 
  User, 
  CheckCircle,
  Monitor, 
  ArrowLeft,
  Loader2
} from 'lucide-react'
import { apiService } from '@/services/api'
import { Product, Category } from '@/types'
import { useAuthStore } from '@/store/authStore'
import { toast } from 'react-hot-toast'

const ProductDetailPage: React.FC = () => {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const { isAuthenticated, user } = useAuthStore()
  const [product, setProduct] = useState<Product | null>(null)
  const [category, setCategory] = useState<Category | null>(null)
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState('')
  const [isPurchasing, setIsPurchasing] = useState(false)
  const isOwnProduct = !!(user && product && product.seller_id === user.id)

  // Fetch product data on component mount
  useEffect(() => {
    const fetchProduct = async () => {
      if (!id) return
      
      try {
        setIsLoading(true)
        const productData = await apiService.getProduct(Number(id))
        setProduct(productData)
        
        // Fetch category data
        if (productData.category_id) {
          try {
            const categoryData = await apiService.getCategory(productData.category_id)
            setCategory(categoryData)
          } catch (err) {
            console.warn('Failed to fetch category:', err)
          }
        }
      } catch (err: unknown) {
        const message =
          typeof err === 'object' && err !== null && 'response' in err
            ? (err as { response?: { data?: { error?: string } } }).response?.data?.error || 'Failed to load product'
            : 'Failed to load product'
        setError(message)
        toast.error('Failed to load product')
      } finally {
        setIsLoading(false)
      }
    }

    fetchProduct()
  }, [id])

  const handlePurchase = async () => {
    if (!isAuthenticated) {
      toast.error('Please log in to purchase this product')
      navigate('/login')
      return
    }

    if (!product) return
    if (isOwnProduct) {
      toast.error("You can't purchase your own product")
      return
    }

    try {
      setIsPurchasing(true)
      
      // Create order
      const order = await apiService.createOrder({
        customer_id: user?.id as number,
        product_id: product.id,
        status: 'pending'
      })
      
      // Create conversation between buyer and seller
      try {
        const conversation = await apiService.createConversation(
          { product_id: product.id },
          [user?.id as number, product.seller_id]
        )
        // Send auto message from seller about purchase
        try {
          await apiService.sendMessage(conversation.id, {
            text: `Purchase: "${product.title}"`,
            sender_id: product.seller_id,
            order_id: order.id,
          })
        } catch (e) {
          console.warn('Failed to send auto purchase message:', e)
        }
        toast.success('Order created successfully! A chat has been created with the seller.')
      } catch (chatError) {
        console.warn('Failed to create chat:', chatError)
        toast.success('Order created successfully!')
      }
      
      // Navigate to orders page
      navigate('/orders')
    } catch (err: unknown) {
      const message =
        typeof err === 'object' && err !== null && 'response' in err
          ? (err as { response?: { data?: { error?: string } } }).response?.data?.error || 'Failed to create order'
          : 'Failed to create order'
      toast.error(message)
    } finally {
      setIsPurchasing(false)
    }
  }

  if (isLoading) {
    return (
      <div className="flex items-center justify-center min-h-[60vh]">
        <Loader2 className="w-8 h-8 animate-spin text-blue-600" />
        <span className="ml-2 text-gray-600">Loading product...</span>
      </div>
    )
  }

  if (error || !product) {
    return (
      <div className="flex items-center justify-center min-h-[60vh] bg-gray-50 dark:bg-gray-950">
        <div className="text-center">
          <h1 className="text-2xl font-bold text-gray-900 mb-4 dark:text-gray-100">Product Not Found</h1>
          <p className="text-gray-600 mb-6 dark:text-gray-400">{error || 'The digital good you\'re looking for doesn\'t exist.'}</p>
          <Link
            to="/products"
            className="bg-blue-600 text-white px-6 py-3 rounded-lg font-medium hover:bg-blue-700 transition-colors dark:bg-blue-500 dark:hover:bg-blue-600"
          >
            Browse Digital Goods
          </Link>
        </div>
      </div>
    )
  }

  // Проверяем, активен ли продукт
  if (!product.is_active) {
    return (
      <div className="flex items-center justify-center min-h-[60vh] bg-gray-50 dark:bg-gray-950">
        <div className="text-center">
          <h1 className="text-2xl font-bold text-gray-900 mb-4 dark:text-gray-100">Product Unavailable</h1>
          <p className="text-gray-600 mb-6 dark:text-gray-400">This product is currently inactive and not available for purchase.</p>
          <Link
            to="/products"
            className="bg-blue-600 text-white px-6 py-3 rounded-lg font-medium hover:bg-blue-700 transition-colors dark:bg-blue-500 dark:hover:bg-blue-600"
          >
            Browse Available Digital Goods
          </Link>
        </div>
      </div>
    )
  }

  return (
    <div className="max-w-7xl mx-auto">
      {/* Breadcrumb */}
      <nav className="mb-8">
        <Link
          to="/products"
          className="flex items-center text-blue-600 hover:text-blue-700 transition-colors"
        >
          <ArrowLeft className="w-4 h-4 mr-2" />
          Back to Digital Goods
        </Link>
      </nav>

      {/* Product Details */}
      <div className="bg-white rounded-lg shadow-sm border border-gray-200 overflow-hidden dark:bg-gray-900 dark:border-gray-800">
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-8 p-8">
          {/* Left Column - Product Info */}
          <div>
            {/* Header */}
            <div className="flex items-center justify-between mb-6">
              <div className="flex items-center space-x-3">
                <Monitor className="w-12 h-12 text-blue-600" />
                <span className="text-sm text-gray-500 bg-gray-100 px-3 py-1 rounded-full capitalize dark:bg-gray-800 dark:text-gray-300">
                  {category?.name || 'Digital Good'}
                </span>
              </div>
              <div className="text-sm text-gray-500 dark:text-gray-400">{product.sold_count} sales</div>
            </div>

            {/* Title and Description */}
            <h1 className="text-3xl font-bold text-gray-900 mb-2 dark:text-gray-100">{product.title}</h1>
            <p className="text-sm text-gray-500 mb-2 dark:text-gray-400">by Seller #{product.seller_id}</p>
            <p className="text-gray-600 text-lg mb-6 dark:text-gray-300">{product.description}</p>

            {/* Seller Info */}
            <div className="border-t border-gray-200 pt-6 dark:border-gray-800">
              <div className="flex items-center justify-between">
                <div className="flex items-center space-x-3">
                  <User className="w-5 h-5 text-gray-400" />
                  <div>
                    <p className="text-sm text-gray-500 dark:text-gray-400">Seller ID</p>
                    <p className="font-medium text-gray-900 dark:text-gray-100">{product.seller_id}</p>
                  </div>
                </div>
                <div className="flex items-center space-x-3">
                  <div className="flex items-center space-x-2">
                    <Shield className="w-5 h-5 text-green-500" />
                    <span className="text-sm text-gray-600 dark:text-gray-300">Trusted seller</span>
                  </div>
                </div>
              </div>
            </div>
          </div>

          {/* Right Column - Purchase */}
          <div className="lg:border-l lg:border-gray-200 lg:pl-8 dark:border-gray-800">
            {/* Price and Purchase */}
            <div className="sticky top-8">
              <div className="bg-blue-50 rounded-lg p-6 mb-6 dark:bg-blue-950">
                <div className="text-center mb-6">
                  <span className="text-4xl font-bold text-blue-600">${(product.price / 100).toFixed(2)}</span>
                </div>

                <div className="space-y-4">
                  <button 
                    onClick={handlePurchase}
                    disabled={isPurchasing || isOwnProduct}
                    className="w-full bg-blue-600 text-white py-3 px-6 rounded-lg font-semibold hover:bg-blue-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center dark:bg-blue-500 dark:hover:bg-blue-600"
                  >
                    {isOwnProduct ? (
                      'Not available for your own product'
                    ) : isPurchasing ? (
                      <>
                        <Loader2 className="w-4 h-4 animate-spin mr-2" />
                        Processing...
                      </>
                    ) : (
                      'Purchase Now'
                    )}
                  </button>
                  
                  <div className="text-center">
                    <p className="text-sm text-gray-600 dark:text-gray-300">
                      Secure payment • Instant delivery • 24/7 support
                    </p>
                  </div>
                </div>
              </div>

              {/* Trust Indicators */}
              <div className="space-y-4">
                <div className="flex items-center space-x-3 p-3 bg-green-50 rounded-lg dark:bg-green-900/30">
                  <CheckCircle className="w-5 h-5 text-green-500" />
                  <span className="text-sm text-green-700 dark:text-green-300">Trusted seller</span>
                </div>
                
                <div className="flex items-center space-x-3 p-3 bg-blue-50 rounded-lg dark:bg-blue-900/30">
                  <Clock className="w-5 h-5 text-blue-500" />
                  <span className="text-sm text-blue-700 dark:text-blue-300">Instant digital delivery</span>
                </div>
                
                <div className="flex items-center space-x-3 p-3 bg-purple-50 rounded-lg dark:bg-purple-900/30">
                  <Shield className="w-5 h-5 text-purple-500" />
                  <span className="text-sm text-purple-700 dark:text-purple-300">Secure transaction guaranteed</span>
                </div>
              </div>
            </div>
          </div>
        </div>

        {/* Description */}
        <div className="border-t border-gray-200 p-8 dark:border-gray-800">
          <h2 className="text-2xl font-bold text-gray-900 mb-4 dark:text-gray-100">Description</h2>
          <div className="prose max-w-none text-gray-700 dark:text-gray-300">
            <p className="whitespace-pre-line">{product.description}</p>
          </div>
        </div>
      </div>
    </div>
  )
}

export default ProductDetailPage



