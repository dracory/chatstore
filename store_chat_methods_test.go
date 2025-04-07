package chatstore

import (
	"testing"

	"github.com/gouniverse/sb"
	_ "modernc.org/sqlite"
)

func TestStore_ChatCount(t *testing.T) {
	store, err := initStore(":memory:")

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	// Create multiple chats
	chat1 := NewChat().
		SetOwnerID(testUser_O1).
		SetStatus(CHAT_STATUS_ACTIVE)

	chat2 := NewChat().
		SetOwnerID(testUser_O1).
		SetStatus(CHAT_STATUS_INACTIVE)

	chat3 := NewChat().
		SetOwnerID(testUser_O1).
		SetStatus(CHAT_STATUS_INACTIVE)

	err = store.ChatCreate(chat1)
	if err != nil {
		t.Fatal("unexpected error creating chat1:", err)
	}

	err = store.ChatCreate(chat2)
	if err != nil {
		t.Fatal("unexpected error creating chat2:", err)
	}

	err = store.ChatCreate(chat3)
	if err != nil {
		t.Fatal("unexpected error creating chat3:", err)
	}

	// Test counting all chats
	allCount, err := store.ChatCount(ChatQuery())
	if err != nil {
		t.Fatal("unexpected error counting all chats:", err)
	}

	if allCount != 3 {
		t.Fatalf("Expected count of 3 chats, got %d", allCount)
	}

	// Test counting by monitor ID
	chatCount, err := store.ChatCount(ChatQuery().SetOwnerID(testUser_O1))
	if err != nil {
		t.Fatal("unexpected error counting monitor chats:", err)
	}

	if chatCount != 3 {
		t.Fatalf("Expected count of 3 chats, got %d", chatCount)
	}

	// Test counting by status
	activeCount, err := store.ChatCount(ChatQuery().SetStatus(CHAT_STATUS_ACTIVE))
	if err != nil {
		t.Fatal("unexpected error counting active chats:", err)
	}

	if activeCount != 1 {
		t.Fatalf("Expected count of 1 active chat, got %d", activeCount)
	}
}

func TestStore_ChatCreate(t *testing.T) {
	store, err := initStore(":memory:")

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	chat := NewChat().
		SetOwnerID(testUser_O1).
		SetStatus(CHAT_STATUS_ACTIVE)

	err = store.ChatCreate(chat)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}
}

func TestStore_ChatCreateDuplicate(t *testing.T) {
	store, err := initStore(":memory:")

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	chat := NewChat().
		SetOwnerID(testUser_O1).
		SetStatus(CHAT_STATUS_ACTIVE)

	err = store.ChatCreate(chat)
	if err != nil {
		t.Fatal("unexpected error on first create:", err)
	}

	// Try to create the same chat again
	err = store.ChatCreate(chat)
	if err == nil {
		t.Fatal("expected error for duplicate chat, but got nil")
	}
}

