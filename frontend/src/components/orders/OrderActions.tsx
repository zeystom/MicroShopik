import React from 'react';
import { useNavigate } from 'react-router-dom';
import { Button } from '@/components/ui/button';
import { apiService } from '@/services/api';
import { Order } from '@/types';
import { toast } from 'react-hot-toast';

interface OrderActionsProps {
  order: Order;
  onOrderUpdated: () => void;
  userRole: string;
}

const OrderActions: React.FC<OrderActionsProps> = ({ order, onOrderUpdated, userRole }) => {
  const navigate = useNavigate();

  const handleProcessOrder = async () => {
    try {
      await apiService.processOrder(order.id);
      toast.success('Order processed successfully! Conversation created.');
      onOrderUpdated();
    } catch (error) {
      toast.error('Failed to process order');
      console.error('Error processing order:', error);
    }
  };

  const handleCancelOrder = async () => {
    if (!order.customer_id) {
      toast.error('Cannot cancel order: customer ID not found');
      return;
    }
    
    try {
      await apiService.cancelOrder(order.id, order.customer_id);
      toast.success('Order cancelled successfully!');
      onOrderUpdated();
    } catch (error) {
      toast.error('Failed to cancel order');
      console.error('Error cancelling order:', error);
    }
  };

  return (
    <div className="flex gap-2">
      {order.status === 'pending' && (
        <>
          {(userRole === 'admin' || userRole === 'seller') && (
            <Button
              onClick={handleProcessOrder}
              className="bg-green-600 hover:bg-green-700 dark:bg-green-500 dark:hover:bg-green-600 text-white transition-colors"
            >
              Process Order
            </Button>
          )}
          <Button
            onClick={handleCancelOrder}
            variant="destructive"
            className="dark:bg-red-500 dark:hover:bg-red-600 transition-colors"
          >
            Cancel Order
          </Button>
        </>
      )}
      
      {order.status === 'completed' && (
        <Button
          onClick={() => {
            navigate('/conversations');
            toast.success('Redirecting to conversations...');
          }}
          className="bg-blue-600 hover:bg-blue-700 dark:bg-blue-500 dark:hover:bg-blue-600 text-white transition-colors"
        >
          View Conversation
        </Button>
      )}
    </div>
  );
};

export default OrderActions;
