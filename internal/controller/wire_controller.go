package controller

import (
	"Art-Design-Backend/internal/service"

	"github.com/google/wire"
)

var ControllerSet = wire.NewSet(
	AuthCtrlSet,
	UserCtrlSet,
	MenuCtrlSet,
	RoleCtrlSet,
	DigitPredictSet,
	AIModelCtrlSet,
	KnowledgeBaseCtrlSet,
	BrowserAgentCtrlSet,
	OperationLogCtrlSet,
)

var AuthCtrlSet = wire.NewSet(
	NewAuthController,
	wire.Struct(new(service.AuthService), "*"),
)

var UserCtrlSet = wire.NewSet(
	NewUserController,
	wire.Struct(new(service.UserService), "*"),
)

var MenuCtrlSet = wire.NewSet(
	NewMenuController,
	wire.Struct(new(service.MenuService), "*"),
)

var RoleCtrlSet = wire.NewSet(
	NewRoleController,
	wire.Struct(new(service.RoleService), "*"),
)

var DigitPredictSet = wire.NewSet(
	NewDigitPredictController,
	wire.Struct(new(service.DigitPredictService), "*"),
)

var AIModelCtrlSet = wire.NewSet(
	NewAIController,
	wire.Struct(new(service.AIService), "*"),
)

var KnowledgeBaseCtrlSet = wire.NewSet(
	NewKnowledgeBaseController,
	wire.Struct(new(service.KnowledgeBaseService), "*"),
)

var BrowserAgentCtrlSet = wire.NewSet(
	NewBrowserAgentController,
	service.NewBrowserAgentService,
	wire.Struct(new(service.BrowserAgentDashboardService), "*"),
)

var OperationLogCtrlSet = wire.NewSet(
	NewOperationLogController,
	wire.Struct(new(service.OperationLogService), "*"),
)
