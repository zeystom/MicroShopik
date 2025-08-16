import React from 'react'
import ConversationsPage from '@/pages/ConversationsPage'

// Simple wrapper to present the same conversations UI with seller-focused title/description handled in ConversationsPage itself.
// We reuse the existing chat implementation so sellers can chat with clients.
const SellerMessagesPage: React.FC = () => {
  return <ConversationsPage />
}

export default SellerMessagesPage




