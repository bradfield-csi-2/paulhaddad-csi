package metrics

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"
	"time"
)

type UserId int
type UserMap map[UserId]*User

type Address struct {
	fullAddress string
	zip         int
}

type DollarAmount struct {
	dollars, cents uint64
}

type Payment struct {
	amount DollarAmount
	time   time.Time
}

type User struct {
	id       UserId
	name     string
	age      int
	address  Address
	payments []Payment
}

func AverageAge(userAges []uint32) float64 {
	count := 0
	var sum uint32 = 0
	for _, age := range userAges {
		count++
		sum += age
	}

	return float64(sum) / float64(count)
}

func AveragePaymentAmount(payments []DollarAmount) float64 {
	count := 0.0
	var dollars uint64 = 0
	var cents uint64 = 0
	for _, p := range payments {
		count++
		dollars += p.dollars
		cents += p.cents
	}
	totalPayments := float64(dollars) + float64(cents)/100
	return totalPayments / count
}

// Compute the standard deviation of payment amounts
// func StdDevPaymentAmount(users UserMap) float64 {
// 	mean := AveragePaymentAmount(users)
// 	squaredDiffs, count := 0.0, 0.0
// 	for _, u := range users {
// 		for _, p := range u.payments {
// 			count += 1
// 			amount := float64(p.amount.dollars) + float64(p.amount.cents)/100
// 			diff := amount - mean
// 			squaredDiffs += diff * diff
// 		}
// 	}
// 	return math.Sqrt(squaredDiffs / count)
// }

func LoadData() (UserMap, []uint32, []DollarAmount) {
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
	userAges := make([]uint32, len(userLines))
	for _, line := range userLines {
		id, _ := strconv.Atoi(line[0])
		name := line[1]
		age, _ := strconv.Atoi(line[2])
		address := line[3]
		zip, _ := strconv.Atoi(line[3])
		users[UserId(id)] = &User{UserId(id), name, age, Address{address, zip}, []Payment{}}
		userAges[id] = uint32(age)
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

	payments := make([]DollarAmount, 0)
	for _, line := range paymentLines {
		userId, _ := strconv.Atoi(line[2])
		paymentCents, _ := strconv.Atoi(line[0])
		datetime, _ := time.Parse(time.RFC3339, line[1])
		users[UserId(userId)].payments = append(users[UserId(userId)].payments, Payment{
			DollarAmount{uint64(paymentCents / 100), uint64(paymentCents % 100)},
			datetime,
		})
		payments = append(payments, DollarAmount{uint64(paymentCents / 100), uint64(paymentCents % 100)})
	}

	return users, userAges, payments
}
