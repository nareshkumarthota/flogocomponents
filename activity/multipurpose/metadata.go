package multipurpose

import "github.com/project-flogo/core/data/coerce"

type Settings struct {
	ASetting string `md:"aSetting,required"`
}

type Input struct {
	MethodName string      `md:"methodName,required"`
	InputData  interface{} `md:"inputData"`
}

func (r *Input) FromMap(values map[string]interface{}) error {
	strVal, _ := coerce.ToString(values["methodName"])
	r.MethodName = strVal

	intrfcVal, _ := coerce.ToAny(values["inputData"])
	r.InputData = intrfcVal

	return nil
}

func (r *Input) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"methodName": r.MethodName,
		"inputData":  r.InputData,
	}
}

type Output struct {
	MethodName string      `md:"methodName"`
	OutputData interface{} `md:"outputData"`
}

func (o *Output) FromMap(values map[string]interface{}) error {
	strVal, _ := coerce.ToString(values["methodName"])
	o.MethodName = strVal

	intrfcVal, _ := coerce.ToAny(values["outputData"])
	o.OutputData = intrfcVal
	return nil
}

func (o *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"methodName": o.MethodName,
		"outputData": o.OutputData,
	}
}
