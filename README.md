# 🎰 GG Poker - WHOP Integration

A fully-featured P2P Texas Hold'em poker application with WHOP integration for user management, subscription handling, and payment processing.

## ✨ Features

- **Complete Texas Hold'em Implementation**: Full poker game logic with hand evaluation
- **P2P Networking**: Decentralized peer-to-peer gameplay
- **WHOP Integration**: User authentication, subscription management, and payments
- **Modern Web UI**: Beautiful, responsive poker table interface
- **Real-time Updates**: WebSocket-based game state synchronization
- **Multi-player Support**: Up to 6 players per table
- **Professional Design**: Casino-quality visual experience

## 🏗️ Architecture

```
ggpoker/
├── deck/           # Card and deck management
├── p2p/           # P2P networking and game logic
├── whop/          # WHOP API integration
├── web/           # React-based frontend
├── config.yaml    # Application configuration
└── main.go        # Main application entry point
```

## 🚀 Quick Start

### Prerequisites

- Go 1.24.6 or later
- Node.js 18+ and npm
- WHOP API credentials

### 1. Clone and Setup

```bash
git clone https://github.com/YOUR_USERNAME/ggpoker.git
cd ggpoker
```

### 2. Configure Environment Variables

Create a `.env` file in the root directory:

```bash
# WHOP Configuration
export WHOP_API_KEY="your_whop_api_key_here"
export WHOP_PRODUCT_ID="your_product_id_here"
export WHOP_WEBHOOK_SECRET="your_webhook_secret_here"

# Security
export JWT_SECRET="your_jwt_secret_here"
```

### 3. Install Dependencies

```bash
# Go dependencies
go mod tidy

# Frontend dependencies
cd web
npm install
cd ..
```

### 4. Build and Run

```bash
# Build the Go application
make build

# Build the frontend
cd web
npm run build
cd ..

# Run the application
make run
```

The application will start with:
- P2P server on port 3000
- API server on port 3001
- Web interface on port 5173

## 🎮 How to Play

### 1. Authentication
- Visit the web interface
- Enter your WHOP token to authenticate
- Verify your subscription status

### 2. Join a Game
- Click "Ready to Play" to join the table
- Wait for other players to join
- Game starts automatically when minimum players are ready

### 3. Gameplay
- **Pre-Flop**: Receive hole cards, post blinds
- **Flop**: 3 community cards dealt
- **Turn**: 4th community card
- **River**: 5th community card
- **Showdown**: Best hand wins the pot

### 4. Actions Available
- **Fold**: Give up your hand
- **Check**: Pass action without betting
- **Call**: Match the current bet
- **Bet/Raise**: Increase the pot

## 🔧 Configuration

Edit `config.yaml` to customize:

- **Blind amounts**: Small/big blind values
- **Starting stack**: Initial chip count per player
- **Network ports**: P2P and API server ports
- **Game rules**: Timeouts, player limits, etc.

## 🌐 WHOP Integration

### Setup WHOP Product

1. Create a product at [WHOP.com](https://whop.com)
2. Set up subscription tiers and pricing
3. Configure webhooks for real-time updates
4. Get your API credentials

### API Endpoints

- `POST /api/whop/validate` - Validate user tokens
- `GET /api/whop/user/{id}` - Get user information
- `GET /api/whop/subscriptions/{id}` - Check subscription status

### Webhook Handling

The application automatically handles WHOP webhooks for:
- Subscription activations
- Payment confirmations
- User cancellations

## 🧪 Testing

```bash
# Run Go tests
make test

# Run frontend tests
cd web
npm test
cd ..

# Run integration tests
make test-integration
```

## 📱 Deployment

### Docker Deployment

```bash
# Build Docker image
docker build -t ggpoker .

# Run container
docker run -p 3000:3000 -p 3001:3001 -p 5173:5173 \
  -e WHOP_API_KEY=your_key \
  -e WHOP_PRODUCT_ID=your_product \
  ggpoker
```

### Production Considerations

- Use HTTPS for production
- Set up proper firewall rules
- Configure logging and monitoring
- Use environment variables for secrets
- Set up database persistence
- Configure rate limiting

## 🔒 Security Features

- JWT-based authentication
- WHOP token validation
- Rate limiting per IP
- Input validation and sanitization
- Secure WebSocket connections
- CORS protection

## 📊 Monitoring and Logging

- Structured logging with logrus
- Game state monitoring
- Player action tracking
- Network performance metrics
- Error reporting and alerting

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🆘 Support

- **Issues**: [GitHub Issues](https://github.com/YOUR_USERNAME/ggpoker/issues)
- **Discussions**: [GitHub Discussions](https://github.com/YOUR_USERNAME/ggpoker/discussions)
- **Documentation**: [Wiki](https://github.com/YOUR_USERNAME/ggpoker/wiki)

## 🙏 Acknowledgments

- WHOP for subscription management platform
- Go community for excellent networking libraries
- React community for frontend framework
- Poker community for game rules and testing

---

**Ready to play?** 🃏 Get your WHOP token and start playing professional-grade poker today!
