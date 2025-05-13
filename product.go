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

type ProductStackProps struct {
	awscdk.StackProps
	s3bucket awss3.Bucket
	version  awsapigateway.Resource
}

func NewProductStack(stack constructs.Construct, props *ProductStackProps) {
	table := awsdynamodb.NewTable(stack, jsii.String("ProductDynamodb"), &awsdynamodb.TableProps{
		TableName:     jsii.String("product"),
		RemovalPolicy: awscdk.RemovalPolicy_DESTROY,
		PartitionKey: &awsdynamodb.Attribute{
			Name: jsii.String("uuid"),
			Type: awsdynamodb.AttributeType_STRING,
		},
	})

	table.AddGlobalSecondaryIndex(&awsdynamodb.GlobalSecondaryIndexProps{
		IndexName: jsii.String("GSI_CATEGORY"),
		PartitionKey: &awsdynamodb.Attribute{
			Name: jsii.String("category_uuid"),
			Type: awsdynamodb.AttributeType_STRING,
		},
		ProjectionType: awsdynamodb.ProjectionType_ALL,
	})

	table.AddGlobalSecondaryIndex(&awsdynamodb.GlobalSecondaryIndexProps{
		IndexName: jsii.String("GSI_QRCODE"),
		PartitionKey: &awsdynamodb.Attribute{
			Name: jsii.String("qrcode"),
			Type: awsdynamodb.AttributeType_STRING,
		},
		ProjectionType: awsdynamodb.ProjectionType_ALL,
	})

	lambdaFunc := awslambda.NewFunction(stack, jsii.String("ProductLambda"), &awslambda.FunctionProps{
		FunctionName: jsii.String("manage-product"),
		Code:         awslambda.AssetCode_FromAsset(jsii.String("./product/assets/lambda.zip"), nil),
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
		products     = props.version.AddResource(jsii.String("products"), nil)
		productsUuid = products.ResourceForPath(jsii.String("{uuid}"))
		options      = &awsapigateway.LambdaIntegrationOptions{
			AllowTestInvoke: jsii.Bool(false),
		}
	)

	products.AddMethod(jsii.String("GET"), awsapigateway.NewLambdaIntegration(lambdaFunc, options), nil)
	products.AddMethod(jsii.String("POST"), awsapigateway.NewLambdaIntegration(lambdaFunc, options), nil)
	productsUuid.AddMethod(jsii.String("PUT"), awsapigateway.NewLambdaIntegration(lambdaFunc, options), nil)
	productsUuid.AddMethod(jsii.String("DELETE"), awsapigateway.NewLambdaIntegration(lambdaFunc, options), nil)
}
