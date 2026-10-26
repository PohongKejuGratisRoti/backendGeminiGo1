package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/alexedwards/scs/v2"
	_ "modernc.org/sqlite"
)

// allowed fucking origin
var allowedFuckingOrigin = map[string]bool{
	"http://localhost:5173":     true,
	"http://172.16.10.167:5173": true,
}

func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if allowedFuckingOrigin[origin] {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		}

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

var mux = http.NewServeMux()

var server = &http.Server{
	Addr:    ":9000",
	Handler: sessionManager.LoadAndSave(cors(mux)),
}

// session manager
var sessionManager = scs.New()

func respondError(errorRes http.ResponseWriter, message string, status int) {
	errorRes.Header().Set("Content-Type", "application/json")
	errorRes.WriteHeader(status)

	json.NewEncoder(errorRes).Encode(map[string]string{
		"error": message,
	})
}

type loginReqModel struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func main() {
	DbErr := initDB()
	if DbErr != nil {
		panic(DbErr)
	}
	defer db.Close()

	mux.HandleFunc("GET /", ping)
	mux.HandleFunc("POST /user", userLogin)
	log.Println("Server running on http://localhost:9000")

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

// backend func endpoint
func ping(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte("Hello from Go!"))

}

func userLogin(res http.ResponseWriter, req *http.Request) {
	var input loginReqModel
	err := json.NewDecoder(req.Body).Decode(&input)

	// errpr jamd;er
	if err != nil {
		respondError(res, "Invalid JSON", http.StatusBadRequest)
		return
	}
	if input.Username == "" {
		respondError(res, "please input username", http.StatusBadRequest)
	}
	if input.Password == "" {
		respondError(res, "please input username", http.StatusBadRequest)
	}

	if !checkUser(input.Username, input.Password) {
		respondError(res, "user tidak ditemukan", http.StatusUnauthorized)
		return
	}

	currentUserID, err := checkID(input.Username)
	if err != nil {
		respondError(res, "id tidak ditemukan", http.StatusBadRequest)
		return
	}

	sessionManager.RenewToken(req.Context())
	sessionManager.Put(req.Context(), "userID", currentUserID)

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(map[string]string{"message": "ok"})
}
