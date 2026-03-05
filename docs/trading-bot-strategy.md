# BTC/USDT Futures Grid Trading Bot — Strategy Document

## Overview

A hedged grid trading bot for BTC/USDT perpetual futures on WhiteBit exchange, written in Go. The bot places maker-only limit orders above and below the grid anchor, profiting from BTC's natural price oscillations. It uses hedge mode to hold long and short positions simultaneously, an EMA trend filter to bias grid direction, and a multi-layer risk management system to survive black swan events.

**One-sentence summary:** Place buy orders below the grid anchor and sell orders above it, profit from BTC bouncing through the grid levels, and protect yourself if price moves outside the grid range.

## Exchange & Fees

- **Exchange:** WhiteBit
- **Pair:** BTC/USDT perpetual futures
- **Mode:** Hedge mode (long and short positions simultaneously)
- **Maker fee:** 0.01% (our edge — most exchanges charge 0.02-0.025%)
- **Taker fee:** 0.055% (avoid at all costs)
- **All orders must be `post_only: true`** to guarantee maker fee

## Grid Terminology

| Term | Meaning |
|---|---|
| **Grid anchor** | The center price around which the grid is built. Set on startup to current market price. Updated by the rebalancing algorithm when price drifts outside the grid range. No order is placed at the anchor itself — it's a reference point, not a level. |
| **Grid step** | The fixed price distance between adjacent grid levels (e.g. \$200). |
| **Grid level** | A specific price where an order sits. Each level is `anchor ± (N × step)` where N is the level number (1, 2, 3...). |
| **Long levels** | Grid levels below the anchor where buy (open long) orders are placed. Labeled L1 (closest to anchor), L2, L3, etc. |
| **Short levels** | Grid levels above the anchor where sell (open short) orders are placed. Labeled S1 (closest to anchor), S2, S3, etc. |
| **Grid range** | The total price span from the lowest long level to the highest short level. With 3 levels per side and \$200 step: range = 6 × \$200 = \$1,200. |
| **Round trip** | A complete entry fill + take-profit fill on the same grid level. One round trip = one profit capture. |

## Account Parameters

| Parameter | Starting Value | Notes |
|---|---|---|
| Account size | \$100 | Only keep trading capital on exchange |
| Leverage | 5x | Buying power: \$500 |
| Grid step | \$200 | Tighten as account grows |
| Grid levels | 3 per side (6 total) | Increase as account grows |
| Position size | ~0.001 BTC per level | Adjust to stay within margin |

### Scaling Guide

```
$100 account  → $200 grid step, 3 levels per side
$200 account  → $150 grid step, 4 levels per side
$500 account  → $100 grid step, 5 levels per side
$2,000+ account → consider adding ETH/USDT as second grid
$10,000+ account → diversify across exchanges and strategies
```

## Core Strategy: Hedged Grid

### Grid Setup

On startup, the bot sets the grid anchor to the current market price and places maker limit orders at each grid level:

```
Example: grid anchor = $60,000, grid step = $200, 3 levels per side

Short levels (limit sells above anchor):
  S1  $60,200 — open short (take profit at $60,000)
  S2  $60,400 — open short (take profit at $60,200)
  S3  $60,600 — open short (take profit at $60,400)

          --- $60,000 (grid anchor, no order here) ---

Long levels (limit buys below anchor):
  L1  $59,800 — open long (take profit at $60,000)
  L2  $59,600 — open long (take profit at $59,800)
  L3  $59,400 — open long (take profit at $59,600)
```

### How It Profits

BTC oscillates through the grid levels — entry orders fill on both sides, take-profits close them at the next level. Each completed round trip captures one grid step as gross profit.

### Profit Per Round Trip

```
Grid step: $200
Position size: 0.001 BTC
Gross profit per round trip: $0.20
Maker fees (0.01% × 2 fills): ~$0.012
Net profit per round trip: ~$0.19
```

### Core Grid Algorithm

