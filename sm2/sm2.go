package sm2

import (
	"fmt"
	"math"
	"time"
)

type SM2Result struct {
	EaseFactor  float64
	Interval    int
	Repetitions int
	NextReview  string
}

func Calculate(easeFactor float64, interval int, repetitions int, grade int) *SM2Result {
	if grade >= 3 {
		switch repetitions {
		case 0:
			interval = 1
		case 1:
			interval = 6
		default:
			interval = int(math.Round(float64(interval) * easeFactor))
		}
		repetitions++
	} else {
		repetitions = 0
		interval = 1
	}

	easeFactor = easeFactor + (0.1 - (5-float64(grade))*(0.08+(5-float64(grade))*0.02))
	if easeFactor < 1.3 {
		easeFactor = 1.3
	}

	nextReview := time.Now().AddDate(0, 0, interval).UTC().Format(time.RFC3339)

	return &SM2Result{
		EaseFactor:  easeFactor,
		Interval:    interval,
		Repetitions: repetitions,
		NextReview:  nextReview,
	}
}

func Reset() *SM2Result {
	now := time.Now().UTC().Format(time.RFC3339)
	return &SM2Result{
		EaseFactor:  2.5,
		Interval:    0,
		Repetitions: 0,
		NextReview:  now,
	}
}

func ValidateGrade(grade int) error {
	switch grade {
	case 0, 3, 4, 5:
		return nil
	default:
		return fmt.Errorf("invalid grade: %d, must be 0, 3, 4, or 5", grade)
	}
}
