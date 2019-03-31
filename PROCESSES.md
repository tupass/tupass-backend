# Processes

## Developing a feature
1. Create a change branch with `git checkout -b <your-cool-change-branch-name>`
2. Commit code with `git add <file1> <file2> && git commit -m "<your-helpful-commit-message>"`
3. Push the change branch upstream with `git push origin <your-cool-change-branch-name>` (you can do this whenever you like)
4. Raise a pull request and describe the new feature
5. Wait for someone to merge it / request changes

## Fixing a bug
1. Create a change branch with `git checkout -b <your-cool-change-branch-name>`
2. Commit a few fixes with `git add <file1> <file2> && git commit -m "<your-helpful-commit-message>"`
3. Push the change branch upstream with `git push origin <your-cool-change-branch-name>`
4. Raise a pull request and describe the bug and your fix there
5. Wait for someone to merge it / request changes

## Releasing a new version
1. Bump the version to the new one
2. Create a new git tag with `git tag <version-number>`
3. Push the new tag with `git push --tags`
