# ARCHITECTURE.md — Membership SaaS

## System Overview

```
┌─────────────────────────────────────────────────────┐
│                      VPS (2c/12GB)                   │
│  ┌────────────┐  ┌────────────┐  ┌──────────────┐  │
│  │  Frontend   │  │  Backend   │  │  PostgreSQL  │  │
│  │  React+Vite │  │  Go (Gin)  │  │   (Docker)   │  │
│  │  Port 3001  │  │  Port 8081 │  │  Port 5432   │  │
│  └──────┬─────┘  └──────┬─────┘  └──────┬───────┘  │
│         │               │                │          │
│         └───────────────┼────────────────┘          │
│                         │                           │
│              ┌──────────┴──────────┐                │
│              │       Nginx         │                │
│              │   Reverse Proxy     │                │
│              └──────────┬──────────┘                │
└─────────────────────────┼───────────────────────────┘
                          │
                    ┌─────┴─────┐
                    │  Client   │
                    │ member.*  │
                    └───────────┘
```

---

## Backend (Go)

### Structure

```
membership-saas/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── config/          # env vars, app config
│   ├── database/        # GORM connection, migrations
│   ├── middleware/      # auth, RBAC, CORS
│   ├── model/           # GORM models
│   ├── handler/         # HTTP handlers (controllers)
│   ├── service/         # Business logic
│   ├── repository/      # DB queries
│   └── dto/             # Request/Response structs
├── migrations/          # SQL migration files (golang-migrate)
├── Dockerfile
├── docker-compose.yml
└── Makefile
```

### Dependencies

| Package | Purpose |
|---------|---------|
| `gin-gonic/gin` | HTTP router |
| `gorm.io/gorm` + `gorm.io/driver/postgres` | ORM |
| `golang-jwt/jwt/v5` | JWT tokens |
| `golang.org/x/oauth2` | Google OAuth |
| `skasquare/go-qrcode` | QR code generation |
| `joho/godotenv` | Environment loading |
| `golang-migrate/migrate` | DB migrations |

### Auth Flow

```
1. User clicks "Login with Google"
   → Frontend redirects to /api/auth/google
   → Backend redirects to Google OAuth consent screen
   → User authenticates with Google
   → Google redirects to /api/auth/google/callback
   → Backend: verify state, exchange code for token
   → Fetch user info (email, name, picture)
   → If new user → create with role=member (default)
   → Issue JWT (access_token + refresh_token)
   → Redirect frontend with token

2. Subsequent requests:
   → Frontend includes Authorization: Bearer <token>
   → Middleware: validate JWT, extract user_id + role
   → RBAC middleware: check role against allowed roles for route
```

### RBAC

| Route Group | Allowed Roles |
|-------------|---------------|
| `/api/admin/*` | super_admin |
| `/api/members/*` | admin, super_admin |
| `/api/invoices/*` | admin, super_admin |
| `/api/invoices/:id/view` | member (own invoices only) |
| `/api/store/*` | admin, super_admin |
| `/api/auth/me` | any authenticated |

### Data Isolation

- **Admin queries** always include `WHERE store_id = ?` from JWT claims
- **Member queries** always scope to `WHERE member_id = ?` from JWT claims
- Super Admin bypasses store_id filter (sees all)

---

## Frontend (React + Vite + TypeScript + Tailwind)

### Structure

```
frontend/
├── public/
├── src/
│   ├── api/             # Axios instance + API functions
│   ├── components/      # Shared components (Card, Table, Modal, etc.)
│   ├── pages/           # Route-level page components
│   ├── hooks/           # Custom hooks (useAuth, useMembers, etc.)
│   ├── store/           # Zustand state management
│   ├── types/           # TypeScript interfaces
│   ├── utils/           # Helpers (formatCurrency, formatDate)
│   ├── layouts/         # MainLayout, AuthLayout
│   └── App.tsx          # Router definition
├── tailwind.config.js
├── postcss.config.js
├── tsconfig.json
├── vite.config.ts
└── Dockerfile
```

### Routing

| Path | Page | Access |
|------|------|--------|
| `/login` | Login | Public |
| `/dashboard` | Dashboard | Admin, Super Admin |
| `/members` | Member list | Admin, Super Admin |
| `/members/:id` | Member detail + card | Admin, Super Admin |
| `/invoices` | Invoice list | Admin, Super Admin |
| `/invoices/create` | Create invoice | Admin, Super Admin |
| `/card` | My membership card | Member |
| `/my-invoices` | My invoices | Member |
| `/payments` | Payment history | Member |
| `/store/settings` | Store settings | Admin |
| `/admin/admins` | Admin management | Super Admin |
| `/admin/dashboard` | Platform dashboard | Super Admin |

### State Management

- **Auth:** Zustand store with user info, token (in httpOnly cookie ideally, or localStorage with caution)
- **Data fetching:** React Query (TanStack Query) for server state caching & refetch
- **UI state:** Local component state or Zustand

### Key Components

1. **MembershipCard** — Displays digital card with member name, store name, QR code, tier badge, join date
2. **InvoiceTable** — Sortable/filterable table with status badges
3. **DashboardStats** — Cards showing key metrics (total members, revenue, pending invoices)

---

## Database Schema

