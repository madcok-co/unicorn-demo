# Sales Module API Endpoints

Complete API documentation for all 24 Sales module endpoints implementing quote-to-cash workflow.

## Base URL
```
http://localhost:8080/api/v1/sales
```

## Authentication
All endpoints require authentication via Bearer token (OAuth2).

## Multi-tenancy
All requests are tenant-isolated based on subdomain:
- `acme.localhost:8080` - Acme Corporation tenant
- `techcorp.localhost:8080` - Tech Corp tenant

---

## Sales Workflow

```
QUOTATION → SALES ORDER → DELIVERY ORDER → SALES INVOICE → PAYMENT
```

---

## 1. Sales Quotations API

### List Sales Quotations
```http
GET /api/v1/sales/quotations
```

**Query Parameters:**
- `page` (int, optional): Page number (default: 1)
- `limit` (int, optional): Items per page (default: 10, max: 100)
- `search` (string, optional): Search by code, customer name, or reference
- `customer_id` (string, optional): Filter by customer
- `status` (string, optional): draft, sent, accepted, rejected, expired, cancelled
- `salesperson_id` (string, optional): Filter by salesperson
- `date_from` (date, optional): Quotation date from
- `date_to` (date, optional): Quotation date to
- `active` (boolean, optional): Filter by active status

**Response:**
```json
{
  "data": [
    {
      "id": "uuid",
      "code": "SQ-2024-001",
      "customer_id": "uuid",
      "customer_name": "Acme Corp",
      "quotation_date": "2024-01-01T00:00:00Z",
      "valid_until": "2024-01-31T00:00:00Z",
      "expected_date": "2024-01-15T00:00:00Z",
      "status": "accepted",
      "currency": "USD",
      "exchange_rate": 1.0,
      "subtotal": 799.95,
      "total_discount": 49.995,
      "total_tax": 75.00,
      "shipping_cost": 50.00,
      "other_cost": 0.00,
      "grandtotal": 875.00,
      "payment_term_days": 30,
      "items": [
        {
          "line_number": 1,
          "product_code": "PROD001",
          "product_name": "Widget A",
          "quantity": 5.0,
          "uom": "PCS",
          "unit_price": 99.99,
          "discount_type": "percentage",
          "discount": 10.0,
          "tax_rate": 10.0,
          "subtotal": 499.95,
          "total": 495.00
        }
      ],
      "created_at": "2024-01-01T00:00:00Z"
    }
  ],
  "total": 10,
  "page": 1,
  "limit": 10,
  "total_pages": 1
}
```

### Get Sales Quotation
```http
GET /api/v1/sales/quotations/:id
```

### Create Sales Quotation
```http
POST /api/v1/sales/quotations
Content-Type: application/json

{
  "code": "SQ-2024-001",
  "customer_id": "uuid",
  "quotation_date": "2024-01-01T00:00:00Z",
  "valid_until": "2024-01-31T00:00:00Z",
  "expected_date": "2024-01-15T00:00:00Z",
  "reference_no": "REF-Q-001",
  "salesperson_id": "uuid",
  "currency": "USD",
  "exchange_rate": 1.0,
  "payment_term_days": 30,
  "shipping_address": "123 Main St",
  "shipping_city": "New York",
  "shipping_state": "NY",
  "shipping_zip": "10001",
  "shipping_country": "USA",
  "shipping_cost": 50.00,
  "other_cost": 0.00,
  "notes": "Sample quotation",
  "internal_notes": "Customer wants discount",
  "terms": "Payment terms: Net 30",
  "tags": ["demo", "vip"],
  "items": [
    {
      "product_id": "uuid",
      "description": "High quality product",
      "quantity": 5.0,
      "unit_price": 99.99,
      "discount_type": "percentage",
      "discount": 10.0,
      "tax_id": "uuid",
      "notes": ""
    }
  ]
}
```

**Business Rules:**
- Items array must have at least 1 item
- Customer must exist and be active
- Product must exist and be active
- Totals are auto-calculated from items
- Status defaults to "draft"

