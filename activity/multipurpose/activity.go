package multipurpose

import (
	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/data/metadata"
)

func init() {
	_ = activity.Register(&Activity{}) //activity.Register(&Activity{}, New) to create instances using factory method 'New'
	methodData = make(map[string]UserDefFunc)
}

// UserDefFunc signature for userdefined functions
type UserDefFunc func(inputs map[string]interface{}) (map[string]interface{}, error)

var activityMd = activity.ToMetadata(&Settings{}, &Input{}, &Output{})
var methodData map[string]UserDefFunc

// RegisterFuncs registers userdefined functions to methodData map
func RegisterFuncs(methodName string, mthd UserDefFunc) {
	methodData[methodName] = mthd
}

//New optional factory method, should be used if one activity instance per configuration is desired
func New(ctx activity.InitContext) (activity.Activity, error) {

	s := &Settings{}
	err := metadata.MapToStruct(ctx.Settings(), s, true)
	if err != nil {
		return nil, err
	}

	ctx.Logger().Debugf("Setting: %s", s.ASetting)

	act := &Activity{} //add aSetting to instance

	return act, nil
}

// Activity is an sample Activity that can be used as a base to create a custom activity
type Activity struct {
}

// Metadata returns the activity's metadata
func (a *Activity) Metadata() *activity.Metadata {
	return activityMd
}

// Eval implements api.Activity.Eval - Logs the Message
func (a *Activity) Eval(ctx activity.Context) (done bool, err error) {

	input := &Input{}
	err = ctx.GetInputObject(input)
	if err != nil {
		return true, err
	}

	var outPutFromFunc interface{}
	funcExeFlag := false
	if len(methodData) != 0 {
		if val, ok := methodData[input.MethodName]; ok {
			outPutFromFunc, err = val(input.ToMap()["inputData"].(map[string]interface{}))
			if err != nil {
				ctx.Logger().Error("error in executing method: ", input.MethodName)
			} else {
				funcExeFlag = true
			}
		} else {
			ctx.Logger().Errorf("method(%s) not registerd to activity", input.MethodName)
		}
	} else {
		ctx.Logger().Error("methods not registerd to activity sending default response")
	}

	if !funcExeFlag {
		resp := make(map[string]interface{})
		resp["response"] = "success message from activity"
		outPutFromFunc = resp
	}

	ctx.Logger().Debugf("Inputmethod name: %s", input.MethodName)

	output := &Output{MethodName: input.MethodName, OutputData: outPutFromFunc}
	err = ctx.SetOutputObject(output)
	if err != nil {
		return true, err
	}

	return true, nil
}
