package chatstore

import (
	"testing"
)

func TestNewChat(t *testing.T) {
	chat := NewChat()

	if chat == nil {
		t.Fatal("NewChat returned nil")
	}

	if chat.ID() == "" {
		t.Error("Expected ID to be set")
	}

	if chat.Status() != CHAT_STATUS_ACTIVE {
		t.Errorf("Expected status %s, got %s", CHAT_STATUS_ACTIVE, chat.Status())
	}

	if chat.Title() != "" {
		t.Error("Expected empty title")
	}

	if chat.Memo() != "" {
		t.Error("Expected empty memo")
	}

	if chat.CreatedAt() == "" {
		t.Error("Expected CreatedAt to be set")
	}

	if chat.UpdatedAt() == "" {
		t.Error("Expected UpdatedAt to be set")
	}

	if chat.SoftDeletedAt() == "" {
		t.Error("Expected SoftDeletedAt to be set")
	}

	metas, err := chat.Metas()
	if err != nil {
		t.Fatalf("Failed to get metas: %v", err)
	}

	if len(metas) != 0 {
		t.Errorf("Expected empty metas, got %d", len(metas))
	}
}

func TestNewChatFromExistingData(t *testing.T) {
	data := map[string]string{
		COLUMN_ID:              "test-id",
		COLUMN_OWNER_ID:        "owner-id",
		COLUMN_STATUS:          CHAT_STATUS_INACTIVE,
		COLUMN_TITLE:           "Test Chat",
		COLUMN_MEMO:            "Test memo",
		COLUMN_CREATED_AT:      "2024-01-01 00:00:00",
		COLUMN_UPDATED_AT:      "2024-01-02 00:00:00",
		COLUMN_SOFT_DELETED_AT: "",
		COLUMN_METAS:           `{"key1":"value1","key2":"value2"}`,
	}

	chat := NewChatFromExistingData(data)

	if chat == nil {
		t.Fatal("NewChatFromExistingData returned nil")
	}

	if chat.ID() != "test-id" {
		t.Errorf("Expected ID test-id, got %s", chat.ID())
	}

	if chat.OwnerID() != "owner-id" {
		t.Errorf("Expected OwnerID owner-id, got %s", chat.OwnerID())
	}

	if chat.Status() != CHAT_STATUS_INACTIVE {
		t.Errorf("Expected status %s, got %s", CHAT_STATUS_INACTIVE, chat.Status())
	}

	if chat.Title() != "Test Chat" {
		t.Errorf("Expected title 'Test Chat', got %s", chat.Title())
	}

	if chat.Memo() != "Test memo" {
		t.Errorf("Expected memo 'Test memo', got %s", chat.Memo())
	}

	if chat.CreatedAt() != "2024-01-01 00:00:00" {
		t.Errorf("Expected CreatedAt 2024-01-01 00:00:00, got %s", chat.CreatedAt())
	}

	if chat.UpdatedAt() != "2024-01-02 00:00:00" {
		t.Errorf("Expected UpdatedAt 2024-01-02 00:00:00, got %s", chat.UpdatedAt())
	}

	if chat.SoftDeletedAt() != "" {
		t.Errorf("Expected empty SoftDeletedAt, got %s", chat.SoftDeletedAt())
	}
}

func TestChatID(t *testing.T) {
	chat := NewChat()

	testID := "test-id-123"
	chat.SetID(testID)

	if chat.ID() != testID {
		t.Errorf("Expected ID %s, got %s", testID, chat.ID())
	}
}

func TestChatOwnerID(t *testing.T) {
	chat := NewChat()

	testOwnerID := "owner-id-456"
	chat.SetOwnerID(testOwnerID)

	if chat.OwnerID() != testOwnerID {
		t.Errorf("Expected OwnerID %s, got %s", testOwnerID, chat.OwnerID())
	}
}

func TestChatStatus(t *testing.T) {
	chat := NewChat()

	if chat.Status() != CHAT_STATUS_ACTIVE {
		t.Errorf("Expected initial status %s, got %s", CHAT_STATUS_ACTIVE, chat.Status())
	}

	chat.SetStatus(CHAT_STATUS_INACTIVE)

	if chat.Status() != CHAT_STATUS_INACTIVE {
		t.Errorf("Expected status %s, got %s", CHAT_STATUS_INACTIVE, chat.Status())
	}
}

func TestChatTitle(t *testing.T) {
	chat := NewChat()

	testTitle := "Test Chat Title"
	chat.SetTitle(testTitle)

	if chat.Title() != testTitle {
		t.Errorf("Expected title %s, got %s", testTitle, chat.Title())
	}
}

func TestChatMemo(t *testing.T) {
	chat := NewChat()

	testMemo := "Test memo content"
	chat.SetMemo(testMemo)

	if chat.Memo() != testMemo {
		t.Errorf("Expected memo %s, got %s", testMemo, chat.Memo())
	}
}

