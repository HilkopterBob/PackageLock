package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"packagelock/structs" // Replace with your actual package path

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	workerCount  = 20   // Number of concurrent workers
	batchSize    = 100 // Number of documents to insert in one batch
	totalBatches = 9999999 // Number of batches to insert
)

// Random string generator
func randomString(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

// Random boolean generator
func randomBool() bool {
	return rand.Intn(2) == 0
}

// Create synthetic data for a Package
func createPackage() structs.Package {
	return structs.Package{
		PackageID:      uuid.New(),
		PackageName:    randomString(10),
		PackageVersion: fmt.Sprintf("%d.%d.%d", rand.Intn(10), rand.Intn(10), rand.Intn(10)),
		Updatable:      randomBool(),
		CreationTime:   time.Now().AddDate(-1, 0, 0),
		UpdateTime:     time.Now(),
	}
}

// Create synthetic data for Network_Info
func createNetworkInfo() structs.Network_Info {
	return structs.Network_Info{
		Ip_addr:      fmt.Sprintf("192.168.%d.%d", rand.Intn(255), rand.Intn(255)),
		Mac_addr:     fmt.Sprintf("00:%02x:%02x:%02x:%02x:%02x", rand.Intn(256), rand.Intn(256), rand.Intn(256), rand.Intn(256), rand.Intn(256), rand.Intn(256)),
		CreationTime: time.Now().AddDate(-1, 0, 0),
		UpdateTime:   time.Now(),
	}
}

// Create synthetic data for Package_Manager
func createPackageManager() structs.Package_Manager {
	return structs.Package_Manager{
		Package_manager_name: randomString(8),
		Package_repos: []string{
			"https://repo1.example.com",
			"https://repo2.example.com",
		},
		CreationTime: time.Now().AddDate(-1, 0, 0),
		UpdateTime:   time.Now(),
	}
}

// Create synthetic data for Host
func createHost() structs.Host {
	var packages []structs.Package
	for i := 0; i < rand.Intn(10)+1; i++ {
		packages = append(packages, createPackage())
	}

	return structs.Host{
		Name:             randomString(12),
		ID:               rand.Intn(1000),
		Current_packages: packages,
		Network_info:     createNetworkInfo(),
		Package_manager:  createPackageManager(),
	}
}

// Create synthetic data for Agent
func createAgent(hostID int) structs.Agent {
	return structs.Agent{
		Agent_name:   randomString(8),
		Agent_secret: randomString(16),
		Host_ID:      hostID,
		Agent_ID:     rand.Intn(1000),
	}
}

// Create synthetic data for ApiKey
func createApiKey() structs.ApiKey {
	return structs.ApiKey{
		KeyValue:         randomString(32),
		Description:      "API Key for synthetic data",
		AccessSeperation: randomBool(),
		AccessRights:     []string{"read", "write", "update"},
		CreationTime:     time.Now().AddDate(-1, 0, 0),
		UpdateTime:       time.Now(),
	}
}

// Create synthetic data for User
func createUser() structs.User {
	var apiKeys []structs.ApiKey
	for i := 0; i < rand.Intn(5)+1; i++ {
		apiKeys = append(apiKeys, createApiKey())
	}

	return structs.User{
		UserID:       uuid.New(),
		Username:     randomString(10),
		Password:     randomString(16),
		Groups:       []string{"admin", "user"},
		CreationTime: time.Now().AddDate(-1, 0, 0),
		UpdateTime:   time.Now(),
		ApiKeys:      apiKeys,
	}
}

// Function to insert a batch of documents into MongoDB
func insertBatch(db *mongo.Database, batch []interface{}, collectionName string) error {
	_, err := db.Collection(collectionName).InsertMany(context.TODO(), batch)
	return err
}

// Worker function that generates and inserts data
func worker(db *mongo.Database, wg *sync.WaitGroup, jobs <-chan struct{}) {
	defer wg.Done()
	for range jobs {
		var users []interface{}
		var hosts []interface{}
		var agents []interface{}

		// Generate a batch of synthetic data
		for i := 0; i < batchSize; i++ {
			user := createUser()
			host := createHost()
			agent := createAgent(host.ID)

			users = append(users, user)
			hosts = append(hosts, host)
			agents = append(agents, agent)
		}

		// Insert the batches concurrently
		if err := insertBatch(db, users, "users"); err != nil {
			log.Printf("Error inserting users: %v", err)
		}
		if err := insertBatch(db, hosts, "hosts"); err != nil {
			log.Printf("Error inserting hosts: %v", err)
		}
		if err := insertBatch(db, agents, "agents"); err != nil {
			log.Printf("Error inserting agents: %v", err)
		}
	}
}

func main() {
	// Seed the random generator for unique data
	rand.Seed(time.Now().UnixNano())

	// MongoDB connection setup
	clientOptions := options.Client().ApplyURI("mongodb://username:password@172.19.0.3:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.TODO())

	// Ping the MongoDB server to check the connection
	err = client.Ping(context.TODO(), readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")

	db := client.Database("packagelock")

	// Create a worker pool
	var wg sync.WaitGroup
	jobs := make(chan struct{}, totalBatches)

	// Launch workers
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go worker(db, &wg, jobs)
	}

	// Queue up jobs
	for i := 0; i < totalBatches; i++ {
		jobs <- struct{}{}
	}
	close(jobs)

	// Wait for all workers to finish
	wg.Wait()

	fmt.Println("Successfully inserted synthetic data into MongoDB.")
}
