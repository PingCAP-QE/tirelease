package fsm

// State is the template of state machine.
// Constructed by template model, see https://refactoring.guru/design-patterns/template-method/go/example
type State[T IStateContext] struct {
	IState[T]
}

func (state *State[T]) onTran(trans IStateTransition[T], context *T) (bool, error) {
	fitConstraints, err := trans.FitConstraints(*context)

	if err != nil {
		return false, nil
	}

	if fitConstraints {
		isTransOk, err := trans.Effect(context)
		if err != nil {
			return false, nil
		}

		return isTransOk, nil

	} else {
		return false, nil
	}
}

func (state *State[T]) onLeave(context *T) (bool, error) {
	return true, nil
}

func (state *State[T]) Dispatch(toState StateText, context *T) (bool, error) {
	if toState == state.IState.getStateText() {
		return false, nil
	}

	transition := state.IState.getTransitionMap()[StateTransitionMeta{state.getStateText(), toState}]

	isTransOK, err := state.IState.onTran(transition, context)
	if err != nil {
		return false, err
	}
	if !isTransOK {
		return false, nil
	} else {
		// TODO: See if the state should be changed before or after the transition is done.
		// If the answer is YES, then maybe there should be a method to rollback the transition.
		state.IState.setStateText(toState)
	}

	isLeaveOK, err := state.IState.onLeave(context)
	if err != nil {
		return false, err
	}
	if !isLeaveOK {
		return false, nil
	}

	return true, nil
}
