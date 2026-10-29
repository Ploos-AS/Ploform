-- Initial multi-tenant schema. All tenant-owned rows carry organization_id.
CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE organizations (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  slug text NOT NULL UNIQUE CHECK (slug ~ '^[a-z0-9][a-z0-9-]{1,62}$'),
  name text NOT NULL CHECK (length(trim(name)) > 0),
  created_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE workspaces (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  organization_id uuid NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
  slug text NOT NULL CHECK (slug ~ '^[a-z0-9][a-z0-9-]{1,62}$'),
  name text NOT NULL CHECK (length(trim(name)) > 0),
  created_at timestamptz NOT NULL DEFAULT now(),
  UNIQUE (organization_id, id),
  UNIQUE (organization_id, slug)
);

CREATE TABLE forms (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  organization_id uuid NOT NULL,
  workspace_id uuid NOT NULL,
  slug text NOT NULL CHECK (slug ~ '^[a-z0-9][a-z0-9-]{1,62}$'),
  title text NOT NULL CHECK (length(trim(title)) > 0),
  status text NOT NULL DEFAULT 'draft' CHECK (status IN ('draft','published','archived')),
  created_at timestamptz NOT NULL DEFAULT now(),
  UNIQUE (organization_id, id),
  UNIQUE (workspace_id, slug),
  FOREIGN KEY (organization_id, workspace_id) REFERENCES workspaces(organization_id, id) ON DELETE CASCADE
);

CREATE TABLE form_versions (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  organization_id uuid NOT NULL,
  form_id uuid NOT NULL,
  version integer NOT NULL CHECK (version > 0),
  schema_json jsonb NOT NULL CHECK (jsonb_typeof(schema_json) = 'object'),
  created_at timestamptz NOT NULL DEFAULT now(),
  UNIQUE (organization_id, id),
  UNIQUE (form_id, version),
  FOREIGN KEY (organization_id, form_id) REFERENCES forms(organization_id, id) ON DELETE CASCADE
);

CREATE TABLE submissions (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  organization_id uuid NOT NULL,
  form_version_id uuid NOT NULL,
  response_json jsonb NOT NULL CHECK (jsonb_typeof(response_json) = 'object'),
  created_at timestamptz NOT NULL DEFAULT now(),
  FOREIGN KEY (organization_id, form_version_id) REFERENCES form_versions(organization_id, id) ON DELETE CASCADE
);
CREATE INDEX submissions_org_created_idx ON submissions (organization_id, created_at DESC);

-- Database-side isolation as defense in depth. Application roles must not own tables
-- or have BYPASSRLS; application must set app.organization_id in each transaction.
ALTER TABLE workspaces ENABLE ROW LEVEL SECURITY;
ALTER TABLE forms ENABLE ROW LEVEL SECURITY;
ALTER TABLE form_versions ENABLE ROW LEVEL SECURITY;
ALTER TABLE submissions ENABLE ROW LEVEL SECURITY;
ALTER TABLE workspaces FORCE ROW LEVEL SECURITY;
ALTER TABLE forms FORCE ROW LEVEL SECURITY;
ALTER TABLE form_versions FORCE ROW LEVEL SECURITY;
ALTER TABLE submissions FORCE ROW LEVEL SECURITY;

CREATE POLICY workspace_tenant ON workspaces
  USING (organization_id = nullif(current_setting('app.organization_id', true), '')::uuid)
  WITH CHECK (organization_id = nullif(current_setting('app.organization_id', true), '')::uuid);
CREATE POLICY form_tenant ON forms
  USING (organization_id = nullif(current_setting('app.organization_id', true), '')::uuid)
  WITH CHECK (organization_id = nullif(current_setting('app.organization_id', true), '')::uuid);
CREATE POLICY version_tenant ON form_versions
  USING (organization_id = nullif(current_setting('app.organization_id', true), '')::uuid)
  WITH CHECK (organization_id = nullif(current_setting('app.organization_id', true), '')::uuid);
CREATE POLICY submission_tenant ON submissions
  USING (organization_id = nullif(current_setting('app.organization_id', true), '')::uuid)
  WITH CHECK (organization_id = nullif(current_setting('app.organization_id', true), '')::uuid);
