import { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { Plus, TrendingUp, Package, DollarSign, Users, Loader2 } from 'lucide-react';
import { Product, Order } from '@/types';
import { apiService } from '@/services/api';
import { useAuthStore } from '@/store/authStore';
import { toast } from 'react-hot-toast';
import ProductStatusBadge from '@/components/ui/ProductStatusBadge';

const SellerDashboardPage = () => {
  const [products, setProducts] = useState<Product[]>([]);
  const [orders, setOrders] = useState<Order[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const { user } = useAuthStore();

  useEffect(() => {
    const fetchData = async () => {
      try {
        setIsLoading(true);
        
        // Fetch seller's products and orders using seller_id parameter
        console.log('Fetching data for seller ID:', user?.id);
        
        // Try to fetch products first
        let productsData = [];
        try {
          productsData = await apiService.getSellerProducts(user?.id || 0);
          setProducts(productsData);
        } catch (error) {
          console.error('Failed to fetch products:', error);
          setProducts([]);
        }
        
        // Try to fetch orders separately
        let ordersData = [];
        try {
          ordersData = await apiService.getSellerOrders();
          setOrders(ordersData);
        } catch (error) {
          console.error('Failed to fetch orders:', error);
          setOrders([]);
        }
      } catch (error: unknown) {
        console.error('Failed to fetch dashboard data:', error);
        toast.error('Failed to load dashboard data');
        // Set empty arrays as fallback
        setProducts([]);
        setOrders([]);
      } finally {
        setIsLoading(false);
      }
    };

    if (user?.id) {
      fetchData();
    } else {
      setIsLoading(false);
    }
  }, [user?.id]);

  const formatPrice = (price: number) => {
    return `$${(price / 100).toFixed(2)}`;
  };

  const totalRevenue = products.reduce((sum, product) => sum + (product.price * product.sold_count), 0);
  const totalSales = products.reduce((sum, product) => sum + product.sold_count, 0);
  const activeProducts = products.filter(product => product.is_active).length;
  const recentOrders = orders.slice(0, 5);

  if (isLoading) {
    return (
      <div className="flex items-center justify-center min-h-[60vh]">
        <Loader2 className="w-8 h-8 animate-spin text-blue-600" />
        <span className="ml-2 text-gray-600 dark:text-gray-300">Loading dashboard...</span>
      </div>
    );
  }

  if (!user) {
    return (
      <div className="text-center py-12">
        <p className="text-gray-600 dark:text-gray-300">Please log in to view your seller dashboard.</p>
      </div>
    );
  }

  if (products.length === 0) {
    return (
      <div className="space-y-6">
        {/* Header */}
        <div className="flex justify-between items-center">
          <div>
            <h1 className="text-3xl font-bold text-gray-900 dark:text-gray-100">Seller Dashboard</h1>
            <p className="text-gray-600 mt-1 dark:text-gray-300">
              Manage your products and track your sales performance
            </p>
          </div>
          <Link
            to="/seller/products/new"
            className="bg-blue-600 text-white px-4 py-2 rounded-lg font-medium hover:bg-blue-700 transition-colors inline-flex items-center"
          >
            <Plus className="mr-2 h-4 w-4" />
            Add New Product
          </Link>
        </div>

        {/* Empty State */}
        <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-10 text-center dark:bg-gray-900 dark:border-gray-800">
          <Package className="w-12 h-12 text-gray-400 mx-auto mb-4" />
          <h3 className="text-lg font-medium text-gray-900 mb-2 dark:text-gray-100">No products yet</h3>
          <p className="text-gray-600 mb-4 dark:text-gray-300">Start selling by adding your first product</p>
          <Link
            to="/seller/products/new"
            className="bg-blue-600 text-white px-4 py-2 rounded-lg font-medium hover:bg-blue-700 transition-colors inline-flex items-center"
          >
            <Plus className="mr-2 h-4 w-4" />
            Add Your First Product
          </Link>
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-3xl font-bold text-gray-900 dark:text-gray-100">Seller Dashboard</h1>
          <p className="text-gray-600 mt-1 dark:text-gray-300">
            Manage your products and track your sales performance
          </p>
        </div>
        <Link
          to="/seller/products/new"
          className="bg-blue-600 text-white px-4 py-2 rounded-lg font-medium hover:bg-blue-700 transition-colors inline-flex items-center"
        >
          <Plus className="mr-2 h-4 w-4" />
          Add New Product
        </Link>
      </div>

      {/* Stats Cards */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6 dark:bg-gray-900 dark:border-gray-800">
          <div className="flex items-center">
            <div className="bg-blue-100 rounded-full p-3">
              <DollarSign className="h-6 w-6 text-blue-600" />
            </div>
            <div className="ml-4">
              <p className="text-sm font-medium text-gray-600 dark:text-gray-300">Total Revenue</p>
              <p className="text-2xl font-bold text-gray-900 dark:text-gray-100">{formatPrice(totalRevenue)}</p>
            </div>
          </div>
        </div>

        <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6 dark:bg-gray-900 dark:border-gray-800">
          <div className="flex items-center">
            <div className="bg-green-100 rounded-full p-3">
              <Package className="h-6 w-6 text-green-600" />
            </div>
            <div className="ml-4">
              <p className="text-sm font-medium text-gray-600 dark:text-gray-300">Total Sales</p>
              <p className="text-2xl font-bold text-gray-900 dark:text-gray-100">{totalSales}</p>
            </div>
          </div>
        </div>

        <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6 dark:bg-gray-900 dark:border-gray-800">
          <div className="flex items-center">
            <div className="bg-purple-100 rounded-full p-3">
              <TrendingUp className="h-6 w-6 text-purple-600" />
            </div>
            <div className="ml-4">
              <p className="text-sm font-medium text-gray-600 dark:text-gray-300">Active Products</p>
              <p className="text-2xl font-bold text-gray-900 dark:text-gray-100">{activeProducts}</p>
            </div>
          </div>
        </div>

        <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6 dark:bg-gray-900 dark:border-gray-800">
          <div className="flex items-center">
            <div className="bg-yellow-100 rounded-full p-3">
              <Users className="h-6 w-6 text-yellow-600" />
            </div>
            <div className="ml-4">
              <p className="text-sm font-medium text-gray-600 dark:text-gray-300">Total Products</p>
              <p className="text-2xl font-bold text-gray-900 dark:text-gray-100">{products.length}</p>
            </div>
          </div>
        </div>
      </div>

      {/* Recent Activity */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Recent Orders */}
        <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6 dark:bg-gray-900 dark:border-gray-800">
          <div className="flex items-center justify-between mb-4">
            <h2 className="text-lg font-semibold text-gray-900 dark:text-gray-100">Recent Orders</h2>
            <Link
              to="/orders"
              className="text-blue-600 hover:text-blue-700 text-sm font-medium"
            >
              View All
            </Link>
          </div>
          
          <div className="space-y-3">
            {recentOrders.length > 0 ? (
              recentOrders.map((order) => (
                <div key={order.id} className="flex items-center justify-between p-3 bg-gray-50 rounded-lg dark:bg-gray-800">
                  <div>
                    <p className="font-medium text-gray-900 dark:text-gray-100">Order #{order.id}</p>
                    <p className="text-sm text-gray-600 dark:text-gray-300">
                      {products.find(p => p.id === order.product_id)?.title || 'Product'}
                    </p>
                  </div>
                  <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${
                    (order.status === 'completed' || order.status === 'confirmed')
                      ? 'bg-green-100 text-green-800'
                      : order.status === 'pending'
                      ? 'bg-yellow-100 text-yellow-800'
                      : 'bg-gray-100 text-gray-800'
                  }`}>
                    {order.status}
                  </span>
                </div>
              ))
            ) : (
              <p className="text-gray-500 text-center py-4">No orders yet</p>
            )}
          </div>
        </div>

        {/* Top Products */}
        <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6 dark:bg-gray-900 dark:border-gray-800">
          <div className="flex items-center justify-between mb-4">
            <h2 className="text-lg font-semibold text-gray-900 dark:text-gray-100">Top Products</h2>
            <Link
              to="/seller/products"
              className="text-blue-600 hover:text-blue-700 text-sm font-medium"
            >
              View All
            </Link>
          </div>
          
          <div className="space-y-3">
            {products
              .sort((a, b) => b.sold_count - a.sold_count)
              .slice(0, 3)
              .map((product) => (
                <div key={product.id} className="flex items-center justify-between p-3 bg-gray-50 rounded-lg dark:bg-gray-800">
                  <div>
                    <p className="font-medium text-gray-900 dark:text-gray-100">{product.title}</p>
                    <p className="text-sm text-gray-600 dark:text-gray-300">{product.sold_count} sold</p>
                  </div>
                  <span className="font-semibold text-blue-600">
                    {formatPrice(product.price)}
                  </span>
                </div>
              ))}
            {products.length === 0 && (
              <p className="text-gray-500 text-center py-4">No products yet</p>
            )}
          </div>
        </div>
      </div>

      {/* Products Table */}
      <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6 dark:bg-gray-900 dark:border-gray-800">
        <div className="flex items-center justify-between mb-6">
          <h2 className="text-lg font-semibold text-gray-900 dark:text-gray-100">Your Products</h2>
          <Link
            to="/seller/products"
            className="text-blue-600 hover:text-blue-700 text-sm font-medium"
          >
            Manage All Products
          </Link>
        </div>

        {products.length > 0 ? (
          <div className="overflow-x-auto">
            <table className="min-w-full divide-y divide-gray-200">
              <thead className="bg-gray-50 dark:bg-gray-800">
                <tr>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Product
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Price
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Sales
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Status
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Actions
                  </th>
                </tr>
              </thead>
              <tbody className="bg-white divide-y divide-gray-200 dark:bg-gray-900 dark:divide-gray-800">
                {products.slice(0, 5).map((product) => (
                  <tr key={product.id}>
                    <td className="px-6 py-4 whitespace-nowrap">
                      <div>
                        <div className="text-sm font-medium text-gray-900 dark:text-gray-100">{product.title}</div>
                        <div className="text-sm text-gray-600 dark:text-gray-300">{product.description.substring(0, 80)}...</div>
                      </div>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-gray-100">
                      {formatPrice(product.price)}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-gray-100">
                      {product.sold_count}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${
                        product.is_active 
                          ? 'bg-green-100 text-green-800' 
                          : 'bg-red-100 text-red-800'
                      }`}>
                        {product.is_active ? 'Active' : 'Inactive'}
                      </span>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm font-medium">
                      <Link
                        to={`/products/${product.id}`}
                        className="text-blue-600 hover:text-blue-900 mr-4"
                      >
                        View
                      </Link>
                      <Link
                        to={`/seller/products/${product.id}/edit`}
                        className="text-green-600 hover:text-green-900"
                      >
                        Edit
                      </Link>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        ) : (
          <div className="text-center py-8">
            <Package className="w-12 h-12 text-gray-400 mx-auto mb-4" />
            <h3 className="text-lg font-medium text-gray-900 mb-2">No products yet</h3>
            <p className="text-gray-600 mb-4">Start selling by adding your first product</p>
            <Link
              to="/seller/products/new"
              className="bg-blue-600 text-white px-4 py-2 rounded-lg font-medium hover:bg-blue-700 transition-colors inline-flex items-center"
            >
              <Plus className="mr-2 h-4 w-4" />
              Add Your First Product
            </Link>
          </div>
        )}
      </div>
    </div>
  );
};

export default SellerDashboardPage;

