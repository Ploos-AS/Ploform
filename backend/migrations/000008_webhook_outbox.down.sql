REVOKE ALL ON webhook_outbox, webhook_destinations FROM ploform_app;
DROP TABLE webhook_outbox;
DROP TABLE webhook_destinations;
ALTER TABLE submissions DROP CONSTRAINT submissions_org_id_unique;
