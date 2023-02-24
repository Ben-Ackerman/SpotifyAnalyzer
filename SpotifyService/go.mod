module github.com/Ben-Ackerman/SpotifyAnalyzer/SpotifyService

go 1.13

require (
	github.com/go-redis/redis v6.15.5+incompatible
	github.com/onsi/ginkgo v1.16.5 // indirect
	github.com/onsi/gomega v1.27.1 // indirect
	golang.org/x/oauth2 v0.0.0-20190604053449-0f29369cfe45
	google.golang.org/appengine v1.5.0 // indirect
)

replace github.com/Ben-Ackerman/SpotifyAnalyzer v0.0.0 => ../
