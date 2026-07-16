# Rooms to Chat

A web app where you can chat with people online and make friends.
https://go-chat.duckdns.org/

## Motivation

I wanted hands-on experience combining Go's networking capabilities with real database CRUD operations, and a chat app felt like the perfect vehicle for that.

## Quick Start

The app is very self-explanatory; you browse the list of public rooms and join one. Or you can create your own where it says "Create Room" by typing into the text box.
After you join a room, simply write a message in the text box below that says "Message" and send!


## Usage

- Real-time messaging with people online.
- Creating rooms to chat with multiple people.
- Right-click on any username in a chat room to add them as a friend and start a private 1-to-1 conversation.

## Tech Stack

- Go
- PostgreSQL
- WebSocket
- HTML
- Hosted on AWS EC2
- UUID for user authentication


## 🤝 Contributing

### Clone the repo
```bash
git clone https://github.com/BetoDev25/chatroom-project@latest
cd chatroom-project
```

### Build the compiled binary

```bash
go build
```

### Submit a pull request

If you'd like to contribute, please fork the repository and open a pull request to the `main` branch.


## Coming Soon
- Private rooms with passwords
- Blocking users
- User profiles with bios and avatars