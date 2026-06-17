# ServerMoveRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**GameID** | Pointer to **string** |  | [optional] 
**X** | Pointer to **int32** |  | [optional] 
**Y** | Pointer to **int32** |  | [optional] 

## Methods

### NewServerMoveRequest

`func NewServerMoveRequest() *ServerMoveRequest`

NewServerMoveRequest instantiates a new ServerMoveRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewServerMoveRequestWithDefaults

`func NewServerMoveRequestWithDefaults() *ServerMoveRequest`

NewServerMoveRequestWithDefaults instantiates a new ServerMoveRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetGameID

`func (o *ServerMoveRequest) GetGameID() string`

GetGameID returns the GameID field if non-nil, zero value otherwise.

### GetGameIDOk

`func (o *ServerMoveRequest) GetGameIDOk() (*string, bool)`

GetGameIDOk returns a tuple with the GameID field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetGameID

`func (o *ServerMoveRequest) SetGameID(v string)`

SetGameID sets GameID field to given value.

### HasGameID

`func (o *ServerMoveRequest) HasGameID() bool`

HasGameID returns a boolean if a field has been set.

### GetX

`func (o *ServerMoveRequest) GetX() int32`

GetX returns the X field if non-nil, zero value otherwise.

### GetXOk

`func (o *ServerMoveRequest) GetXOk() (*int32, bool)`

GetXOk returns a tuple with the X field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetX

`func (o *ServerMoveRequest) SetX(v int32)`

SetX sets X field to given value.

### HasX

`func (o *ServerMoveRequest) HasX() bool`

HasX returns a boolean if a field has been set.

### GetY

`func (o *ServerMoveRequest) GetY() int32`

GetY returns the Y field if non-nil, zero value otherwise.

### GetYOk

`func (o *ServerMoveRequest) GetYOk() (*int32, bool)`

GetYOk returns a tuple with the Y field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetY

`func (o *ServerMoveRequest) SetY(v int32)`

SetY sets Y field to given value.

### HasY

`func (o *ServerMoveRequest) HasY() bool`

HasY returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