```
on_startup:
  anchor = get_current_price()
  for i = 1 to NUM_LEVELS:
    place_buy(anchor - step * i)     // long levels
    place_sell(anchor + step * i)    // short levels

on_order_filled(order):
  if order.side == BUY:
    next_side = SELL
    next_price = order.price + step
  if order.side == SELL:
    next_side = BUY
    next_price = order.price - step

  if order.is_grid:
    // Entry filled → place take profit at the next level
    place_order(next_side, next_price, type=TAKE_PROFIT)
  if order.is_take_profit:
    // Take profit filled → restore the consumed grid level
    place_order(next_side, next_price, type=GRID)
```

The fill/take-profit cycle automatically restores consumed grid levels as long as price stays within the grid range. Each completed round trip returns the level to its original state with profit captured. The bot does not predict price direction — it profits from price oscillating through the grid levels.

When price drifts outside the grid range, the core algorithm can no longer generate fills. The rebalancing algorithm (see Grid Rebalancing section) handles this by moving the grid anchor to the current price.

## Trend Filter: EMA(50) on 15-Minute Candles

The grid runs with a directional bias based on a 50-period EMA calculated on 15-minute candlestick closes.

### Logic

```
Price > EMA → Bullish bias
  Long levels:  3 (full allocation)
  Short levels: 2 (reduced)

Price < EMA → Bearish bias
  Short levels: 3 (full allocation)
  Long levels:  2 (reduced)

Price ≈ EMA → Neutral
  Both sides: equal levels
```

### Why EMA

- SMA is too slow — confirms trends too late
- DEMA/HMA are too fast — flip bias too often, causing unnecessary order churn
- EMA is the sweet spot for a trend filter — reacts fast enough to catch real trends, smooth enough to ignore noise

### EMA Calculation

```
multiplier = 2 / (period + 1)
EMA = (price - previous_EMA) × multiplier + previous_EMA
```

On startup: fetch 50 candles via REST API, seed EMA by iterating through them. Then update with each new 15-minute candle close from WebSocket.

### Data Source

WhiteBit API provides raw kline data:
```
GET /api/v4/public/kline?market=BTC_USDT&interval=15m&limit=50
Response: [timestamp, open, close, high, low, volume]
```

EMA is calculated locally — no exchange provides pre-calculated indicators.

## Risk Management

### Three-Layer Stop-Loss System

```
Layer 1 — Bot logic (smart):
  At -$59,000: hedge lock with opposite position (maker order)
  Tries limit orders first, intelligent exit

Layer 2 — Exchange stop-loss (dumb but reliable):
  At -$58,800: market sell order sitting in exchange matching engine
  Fires only if bot fails to act (API down, bot crash)

Layer 3 — Low leverage (airbag):
  5x leverage, liquidation at ~$48,000
  Survives even if both Layer 1 and Layer 2 fail temporarily
```

### Hedge Lock Mechanism (Black Swan Protection)

Instead of a hard stop-loss that realizes a loss permanently, the bot opens an equal opposite position to freeze the loss.

```
Normal stop-loss:
  Price hits SL → market sell → loss is REALIZED and FINAL
  If price bounces back, you already sold. Missed recovery.

Hedge lock:
  Price hits threshold → open equal short (maker order)
  Loss is FROZEN but NOT realized
  Wait up to 48 hours for bounce:
    - Bounce happens → close short, ride long back up (potential full recovery)
    - No bounce in 48h → close both, accept loss (same outcome as stop-loss)
```

**Hedge lock wins when:** price bounces (most common for BTC), avoids selling at worst moment, avoids taker fees.

**Stop-loss wins when:** price never recovers, extended sideways (funding fees accumulate).

**For BTC specifically, hedge lock has an edge** because BTC almost always bounces to some degree after crashes.

### Circuit Breakers

```
Level 1: Unrealized PnL < -3% of account
  → Stop opening new positions
  → Hedge lock any open positions

Level 2: Realized + Unrealized PnL < -5% daily
  → Close all positions
  → Bot pauses 24 hours

Level 3: Weekly drawdown > 12%
  → Bot pauses 7 days
  → Send alert (Telegram notification)
```

### Position Sizing Rules

