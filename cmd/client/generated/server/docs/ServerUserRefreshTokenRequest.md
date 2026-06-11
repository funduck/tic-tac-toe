# ServerUserRefreshTokenRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**RefreshToken** | Pointer to **string** |  | [optional] 
**UserID** | Pointer to **string** |  | [optional] 

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

### GetUserID

`func (o *ServerUserRefreshTokenRequest) GetUserID() string`

GetUserID returns the UserID field if non-nil, zero value otherwise.

### GetUserIDOk

`func (o *ServerUserRefreshTokenRequest) GetUserIDOk() (*string, bool)`

GetUserIDOk returns a tuple with the UserID field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUserID

`func (o *ServerUserRefreshTokenRequest) SetUserID(v string)`

SetUserID sets UserID field to given value.

### HasUserID

`func (o *ServerUserRefreshTokenRequest) HasUserID() bool`

HasUserID returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


