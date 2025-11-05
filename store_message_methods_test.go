package chatstore

import (
	"testing"

	"github.com/dracory/sb"
	// _ "modernc.org/sqlite"
)

func TestStore_MessageCount(t *testing.T) {
	store, err := initStore(":memory:")

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	// Create multiple messages
	message1 := NewMessage().
		SetChatID(testChat_O1).
		SetSenderID(testUser_O1).
		SetRecipientID(testUser_O2).
		SetText("Message 1")

	message2 := NewMessage().
		SetStatus(MESSAGE_STATUS_INACTIVE).
		SetChatID(testChat_O1).
		SetSenderID(testUser_O1).
		SetRecipientID(testUser_O2).
		SetText("Message 2")

	message3 := NewMessage().
		SetChatID(testChat_O1).
		SetSenderID(testUser_O1).
		SetRecipientID(testUser_O2).
		SetText("Message 3")

	err = store.MessageCreate(message1)
	if err != nil {
		t.Fatal("unexpected error creating message1:", err)
	}

	err = store.MessageCreate(message2)
	if err != nil {
		t.Fatal("unexpected error creating message2:", err)
	}

	err = store.MessageCreate(message3)
	if err != nil {
		t.Fatal("unexpected error creating message3:", err)
	}

	// Test counting all messages
	allCount, err := store.MessageCount(MessageQuery())
	if err != nil {
		t.Fatal("unexpected error counting all messages:", err)
	}

	if allCount != 3 {
		t.Fatalf("Expected count of 3 chats, got %d", allCount)
	}

	// Test counting by monitor ID
	messageCount, err := store.MessageCount(MessageQuery().SetChatID(testChat_O1))
	if err != nil {
		t.Fatal("unexpected error counting monitor messages:", err)
	}

	if messageCount != 3 {
		t.Fatalf("Expected count of 3 messages, got %d", messageCount)
	}

	// Test counting by status
	activeCount, err := store.MessageCount(MessageQuery().SetStatus(MESSAGE_STATUS_ACTIVE))
	if err != nil {
		t.Fatal("unexpected error counting active messages:", err)
	}

	if activeCount != 2 {
		t.Fatalf("Expected count of 2 active messages, got %d", activeCount)
	}
}

func TestStore_MessageCreate(t *testing.T) {
	store, err := initStore(":memory:")

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	message := NewMessage().
		SetChatID(testChat_O1).
		SetSenderID(testUser_O1).
		SetRecipientID(testUser_O2).
		SetText("Message 1")

	err = store.MessageCreate(message)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}
}

func TestStore_MessageCreateDuplicate(t *testing.T) {
	store, err := initStore(":memory:")

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	message := NewMessage().
		SetChatID(testChat_O1).
		SetSenderID(testUser_O1).
		SetRecipientID(testUser_O2).
		SetText("Message 1")

	err = store.MessageCreate(message)
	if err != nil {
		t.Fatal("unexpected error on first create:", err)
	}

	// Try to create the same message again
	err = store.MessageCreate(message)
	if err == nil {
		t.Fatal("expected error for duplicate message, but got nil")
	}
}

func TestStore_MessageFindByID(t *testing.T) {
	store, err := initStore(":memory:")

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	message := NewMessage().
		SetChatID(testChat_O1).
		SetSenderID(testUser_O1).
		SetRecipientID(testUser_O2).
		SetText("Message 1")

	err = message.SetMetas(map[string]string{
		"severity":         "high",
		"affected_service": "web",
		"team":             "backend",
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	err = store.MessageCreate(message)
	if err != nil {
		t.Error("unexpected error:", err)
	}

	messageFound, errFind := store.MessageFindByID(message.ID())

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
	}

	if messageFound == nil {
		t.Fatal("Message MUST NOT be nil")
	}

	if messageFound.ID() != message.ID() {
		t.Fatal("IDs do not match")
	}

	if messageFound.Status() != message.Status() {
		t.Fatal("Statuses do not match")
	}

	if messageFound.Status() != MESSAGE_STATUS_ACTIVE {
		t.Fatal("Statuses do not match")
	}

	severity, err := message.Meta("severity")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	foundSeverity, err := messageFound.Meta("severity")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if foundSeverity != severity {
		t.Fatal("Metas do not match")
	}

	affectedService, err := message.Meta("affected_service")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	foundAffectedService, err := messageFound.Meta("affected_service")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if foundAffectedService != affectedService {
		t.Fatal("Metas do not match")
	}

	team, err := message.Meta("team")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	foundTeam, err := messageFound.Meta("team")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if foundTeam != team {
		t.Fatal("Metas do not match")
	}
}

func TestStore_MessageFindByIDNotFound(t *testing.T) {
	store, err := initStore(":memory:")

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	messageFound, errFind := store.MessageFindByID("non-existent-id")

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
	}

	if messageFound != nil {
		t.Fatal("Message MUST be nil for non-existent ID")
	}
}

