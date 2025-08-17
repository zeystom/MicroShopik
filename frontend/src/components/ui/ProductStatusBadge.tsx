import React from 'react';

interface ProductStatusBadgeProps {
  isActive: boolean;
  className?: string;
}

const ProductStatusBadge: React.FC<ProductStatusBadgeProps> = ({ isActive, className = '' }) => {
  if (isActive) {
    return null; // Не показываем бейдж для активных продуктов
  }

  return (
    <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200 ${className}`}>
      Inactive
    </span>
  );
};

export default ProductStatusBadge;
