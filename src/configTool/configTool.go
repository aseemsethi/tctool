package configTool

import (
	"fmt"
	"github.com/aseemsethi/tctool/src/tcGlobals"
	"github.com/aws/aws-sdk-go/service/configservice"
	"github.com/sirupsen/logrus"
)

type ConfigTool struct {
	Name                    string
	svc                     *configservice.ConfigService
	configRules             []*configservice.ConfigRule
	complianceDetailsResult []*configservice.EvaluationResult
}

var sLog *logrus.Logger

func (i *ConfigTool) Initialize() bool {
	sLog = tcGlobals.Tcg.Log
	i.svc = configservice.New(tcGlobals.Tcg.Sess, &tcGlobals.Tcg.GConf)
	i.configRules = make([]*configservice.ConfigRule, 0)
	i.complianceDetailsResult = make([]*configservice.EvaluationResult, 0)

	sLog.WithFields(logrus.Fields{"Test": "Config"}).Info("Enabled")
	return true
}

func getComplianceDetails(i *ConfigTool) {
	for _, configRule := range i.configRules {
		nextToken := ""
		for {
			output, err := i.svc.GetComplianceDetailsByConfigRule(
				&configservice.GetComplianceDetailsByConfigRuleInput{
					ConfigRuleName: configRule.ConfigRuleName,
					NextToken:      &nextToken,
				})
			if err != nil {
				fmt.Println(err)
				continue
			}
			i.complianceDetailsResult = append(i.complianceDetailsResult, output.EvaluationResults...)
			fmt.Println("Results: ", output.EvaluationResults)
			sLog.WithFields(logrus.Fields{"Test": "Config"}).Info("getComplianceRules: ", output.EvaluationResults)
			if output.NextToken == nil {
				break
			}
			nextToken = *output.NextToken
		}
	}
}

func getConfigRules(i *ConfigTool) {
	nextToken := ""
	for {
		output, err := i.svc.DescribeConfigRules(&configservice.DescribeConfigRulesInput{
			NextToken: &nextToken,
		})
		if err != nil {
			sLog.WithFields(logrus.Fields{"Test": "Config"}).Info("Error in getConfigRules: ", err)
		}
		fmt.Println("Rule: ", output.ConfigRules)
		sLog.WithFields(logrus.Fields{"Test": "Config"}).Info("getConfigRules: ", output.ConfigRules)
		i.configRules = append(i.configRules, output.ConfigRules...)

		if output.NextToken == nil {
			break
		}
		nextToken = *output.NextToken
	}
}

func (i *ConfigTool) Run() {
	sLog.WithFields(logrus.Fields{
		"Test": "Config"}).Info("ConfigTool Run...")
	getConfigRules(i)
	getComplianceDetails(i)
}
