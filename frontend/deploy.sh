#!/bin/bash

echo "🚀 Preparing Monkey Playground for Vercel deployment..."

# Build the API functions first
echo "🔧 Preparing API functions..."
cd api
go mod tidy
cd ..

# Build the project
echo "📦 Building frontend..."
npm run build

# Check if WASM file exists
if [ -f "public/monkey.wasm" ]; then
    echo "✅ WASM file found"
else
    echo "❌ WASM file missing! Building WASM..."
    cd src/wasm
    GOOS=js GOARCH=wasm go build -o monkey.wasm main.go
    cp monkey.wasm ../../public/
    cd ../..
fi

# Check if wasm_exec.js exists
if [ -f "public/wasm_exec.js" ]; then
    echo "✅ WASM exec script found"
else
    echo "❌ WASM exec script missing! Copying from Go installation..."
    cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" public/
fi

echo "✅ Ready for deployment!"
echo ""
echo "Next steps:"
echo "1. git add ."
echo "2. git commit -m 'Ready for Vercel deployment'"
echo "3. git push origin main"
echo "4. Deploy on vercel.com"
