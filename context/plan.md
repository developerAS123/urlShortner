# AI-Powered URL Shortener - Project Plan

## Week 1 — Foundations & Core Shortening (Current Focus)
- Set up Go project (Gin)
- Set up Postgres schema (users, links, click_events, ai_summaries)
- docker-compose.yml for local development (Go app, Postgres, Redis)
- Implement JWT auth (register/login)
- POST `/api/shorten` (slug generation, store in Postgres, cache in Redis TTL 24h)
- GET `/{slug}` (Redis-first lookup, redirect, Postgres fallback)
- Async click logging (goroutine inserting click event)
- CORS middleware
- Deploy Go API to Render
- Connect to Neon/Supabase Postgres and Upstash Redis

## Week 2 — Analytics Data Pipeline
- Expand click logging (User-Agent parsing, maxmind GeoLite2 integration)
- GET `/api/links/{slug}/analytics`
- GET `/api/links`
- Rate limiting middleware

## Week 3 — AI Summary Worker
- Background goroutine scheduler (nightly)
- Groq API integration for LLM summaries
- Cache summaries in DB
- GET `/api/links/{slug}/summary`

## Week 4 — React Frontend
- Vite + React + Tailwind CSS
- Auth pages (login/register)
- Dashboard page (list links)
- Analytics page (Recharts + AI Summary)
- Deploy frontend to Vercel

## Week 5 — Polish, Demo & CI/CD
- Multi-stage Dockerfile
- GitHub Actions CI pipeline
- README
- Demo account setup (seed script)
