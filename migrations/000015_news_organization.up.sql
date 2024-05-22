ALTER TABLE news
ADD COLUMN organization_id INT NOT NULL REFERENCES organizations(id) default 1;
