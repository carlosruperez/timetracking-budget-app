# Deployment Guide

This document covers the full deployment lifecycle for the timetracking-budget-app monorepo:

- **Backend (Go)** → Fly.io
- **Web (Next.js)** → Vercel
- **Mobile (Expo)** → EAS (Expo Application Services)
- **Database (PostgreSQL)** → Supabase
- **CI/CD** → GitHub Actions

---

## Table of Contents

1. [Prerequisites](#1-prerequisites)
2. [GitHub Repository Setup](#2-github-repository-setup)
3. [Database Setup — Supabase](#3-database-setup--supabase)
4. [Backend Deployment — Fly.io](#4-backend-deployment--flyio)
5. [Web Deployment — Vercel](#5-web-deployment--vercel)
6. [Mobile Build — EAS](#6-mobile-build--eas)
7. [GitHub Actions Secrets](#7-github-actions-secrets)
8. [Local Development Quick-Start](#8-local-development-quick-start)

---

## 1. Prerequisites

Install the following CLIs before proceeding.

### Required tools

| Tool | Install | Version |
|------|---------|---------|
| Go | https://go.dev/dl/ | 1.22+ |
| Node.js | https://nodejs.org/ | 20 LTS |
| pnpm | `npm install -g pnpm@9` | 9.x |
| Docker | https://docs.docker.com/get-docker/ | 24+ |
| flyctl | `curl -L https://fly.io/install.sh \| sh` | latest |
| Vercel CLI | `pnpm add -g vercel` | latest |
| EAS CLI | `pnpm add -g eas-cli` | 10+ |

### Accounts needed

- [GitHub](https://github.com) — source control and CI/CD
- [Supabase](https://supabase.com) — managed PostgreSQL
- [Fly.io](https://fly.io) — Go backend hosting
- [Vercel](https://vercel.com) — Next.js hosting
- [Expo](https://expo.dev) — mobile build service (EAS)
- [Apple Developer Program](https://developer.apple.com) — iOS builds ($99/year)
- [Google Play Console](https://play.google.com/console) — Android builds

---

## 2. GitHub Repository Setup

### Create and push the repository

```bash
# From the project root
git init
git add .
git commit -m "chore: initial commit"

# Create repository on GitHub (using the GitHub CLI)
gh repo create timetracking-budget-app --private --source=. --push

# Or push to an existing repository
git remote add origin https://github.com/YOUR_USERNAME/timetracking-budget-app.git
git branch -M main
git push -u origin main
```

### Branch strategy

The CI/CD pipelines are configured for:

- `main` → production deployments trigger automatically on push
- Pull requests → run lint, type-check, and build checks only (no deploy)

---

## 3. Database Setup — Supabase

### Create a Supabase project

1. Go to https://supabase.com and sign in.
2. Click **New project**.
3. Fill in:
   - **Name**: `timetracking`
   - **Database Password**: generate a strong password and store it securely
   - **Region**: pick the closest to your Fly.io region (e.g., `eu-west-1` for Madrid)
4. Wait for the project to finish provisioning (~2 minutes).

### Get the connection string

1. In your Supabase project, go to **Settings → Database**.
2. Under **Connection string**, select **URI** mode.
3. Copy the string — it looks like:

```
postgres://postgres:[YOUR-PASSWORD]@db.xxxxxxxxxxxx.supabase.co:5432/postgres
```

4. Append `?sslmode=require` to the end:

```
postgres://postgres:[YOUR-PASSWORD]@db.xxxxxxxxxxxx.supabase.co:5432/postgres?sslmode=require
```

Keep this string — you will need it for Fly.io secrets and GitHub Actions secrets.

### Run migrations

The Go backend runs migrations automatically on startup. Once the backend is deployed and `DATABASE_URL` is set, migrations will run on the first boot.

To run migrations manually against Supabase during development:

```bash
export DATABASE_URL="postgres://postgres:[PASSWORD]@db.xxxxxxxxxxxx.supabase.co:5432/postgres?sslmode=require"
cd backend
go run ./cmd/server/main.go
# The server will migrate and then start listening
```

---

## 4. Backend Deployment — Fly.io

### Authenticate

```bash
flyctl auth login
```

### Launch the app (first time only)

Run this from the `backend/` directory. When prompted, do **not** let Fly overwrite the existing `fly.toml` — answer **No**.

```bash
cd backend
flyctl launch --no-deploy
```

The `fly.toml` already exists and is pre-configured with:
- App name: `timetracking-api`
- Region: `mad` (Madrid) — change `primary_region` in `fly.toml` if needed
- 256 MB RAM, 1 shared CPU
- Health check at `GET /health`
- Auto-stop when idle (scale to zero)

### Set secrets

Fly.io secrets are environment variables that are encrypted at rest and injected at runtime. Never commit these to git.

```bash
cd backend

# Set the production database URL from Supabase
flyctl secrets set DATABASE_URL="postgres://postgres:[PASSWORD]@db.xxxxxxxxxxxx.supabase.co:5432/postgres?sslmode=require"

# Generate and set a strong JWT secret
flyctl secrets set JWT_SECRET="$(openssl rand -hex 32)"
```

Verify secrets are set (values are hidden):

```bash
flyctl secrets list
```

### Deploy manually

```bash
cd backend
flyctl deploy
```

Fly.io will build the Docker image remotely using `backend/Dockerfile` and deploy it.

### Verify the deployment

```bash
# Check deployment status
flyctl status

# Tail live logs
flyctl logs

# Check the health endpoint
curl https://timetracking-api.fly.dev/health
```

### Get the Fly API token for CI

```bash
flyctl tokens create deploy -x 999999h
```

Copy the printed token — you will add it as `FLY_API_TOKEN` in GitHub Secrets (see Section 7).

### Changing the region

Edit `primary_region` in `backend/fly.toml` to one of the Fly.io region codes:

| Code | Location |
|------|----------|
| `mad` | Madrid, Spain |
| `lhr` | London, UK |
| `cdg` | Paris, France |
| `iad` | Ashburn, VA (US) |
| `sin` | Singapore |
| `syd` | Sydney, Australia |

Full list: https://fly.io/docs/reference/regions/

---

## 5. Web Deployment — Vercel

### Authenticate

```bash
vercel login
```

### Link the project (first time only)

```bash
cd apps/web
vercel link
```

Follow the prompts:
- Select your Vercel account/team
- Create a new project named `timetracking-web` (or any name you prefer)

This creates `.vercel/project.json` locally with `projectId` and `orgId`. You will need these for CI.

### Set environment variables on Vercel

Go to your Vercel project dashboard → **Settings → Environment Variables** and add:

| Variable | Value | Environments |
|----------|-------|--------------|
| `NEXT_PUBLIC_API_URL` | `https://timetracking-api.fly.dev` | Production, Preview |
| `BACKEND_URL` | `https://timetracking-api.fly.dev` | Production, Preview |

`NEXT_PUBLIC_API_URL` is embedded in the client bundle at build time.
`BACKEND_URL` is used server-side (API routes / rewrites).

### Deploy manually

```bash
cd apps/web
vercel --prod
```

### Get IDs for CI

```bash
# After running vercel link, read the generated file
cat apps/web/.vercel/project.json
```

You will see:
```json
{
  "orgId": "team_xxxxxxxxxxxx",
  "projectId": "prj_xxxxxxxxxxxx"
}
```

You will also need a Vercel API token:
1. Go to https://vercel.com/account/tokens
2. Create a new token with **Full Account** scope
3. Copy the token — add it as `VERCEL_TOKEN` in GitHub Secrets

---

## 6. Mobile Build — EAS

### Authenticate

```bash
eas login
```

### Configure EAS (first time only)

```bash
cd apps/mobile
eas build:configure
```

This updates `app.json` with your EAS project ID and creates `eas.json` if it doesn't exist. Since `eas.json` already exists in this repo, only `app.json` needs updating.

### Build profiles

The `eas.json` defines three profiles:

| Profile | Purpose | Distribution |
|---------|---------|--------------|
| `development` | Local dev client with hot reload | Internal (simulator on iOS) |
| `preview` | Internal testing build | Internal (install via QR code) |
| `production` | App Store / Play Store release | Store |

### Trigger a build manually

```bash
cd apps/mobile

# Build for both platforms (preview profile)
eas build --platform all --profile preview

# Build iOS only (development — runs in simulator)
eas build --platform ios --profile development

# Build production release
eas build --platform all --profile production
```

Builds run on Expo's cloud infrastructure. You will receive a link to download the artifact when the build completes.

### iOS setup

For iOS builds you need:
1. An Apple Developer account enrolled in the Apple Developer Program
2. Run `eas credentials` to let EAS manage your certificates and provisioning profiles automatically

Update `eas.json → submit.production.ios` with your Apple credentials:

```json
"ios": {
  "appleId": "you@example.com",
  "ascAppId": "1234567890",
  "appleTeamId": "ABCDE12345"
}
```

- `ascAppId`: found in App Store Connect under your app's **App Information**
- `appleTeamId`: found at https://developer.apple.com/account under **Membership**

### Android setup

For Android production builds you need a Google Play service account:
1. Go to Google Play Console → **Setup → API access**
2. Create a service account with **Release Manager** permissions
3. Download the JSON key file
4. Place it at `apps/mobile/google-service-account.json` (this file is gitignored)

### Submit to stores

After a successful production build:

```bash
# Submit to both stores
eas submit --platform all --profile production

# Submit to App Store only
eas submit --platform ios --profile production

# Submit to Play Store only
eas submit --platform android --profile production
```

### EAS token for CI

1. Go to https://expo.dev/accounts/[your-account]/settings/access-tokens
2. Create a new token
3. Copy it — add it as `EXPO_TOKEN` in GitHub Secrets

---

## 7. GitHub Actions Secrets

Go to your GitHub repository → **Settings → Secrets and variables → Actions → New repository secret**.

Add all of the following secrets:

### Backend (Fly.io)

| Secret | How to get it |
|--------|---------------|
| `FLY_API_TOKEN` | `flyctl tokens create deploy -x 999999h` |

### Web (Vercel)

| Secret | How to get it |
|--------|---------------|
| `VERCEL_TOKEN` | Vercel dashboard → Account Settings → Tokens |
| `VERCEL_ORG_ID` | `cat apps/web/.vercel/project.json` → `orgId` |
| `VERCEL_PROJECT_ID` | `cat apps/web/.vercel/project.json` → `projectId` |
| `NEXT_PUBLIC_API_URL` | `https://timetracking-api.fly.dev` |
| `BACKEND_URL` | `https://timetracking-api.fly.dev` |

### Mobile (EAS)

| Secret | How to get it |
|--------|---------------|
| `EXPO_TOKEN` | expo.dev → Account Settings → Access Tokens |

### Summary table

| Secret | Used by workflow |
|--------|-----------------|
| `FLY_API_TOKEN` | `backend.yml` |
| `VERCEL_TOKEN` | `web.yml` |
| `VERCEL_ORG_ID` | `web.yml` |
| `VERCEL_PROJECT_ID` | `web.yml` |
| `NEXT_PUBLIC_API_URL` | `web.yml` |
| `BACKEND_URL` | `web.yml` |
| `EXPO_TOKEN` | `mobile.yml` |

---

## 8. Local Development Quick-Start

### Option A — Docker Compose (recommended for backend parity)

Runs PostgreSQL, the Go backend, and the Next.js web app in containers.

```bash
# From the project root
cp .env.example .env
# Edit .env if needed (defaults work out of the box for local dev)

docker compose up --build
```

Services:
- PostgreSQL: `localhost:5432`
- Backend API: `http://localhost:8080`
- Web app: `http://localhost:3000`

Stop everything:

```bash
docker compose down
```

Wipe the database volume:

```bash
docker compose down -v
```

### Option B — Run services natively

**1. Start PostgreSQL**

```bash
docker run --rm \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=timetracking \
  -p 5432:5432 \
  postgres:16-alpine
```

**2. Start the Go backend**

```bash
cd backend
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/timetracking?sslmode=disable"
export JWT_SECRET="dev-secret-change-in-production"
export PORT=8080
export ENV=development
go run ./cmd/server/main.go
```

**3. Start the Next.js web app**

```bash
# From the monorepo root
pnpm install

# Build the shared core package first
pnpm --filter @timetracking/core build

# Start the web dev server
pnpm --filter @timetracking/web dev
```

Web app available at http://localhost:3000.

**4. Start the Expo mobile app**

```bash
cd apps/mobile
pnpm install

# Start the Metro bundler
npx expo start

# Press 'i' for iOS simulator, 'a' for Android emulator
```

### Environment variables for local dev

Copy `.env.example` to `.env` at the project root:

```bash
cp .env.example .env
```

The defaults in `.env.example` work for local development without modification. For production values, see the relevant sections above — never commit production secrets.

### Useful commands

```bash
# Run all backend tests
cd backend && go test ./... -v

# Type-check the web app
pnpm --filter @timetracking/web exec tsc --noEmit

# Build everything
pnpm build

# Lint everything
pnpm lint

# Check backend health
curl http://localhost:8080/health
```
