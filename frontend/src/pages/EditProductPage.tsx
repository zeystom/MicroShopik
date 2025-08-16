import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { ArrowLeft, Save, Loader2, AlertCircle } from 'lucide-react';
import { Product, Category } from '@/types';
import { apiService } from '@/services/api';
import { useAuthStore } from '@/store/authStore';
import { useTheme } from '@/hooks/useTheme';
import { toast } from 'react-hot-toast';

const EditProductPage = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const { user } = useAuthStore();
  const { theme: _theme } = useTheme();
  
  const [product, setProduct] = useState<Product | null>(null);
  const [categories, setCategories] = useState<Category[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [isSaving, setIsSaving] = useState(false);
  const [isCategoriesLoading, setIsCategoriesLoading] = useState(true);
  
  const [formData, setFormData] = useState({
    title: '',
    description: '',
    price: '',
    category_id: '',
    max_sales: '',
    disposable: false,
    is_active: true
  });

  useEffect(() => {
    const fetchData = async () => {
      if (!id) return;
      
      try {
        setIsLoading(true);
        const [productData, categoriesData] = await Promise.all([
          apiService.getProduct(Number(id)),
          apiService.getCategories()
        ]);
        
        // Verify the product belongs to the current user
        if (productData.seller_id !== user?.id) {
          toast.error('You can only edit your own products');
          navigate('/seller/products');
          return;
        }
        
        setProduct(productData);
        setCategories(categoriesData);
        
        // Pre-fill form with existing data
        setFormData({
          title: productData.title,
          description: productData.description,
          price: (productData.price / 100).toFixed(2), // Convert from cents
          category_id: productData.category_id.toString(),
          max_sales: productData.max_sales.toString(),
          disposable: productData.disposable,
          is_active: productData.is_active
        });
      } catch (error: any) {
        console.error('Failed to fetch product:', error);
        toast.error('Failed to load product');
        navigate('/seller/products');
      } finally {
        setIsLoading(false);
        setIsCategoriesLoading(false);
      }
    };

    if (user?.id) {
      fetchData();
    }
  }, [id, user?.id, navigate]);

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement>) => {
    const { name, value, type } = e.target;
    
    if (type === 'checkbox') {
      const checked = (e.target as HTMLInputElement).checked;
      setFormData(prev => ({ ...prev, [name]: checked }));
    } else {
      setFormData(prev => ({ ...prev, [name]: value }));
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!product || !user?.id) return;

    // Validation
    if (!formData.title.trim() || !formData.description.trim() || !formData.price || !formData.category_id) {
      toast.error('Please fill in all required fields');
      return;
    }

    const price = parseFloat(formData.price);
    if (isNaN(price) || price <= 0) {
      toast.error('Please enter a valid price');
      return;
    }

    const maxSales = parseInt(formData.max_sales);
    if (isNaN(maxSales) || maxSales <= 0) {
      toast.error('Please enter a valid maximum sales number');
      return;
    }

    try {
      setIsSaving(true);
      
      const updateData = {
        title: formData.title.trim(),
        description: formData.description.trim(),
        price: Math.round(price * 100), // Convert to cents
        category_id: parseInt(formData.category_id),
        max_sales: maxSales,
        disposable: formData.disposable,
        is_active: formData.is_active
      };

      await apiService.updateProduct(product.id, updateData);
      
      toast.success('Product updated successfully!');
      navigate('/seller/products');
    } catch (error: any) {
      console.error('Failed to update product:', error);
      toast.error(error.response?.data?.error || 'Failed to update product');
    } finally {
      setIsSaving(false);
    }
  };

  if (isLoading) {
    return (
      <div className="flex items-center justify-center min-h-screen bg-gray-50 dark:bg-gray-950">
        <Loader2 className="w-8 h-8 animate-spin text-blue-600 dark:text-blue-400" />
        <span className="ml-2 text-gray-600 dark:text-gray-300">Loading product...</span>
      </div>
    );
  }

  if (!product) {
    return (
      <div className="flex items-center justify-center min-h-screen bg-gray-50 dark:bg-gray-950">
        <div className="text-center">
          <AlertCircle className="w-16 h-16 text-red-500 dark:text-red-400 mx-auto mb-4" />
          <h3 className="text-lg font-medium text-gray-900 dark:text-gray-100 mb-2">Product not found</h3>
          <p className="text-gray-600 dark:text-gray-400 mb-4">The product you're looking for doesn't exist or you don't have permission to edit it.</p>
          <button
            onClick={() => navigate('/seller/products')}
            className="btn-primary-dark"
          >
            Back to Products
          </button>
        </div>
      </div>
    );
  }

  if (isCategoriesLoading) {
    return (
      <div className="flex items-center justify-center min-h-screen bg-gray-50 dark:bg-gray-950">
        <Loader2 className="w-8 h-8 animate-spin text-blue-600 dark:text-blue-400" />
        <span className="ml-2 text-gray-600 dark:text-gray-300">Loading categories...</span>
      </div>
    );
  }

  return (
    <div className="container mx-auto px-4 py-8 bg-gray-50 dark:bg-gray-950 min-h-screen">
      <div className="max-w-2xl mx-auto">
        {/* Header */}
        <div className="flex items-center justify-between mb-8">
          <div className="flex items-center">
            <button
              onClick={() => navigate('/seller/products')}
              className="mr-4 p-2 text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-gray-100 hover:bg-gray-100 dark:hover:bg-gray-800 rounded-lg transition-colors"
            >
              <ArrowLeft className="w-5 h-5" />
            </button>
            <div>
              <h1 className="text-3xl font-bold text-gray-900 dark:text-gray-100">Edit Product</h1>
              <p className="text-gray-600 dark:text-gray-400 mt-2">Update your product information</p>
            </div>
          </div>
        </div>

        {/* Form */}
        <form onSubmit={handleSubmit} className="bg-white dark:bg-gray-900 rounded-lg shadow-sm border border-gray-200 dark:border-gray-800 p-6">
          <div className="space-y-6">
            {/* Title */}
            <div>
              <label htmlFor="title" className="form-label">
                Product Title *
              </label>
              <input
                type="text"
                id="title"
                name="title"
                value={formData.title}
                onChange={handleInputChange}
                required
                className="form-input"
                placeholder="Enter product title"
              />
            </div>

            {/* Description */}
            <div>
              <label htmlFor="description" className="form-label">
                Description *
              </label>
              <textarea
                id="description"
                name="description"
                value={formData.description}
                onChange={handleInputChange}
                required
                rows={4}
                className="form-textarea"
                placeholder="Describe your product"
              />
            </div>

            {/* Price and Category */}
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div>
                <label htmlFor="price" className="form-label">
                  Price (USD) *
                </label>
                <input
                  type="number"
                  id="price"
                  name="price"
                  value={formData.price}
                  onChange={handleInputChange}
                  required
                  min="0"
                  step="0.01"
                  className="form-input"
                  placeholder="0.00"
                />
              </div>

              <div>
                <label htmlFor="category_id" className="form-label">
                  Category *
                </label>
                <select
                  id="category_id"
                  name="category_id"
                  value={formData.category_id}
                  onChange={handleInputChange}
                  required
                  className="form-select"
                >
                  <option value="">Select a category</option>
                  {categories.map((category) => (
                    <option key={category.id} value={category.id}>
                      {category.name}
                    </option>
                  ))}
                </select>
              </div>
            </div>

            {/* Max Sales */}
            <div>
              <label htmlFor="max_sales" className="form-label">
                Maximum Sales *
              </label>
              <input
                type="number"
                id="max_sales"
                name="max_sales"
                value={formData.max_sales}
                onChange={handleInputChange}
                required
                min="1"
                className="form-input"
                placeholder="Enter maximum number of sales"
              />
            </div>

            {/* Checkboxes */}
            <div className="space-y-4">
              <div className="flex items-center">
                <input
                  type="checkbox"
                  id="disposable"
                  name="disposable"
                  checked={formData.disposable}
                  onChange={handleInputChange}
                  className="form-checkbox"
                />
                <label htmlFor="disposable" className="ml-2 block text-sm text-gray-900 dark:text-gray-100">
                  Disposable (one-time use)
                </label>
              </div>

              <div className="flex items-center">
                <input
                  type="checkbox"
                  id="is_active"
                  name="is_active"
                  checked={formData.is_active}
                  onChange={handleInputChange}
                  className="form-checkbox"
                />
                <label htmlFor="is_active" className="ml-2 block text-sm text-gray-900 dark:text-gray-100">
                  Active (available for purchase)
                </label>
              </div>
            </div>

            {/* Submit Button */}
            <div className="flex justify-end space-x-4 pt-6 border-t border-gray-200 dark:border-gray-700">
              <button
                type="button"
                onClick={() => navigate('/seller/products')}
                className="btn-secondary-dark"
              >
                Cancel
              </button>
              <button
                type="submit"
                disabled={isSaving}
                className="btn-primary-dark disabled:opacity-50 disabled:cursor-not-allowed inline-flex items-center"
              >
                {isSaving ? (
                  <>
                    <Loader2 className="w-4 h-4 mr-2 animate-spin" />
                    Saving...
                  </>
                ) : (
                  <>
                    <Save className="w-4 h-4 mr-2" />
                    Save Changes
                  </>
                )}
              </button>
            </div>
          </div>
        </form>
      </div>
    </div>
  );
};

export default EditProductPage;
