package behavior3go

import (
	_ "github.com/behavior3/behavior3go/actions"
	actionspkg "github.com/behavior3/behavior3go/actions"
	_ "github.com/behavior3/behavior3go/composites"
	compositespkg "github.com/behavior3/behavior3go/composites"
	"github.com/behavior3/behavior3go/core"
	_ "github.com/behavior3/behavior3go/decorators"
	decoratorspkg "github.com/behavior3/behavior3go/decorators"
)

const (
	VERSION   = core.VERSION
	SUCCESS   = core.SUCCESS
	FAILURE   = core.FAILURE
	RUNNING   = core.RUNNING
	ERROR     = core.ERROR
	COMPOSITE = core.COMPOSITE
	DECORATOR = core.DECORATOR
	ACTION    = core.ACTION
	CONDITION = core.CONDITION
)

type Status = core.Status
type Node = core.Node
type NodeConstructor = core.NodeConstructor
type TreeData = core.TreeData
type NodeData = core.NodeData
type CustomNodeData = core.CustomNodeData

type BaseNode = core.BaseNode
type Action = core.Action
type Composite = core.Composite
type Decorator = core.Decorator
type Condition = core.Condition
type BehaviorTree = core.BehaviorTree
type Blackboard = core.Blackboard
type Tick = core.Tick

type Error = actionspkg.Error
type Failer = actionspkg.Failer
type Runner = actionspkg.Runner
type Succeeder = actionspkg.Succeeder
type Wait = actionspkg.Wait

type Priority = compositespkg.Priority
type MemPriority = compositespkg.MemPriority
type MemSequence = compositespkg.MemSequence
type Sequence = compositespkg.Sequence

type Inverter = decoratorspkg.Inverter
type Limiter = decoratorspkg.Limiter
type MaxTime = decoratorspkg.MaxTime
type Repeater = decoratorspkg.Repeater
type RepeatUntilFailure = decoratorspkg.RepeatUntilFailure
type RepeatUntilSuccess = decoratorspkg.RepeatUntilSuccess

var NewBaseNode = core.NewBaseNode
var NewAction = core.NewAction
var NewComposite = core.NewComposite
var NewDecorator = core.NewDecorator
var NewCondition = core.NewCondition
var NewBehaviorTree = core.NewBehaviorTree
var NewBlackboard = core.NewBlackboard
var NewTick = core.NewTick

var NewError = actionspkg.NewError
var NewFailer = actionspkg.NewFailer
var NewRunner = actionspkg.NewRunner
var NewSucceeder = actionspkg.NewSucceeder
var NewWait = actionspkg.NewWait

var NewPriority = compositespkg.NewPriority
var NewMemPriority = compositespkg.NewMemPriority
var NewMemSequence = compositespkg.NewMemSequence
var NewSequence = compositespkg.NewSequence

var NewInverter = decoratorspkg.NewInverter
var NewLimiter = decoratorspkg.NewLimiter
var NewMaxTime = decoratorspkg.NewMaxTime
var NewRepeater = decoratorspkg.NewRepeater
var NewRepeatUntilFailure = decoratorspkg.NewRepeatUntilFailure
var NewRepeatUntilSuccess = decoratorspkg.NewRepeatUntilSuccess
