# 🔗 Linklytics: AI-Powered URL Shortener

An enterprise-grade, high-performance URL shortener built with **Go (Gin)** and **React**. 

Unlike standard URL shorteners, Linklytics features an asynchronous AI worker that automatically generates summaries of the content you are linking to, robust geographic click analytics, and distributed rate limiting to prevent abuse.

## ✨ Features

- **Blazing Fast Redirects**: Powered by Go and Redis caching for sub-10ms latency.
- **Advanced Analytics**: Tracks total clicks, unique geographic locations (via MaxMind GeoIP), and device types.
- **AI Link Summarization**: Uses the Groq LLM API to asynchronously summarize the destination webpage in the background without blocking the user.
- **Distributed Rate Limiting**: Implements a Token Bucket algorithm in Redis to protect public endpoints from DDoS and spam.
- **Secure Authentication**: JWT-based user authentication and bcrypt password hashing.
- **Premium UI**: Custom-built, responsive frontend using Vite, React, Recharts, and a Glassmorphism design system.

## 🛠️ Tech Stack

- **Backend**: Go, Gin, GORM
- **Frontend**: React (Vite), React Router, Recharts, Lucide Icons
- **Database**: PostgreSQL
- **Cache/Rate Limiting**: Redis
- **AI Integration**: Groq (Llama 3 70B)
- **Infrastructure**: Docker, Docker Compose, GitHub Actions (CI/CD)

## 🚀 Quick Start (Local Development)

The entire full-stack application is containerized. You can spin up the Go backend, React frontend, PostgreSQL database, and Redis cache with a single command.

### Prerequisites
- Docker and Docker Compose installed.
- A free API key from [Groq](https://console.groq.com/).

### Installation

1. **Clone the repository:**
   ```bash
   git clone https://github.com/developerAS123/urlShortner.git
   cd urlShortner
   ```

2. **Configure Environment Variables:**
   Create a `.env` file in the `backend/` directory:
   ```bash
   PORT=8080
   DATABASE_URL=postgres://postgres:password@db:5432/urlshortener?sslmode=disable
   REDIS_URL=redis://redis:6379/0
   JWT_SECRET=your_super_secret_key_here
   GROQ_API_KEY=your_groq_api_key_here
   ```

3. **Run the stack:**
   ```bash
   docker compose up --build
   ```

4. **Access the Application:**
   - **Frontend**: [http://localhost:3000](http://localhost:3000)
   - **Backend API**: [http://localhost:8080/api](http://localhost:8080/api)

## 📖 Architecture

Curious about how the background workers, rate limiters, and telemetry pipelines are designed? 
Check out the [Architecture Documentation](./ARCHITECTURE.md) for a deep dive into the system design.

## 🤝 Contributing
Contributions, issues, and feature requests are welcome! Feel free to check the issues page.
