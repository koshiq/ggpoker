# GGPoker

A Go implementation of a poker deck and card system.

## Features

- Card representation with suits (Spades, Hearts, Diamonds, Clubs)
- Unicode suit symbols (♠, ♥, ♦, ♣)
- Card validation (values 1-13)
- Clean Go package structure

## Project Structure

```
ggpoker/
├── main.go          # Main application entry point
├── deck/            # Card and deck package
│   └── deck.go     # Card implementation
├── bin/            # Binary output directory
├── go.mod          # Go module definition
└── Makefile        # Build and run commands
```

## Getting Started

### Prerequisites

- Go 1.24.6 or later

### Installation

1. Clone the repository:
```bash
git clone https://github.com/YOUR_USERNAME/ggpoker.git
cd ggpoker
```

2. Run the application:
```bash
make run
```

Or manually:
```bash
go build -o bin/ggpoker
./bin/ggpoker
```

### Running Tests

```bash
make test
```

## Usage

The current implementation creates a single card and displays it:

```go
card := deck.NewCard(deck.Spades, 1)
fmt.Println(card)
// Output: 1 of Spades ♠
```

## Development

- `make build` - Build the binary
- `make run` - Build and run the application
- `make test` - Run all tests

## License

This project is open source and available under the [MIT License](LICENSE).
