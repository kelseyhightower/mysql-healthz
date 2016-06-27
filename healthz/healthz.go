package healthz

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-sql-driver/mysql"
)

type Config struct {
	Hostname string
	Database DatabaseConfig
}

type DatabaseConfig struct {
	DriverName     string
	DataSourceName string
	DatabaseName   string
	Tables         []string
}

type handler struct {
	dc           *DBChecker
	databaseName string
	hostname     string
	metadata     map[string]string
	tables       []string
}

func Handler(hc *Config) (http.Handler, error) {
	dc, err := NewDBChecker(hc.Database.DriverName, hc.Database.DataSourceName)
	if err != nil {
		return nil, err
	}

	config, err := mysql.ParseDSN(hc.Database.DataSourceName)
	if err != nil {
		return nil, err
	}

	metadata := make(map[string]string)
	metadata["database_url"] = config.Addr
	metadata["database_user"] = config.User

	h := &handler{dc, hc.Database.DatabaseName, hc.Hostname, metadata, hc.Database.Tables}
	return h, nil
}

type Response struct {
	Hostname string            `json:"hostname"`
	Metadata map[string]string `json:"metadata"`
	Errors   []Error           `json:"errors"`
}

type Error struct {
	Description string            `json:"description"`
	Error       string            `json:"error"`
	Metadata    map[string]string `json:"metadata"`
	Type        string            `json:"type"`
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	response := Response{
		Hostname: h.hostname,
		Metadata: h.metadata,
	}

	statusCode := http.StatusOK

	errors := make([]Error, 0)

	err := h.dc.Reachable()
	if err != nil {
		errors = append(errors, Error{
			Type:        "DatabasePing",
			Description: "Database liveliness check.",
			Error:       err.Error(),
		})
	}

	for _, table := range h.tables {
		err = h.dc.TableExist(h.databaseName, table)
		if err != nil {
			metadata := make(map[string]string)
			metadata["table_name"] = table

			errors = append(errors, Error{
				Type:        "DatabaseTableExist",
				Description: "Database table exist check.",
				Error:       err.Error(),
				Metadata:    metadata,
			})
		}
	}

	response.Errors = errors
	if len(response.Errors) > 0 {
		statusCode = http.StatusInternalServerError
		for _, e := range response.Errors {
			log.Println(e.Error)
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)
	data, err := json.MarshalIndent(&response, "", "  ")
	if err != nil {
		log.Println(err)
	}
	w.Write(data)
}
