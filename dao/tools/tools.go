package tools

import (
	"context"
	"errors"
	"fmt"
	"github.com/Jackzode/painting/commons/handler"
	"reflect"
	"xorm.io/xorm"
)

func Help(page, pageSize int, rowsSlicePtr interface{}, rowElement interface{}, session *xorm.Session) (total int64, err error) {
	page, pageSize = ValPageAndPageSize(page, pageSize)

	sliceValue := reflect.Indirect(reflect.ValueOf(rowsSlicePtr))
	if sliceValue.Kind() != reflect.Slice {
		return 0, errors.New("not a slice")
	}

	startNum := (page - 1) * pageSize
	return session.Limit(pageSize, startNum).FindAndCount(rowsSlicePtr, rowElement)
}

func ValPageAndPageSize(page, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	return page, pageSize
}

const (
	QuestionObjectType   = "question"
	AnswerObjectType     = "answer"
	TagObjectType        = "tag"
	UserObjectType       = "user"
	CollectionObjectType = "collection"
	CommentObjectType    = "comment"
	ReportObjectType     = "report"
)

var (
	ObjectTypeStrMapping = map[string]int{
		QuestionObjectType:   1,
		AnswerObjectType:     2,
		TagObjectType:        3,
		UserObjectType:       4,
		CollectionObjectType: 6,
		CommentObjectType:    7,
		ReportObjectType:     8,
	}
)

func GenUniqueIDStr(ctx context.Context, key string) (uniqueID string, err error) {
	objectType := ObjectTypeStrMapping[key]
	luaScript := `
    local current_value = redis.call("GET", KEYS[1])
    if not current_value then
        current_value = 1
    else
        current_value = tonumber(current_value)
    end
    local new_value = current_value + 1
    redis.call("SET", KEYS[1], new_value)
    return new_value
    `
	uniqKey := fmt.Sprintf("_lawyer_uniq_%v", objectType)
	eval := handler.RedisClient.Eval(ctx, luaScript, []string{uniqKey})
	if eval.Err() != nil {
		return "", eval.Err()
	}
	return fmt.Sprintf("1%03d%013d", objectType, eval.Val()), nil
}
