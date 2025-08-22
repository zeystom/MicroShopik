import axios, { AxiosInstance, AxiosResponse } from 'axios';
import { useAuthStore } from '@/store/authStore';
import { config } from '@/config';
import { Product, Order, Conversation, User, Role, Category } from '@/types';

class ApiService {
  private api: AxiosInstance;

  constructor() {
    this.api = axios.create({
      baseURL: config.apiUrl,
      headers: {
        'Content-Type': 'application/json',
      },
    });

    // Request interceptor to add auth token
    this.api.interceptors.request.use(
      (config) => {
        const token = useAuthStore.getState().token;
        if (token) {
          config.headers.Authorization = `Bearer ${token}`;
        }
        
        // Log request data for debugging
        if (config.method === 'post' && config.url === '/products') {
          console.log('Axios interceptor - Request data:', config.data);
          console.log('Axios interceptor - Request URL:', config.url);
          console.log('Axios interceptor - Request method:', config.method);
        }
        
        return config;
      },
      (error) => {
        return Promise.reject(error);
      }
    );

    // Response interceptor to handle auth errors
    this.api.interceptors.response.use(
      (response: AxiosResponse) => response,
      (error) => {
        if (error.response?.status === 401) {
          useAuthStore.getState().logout();
          // Use hash-based route to avoid server 404s on SPA
          window.location.href = '/#/login';
        }
        return Promise.reject(error);
      }
    );
  }

  // Normalize server message shape (content -> text)
  private mapMessageFromServer = (m: any) => {
    if (!m) return m;
    if (m.text == null && m.content != null) {
      return { ...m, text: m.content };
    }
    return m;
  };

  private mapConversationFromServer = (c: any) => {
    if (!c) return c;
    if (Array.isArray(c.messages)) {
      return { ...c, messages: c.messages.map(this.mapMessageFromServer) };
    }
    return c;
  };

  // Auth endpoints
  async login(credentials: { email: string; password: string }) {
    const response = await this.api.post('/auth/login', credentials);
    return response.data;
  }

  async register(userData: { username: string; email: string; password: string }) {
    const response = await this.api.post('/auth/register', userData);
    return response.data;
  }

  // Product endpoints
  async getProducts(params?: { category_id?: number; search?: string; page?: number; limit?: number; seller_id?: number; is_active?: boolean }) {
    console.log('API getProducts called with params:', params);
    try {
      // По умолчанию показываем только активные продукты для обычных пользователей
      const defaultParams = { is_active: true, ...params };
      const response = await this.api.get('/products', { params: defaultParams });
      console.log('API getProducts response:', response.data);
      return response.data;
    } catch (error: unknown) {
      if (typeof error === 'object' && error !== null && 'response' in error) {
        const e = error as { response?: { data?: unknown; status?: number } };
        console.error('API getProducts error:', e.response?.data);
        console.error('API getProducts error status:', e.response?.status);
      }
      throw error;
    }
  }

  // Метод для получения продуктов без дефолтного фильтра is_active
  async getProductsRaw(params?: { category_id?: number; search?: string; page?: number; limit?: number; seller_id?: number; is_active?: boolean }) {
    console.log('API getProductsRaw called with params:', params);
    try {
      const response = await this.api.get('/products', { params });
      console.log('API getProductsRaw response:', response.data);
      return response.data;
    } catch (error: unknown) {
      if (typeof error === 'object' && error !== null && 'response' in error) {
        const e = error as { response?: { data?: unknown; status?: number } };
        console.error('API getProductsRaw error:', e.response?.data);
        console.error('API getProductsRaw error status:', e.response?.status);
      }
      throw error;
    }
  }

  async getProduct(id: number) {
    const response = await this.api.get(`/products/${id}`);
    return response.data;
  }



  async createProduct(productData: Partial<Product>) {
    console.log('API createProduct called with:', productData);
    console.log('API createProduct categoryId type:', typeof productData.category_id, 'value:', productData.category_id);
    console.log('API createProduct full data:', JSON.stringify(productData, null, 2));
    
    try {
      const response = await this.api.post('/products', productData);
      console.log('API createProduct response:', response.data);
      return response.data;
    } catch (error: unknown) {
      if (typeof error === 'object' && error !== null && 'response' in error) {
        const e = error as { response?: { data?: unknown; status?: number; headers?: unknown } };
        console.error('API createProduct error response:', e.response?.data);
        console.error('API createProduct error status:', e.response?.status);
        console.error('API createProduct error headers:', e.response?.headers);
      }
      throw error;
    }
  }

