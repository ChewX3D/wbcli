# wbcli

[![coverage](https://img.shields.io/endpoint?url=https://raw.githubusercontent.com/ChewX3D/wbcli/main/configs/badges/test-coverage.json)](https://github.com/ChewX3D/wbcli/actions/workflows/badges.yml)
![ARM size](https://img.shields.io/endpoint?url=https://raw.githubusercontent.com/ChewX3D/wbcli/main/configs/badges/binary-size-arm64.json)
![AMD size](https://img.shields.io/endpoint?url=https://raw.githubusercontent.com/ChewX3D/wbcli/main/configs/badges/binary-size-amd64.json)

Lightweight CLI project for placing WhiteBIT trading orders safely, with a design that can later be wrapped by a UI.

## Scope

This repository currently defines the project plan and operating docs for:

- secure API key management
- single collateral limit order placement
- range/batch order placement (step-based ladders)
- clear architecture boundaries so a web/desktop UI can call the same core logic later

## Planned Core Features

- `auth`:
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

## Build

```bash
go build -o bin/wbcli .
```

## Install

```bash
go install github.com/ChewX3D/wbcli@latest
```

## Run

```bash
go run . --help
```

## Suggested Build Direction

Keep core logic reusable:

- `internal/domain`: order plan models and validation
- `internal/app`: use-cases (`PlaceOrder`, `PlaceOrderRange`, `PreviewRange`)
- `internal/adapters/whitebit`: signed HTTP client + endpoint mapping
- `internal/adapters/secretstore`: keychain integration
- `main.go` + `cmd/`: CLI surface and command wiring only

This split lets you add a UI later without duplicating order logic.

## AWS forever free resources

### Compute & Serverless

- AWS Lambda — 1M requests/month + 400,000 GB-seconds
- AWS Step Functions — 4,000 state transitions/month

### Storage & Database

- Amazon DynamoDB — 25 GB storage + 25 read/write capacity units (~200M requests/month)
- Amazon S3 — 5 GB standard storage (note: some sources list this as 12-month only — double-check your billing console)
- Amazon SimpleDB — 25 machine hours + 1 GB storage

### API & Networking

- Amazon API Gateway — 1M REST API calls/month
- Amazon CloudFront — 1 TB data out + 10M HTTP/S requests/month

### Messaging & Queues

- Amazon SQS — 1M requests/month
- Amazon SNS — 1M publishes + 100K HTTP deliveries + 1K email notifications/month
- Amazon SES — 62,000 outbound emails/month (when called from EC2/Beanstalk)

### Auth & Identity

- Amazon Cognito — 50,000 monthly active users

### Monitoring & Management

- Amazon CloudWatch — 10 custom metrics, 10 alarms, 1M API requests, 5 GB log ingestion, 3 dashboards (50 metrics each)
- AWS X-Ray — 100K traces recorded + 1M traces scanned/month
- AWS Config — 7 rules per region + 10K config items recorded/month
- AWS Budgets — 2 budgets (action-enabled)

### Security

- AWS WAF Bot Control — 10M requests
- AWS Security Hub — 30-day trial (per new enablement)
- AWS Resource Access Manager — unlimited

### Other

- Amazon SWF — 1,000 workflow executions + 10K tasks/signals/month
- AWS License Manager — unlimited
- AWS Well-Architected Tool — unlimited
- AWS Service Catalog — 1,000 API calls/month
- Amazon Glacier — 10 GB retrieval/month (storage itself costs money)
