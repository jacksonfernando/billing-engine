-- Deploy billing_engine:0001-create-all-tables to mysql
BEGIN;

-- Create users table (for customer reference)
CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    customer_id VARCHAR(36) UNIQUE NOT NULL,
    name VARCHAR(255),
    email VARCHAR(255),
    phone VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_customer_id (customer_id),
    INDEX idx_deleted_at (deleted_at)
);

-- Create disbursement_details table
CREATE TABLE IF NOT EXISTS disbursement_details (
    id INT AUTO_INCREMENT PRIMARY KEY,
    loan_id VARCHAR(50) UNIQUE NOT NULL,
    customer_id VARCHAR(36) NOT NULL,
    disbursement_date DATE NOT NULL,
    disbursed_amount DECIMAL(15,2) NOT NULL,
    disbursed_currency CHAR(3) DEFAULT 'IDR',
    status VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(255),
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    updated_by VARCHAR(255),
    deleted_at TIMESTAMP NULL,
    INDEX idx_loan_id (loan_id),
    INDEX idx_customer_id (customer_id),
    INDEX idx_deleted_at (deleted_at)
);

-- Create loan_summaries table
CREATE TABLE IF NOT EXISTS loan_summaries (
    id INT AUTO_INCREMENT PRIMARY KEY,
    loan_id VARCHAR(50) UNIQUE NOT NULL,
    customer_id VARCHAR(36) NOT NULL,
    principal_amount DECIMAL(15,2) NOT NULL,
    interest_amount DECIMAL(15,2) NOT NULL,
    outstanding_amount DECIMAL(15,2) NOT NULL,
    no_of_installment INT NOT NULL,
    installment_unit VARCHAR(100) NOT NULL,
    installment_amount DECIMAL(15,2) NOT NULL,
    effective_interest_rate DECIMAL(5,4) NOT NULL,
    status VARCHAR(100) NOT NULL,
    loan_start_date DATE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(255),
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    updated_by VARCHAR(255),
    deleted_at TIMESTAMP NULL,
    INDEX idx_loan_id (loan_id),
    INDEX idx_customer_id (customer_id),
    INDEX idx_status (status),
    INDEX idx_deleted_at (deleted_at)
);

-- Create payment_schedule table
CREATE TABLE IF NOT EXISTS payment_schedules (
    id INT AUTO_INCREMENT PRIMARY KEY,
    loan_id VARCHAR(50) NOT NULL,
    installment_number INT NOT NULL,
    installment_amount DECIMAL(15,2) NOT NULL,
    installment_due_date DATE NOT NULL,
    installment_paid DECIMAL(15,2) NOT NULL DEFAULT 0,
    status VARCHAR(100) DEFAULT 'PENDING',
    currency CHAR(3) DEFAULT 'IDR',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(255),
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    updated_by VARCHAR(255),
    deleted_at TIMESTAMP NULL,
    deleted_by VARCHAR(255),
    INDEX idx_loan_id (loan_id),
    INDEX idx_installment_due_date (installment_due_date),
    INDEX idx_status (status),
    INDEX idx_deleted_at (deleted_at),
    UNIQUE KEY unique_loan_installment (loan_id, installment_number)
);

-- Create payment_schedule_history table
CREATE TABLE IF NOT EXISTS payment_schedule_histories (
    id INT AUTO_INCREMENT PRIMARY KEY,
    schedule_id INT NOT NULL,
    loan_id VARCHAR(50) NOT NULL,
    action VARCHAR(100) NOT NULL,
    installment_number INT NOT NULL,
    installment_amount DECIMAL(15,2) NOT NULL,
    installment_due_date DATE NOT NULL,
    status VARCHAR(100),
    currency CHAR(3) DEFAULT 'IDR',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(255),
    INDEX idx_schedule_id (schedule_id),
    INDEX idx_loan_id (loan_id),
    INDEX idx_created_at (created_at)
);

COMMIT;
