#!/usr/bin/env python3
import re

# Read the file
with open('internal/repository/event_repository.go', 'r') as f:
    lines = f.readlines()

# Functions that need parquetSource
functions_needing_fix = [
    'GetEntryExitPages',
    'GetTopCountries',
    'GetTopSources',
    'GetTopEvents',
    'GetBrowsersDevicesOS',
    'GetChannels',
    'GetProjects'
]

# Process line by line
i = 0
while i < len(lines):
    line = lines[i]
    
    # Check if this is a function definition that needs fixing
    for func_name in functions_needing_fix:
        if f'func (r *eventRepository) {func_name}' in line:
            # Find the first line after the function signature
            j = i + 1
            # Skip to the opening brace and first line of function body
            while j < len(lines) and '{' not in lines[j]:
                j += 1
            j += 1  # Move past the opening brace
            
            # Check if parquetSource is already defined
            if j < len(lines) and 'parquetSource :=' not in lines[j]:
                # Insert parquetSource definition
                indent = '\t'
                lines.insert(j, f'{indent}parquetSource := r.getParquetSource()\n')
                print(f"Added parquetSource to {func_name}")
            break
    
    i += 1

# Write back
with open('internal/repository/event_repository.go', 'w') as f:
    f.writelines(lines)

print("Fixed all functions to include parquetSource")