### Update Sales Quotation
```http
PUT /api/v1/sales/quotations/:id
Content-Type: application/json

{
  "valid_until": "2024-02-28T00:00:00Z",
  "status": "sent",
  "shipping_cost": 75.00,
  "notes": "Updated notes"
}
```

**Business Rules:**
- Only draft quotations can be freely updated
- Status changes update workflow timestamps:
  - sent → sent_at
  - accepted → accepted_at
  - rejected → rejected_at
  - expired → expired_at
- Changing costs recalculates grand_total

### Delete Sales Quotation
```http
DELETE /api/v1/sales/quotations/:id
```

**Business Rules:**
- Only draft or rejected quotations can be deleted
- Soft delete (sets active=false)
- Cascade deletes quotation items

### Get Quotation Statistics
```http
GET /api/v1/sales/quotations/stats
```

**Response:**
```json
{
  "total_quotations": 50,
  "total_value": 125000.00,
  "accepted_value": 75000.00,
  "quotations_by_status": [
    {"status": "draft", "count": 10, "total": 25000.00},
    {"status": "sent", "count": 15, "total": 37500.00},
    {"status": "accepted", "count": 20, "total": 50000.00},
    {"status": "rejected", "count": 5, "total": 12500.00}
  ]
}
```

---

## 2. Sales Orders API

### List Sales Orders
```http
GET /api/v1/sales/orders
```

**Query Parameters:**
- `page`, `limit`, `search`, `customer_id`, `salesperson_id`, `active` (same as quotations)
- `status` (string, optional): draft, confirmed, processing, completed, cancelled
- `priority` (string, optional): low, normal, high, urgent
- `payment_status` (string, optional): unpaid, partial, paid
- `fulfillment_status` (string, optional): pending, partial, completed
- `warehouse_id` (string, optional): Filter by warehouse
- `date_from`, `date_to` (date, optional): Order date range

### Get Sales Order
```http
GET /api/v1/sales/orders/:id
```

### Create Sales Order
```http
POST /api/v1/sales/orders
Content-Type: application/json

{
  "code": "SO-2024-001",
  "quotation_id": "uuid",
  "customer_id": "uuid",
  "order_date": "2024-01-01T00:00:00Z",
  "expected_date": "2024-01-15T00:00:00Z",
  "priority": "normal",
  "reference_no": "REF-SO-001",
  "customer_po": "PO-CUST-123",
  "salesperson_id": "uuid",
  "currency": "USD",
  "exchange_rate": 1.0,
  "payment_term_days": 30,
  "warehouse_id": "uuid",
  "shipping_method": "Standard Shipping",
  "shipping_address": "123 Main St",
  "shipping_city": "New York",
  "shipping_state": "NY",
  "shipping_zip": "10001",
  "shipping_country": "USA",
  "shipping_cost": 50.00,
  "other_cost": 0.00,
  "notes": "Important order",
  "internal_notes": "Priority customer",
  "terms": "Payment terms: Net 30",
  "tags": ["demo", "urgent"],
  "items": [
    {
      "product_id": "uuid",
      "description": "High quality product",
      "quantity": 5.0,
      "unit_price": 99.99,
      "discount_type": "percentage",
      "discount": 10.0,
      "tax_id": "uuid",
      "notes": ""
    }
  ]
}
```

**Business Rules:**
- Links to quotation (optional)
- Warehouse must exist and be active
- Status defaults to "draft"
- Payment status defaults to "unpaid"
- Fulfillment status defaults to "pending"
- delivered_qty and invoiced_qty initialize to 0

### Update Sales Order
```http
PUT /api/v1/sales/orders/:id
Content-Type: application/json

{
  "status": "confirmed",
  "priority": "high",
  "payment_status": "partial",
  "fulfillment_status": "partial"
}
```

**Business Rules:**
- Only draft or confirmed orders can be updated
- Confirming order sets confirmed_date
- Changing costs recalculates grand_total

### Delete Sales Order
```http
DELETE /api/v1/sales/orders/:id
```

**Business Rules:**
- Only draft orders can be deleted
- Cannot delete if deliveries or invoices exist

### Get Order Statistics
```http
GET /api/v1/sales/orders/stats
```

