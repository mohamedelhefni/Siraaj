import re

with open('internal/repository/event_repository.go', 'r') as f:
    content = f.read()

# Fix line 452 (month timeline query)
content = re.sub(
    r"(strftime\(DATE_TRUNC\('month', timestamp\), '%%Y-%%m-01'\) as date, \n\s+%s\n\s+FROM %s \n\s+WHERE %s\n\s+GROUP BY date \n\s+ORDER BY date\n\s+`), (selectClause), (whereClause)\)",
    r"\1, \2, parquetSource, \3)",
    content
)

print("Fixed sprintf issues")

with open('internal/repository/event_repository.go', 'w') as f:
    f.write(content)
