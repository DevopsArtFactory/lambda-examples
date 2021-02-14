package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
)

var AwsRegion = "ap-northeast-2"
var DBClusters = os.Getenv("db_clusters")

var EMPTYSTRING = ""

func getAwsSession() *session.Session {
	mySession := session.Must(session.NewSession(&aws.Config{Region: aws.String(AwsRegion)}))
	return mySession
}

func startDBCluster(session *session.Session, dbClusterID string) error {
	svc := rds.New(session)

	input := &rds.StartDBClusterInput{
		DBClusterIdentifier: aws.String(dbClusterID),
	}

	_, err := svc.StartDBCluster(input)

	if err != nil {
		return err
	}

	return nil
}

func stopDBCluster(session *session.Session, dbClusterID string) error {
	svc := rds.New(session)

	input := &rds.StopDBClusterInput{
		DBClusterIdentifier: aws.String(dbClusterID),
	}

	_, err := svc.StopDBCluster(input)

	if err != nil {
		return err
	}

	return nil
}

func getModeFromEventResource(eventResources []string) string {
	for _, resource := range eventResources {
		if strings.Contains(resource, "startDBClusterEvent") {
			return "startDBClusterEvent"
		} else if strings.Contains(resource, "stopDBClusterEvent") {
			return "stopDBClusterEvent"
		}
	}

	return EMPTYSTRING
}

func startDBClusterHandler() {
	session := getAwsSession()

	for _, dbCluster := range strings.Split(DBClusters, ",") {
		startDBCluster(session, dbCluster)
	}
}

func stopDBClusterHandler() {
	session := getAwsSession()

	for _, dbCluster := range strings.Split(DBClusters, ",") {
		stopDBCluster(session, dbCluster)
	}
}

func handler(event events.CloudWatchEvent) {
	fmt.Println(event.Resources)
	switch getModeFromEventResource(event.Resources) {
	case "startDBClusterEvent":
		fmt.Println("Start")
		startDBClusterHandler()
	case "stopDBClusterEvent":
		fmt.Println("Stop")
		stopDBClusterHandler()
	default:
		fmt.Println("N/A")
	}
}

func main() {
	lambda.Start(handler)
}