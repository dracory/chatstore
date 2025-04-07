# ChatStore <a href="https://gitpod.io/#https://github.com/dracory/chatstore" style="float:right:"><img src="https://gitpod.io/button/open-in-gitpod.svg" alt="Open in Gitpod" loading="lazy"></a>

[![Tests Status](https://github.com/dracory/chatstore/actions/workflows/tests.yml/badge.svg?branch=main)](https://github.com/dracory/chatstore/actions/workflows/tests.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/dracory/chatstore)](https://goreportcard.com/report/github.com/dracory/chatstore)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/dracory/chatstore)](https://pkg.go.dev/github.com/dracory/chatstore)

A Go library for storing and retrieving chat messages.

## License

This project is dual-licensed under the following terms:

-   For non-commercial use, you may choose either the GNU Affero General Public License v3.0 (AGPLv3) *or* a separate commercial license (see below). You can find a copy of the AGPLv3 at: https://www.gnu.org/licenses/agpl-3.0.txt

-   For commercial use, a separate commercial license is required. Commercial licenses are available for various use cases. Please contact me via my [contact page](https://lesichkov.co.uk/contact) to obtain a commercial license.

## Installation

```
go get github.com/dracory/chatstore
```


## Usage Examples

Here are some examples demonstrating how to use the `chatstore` library.

### Example 1: Creating a Chat Store

This example shows how to create a chat store.

```go
// Initialize the database connection.
db, err := sql.Open("sqlite3", "./chatstore.db") // Replace with your database details
if err != nil {
    log.Fatalf("Failed to open database: %v", err)
}
defer db.Close()

// Create the store with database options.
store, err := chatstore.CreateStore(chatstore.NewStoreOptions{
    DB:                 db,
    TableChatName:      "chat_table",
    TableMessageName:   "message_table",
    AutomigrateEnabled: true,
})
if err != nil {
    log.Fatalf("Failed to create store: %v", err)
}
```

### Example 2. Creating a Chat

This example shows how to create a chat.

```go
chat := chatstore.NewChat().
		SetName("Test Chat").
		SetOwnerID(testUser_O1).
		SetStatus(CHAT_STATUS_ACTIVE)

err = store.ChatCreate(chat)
if err != nil {
    log.Fatalf("Failed to create chat: %v", err)
}

fmt.Println("Chat created successfully!")
```

### Example 3: Creating a Chat Message

This example shows how to create and store a single chat message.

```go
message := chatstore.NewMessage().
		SetChatID(chat.ID()).
		SetSenderID(user1.ID()).
		SetRecipientID(user2.ID()).
		SetText("Message 1")

err = store.MessageCreate(message)
if err != nil {
    log.Fatalf("Failed to create message: %v", err)
}

fmt.Println("Message stored successfully!")
```