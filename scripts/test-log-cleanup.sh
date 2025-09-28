#!/bin/bash

# Script to test log cleanup functionality
# Usage: ./test-log-cleanup.sh

echo "🧪 Testing Log Cleanup Functionality"
echo "===================================="

# Test 1: Check if API endpoint works
echo "📡 Testing API endpoint..."
response=$(curl -s -X POST "http://localhost:8080/api/v1/admin/cron/log-cleanup?days=30" -H "Content-Type: application/json")
echo "Response: $response"

# Test 2: Check scheduled jobs
echo ""
echo "📋 Checking scheduled jobs..."
jobs_response=$(curl -s "http://localhost:8080/api/v1/admin/cron/jobs" -H "Content-Type: application/json")
echo "Scheduled jobs: $jobs_response"

# Test 3: Check log files before cleanup
echo ""
echo "📁 Log files before cleanup:"
ls -la ./storage/logs/*.log 2>/dev/null || echo "No log files found"

# Test 4: Create a test log file with old timestamp
echo ""
echo "📝 Creating test log file with old timestamp..."
test_log_file="./storage/logs/test-old-$(date +%Y%m%d).log"
echo '{"level":"INFO","time":"2024-01-01T00:00:00.000+0700","caller":"test","msg":"This is a test log file older than 30 days"}' > "$test_log_file"

# Set file modification time to 35 days ago
if [[ "$OSTYPE" == "darwin"* ]]; then
    # macOS
    touch -t $(date -v-35d +%Y%m%d%H%M) "$test_log_file"
else
    # Linux
    touch -d "35 days ago" "$test_log_file"
fi

echo "Created test file: $test_log_file"
ls -la "$test_log_file"

# Test 5: Run cleanup
echo ""
echo "🗑️ Running log cleanup..."
cleanup_response=$(curl -s -X POST "http://localhost:8080/api/v1/admin/cron/log-cleanup?days=30" -H "Content-Type: application/json")
echo "Cleanup response: $cleanup_response"

# Test 6: Check if test file was deleted
echo ""
echo "📁 Log files after cleanup:"
ls -la ./storage/logs/*.log 2>/dev/null || echo "No log files found"

if [ -f "$test_log_file" ]; then
    echo "❌ Test file still exists - cleanup may not be working"
else
    echo "✅ Test file was deleted - cleanup is working!"
fi

echo ""
echo "🎯 Test completed!"
