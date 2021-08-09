package main

import (
	"os"
	"os/exec"

	"github.com/aws/aws-cdk-go/awscdk"
	"github.com/aws/aws-cdk-go/awscdk/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/awslambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/constructs-go/constructs/v3"
	"github.com/aws/jsii-runtime-go"
)

type PrototypeCdkGoStackProps struct {
	awscdk.StackProps
}

func NewComprehendLambdaStack(scope constructs.Construct, id string, props *PrototypeCdkGoStackProps) (awscdk.Stack, error) {
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
	var role awsiam.Role = awsiam.NewRole(stack, jsii.String("go-cdk-lambda-role"), &awsiam.RoleProps{
		AssumedBy:       awsiam.NewServicePrincipal(aws.String("lambda.amazonaws.com"), nil),
		ManagedPolicies: &[]awsiam.IManagedPolicy{basicExecutionPolicy, comprehendFullAccessPolicy},
	})

	awslambda.NewFunction(stack, jsii.String("go-cdk-comprehend-lambda"), &awslambda.FunctionProps{
		FunctionName: jsii.String("go-cdk-comprehend-function"),
		Runtime:      awslambda.Runtime_GO_1_X(),
		Code:         awslambda.Code_Asset(jsii.String("bin/handler/")),
		Handler:      jsii.String("main"),
		Role:         role,
	})

	return stack, nil
}

func main() {
	app := awscdk.NewApp(nil)

	_, err := NewComprehendLambdaStack(app, "comprehendLambdaStack", &PrototypeCdkGoStackProps{
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
