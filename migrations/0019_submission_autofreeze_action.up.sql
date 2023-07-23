ALTER TABLE submission ADD fk_autofreeze_action_id BIGINT DEFAULT NULL;
ALTER TABLE submission ADD CONSTRAINT fk_fk_autofreeze_action_id FOREIGN KEY (fk_autofreeze_action_id) REFERENCES action (id);