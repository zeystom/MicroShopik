import React from 'react'
import { Routes, Route } from 'react-router-dom'
import { Toaster } from 'react-hot-toast'
import Layout from '@/components/common/Layout'
import ProtectedRoute from '@/components/common/ProtectedRoute'

// Public pages
import LoginPage from '@/pages/LoginPage'
import RegisterPage from '@/pages/RegisterPage'
import HomePage from '@/pages/HomePage'
import ProductsPage from '@/pages/ProductsPage'
import ProductDetailPage from '@/pages/ProductDetailPage'
import NotFoundPage from '@/pages/NotFoundPage'

// Protected pages
import ProfilePage from '@/pages/ProfilePage'
import OrdersPage from '@/pages/OrdersPage'
import OrderDetailPage from '@/pages/OrderDetailPage'
import SellerDashboardPage from '@/pages/SellerDashboardPage'
import AdminDashboardPage from '@/pages/AdminDashboardPage'
import AdminUsersPage from '@/pages/AdminUsersPage'
import AdminRolesPage from '@/pages/AdminRolesPage'
import AdminSettingsPage from '@/pages/AdminSettingsPage'
import AdminReportsPage from '@/pages/AdminReportsPage'
import AdminCategoriesPage from '@/pages/AdminCategoriesPage'
import AdminOrdersPage from '@/pages/AdminOrdersPage'

import ConversationsPage from '@/pages/ConversationsPage'
import CreateProductPage from '@/pages/CreateProductPage'
import ApiTestPage from '@/pages/ApiTestPage'
import SellerProductsPage from '@/pages/SellerProductsPage'
import EditProductPage from '@/pages/EditProductPage'
import SellerMessagesPage from '@/pages/SellerMessagesPage'

const App: React.FC = () => {
  return (
    <>
      <Toaster position="top-right" />
      <Routes>
        {/* Standalone public routes (no header) */}
        <Route path="/" element={<HomePage />} />
        <Route path="/login" element={<LoginPage />} />
        <Route path="/register" element={<RegisterPage />} />

        {/* Routes with layout (includes header) */}
        <Route path="/" element={<Layout />}>
          <Route path="/products" element={<ProductsPage />} />
          <Route path="/products/:id" element={<ProductDetailPage />} />
          
          {/* Protected routes */}
          <Route path="/profile" element={
            <ProtectedRoute>
              <ProfilePage />
            </ProtectedRoute>
          } />
          <Route path="/orders" element={
            <ProtectedRoute>
              <OrdersPage />
            </ProtectedRoute>
          } />
          <Route path="/orders/:id" element={
            <ProtectedRoute>
              <OrderDetailPage />
            </ProtectedRoute>
          } />
          <Route path="/conversations" element={
            <ProtectedRoute>
              <ConversationsPage />
            </ProtectedRoute>
          } />
          <Route path="/api-test" element={
            <ProtectedRoute requiredRoles={['admin']}>
              <ApiTestPage />
            </ProtectedRoute>
          } />
          
          {/* Seller routes */}
          <Route path="/seller-dashboard" element={
            <ProtectedRoute requiredRoles={['seller']}>
              <SellerDashboardPage />
            </ProtectedRoute>
          } />
          <Route path="/seller/messages" element={
            <ProtectedRoute requiredRoles={['seller']}>
              <SellerMessagesPage />
            </ProtectedRoute>
          } />
          <Route path="/seller/products" element={
            <ProtectedRoute requiredRoles={['seller']}>
              <SellerProductsPage />
            </ProtectedRoute>
          } />
          <Route path="/seller/products/new" element={
            <ProtectedRoute requiredRoles={['seller']}>
              <CreateProductPage />
            </ProtectedRoute>
          } />
          <Route path="/seller/products/:id/edit" element={
            <ProtectedRoute requiredRoles={['seller']}>
              <EditProductPage />
            </ProtectedRoute>
          } />
          
          {/* Admin routes */}
          <Route path="/admin-dashboard" element={
            <ProtectedRoute requiredRoles={['admin']}>
              <AdminDashboardPage />
            </ProtectedRoute>
          } />
          <Route path="/admin/users" element={
            <ProtectedRoute requiredRoles={['admin']}>
              <AdminUsersPage />
            </ProtectedRoute>
          } />
          <Route path="/admin/roles" element={
            <ProtectedRoute requiredRoles={['admin']}>
              <AdminRolesPage />
            </ProtectedRoute>
          } />
          <Route path="/admin/settings" element={
            <ProtectedRoute requiredRoles={['admin']}>
              <AdminSettingsPage />
            </ProtectedRoute>
          } />
          <Route path="/admin/reports" element={
            <ProtectedRoute requiredRoles={['admin']}>
              <AdminReportsPage />
            </ProtectedRoute>
          } />
          <Route path="/admin/categories" element={
            <ProtectedRoute requiredRoles={['admin']}>
              <AdminCategoriesPage />
            </ProtectedRoute>
          } />
          <Route path="/admin/orders" element={
            <ProtectedRoute requiredRoles={['admin']}>
              <AdminOrdersPage />
            </ProtectedRoute>
          } />

        </Route>

        {/* 404 page */}
        <Route path="*" element={<NotFoundPage />} />
      </Routes>
    </>
  )
}

export default App

