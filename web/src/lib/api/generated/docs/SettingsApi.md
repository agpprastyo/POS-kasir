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
> POSKasirInternalDtoBrandingSettingsResponse settingsBrandingGet()

Get branding settings like app name, logo, footer text

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

**POSKasirInternalDtoBrandingSettingsResponse**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | OK |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **settingsBrandingLogoPost**
> { [key: string]: string; } settingsBrandingLogoPost()

Upload and update app logo

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

**{ [key: string]: string; }**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: multipart/form-data
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | OK |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **settingsBrandingPut**
> POSKasirInternalDtoBrandingSettingsResponse settingsBrandingPut(request)

Update branding settings

### Example

```typescript
import {
    SettingsApi,
    Configuration,
    POSKasirInternalDtoUpdateBrandingRequest
} from 'restClient';

const configuration = new Configuration();
const apiInstance = new SettingsApi(configuration);

let request: POSKasirInternalDtoUpdateBrandingRequest; //Update Branding Request

const { status, data } = await apiInstance.settingsBrandingPut(
    request
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **request** | **POSKasirInternalDtoUpdateBrandingRequest**| Update Branding Request | |


### Return type

**POSKasirInternalDtoBrandingSettingsResponse**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | OK |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **settingsPrinterGet**
> POSKasirInternalDtoPrinterSettingsResponse settingsPrinterGet()

Get printer settings like connection string, paper width

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

**POSKasirInternalDtoPrinterSettingsResponse**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | OK |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **settingsPrinterPut**
> POSKasirInternalDtoPrinterSettingsResponse settingsPrinterPut(request)

Update printer settings

### Example

```typescript
import {
    SettingsApi,
    Configuration,
    POSKasirInternalDtoUpdatePrinterSettingsRequest
} from 'restClient';

const configuration = new Configuration();
const apiInstance = new SettingsApi(configuration);

let request: POSKasirInternalDtoUpdatePrinterSettingsRequest; //Update Printer Settings Request

const { status, data } = await apiInstance.settingsPrinterPut(
    request
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **request** | **POSKasirInternalDtoUpdatePrinterSettingsRequest**| Update Printer Settings Request | |


### Return type

**POSKasirInternalDtoPrinterSettingsResponse**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | OK |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

