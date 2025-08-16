# Admin Panel - MicroShopik

## Overview
The Admin Panel provides comprehensive system management capabilities for administrators, including user management, role management, system settings, and detailed analytics.

## Features

### 🏠 Admin Dashboard (`/admin-dashboard`)
- **System Overview**: Real-time statistics for users, products, orders, and revenue
- **Quick Actions**: Direct access to key administrative functions
- **System Status**: Monitor database, API services, and order processing
- **Recent Activity**: View latest system activities and orders

### 👥 User Management (`/admin/users`)
- **View All Users**: Complete list of registered users with search and filtering
- **User Details**: Username, email, roles, and creation date
- **Role Management**: Assign/remove roles from users
- **User Actions**: Edit user information and delete users
- **Create Users**: Add new user accounts with role assignments

### 🛡️ Role Management (`/admin/roles`)
- **View All Roles**: Complete list of system roles with descriptions
- **Role Details**: Role name, description, and associated users
- **Create Roles**: Add new roles to the system
- **Edit Roles**: Modify role names and descriptions
- **Delete Roles**: Remove roles (with user impact warnings)

### ⚙️ System Settings (`/admin/settings`)
- **General Settings**: Maintenance mode, registration control, user limits
- **System Configuration**: Notifications, backup settings, log levels
- **System Status**: Real-time monitoring of database and services
- **System Logs**: View recent system activity and error logs

### 📊 Reports & Analytics (`/admin/reports`)
- **Key Metrics**: Revenue, orders, users, and conversion rates
- **Time Periods**: Filter data by 7 days, 30 days, 90 days, or 1 year
- **Export Reports**: Download CSV reports for external analysis
- **Trend Analysis**: Revenue trends and order statistics
- **User Analytics**: New user registrations and activity

## API Endpoints

### User Management
- `GET /admin/users` - Get all users
- `GET /admin/users/:id` - Get specific user
- `PUT /admin/users/:id` - Update user
- `DELETE /admin/users/:id` - Delete user

### Product Management
- `GET /admin/products` - Get all products
- `PUT /admin/products/:id/status` - Update product status
- `DELETE /admin/products/:id` - Delete product

### Order Management
- `GET /admin/orders` - Get all orders
- `PUT /admin/orders/:id/status` - Update order status
- `DELETE /admin/orders/:id` - Delete order

### System Management
- `GET /admin/stats` - Get system statistics
- `GET /admin/logs` - Get system logs
- `PUT /admin/settings` - Update system settings

### Role Management
- `GET /roles` - Get all roles
- `POST /roles` - Create new role
- `PUT /roles/:id` - Update role
- `DELETE /roles/:id` - Delete role
- `POST /users/:user_id/roles/:role_name` - Assign role to user
- `DELETE /users/:user_id/roles/:role_name` - Remove role from user

## Usage Instructions

### Accessing Admin Panel
1. Login with an account that has the `admin` role
2. Navigate to `/admin-dashboard`
3. Use the navigation buttons to access different sections

### Managing Users
1. Go to `/admin/users`
2. Use search and filters to find specific users
3. Click the edit button to modify user information
4. Use role checkboxes to assign/remove roles
5. Click delete button to remove users (with confirmation)

### Managing Roles
1. Go to `/admin/roles`
2. Click "Create New Role" to add roles
3. Fill in role name and description
4. Use edit button to modify existing roles
5. Delete roles (warning: affects all users with that role)

### System Configuration
1. Go to `/admin/settings`
2. Toggle switches for boolean settings
3. Enter values for numeric settings
4. Select options from dropdowns
5. Click "Save Settings" to apply changes

### Generating Reports
1. Go to `/admin/reports`
2. Select time period and report type
3. View key metrics and detailed breakdowns
4. Click "Export Report" to download CSV

## Security Features

- **Role-based Access Control**: Only users with `admin` role can access
- **JWT Authentication**: Secure token-based authentication
- **Input Validation**: Server-side validation of all inputs
- **Audit Logging**: System logs track all administrative actions

## Theme Support

The admin panel supports both light and dark themes:
- **Light Theme**: Clean, professional appearance
- **Dark Theme**: Easy on the eyes for extended use
- **Auto-switching**: Respects system preference
- **Manual Toggle**: Users can override system preference

## Responsive Design

- **Mobile First**: Optimized for mobile devices
- **Tablet Support**: Responsive layouts for medium screens
- **Desktop Enhanced**: Full feature set on large screens
- **Touch Friendly**: Optimized for touch interactions

## Browser Compatibility

- **Modern Browsers**: Chrome, Firefox, Safari, Edge
- **ES6+ Support**: Requires modern JavaScript features
- **CSS Grid/Flexbox**: Modern CSS layout support
- **Progressive Enhancement**: Graceful degradation for older browsers

## Troubleshooting

### Common Issues

1. **Access Denied**: Ensure your account has the `admin` role
2. **Data Not Loading**: Check network connection and API status
3. **Permission Errors**: Verify role assignments and permissions
4. **Theme Issues**: Clear browser cache and localStorage

### Performance Tips

1. **Refresh Data**: Use refresh buttons to get latest information
2. **Filter Results**: Use search and filters to reduce data load
3. **Export Reports**: Download large datasets instead of viewing in browser
4. **Close Modals**: Close unused modals to free up memory

## Development Notes

### Frontend Technologies
- **React 18**: Modern React with hooks
- **TypeScript**: Type-safe development
- **Tailwind CSS**: Utility-first CSS framework
- **Lucide Icons**: Modern icon library

### State Management
- **React Hooks**: useState, useEffect for local state
- **Context API**: Theme and authentication context
- **Local Storage**: Theme preferences and user settings

### API Integration
- **Axios**: HTTP client with interceptors
- **JWT Tokens**: Secure authentication
- **Error Handling**: Comprehensive error management
- **Loading States**: User feedback during operations

## Future Enhancements

- **Real-time Updates**: WebSocket integration for live data
- **Advanced Charts**: Interactive data visualizations
- **Bulk Operations**: Mass user/role management
- **Audit Trail**: Detailed action logging
- **API Documentation**: Swagger/OpenAPI integration
- **Mobile App**: Native mobile application
- **Multi-language**: Internationalization support
- **Advanced Analytics**: Machine learning insights
