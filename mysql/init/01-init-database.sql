-- Initialize Billing Engine Database
-- This script runs automatically when MySQL container starts

-- Create database
CREATE DATABASE IF NOT EXISTS billing_engine;

-- Create user and grant privileges
CREATE USER IF NOT EXISTS 'billing_admin'@'%' IDENTIFIED BY 'billing_password';
GRANT ALL PRIVILEGES ON billing_engine.* TO 'billing_admin'@'%';
FLUSH PRIVILEGES;

-- Use the billing_engine database
USE billing_engine;
