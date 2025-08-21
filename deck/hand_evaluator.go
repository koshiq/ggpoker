package deck

import (
	"sort"
)

type HandRank int

const (
	HighCard HandRank = iota
	OnePair
	TwoPair
	ThreeOfAKind
	Straight
	Flush
	FullHouse
	FourOfAKind
	StraightFlush
	RoyalFlush
)

func (hr HandRank) String() string {
	switch hr {
	case HighCard:
		return "High Card"
	case OnePair:
		return "One Pair"
	case TwoPair:
		return "Two Pair"
	case ThreeOfAKind:
		return "Three of a Kind"
	case Straight:
		return "Straight"
	case Flush:
		return "Flush"
	case FullHouse:
		return "Full House"
	case FourOfAKind:
		return "Four of a Kind"
	case StraightFlush:
		return "Straight Flush"
	case RoyalFlush:
		return "Royal Flush"
	default:
		return "Unknown"
	}
}

type Hand struct {
	Cards []Card
	Rank  HandRank
	Value int // For comparing hands of the same rank
}

func EvaluateHand(cards []Card) Hand {
	if len(cards) < 5 {
		return Hand{Cards: cards, Rank: HighCard, Value: 0}
	}

	// Get best 5-card combination
	bestHand := getBestFiveCardHand(cards)

	// Evaluate the hand
	rank, value := evaluateFiveCardHand(bestHand)

	return Hand{
		Cards: bestHand,
		Rank:  rank,
		Value: value,
	}
}

func getBestFiveCardHand(cards []Card) []Card {
	if len(cards) == 5 {
		return cards
	}

	// For now, return first 5 cards
	// TODO: Implement proper 5-card selection for 7-card hands
	return cards[:5]
}

func evaluateFiveCardHand(cards []Card) (HandRank, int) {
	if len(cards) != 5 {
		return HighCard, 0
	}

	// Check for flush
	isFlush := isFlush(cards)

	// Check for straight
	isStraight, straightValue := isStraight(cards)

	// Check for straight flush
	if isFlush && isStraight {
		if straightValue == 14 { // Ace high straight
			return RoyalFlush, 14
		}
		return StraightFlush, straightValue
	}

	// Check for four of a kind
	if fourValue := hasFourOfAKind(cards); fourValue > 0 {
		return FourOfAKind, fourValue
	}

	// Check for full house
	if fullHouseValue := hasFullHouse(cards); fullHouseValue > 0 {
		return FullHouse, fullHouseValue
	}

	// Check for flush
	if isFlush {
		return Flush, getHighCardValue(cards)
	}

	// Check for straight
	if isStraight {
		return Straight, straightValue
	}

	// Check for three of a kind
	if threeValue := hasThreeOfAKind(cards); threeValue > 0 {
		return ThreeOfAKind, threeValue
	}

	// Check for two pair
	if twoPairValue := hasTwoPair(cards); twoPairValue > 0 {
		return TwoPair, twoPairValue
	}

	// Check for one pair
	if pairValue := hasOnePair(cards); pairValue > 0 {
		return OnePair, pairValue
	}

	// High card
	return HighCard, getHighCardValue(cards)
}

func isFlush(cards []Card) bool {
	if len(cards) != 5 {
		return false
	}
	suit := cards[0].Suit
	for _, card := range cards {
		if card.Suit != suit {
			return false
		}
	}
	return true
}

func isStraight(cards []Card) (bool, int) {
	if len(cards) != 5 {
		return false, 0
	}

	// Sort cards by value
	sortedCards := make([]Card, len(cards))
	copy(sortedCards, cards)
	sort.Slice(sortedCards, func(i, j int) bool {
		return sortedCards[i].Value < sortedCards[j].Value
	})

	// Check for regular straight
	for i := 1; i < 5; i++ {
		if sortedCards[i].Value != sortedCards[i-1].Value+1 {
			// Check for Ace-low straight (A,2,3,4,5)
			if i == 4 && sortedCards[0].Value == 1 && sortedCards[1].Value == 2 &&
				sortedCards[2].Value == 3 && sortedCards[3].Value == 4 && sortedCards[4].Value == 5 {
				return true, 5
			}
			return false, 0
		}
	}

	return true, sortedCards[4].Value
}

func hasFourOfAKind(cards []Card) int {
	valueCount := make(map[int]int)
	for _, card := range cards {
		valueCount[card.Value]++
	}

	for value, count := range valueCount {
		if count == 4 {
			return value
		}
	}
	return 0
}

func hasFullHouse(cards []Card) int {
	valueCount := make(map[int]int)
	for _, card := range cards {
		valueCount[card.Value]++
	}

	var threeValue, twoValue int
	for value, count := range valueCount {
		if count == 3 {
			threeValue = value
		} else if count == 2 {
			twoValue = value
		}
	}

	if threeValue > 0 && twoValue > 0 {
		return threeValue*100 + twoValue
	}
	return 0
}

func hasThreeOfAKind(cards []Card) int {
	valueCount := make(map[int]int)
	for _, card := range cards {
		valueCount[card.Value]++
	}

	for value, count := range valueCount {
		if count == 3 {
			return value
		}
	}
	return 0
}

func hasTwoPair(cards []Card) int {
	valueCount := make(map[int]int)
	for _, card := range cards {
		valueCount[card.Value]++
	}

	var pairs []int
	for value, count := range valueCount {
		if count == 2 {
			pairs = append(pairs, value)
		}
	}

	if len(pairs) == 2 {
		// Sort pairs in descending order
		if pairs[0] < pairs[1] {
			pairs[0], pairs[1] = pairs[1], pairs[0]
		}
		return pairs[0]*100 + pairs[1]
	}
	return 0
}

func hasOnePair(cards []Card) int {
	valueCount := make(map[int]int)
	for _, card := range cards {
		valueCount[card.Value]++
	}

	for value, count := range valueCount {
		if count == 2 {
			return value
		}
	}
	return 0
}

func getHighCardValue(cards []Card) int {
	maxValue := 0
	for _, card := range cards {
		if card.Value > maxValue {
			maxValue = card.Value
		}
	}
	return maxValue
}

// CompareHands returns 1 if hand1 wins, -1 if hand2 wins, 0 if tie
func CompareHands(hand1, hand2 Hand) int {
	if hand1.Rank > hand2.Rank {
		return 1
	}
	if hand1.Rank < hand2.Rank {
		return -1
	}

	// Same rank, compare values
	if hand1.Value > hand2.Value {
		return 1
	}
	if hand1.Value < hand2.Value {
		return -1
	}

	// Same rank and value, compare high cards
	return compareHighCards(hand1.Cards, hand2.Cards)
}

func compareHighCards(cards1, cards2 []Card) int {
	// Sort both hands by value in descending order
	sort.Slice(cards1, func(i, j int) bool {
		return cards1[i].Value > cards1[j].Value
	})
	sort.Slice(cards2, func(i, j int) bool {
		return cards2[i].Value > cards2[j].Value
	})

	// Compare high cards
	for i := 0; i < len(cards1) && i < len(cards2); i++ {
		if cards1[i].Value > cards2[i].Value {
			return 1
		}
		if cards1[i].Value < cards2[i].Value {
			return -1
		}
	}

	return 0 // Tie
}
