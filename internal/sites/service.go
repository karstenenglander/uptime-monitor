package sites

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"
	repo "uptime-monitor/internal/adapters/postgresql/sqlc"
	tasks "uptime-monitor/internal/tasks"
	types "uptime-monitor/internal/types"

	"cloud.google.com/go/cloudtasks/apiv2/cloudtaskspb"
	"github.com/jackc/pgx/v5/pgtype"
)

type Service interface {
	ListSites(ctx context.Context) ([]repo.Site, error)
	EnqueuePollSites(ctx context.Context) ([]*cloudtaskspb.Task, error)
	PollSite(ctx context.Context, params types.PollParams) (*http.Response, error)
	AddSite(ctx context.Context, params types.CreateAddParams) (int64, error)
	RemoveSite(ctx context.Context, params types.CreateIdParams) (string, error)
	FindSitesByID(ctx context.Context, params types.CreateIdParams) (repo.Site, error)
}

type svc struct {
	repo repo.Querier
}

func NewService(repo repo.Querier) Service {
	return &svc{repo: repo}
}

func (s *svc) ListSites(ctx context.Context) ([]repo.Site, error) {
	sites, err := s.repo.ListSites(ctx)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return sites, nil
}

func (s *svc) updateSitePolled(ctx context.Context, id int64, latency int64, lastStatusCode int) (int64, error) {
	currentTime := pgtype.Timestamptz{Time: time.Now(), InfinityModifier: 0, Valid: true}
	pgLatency := pgtype.Int8{Int64: latency, Valid: true}
	pgLastStatusCode := pgtype.Int4{Int32: int32(lastStatusCode), Valid: true}
	params := repo.UpdateSitePolledParams{ID: id, PolledAt: currentTime, Latency: pgLatency, LastStatusCode: pgLastStatusCode}
	returnId, err := s.repo.UpdateSitePolled(ctx, params)
	if err != nil {
		log.Println(err)
		return returnId, err
	}
	return returnId, nil
}

func (s *svc) EnqueuePollSites(ctx context.Context) ([]*cloudtaskspb.Task, error) {
	sites, err := s.repo.ListSites(ctx)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	projectID, projectIDExists := os.LookupEnv("PROJECT_ID")
	if !projectIDExists {
		projectID = ""
	}
	locationID, locationIDExists := os.LookupEnv("LOCATION_ID")
	if !locationIDExists {
		locationID = ""
	}
	queueID, queueIDExists := os.LookupEnv("QUEUE_ID")
	if !queueIDExists {
		queueID = ""
	}
	endpointURL, endpointURLExists := os.LookupEnv("ENDPOINT_URL")
	if !endpointURLExists {
		endpointURL = ""
	}

	var enqueued []*cloudtaskspb.Task
	for _, v := range sites {
		// The message is the body of the task
		message := types.PollParams{Id: v.ID, Url: v.Url}
		task, err := tasks.CreateHTTPTaskWithToken(projectID, locationID, queueID, endpointURL, &message)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		enqueued = append(enqueued, task)
	}

	return enqueued, nil
}

func (s *svc) PollSite(ctx context.Context, params types.PollParams) (*http.Response, error) {
	url := params.Url

	timeout := 5 * time.Second
	start := time.Now()
	client := &http.Client{
		Timeout: timeout,
	}
	resp, err := client.Get(url)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	end := time.Now()
	latency := end.Sub(start)
	msLatency := latency.Milliseconds()
	// TO-DO: Add checking of return code to determine action
	updateId := params.Id
	_, err = s.updateSitePolled(ctx, updateId, msLatency, resp.StatusCode)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return resp, nil
}

func (s *svc) FindSitesByID(ctx context.Context, params types.CreateIdParams) (repo.Site, error) {
	id := params.Id
	return s.repo.FindSitesByID(ctx, id)
}

func (s *svc) AddSite(ctx context.Context, params types.CreateAddParams) (int64, error) {
	name := params.Name
	url := params.Url
	return s.repo.AddSite(ctx, repo.AddSiteParams{Name: name, Url: url})
}

func (s *svc) RemoveSite(ctx context.Context, params types.CreateIdParams) (string, error) {
	id := params.Id
	return s.repo.RemoveSiteByID(ctx, id)
}
