-- Revert billing_engine:0001-create-all-tables from mysql
BEGIN;

DROP TABLE IF EXISTS payment_schedule_histories;
DROP TABLE IF EXISTS payment_schedules;
DROP TABLE IF EXISTS loan_summaries;
DROP TABLE IF EXISTS disbursement_details;
DROP TABLE IF EXISTS users;

COMMIT;
