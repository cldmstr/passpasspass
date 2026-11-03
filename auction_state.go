package passpasspass

// AuctionState is what the auction has made public: the calls so far, the current contract, and
// what each player has shown by their calls alone. Every seat's call is read with the
// interpreter and that seat's convention sheet.
type AuctionState struct {
	interp   Interpreter
	sheets   [4]*ConventionSheet // indexed by Direction
	history  History
	shown    [4]HandModel // indexed by Direction
	contract *Bid
}

// NewAuctionState returns the state before the first call: interp reads the calls, sheets holds
// each seat's convention sheet, and shown what each player is known to hold at the start.
func NewAuctionState(interp Interpreter, sheets [4]*ConventionSheet, shown [4]HandModel) *AuctionState {
	return &AuctionState{interp: interp, sheets: sheets, shown: shown}
}

// Record adds the call bid by d. The caller's model is narrowed by what the call shows; when the
// call contradicts everything the caller had shown, the latest call wins and its meaning
// replaces the model.
func (s *AuctionState) Record(d Direction, bid *Bid) {
	read := s.interp.Read(s.sheets[d], &CallContext{Caller: d, History: s.history, Shown: s.shown}, bid)
	shown := s.shown[d].Intersect(read)
	if shown.IsEmpty() {
		shown = read
	}
	s.shown[d] = shown
	s.history = append(s.history, BidAction{Direction: d, Bid: bid})
	if bid != Pass {
		s.contract = bid
	}
}

// History returns the calls so far.
func (s *AuctionState) History() History {
	return s.history
}

// Contract returns the highest bid so far, or nil before the opening.
func (s *AuctionState) Contract() *Bid {
	return s.contract
}

// Shown returns what d has shown by their calls.
func (s *AuctionState) Shown(d Direction) HandModel {
	return s.shown[d]
}

// Sheet returns d's convention sheet.
func (s *AuctionState) Sheet(d Direction) *ConventionSheet {
	return s.sheets[d]
}

// Perspective returns the view of the player at d, who holds hand.
func (s *AuctionState) Perspective(d Direction, hand *Hand) Perspective {
	return Perspective{
		Hand:     hand,
		Seat:     d,
		History:  s.history,
		Contract: s.contract,
		Partner:  seenWith(s.shown[partnerOf(d)], hand),
		LHO:      s.shown[lhoOf(d)],
		RHO:      s.shown[rhoOf(d)],
	}
}

// Perspective is one player's view of the auction: their own hand, partner's model as it looks
// from that hand, the opponents' models as shown, and the calls and contract so far.
type Perspective struct {
	Hand     *Hand
	Seat     Direction
	History  History
	Contract *Bid
	Partner  HandModel
	LHO      HandModel
	RHO      HandModel
}

// Context returns the AuctionContext the player's convention sheet bids from.
func (p Perspective) Context() *AuctionContext {
	return &AuctionContext{
		SelfHand:    p.Hand,
		Position:    p.Seat,
		History:     p.History,
		Contract:    p.Contract,
		Partner:     p.Partner,
		LHO:         p.LHO,
		RHO:         p.RHO,
		Partnership: NewPartnershipEval(p.Hand, boundsOf(p.Partner)),
	}
}

// seenWith narrows partner's model by the cards in hand: partner holds at most 13 minus own
// length in each suit and at most 40 minus own HCP. When no hand partner could hold is left,
// partner has misbid or is lying, so the model is kept as shown and marked as a conflict.
func seenWith(partner HandModel, hand *Hand) HandModel {
	room := DefaultHandRange()
	room.MaxHCP = MaxDeckHCP - hand.HCP()
	for s, n := range hand.Distribution() {
		room.SuitLength[s].Max = HandSize - n
	}
	seen := partner.Intersect(HandModel{Cases: []HandRange{room}})
	if seen.IsEmpty() {
		partner.Conflict = true
		return partner
	}
	return seen
}
