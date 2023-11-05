package collectors

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/gocolly/colly"
)

const (
	targetURL              = "https://www.nnmarathonrotterdam.nl/info/vraag-en-aanbod-platform/"
	dynamoDBEntryTableName = "nn_entries"
	dynamoDBFormTableName  = "nn_forms"
)

type NNentry struct {
	DataRecordId   string  `json:"dataRecordId"`
	Datum          string  `json:"datum"`
	ID             int     `json:"id"`
	Onderdeel      string  `json:"onderdeel"`
	ExtraProducten string  `json:"extraProducten,omitempty"`
	Prijs          float64 `json:"prijs"`
}

type NNform struct {
	RequestVerificationToken string `json:"__RequestVerificationToken"`
	FormId                   string `json:"FormId"`
	FormName                 string `json:"FormName"`
	RecordId                 string `json:"RecordId"`
	PreviousClicked          string `json:"PreviousClicked"`
	FormStep                 int    `json:"FormStep"`
	RecordState              string `json:"RecordState"`
	Ufprt                    string `json:"ufprt"`

	// Variable Form Names
	EventID         string `json:"eventid"`
	Firstname       string `json:"firstname"`
	Lastname        string `json:"lastname"`
	Email           string `json:"email"`
	Disclaimer      string `json:"disclaimer"`
	RelatedRecordID string `json:"relatedrecordid"`
}

type NNformUpdate struct {
	Ufprt                    string `json:":uf"`
	RequestVerificationToken string `json:":rv"`
	RecordState              string `json:":rs"`
}

func HandleLambdaEvent() error {
	c := colly.NewCollector()

	var entries []NNentry
	var currentForm NNform

	//
	// First: Scrape all info from target site
	//
	c.OnHTML("div.marketplace-buy", func(s *colly.HTMLElement) {
		// First Part: Scrape all NN entries from HTML buy Table
		s.ForEach(".marketplace-buy-table > tbody", func(i int, e *colly.HTMLElement) {
			e.ForEach("tr", func(_ int, el *colly.HTMLElement) {
				// Convert id variable naar Int
				id, err := strconv.Atoi(el.ChildText("td:nth-child(2)"))
				if err != nil {
					return
				}

				// Cleanup prijs variable en convert naar float
				tmp := strings.TrimPrefix(el.ChildText("td:nth-child(5)"), "€")
				tmp = strings.TrimSpace(tmp)
				tmp = strings.Replace(tmp, ",", ".", 1)
				prijs, err := strconv.ParseFloat(tmp, 32)
				if err != nil {
					return
				}

				entries = append(entries, NNentry{
					DataRecordId:   el.Attr("data-record-id"),
					Datum:          el.ChildText("td:nth-child(1)"),
					ID:             id,
					Onderdeel:      el.ChildText("td:nth-child(3)"),
					ExtraProducten: el.ChildText("td:nth-child(4)"),
					Prijs:          prijs,
				})
			})
		})

		// Second Part: Scrape all variable form input names and values from HTML buy form
		s.ForEach(".nnmarathonrotterdammarketplacebuy > form", func(i int, f *colly.HTMLElement) {
			// Scrape non-variable form names and values
			currentForm = NNform{
				RequestVerificationToken: f.ChildAttr("input[name='__RequestVerificationToken']", "value"),
				FormId:                   f.ChildAttr("input[name='FormId']", "value"),
				FormName:                 f.ChildAttr("input[name='FormName']", "value"),
				RecordId:                 f.ChildAttr("input[name='RecordId']", "value"),
				PreviousClicked:          f.ChildAttr("input[name='PreviousClicked']", "value"),
				RecordState:              f.ChildAttr("input[name='RecordState']", "value"),
				Ufprt:                    f.ChildAttr("input[name='ufprt']", "value"),
				EventID:                  f.ChildAttr(".eventid input", "id"),
				Firstname:                f.ChildAttr(".marketplaceformsfirstname input", "id"),
				Lastname:                 f.ChildAttr(".marketplaceformslastname input", "id"),
				Email:                    f.ChildAttr(".marketplaceformsemail input", "id"),
				Disclaimer:               f.ChildAttr(".marketplaceformsdisclaimertext input", "id"),
				RelatedRecordID:          f.ChildAttr(".relatedrecordid input", "id"),
			}
		})
	})
	c.Visit(targetURL)

	//
	//	Second: Upload data to dynamoDB
	//
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	svc := dynamodb.New(sess)

	for _, entry := range entries {
		av, err := dynamodbattribute.MarshalMap(entry)
		if err != nil {
			return err
		}

		entryInput := &dynamodb.PutItemInput{
			Item:                av,
			TableName:           aws.String(dynamoDBEntryTableName),
			ConditionExpression: aws.String("attribute_not_exists(dataRecordId)"),
		}
		svc.PutItem(entryInput)
	}

	av, err := dynamodbattribute.MarshalMap(currentForm)
	if err != nil {
		return err
	}

	entryInput := &dynamodb.PutItemInput{
		Item:                av,
		TableName:           aws.String(dynamoDBFormTableName),
		ConditionExpression: aws.String("attribute_not_exists(FormId)"),
	}
	_, err = svc.PutItem(entryInput)
	if err != nil {
		if err.Error() == "ConditionalCheckFailedException: The conditional request failed" {
			// If form already exists, only update the differential fields
			update, err := dynamodbattribute.MarshalMap(NNformUpdate{
				Ufprt:                    currentForm.Ufprt,
				RequestVerificationToken: currentForm.RequestVerificationToken,
				RecordState:              currentForm.RecordState,
			})
			if err != nil {
				fmt.Println(err.Error())
			}

			input := &dynamodb.UpdateItemInput{
				Key: map[string]*dynamodb.AttributeValue{
					"FormId": {
						S: &currentForm.FormId,
					},
				},
				TableName:                 aws.String(dynamoDBFormTableName),
				UpdateExpression:          aws.String("set ufprt = :uf, #rv = :rv, RecordState = :rs"),
				ExpressionAttributeValues: update,
				ExpressionAttributeNames: map[string]*string{
					"#rv": aws.String("__RequestVerificationToken"),
				},
				ReturnValues: aws.String("UPDATED_NEW"),
			}
			_, err = svc.UpdateItem(input)
			if err != nil {
				fmt.Println(err.Error())
			}
		}
	}
	return nil
}

func main() {
	lambda.Start(HandleLambdaEvent)
}
