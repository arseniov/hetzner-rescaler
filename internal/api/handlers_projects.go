package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/jonamat/hetzner-rescaler/internal/crypto"
	"github.com/jonamat/hetzner-rescaler/internal/store"
)

// KeyFromEnv reads the AES-GCM key from the environment. In production
// the cmd layer passes a keyring built from RESCALER_ENCRYPTION_KEY
// (via crypto.LoadKeyring). Tests build a fresh keyring via
// crypto.NewKeyring() and inject it via Deps.Keyring.
//
// This wrapper exists so handlers can call into a uniform interface
// without caring how the key was sourced.
func KeyFromEnv() (*crypto.Keyring, error) {
	return crypto.LoadKeyring()
}

// keyring resolves the active keyring for a Deps: prefer the injected
// one (set in tests), fall back to env (production).
func (d Deps) keyring() (*crypto.Keyring, error) {
	if d.Keyring != nil {
		return d.Keyring, nil
	}
	return KeyFromEnv()
}

func (d Deps) handleListProjects(w http.ResponseWriter, r *http.Request) {
	projects, err := d.Store.ListProjects()
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	out := make([]ProjectResponse, 0, len(projects))
	for _, p := range projects {
		out = append(out, projectToResponse(p))
	}
	writeJSON(w, http.StatusOK, out)
}

func (d Deps) handleCreateProject(w http.ResponseWriter, r *http.Request) {
	var req CreateProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	req.Name = strings.TrimSpace(req.Name)
	req.HCloudToken = strings.TrimSpace(req.HCloudToken)
	if req.Name == "" || req.HCloudToken == "" {
		writeJSONError(w, http.StatusBadRequest, "name and hcloud_token are required")
		return
	}

	key, err := d.keyring()
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "key unavailable: "+err.Error())
		return
	}
	encToken, nonce, err := key.Seal([]byte(req.HCloudToken))
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "encrypt token: "+err.Error())
		return
	}

	p, err := d.Store.CreateProject(req.Name, encToken, nonce)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Auto-populate servers from the Hetzner API for the newly created
	// project. This restores the behavior users expect: entering a new
	// project + token immediately surfaces the existing cloud servers in
	// the UI without a separate manual "Refresh" click.
	//
	// Fetch errors are non-fatal: a bad token is logged in the response
	// and we still return the created project so the UI can show a
	// helpful error state. The user can retry via POST /refresh.
	added, skipped, fetchErr := d.syncProjectServers(r.Context(), p.ID)

	resp := createProjectResponse{
		ProjectResponse: projectToResponse(p),
		Added:           added,
		Skipped:         skipped,
	}
	if fetchErr != nil {
		resp.LastError = fetchErr.Error()
	}
	writeJSON(w, http.StatusCreated, resp)
}

func (d Deps) handleDeleteProject(w http.ResponseWriter, r *http.Request) {
	id, ok := pathInt64(r, "id")
	if !ok {
		writeJSONError(w, http.StatusBadRequest, "invalid id")
		return
	}
	// Confirm the project exists so the second-DELETE case returns 404
	// rather than a silent no-op.
	if _, err := d.Store.GetProject(id); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeJSONError(w, http.StatusNotFound, "project not found")
			return
		}
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if err := d.Store.DeleteProject(id); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (d Deps) handleRefreshProject(w http.ResponseWriter, r *http.Request) {
	id, ok := pathInt64(r, "id")
	if !ok {
		writeJSONError(w, http.StatusBadRequest, "invalid id")
		return
	}
	if _, err := d.Store.GetProject(id); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeJSONError(w, http.StatusNotFound, "project not found")
			return
		}
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	added, skipped, fetchErr := d.syncProjectServers(r.Context(), id)
	if fetchErr != nil {
		// Match the legacy behavior: a Hetzner list failure bubbles up
		// as 502 so the user can act on it (e.g. invalid token).
		writeJSONError(w, http.StatusBadGateway, fetchErr.Error())
		return
	}

	writeJSON(w, http.StatusOK, RefreshProjectResponse{Added: added, Skipped: skipped})
}

// syncProjectServers fetches the project from Hetzner and inserts any
// new servers into the store. It is shared by handleCreateProject (for
// auto-population on creation) and handleRefreshProject (manual sync).
//
// Errors come from three sources:
//   - APIFor: client construction failed (e.g. missing/invalid token)
//   - ListServers: Hetzner API call failed
//   - CreateServer: store insert failed mid-loop
//
// The caller decides how to surface these — handleCreateProject
// returns 201 with LastError populated, handleRefreshProject returns 502.
func (d Deps) syncProjectServers(ctx context.Context, projectID int64) ([]ServerResponse, []ServerResponse, error) {
	api, err := d.APIFor(projectID)
	if err != nil {
		return nil, nil, fmt.Errorf("hetzner client: %w", err)
	}

	remoteServers, err := api.ListServers(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("hetzner list: %w", err)
	}

	existing, err := d.Store.ListServersByProject(projectID)
	if err != nil {
		return nil, nil, fmt.Errorf("list existing: %w", err)
	}
	existingByID := make(map[int]*store.Server, len(existing))
	for _, s := range existing {
		existingByID[s.HCloudServerID] = s
	}

	var added, skipped []ServerResponse
	for _, hs := range remoteServers {
		if hs == nil {
			continue
		}
		if existing, ok := existingByID[hs.ID]; ok {
			skipped = append(skipped, serverToResponse(existing))
			continue
		}
		baseName := ""
		if st := hs.ServerType; st != nil {
			baseName = st.Name
		}
		srv, err := d.Store.CreateServer(projectID, store.Server{
			HCloudServerID: hs.ID,
			Name:           hs.Name,
			Label:          hs.Name,
			BaseServerType: baseName,
			TopServerType:  baseName,
			FallbackChain:  []string{baseName},
			Mode:           "manual",
			Timezone:       "UTC",
		})
		if err != nil {
			return added, skipped, fmt.Errorf("create server: %w", err)
		}
		added = append(added, serverToResponse(srv))
	}

	return added, skipped, nil
}

// projectToResponse converts a *store.Project to its API projection.
// The encrypted token is never included — only a boolean indicating
// that one is stored.
func projectToResponse(p *store.Project) ProjectResponse {
	if p == nil {
		return ProjectResponse{}
	}
	return ProjectResponse{
		ID:        p.ID,
		Name:      p.Name,
		HasToken:  len(p.HCloudTokenEncrypted) > 0,
		LastError: "",
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}

// pathInt64 extracts an int64 path parameter by name. Returns (0, false)
// if the parameter is missing or unparseable.
func pathInt64(r *http.Request, name string) (int64, bool) {
	raw := r.PathValue(name)
	if raw == "" {
		return 0, false
	}
	v, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		return 0, false
	}
	return v, true
}

// writeJSON writes a JSON response with the given status code.
func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
