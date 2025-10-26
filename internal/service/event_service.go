package service

import (
	"time"

	"github.com/mohamedelhefni/siraaj/internal/domain"
	"github.com/mohamedelhefni/siraaj/internal/repository"
)

type EventService interface {
	TrackEvent(event domain.Event) error
	TrackEventBatch(events []domain.Event) error
	GetEvents(startDate, endDate time.Time, limit, offset int) (map[string]interface{}, error)
	GetStats(startDate, endDate time.Time, limit int, filters map[string]string) (map[string]interface{}, error)
	GetOnlineUsers(timeWindow int) (map[string]interface{}, error)
	GetProjects() ([]string, error)
	GetFunnelAnalysis(request domain.FunnelRequest) (*domain.FunnelAnalysisResult, error)
}

type eventService struct {
	repo repository.EventRepository
}

func NewEventService(repo repository.EventRepository) EventService {
	return &eventService{repo: repo}
}

func (s *eventService) TrackEvent(event domain.Event) error {
	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now()
	}
	return s.repo.Create(event)
}

func (s *eventService) TrackEventBatch(events []domain.Event) error {
	now := time.Now()
	for i := range events {
		if events[i].Timestamp.IsZero() {
			events[i].Timestamp = now
		}
	}
	return s.repo.CreateBatch(events)
}

func (s *eventService) GetEvents(startDate, endDate time.Time, limit, offset int) (map[string]interface{}, error) {
	return s.repo.GetEvents(startDate, endDate, limit, offset)
}

func (s *eventService) GetStats(startDate, endDate time.Time, limit int, filters map[string]string) (map[string]interface{}, error) {
	return s.repo.GetStats(startDate, endDate, limit, filters)
}

func (s *eventService) GetOnlineUsers(timeWindow int) (map[string]interface{}, error) {
	return s.repo.GetOnlineUsers(timeWindow)
}

func (s *eventService) GetProjects() ([]string, error) {
	return s.repo.GetProjects()
}

func (s *eventService) GetFunnelAnalysis(request domain.FunnelRequest) (*domain.FunnelAnalysisResult, error) {
	return s.repo.GetFunnelAnalysis(request)
}
