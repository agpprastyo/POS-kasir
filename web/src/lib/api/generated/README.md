## restClient@1.0

This generator creates TypeScript/JavaScript client that utilizes [axios](https://github.com/axios/axios). The generated Node module can be used in the following environments:

Environment
* Node.js
* Webpack
* Browserify

Language level
* ES5 - you must have a Promises/A+ library installed
* ES6

Module system
* CommonJS
* ES6 module system

It can be used in both TypeScript and JavaScript. In TypeScript, the definition will be automatically resolved via `package.json`. ([Reference](https://www.typescriptlang.org/docs/handbook/declaration-files/consumption.html))

### Building

To build and compile the typescript sources to javascript use:
```
npm install
npm run build
```

### Publishing

First build the package then run `npm publish`

### Consuming

navigate to the folder of your consuming project and run one of the following commands.

_published:_

```
npm install restClient@1.0 --save
```

_unPublished (not recommended):_

```
npm install PATH_TO_GENERATED_PACKAGE --save
```

### Documentation for API Endpoints

All URIs are relative to *http://localhost:8080/api/v1*

Class | Method | HTTP request | Description
------------ | ------------- | ------------- | -------------
*AuthApi* | [**authLoginPost**](docs/AuthApi.md#authloginpost) | **POST** /auth/login | Login
*AuthApi* | [**authLogoutPost**](docs/AuthApi.md#authlogoutpost) | **POST** /auth/logout | Logout
*AuthApi* | [**authMeAvatarPut**](docs/AuthApi.md#authmeavatarput) | **PUT** /auth/me/avatar | Update avatar
*AuthApi* | [**authMeGet**](docs/AuthApi.md#authmeget) | **GET** /auth/me | Get profile
*AuthApi* | [**authRegisterPost**](docs/AuthApi.md#authregisterpost) | **POST** /auth/register | Register
*AuthApi* | [**authUpdatePasswordPost**](docs/AuthApi.md#authupdatepasswordpost) | **POST** /auth/update-password | Update password
*CancellationReasonsApi* | [**apiV1CancellationReasonsGet**](docs/CancellationReasonsApi.md#apiv1cancellationreasonsget) | **GET** /api/v1/cancellation-reasons | List cancellation reasons
*UsersApi* | [**usersGet**](docs/UsersApi.md#usersget) | **GET** /users | Get all users
*UsersApi* | [**usersIdDelete**](docs/UsersApi.md#usersiddelete) | **DELETE** /users/{id} | Delete user
*UsersApi* | [**usersIdGet**](docs/UsersApi.md#usersidget) | **GET** /users/{id} | Get user by ID
*UsersApi* | [**usersIdPut**](docs/UsersApi.md#usersidput) | **PUT** /users/{id} | Update user
*UsersApi* | [**usersIdTogglePut**](docs/UsersApi.md#usersidtoggleput) | **PUT** /users/{id}/toggle | Toggle user status
*UsersApi* | [**usersPost**](docs/UsersApi.md#userspost) | **POST** /users | Create user


### Documentation For Models

 - [ApiV1CancellationReasonsGet200Response](docs/ApiV1CancellationReasonsGet200Response.md)
 - [AuthLoginPost200Response](docs/AuthLoginPost200Response.md)
 - [AuthLogoutPost200Response](docs/AuthLogoutPost200Response.md)
 - [AuthMeGet200Response](docs/AuthMeGet200Response.md)
 - [POSKasirInternalCommonErrorResponse](docs/POSKasirInternalCommonErrorResponse.md)
 - [POSKasirInternalCommonSuccessResponse](docs/POSKasirInternalCommonSuccessResponse.md)
 - [POSKasirInternalDtoCancellationReasonResponse](docs/POSKasirInternalDtoCancellationReasonResponse.md)
 - [POSKasirInternalDtoCreateUserRequest](docs/POSKasirInternalDtoCreateUserRequest.md)
 - [POSKasirInternalDtoLoginRequest](docs/POSKasirInternalDtoLoginRequest.md)
 - [POSKasirInternalDtoLoginResponse](docs/POSKasirInternalDtoLoginResponse.md)
 - [POSKasirInternalDtoProfileResponse](docs/POSKasirInternalDtoProfileResponse.md)
 - [POSKasirInternalDtoRegisterRequest](docs/POSKasirInternalDtoRegisterRequest.md)
 - [POSKasirInternalDtoUpdatePasswordRequest](docs/POSKasirInternalDtoUpdatePasswordRequest.md)
 - [POSKasirInternalDtoUpdateUserRequest](docs/POSKasirInternalDtoUpdateUserRequest.md)
 - [POSKasirInternalDtoUsersResponse](docs/POSKasirInternalDtoUsersResponse.md)
 - [POSKasirInternalRepositoryUserRole](docs/POSKasirInternalRepositoryUserRole.md)
 - [POSKasirPkgPaginationPagination](docs/POSKasirPkgPaginationPagination.md)
 - [UsersGet200Response](docs/UsersGet200Response.md)


<a id="documentation-for-authorization"></a>
## Documentation For Authorization

Endpoints do not require authorization.

