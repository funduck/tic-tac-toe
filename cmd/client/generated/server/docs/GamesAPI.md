# \GamesAPI

All URIs are relative to *http://localhost:8080*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CreateGame**](GamesAPI.md#CreateGame) | **Post** /api/games | Create a new game
[**GetGame**](GamesAPI.md#GetGame) | **Get** /api/games/{gameID} | Get game state
[**GiveUpGame**](GamesAPI.md#GiveUpGame) | **Post** /api/games/{gameID}/giveup | Give up the game
[**JoinAnyGame**](GamesAPI.md#JoinAnyGame) | **Post** /api/games/join | Join any available game
[**JoinGame**](GamesAPI.md#JoinGame) | **Post** /api/games/{gameID}/join | Join a waiting game
[**MakeMove**](GamesAPI.md#MakeMove) | **Post** /api/games/{gameID}/move | Make a move
[**QuitGame**](GamesAPI.md#QuitGame) | **Post** /api/games/{gameID}/quit | Quit a game that is still waiting for an opponent



## CreateGame

> GameGame CreateGame(ctx).Request(request).Execute()

Create a new game

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	request := *openapiclient.NewServerCreateGameRequest() // ServerCreateGameRequest | Create game request

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.GamesAPI.CreateGame(context.Background()).Request(request).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `GamesAPI.CreateGame``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `CreateGame`: GameGame
	fmt.Fprintf(os.Stdout, "Response from `GamesAPI.CreateGame`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiCreateGameRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **request** | [**ServerCreateGameRequest**](ServerCreateGameRequest.md) | Create game request | 

### Return type

[**GameGame**](GameGame.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetGame

> GameGame GetGame(ctx, gameID).Execute()

Get game state

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	gameID := "gameID_example" // string | Game ID

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.GamesAPI.GetGame(context.Background(), gameID).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `GamesAPI.GetGame``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetGame`: GameGame
	fmt.Fprintf(os.Stdout, "Response from `GamesAPI.GetGame`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**gameID** | **string** | Game ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetGameRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**GameGame**](GameGame.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GiveUpGame

> GameGame GiveUpGame(ctx, gameID).Request(request).Execute()

Give up the game

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	gameID := "gameID_example" // string | Game ID
	request := *openapiclient.NewServerGiveUpRequest() // ServerGiveUpRequest | Give up request

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.GamesAPI.GiveUpGame(context.Background(), gameID).Request(request).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `GamesAPI.GiveUpGame``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GiveUpGame`: GameGame
	fmt.Fprintf(os.Stdout, "Response from `GamesAPI.GiveUpGame`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**gameID** | **string** | Game ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiGiveUpGameRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **request** | [**ServerGiveUpRequest**](ServerGiveUpRequest.md) | Give up request | 

### Return type

[**GameGame**](GameGame.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## JoinAnyGame

> GameGame JoinAnyGame(ctx).Execute()

Join any available game

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.GamesAPI.JoinAnyGame(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `GamesAPI.JoinAnyGame``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `JoinAnyGame`: GameGame
	fmt.Fprintf(os.Stdout, "Response from `GamesAPI.JoinAnyGame`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiJoinAnyGameRequest struct via the builder pattern


### Return type

[**GameGame**](GameGame.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## JoinGame

> GameGame JoinGame(ctx, gameID).Request(request).Execute()

Join a waiting game

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	gameID := "gameID_example" // string | Game ID
	request := *openapiclient.NewServerJoinGameRequest() // ServerJoinGameRequest | Join game request

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.GamesAPI.JoinGame(context.Background(), gameID).Request(request).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `GamesAPI.JoinGame``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `JoinGame`: GameGame
	fmt.Fprintf(os.Stdout, "Response from `GamesAPI.JoinGame`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**gameID** | **string** | Game ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiJoinGameRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **request** | [**ServerJoinGameRequest**](ServerJoinGameRequest.md) | Join game request | 

### Return type

[**GameGame**](GameGame.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## MakeMove

> GameGame MakeMove(ctx, gameID).Request(request).Execute()

Make a move

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	gameID := "gameID_example" // string | Game ID
	request := *openapiclient.NewServerMoveRequest() // ServerMoveRequest | Move request

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.GamesAPI.MakeMove(context.Background(), gameID).Request(request).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `GamesAPI.MakeMove``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `MakeMove`: GameGame
	fmt.Fprintf(os.Stdout, "Response from `GamesAPI.MakeMove`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**gameID** | **string** | Game ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiMakeMoveRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **request** | [**ServerMoveRequest**](ServerMoveRequest.md) | Move request | 

### Return type

[**GameGame**](GameGame.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## QuitGame

> GameGame QuitGame(ctx, gameID).Request(request).Execute()

Quit a game that is still waiting for an opponent

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	gameID := "gameID_example" // string | Game ID
	request := *openapiclient.NewServerQuitRequest() // ServerQuitRequest | Quit request

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.GamesAPI.QuitGame(context.Background(), gameID).Request(request).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `GamesAPI.QuitGame``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `QuitGame`: GameGame
	fmt.Fprintf(os.Stdout, "Response from `GamesAPI.QuitGame`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**gameID** | **string** | Game ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiQuitGameRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **request** | [**ServerQuitRequest**](ServerQuitRequest.md) | Quit request | 

### Return type

[**GameGame**](GameGame.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

