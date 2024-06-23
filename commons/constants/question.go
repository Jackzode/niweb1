package constants

const (
	QuestionStatusAvailable = 1
	QuestionStatusClosed    = 2
	QuestionStatusDeleted   = 10
	QuestionUnPin           = 1
	QuestionPin             = 2
	QuestionShow            = 1
	QuestionHide            = 2
)

const (
	QuestionOperationPin   = "pin"
	QuestionOperationUnPin = "unpin"
	QuestionOperationHide  = "hide"
	QuestionOperationShow  = "show"
)

const (
	QuestionOrderCondNewest     = "newest"
	QuestionOrderCondActive     = "active"
	QuestionOrderCondFrequent   = "frequent"
	QuestionOrderCondScore      = "score"
	QuestionOrderCondUnanswered = "unanswered"
)

const (
	OperationLevelInfo    = "info"
	OperationLevelDanger  = "danger"
	OperationLevelWarning = "warning"
)

const (
	QuestionPageRespOperationTypeAsked    = "asked"
	QuestionPageRespOperationTypeAnswered = "answered"
	QuestionPageRespOperationTypeModified = "modified"
)

var AdminQuestionSearchStatus = map[string]int{
	"available": QuestionStatusAvailable,
	"closed":    QuestionStatusClosed,
	"deleted":   QuestionStatusDeleted,
}

var AdminQuestionSearchStatusIntToString = map[int]string{
	QuestionStatusAvailable: "available",
	QuestionStatusClosed:    "closed",
	QuestionStatusDeleted:   "deleted",
}