- Max 1-2% loss per trade
- Max 5% daily drawdown
- Max 12% weekly drawdown
- Never risk more than you can lose without it affecting your life

## Architecture

```
┌──────────────────────────────────────────────────┐
│                WebSocket Feed                     │
│               BTC/USDT price                      │
└──────────┬───────────────────────┬───────────────┘
           │                       │
    FAST LOOP (every tick)   SLOW LOOP (15-min candle close)
           │                       │
    ┌──────▼──────┐         ┌──────▼──────┐
    │ Order Fill  │         │ Update EMA  │
    │ Monitor     │         │ Trend Filter│
    │ → place TPs │         └──────┬──────┘
    └──────┬──────┘                │
           │                ┌──────▼──────┐
    ┌──────▼──────┐         │ Rebalance?  │
    │ Circuit     │         │ Adjust bias │
    │ Breaker     │         │ Move anchor │
    │ Monitor     │         └─────────────┘
    └──────┬──────┘
           │
    ┌──────▼──────┐
    │ Flash Crash │
    │ Detector    │
    │ (3+ levels) │
    │ → emergency │
    │   rebalance │
    └──────┬──────┘
           │
    ┌──────▼──────┐
    │ Hedge Lock  │
    │ (if needed) │
    └──────┬──────┘
           │
    ┌──────▼──────┐
    │ PnL Tracker │
    └─────────────┘
```

## Technical Implementation Notes

### WhiteBit API

- REST API for order management, balance queries, kline data
- WebSocket API for real-time price feed and kline stream
- Auth: API key + secret with HMAC-SHA512 signing, nonce-based
- No official Go SDK — wrap REST/WS endpoints manually
- Docs: https://docs.whitebit.com/

### Key Implementation Details

- **Maker-only enforcement:** Set `post_only: true` on every order. Order gets rejected rather than filled as taker.
- **Rate limiting:** Use token bucket or `time.Ticker` to respect API limits.
- **WebSocket reconnection:** Auto-reconnect with exponential backoff. Connections will drop.
- **State persistence:** Track open positions and grid state in SQLite or Postgres. Must survive restarts.
- **Paper trading mode:** Run against real market data with simulated orders before going live.

### Grid Rebalancing

When price drifts outside the grid range, orders on the far side become stale and the grid stops generating round trips. The bot must detect this and move the grid anchor to the current price.

#### Why Rebalancing Is Needed

```
Grid anchor: $60,000, grid step: $200. BTC trends to $62,000.

  S3  $60,600 — filled, TP waiting at $60,400 (unrealized loss)
  S2  $60,400 — filled, TP waiting at $60,200 (unrealized loss)
  S1  $60,200 — filled, TP waiting at $60,000 (unrealized loss)
      ------- price is $62,000 (outside grid range) -------
  L1  $59,800 — will never fill at this price
  L2  $59,600 — will never fill
  L3  $59,400 — will never fill

Grid is dead. All short positions are losing. Long levels are unreachable.
No round trips can complete.
```

#### Two Processing Speeds

The bot operates at two different speeds:

**Fast loop (every WebSocket tick):**
- Monitor order fills → place take-profits instantly
- Monitor circuit breaker thresholds → react instantly
- Detect flash crash / price gaps 3+ grid steps beyond the grid range → emergency rebalance instantly

**Slow loop (every 15-minute candle close):**
- Update EMA
- Check if price is outside the grid range
- Adjust long/short level allocation based on trend filter
- Move anchor if needed

This separation keeps the bot responsive to fills and emergencies while avoiding unnecessary grid restructuring on every price tick.

#### Rebalancing Algorithm

