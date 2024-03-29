module github.com/Ben-Ackerman/SpotifyAnalyzer

go 1.13

require (
	github.com/golang/protobuf v1.3.2
	github.com/lib/pq v1.2.0
	golang.org/x/net v0.0.0-20190827160401-ba9fcec4b297
	golang.org/x/sys v0.0.0-20190904154756-749cb33beabd // indirect
	golang.org/x/text v0.3.2 // indirect
	google.golang.org/genproto v0.0.0-20190905072037-92dd089d5514 // indirect
	google.golang.org/grpc v1.23.0
)

replace (
	github.com/Ben-Ackerman/SpotifyAnalyzer/AnalyzerService => ./AnalyzerService
	github.com/Ben-Ackerman/SpotifyAnalyzer/LyricsService => ./LyricsService
	github.com/Ben-Ackerman/SpotifyAnalyzer/SpotifyService => ./SpotifyService
)
