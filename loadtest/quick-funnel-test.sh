#!/bin/bash
# Quick Funnel Data Generator and Tester
# This script generates funnel data and opens the funnel analysis page

set -e

echo "🎯 Siraaj Funnel Quick Test Script"
echo "=================================="
echo ""

# Default values
USERS=${1:-10000}
DAYS=${2:-30}
DB_PATH="../data/analytics.db"

echo "📊 Configuration:"
echo "   Users: $USERS"
echo "   Days: $DAYS"
echo "   Database: $DB_PATH"
echo ""

# Check if database exists
if [ ! -f "$DB_PATH" ]; then
    echo "❌ Database not found at $DB_PATH"
    echo "   Please ensure the database exists or update DB_PATH"
    exit 1
fi

# Generate funnel data
echo "🔄 Generating funnel data..."
cd funnel
go run main.go -mode=db -users=$USERS -days=$DAYS -db="../../data/analytics.db"

if [ $? -eq 0 ]; then
    echo ""
    echo "✅ Funnel data generated successfully!"
    echo ""
    echo "📊 Next Steps:"
    echo "   1. Start the server: make run (in project root)"
    echo "   2. Open: http://localhost:8080/dashboard/funnel"
    echo ""
    echo "🧪 Try these funnel configurations:"
    echo ""
    echo "   E-commerce Purchase Funnel:"
    echo "   ---------------------------"
    echo "   Step 1: event_name = \"page_view\", url = \"/\""
    echo "   Step 2: event_name = \"page_view\", url = \"/product/123\""
    echo "   Step 3: event_name = \"add_to_cart\", url = \"/product/123\""
    echo "   Step 4: event_name = \"checkout_started\", url = \"/checkout\""
    echo "   Step 5: event_name = \"purchase\", url = \"/confirmation\""
    echo ""
    echo "   SaaS Activation Funnel:"
    echo "   -----------------------"
    echo "   Step 1: event_name = \"page_view\", url = \"/\""
    echo "   Step 2: event_name = \"page_view\", url = \"/pricing\""
    echo "   Step 3: event_name = \"signup\", url = \"/signup\""
    echo "   Step 4: event_name = \"page_view\", url = \"/dashboard\""
    echo ""
    echo "   Content Newsletter Funnel:"
    echo "   --------------------------"
    echo "   Step 1: event_name = \"page_view\", url = \"/blog\""
    echo "   Step 2: event_name = \"page_view\", url = \"/blog/article-1\""
    echo "   Step 3: event_name = \"button_click\", url = \"/blog/article-1\""
    echo "   Step 4: event_name = \"form_submit\", url = \"/blog/article-1\""
    echo ""
else
    echo ""
    echo "❌ Failed to generate funnel data"
    exit 1
fi
