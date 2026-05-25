package main

import (
	"archive/zip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/atulX7/pqc/internal/app"
)

const maxUploadBytes = 25 * 1024 * 1024

func main() {
	port := env("PORT", "8080")
	server := &http.Server{
		Addr:    ":" + port,
		Handler: routes(),
	}
	log.Printf("pqc web listening on :%s", port)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}

func routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", indexHandler)
	mux.HandleFunc("GET /healthz", healthHandler)
	mux.HandleFunc("POST /api/scan/sample", sampleScanHandler)
	mux.HandleFunc("POST /api/scan/upload", uploadScanHandler)
	return mux
}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func sampleScanHandler(w http.ResponseWriter, r *http.Request) {
	options, err := optionsFromRequest(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	options.Path = env("SAMPLE_REPO_PATH", "sample_repo")
	options.RulesPath = env("RULES_PATH", "rules/crypto_rules.yaml")
	scanAndWrite(w, options)
}

func uploadScanHandler(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadBytes)
	if err := r.ParseMultipartForm(maxUploadBytes); err != nil {
		writeError(w, http.StatusBadRequest, fmt.Errorf("read upload: %w", err))
		return
	}

	file, header, err := r.FormFile("repo")
	if err != nil {
		writeError(w, http.StatusBadRequest, fmt.Errorf("repo zip is required: %w", err))
		return
	}
	defer file.Close()
	if !strings.HasSuffix(strings.ToLower(header.Filename), ".zip") {
		writeError(w, http.StatusBadRequest, fmt.Errorf("upload must be a .zip file"))
		return
	}

	workDir, err := os.MkdirTemp("", "pqc-upload-*")
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	defer os.RemoveAll(workDir)

	zipPath := filepath.Join(workDir, "repo.zip")
	if err := saveUpload(zipPath, file); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	extractDir := filepath.Join(workDir, "repo")
	if err := unzip(zipPath, extractDir); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	options, err := optionsFromRequest(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	options.Path = extractDir
	options.RulesPath = env("RULES_PATH", "rules/crypto_rules.yaml")
	scanAndWrite(w, options)
}

func scanAndWrite(w http.ResponseWriter, options app.ScanOptions) {
	result, err := app.RunScan(options)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func optionsFromRequest(r *http.Request) (app.ScanOptions, error) {
	options := app.DefaultScanOptions()
	if err := r.ParseMultipartForm(maxUploadBytes); err != nil && !errors.Is(err, http.ErrNotMultipart) {
		return options, err
	}
	options.ApplicationName = formValue(r, "app_name", options.ApplicationName)
	options.Sensitivity = formValue(r, "sensitivity", options.Sensitivity)
	options.Exposure = formValue(r, "exposure", options.Exposure)
	options.BusinessCriticality = formValue(r, "business_criticality", options.BusinessCriticality)
	options.CryptoAgility = formValue(r, "crypto_agility", options.CryptoAgility)
	options.VendorDependency = formValue(r, "vendor_dependency", options.VendorDependency)
	options.MigrationComplexity = formValue(r, "migration_complexity", options.MigrationComplexity)
	years := formValue(r, "secrecy_lifetime_years", strconv.Itoa(options.SecrecyLifetimeYears))
	parsedYears, err := strconv.Atoi(years)
	if err != nil || parsedYears < 0 {
		return options, fmt.Errorf("secrecy_lifetime_years must be a non-negative integer")
	}
	options.SecrecyLifetimeYears = parsedYears
	return options, nil
}

func formValue(r *http.Request, key string, fallback string) string {
	value := strings.TrimSpace(r.FormValue(key))
	if value == "" {
		return fallback
	}
	return value
}

func saveUpload(path string, source io.Reader) error {
	target, err := os.Create(path)
	if err != nil {
		return err
	}
	defer target.Close()
	_, err = io.Copy(target, source)
	return err
}

func unzip(src string, dest string) error {
	reader, err := zip.OpenReader(src)
	if err != nil {
		return fmt.Errorf("open zip: %w", err)
	}
	defer reader.Close()

	for _, file := range reader.File {
		target := filepath.Join(dest, file.Name)
		if !strings.HasPrefix(target, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("zip contains unsafe path: %s", file.Name)
		}
		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(target, 0o755); err != nil {
				return err
			}
			continue
		}
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return err
		}
		if err := extractFile(file, target); err != nil {
			return err
		}
	}
	return nil
}

func extractFile(file *zip.File, target string) error {
	source, err := file.Open()
	if err != nil {
		return err
	}
	defer source.Close()

	dest, err := os.OpenFile(target, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
	if err != nil {
		return err
	}
	defer dest.Close()
	_, err = io.Copy(dest, source)
	return err
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(value); err != nil {
		log.Printf("write response: %v", err)
	}
}

func writeError(w http.ResponseWriter, status int, err error) {
	writeJSON(w, status, map[string]string{"error": err.Error()})
}

func env(key string, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}
