# InternalUserCreateUserRequest


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**email** | **string** |  | [default to undefined]
**is_active** | **boolean** |  | [optional] [default to undefined]
**password** | **string** |  | [default to undefined]
**role** | [**POSKasirInternalUserRepositoryUserRole**](POSKasirInternalUserRepositoryUserRole.md) |  | [default to undefined]
**username** | **string** |  | [default to undefined]

## Example

```typescript
import { InternalUserCreateUserRequest } from 'restClient';

const instance: InternalUserCreateUserRequest = {
    email,
    is_active,
    password,
    role,
    username,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
