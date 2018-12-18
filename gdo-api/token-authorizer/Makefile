# Makefile for demo-authorizer
# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
AWSCMD=aws
BINARY_NAME=token-authorizer-gdo
FUNC_NAME=$(BINARY_NAME)
BINARY_LINUX=$(BINARY_NAME)
BINARY_WINDOWS=$(BINARY_NAME).exe
ZIPPER=zip 
ZIP_NAME=$(BINARY_LINUX).zip
# get a fresh token:
# curl --request POST   --url https://jerman3.auth0.com/oauth/token   --header 'content-type: application/json'   --data '{"client_id":"Os36QLz24NAuwgbCV8TnPqqh19A24AQz","client_secret":"rdyucChuSfeXMs_Zx1yugM-uRaN5GVnFa5_SM6Nb6negQDuRxfiwionCMgCGl2De","audience":"https://gdo.heather.com","grant_type":"client_credentials"}'
TEST_URL=https://l8bh3bz5m3.execute-api.us-east-1.amazonaws.com/demo/page?what

TEST_TOKEN=eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsImtpZCI6Ik16VTVOa00wTjBFeU5EQTBRelkwTkRrME16azVNMFV3TVRBME1URXhOVVV3TnpKQ00wVXhNZyJ9.eyJpc3MiOiJodHRwczovL2plcm1hbjMuYXV0aDAuY29tLyIsInN1YiI6Ik9zMzZRTHoyNE5BdXdnYkNWOFRuUHFxaDE5QTI0QVF6QGNsaWVudHMiLCJhdWQiOiJodHRwczovL2dkby5oZWF0aGVyLmNvbSIsImlhdCI6MTU0NDgwMDc0MCwiZXhwIjoxNTQ0ODg3MTQwLCJhenAiOiJPczM2UUx6MjROQXV3Z2JDVjhUblBxcWgxOUEyNEFReiIsImd0eSI6ImNsaWVudC1jcmVkZW50aWFscyJ9.MiFh42XwU1NduDLXpn-LcFEldj_JtnWjiVcAr3Bih-pL9xJP4Ok39uMW6w7glH2vh1GKB9pQ4essxj91cP_d_RzjCxebahC58tkItUhEjla2-64yO-dDVPYgCZ0oqsIC4uxcgUgz_DGDBwKJD2n2At2xZ5Q_yWcUUP01qvJtqdDHyToHGdjq_6p0k-jmK8oyI4JEoNmjXW37xrfVWbHUczfZucyfk0KGbejGpHZxG5m4YBXBoT1VlzNGDsTyivWMzlCXAiCV887phZ9fHm6cCJbeS7PbPliJHiizpAeE0ipol4ttdAWOjtCwDlSMlZCHkT0WSzXA9F8Y0fbiKDCYKQ

.DEFAULT=zip

# test runs a pass
test: $(ZIP_NAME) upload.log
	curl --request GET   --url $(TEST_URL)   --header 'authorization:$(TEST_TOKEN)'

upload.log: $(ZIP_NAME)
	$(AWSCMD) lambda update-function-code --function-name $(FUNC_NAME) --zip-file fileb://$(ZIP_NAME) >upload.log

zip: $(ZIP_NAME)

$(ZIP_NAME): $(BINARY_LINUX)
	$(ZIPPER) -o $(ZIP_NAME) $(BINARY_LINUX)

$(BINARY_LINUX):  *.go
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_LINUX) -v

tidy:
	$(GOCMD) mod tidy