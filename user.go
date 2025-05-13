package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"

	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type UserStackProps struct {
	awscdk.StackProps
	version awsapigateway.Resource
}

func NewUserStack(stack constructs.Construct, props *UserStackProps) {
	table := awsdynamodb.NewTable(stack, jsii.String("UserDynamodb"), &awsdynamodb.TableProps{
		TableName:     jsii.String("user"),
		RemovalPolicy: awscdk.RemovalPolicy_DESTROY,
		PartitionKey: &awsdynamodb.Attribute{
			Name: jsii.String("email"),
			Type: awsdynamodb.AttributeType_STRING,
		},
	})

	lambdaFunc := awslambda.NewFunction(stack, jsii.String("UserLambda"), &awslambda.FunctionProps{
		FunctionName: jsii.String("manage-user"),
		Code:         awslambda.AssetCode_FromAsset(jsii.String("./user/assets/lambda.zip"), nil),
		Runtime:      awslambda.Runtime_PROVIDED_AL2023(),
		Architecture: awslambda.Architecture_ARM_64(),
		MemorySize:   jsii.Number(128),
		Handler:      jsii.String("bootstrap"),
		Timeout:      awscdk.Duration_Seconds(jsii.Number(30)),
		Environment: &map[string]*string{
			"TOKEN_KEY": jsii.String("secret"),
			"TOKEN_EXP": jsii.String("24"),
		},
	})
	table.GrantReadWriteData(lambdaFunc)

	var (
		users      = props.version.AddResource(jsii.String("users"), nil)
		usersEmail = users.ResourceForPath(jsii.String("{email}"))
		options    = &awsapigateway.LambdaIntegrationOptions{
			AllowTestInvoke: jsii.Bool(false),
		}
	)

	users.AddMethod(jsii.String("POST"), awsapigateway.NewLambdaIntegration(lambdaFunc, options), nil)
	users.AddMethod(jsii.String("PUT"), awsapigateway.NewLambdaIntegration(lambdaFunc, options), nil)
	usersEmail.AddMethod(jsii.String("DELETE"), awsapigateway.NewLambdaIntegration(lambdaFunc, options), nil)
}
