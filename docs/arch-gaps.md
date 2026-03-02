# Architectural Gap Analysis

Reference document for planned refactoring. Each finding describes a violation,
its severity, and the fix direction. Implementation order is listed at the bottom.
Resolved findings are moved to the bottom for reference.

Severity scale: **high** | **medium** | **low**

---

## Dependency Direction Violations

### GAP-003 — `cmd/auth/errors.go` imports `authservice.ErrNotLoggedIn` from service package
**Severity:** low

```go
// cmd/auth/errors.go
import authservice "github.com/ChewX3D/wbcli/internal/app/services/auth"

{match: authservice.ErrNotLoggedIn, message: "..."},
```

`ErrNotLoggedIn` is a sentinel the cmd layer needs to handle in user-facing messages.
Sentinels at this boundary belong in `internal/app/ports/auth.go` alongside
`ErrCredentialNotFound` and friends — not inside the service implementation package.

**Fix:** Move `ErrNotLoggedIn` declaration from `internal/app/services/auth` to
`internal/app/ports/auth.go`. Update service and cmd imports accordingly.

---

### GAP-004 — `SystemClock` concrete type lives in `internal/app/services/auth/clock.go`
**Severity:** low

```go
// internal/app/services/auth/clock.go
package auth

type SystemClock struct{}

func (SystemClock) Now() time.Time { return time.Now() }
```

`SystemClock` is a concrete implementation of `ports.Clock`. Concrete implementations
of port interfaces belong in `internal/adapters/`, not in a service package.
The side effect is visible in the composition root:

```go
// factory.go — reaching into a service package for a concrete infrastructure type
clock := authservice.SystemClock{}
```

**Fix:** Move `SystemClock` to `internal/adapters/clock/system_clock.go`. Update
`factory.go` (or its replacement in cmd) to import from the new adapter path.

---

## WhiteBIT Transport Client Mirror Rule Violations

### GAP-005 — `PostOnly + IOC` conflict check is a business rule in the transport client
**Severity:** medium

```go
// internal/adapters/whitebit/collateral.go — CollateralLimitOrderRequest.validate()
if request.PostOnly != nil && request.IOC != nil && *request.PostOnly && *request.IOC {
    return ErrPostOnlyIOCConflict
}
```

The transport client mirror rule states the client must be a strict mirror of API
documentation only — no business decisions. The `PostOnly+IOC` mutual exclusion is an
order-semantics constraint, not a transport field shape constraint. Knowing that these two
flags cannot coexist is business/domain knowledge.

The remaining validation in `validate()` (non-empty `Market`, `Amount`, `Price`) is
acceptable as a transport-level guard against malformed HTTP requests. The enum validation
(`ErrInvalidOrderSide`, `ErrInvalidPositionSide`) is borderline — it is currently
called only from `CollateralOrderExecutorAdapter` which already receives typed constants,
so in practice it is a no-op guard.

**Fix:** Remove `ErrPostOnlyIOCConflict` check from `CollateralLimitOrderRequest.validate()`.
Move the constraint check into `CollateralOrderExecutorAdapter.PlaceCollateralLimitOrder`
or into the collateral service before calling the adapter.

---

### GAP-006 — `PlaceCollateralBulkLimitOrder` is orphaned dead API surface
**Severity:** medium

`PlaceCollateralBulkLimitOrder` exists on `*Client` and is covered by tests, but:

- no port interface method exists for bulk orders
- no `CollateralOrderExecutorAdapter` method wraps it
- no service or command uses it

The mirror rule says "nothing more and nothing less than what the API documents".
An unconnected transport method with no product path violates this: it adds maintenance
surface for a feature that does not exist in the product yet.

**Fix:** Either fully wire it up (add port method → adapter method → service → command),
or remove the transport method and its tests until the feature is planned. Do not leave
partial vertical slices dangling at the transport layer.

---

## DRY Violations

### GAP-008 — `boolRef` helper duplicated across service packages
**Severity:** low

```go
// internal/app/services/auth/login.go
// internal/app/services/collateral/place_order.go — identical in both
func boolRef(value bool) *bool {
    allocated := value
    return &allocated
}
```

