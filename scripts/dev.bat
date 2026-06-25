@echo off
REM Development setup script for Animas (Windows)

echo 🚀 Setting up Animas development environment...

REM Check if Docker is installed
docker --version >nul 2>&1
if %errorlevel% neq 0 (
    echo ❌ Docker is not installed. Please install Docker Desktop first.
    exit /b 1
)

REM Check if Docker Compose is installed
docker-compose --version >nul 2>&1
if %errorlevel% neq 0 (
    echo ❌ Docker Compose is not installed. Please install Docker Compose first.
    exit /b 1
)

REM Start MongoDB and Redis
echo 📦 Starting MongoDB and Redis...
docker-compose -f docker-compose.dev.yml up -d

REM Wait for services to be ready
echo ⏳ Waiting for services to be ready...
timeout /t 5 /nobreak >nul

REM Copy example env if not exists
if not exist api\.env (
    echo 📝 Creating .env file from example...
    copy api\.env.example api\.env
)

if not exist web\.env (
    echo 📝 Creating .env file for web...
    echo VITE_API_URL=http://localhost:8080 > web\.env
)

REM Install Go dependencies
echo 🔧 Installing Go dependencies...
cd api
go mod download

REM Install Node dependencies
echo 🔧 Installing Node dependencies...
cd ..\web
call npm install

echo ✅ Development environment ready!
echo.
echo To start the API server:
echo   cd api ^&^& go run cmd\animas\main.go
echo.
echo To start the frontend:
echo   cd web ^&^& npm run dev
