#!/usr/bin/env python3
import re

# Read the file
with open('internal/repository/event_repository.go', 'r') as f:
    content = f.read()

# Pattern to find fmt.Sprintf with FROM %s that's missing parquetSource
# This regex looks for fmt.Sprintf calls with FROM %s where the first arg is whereClause (not parquetSource)
pattern = r'(fmt\.Sprintf\(`[^`]*FROM %s[^`]*`), whereClause\)'

# Replace with parquetSource as first parameter
replacement = r'\1, parquetSource, whereClause)'

# Perform the replacement
new_content = re.sub(pattern, replacement, content)

# Write back
with open('internal/repository/event_repository.go', 'w') as f:
    f.write(new_content)

print("Fixed all fmt.Sprintf calls to include parquetSource")
