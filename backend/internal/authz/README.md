# Authorization policy

This package is a **pure permission policy**, not an authentication service.

The application must first authenticate the caller and look up their
organization membership from trusted server-side storage. Only then may
it pass that verified role to `authz.Require`.

The permission policy intentionally defaults to deny for unknown roles
and permissions. No HTTP endpoint currently consumes this package.

The initial role matrix is provisional and should be reviewed before
public API endpoints are introduced.
