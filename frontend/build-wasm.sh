#!/bin/bash

echo "Building WASM..."

cd src/wasm

# Build WASM
GOOS=js GOARCH=wasm go build -o monkey.wasm main.go

if [ $? -eq 0 ]; then
    echo "✅ WASM build successful"
    
    # Copy to public directory
    cp monkey.wasm ../../public/
    # check if copied to public directory
    if [ -f ../../public/monkey.wasm ]; then
        echo "Copied to public directory"
    else
        echo "Failed to copy to public directory"
        exit 1
    fi
    
    # Clean up source directory
    rm monkey.wasm
    echo "Cleaned up source directory"
    
    echo "✅ WASM ready for serving!"
else
    echo "❌ WASM build failed"
    exit 1
fi