func TestChatCreatedAt(t *testing.T) {
	chat := NewChat()

	testTime := "2024-03-15 10:30:00"
	chat.SetCreatedAt(testTime)

	if chat.CreatedAt() != testTime {
		t.Errorf("Expected CreatedAt %s, got %s", testTime, chat.CreatedAt())
	}

	carbonTime := chat.CreatedAtCarbon()
	if carbonTime == nil {
		t.Fatal("CreatedAtCarbon returned nil")
	}

	if carbonTime.Format("Y-m-d H:i:s") != testTime {
		t.Errorf("Expected Carbon time %s, got %s", testTime, carbonTime.Format("Y-m-d H:i:s"))
	}
}

func TestChatUpdatedAt(t *testing.T) {
	chat := NewChat()

	testTime := "2024-03-15 10:30:00"
	chat.SetUpdatedAt(testTime)

	if chat.UpdatedAt() != testTime {
		t.Errorf("Expected UpdatedAt %s, got %s", testTime, chat.UpdatedAt())
	}

	carbonTime := chat.UpdatedAtCarbon()
	if carbonTime == nil {
		t.Fatal("UpdatedAtCarbon returned nil")
	}

	if carbonTime.Format("Y-m-d H:i:s") != testTime {
		t.Errorf("Expected Carbon time %s, got %s", testTime, carbonTime.Format("Y-m-d H:i:s"))
	}
}

func TestChatSoftDeletedAt(t *testing.T) {
	chat := NewChat()

	// Initially has MAX_DATE set, so IsSoftDeleted returns true
	if !chat.IsSoftDeleted() {
		t.Error("Expected IsSoftDeleted to be true initially (has MAX_DATE)")
	}

	testTime := "2024-03-15 10:30:00"
	chat.SetSoftDeletedAt(testTime)

	if chat.SoftDeletedAt() != testTime {
		t.Errorf("Expected SoftDeletedAt %s, got %s", testTime, chat.SoftDeletedAt())
	}

	if !chat.IsSoftDeleted() {
		t.Error("Expected IsSoftDeleted to be true after setting SoftDeletedAt")
	}

	carbonTime := chat.SoftDeletedAtCarbon()
	if carbonTime == nil {
		t.Fatal("SoftDeletedAtCarbon returned nil")
	}

	if carbonTime.Format("Y-m-d H:i:s") != testTime {
		t.Errorf("Expected Carbon time %s, got %s", testTime, carbonTime.Format("Y-m-d H:i:s"))
	}

	// Test setting to empty string (not soft deleted)
	chat.SetSoftDeletedAt("")
	if chat.IsSoftDeleted() {
		t.Error("Expected IsSoftDeleted to be false when SoftDeletedAt is empty")
	}
}

func TestChatMetas(t *testing.T) {
	chat := NewChat()

	// Test setting metas
	testMetas := map[string]string{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
	}

	err := chat.SetMetas(testMetas)
	if err != nil {
		t.Fatalf("Failed to set metas: %v", err)
	}

	// Test getting metas
	metas, err := chat.Metas()
	if err != nil {
		t.Fatalf("Failed to get metas: %v", err)
	}

	if len(metas) != 3 {
		t.Errorf("Expected 3 metas, got %d", len(metas))
	}

	if metas["key1"] != "value1" {
		t.Errorf("Expected key1 to be value1, got %s", metas["key1"])
	}

	if metas["key2"] != "value2" {
		t.Errorf("Expected key2 to be value2, got %s", metas["key2"])
	}

	if metas["key3"] != "value3" {
		t.Errorf("Expected key3 to be value3, got %s", metas["key3"])
	}
}

func TestChatMeta(t *testing.T) {
	chat := NewChat()

	// Set metas
	testMetas := map[string]string{
		"key1": "value1",
		"key2": "value2",
	}

	err := chat.SetMetas(testMetas)
	if err != nil {
		t.Fatalf("Failed to set metas: %v", err)
	}

	// Test getting single meta
	value, err := chat.Meta("key1")
	if err != nil {
		t.Fatalf("Failed to get meta: %v", err)
	}

	if value != "value1" {
		t.Errorf("Expected value1, got %s", value)
	}

	// Test getting non-existent meta
	value, err = chat.Meta("nonexistent")
	if err != nil {
		t.Fatalf("Failed to get non-existent meta: %v", err)
	}

	if value != "" {
		t.Errorf("Expected empty string for non-existent meta, got %s", value)
	}
}

func TestChatSetMeta(t *testing.T) {
	chat := NewChat()

	// Set single meta
	err := chat.SetMeta("key1", "value1")
	if err != nil {
		t.Fatalf("Failed to set meta: %v", err)
	}

	// Verify it was set
	value, err := chat.Meta("key1")
	if err != nil {
		t.Fatalf("Failed to get meta: %v", err)
	}

	if value != "value1" {
		t.Errorf("Expected value1, got %s", value)
	}
}

