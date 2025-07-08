---
title: My Project
language_tabs:
  - shell: Shell
  - http: HTTP
  - javascript: JavaScript
  - ruby: Ruby
  - python: Python
  - php: PHP
  - java: Java
  - go: Go
toc_footers: []
includes: []
search: true
code_clipboard: true
highlight_theme: darkula
headingLevel: 2
generator: "@tarslib/widdershins v4.0.30"

---

# My Project

Base URLs:

# Authentication

- HTTP Authentication, scheme: bearer

# POS kasir/Users

## GET List All Users

GET /api/v1/users

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|page|query|integer| no |Nomor halaman untuk pagination. Default: 1.|
|limit|query|integer| no |Jumlah item per halaman. Default: 10.|
|sortBy|query|string| no |Field untuk pengurutan. Pilihan: username, email, createdAt. Default: createdAt.|
|sortOrder|query|string| no |Arah pengurutan. Pilihan: asc, desc. Default: desc.|
|search|query|string| no |Mencari pengguna berdasarkan username atau email.|
|filter[role]|query|string| no |Filter pengguna berdasarkan peran (admin, cashier, manager).|
|filter[isActive]|query|boolean| no |Filter pengguna berdasarkan status aktif (true atau false).|

#### Enum

|Name|Value|
|---|---|
|sortBy|username|
|sortBy|email|
|sortBy|createdAt|
|sortOrder|asc|
|sortOrder|desc|
|filter[role]|admin|
|filter[role]|cashier|
|filter[role]|manager|

> Response Examples

> 200 Response

```json
{
  "message": "Ante aeneus mollitia speciosus. Dolor curatio amissio colo thymum dedico aptus. Vinitor volaticus verecundia condico vacuus canto.",
  "data": [
    {
      "id": "NgrvHEfBSdmB-HzPCccmF",
      "username": "Terrence MacGyver",
      "email": "Tyrell.Conn45@hotmail.com",
      "avatar": "https://avatars.githubusercontent.com/u/71259530",
      "role": "sint pariatur consectetur nostrud",
      "is_active": false,
      "created_at": "2025-07-07T19:26:54.764Z"
    },
    {
      "id": "qK8D_OR76grX0gayoPJgy",
      "username": "Gwen Heller",
      "email": "Dillan46@hotmail.com",
      "avatar": "https://avatars.githubusercontent.com/u/67033657",
      "role": "occaecat sunt",
      "is_active": false,
      "created_at": "2025-07-07T18:35:27.792Z"
    },
    {
      "id": "rK0_st8HtaoAAZQB0Vrf-",
      "username": "Jonathon Hilll",
      "email": "Jailyn_Zemlak47@hotmail.com",
      "avatar": "https://avatars.githubusercontent.com/u/899799",
      "role": "ut deserunt sit",
      "is_active": false,
      "created_at": "2025-07-08T02:10:51.432Z"
    }
  ],
  "pagination": {
    "currentPage": 90400544,
    "totalPages": 83392510,
    "totalItems": -36613553,
    "limit": 70253072
  }
}
```

> 401 Response

```json
{
  "error": "Unauthorized",
  "message": "Authentication token is required or invalid."
}
```

> 403 Response

```json
{
  "error": "Forbidden",
  "message": "You do not have permission to access this resource."
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|none|Inline|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» message|string|true|none||none|
|» data|[object]|true|none||none|
|»» id|string|true|none||none|
|»» username|string|true|none||none|
|»» email|string|true|none||none|
|»» avatar|string¦null|true|none||none|
|»» role|string|true|none||none|
|»» is_active|boolean|true|none||none|
|»» created_at|string|true|none||none|
|» pagination|object|true|none||none|
|»» currentPage|integer|true|none||none|
|»» totalPages|integer|true|none||none|
|»» totalItems|integer|true|none||none|
|»» limit|integer|true|none||none|

HTTP Status Code **401**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» error|string|true|none||none|
|» message|string|true|none||none|

HTTP Status Code **403**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» message|string|true|none||none|
|» data|[object]|true|none||none|
|»» id|string|true|none||none|
|»» username|string|true|none||none|
|»» email|string|true|none||none|
|»» avatar|string¦null|true|none||none|
|»» role|string|true|none||none|
|»» is_active|boolean|true|none||none|
|»» created_at|string|true|none||none|
|» pagination|object|true|none||none|
|»» currentPage|integer|true|none||none|
|»» totalPages|integer|true|none||none|
|»» totalItems|integer|true|none||none|
|»» limit|integer|true|none||none|
|» error|string|true|none||none|

# Data Schema

