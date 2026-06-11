# ServerUserSignupRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Password** | Pointer to **string** |  | [optional] 
**UserID** | Pointer to **string** |  | [optional] 

## Methods

### NewServerUserSignupRequest

`func NewServerUserSignupRequest() *ServerUserSignupRequest`

NewServerUserSignupRequest instantiates a new ServerUserSignupRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewServerUserSignupRequestWithDefaults

`func NewServerUserSignupRequestWithDefaults() *ServerUserSignupRequest`

NewServerUserSignupRequestWithDefaults instantiates a new ServerUserSignupRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetPassword

`func (o *ServerUserSignupRequest) GetPassword() string`

GetPassword returns the Password field if non-nil, zero value otherwise.

### GetPasswordOk

`func (o *ServerUserSignupRequest) GetPasswordOk() (*string, bool)`

GetPasswordOk returns a tuple with the Password field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPassword

`func (o *ServerUserSignupRequest) SetPassword(v string)`

SetPassword sets Password field to given value.

### HasPassword

`func (o *ServerUserSignupRequest) HasPassword() bool`

HasPassword returns a boolean if a field has been set.

### GetUserID

`func (o *ServerUserSignupRequest) GetUserID() string`

GetUserID returns the UserID field if non-nil, zero value otherwise.

### GetUserIDOk

`func (o *ServerUserSignupRequest) GetUserIDOk() (*string, bool)`

GetUserIDOk returns a tuple with the UserID field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUserID

`func (o *ServerUserSignupRequest) SetUserID(v string)`

SetUserID sets UserID field to given value.

### HasUserID

`func (o *ServerUserSignupRequest) HasUserID() bool`

HasUserID returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