// func TestStoreChatUpdateNonExistent(t *testing.T) {
// 	store, err := initStore(":memory:")

// 	if err != nil {
// 		t.Fatal("unexpected error:", err)
// 	}

// 	chat := NewChat().
// 		SetStatus(CHAT_STATUS_ACTIVE)

// 	// Try to update a non-existent chat
// 	err = store.ChatUpdate(chat)
// 	if err == nil {
// 		t.Fatal("expected error for updating non-existent chat, but got nil")
// 	}
// }

func TestStore_MessageDelete(t *testing.T) {
	store, err := initStore(":memory:")

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	message := NewMessage().
		SetChatID(testChat_O1).
		SetSenderID(testUser_O1).
		SetRecipientID(testUser_O2).
		SetText("Message 1")

	err = store.MessageCreate(message)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	// Delete the message
	err = store.MessageDelete(message)
	if err != nil {
		t.Fatal("unexpected error on delete:", err)
	}

	// Verify the message is deleted
	deletedMessage, err := store.MessageFindByID(message.ID())
	if err != nil {
		t.Fatal("unexpected error finding deleted message:", err)
	}

	if deletedMessage != nil {
		t.Fatal("Message should be nil after deletion")
	}
}

func TestStore_MessageDeleteByID(t *testing.T) {
	store, err := initStore(":memory:")

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	message := NewMessage().
		SetChatID(testChat_O1).
		SetSenderID(testUser_O1).
		SetRecipientID(testUser_O2).
		SetText("Message 1")

	err = store.MessageCreate(message)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	// Delete the message by ID
	err = store.MessageDeleteByID(message.ID())
	if err != nil {
		t.Fatal("unexpected error on delete by ID:", err)
	}

	// Verify the message is deleted
	deletedMessage, err := store.MessageFindByID(message.ID())
	if err != nil {
		t.Fatal("unexpected error finding deleted message:", err)
	}

	if deletedMessage != nil {
		t.Fatal("Message should be nil after deletion by ID")
	}
}

func TestStore_MessageList(t *testing.T) {
	store, err := initStore(":memory:")

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	// Create multiple messages
	message1 := NewMessage().
		SetChatID(testChat_O1).
		SetSenderID(testUser_O1).
		SetRecipientID(testUser_O2).
		SetText("Message 1")

	message2 := NewMessage().
		SetChatID(testChat_O1).
		SetSenderID(testUser_O1).
		SetRecipientID(testUser_O2).
		SetText("Message 2")

	err = store.MessageCreate(message1)
	if err != nil {
		t.Fatal("unexpected error creating message1:", err)
	}

	err = store.MessageCreate(message2)
	if err != nil {
		t.Fatal("unexpected error creating message2:", err)
	}

	// Test listing all messages
	allMessages, err := store.MessageList(MessageQuery())
	if err != nil {
		t.Fatal("unexpected error listing all messages:", err)
	}

	if len(allMessages) != 2 {
		t.Fatalf("Expected 2 messages, got %d", len(allMessages))
	}

	// Test filtering by chat ID
	chatMessages, err := store.MessageList(MessageQuery().SetChatID(testChat_O1))
	if err != nil {
		t.Fatal("unexpected error listing chat messages:", err)
	}

	if len(chatMessages) != 2 {
		t.Fatalf("Expected 2 messages for chat %s, got %d", testChat_O1, len(chatMessages))
	}

	// Test filtering by status
	activeMessages, err := store.MessageList(MessageQuery().SetStatus(MESSAGE_STATUS_ACTIVE))
	if err != nil {
		t.Fatal("unexpected error listing active messages:", err)
	}

	if len(activeMessages) != 2 {
		t.Fatalf("Expected 2 active messages, got %d", len(activeMessages))
	}

	// Test limit and offset
	limitedMessages, err := store.MessageList(MessageQuery().SetLimit(1))
	if err != nil {
		t.Fatal("unexpected error listing limited messages:", err)
	}

	if len(limitedMessages) != 1 {
		t.Fatalf("Expected 1 message with limit, got %d", len(limitedMessages))
	}

	offsetMessages, err := store.MessageList(MessageQuery().SetOffset(1).SetLimit(2))
	if err != nil {
		t.Fatal("unexpected error listing offset messages:", err)
	}

	if len(offsetMessages) != 1 {
		t.Fatalf("Expected 1 message with offset, got %d", len(offsetMessages))
	}
}

