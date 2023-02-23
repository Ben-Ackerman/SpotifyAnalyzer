module github.com/Ben-Ackerman/SpotifyAnalyzer

go 1.13

require (
	github.com/golang/protobuf v1.3.2
	golang.org/x/text v0.3.8 // indirect
	google.golang.org/genproto v0.0.0-20190905072037-92dd089d5514 // indirect
	google.golang.org/grpc v1.23.0
)

replace (
	github.com/Ben-Ackerman/SpotifyAnalyzer/AnalyzerService => ./AnalyzerService
	github.com/Ben-Ackerman/SpotifyAnalyzer/LyricsService => ./LyricsService
	github.com/Ben-Ackerman/SpotifyAnalyzer/SpotifyService => ./SpotifyService
)
