package main

import (
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/comprehend"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	svc := comprehend.New(session.New(), &aws.Config{
		Region: aws.String("ap-northeast-1"),
	})

	text, ok := request.QueryStringParameters["text"]

	if !ok {
		return events.APIGatewayProxyResponse{
			Body:       "text is not specified.",
			StatusCode: 400,
		}, nil
	}

	log.Printf("TEXT: %s", text)

	languageCode, ok := request.QueryStringParameters["lg"]

	if !ok {
		languageCode = "ja"
	}

	input := &comprehend.DetectSentimentInput{
		LanguageCode: aws.String(languageCode),
		Text:         aws.String(text),
	}
	res, err := svc.DetectSentiment(input)

	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       res.String(),
			StatusCode: 500,
		}, err
	}

	return events.APIGatewayProxyResponse{
		Body:       res.String(),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}