func TestStore_MessageSoftDelete(t *testing.T) {
	store, err := initStore(":memory:")

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	message := NewMessage().
		SetChatID(testChat_O1).
		SetSenderID(testUser_O1).
		SetRecipientID(testUser_O2).
		SetText("Message 1")

	err = store.MessageCreate(message)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	// Soft delete the message
	err = store.MessageSoftDelete(message)
	if err != nil {
		t.Fatal("unexpected error on soft delete:", err)
	}

	// Verify the message is soft deleted (not found by default)
	softDeletedMessage, err := store.MessageFindByID(message.ID())
	if err != nil {
		t.Fatal("unexpected error finding soft deleted message:", err)
	}

	if softDeletedMessage != nil {
		t.Fatal("Message should not be found after soft deletion")
	}

	// Verify the message can be found when including soft deleted
	query := MessageQuery().
		SetWithSoftDeleted(true).
		SetID(message.ID()).
		SetLimit(1)

	messageFindWithDeleted, err := store.MessageList(query)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if len(messageFindWithDeleted) == 0 {
		t.Fatal("Message should be found when including soft deleted")
	}

	if messageFindWithDeleted[0].SoftDeletedAt() == sb.MAX_DATETIME {
		t.Fatal("Message should be soft deleted, but SoftDeletedAt is MAX_DATETIME:", messageFindWithDeleted[0].SoftDeletedAt())
	}

	if !messageFindWithDeleted[0].IsSoftDeleted() {
		t.Fatal("Message should be marked as soft deleted")
	}
}

func TestStore_MessageSoftDeleteByID(t *testing.T) {
	store, err := initStore(":memory:")

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	message := NewMessage().
		SetChatID(testChat_O1).
		SetSenderID(testUser_O1).
		SetRecipientID(testUser_O2).
		SetText("Message 1")

	err = store.MessageCreate(message)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	// Soft delete the message by ID
	err = store.MessageSoftDeleteByID(message.ID())
	if err != nil {
		t.Fatal("unexpected error on soft delete by ID:", err)
	}

	// Verify the message is soft deleted (not found by default)
	softDeletedMessage, err := store.MessageFindByID(message.ID())
	if err != nil {
		t.Fatal("unexpected error finding soft deleted message:", err)
	}

	if softDeletedMessage != nil {
		t.Fatal("Message should not be found after soft deletion by ID")
	}

	// Verify the message can be found when including soft deleted
	query := MessageQuery().
		SetWithSoftDeleted(true).
		SetID(message.ID()).
		SetLimit(1)

	messageFindWithDeleted, err := store.MessageList(query)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if len(messageFindWithDeleted) == 0 {
		t.Fatal("Message should be found when including soft deleted")
	}

	if messageFindWithDeleted[0].SoftDeletedAt() == sb.MAX_DATETIME {
		t.Fatal("Message should be soft deleted, but SoftDeletedAt is MAX_DATETIME:", messageFindWithDeleted[0].SoftDeletedAt())
	}

	if !messageFindWithDeleted[0].IsSoftDeleted() {
		t.Fatal("Message should be marked as soft deleted")
	}
}

func TestStore_MessageUpdate(t *testing.T) {
	store, err := initStore(":memory:")

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	message := NewMessage().
		SetChatID(testChat_O1).
		SetSenderID(testUser_O1).
		SetRecipientID(testUser_O2).
		SetText("Message 1")

	err = store.MessageCreate(message)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	// Update the message
	message.SetText("Message 2")

	err = store.MessageUpdate(message)
	if err != nil {
		t.Fatal("unexpected error on update:", err)
	}

	// Verify the update
	updatedMessage, err := store.MessageFindByID(message.ID())
	if err != nil {
		t.Fatal("unexpected error finding updated message:", err)
	}

	if updatedMessage.Text() != "Message 2" {
		t.Fatalf("Text not updated. Expected 'Message 2', got '%s'", updatedMessage.Text())
	}
}
