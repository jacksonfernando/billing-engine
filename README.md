# Loan Billing System
## Overview
A comprehensive billing system for managing dynamic loan products with flexible repayment schedules. The system supports configurable loan terms, payment frequencies (weekly/monthly), and automated delinquency tracking through schedulers. **This service focuses purely on billing operations and does not perform customer validation - it processes loan requests as received from downstream services.**

## Quick Start

### Prerequisites
- Go 1.21 or later
- Docker and Docker Compose
- Make (optional, but recommended)

### Setup with Makefile (Recommended)
```bash
# Complete setup - provisions all resources, database, and starts services
make setup

# Check status of all services
make status

# Run tests
make test
```

### Manual Setup
```bash
# 1. Start all services (MySQL + Billing Engine)
docker compose up -d --build

# 2. Check status
docker compose ps

# 3. Test API
curl http://localhost:9006/health
```

### Available Make Commands
```bash
make help                 # Show all available commands
make setup               # Complete setup (recommended for first time)
make dev-setup           # Setup development environment
make docker-up           # Start Docker containers
make docker-down         # Stop Docker containers
make db-setup            # Setup database and tables
make db-reset            # Reset database
make build               # Build the application
make run                 # Run the application
make test                # Run all tests
make test-coverage       # Run tests with coverage
make clean               # Clean up build artifacts
make status              # Show status of all services
```

### Configuration
The system uses the following default configuration:
- **Database**: MySQL 8.0 on port 3306
- **Application**: Runs on port 9006
- **Database Name**: billing_engine
- **Database User**: billing_admin

You can override these by setting environment variables before running make commands.

## Business Rules
### Loan Structure (Dynamic)
- **Principal Amount**: Configurable per loan (amount disbursed to customer)
- **Interest Amount**: Configurable per loan
- **Loan Duration**: Configurable number of installments
- **Payment Frequency**: Weekly or Monthly
- **Installment Amount**: (Principal + Interest) ÷ Number of Installments

### Payment Rules
- Borrowers must pay exact amounts to cover overdue installments
- Payments must cover ALL outstanding/overdue installments
- Partial payments are rejected
- Payments are applied to oldest outstanding installments first

### Delinquency Rules
- 2 Overdue repayments

## Database Design (ERD)
```mermaid
erDiagram
    users {
        BIGINT id PK
        VARCHAR customer_id UK
        VARCHAR name
        VARCHAR email
        VARCHAR phone
        TIMESTAMP created_at
        TIMESTAMP updated_at
        TIMESTAMP deleted_at
    }

    disbursement_details {
        BIGINT id PK
        VARCHAR loan_id UK
        VARCHAR customer_id
        TIMESTAMP disbursement_date
        DECIMAL disbursed_amount
        CHAR disbursed_currency
        VARCHAR status
        TIMESTAMP created_at
        VARCHAR created_by
        TIMESTAMP updated_at
        VARCHAR updated_by
        TIMESTAMP deleted_at
    }

    loan_summaries {
        BIGINT id PK
        VARCHAR loan_id UK
        VARCHAR customer_id
        DECIMAL principal_amount
        DECIMAL interest_amount
        DECIMAL outstanding_amount
        INT no_of_installment
        VARCHAR installment_unit
        DECIMAL installment_amount
        DECIMAL effective_interest_rate
        INT dpd
        VARCHAR collectability
        VARCHAR status
        DATE loan_start_date
        TIMESTAMP created_at
        VARCHAR created_by
        TIMESTAMP updated_at
        VARCHAR updated_by
        TIMESTAMP deleted_at
    }

    payment_schedule {
        BIGINT id PK
        VARCHAR loan_id
        INT installment_number
        DECIMAL installment_amount
        DATE installment_due_date
        DECIMAL outstanding_amount
        DECIMAL principal_amount
        DECIMAL principal_paid
        DECIMAL principal_written_off
        DECIMAL interest_amount
        DECIMAL interest_accrued
        DECIMAL interest_paid
        DECIMAL interest_waived_off
        VARCHAR status
        CHAR currency
        INT dpd
        TIMESTAMP created_at
        VARCHAR created_by
        TIMESTAMP updated_at
        VARCHAR updated_by
        TIMESTAMP deleted_at
        VARCHAR deleted_by
    }

    payment_schedule_history {
        BIGINT id PK
        BIGINT schedule_id
        VARCHAR loan_id
        VARCHAR action
        INT installment_number
        DECIMAL installment_amount
        DATE installment_due_date
        DECIMAL outstanding_amount
        DECIMAL principal_amount
        DECIMAL principal_paid
        DECIMAL principal_written_off
        DECIMAL interest_amount
        DECIMAL interest_accrued
        DECIMAL interest_paid
        DECIMAL interest_waived_off
        VARCHAR status
        CHAR currency
        INT dpd
        TIMESTAMP created_at
        VARCHAR created_by
    }

    loan_summaries ||--|| disbursement_details : "loan_id"
    payment_schedule }o--|| disbursement_details : "loan_id"
    payment_schedule_history }o--|| payment_schedule : "schedule_id"
    payment_schedule_history }o--|| disbursement_details : "loan_id"
    users ||--o{ disbursement_details : "customer_id (optional)"
    users ||--o{ loan_summaries : "customer_id (optional)"
```
## Database Schema
### 1. Users Table ( For Reference Only)
```sql
CREATE TABLE users (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    customer_id VARCHAR(36) UNIQUE NOT NULL,
    name VARCHAR(255),
    email VARCHAR(255),
    phone VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

CREATE INDEX idx_users_customer_id ON users (customer_id);
```
*Note: This table is optional and only for reference. The billing system accepts any customer_id provided by downstream services.*