  async updateProduct(id: number, productData: Partial<Product>) {
    const response = await this.api.put(`/products/${id}`, productData);
    return response.data;
  }

  async deleteProduct(id: number) {
    const response = await this.api.delete(`/products/${id}`);
    return response.data;
  }

  // Category endpoints
  async getCategories() {
    console.log('API getCategories called');
    const response = await this.api.get('/categories');
    console.log('API getCategories response:', response.data);
    return response.data;
  }

  async getCategory(id: number) {
    const response = await this.api.get(`/categories/${id}`);
    return response.data;
  }

  // Order endpoints
  async createOrder(orderData: Partial<Order>) {
    const response = await this.api.post('/orders', orderData);
    return response.data;
  }

  async getOrders() {
    console.log('API getOrders called');
    try {
      const response = await this.api.get('/orders');
      console.log('API getOrders response:', response.data);
      return response.data;
    } catch (error: unknown) {
      if (typeof error === 'object' && error !== null && 'response' in error) {
        const e = error as { response?: { data?: unknown; status?: number } };
        console.error('API getOrders error:', e.response?.data);
        console.error('API getOrders error status:', e.response?.status);
      }
      throw error;
    }
  }

  async getSellerOrders() {
    console.log('API getSellerOrders called');
    try {
      const response = await this.api.get('/orders/seller');
      console.log('API getSellerOrders response:', response.data);
      return response.data;
    } catch (error: unknown) {
      if (typeof error === 'object' && error !== null && 'response' in error) {
        const e = error as { response?: { data?: unknown; status?: number } };
        console.error('API getSellerOrders error:', e.response?.data);
        console.error('API getSellerOrders error status:', e.response?.status);
      }
      throw error;
    }
  }

  // Метод для получения всех продуктов продавца (включая неактивные)
  async getSellerProducts(sellerId: number) {
    console.log('API getSellerProducts called for seller ID:', sellerId);
    try {
      // Используем getProductsRaw для получения всех продуктов продавца без фильтра is_active
      const response = await this.getProductsRaw({ seller_id: sellerId });
      console.log('API getSellerProducts response:', response);
      return response;
    } catch (error: unknown) {
      if (typeof error === 'object' && error !== null && 'response' in error) {
        const e = error as { response?: { data?: unknown; status?: number } };
        console.error('API getSellerProducts error:', e.response?.data);
        console.error('API getSellerProducts error status:', e.response?.status);
      }
      throw error;
    }
  }

  async getOrder(id: number) {
    const response = await this.api.get(`/orders/${id}`);
    return response.data;
  }

  async processOrder(id: number) {
    const response = await this.api.post(`/orders/${id}/process`);
    return response.data;
  }

  async cancelOrder(id: number, customerId: number) {
    const response = await this.api.post(`/orders/${id}/cancel/${customerId}`);
    return response.data;
  }

  async confirmOrder(id: number) {
    const response = await 
    this.api.post(`/orders/${id}/confirm`);
    return response.data;
  }

  // Conversation endpoints
  async getConversations() {
    const currentUserId = useAuthStore.getState().user?.id;
    if (!currentUserId) throw new Error('User not authenticated');
    const response = await this.api.get(`/users/${currentUserId}/conversations`);
    const data = response.data;
    if (Array.isArray(data)) {
      return data.map(this.mapConversationFromServer);
    }
    return data;
  }

  async getConversation(id: number) {
    const response = await this.api.get(`/conversations/${id}`);
    const data = response.data;
    // Some endpoints return { conversation, messages }
    if (data && data.conversation && data.messages) {
      return {
        conversation: this.mapConversationFromServer(data.conversation),
        messages: Array.isArray(data.messages) ? data.messages.map(this.mapMessageFromServer) : data.messages,
      };
    }
    return this.mapConversationFromServer(data);
  }

  async createConversation(conversationData: Partial<Conversation>, participantIds: number[]) {
    const response = await this.api.post('/conversations', {
      conversation: conversationData,
      participant_ids: participantIds
    });
    return response.data;
  }

  // Message endpoints
  async getMessages(conversationId: number) {
    const response = await this.api.get(`/conversations/${conversationId}/messages`);
    const data = response.data;
    return Array.isArray(data) ? data.map(this.mapMessageFromServer) : data;
  }

  async sendMessage(conversationId: number, messageData: { text: string; sender_id: number; order_id?: number }) {
    const response = await this.api.post(`/conversations/${conversationId}/messages`, {
      conversation_id: conversationId,
      content: messageData.text, // backend expects `content`
      sender_id: messageData.sender_id,
      order_id: messageData.order_id,
    });
    return response.data;
  }

