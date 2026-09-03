# Members — Membership Platform

A modern membership platform for managing members, subscriptions, and access control.

## Quick Start

```bash
# Clone
git clone https://github.com/lovelymondayz/members.git
cd members

# Start all services
docker compose up -d --build

# Frontend: http://localhost:3003
# Backend API: http://localhost:8082
# DB: localhost:5435 (user: members, pass: members_secret)
```

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                        NGINX (80/443)                        │
│                   members.arjism.com → :3003                │
├─────────────────────────────────────────────────────────────┤
│  React + Vite + TS + Tailwind  │  Go + GIN + pgx + Postgres │
│        (Frontend :3003)        │       (Backend :8082)      │
├─────────────────────────────────────────────────────────────┤
│              PostgreSQL :5435  │  Local Storage (/data)     │
└─────────────────────────────────────────────────────────────┘
```

## Features

- **Member Management**: Create, update, and manage member profiles
- **Subscription Tracking**: Monitor subscription status and billing
- **Access Control**: Role-based permissions and access management
- **Google OAuth**: Sign in with Google integration
- **JWT Authentication**: Secure token-based authentication
- **Responsive UI**: Mobile-first design with Tailwind CSS

## API Endpoints

### Public
- `GET /api/health` — Health check
- `POST /api/auth/google` — Google OAuth login
- `GET /api/auth/google/callback` — OAuth callback

### Authenticated (JWT required)
- `GET /api/members` — List members
- `POST /api/members` — Create member
- `GET /api/members/:id` — Get member details
- `PUT /api/members/:id` — Update member
- `DELETE /api/members/:id` — Delete member

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| APP_ENV | production | Environment mode |
| PORT | 8082 | Backend port |
| DB_HOST | db | Database host |
| DB_PORT | 5432 | Database port |
| DB_USER | members | Database user |
| DB_PASSWORD | members_secret | Database password |
| DB_NAME | members | Database name |
| JWT_SECRET | - | JWT signing key |
| GOOGLE_CLIENT_ID | - | Google OAuth client ID |
| GOOGLE_CLIENT_SECRET | - | Google OAuth client secret |
| GOOGLE_REDIRECT_URL | - | OAuth redirect URL |

## Development

```bash
# Backend only
cd backend
go run .

# Frontend only
cd frontend
npm install
npm run dev
```

## Deployment

1. Push to `main` → GitHub Action auto-deploys
2. Or manually: `ssh vps && cd /root/members && ./update.sh`

## License

MIT
