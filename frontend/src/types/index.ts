export interface User {
  id: number
  username: string
  email: string
  roles: Role[]
  created_at: string
}

export interface Role {
  id: number
  name: string
  description: string
}

export interface Product {
  id: number
  seller_id: number
  title: string
  description: string
  price: number // Price in cents, will be converted to dollars in UI
  category_id: number
  category: Category
  is_active: boolean
  disposable: boolean
  max_sales: number
  sold_count: number
  created_at: string
  updated_at: string
}

export interface Category {
  id: number
  name: string
  description: string
}

export interface Order {
  id: number
  customer_id?: number
  product_id?: number
  product?: Product
  status: string
  created_at: string
  updated_at: string
}

export interface Message {
  id: number
  conversation_id: number
  sender_id: number
  sender?: User
  order_id?: number
  text: string
  is_system?: boolean
  created_at: string
}

export interface Conversation {
  id: number
  product_id?: number
  product?: Product
  participants: Participant[]
  messages: Message[]
  created_at: string
  updated_at: string
}

export interface Participant {
  id: number
  conversation_id: number
  user_id: number
  user: User
  role: string
  created_at: string
}

export interface AuthRequest {
  username: string
  password: string
}

export interface AuthResponse {
  token: string
  user: User
}

export interface ApiResponse<T> {
  success: boolean
  data?: T
  message?: string
}

// Digital goods specific types
export interface GameAccount {
  id: number
  product_id: number
  product: Product
  game_name: string
  account_level: string
  platform: string
  region: string
  additional_info: string
  is_available: boolean
}

export interface StreamingAccount {
  id: number
  product_id: number
  product: Product
  service_name: string
  account_type: string
  subscription_tier: string
  region: string
  additional_info: string
  is_available: boolean
}

