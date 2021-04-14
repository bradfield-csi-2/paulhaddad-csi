package metrics

import (
	"encoding/csv"
	"log"
	"math"
	"os"
	"strconv"
	"time"
)

type UserId int
type UserMap map[UserId]*User

type Address struct {
	fullAddress string
	zip         uint16
}

type DollarAmount struct {
	dollars, cents uint32
}

type Payment struct {
	amount DollarAmount
	time   time.Time
}

type User struct {
	id       UserId
	name     string
	age      uint8
	address  Address
	payments []Payment
}

func AverageAge(ages []uint8) float64 {
	i := 0
	sum1, sum2, sum3 := 0, 0, 0
	numUsers := len(ages)

	for i < numUsers-3 {
		sum1 += int(ages[i])
		sum2 += int(ages[i+1])
		sum3 += int(ages[i+2])
		i += 3
	}

	for i < numUsers-1 {
		sum1 += int(ages[i])
		i++
	}

	return float64(sum1+sum2+sum3) / float64(i)
}

func AveragePaymentAmount(payments []uint32) float64 {
	amount := 0
	count := 0
	for _, p := range payments {
		count++
		amount += int(p)
	}

	return float64(amount) / float64(count) / 100.0
}

// Compute the standard deviation of payment amounts
func StdDevPaymentAmount(payments []uint32) float64 {
	mean := AveragePaymentAmount(payments) * 100
	squaredDiffs, count := 0.0, 0.0
	for _, p := range payments {
		count++
		amount := float64(p)
		diff := amount - mean
		squaredDiffs += diff * diff
	}
	return math.Sqrt(squaredDiffs/count) / 100.0
}

func LoadData() ([]uint8, []uint32) {
	f, err := os.Open("users.csv")
	if err != nil {
		log.Fatalln("Unable to read users.csv", err)
	}
	reader := csv.NewReader(f)
	userLines, err := reader.ReadAll()
	if err != nil {
		log.Fatalln("Unable to parse users.csv as csv", err)
	}

	users := make(UserMap, len(userLines))
	userAges := make([]uint8, len(userLines))
	for _, line := range userLines {
		id, _ := strconv.Atoi(line[0])
		name := line[1]
		age, _ := strconv.ParseUint(line[2], 10, 8)
		address := line[3]
		zip, _ := strconv.ParseUint(line[3], 10, 16)
		users[UserId(id)] = &User{UserId(id), name, uint8(age), Address{address, uint16(zip)}, []Payment{}}
		userAges[id] = uint8(age)
	}

	f, err = os.Open("payments.csv")
	if err != nil {
		log.Fatalln("Unable to read payments.csv", err)
	}
	reader = csv.NewReader(f)
	paymentLines, err := reader.ReadAll()
	if err != nil {
		log.Fatalln("Unable to parse payments.csv as csv", err)
	}

	paymentAmounts := make([]uint32, len(paymentLines))
	for i, line := range paymentLines {
		userId, _ := strconv.Atoi(line[2])
		paymentCents, _ := strconv.Atoi(line[0])
		datetime, _ := time.Parse(time.RFC3339, line[1])
		users[UserId(userId)].payments = append(users[UserId(userId)].payments, Payment{
			DollarAmount{uint32(paymentCents / 100), uint32(paymentCents % 100)},
			datetime,
		})
		paymentAmounts[i] = uint32(paymentCents)
	}

	return userAges, paymentAmounts
}
