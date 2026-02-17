package handlers

import (
	"testing"
	"time"

	"github.com/madcok-co/unicorn-demo/internal/domain"
)

func TestCreatePurchaseRequest_Success(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	defer teardownTestDB(db)

	ctx := setupTestContext("acme", "user-123")
	ctx.Set("db", db)

	// Create test data
	user := createTestUser(db, "acme")
	product := createTestProduct(db, "acme")

	// Test data
	req := domain.CreatePurchaseRequestDTO{
		Code:         "PR-2024-001",
		RequesterID:  user.ID,
		RequestDate:  time.Now(),
		RequiredDate: time.Now().AddDate(0, 0, 14),
		Priority:     "normal",
		Purpose:      "Need materials for production",
		Items: []domain.CreatePurchaseRequestItemDTO{
			{
				ProductID:      product.ID,
				Description:    "Raw materials",
				Quantity:       50,
				EstimatedPrice: 75000,
			},
		},
	}

	// Execute
	result, err := CreatePurchaseRequest(ctx, req)

	// Assert
	assertNoError(t, err, "CreatePurchaseRequest should not return error")
	assertNotNil(t, result, "CreatePurchaseRequest should return request")
	assertEqual(t, result.Code, req.Code, "Request code should match")
	assertEqual(t, result.RequesterID, user.ID, "Requester ID should match")
	assertEqual(t, result.RequesterName, user.Name, "Requester name should match")
	assertEqual(t, result.Status, "draft", "Initial status should be draft")
	assertEqual(t, result.TenantID, "acme", "Tenant ID should match")
	assertEqual(t, len(result.Items), 1, "Should have 1 item")

	// Verify item
	item := result.Items[0]
	assertEqual(t, item.ProductID, product.ID, "Item product ID should match")
	assertEqual(t, item.Quantity, 50.0, "Item quantity should match")
	assertEqual(t, item.OrderedQty, 0.0, "Initial ordered qty should be 0")
}

func TestCreatePurchaseRequest_InvalidRequester(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	defer teardownTestDB(db)

	ctx := setupTestContext("acme", "user-123")
	ctx.Set("db", db)

	product := createTestProduct(db, "acme")

	// Test data with non-existent requester
	req := domain.CreatePurchaseRequestDTO{
		Code:         "PR-2024-001",
		RequesterID:  "non-existent-user",
		RequestDate:  time.Now(),
		RequiredDate: time.Now().AddDate(0, 0, 14),
		Items: []domain.CreatePurchaseRequestItemDTO{
			{
				ProductID: product.ID,
				Quantity:  50,
			},
		},
	}

	// Execute
	_, err := CreatePurchaseRequest(ctx, req)

	// Assert
	assertError(t, err, "CreatePurchaseRequest should fail with invalid requester")
}

func TestUpdatePurchaseRequest_ApprovalWorkflow(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	defer teardownTestDB(db)

	ctx := setupTestContext("acme", "user-123")
	ctx.Set("db", db)

	// Create test data
	user := createTestUser(db, "acme")
	product := createTestProduct(db, "acme")
	request := createTestPurchaseRequest(db, "acme", user, product)

	// Test workflow: draft -> submitted
	reqSubmit := domain.UpdatePurchaseRequestDTO{
		Status: "submitted",
	}
	resultSubmit, err := UpdatePurchaseRequest(ctx, request.ID, reqSubmit)
	assertNoError(t, err, "Should transition to submitted")
	assertEqual(t, resultSubmit.Status, "submitted", "Status should be submitted")
	assertNotNil(t, resultSubmit.SubmittedAt, "SubmittedAt should be set")

	// Test workflow: submitted -> approved
	reqApprove := domain.UpdatePurchaseRequestDTO{
		Status: "approved",
	}
	resultApprove, err := UpdatePurchaseRequest(ctx, request.ID, reqApprove)
	assertNoError(t, err, "Should transition to approved")
	assertEqual(t, resultApprove.Status, "approved", "Status should be approved")
	assertNotNil(t, resultApprove.ApprovedDate, "ApprovedDate should be set")
}

