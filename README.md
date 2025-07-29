# Full Stack Chat App

A simple chat application I built to improve my skills with Gorilla WebSocket.
Technologies used: GoFiber, Gorilla WebSocket, GORM, PostgreSQL, and SvelteK

### 📸 Gallery
![Chat Example](./galery/chat.png)

### 🔧 Features
- Direct Messaging (DM)
- Group Creation (to be added on the frontend)
- User Registration & Login

### 🛠 Upcoming Features
- Email Authentication
- Update Password / Email / Username
- Meeting Feature
- Responsive Design
- CDN (File Upload Support)
- Emoji / GIF Picker

### 🚀 How to Run
**1. Backend Configuration***
Update your backend/configs/service.ini file with your own settings:
```ini
[db]
dsn="postgres://username:password@localhost:5432/database?sslmode=disable"
redisAddr="localhost:6379"

[server]
port="127.0.0.1:4000"

[smtp]
host="smtp.gmail.com"
port="587"
password="your-email-password"
from=""
```
> **Note**: Replace the placeholders (username, password, database, your-email-password, etc.) with your actual credentials.,

**2.Start the Backend**
```bash
go get .
go run .
```

**3.Start the Frontend**
```bash
cd frontend
npm install
npm run dev
```

