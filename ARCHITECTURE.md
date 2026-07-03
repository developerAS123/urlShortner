# System Architecture

Linklytics is designed as a decoupled, microservices-oriented application. The architecture prioritizes **low latency** for core operations (redirecting users) while offloading heavy computations (AI summarization, telemetry aggregation) to asynchronous background workers.

## High-Level Architecture Diagram

```mermaid
graph TD
    Client[Web Browser / Client]
    
    subgraph Frontend [React Frontend]
        UI[Dashboard & Analytics UI]
    end

    subgraph Backend [Go (Gin) Backend]
        API[REST API Handlers]
        RateLimiter[Token Bucket Rate Limiter]
        RedirectEngine[Redirect Engine]
        
        subgraph Workers [Asynchronous Workers]
            Telemetry[GeoIP Telemetry Logger]
            AIWorker[Groq AI Summarization Job]
        end
    end

    subgraph Data Layer
        Postgres[(PostgreSQL)]
        Redis[(Redis Cache)]
        MaxMind[(MaxMind GeoIP DB)]
    end
    
    subgraph External APIs
        Groq[Groq LLM API]
    end

    Client <--> UI
    UI <--> |HTTP/JSON| RateLimiter
    Client --> |GET /:slug| RateLimiter
    
    RateLimiter <--> |Check Quota| Redis
    RateLimiter --> API
    RateLimiter --> RedirectEngine
    
    API <--> Postgres
    RedirectEngine <--> |Lookup ShortURL| Postgres
    
    RedirectEngine -.-> |Fire & Forget| Telemetry
    Telemetry <--> MaxMind
    Telemetry --> |Insert Click Event| Postgres
    
    API -.-> |Trigger Job| AIWorker
    AIWorker <--> Groq
    AIWorker --> |Cache Summary| Postgres
```

## Key Technical Decisions

### 1. Token Bucket Rate Limiting (Redis)
To protect the application from brute-force attacks and DDoS, we implemented a custom Rate Limiter using Redis. 
- **Why Redis?** Checking rate limits requires an operation on *every single incoming request*. Querying a SQL database would introduce unacceptable latency. Redis operates entirely in memory, allowing O(1) time complexity checks that add <1ms to the request lifecycle.

### 2. Asynchronous AI Workers
Generating a summary for a target URL using an LLM can take anywhere from 1 to 5 seconds.
- **The Problem**: If this was done synchronously on the main HTTP thread, the user would be left staring at a loading spinner when creating a link.
- **The Solution**: The API immediately returns a `201 Created` response. In the background, a Go goroutine fetches the webpage content, calls the Groq API (with exponential backoff retries for `429 Too Many Requests`), and persists the AI summary to PostgreSQL. 

### 3. Non-Blocking Telemetry (GeoIP)
When a user clicks a short link, the system must parse their User-Agent, look up their IP address in the MaxMind GeoLite2 database to find their country, and log a `click_event`.
- **Optimization**: The redirect engine (`GET /:slug`) immediately responds with a `301 Moved Permanently`. The telemetry logging is dispatched asynchronously via Go channels, ensuring that tracking analytics has **0% impact** on the user's redirect speed.

## Database Schema

- **`users`**: Manages authentication (email, bcrypt password_hash).
- **`links`**: Core domain model (slug, original_url, user_id).
- **`click_events`**: Append-only telemetry log (link_id, ip_address, country, device_type, clicked_at).
- **`ai_summaries`**: Caches the generated insights (link_id, summary text, week_start).
