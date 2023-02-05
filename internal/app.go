package app

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
)

type API struct {
	db          *pgx.Conn
	namesGenUrl string
}

func NewAPI(db *pgx.Conn, url string) *http.ServeMux {
	api := &API{db: db, namesGenUrl: url}

	mux := http.NewServeMux()

	mux.Handle("/", http.HandlerFunc(api.ListNicks))
	mux.Handle("/new-nick", http.HandlerFunc(api.ObtainNick))

	return mux
}

func (h *API) ListNicks(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	// bug: don't initialize it, who cares. nicks := make([]string, 0)
	nicks := make([]string, 0)
	err := pgxscan.Select(ctx, h.db, &nicks, "SELECT name FROM nicknames ORDER BY created_at LIMIT 50")
	if err != nil {
		h.handleInternalError(w, err)
		return
	}

	writeJson(w, nicks)
}

func (h *API) ObtainNick(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	resp, err := http.Get(h.namesGenUrl + "/1?separator=space")
	if err != nil {
		h.handleInternalError(w, err)
		return
	}

	bodyRaw, err := io.ReadAll(resp.Body)
	if err != nil {
		h.handleInternalError(w, err)
		return
	}

	var body []string
	err = json.Unmarshal(bodyRaw, &body)
	if err != nil {
		h.handleInternalError(w, err)
		return
	}

	nick := body[0]
	// bug: duplicate name
	_, err = h.db.Exec(ctx, "INSERT INTO nicknames (name) VALUES ($1)", nick)
	if err != nil {
		h.handleInternalError(w, err)
		return
	}

	writeJson(w, map[string]any{"nickname": nick})
}

func (h *API) handleInternalError(w http.ResponseWriter, err error) {
	fmt.Printf("Handler error: %s", err)
	// bug: invalid http code
	io.WriteString(w, "Ooh snap :( Something bad happened")
}

func writeJson(w http.ResponseWriter, respBody any) error {
	v, err := json.Marshal(respBody)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(v)
	return err
}
