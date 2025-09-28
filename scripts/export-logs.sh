#!/bin/bash

# Script to export log data in a more readable format
# Usage: ./export-logs.sh [output_file]

LOG_DIR="./storage/logs"
OUTPUT_FILE=${1:-"exported-logs-$(date +%Y%m%d-%H%M%S).txt"}

echo "📊 Veloras API Log Exporter"
echo "=========================="
echo "Log directory: $LOG_DIR"
echo "Output file: $OUTPUT_FILE"
echo ""

# Check if log directory exists
if [ ! -d "$LOG_DIR" ]; then
    echo "❌ Log directory not found: $LOG_DIR"
    exit 1
fi

# Create output file
echo "📋 Veloras API Log Export - $(date)" > "$OUTPUT_FILE"
echo "===============================================" >> "$OUTPUT_FILE"
echo "" >> "$OUTPUT_FILE"

# Process each log file
for log_file in "$LOG_DIR"/*.log; do
    if [ -f "$log_file" ]; then
        filename=$(basename "$log_file")
        echo "📄 Processing: $filename"
        
        echo "File: $filename" >> "$OUTPUT_FILE"
        echo "Size: $(stat -f%z "$log_file" 2>/dev/null || stat -c%s "$log_file" 2>/dev/null) bytes" >> "$OUTPUT_FILE"
        echo "Modified: $(stat -f%Sm "$log_file" 2>/dev/null || stat -c%y "$log_file" 2>/dev/null)" >> "$OUTPUT_FILE"
        echo "----------------------------------------" >> "$OUTPUT_FILE"
        
        # Parse JSON logs and format them nicely
        while IFS= read -r line; do
            if [ -n "$line" ]; then
                # Try to parse as JSON and format
                echo "$line" | python3 -c "
import json
import sys
try:
    data = json.loads(sys.stdin.read().strip())
    level = data.get('level', 'UNKNOWN')
    time = data.get('time', '')
    caller = data.get('caller', '')
    msg = data.get('msg', '')
    
    # Color coding for different levels
    if level == 'ERROR' or level == 'FATAL':
        level_color = '🔴'
    elif level == 'WARN':
        level_color = '🟡'
    elif level == 'INFO':
        level_color = '🔵'
    elif level == 'DEBUG':
        level_color = '⚪'
    else:
        level_color = '⚫'
    
    print(f'{level_color} [{level}] {time}')
    if caller:
        print(f'   📍 {caller}')
    print(f'   💬 {msg}')
    
    # Print additional fields
    for key, value in data.items():
        if key not in ['level', 'time', 'caller', 'msg']:
            print(f'   📊 {key}: {value}')
    print()
    
except:
    # If not JSON, just print as is
    print(f'📝 {sys.stdin.read().strip()}')
    print()
" >> "$OUTPUT_FILE"
            fi
        done < "$log_file"
        
        echo "" >> "$OUTPUT_FILE"
        echo "===============================================" >> "$OUTPUT_FILE"
        echo "" >> "$OUTPUT_FILE"
    fi
done

echo "✅ Log export completed!"
echo "📁 Output saved to: $OUTPUT_FILE"
echo ""
echo "📖 To view the exported logs:"
echo "   cat $OUTPUT_FILE"
echo "   less $OUTPUT_FILE"
echo "   open $OUTPUT_FILE"
