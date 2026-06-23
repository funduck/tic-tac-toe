# ServerUserSignupRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Password** | Pointer to **string** |  | [optional] 
**UserId** | Pointer to **string** |  | [optional] 

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

### GetUserId

`func (o *ServerUserSignupRequest) GetUserId() string`

GetUserId returns the UserId field if non-nil, zero value otherwise.

### GetUserIdOk

`func (o *ServerUserSignupRequest) GetUserIdOk() (*string, bool)`

GetUserIdOk returns a tuple with the UserId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUserId

`func (o *ServerUserSignupRequest) SetUserId(v string)`

SetUserId sets UserId field to given value.

### HasUserId

`func (o *ServerUserSignupRequest) HasUserId() bool`

HasUserId returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


