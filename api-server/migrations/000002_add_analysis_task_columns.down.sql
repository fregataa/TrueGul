-- Revert MQ + Callback columns
ALTER TABLE analyses DROP COLUMN task_id;
ALTER TABLE analyses DROP COLUMN error_code;
ALTER TABLE analyses DROP COLUMN retry_count;
