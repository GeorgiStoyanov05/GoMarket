# GoMarket 📈

A full-stack stock market web app for tracking prices, viewing stock details (charts/news/fundamentals), and managing watchlists & alerts.

> **Status:** Work in progress.

---

## ✨ Features
- Public stock browsing (no account required)
- Stock details page (chart + news + fundamentals)
- JWT authentication (register/login)
- Watchlists (save and track symbols)
- Price alerts (get notified when conditions hit)
- Real-time updates via WebSockets (planned)

---

## 🛠 Tech Stack
**Frontend**
- Next.js (React + TypeScript)

**Backend**
- Go
- Gin/Gonic (REST API)
- JWT authentication
- WebSockets (gorilla/websocket)

**Data / Integrations**
- MongoDB
- Finnhub API (market data)
- TradingView widgets (UI embedding)

---

## 📁 Project Structure
```txt
GoMarket/
  backend/     # Go API + WebSocket server
  frontend/    # Next.js client
```

---

## 🚀 Getting Started

### Prerequisites
- Go
- Node.js + npm
- MongoDB (local or Atlas)
- Finnhub API key

### Setup

1) Clone the repo:
```bash
git clone https://github.com/GeorgiStoyanov05/GoMarket.git
cd GoMarket
```

2) Create backend env file: `backend/.env`
```env
PORT=8080
MONGODB_URI=your_mongo_uri
FINNHUB_API_KEY=your_finnhub_key
JWT_SECRET=your_jwt_secret
```

3) Start the backend:
```bash
cd backend
go mod download
go run .
```

4) Start the frontend:
```bash
cd ../frontend
npm install
npm run dev
```

### Default URLs
- Frontend: `http://localhost:3000`
- Backend: `http://localhost:8080`

---

## 📌 Roadmap
- [ ] Finish auth + protected routes
- [ ] Complete watchlists + alerts
- [ ] Real-time prices (WS) on dashboard/watchlists
- [ ] Improve UI/UX + add tests

---

## 📄 License
MIT License — see `LICENSE`.
