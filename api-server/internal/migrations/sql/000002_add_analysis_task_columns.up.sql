-- Add columns for MQ + Callback architecture
ALTER TABLE analyses ADD COLUMN task_id UUID UNIQUE;
ALTER TABLE analyses ADD COLUMN error_code VARCHAR(50);
ALTER TABLE analyses ADD COLUMN retry_count INTEGER NOT NULL DEFAULT 0;
