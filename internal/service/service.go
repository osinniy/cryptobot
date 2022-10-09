package service

import (
	d "osinniy/cryptobot/internal/data"
	"osinniy/cryptobot/internal/models"
	"osinniy/cryptobot/internal/store"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const MODULE = "service"

type Service struct {
	running         bool
	updateInterval  time.Duration
	keepAlive       time.Duration
	cleanupInterval time.Duration

	ticker        *time.Ticker
	cleanupTicker *time.Ticker

	repo   store.DataRepository
	logger zerolog.Logger
}

// Creates new service with provided interval
func New(repo store.DataRepository, interval time.Duration, keepAlive time.Duration) *Service {
	return &Service{
		updateInterval:  interval,
		keepAlive:       keepAlive,
		cleanupInterval: keepAlive / 100,
		repo:            repo,
		logger:          log.With().Str("module", MODULE).Logger(),
	}
}

// Starts service with interval.
// Corrects starting time if service was restarted
func (s *Service) Start() {
	if s.running {
		return
	}
	s.running = true

	// Cleanup job
	go func() {
		s.cleanup()

		s.cleanupTicker = time.NewTicker(s.cleanupInterval)
		for range s.cleanupTicker.C {
			s.cleanup()
		}
	}()

	// Obtain last run time
	data, _ := s.repo.Latest()

	// Check last run time and correct first run time if needed
	if data != nil {
		delay := s.updateInterval - time.Since(time.Unix(data.Upd, 0))
		if delay > 0 {
			time.Sleep(time.Duration(delay))
		}
	}

	if !s.running {
		return
	}

	// First run
	s.run()

	if !s.running {
		return
	}

	// Then run task every timer tick
	s.ticker = time.NewTicker(s.updateInterval)
	for range s.ticker.C {
		s.run()
	}
}

func (s *Service) Stop() {
	if !s.running {
		return
	}

	if s.ticker != nil {
		s.ticker.Stop()
		s.ticker = nil
	}
	if s.cleanupTicker != nil {
		s.cleanupTicker.Stop()
		s.cleanupTicker = nil
	}

	s.running = false
}

func (s Service) IsRunning() bool {
	return s.running
}

// Do single job
func (s *Service) run() {
	timestamp := time.Now()

	cmc, coinglass, err := s.update()
	if err != nil {
		s.logger.Error().Dur("elapsed_ms", time.Since(timestamp)).Msg("update failed")
		return
	}

	latestData := models.BuildData(*cmc, *coinglass)
	if err := s.repo.Save(latestData); err == nil {
		s.logger.Info().Dur("elapsed_ms", time.Since(timestamp)).Msg("data updated")
	}
}

// Triggers data endpoints and returns their responses
func (s *Service) update() (cmc *d.CMCMetricsResponse, coinglass *d.CoinglassResponse, err error) {
	cmc, err = d.LatestMarketStats()
	if err != nil {
		return
	}
	coinglass, err = d.CoinglassStats()
	if err != nil {
		return
	}
	return
}

// Toggles [store.DataRepository.Cleanup]
func (s *Service) cleanup() {
	affected, err := s.repo.Cleanup(time.Now().Add(-s.keepAlive))
	if err != nil {
		s.logger.Error().Err(err).Msg("cleanup failed")
		return
	}

	if affected > 0 {
		s.logger.Info().Int64("entries_removed", affected).Msgf("cleanup done. next in %s", s.cleanupInterval.String())
	}
}
