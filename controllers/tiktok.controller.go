package controllers

import (
	"math/rand"
	"tiktok-arena/models"
	"time"
)

func findDifferenceOfTwoTiktokSlices(s1 []models.Tiktok, s2 []models.Tiktok) []models.Tiktok {
	var dif []models.Tiktok
	for _, t1 := range s1 {
		existsInS2 := false
		for _, t2 := range s2 {
			if t1.TournamentID.String() == t2.TournamentID.String() && t1.URL == t2.URL {
				existsInS2 = true
				break
			}
		}
		if !existsInS2 {
			dif = append(dif, t1)
		}
	}
	return dif
}

func containsTiktok(slice []models.Tiktok, t models.Tiktok) bool {
	for _, item := range slice {
		if item.TournamentID == t.TournamentID && item.URL == t.URL {
			return true
		}
	}
	return false
}

func shuffleTiktok(t []models.Tiktok) {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(t), func(i, j int) { t[i], t[j] = t[j], t[i] })
}
