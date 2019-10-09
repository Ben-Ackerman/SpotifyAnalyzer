module github.com/Ben-Ackerman/SpotifyAnalyzer/SpotifyService

go 1.13

require (
	cloud.google.com/go v0.38.0
	github.com/go-redis/redis v6.15.5+incompatible
	github.com/googleapis/gax-go v2.0.2+incompatible // indirect
	github.com/gorilla/sessions v1.2.0
	golang.org/x/net v0.0.0-20190827160401-ba9fcec4b297 // indirect
	golang.org/x/oauth2 v0.0.0-20190604053449-0f29369cfe45
	google.golang.org/api v0.10.0 // indirect
	google.golang.org/genproto v0.0.0-20191002211648-c459b9ce5143 // indirect
	google.golang.org/grpc v1.24.0 // indirect
)

replace github.com/Ben-Ackerman/SpotifyAnalyzer v0.0.0 => ../
