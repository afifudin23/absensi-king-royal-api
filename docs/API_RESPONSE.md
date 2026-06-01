# Standard API Response

## Format Umum

Semua response menggunakan struktur yang sama:

```json
{
  "success": true | false,
  "data": { ... } | [ ... ] | null,
  "error": null | { "code": "...", "message": "...", "details": ... }
}
```

---

## Success Responses

### GET — Single Object
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "field": "value"
  },
  "error": null
}
```

### GET — List
```json
{
  "success": true,
  "data": [
    { "id": "uuid", "field": "value" },
    { "id": "uuid", "field": "value" }
  ],
  "error": null
}
```
> Tidak ada pagination — semua list return array penuh.

### POST / PUT / PATCH — Return ID
```json
{
  "success": true,
  "data": { "id": "uuid" },
  "error": null
}
```

### DELETE
```json
{
  "success": true,
  "data": { "deleted_count": 2 },
  "error": null
}
```

---

## Error Responses

### Validation Error (400)
```json
{
  "success": false,
  "data": null,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Validation failed",
    "details": {
      "field": "error message"
    }
  }
}
```

### Bad Request (400)
```json
{
  "success": false,
  "data": null,
  "error": {
    "code": "BAD_REQUEST",
    "message": "...",
    "details": null
  }
}
```

### Unauthorized (401)
```json
{
  "success": false,
  "data": null,
  "error": {
    "code": "UNAUTHORIZED",
    "message": "...",
    "details": null
  }
}
```

### Forbidden (403)
```json
{
  "success": false,
  "data": null,
  "error": {
    "code": "FORBIDDEN",
    "message": "...",
    "details": null
  }
}
```

### Not Found (404)
```json
{
  "success": false,
  "data": null,
  "error": {
    "code": "NOT_FOUND",
    "message": "...",
    "details": null
  }
}
```

### Internal Server Error (500)
```json
{
  "success": false,
  "data": null,
  "error": {
    "code": "INTERNAL_SERVER_ERROR",
    "message": "Internal server error",
    "details": null
  }
}
```

---

## Auth Header

Semua endpoint kecuali `/auth/login`, `/auth/register`, `/auth/forgot-password`, `/auth/reset-password` wajib menyertakan:

```
Authorization: Bearer <access_token>
```
