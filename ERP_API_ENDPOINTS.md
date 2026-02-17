# ERP Master Data API Endpoints

Complete API documentation for all 48 ERP master data endpoints.

## Base URL
```
http://localhost:8080/api/v1/erp
```

## Authentication
All endpoints require authentication via Bearer token (OAuth2).

## Multi-tenancy
All requests are tenant-isolated based on subdomain:
- `acme.localhost:8080` - Acme Corporation tenant
- `techcorp.localhost:8080` - Tech Corp tenant

---

## 1. Customers API

### List Customers
```http
GET /api/v1/erp/customers
```

**Query Parameters:**
- `page` (int, optional): Page number (default: 1)
- `limit` (int, optional): Items per page (default: 10, max: 100)
- `search` (string, optional): Search by name, code, or email
- `type` (string, optional): Filter by type (`individual` or `company`)
- `active` (boolean, optional): Filter by active status

**Response:**
```json
{
  "data": [
    {
      "id": "uuid",
      "code": "CUST001",
      "name": "John Doe",
      "type": "individual",
      "email": "john@example.com",
      "phone": "+1234567890",
      "credit_limit": 10000.00,
      "payment_term_days": 30,
      "currency": "USD",
      "active": true,
      "created_at": "2024-01-01T00:00:00Z"
    }
  ],
  "total": 100,
  "page": 1,
  "limit": 10,
  "total_pages": 10
}
```

### Get Customer
```http
GET /api/v1/erp/customers/:id
```

### Create Customer
```http
POST /api/v1/erp/customers
Content-Type: application/json

{
  "code": "CUST001",
  "name": "John Doe",
  "type": "individual",
  "email": "john@example.com",
  "phone": "+1234567890",
  "mobile": "+1234567891",
  "tax_id": "123456789",
  "billing_address": "123 Main St",
  "billing_city": "New York",
  "billing_state": "NY",
  "billing_zip": "10001",
  "billing_country": "USA",
  "shipping_address": "123 Main St",
  "shipping_city": "New York",
  "shipping_state": "NY",
  "shipping_zip": "10001",
  "shipping_country": "USA",
  "credit_limit": 10000.00,
  "payment_term_days": 30,
  "currency": "USD",
  "contact_person": "Jane Doe",
  "contact_person_job": "Purchasing Manager",
  "notes": "VIP customer",
  "tags": ["vip", "retail"]
}
```

### Update Customer
```http
PUT /api/v1/erp/customers/:id
Content-Type: application/json

{
  "name": "John Doe Updated",
  "email": "john.new@example.com",
  "credit_limit": 15000.00
}
```

### Delete Customer (Soft Delete)
```http
DELETE /api/v1/erp/customers/:id
```

### Get Customer Statistics
```http
GET /api/v1/erp/customers/stats
```

**Response:**
```json
{
  "total_customers": 150,
  "customers_by_type": [
    {"type": "individual", "count": 100},
    {"type": "company", "count": 50}
  ]
}
```

---

## 2. Vendors API

### List Vendors
```http
GET /api/v1/erp/vendors
```

**Query Parameters:**
- `page`, `limit`, `search`, `type`, `active` (same as Customers)
- `category` (string, optional): Filter by vendor category

### Get Vendor
```http
GET /api/v1/erp/vendors/:id
```

### Create Vendor
```http
POST /api/v1/erp/vendors
Content-Type: application/json

{
  "code": "VEND001",
  "name": "ABC Supplies Inc",
  "type": "company",
  "email": "contact@abcsupplies.com",
  "phone": "+1234567890",
  "address": "456 Vendor St",
  "city": "Los Angeles",
  "state": "CA",
  "zip": "90001",
  "country": "USA",
  "payment_term_days": 30,
  "currency": "USD",
  "bank_name": "Bank of America",
  "bank_account": "1234567890",
  "bank_account_name": "ABC Supplies Inc",
  "contact_person": "Tom Smith",
  "contact_person_job": "Sales Manager",
  "rating": 5,
  "category": "Raw Materials",
  "notes": "Reliable supplier",
  "tags": ["preferred", "local"]
}
```

### Update Vendor
```http
PUT /api/v1/erp/vendors/:id
```

### Delete Vendor
```http
DELETE /api/v1/erp/vendors/:id
```

### Get Vendor Statistics
```http
GET /api/v1/erp/vendors/stats
```

---

## 3. Products API

### List Products
```http
GET /api/v1/erp/products
```