func TestUpdatePurchaseRequest_Rejection(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	defer teardownTestDB(db)

	ctx := setupTestContext("acme", "user-123")
	ctx.Set("db", db)

	// Create test data
	user := createTestUser(db, "acme")
	product := createTestProduct(db, "acme")
	request := createTestPurchaseRequest(db, "acme", user, product)

	// Submit first
	request.Status = "submitted"
	db.Save(request)

	// Test workflow: submitted -> rejected
	reqReject := domain.UpdatePurchaseRequestDTO{
		Status: "rejected",
	}
	resultReject, err := UpdatePurchaseRequest(ctx, request.ID, reqReject)
	assertNoError(t, err, "Should transition to rejected")
	assertEqual(t, resultReject.Status, "rejected", "Status should be rejected")
	assertNotNil(t, resultReject.RejectedAt, "RejectedAt should be set")
}

func TestDeletePurchaseRequest_Draft(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	defer teardownTestDB(db)

	ctx := setupTestContext("acme", "user-123")
	ctx.Set("db", db)

	// Create test data
	user := createTestUser(db, "acme")
	product := createTestProduct(db, "acme")
	request := createTestPurchaseRequest(db, "acme", user, product)

	// Execute
	err := DeletePurchaseRequest(ctx, request.ID)

	// Assert
	assertNoError(t, err, "DeletePurchaseRequest should not return error")

	// Verify request is soft deleted
	var deletedRequest domain.PurchaseRequest
	db.Where("id = ?", request.ID).First(&deletedRequest)
	assertEqual(t, deletedRequest.Active, false, "Request should be soft deleted")
}

func TestDeletePurchaseRequest_Approved(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	defer teardownTestDB(db)

	ctx := setupTestContext("acme", "user-123")
	ctx.Set("db", db)

	// Create test data
	user := createTestUser(db, "acme")
	product := createTestProduct(db, "acme")
	request := createTestPurchaseRequest(db, "acme", user, product)

	// Change status to approved
	request.Status = "approved"
	db.Save(request)

	// Execute
	err := DeletePurchaseRequest(ctx, request.ID)

	// Assert
	assertError(t, err, "DeletePurchaseRequest should fail for approved request")
}

func TestListPurchaseRequests_WithFilters(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	defer teardownTestDB(db)

	ctx := setupTestContext("acme", "user-123")
	ctx.Set("db", db)

	// Create test data
	user1 := createTestUser(db, "acme")
	user2 := createTestUser(db, "acme")
	product := createTestProduct(db, "acme")

	request1 := createTestPurchaseRequest(db, "acme", user1, product)
	request1.Status = "draft"
	request1.Priority = "normal"
	db.Save(request1)

	request2 := createTestPurchaseRequest(db, "acme", user2, product)
	request2.Status = "submitted"
	request2.Priority = "urgent"
	db.Save(request2)

	request3 := createTestPurchaseRequest(db, "acme", user1, product)
	request3.Status = "approved"
	request3.Priority = "high"
	db.Save(request3)

	// Test: List all requests
	req := domain.ListPurchaseRequestsDTO{
		Page:  1,
		Limit: 10,
	}

	result, err := ListPurchaseRequests(ctx, req)
	assertNoError(t, err, "ListPurchaseRequests should not return error")
	assertNotNil(t, result, "ListPurchaseRequests should return result")

	// Test: Filter by requester
	reqRequester := domain.ListPurchaseRequestsDTO{
		Page:        1,
		Limit:       10,
		RequesterID: user1.ID,
	}

	resultRequester, err := ListPurchaseRequests(ctx, reqRequester)
	assertNoError(t, err, "ListPurchaseRequests with requester filter should not return error")
	assertNotNil(t, resultRequester, "ListPurchaseRequests with requester filter should return result")

	// Test: Filter by status
	reqStatus := domain.ListPurchaseRequestsDTO{
		Page:   1,
		Limit:  10,
		Status: "submitted",
	}

	resultStatus, err := ListPurchaseRequests(ctx, reqStatus)
	assertNoError(t, err, "ListPurchaseRequests with status filter should not return error")
	assertNotNil(t, resultStatus, "ListPurchaseRequests with status filter should return result")

	// Test: Filter by priority
	reqPriority := domain.ListPurchaseRequestsDTO{
		Page:     1,
		Limit:    10,
		Priority: "urgent",
	}

	resultPriority, err := ListPurchaseRequests(ctx, reqPriority)
	assertNoError(t, err, "ListPurchaseRequests with priority filter should not return error")
	assertNotNil(t, resultPriority, "ListPurchaseRequests with priority filter should return result")
}

