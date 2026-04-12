# Repository Guidelines

## Project Structure & Module Organization
`src/` contains the library source, split by concern: `core/` for shared runtime classes, `actions/`, `composites/`, and `decorators/` for node implementations, plus `index.js` barrel exports. `test/` mirrors that structure with Mocha specs such as `test/core/Blackboard.js` and `test/composites/Sequence.js`. `docs/theme/` holds the YUIDoc theme, while generated docs under `docs/behavior3-*` and bundled files in `libs/` are build artifacts and are ignored by git.

## Build, Test, and Development Commands
Run `npm install` once to install the legacy toolchain. Use `npm test` to execute the Mocha suite with Babel transpilation. Use `npx gulp build` to bundle the library into `libs/` and generate versioned YUIDoc output in `docs/`. Use `npx gulp dev` during edits to watch `src/**/*.js` and run JSHint.

## Coding Style & Naming Conventions
Follow the existing ES6 style in `src/`: 2-space indentation, semicolons, and `import`/`export` syntax. Keep class and file names in PascalCase for nodes and core types (`BehaviorTree.js`, `RepeatUntilSuccess.js`); reserve `index.js` for module aggregators. The repository still uses `var` in `gulpfile.js` and older test files, so match the surrounding file instead of mixing styles within one file. Linting is driven by JSHint with `esversion: 6`.

## Testing Guidelines
Tests use Mocha’s TDD interface (`suite`, `test`) with Chai assertions. Add or update tests in the matching path under `test/` whenever behavior changes. Mirror source names when possible, and use focused filenames for special cases, for example `test/core/BehaviorTree-Serialization.js`. Run `npm test` before opening a PR; keep coverage centered on tree execution, blackboard scope, and node edge cases.

## Commit & Pull Request Guidelines
History is minimal, but the existing commit format uses a short imperative subject (`Update README.md`). Keep commit titles concise, capitalized, and under roughly 72 characters. PRs should explain the behavioral change, list affected modules, and include the exact verification command run. Attach screenshots only when documentation output or other rendered artifacts change.

## Build Artifacts & Docs
Do not hand-edit generated files in `libs/` or versioned `docs/behavior3-*`; regenerate them with `npx gulp build`. Edit source files in `src/` and theme assets in `docs/theme/` instead.
