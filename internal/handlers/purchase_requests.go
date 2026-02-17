package handlers

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/madcok-co/unicorn-demo/internal/domain"
	"github.com/madcok-co/unicorn/contrib/pagination"
	"github.com/madcok-co/unicorn/core/pkg/context"
)

// ListPurchaseRequests returns paginated list of purchase requests
func ListPurchaseRequests(ctx *context.Context, req domain.ListPurchaseRequestsDTO) (*pagination.OffsetResult, error) {
	db := getDB(ctx)
	tenantID, err := getTenantID(ctx)
	if err != nil {
		return nil, err
	}

	var requests []domain.PurchaseRequest
	var total int64

	query := db.Model(&domain.PurchaseRequest{}).
		Preload("Items").
		Where("tenant_id = ?", tenantID)

	// Apply filters
	if req.Search != "" {
		search := "%" + req.Search + "%"
		query = query.Where("code LIKE ? OR requester_name LIKE ? OR reference_no LIKE ? OR purpose LIKE ?", search, search, search, search)
	}
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}
	if req.RequesterID != "" {
		query = query.Where("requester_id = ?", req.RequesterID)
	}
	if req.ApproverID != "" {
		query = query.Where("approver_id = ?", req.ApproverID)
	}
	if req.DepartmentID != "" {
		query = query.Where("department_id = ?", req.DepartmentID)
	}
	if req.Priority != "" {
		query = query.Where("priority = ?", req.Priority)
	}
	if req.DateFrom != nil {
		query = query.Where("request_date >= ?", req.DateFrom)
	}
	if req.DateTo != nil {
		query = query.Where("request_date <= ?", req.DateTo)
	}
	if req.Active != nil {
		query = query.Where("active = ?", *req.Active)
	}

	// Count total
	query.Count(&total)

	// Apply pagination
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 {
		req.Limit = 10
	}
	offset := (req.Page - 1) * req.Limit
	if err := query.Offset(offset).Limit(req.Limit).Order("request_date DESC, created_at DESC").Find(&requests).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch purchase requests: %w", err)
	}

	return &pagination.OffsetResult{
		Data:       requests,
		Total:      total,
		Page:       req.Page,
		Limit:      req.Limit,
		TotalPages: int((total + int64(req.Limit) - 1) / int64(req.Limit)),
	}, nil
}

// GetPurchaseRequest returns a single purchase request by ID
func GetPurchaseRequest(ctx *context.Context, id string) (*domain.PurchaseRequest, error) {
	db := getDB(ctx)
	tenantID, err := getTenantID(ctx)
	if err != nil {
		return nil, err
	}

	var request domain.PurchaseRequest
	if err := db.Preload("Items").
		Where("id = ? AND tenant_id = ?", id, tenantID).
		First(&request).Error; err != nil {
		return nil, fmt.Errorf("purchase request not found: %w", err)
	}

	return &request, nil
}

