# ServerUserRefreshTokenRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**RefreshToken** | Pointer to **string** |  | [optional] 
**UserId** | Pointer to **string** |  | [optional] 

## Methods

### NewServerUserRefreshTokenRequest

`func NewServerUserRefreshTokenRequest() *ServerUserRefreshTokenRequest`

NewServerUserRefreshTokenRequest instantiates a new ServerUserRefreshTokenRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewServerUserRefreshTokenRequestWithDefaults

`func NewServerUserRefreshTokenRequestWithDefaults() *ServerUserRefreshTokenRequest`

NewServerUserRefreshTokenRequestWithDefaults instantiates a new ServerUserRefreshTokenRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetRefreshToken

`func (o *ServerUserRefreshTokenRequest) GetRefreshToken() string`

GetRefreshToken returns the RefreshToken field if non-nil, zero value otherwise.

### GetRefreshTokenOk

`func (o *ServerUserRefreshTokenRequest) GetRefreshTokenOk() (*string, bool)`

GetRefreshTokenOk returns a tuple with the RefreshToken field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRefreshToken

`func (o *ServerUserRefreshTokenRequest) SetRefreshToken(v string)`

SetRefreshToken sets RefreshToken field to given value.

### HasRefreshToken

`func (o *ServerUserRefreshTokenRequest) HasRefreshToken() bool`

HasRefreshToken returns a boolean if a field has been set.

### GetUserId

`func (o *ServerUserRefreshTokenRequest) GetUserId() string`

GetUserId returns the UserId field if non-nil, zero value otherwise.

### GetUserIdOk

`func (o *ServerUserRefreshTokenRequest) GetUserIdOk() (*string, bool)`

GetUserIdOk returns a tuple with the UserId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUserId

`func (o *ServerUserRefreshTokenRequest) SetUserId(v string)`

SetUserId sets UserId field to given value.

### HasUserId

`func (o *ServerUserRefreshTokenRequest) HasUserId() bool`

HasUserId returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


