@echo off
REM Build script for Windows

echo Building AI Terminal Pro...

REM Build frontend
cd frontend
call npm install
call npm run build
cd ..

REM Build Go application
call wails build

echo Build complete! Check build/bin/ directory.
pause
