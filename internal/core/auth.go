package core

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func AuthCheck(w http.ResponseWriter, r *http.Request) {
	var username string
	var usertype int
	var token string
	var err error
	if r.Method == http.MethodPost {
		c, err := r.Cookie("session_token")
		if err != nil {
			if val := r.URL.Query().Get("token"); val != "" {
				username, usertype, err = Engine.UDb.ValidateToken(val)
				if err != nil {
					username, usertype, token, err = rauthHelper(w, r)
					if err != nil {
						Warn.Printf("not authorized (%s)\n", r.RemoteAddr)
						http.Error(w, "invalid credentials", http.StatusBadRequest)
						return
					}
				} else {
					token = val
				}
			} else {
				username, usertype, token, err = rauthHelper(w, r)
				if err != nil {
					Warn.Printf("not authorized (%s)\n", r.RemoteAddr)
					http.Error(w, "invalid credentials", http.StatusBadRequest)
					return
				}
			}
		} else {
			username, usertype, err = Engine.UDb.ValidateToken(c.Value)
			if err != nil {
				username, usertype, token, err = rauthHelper(w, r)
				if err != nil {
					Warn.Printf("not authorized (%s)\n", r.RemoteAddr)
					http.Error(w, "invalid credentials", http.StatusBadRequest)
					return
				}
			} else {
				token = c.Value
			}
		}
	} else {
		username, usertype, token, err = authHelper(w, r)
		if err != nil {
			Warn.Printf("%s (%s)\n", err, r.RemoteAddr)
			return
		}
	}
	var cmsg []byte
	if usertype == 0 {
		cmsg, err = json.Marshal(&ConnectionMsg{Type: "user", Session: token})
	} else if usertype == 1 {
		cmsg, err = json.Marshal(&ConnectionMsg{Type: "admin", Session: token})
	} else if usertype == -1 {
		cmsg, err = json.Marshal(&ConnectionMsg{Type: "disabled", Session: token})
	}
	if err != nil {
		return
	}
	_, err = w.Write(cmsg)
	if err != nil {
		return
	}
	Info.Printf("User %s (%s) Authenticated", username, r.RemoteAddr)

}

func authHelper(w http.ResponseWriter, r *http.Request) (username string, usertype int, token string, err error) {
	c, err := r.Cookie("session_token")
	if err != nil {
		if val := r.URL.Query().Get("token"); val != "" {
			username, usertype, err = Engine.UDb.ValidateToken(val)
			if err != nil {
				return bauthHelper(w, r)
			}
			token = val
		} else {
			return bauthHelper(w, r)
		}
	} else {
		username, usertype, err = Engine.UDb.ValidateToken(c.Value)
		if err != nil {
			return bauthHelper(w, r)
		}
		token = c.Value
	}
	return username, usertype, token, nil
}

func bauthHelper(w http.ResponseWriter, r *http.Request) (username string, usertype int, token string, err error) {
	defer func() {
		if r := recover(); r != nil { // uuid may panic
			username = ""
			usertype = -1
			token = ""
			err = fmt.Errorf("uuid error")
			return
		}
	}()
	var password string
	var ok bool
	w.Header().Set("WWW-Authenticate", `Basic realm="Protected"`)
	username, password, ok = r.BasicAuth()
	if !ok {
		http.Error(w, "Not authorized", http.StatusUnauthorized)
		return "", -1, "", fmt.Errorf("not authorized")
	}

	if len(username) < 5 || len(password) < 5 {
		http.Error(w, "Not authorized", http.StatusUnauthorized)
		return "", -1, "", fmt.Errorf("not authorized")
	}

	usertype, ok = Engine.UDb.Validate(username, password)
	if !ok {
		http.Error(w, "Not authorized", http.StatusUnauthorized)
		return "", -1, "", fmt.Errorf("not authorized")
	}
	token = uuid.New().String()
	err = Engine.UDb.SetToken(username, token)
	if err != nil {
		http.Error(w, "Not authorized", http.StatusUnauthorized)
		return "", -1, "", fmt.Errorf("not authorized")
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Path:     "/",
		Value:    token,
		Expires:  time.Now().Add(48 * time.Hour),
		SameSite: http.SameSiteStrictMode,
	})
	return username, usertype, token, nil
}

func rauthHelper(w http.ResponseWriter, r *http.Request) (username string, usertype int, token string, err error) {
	defer func() {
		if r := recover(); r != nil { // uuid may panic
			username = ""
			usertype = -1
			token = ""
			err = fmt.Errorf("uuid error")
			return
		}
	}()
	var ok bool
	var cred ConReq
	reqreader := io.LimitReader(r.Body, 1048576) // 1MB limiter
	err = json.NewDecoder(reqreader).Decode(&cred)
	if err != nil {
		return "", -1, "", err
	}
	if len(cred.Data1) < 5 || len(cred.Data2) < 5 {
		return "", -1, "", fmt.Errorf("invalid credentials")
	}
	usertype, ok = Engine.UDb.Validate(cred.Data1, cred.Data2)
	if !ok {
		return "", -1, "", fmt.Errorf("invalid credentials")
	}
	token = uuid.New().String()
	err = Engine.UDb.SetToken(cred.Data1, token)
	if err != nil {
		return "", -1, "", fmt.Errorf("not authorized")
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Path:     "/",
		Value:    token,
		Expires:  time.Now().Add(48 * time.Hour),
		SameSite: http.SameSiteStrictMode,
	})
	username = cred.Data1
	return username, usertype, token, nil
}
