Webpush Example Project
================

This is a demo project for understanding how web pushes through FCM can be implemented.

Project contains following endpoints:
* `POST /api/v1/users/login` - authorize user and set cookie
* `POST /api/v1/users/logout` - logout user and unset cookie
* `POST /api/v1/users/me` - get current authorized user
* `POST /api/v1/tokens` - save token for authorized user
* `GET /api/v1/tokens?user_id=<user_id>` - get push-tokens for user
and static content:
`simple_example.html` - page that only gets your FCM push-token and displays it.
Webpush can be sent through Firebase console or manually.
`user_example.html` - more complex example with "authorization" using backend.
`firebase-messaging-sw.js` - simple service worker for receiving and processing push notifications.

Token invalidation occurs in following cases:
* after user logout
* when unauthorized user loads page with valid token

# How to run

## Preparation
1. Install [https://golang.org/dl/](go).

1. Create a project in [https://console.firebase.google.com](Firebase console).

1. In project console you have to register web application to receive client credentials snippet and service account credentials file.
This snippet should be inserted into `user_example.html` and `simple_example.html` (`firebaseConfig` variable).
Service account credentials file (named like `your-project-adminsdk-v1ab4-aaaaaaa.json`) should be placed in root of this project.

## Running
Run
```bash
go run ./cmd/server/main.go -firebase-service-account your-project-adminsdk-v1ab4-aaaaaaa.json
```

After successful start program displays a listen address which can be opened from browser.

Push to authorized user (`user_example.html`) can be sent using command like
```bash
curl http://<listen-addr>/api/v1/notifications -d '{"user_id": 1235, "notification":{"title": "hello", "body": "world"}}'
```
