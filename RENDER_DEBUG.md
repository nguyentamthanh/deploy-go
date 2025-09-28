# Render Deployment Debug Guide

## Common Issues và Solutions

### 1. Database Connection Error

**Error**: `failed to connect to database, got error failed to connect to 'user=postgres database=deploy_go'`

**Solutions**:

1. ✅ **Fixed**: Updated `render.yaml` to use correct database user (`postgres` instead of `deploy_go_user`)
2. ✅ **Added**: Retry mechanism with 5 attempts and 2-second intervals
3. ✅ **Added**: Better error logging and connection details

### 2. Environment Variables

Make sure these environment variables are properly set in Render:

```yaml
envVars:
  - key: GIN_MODE
    value: release
  - key: DB_HOST
    fromDatabase:
      name: deploy-go-db
      property: host
  - key: DB_PORT
    fromDatabase:
      name: deploy-go-db
      property: port
  - key: DB_USER
    fromDatabase:
      name: deploy-go-db
      property: user
  - key: DB_PASSWORD
    fromDatabase:
      name: deploy-go-db
      property: password
  - key: DB_NAME
    fromDatabase:
      name: deploy-go-db
      property: database
  - key: DB_SSLMODE
    value: require
```

### 3. Database Service Configuration

```yaml
- type: pserv
  name: deploy-go-db
  plan: free
  databaseName: deploy_go
  user: postgres # ✅ Fixed: was deploy_go_user
```

### 4. Testing Locally

```bash
# Start PostgreSQL with Docker
docker-compose up -d postgres

# Test database connection
go run test-db-connection.go

# Run application
go run .
```

### 5. Render Deployment Steps

1. Push code to GitHub repository
2. Connect repository to Render
3. Render will automatically:
   - Create PostgreSQL database service
   - Build Docker image
   - Set environment variables
   - Deploy application

### 6. Debug Logs

The application now includes detailed logging:

- 🔗 Connection details (without password)
- ⏳ Retry attempts
- ✅ Success confirmations
- ❌ Error details

### 7. Health Check

The application provides health check endpoint:

- `GET /health` - Returns database connection status

### 8. Common Render Issues

- **Database not ready**: Added retry mechanism
- **Wrong credentials**: Fixed user configuration
- **SSL issues**: Set `DB_SSLMODE=require` for production
- **Timeout**: Increased retry attempts to 5

## Next Steps

1. Commit and push these changes
2. Redeploy on Render
3. Check Render logs for detailed connection info
4. Test `/health` endpoint to verify database connection