  async sendSystemMessage(conversationId: number, text: string, orderId?: number) {
    const response = await this.api.post(`/conversations/${conversationId}/messages/system`, {
      text,
      order_id: orderId ?? null,
    });
    return response.data;
  }

  // User endpoints
  async getProfile() {
    const response = await this.api.get('/users/profile');
    return response.data;
  }

  async updateProfile(userData: Partial<User>) {
    const response = await this.api.put('/users/profile', userData);
    return response.data;
  }

  // Role management endpoints
  async getRoles() {
    const response = await this.api.get('/roles');
    return response.data;
  }

  async createRole(roleData: Partial<Role>) {
    const response = await this.api.post('/roles', roleData);
    return response.data;
  }

  async getUserRoles(userId: number) {
    const response = await this.api.get(`/users/${userId}/roles`);
    return response.data;
  }

  // Admin endpoints
  async getUsers() {
    const response = await this.api.get('/admin/users');
    return response.data;
  }

  async assignRole(userId: number, roleName: string) {
    const response = await this.api.post(`/users/${userId}/roles/${roleName}`);
    return response.data;
  }

  // Admin user management
  async deleteUser(userId: number) {
    const response = await this.api.delete(`/admin/users/${userId}`);
    return response.data;
  }

  async updateUser(userId: number, userData: Partial<User>) {
    const response = await this.api.put(`/admin/users/${userId}`, userData);
    return response.data;
  }

  async getUserById(userId: number) {
    const response = await this.api.get(`/admin/users/${userId}`);
    return response.data;
  }

  // Admin category management
  async createCategory(categoryData: Partial<Category>) {
    const response = await this.api.post('/categories', categoryData);
    return response.data;
  }

  async updateCategory(categoryId: number, categoryData: Partial<Category>) {
    const response = await this.api.put(`/categories/${categoryId}`, categoryData);
    return response.data;
  }

  async deleteCategory(categoryId: number) {
    const response = await this.api.delete(`/categories/${categoryId}`);
    return response.data;
  }

  // Admin role management
  async updateRole(roleId: number, roleData: Partial<Role>) {
    const response = await this.api.put(`/roles/${roleId}`, roleData);
    return response.data;
  }

  async deleteRole(roleId: number) {
    const response = await this.api.delete(`/roles/${roleId}`);
    return response.data;
  }

  async removeRoleFromUser(userId: number, roleName: string) {
    const response = await this.api.delete(`/users/${userId}/roles/${roleName}`);
    return response.data;
  }

  // Admin system management
  async getSystemStats() {
    const response = await this.api.get('/admin/stats');
    return response.data;
  }

  async getSystemLogs() {
    const response = await this.api.get('/admin/logs');
    return response.data;
  }

  async updateSystemSettings(settings: Partial<{ maintenanceMode: boolean; allowRegistration: boolean; maxProductsPerUser: number; maxUsersPerRole: number; systemNotifications: boolean; autoBackup: boolean; backupFrequency: string; logLevel: string }>) {
    const response = await this.api.put('/admin/settings', settings);
    return response.data;
  }

  // Admin order management
  async getAllOrders() {
    const response = await this.api.get('/admin/orders');
    return response.data;
  }

  async updateOrderStatus(orderId: number, status: string) {
    const response = await this.api.put(`/admin/orders/${orderId}/status`, { status });
    return response.data;
  }

  async deleteOrder(orderId: number) {
    const response = await this.api.delete(`/admin/orders/${orderId}`);
    return response.data;
  }

  // Admin product management
  async getAllProducts(params?: { category_id?: number; search?: string; page?: number; limit?: number; seller_id?: number }) {
    console.log('API getAllProducts called with params:', params);
    try {
      const response = await this.api.get('/admin/products', { params });
      console.log('API getAllProducts response:', response.data);
      return response.data;
    } catch (error: unknown) {
      if (typeof error === 'object' && error !== null && 'response' in error) {
        const e = error as { response?: { data?: unknown; status?: number } };
        console.error('API getAllProducts error:', e.response?.data);
        console.error('API getAllProducts error status:', e.response?.status);
      }
      throw error;
    }
  }

  async updateProductStatus(productId: number, isActive: boolean) {
    const response = await this.api.put(`/admin/products/${productId}/status`, { is_active: isActive });
    return response.data;
  }

  async deleteProductAdmin(productId: number) {
    const response = await this.api.delete(`/admin/products/${productId}`);
    return response.data;
  }

  // Health check
  async healthCheck() {
    const response = await this.api.get('/health');
    return response.data;
  }
}

export const apiService = new ApiService();
export default apiService;

