package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/gorilla/mux"
	"github.com/jeremyKisner/streaming-daemon/internal/handler"
)

func ListTables() {
	cfg := aws.Config{
		Region:                        aws.String("us-west-2"),
		Endpoint:                      aws.String("http://localhost:8000"),
		CredentialsChainVerboseErrors: aws.Bool(false),
	}
	sess := session.Must(session.NewSession(&cfg))
	svc := dynamodb.New(sess)

	input := &dynamodb.ListTablesInput{}
	result, err := svc.ListTablesWithContext(context.Background(), input)
	if err != nil {
		fmt.Println("Error listing tables:", err)
		return
	}

	fmt.Println("Tables:")
	for _, tableName := range result.TableNames {
		fmt.Println(*tableName)
	}
}

func main() {
	port := ":8080"
	ListTables()
	r := mux.NewRouter()
	r.HandleFunc("/", handler.Healthz)
	r.HandleFunc("/beepstream", handler.BeepStream)
	fmt.Printf("server started at http://localhost%s/\n", port)
	http.ListenAndServe(port, r)
}
