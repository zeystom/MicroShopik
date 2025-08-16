# MicroShopik Setup Guide

## Prerequisites

- Go 1.21 or higher
- Node.js 18 or higher
- PostgreSQL database
- Git

## Quick Start

### 1. Database Setup

1. Create a PostgreSQL database
2. Update the database configuration in `configs/config.go` if needed
3. Run migrations:
   ```bash
   # Install golang-migrate if you haven't already
   go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
   
   # Run migrations
   migrate -path migrations -database "postgres://username:password@localhost:5432/microshopik?sslmode=disable" up
   ```

### 2. Seed Initial Data

```bash
# Seed roles
go run scripts/seed_roles.go

# Seed test data (users, categories, products)
go run scripts/seed_data.go
```

### 3. Start Development Servers

#### Option 1: Using the batch script (Windows)
```bash
start_dev.bat
```

#### Option 2: Manual start

**Backend:**
```bash
go run main.go
```

**Frontend:**
```bash
cd frontend
npm install
npm run dev
```

## Test Accounts

After running the seed script, you can use these test accounts:

- **Admin**: admin@microshopik.com / password
- **Seller**: seller1@microshopik.com / password  
- **Buyer**: buyer1@microshopik.com / password

## API Endpoints

- Backend API: http://localhost:8080
- Frontend: http://localhost:5173
- Swagger Docs: http://localhost:8080/swagger/

## Features

### Authentication & Authorization
- JWT-based authentication
- Role-based access control (admin, seller, buyer, moderator)
- Protected routes with middleware

### Product Management
- Create, read, update, delete products
- Category-based organization
- Product availability checking

### Order System
- Create orders
- Process orders (creates conversations)
- Cancel orders
- Order status management

### Messaging System
- Conversations between buyers and sellers
- Real-time messaging (WebSocket ready)
- Message history

### User Management
- User registration and login
- Profile management
- Role assignment (admin only)

## Project Structure

```
MicroShopik/
├── configs/           # Configuration
├── docs/             # Documentation and Swagger
├── frontend/         # React + TypeScript frontend
├── internal/         # Go backend
│   ├── controllers/  # HTTP handlers
│   ├── database/     # Database connection
│   ├── domain/       # Data models
│   ├── middleware/   # HTTP middleware
│   ├── repositories/ # Data access layer
│   └── services/     # Business logic
├── migrations/       # Database migrations
├── scripts/          # Seed scripts
└── main.go          # Application entry point
```

## Development

### Backend Development
- Uses Echo framework for HTTP routing
- GORM for database operations
- JWT for authentication
- Structured logging

### Frontend Development
- React 18 with TypeScript
- Vite for build tooling
- Tailwind CSS for styling
- Zustand for state management
- React Router for navigation

## Troubleshooting

### Common Issues

1. **Database Connection Error**
   - Check PostgreSQL is running
   - Verify connection string in config
   - Ensure database exists

2. **CORS Issues**
   - Backend has CORS middleware configured
   - Frontend should connect to http://localhost:8080

3. **Authentication Issues**
   - Check JWT secret in config
   - Verify user exists and has roles assigned

4. **Frontend Build Issues**
   - Run `npm install` in frontend directory
   - Check Node.js version compatibility

## Production Deployment

1. Build frontend: `cd frontend && npm run build`
2. Set environment variables for production
3. Use a production database
4. Configure reverse proxy (nginx recommended)
5. Set up SSL certificates

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request
