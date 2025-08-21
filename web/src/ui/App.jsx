import React, { useEffect, useState } from 'react'
import PokerTable from './PokerTable'
import './App.css'

export const App = () => {
  const [messages, setMessages] = useState([])
  const [connected, setConnected] = useState(false)
  const [gameState, setGameState] = useState(null)
  const [playerAddress, setPlayerAddress] = useState('')
  const [whopToken, setWhopToken] = useState('')
  const [isAuthenticated, setIsAuthenticated] = useState(false)
  const [showAuth, setShowAuth] = useState(true)

  useEffect(() => {
    // Check for existing WHOP token
    const token = localStorage.getItem('whop_token')
    if (token) {
      setWhopToken(token)
      setIsAuthenticated(true)
      setShowAuth(false)
    }

    const proto = location.protocol === 'https:' ? 'wss' : 'ws'
    const ws = new WebSocket(`${proto}://${location.host.replace(/:\d+$/, '')}:3001/ws`)
    
    ws.onopen = () => setConnected(true)
    ws.onclose = () => setConnected(false)
    ws.onmessage = (evt) => {
      try {
        const data = JSON.parse(evt.data)
        if (data.type === 'game_state') {
          setGameState(data.data)
        } else {
          setMessages((m) => [...m, evt.data])
        }
      } catch (e) {
        setMessages((m) => [...m, evt.data])
      }
    }
    
    return () => ws.close()
  }, [])

  const handleWhopAuth = async (token) => {
    try {
      // Validate token with WHOP API
      const response = await fetch('/api/whop/validate', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ token }),
      })

      if (response.ok) {
        const userData = await response.json()
        setWhopToken(token)
        setIsAuthenticated(true)
        setShowAuth(false)
        setPlayerAddress(userData.id)
        localStorage.setItem('whop_token', token)
      } else {
        alert('Invalid WHOP token. Please check your token and try again.')
      }
    } catch (error) {
      console.error('Authentication error:', error)
      alert('Authentication failed. Please try again.')
    }
  }

  const handleGameAction = async (action, amount = 0) => {
    if (!isAuthenticated) {
      alert('Please authenticate with WHOP first')
      return
    }

    try {
      let endpoint = ''
      let body = {}

      switch (action) {
        case 'fold':
          endpoint = '/fold'
          break
        case 'check':
          endpoint = '/check'
          break
        case 'call':
          endpoint = '/call'
          break
        case 'bet':
          endpoint = `/bet/${amount}`
          break
        case 'raise':
          endpoint = `/bet/${amount}`
          break
        default:
          console.error('Unknown action:', action)
          return
      }

      const response = await fetch(endpoint, {
        method: 'GET',
        headers: {
          'Authorization': `Bearer ${whopToken}`,
        },
      })

      if (!response.ok) {
        throw new Error(`Action failed: ${response.statusText}`)
      }

      console.log(`${action} action successful`)
    } catch (error) {
      console.error('Game action error:', error)
      alert(`Failed to ${action}: ${error.message}`)
    }
  }

  const handleReady = async () => {
    if (!isAuthenticated) {
      alert('Please authenticate with WHOP first')
      return
    }

    try {
      const response = await fetch('/ready', {
        method: 'GET',
        headers: {
          'Authorization': `Bearer ${whopToken}`,
        },
      })

      if (!response.ok) {
        throw new Error(`Ready failed: ${response.statusText}`)
      }

      console.log('Ready action successful')
    } catch (error) {
      console.error('Ready action error:', error)
      alert(`Failed to ready: ${error.message}`)
    }
  }

  const logout = () => {
    setWhopToken('')
    setIsAuthenticated(false)
    setShowAuth(true)
    setPlayerAddress('')
    localStorage.removeItem('whop_token')
  }

  if (showAuth) {
    return (
      <div className="auth-container">
        <div className="auth-card">
          <h1>🎰 GG Poker</h1>
          <p>Connect your WHOP account to start playing</p>
          
          <div className="auth-form">
            <input
              type="text"
              placeholder="Enter your WHOP token"
              value={whopToken}
              onChange={(e) => setWhopToken(e.target.value)}
              className="token-input"
            />
            <button 
              onClick={() => handleWhopAuth(whopToken)}
              className="auth-button"
              disabled={!whopToken.trim()}
            >
              Connect Account
            </button>
          </div>
          
          <div className="auth-help">
            <p>Don't have a WHOP token?</p>
            <a 
              href="https://whop.com" 
              target="_blank" 
              rel="noopener noreferrer"
              className="whop-link"
            >
              Get one at WHOP.com
            </a>
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className="app-container">
      <header className="app-header">
        <div className="header-content">
          <h1>🎰 GG Poker</h1>
          <div className="header-info">
            <div className="connection-status">
              WS: {connected ? '🟢 Connected' : '🔴 Disconnected'}
            </div>
            <div className="player-info">
              Player: {playerAddress.slice(0, 8)}...
            </div>
            <button onClick={logout} className="logout-btn">
              Logout
            </button>
          </div>
        </div>
      </header>

      <div className="game-controls">
        <button 
          onClick={handleReady}
          className="ready-btn"
          disabled={!connected || !isAuthenticated}
        >
          Ready to Play
        </button>
      </div>

      {gameState ? (
        <PokerTable 
          gameState={gameState}
          onAction={handleGameAction}
          playerAddress={playerAddress}
        />
      ) : (
        <div className="waiting-message">
          <h2>Waiting for game to start...</h2>
          <p>Click "Ready to Play" when you're ready to join the table</p>
        </div>
      )}

      <div className="debug-section">
        <h3>Debug Info</h3>
        <div className="debug-content">
          <div>Connected: {connected ? 'Yes' : 'No'}</div>
          <div>Authenticated: {isAuthenticated ? 'Yes' : 'No'}</div>
          <div>Player Address: {playerAddress || 'None'}</div>
        </div>
        
        <details className="message-log">
          <summary>Message Log</summary>
          <pre>{messages.join('\n')}</pre>
        </details>
      </div>
    </div>
  )
}


