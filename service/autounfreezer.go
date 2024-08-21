package service

import (
	"context"
	"fmt"
	"github.com/FlashpointProject/flashpoint-submission-system/constants"
	"github.com/FlashpointProject/flashpoint-submission-system/utils"
	"github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"
	"sync"
	"time"
)

func (s *SiteService) RunAutounfreezer(logger *logrus.Entry, ctx context.Context, wg *sync.WaitGroup, dataPacksPath, frozenPacksPath, deletedPacksPath string, doNotUnfreezeGameList []string) {
	defer wg.Done()
	l := logger.WithField("serviceName", "autounfreezer")
	defer l.Info("autounfreezer stopped")

	ticker := time.NewTicker(time.Minute) // TODO 12 hours or whatever
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			l.Info("context cancelled, stopping autounfreezer")
			return
		case <-ticker.C:

			loop := func() {
				ctx := context.WithValue(ctx, utils.CtxKeys.Log, l)

				dbs, err := s.pgdal.NewSession(ctx)
				if err != nil {
					l.Error(err)
					return
				}
				defer func() {
					err = dbs.Rollback()
					if err != nil {
						l.Error(err)
					}
				}()

				games, err := s.pgdal.GetFrozenGames(dbs)
				if err != nil {
					l.Error(err)
					return
				}

				ageThreshold := time.Now().Add(-time.Hour * 24 * 365 * 3)

				for _, game := range games {
					if slices.Contains(doNotUnfreezeGameList, game.GameID) {
						l.Infof("game %s with release date '%s' is on the do-not-unfreeze list and will be skipped", game.GameID, game.ReleaseDate)
						continue
					}

					releaseTime, err := parseDate(game.ReleaseDate)
					if err != nil {
						l.Errorf("game %s has unexpected release date format '%s' and will be skipped", game.GameID, game.ReleaseDate)
						continue
					}

					if releaseTime.Before(ageThreshold) {
						l.Infof("game %s with release date '%s' will be unfrozen", game.GameID, game.ReleaseDate)
						err = s.UnfreezeGame(ctx, game.GameID, constants.SystemID, dataPacksPath, frozenPacksPath, deletedPacksPath)
						if err != nil {
							l.Error(err)
							return
						}
					}
				}
			}

			loop()
		}
	}
}

// parseDate parses date string to up to day resolution.
// If the day is missing, then it is assumed it is the last day of the month.
// If the month is also missing, then it is assumed it is the last day of the year.
func parseDate(dateStr string) (time.Time, error) {
	var layouts = []string{
		"2006",       // yyyy
		"2006-01",    // yyyy-mm
		"2006-01-02", // yyyy-mm-dd
	}

	for _, layout := range layouts {
		if t, err := time.Parse(layout, dateStr); err == nil {
			switch layout {
			case "2006": // Move to the last day of the year
				t = time.Date(t.Year(), time.December, 31, 0, 0, 0, 0, time.UTC)
			case "2006-01": // Move to the last day of the month
				t = t.AddDate(0, 1, -1)
			}
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("invalid date format")
}