func TestStore_ChatFindByID(t *testing.T) {
	store, err := initStore(":memory:")

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	chat := NewChat().
		SetOwnerID(testChat_O1).
		SetStatus(CHAT_STATUS_ACTIVE)

	err = chat.SetMetas(map[string]string{
		"severity":         "high",
		"affected_service": "web",
		"team":             "backend",
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	err = store.ChatCreate(chat)
	if err != nil {
		t.Error("unexpected error:", err)
	}

	chatFound, errFind := store.ChatFindByID(chat.ID())

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
	}

	if chatFound == nil {
		t.Fatal("Chat MUST NOT be nil")
	}

	if chatFound.ID() != chat.ID() {
		t.Fatal("IDs do not match")
	}

	if chatFound.Status() != chat.Status() {
		t.Fatal("Statuses do not match")
	}

	if chatFound.Status() != CHAT_STATUS_ACTIVE {
		t.Fatal("Statuses do not match")
	}

	severity, err := chat.Meta("severity")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	foundSeverity, err := chatFound.Meta("severity")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if foundSeverity != severity {
		t.Fatal("Metas do not match")
	}

	affectedService, err := chat.Meta("affected_service")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	foundAffectedService, err := chatFound.Meta("affected_service")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if foundAffectedService != affectedService {
		t.Fatal("Metas do not match")
	}

	team, err := chat.Meta("team")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	foundTeam, err := chatFound.Meta("team")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if foundTeam != team {
		t.Fatal("Metas do not match")
	}
}

func TestStore_ChatFindByIDNotFound(t *testing.T) {
	store, err := initStore(":memory:")

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	chatFound, errFind := store.ChatFindByID("non-existent-id")

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
	}

	if chatFound != nil {
		t.Fatal("Chat MUST be nil for non-existent ID")
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

func TestStore_ChatDelete(t *testing.T) {
	store, err := initStore(":memory:")

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	chat := NewChat().
		SetOwnerID(testUser_O1).
		SetStatus(CHAT_STATUS_ACTIVE)

	err = store.ChatCreate(chat)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	// Delete the chat
	err = store.ChatDelete(chat)
	if err != nil {
		t.Fatal("unexpected error on delete:", err)
	}

	// Verify the chat is deleted
	deletedChat, err := store.ChatFindByID(chat.ID())
	if err != nil {
		t.Fatal("unexpected error finding deleted chat:", err)
	}

	if deletedChat != nil {
		t.Fatal("Chat should be nil after deletion")
	}
}

func TestStore_ChatDeleteByID(t *testing.T) {
	store, err := initStore(":memory:")

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	chat := NewChat().
		SetOwnerID(testUser_O1).
		SetStatus(CHAT_STATUS_ACTIVE)

	err = store.ChatCreate(chat)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	// Delete the chat by ID
	err = store.ChatDeleteByID(chat.ID())
	if err != nil {
		t.Fatal("unexpected error on delete by ID:", err)
	}

	// Verify the chat is deleted
	deletedChat, err := store.ChatFindByID(chat.ID())
	if err != nil {
		t.Fatal("unexpected error finding deleted chat:", err)
	}

	if deletedChat != nil {
		t.Fatal("Chat should be nil after deletion by ID")
	}
}

func TestStore_ChatList(t *testing.T) {
	store, err := initStore(":memory:")

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	// Create multiple chats
	chat1 := NewChat().
		SetOwnerID(testUser_O1).
		SetStatus(CHAT_STATUS_ACTIVE)

	chat2 := NewChat().
		SetOwnerID(testUser_O1).
		SetStatus(CHAT_STATUS_INACTIVE)

	err = store.ChatCreate(chat1)
	if err != nil {
		t.Fatal("unexpected error creating chat1:", err)
	}

	err = store.ChatCreate(chat2)
	if err != nil {
		t.Fatal("unexpected error creating chat2:", err)
	}

	// Test listing all chats
	allChats, err := store.ChatList(ChatQuery())
	if err != nil {
		t.Fatal("unexpected error listing all chats:", err)
	}

	if len(allChats) != 2 {
		t.Fatalf("Expected 2 chats, got %d", len(allChats))
	}

	// Test filtering by owner ID
	ownerChats, err := store.ChatList(ChatQuery().SetOwnerID(testUser_O1))
	if err != nil {
		t.Fatal("unexpected error listing owner chats:", err)
	}

	if len(ownerChats) != 2 {
		t.Fatalf("Expected 2 chats for MONITOR_01, got %d", len(ownerChats))
	}

	// Test filtering by status
	activeChats, err := store.ChatList(ChatQuery().SetStatus(CHAT_STATUS_ACTIVE))
	if err != nil {
		t.Fatal("unexpected error listing active chats:", err)
	}

	if len(activeChats) != 1 {
		t.Fatalf("Expected 1 active chat, got %d", len(activeChats))
	}

	// Test limit and offset
	limitedChats, err := store.ChatList(ChatQuery().SetLimit(1))
	if err != nil {
		t.Fatal("unexpected error listing limited chats:", err)
	}

	if len(limitedChats) != 1 {
		t.Fatalf("Expected 1 chat with limit, got %d", len(limitedChats))
	}

	offsetChats, err := store.ChatList(ChatQuery().SetOffset(1).SetLimit(2))
	if err != nil {
		t.Fatal("unexpected error listing offset chats:", err)
	}

	if len(offsetChats) != 1 {
		t.Fatalf("Expected 1 chat with offset, got %d", len(offsetChats))
	}
}

func TestStore_ChatSoftDelete(t *testing.T) {
	store, err := initStore(":memory:")

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	chat := NewChat().
		SetOwnerID(testUser_O1).
		SetStatus(CHAT_STATUS_ACTIVE)

	err = store.ChatCreate(chat)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	// Soft delete the chat
	err = store.ChatSoftDelete(chat)
	if err != nil {
		t.Fatal("unexpected error on soft delete:", err)
	}

	// Verify the chat is soft deleted (not found by default)
	softDeletedChat, err := store.ChatFindByID(chat.ID())
	if err != nil {
		t.Fatal("unexpected error finding soft deleted chat:", err)
	}

	if softDeletedChat != nil {
		t.Fatal("Chat should not be found after soft deletion")
	}

	// Verify the chat can be found when including soft deleted
	query := ChatQuery().
		SetWithSoftDeleted(true).
		SetID(chat.ID()).
		SetLimit(1)

	chatFindWithDeleted, err := store.ChatList(query)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if len(chatFindWithDeleted) == 0 {
		t.Fatal("Chat should be found when including soft deleted")
	}

	if chatFindWithDeleted[0].SoftDeletedAt() == sb.MAX_DATETIME {
		t.Fatal("Chat should be soft deleted, but SoftDeletedAt is MAX_DATETIME:", chatFindWithDeleted[0].SoftDeletedAt())
	}

	if !chatFindWithDeleted[0].IsSoftDeleted() {
		t.Fatal("Chat should be marked as soft deleted")
	}
}

func TestStore_ChatSoftDeleteByID(t *testing.T) {
	store, err := initStore(":memory:")

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	chat := NewChat().
		SetOwnerID(testUser_O1).
		SetStatus(CHAT_STATUS_ACTIVE)

	err = store.ChatCreate(chat)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	// Soft delete the chat by ID
	err = store.ChatSoftDeleteByID(chat.ID())
	if err != nil {
		t.Fatal("unexpected error on soft delete by ID:", err)
	}

	// Verify the chat is soft deleted (not found by default)
	softDeletedChat, err := store.ChatFindByID(chat.ID())
	if err != nil {
		t.Fatal("unexpected error finding soft deleted chat:", err)
	}

	if softDeletedChat != nil {
		t.Fatal("Chat should not be found after soft deletion by ID")
	}

	// Verify the chat can be found when including soft deleted
	query := ChatQuery().
		SetWithSoftDeleted(true).
		SetID(chat.ID()).
		SetLimit(1)

	chatFindWithDeleted, err := store.ChatList(query)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if len(chatFindWithDeleted) == 0 {
		t.Fatal("Chat should be found when including soft deleted")
	}

	if chatFindWithDeleted[0].SoftDeletedAt() == sb.MAX_DATETIME {
		t.Fatal("Chat should be soft deleted, but SoftDeletedAt is MAX_DATETIME:", chatFindWithDeleted[0].SoftDeletedAt())
	}

	if !chatFindWithDeleted[0].IsSoftDeleted() {
		t.Fatal("Chat should be marked as soft deleted")
	}
}

func TestStore_ChatUpdate(t *testing.T) {
	store, err := initStore(":memory:")

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	chat := NewChat().
		SetOwnerID(testUser_O1).
		SetStatus(CHAT_STATUS_ACTIVE)

	err = store.ChatCreate(chat)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	// Update the chat
	chat.SetStatus(CHAT_STATUS_INACTIVE).
		SetMemo("Resolved by ops team")

	err = store.ChatUpdate(chat)
	if err != nil {
		t.Fatal("unexpected error on update:", err)
	}

	// Verify the update
	updatedChat, err := store.ChatFindByID(chat.ID())
	if err != nil {
		t.Fatal("unexpected error finding updated chat:", err)
	}

	if updatedChat.Status() != CHAT_STATUS_INACTIVE {
		t.Fatalf("Status not updated. Expected %s, got %s", CHAT_STATUS_INACTIVE, updatedChat.Status())
	}

	if updatedChat.Memo() != "Resolved by ops team" {
		t.Fatalf("Memo not updated. Expected 'Resolved by ops team', got '%s'", updatedChat.Memo())
	}
}
