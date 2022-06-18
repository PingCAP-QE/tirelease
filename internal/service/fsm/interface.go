package fsm

type StateText string

type StateTransitionMeta struct {
	FromState StateText
	ToState   StateText
}

type IStateTransition[T IStateContext] interface {
	FitConstraints(context T) (bool, error)
	Effect(context *T) (bool, error)
}

type IStateContext interface {
	Trans(toState StateText) (bool, error)
}

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
