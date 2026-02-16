# ActivityLogsApi

All URIs are relative to *http://localhost:8080/api/v1*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**activityLogsGet**](#activitylogsget) | **GET** /activity-logs | Get activity logs|

# **activityLogsGet**
> ActivityLogsGet200Response activityLogsGet()

Get a list of activity logs with filtering and pagination (Roles: admin)

### Example

```typescript
import {
    ActivityLogsApi,
    Configuration
} from 'restClient';

const configuration = new Configuration();
const apiInstance = new ActivityLogsApi(configuration);

let page: number; //Page number (optional) (default to 1)
let limit: number; //Items per page (optional) (default to 10)
let search: string; //Search term (optional) (default to undefined)
let startDate: string; //Start date (YYYY-MM-DD) (optional) (default to undefined)
let endDate: string; //End date (YYYY-MM-DD) (optional) (default to undefined)
let userId: string; //User ID (optional) (default to undefined)
let entityType: 'PRODUCT' | 'CATEGORY' | 'PROMOTION' | 'ORDER' | 'USER'; //Entity Type (optional) (default to undefined)
let actionType: 'CREATE' | 'UPDATE' | 'DELETE' | 'CANCEL' | 'APPLY_PROMOTION' | 'PROCESS_PAYMENT' | 'REGISTER' | 'UPDATE_PASSWORD' | 'UPDATE_AVATAR' | 'LOGIN_SUCCESS' | 'LOGIN_FAILED'; //Action Type (optional) (default to undefined)

const { status, data } = await apiInstance.activityLogsGet(
    page,
    limit,
    search,
    startDate,
    endDate,
    userId,
    entityType,
    actionType
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **page** | [**number**] | Page number | (optional) defaults to 1|
| **limit** | [**number**] | Items per page | (optional) defaults to 10|
| **search** | [**string**] | Search term | (optional) defaults to undefined|
| **startDate** | [**string**] | Start date (YYYY-MM-DD) | (optional) defaults to undefined|
| **endDate** | [**string**] | End date (YYYY-MM-DD) | (optional) defaults to undefined|
| **userId** | [**string**] | User ID | (optional) defaults to undefined|
| **entityType** | [**&#39;PRODUCT&#39; | &#39;CATEGORY&#39; | &#39;PROMOTION&#39; | &#39;ORDER&#39; | &#39;USER&#39;**]**Array<&#39;PRODUCT&#39; &#124; &#39;CATEGORY&#39; &#124; &#39;PROMOTION&#39; &#124; &#39;ORDER&#39; &#124; &#39;USER&#39;>** | Entity Type | (optional) defaults to undefined|
| **actionType** | [**&#39;CREATE&#39; | &#39;UPDATE&#39; | &#39;DELETE&#39; | &#39;CANCEL&#39; | &#39;APPLY_PROMOTION&#39; | &#39;PROCESS_PAYMENT&#39; | &#39;REGISTER&#39; | &#39;UPDATE_PASSWORD&#39; | &#39;UPDATE_AVATAR&#39; | &#39;LOGIN_SUCCESS&#39; | &#39;LOGIN_FAILED&#39;**]**Array<&#39;CREATE&#39; &#124; &#39;UPDATE&#39; &#124; &#39;DELETE&#39; &#124; &#39;CANCEL&#39; &#124; &#39;APPLY_PROMOTION&#39; &#124; &#39;PROCESS_PAYMENT&#39; &#124; &#39;REGISTER&#39; &#124; &#39;UPDATE_PASSWORD&#39; &#124; &#39;UPDATE_AVATAR&#39; &#124; &#39;LOGIN_SUCCESS&#39; &#124; &#39;LOGIN_FAILED&#39;>** | Action Type | (optional) defaults to undefined|


### Return type

**ActivityLogsGet200Response**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | OK |  -  |
|**400** | Bad Request |  -  |
|**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

