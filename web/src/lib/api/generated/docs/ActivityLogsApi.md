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
let entityType: string; //Entity Type (optional) (default to undefined)
let actionType: string; //Action Type (optional) (default to undefined)

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
| **entityType** | [**string**] | Entity Type | (optional) defaults to undefined|
| **actionType** | [**string**] | Action Type | (optional) defaults to undefined|


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

