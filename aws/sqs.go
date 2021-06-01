package awshelper

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

// SQSCreateQueueAPI defines the interface for the CreateQueue function.
// We use this interface to test the function using a mocked service.
type SQSQueueAPI interface {
	CreateQueue(ctx context.Context, params *sqs.CreateQueueInput, optFns ...func(*sqs.Options)) (*sqs.CreateQueueOutput, error)
	GetQueueUrl(ctx context.Context, params *sqs.GetQueueUrlInput, optFns ...func(*sqs.Options)) (*sqs.GetQueueUrlOutput, error)
	SendMessage(ctx context.Context, params *sqs.SendMessageInput, optFns ...func(*sqs.Options)) (*sqs.SendMessageOutput, error)
	ReceiveMessage(ctx context.Context, params *sqs.ReceiveMessageInput, optFns ...func(*sqs.Options)) (*sqs.ReceiveMessageOutput, error)
}

type SQS struct {
	client *sqs.Client
}

func NewSQS(client *sqs.Client) *SQS {
	return &SQS{
		client: client,
	}
}

// CreateQueue creates an Amazon SQS queue.
// Inputs:
//     c is the context of the method call, which includes the AWS Region.
//     api is the interface that defines the method call.
//     input defines the input arguments to the service call.
// Output:
//     If success, a CreateQueueOutput object containing the result of the service call and nil.
//     Otherwise, nil and an error from the call to CreateQueue.
func CreateQueue(c context.Context, api SQSQueueAPI, input *sqs.CreateQueueInput) (*sqs.CreateQueueOutput, error) {
	return api.CreateQueue(c, input)
}
func GetQueueURL(c context.Context, api SQSQueueAPI, input *sqs.GetQueueUrlInput) (*sqs.GetQueueUrlOutput, error) {
	return api.GetQueueUrl(c, input)
}
func SendMsg(c context.Context, api SQSQueueAPI, input *sqs.SendMessageInput) (*sqs.SendMessageOutput, error) {
	return api.SendMessage(c, input)
}
func GetMessages(c context.Context, api SQSQueueAPI, input *sqs.ReceiveMessageInput) (*sqs.ReceiveMessageOutput, error) {
	return api.ReceiveMessage(c, input)
}

//Create New Queue
func (s *SQS) CreateNewQueue(queue *string) (string, error) {

	if *queue == "" {
		return "", errors.New("you must supply a queue Name")
	}

	inputSearch := &sqs.GetQueueUrlInput{
		QueueName: queue,
	}
	//check if the que name already exists before
	find, err := GetQueueURL(context.TODO(), s.client, inputSearch)
	if err == nil && *find.QueueUrl != "" {
		return *find.QueueUrl, errors.New(fmt.Sprintf("Que name already exist"))
	}

	input := &sqs.CreateQueueInput{
		QueueName: queue,
		Attributes: map[string]string{
			"DelaySeconds":           "60",
			"MessageRetentionPeriod": "86400",
		},
	}

	result, err := CreateQueue(context.TODO(), s.client, input)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Got an error creating the queue:%s", err.Error()))
	}

	return *result.QueueUrl, nil
}

//Get  Queue ID
func (s *SQS) GetQueueID(queue string) (string, error) {

	if queue == "" {
		return "", errors.New("you must supply a queue Name")
	}

	input := &sqs.GetQueueUrlInput{
		QueueName: &queue,
	}
	find, err := GetQueueURL(context.TODO(), s.client, input)
	if err != nil {
		fmt.Println("Got an error getting the queue URL:")
		fmt.Println(err)
		return "", errors.New(fmt.Sprintf("Got an error getting the queue URL:%s", err.Error()))
	}

	return *find.QueueUrl, nil
}

func (s *SQS) SendMessage(queue *string, message string) (string, error) {
	if *queue == "" {
		return "", errors.New("you must supply a queue Name")
	}

	// Get URL of queue
	input := &sqs.GetQueueUrlInput{
		QueueName: queue,
	}

	result, err := GetQueueURL(context.TODO(), s.client, input)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Got an error getting the queue URL:%s", err.Error()))
	}
	queueURL := result.QueueUrl

	sMInput := &sqs.SendMessageInput{
		DelaySeconds: 10,
		MessageAttributes: map[string]types.MessageAttributeValue{
			"Title": {
				DataType:    aws.String("String"),
				StringValue: aws.String("Send Employee Info"),
			},
			"Author": {
				DataType:    aws.String("String"),
				StringValue: aws.String("Ahmed Mahmoud"),
			},
			"WeeksOn": {
				DataType:    aws.String("Number"),
				StringValue: aws.String("6"),
			},
		},
		MessageBody: aws.String(message),
		QueueUrl:    queueURL,
	}

	resp, err := SendMsg(context.TODO(), s.client, sMInput)
	if err != nil {

		return "", errors.New("Configuration error, " + err.Error())
	}

	return *resp.MessageId, nil
}

//Recieve Message
//Input
//		queue   >> queue name
//		timeout >> timeout in seconds
func (s *SQS) Recievemessage(queue string, timeout int32) (string, error) {

	if queue == "" {
		fmt.Println("You must supply the name of a queue  ")
		return "", errors.New("you must supply the name of a queue ")
	}

	if timeout < 0 {
		timeout = 0
	}

	if timeout > 12*60*60 {
		timeout = 12 * 60 * 60
	}

	gQInput := &sqs.GetQueueUrlInput{
		QueueName: &queue,
	}

	// Get URL of queue
	urlResult, err := GetQueueURL(context.TODO(), s.client, gQInput)
	if err != nil {
		return "", errors.New("Got an error getting the queue URL:" + err.Error())
	}

	queueURL := urlResult.QueueUrl

	gMInput := &sqs.ReceiveMessageInput{
		MessageAttributeNames: []string{
			string(types.QueueAttributeNameAll),
		},
		QueueUrl:            queueURL,
		MaxNumberOfMessages: 1,
		VisibilityTimeout:   timeout,
	}

	msgResult, err := GetMessages(context.TODO(), s.client, gMInput)
	if err != nil {
		return "", errors.New("Got an error receiving messages:" + err.Error())
	}
	fmt.Println(*msgResult.Messages[0].Body)
	return *msgResult.Messages[0].ReceiptHandle, nil

}
