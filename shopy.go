package main

import (
	"os"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"

	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type ShopyStackProps struct {
	awscdk.StackProps
}

func NewShopyStack(scope constructs.Construct, id string, props *ShopyStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}

	stack := awscdk.NewStack(scope, &id, &sprops)
	s3bucket := awss3.NewBucket(stack, jsii.String("ShopyBucket"), &awss3.BucketProps{
		BucketName: jsii.String("shopy-bucket"),
		BlockPublicAccess: awss3.NewBlockPublicAccess(&awss3.BlockPublicAccessOptions{
			BlockPublicAcls:       jsii.Bool(false),
			BlockPublicPolicy:     jsii.Bool(false),
			IgnorePublicAcls:      jsii.Bool(false),
			RestrictPublicBuckets: jsii.Bool(false),
		}),
		AccessControl: awss3.BucketAccessControl_BUCKET_OWNER_FULL_CONTROL,
		RemovalPolicy: awscdk.RemovalPolicy_RETAIN,
	})

	s3bucket.AddToResourcePolicy(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Actions:    jsii.Strings("s3:GetObject"),
		Effect:     awsiam.Effect_ALLOW,
		Principals: &[]awsiam.IPrincipal{awsiam.NewAnyPrincipal()},
		Resources:  jsii.Strings(*s3bucket.BucketArn() + "/*"),
	}))

	restapi := awsapigateway.NewRestApi(stack, jsii.String("ShopyApigateway"), &awsapigateway.RestApiProps{
		RestApiName: jsii.String("shopy-restapi"),
		DefaultCorsPreflightOptions: &awsapigateway.CorsOptions{
			AllowHeaders: awsapigateway.Cors_DEFAULT_HEADERS(),
			AllowMethods: awsapigateway.Cors_ALL_METHODS(),
			AllowOrigins: awsapigateway.Cors_ALL_ORIGINS(),
		},
		DeployOptions: &awsapigateway.StageOptions{
			StageName: jsii.String("prod"),
		},
	})

	version := restapi.Root().AddResource(jsii.String("v1"), nil)

	NewCategoryStack(stack, &CategoryStackProps{
		StackProps: sprops,
		s3bucket:   s3bucket,
		version:    version,
	})

	NewProductStack(stack, &ProductStackProps{
		StackProps: sprops,
		s3bucket:   s3bucket,
		version:    version,
	})

	NewUserStack(stack, &UserStackProps{
		StackProps: sprops,
		version:    version,
	})

	return stack
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	NewShopyStack(app, "ShopyStack", &ShopyStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

// env determines the AWS environment (account+region) in which our stack is to
// be deployed. For more information see: https://docs.aws.amazon.com/cdk/latest/guide/environments.html
func env() *awscdk.Environment {
	// If unspecified, this stack will be "environment-agnostic".
	// Account/Region-dependent features and context lookups will not work, but a
	// single synthesized template can be deployed anywhere.
	//---------------------------------------------------------------------------
	// return nil

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
	return &awscdk.Environment{
		Account: jsii.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
		Region:  jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
	}
}
