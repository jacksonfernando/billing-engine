-- Verify billing_engine:0001-create-all-tables on mysql
BEGIN;

-- Verify that all tables exist
SELECT 1/COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = 'users';
SELECT 1/COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = 'disbursement_details';
SELECT 1/COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = 'loan_summaries';
SELECT 1/COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = 'payment_schedule';
SELECT 1/COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = 'payment_schedule_history';

-- Verify key indexes exist
SELECT 1/COUNT(*) FROM information_schema.statistics WHERE table_schema = DATABASE() AND table_name = 'users' AND index_name = 'idx_customer_id';
SELECT 1/COUNT(*) FROM information_schema.statistics WHERE table_schema = DATABASE() AND table_name = 'disbursement_details' AND index_name = 'idx_loan_id';
SELECT 1/COUNT(*) FROM information_schema.statistics WHERE table_schema = DATABASE() AND table_name = 'loan_summaries' AND index_name = 'idx_loan_id';
SELECT 1/COUNT(*) FROM information_schema.statistics WHERE table_schema = DATABASE() AND table_name = 'payment_schedule' AND index_name = 'idx_loan_id';

ROLLBACK;
