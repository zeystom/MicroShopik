import { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { 
  Users, Package, DollarSign, TrendingUp, Shield, Settings, BarChart3, Loader2, 
  UserPlus, Trash2, 
  Activity, Database, Server, RefreshCw
} from 'lucide-react';
import { User, Product, Order, Category, Role } from '@/types';
import { apiService } from '@/services/api';
import { toast } from 'react-hot-toast';
import { useTheme } from '@/hooks/useTheme';

const AdminDashboardPage = () => {
  const { theme } = useTheme();
  const [users, setUsers] = useState<User[]>([]);
  const [products, setProducts] = useState<Product[]>([]);
  const [orders, setOrders] = useState<Order[]>([]);
  const [categories, setCategories] = useState<Category[]>([]);
  const [roles, setRoles] = useState<Role[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [showCategoryModal, setShowCategoryModal] = useState(false);
  const [newCategory, setNewCategory] = useState({ name: '', description: '' });

  useEffect(() => {
    fetchData();
  }, []);

  const fetchData = async () => {
    try {
      setIsLoading(true);
      
      const [usersData, productsData, ordersData, categoriesData, rolesData] = await Promise.all([
        apiService.getUsers(),
        apiService.getAllProducts(),
        apiService.getAllOrders(),
        apiService.getCategories(),
        apiService.getRoles()
      ]);
      
      setUsers(usersData);
      setProducts(productsData);
      setOrders(ordersData);
      setCategories(categoriesData);
      setRoles(rolesData);
    } catch (error) {
      console.error('Failed to fetch admin dashboard data:', error);
      toast.error('Failed to load dashboard data');
      setUsers([]);
      setProducts([]);
      setOrders([]);
      setCategories([]);
      setRoles([]);
    } finally {
      setIsLoading(false);
    }
  };

  const formatPrice = (price: number) => {
    return `$${(price / 100).toFixed(2)}`;
  };

  const handleDeleteUser = async (userId: number) => {
    if (window.confirm('Are you sure you want to delete this user?')) {
      try {
        await apiService.deleteUser(userId);
        setUsers(users.filter(u => u.id !== userId));
        toast.success('User deleted successfully');
      } catch (error) {
        toast.error('Failed to delete user');
      }
    }
  };

  const handleDeleteProduct = async (productId: number) => {
    if (window.confirm('Are you sure you want to delete this product?')) {
      try {
        await apiService.deleteProductAdmin(productId);
        setProducts(products.filter(p => p.id !== productId));
        toast.success('Product deleted successfully');
      } catch (error) {
        toast.error('Failed to delete product');
      }
    }
  };

  const handleDeleteCategory = async (categoryId: number) => {
    if (window.confirm('Are you sure you want to delete this category?')) {
      try {
        await apiService.deleteCategory(categoryId);
        setCategories(categories.filter(c => c.id !== categoryId));
        toast.success('Category deleted successfully');
      } catch (error) {
        toast.error('Failed to delete category');
      }
    }
  };

  const handleCreateCategory = async () => {
    try {
      const category = await apiService.createCategory(newCategory);
      setCategories([...categories, category]);
      setNewCategory({ name: '', description: '' });
      setShowCategoryModal(false);
      toast.success('Category created successfully');
    } catch (error) {
      toast.error('Failed to create category');
    }
  };

  const handleToggleProductStatus = async (productId: number, currentStatus: boolean) => {
    try {
      await apiService.updateProductStatus(productId, !currentStatus);
      setProducts(products.map(p => 
        p.id === productId ? { ...p, is_active: !currentStatus } : p
      ));
      toast.success('Product status updated');
    } catch (error) {
      toast.error('Failed to update product status');
    }
  };

  const totalRevenue = products.reduce((sum, product) => sum + (product.price * product.sold_count), 0);
  const totalUsers = users.length;
  const totalProducts = products.length;
  const activeProducts = products.filter(product => product.is_active).length;
  const totalOrders = orders.length;
  const completedOrders = orders.filter(order => order.status === 'completed').length;
  const pendingOrders = orders.filter(order => order.status === 'pending').length;

  if (isLoading) {
    return (
      <div className="flex items-center justify-center min-h-[60vh]">
        <Loader2 className="w-8 h-8 animate-spin text-blue-600 dark:text-blue-400" />
        <span className="ml-2 text-gray-600 dark:text-gray-300">Loading admin dashboard...</span>
      </div>
    );
  }

  return (
    <div className="space-y-6 p-6 bg-gray-50 dark:bg-gray-900 min-h-screen">
      {/* Header */}
      <div className="flex flex-col lg:flex-row justify-between items-start lg:items-center gap-4">
        <div>
          <h1 className="text-3xl font-bold text-gray-900 dark:text-white">Admin Dashboard</h1>
          <p className="text-gray-600 dark:text-gray-300 mt-1">
            System overview and administration tools
          </p>
        </div>
        <div className="flex flex-wrap gap-3">
          <button
            onClick={fetchData}
            className="bg-blue-600 text-white px-4 py-2 rounded-lg font-medium hover:bg-blue-700 transition-colors inline-flex items-center"
          >
            <RefreshCw className="mr-2 h-4 w-4" />
            Refresh
          </button>
          <Link
            to="/admin/categories"
            className="bg-emerald-600 text-white px-4 py-2 rounded-lg font-medium hover:bg-emerald-700 transition-colors inline-flex items-center"
          >
            <Package className="mr-2 h-4 w-4" />
            Manage Categories
          </Link>
          <Link
            to="/admin/users"
            className="bg-gray-600 text-white px-4 py-2 rounded-lg font-medium hover:bg-gray-700 transition-colors inline-flex items-center"
          >
            <Users className="mr-2 h-4 w-4" />
            Manage Users
          </Link>
          <Link
            to="/admin/roles"
            className="bg-purple-600 text-white px-4 py-2 rounded-lg font-medium hover:bg-purple-700 transition-colors inline-flex items-center"
          >
            <Shield className="mr-2 h-4 w-4" />
            Manage Roles
          </Link>
          <Link
            to="/admin/settings"
            className="bg-green-600 text-white px-4 py-2 rounded-lg font-medium hover:bg-green-700 transition-colors inline-flex items-center"
          >
            <Settings className="mr-2 h-4 w-4" />
            System Settings
          </Link>
        </div>
      </div>

      {/* Stats Cards */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        <div className="bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700 p-6">
          <div className="flex items-center">
            <div className="bg-blue-100 dark:bg-blue-900 rounded-full p-3">
              <Users className="h-6 w-6 text-blue-600 dark:text-blue-400" />
            </div>
            <div className="ml-4">
              <p className="text-sm font-medium text-gray-600 dark:text-gray-300">Total Users</p>
              <p className="text-2xl font-bold text-gray-900 dark:text-white">{totalUsers}</p>
            </div>
          </div>
        </div>

        <div className="bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700 p-6">
          <div className="flex items-center">
            <div className="bg-green-100 dark:bg-green-900 rounded-full p-3">
              <Package className="h-6 w-6 text-green-600 dark:text-green-400" />
            </div>
            <div className="ml-4">
              <p className="text-sm font-medium text-gray-600 dark:text-gray-300">Total Products</p>
              <p className="text-2xl font-bold text-gray-900 dark:text-white">{totalProducts}</p>
            </div>
          </div>
        </div>

        <div className="bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700 p-6">
          <div className="flex items-center">
            <div className="bg-purple-100 dark:bg-purple-900 rounded-full p-3">
              <TrendingUp className="h-6 w-6 text-purple-600 dark:text-purple-400" />
            </div>
            <div className="ml-4">
              <p className="text-sm font-medium text-gray-600 dark:text-gray-300">Active Products</p>
              <p className="text-2xl font-bold text-gray-900 dark:text-white">{activeProducts}</p>
            </div>
          </div>
        </div>

        <div className="bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700 p-6">
          <div className="flex items-center">
            <div className="bg-yellow-100 dark:bg-yellow-900 rounded-full p-3">
              <DollarSign className="h-6 w-6 text-yellow-600 dark:text-yellow-400" />
            </div>
            <div className="ml-4">
              <p className="text-sm font-medium text-gray-600 dark:text-gray-300">Platform Revenue</p>
              <p className="text-2xl font-bold text-gray-900 dark:text-white">{formatPrice(totalRevenue)}</p>
            </div>
          </div>
        </div>
      </div>

      {/* Quick Actions */}
      <div className="bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700 p-6">
        <h2 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">Quick Actions</h2>
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
          <Link
            to="/admin/users"
            className="p-4 border-2 border-dashed border-gray-300 dark:border-gray-600 rounded-lg text-center hover:border-blue-400 hover:bg-blue-50 dark:hover:bg-blue-900/20 transition-colors"
          >
            <UserPlus className="h-8 w-8 text-gray-400 dark:text-gray-500 mx-auto mb-2" />
            <p className="font-medium text-gray-900 dark:text-white">Manage Users</p>
            <p className="text-sm text-gray-600 dark:text-gray-400">View and edit user accounts</p>
          </Link>
          
          <button
            onClick={() => setShowCategoryModal(true)}
            className="p-4 border-2 border-dashed border-gray-300 dark:border-gray-600 rounded-lg text-center hover:border-green-400 hover:bg-green-50 dark:hover:bg-green-900/20 transition-colors"
          >
            <Package className="h-8 w-8 text-gray-400 dark:text-gray-500 mx-auto mb-2" />
            <p className="font-medium text-gray-900 dark:text-white">Add Category</p>
            <p className="text-sm text-gray-600 dark:text-gray-400">Create new categories</p>
          </button>

          

          <Link
            to="/admin/roles"
            className="p-4 border-2 border-dashed border-gray-300 dark:border-gray-600 rounded-lg text-center hover:border-purple-400 hover:bg-purple-50 dark:hover:bg-purple-900/20 transition-colors"
          >
            <Shield className="h-8 w-8 text-gray-400 dark:text-gray-500 mx-auto mb-2" />
            <p className="font-medium text-gray-900 dark:text-white">Manage Roles</p>
            <p className="text-sm text-gray-600 dark:text-gray-400">Configure permissions</p>
          </Link>
          
          <Link
            to="/admin/reports"
            className="p-4 border-2 border-dashed border-gray-300 dark:border-gray-600 rounded-lg text-center hover:border-orange-400 hover:bg-orange-50 dark:hover:bg-orange-900/20 transition-colors"
          >
            <BarChart3 className="h-8 w-8 text-gray-400 dark:text-gray-500 mx-auto mb-2" />
            <p className="font-medium text-gray-900 dark:text-white">View Reports</p>
            <p className="text-sm text-gray-600 dark:text-gray-400">Analytics & insights</p>
          </Link>
        </div>
      </div>

      {/* System Overview */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* User Management */}
        <div className="bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700 p-6">
          <div className="flex items-center justify-between mb-4">
            <h2 className="text-lg font-semibold text-gray-900 dark:text-white">User Management</h2>
            <Link
              to="/admin/users"
              className="text-blue-600 dark:text-blue-400 hover:text-blue-700 dark:hover:text-blue-300 text-sm font-medium"
            >
              View All
            </Link>
          </div>
          
          <div className="space-y-3">
            {users.slice(0, 5).map((user) => (
              <div key={user.id} className="flex items-center justify-between p-3 bg-gray-50 dark:bg-gray-700 rounded-lg">
                <div>
                  <p className="font-medium text-gray-900 dark:text-white">{user.username}</p>
                  <p className="text-sm text-gray-600 dark:text-gray-300">{user.email}</p>
                </div>
                <div className="flex items-center space-x-2">
                  <div className="flex space-x-1">
                    {user.roles.map((role) => (
                      <span
                        key={role.id}
                        className={`inline-flex items-center px-2 py-1 rounded-full text-xs font-medium ${
                          role.name === 'admin' 
                            ? 'bg-red-100 dark:bg-red-900 text-red-800 dark:text-red-200'
                            : role.name === 'seller'
                            ? 'bg-blue-100 dark:bg-blue-900 text-blue-800 dark:text-blue-200'
                            : 'bg-gray-100 dark:bg-gray-700 text-gray-800 dark:text-gray-200'
                        }`}
                      >
                        {role.name}
                      </span>
                    ))}
                  </div>
                  <div className="flex space-x-1">
                    <button
                      onClick={() => handleDeleteUser(user.id)}
                      className="p-1 text-gray-400 hover:text-red-600 dark:hover:text-red-400"
                    >
                      <Trash2 className="h-4 w-4" />
                    </button>
                  </div>
                </div>
              </div>
            ))}
            {users.length === 0 && (
              <p className="text-gray-500 dark:text-gray-400 text-center py-4">No users found</p>
            )}
          </div>
        </div>

        {/* System Status */}
        <div className="bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700 p-6">
          <div className="flex items-center justify-between mb-4">
            <h2 className="text-lg font-semibold text-gray-900 dark:text-white">System Status</h2>
            <Shield className="h-5 w-5 text-green-600 dark:text-green-400" />
          </div>
          
          <div className="space-y-3">
            <div className="flex items-center justify-between p-3 bg-green-50 dark:bg-green-900/20 rounded-lg">
              <div className="flex items-center">
                <Database className="h-4 w-4 text-green-600 dark:text-green-400 mr-2" />
                <div>
                  <p className="font-medium text-gray-900 dark:text-white">Database</p>
                  <p className="text-sm text-gray-600 dark:text-gray-300">All systems operational</p>
                </div>
              </div>
              <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 dark:bg-green-900 text-green-800 dark:text-green-200">
                Online
              </span>
            </div>
            
            <div className="flex items-center justify-between p-3 bg-green-50 dark:bg-green-900/20 rounded-lg">
              <div className="flex items-center">
                <Server className="h-4 w-4 text-green-600 dark:text-green-400 mr-2" />
                <div>
                  <p className="font-medium text-gray-900 dark:text-white">API Services</p>
                  <p className="text-sm text-gray-600 dark:text-gray-300">Running smoothly</p>
                </div>
              </div>
              <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 dark:bg-green-900 text-green-800 dark:text-green-200">
                Online
              </span>
            </div>
            
            <div className="flex items-center justify-between p-3 bg-blue-50 dark:bg-blue-900/20 rounded-lg">
              <div className="flex items-center">
                <Activity className="h-4 w-4 text-blue-600 dark:text-blue-400 mr-2" />
                <div>
                  <p className="font-medium text-gray-900 dark:text-white">Orders</p>
                  <p className="text-sm text-gray-600 dark:text-gray-300">{completedOrders} completed, {pendingOrders} pending</p>
                </div>
              </div>
              <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-blue-100 dark:bg-blue-900 text-blue-800 dark:text-blue-200">
                Active
              </span>
            </div>
          </div>
        </div>
      </div>

      {/* Recent Activity */}
      <div className="bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700 p-6">
        <div className="flex items-center justify-between mb-6">
          <h2 className="text-lg font-semibold text-gray-900 dark:text-white">Recent System Activity</h2>
          <Link
            to="/admin/orders"
            className="text-blue-600 dark:text-blue-400 hover:text-blue-700 dark:hover:text-blue-300 text-sm font-medium"
          >
            View All Orders
          </Link>
        </div>

        <div className="space-y-3">
          {orders.slice(0, 5).map((order) => (
            <div key={order.id} className="flex items-center justify-between p-3 bg-gray-50 dark:bg-gray-700 rounded-lg">
              <div className="flex items-center space-x-3">
                <div className="bg-blue-100 dark:bg-blue-900 rounded-full p-2">
                  <Package className="h-4 w-4 text-blue-600 dark:text-blue-400" />
                </div>
                <div>
                  <p className="font-medium text-gray-900 dark:text-white">Order #{order.id}</p>
                  <p className="text-sm text-gray-600 dark:text-gray-300">
                    {products.find(p => p.id === order.product_id)?.title || 'Product'} - {order.status}
                  </p>
                </div>
              </div>
              <div className="flex items-center space-x-2">
                <span className="text-sm text-gray-500 dark:text-gray-400">
                  {new Date(order.created_at).toLocaleDateString()}
                </span>
                <span className={`inline-flex items-center px-2 py-1 rounded-full text-xs font-medium ${
                  order.status === 'completed' 
                    ? 'bg-green-100 dark:bg-green-900 text-green-800 dark:text-green-200'
                    : order.status === 'pending'
                    ? 'bg-yellow-100 dark:bg-yellow-900 text-yellow-800 dark:text-yellow-200'
                    : 'bg-gray-100 dark:bg-gray-700 text-gray-800 dark:text-gray-200'
                }`}>
                  {order.status}
                </span>
              </div>
            </div>
          ))}
          {orders.length === 0 && (
            <p className="text-gray-500 dark:text-gray-400 text-center py-4">No recent activity</p>
          )}
        </div>
      </div>

      {/* Modals */}
      {/* Category Modal */}
      {showCategoryModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <div className="bg-white dark:bg-gray-800 rounded-lg p-6 w-full max-w-md">
            <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">Create New Category</h3>
            <div className="space-y-4">
              <input
                type="text"
                placeholder="Category name"
                value={newCategory.name}
                onChange={(e) => setNewCategory({...newCategory, name: e.target.value})}
                className="w-full p-2 border border-gray-300 dark:border-gray-600 rounded-lg dark:bg-gray-700 dark:text-white"
              />
              <textarea
                placeholder="Description"
                value={newCategory.description}
                onChange={(e) => setNewCategory({...newCategory, description: e.target.value})}
                className="w-full p-2 border border-gray-300 dark:border-gray-600 rounded-lg dark:bg-gray-700 dark:text-white"
                rows={3}
              />
              <div className="flex space-x-3">
                <button
                  onClick={handleCreateCategory}
                  className="bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700"
                >
                  Create
                </button>
                <button
                  onClick={() => setShowCategoryModal(false)}
                  className="bg-gray-600 text-white px-4 py-2 rounded-lg hover:bg-gray-700"
                >
                  Cancel
                </button>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default AdminDashboardPage;

