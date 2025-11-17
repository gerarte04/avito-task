package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"
	"time"

	vegeta "github.com/tsenart/vegeta/lib"
	"github.com/tsenart/vegeta/lib/plot"
)

type Team struct {
	Name    string  `json:"team_name" db:"name"`
	Members []*User `json:"members"`
}

type User struct {
	ID       string `json:"user_id" db:"id"`
	Name     string `json:"username" db:"name"`
	IsActive bool   `json:"is_active" db:"is_active"`

	TeamName string `json:"team_name,omitempty"`
}

type PullRequest struct {
	ID        string     `json:"pull_request_id" db:"id"`
	Name      string     `json:"pull_request_name" db:"name"`
	AuthorID  string     `json:"author_id" db:"author_id"`
	Status    string     `json:"status" db:"status"`
	Reviewers []string   `json:"assigned_reviewers"`
	CreatedAt *time.Time `json:"createdAt" db:"created_at"`
	MergedAt  *time.Time `json:"mergedAt" db:"merged_at"`
}

type SetIsActive struct {
	UserID   string `json:"user_id"`
	IsActive bool   `json:"is_active"`
}

type Reassign struct {
	PRID   string `json:"pull_request_id"`
	UserID string `json:"old_reviewer_id"`
}

const (
	MaxTeams = 20
	MaxUsers = 200
	MaxPRs   = 100
	RPS      = 5
	Duration = 120

	ResultsFile = "results.bin"
	ReportFile = "report.txt"
	PlotFile = "plot.html"

	Address = "http://localhost:8080"
)

func Rnd(n int) int {
	num, _ := rand.Int(rand.Reader, big.NewInt(int64(n)))
	return int(num.Int64())
}

var (
	nTeams = 0
	nPRs   = 0

	generators = []func(*vegeta.Target) {
		func(t *vegeta.Target) {
			teamName := fmt.Sprintf("team-%d", nTeams)
			nTeams = min(nTeams+1, MaxTeams)

			s := Team{
				Name: teamName,
			}

			for i := 0; i < MaxUsers; i++ {
				s.Members = append(s.Members, &User{
					ID:       fmt.Sprintf("user-%d", i),
					Name:     "aaa",
					IsActive: true,
				})
			}

			body, _ := json.Marshal(&s)

			t.Method = http.MethodPost
			t.URL = fmt.Sprintf("%s/team/add", Address)
			t.Body = body
		},
		func(t *vegeta.Target) {
			id := 0

			if nTeams > 0 {
				id = Rnd(nTeams)
			}

			t.Method = http.MethodGet
			t.URL = fmt.Sprintf("%s/team/get?team_name=team-%d", Address, id)
		},
		func(t *vegeta.Target) {
			id := Rnd(MaxUsers)

			body, _ := json.Marshal(&SetIsActive{
				UserID:   fmt.Sprintf("user-%d", id),
				IsActive: true,
			})

			t.Method = http.MethodPost
			t.URL = fmt.Sprintf("%s/users/setIsActive", Address)
			t.Body = body
		},
		func(t *vegeta.Target) {
			id := Rnd(MaxUsers)

			t.Method = http.MethodGet
			t.URL = fmt.Sprintf("%s/users/getReview?user-id=user-%d", Address, id)
		},
		func(t *vegeta.Target) {
			prID := fmt.Sprintf("pr-%d", nPRs)
			nPRs = min(nPRs+1, MaxPRs)

			body, _ := json.Marshal(&PullRequest{
				ID:       prID,
				Name:     "aaa",
				AuthorID: fmt.Sprintf("user-%d", Rnd(MaxUsers)),
			})

			t.Method = http.MethodPost
			t.URL = fmt.Sprintf("%s/pullRequest/create", Address)
			t.Body = body
		},
		func(t *vegeta.Target) {
			prID := 0

			if nPRs > 0 {
				prID = Rnd(nPRs)
			}

			body, _ := json.Marshal(&Reassign{
				PRID: fmt.Sprintf("pr-%d", prID),
				UserID: fmt.Sprintf("user-%d", Rnd(MaxUsers)),
			})

			t.Method = http.MethodPost
			t.URL = fmt.Sprintf("%s/pullRequest/reassign", Address)
			t.Body = body
		},
	}
)

func WriteTarget(t *vegeta.Target) error {
	generators[Rnd(len(generators))](t)
	return nil
}

func main() {
	resultsBin, err := os.Create(ResultsFile)
	if err != nil {
		log.Panicf("error creating .bin: %v", err)
	}
	defer resultsBin.Close()
	
	enc := vegeta.NewEncoder(resultsBin)
	p := plot.New()
	
	var metrics vegeta.Metrics
	fmt.Printf("starting attack...")
	
	attacker := vegeta.NewAttacker()
	
	for res := range attacker.Attack(
		WriteTarget,
		vegeta.Rate{Freq: RPS, Per: time.Second},
		Duration * time.Second,
		"load-test",
	) {
		metrics.Add(res)
		_ = p.Add(res)
		
		if err = enc.Encode(res); err != nil {
			log.Panicf("error encoding result: %v", err)
		}
	}

	metrics.Close()
	p.Close()
		
	reportTxt, err := os.Create(ReportFile)
	if err != nil {
		log.Panicf("error creating .txt: %v", err)
	}
	defer reportTxt.Close()
	
	plotHTML, err := os.Create(PlotFile)
	if err != nil {
		log.Panicf("error creating .html: %v", err)
	}
	defer plotHTML.Close()

	rep := vegeta.NewTextReporter(&metrics)
	if err = rep.Report(reportTxt); err != nil {
		log.Panicf("error writing to .txt: %v", err)
	}
	
	if _, err = p.WriteTo(plotHTML); err != nil {
		log.Panicf("error writing to .html: %v", err)
	}
}
