package metrics

import (
	"encoding/csv"
	"log"
	"math"
	"os"
	"strconv"
)

type UserData struct {
	ages     []uint8
	payments []uint32
}

func AverageAge(ages []uint8) float64 {
	sum := uint64(0)
	numUsers := len(ages)

	for _, age := range ages {
		sum += uint64(age)
	}

	return float64(sum) / float64(numUsers)
}

func AveragePaymentAmount(payments []uint32) float64 {
	amount := uint64(0)
	count := len(payments)
	for _, p := range payments {
		amount += uint64(p)
	}

	return float64(amount) / float64(count) / 100.0
}

func StdDevPaymentAmount(payments []uint32) float64 {
	sum := 0.0
	sumSquares := 0.0
	for _, p := range payments {
		x := float64(p) / 100.0
		sumSquares += x * x
		sum += x
	}
	count := float64(len(payments))
	avgSquare := sumSquares / count
	avg := sum / count

	return math.Sqrt(avgSquare - avg*avg)
}

func LoadData() UserData {
	f, err := os.Open("users.csv")
	if err != nil {
		log.Fatalln("Unable to read users.csv", err)
	}
	reader := csv.NewReader(f)
	userLines, err := reader.ReadAll()
	if err != nil {
		log.Fatalln("Unable to parse users.csv as csv", err)
	}

	userData := UserData{}
	ages := make([]uint8, len(userLines))
	for i, line := range userLines {
		age, _ := strconv.Atoi(line[2])
		ages[i] = uint8(age)
	}
	userData.ages = ages

	f, err = os.Open("payments.csv")
	if err != nil {
		log.Fatalln("Unable to read payments.csv", err)
	}
	reader = csv.NewReader(f)
	paymentLines, err := reader.ReadAll()
	if err != nil {
		log.Fatalln("Unable to parse payments.csv as csv", err)
	}

	payments := make([]uint32, len(paymentLines))
	for i, line := range paymentLines {
		paymentCents, _ := strconv.Atoi(line[0])
		payments[i] = uint32(paymentCents)
	}
	userData.payments = payments

	return userData
}
