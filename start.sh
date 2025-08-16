#!/bin/bash

echo "Starting MicroShopik..."
echo

echo "Checking if Go is installed..."
if ! command -v go &> /dev/null; then
    echo "Error: Go is not installed or not in PATH"
    echo "Please install Go from https://golang.org/"
    exit 1
fi

echo "Checking if Node.js is installed..."
if ! command -v node &> /dev/null; then
    echo "Error: Node.js is not installed or not in PATH"
    echo "Please install Node.js from https://nodejs.org/"
    exit 1
fi

echo "Installing dependencies..."
npm run install:all

echo
echo "Starting MicroShopik in development mode..."
echo "Backend will be available at: http://localhost:8080"
echo "Frontend will be available at: http://localhost:5173"
echo
echo "Press Ctrl+C to stop the application"
echo

npm run dev
