package chatstore

import (
	"testing"

	"github.com/dracory/sb"
)

func TestNewMessage(t *testing.T) {
	message := NewMessage()

	if message == nil {
		t.Fatal("NewMessage returned nil")
	}

	if message.ID() == "" {
		t.Error("Expected ID to be set")
	}

	if message.Status() != CHAT_MESSAGE_STATUS_ACTIVE {
		t.Errorf("Expected status %s, got %s", CHAT_MESSAGE_STATUS_ACTIVE, message.Status())
	}

	if message.ChatID() != "" {
		t.Error("Expected empty ChatID")
	}

	if message.SenderID() != "" {
		t.Error("Expected empty SenderID")
	}

	if message.RecipientID() != "" {
		t.Error("Expected empty RecipientID")
	}

	if message.Text() != "" {
		t.Error("Expected empty Text")
	}

	if message.Memo() != "" {
		t.Error("Expected empty Memo")
	}

	if message.CreatedAt() == "" {
		t.Error("Expected CreatedAt to be set")
	}

	if message.UpdatedAt() == "" {
		t.Error("Expected UpdatedAt to be set")
	}

	if message.SoftDeletedAt() == "" {
		t.Error("Expected SoftDeletedAt to be set")
	}

	metas, err := message.Metas()
	if err != nil {
		t.Fatalf("Failed to get metas: %v", err)
	}

	if len(metas) != 0 {
		t.Errorf("Expected empty metas, got %d", len(metas))
	}
}

func TestNewMessageFromExistingData(t *testing.T) {
	data := map[string]string{
		COLUMN_ID:              "test-id",
		COLUMN_CHAT_ID:         "chat-id",
		COLUMN_SENDER_ID:       "sender-id",
		COLUMN_RECIPIENT_ID:    "recipient-id",
		COLUMN_STATUS:          CHAT_MESSAGE_STATUS_INACTIVE,
		COLUMN_TEXT:            "Test message",
		COLUMN_MEMO:            "Test memo",
		COLUMN_CREATED_AT:      "2024-01-01 00:00:00",
		COLUMN_UPDATED_AT:      "2024-01-02 00:00:00",
		COLUMN_SOFT_DELETED_AT: "2024-01-03 00:00:00",
		COLUMN_METAS:           `{"key1":"value1","key2":"value2"}`,
	}

	message := NewMessageFromExistingData(data)

	if message == nil {
		t.Fatal("NewMessageFromExistingData returned nil")
	}

	if message.ID() != "test-id" {
		t.Errorf("Expected ID test-id, got %s", message.ID())
	}

	if message.ChatID() != "chat-id" {
		t.Errorf("Expected ChatID chat-id, got %s", message.ChatID())
	}

	if message.SenderID() != "sender-id" {
		t.Errorf("Expected SenderID sender-id, got %s", message.SenderID())
	}

	if message.RecipientID() != "recipient-id" {
		t.Errorf("Expected RecipientID recipient-id, got %s", message.RecipientID())
	}

	if message.Status() != CHAT_MESSAGE_STATUS_INACTIVE {
		t.Errorf("Expected status %s, got %s", CHAT_MESSAGE_STATUS_INACTIVE, message.Status())
	}

	if message.Text() != "Test message" {
		t.Errorf("Expected text 'Test message', got %s", message.Text())
	}

	if message.Memo() != "Test memo" {
		t.Errorf("Expected memo 'Test memo', got %s", message.Memo())
	}

	if message.CreatedAt() != "2024-01-01 00:00:00" {
		t.Errorf("Expected CreatedAt 2024-01-01 00:00:00, got %s", message.CreatedAt())
	}

	if message.UpdatedAt() != "2024-01-02 00:00:00" {
		t.Errorf("Expected UpdatedAt 2024-01-02 00:00:00, got %s", message.UpdatedAt())
	}

	if message.SoftDeletedAt() != "2024-01-03 00:00:00" {
		t.Errorf("Expected SoftDeletedAt 2024-01-03 00:00:00, got %s", message.SoftDeletedAt())
	}
}

