#!/bin/bash
# Development setup script for Animas

set -e

echo "🚀 Setting up Animas development environment..."

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo "❌ Docker is not installed. Please install Docker first."
    exit 1
fi

# Check if Docker Compose is installed
if ! command -v docker-compose &> /dev/null; then
    echo "❌ Docker Compose is not installed. Please install Docker Compose first."
    exit 1
fi

# Start MongoDB and Redis
echo "📦 Starting MongoDB and Redis..."
docker-compose -f docker-compose.dev.yml up -d

# Wait for services to be ready
echo "⏳ Waiting for services to be ready..."
sleep 5

# Copy example env if not exists
if [ ! -f api/.env ]; then
    echo "📝 Creating .env file from example..."
    cp api/.env.example api/.env
fi

if [ ! -f web/.env ]; then
    echo "📝 Creating .env file for web..."
    echo "VITE_API_URL=http://localhost:8080" > web/.env
fi

# Install Go dependencies
echo "🔧 Installing Go dependencies..."
cd api
go mod download

# Install Node dependencies
echo "🔧 Installing Node dependencies..."
cd ../web
npm install

echo "✅ Development environment ready!"
echo ""
echo "To start the API server:"
echo "  cd api && go run cmd/animas/main.go"
echo ""
echo "To start the frontend:"
echo "  cd web && npm run dev"
