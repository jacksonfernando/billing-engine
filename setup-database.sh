#!/bin/bash

# Billing Engine Database Setup Script
echo "Setting up Billing Engine Database..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Database configuration
DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-3306}
DB_NAME=${DB_NAME:-billing_engine}
DB_USER=${DB_USER:-billing_admin}
DB_PASSWORD=${DB_PASSWORD:-billing_password}
ROOT_PASSWORD=${ROOT_PASSWORD:-root}

echo -e "${YELLOW}Database Configuration:${NC}"
echo "Host: $DB_HOST"
echo "Port: $DB_PORT"
echo "Database: $DB_NAME"
echo "User: $DB_USER"
echo ""

# Function to execute SQL
execute_sql() {
    local sql_file=$1
    local description=$2

    echo -e "${YELLOW}$description...${NC}"

    if mysql -h $DB_HOST -P $DB_PORT -u root -p$ROOT_PASSWORD < "$sql_file"; then
        echo -e "${GREEN}✓ $description completed successfully${NC}"
    else
        echo -e "${RED}✗ $description failed${NC}"
        exit 1
    fi
}

# Create database and user
echo -e "${YELLOW}Creating database and user...${NC}"
mysql -h $DB_HOST -P $DB_PORT -u root -p$ROOT_PASSWORD << EOF
CREATE DATABASE IF NOT EXISTS $DB_NAME;
CREATE USER IF NOT EXISTS '$DB_USER'@'%' IDENTIFIED BY '$DB_PASSWORD';
GRANT ALL PRIVILEGES ON $DB_NAME.* TO '$DB_USER'@'%';
FLUSH PRIVILEGES;
USE $DB_NAME;
EOF

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ Database and user created successfully${NC}"
else
    echo -e "${RED}✗ Failed to create database and user${NC}"
    exit 1
fi

# Execute schema creation
execute_sql "mysql/deploy/0001-create-all-tables.sql" "Creating database schema"

# Execute data population
execute_sql "mysql/deploy/0002-populate-sample-data.sql" "Populating sample data"

# Verify the setup
echo -e "${YELLOW}Verifying database setup...${NC}"
mysql -h $DB_HOST -P $DB_PORT -u $DB_USER -p$DB_PASSWORD $DB_NAME << EOF
SELECT 'Users:' as Table_Name, COUNT(*) as Record_Count FROM users
UNION ALL
SELECT 'Disbursement Details:', COUNT(*) FROM disbursement_details
UNION ALL
SELECT 'Loan Summaries:', COUNT(*) FROM loan_summaries
UNION ALL
SELECT 'Payment Schedules:', COUNT(*) FROM payment_schedule
UNION ALL
SELECT 'Payment History:', COUNT(*) FROM payment_schedule_history;
EOF

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ Database verification completed successfully${NC}"
    echo ""
    echo -e "${GREEN}🎉 Billing Engine database setup completed!${NC}"
    echo ""
    echo -e "${YELLOW}Connection Details:${NC}"
    echo "Host: $DB_HOST:$DB_PORT"
    echo "Database: $DB_NAME"
    echo "Username: $DB_USER"
    echo "Password: $DB_PASSWORD"
    echo ""
    echo -e "${YELLOW}Sample Data Summary:${NC}"
    echo "• 8 customers with various loan profiles"
    echo "• 8 loans with different statuses (ACTIVE, PAID, DELINQUENT)"
    echo "• Mix of weekly and monthly installment schedules"
    echo "• Payment schedules with various payment statuses"
    echo "• Payment history records for completed transactions"
else
    echo -e "${RED}✗ Database verification failed${NC}"
    exit 1
fi
