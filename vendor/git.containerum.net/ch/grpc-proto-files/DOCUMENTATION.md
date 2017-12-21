# Protocol Documentation
<a name="top"/>

## Table of Contents

- [uuid.proto](#uuid.proto)
    - [UUID](#.UUID)
  
  
  
  

- [auth_types.proto](#auth_types.proto)
    - [AccessObject](#.AccessObject)
    - [ResourcesAccess](#.ResourcesAccess)
    - [StoredToken](#.StoredToken)
    - [StoredTokenForUser](#.StoredTokenForUser)
  
    - [AccessLevel](#.AccessLevel)
    - [Role](#.Role)
  
  
  

- [empty.proto](#empty.proto)
    - [Empty](#google.protobuf.Empty)
  
  
  
  

- [auth.proto](#auth.proto)
    - [CheckTokenRequest](#.CheckTokenRequest)
    - [CheckTokenResponse](#.CheckTokenResponse)
    - [CreateTokenRequest](#.CreateTokenRequest)
    - [CreateTokenResponse](#.CreateTokenResponse)
    - [DeleteTokenRequest](#.DeleteTokenRequest)
    - [DeleteUserTokensRequest](#.DeleteUserTokensRequest)
    - [ExtendTokenRequest](#.ExtendTokenRequest)
    - [ExtendTokenResponse](#.ExtendTokenResponse)
    - [GetUserTokensRequest](#.GetUserTokensRequest)
    - [GetUserTokensResponse](#.GetUserTokensResponse)
    - [UpdateAccessRequest](#.UpdateAccessRequest)
  
  
  
    - [Auth](#.Auth)
  

- [Scalar Value Types](#scalar-value-types)



<a name="uuid.proto"/>
<p align="right"><a href="#top">Top</a></p>

## uuid.proto



<a name=".UUID"/>

### UUID
Represents UUID in standart format


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| value | [string](#string) |  |  |





 

 

 

 



<a name="auth_types.proto"/>
<p align="right"><a href="#top">Top</a></p>

## auth_types.proto



<a name=".AccessObject"/>

### AccessObject



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| label | [string](#string) |  |  |
| id | [string](#string) |  |  |
| access | [.AccessLevel](#..AccessLevel) |  |  |






<a name=".ResourcesAccess"/>

### ResourcesAccess



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| namespace | [.AccessObject](#..AccessObject) | repeated |  |
| volume | [.AccessObject](#..AccessObject) | repeated |  |






<a name=".StoredToken"/>

### StoredToken



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| token_id | [.UUID](#..UUID) |  |  |
| user_agent | [string](#string) |  |  |
| platform | [string](#string) |  |  |
| fingerprint | [string](#string) |  |  |
| user_id | [.UUID](#..UUID) |  |  |
| user_role | [.Role](#..Role) |  |  |
| user_namespace | [string](#string) |  |  |
| user_volume | [string](#string) |  |  |
| rw_access | [bool](#bool) |  |  |
| user_ip | [string](#string) |  |  |
| part_token_id | [.UUID](#..UUID) |  |  |






<a name=".StoredTokenForUser"/>

### StoredTokenForUser



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| token_id | [.UUID](#..UUID) |  |  |
| user_agent | [string](#string) |  |  |
| ip | [string](#string) |  |  |
| created_at | [string](#string) |  |  |





 


<a name=".AccessLevel"/>

### AccessLevel


| Name | Number | Description |
| ---- | ------ | ----------- |
| OWNER | 0 |  |



<a name=".Role"/>

### Role


| Name | Number | Description |
| ---- | ------ | ----------- |
| USER | 0 |  |


 

 

 



<a name="empty.proto"/>
<p align="right"><a href="#top">Top</a></p>

## empty.proto



<a name="google.protobuf.Empty"/>

### Empty
A generic empty message that you can re-use to avoid defining duplicated
empty messages in your APIs. A typical example is to use it as the request
or the response type of an API method. For instance:

service Foo {
rpc Bar(google.protobuf.Empty) returns (google.protobuf.Empty);
}

The JSON representation for `Empty` is empty JSON object `{}`.





 

 

 

 



<a name="auth.proto"/>
<p align="right"><a href="#top">Top</a></p>

## auth.proto



<a name=".CheckTokenRequest"/>

### CheckTokenRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| access_token | [string](#string) |  |  |
| user_agent | [string](#string) |  |  |
| finger_print | [string](#string) |  |  |
| user_ip | [string](#string) |  |  |






<a name=".CheckTokenResponse"/>

### CheckTokenResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| access | [.ResourcesAccess](#..ResourcesAccess) |  |  |
| user_id | [.UUID](#..UUID) |  |  |
| user_role | [.Role](#..Role) |  |  |
| token_id | [.UUID](#..UUID) |  |  |
| part_token_id | [.UUID](#..UUID) |  |  |






<a name=".CreateTokenRequest"/>

### CreateTokenRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user_agent | [string](#string) |  |  |
| fingerprint | [string](#string) |  |  |
| user_id | [.UUID](#..UUID) |  |  |
| user_ip | [string](#string) |  |  |
| user_role | [.Role](#..Role) |  |  |
| rw_access | [bool](#bool) |  |  |
| access | [.ResourcesAccess](#..ResourcesAccess) |  |  |
| part_token_id | [.UUID](#..UUID) |  |  |






<a name=".CreateTokenResponse"/>

### CreateTokenResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| access_token | [string](#string) |  |  |
| refresh_token | [string](#string) |  |  |






<a name=".DeleteTokenRequest"/>

### DeleteTokenRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| token_id | [.UUID](#..UUID) |  |  |
| user_id | [.UUID](#..UUID) |  |  |






<a name=".DeleteUserTokensRequest"/>

### DeleteUserTokensRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user_id | [.UUID](#..UUID) |  |  |






<a name=".ExtendTokenRequest"/>

### ExtendTokenRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| refresh_token | [string](#string) |  |  |
| fingerprint | [string](#string) |  |  |






<a name=".ExtendTokenResponse"/>

### ExtendTokenResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| access_token | [string](#string) |  |  |
| refresh_token | [string](#string) |  |  |






<a name=".GetUserTokensRequest"/>

### GetUserTokensRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user_id | [.UUID](#..UUID) |  |  |






<a name=".GetUserTokensResponse"/>

### GetUserTokensResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| tokens | [.StoredTokenForUser](#..StoredTokenForUser) | repeated |  |






<a name=".UpdateAccessRequest"/>

### UpdateAccessRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user_id | [.UUID](#..UUID) |  |  |





 

 

 


<a name=".Auth"/>

### Auth
The Auth API project is an OAuth authentication server that is used to authenticate users.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| CreateToken | [CreateTokenRequest](#CreateTokenRequest) | [CreateTokenResponse](#CreateTokenRequest) |  |
| CheckToken | [CheckTokenRequest](#CheckTokenRequest) | [CheckTokenResponse](#CheckTokenRequest) |  |
| ExtendToken | [ExtendTokenRequest](#ExtendTokenRequest) | [ExtendTokenResponse](#ExtendTokenRequest) |  |
| UpdateAccess | [UpdateAccessRequest](#UpdateAccessRequest) | [google.protobuf.Empty](#UpdateAccessRequest) |  |
| GetUserTokens | [GetUserTokensRequest](#GetUserTokensRequest) | [GetUserTokensResponse](#GetUserTokensRequest) |  |
| DeleteToken | [DeleteTokenRequest](#DeleteTokenRequest) | [google.protobuf.Empty](#DeleteTokenRequest) |  |
| DeleteUserTokens | [DeleteUserTokensRequest](#DeleteUserTokensRequest) | [google.protobuf.Empty](#DeleteUserTokensRequest) |  |

 



## Scalar Value Types

| .proto Type | Notes | C++ Type | Java Type | Python Type |
| ----------- | ----- | -------- | --------- | ----------- |
| <a name="double" /> double |  | double | double | float |
| <a name="float" /> float |  | float | float | float |
| <a name="int32" /> int32 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int | int |
| <a name="int64" /> int64 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | long | int/long |
| <a name="uint32" /> uint32 | Uses variable-length encoding. | uint32 | int | int/long |
| <a name="uint64" /> uint64 | Uses variable-length encoding. | uint64 | long | int/long |
| <a name="sint32" /> sint32 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int | int |
| <a name="sint64" /> sint64 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | long | int/long |
| <a name="fixed32" /> fixed32 | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | int | int |
| <a name="fixed64" /> fixed64 | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | long | int/long |
| <a name="sfixed32" /> sfixed32 | Always four bytes. | int32 | int | int |
| <a name="sfixed64" /> sfixed64 | Always eight bytes. | int64 | long | int/long |
| <a name="bool" /> bool |  | bool | boolean | boolean |
| <a name="string" /> string | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | String | str/unicode |
| <a name="bytes" /> bytes | May contain any arbitrary sequence of bytes. | string | ByteString | str |

