# MicroShopik - Digital Goods Marketplace

A full-stack digital goods marketplace built with Go (backend) and React TypeScript (frontend).

## Features

- **User Authentication & Authorization**: JWT-based authentication with role-based access control
- **Digital Goods Management**: Products, categories, and inventory management
- **Order Processing**: Complete order lifecycle with status tracking
- **Messaging System**: Real-time conversations between buyers and sellers
- **Role-Based Access**: Admin, Seller, and Customer roles with different permissions
- **RESTful API**: Comprehensive API with Swagger documentation

## Tech Stack

### Backend
- **Go** with Echo framework
- **GORM** for database ORM
- **JWT** for authentication
- **PostgreSQL** database
- **Swagger** for API documentation

### Frontend
- **React 18** with TypeScript
- **Vite** for build tooling
- **Tailwind CSS** for styling
- **React Router** for navigation
- **Zustand** for state management
- **Axios** for API communication
- **React Hot Toast** for notifications

## Prerequisites

- Go 1.21+
- Node.js 18+
- PostgreSQL 13+
- Git

## Quick Start

### 1. Clone the Repository
```bash
git clone <repository-url>
cd MicroShopik
```

### 2. Install Dependencies
```bash
# Install all dependencies (backend and frontend)
npm run install:all
```

### 3. Database Setup
```bash
# Create PostgreSQL database
createdb microshopik

# Run migrations
migrate -path migrations -database "postgres://username:password@localhost:5432/microshopik?sslmode=disable" up

# Seed initial roles
go run scripts/seed_roles.go
```

### 4. Environment Configuration

Create a `.env` file in the frontend directory:
```bash
cd frontend
echo "VITE_API_URL=http://localhost:8080" > .env
```

### 5. Run the Application

#### Option A: Run Both Backend and Frontend Together
```bash
npm run dev
```

#### Option B: Run Separately
```bash
# Terminal 1 - Backend
npm run dev:backend

# Terminal 2 - Frontend
npm run dev:frontend
```

### 6. Access the Application

- **Frontend**: http://localhost:5173
- **Backend API**: http://localhost:8080
- **API Documentation**: http://localhost:8080/swagger/

## Project Structure

```
MicroShopik/
├── configs/                 # Configuration files
├── docs/                   # API documentation
├── frontend/               # React TypeScript frontend
│   ├── src/
│   │   ├── components/     # React components
│   │   ├── pages/         # Page components
│   │   ├── services/      # API services
│   │   ├── store/         # State management
│   │   └── types/         # TypeScript types
├── internal/               # Go backend
│   ├── controllers/       # HTTP controllers
│   ├── domain/           # Domain models
│   ├── middleware/       # HTTP middleware
│   ├── repositories/     # Data access layer
│   └── services/         # Business logic
├── migrations/            # Database migrations
└── scripts/              # Utility scripts
```

## API Endpoints

### Authentication
- `POST /auth/register` - User registration
- `POST /auth/login` - User login

### Products
- `GET /products` - List products
- `GET /products/:id` - Get product details
- `POST /products` - Create product (Seller/Admin)
- `PUT /products/:id` - Update product (Seller/Admin)
- `DELETE /products/:id` - Delete product (Seller/Admin)

### Orders
- `GET /orders` - List user orders
- `POST /orders` - Create order
- `GET /orders/:id` - Get order details
- `PUT /orders/:id/status` - Update order status
- `POST /orders/:id/process` - Process order
- `POST /orders/:id/cancel/:customerID` - Cancel order

### Conversations
- `GET /conversations` - List conversations
- `POST /conversations` - Create conversation
- `GET /conversations/:id` - Get conversation details
- `GET /conversations/:id/messages` - Get messages
- `POST /conversations/:id/messages` - Send message

### Admin Routes
- `GET /admin/dashboard` - Admin dashboard
- `GET /admin/users` - List all users
- `GET /roles` - List roles
- `POST /users/:user_id/roles/:role_name` - Assign role

## Role-Based Access Control

### Admin
- Full access to all features
- User management
- Role assignment
- System configuration

### Seller
- Product management
- Order processing
- Customer communication
- Sales analytics

### Customer
- Browse products
- Place orders
- View order history
- Communicate with sellers

## Development

### Backend Development
```bash
# Run with hot reload
go run main.go

# Run tests
go test ./...

# Generate API docs
swag init
```

### Frontend Development
```bash
cd frontend

# Start development server
npm run dev

# Build for production
npm run build

# Run tests
npm test
```

## Database Migrations

```bash
# Create new migration
migrate create -ext sql -dir migrations -seq migration_name

# Run migrations
migrate -path migrations -database "postgres://username:password@localhost:5432/microshopik?sslmode=disable" up

# Rollback migrations
migrate -path migrations -database "postgres://username:password@localhost:5432/microshopik?sslmode=disable" down
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

MIT License - see LICENSE file for details