// CreatePurchaseRequest creates a new purchase request
func CreatePurchaseRequest(ctx *context.Context, req domain.CreatePurchaseRequestDTO) (*domain.PurchaseRequest, error) {
	db := getDB(ctx)
	tenantID, err := getTenantID(ctx)
	if err != nil {
		return nil, err
	}
	userID, _ := getUserID(ctx)

	// Check duplicate code
	var existing domain.PurchaseRequest
	if err := db.Where("code = ? AND tenant_id = ?", req.Code, tenantID).First(&existing).Error; err == nil {
		return nil, fmt.Errorf("purchase request with code %s already exists", req.Code)
	}

	// Get requester info
	var requester domain.User
	if err := db.Where("id = ? AND tenant_id = ? AND active = ?", req.RequesterID, tenantID, true).First(&requester).Error; err != nil {
		return nil, fmt.Errorf("requester not found: %w", err)
	}

	// Get department info if provided
	var departmentName string
	if req.DepartmentID != "" {
		// TODO: Add Department model lookup when available
		departmentName = "Department" // Placeholder until Department model exists
	}

	// Get approver info if provided
	var approverName string
	if req.ApproverID != "" {
		var approver domain.User
		if err := db.Where("id = ? AND tenant_id = ? AND active = ?", req.ApproverID, tenantID, true).First(&approver).Error; err == nil {
			approverName = approver.Name
		}
	}

	// Process items and calculate estimated costs
	items := make([]domain.PurchaseRequestItem, len(req.Items))
	for i, itemReq := range req.Items {
		// Get product info
		var product domain.Product
		if err := db.Where("id = ? AND tenant_id = ? AND active = ?", itemReq.ProductID, tenantID, true).First(&product).Error; err != nil {
			return nil, fmt.Errorf("product not found for item %d: %w", i+1, err)
		}

		items[i] = domain.PurchaseRequestItem{
			ID:             uuid.New().String(),
			TenantID:       tenantID,
			LineNumber:     i + 1,
			ProductID:      product.ID,
			ProductCode:    product.Code,
			ProductName:    product.Name,
			ProductType:    product.Type,
			Description:    itemReq.Description,
			Quantity:       itemReq.Quantity,
			OrderedQty:     0,
			UOM:            product.UOM,
			EstimatedPrice: itemReq.EstimatedPrice,
			Notes:          itemReq.Notes,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}
	}

	// Set default priority if not provided
	priority := req.Priority
	if priority == "" {
		priority = "normal"
	}

	request := &domain.PurchaseRequest{
		ID:             uuid.New().String(),
		Code:           req.Code,
		TenantID:       tenantID,
		RequesterID:    requester.ID,
		RequesterName:  requester.Name,
		DepartmentID:   req.DepartmentID,
		DepartmentName: departmentName,
		RequestDate:    req.RequestDate,
		RequiredDate:   req.RequiredDate,
		Status:         "draft",
		Priority:       priority,
		ReferenceNo:    req.ReferenceNo,
		Purpose:        req.Purpose,
		ApproverID:     req.ApproverID,
		ApproverName:   approverName,
		Notes:          req.Notes,
		InternalNotes:  req.InternalNotes,
		Tags:           req.Tags,
		Items:          items,
		Active:         true,
		CreatedBy:      userID,
		UpdatedBy:      userID,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := db.Create(request).Error; err != nil {
		return nil, fmt.Errorf("failed to create purchase request: %w", err)
	}

	return request, nil
}

// UpdatePurchaseRequest updates an existing purchase request
func UpdatePurchaseRequest(ctx *context.Context, id string, req domain.UpdatePurchaseRequestDTO) (*domain.PurchaseRequest, error) {
	db := getDB(ctx)
	tenantID, err := getTenantID(ctx)
	if err != nil {
		return nil, err
	}
	userID, _ := getUserID(ctx)

	var request domain.PurchaseRequest
	if err := db.Preload("Items").Where("id = ? AND tenant_id = ?", id, tenantID).First(&request).Error; err != nil {
		return nil, fmt.Errorf("purchase request not found: %w", err)
	}

	// Only draft or submitted requests can be updated (except for status changes)
	if request.Status != "draft" && request.Status != "submitted" && req.Status == "" {
		return nil, fmt.Errorf("only draft or submitted purchase requests can be updated")
	}

	// Update fields if provided
	if req.RequiredDate != nil {
		request.RequiredDate = *req.RequiredDate
	}
	if req.Status != "" {
		request.Status = req.Status
		// Update workflow timestamps
		now := time.Now()
		switch req.Status {
		case "submitted":
			if request.SubmittedAt == nil {
				request.SubmittedAt = &now
			}
		case "approved":
			if request.ApprovedDate == nil {
				request.ApprovedDate = &now
			}
		case "rejected":
			if request.RejectedAt == nil {
				request.RejectedAt = &now
			}
		case "completed":
			if request.CompletedAt == nil {
				request.CompletedAt = &now
			}
		}
	}
	if req.Priority != "" {
		request.Priority = req.Priority
	}
	if req.ReferenceNo != "" {
		request.ReferenceNo = req.ReferenceNo
	}
	if req.Purpose != "" {
		request.Purpose = req.Purpose
	}
	if req.ApproverID != "" {
		request.ApproverID = req.ApproverID
		// Update approver name
		var approver domain.User
		if err := db.Where("id = ? AND tenant_id = ?", req.ApproverID, tenantID).First(&approver).Error; err == nil {
			request.ApproverName = approver.Name
		}
	}
	if req.Notes != "" {
		request.Notes = req.Notes
	}
	if req.InternalNotes != "" {
		request.InternalNotes = req.InternalNotes
	}
	if req.Tags != nil {
		request.Tags = req.Tags
	}

	request.UpdatedBy = userID
	request.UpdatedAt = time.Now()

	if err := db.Save(&request).Error; err != nil {
		return nil, fmt.Errorf("failed to update purchase request: %w", err)
	}

	return &request, nil
}

// DeletePurchaseRequest soft deletes a purchase request
func DeletePurchaseRequest(ctx *context.Context, id string) error {
	db := getDB(ctx)
	tenantID, err := getTenantID(ctx)
	if err != nil {
		return err
	}
	userID, _ := getUserID(ctx)

	var request domain.PurchaseRequest
	if err := db.Where("id = ? AND tenant_id = ?", id, tenantID).First(&request).Error; err != nil {
		return fmt.Errorf("purchase request not found: %w", err)
	}

	// Only draft or rejected requests can be deleted
	if request.Status != "draft" && request.Status != "rejected" {
		return fmt.Errorf("only draft or rejected purchase requests can be deleted")
	}

	request.Active = false
	request.UpdatedBy = userID
	request.UpdatedAt = time.Now()

	if err := db.Save(&request).Error; err != nil {
		return fmt.Errorf("failed to delete purchase request: %w", err)
	}

	return nil
}

// GetPurchaseRequestStats returns purchase request statistics
func GetPurchaseRequestStats(ctx *context.Context) (map[string]interface{}, error) {
	db := getDB(ctx)
	tenantID, err := getTenantID(ctx)
	if err != nil {
		return nil, err
	}

	var totalRequests int64
	db.Model(&domain.PurchaseRequest{}).Where("tenant_id = ? AND active = ?", tenantID, true).Count(&totalRequests)

	var statusStats []struct {
		Status string
		Count  int64
	}
	db.Model(&domain.PurchaseRequest{}).
		Select("status, COUNT(*) as count").
		Where("tenant_id = ? AND active = ?", tenantID, true).
		Group("status").
		Scan(&statusStats)

	var priorityStats []struct {
		Priority string
		Count    int64
	}
	db.Model(&domain.PurchaseRequest{}).
		Select("priority, COUNT(*) as count").
		Where("tenant_id = ? AND active = ?", tenantID, true).
		Group("priority").
		Scan(&priorityStats)

	var pendingApproval int64
	db.Model(&domain.PurchaseRequest{}).
		Where("tenant_id = ? AND active = ? AND status = ?", tenantID, true, "submitted").
		Count(&pendingApproval)

	var approvedRequests int64
	db.Model(&domain.PurchaseRequest{}).
		Where("tenant_id = ? AND active = ? AND status = ?", tenantID, true, "approved").
		Count(&approvedRequests)

	return map[string]interface{}{
		"total_requests":       totalRequests,
		"pending_approval":     pendingApproval,
		"approved_requests":    approvedRequests,
		"requests_by_status":   statusStats,
		"requests_by_priority": priorityStats,
	}, nil
}