**Response:**
```json
{
  "total_orders": 100,
  "total_value": 250000.00,
  "orders_by_status": [...],
  "orders_by_payment_status": [...],
  "orders_by_fulfillment_status": [...]
}
```

---

## 3. Delivery Orders API

### List Delivery Orders
```http
GET /api/v1/sales/deliveries
```

**Query Parameters:**
- `page`, `limit`, `search`, `active` (standard)
- `order_id` (string, optional): Filter by sales order
- `customer_id` (string, optional): Filter by customer
- `warehouse_id` (string, optional): Filter by warehouse
- `status` (string, optional): draft, ready, in_transit, delivered, cancelled
- `date_from`, `date_to` (date, optional): Delivery date range

### Get Delivery Order
```http
GET /api/v1/sales/deliveries/:id
```

### Create Delivery Order
```http
POST /api/v1/sales/deliveries
Content-Type: application/json

{
  "code": "DO-2024-001",
  "order_id": "uuid",
  "delivery_date": "2024-01-10T00:00:00Z",
  "scheduled_date": "2024-01-10T00:00:00Z",
  "reference_no": "REF-DO-001",
  "tracking_number": "TRACK-123456",
  "warehouse_id": "uuid",
  "shipping_method": "Standard Shipping",
  "shipping_cost": 50.00,
  "courier_name": "Express Courier",
  "driver_name": "John Driver",
  "vehicle_number": "B-1234-XYZ",
  "shipping_address": "123 Main St",
  "shipping_city": "New York",
  "shipping_state": "NY",
  "shipping_zip": "10001",
  "shipping_country": "USA",
  "recipient_name": "Jane Doe",
  "recipient_phone": "+1234567890",
  "recipient_email": "jane@example.com",
  "notes": "Fragile items",
  "internal_notes": "Deliver morning only",
  "tags": ["urgent", "fragile"],
  "items": [
    {
      "order_item_id": "uuid",
      "product_id": "uuid",
      "delivered_qty": 3.0,
      "serial_numbers": ["SN001", "SN002", "SN003"],
      "batch_numbers": ["BATCH-2024-001"],
      "notes": "Partial delivery"
    }
  ]
}
```

**Business Rules:**
- **CRITICAL:** Updates sales_order_items.delivered_qty in transaction
- Validates order_item_id exists
- Validates product_id matches order item
- Checks remaining quantity: delivered_qty <= (quantity - previous_delivered_qty)
- Updates sales_order.fulfillment_status:
  - pending: no items delivered
  - partial: some items delivered
  - completed: all items delivered
- Status defaults to "draft"

### Update Delivery Order
```http
PUT /api/v1/sales/deliveries/:id
Content-Type: application/json

{
  "status": "delivered",
  "actual_date": "2024-01-10T14:30:00Z",
  "received_by": "Jane Doe",
  "receiver_signature": "https://cdn.example.com/signatures/sig-001.jpg",
  "delivery_proof": ["https://cdn.example.com/proof/photo-001.jpg"]
}
```

**Business Rules:**
- Only draft or ready deliveries can be updated
- Status "delivered" requires actual_date

### Delete Delivery Order
```http
DELETE /api/v1/sales/deliveries/:id
```

**Business Rules:**
- **CRITICAL:** Only draft deliveries can be deleted
- **Reverts sales_order_items.delivered_qty in transaction**
- **Recalculates sales_order.fulfillment_status**

### Get Delivery Statistics
```http
GET /api/v1/sales/deliveries/stats
```

**Response:**
```json
{
  "total_deliveries": 75,
  "today_deliveries": 5,
  "deliveries_by_status": [...],
  "pending_deliveries": 10,
  "completed_deliveries": 60
}
```

---

## 4. Sales Invoices API

### List Sales Invoices
```http
GET /api/v1/sales/invoices
```

**Query Parameters:**
- `page`, `limit`, `search`, `customer_id`, `active` (standard)
- `order_id` (string, optional): Filter by sales order
- `status` (string, optional): draft, sent, partial_paid, paid, overdue, cancelled
- `payment_status` (string, optional): unpaid, partial, paid, overdue
- `date_from`, `date_to` (date, optional): Invoice date range
- `due_date_from`, `due_date_to` (date, optional): Due date range

