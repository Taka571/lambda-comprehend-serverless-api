package main

import (
	"os"
	"os/exec"

	"github.com/aws/aws-cdk-go/awscdk"
	"github.com/aws/aws-cdk-go/awscdk/awsapigateway"
	"github.com/aws/aws-cdk-go/awscdk/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/awslambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/constructs-go/constructs/v3"
	"github.com/aws/jsii-runtime-go"
)

type PrototypeCdkGoStackProps struct {
	awscdk.StackProps
}

func NewComprehendLambdaApiStack(scope constructs.Construct, id string, props *PrototypeCdkGoStackProps) (awscdk.Stack, error) {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	Cmd := exec.Command("go", "build", "-o", "bin/handler/main", "lambda/main.go")
	Cmd.Env = append(os.Environ(), "GOOS=linux", "CGO_ENABLED=0")
	_, err := Cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	var basicExecutionPolicy awsiam.IManagedPolicy = awsiam.ManagedPolicy_FromAwsManagedPolicyName(aws.String("service-role/AWSLambdaBasicExecutionRole"))
	var comprehendFullAccessPolicy awsiam.IManagedPolicy = awsiam.ManagedPolicy_FromAwsManagedPolicyName(aws.String("ComprehendFullAccess"))
	var role awsiam.Role = awsiam.NewRole(stack, jsii.String("lambda-role"), &awsiam.RoleProps{
		AssumedBy:       awsiam.NewServicePrincipal(aws.String("lambda.amazonaws.com"), nil),
		ManagedPolicies: &[]awsiam.IManagedPolicy{basicExecutionPolicy, comprehendFullAccessPolicy},
	})

	var function awslambda.Function = awslambda.NewFunction(stack, jsii.String("comprehend-lambda"), &awslambda.FunctionProps{
		FunctionName: jsii.String("comprehend-function"),
		Runtime:      awslambda.Runtime_GO_1_X(),
		Code:         awslambda.Code_Asset(jsii.String("bin/handler/")),
		Handler:      jsii.String("main"),
		Role:         role,
	})

	var api awsapigateway.LambdaRestApi = awsapigateway.NewLambdaRestApi(stack, jsii.String("lambda-rest-api"), &awsapigateway.LambdaRestApiProps{
		RestApiName:      jsii.String("comprehend-api"),
		Handler:          function,
		ApiKeySourceType: awsapigateway.ApiKeySourceType_HEADER,
		EndpointTypes:    &[]awsapigateway.EndpointType{awsapigateway.EndpointType_REGIONAL},
		DefaultMethodOptions: &awsapigateway.MethodOptions{
			ApiKeyRequired: jsii.Bool(true),
		},
		Proxy: jsii.Bool(true),
	})

	apiKey := awsapigateway.NewApiKey(stack, jsii.String("api-key"), &awsapigateway.ApiKeyProps{
		ApiKeyName: jsii.String("comprehend-api-key"),
		Enabled:    jsii.Bool(true),
	})

	slice := []*awsapigateway.UsagePlanPerApiStage{}
	ApiStages := append(slice, &awsapigateway.UsagePlanPerApiStage{
		Api:   api,
		Stage: api.DeploymentStage(),
	})

	usagePlan := api.AddUsagePlan(jsii.String("api-usage-plan"), &awsapigateway.UsagePlanProps{
		Name:      jsii.String("comprehend-usage-plan"),
		ApiStages: &ApiStages,
	})

	usagePlan.AddApiKey(apiKey, &awsapigateway.AddApiKeyOptions{OverrideLogicalId: jsii.String("comprehendApiKey")})

	return stack, nil
}

func main() {
	app := awscdk.NewApp(nil)

	_, err := NewComprehendLambdaApiStack(app, "comprehendLambdaApiStack", &PrototypeCdkGoStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	if err != nil {
		panic(err)
	}

	app.Synth(nil)
}

// env determines the AWS environment (account+region) in which our stack is to
// be deployed. For more information see: https://docs.aws.amazon.com/cdk/latest/guide/environments.html
func env() *awscdk.Environment {
	// If unspecified, this stack will be "environment-agnostic".
	// Account/Region-dependent features and context lookups will not work, but a
	// single synthesized template can be deployed anywhere.
	//---------------------------------------------------------------------------
	return nil

	// Uncomment if you know exactly what account and region you want to deploy
	// the stack to. This is the recommendation for production stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String("123456789012"),
	//  Region:  jsii.String("us-east-1"),
	// }

	// Uncomment to specialize this stack for the AWS Account and Region that are
	// implied by the current CLI configuration. This is recommended for dev
	// stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
	//  Region:  jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
	// }
}
