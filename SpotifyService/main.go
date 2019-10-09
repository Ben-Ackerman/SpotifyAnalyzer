package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/Ben-Ackerman/SpotifyAnalyzer/SpotifyService/app"
	"github.com/Ben-Ackerman/SpotifyAnalyzer/SpotifyService/spotifyapi"
	"github.com/go-redis/redis"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func run() error {
	clientID := os.Getenv("SpotifyClientID")
	clientSecret := os.Getenv("SpotifyClientSecret")
	redirectURL := os.Getenv("SpotifyRedirectURL")
	redisPort := os.Getenv("REDIS_PORT")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	s := &app.Server{
		Router:      http.DefaultServeMux,
		SpotifyAuth: spotifyapi.NewAuthenticator(redirectURL, clientID, clientSecret, spotifyapi.ScopeUserTopRead),
		RedisClient: redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("redis:%s", redisPort),
			Password: redisPassword,
			DB:       0, // use default DB
		}),
	}
	s.Init()
	pong, err := s.RedisClient.Ping().Result()
	fmt.Println(pong, err)

	servicePort := os.Getenv("SpotifyServicePort")
	err = http.ListenAndServe(fmt.Sprintf(":%s", servicePort), s)
	return err
}
