# PLAN.md — Membership SaaS

## Overview

A multi-tenant SaaS membership platform where business owners (Admins) manage their members, stores, and billing. Super Admin manages the platform and onboards Admins. Members view their digital membership card, invoices, and payment history.

**Client:** [TBD]
**Timeline:** 4-6 weeks (MVP)
**Stack:** Go + PostgreSQL (backend), React + Vite + TypeScript + Tailwind (frontend), Docker

---

## Roles

| Role | Description |
|------|-------------|
| **Super Admin** | Platform owner. Manages Admin accounts, global settings, billing overview |
| **Admin** | Business owner. Manages members, owns exactly one store (toko), creates invoices |
| **Member** | End user. Views membership card, sees invoices, payment history |

---

## Entities & Relationships

```
Super Admin (-platform level)
  └── Admin (1 user → 1 store)
        └── Store (toko)
              └── Members (many)
                    ├── Invoices (many)
                    └── Payment History (many)
```

- **1 Admin = 1 Store (toko)**. If the client later wants multi-store, we add a `stores` table and many-to-many. For MVP: strict 1:1.
- Members belong to a store. A member can only see their own store's data.

---

## Core Features

### Phase 1 — Foundation (Week 1-2)
- [ ] Auth (Google OAuth + email/password fallback)
- [ ] Role-based access control (RBAC) middleware
- [ ] Admin onboarding (Super Admin creates Admin account)
- [ ] Store creation (auto-created on Admin onboarding)
- [ ] Member management CRUD (Admin creates/imports members)

### Phase 2 — Core Business (Week 2-3)
- [ ] Membership card (digital display with QR code)
- [ ] Invoice CRUD (Admin creates tagihan for member)
- [ ] Payment invoice (Admin records payment against invoice)
- [ ] Payment history (Member view)
- [ ] Invoice status lifecycle: `draft` → `sent` → `paid` → `overdue` → `cancelled`

### Phase 3 — Polish (Week 4)
- [ ] Dashboard summary (Super Admin: total stores, members, revenue; Admin: store stats)
- [ ] Member self-registration (optional, via invitation link from Admin)
- [ ] Export CSV (members list, payment history)
- [ ] Responsive mobile-first UI (members view card on phone)
- [ ] Basic audit log (who created what)

### Phase 4 — Optional Extensions (Week 5+)
- [ ] WhatsApp OTP auth (add beside Google)
- [ ] Payment gateway integration (Midtrans/Xendit) — auto-generate VA
- [ ] Recurring invoice automation (subscription model)
- [ ] Store branding (logo upload, custom card colors)

---

## API Design

### Auth
```
POST /api/auth/register
POST /api/auth/login
GET  /api/auth/google
GET  /api/auth/google/callback
POST /api/auth/logout
GET  /api/auth/me
```

### Admin (Super Admin only)
```
GET    /api/admin/admins
POST   /api/admin/admins              → auto-create store
PUT    /api/admin/admins/:id
DELETE /api/admin/admins/:id
GET    /api/admin/dashboard           → platform-wide stats
```

### Members (Admin scope)
```
GET    /api/members?store_id=:id
POST   /api/members
GET    /api/members/:id
PUT    /api/members/:id
DELETE /api/members/:id
GET    /api/members/:id/card         → membership card data (QR, name, store, tier)
```

### Invoices (Admin creates, Member views own)
```
GET    /api/invoices?store_id=:id&member_id=:id&status=:status
POST   /api/invoices
GET    /api/invoices/:id
PUT    /api/invoices/:id
DELETE /api/invoices/:id
POST   /api/invoices/:id/pay         → record payment
```

### Store
```
GET    /api/store
PUT    /api/store                    → update branding/settings
```

---

## Data Model (Summary)

| Table | Key fields |
|-------|-----------|
| `users` | id, role_id (super_admin/admin/member), email, password_hash, google_id, name, is_active |
| `roles` | id, name (super_admin, admin, member) |
| `stores` | id, admin_id (FK users), name, logo_url, address, card_color_theme |
| `members` | id, store_id (FK stores), user_id (FK users), member_code (unique), tier, joined_at |
| `invoices` | id, store_id, member_id, invoice_number, amount, status, due_date, created_at |
| `payments` | id, invoice_id, amount, method, reference, paid_at |

---

## Tech Decisions

| Decision | Choice | Reason |
|----------|--------|--------|
| Auth | Google OAuth + email/password | Cheaper, stable, easy to maintain |
| DB | PostgreSQL | Standard, GORM support, relational fits membership |
| ORM | GORM | Team reference project uses it, solid migration support |
| Router | chi or gin | Lightweight, idiomatic Go |
| Frontend | React + Vite + TS + Tailwind | Standard stack |
| State | Zustand or React Query | Lightweight |
| QR Code | go-qrcode library | Generate membership card QR |
| Deployment | Docker + VPS | Standard, already set up |

---

## MVP Scope (What's OUT)

- ❌ Multi-store per Admin (add later)
- ❌ Payment gateway auto-integration (manual payment records first)
- ❌ Mobile app (responsive web is enough)
- ❌ Real-time notifications (poll or add later)
- ❌ Advanced reporting (basic dashboard only)
- ❌ Member self-service password reset (Admin resets manually)

---

## Acceptance Criteria

1. Super Admin can create Admin → auto-creates store
2. Admin can login → see only their store's members
3. Admin can create member → member can login
4. Member sees digital card with their info + QR
5. Admin creates invoice → member sees it in their dashboard
6. Admin records payment → invoice status changes to paid
7. All RBAC enforced at API level (cannot access other store's data)

---

## Next Steps

1. [ ] Confirm client name and domain
2. [ ] Confirm branding requirements (logo, colors, card design)
3. [ ] Set up GitHub repo
4. [ ] Create project skeleton (Go backend + React frontend)
5. [ ] Begin Phase 1 implementation
