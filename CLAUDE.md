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

## Team conventions (suggested by Agent Observatory)

- Before opening a PR, run the full checks for the touched service again. If a check fails, fix the cause and rerun the failing check before you continue. Do not commit or report success while checks fail.
- Before opening the PR, review the diff with the reviewer subagent, or do a self-review if none is available, and list the remaining risks in the PR description.
- Only push the branch and open the PR after both steps pass.
