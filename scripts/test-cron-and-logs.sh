#!/bin/bash

# Test script for Cron Jobs and Log Viewer
# Usage: ./scripts/test-cron-and-logs.sh

BASE_URL="http://localhost:8080/api/v1"
WEB_URL="http://localhost:8080"

echo "🚀 Testing Veloras API Cron Jobs and Log Viewer"
echo "================================================"

# Check if server is running
echo "📡 Checking if server is running..."
if ! curl -s "$BASE_URL/admin/cron/jobs" > /dev/null; then
    echo "❌ Server is not running. Please start the server first:"
    echo "   go run cmd/server/main.go"
    exit 1
fi
echo "✅ Server is running"

echo ""
echo "🔍 Testing Cron Job Endpoints"
echo "-----------------------------"

# Test 1: Get scheduled jobs
echo "1. Getting scheduled cron jobs..."
response=$(curl -s "$BASE_URL/admin/cron/jobs")
echo "Response: $response"
echo ""

# Test 2: Run manual cleanup
echo "2. Running manual cleanup..."
response=$(curl -s -X POST "$BASE_URL/admin/cron/cleanup")
echo "Response: $response"
echo ""

echo "📊 Testing Log Viewer Endpoints"
echo "-------------------------------"

# Test 3: Get log files
echo "3. Getting available log files..."
response=$(curl -s "$BASE_URL/admin/logs/files")
echo "Response: $response"
echo ""

# Test 4: Get log content (if any log files exist)
echo "4. Getting log content..."
# First, let's get the list of files and extract the first one
files_response=$(curl -s "$BASE_URL/admin/logs/files")
if echo "$files_response" | grep -q '"name"'; then
    # Extract first filename from JSON response
    filename=$(echo "$files_response" | grep -o '"name":"[^"]*"' | head -1 | cut -d'"' -f4)
    if [ ! -z "$filename" ]; then
        echo "Getting content for file: $filename"
        response=$(curl -s "$BASE_URL/admin/logs/content?filename=$filename&lines=10")
        echo "Response: $response"
    else
        echo "No log files found"
    fi
else
    echo "No log files available"
fi
echo ""

echo "🌐 Web Interface"
echo "---------------"
echo "Log Viewer Web Interface: $WEB_URL"
echo "Direct link: $WEB_URL/web/logs.html"
echo ""

echo "📝 Log Files Location"
echo "--------------------"
echo "Log files are stored in: ./storage/logs/"
echo "Daily log files format: dev.001-YYYY-MM-DD.log"
echo ""

echo "✅ Test completed!"
echo ""
echo "📋 Summary:"
echo "- Cron jobs endpoint: $BASE_URL/admin/cron/jobs"
echo "- Manual cleanup: $BASE_URL/admin/cron/cleanup (POST)"
echo "- Log files list: $BASE_URL/admin/logs/files"
echo "- Log content: $BASE_URL/admin/logs/content?filename=<filename>&lines=<number>"
echo "- Web interface: $WEB_URL"
echo ""
echo "💡 Tips:"
echo "- Check the web interface for a better log viewing experience"
echo "- Log files are automatically rotated daily"
echo "- Use the manual cleanup endpoint to test cron functionality"
