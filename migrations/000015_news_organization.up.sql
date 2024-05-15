ALTER TABLE news
ADD COLUMN organization_id INT NOT NULL REFERENCES organizations(id) default 1,

ALTER COLUMN thumbnail SET NOT NULL,
ALTER COLUMN thumbnail SET DEFAULT 'https://id.pufacomputing.live/news/default.jpg';
