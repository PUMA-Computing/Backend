CREATE TABLE IF NOT EXISTS version_info (
    id SERIAL PRIMARY KEY,
    latest_version VARCHAR(20) NOT NULL,
    changelog JSONB
);

INSERT INTO version_info (latest_version, changelog)
VALUES (
        '1.0.0',
        '[
          {
            "1.0.0": [
              "## [1.0.10](https://github.com/PUFA-Computing/Frontend/compare/v1.0.9...v1.0.10) (2024-05-22)\n\n\n### Bug Fixes\n\n* using api github get latest release ([b604537](https://github.com/PUFA-Computing/Frontend/commit/b60453758aa7addc45a962500c5cc124f6b24e71))\n\n\n\n"
            ],
            "1.0.1": [
              "Fixed a bug where the app would crash on startup"
            ]
          }
        ]'
       );