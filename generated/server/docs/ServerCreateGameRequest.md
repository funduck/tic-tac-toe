# ServerCreateGameRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**GameID** | Pointer to **string** |  | [optional] 
**UserID** | Pointer to **string** |  | [optional] 

## Methods

### NewServerCreateGameRequest

`func NewServerCreateGameRequest() *ServerCreateGameRequest`

NewServerCreateGameRequest instantiates a new ServerCreateGameRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewServerCreateGameRequestWithDefaults

`func NewServerCreateGameRequestWithDefaults() *ServerCreateGameRequest`

NewServerCreateGameRequestWithDefaults instantiates a new ServerCreateGameRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetGameID

`func (o *ServerCreateGameRequest) GetGameID() string`

GetGameID returns the GameID field if non-nil, zero value otherwise.

### GetGameIDOk

`func (o *ServerCreateGameRequest) GetGameIDOk() (*string, bool)`

GetGameIDOk returns a tuple with the GameID field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetGameID

`func (o *ServerCreateGameRequest) SetGameID(v string)`

SetGameID sets GameID field to given value.

### HasGameID

`func (o *ServerCreateGameRequest) HasGameID() bool`

HasGameID returns a boolean if a field has been set.

### GetUserID

`func (o *ServerCreateGameRequest) GetUserID() string`

GetUserID returns the UserID field if non-nil, zero value otherwise.

### GetUserIDOk

`func (o *ServerCreateGameRequest) GetUserIDOk() (*string, bool)`

GetUserIDOk returns a tuple with the UserID field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUserID

`func (o *ServerCreateGameRequest) SetUserID(v string)`

SetUserID sets UserID field to given value.

### HasUserID

`func (o *ServerCreateGameRequest) HasUserID() bool`

HasUserID returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


