# SettingsApi

All URIs are relative to *http://localhost:8080/api/v1*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**settingsBrandingGet**](#settingsbrandingget) | **GET** /settings/branding | Get branding settings|
|[**settingsBrandingLogoPost**](#settingsbrandinglogopost) | **POST** /settings/branding/logo | Update app logo|
|[**settingsBrandingPut**](#settingsbrandingput) | **PUT** /settings/branding | Update branding settings|
|[**settingsPrinterGet**](#settingsprinterget) | **GET** /settings/printer | Get printer settings|
|[**settingsPrinterPut**](#settingsprinterput) | **PUT** /settings/printer | Update printer settings|

# **settingsBrandingGet**
> SettingsBrandingGet200Response settingsBrandingGet()

Retrieve branding settings (app name, logo, footer text, theme colors) for the application (Roles: authenticated)

### Example

```typescript
import {
    SettingsApi,
    Configuration
} from 'restClient';

const configuration = new Configuration();
const apiInstance = new SettingsApi(configuration);

const { status, data } = await apiInstance.settingsBrandingGet();
```

### Parameters
This endpoint does not have any parameters.


### Return type

**SettingsBrandingGet200Response**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Branding settings fetched successfully |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **settingsBrandingLogoPost**
> SettingsBrandingLogoPost200Response settingsBrandingLogoPost()

Upload and update the application logo image (Roles: admin)

### Example

```typescript
import {
    SettingsApi,
    Configuration
} from 'restClient';

const configuration = new Configuration();
const apiInstance = new SettingsApi(configuration);

let logo: File; //Logo image file (default to undefined)

const { status, data } = await apiInstance.settingsBrandingLogoPost(
    logo
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **logo** | [**File**] | Logo image file | defaults to undefined|


### Return type

**SettingsBrandingLogoPost200Response**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: multipart/form-data
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Logo updated successfully |  -  |
|**400** | Logo file is required or file size exceeds limit |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **settingsBrandingPut**
> SettingsBrandingGet200Response settingsBrandingPut(request)

Update application branding settings (Roles: admin)

### Example

```typescript
import {
    SettingsApi,
    Configuration,
    InternalSettingsUpdateBrandingRequest
} from 'restClient';

const configuration = new Configuration();
const apiInstance = new SettingsApi(configuration);

let request: InternalSettingsUpdateBrandingRequest; //Branding update request

const { status, data } = await apiInstance.settingsBrandingPut(
    request
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **request** | **InternalSettingsUpdateBrandingRequest**| Branding update request | |


### Return type

**SettingsBrandingGet200Response**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Branding settings updated successfully |  -  |
|**400** | Invalid request body or validation failure |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **settingsPrinterGet**
> SettingsPrinterGet200Response settingsPrinterGet()

Retrieve printer settings like connection string and paper width (Roles: authenticated)

### Example

```typescript
import {
    SettingsApi,
    Configuration
} from 'restClient';

const configuration = new Configuration();
const apiInstance = new SettingsApi(configuration);

const { status, data } = await apiInstance.settingsPrinterGet();
```

### Parameters
This endpoint does not have any parameters.


### Return type

**SettingsPrinterGet200Response**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Printer settings fetched successfully |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **settingsPrinterPut**
> SettingsPrinterGet200Response settingsPrinterPut(request)

Update global printer configuration (Roles: admin)

### Example

```typescript
import {
    SettingsApi,
    Configuration,
    InternalSettingsUpdatePrinterSettingsRequest
} from 'restClient';

const configuration = new Configuration();
const apiInstance = new SettingsApi(configuration);

let request: InternalSettingsUpdatePrinterSettingsRequest; //Printer settings update request

const { status, data } = await apiInstance.settingsPrinterPut(
    request
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **request** | **InternalSettingsUpdatePrinterSettingsRequest**| Printer settings update request | |


### Return type

**SettingsPrinterGet200Response**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Printer settings updated successfully |  -  |
|**400** | Invalid request body or validation failure |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

