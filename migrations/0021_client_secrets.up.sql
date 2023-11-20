CREATE TABLE IF NOT EXISTS oauth_client 
(
  client_id VARCHAR(36) PRIMARY KEY,
  client_secret TEXT NOT NULL
);