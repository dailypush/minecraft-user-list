package main

import (
	"bufio"
	"context"
	"log"
	"os"
	"regexp"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/go-redis/redis/v8"
)

const (
	redisHost = "redis:6379"
	redisDB   = 0
	redisKey  = "minecraft_users"
)

var logFile string
var ctx = context.Background()
var redisClient *redis.Client
var mu sync.Mutex

func Initialize(lf string) error {
	logFile = lf
	redisClient = createRedisClient()
	usernames, err := extractUsernames(logFile)
	if err != nil {
		return err
	}

	return storeUsernamesInRedis(usernames)
}

func GetUsernames() ([]string, error) {
	mu.Lock()
	defer mu.Unlock()

	usernames, err := redisClient.SMembers(ctx, redisKey).Result()
	if err != nil {
		return nil, err
	}

	return usernames, nil
}

func WatchLogFile() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalf("Error creating watcher: %v", err)
	}
	defer watcher.Close()

	done := make(chan bool)

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("Log file modified, updating Redis...")
					usernames, err := extractUsernames(logFile)
					if err != nil {
						log.Printf("Error extracting usernames: %v", err)
						continue
					}
					err = storeUsernamesInRedis(usernames)
					if err != nil {
						log.Printf("Error storing usernames in Redis: %v", err)
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Printf("Watcher error: %v", err)
			}
		}
	}()

	err = watcher.Add(logFile)
	if err != nil {
		log.Fatalf("Error adding log file to watcher: %v", err)
	}

	<-done
}

func createRedisClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: redisHost,
		DB:   redisDB,
	})

	return client
}

func extractUsernames(logFile string) ([]string, error) {
	file, err := os.Open(logFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	usernames := make(map[string]bool)
	scanner := bufio.NewScanner(file)
	pattern := regexp.MustCompile(`UUID of player (\w+) is`)

	for scanner.Scan() {
		line := scanner.Text()
		match := pattern.FindStringSubmatch(line)
		if len(match) > 1 {
			username := match[1]
			usernames[username] = true
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	var uniqueUsernames []string
	for username := range usernames {
		uniqueUsernames = append(uniqueUsernames, username)
	}

	return uniqueUsernames, nil
}

func storeUsernamesInRedis(usernames []string) error {
	mu.Lock()
	defer mu.Unlock()

	_, err := redisClient.Del(ctx, redisKey).Result()
	if err != nil {
		return err
	}

	for _, username := range usernames {
		_, err := redisClient.SAdd(ctx, redisKey, username).Result()
		if err != nil {
			return err
		}
	}

	return nil
}

func logUsernamesFromRedis() {
	usernames, err := GetUsernames()
	if err != nil {
		log.Printf("Error fetching usernames from Redis: %v", err)
		return
	}

	log.Printf("Usernames in Redis: %v", usernames)
}
