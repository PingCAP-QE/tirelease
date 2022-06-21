package fsm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	TestState1 = StateText("test_state_1")
	TestState2 = StateText("test_state_2")
	TestState3 = StateText("test_state_3")
)

type TestStateTransition_1_2 struct {
}

func (trans TestStateTransition_1_2) FitConstraints(context *testStateContext) (bool, error) {
	if context.Var1 == "test1" {
		return true, nil
	} else {
		return false, nil
	}
}

func (trans TestStateTransition_1_2) Effect(context **testStateContext) (bool, error) {
	(*context).Var2 = "testEffect"
	return true, nil
}

type testStateContext struct {
	State *testState
	ID    string
	Var1  string
	Var2  string
}

func NewStateContext(stateText StateText, ID string) (*testStateContext, error) {
	context := &testStateContext{}
	context.ID = "test_id"

	State, err := NewState(stateText)
	if err != nil {
		return nil, err
	}
	context.State = State

	context.Var1 = "test1"
	context.Var2 = "test2"

	return context, nil
}

func (context *testStateContext) Trans(toState StateText) (bool, error) {
	isSuccess, err := context.State.Dispatch(toState, &context)
	if err != nil {
		return false, err
	}
	return isSuccess, nil
}

// Make the State struct private to force the only entrance be NewState func.
type testState struct {
	State[*testStateContext]
	StateText StateText
	transMap  TransitionMap[*testStateContext]
}

func NewState(stateText StateText) (*testState, error) {
	if stateText == "" {
		stateText = TestState3
	}
	testState := &testState{
		StateText: stateText,
	}
	testState.IState = interface{}(testState).(IState[*testStateContext])
	testState.init()

	return testState, nil
}

func (state *testState) getStateText() StateText {
	return state.StateText
}

func (state *testState) setStateText(stateText StateText) {
	state.StateText = stateText
}

func (state *testState) getTransitionMap() TransitionMap[*testStateContext] {
	if state.transMap == nil {
		state.transMap = make(TransitionMap[*testStateContext])
	}
	return state.transMap
}

func (state *testState) init() error {
	if len(state.getTransitionMap()) > 0 {
		return nil
	}
	state.getTransitionMap()[StateTransitionMeta{FromState: StateText(TestState1), ToState: StateText(TestState2)}] = TestStateTransition_1_2{}
	return nil
}

// The codes above are used to construct the following test cases:
// ---------------------------------------------
// The functions below are used to test the FSM.
//     Test Transition
func TestTransition(t *testing.T) {
	context, err := NewStateContext(TestState1, "test_id")
	assert.Equal(t, nil, err)
	assert.Equal(t, TestState1, context.State.getStateText())
	assert.Equal(t, "test_id", context.ID)

	transition := TestStateTransition_1_2{}
	isFitConstrains, err := transition.FitConstraints(context)
	assert.Equal(t, true, isFitConstrains)
	assert.Equal(t, nil, err)

	isTransOK, err := transition.Effect(&context)
	assert.Equal(t, true, isTransOK)
	assert.Equal(t, nil, err)
	assert.Equal(t, "testEffect", context.Var2)
}

//     Test State
func TestStateInit(t *testing.T) {
	testState, err := NewState(TestState1)
	assert.Equal(t, nil, err)
	assert.Equal(t, TestState1, testState.getStateText())
	assert.Equal(t, 1, len(testState.getTransitionMap()))
}

func TestStateDispatch(t *testing.T) {
	context, err := NewStateContext(TestState1, "test_id")
	assert.Equal(t, nil, err)

	isTransOK, err := context.State.Dispatch(TestState2, &context)

	assert.Equal(t, true, isTransOK)
	assert.Equal(t, nil, err)
	assert.Equal(t, TestState2, context.State.getStateText())
	assert.Equal(t, "testEffect", context.Var2)
}

//     Test State Context
func TestContextTrans(t *testing.T) {
	context, err := NewStateContext(TestState1, "test_id")
	assert.Equal(t, nil, err)

	isTransOK, err := context.Trans(TestState2)
	assert.Equal(t, nil, err)
	assert.Equal(t, true, isTransOK)
	assert.Equal(t, TestState2, context.State.getStateText())
}
