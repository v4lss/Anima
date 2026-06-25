#!/bin/bash
# Deployment script for Animas

set -e

echo "🚀 Deploying Animas..."

# Build API
echo "🔨 Building API..."
cd api
go build -o animas cmd/animas/main.go

# Build Frontend
echo "🔨 Building Frontend..."
cd ../web
npm run build

# Copy frontend build to api static
echo "📦 Copying frontend build..."
cd ..
rm -rf api/static
mkdir -p api/static
cp -r web/dist/* api/static/

echo "✅ Build complete!"
echo ""
echo "To deploy, copy the following:"
echo "  - api/animas (binary)"
echo "  - api/configs/ (configuration)"
echo "  - api/static/ (frontend)"
echo ""
echo "Then run: ./animas"
