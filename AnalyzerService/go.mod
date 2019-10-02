module github.com/Ben-Ackerman/SpotifyAnalyzer/AnalyzerService

go 1.13

require (
	github.com/Ben-Ackerman/SpotifyAnalyzer v0.0.0
	github.com/gorilla/sessions v1.2.0
	github.com/lib/pq v1.2.0
	google.golang.org/grpc v1.23.0
)

replace github.com/Ben-Ackerman/SpotifyAnalyzer v0.0.0 => ../
