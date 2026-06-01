ALTER TABLE payrolls
    ADD COLUMN gross_salary DECIMAL(15,2) NULL DEFAULT 0 AFTER additional_data;
