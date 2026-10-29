-- Tenant-isolated webhook outbox. No network delivery or API is enabled.
-- Composite FK requires a unique tenant/id pair on submissions.
ALTER TABLE submissions ADD CONSTRAINT submissions_org_id_unique UNIQUE (organization_id, id);

CREATE TABLE webhook_destinations (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  organization_id uuid NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
  endpoint_url text NOT NULL CHECK (length(endpoint_url) BETWEEN 9 AND 2048 AND endpoint_url LIKE 'https://%'),
  enabled boolean NOT NULL DEFAULT false,
  created_at timestamptz NOT NULL DEFAULT now(),
  UNIQUE (organization_id, id)
);

-- Only metadata and references; no form answers or signing secrets.
CREATE TABLE webhook_outbox (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  organization_id uuid NOT NULL,
  destination_id uuid NOT NULL,
  event_type text NOT NULL CHECK (event_type IN ('submission.created')),
  submission_id uuid NOT NULL,
  status text NOT NULL DEFAULT 'pending' CHECK (status IN ('pending','processing','delivered','failed')),
  attempts integer NOT NULL DEFAULT 0 CHECK (attempts BETWEEN 0 AND 20),
  next_attempt_at timestamptz NOT NULL DEFAULT now(),
  locked_until timestamptz,
  created_at timestamptz NOT NULL DEFAULT now(),
  delivered_at timestamptz,
  last_error_code text,
  UNIQUE (organization_id, destination_id, event_type, submission_id),
  FOREIGN KEY (organization_id, destination_id) REFERENCES webhook_destinations(organization_id, id) ON DELETE CASCADE,
  FOREIGN KEY (organization_id, submission_id) REFERENCES submissions(organization_id, id) ON DELETE CASCADE
);
CREATE INDEX webhook_outbox_pending_idx ON webhook_outbox (next_attempt_at, id) WHERE status = 'pending';

ALTER TABLE webhook_destinations ENABLE ROW LEVEL SECURITY;
ALTER TABLE webhook_destinations FORCE ROW LEVEL SECURITY;
ALTER TABLE webhook_outbox ENABLE ROW LEVEL SECURITY;
ALTER TABLE webhook_outbox FORCE ROW LEVEL SECURITY;
CREATE POLICY webhook_destinations_tenant ON webhook_destinations
 USING (organization_id = nullif(current_setting('app.organization_id', true), '')::uuid)
 WITH CHECK (organization_id = nullif(current_setting('app.organization_id', true), '')::uuid);
CREATE POLICY webhook_outbox_tenant ON webhook_outbox
 USING (organization_id = nullif(current_setting('app.organization_id', true), '')::uuid)
 WITH CHECK (organization_id = nullif(current_setting('app.organization_id', true), '')::uuid);

GRANT SELECT, INSERT, UPDATE, DELETE ON webhook_destinations, webhook_outbox TO ploform_app;
-- A future dedicated worker role must use scoped claims and least-privilege grants.
