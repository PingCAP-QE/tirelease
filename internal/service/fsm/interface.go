package fsm

type StateText string

type StateTransitionMeta struct {
	FromState StateText
	ToState   StateText
}

// A transition is used to move from one state to another.
// FitConstraints is used to check if the transition is valid.
// Effect is used to take the effect of the transition.
type IStateTransition[T IStateContext] interface {
	FitConstraints(context T) (bool, error)
	Effect(context *T) (bool, error)
}

// A StateContext is the context object of the state machine. Such as the Version、the Triage Result...
// It contains all the information about the state context.
// It is corresponding to an entity of DDD in TiRelease.
// Trans is used to make the state transition.
type IStateContext interface {
	Trans(toState StateText) (bool, error)
}

// State is the key object of the state machine.
// init() is used to initialize the State, like defining the TransitionMap
// Dispatch() is the key function to deliver the transition signal
//     and make the state transition by using onTran and on Leave.
// onTran() is used to deal the 1 to 1 dispatch by certain transition.
// onLeave() is used when a common logic should be applied when the state is leaving.
type IState[T IStateContext] interface {
	init() error
	onTran(tran IStateTransition[T], context *T) (bool, error)
	onLeave(context *T) (bool, error)
	getStateText() StateText
	setStateText(stateText StateText)
	getTransitionMap() TransitionMap[T]
	Dispatch(toState StateText, context *T) (bool, error)
}

type TransitionMap[T IStateContext] map[StateTransitionMeta]IStateTransition[T]
