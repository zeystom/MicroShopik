import React, { useState, useEffect } from 'react'
import { Link } from 'react-router-dom'
import { ShoppingBag, Loader2, ChevronDown, ChevronRight } from 'lucide-react'
import { Order } from '@/types'
import { apiService } from '@/services/api'
import { useAuthStore } from '@/store/authStore'
import { toast } from 'react-hot-toast'

const OrdersPage: React.FC = () => {
  const [purchases, setPurchases] = useState<Order[]>([])
  const [sales, setSales] = useState<Order[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const { user } = useAuthStore()
  const [showPurchases, setShowPurchases] = useState(true)
  const [showSales, setShowSales] = useState(true)

  useEffect(() => {
    const fetchData = async () => {
      try {
        setIsLoading(true)
        const [purchasesData, salesData] = await Promise.all([
          apiService.getOrders(),
          apiService.getSellerOrders().catch(() => [])
        ])
        setPurchases(purchasesData || [])
        setSales(salesData || [])
      } catch (error: any) {
        console.error('Failed to fetch data:', error)
        toast.error('Failed to load data')
        setPurchases([])
        setSales([])
      } finally {
        setIsLoading(false)
      }
    }

    if (user?.id) {
      fetchData()
    } else {
      setIsLoading(false)
    }
  }, [user])

  if (isLoading) {
    return (
      <div className="flex items-center justify-center min-h-[60vh]">
        <Loader2 className="w-8 h-8 animate-spin text-blue-600" />
        <span className="ml-2 text-gray-600 dark:text-gray-300">Loading...</span>
      </div>
    )
  }

  if (!user) {
    return (
      <div className="text-center py-12">
        <p className="text-gray-600 dark:text-gray-300">Please log in to view your orders.</p>
      </div>
    )
  }

  // Customer + Seller view with sections
  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-gray-900 dark:text-gray-100">Orders</h1>
        <p className="text-gray-600 mt-1 dark:text-gray-300">Your purchases and your sales</p>
      </div>

      {/* Purchases Section */}
      <div className="mb-8">
        <button
          onClick={() => setShowPurchases(!showPurchases)}
          className="w-full flex items-center justify-between bg-white dark:bg-gray-900 border border-gray-200 dark:border-gray-800 rounded-lg px-4 py-3"
        >
          <div className="flex items-center space-x-2">
            {showPurchases ? <ChevronDown className="w-5 h-5" /> : <ChevronRight className="w-5 h-5" />}
            <h2 className="text-lg font-semibold text-gray-900 dark:text-gray-100">My Purchases</h2>
          </div>
          <span className="text-sm text-gray-500 dark:text-gray-400">{purchases.length} orders</span>
        </button>
        {showPurchases && (
          purchases.length > 0 ? (
            <div className="mt-4 bg-white rounded-lg shadow border border-gray-200 overflow-x-auto dark:bg-gray-900 dark:border-gray-800">
              <table className="min-w-full divide-y divide-gray-200 dark:divide-gray-800">
                <thead className="bg-gray-50 dark:bg-gray-800">
                  <tr>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Order ID</th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Product</th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Status</th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Date</th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Actions</th>
                  </tr>
                </thead>
                <tbody className="bg-white divide-y divide-gray-200 dark:bg-gray-900 dark:divide-gray-800">
                  {purchases.map((order) => (
                    <tr key={order.id}>
                      <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900 dark:text-gray-100">#{order.id}</td>
                      <td className="px-6 py-4 whitespace-nowrap">
                        <div className="text-sm text-gray-900 dark:text-gray-100">{order.product?.title || 'Product not found'}</div>
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap">
                        <span className={`inline-flex px-2 py-1 text-xs font-semibold rounded-full ${
                          order.status === 'completed' || order.status === 'confirmed'
                            ? 'bg-green-100 text-green-800'
                            : order.status === 'pending'
                            ? 'bg-yellow-100 text-yellow-800'
                            : 'bg-red-100 text-red-800'
                        }`}>{order.status}</span>
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500 dark:text-gray-400">{new Date(order.created_at).toLocaleDateString()}</td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm font-medium">
                        <Link to={`/orders/${order.id}`} className="text-blue-600 hover:text-blue-900">View Details</Link>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          ) : (
            <div className="mt-4 text-center py-8 bg-white border border-gray-200 rounded-lg dark:bg-gray-900 dark:border-gray-800">
              <ShoppingBag className="w-12 h-12 text-gray-400 mx-auto mb-4" />
              <h3 className="text-lg font-medium text-gray-900 mb-2 dark:text-gray-100">No purchases yet</h3>
              <p className="text-gray-600 mb-4 dark:text-gray-300">Start shopping to see your orders here</p>
              <Link to="/products" className="bg-blue-600 text-white px-4 py-2 rounded-lg font-medium hover:bg-blue-700 transition-colors inline-flex items-center">
                <ShoppingBag className="mr-2 h-4 w-4" /> Browse Products
              </Link>
            </div>
          )
        )}
      </div>

      {/* Sales Section */}
      <div className="mb-2">
        <button
          onClick={() => setShowSales(!showSales)}
          className="w-full flex items-center justify-between bg-white dark:bg-gray-900 border border-gray-200 dark:border-gray-800 rounded-lg px-4 py-3"
        >
          <div className="flex items-center space-x-2">
            {showSales ? <ChevronDown className="w-5 h-5" /> : <ChevronRight className="w-5 h-5" />}
            <h2 className="text-lg font-semibold text-gray-900 dark:text-gray-100">My Sales</h2>
          </div>
          <span className="text-sm text-gray-500 dark:text-gray-400">{sales.length} orders</span>
        </button>
        {showSales && sales.length > 0 && (
          <div className="mt-4 bg-white rounded-lg shadow border border-gray-200 overflow-x-auto dark:bg-gray-900 dark:border-gray-800">
            <table className="min-w-full divide-y divide-gray-200 dark:divide-gray-800">
              <thead className="bg-gray-50 dark:bg-gray-800">
                <tr>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Order ID</th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Product</th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Buyer</th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Status</th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Date</th>
                </tr>
              </thead>
              <tbody className="bg-white divide-y divide-gray-200 dark:bg-gray-900 dark:divide-gray-800">
                {sales.map((order) => (
                  <tr key={order.id}>
                    <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900 dark:text-gray-100">#{order.id}</td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-gray-100">{order.product?.title || 'Product'}</td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500 dark:text-gray-400">User #{order.customer_id}</td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      <span className={`inline-flex px-2 py-1 text-xs font-semibold rounded-full ${
                        order.status === 'completed' || order.status === 'confirmed'
                          ? 'bg-green-100 text-green-800'
                          : order.status === 'pending'
                          ? 'bg-yellow-100 text-yellow-800'
                          : 'bg-red-100 text-red-800'
                      }`}>{order.status}</span>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500 dark:text-gray-400">{new Date(order.created_at).toLocaleDateString()}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>
    </div>
  )
}

export default OrdersPage