### Get Sales Invoice
```http
GET /api/v1/sales/invoices/:id
```

### Create Sales Invoice
```http
POST /api/v1/sales/invoices
Content-Type: application/json

{
  "code": "INV-2024-001",
  "order_id": "uuid",
  "delivery_order_id": "uuid",
  "customer_id": "uuid",
  "invoice_date": "2024-01-12T00:00:00Z",
  "due_date": "2024-02-11T00:00:00Z",
  "reference_no": "REF-INV-001",
  "customer_po": "PO-CUST-123",
  "tax_invoice_no": "FP-001-2024",
  "currency": "USD",
  "exchange_rate": 1.0,
  "payment_term_days": 30,
  "income_account_id": "uuid",
  "ar_account_id": "uuid",
  "billing_address": "123 Main St",
  "billing_city": "New York",
  "billing_state": "NY",
  "billing_zip": "10001",
  "billing_country": "USA",
  "shipping_cost": 50.00,
  "other_cost": 0.00,
  "notes": "Thank you for your business",
  "internal_notes": "Payment expected on time",
  "terms": "Payment terms: Net 30. Late payment subject to 2% interest.",
  "tags": ["demo"],
  "items": [
    {
      "order_item_id": "uuid",
      "product_id": "uuid",
      "description": "High quality product",
      "quantity": 5.0,
      "unit_price": 99.99,
      "discount_type": "percentage",
      "discount": 10.0,
      "tax_id": "uuid",
      "account_id": "uuid",
      "notes": ""
    }
  ]
}
```

**Business Rules:**
- **CRITICAL:** Updates sales_order_items.invoiced_qty in transaction
- Links to order_id or delivery_order_id
- Validates remaining invoiced quantity
- Auto-calculates amount_due = grandtotal - amount_paid
- Status defaults to "draft"
- Payment status defaults to "unpaid"

### Update Sales Invoice
```http
PUT /api/v1/sales/invoices/:id
Content-Type: application/json

{
  "status": "sent",
  "payment_status": "partial",
  "amount_paid": 250.00
}
```

**Business Rules:**
- Only draft invoices can be freely updated
- **Auto-calculates payment_status:**
  - amount_paid == 0 → unpaid
  - 0 < amount_paid < grandtotal → partial
  - amount_paid >= grandtotal → paid
  - past due_date && amount_due > 0 → overdue
- **Recalculates amount_due = grandtotal - amount_paid**
- Status "sent" sets sent_at timestamp
- Status "paid" sets paid_at timestamp

### Delete Sales Invoice
```http
DELETE /api/v1/sales/invoices/:id
```

**Business Rules:**
- **CRITICAL:** Only draft invoices can be deleted
- **Reverts sales_order_items.invoiced_qty in transaction**

### Get Invoice Statistics
```http
GET /api/v1/sales/invoices/stats
```

**Response:**
```json
{
  "total_invoices": 120,
  "total_value": 300000.00,
  "total_paid": 200000.00,
  "total_due": 100000.00,
  "overdue_invoices": 5,
  "overdue_amount": 15000.00,
  "invoices_by_status": [...],
  "invoices_by_payment_status": [...]
}
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
- `403 Forbidden` - Insufficient permissions or business rule violation
- `404 Not Found` - Resource not found
- `409 Conflict` - Duplicate code or constraint violation

### Server Errors (5xx)
- `500 Internal Server Error` - Server-side error

**Error Response Format:**
```json
{
  "error": "customer not found: record not found",
  "details": {}
}
```

---

## Business Rule Violations

Common business rule errors:

**Quotations:**
```json
{"error": "only draft or rejected quotations can be deleted"}
```

**Orders:**
```json
{"error": "only draft orders can be deleted"}
```

**Deliveries:**
```json
{"error": "delivered quantity exceeds remaining quantity"}
{"error": "only draft delivery orders can be deleted"}
```

**Invoices:**
```json
{"error": "invoiced quantity exceeds remaining quantity"}
{"error": "only draft invoices can be deleted"}
```

---

## Testing with cURL

### Example 1: Create Sales Quotation
```bash
curl -X POST http://acme.localhost:8080/api/v1/sales/quotations \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "code": "SQ-2024-002",
    "customer_id": "CUSTOMER_UUID",
    "quotation_date": "2024-01-01T00:00:00Z",
    "valid_until": "2024-01-31T00:00:00Z",
    "currency": "USD",
    "exchange_rate": 1.0,
    "payment_term_days": 30,
    "items": [
      {
        "product_id": "PRODUCT_UUID",
        "quantity": 10.0,
        "unit_price": 99.99,
        "discount_type": "percentage",
        "discount": 10.0,
        "tax_id": "TAX_UUID"
      }
    ]
  }'
