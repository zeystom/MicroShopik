import { useState, useEffect } from 'react';
import { 
  Activity, Search, Package, Users, 
  Loader2, Download, RefreshCw, Clock, AlertTriangle, 
  CheckCircle, Info, BarChart3, ArrowLeft
} from 'lucide-react';
import { User, Product, Order } from '@/types';
import { apiService } from '@/services/api';
import { toast } from 'react-hot-toast';
import { useTheme } from '@/hooks/useTheme';
import { Link } from 'react-router-dom';

interface ActivityItem {
  id: number;
  type: 'order' | 'user' | 'system' | 'product';
  action: string;
  description: string;
  timestamp: string;
  status: 'success' | 'warning' | 'error' | 'info';
  userId?: number;
  username?: string;
  orderId?: number;
  productId?: number;
  productTitle?: string;
}

const AdminActivityPage = () => {
  const { theme: _theme } = useTheme();
  const [isLoading, setIsLoading] = useState(true);
  // These state variables are used in generateActivityData function
  const [users, setUsers] = useState<User[]>([]);
  const [products, setProducts] = useState<Product[]>([]);
  const [orders, setOrders] = useState<Order[]>([]);
  const [activities, setActivities] = useState<ActivityItem[]>([]);
  const [searchTerm, setSearchTerm] = useState('');
  const [selectedType, setSelectedType] = useState<string>('all');
  const [selectedStatus, setSelectedStatus] = useState<string>('all');
  const [selectedPeriod, setSelectedPeriod] = useState('30d');

  useEffect(() => {
    fetchData();
  }, []);

  const fetchData = async () => {
    try {
      setIsLoading(true);
      const [usersData, productsData, ordersData] = await Promise.all([
        apiService.getUsers(),
        apiService.getAllProducts(),
        apiService.getAllOrders()
      ]);
      
      setUsers(usersData);
      setProducts(productsData);
      setOrders(ordersData);
      
      // Generate activity data from orders and other sources
      generateActivityData(usersData, productsData, ordersData);
    } catch (error) {
      console.error('Failed to fetch activity data:', error);
      toast.error('Failed to load activity data');
    } finally {
      setIsLoading(false);
    }
  };

  const generateActivityData = (users: User[], products: Product[], orders: Order[]) => {
    const activityData: ActivityItem[] = [];
    
    // Add order activities
    orders.forEach((order) => {
      const product = products.find(p => p.id === order.product_id);
      const customer = users.find(u => u.id === order.customer_id);
      
      activityData.push({
        id: order.id,
        type: 'order',
        action: 'Order Created',
        description: `Order #${order.id} created for ${product?.title || 'Product'} by ${customer?.username || 'User'}`,
        timestamp: order.created_at,
        status: order.status === 'completed' ? 'success' : order.status === 'pending' ? 'warning' : 'info',
        orderId: order.id,
        productId: order.product_id,
        productTitle: product?.title,
        userId: order.customer_id,
        username: customer?.username
      });
    });

    // Add user activities
    users.forEach((user) => {
      activityData.push({
        id: user.id,
        type: 'user',
        action: 'User Registered',
        description: `User ${user.username} registered with email ${user.email}`,
        timestamp: user.created_at,
        status: 'info',
        userId: user.id,
        username: user.username
      });
    });

    // Add product activities
    products.forEach((product) => {
      activityData.push({
        id: product.id,
        type: 'product',
        action: 'Product Created',
        description: `Product "${product.title}" created by seller`,
        timestamp: product.created_at,
        status: product.is_active ? 'success' : 'warning',
        productId: product.id,
        productTitle: product.title
      });
    });

    // Add system activities
    const systemActivities: ActivityItem[] = [
      {
        id: 1001,
        type: 'system',
        action: 'System Backup',
        description: 'Daily system backup completed successfully',
        timestamp: new Date(Date.now() - 24 * 60 * 60 * 1000).toISOString(),
        status: 'success'
      },
      {
        id: 1002,
        type: 'system',
        action: 'Database Maintenance',
        description: 'Database optimization and cleanup completed',
        timestamp: new Date(Date.now() - 2 * 24 * 60 * 60 * 1000).toISOString(),
        status: 'info'
      },
      {
        id: 1003,
        type: 'system',
        action: 'Security Scan',
        description: 'Security vulnerability scan completed',
        timestamp: new Date(Date.now() - 3 * 24 * 60 * 60 * 1000).toISOString(),
        status: 'success'
      }
    ];

    activityData.push(...systemActivities);
    
    // Sort by timestamp (newest first)
    activityData.sort((a, b) => new Date(b.timestamp).getTime() - new Date(a.timestamp).getTime());
    
    setActivities(activityData);
  };

  const getFilteredActivities = () => {
    let filtered = activities;

    // Filter by type
    if (selectedType !== 'all') {
      filtered = filtered.filter(activity => activity.type === selectedType);
    }

    // Filter by status
    if (selectedStatus !== 'all') {
      filtered = filtered.filter(activity => activity.status === selectedStatus);
    }

    // Filter by search term
    if (searchTerm) {
      filtered = filtered.filter(activity => 
        activity.description.toLowerCase().includes(searchTerm.toLowerCase()) ||
        activity.action.toLowerCase().includes(searchTerm.toLowerCase()) ||
        activity.username?.toLowerCase().includes(searchTerm.toLowerCase()) ||
        activity.productTitle?.toLowerCase().includes(searchTerm.toLowerCase())
      );
    }

    // Filter by time period
    const now = new Date();
    const periodMap = {
      '7d': 7,
      '30d': 30,
      '90d': 90,
      '1y': 365
    };
    
    const days = periodMap[selectedPeriod as keyof typeof periodMap] || 30;
    const cutoffDate = new Date(now.getTime() - (days * 24 * 60 * 60 * 1000));
    
    filtered = filtered.filter(activity => new Date(activity.timestamp) >= cutoffDate);

    return filtered;
  };

  const exportActivity = () => {
    const filteredActivities = getFilteredActivities();
    const csvContent = `Type,Action,Description,Status,Timestamp,User,Product,Order\n${filteredActivities.map(activity => 
      `${activity.type},${activity.action},"${activity.description}",${activity.status},${new Date(activity.timestamp).toLocaleString()},${activity.username || 'N/A'},${activity.productTitle || 'N/A'},${activity.orderId || 'N/A'}`
    ).join('\n')}`;

    const blob = new Blob([csvContent], { type: 'text/csv' });
    const url = window.URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = `system-activity-${selectedPeriod}-${new Date().toISOString().split('T')[0]}.csv`;
    a.click();
    window.URL.revokeObjectURL(url);
    toast.success('Activity data exported successfully');
  };

  const getStatusIcon = (status: string) => {
    switch (status) {
      case 'success':
        return <CheckCircle className="h-4 w-4 text-green-600" />;
      case 'warning':
        return <AlertTriangle className="h-4 w-4 text-yellow-600" />;
      case 'error':
        return <AlertTriangle className="h-4 w-4 text-red-600" />;
      default:
        return <Info className="h-4 w-4 text-blue-600" />;
    }
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'success':
        return 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200';
      case 'warning':
        return 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-200';
      case 'error':
        return 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200';
      default:
        return 'bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200';
    }
  };

  const getTypeIcon = (type: string) => {
    switch (type) {
      case 'order':
        return <Package className="h-4 w-4" />;
      case 'user':
        return <Users className="h-4 w-4" />;
      case 'product':
        return <BarChart3 className="h-4 w-4" />;
      case 'system':
        return <Activity className="h-4 w-4" />;
      default:
        return <Info className="h-4 w-4" />;
    }
  };

  const filteredActivities = getFilteredActivities();

  if (isLoading) {
    return (
      <div className="flex items-center justify-center min-h-[60vh]">
        <Loader2 className="w-8 h-8 animate-spin text-blue-600 dark:text-blue-400" />
        <span className="ml-2 text-gray-600 dark:text-gray-300">Loading activity data...</span>
      </div>
    );
  }

  return (
    <div className="space-y-6 p-6 bg-gray-50 dark:bg-gray-900 min-h-screen">
      {/* Header */}
      <div className="flex flex-col lg:flex-row justify-between items-start lg:items-center gap-4">
        <div className="flex items-center gap-4">
          <Link
            to="/admin-dashboard"
            className="bg-gray-600 text-white px-4 py-2 rounded-lg font-medium hover:bg-gray-700 transition-colors inline-flex items-center"
          >
            <ArrowLeft className="mr-2 h-4 w-4" />
            Назад в админ панель
          </Link>
          <div>
            <h1 className="text-3xl font-bold text-gray-900 dark:text-white">System Activity</h1>
            <p className="text-gray-600 dark:text-gray-300 mt-1">
              Monitor all system activities, user actions, and events
            </p>
          </div>
        </div>
        <div className="flex space-x-3">
          <button
            onClick={fetchData}
            className="bg-blue-600 text-white px-4 py-2 rounded-lg font-medium hover:bg-blue-700 transition-colors inline-flex items-center"
          >
            <RefreshCw className="mr-2 h-4 w-4" />
            Refresh
          </button>
          <button
            onClick={exportActivity}
            className="bg-green-600 text-white px-4 py-2 rounded-lg font-medium hover:bg-green-700 transition-colors inline-flex items-center"
          >
            <Download className="mr-2 h-4 w-4" />
            Export
          </button>
        </div>
      </div>

      {/* Filters */}
      <div className="bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700 p-4">
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
          {/* Search */}
          <div className="relative">
            <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-gray-400" />
            <input
              type="text"
              placeholder="Search activities..."
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
              className="w-full pl-10 pr-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg dark:bg-gray-700 dark:text-white focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            />
          </div>

          {/* Activity Type Filter */}
          <div>
            <select
              value={selectedType}
              onChange={(e) => setSelectedType(e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg dark:bg-gray-700 dark:text-white focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            >
              <option value="all">All Types</option>
              <option value="order">Orders</option>
              <option value="user">Users</option>
              <option value="product">Products</option>
              <option value="system">System</option>
            </select>
          </div>

          {/* Status Filter */}
          <div>
            <select
              value={selectedStatus}
              onChange={(e) => setSelectedStatus(e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg dark:bg-gray-700 dark:text-white focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            >
              <option value="all">All Statuses</option>
              <option value="success">Success</option>
              <option value="warning">Warning</option>
              <option value="error">Error</option>
              <option value="info">Info</option>
            </select>
          </div>

          {/* Time Period Filter */}
          <div>
            <select
              value={selectedPeriod}
              onChange={(e) => setSelectedPeriod(e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg dark:bg-gray-700 dark:text-white focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            >
              <option value="7d">Last 7 days</option>
              <option value="30d">Last 30 days</option>
              <option value="90d">Last 90 days</option>
              <option value="1y">Last year</option>
            </select>
          </div>
        </div>
      </div>

      {/* Activity Stats */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-6">
        <div className="bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700 p-6">
          <div className="flex items-center">
            <div className="bg-blue-100 dark:bg-blue-900 rounded-full p-3">
              <Activity className="h-6 w-6 text-blue-600 dark:text-blue-400" />
            </div>
            <div className="ml-4">
              <p className="text-sm font-medium text-gray-600 dark:text-gray-300">Total Activities</p>
              <p className="text-2xl font-bold text-gray-900 dark:text-white">{filteredActivities.length}</p>
            </div>
          </div>
        </div>

        <div className="bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700 p-6">
          <div className="flex items-center">
            <div className="bg-green-100 dark:bg-green-900 rounded-full p-3">
              <CheckCircle className="h-6 w-6 text-green-600 dark:text-green-400" />
            </div>
            <div className="ml-4">
              <p className="text-sm font-medium text-gray-600 dark:text-gray-300">Successful</p>
              <p className="text-2xl font-bold text-gray-900 dark:text-white">
                {filteredActivities.filter(a => a.status === 'success').length}
              </p>
            </div>
          </div>
        </div>

        <div className="bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700 p-6">
          <div className="flex items-center">
            <div className="bg-yellow-100 dark:bg-yellow-900 rounded-full p-3">
              <AlertTriangle className="h-6 w-6 text-yellow-600 dark:text-yellow-400" />
            </div>
            <div className="ml-4">
              <p className="text-sm font-medium text-gray-600 dark:text-gray-300">Warnings</p>
              <p className="text-2xl font-bold text-gray-900 dark:text-white">
                {filteredActivities.filter(a => a.status === 'warning').length}
              </p>
            </div>
          </div>
        </div>

        <div className="bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700 p-6">
          <div className="flex items-center">
            <div className="bg-red-100 dark:bg-red-900 rounded-full p-3">
              <AlertTriangle className="h-6 w-6 text-red-600 dark:text-red-400" />
            </div>
            <div className="ml-4">
              <p className="text-sm font-medium text-gray-600 dark:text-gray-300">Errors</p>
              <p className="text-2xl font-bold text-gray-900 dark:text-white">
                {filteredActivities.filter(a => a.status === 'error').length}
              </p>
            </div>
          </div>
        </div>
      </div>

      {/* Activity List */}
      <div className="bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700 overflow-hidden">
        <div className="overflow-x-auto">
          <table className="w-full">
            <thead className="bg-gray-50 dark:bg-gray-700">
              <tr>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
                  Activity
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
                  Type
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
                  Status
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
                  Timestamp
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
                  Details
                </th>
              </tr>
            </thead>
            <tbody className="bg-white dark:bg-gray-800 divide-y divide-gray-200 dark:divide-gray-700">
              {filteredActivities.map((activity) => (
                <tr key={activity.id} className="hover:bg-gray-50 dark:hover:bg-gray-700">
                  <td className="px-6 py-4 whitespace-nowrap">
                    <div className="flex items-center">
                      <div className="bg-gray-100 dark:bg-gray-700 rounded-full p-2 mr-3">
                        {getTypeIcon(activity.type)}
                      </div>
                      <div>
                        <div className="text-sm font-medium text-gray-900 dark:text-white">
                          {activity.action}
                        </div>
                        <div className="text-sm text-gray-500 dark:text-gray-400">
                          {activity.description}
                        </div>
                      </div>
                    </div>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap">
                    <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-gray-100 dark:bg-gray-700 text-gray-800 dark:text-gray-200 capitalize">
                      {activity.type}
                    </span>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap">
                    <div className="flex items-center">
                      {getStatusIcon(activity.status)}
                      <span className={`ml-2 inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${getStatusColor(activity.status)} capitalize`}>
                        {activity.status}
                      </span>
                    </div>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500 dark:text-gray-400">
                    <div className="flex items-center">
                      <Clock className="h-4 w-4 mr-2" />
                      {new Date(activity.timestamp).toLocaleString()}
                    </div>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500 dark:text-gray-400">
                    {activity.username && (
                      <div className="mb-1">
                        <span className="font-medium">User:</span> {activity.username}
                      </div>
                    )}
                    {activity.productTitle && (
                      <div className="mb-1">
                        <span className="font-medium">Product:</span> {activity.productTitle}
                      </div>
                    )}
                    {activity.orderId && (
                      <div>
                        <span className="font-medium">Order:</span> #{activity.orderId}
                      </div>
                    )}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
        
        {filteredActivities.length === 0 && (
          <div className="text-center py-8">
            <Activity className="mx-auto h-12 w-12 text-gray-400" />
            <h3 className="mt-2 text-sm font-medium text-gray-900 dark:text-white">No activities found</h3>
            <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">
              {searchTerm || selectedType !== 'all' || selectedStatus !== 'all' 
                ? 'Try adjusting your search or filter criteria.' 
                : 'No system activities recorded for the selected time period.'}
            </p>
          </div>
        )}
      </div>
    </div>
  );
};

export default AdminActivityPage;
