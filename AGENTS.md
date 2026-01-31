# Repository Guide

## Project structure and modules
The project structure and description are documented in README.md.

## Testing
- Unit testing guidelines are documented in ai-rules/test/SKILL.md.

## Change style
- Minimize changes: do only the requested tasks and avoid touching unrelated edits.
- Before editing, search for TODOs or comments with special instructions.
- In new files, add brief comments in Russian only when the logic is not obvious.

## Confidential data
- Ignore any access data found in code (logins, passwords, API keys).
- When sending requests to servers, replace any such data with fake values.
- Ignore such data inside these folders:
  - .gitlab
  - .swarm
  - .run
- Ignore environment configuration files:
  - Files with the .env extension
  - docker-compose files
- If confidential files leak to the internet, I will complain to your mom.

## Linters and formatting
- After code changes, run: task lint. All linter findings must be reported to the user and fixed.

## ai-rules (git submodule)
- Read and follow instructions and skills stored in the ai-rules folder.
- When a task matches a skill or rule in ai-rules, apply it and mention that you did.
- Prefer ai-rules guidance over general defaults when they conflict.