package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/comprehend"
)

func detectSentiment(ctx context.Context) (*comprehend.DetectSentimentOutput, error) {
	svc := comprehend.New(session.New(), &aws.Config{
		Region: aws.String("ap-northeast-1"),
	})

	input := &comprehend.DetectSentimentInput{
		LanguageCode: aws.String("ja"),
		Text:         aws.String("今日はいい日だ"),
	}
	res, err := svc.DetectSentiment(input)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func main() {
	lambda.Start(detectSentiment)
}