```
on_candle_close_15min(price):
  update_ema(candle.close)
  levels_beyond = calculate_levels_beyond_grid_range(price)

  if levels_beyond == 0:
    // Price inside grid range. Only adjust trend bias.
    adjust_grid_bias(ema)
    return

  if levels_beyond == 1:
    // Slightly outside grid range. Shift anchor by 1 step.
    // Cancel farthest level on opposite side.
    // Place new level 1 step closer to current price.
    // Don't close losing positions yet.
    shift_anchor_one_step(direction)
    return

  if levels_beyond >= 2:
    // Clearly trending. Full rebalance — move anchor to current price.
    cancel_all_stale_orders()
    if losing_positions_exceed(2% of account):
      hedge_lock(losing_positions)
    else:
      close(losing_positions)  // accept small loss
    place_new_grid(current_price)  // current price becomes the new anchor
    return
```

**Emergency rebalance (on any WebSocket tick):**
```
on_price_update(price):
  levels_beyond = calculate_levels_beyond_grid_range(price)

  if levels_beyond >= 3:
    // Flash crash or pump. Don't wait for candle.
    emergency_rebalance(price)  // move anchor to current price immediately
    return
```

#### Gradual Shift Example

```
Grid anchor: $60,000, grid step: $200

15-min candle closes at $60,300 (1 step beyond S1):
  → Cancel L3 at $59,400 (farthest stale long level)
  → Place new level at $60,200
  → Anchor shifts to $60,200 without closing any positions

Next candle closes at $60,500 (still 1 step beyond new grid edge):
  → Cancel farthest stale long level
  → Place new level near current price
  → Anchor shifts again — grid gradually follows the trend

This avoids realizing losses on temporary moves while keeping
the grid alive during real trends.
```

#### Why Not React Instantly to Every Tick

```
BTC at $60,000:
  12:00:00 — price spikes to $60,250 (beyond S1)
  12:00:03 — price drops back to $60,150 (inside grid range)

  Instant reaction: moved anchor for no reason, wasted API calls
  Wait for 15-min candle: did nothing. Correct decision.
```

Tying rebalancing to 15-minute candle closes naturally filters out wicks and noise. The only exception is flash crashes (3+ grid steps beyond the grid range) which require immediate action.

## Expected Performance

### Per-Trade Metrics

```
Win rate per round trip: ~85%
Profit per winning round trip: ~$0.19 (at $100 account)
Loss per losing trade: capped by circuit breakers
```

### Monthly Estimates (not compounding)

```
Conservative (low volatility):  $18-$36/month  (MRR 20-35%)
Average month:                  $48-$90/month  (MRR 50-90%)
Great month (high volatility):  $112-$187/month (MRR 110-180%)
Bad month (black swan):         -$5 to +$10    (MRR -5% to +10%)
```

### Growth Projection (with compounding)

```
Month 1:  $100  → $150
Month 2:  $150  → $210
Month 3:  $210  → $180  (bad month)
Month 4:  $180  → $260
Month 5:  $260  → $370
Month 6:  $370  → $320  (pullback)
Month 7:  $320  → $450
Month 8:  $450  → $600
Month 12: ~$800-$1,500

Target: $100 → $500 in 8-12 months
Realistic ARR: 400-500% (accounting for bad months)
```

### Important Caveats

- These estimates assume average BTC volatility. Low-volatility periods will underperform.
- You will have losing months. The strategy wins in aggregate, not on every trade.
- Percentage returns decrease as account size grows (fewer fills relative to capital, orderbook depth limits).
- Past BTC volatility patterns may not continue.

## Three Outcomes

```
1. BTC oscillates within grid range (70% of days) → round trips complete on both sides → profit
2. BTC trends slowly (20% of days)  → rebalancing moves anchor, trend filter helps, reduced profit
3. BTC crashes/pumps hard (10%)     → hedge lock + circuit breaker → small locked loss
```

## Pre-Launch Checklist

- [ ] Implement paper trading mode and run for 2+ weeks
- [ ] Verify all orders are maker-only (check fill reports)
- [ ] Test circuit breakers with simulated crashes
- [ ] Test hedge lock mechanism manually
- [ ] Verify WebSocket reconnection works reliably
- [ ] Set up Telegram alerts for circuit breaker triggers
- [ ] Confirm WhiteBit minimum order sizes for BTC/USDT futures
- [ ] Start with minimum position sizes, scale up gradually
- [ ] Keep only trading capital on exchange, rest in cold wallet