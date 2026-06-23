# GameGame

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Board** | Pointer to **[][]int32** |  | [optional] 
**CurrentPlayerId** | Pointer to **string** | user ID of the player whose turn it is | [optional] 
**Id** | Pointer to **string** |  | [optional] 
**Private** | Pointer to **bool** | optional field to indicate if the game is private or public | [optional] 
**Result** | Pointer to [**GameGameResult**](GameGameResult.md) |  | [optional] 
**Status** | Pointer to [**GameGameStatus**](GameGameStatus.md) |  | [optional] 
**UserId1** | Pointer to **string** |  | [optional] 
**UserId2** | Pointer to **string** |  | [optional] 
**WinnerId** | Pointer to **string** |  | [optional] 

## Methods

### NewGameGame

`func NewGameGame() *GameGame`

NewGameGame instantiates a new GameGame object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewGameGameWithDefaults

`func NewGameGameWithDefaults() *GameGame`

NewGameGameWithDefaults instantiates a new GameGame object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetBoard

`func (o *GameGame) GetBoard() [][]int32`

GetBoard returns the Board field if non-nil, zero value otherwise.

### GetBoardOk

`func (o *GameGame) GetBoardOk() (*[][]int32, bool)`

GetBoardOk returns a tuple with the Board field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBoard

`func (o *GameGame) SetBoard(v [][]int32)`

SetBoard sets Board field to given value.

### HasBoard

`func (o *GameGame) HasBoard() bool`

HasBoard returns a boolean if a field has been set.

### GetCurrentPlayerId

`func (o *GameGame) GetCurrentPlayerId() string`

GetCurrentPlayerId returns the CurrentPlayerId field if non-nil, zero value otherwise.

### GetCurrentPlayerIdOk

`func (o *GameGame) GetCurrentPlayerIdOk() (*string, bool)`

GetCurrentPlayerIdOk returns a tuple with the CurrentPlayerId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCurrentPlayerId

`func (o *GameGame) SetCurrentPlayerId(v string)`

SetCurrentPlayerId sets CurrentPlayerId field to given value.

### HasCurrentPlayerId

`func (o *GameGame) HasCurrentPlayerId() bool`

HasCurrentPlayerId returns a boolean if a field has been set.

### GetId

`func (o *GameGame) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *GameGame) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *GameGame) SetId(v string)`

SetId sets Id field to given value.

### HasId

`func (o *GameGame) HasId() bool`

HasId returns a boolean if a field has been set.

### GetPrivate

`func (o *GameGame) GetPrivate() bool`

GetPrivate returns the Private field if non-nil, zero value otherwise.

### GetPrivateOk

`func (o *GameGame) GetPrivateOk() (*bool, bool)`

GetPrivateOk returns a tuple with the Private field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPrivate

`func (o *GameGame) SetPrivate(v bool)`

SetPrivate sets Private field to given value.

### HasPrivate

`func (o *GameGame) HasPrivate() bool`

HasPrivate returns a boolean if a field has been set.

### GetResult

`func (o *GameGame) GetResult() GameGameResult`

GetResult returns the Result field if non-nil, zero value otherwise.

### GetResultOk

`func (o *GameGame) GetResultOk() (*GameGameResult, bool)`

GetResultOk returns a tuple with the Result field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetResult

`func (o *GameGame) SetResult(v GameGameResult)`

SetResult sets Result field to given value.

### HasResult

`func (o *GameGame) HasResult() bool`

HasResult returns a boolean if a field has been set.

### GetStatus

`func (o *GameGame) GetStatus() GameGameStatus`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *GameGame) GetStatusOk() (*GameGameStatus, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *GameGame) SetStatus(v GameGameStatus)`

SetStatus sets Status field to given value.

### HasStatus

`func (o *GameGame) HasStatus() bool`

HasStatus returns a boolean if a field has been set.

### GetUserId1

`func (o *GameGame) GetUserId1() string`

GetUserId1 returns the UserId1 field if non-nil, zero value otherwise.

### GetUserId1Ok

`func (o *GameGame) GetUserId1Ok() (*string, bool)`

GetUserId1Ok returns a tuple with the UserId1 field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUserId1

`func (o *GameGame) SetUserId1(v string)`

SetUserId1 sets UserId1 field to given value.

### HasUserId1

`func (o *GameGame) HasUserId1() bool`

HasUserId1 returns a boolean if a field has been set.

### GetUserId2

`func (o *GameGame) GetUserId2() string`

GetUserId2 returns the UserId2 field if non-nil, zero value otherwise.

### GetUserId2Ok

`func (o *GameGame) GetUserId2Ok() (*string, bool)`

GetUserId2Ok returns a tuple with the UserId2 field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUserId2

`func (o *GameGame) SetUserId2(v string)`

SetUserId2 sets UserId2 field to given value.

### HasUserId2

`func (o *GameGame) HasUserId2() bool`

HasUserId2 returns a boolean if a field has been set.

### GetWinnerId

`func (o *GameGame) GetWinnerId() string`

GetWinnerId returns the WinnerId field if non-nil, zero value otherwise.

### GetWinnerIdOk

`func (o *GameGame) GetWinnerIdOk() (*string, bool)`

GetWinnerIdOk returns a tuple with the WinnerId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWinnerId

`func (o *GameGame) SetWinnerId(v string)`

SetWinnerId sets WinnerId field to given value.

### HasWinnerId

`func (o *GameGame) HasWinnerId() bool`

HasWinnerId returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


