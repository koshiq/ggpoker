import React, { useState, useEffect } from 'react';
import './PokerTable.css';

const PokerTable = ({ gameState, onAction, playerAddress }) => {
  const [betAmount, setBetAmount] = useState(0);
  const [showBetInput, setShowBetInput] = useState(false);

  if (!gameState) {
    return (
      <div className="poker-table-container">
        <div className="loading">Loading game...</div>
      </div>
    );
  }

  const { players, communityCards, pot, currentBet, minRaise, currentRound, handNumber } = gameState;
  const currentPlayer = players[playerAddress];

  const handleAction = (action, amount = 0) => {
    onAction(action, amount);
    if (action === 'bet' || action === 'raise') {
      setShowBetInput(false);
      setBetAmount(0);
    }
  };

  const canCheck = currentPlayer && currentPlayer.totalBet >= currentBet;
  const canCall = currentPlayer && currentPlayer.totalBet < currentBet;
  const callAmount = currentBet - (currentPlayer?.totalBet || 0);

  const renderCard = (card) => {
    if (!card) return null;
    
    const suitSymbols = {
      'SPADES': '♠',
      'HARTS': '♥', 
      'DIAMONDS': '♦',
      'CLUBS': '♣'
    };
    
    const value = card.Value === 1 ? 'A' : card.Value === 11 ? 'J' : card.Value === 12 ? 'Q' : card.Value === 13 ? 'K' : card.Value;
    const suit = suitSymbols[card.Suit] || card.Suit;
    
    return (
      <div className={`card ${card.Suit === 'HARTS' || card.Suit === 'DIAMONDS' ? 'red' : 'black'}`}>
        <div className="card-value">{value}</div>
        <div className="card-suit">{suit}</div>
      </div>
    );
  };

  const renderPlayer = (addr, player) => {
    const isCurrentPlayer = addr === playerAddress;
    const isActive = !player.folded && !player.allIn;
    
    return (
      <div key={addr} className={`player-seat ${isCurrentPlayer ? 'current-player' : ''} ${!isActive ? 'inactive' : ''}`}>
        <div className="player-info">
          <div className="player-address">{addr}</div>
          <div className="player-stack">${player.stack}</div>
          {player.bet > 0 && (
            <div className="player-bet">${player.bet}</div>
          )}
          {player.folded && <div className="player-status folded">FOLDED</div>}
          {player.allIn && <div className="player-status allin">ALL IN</div>}
          {player.isDealer && <div className="player-status dealer">D</div>}
          {player.isSmallBlind && <div className="player-status sb">SB</div>}
          {player.isBigBlind && <div className="player-status bb">BB</div>}
        </div>
        
        {isCurrentPlayer && player.holeCards && (
          <div className="hole-cards">
            {player.holeCards.map((card, index) => (
              <div key={index} className="hole-card">
                {renderCard(card)}
              </div>
            ))}
          </div>
        )}
      </div>
    );
  };

  return (
    <div className="poker-table-container">
      <div className="game-info">
        <div className="hand-number">Hand #{handNumber}</div>
        <div className="current-round">{currentRound}</div>
        <div className="pot-info">
          <div className="pot-label">Pot:</div>
          <div className="pot-amount">${pot.reduce((sum, p) => sum + p.Amount, 0)}</div>
        </div>
        {currentBet > 0 && (
          <div className="current-bet">Current Bet: ${currentBet}</div>
        )}
      </div>

      <div className="poker-table">
        <div className="community-cards">
          {communityCards.map((card, index) => (
            <div key={index} className="community-card">
              {renderCard(card)}
            </div>
          ))}
        </div>

        <div className="player-positions">
          {Object.entries(players).map(([addr, player]) => renderPlayer(addr, player))}
        </div>
      </div>

      {currentPlayer && !currentPlayer.folded && !currentPlayer.allIn && (
        <div className="action-controls">
          <div className="action-buttons">
            <button 
              className="action-btn fold"
              onClick={() => handleAction('fold')}
            >
              Fold
            </button>
            
            {canCheck ? (
              <button 
                className="action-btn check"
                onClick={() => handleAction('check')}
              >
                Check
              </button>
            ) : (
              <button 
                className="action-btn call"
                onClick={() => handleAction('call')}
                disabled={!canCall}
              >
                Call ${callAmount}
              </button>
            )}
            
            <button 
              className="action-btn bet"
              onClick={() => setShowBetInput(!showBetInput)}
            >
              {currentBet > 0 ? 'Raise' : 'Bet'}
            </button>
          </div>

          {showBetInput && (
            <div className="bet-input">
              <input
                type="number"
                min={minRaise}
                max={currentPlayer.stack}
                value={betAmount}
                onChange={(e) => setBetAmount(parseInt(e.target.value) || 0)}
                placeholder={`Min: $${minRaise}`}
              />
              <button 
                className="confirm-bet"
                onClick={() => handleAction(currentBet > 0 ? 'raise' : 'bet', betAmount)}
                disabled={betAmount < minRaise || betAmount > currentPlayer.stack}
              >
                Confirm
              </button>
            </div>
          )}
        </div>
      )}

      <div className="game-log">
        <h3>Game Log</h3>
        <div className="log-entries">
          {Object.entries(players).map(([addr, player]) => {
            if (player.lastAction && player.lastAction !== 'NONE') {
              return (
                <div key={addr} className="log-entry">
                  {addr}: {player.lastAction} {player.bet > 0 ? `$${player.bet}` : ''}
                </div>
              );
            }
            return null;
          })}
        </div>
      </div>
    </div>
  );
};

export default PokerTable;
