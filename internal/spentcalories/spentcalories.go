package spentcalories

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

const (
	lenStep                    = 0.65
	mInKm                      = 1000
	minInH                     = 60
	stepLengthCoefficient      = 0.45
	walkingCaloriesCoefficient = 0.5
)

func parseTraining(data string) (int, string, time.Duration, error) {
	parts := strings.Split(data, ",")
	if len(parts) != 3 {
		return 0, "", 0, errors.New("некорректный формат данных тренировки")
	}
	if parts[1] == "" {
		return 0, "", 0, errors.New("не указан тип тренировки")
	}

	steps, err := strconv.Atoi(parts[0])
	if err != nil || steps <= 0 {
		return 0, "", 0, errors.New("количество шагов должно быть положительным числом")
	}

	duration, err := time.ParseDuration(parts[2])
	if err != nil || duration <= 0 {
		return 0, "", 0, errors.New("продолжительность должна быть положительной")
	}

	return steps, parts[1], duration, nil
}

func distance(steps int, height float64) float64 {
	return float64(steps) * height * stepLengthCoefficient / mInKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}
	return distance(steps, height) / duration.Hours()
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, trainingType, duration, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", err
	}

	var calories float64
	switch trainingType {
	case "Бег":
		calories, err = RunningSpentCalories(steps, weight, height, duration)
	case "Ходьба":
		calories, err = WalkingSpentCalories(steps, weight, height, duration)
	default:
		return "", errors.New("неизвестный тип тренировки")
	}
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
		trainingType, duration.Hours(), distance(steps, height), meanSpeed(steps, height, duration), calories), nil
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if err := validateTrainingValues(steps, weight, height, duration); err != nil {
		return 0, err
	}

	return weight * meanSpeed(steps, height, duration) * duration.Minutes() / minInH, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	calories, err := RunningSpentCalories(steps, weight, height, duration)
	if err != nil {
		return 0, err
	}
	return calories * walkingCaloriesCoefficient, nil
}

func validateTrainingValues(steps int, weight, height float64, duration time.Duration) error {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return errors.New("некорректные параметры тренировки")
	}
	return nil
}