### 2. Disbursement Detail Table
```sql
CREATE TABLE disbursement_details (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    loan_id VARCHAR(36) UNIQUE NOT NULL,
    customer_id VARCHAR(36) NOT NULL,
    disbursement_date TIMESTAMP NOT NULL,
    disbursed_amount DECIMAL(15,2) NOT NULL,
    disbursed_currency CHAR(3) DEFAULT 'IDR',
    status VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(255),
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    updated_by VARCHAR(255),
    deleted_at TIMESTAMP NULL
);

CREATE INDEX idx_disbursement_details_created_at ON disbursement_details (created_at);
CREATE INDEX idx_disbursement_details_status ON disbursement_details (status);
CREATE INDEX idx_disbursement_details_customer_id ON disbursement_details (customer_id);
```

### 3. Loan Summary Table
```sql
CREATE TABLE loan_summaries (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    loan_id VARCHAR(36) UNIQUE NOT NULL,
    customer_id VARCHAR(36) NOT NULL,
    principal_amount DECIMAL(15,2) NOT NULL,
    interest_amount DECIMAL(15,2) NOT NULL,
    outstanding_amount DECIMAL(15,2) NOT NULL,
    no_of_installment INT NOT NULL,
    installment_unit VARCHAR(100) NOT NULL, -- 'week' or 'month'
    installment_amount DECIMAL(15,2) NOT NULL,
    effective_interest_rate DECIMAL(5,4) NOT NULL,
    dpd INT DEFAULT 0,
    status VARCHAR(100) NOT NULL, -- 'PENDING', 'PAID' and 'DELINQUENT'
    loan_start_date DATE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(255),
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    updated_by VARCHAR(255),
    deleted_at TIMESTAMP NULL,
    FOREIGN KEY (loan_id) REFERENCES disbursement_details(loan_id)
);

CREATE INDEX idx_loan_summaries_created_at ON loan_summaries (created_at);
CREATE INDEX idx_loan_summaries_customer_id ON loan_summaries (customer_id);
CREATE INDEX idx_loan_summaries_status ON loan_summaries (status);
CREATE INDEX idx_loan_summaries_dpd ON loan_summaries (dpd);
CREATE INDEX idx_loan_summaries_installment_unit ON loan_summaries (installment_unit);
```

### 4. Payment Schedule Table
```sql
CREATE TABLE payment_schedules (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    loan_id VARCHAR(36) NOT NULL,
    installment_number INT NOT NULL,
    installment_amount DECIMAL(15,2) NOT NULL,
    installment_due_date DATE NOT NULL,
    outstanding_amount DECIMAL(15,2) NOT NULL,
    outstanding_paid DECIMAL(15,2) NOT NULL,
    status VARCHAR(100) DEFAULT 'PENDING', -- 'PENDING' or 'PAID'
    currency CHAR(3) DEFAULT 'IDR',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(255),
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    updated_by VARCHAR(255),
    deleted_at TIMESTAMP NULL,
    deleted_by VARCHAR(255),
    FOREIGN KEY (loan_id) REFERENCES disbursement_details(loan_id)
);

CREATE INDEX idx_payment_schedules_loan_id ON payment_schedules (loan_id);
CREATE INDEX idx_payment_schedules_due_date ON payment_schedules (installment_due_date);
CREATE INDEX idx_payment_schedules_status ON payment_schedules (status);
CREATE INDEX idx_payment_schedules_outstanding ON payment_schedules (outstanding_amount);
CREATE UNIQUE INDEX idx_payment_schedules_loan_installment ON payment_schedules (loan_id, installment_number);
```

### 5. Payment Schedule History Table
```sql
CREATE TABLE payment_schedule_histories (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    schedule_id BIGINT NOT NULL,
    loan_id VARCHAR(36) NOT NULL,
    action VARCHAR(100) NOT NULL,
    installment_number INT NOT NULL,
    installment_amount DECIMAL(15,2) NOT NULL,
    installment_due_date DATE NOT NULL,
    outstanding_amount DECIMAL(15,2) NOT NULL,
    status VARCHAR(100),
    currency CHAR(3) DEFAULT 'IDR',
    dpd INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(255),
    FOREIGN KEY (schedule_id) REFERENCES payment_schedule(id),
    FOREIGN KEY (loan_id) REFERENCES disbursement_details(loan_id)
);

CREATE INDEX idx_payment_schedule_histories_schedule_id ON payment_schedule_histories (schedule_id);
CREATE INDEX idx_payment_schedule_histories_loan_id ON payment_schedule_histories (loan_id);
CREATE INDEX idx_payment_schedule_histories_created_at ON payment_schedule_histories (created_at);
```

