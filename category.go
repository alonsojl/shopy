package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"

	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type CategoryStackProps struct {
	awscdk.StackProps
	s3bucket awss3.Bucket
	version  awsapigateway.Resource
}

func NewCategoryStack(stack constructs.Construct, props *CategoryStackProps) {
	table := awsdynamodb.NewTable(stack, jsii.String("CategoryDynamodb"), &awsdynamodb.TableProps{
		TableName:     jsii.String("category"),
		RemovalPolicy: awscdk.RemovalPolicy_DESTROY,
		PartitionKey: &awsdynamodb.Attribute{
			Name: jsii.String("uuid"),
			Type: awsdynamodb.AttributeType_STRING,
		},
	})

	lambdaFunc := awslambda.NewFunction(stack, jsii.String("CategoryLambda"), &awslambda.FunctionProps{
		FunctionName: jsii.String("manage-category"),
		Code:         awslambda.AssetCode_FromAsset(jsii.String("./category/assets/lambda.zip"), nil),
		Runtime:      awslambda.Runtime_PROVIDED_AL2023(),
		Architecture: awslambda.Architecture_ARM_64(),
		MemorySize:   jsii.Number(128),
		Handler:      jsii.String("bootstrap"),
		Timeout:      awscdk.Duration_Seconds(jsii.Number(30)),
		Environment: &map[string]*string{
			"BUCKET_NAME": props.s3bucket.BucketName(),
		},
	})

	table.GrantReadWriteData(lambdaFunc)
	props.s3bucket.GrantReadWrite(lambdaFunc, nil)

	var (
		categories     = props.version.AddResource(jsii.String("categories"), nil)
		categoriesUuid = categories.ResourceForPath(jsii.String("{uuid}"))
		options        = &awsapigateway.LambdaIntegrationOptions{
			AllowTestInvoke: jsii.Bool(false),
		}
	)
	categories.AddMethod(jsii.String("GET"), awsapigateway.NewLambdaIntegration(lambdaFunc, options), nil)
	categories.AddMethod(jsii.String("POST"), awsapigateway.NewLambdaIntegration(lambdaFunc, options), nil)
	categoriesUuid.AddMethod(jsii.String("DELETE"), awsapigateway.NewLambdaIntegration(lambdaFunc, options), nil)
}
