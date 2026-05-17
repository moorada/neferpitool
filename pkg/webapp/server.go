package webapp

import (
	"encoding/json"
	"fmt"
	"html/template"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/moorada/neferpitool/pkg/app"
	"github.com/moorada/neferpitool/pkg/configuration"
	"github.com/moorada/neferpitool/pkg/domains"
	"github.com/moorada/neferpitool/pkg/format"
	"github.com/moorada/neferpitool/pkg/log"
	"golang.org/x/net/idna"
)

type Server struct {
	service *app.Service
	bg      *backgroundManager
	mux     *http.ServeMux
	index   *template.Template
	opMu    sync.Mutex
	opSeq   int64
	ops     map[string]*operationStatus
}

type apiResponse struct {
	OK      bool        `json:"ok"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type typoDomainDTO struct {
	Name         string `json:"name"`
	Status       string `json:"status"`
	Algorithm    string `json:"algorithm"`
	LegitDomain  string `json:"legit_domain"`
	Registrar    string `json:"registrar"`
	ExpiryDate   string `json:"expiry_date"`
	DateCheck    string `json:"date_check"`
	Unicode      string `json:"unicode"`
	ErrorWhois   string `json:"error_whois"`
	ErrorStatus  string `json:"error_status"`
	TimeWhois    string `json:"time_whois"`
	TimeStatus   string `json:"time_status"`
	CreatedAtISO string `json:"created_at_iso"`
}

type operationStatus struct {
	ID              string            `json:"id"`
	Type            string            `json:"type"`
	Domain          string            `json:"domain"`
	State           string            `json:"state"`
	Message         string            `json:"message"`
	Done            int               `json:"done"`
	Total           int               `json:"total"`
	ProgressPercent int               `json:"progress_percent"`
	TypoCount       int               `json:"typodomains_count"`
	Errors          map[string]string `json:"errors,omitempty"`
	Error           string            `json:"error,omitempty"`
	StartedAt       time.Time         `json:"started_at"`
	UpdatedAt       time.Time         `json:"updated_at"`
}

func NewServer() *Server {
	_ = log.ActiveConsoleLog()
	service := app.NewService()
	s := &Server{
		service: service,
		bg:      newBackgroundManager(service),
		mux:     http.NewServeMux(),
		index:   template.Must(template.New("page").Parse(dashboardPage)),
		ops:     map[string]*operationStatus{},
	}
	s.registerRoutes()
	return s
}

func (s *Server) Handler() http.Handler {
	return s.mux
}

func (s *Server) Start(addr string) error {
	if strings.TrimSpace(addr) == "" {
		addr = ":8080"
	}
	log.Info("Web app listening on %s", addr)
	return http.ListenAndServe(addr, s.mux)
}

func (s *Server) registerRoutes() {
	s.mux.HandleFunc("GET /", s.handleIndex)
	s.mux.HandleFunc("GET /logo.png", s.handleLogo)

	s.mux.HandleFunc("GET /api/overview", s.handleOverview)
	s.mux.HandleFunc("GET /api/config", s.handleConfig)
	s.mux.HandleFunc("GET /api/stats", s.handleStats)
	s.mux.HandleFunc("GET /api/expiration", s.handleExpiration)
	s.mux.HandleFunc("GET /api/reliable", s.handleReliableChanges)
	s.mux.HandleFunc("POST /api/presence", s.handlePresence)

	s.mux.HandleFunc("POST /api/domain/add", s.handleAddDomain)
	s.mux.HandleFunc("POST /api/domain/add/start", s.handleAddDomainStart)
	s.mux.HandleFunc("POST /api/domain/import", s.handleImportTypos)
	s.mux.HandleFunc("GET /api/domain/{domain}/typos", s.handleDomainTypos)
	s.mux.HandleFunc("GET /api/domain/{domain}/typos/page", s.handleDomainTyposPage)
	s.mux.HandleFunc("POST /api/domain/{domain}/remove", s.handleRemoveDomain)
	s.mux.HandleFunc("POST /api/domain/{domain}/check", s.handleCheckDomain)
	s.mux.HandleFunc("POST /api/check/all", s.handleCheckAll)

	s.mux.HandleFunc("GET /api/typodomain/{name}/history", s.handleTypoHistory)
	s.mux.HandleFunc("POST /api/typodomain/{name}/update", s.handleTypoUpdate)
	s.mux.HandleFunc("POST /api/typodomain/{name}/delete", s.handleTypoDelete)

	s.mux.HandleFunc("POST /api/background/start", s.handleBackgroundStart)
	s.mux.HandleFunc("POST /api/background/stop", s.handleBackgroundStop)
	s.mux.HandleFunc("POST /api/background/run", s.handleBackgroundRunOnce)

	s.mux.HandleFunc("GET /api/ops/{id}", s.handleOperationStatus)
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	_ = s.index.Execute(w, nil)
}

func (s *Server) handleLogo(w http.ResponseWriter, r *http.Request) {
	candidates := []string{"logo.png", "../logo.png"}
	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			http.ServeFile(w, r, filepath.Clean(p))
			return
		}
	}
	http.NotFound(w, r)
}

func (s *Server) handleOverview(w http.ResponseWriter, r *http.Request) {
	domainsList := s.service.ListDomains()
	sort.Slice(domainsList, func(i, j int) bool { return domainsList[i].Name < domainsList[j].Name })

	var names []string
	for _, d := range domainsList {
		names = append(names, d.Name)
	}

	s.writeOK(w, map[string]interface{}{
		"domains":           names,
		"domains_count":     len(names),
		"background_status": s.bg.Status(),
	})
}

func (s *Server) handleConfig(w http.ResponseWriter, r *http.Request) {
	conf := configuration.GetConf()
	s.writeOK(w, conf)
}

func (s *Server) handleStats(w http.ResponseWriter, r *http.Request) {
	tds := s.service.ListAllTypoDomains()
	errorsCount := 0
	statusCount := map[string]int{}
	for _, td := range tds {
		if td.ErrorStatus != "" || td.ErrorWhois != "" {
			errorsCount++
		}
		statusCount[td.StatusToString()]++
	}
	s.writeOK(w, map[string]interface{}{
		"typodomains_total":  len(tds),
		"typodomains_errors": errorsCount,
		"status_count":       statusCount,
	})
}

func (s *Server) handleExpiration(w http.ResponseWriter, r *http.Request) {
	tds := s.service.GetTypoDomainsInExpiration()
	s.writeOK(w, s.toTypoDTOList(tds))
}

func (s *Server) handleReliableChanges(w http.ResponseWriter, r *http.Request) {
	changesList, err := s.service.ListReliableChanges()
	if err != nil {
		s.writeErr(w, http.StatusInternalServerError, err)
		return
	}
	s.writeOK(w, changesList)
}

func (s *Server) handlePresence(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Domains []string `json:"domains"`
	}
	if err := decodeJSON(r, &req); err != nil {
		s.writeErr(w, http.StatusBadRequest, err)
		return
	}
	s.writeOK(w, s.service.DomainPresence(req.Domains))
}

func (s *Server) handleAddDomain(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Domain           string `json:"domain"`
		IncludeTypoItems *bool  `json:"include_typodomains"`
	}
	if err := decodeJSON(r, &req); err != nil {
		s.writeErr(w, http.StatusBadRequest, err)
		return
	}
	req.Domain = strings.TrimSpace(req.Domain)
	if req.Domain == "" {
		s.writeErr(w, http.StatusBadRequest, fmt.Errorf("domain is required"))
		return
	}

	tds, errs, err := s.service.AddDomainAndTypos(req.Domain, nil)
	if err != nil {
		s.writeErr(w, http.StatusInternalServerError, err)
		return
	}
	resp := map[string]interface{}{
		"typodomains_count": len(tds),
		"errors":            formatErrMap(errs),
	}
	includeItems := true
	if req.IncludeTypoItems != nil {
		includeItems = *req.IncludeTypoItems
	}
	if includeItems {
		resp["typodomains"] = s.toTypoDTOList(tds)
	}
	s.writeOK(w, resp)
}

func (s *Server) handleImportTypos(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Domain           string `json:"domain"`
		Typos            string `json:"typos"`
		IncludeTypoItems *bool  `json:"include_typodomains"`
	}
	if err := decodeJSON(r, &req); err != nil {
		s.writeErr(w, http.StatusBadRequest, err)
		return
	}
	req.Domain = strings.TrimSpace(req.Domain)
	if req.Domain == "" {
		s.writeErr(w, http.StatusBadRequest, fmt.Errorf("domain is required"))
		return
	}
	lines := strings.Split(req.Typos, "\n")
	tds, errs, err := s.service.ImportTyposFromLines(req.Domain, lines, nil)
	if err != nil {
		s.writeErr(w, http.StatusBadRequest, err)
		return
	}

	resp := map[string]interface{}{
		"typodomains_count": len(tds),
		"errors":            formatErrMap(errs),
	}
	includeItems := true
	if req.IncludeTypoItems != nil {
		includeItems = *req.IncludeTypoItems
	}
	if includeItems {
		resp["typodomains"] = s.toTypoDTOList(tds)
	}

	s.writeOK(w, resp)
}

func (s *Server) handleAddDomainStart(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Domain string `json:"domain"`
	}
	if err := decodeJSON(r, &req); err != nil {
		s.writeErr(w, http.StatusBadRequest, err)
		return
	}
	domain := strings.TrimSpace(req.Domain)
	if domain == "" {
		s.writeErr(w, http.StatusBadRequest, fmt.Errorf("domain is required"))
		return
	}

	op := s.newOperation("add_domain", domain, "Preparing typodomain generation...")
	go s.runAddDomainOperation(op.ID, domain)
	s.writeOK(w, map[string]string{"operation_id": op.ID})
}

func (s *Server) handleDomainTypos(w http.ResponseWriter, r *http.Request) {
	domain := strings.TrimSpace(r.PathValue("domain"))
	if domain == "" {
		s.writeErr(w, http.StatusBadRequest, fmt.Errorf("domain is required"))
		return
	}
	tds := s.service.ListTypoDomains(domain)
	sort.Slice(tds, func(i, j int) bool { return tds[i].Name < tds[j].Name })
	s.writeOK(w, s.toTypoDTOList(tds))
}

func (s *Server) handleDomainTyposPage(w http.ResponseWriter, r *http.Request) {
	domain := strings.TrimSpace(r.PathValue("domain"))
	if domain == "" {
		s.writeErr(w, http.StatusBadRequest, fmt.Errorf("domain is required"))
		return
	}

	page := parsePositiveInt(r.URL.Query().Get("page"), 1)
	size := parsePositiveInt(r.URL.Query().Get("size"), 10)
	if size > 200 {
		size = 200
	}
	sortBy := strings.TrimSpace(r.URL.Query().Get("sort_by"))
	sortDir := strings.TrimSpace(r.URL.Query().Get("sort_dir"))
	filter := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("q")))
	statusFilter := parseFilterValues(r.URL.Query().Get("status"))
	algorithmFilter := parseFilterValues(r.URL.Query().Get("algorithm"))
	registrarFilter := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("registrar")))
	nameFilter := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("name")))
	legitDomainFilter := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("legit_domain")))
	hasErrorsFilter, hasErrorsFilterSet := parseOptionalBool(r.URL.Query().Get("has_errors"))
	expiryFrom, hasExpiryFrom := parseOptionalDate(r.URL.Query().Get("expiry_from"))
	expiryTo, hasExpiryTo := parseOptionalDate(r.URL.Query().Get("expiry_to"))

	tds := s.service.ListTypoDomains(domain)
	items := s.toTypoDTOList(tds)
	availableStatuses, availableAlgorithms := collectAvailableFilterValues(items)

	filtered := make([]typoDomainDTO, 0, len(items))
	for _, td := range items {
		if len(statusFilter) > 0 {
			if _, ok := statusFilter[strings.ToLower(td.Status)]; !ok {
				continue
			}
		}
		if len(algorithmFilter) > 0 {
			if _, ok := algorithmFilter[strings.ToLower(td.Algorithm)]; !ok {
				continue
			}
		}
		if registrarFilter != "" && !strings.Contains(strings.ToLower(td.Registrar), registrarFilter) {
			continue
		}
		if nameFilter != "" && !strings.Contains(strings.ToLower(td.Name), nameFilter) {
			continue
		}
		if legitDomainFilter != "" && !strings.Contains(strings.ToLower(td.LegitDomain), legitDomainFilter) {
			continue
		}
		if hasErrorsFilterSet {
			hasErrors := td.ErrorWhois != "" || td.ErrorStatus != ""
			if hasErrors != hasErrorsFilter {
				continue
			}
		}
		if hasExpiryFrom || hasExpiryTo {
			expiryDate, ok := parseTypoExpiryDate(td.ExpiryDate)
			if !ok {
				continue
			}
			if hasExpiryFrom && expiryDate.Before(expiryFrom) {
				continue
			}
			if hasExpiryTo && expiryDate.After(expiryTo) {
				continue
			}
		}
		if filter != "" {
			candidate := strings.ToLower(strings.Join([]string{
				td.Name,
				td.Status,
				td.Algorithm,
				td.Registrar,
				td.Unicode,
				td.LegitDomain,
			}, " "))
			if !strings.Contains(candidate, filter) {
				continue
			}
		}
		filtered = append(filtered, td)
	}
	items = filtered

	applyTypoSort(items, sortBy, sortDir)

	total := len(items)
	pages := 0
	if total > 0 {
		pages = int(math.Ceil(float64(total) / float64(size)))
	}
	if pages > 0 && page > pages {
		page = pages
	}
	start := 0
	end := 0
	if total > 0 {
		start = (page - 1) * size
		if start < 0 {
			start = 0
		}
		if start > total {
			start = total
		}
		end = start + size
		if end > total {
			end = total
		}
	}

	s.writeOK(w, map[string]interface{}{
		"items": items[start:end],
		"page":  page,
		"size":  size,
		"pages": pages,
		"total": total,
		"q":     filter,
		"sort": map[string]string{
			"by":  normalizeSortBy(sortBy),
			"dir": normalizeSortDir(sortDir),
		},
		"available_statuses":   availableStatuses,
		"available_algorithms": availableAlgorithms,
	})
}

func (s *Server) handleRemoveDomain(w http.ResponseWriter, r *http.Request) {
	domain := strings.TrimSpace(r.PathValue("domain"))
	if domain == "" {
		s.writeErr(w, http.StatusBadRequest, fmt.Errorf("domain is required"))
		return
	}
	s.service.RemoveDomain(domain)
	s.writeOK(w, map[string]string{"domain": domain})
}

func (s *Server) handleCheckDomain(w http.ResponseWriter, r *http.Request) {
	domain := strings.TrimSpace(r.PathValue("domain"))
	if domain == "" {
		s.writeErr(w, http.StatusBadRequest, fmt.Errorf("domain is required"))
		return
	}
	changesList, errs := s.service.CheckChangesForDomain(domain, nil)
	s.writeOK(w, map[string]interface{}{
		"domain":  domain,
		"changes": changesList,
		"errors":  formatErrMap(errs),
	})
}

func (s *Server) handleCheckAll(w http.ResponseWriter, r *http.Request) {
	changesList, errs := s.service.CheckChangesForAll(nil)
	s.writeOK(w, map[string]interface{}{
		"changes": changesList,
		"errors":  formatErrMap(errs),
	})
}

func (s *Server) handleTypoHistory(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimSpace(r.PathValue("name"))
	if name == "" {
		s.writeErr(w, http.StatusBadRequest, fmt.Errorf("typodomain is required"))
		return
	}
	tds := s.service.ListTypoDomainHistory(name)
	sort.Slice(tds, func(i, j int) bool { return tds[i].CreatedAt.Before(tds[j].CreatedAt) })
	s.writeOK(w, s.toTypoDTOList(tds))
}

func (s *Server) handleTypoUpdate(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimSpace(r.PathValue("name"))
	if name == "" {
		s.writeErr(w, http.StatusBadRequest, fmt.Errorf("typodomain is required"))
		return
	}
	td, err := s.service.UpdateTypoDomain(name)
	if err != nil {
		s.writeErr(w, http.StatusInternalServerError, err)
		return
	}
	s.writeOK(w, s.toTypoDTO(td))
}

func (s *Server) handleTypoDelete(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimSpace(r.PathValue("name"))
	if name == "" {
		s.writeErr(w, http.StatusBadRequest, fmt.Errorf("typodomain is required"))
		return
	}
	s.service.RemoveTypoDomain(name)
	s.writeOK(w, map[string]string{"typodomain": name})
}

func (s *Server) handleBackgroundStart(w http.ResponseWriter, r *http.Request) {
	if err := s.bg.Start(); err != nil {
		s.writeErr(w, http.StatusInternalServerError, err)
		return
	}
	s.writeOK(w, s.bg.Status())
}

func (s *Server) handleBackgroundStop(w http.ResponseWriter, r *http.Request) {
	s.bg.Stop()
	s.writeOK(w, s.bg.Status())
}

func (s *Server) handleBackgroundRunOnce(w http.ResponseWriter, r *http.Request) {
	err := s.bg.RunOnce()
	if err != nil {
		s.writeErr(w, http.StatusInternalServerError, err)
		return
	}
	s.writeOK(w, s.bg.Status())
}

func (s *Server) handleOperationStatus(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.PathValue("id"))
	if id == "" {
		s.writeErr(w, http.StatusBadRequest, fmt.Errorf("operation id is required"))
		return
	}
	op, ok := s.getOperation(id)
	if !ok {
		s.writeErr(w, http.StatusNotFound, fmt.Errorf("operation not found"))
		return
	}
	s.writeOK(w, op)
}

func (s *Server) toTypoDTO(td domains.TypoDomain) typoDomainDTO {
	var p *idna.Profile
	p = idna.New()
	nameUnicode, err := p.ToUnicode(td.Name)
	if err != nil {
		nameUnicode = td.Name
	}
	return typoDomainDTO{
		Name:         td.Name,
		Status:       td.StatusToString(),
		Algorithm:    td.Algorithm,
		LegitDomain:  td.LegitDomain,
		Registrar:    format.NormalizeText(td.Whois.Parsed.Registrar.RegistrarName),
		ExpiryDate:   td.GetExpiryDateString(),
		DateCheck:    td.CreatedAt.Format("02/01/2006 15:04"),
		Unicode:      format.NormalizeText(nameUnicode),
		ErrorWhois:   format.NormalizeText(td.ErrorWhois),
		ErrorStatus:  format.NormalizeText(td.ErrorStatus),
		TimeWhois:    td.TimeWhois,
		TimeStatus:   td.TimeStatus,
		CreatedAtISO: td.CreatedAt.Format(time.RFC3339),
	}
}

func (s *Server) toTypoDTOList(tds domains.TypoList) []typoDomainDTO {
	out := make([]typoDomainDTO, 0, len(tds))
	for _, td := range tds {
		out = append(out, s.toTypoDTO(td))
	}
	return out
}

func (s *Server) writeOK(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(apiResponse{
		OK:   true,
		Data: data,
	})
}

func (s *Server) writeErr(w http.ResponseWriter, code int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(apiResponse{
		OK:    false,
		Error: err.Error(),
	})
}

func decodeJSON(r *http.Request, dst interface{}) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	return dec.Decode(dst)
}

func formatErrMap(errs map[string]error) map[string]string {
	out := map[string]string{}
	for k, v := range errs {
		if v != nil {
			out[k] = v.Error()
		}
	}
	return out
}

func parsePositiveInt(raw string, fallback int) int {
	if strings.TrimSpace(raw) == "" {
		return fallback
	}
	n, err := strconv.Atoi(raw)
	if err != nil || n <= 0 {
		return fallback
	}
	return n
}

func parseFilterValues(raw string) map[string]struct{} {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil
	}
	values := map[string]struct{}{}
	for _, part := range strings.Split(raw, ",") {
		v := strings.ToLower(strings.TrimSpace(part))
		if v == "" {
			continue
		}
		values[v] = struct{}{}
	}
	if len(values) == 0 {
		return nil
	}
	return values
}

func collectAvailableFilterValues(items []typoDomainDTO) ([]string, []string) {
	statuses := map[string]struct{}{}
	algorithms := map[string]struct{}{}
	for _, item := range items {
		s := strings.TrimSpace(item.Status)
		if s != "" {
			statuses[s] = struct{}{}
		}
		a := strings.TrimSpace(item.Algorithm)
		if a != "" {
			algorithms[a] = struct{}{}
		}
	}

	statusList := make([]string, 0, len(statuses))
	for v := range statuses {
		statusList = append(statusList, v)
	}
	sort.Slice(statusList, func(i, j int) bool {
		return strings.ToLower(statusList[i]) < strings.ToLower(statusList[j])
	})

	algorithmList := make([]string, 0, len(algorithms))
	for v := range algorithms {
		algorithmList = append(algorithmList, v)
	}
	sort.Slice(algorithmList, func(i, j int) bool {
		return strings.ToLower(algorithmList[i]) < strings.ToLower(algorithmList[j])
	})

	return statusList, algorithmList
}

func parseOptionalBool(raw string) (bool, bool) {
	v := strings.ToLower(strings.TrimSpace(raw))
	switch v {
	case "":
		return false, false
	case "true", "1", "yes", "y":
		return true, true
	case "false", "0", "no", "n":
		return false, true
	default:
		return false, false
	}
}

func parseOptionalDate(raw string) (time.Time, bool) {
	value := strings.TrimSpace(raw)
	if value == "" {
		return time.Time{}, false
	}
	parsed, err := format.StringToTime(value)
	if err != nil {
		return time.Time{}, false
	}
	return parsed, true
}

func parseTypoExpiryDate(raw string) (time.Time, bool) {
	return parseOptionalDate(raw)
}

func normalizeSortBy(raw string) string {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "name", "status", "algorithm", "legit_domain", "registrar", "expiry_date", "date_check", "unicode", "error_status", "error_whois", "time_whois", "time_status":
		return strings.ToLower(strings.TrimSpace(raw))
	default:
		return "name"
	}
}

func normalizeSortDir(raw string) string {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "desc":
		return "desc"
	default:
		return "asc"
	}
}

func applyTypoSort(items []typoDomainDTO, sortBy string, sortDir string) {
	by := normalizeSortBy(sortBy)
	desc := normalizeSortDir(sortDir) == "desc"

	sort.SliceStable(items, func(i, j int) bool {
		cmp := compareTypoForSort(items[i], items[j], by)
		if desc {
			return cmp > 0
		}
		return cmp < 0
	})
}

func compareTypoForSort(a typoDomainDTO, b typoDomainDTO, sortBy string) int {
	compareString := func(x string, y string) int {
		xs := strings.ToLower(strings.TrimSpace(x))
		ys := strings.ToLower(strings.TrimSpace(y))
		switch {
		case xs < ys:
			return -1
		case xs > ys:
			return 1
		default:
			return 0
		}
	}

	compareDate := func(rawA string, rawB string) int {
		ta, oka := parseOptionalDate(rawA)
		tb, okb := parseOptionalDate(rawB)
		switch {
		case oka && okb:
			if ta.Before(tb) {
				return -1
			}
			if ta.After(tb) {
				return 1
			}
			return 0
		case oka && !okb:
			return -1
		case !oka && okb:
			return 1
		default:
			return 0
		}
	}

	var cmp int
	switch sortBy {
	case "status":
		cmp = compareString(a.Status, b.Status)
	case "algorithm":
		cmp = compareString(a.Algorithm, b.Algorithm)
	case "legit_domain":
		cmp = compareString(a.LegitDomain, b.LegitDomain)
	case "registrar":
		cmp = compareString(a.Registrar, b.Registrar)
	case "expiry_date":
		cmp = compareDate(a.ExpiryDate, b.ExpiryDate)
	case "date_check":
		cmp = compareDate(a.CreatedAtISO, b.CreatedAtISO)
	case "unicode":
		cmp = compareString(a.Unicode, b.Unicode)
	case "error_status":
		cmp = compareString(a.ErrorStatus, b.ErrorStatus)
	case "error_whois":
		cmp = compareString(a.ErrorWhois, b.ErrorWhois)
	case "time_whois":
		cmp = compareString(a.TimeWhois, b.TimeWhois)
	case "time_status":
		cmp = compareString(a.TimeStatus, b.TimeStatus)
	default:
		cmp = compareString(a.Name, b.Name)
	}
	if cmp != 0 {
		return cmp
	}
	return compareString(a.Name, b.Name)
}

func (s *Server) newOperation(kind, domain, message string) *operationStatus {
	s.opMu.Lock()
	defer s.opMu.Unlock()

	s.opSeq++
	id := fmt.Sprintf("%d-%d", time.Now().UnixNano(), s.opSeq)
	now := time.Now()
	op := &operationStatus{
		ID:        id,
		Type:      kind,
		Domain:    domain,
		State:     "running",
		Message:   message,
		StartedAt: now,
		UpdatedAt: now,
	}
	s.ops[id] = op
	s.trimOperationsLocked(150)
	return cloneOperation(op)
}

func (s *Server) getOperation(id string) (operationStatus, bool) {
	s.opMu.Lock()
	defer s.opMu.Unlock()

	op, ok := s.ops[id]
	if !ok {
		return operationStatus{}, false
	}
	return *cloneOperation(op), true
}

func (s *Server) runAddDomainOperation(id, domain string) {
	progress := func(done, total int) {
		s.updateOperation(id, func(op *operationStatus) {
			op.Done = done
			op.Total = total
			op.Message = fmt.Sprintf("Scanning typodomains: %d/%d", done, total)
		})
	}

	tds, errs, err := s.service.AddDomainAndTypos(domain, progress)
	if err != nil {
		s.updateOperation(id, func(op *operationStatus) {
			op.State = "failed"
			op.Error = err.Error()
			op.Message = "Failed to add domain"
			op.Errors = formatErrMap(errs)
			op.TypoCount = len(tds)
		})
		return
	}

	s.updateOperation(id, func(op *operationStatus) {
		op.State = "done"
		op.Done = len(tds)
		op.Total = len(tds)
		op.Message = fmt.Sprintf("Done. Generated and scanned %d typodomains.", len(tds))
		op.Errors = formatErrMap(errs)
		op.TypoCount = len(tds)
	})
}

func (s *Server) updateOperation(id string, fn func(op *operationStatus)) {
	s.opMu.Lock()
	defer s.opMu.Unlock()

	op, ok := s.ops[id]
	if !ok {
		return
	}
	fn(op)
	op.ProgressPercent = calcPercent(op.Done, op.Total)
	op.UpdatedAt = time.Now()
}

func calcPercent(done, total int) int {
	if total <= 0 {
		return 0
	}
	if done <= 0 {
		return 0
	}
	if done >= total {
		return 100
	}
	return int(math.Round((float64(done) / float64(total)) * 100))
}

func cloneOperation(op *operationStatus) *operationStatus {
	c := *op
	if op.Errors != nil {
		c.Errors = map[string]string{}
		for k, v := range op.Errors {
			c.Errors[k] = v
		}
	}
	return &c
}

func (s *Server) trimOperationsLocked(limit int) {
	if len(s.ops) <= limit {
		return
	}
	type item struct {
		id string
		ts time.Time
	}
	items := make([]item, 0, len(s.ops))
	for id, op := range s.ops {
		items = append(items, item{id: id, ts: op.UpdatedAt})
	}
	sort.Slice(items, func(i, j int) bool { return items[i].ts.Before(items[j].ts) })
	toDrop := len(items) - limit
	for i := 0; i < toDrop; i++ {
		delete(s.ops, items[i].id)
	}
}
