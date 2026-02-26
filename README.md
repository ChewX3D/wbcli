# whitbit_cli

Lightweight CLI project for placing WhiteBIT trading orders safely, with a design that can later be wrapped by a UI.

## Scope

This repository currently defines the project plan and operating docs for:

- secure API key management
- single collateral limit order placement
- range/batch order placement (step-based ladders)
- clear architecture boundaries so a web/desktop UI can call the same core logic later

## Planned Core Features

- `keys`:
  - store API credentials in OS keychain/secret store where possible
  - never persist raw secrets in git-tracked files
  - support profile-based credentials (for multiple accounts/environments)
- `order place`:
  - place one collateral limit order via WhiteBIT authenticated API
- `order range`:
  - generate and place a range of limit orders from `start` to `end` with `step`
  - amount modes:
    - constant amount per step
    - arithmetic progression (`x1 x2 x3 ...`)
    - geometric progression (`x1 x2 x4 x8 ...`)
    - capped geometric progression (risk cap)
    - optional later variants: fibonacci, custom multiplier list
  - `--dry-run` preview before submission

## Docs

- [Product + CLI design](docs/cli-design.md)
- [WhiteBIT API integration notes](docs/whitebit-integration.md)
- [Project operating system / workflow rules](AGENTS.md)
- [Ticket workflow and commands](tickets/README.md)

## Suggested Build Direction

Keep core logic reusable:

- `internal/domain`: order plan models and validation
- `internal/app`: use-cases (`PlaceOrder`, `PlaceOrderRange`, `PreviewRange`)
- `internal/adapters/whitebit`: signed HTTP client + endpoint mapping
- `internal/adapters/secretstore`: keychain integration
- `cmd/whitbit`: CLI surface only

This split lets you add a UI later without duplicating order logic.
