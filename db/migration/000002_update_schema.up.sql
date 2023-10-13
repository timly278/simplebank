
ALTER TABLE accounts ALTER COLUMN id TYPE bigint;
ALTER TABLE transfers ALTER COLUMN from_account_id SET NOT NULL;
ALTER TABLE transfers ALTER COLUMN to_account_id SET NOT NULL;
