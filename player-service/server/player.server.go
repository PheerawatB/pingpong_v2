package server

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/exp/rand"
)

// Initialize variables
var LogMatch string
var CountMatch uint = 0
var turnOfMatch uint = 0

// Start a new match
func NewMatch() string {
	turnOfMatch = 0
	LogMatch = ""
	CountMatch++
	LogToCSV(fmt.Sprintf("------------------ New Match %d ------------------", CountMatch))
	LogToCSV("Player A 🧍🏻 & Player B 🧍🏻 on the court")
	time.Sleep(1 * time.Second)

	ach := make(chan uint)
	bch := make(chan uint)
	done := make(chan string)

	// Start the game
	go func() {
		initialPower := uint(rand.Intn(51)) + 50
		ach <- initialPower
	}()

	closeCh := false
	go func() {
		for {
			if closeCh {
				break
			}

			select {
			case power, ok := <-ach:
				turnOfMatch += 1
				if !ok {
					closeCh = true
					continue
				}
				go func() {
					time.Sleep(1 * time.Second)
					if !Player(power, bch, "Player A", "Player B") {
						done <- "Player A"
						return
					}
				}()

			case power, ok := <-bch:
				turnOfMatch += 1
				if !ok {
					closeCh = true
					continue
				}
				go func() {
					time.Sleep(1 * time.Second)
					if !Player(power, ach, "Player B", "Player A") {
						done <- "Player B"
						return
					}
				}()
			}
		}
	}()
	winner := <-done

	close(ach)
	close(bch)
	LogToCSV(fmt.Sprintf("[Alert] %s wins!", winner))
	LogToCSV("------------------- Game Over -------------------")

	return winner
}

func Player(power uint, wakeCh chan uint, name string, opponent string) bool {
	url := fmt.Sprintf("http://table-service:8889/ping-power?power=%d&name=%s", power, name)
	res, err := http.Get(url)
	if err != nil {
		fmt.Println("Error:", err)
		return false
	}
	defer res.Body.Close()
	time.Sleep(1 * time.Second)

	if res.StatusCode == http.StatusOK {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			fmt.Println("Error reading response:", err)
			return false
		}
		bodyInt, err := strconv.Atoi(string(body))
		if err != nil {
			fmt.Println("Error converting response to integer:", err)
			return false
		}

		newPower := uint(rand.Intn(51)) + 50
		time.Sleep(1 * time.Second)

		if newPower > uint(bodyInt) {
			wakeCh <- newPower
			LogToCSV(fmt.Sprintf("(t%d) [%s] 🏓 💥 {%d} ========== [%d] ==========> 🏓 [%s] ", turnOfMatch, name, power, bodyInt, opponent))
			return true
		} else {
			LogToCSV(fmt.Sprintf("(t%d) [%s] 🏓 💥 {%d} ========== [%d] ==========> 💀 [%d] [%s] ", turnOfMatch, name, power, bodyInt, newPower, opponent))
			return false
		}
	}
	return false
}

func LogToCSV(message string) {

	fmt.Println(message)
	logDir := "./logs"
	err := os.MkdirAll(logDir, os.ModePerm)
	if err != nil {
		fmt.Println("Error creating log directory:", err)
		return
	}

	date := time.Now().Format("20060102_15")
	logFileName := fmt.Sprintf("%s_%s.csv", "match", date)
	logFilePath := filepath.Join(logDir, logFileName)

	file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening/creating CSV file:", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	record := []string{
		time.Now().Format(time.RFC3339),
		message,
	}
	LogMatch += time.Now().Format(time.RFC3339) + ":" + message + "\n"

	err = writer.Write(record)
	if err != nil {
		fmt.Println("Error writing to CSV file:", err)
	}
}

func LogMatchResultToMongoDB(matchID uint, logMatch string, mongoClient *mongo.Client) {
	collection := mongoClient.Database("match_results").Collection("results")

	matchResult := MatchLog{
		MatchID:  matchID,
		MatchLog: logMatch,
		Time:     time.Now(),
	}

	_, err := collection.InsertOne(context.TODO(), matchResult)
	if err != nil {
		fmt.Printf("Error inserting match result to MongoDB: %v", err)
	} else {
		fmt.Println("Match result inserted successfully:", matchResult)
	}

}

func GetMatchID(mongoClient *mongo.Client, id uint) (MatchLog, error) {
	collection := mongoClient.Database("match_results").Collection("results")

	var match MatchLog
	err := collection.FindOne(
		context.TODO(),
		bson.D{{Key: "match_id", Value: id}},
	).Decode(&match)
	if err != nil {
		return MatchLog{}, err
	}

	return match, nil
}

func GetLastMatchID(mongoClient *mongo.Client) (uint, error) {
	collection := mongoClient.Database("match_results").Collection("results")

	var lastMatch MatchLog
	err := collection.FindOne(
		context.TODO(),
		bson.D{},
		options.FindOne().SetSort(bson.D{{Key: "match_id", Value: -1}}),
	).Decode(&lastMatch)
	if err != nil {
		return 0, err
	}

	return lastMatch.MatchID, nil
}

func GetAllMatches(mongoClient *mongo.Client) ([]MatchLog, error) {
	collection := mongoClient.Database("match_results").Collection("results")
	var matches []MatchLog

	cursor, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var match MatchLog
		if err := cursor.Decode(&match); err != nil {
			return nil, err
		}
		matches = append(matches, match)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return matches, nil
}

type MatchLog struct {
	MatchID  uint      `json:"match_id" bson:"match_id"`
	MatchLog string    `json:"match_log" bson:"match_log"`
	Time     time.Time `json:"time" bson:"time"`
}