func TestGetPurchaseRequestStats_Success(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	defer teardownTestDB(db)

	ctx := setupTestContext("acme", "user-123")
	ctx.Set("db", db)

	// Create test data
	user := createTestUser(db, "acme")
	product := createTestProduct(db, "acme")

	request1 := createTestPurchaseRequest(db, "acme", user, product)
	request1.Status = "draft"
	request1.Priority = "normal"
	db.Save(request1)

	request2 := createTestPurchaseRequest(db, "acme", user, product)
	request2.Status = "submitted"
	request2.Priority = "urgent"
	db.Save(request2)

	request3 := createTestPurchaseRequest(db, "acme", user, product)
	request3.Status = "approved"
	request3.Priority = "high"
	db.Save(request3)

	// Execute
	stats, err := GetPurchaseRequestStats(ctx)

	// Assert
	assertNoError(t, err, "GetPurchaseRequestStats should not return error")
	assertNotNil(t, stats, "GetPurchaseRequestStats should return stats")
	assertNotNil(t, stats["total_requests"], "Stats should contain total_requests")
	assertNotNil(t, stats["pending_approval"], "Stats should contain pending_approval")
	assertNotNil(t, stats["approved_requests"], "Stats should contain approved_requests")
	assertNotNil(t, stats["requests_by_status"], "Stats should contain requests_by_status")
	assertNotNil(t, stats["requests_by_priority"], "Stats should contain requests_by_priority")
}

func TestPurchaseRequest_MultiTenantIsolation(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	defer teardownTestDB(db)

	// Create requests for different tenants
	userAcme := createTestUser(db, "acme")
	productAcme := createTestProduct(db, "acme")
	_ = createTestPurchaseRequest(db, "acme", userAcme, productAcme)

	userGlobex := createTestUser(db, "globex")
	productGlobex := createTestProduct(db, "globex")
	requestGlobex := createTestPurchaseRequest(db, "globex", userGlobex, productGlobex)

	// Test: Tenant "acme" should only see their request
	ctxAcme := setupTestContext("acme", "user-123")
	ctxAcme.Set("db", db)

	req := domain.ListPurchaseRequestsDTO{
		Page:  1,
		Limit: 100,
	}

	resultAcme, err := ListPurchaseRequests(ctxAcme, req)
	assertNoError(t, err, "ListPurchaseRequests should not return error")
	assertEqual(t, resultAcme.Total, int64(1), "Tenant acme should see 1 request")

	// Test: Tenant "globex" should only see their request
	ctxGlobex := setupTestContext("globex", "user-456")
	ctxGlobex.Set("db", db)

	resultGlobex, err := ListPurchaseRequests(ctxGlobex, req)
	assertNoError(t, err, "ListPurchaseRequests should not return error")
	assertEqual(t, resultGlobex.Total, int64(1), "Tenant globex should see 1 request")

	// Test: Cannot get request from different tenant
	_, err = GetPurchaseRequest(ctxAcme, requestGlobex.ID)
	assertError(t, err, "Should not be able to get request from different tenant")
}