**Query Parameters:**
- `page`, `limit`, `search`, `active` (standard)
- `type` (string, optional): `product`, `service`, or `consumable`
- `category` (string, optional): Product category
- `can_be_sold` (boolean, optional)
- `can_be_purchased` (boolean, optional)
- `track_inventory` (boolean, optional)

### Get Product
```http
GET /api/v1/erp/products/:id
```

### Create Product
```http
POST /api/v1/erp/products
Content-Type: application/json

{
  "code": "PROD001",
  "name": "Widget A",
  "type": "product",
  "category": "Electronics",
  "uom": "PCS",
  "barcode": "1234567890123",
  "sku": "WGT-A-001",
  "can_be_sold": true,
  "can_be_purchased": true,
  "track_inventory": true,
  "min_stock": 10,
  "max_stock": 100,
  "reorder_level": 20,
  "sale_price": 99.99,
  "purchase_price": 50.00,
  "cost": 45.00,
  "currency": "USD",
  "tax_category": "Standard",
  "tax_rate": 10.0,
  "weight": 1.5,
  "volume": 0.1,
  "length": 10,
  "width": 5,
  "height": 3,
  "warranty_days": 365,
  "description": "High-quality widget",
  "internal_notes": "Popular item",
  "image_url": "https://example.com/image.jpg",
  "income_account_id": "uuid",
  "expense_account_id": "uuid",
  "asset_account_id": "uuid",
  "tags": ["electronics", "popular"]
}
```

### Update Product
```http
PUT /api/v1/erp/products/:id
```

### Delete Product
```http
DELETE /api/v1/erp/products/:id
```

### Get Product Statistics
```http
GET /api/v1/erp/products/stats
```

---

## 4. Warehouses API

### List Warehouses
```http
GET /api/v1/erp/warehouses
```

**Query Parameters:**
- `page`, `limit`, `search`, `active` (standard)
- `type` (string, optional): `physical`, `virtual`, or `transit`

### Get Warehouse
```http
GET /api/v1/erp/warehouses/:id
```

### Create Warehouse
```http
POST /api/v1/erp/warehouses
Content-Type: application/json

{
  "code": "WH001",
  "name": "Main Warehouse",
  "type": "physical",
  "address": "789 Warehouse Blvd",
  "city": "Chicago",
  "state": "IL",
  "zip": "60601",
  "country": "USA",
  "phone": "+1234567890",
  "email": "warehouse@example.com",
  "manager_id": "uuid",
  "manager_name": "Mike Manager",
  "total_area": 10000.0,
  "total_capacity": 50000.0,
  "allow_negative_stock": false,
  "is_default": true,
  "description": "Primary storage facility",
  "notes": "Climate controlled"
}
```

### Update Warehouse
```http
PUT /api/v1/erp/warehouses/:id
```

### Delete Warehouse
```http
DELETE /api/v1/erp/warehouses/:id
```

### Get Warehouse Statistics
```http
GET /api/v1/erp/warehouses/stats
```

---

## 5. Chart of Accounts API

### List Chart of Accounts
```http
GET /api/v1/erp/chart-of-accounts
```

**Query Parameters:**
- `page`, `limit`, `search`, `active` (standard)
- `type` (string, optional): `asset`, `liability`, `equity`, `income`, or `expense`
- `category` (string, optional)
- `parent_id` (string, optional): Filter by parent account
- `is_group` (boolean, optional): Filter group accounts

### Get Chart of Account
```http
GET /api/v1/erp/chart-of-accounts/:id
```

### Create Chart of Account
```http
POST /api/v1/erp/chart-of-accounts
Content-Type: application/json

{
  "code": "1000",
  "name": "Assets",
  "type": "asset",
  "category": "Current Assets",
  "parent_id": null,
  "is_group": true,
  "currency": "USD",
  "allow_reconciliation": true,
  "require_tax_reporting": false,
  "description": "All company assets",
  "notes": "Root account for assets"
}
```

### Update Chart of Account
```http
PUT /api/v1/erp/chart-of-accounts/:id
```

### Delete Chart of Account
```http
DELETE /api/v1/erp/chart-of-accounts/:id
```

### Get Chart of Account Statistics
```http
GET /api/v1/erp/chart-of-accounts/stats
```

---

## 6. Taxes API

### List Taxes
```http
GET /api/v1/erp/taxes
```

**Query Parameters:**
- `page`, `limit`, `search` (standard)
- `type` (string, optional): Tax type
- `scope` (string, optional): `sales`, `purchase`, or `both`
- `is_default` (boolean, optional)

### Get Tax
```http
GET /api/v1/erp/taxes/:id
```

