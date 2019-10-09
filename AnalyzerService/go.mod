module github.com/Ben-Ackerman/SpotifyAnalyzer/AnalyzerService

go 1.13

require (
	github.com/Ben-Ackerman/SpotifyAnalyzer v0.0.0
	github.com/gorilla/sessions v1.2.0
	github.com/lib/pq v1.2.0
	golang.org/x/oauth2 v0.0.0-20180821212333-d2e6202438be
	google.golang.org/grpc v1.23.0
)

replace github.com/Ben-Ackerman/SpotifyAnalyzer v0.0.0 => ../
