package main

import (
	"encoding/json"
	"log"
	"os"
	"sort"
	"time"
)

// Entry adalah satu temuan/aksi tercatat selama migrasi.
type Entry struct {
	Phase  string `json:"phase"`
	Table  string `json:"table"`
	V1ID   string `json:"v1_id,omitempty"`
	Action string `json:"action"`
	Detail string `json:"detail,omitempty"`
}

type Report struct {
	DryRun     bool           `json:"dry_run"`
	StartedAt  time.Time      `json:"started_at"`
	FinishedAt time.Time      `json:"finished_at"`
	Counters   map[string]int `json:"counters"`
	Entries    []Entry        `json:"entries"`
}

func NewReport(dryRun bool) *Report {
	return &Report{
		DryRun:    dryRun,
		StartedAt: time.Now().UTC(),
		Counters:  map[string]int{},
	}
}

func (r *Report) Count(key string) {
	r.Counters[key]++
}

func (r *Report) Add(phase, table, v1ID, action, detail string) {
	r.Counters[table+"."+action]++
	r.Entries = append(r.Entries, Entry{Phase: phase, Table: table, V1ID: v1ID, Action: action, Detail: detail})
}

func (r *Report) WriteFile(path string) error {
	b, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, 0644)
}

func (r *Report) PrintSummary() {
	keys := make([]string, 0, len(r.Counters))
	for k := range r.Counters {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	log.Print("=== Ringkasan ===")
	for _, k := range keys {
		log.Printf("  %-55s %d", k, r.Counters[k])
	}
	log.Printf("  total entri report: %d", len(r.Entries))
}