func TestMessageID(t *testing.T) {
	message := NewMessage()

	testID := "test-id-123"
	message.SetID(testID)

	if message.ID() != testID {
		t.Errorf("Expected ID %s, got %s", testID, message.ID())
	}
}

func TestMessageChatID(t *testing.T) {
	message := NewMessage()

	testChatID := "chat-id-456"
	message.SetChatID(testChatID)

	if message.ChatID() != testChatID {
		t.Errorf("Expected ChatID %s, got %s", testChatID, message.ChatID())
	}
}

func TestMessageSenderID(t *testing.T) {
	message := NewMessage()

	testSenderID := "sender-id-789"
	message.SetSenderID(testSenderID)

	if message.SenderID() != testSenderID {
		t.Errorf("Expected SenderID %s, got %s", testSenderID, message.SenderID())
	}
}

func TestMessageRecipientID(t *testing.T) {
	message := NewMessage()

	testRecipientID := "recipient-id-012"
	message.SetRecipientID(testRecipientID)

	if message.RecipientID() != testRecipientID {
		t.Errorf("Expected RecipientID %s, got %s", testRecipientID, message.RecipientID())
	}
}

func TestMessageStatus(t *testing.T) {
	message := NewMessage()

	if message.Status() != CHAT_MESSAGE_STATUS_ACTIVE {
		t.Errorf("Expected initial status %s, got %s", CHAT_MESSAGE_STATUS_ACTIVE, message.Status())
	}

	message.SetStatus(CHAT_MESSAGE_STATUS_INACTIVE)

	if message.Status() != CHAT_MESSAGE_STATUS_INACTIVE {
		t.Errorf("Expected status %s, got %s", CHAT_MESSAGE_STATUS_INACTIVE, message.Status())
	}
}

func TestMessageText(t *testing.T) {
	message := NewMessage()

	testText := "Test message text"
	message.SetText(testText)

	if message.Text() != testText {
		t.Errorf("Expected text %s, got %s", testText, message.Text())
	}
}

func TestMessageMemo(t *testing.T) {
	message := NewMessage()

	testMemo := "Test memo content"
	message.SetMemo(testMemo)

	if message.Memo() != testMemo {
		t.Errorf("Expected memo %s, got %s", testMemo, message.Memo())
	}
}

func TestMessageCreatedAt(t *testing.T) {
	message := NewMessage()

	testTime := "2024-03-15 10:30:00"
	message.SetCreatedAt(testTime)

	if message.CreatedAt() != testTime {
		t.Errorf("Expected CreatedAt %s, got %s", testTime, message.CreatedAt())
	}

	carbonTime := message.CreatedAtCarbon()
	if carbonTime == nil {
		t.Fatal("CreatedAtCarbon returned nil")
	}

	if carbonTime.Format("Y-m-d H:i:s") != testTime {
		t.Errorf("Expected Carbon time %s, got %s", testTime, carbonTime.Format("Y-m-d H:i:s"))
	}
}

func TestMessageUpdatedAt(t *testing.T) {
	message := NewMessage()

	testTime := "2024-03-15 10:30:00"
	message.SetUpdatedAt(testTime)

	if message.UpdatedAt() != testTime {
		t.Errorf("Expected UpdatedAt %s, got %s", testTime, message.UpdatedAt())
	}

	carbonTime := message.UpdatedAtCarbon()
	if carbonTime == nil {
		t.Fatal("UpdatedAtCarbon returned nil")
	}

	if carbonTime.Format("Y-m-d H:i:s") != testTime {
		t.Errorf("Expected Carbon time %s, got %s", testTime, carbonTime.Format("Y-m-d H:i:s"))
	}
}

