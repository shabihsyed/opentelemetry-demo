---
name: pr-ready
description: Final checklist before pushing a branch for a PR in this repo (branch name, commit format, changelog entry, scope). Use right before git push.
---

# PR ready checklist

Check each item and fix anything that fails before pushing:

- Branch is named `<issue-number>-<slug>` and is not `main`.
- `git diff main --stat` only touches the service the issue names (plus its lockfiles or tests).
- Each commit message is Conventional Commits with the service as scope, has `Closes #<issue>` in the body,
  and uses an `Assisted-by:` trailer, never `Co-authored-by:`.
- Changelog: user-facing changes need a `.chloggen/<branch>.yaml` entry (copy the shape of an existing file there).
  Test-only, CI-only or comment-only changes skip it; mention in your summary that the PR needs the `Skip Changelog` label.
- The checks for the language you touched were run and passed (for Go: the `go-service-checks` skill).

Then push with `git push -u origin <branch>` and report the branch name.
