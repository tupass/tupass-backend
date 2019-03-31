# Contributing

This project uses the Common Flow workflow which is described at https://commonflow.org in its newest version (currently [`1.0.0-rc.5`](https://commonflow.org/spec/1.0.0-rc.5.html)).

The following is a short summary of the most important core principles.


## General
* Commit early and often
* Use descriptive commit messages

## Master branch
* Is protected (no direct commits)
* Only takes changes via pull requests (or merge requests in GitLab)

## Change branches
* Development happens in so-called change branches
* Based on the current `master` (or on other change branches for bigger changes)
* Give them a meaningful and descriptive name
* Rebase early and often (only fast-forward merges allowed)
* Can be pushed as soon as it's created (good for discussion along the code)
  * Start the title with a `WIP: ` while you're still working on it

## Change merge
* Only merge when CI/CD pipeline ran through
* Only fast-forward merge (if merging is blocked, rebase the change branch on the current `master`)
  * Makes sure that the CI/CD pipeline will run after merging
* Shouldn't be merged by the developer who opened the pull request
  * The merging developer should review

## Releases
* Tagged commits on `master`
* Tag is the exact version number (e.g. `1.0.0`)


---
For the precise processes see [`PROCESSES.md`](PROCESSES.md).

For some checklists for processes see [`CHECKLISTS.md`](CHECKLISTS.md).