## API Specifications
### 1. Disbursement API
**Endpoint**: `POST /v1/disbursement`
**Request Body**:
```json
{
  "principal_amount": 5000000.00,
  "interest_amount": 500000.00,
  "installment_unit": "week",
  "number_of_installment": 50,
  "start_date": "2025-08-31T11:43:00",
  "customer_id": "12312312"
}
```
**Response (Success)**:
```json
{
  "status": "success",
  "data": {
    "loan_id": "loan_123456789",
    "customer_id": "12312312",
    "disbursed_amount": 5000000.00,
    "installment_amount": 110000.00,
    "outstanding_amount": 5500000.00,
    "installment_unit": "week",
    "number_of_installment": 50,
    "disbursement_date": "2025-08-31T11:43:00",
    "first_due_date": "2025-09-07",
    "final_due_date": "2026-08-23"
  }
}
```
**Business Logic**:
1. Validate loan parameters (amounts, installment_unit, number_of_installments)
2. Generate unique loan_id
4. Calculate due dates based on installment_unit
5. Create records in `disbursement_details` and `loan_summaries`
6. Generate payment schedules in `payment_schedule` table with outstanding_amount tracking
7. Set initial outstanding_amount to (principal_amount + interest_amount)

### 2. Repayment API
**Endpoint**: `POST /v1/repayment`
**Request Body**:
```json
{
  "loan_id": "loan_123456789",
  "payment_amount": 220000.00
}
```
**Response (Success)**:
```json
{
  "status": "success",
  "data": {
    "loan_id": "loan_123456789",
    "payment_amount": 220000.00,
    "installments_paid": 2,
    "installment_amount": 110000.00,
    "remaining_installments": 48,
    "outstanding_amount": 5280000.00,
    "next_due_date": "2025-09-21",
    "payment_date": "2025-09-15T10:30:00"
  }
}
```

**Business Logic**:
1. Validate loan exists in billing system
2. Get all overdue/pending installments for the loan (ordered by due date)
3. Calculate total required amount for all overdue installments
4. Validate payment amount equals required amount exactly
5. If valid:
   - Update `payment_schedule` records (mark as PAID)
   - Update outstanding_amount in both `payment_schedule` and `loan_summaries`
   - Create history records in `payment_schedule_history`
6. Each installment in `payment_schedule` tracks its own outstanding amount.
7. Update status from `PENDING` to `PAID` for all the paid installments.
8. If all installment statuses are marked as `PAID`, the loan summary status will be  updated to `PAID`

### Get Outstanding Balance
**Endpoint**: `GET /v1/loans/{loan_id}/outstanding`

**Response**:
```json
{
  "status": "success",
  "data": {
    "loan_id": "loan_123456789",
    "customer_id": "12312312",
    "loan_details": {
      "installment_unit": "week",
      "installment_amount": 110000.00,
      "total_installments": 50
    },
    "outstanding_amount": 3300000.00,
    "overdue_amount": 220000.00,
    "overdue_installments": 2,
    "paid_installments": 20,
    "remaining_installments": 30
  }
}
```

### Check Delinquency Status
**Endpoint**: `GET /v1/loans/{loan_id}/delinquency`

**Response**:
```json
{
  "status": "success",
  "data": {
    "loan_id": "loan_123456789",
    "customer_id": "12312312",
    "is_delinquent": true,
    "installment_unit": "week",
    "overdue_installments": 2,
    "overdue_amount": 220000.00,
    "outstanding_amount": 3300000.00,
    "required_payment_amount": 220000.00
  }
}
```

### Get Loan Schedule
**Endpoint**: `GET /v1/loans/{loan_id}/schedule`

**Response**:
```json
{
  "status": "success",
  "data": {
    "loan_id": "loan_123456789",
    "loan_summaries": {
      "installment_unit": "week",
      "total_installments": 50,
      "installment_amount": 110000.00,
      "disbursed_amount": 5000000.00,
      "outstanding_amount": 3300000.00
    },
    "schedule": [
      {
        "installment_number": 1,
        "due_date": "2025-09-07",
        "installment_amount": 110000.00,
        "outstanding_amount": 5390000.00,
        "status": "PAID",
        "paid_date": "2025-09-07",
      },
      {
        "installment_number": 2,
        "due_date": "2025-09-14",
        "installment_amount": 110000.00,
        "outstanding_amount": 5280000.00,
        "status": "PENDING",
        "paid_date": null,
       }
    ]
  }
}
```
