CREATE TYPE job_status AS ENUM (
  'pending',
  'processing',
  'done',
  'failed'
);

CREATE TABLE job (
  id BIGSERIAL PRIMARY KEY,
  payload JSONB NOT NULL DEFAULT '{}'::jsonb,
  result JSONB,
  status job_status NOT NULL DEFAULT 'pending',
  error TEXT,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  started_at TIMESTAMP,
  finished_at TIMESTAMP
);

-- useful indexes for worker polling
CREATE INDEX idx_jobs_status ON job(status);
CREATE INDEX idx_jobs_created_at ON job(created_at);