```

### Example 2: Create Delivery Order
```bash
curl -X POST http://acme.localhost:8080/api/v1/sales/deliveries \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "code": "DO-2024-002",
    "order_id": "ORDER_UUID",
    "delivery_date": "2024-01-15T00:00:00Z",
    "scheduled_date": "2024-01-15T00:00:00Z",
    "warehouse_id": "WAREHOUSE_UUID",
    "items": [
      {
        "order_item_id": "ORDER_ITEM_UUID",
        "product_id": "PRODUCT_UUID",
        "delivered_qty": 5.0
      }
    ]
  }'
```

### Example 3: Update Invoice Payment
```bash
curl -X PUT http://acme.localhost:8080/api/v1/sales/invoices/INVOICE_UUID \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "amount_paid": 500.00
  }'
```

---

## Workflow Examples

### Complete Sales Flow

#### 1. Create Quotation
```bash
POST /api/v1/sales/quotations
# Returns quotation_id
```

#### 2. Send Quotation
```bash
PUT /api/v1/sales/quotations/:quotation_id
{"status": "sent"}
```

#### 3. Accept Quotation
```bash
PUT /api/v1/sales/quotations/:quotation_id
{"status": "accepted"}
```

#### 4. Create Sales Order
```bash
POST /api/v1/sales/orders
{
  "quotation_id": "quotation_id",
  "code": "SO-2024-XXX",
  ...
}
# Returns order_id
```

#### 5. Confirm Order
```bash
PUT /api/v1/sales/orders/:order_id
{"status": "confirmed"}
```

#### 6. Create Delivery (Partial)
```bash
POST /api/v1/sales/deliveries
{
  "order_id": "order_id",
  "items": [{"order_item_id": "...", "delivered_qty": 3.0}]
}
# Updates order.delivered_qty = 3.0
# Updates order.fulfillment_status = "partial"
```

#### 7. Create Invoice
```bash
POST /api/v1/sales/invoices
{
  "order_id": "order_id",
  "items": [{"order_item_id": "...", "quantity": 5.0}]
}
# Updates order.invoiced_qty = 5.0
```

#### 8. Record Payment
```bash
PUT /api/v1/sales/invoices/:invoice_id
{"amount_paid": 500.00}
# Auto-updates payment_status to "paid"
# Recalculates amount_due
```

---

## Notes

1. **Multi-tenant Isolation**: All data is automatically isolated by tenant based on subdomain
2. **Soft Deletes**: DELETE operations set `active=false` instead of removing records
3. **Audit Trail**: All records track `created_by`, `updated_by`, `created_at`, `updated_at`
4. **Code Uniqueness**: Codes must be unique per tenant (not globally)
5. **Pagination**: Default limit is 10, maximum is 100 items per page
6. **Partial Updates**: PUT requests only update provided fields
7. **Auto Calculations**: Totals, payment status, fulfillment status are auto-calculated
8. **Transactions**: Critical operations use DB transactions for data consistency
9. **Workflow Enforcement**: Status transitions are validated per business rules
10. **Quantity Tracking**: delivered_qty and invoiced_qty are cumulative and validated

---

**Generated:** 2024-02-17  
**Version:** 1.0.0  
**API Version:** v1  
**Total Endpoints:** 24
