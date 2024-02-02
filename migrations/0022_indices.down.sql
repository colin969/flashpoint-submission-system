DROP INDEX idx_fk_submission_id_deleted_at ON submission_file;
DROP INDEX idx_bot_action ON submission_cache;
DROP INDEX idx_active_assigned_testing_ids ON submission_cache;
DROP INDEX idx_active_assigned_verification_ids ON submission_cache;
DROP INDEX idx_active_requested_changes_ids ON submission_cache;
DROP INDEX idx_active_approved_ids ON submission_cache;
DROP INDEX idx_active_verified_ids ON submission_cache;
DROP INDEX idx_distinct_actions ON submission_cache;
