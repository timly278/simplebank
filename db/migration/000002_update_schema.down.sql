

ALTER TABLE accounts ALTER COLUMN id TYPE INT;
ALTER TABLE transfers ALTER COLUMN from_account_id TYPE SET NULL;
ALTER TABLE transfers ALTER COLUMN to_account_id TYPE SET NULL;