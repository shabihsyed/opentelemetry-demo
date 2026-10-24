---
name: dependency-bump
description: Find and bump one pinned dependency for a service (base image digest, build ARG, .env version, go.mod/package manifest). Use for dependency-update issues.
---

1. Find where the service's versions are pinned. Look in this order:
   - `src/<service>/Dockerfile` (`FROM ...@sha256:...`, `ARG ..._VERSION`)
   - `.env` and `.env.override` (`*_VERSION` variables passed as build args)
   - `src/<service>/go.mod`, `package.json`, `requirements*.txt`, `*.csproj`, `CMakeLists.txt`
   - `.github/dependabot.yml` (shows which pins Dependabot already manages)
2. Pick one pin that Dependabot does not already update, unless the issue names one.
3. Check the newest stable release from the upstream project's release notes before choosing a version. Prefer the latest patch of the current minor; call out any breaking change.
4. Change only that pin (plus lockfiles it regenerates). If it is a Go module: `go get <module>@<version> && go mod tidy`, then run the `go-service-checks` skill.
5. If the service can't be built locally (for example C++ services need Docker), say so in the commit body and in your summary.
6. Commit as `chore(deps): bump <dependency> to <version> in <service>`.