func TestPurchaseRequest_PriorityHandling(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	defer teardownTestDB(db)

	ctx := setupTestContext("acme", "user-123")
	ctx.Set("db", db)

	// Create test data
	user := createTestUser(db, "acme")
	product := createTestProduct(db, "acme")

	tests := []struct {
		name     string
		priority string
		expected string
	}{
		{
			name:     "Low priority",
			priority: "low",
			expected: "low",
		},
		{
			name:     "Normal priority",
			priority: "normal",
			expected: "normal",
		},
		{
			name:     "High priority",
			priority: "high",
			expected: "high",
		},
		{
			name:     "Urgent priority",
			priority: "urgent",
			expected: "urgent",
		},
		{
			name:     "Empty priority defaults to normal",
			priority: "",
			expected: "normal",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := domain.CreatePurchaseRequestDTO{
				Code:         "PR-" + tt.name,
				RequesterID:  user.ID,
				RequestDate:  time.Now(),
				RequiredDate: time.Now().AddDate(0, 0, 14),
				Priority:     tt.priority,
				Items: []domain.CreatePurchaseRequestItemDTO{
					{
						ProductID: product.ID,
						Quantity:  10,
					},
				},
			}

			result, err := CreatePurchaseRequest(ctx, req)
			assertNoError(t, err, "CreatePurchaseRequest should not return error")
			assertEqual(t, result.Priority, tt.expected, "Priority should match expected")
		})
	}
}

func TestUpdatePurchaseRequest_StatusRestrictions(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	defer teardownTestDB(db)

	ctx := setupTestContext("acme", "user-123")
	ctx.Set("db", db)

	// Create test data
	user := createTestUser(db, "acme")
	product := createTestProduct(db, "acme")
	request := createTestPurchaseRequest(db, "acme", user, product)

	// Change status to completed
	request.Status = "completed"
	db.Save(request)

	// Try to update completed request
	req := domain.UpdatePurchaseRequestDTO{
		Notes: "Trying to update completed request",
	}

	// Execute
	_, err := UpdatePurchaseRequest(ctx, request.ID, req)

	// Assert
	assertError(t, err, "UpdatePurchaseRequest should fail for completed request")
}

func TestPurchaseRequest_ItemTracking(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	defer teardownTestDB(db)

	ctx := setupTestContext("acme", "user-123")
	ctx.Set("db", db)

	// Create test data
	user := createTestUser(db, "acme")
	product := createTestProduct(db, "acme")

	// Create purchase request
	req := domain.CreatePurchaseRequestDTO{
		Code:         "PR-2024-001",
		RequesterID:  user.ID,
		RequestDate:  time.Now(),
		RequiredDate: time.Now().AddDate(0, 0, 14),
		Items: []domain.CreatePurchaseRequestItemDTO{
			{
				ProductID:      product.ID,
				Quantity:       100,
				EstimatedPrice: 75000,
			},
		},
	}

	result, err := CreatePurchaseRequest(ctx, req)
	assertNoError(t, err, "CreatePurchaseRequest should not return error")

	// Verify initial item tracking
	assertEqual(t, len(result.Items), 1, "Should have 1 item")
	item := result.Items[0]
	assertEqual(t, item.Quantity, 100.0, "Item quantity should be 100")
	assertEqual(t, item.OrderedQty, 0.0, "Initial ordered qty should be 0")

	// Simulate converting to PO by updating ordered_qty
	item.OrderedQty = 100.0
	db.Save(&item)

	// Reload and verify
	var updatedRequest domain.PurchaseRequest
	db.Preload("Items").Where("id = ?", result.ID).First(&updatedRequest)
	assertEqual(t, updatedRequest.Items[0].OrderedQty, 100.0, "Ordered qty should be updated to 100")
}