func TestMessageSoftDeletedAt(t *testing.T) {
	message := NewMessage()

	// Initially has MAX_DATETIME set, so IsSoftDeleted returns false
	if message.IsSoftDeleted() {
		t.Error("Expected IsSoftDeleted to be false initially (has MAX_DATETIME)")
	}

	testTime := "2024-03-15 10:30:00"
	message.SetSoftDeletedAt(testTime)

	if message.SoftDeletedAt() != testTime {
		t.Errorf("Expected SoftDeletedAt %s, got %s", testTime, message.SoftDeletedAt())
	}

	if !message.IsSoftDeleted() {
		t.Error("Expected IsSoftDeleted to be true after setting SoftDeletedAt to non-MAX_DATETIME")
	}

	carbonTime := message.SoftDeletedAtCarbon()
	if carbonTime == nil {
		t.Fatal("SoftDeletedAtCarbon returned nil")
	}

	if carbonTime.Format("Y-m-d H:i:s") != testTime {
		t.Errorf("Expected Carbon time %s, got %s", testTime, carbonTime.Format("Y-m-d H:i:s"))
	}

	// Test setting to MAX_DATETIME (not soft deleted)
	message.SetSoftDeletedAt(sb.MAX_DATETIME)
	if message.IsSoftDeleted() {
		t.Error("Expected IsSoftDeleted to be false when SoftDeletedAt is MAX_DATETIME")
	}
}

