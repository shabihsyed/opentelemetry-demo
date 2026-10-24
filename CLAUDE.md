# CLAUDE.md

@AGENTS.md

## Working on issues in this fork

- One branch per issue, named `<issue-number>-<short-slug>` (for example `6-bump-currency-otel-cpp`).
- Stay inside the service the issue names (`src/<service>/`). If the fix needs changes elsewhere, stop and say so.
- Commit messages use Conventional Commits with the service as scope (`fix(checkout): ...`, `chore(deps): ...`),
  a body line `Closes #<issue>`, and an `Assisted-by:` trailer as AGENTS.md describes.
- Never push to `main`. Push the branch; a person opens and merges the PR.
- Before committing, run the checks for what you touched. For Go services use the `go-service-checks` skill;
  for version bumps use the `dependency-bump` skill; before pushing use the `pr-ready` skill.
- Changes to the proto in `pb/` regenerate code for every service. Avoid them unless the issue asks for it.
