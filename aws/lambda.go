package awshelper


import (
	"context"
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"io/ioutil"
	"os"
)

type LambdaAPI interface {
	CreateFunction(ctx context.Context, params *lambda.CreateFunctionInput, optFns ...func(*lambda.Options)) (*lambda.CreateFunctionOutput, error)
}

type Lambda struct {
	client *lambda.Client
}

func NewLambda(client *lambda.Client) *Lambda {
	return &Lambda{
		client: client,
	}
}
func CreateFunction(c context.Context, api LambdaAPI, input *lambda.CreateFunctionInput) (*lambda.CreateFunctionOutput, error) {
	return api.CreateFunction(c, input)
}
func (l *Lambda) Create()  {
	zipFilePtr := flag.String("z", "", "The name of the ZIP file (without the .zip extension)")
	bucketPtr := flag.String("b", "", "the name of bucket to which the ZIP file is uploaded")
	functionPtr := flag.String("f", "", "The name of the Lambda function")
	handlerPtr := flag.String("h", "", "The name of the package.class handling the call")
	resourcePtr := flag.String("a", "", "The ARN of the role that calls the function")
	runtimePtr := flag.String("r", "", "The runtime for the function.")

	flag.Parse()

	if *zipFilePtr == "" || *bucketPtr == "" || *functionPtr == "" || *handlerPtr == "" || *resourcePtr == "" || *runtimePtr == "" {
		fmt.Println("You must supply a zip file name, bucket name, function name, handler, ARN, and runtime.")
		os.Exit(0)
	}
	_, err := ioutil.ReadFile(*zipFilePtr + ".zip")
	if err != nil {
		fmt.Println("Could not read " + *zipFilePtr + ".zip")
		os.Exit(0)
	}


}