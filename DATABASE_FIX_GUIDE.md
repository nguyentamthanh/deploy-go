# 🔧 Database Connection Fix Guide

## 🎯 **Root Cause Analysis**

Lỗi `failed to connect to 'user=postgres database=deploy_go'` xảy ra vì:

1. **Environment Variables**: Render không set đúng environment variables
2. **Connection Method**: Cần sử dụng DATABASE_URL thay vì individual variables
3. **SSL Configuration**: Production cần SSL mode khác với development

## ✅ **Solutions Applied**

### **1. Dual Connection Strategy**

```go
// Try DATABASE_URL first (Render's preferred method)
if databaseURL := os.Getenv("DATABASE_URL"); databaseURL != "" {
    return connectWithDatabaseURL(databaseURL)
}
// Fallback to individual environment variables
```

### **2. Updated render.yaml**

```yaml
envVars:
  - key: DATABASE_URL
    fromDatabase:
      name: deploy-go-db
      property: connectionString
  # Plus individual variables as fallback
```

### **3. Enhanced Error Handling**

- ✅ Detailed logging without exposing passwords
- ✅ Retry mechanism (5 attempts, 2-second intervals)
- ✅ Connection test before auto migration
- ✅ Fallback mechanisms

### **4. Production-Ready Configuration**

- ✅ SSL mode: `require` for production
- ✅ Connection pooling optimized
- ✅ Proper timeout handling

## 🚀 **Deployment Steps**

### **Step 1: Commit and Push**

```bash
git add .
git commit -m "Fix database connection with dual strategy"
git push origin main
```

### **Step 2: Monitor Render Logs**

Look for these log messages:

- `🔗 Using DATABASE_URL for connection` ✅
- `🔗 Connecting to database: host=...` ✅
- `✅ Database connection test successful` ✅
- `✅ Database tables auto-migrated successfully` ✅

### **Step 3: Test Health Endpoint**

```bash
curl https://your-app.onrender.com/health
```

Expected response:

```json
{
  "status": "healthy",
  "database": "connected"
}
```

## 🧪 **Testing Locally**

### **Using Docker Compose**

```bash
# Start PostgreSQL
docker-compose up -d postgres

# Set local environment
export DATABASE_URL="postgres://postgres:password@localhost:5432/deploy_go?sslmode=disable"

# Run application
go run .
```

### **Using Individual Variables**

```bash
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=password
export DB_NAME=deploy_go
export DB_SSLMODE=disable

go run .
```

## 🔍 **Debug Information**

### **Environment Variables Priority**

1. `DATABASE_URL` (Render's preferred)
2. Individual `DB_*` variables (fallback)

### **Connection String Format**

```
postgres://username:password@host:port/database?sslmode=require
```

### **Common Issues & Solutions**

| Issue                              | Solution                             |
| ---------------------------------- | ------------------------------------ |
| `user=postgres database=deploy_go` | Use DATABASE_URL instead             |
| SSL connection failed              | Set `sslmode=require` for production |
| Connection timeout                 | Added retry mechanism                |
| Environment variables not set      | Added fallback to individual vars    |

## 📊 **Monitoring**

### **Health Check Endpoint**

- `GET /health` - Returns database status
- `GET /` - Main endpoint with database info

### **Log Messages to Watch**

- ✅ `🔗 Using DATABASE_URL for connection`
- ✅ `✅ Database connection test successful`
- ✅ `✅ Database tables auto-migrated successfully`
- ❌ `❌ Database connection failed: ...`

## 🎉 **Expected Results**

After deployment, you should see:

1. ✅ Application starts successfully
2. ✅ Database connection established
3. ✅ Tables auto-migrated
4. ✅ Health check returns "healthy"
5. ✅ API endpoints working

**Status: 🚀 Ready for Production Deployment!**