**Fix:** Extract to a shared internal utility. Options:
- `internal/ptrutil/ptrutil.go` (preferred, explicit package)
- or inline in a shared `internal/app/services/util.go`

---

## Concurrency Safety

### GAP-009 — `applicationFactory` package-level global not mutex-protected
**Severity:** low

```go
// cmd/application_runtime.go
var applicationFactory = appcontainer.NewDefault
```

`SetApplicationFactoryForTest` replaces this global without a lock. Safe for sequential
tests today, but will fail `go test -race ./...` if tests ever run in parallel.

**Fix:** Remove the global. Pass the factory as a parameter into `newRootCmd`. This is
consistent with `NewRootCmdForTest()` which already passes a factory at construction.
The `NewDefault` wiring (after GAP-001 fix) would live directly in `main.go` or
`cmd/application_runtime.go` and be passed in, never stored as a mutable global.

---

## Implementation Order

Fixes are ordered by impact and dependency between findings:

| Order | Finding | Reason |
|-------|---------|--------|
| 1 | GAP-009 | Remove mutable global; pass factory as parameter |
| 2 | GAP-003 | Promote `ErrNotLoggedIn` to ports; small, isolated change |
| 3 | GAP-005 | Move `PostOnly+IOC` rule out of transport client |
| 4 | GAP-006 | Decide: wire bulk orders fully or remove dead surface |
| 5 | GAP-004 | Move `SystemClock` to adapters |
| 6 | GAP-008 | Extract `boolRef` to shared util |

---

## Resolved

### GAP-001 — `internal/app/application/factory.go` imports concrete adapters
**Severity:** ~~high~~ — **closed, accepted as-is**

`internal/app/application/` is the composition root of this project. The composition
root must import all concrete adapter packages — that is its job. `NewDefault()` wires
`configstore`, `secretstore`, and `whitebit` adapters into services; this is correct
and intentional.

The remaining constructors (`New`, `NewWithUseCases`, `NewWithAuthServices`,
`NewWithServices`) accept interfaces only and are fully clean.

**Decision:** Keep `NewDefault()` in `factory.go`. `internal/app/application/` is the
one package inside `internal/` that is explicitly permitted to import concrete adapters,
because it IS the composition root. This exception must be stated in CLAUDE.md:

> `internal/app/application/` — composition root; the only package inside `internal/`
> permitted to import concrete adapter packages; wires adapters into services for the
> runtime application container.

---

### GAP-002 — `cmd/order/errors.go` imports `internal/adapters/whitebit` directly
**Severity:** ~~high~~ — **resolved** (commit `c8f0b0f`)

`cmd/order/errors.go` imported the WhiteBIT adapter package to check what kind of error
came back from a failed order:

```go
// cmd/order/errors.go (before fix)
import "github.com/ChewX3D/wbcli/internal/adapters/whitebit"

case errors.Is(err, whitebit.ErrForbidden):
case errors.Is(err, whitebit.ErrUnauthorized):
```

The problem: the command layer was directly reading WhiteBIT-specific error values.
If the adapter is ever swapped out, renamed, or wrapped behind a different backend,
this import breaks and the cmd code has to change — even though the command itself
has nothing to do with HTTP or WhiteBIT internals.

**Resolution:** Introduced unified `ports.APIError` type. The WhiteBIT adapter now
translates transport errors into `*ports.APIError` at the boundary via `buildAPIError`.
`cmd/order/errors.go` checks `errors.As(err, &apiErr)` — no adapter import needed.

---

### GAP-007 — `indicatesMissingEndpointAccess` and `extractErrorDetail` duplicated
**Severity:** ~~high~~ — **resolved** as part of GAP-002 fix (commit `c8f0b0f`)

Identical string-scraping logic appeared in two packages:

- `cmd/order/errors.go`
- `internal/adapters/whitebit/credential_verifier.go`

Both scraped `"not authorized to perform this action"` from error message strings to
classify missing API endpoint permission errors. The duplication was a direct consequence
of GAP-002: because `cmd/order` bypassed the port abstraction, it had to re-implement
classification already done inside the adapter.

**Resolution:** Fixing GAP-002 eliminated this automatically. The helpers now live only
in `internal/adapters/whitebit/apierror.go`.
