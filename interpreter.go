package passpasspass

// CallContext is what every player at the table knows when a call is made: who makes it, the
// calls before it, and what each player had shown by then. It never holds the caller's hand.
type CallContext struct {
	Caller  Direction
	History History
	Shown   [4]HandModel // indexed by Direction
}

// auctionContext returns the part of an AuctionContext that conditions on the auction read: the
// caller's seat and the history. Conditions on the hand must not be evaluated with it.
func (c *CallContext) auctionContext() *AuctionContext {
	return &AuctionContext{Position: c.Caller, History: c.History}
}

// Interpreter reads what a call shows about the caller's hand.
type Interpreter interface {
	// Read returns the model of the hands with which the caller, playing sheet, could have made
	// call in ctx.
	Read(sheet *ConventionSheet, ctx *CallContext, call *Bid) HandModel
}

// sheetInterpreter reads a call from the caller's convention sheet.
type sheetInterpreter struct{}

// NewSheetInterpreter returns an Interpreter that reads a call as the union of the models of
// every convention in the sheet whose bid equals the call. A convention whose conditions cannot
// hold in the auction contributes nothing. A Pass, or a call no convention explains, shows
// nothing about the hand.
func NewSheetInterpreter() Interpreter {
	return sheetInterpreter{}
}

func (sheetInterpreter) Read(sheet *ConventionSheet, ctx *CallContext, call *Bid) HandModel {
	if call == Pass {
		return DefaultHandModel()
	}
	var shown HandModel
	for _, c := range sheet.sheet {
		if c.Bid.Value() == call.Value() {
			shown = unionModels(shown, c.Shown(ctx))
		}
	}
	if shown.IsEmpty() {
		return DefaultHandModel()
	}
	return shown
}
