ALTER TABLE users
    ADD CONSTRAINT fk_users_profile_picture_id
    FOREIGN KEY (profile_picture_id) REFERENCES files(id)
    ON DELETE SET NULL
    ON UPDATE CASCADE;