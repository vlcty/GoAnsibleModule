package GoAnsibleModule

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	// "strings"
)

type AnsibleResult struct {
	Changed   bool                   `json:"changed"`
	Failed    bool                   `json:"failed"`
	Message   string                 `json:"msg"`
	Arguments map[string]interface{} `json:"arguments"`
}

func (result *AnsibleResult) mashalAndExit() {
	json, err := json.Marshal(result)

	if err != nil {
		panic(err)
	} else {
		fmt.Printf("%s\n", json)
	}

	os.Exit(0)
}

type AnsibleModule struct {
	Arguments         map[string]interface{}
	RequiredArguments []string
}

func NewAnsibleModule() *AnsibleModule {
	ansibleModule := &AnsibleModule{}
	ansibleModule.parseArgumentsFile()

	return ansibleModule
}

func (module *AnsibleModule) parseArgumentsFile() {
	if len(os.Args) < 2 {
		return
	}

	data, err := ioutil.ReadFile(os.Args[1])

	if err != nil {
		module.ExitFailed(err.Error())
	}

	module.Arguments = make(map[string]interface{})

	json.Unmarshal(data, &module.Arguments)

	if err != nil {
		module.ExitFailed(err.Error())
	}
}

func (module *AnsibleModule) CheckRequiredArguments() {
	for _, argumentToLookFor := range module.RequiredArguments {
		_, isThere := module.Arguments[argumentToLookFor]

		if !isThere {
			module.ExitFailed(fmt.Sprintf("Required argument %s is missing",
				argumentToLookFor))
		}
	}
}

func (module *AnsibleModule) ExitUnchanged() {
	result := AnsibleResult{}
	result.Arguments = module.Arguments

	result.mashalAndExit()
}

func (module *AnsibleModule) ExitChanged() {
	result := AnsibleResult{}
	result.Changed = true
	result.Arguments = module.Arguments

	result.mashalAndExit()
}

func (module *AnsibleModule) ExitFailed(message string) {
	result := AnsibleResult{}
	result.Failed = true
	result.Message = message
	result.Arguments = module.Arguments

	result.mashalAndExit()
}
