module github.com/Ben-Ackerman/SpotifyAnalyzer/SpotifyService

go 1.13

require (
	github.com/Ben-Ackerman/SpotifyAnalyzer v0.0.0-20190906184055-9af39206995f
	golang.org/x/oauth2 v0.0.0-20190604053449-0f29369cfe45
	google.golang.org/grpc v1.23.0
)

replace github.com/Ben-Ackerman/SpotifyAnalyzer => ../