func TestMessageMetas(t *testing.T) {
	message := NewMessage()

	// Test setting metas
	testMetas := map[string]string{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
	}

	err := message.SetMetas(testMetas)
	if err != nil {
		t.Fatalf("Failed to set metas: %v", err)
	}

	// Test getting metas
	metas, err := message.Metas()
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

func TestMessageMeta(t *testing.T) {
	message := NewMessage()

	// Set metas
	testMetas := map[string]string{
		"key1": "value1",
		"key2": "value2",
	}

	err := message.SetMetas(testMetas)
	if err != nil {
		t.Fatalf("Failed to set metas: %v", err)
	}

	// Test getting single meta
	value, err := message.Meta("key1")
	if err != nil {
		t.Fatalf("Failed to get meta: %v", err)
	}

	if value != "value1" {
		t.Errorf("Expected value1, got %s", value)
	}

	// Test getting non-existent meta
	value, err = message.Meta("nonexistent")
	if err != nil {
		t.Fatalf("Failed to get non-existent meta: %v", err)
	}

	if value != "" {
		t.Errorf("Expected empty string for non-existent meta, got %s", value)
	}
}

func TestMessageSetMeta(t *testing.T) {
	message := NewMessage()

	// Set single meta
	err := message.SetMeta("key1", "value1")
	if err != nil {
		t.Fatalf("Failed to set meta: %v", err)
	}

	// Verify it was set
	value, err := message.Meta("key1")
	if err != nil {
		t.Fatalf("Failed to get meta: %v", err)
	}

	if value != "value1" {
		t.Errorf("Expected value1, got %s", value)
	}
}

func TestMessageUpsertMetas(t *testing.T) {
	message := NewMessage()

	// Set initial metas
	initialMetas := map[string]string{
		"key1": "value1",
		"key2": "value2",
	}

	err := message.SetMetas(initialMetas)
	if err != nil {
		t.Fatalf("Failed to set initial metas: %v", err)
	}

	// Upsert new metas (update existing, add new)
	upsertMetas := map[string]string{
		"key1": "updated-value1",
		"key3": "value3",
	}

	err = message.UpsertMetas(upsertMetas)
	if err != nil {
		t.Fatalf("Failed to upsert metas: %v", err)
	}

	// Verify results
	metas, err := message.Metas()
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

func TestMessageMetasEmpty(t *testing.T) {
	message := NewMessage()

	// Get metas when empty
	metas, err := message.Metas()
	if err != nil {
		t.Fatalf("Failed to get metas: %v", err)
	}

	if len(metas) != 0 {
		t.Errorf("Expected empty metas, got %d", len(metas))
	}
}

func TestMessageSetMetasEmpty(t *testing.T) {
	message := NewMessage()

	// Set some metas first
	testMetas := map[string]string{
		"key1": "value1",
	}

	err := message.SetMetas(testMetas)
	if err != nil {
		t.Fatalf("Failed to set metas: %v", err)
	}

	// Set empty metas
	err = message.SetMetas(map[string]string{})
	if err != nil {
		t.Fatalf("Failed to set empty metas: %v", err)
	}

	// Verify they're empty
	metas, err := message.Metas()
	if err != nil {
		t.Fatalf("Failed to get metas: %v", err)
	}

	if len(metas) != 0 {
		t.Errorf("Expected empty metas after setting empty, got %d", len(metas))
	}
}

func TestMessageSettersReturnMessageInterface(t *testing.T) {
	message := NewMessage()

	// Test that all setters return MessageInterface for chaining
	result := message.SetChatID("chat-id")
	if result == nil {
		t.Error("SetChatID returned nil")
	}

	result = message.SetSenderID("sender-id")
	if result == nil {
		t.Error("SetSenderID returned nil")
	}

	result = message.SetRecipientID("recipient-id")
	if result == nil {
		t.Error("SetRecipientID returned nil")
	}

	result = message.SetStatus(CHAT_MESSAGE_STATUS_INACTIVE)
	if result == nil {
		t.Error("SetStatus returned nil")
	}

	result = message.SetText("text")
	if result == nil {
		t.Error("SetText returned nil")
	}

	result = message.SetMemo("memo")
	if result == nil {
		t.Error("SetMemo returned nil")
	}

	result = message.SetCreatedAt("2024-01-01 00:00:00")
	if result == nil {
		t.Error("SetCreatedAt returned nil")
	}

	result = message.SetUpdatedAt("2024-01-01 00:00:00")
	if result == nil {
		t.Error("SetUpdatedAt returned nil")
	}

	result = message.SetSoftDeletedAt("2024-01-01 00:00:00")
	if result == nil {
		t.Error("SetSoftDeletedAt returned nil")
	}
}

func TestMessageChaining(t *testing.T) {
	message := NewMessage()

	// Set ID separately since it doesn't return MessageInterface
	message.SetID("test-id")

	// Test method chaining for other setters
	result := message.SetChatID("chat-id").
		SetSenderID("sender-id").
		SetRecipientID("recipient-id").
		SetStatus(CHAT_MESSAGE_STATUS_INACTIVE).
		SetText("Test Text").
		SetMemo("Test Memo").
		SetCreatedAt("2024-01-01 00:00:00").
		SetUpdatedAt("2024-01-01 00:00:00").
		SetSoftDeletedAt("")

	if result == nil {
		t.Fatal("Chaining returned nil")
	}

	if message.ID() != "test-id" {
		t.Errorf("Chaining failed: expected ID test-id, got %s", message.ID())
	}

	if message.ChatID() != "chat-id" {
		t.Errorf("Chaining failed: expected ChatID chat-id, got %s", message.ChatID())
	}

	if message.SenderID() != "sender-id" {
		t.Errorf("Chaining failed: expected SenderID sender-id, got %s", message.SenderID())
	}

	if message.RecipientID() != "recipient-id" {
		t.Errorf("Chaining failed: expected RecipientID recipient-id, got %s", message.RecipientID())
	}

	if message.Status() != CHAT_MESSAGE_STATUS_INACTIVE {
		t.Errorf("Chaining failed: expected status %s, got %s", CHAT_MESSAGE_STATUS_INACTIVE, message.Status())
	}

	if message.Text() != "Test Text" {
		t.Errorf("Chaining failed: expected text 'Test Text', got %s", message.Text())
	}

	if message.Memo() != "Test Memo" {
		t.Errorf("Chaining failed: expected memo 'Test Memo', got %s", message.Memo())
	}
}
