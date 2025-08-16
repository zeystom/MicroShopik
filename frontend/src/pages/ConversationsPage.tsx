import React, { useState, useEffect, useRef, useCallback, useMemo } from 'react';
import { useAuthStore } from '@/store/authStore';
import { apiService } from '@/services/api';
import { Conversation, Message, Order, Product } from '@/types';
import ConversationList from '@/components/messaging/ConversationList';
import { Send, ArrowLeft, CheckCircle2, Loader2, Package, User as UserIcon, MessageCircle } from 'lucide-react';
import toast from 'react-hot-toast';
import { useSearchParams } from 'react-router-dom';
import { useAutoRefresh } from '@/hooks/useAutoRefresh';

const ConversationsPage: React.FC = () => {
  const { user } = useAuthStore();
  const [searchParams, setSearchParams] = useSearchParams();
  const [conversations, setConversations] = useState<Conversation[]>([]);
  const [selectedConversation, setSelectedConversation] = useState<Conversation | null>(null);
  const [messages, setMessages] = useState<Message[]>([]);
  const [newMessage, setNewMessage] = useState('');
  const [isSending, setIsSending] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const [orderId, setOrderId] = useState<number | null>(null);
  const [order, setOrder] = useState<Order | null>(null);
  const [isConfirming, setIsConfirming] = useState(false);
  const [isLoadingMessages, setIsLoadingMessages] = useState(false);
  const [isLoadingOrder, setIsLoadingOrder] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [isOnline, setIsOnline] = useState(navigator.onLine);
  const [isMobile, setIsMobile] = useState(false);
  const [mobileView, setMobileView] = useState<'list' | 'chat'>('list');
  const [isSwitchingConversation, setIsSwitchingConversation] = useState(false);
  const [listSearch, setListSearch] = useState('');
  
  const messagesContainerRef = useRef<HTMLDivElement>(null);
  const messageInputRef = useRef<HTMLInputElement>(null);
  const lastScrollTopRef = useRef<number>(0);
  const conversationsScrollRef = useRef<HTMLDivElement>(null);
  const wasNearBottomRef = useRef<boolean>(true);
  const activeConversationIdRef = useRef<number | null>(null);
  const isFetchingConversationsRef = useRef(false);
  const isFetchingMessagesRef = useRef(false);

  // Load conversations
  const loadConversations = useCallback(async (options?: { silent?: boolean; syncSelection?: boolean }) => {
    if (!user) return;
    
    try {
      if (isFetchingConversationsRef.current) return;
      isFetchingConversationsRef.current = true;
      if (!options?.silent) {
        setIsLoading(true);
      }
      setError(null);
      const data = await apiService.getConversations();

      // Preserve scroll position to avoid jank in the list
      const listEl = conversationsScrollRef.current;
      const prevScrollTop = listEl ? listEl.scrollTop : 0;
      const prevScrollHeight = listEl ? listEl.scrollHeight : 0;

      if (options?.silent) {
        // Merge without resorting to prevent jumping
        setConversations((prev) => {
          if (!prev || prev.length === 0) {
            return data;
          }
          const incomingById = new Map<number, Conversation>(data.map((c: Conversation) => [c.id, c]));
          const updated: Conversation[] = prev.map((oldConv: Conversation) => {
            const inc: Conversation | undefined = incomingById.get(oldConv.id);
            return inc ? { ...oldConv, ...inc } : oldConv;
          });
          // Append any new conversations at the end
          for (const conv of data as Conversation[]) {
            if (!prev.find((p: Conversation) => p.id === conv.id)) {
              updated.push(conv);
            }
          }
          return updated;
        });

        // After DOM updates, restore scroll offset relative to height change
        setTimeout(() => {
          const newEl = conversationsScrollRef.current;
          if (newEl) {
            const newHeight = newEl.scrollHeight;
            const delta = newHeight - prevScrollHeight;
            newEl.scrollTop = prevScrollTop + (delta > 0 ? delta : 0);
          }
        }, 0);
      } else {
        // Initial/full load: sort by last activity
        const sorted: Conversation[] = [...(data as Conversation[])].sort((a: Conversation, b: Conversation) => {
          const aLastMessage = a.messages[a.messages.length - 1];
          const bLastMessage = b.messages[b.messages.length - 1];
          if (aLastMessage && bLastMessage) {
            return new Date(bLastMessage.created_at).getTime() - new Date(aLastMessage.created_at).getTime();
          }
          return new Date(b.updated_at).getTime() - new Date(a.updated_at).getTime();
        });
        setConversations(sorted);
      }
      // Optionally keep selected conversation in sync (skip during silent refresh to reduce flicker)
      if (options?.syncSelection !== false) {
        setSelectedConversation((prev) => {
          if (!prev) return prev;
          const pool: Conversation[] = Array.isArray(conversations) ? conversations : [];
          const updated = pool.find((c: Conversation) => c.id === prev.id) || null;
          return updated || prev;
        });
      }
    } catch (error) {
      console.error('Error loading conversations:', error);
      setError('Failed to load conversations');
      toast.error('Failed to load conversations');
    } finally {
      if (!options?.silent) {
        setIsLoading(false);
      }
      isFetchingConversationsRef.current = false;
    }
  }, [user]);

  // Load messages for a conversation
  const loadMessages = useCallback(async (conversationId: number, options?: { silent?: boolean }): Promise<boolean> => {
    try {
      // Allow concurrent loads; ignore stale responses via activeConversationIdRef
      isFetchingMessagesRef.current = true;
      if (!options?.silent) {
        setIsLoadingMessages(true);
      }
      setError(null);
      // do not override user's manual scroll state here; it is tracked in onScroll
      const data = await apiService.getMessages(conversationId);
      
      // Sort messages chronologically
      const sorted = [...data].sort((a, b) => 
        new Date(a.created_at).getTime() - new Date(b.created_at).getTime()
      );
      
      const isForActiveConversation = activeConversationIdRef.current === conversationId;
      if (isForActiveConversation) {
        // Merge by id, keep stable order (avoid jank by not re-sorting existing messages)
        setMessages((prev) => {
          if (!options?.silent || prev.length === 0) return sorted;
          const existingById = new Map(prev.map(m => [m.id, m]));
          const merged = [...prev];
          for (const m of sorted) {
            if (!existingById.has(m.id)) merged.push(m);
          }
          return merged;
        });
        // scroll only if user was near bottom before update (stable check)
        if (!options?.silent && wasNearBottomRef.current) {
          setTimeout(() => {
            if (messagesContainerRef.current) {
              messagesContainerRef.current.scrollTop = messagesContainerRef.current.scrollHeight;
            }
          }, 0);
        }
        return true;
      }
      return false;
    } catch (error) {
      console.error('Error loading messages:', error);
      setError('Failed to load messages');
      toast.error('Failed to load messages');
      return false;
    } finally {
      if (!options?.silent) {
        setIsLoadingMessages(false);
      }
      isFetchingMessagesRef.current = false;
    }
  }, []);

  // Load order information
  const loadOrder = useCallback(async (id: number) => {
    try {
      setIsLoadingOrder(true);
      setError(null);
      const orderData = await apiService.getOrder(id);
      setOrder(orderData);
    } catch (error) {
      console.error('Error loading order:', error);
      setError('Failed to load order');
    } finally {
      setIsLoadingOrder(false);
    }
  }, []);

  // Send message
  const sendMessage = useCallback(async () => {
    if (!newMessage.trim() || !selectedConversation || !user || isSending) return;

    try {
      setIsSending(true);
      setError(null);
      
      const messageData = {
        text: newMessage.trim(),
        sender_id: user.id,
        order_id: order?.id || undefined
      };

      await apiService.sendMessage(selectedConversation.id, messageData);
      
      // Clear input and retain focus for continued typing
      setNewMessage('');
      
      // Reload messages and conversations
      await Promise.all([
        loadMessages(selectedConversation.id),
        loadConversations()
      ]);
      
      toast.success('Message sent');
      // Restore focus after sending
      if (messageInputRef.current) {
        messageInputRef.current.focus();
      }
    } catch (error) {
      console.error('Error sending message:', error);
      setError('Failed to send message');
      toast.error('Failed to send message');
    } finally {
      setIsSending(false);
    }
  }, [newMessage, selectedConversation, user, order?.id, loadMessages, loadConversations]);

  // Handle conversation selection
  const handleConversationSelect = useCallback((conversation: Conversation) => {
    const isSameConversation = selectedConversation?.id === conversation.id;

    if (!isSameConversation) {
      setIsSwitchingConversation(true);
      setSelectedConversation(conversation);
      setMessages([]);
      // Ensure first load scrolls to bottom once
      wasNearBottomRef.current = true;
      activeConversationIdRef.current = conversation.id;
      setNewMessage('');
    }
    setError(null);

    // Derive order from selected conversation and sync URL + state
    const derivedOrderId = conversation.messages.find(msg => msg.order_id)?.order_id || null;
    if (derivedOrderId) {
      if (derivedOrderId !== orderId) {
        setOrderId(derivedOrderId);
        loadOrder(derivedOrderId);
      }
      setSearchParams({ orderId: String(derivedOrderId) });
    } else {
      setOrderId(null);
      setOrder(null);
      setSearchParams({});
    }

    // On mobile, navigate to chat view
    if (isMobile) {
      setMobileView('chat');
    }
  }, [orderId, loadOrder, setSearchParams, selectedConversation?.id, isMobile]);

  // Handle order confirmation
  const handleConfirmOrder = useCallback(async () => {
    if (!order?.id) return;
    
    try {
      setIsConfirming(true);
      setError(null);
      await apiService.confirmOrder(order.id);
      toast.success('Order confirmed successfully!');
      await Promise.all([
        loadOrder(order.id),
        loadConversations()
      ]);
    } catch (error) {
      console.error('Error confirming order:', error);
      setError('Failed to confirm order');
      toast.error('Failed to confirm order');
    } finally {
      setIsConfirming(false);
    }
  }, [order?.id, loadOrder, loadConversations]);

  // Removed unused scrollToBottom helper; we scroll on load explicitly

  // Track user scroll position to prevent auto-scroll when user reads history
  const handleMessagesScroll = useCallback(() => {
    const el = messagesContainerRef.current;
    if (!el) return;
    const previousTop = lastScrollTopRef.current;
    const currentTop = el.scrollTop;
    const distanceFromBottom = el.scrollHeight - currentTop - el.clientHeight;

    // If user scrolls upward even slightly, disable auto-scroll immediately
    if (currentTop < previousTop) {
      wasNearBottomRef.current = false;
    } else {
      // Only consider "near bottom" when within a very small threshold
      const nearBottomThreshold = 8;
      wasNearBottomRef.current = distanceFromBottom <= nearBottomThreshold;
    }

    lastScrollTopRef.current = currentTop;
  }, []);

  // Handle key down in message input (onKeyPress is deprecated in React 18)
  const handleKeyDown = useCallback((e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      if (!isSending) {
        sendMessage();
      }
    }
  }, [sendMessage, isSending]);

  // Get product info
  const getProductInfo = useCallback((): Product | null => {
    if (selectedConversation?.product) {
      return selectedConversation.product;
    }
    if (order?.product) {
      return order.product;
    }
    return null;
  }, [selectedConversation, order]);

  // Get order info
  const getOrderInfo = useCallback((): Order | null => {
    return order;
  }, [order]);

  // Auto-refresh conversations
  const { pause: pauseAutoRefresh, resume: resumeAutoRefresh } = useAutoRefresh({
    interval: 15000,
    enabled: isOnline && !!user,
    onRefresh: () => loadConversations({ silent: true, syncSelection: false }),
    dependencies: [user?.id]
  });

  // Auto-refresh messages for selected conversation
  const { pause: pauseMessageRefresh, resume: resumeMessageRefresh } = useAutoRefresh({
    interval: 5000,
    enabled: isOnline && !!selectedConversation,
    onRefresh: async () => {
      if (selectedConversation) {
        await loadMessages(selectedConversation.id, { silent: true });
        if (order?.id) {
          await loadOrder(order.id);
        }
      }
    },
    dependencies: [selectedConversation?.id, order?.id]
  });

  // Handle online/offline status
  useEffect(() => {
    const handleOnline = () => {
      setIsOnline(true);
      resumeAutoRefresh();
      resumeMessageRefresh();
      // Reload data when coming back online
      if (user) {
        loadConversations();
        if (selectedConversation) {
          loadMessages(selectedConversation.id);
        }
      }
    };

    const handleOffline = () => {
      setIsOnline(false);
      pauseAutoRefresh();
      pauseMessageRefresh();
      setError('You are offline. Some features may not work.');
    };

    window.addEventListener('online', handleOnline);
    window.addEventListener('offline', handleOffline);

    return () => {
      window.removeEventListener('online', handleOnline);
      window.removeEventListener('offline', handleOffline);
    };
  }, [user, selectedConversation, loadConversations, loadMessages, resumeAutoRefresh, resumeMessageRefresh, pauseAutoRefresh, pauseMessageRefresh]);

  // Main effects
  useEffect(() => {
    // Load conversations when user changes
    if (user) {
      loadConversations();
    }
  }, [user, loadConversations]);

  useEffect(() => {
    // Handle orderId from URL params
    const orderIdParam = searchParams.get('orderId');
    if (orderIdParam) {
      const id = parseInt(orderIdParam);
      setOrderId(id);
      loadOrder(id);
    }
  }, [searchParams, loadOrder]);

  useEffect(() => {
    // Load messages when conversation changes
    if (selectedConversation) {
      // mark active conversation id
      activeConversationIdRef.current = selectedConversation.id;
      // First load should magnet to bottom once, then lock near-bottom for initial render
      wasNearBottomRef.current = true;
      (async () => {
        await loadMessages(selectedConversation.id);
        // Force scroll to bottom after first paint to avoid starting mid-list
        requestAnimationFrame(() => {
          if (messagesContainerRef.current) {
            messagesContainerRef.current.scrollTop = messagesContainerRef.current.scrollHeight;
          }
        });
        setIsSwitchingConversation(false);
      })();
      resumeMessageRefresh();
    } else {
      setMessages([]);
      pauseMessageRefresh();
    }
  }, [selectedConversation?.id, loadMessages, resumeMessageRefresh, pauseMessageRefresh]);

  // Detect mobile and keep mobile view in sync
  useEffect(() => {
    const handleResize = () => setIsMobile(window.innerWidth < 768);
    handleResize();
    window.addEventListener('resize', handleResize);
    return () => window.removeEventListener('resize', handleResize);
  }, []);

  useEffect(() => {
    if (!isMobile) {
      // Reset to list when leaving mobile (desktop shows both panes)
      setMobileView('list');
    } else if (!selectedConversation) {
      // If nothing selected on mobile, show list
      setMobileView('list');
    }
  }, [isMobile, selectedConversation]);

  // Removed global auto-scroll to bottom to prevent jank; handled in loader

  useEffect(() => {
    // Auto-select conversation if orderId is provided and conversations are loaded
    if (orderId && conversations.length > 0 && !selectedConversation) {
      const relatedConversation = conversations.find(conv => 
        conv.messages.some(msg => msg.order_id === orderId)
      );
      
      if (relatedConversation) {
        setSelectedConversation(relatedConversation);
      }
    }
  }, [orderId, conversations, selectedConversation]);

  // Refresh conversations
  const handleRefresh = useCallback(() => {
    loadConversations();
    if (selectedConversation) {
      loadMessages(selectedConversation.id);
    }
  }, [loadConversations, selectedConversation, loadMessages]);

  // Filter conversations for mobile list
  const filteredConversations = useMemo(() => {
    if (!listSearch.trim()) return conversations;
    const q = listSearch.toLowerCase();
    return conversations.filter((c) => {
      const title = c.product?.title?.toLowerCase() || '';
      const participants = c.participants
        .filter(p => p.user.id !== user?.id)
        .map(p => p.user.username.toLowerCase())
        .join(' ');
      const lastText = (c.messages[c.messages.length - 1]?.text || '').toLowerCase();
      return title.includes(q) || participants.includes(q) || lastText.includes(q);
    });
  }, [conversations, listSearch, user?.id]);

  if (!user) {
    return (
      <div className="min-h-screen bg-gray-50 dark:bg-gray-950 flex items-center justify-center">
        <div className="text-center">
          <h1 className="text-2xl font-bold text-gray-900 mb-4 dark:text-gray-100">Please log in to view messages</h1>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50 dark:bg-gray-950">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4 sm:py-8">
        <div className="bg-white rounded-none shadow-none border-0 sm:rounded-lg sm:shadow dark:bg-gray-900 sm:dark:border sm:dark:border-gray-800">
          {/* Status Bar */}
          {/* Status Bar removed */}

          {/* Error Display */}
          {error && (
            <div className="bg-red-50 border border-red-200 p-4 rounded-t-lg">
              <div className="flex items-center justify-between">
                <div className="flex items-center">
                  <div className="text-red-400">
                    <svg className="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                      <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clipRule="evenodd" />
                    </svg>
                  </div>
                  <p className="ml-3 text-sm text-red-800">{error}</p>
                </div>
                <button
                  onClick={() => setError(null)}
                  className="text-red-400 hover:text-red-600"
                >
                  <svg className="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                    <path fillRule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clipRule="evenodd" />
                  </svg>
                </button>
              </div>
            </div>
          )}

          {/* Responsive layout: stack on mobile, split on desktop */}
          {!isMobile ? (
            <div className="flex h-[80vh]">
              {/* Left Panel - Product Info & Conversations */}
              <div className="w-2/5 border-r border-gray-200 flex flex-col min-h-0 overflow-hidden dark:border-slate-800">
              {/* Product Information */}
              <div className="p-4 border-b border-gray-200 bg-gray-50 flex-shrink-0 dark:bg-slate-900 dark:border-slate-800">
                <div className="flex items-center justify-between mb-4">
                  <h2 className="text-lg font-semibold text-gray-900 flex items-center dark:text-slate-100">
                    <Package className="w-5 h-5 mr-2" />
                    Product Information
                  </h2>
                </div>
                
                {getProductInfo() ? (
                  <div className="bg-white rounded-lg p-4 border border-gray-200 dark:bg-slate-900 dark:border-slate-800">
                    <h3 className="font-medium text-gray-900 mb-2 dark:text-slate-100">{getProductInfo()?.title}</h3>
                    <p className="text-gray-600 text-sm mb-3 dark:text-slate-300">{getProductInfo()?.description}</p>
                    <div className="text-lg font-bold text-green-600 mb-3">
                      ${((getProductInfo()?.price || 0) / 100).toFixed(2)}
                    </div>
                    
                    {/* Product Status */}
                    <div className="flex items-center justify-between mb-3">
                      <span className="text-sm text-gray-600 dark:text-slate-300">Status:</span>
                      <span className={`px-2 py-1 rounded-full text-xs font-medium ${
                        getProductInfo()?.is_active ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'
                      }`}>
                        {getProductInfo()?.is_active ? 'Active' : 'Inactive'}
                      </span>
                    </div>
                    
                    {/* Product Type */}
                    {getProductInfo()?.disposable && (
                      <div className="flex items-center justify-between mb-3">
                        <span className="text-sm text-gray-600 dark:text-slate-300">Type:</span>
                        <span className="px-2 py-1 rounded-full text-xs font-medium bg-indigo-100 text-indigo-800 dark:bg-indigo-900/40 dark:text-indigo-300">
                          Disposable
                        </span>
                      </div>
                    )}
                    
                    {getOrderInfo() && (
                      <div className="border-t border-gray-200 pt-3 mt-3 dark:border-slate-800">
                        <h4 className="font-medium text-gray-900 mb-2 dark:text-slate-100">Order Details</h4>
                        {isLoadingOrder ? (
                          <div className="flex items-center justify-center py-2">
                            <Loader2 className="animate-spin rounded-full h-4 w-4 border-b-2 border-blue-600" />
                            <span className="ml-2 text-sm text-gray-600 dark:text-slate-300">Loading order...</span>
                          </div>
                        ) : (
                          <>
                            <div className="text-sm text-gray-600 mb-2 dark:text-slate-300">
                              Order #{getOrderInfo()?.id}
                            </div>
                            <div className="text-sm text-gray-600 mb-3 dark:text-slate-300">
                              Status: <span className={`font-medium ${
                                getOrderInfo()?.status === 'completed' ? 'text-green-600' :
                                getOrderInfo()?.status === 'confirmed' ? 'text-blue-600' :
                                getOrderInfo()?.status === 'pending' ? 'text-yellow-600' :
                                'text-red-600'
                              }`}>{getOrderInfo()?.status}</span>
                            </div>
                            
                            {getOrderInfo()?.status === 'pending' && getOrderInfo()?.customer_id === user.id && (
                              <button
                                onClick={handleConfirmOrder}
                                disabled={isConfirming}
                                className="w-full px-4 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700 disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center space-x-2 dark:bg-green-500 dark:hover:bg-green-600"
                              >
                                {isConfirming ? (
                                  <>
                                    <Loader2 className="h-4 w-4 animate-spin" />
                                    <span>Confirming...</span>
                                  </>
                                ) : (
                                  <>
                                    <CheckCircle2 className="h-4 w-4" />
                                    <span>Confirm Payment</span>
                                  </>
                                )}
                              </button>
                            )}
                          </>
                        )}
                      </div>
                    )}
                  </div>
                ) : (
                  <div className="bg-white rounded-lg p-4 border border-gray-200 text-center dark:bg-gray-900 dark:border-gray-800">
                    <Package className="w-12 h-12 text-gray-400 mx-auto mb-2" />
                    <p className="text-gray-500 dark:text-gray-400">Select a conversation to view product details</p>
                  </div>
                )}
              </div>

              {/* Conversations List */}
              <div className="flex-1 p-4 min-h-0">
                <div className="mb-4 flex-shrink-0">
                  <h3 className="text-lg font-semibold text-gray-900 flex items-center dark:text-slate-100">
                    <UserIcon className="w-5 h-5 mr-2" />
                    Conversations
                  </h3>
                </div>
                
                <div className="h-full overflow-y-scroll conversation-list pb-4" ref={conversationsScrollRef}>
                  {isLoading ? (
                    <div className="flex items-center justify-center h-32">
                      <Loader2 className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600" />
                    </div>
                  ) : (
                    <ConversationList
                      onConversationSelect={handleConversationSelect}
                      currentUser={user}
                      selectedConversationId={selectedConversation?.id}
                      conversations={conversations}
                    />
                  )}
                </div>
              </div>
            </div>

            {/* Right Panel - Chat */}
            <div className="flex-1 flex flex-col chat-container min-h-0 overflow-hidden">
              {selectedConversation ? (
                <>
                  {/* Chat Header */}
                  <div className="border-b border-gray-200 p-4 bg-white dark:bg-slate-900 dark:border-slate-800">
                    <div className="flex items-center justify-between">
                      <div className="flex items-center space-x-3">
                        <div>
                          <h3 className="text-lg font-semibold text-gray-900 dark:text-slate-100">
                            Chat with {selectedConversation.participants
                              .filter(p => p.user.id !== user.id)
                              .map(p => p.user.username)
                              .join(', ')}
                          </h3>
                          <p className="text-sm text-gray-500 dark:text-slate-400">
                            {messages.length} messages
                            {getProductInfo() && (
                              <span className="ml-2">• {getProductInfo()?.title}</span>
                            )}
                          </p>
                        </div>
                      </div>
                      
                    </div>
                  </div>

                  {/* Messages Container */}
                  <div 
                    ref={messagesContainerRef}
                    className={`chat-messages transition-opacity duration-150 ${isSwitchingConversation ? 'opacity-0' : 'opacity-100'}`}
                    onScroll={handleMessagesScroll}
                  >
                    {isLoadingMessages && messages.length === 0 ? (
                      <div className="flex items-center justify-center h-32">
                        <Loader2 className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600" />
                        <span className="ml-2 text-gray-600 dark:text-slate-300">Loading messages...</span>
                      </div>
                    ) : isSwitchingConversation ? (
                      <div className="flex items-center justify-center h-32">
                        <Loader2 className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600" />
                        <span className="ml-2 text-gray-600 dark:text-slate-300">Opening conversation...</span>
                      </div>
                    ) : messages.length === 0 ? (
                      <div className="flex items-center justify-center h-32">
                        <div className="text-center">
                          <MessageCircle className="w-12 h-12 text-gray-400 mx-auto mb-2" />
                          <p className="text-gray-500 dark:text-slate-400">No messages yet</p>
                          <p className="text-sm text-gray-400">Start the conversation!</p>
                        </div>
                      </div>
                    ) : (
                      messages.map((message) => (
                        <div
                          key={message.id}
                          className={`flex ${
                            message.sender_id === user.id ? 'justify-end' : 'justify-start'
                          }`}
                        >
                          <div
                            className={`message-bubble ${
                              message.sender_id === user.id ? 'sent' : 'received'
                            }`}
                          >
                            <div className="text-sm font-medium">
                              {message.sender?.username || 'System'}
                            </div>
                            <div className="mt-1">{message.text}</div>
                            <div className={`text-xs mt-1 ${
                              message.sender_id === user.id ? 'text-blue-100' : 'text-gray-500'
                            }`}>
                              {new Date(message.created_at).toLocaleTimeString()}
                            </div>
                          </div>
                        </div>
                      ))
                    )}
                  </div>

                  {/* Message Input */}
                  <div className="chat-input dark:bg-gray-900 dark:border-gray-800">
                    <div className="flex space-x-2">
                      <input
                        ref={messageInputRef}
                        type="text"
                        value={newMessage}
                        onChange={(e) => setNewMessage(e.target.value)}
                        onKeyDown={handleKeyDown}
                        placeholder="Type your message..."
                        className="flex-1 px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 transition-all dark:bg-gray-800 dark:border-gray-700 dark:text-gray-100"
                        disabled={!isOnline}
                      />
                      <button
                        onClick={sendMessage}
                        disabled={!newMessage.trim() || isSending || !isOnline}
                        className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center space-x-2 transition-all dark:bg-blue-500 dark:hover:bg-blue-600"
                      >
                        {isSending ? (
                          <Loader2 className="h-4 w-4 animate-spin" />
                        ) : (
                          <Send className="h-4 w-4" />
                        )}
                        <span>{isSending ? 'Sending...' : 'Send'}</span>
                      </button>
                    </div>
                  </div>
                </>
              ) : (
                <div className="flex-1 flex items-center justify-center">
                  <div className="w-full max-w-md bg-white border border-slate-200 rounded-xl p-8 shadow-sm text-center dark:bg-slate-900 dark:border-slate-800">
                    <div className="w-16 h-16 mx-auto mb-4 rounded-full bg-indigo-100 flex items-center justify-center dark:bg-indigo-900/40">
                      <MessageCircle className="w-8 h-8 text-indigo-600 dark:text-indigo-300" />
                    </div>
                    <h3 className="text-xl font-semibold text-slate-900 mb-2 dark:text-slate-100">No conversation selected</h3>
                    <p className="text-slate-600 mb-6 dark:text-slate-400">Choose a chat from the list on the left to start messaging</p>
                    {orderId && (
                      <div className="bg-indigo-50 border border-indigo-200 rounded-lg p-4 mb-2 text-left dark:bg-indigo-950 dark:border-indigo-900">
                        <p className="text-sm text-indigo-800 mb-3 dark:text-indigo-300">
                          You are viewing conversations related to order #{orderId}
                        </p>
                        <button
                          onClick={() => {
                            toast.success('Creating conversation for this order...');
                          }}
                          className="w-full px-4 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 text-sm transition-colors dark:bg-indigo-500 dark:hover:bg-indigo-600"
                        >
                          Start new conversation for this order
                        </button>
                      </div>
                    )}
                  </div>
                </div>
              )}
            </div>
            </div>
          ) : (
            // Mobile layout (polished)
            <div className="flex flex-col min-h-[80vh]">
              {mobileView === 'list' ? (
                <div className="flex-1 flex flex-col min-h-0">
                  {/* Sticky mobile header with search */}
                  <div className="sticky top-0 z-10 bg-white/95 dark:bg-slate-900/95 backdrop-blur supports-[backdrop-filter]:bg-white/70 dark:supports-[backdrop-filter]:bg-slate-900/70 border-b border-gray-200 dark:border-slate-800">
                    <div className="px-4 py-3">
                      <div className="flex items-center justify-between mb-2">
                        <h3 className="text-base font-semibold text-gray-900 dark:text-slate-100 flex items-center">
                          <UserIcon className="w-5 h-5 mr-2" />
                          Conversations
                        </h3>
                        <button onClick={handleRefresh} className="text-xs text-blue-600 dark:text-blue-400">Refresh</button>
                      </div>
                      <div className="flex items-center space-x-2">
                        <input
                          type="search"
                          value={listSearch}
                          onChange={(e) => setListSearch(e.target.value)}
                          placeholder="Search by title, user, or message..."
                          className="flex-1 px-3 py-2 rounded-md border border-gray-300 dark:border-gray-700 dark:bg-slate-800 dark:text-slate-100 text-sm"
                        />
                      </div>
                    </div>
                  </div>

                  {/* Conversation list */}
                  <div className="flex-1 overflow-y-auto conversation-list px-3 pb-3" ref={conversationsScrollRef}>
                    {isLoading ? (
                      <div className="flex items-center justify-center h-32">
                        <Loader2 className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600" />
                      </div>
                    ) : (
                      <ConversationList
                        onConversationSelect={handleConversationSelect}
                        currentUser={user}
                        selectedConversationId={selectedConversation?.id}
                        conversations={filteredConversations}
                      />
                    )}
                  </div>
                </div>
              ) : (
                <div className="flex-1 flex flex-col min-h-0">
                  {/* Chat header */}
                  <div className="sticky top-0 z-10 border-b border-gray-200 p-3 bg-white/95 dark:bg-slate-900/95 backdrop-blur dark:border-slate-800">
                    <div className="flex items-center">
                      <button
                        onClick={() => setMobileView('list')}
                        className="mr-2 p-2 rounded-md hover:bg-gray-100 dark:hover:bg-gray-800"
                        aria-label="Back to conversations"
                      >
                        <ArrowLeft className="w-5 h-5" />
                      </button>
                      <div className="min-w-0">
                        <h3 className="text-base font-semibold text-gray-900 truncate dark:text-slate-100 max-w-[55vw]">
                          {selectedConversation?.participants
                            .filter(p => p.user.id !== user.id)
                            .map(p => p.user.username)
                            .join(', ')}
                        </h3>
                        <p className="text-xs text-gray-500 truncate dark:text-slate-400 max-w-[55vw]">
                          {getProductInfo()?.title || 'Direct chat'}
                        </p>
                      </div>
                      <div className="ml-auto flex items-center space-x-2 flex-shrink-0">
                        {getOrderInfo() && getOrderInfo()?.status === 'pending' && getOrderInfo()?.customer_id === user?.id && (
                          <button
                            onClick={handleConfirmOrder}
                            disabled={isConfirming}
                            className="text-xs px-2 py-1 rounded-md bg-emerald-600 text-white hover:bg-emerald-700 disabled:opacity-50"
                          >
                            {isConfirming ? '...' : 'Confirm'}
                          </button>
                        )}
                        {getProductInfo() && (
                          <button
                            onClick={() => {
                              const id = getProductInfo()?.id
                              if (id) window.location.hash = `#/products/${id}`
                            }}
                            className="text-xs px-2 py-1 rounded-md border border-gray-300 text-gray-700 hover:bg-gray-50 dark:border-gray-700 dark:text-gray-300 dark:hover:bg-gray-800"
                          >
                            View product
                          </button>
                        )}
                      </div>
                    </div>
                  </div>

                  {/* Messages */}
                  <div
                    ref={messagesContainerRef}
                    className={`chat-messages transition-opacity duration-150 ${isSwitchingConversation ? 'opacity-0' : 'opacity-100'}`}
                    onScroll={handleMessagesScroll}
                  >
                    {isLoadingMessages && messages.length === 0 ? (
                      <div className="flex items-center justify-center h-32">
                        <Loader2 className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600" />
                        <span className="ml-2 text-gray-600 dark:text-slate-300">Loading messages...</span>
                      </div>
                    ) : isSwitchingConversation ? (
                      <div className="flex items-center justify-center h-32">
                        <Loader2 className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600" />
                        <span className="ml-2 text-gray-600 dark:text-slate-300">Opening conversation...</span>
                      </div>
                    ) : messages.length === 0 ? (
                      <div className="flex items-center justify-center h-32">
                        <div className="text-center">
                          <MessageCircle className="w-12 h-12 text-gray-400 mx-auto mb-2" />
                          <p className="text-gray-500 dark:text-slate-400">No messages yet</p>
                          <p className="text-sm text-gray-400">Start the conversation!</p>
                        </div>
                      </div>
                    ) : (
                      messages.map((message) => (
                        <div
                          key={message.id}
                          className={`flex ${message.sender_id === user.id ? 'justify-end' : 'justify-start'}`}
                        >
                          <div className={`message-bubble ${message.sender_id === user.id ? 'sent' : 'received'}`}>
                            <div className="text-sm font-medium">
                              {message.sender?.username || 'System'}
                            </div>
                            <div className="mt-1">{message.text}</div>
                            <div className={`text-xs mt-1 ${message.sender_id === user.id ? 'text-blue-100' : 'text-gray-500'}`}>
                              {new Date(message.created_at).toLocaleTimeString()}
                            </div>
                          </div>
                        </div>
                      ))
                    )}
                  </div>

                  {/* Input bar */}
                  <div className="sticky bottom-0 bg-white/95 dark:bg-slate-900/95 backdrop-blur border-t border-gray-200 dark:border-gray-800">
                    <div className="chat-input dark:bg-transparent dark:border-transparent">
                      <div className="flex items-end space-x-2">
                        <input
                          ref={messageInputRef}
                          type="text"
                          value={newMessage}
                          onChange={(e) => setNewMessage(e.target.value)}
                          onKeyDown={handleKeyDown}
                          placeholder="Type a message..."
                          className="flex-1 px-3 py-3 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 transition-all dark:bg-gray-800 dark:border-gray-700 dark:text-gray-100"
                          disabled={!isOnline}
                        />
                        <button
                          onClick={sendMessage}
                          disabled={!newMessage.trim() || isSending || !isOnline}
                          className="w-10 h-10 rounded-full bg-blue-600 text-white hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center transition-all dark:bg-blue-500 dark:hover:bg-blue-600 flex-shrink-0"
                        >
                          {isSending ? <Loader2 className="h-4 w-4 animate-spin" /> : <Send className="h-4 w-4" />}
                        </button>
                      </div>
                    </div>
                  </div>
                </div>
              )}
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

export default ConversationsPage;