func TestChatUpsertMetas(t *testing.T) {
	chat := NewChat()

	// Set initial metas
	initialMetas := map[string]string{
		"key1": "value1",
		"key2": "value2",
	}

	err := chat.SetMetas(initialMetas)
	if err != nil {
		t.Fatalf("Failed to set initial metas: %v", err)
	}

	// Upsert new metas (update existing, add new)
	upsertMetas := map[string]string{
		"key1": "updated-value1",
		"key3": "value3",
	}

	err = chat.UpsertMetas(upsertMetas)
	if err != nil {
		t.Fatalf("Failed to upsert metas: %v", err)
	}

	// Verify results
	metas, err := chat.Metas()
	if err != nil {
		t.Fatalf("Failed to get metas: %v", err)
	}

	if len(metas) != 3 {
		t.Errorf("Expected 3 metas after upsert, got %d", len(metas))
	}

	if metas["key1"] != "updated-value1" {
		t.Errorf("Expected key1 to be updated-value1, got %s", metas["key1"])
	}

	if metas["key2"] != "value2" {
		t.Errorf("Expected key2 to remain value2, got %s", metas["key2"])
	}

	if metas["key3"] != "value3" {
		t.Errorf("Expected key3 to be value3, got %s", metas["key3"])
	}
}

func TestChatMetasEmpty(t *testing.T) {
	chat := NewChat()

	// Get metas when empty
	metas, err := chat.Metas()
	if err != nil {
		t.Fatalf("Failed to get metas: %v", err)
	}

	if len(metas) != 0 {
		t.Errorf("Expected empty metas, got %d", len(metas))
	}
}

func TestChatSetMetasEmpty(t *testing.T) {
	chat := NewChat()

	// Set some metas first
	testMetas := map[string]string{
		"key1": "value1",
	}

	err := chat.SetMetas(testMetas)
	if err != nil {
		t.Fatalf("Failed to set metas: %v", err)
	}

	// Set empty metas
	err = chat.SetMetas(map[string]string{})
	if err != nil {
		t.Fatalf("Failed to set empty metas: %v", err)
	}

	// Verify they're empty
	metas, err := chat.Metas()
	if err != nil {
		t.Fatalf("Failed to get metas: %v", err)
	}

	if len(metas) != 0 {
		t.Errorf("Expected empty metas after setting empty, got %d", len(metas))
	}
}

func TestChatSettersReturnChatInterface(t *testing.T) {
	chat := NewChat()

	// Test that all setters return ChatInterface for chaining
	result := chat.SetID("test-id")
	if result == nil {
		t.Error("SetID returned nil")
	}

	result = chat.SetOwnerID("owner-id")
	if result == nil {
		t.Error("SetOwnerID returned nil")
	}

	result = chat.SetStatus(CHAT_STATUS_INACTIVE)
	if result == nil {
		t.Error("SetStatus returned nil")
	}

	result = chat.SetTitle("title")
	if result == nil {
		t.Error("SetTitle returned nil")
	}

	result = chat.SetMemo("memo")
	if result == nil {
		t.Error("SetMemo returned nil")
	}

	result = chat.SetCreatedAt("2024-01-01 00:00:00")
	if result == nil {
		t.Error("SetCreatedAt returned nil")
	}

	result = chat.SetUpdatedAt("2024-01-01 00:00:00")
	if result == nil {
		t.Error("SetUpdatedAt returned nil")
	}

	result = chat.SetSoftDeletedAt("2024-01-01 00:00:00")
	if result == nil {
		t.Error("SetSoftDeletedAt returned nil")
	}
}

func TestChatMarkAsNotDirty(t *testing.T) {
	chat := NewChat()

	// This should not panic
	chat.MarkAsNotDirty()
}

func TestChatChaining(t *testing.T) {
	chat := NewChat()

	// Test method chaining
	result := chat.SetID("test-id").
		SetOwnerID("owner-id").
		SetStatus(CHAT_STATUS_INACTIVE).
		SetTitle("Test Title").
		SetMemo("Test Memo").
		SetCreatedAt("2024-01-01 00:00:00").
		SetUpdatedAt("2024-01-01 00:00:00").
		SetSoftDeletedAt("")

	if result == nil {
		t.Fatal("Chaining returned nil")
	}

	if chat.ID() != "test-id" {
		t.Errorf("Chaining failed: expected ID test-id, got %s", chat.ID())
	}

	if chat.OwnerID() != "owner-id" {
		t.Errorf("Chaining failed: expected OwnerID owner-id, got %s", chat.OwnerID())
	}

	if chat.Status() != CHAT_STATUS_INACTIVE {
		t.Errorf("Chaining failed: expected status %s, got %s", CHAT_STATUS_INACTIVE, chat.Status())
	}

	if chat.Title() != "Test Title" {
		t.Errorf("Chaining failed: expected title 'Test Title', got %s", chat.Title())
	}

	if chat.Memo() != "Test Memo" {
		t.Errorf("Chaining failed: expected memo 'Test Memo', got %s", chat.Memo())
	}
}
