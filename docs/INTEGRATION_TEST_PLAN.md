# Internal integration test plan

The development Compose stack must not be exposed publicly. Run tests on an operator-controlled Docker host, with local secret files outside version control.

Before calling the stack ready, verify: all eight migrations applied; API connects as restricted app and auth roles; `/readyz` succeeds only with a working database; Caddy routes `/healthz`; no host ports are published; an unknown proxy peer cannot pass same-origin validation; cross-organization reads fail; a restart does not rerun completed migrations.

The current smoke script checks only a subset of these requirements. It has not been executed in this session because Docker is unavailable in the execution environment. Do not describe the deployment as validated or production-ready.
