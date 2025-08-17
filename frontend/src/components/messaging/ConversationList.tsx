import React, { useState, useCallback, useMemo, useEffect } from 'react';
import { apiService } from '@/services/api';
import { Conversation, User, Order } from '@/types';
import { MessageCircle, Users, Clock, AlertCircle } from 'lucide-react';
import ProductStatusBadge from '@/components/ui/ProductStatusBadge';

interface ConversationListProps {
  onConversationSelect: (conversation: Conversation) => void;
  currentUser: User;
  selectedConversationId?: number;
  conversations?: Conversation[];
  onConversationsUpdate?: () => void;
}

const ConversationList: React.FC<ConversationListProps> = ({ 
  onConversationSelect, 
  currentUser,
  selectedConversationId,
  conversations: externalConversations,
  onConversationsUpdate: _onConversationsUpdate
}) => {
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [orderStatusByConv, setOrderStatusByConv] = useState<Record<number, string>>({});

  // Use external conversations if provided
  const displayConversations = externalConversations || [];

  // Fetch order statuses per conversation using first message's order_id
  useEffect(() => {
    const pending: Array<{ convId: number; orderId: number }> = [];
    for (const conv of displayConversations) {
      if (!conv || !conv.messages || conv.messages.length === 0) continue;
      const withOrder = conv.messages.find(m => typeof m.order_id === 'number' && m.order_id > 0);
      if (withOrder && orderStatusByConv[conv.id] === undefined) {
        pending.push({ convId: conv.id, orderId: withOrder.order_id as number });
      }
    }
    if (pending.length === 0) return;
    (async () => {
      try {
        const results = await Promise.all(
          pending.map(async ({ convId, orderId }) => {
            try {
              const order: Order = await apiService.getOrder(orderId);
              return { convId, status: order?.status || '' };
            } catch {
              return { convId, status: '' };
            }
          })
        );
        setOrderStatusByConv(prev => {
          const next = { ...prev };
          for (const r of results) {
            if (r.status) next[r.convId] = r.status;
          }
          return next;
        });
      } catch {}
    })();
  }, [displayConversations]);

  // Only show loading if no external conversations and we're loading
  if (!externalConversations && isLoading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
      </div>
    );
  }

  if (!displayConversations || displayConversations.length === 0) {
    return (
      <div className="text-center py-8">
        <MessageCircle className="h-12 w-12 text-gray-400 mx-auto mb-4" />
        <h3 className="text-lg font-medium text-gray-900 mb-2">No conversations yet</h3>
        <p className="text-gray-600">
          Start a conversation when you process an order
        </p>
      </div>
    );
  }

  const formatDate = (dateString: string) => {
    try {
      const date = new Date(dateString);
      if (isNaN(date.getTime())) {
        return 'Invalid date';
      }
      
      const now = new Date();
      const diffInHours = (now.getTime() - date.getTime()) / (1000 * 60 * 60);
      
      if (diffInHours < 24) {
        return date.toLocaleTimeString('en-US', {
          hour: '2-digit',
          minute: '2-digit'
        });
      } else if (diffInHours < 48) {
        return 'Yesterday';
      } else {
        return date.toLocaleDateString('en-US', {
          month: 'short',
          day: 'numeric'
        });
      }
    } catch (error) {
      return 'Invalid date';
    }
  };

  const getOtherParticipants = (conversation: Conversation) => {
    return conversation.participants.filter(
      p => p.user.id !== currentUser.id
    );
  };

  const getProductFromFirstMessage = (conversation: Conversation) => {
    const first = conversation.messages && conversation.messages[0];
    if (!first || !first.text) return null;
    // Expecting pattern: "Purchase: product name"
    const match = first.text.match(/^Purchase:\s*(.+)$/i);
    return match ? match[1] : null;
  };

  const getOrderStatusFromMessages = (conversation: Conversation): string | null => {
    // Prefer last message with apparent status keyword
    const reversed = [...conversation.messages].reverse();
    for (const m of reversed) {
      const t = (m.text || '').toLowerCase();
      if (t.includes('confirmed')) return 'confirmed';
      if (t.includes('completed')) return 'completed';
      if (t.includes('pending')) return 'pending';
      if (t.includes('canceled') || t.includes('cancelled')) return 'canceled';
    }
    return null;
  };

  const getLastMessagePreview = (conversation: Conversation) => {
    if (!conversation.messages || conversation.messages.length === 0) {
      return 'No messages yet';
    }
    
    const lastMessage = conversation.messages[conversation.messages.length - 1];
    if (!lastMessage) return 'No messages yet';
    
    const senderName = lastMessage.sender?.username || 'System';
    const messageText = lastMessage.text || '';
    
    // Truncate long messages
    const truncatedText = messageText.length > 50 
      ? messageText.substring(0, 50) + '...' 
      : messageText;
    
    return `${senderName}: ${truncatedText}`;
  };

  const handleConversationClick = (conversation: Conversation) => {
    try {
      onConversationSelect(conversation);
    } catch (error) {
      console.error('Error selecting conversation:', error);
      setError('Failed to select conversation');
    }
  };

  return (
    <div className="pr-1 sm:pr-2 pb-3 sm:pb-4 space-y-2 conversation-list">
      {/* Error Display */}
      {error && (
        <div className="bg-red-50 border border-red-200 p-3 rounded-lg mb-3">
          <div className="flex items-center">
            <AlertCircle className="h-4 w-4 text-red-400 mr-2" />
            <p className="text-sm text-red-800">{error}</p>
            <button
              onClick={() => setError(null)}
              className="ml-auto text-red-400 hover:text-red-600"
            >
              <svg className="h-4 w-4" viewBox="0 0 20 20" fill="currentColor">
                <path fillRule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clipRule="evenodd" />
              </svg>
            </button>
          </div>
        </div>
      )}

      {displayConversations.map((conversation) => {
        const otherParticipants = getOtherParticipants(conversation);
        const lastMessage = conversation.messages[conversation.messages.length - 1];
        const isSelected = selectedConversationId === conversation.id;
        const hasUnreadMessages = lastMessage && 
          lastMessage.sender_id !== currentUser.id && 
          !lastMessage.is_system;
        const chatProduct = conversation.product?.title || getProductFromFirstMessage(conversation);
        const chatTitle = chatProduct
          ? `${chatProduct} · ${otherParticipants.map(p => p.user.username).join(', ')}`
          : otherParticipants.map(p => p.user.username).join(', ');
        const derivedStatus = getOrderStatusFromMessages(conversation);
        const orderStatus = orderStatusByConv[conversation.id] || derivedStatus;
        
        return (
          <div
            key={conversation.id}
            onClick={() => handleConversationClick(conversation)}
            className={`p-4 border rounded-lg cursor-pointer transition-all duration-200 border-l-4 ${
              isSelected 
                ? 'bg-indigo-50 border-indigo-200 shadow-sm ring-2 ring-indigo-200 dark:bg-indigo-950 dark:border-indigo-900 dark:ring-indigo-900' 
                : 'border-slate-200 hover:bg-slate-50 hover:border-slate-300 dark:border-slate-800 dark:hover:bg-slate-800'
            } ${hasUnreadMessages ? 'border-l-indigo-500' : 'border-l-transparent'}`}
          >
            <div className="flex items-center justify-between mb-1 sm:mb-2">
              <div className="flex items-center space-x-2">
                <Users className={`h-4 w-4 ${
                  isSelected ? 'text-indigo-600 dark:text-indigo-300' : 'text-slate-500 dark:text-slate-400'
                }`} />
                <span className={`font-medium truncate max-w-[65vw] sm:max-w-56 ${
                  isSelected ? 'text-indigo-900 dark:text-indigo-200' : 'text-slate-900 dark:text-slate-100'
                }`}>
                  {chatTitle}
                </span>
                {hasUnreadMessages && (
                  <div className="w-2 h-2 bg-indigo-500 rounded-full flex-shrink-0"></div>
                )}
              </div>
              <span className={`text-xs flex-shrink-0 ${
                isSelected ? 'text-indigo-600 dark:text-indigo-400' : 'text-slate-500 dark:text-slate-400'
              }`}>
                {formatDate(conversation.updated_at)}
              </span>
            </div>
            
            <div className="flex items-center space-x-2 mb-1 sm:mb-2">
              <span className={`text-sm truncate flex-1 ${
                isSelected ? 'text-indigo-700 dark:text-indigo-300' : 'text-slate-600 dark:text-slate-400'
              }`}>
                {getLastMessagePreview(conversation)}
              </span>
            </div>
            
            <div className="flex items-center justify-between">
              <span className={`text-xs ${
                isSelected ? 'text-indigo-500 dark:text-indigo-400' : 'text-slate-500 dark:text-slate-400'
              }`}>
                {conversation.messages.length} messages
              </span>
              <div className="flex items-center space-x-2">
                {/* Показываем статус продукта если он неактивен */}
                {conversation.product && !conversation.product.is_active && (
                  <ProductStatusBadge isActive={conversation.product.is_active} />
                )}
                {orderStatus && (
                  <span className={`inline-flex px-2 py-0.5 rounded-full text-[10px] font-medium capitalize ${
                    orderStatus === 'completed' || orderStatus === 'confirmed'
                      ? 'bg-emerald-100 text-emerald-700 dark:bg-emerald-900/40 dark:text-emerald-300'
                      : orderStatus === 'pending'
                      ? 'bg-amber-100 text-amber-800 dark:bg-amber-900/40 dark:text-amber-300'
                      : 'bg-rose-100 text-rose-800 dark:bg-rose-900/40 dark:text-rose-300'
                  }`}>
                    {orderStatus}
                  </span>
                )}
                <Clock className={`h-3 w-3 ${
                  isSelected ? 'text-indigo-400 dark:text-indigo-400' : 'text-slate-400 dark:text-slate-500'
                }`} />
              </div>
            </div>
          </div>
        );
      })}
    </div>
  );
};

export default ConversationList;
