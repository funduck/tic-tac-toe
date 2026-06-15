# GameGame

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Board** | Pointer to **[][]int32** |  | [optional] 
**CurrentPlayerID** | Pointer to **string** | user ID of the player whose turn it is | [optional] 
**Id** | Pointer to **string** |  | [optional] 
**Private** | Pointer to **bool** | optional field to indicate if the game is private or public | [optional] 
**Result** | Pointer to [**GameGameResult**](GameGameResult.md) |  | [optional] 
**Status** | Pointer to [**GameGameStatus**](GameGameStatus.md) |  | [optional] 
**UserID1** | Pointer to **string** |  | [optional] 
**UserID2** | Pointer to **string** |  | [optional] 
**WinnerID** | Pointer to **string** |  | [optional] 

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

### GetCurrentPlayerID

`func (o *GameGame) GetCurrentPlayerID() string`

GetCurrentPlayerID returns the CurrentPlayerID field if non-nil, zero value otherwise.

### GetCurrentPlayerIDOk

`func (o *GameGame) GetCurrentPlayerIDOk() (*string, bool)`

GetCurrentPlayerIDOk returns a tuple with the CurrentPlayerID field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCurrentPlayerID

`func (o *GameGame) SetCurrentPlayerID(v string)`

SetCurrentPlayerID sets CurrentPlayerID field to given value.

### HasCurrentPlayerID

`func (o *GameGame) HasCurrentPlayerID() bool`

HasCurrentPlayerID returns a boolean if a field has been set.

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

### GetUserID1

`func (o *GameGame) GetUserID1() string`

GetUserID1 returns the UserID1 field if non-nil, zero value otherwise.

### GetUserID1Ok

`func (o *GameGame) GetUserID1Ok() (*string, bool)`

GetUserID1Ok returns a tuple with the UserID1 field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUserID1

`func (o *GameGame) SetUserID1(v string)`

SetUserID1 sets UserID1 field to given value.

### HasUserID1

`func (o *GameGame) HasUserID1() bool`

HasUserID1 returns a boolean if a field has been set.

### GetUserID2

`func (o *GameGame) GetUserID2() string`

GetUserID2 returns the UserID2 field if non-nil, zero value otherwise.

### GetUserID2Ok

`func (o *GameGame) GetUserID2Ok() (*string, bool)`

GetUserID2Ok returns a tuple with the UserID2 field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUserID2

`func (o *GameGame) SetUserID2(v string)`

SetUserID2 sets UserID2 field to given value.

### HasUserID2

`func (o *GameGame) HasUserID2() bool`

HasUserID2 returns a boolean if a field has been set.

### GetWinnerID

`func (o *GameGame) GetWinnerID() string`

GetWinnerID returns the WinnerID field if non-nil, zero value otherwise.

### GetWinnerIDOk

`func (o *GameGame) GetWinnerIDOk() (*string, bool)`

GetWinnerIDOk returns a tuple with the WinnerID field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWinnerID

`func (o *GameGame) SetWinnerID(v string)`

SetWinnerID sets WinnerID field to given value.

### HasWinnerID

`func (o *GameGame) HasWinnerID() bool`

HasWinnerID returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


