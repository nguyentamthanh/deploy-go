# API Documentation - GORM ORM

## Database Models (GORM)

### User

- `id` (uint): Primary key (auto-increment)
- `name` (string): User's full name (required)
- `email` (string): User's email address (unique, required)
- `posts` ([]Post): Related posts (one-to-many relationship)
- `created_at` (timestamp): Creation timestamp (auto-managed)
- `updated_at` (timestamp): Last update timestamp (auto-managed)
- `deleted_at` (timestamp): Soft delete timestamp (auto-managed)

### Post

- `id` (uint): Primary key (auto-increment)
- `title` (string): Post title (required)
- `content` (string): Post content (text type)
- `user_id` (uint): Foreign key to users table (required)
- `user` (\*User): Related user (many-to-one relationship)
- `created_at` (timestamp): Creation timestamp (auto-managed)
- `updated_at` (timestamp): Last update timestamp (auto-managed)
- `deleted_at` (timestamp): Soft delete timestamp (auto-managed)

## API Endpoints

### Users

#### Get All Users

- **GET** `/api/users`
- **Response**: Array of user objects

#### Get User by ID

- **GET** `/api/users/{id}`
- **Response**: Single user object

#### Create User

- **POST** `/api/users`
- **Body**:

```json
{
  "name": "John Doe",
  "email": "john@example.com"
}
```

#### Update User

- **PUT** `/api/users/{id}`
- **Body** (all fields optional):

```json
{
  "name": "Jane Doe",
  "email": "jane@example.com"
}
```

#### Delete User

- **DELETE** `/api/users/{id}`
- **Response**: Success message

### Posts

#### Get All Posts

- **GET** `/api/posts`
- **Response**: Array of post objects with user information

#### Get Post by ID

- **GET** `/api/posts/{id}`
- **Response**: Single post object with user information

#### Create Post

- **POST** `/api/posts`
- **Body**:

```json
{
  "title": "My First Post",
  "content": "This is the content of my post",
  "user_id": 1
}
```

#### Update Post

- **PUT** `/api/posts/{id}`
- **Body** (all fields optional):

```json
{
  "title": "Updated Title",
  "content": "Updated content"
}
```

#### Delete Post

- **DELETE** `/api/posts/{id}`
- **Response**: Success message

#### Get Posts by User

- **GET** `/api/posts/user/{userId}`
- **Response**: Array of post objects for specific user

## Environment Variables

### Database Configuration

- `DB_HOST`: Database host (default: localhost)
- `DB_PORT`: Database port (default: 5432)
- `DB_USER`: Database username (default: postgres)
- `DB_PASSWORD`: Database password (default: password)
- `DB_NAME`: Database name (default: deploy_go)
- `DB_SSLMODE`: SSL mode (default: disable)

### Application Configuration

- `GIN_MODE`: Gin mode (debug/release)
- `PORT`: Application port (default: 8080)

## Response Format

All API responses follow this format:

```json
{
  "success": true,
  "message": "Operation successful",
  "data": {} // Optional data object
}
```

## Error Responses

Error responses also follow the same format:

```json
{
  "success": false,
  "message": "Error description"
}
```

## Local Development

1. Start PostgreSQL database
2. Set environment variables (see .env.example)
3. Run the application:

```bash
go run .
```

## GORM Features

The application now uses GORM ORM which provides:

- **Auto Migration**: Tables are automatically created/updated based on model definitions
- **Soft Deletes**: Records are not physically deleted, marked with `deleted_at`
- **Relationships**: Automatic handling of foreign keys and relationships
- **Validation**: Built-in validation based on struct tags
- **Hooks**: Automatic timestamps (`created_at`, `updated_at`)
- **Preloading**: Efficient loading of related data

## Deployment

The application is configured to work with Render.com and includes:

- PostgreSQL database service
- GORM ORM for database operations
- Automatic environment variable configuration
- Health check endpoint at `/health`
