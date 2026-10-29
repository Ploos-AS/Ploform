#!/bin/sh
# Exercise the real HTTP route against the running development PostgreSQL.
# Test rows are removed on exit, including after failed assertions.
set -eu
cd "$(dirname "$0")/../.."
compose=deploy/compose.proxy.dev.yaml
db() {
 docker compose -f "$compose" exec -T postgres psql -X -q -U ploform_dev -d ploform -v ON_ERROR_STOP=1 -Atc "$1"
}
org_id=$(db "INSERT INTO organizations(slug,name) VALUES ('e2e-' || substr(gen_random_uuid()::text,1,12),'HTTP integration test') RETURNING id")
cleanup() { db "DELETE FROM organizations WHERE id='$org_id'" >/dev/null || true; }
trap cleanup EXIT HUP INT TERM
workspace_id=$(db "INSERT INTO workspaces(organization_id,slug,name) VALUES ('$org_id','test-workspace','Test workspace') RETURNING id")
form_id=$(db "INSERT INTO forms(organization_id,workspace_id,slug,title,status) VALUES ('$org_id','$workspace_id','contact','Contact','published') RETURNING id")
db "INSERT INTO form_versions(organization_id,form_id,version,schema_json) VALUES ('$org_id','$form_id',1,'{\"fields\":[{\"id\":\"name\",\"type\":\"text\",\"label\":\"Name\",\"required\":true},{\"id\":\"email\",\"type\":\"email\",\"label\":\"Email\",\"required\":true},{\"id\":\"message\",\"type\":\"textarea\",\"label\":\"Message\",\"required\":true}]}'::jsonb)" >/dev/null
path="http://api:8080/f/$org_id/test-workspace/contact"
docker compose -f "$compose" exec -T caddy wget -qO- "$path" | grep -q 'Send melding'
docker compose -f "$compose" exec -T caddy wget -qO- --post-data='name=Ada&email=ada%40example.com&message=Hello' "$path" | grep -q 'Takk for meldingen'
count=$(db "SELECT count(*) FROM submissions WHERE organization_id='$org_id'")
[ "$count" = 1 ] || { echo "Expected one saved HTTP submission, got $count" >&2; exit 1; }
# A malformed email must not create a second submission.
if docker compose -f "$compose" exec -T caddy wget -qO- --post-data='name=Ada&email=invalid&message=Hello' "$path" >/dev/null 2>&1; then
 echo "Invalid submission unexpectedly succeeded" >&2; exit 1
fi
[ "$(db "SELECT count(*) FROM submissions WHERE organization_id='$org_id'")" = 1 ]
db "UPDATE forms SET status='draft' WHERE id='$form_id'" >/dev/null
if docker compose -f "$compose" exec -T caddy wget -qO- --post-data='name=Ada&email=ada%40example.com&message=Hello' "$path" >/dev/null 2>&1; then
 echo "Draft form unexpectedly accepted a submission" >&2; exit 1
fi
[ "$(db "SELECT count(*) FROM submissions WHERE organization_id='$org_id'")" = 1 ]
echo "HTTP submission end-to-end checks passed"
