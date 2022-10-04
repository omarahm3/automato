module github.com/omarahm3/automato/uploader

go 1.19

require (
	github.com/joho/godotenv v1.4.0
	golang.org/x/oauth2 v0.0.0-20220909003341-f21342109be1
	google.golang.org/api v0.98.0
)

require (
	github.com/golang/snappy v0.0.3 // indirect
	github.com/klauspost/compress v1.13.6 // indirect
	github.com/montanaflynn/stats v0.0.0-20171201202039-1bf9dbcd8cbe // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.1 // indirect
	github.com/xdg-go/stringprep v1.0.3 // indirect
	github.com/youmark/pkcs8 v0.0.0-20181117223130-1be2e3e5546d // indirect
	go.mongodb.org/mongo-driver v1.10.2 // indirect
	golang.org/x/crypto v0.0.0-20220622213112-05595931fe9d // indirect
	golang.org/x/sync v0.0.0-20220601150217-0de741cfad7f // indirect
)

require (
	cloud.google.com/go/compute v1.7.0 // indirect
	github.com/golang/groupcache v0.0.0-20200121045136-8c9f03a8e57e // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.1.0 // indirect
	github.com/googleapis/gax-go/v2 v2.4.0 // indirect
	github.com/omarahm3/automato/database v0.0.0
	github.com/omarahm3/automato/types v0.0.0
	go.opencensus.io v0.23.0 // indirect
	golang.org/x/net v0.0.0-20220909164309-bea034e7d591 // indirect
	golang.org/x/sys v0.0.0-20220728004956-3c1f35247d10 // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20220624142145-8cd45d7dbd1f // indirect
	google.golang.org/grpc v1.47.0 // indirect
	google.golang.org/protobuf v1.28.0 // indirect
)

replace github.com/omarahm3/automato/database v0.0.0 => ../database/

replace github.com/omarahm3/automato/types v0.0.0 => ../types/
