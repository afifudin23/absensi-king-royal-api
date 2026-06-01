-- Delete duplicate user_profiles, keep only one (latest by updated_at) per user_id
DELETE FROM user_profiles
WHERE id NOT IN (
    SELECT keep_id FROM (
        SELECT SUBSTRING_INDEX(
            GROUP_CONCAT(id ORDER BY updated_at DESC SEPARATOR ','),
            ',', 1
        ) AS keep_id
        FROM user_profiles
        GROUP BY user_id
    ) AS to_keep
);

-- Add unique constraint to prevent future duplicates
ALTER TABLE user_profiles
    ADD CONSTRAINT uq_user_profiles_user_id UNIQUE (user_id);
