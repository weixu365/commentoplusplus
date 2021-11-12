-- Default sort policy for each domain

-- DROP TYPE if exists sortPolicy;

-- CREATE TYPE sortPolicy AS ENUM (
--   'score-desc',
--   'creationdate-desc',
--   'creationdate-asc'
-- );

ALTER TABLE domains
  ADD defaultSortPolicy text NOT NULL DEFAULT 'score-desc';
