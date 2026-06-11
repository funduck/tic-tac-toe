module github.com/funduck/tic-tac-toe/client

go 1.26.2

replace github.com/GIT_USER_ID/GIT_REPO_ID => ./generated/server

require (
	github.com/GIT_USER_ID/GIT_REPO_ID v0.0.0
	github.com/funduck/tic-tac-toe v0.0.0-20260611134052-557dae18aa36
)

require github.com/google/uuid v1.6.0 // indirect
