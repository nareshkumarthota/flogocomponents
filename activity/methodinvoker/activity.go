package methodinvoker

import (
	"github.com/project-flogo/core/activity"
)

func init() {
	_ = activity.Register(&Activity{})
	methodData = make(map[string]UserDefinedMethod)
}

// UserDefinedMethod signature for userdefined methods
type UserDefinedMethod func(inputs interface{}) (map[string]interface{}, error)

var activityMd = activity.ToMetadata(&Input{}, &Output{})
var methodData map[string]UserDefinedMethod

// RegisterMethods registers userdefined functions to methodData map
func RegisterMethods(methodName string, mthd UserDefinedMethod) {
	methodData[methodName] = mthd
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
			outPutFromFunc, err = val(input.ToMap()["inputData"])
			if err != nil {
				ctx.Logger().Error("error in executing method: ", input.MethodName)
			} else {
				funcExeFlag = true
			}
		} else {
			ctx.Logger().Errorf("method[%s] not registerd to activity", input.MethodName)
		}
	} else {
		ctx.Logger().Error("methods not registerd to activity sending default response")
	}

	if !funcExeFlag {
		resp := make(map[string]interface{})
		resp["response"] = "success message from method invoke activity"
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
