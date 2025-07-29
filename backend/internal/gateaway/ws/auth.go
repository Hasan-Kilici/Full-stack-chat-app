package ws

import (
    "encoding/json"
    "net/http"
    "time"

    "github.com/hasan-kilici/chat/internal/service/repository"
    "github.com/hasan-kilici/chat/pkg/auth"
)

type authResult struct {
    UserID      string
    Name        string
    SessionID   string
    Err         error
}

func checkJWTAuth(r *http.Request, w http.ResponseWriter) (*authResult, bool) {
    token := r.Header.Get("Sec-WebSocket-Protocol")
    if token == "" {
        writeJSONError(w, http.StatusUnauthorized, "WebSocket protocol header (token) missing")
        return nil, false
    }

    claims, err := auth.ParseJWT(token)
    if err != nil {
        writeJSONError(w, http.StatusUnauthorized, "Invalid or expired token")
        return nil, false
    }

    session, err := repository.GetSessionByID(claims.SessionID)
    if err != nil || session == nil || session.ExpiresAt.Before(time.Now()) {
        writeJSONError(w, http.StatusUnauthorized, "Session expired or invalid")
        return nil, false
    }

    user, err := repository.GetUserByID(claims.UserID)
    if err != nil || user == nil {
        writeJSONError(w, http.StatusUnauthorized, "User not found")
        return nil, false
    }

    return &authResult{
        UserID:     claims.UserID.String(),
        SessionID:  claims.SessionID.String(),
        Name:       user.Name,
        Err:        nil,
    }, true
}

func writeJSONError(w http.ResponseWriter, status int, msg string) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    _ = json.NewEncoder(w).Encode(map[string]string{"error": msg})
}