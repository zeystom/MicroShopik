import React, { useEffect, useState } from 'react'
import { useParams, Link, useNavigate } from 'react-router-dom'
import { apiService } from '@/services/api'
import { Order } from '@/types'
import { useAuthStore } from '@/store/authStore'
import { useTheme } from '@/hooks/useTheme'
import { Loader2, ArrowLeft, MessageCircle, CheckCircle2, XCircle } from 'lucide-react'
import { toast } from 'react-hot-toast'
import ProductStatusBadge from '@/components/ui/ProductStatusBadge'

const OrderDetailPage: React.FC = () => {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const { user } = useAuthStore()
  const { theme: _theme } = useTheme()
  const [order, setOrder] = useState<Order | null>(null)
  const [isLoading, setIsLoading] = useState(true)
  const [isWorking, setIsWorking] = useState(false)

  const orderId = Number(id)

  useEffect(() => {
    const load = async () => {
      if (!orderId) return
      try {
        setIsLoading(true)
        const data = await apiService.getOrder(orderId)
        setOrder(data)
      } catch (err: any) {
        toast.error(err.response?.data?.error || 'Failed to load order')
      } finally {
        setIsLoading(false)
      }
    }
    load()
  }, [orderId])

  const isSeller = !!user?.roles?.some(r => r.name === 'seller')
  const isCustomerOwner = !!(user && order?.customer_id === user.id)

  const refresh = async () => {
    try {
      const data = await apiService.getOrder(orderId)
      setOrder(data)
    } catch {}
  }

  const handleProcess = async () => {
    try {
      setIsWorking(true)
      await apiService.processOrder(orderId)
      toast.success('Order processed')
      await refresh()
    } catch (err: any) {
      toast.error(err.response?.data?.error || 'Failed to process order')
    } finally {
      setIsWorking(false)
    }
  }

  const handleCancel = async () => {
    if (!user) return
    try {
      setIsWorking(true)
      await apiService.cancelOrder(orderId, user.id)
      toast.success('Order cancelled')
      await refresh()
    } catch (err: any) {
      toast.error(err.response?.data?.error || 'Failed to cancel order')
    } finally {
      setIsWorking(false)
    }
  }

  const getStatusBadgeClass = (status: string) => {
    switch (status) {
      case 'completed':
      case 'confirmed':
        return 'status-badge success'
      case 'pending':
        return 'status-badge warning'
      default:
        return 'status-badge error'
    }
  }

  if (isLoading) {
    return (
      <div className="flex items-center justify-center min-h-[60vh] bg-gray-50 dark:bg-gray-950">
        <Loader2 className="w-8 h-8 animate-spin text-blue-600 dark:text-blue-400" />
        <span className="ml-2 text-gray-600 dark:text-gray-300">Loading order...</span>
      </div>
    )
  }

  if (!order) {
    return (
      <div className="flex items-center justify-center min-h-[60vh] bg-gray-50 dark:bg-gray-950">
        <div className="text-center">
          <h1 className="text-2xl font-bold text-gray-900 dark:text-gray-100 mb-4">Order Not Found</h1>
          <Link to="/orders" className="text-blue-600 dark:text-blue-400 hover:text-blue-800 dark:hover:text-blue-300">Go back to Orders</Link>
        </div>
      </div>
    )
  }

  return (
    <div className="max-w-5xl mx-auto px-4 sm:px-6 lg:px-8 py-8 bg-gray-50 dark:bg-gray-950 min-h-screen">
      <nav className="mb-6">
        <button onClick={() => navigate(-1)} className="flex items-center text-blue-600 dark:text-blue-400 hover:text-blue-700 dark:hover:text-blue-300 transition-colors">
          <ArrowLeft className="w-4 h-4 mr-2" />
          Back
        </button>
      </nav>

      <div className="card-dark">
        <div className="px-6 py-4 border-b border-gray-200 dark:border-gray-700 flex items-center justify-between">
          <div>
            <h1 className="text-2xl font-bold text-gray-900 dark:text-gray-100">Order #{order.id}</h1>
            <p className="text-gray-600 dark:text-gray-400">Placed on {new Date(order.created_at).toLocaleString()}</p>
          </div>
          <span className={getStatusBadgeClass(order.status)}>
            {order.status}
          </span>
        </div>

        <div className="p-6 space-y-6">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div className="card-section">
              <h2 className="font-semibold text-gray-900 dark:text-gray-100 mb-3">Product</h2>
              {order.product ? (
                <div>
                  <div className="text-gray-900 dark:text-gray-100 font-medium">{order.product.title}</div>
                  <div className="text-gray-600 dark:text-gray-400">Seller ID: {order.product.seller_id}</div>
                  <div className="text-gray-900 dark:text-gray-100 mt-2">Price: ${((order.product.price || 0) / 100).toFixed(2)}</div>
                  {!order.product.is_active && (
                    <div className="mt-2">
                      <ProductStatusBadge isActive={order.product.is_active} />
                    </div>
                  )}
                </div>
              ) : (
                <div className="text-gray-600 dark:text-gray-400">No product info</div>
              )}
            </div>

            <div className="card-section">
              <h2 className="font-semibold text-gray-900 dark:text-gray-100 mb-3">Actions</h2>
              <div className="flex flex-wrap gap-3">
                <Link 
                  to={`/conversations?orderId=${order.id}`} 
                  className="btn-action-blue"
                >
                  <MessageCircle className="w-4 h-4 mr-2" /> Open Chat
                </Link>
                {isSeller && order.status === 'pending' && (
                  <button 
                    onClick={handleProcess} 
                    disabled={isWorking} 
                    className="btn-action-green"
                  >
                    {isWorking ? <Loader2 className="w-4 h-4 mr-2 animate-spin" /> : <CheckCircle2 className="w-4 h-4 mr-2" />} Process
                  </button>
                )}
                {isCustomerOwner && order.status === 'pending' && (
                  <button 
                    onClick={handleCancel} 
                    disabled={isWorking} 
                    className="btn-action-red"
                  >
                    {isWorking ? <Loader2 className="w-4 h-4 mr-2 animate-spin" /> : <XCircle className="w-4 h-4 mr-2" />} Cancel
                  </button>
                )}
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}

export default OrderDetailPage





