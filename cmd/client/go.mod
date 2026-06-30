module github.com/funduck/tic-tac-toe/client

go 1.26.2

replace github.com/GIT_USER_ID/GIT_REPO_ID => ./generated/server

require (
	github.com/GIT_USER_ID/GIT_REPO_ID v0.0.0
	golang.org/x/term v0.44.0
)

require golang.org/x/sys v0.46.0 // indirect
