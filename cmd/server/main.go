package main

import (
	"context"
	"flag"
	"log"
	"net"
	"net/http"
	"path"
	"path/filepath"
	"runtime"

	"push_article/internal/routes"
	"push_article/pkg/token"

	firebase "firebase.google.com/go"
	"github.com/go-chi/chi"
	"google.golang.org/api/option"
)

var _, filePath, _, _ = runtime.Caller(0)
var projectRoot, _ = filepath.Abs(path.Join(filepath.Dir(filePath), "..", ".."))

var (
	firebaseServiceAccountFlag = flag.String("firebase-service-account", "Path to file with firebase client config", "")
	listenAddrFlag             = flag.String("listen-addr", "HTTP server listen address", ":0")
)

func main() {
	flag.Parse()

	tokenStorage := token.NewMemoryStorage()
	var fbOpts []option.ClientOption
	if firebaseServiceAccountFlag != nil {
		fbOpts = append(fbOpts, option.WithCredentialsFile(*firebaseServiceAccountFlag))
	}

	fbApp, err := firebase.NewApp(context.Background(), nil, fbOpts...)
	if err != nil {
		panic(err)
	}

	fbMessaging, err := fbApp.Messaging(context.Background())
	if err != nil {
		panic(err)
	}

	mux := chi.NewMux()
	mux.Mount("/", http.FileServer(http.Dir(filepath.Join(projectRoot, "statics"))))
	mux.Mount("/api/v1/users", chi.NewMux().Group((&routes.UserService{}).AddToRouter))
	mux.Mount("/api/v1/tokens", chi.NewMux().Group((&routes.TokenService{Storage: tokenStorage}).AddToRouter))
	mux.Mount("/api/v1/notifications", chi.NewMux().Group((&routes.NotificationService{
		Storage: tokenStorage,
		Client:  fbMessaging,
	}).AddToRouter))
	mux.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, r.URL.String()+" not found", http.StatusNotFound)
	})

	listener, err := net.Listen("tcp", *listenAddrFlag)
	if err != nil {
		panic(err)
	}

	log.Printf("Listening on: http://localhost:%d", listener.Addr().(*net.TCPAddr).Port)

	panic(http.Serve(listener, mux))
}