```sql
-- Roles
CREATE TABLE roles (
    role_id   SERIAL PRIMARY KEY,
    name      VARCHAR(50) UNIQUE NOT NULL,  -- 'super_admin', 'admin', 'member'
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Users
CREATE TABLE users (
    user_id      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    role_id      INT REFERENCES roles(role_id),
    email        VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255),  -- nullable for Google-only users
    google_id    VARCHAR(255) UNIQUE,
    name         VARCHAR(255) NOT NULL,
    avatar_url   TEXT,
    is_active    BOOLEAN DEFAULT true,
    created_at   TIMESTAMPTZ DEFAULT NOW(),
    updated_at   TIMESTAMPTZ DEFAULT NOW(),
    deleted_at   TIMESTAMPTZ
);

-- Stores (1:1 with Admin)
CREATE TABLE stores (
    store_id       UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    admin_id       UUID UNIQUE REFERENCES users(user_id),
    name           VARCHAR(255) NOT NULL,
    logo_url       TEXT,
    address        TEXT,
    phone         VARCHAR(50),
    card_color_hex VARCHAR(7) DEFAULT '#1E40AF',
    created_at     TIMESTAMPTZ DEFAULT NOW(),
    updated_at     TIMESTAMPTZ DEFAULT NOW()
);

-- Members
CREATE TABLE members (
    member_id    UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    store_id     UUID NOT NULL REFERENCES stores(store_id),
    user_id      UUID UNIQUE REFERENCES users(user_id,
    member_code  VARCHAR(50) UNIQUE NOT NULL,
    tier         VARCHAR(50) DEFAULT 'standard',  -- standard, premium, vip
    joined_at    DATE DEFAULT CURRENT_DATE,
    created_at   TIMESTAMPTZ DEFAULT NOW(),
    updated_at   TIMESTAMPTZ DEFAULT NOW(),
    deleted_at   TIMESTAMPTZ
);

-- Invoices
CREATE TYPE invoice_status AS ENUM ('draft', 'sent', 'paid', 'overdue', 'cancelled');

CREATE TABLE invoices (
    invoice_id     UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    store_id       UUID NOT NULL REFERENCES stores(store_id),
    member_id      UUID NOT NULL REFERENCES members(member_id),
    invoice_number VARCHAR(50) UNIQUE NOT NULL,
    amount         DECIMAL(12,2) NOT NULL,
    description    TEXT,
    status         invoice_status DEFAULT 'draft',
    due_date       DATE,
    created_at     TIMESTAMPTZ DEFAULT NOW(),
    updated_at     TIMESTAMPTZ DEFAULT NOW()
);

-- Payments
CREATE TABLE payments (
    payment_id   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    invoice_id   UUID NOT NULL REFERENCES invoices(invoice_id),
    amount       DECIMAL(12,2) NOT NULL,
    method       VARCHAR(50) DEFAULT 'manual',  -- manual, bank_transfer, etc.
    reference    VARCHAR(255),
    note         TEXT,
    paid_at       TIMESTAMPTZ NOT NULL,
    created_at   TIMESTAMPTZ DEFAULT NOW()
);
```

---

## Environment Variables

```env
# Backend
APP_ENV=production
PORT=8081

DB_HOST=localhost
DB_PORT=5432
DB_USER=membership
DB_PASSWORD=membership_secret
DB_NAME=membership

JWT_SECRET=<generate-random>
JWT_EXPIRY_HOURS=24

GOOGLE_CLIENT_ID=<from-google-console>
GOOGLE_CLIENT_SECRET=<from-google-console>
GOOGLE_REDIRECT_URL=https://client.arjism.com/api/auth/google/callback

# Frontend
VITE_API_URL=https://client.arjism.com/api
```

---

## Deployment

1. Build Go binary: `go build -o main ./cmd/server`
2. Build React: `cd frontend && npm run build` → static files copied to nginx
3. `docker compose up -d` (backend + postgres)
4. Nginx serves frontend static files, proxies `/api/*` to backend

### Nginx Config (simplified)

```nginx
server {
    listen 443 ssl;
    server_name client.arjism.com;

    # Frontend static files
    root /var/www/membership-saas/frontend/dist;
    index index.html;

    location / {
        try_files $uri $uri/ /index.html;
    }

    # API proxy
    location /api/ {
        proxy_pass http://localhost:8081;
        proxy_set_header Host $host;
        proxy_set_header X-Forwarded-For $remote_addr;
    }
}
```

---

## Security Considerations

| Concern | Mitigation |
|---------|-----------|
| JWT theft | Short expiry (24h), consider httpOnly cookies |
| SQL injection | GORM parameterized queries (no raw SQL) |
| RBAC bypass | Every handler validates `store_id` against JWT claims |
| Data leakage | CORS restricted to known origin |
| Password storage | bcrypt hashing |
| Google OAuth state | Random state parameter + verification |
| Rate limiting | Add later for auth endpoints if needed |

---

## Trade-offs & Decisions

| Decision | Choice over | Why |
|----------|-----------|-----|
| Google Auth over WA Auth | WhatsApp API | Cheaper, stable, no verification hassle |
| JWT over session | Server-side sessions | Stateless, simpler for single backend |
| Go over Node.js | Express/NestJS | Smaller binary, faster on 2c/12GB VPS |
| Manual payment recording | Midtrans auto-integration | MVP scope — client can record payments manually first |
| Nginx static file serving | Serve from Go | Nginx is faster at static files, frees Go for API |
| Single DB | Separate DBs per tenant | Simpler for MVP, migration path to schema-per-tenant later |

---

## Future Scaling Path

If the client grows:
1. **Multi-store** → add `stores` table with many-to-many to admins
2. **Payment gateway** → add Midtrans/Xendit webhook handlers
3. **Subscription automation** → cron job to auto-generate recurring invoices
4. **WA auth** → add as second OAuth provider (same user_id, linked accounts)
5. **Read replica** → if DB becomes bottleneck, add PG read replica for reports
6. **CDN** → move static assets to Cloudflare R2
