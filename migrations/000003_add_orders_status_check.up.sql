DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'orders_status_check'
    ) THEN
        ALTER TABLE orders
            ADD CONSTRAINT orders_status_check
            CHECK (status IN ('NEW', 'PROCESSING', 'INVALID', 'PROCESSED'));
    END IF;
END $$;
