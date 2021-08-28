# Example

- `$ npm install -g aws-cdk` If you haven't installed it yet
- `$ cdk deploy comprehendLambdaApiStack --profile [aws profile]`

```
$ curl -G https://api-gateway-id.execute-api.ap-northeast-1.amazonaws.com/prod \
    --header 'x-api-key:your-api-key' \
    --data-urlencode "text=完全に理解した"
{
  Sentiment: "POSITIVE",
  SentimentScore: {
    Mixed: 0.00018109637312591076,
    Negative: 0.00958448089659214,
    Neutral: 0.1103581041097641,
    Positive: 0.8798763155937195
  }
}

$ curl -G https://api-gateway-id.execute-api.ap-northeast-1.amazonaws.com/prod \
    --header 'x-api-key:your-api-key' \
    --data-urlencode "text=何にもわからなかった"
{
  Sentiment: "NEGATIVE",
  SentimentScore: {
    Mixed: 4.52065905847121e-05,
    Negative: 0.6747143268585205,
    Neutral: 0.32068878412246704,
    Positive: 0.004551670514047146
  }
}
```


