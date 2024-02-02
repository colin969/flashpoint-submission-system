ALTER TABLE submission_file ADD INDEX idx_fk_submission_id_deleted_at (fk_submission_id, deleted_at);
ALTER TABLE submission_cache ADD INDEX idx_bot_action (bot_action);
ALTER TABLE submission_cache ADD INDEX idx_active_assigned_testing_ids (active_assigned_testing_ids);
ALTER TABLE submission_cache ADD INDEX idx_active_assigned_verification_ids (active_assigned_verification_ids);
ALTER TABLE submission_cache ADD INDEX idx_active_requested_changes_ids (active_requested_changes_ids);
ALTER TABLE submission_cache ADD INDEX idx_active_approved_ids (active_approved_ids);
ALTER TABLE submission_cache ADD INDEX idx_active_verified_ids (active_verified_ids);
ALTER TABLE submission_cache ADD INDEX idx_distinct_actions (distinct_actions);
