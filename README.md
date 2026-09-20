# Orange Mango — QR Ordering System

A QR-code based ordering system for a single burger shop, built in Go. Customers scan a QR code at their table, browse the menu, place an order, and track its status live — no app install, no login required. Staff manage the menu, tables, and kitchen workflow through authenticated dashboards.

## Features

- **QR-based table ordering** — each table has a unique QR code linking to its own order page
- **Live menu** with categories and products, admin-managed
- **Cart & checkout** with UPI (deep link) or Cash payment options
- **WhatsApp order confirmation** for cash payments (click-to-chat, no API cost)
- **Live order tracking** for customers — Pending → Preparing → Ready → Served, polling-based
- **Kitchen dashboard** — live board of active orders, staff can advance order status
- **Admin dashboard** — add/pause/delete products, manage table count, view order history and revenue
- **Table availability system** — "Order Now" on the landing page auto-assigns the next free table; shows an apology message if all tables are occupied
- **One active order per table** — enforced server-side to prevent duplicate/overlapping orders
- **JWT-based staff authentication** — separate `admin` and `kitchen` roles
- **Rate limiting** on login and order placement to prevent abuse

## Tech Stack

- **Backend:** Go, [chi](https://github.com/go-chi/chi) router
- **Database:** SQLite (via `modernc.org/sqlite`) — Postgres migration in progress
- **Auth:** JWT (`golang-jwt/jwt`) + bcrypt password hashing
- **QR generation:** `skip2/go-qrcode`
- **Rate limiting:** `go-chi/httprate`
- **Frontend:** Plain HTML/CSS/JS (no framework) — server-rendered static pages with fetch-based API calls

## Project Structure

```
qr-ordering-system/
├── cmd/
│   ├── server/          # main application entrypoint
│   └── seedadmin/       # CLI script to create staff (admin/kitchen) accounts
├── internal/
│   ├── admin/           # admin dashboard service + handlers
│   ├── config/          # environment-based configuration loading
│   ├── database/        # DB connection + schema migrations
│   ├── handlers/        # HTTP handlers (product, order, kitchen, auth, table, config)
│   ├── kitchen/         # kitchen dashboard service
│   ├── middleware/      # JWT auth middleware
│   ├── models/          # data models (Order, Product, Category, Table, User, etc.)
│   ├── repository/      # DB access layer, one file per entity
│   ├── routes/          # route wiring
│   ├── services/        # business logic layer
│   └── utils/           # shared response helpers
├── web/
│   ├── static/          # generated QR code images, static assets
│   └── templates/       # HTML pages (home, menu, login, kitchen, admin dashboards)
├── go.mod / go.sum
└── .env                 # local environment config (not committed)
```

## Setup

### 1. Install dependencies
```bash
go mod tidy
```

### 2. Configure environment
Create a `.env` file in the project root (or export these as shell variables):
```env
PORT=8080
DATABASE_URL=postgres://orangemango:binarySearch_3@localhost:5432/qr_ordering?sslmode=disable
BASE_URL=http://localhost:8080
JWT_SECRET=replace-with-a-long-random-string
UPI_ID=yourshop@upi
WHATSAPP_NUMBER=911234567890
```

### 3. Run the server
```bash
go run ./cmd/server
```

On first run, the database file and schema are created automatically.

### 4. Create a staff account
```bash
go run ./cmd/seedadmin <username> <password> admin
go run ./cmd/seedadmin <username> <password> kitchen
```

### 5. Seed a menu
Log in to get a token, then use the admin dashboard (`/admin-dashboard`) or the API directly:
```bash
curl -X POST http://localhost:8080/api/admin/categories \
  -H "Authorization: Bearer <token>" -H "Content-Type: application/json" \
  -d '{"name": "Burgers", "display_order": 1}'
```

## Key URLs

| Path | Purpose |
|---|---|
| `/` | Customer landing page — "Order Now" button |
| `/order/{tableID}` | Menu + cart + checkout for a specific table |
| `/login` | Staff login |
| `/kitchen-dashboard` | Live kitchen order board (admin + kitchen roles) |
| `/admin-dashboard` | Menu management, table settings, revenue & order history (admin only) |

## Core API Endpoints

**Public**
- `GET /api/menu` — categories with nested products
- `GET /api/config` — public config (UPI ID, WhatsApp number)
- `GET /api/tables/next-available` — next free table, or 409 if none
- `POST /api/orders` — place an order
- `GET /api/orders/{orderID}` — order status (customer tracking; the order ID itself is the access token)
- `PATCH /api/orders/{orderID}/confirm-served` — customer confirms receipt

**Admin only**
- `POST /api/admin/products` / `GET /api/admin/products` / `DELETE /api/admin/products/{id}`
- `PATCH /api/admin/products/{id}/toggle-availability` — pause/resume an item
- `POST /api/admin/categories`
- `PATCH /api/admin/settings/tables` — set table count (5–20)
- `GET /api/admin/dashboard` — revenue summary + order history
- `DELETE /api/admin/orders/flush` — wipe all orders (testing/reset only)

**Admin + Kitchen**
- `GET /api/kitchen/orders` — active orders (pending/preparing/ready)
- `PATCH /api/kitchen/orders/{orderID}/status` — advance order status

## Design Notes

## Token

TOKEN=$(curl -s -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "himanshu", "password": "yourpassword123"}' \
  | grep -o '"token":"[^"]*' | cut -d'"' -f4)

- **Order isolation without customer login:** the order ID (a UUID) doubles as an access token. Anyone with the ID can view/confirm that order, but IDs aren't guessable or enumerable.
- **One active order per table:** enforced in `OrderService.PlaceOrder` by checking for any existing order in `pending`/`preparing`/`ready` status before allowing a new one.
- **Table availability:** `GetNextAvailableTable` scans table numbers 1 through the configured total and returns the first without an active order; returns a clear error if none are free.
- **Payment tracking:** orders store `payment_method` (`upi` or `cash`). UPI uses a `upi://pay` deep link (no gateway integration, no automatic payment verification — see Known Limitations). Cash orders trigger a pre-filled WhatsApp message via `wa.me` click-to-chat.

## Known Limitations / Not Yet Implemented

- **UPI payments are unverified** — the current flow is honor-system (customer confirms they paid); there's no webhook or gateway integration confirming funds actually landed. A real payment gateway (e.g. Razorpay) would be needed for verified payments.
- **WhatsApp confirmation is manual** — `wa.me` links require the customer to tap "Send" themselves; there's no server-triggered automatic message (would require WhatsApp Business API).
- **SQLite** — fine for single-shop scale, but only supports one concurrent writer. A migration to PostgreSQL is in progress for better concurrency and to support future multi-location deployments.
- **No order history pagination** — admin dashboard shows the last 100 served orders only.
- **No audit log** for admin actions (product deletions, table count changes, order flushes).
- **`.env` loading** is manual (exported shell variables); `godotenv` auto-loading is not yet wired in.

## Roadmap (Multi-Kitchen / Scale-Up)

For expanding beyond a single shop (discussed for a 15–20 kitchen, ~100 table scale):
- Migrate to PostgreSQL (in progress)
- Add multi-tenancy (`kitchen_id`/`location_id` across schema)
- Replace polling with WebSockets for real-time order push
- Add Redis for pub/sub (multi-instance websocket fan-out) and rate limiting
- Proper role-based JWT auth per location
- Structured logging and basic operational metrics
