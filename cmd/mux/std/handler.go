package main

import (
	"encoding/json"
	"math/rand/v2"
	"net/http"
	"strconv"
	"strings"
)

type user struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email,omitempty"`
}

func createUser(w http.ResponseWriter, r *http.Request) {
	var req user
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "请求体不是合法 JSON")
		return
	}
	req.ID = strconv.FormatUint(rand.Uint64(), 10)
	writeJSON(w, http.StatusCreated, req)
}

func echoPath(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"path": r.URL.Path,
		"rest": strings.TrimPrefix(r.URL.Path, "/"),
	})
}

func panicHandler(w http.ResponseWriter, r *http.Request) {
	panic("故意 panic，演示 Recoverer")
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		logger.Error("写响应失败", "err", err)
	}
}

func writeError(w http.ResponseWriter, code int, msg string) {
	writeJSON(w, code, map[string]string{"error": msg})
}
