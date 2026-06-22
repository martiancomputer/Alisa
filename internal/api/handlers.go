package api

import (
	// "database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/martiancomputer/Alisa/internal/database/queries"
)

type APIHandler struct {
	CaseRepo  *queries.CaseRepository
	LogRepo   *queries.LogRepository
	StatsRepo *queries.StatsRepository
}

// HandleGetRecentCases marshals and returns historical infraction records.
func (h *APIHandler) HandleGetRecentCases(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "Bad Request: Missing user_id tracking constraint", http.StatusBadRequest)
		return
	}

	cases, err := h.CaseRepo.GetUserHistory(userID)
	if err != nil {
		log.Printf("API Error: Failed to poll infraction history: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(cases); err != nil {
		log.Printf("API Serialization Error: %v", err)
	}
}

// HandleGetAuditFeed extracts paginated server event telemetry streams directly to the client interface.
func (h *APIHandler) HandleGetAuditFeed(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Enforce default limit bounds to control heap allocation volume
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	rows, err := h.LogRepo.FetchAuditFeed(limit, offset)
	if err != nil {
		log.Printf("API Error: Failed to step through audit logs: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Stream formatting directly down the network wire row-by-row to bypass bulk array building in RAM
	w.Write([]byte("["))
	first := true
	for rows.Next() {
		if !first {
			w.Write([]byte(","))
		}
		first = false

		var id int64
		var eventType, uID, modID, chID, payload string
		var timestamp int64

		if err := rows.Scan(&id, &eventType, &uID, &modID, &chID, &timestamp, &payload); err != nil {
			log.Printf("API Streaming Scan Failure: %v", err)
			return
		}

		// Reconstruct raw JSON components directly into output stream
		item, _ := json.Marshal(map[string]interface{}{
			"log_id":        id,
			"event_type":    eventType,
			"user_id":       uID,
			"moderator_id":  modID,
			"channel_id":    chID,
			"timestamp":     timestamp,
			"payload":       json.RawMessage(payload), // Keeps JSON structural formatting intact
		})
		w.Write(item)
	}
	w.Write([]byte("]"))
}