### Create Tax
```http
POST /api/v1/erp/taxes
Content-Type: application/json

{
  "code": "VAT10",
  "name": "VAT 10%",
  "description": "Value Added Tax 10%",
  "type": "percentage",
  "scope": "sales",
  "rate": 10.0,
  "is_default": true
}
```

### Update Tax
```http
PUT /api/v1/erp/taxes/:id
```

### Delete Tax
```http
DELETE /api/v1/erp/taxes/:id
```

### Get Tax Statistics
```http
GET /api/v1/erp/taxes/stats
```

---

## 7. Payment Terms API

### List Payment Terms
```http
GET /api/v1/erp/payment-terms
```

**Query Parameters:**
- `page`, `limit`, `search` (standard)
- `is_default` (boolean, optional)

### Get Payment Term
```http
GET /api/v1/erp/payment-terms/:id
```

### Create Payment Term
```http
POST /api/v1/erp/payment-terms
Content-Type: application/json

{
  "code": "NET30",
  "name": "Net 30",
  "description": "Payment due in 30 days",
  "days": 30,
  "discount_percent": 0.0,
  "discount_days": 0,
  "is_default": true
}
```

### Update Payment Term
```http
PUT /api/v1/erp/payment-terms/:id
```

### Delete Payment Term
```http
DELETE /api/v1/erp/payment-terms/:id
```

### Get Payment Term Statistics
```http
GET /api/v1/erp/payment-terms/stats
```

---

## 8. Currencies API

### List Currencies
```http
GET /api/v1/erp/currencies
```

**Query Parameters:**
- `page`, `limit`, `search` (standard)
- `is_default` (boolean, optional)

### Get Currency
```http
GET /api/v1/erp/currencies/:id
```

### Create Currency
```http
POST /api/v1/erp/currencies
Content-Type: application/json

{
  "code": "USD",
  "name": "US Dollar",
  "symbol": "$",
  "exchange_rate": 1.0,
  "decimal_places": 2,
  "is_default": true
}
```

### Update Currency
```http
PUT /api/v1/erp/currencies/:id
```

### Delete Currency
```http
DELETE /api/v1/erp/currencies/:id
```

### Get Currency Statistics
```http
GET /api/v1/erp/currencies/stats
```

---

## Error Responses

All endpoints follow standard HTTP status codes:

### Success (2xx)
- `200 OK` - Request successful
- `201 Created` - Resource created successfully

### Client Errors (4xx)
- `400 Bad Request` - Invalid request body or parameters
- `401 Unauthorized` - Missing or invalid authentication
- `403 Forbidden` - Insufficient permissions
- `404 Not Found` - Resource not found
- `409 Conflict` - Duplicate code or constraint violation

### Server Errors (5xx)
- `500 Internal Server Error` - Server-side error

**Error Response Format:**
```json
{
  "error": "Error message",
  "code": "ERROR_CODE",
  "details": {}
}
```

---

## Testing with cURL

### Example: Create Customer
```bash
curl -X POST http://acme.localhost:8080/api/v1/erp/customers \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "code": "CUST001",
    "name": "Test Customer",
    "type": "individual",
    "email": "test@example.com",
    "credit_limit": 5000,
    "payment_term_days": 30,
    "currency": "USD"
  }'
```

### Example: List Products with Filters
```bash
curl -X GET "http://acme.localhost:8080/api/v1/erp/products?page=1&limit=20&type=product&can_be_sold=true" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### Example: Update Vendor
```bash
curl -X PUT http://acme.localhost:8080/api/v1/erp/vendors/UUID_HERE \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "rating": 5,
    "notes": "Excellent supplier"
  }'
```

---

## Notes

1. **Multi-tenant Isolation**: All data is automatically isolated by tenant based on subdomain
2. **Soft Deletes**: DELETE operations set `active=false` instead of removing records
3. **Audit Trail**: All records track `created_by`, `updated_by`, `created_at`, `updated_at`
4. **Code Uniqueness**: Codes must be unique per tenant (not globally)
5. **Pagination**: Default limit is 10, maximum is 100 items per page
6. **Partial Updates**: PUT requests only update provided fields (null values are ignored)

---

## Postman Collection

A Postman collection with all 48 endpoints can be imported for easy testing. Create a collection with:
- Environment variables for `base_url`, `tenant`, and `auth_token`
- Pre-request scripts for token refresh
- Tests for response validation

**Environment Variables:**
```json
{
  "base_url": "http://localhost:8080",
  "tenant": "acme",
  "auth_token": "your_oauth2_token"
}
```

---

**Generated:** 2024-02-17  
**Version:** 1.0.0  
**API Version:** v1
