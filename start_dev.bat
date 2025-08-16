@echo off
echo Starting MicroShopik Development Environment...

echo.
echo 1. Starting Backend Server...
start "Backend" cmd /k "go run main.go"

echo.
echo 2. Starting Frontend Development Server...
cd frontend
start "Frontend" cmd /k "npm run dev"

echo.
echo Development servers are starting...
echo Backend: http://localhost:8080
echo Frontend: http://localhost:5173
echo.
echo Press any key to exit this script (servers will continue running)
pause > nul
