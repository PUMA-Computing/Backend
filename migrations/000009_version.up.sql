CREATE TABLE IF NOT EXISTS version_info (
    id SERIAL PRIMARY KEY,
    latest_version VARCHAR(20) NOT NULL,
    changelog JSONB
);

INSERT INTO version_info (latest_version, changelog)
VALUES (
        '0.6.0',
        '[{"version": "0.6.0", "changes": ["Initial release"]}]'
       );