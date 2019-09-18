package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/sts"
)

// AWSKeySes is a struct holding the new key and secret
// along with the session used to allocate the new key.  The sessions
// can be used later to disable the old key
type AWSKeySes struct {
	key, secret string
	iamclient   *iam.IAM
}

func pluralize(n int) string {
	if n == 1 {
		return ""
	}
	return "s"
}

func checkErr(context string, err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, context)
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

// checkArgs test the command line arguments to warrant that either
// a key and secret are provided on the command line or a profile within a credentials ini is referenced
func checkArgs(rtc RTC) {

	if *rtc.Opts["profile"] == "none" {
		if _, ok := rtc.Opts["key"]; !ok {
			fmt.Print("\nError! No AWS IAM Key provided. Please provide key with -k <AKIA*> or by -profile <profile name>.\n\n")
			flag.Usage()
			os.Exit(1)
		} else { // pull from env
			if strings.HasPrefix(*rtc.Opts["key"], "env") {
				rtc.Opts["key"] = osenvvar(strings.TrimPrefix(*rtc.Opts["key"], "env"))
			}
		}
		if _, ok := rtc.Opts["secret"]; !ok {
			fmt.Print("\nError! No AWS IAM Secret provided.  Please provide secret with -s <Secret> or by -profile <profile name>.\n\n")
			flag.Usage()
			os.Exit(1)
		} else { // pull form env
			if strings.HasPrefix(*rtc.Opts["secret"], "env") {
				rtc.Opts["secret"] = osenvvar(strings.TrimPrefix(*rtc.Opts["secret"], "env"))
			}
		}
	} else { // load key and secret from credentials INI also set output to ini if not otherwise specified
		credentialsProvider := credentials.NewSharedCredentials(*rtc.Opts["credsini"], *rtc.Opts["profile"])
		creds, err := credentialsProvider.Get()
		if err != nil {
			fmt.Printf("Error retrieving credentals: %s\n", err)
			os.Exit(1)
		}
		rtc.Opts["key"] = &creds.AccessKeyID
		rtc.Opts["secret"] = &creds.SecretAccessKey
		if *rtc.Opts["output"] == "none" {
			rtc.Opts["output"] = rtc.Opts["credsini"]
		}
		return
	}

}
func rotateCreds(rtc RTC) AWSKeySes {
	creds := credentials.NewStaticCredentials(*rtc.Opts["key"], *rtc.Opts["secret"], "")
	sess, err := session.NewSession(&aws.Config{Credentials: creds})
	if err != nil {
		fmt.Printf("Error authenticating: %s\n", err)
		os.Exit(1)
	}

	stsClient := sts.New(sess)
	respGetCallerIdentity, err := stsClient.GetCallerIdentity(&sts.GetCallerIdentityInput{})
	checkErr("Error getting caller identity. Is the key disabled?", err)

	username := strings.Split(*respGetCallerIdentity.Arn, "/")
	iamClient := iam.New(sess)
	respListAccessKeys, err := iamClient.ListAccessKeys(&iam.ListAccessKeysInput{UserName: &username[len(username)-1]})
	checkErr("Failed to list Access Keys.  Does the user have the necessary IAM policy attached?", err)

	if len(respListAccessKeys.AccessKeyMetadata) == 2 {
		keyIndex := 0
		if *respListAccessKeys.AccessKeyMetadata[0].AccessKeyId == *rtc.Opts["key"] {
			keyIndex = 1
		}

		_, err2 := iamClient.DeleteAccessKey(&iam.DeleteAccessKeyInput{
			AccessKeyId: respListAccessKeys.AccessKeyMetadata[keyIndex].AccessKeyId,
		})
		checkErr("Failed to delete secondary access key!", err2)

	}

	respCreateAccessKey, err := iamClient.CreateAccessKey(&iam.CreateAccessKeyInput{UserName: &username[len(username)-1]})
	checkErr("New key creation failed!", err)

	//return *respCreateAccessKey.AccessKey.AccessKeyId, *respCreateAccessKey.AccessKey.SecretAccessKey
	return AWSKeySes{
		key:       *respCreateAccessKey.AccessKey.AccessKeyId,
		secret:    *respCreateAccessKey.AccessKey.SecretAccessKey,
		iamclient: iamClient,
	}
}
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
func sedCreds(rtc RTC, newcreds AWSKeySes) {
	// Read credentials file
	bytes, err := ioutil.ReadFile(*rtc.Opts["output"])
	checkErr(fmt.Sprintf("Failed to read file %s before replacing key and secret!", *rtc.Opts["output"]), err)
	text := string(bytes)

	re := regexp.MustCompile(fmt.Sprintf(`%s`, regexp.QuoteMeta(*rtc.Opts["key"])))
	text = re.ReplaceAllString(text, newcreds.key)
	re = regexp.MustCompile(fmt.Sprintf(`%s`, regexp.QuoteMeta(*rtc.Opts["secret"])))
	text = re.ReplaceAllString(text, newcreds.secret)

	// Verify that the regexp actually replaced something
	if !strings.Contains(text, newcreds.key) || !strings.Contains(text, newcreds.secret) {
		fmt.Printf("Erorr: Failed to replace old access key in %s", *rtc.Opts["output"])
		os.Exit(1)
	}

	// rewrite file
	err = ioutil.WriteFile(*rtc.Opts["output"], []byte(text), 0600)
	checkErr(fmt.Sprintf("Failed to write file %s before replacing key and secret!", *rtc.Opts["output"]), err)
	fmt.Printf("Wrote new key pair to %s\n", *rtc.Opts["output"])

}

func disableKey(rtc RTC, newcreds AWSKeySes) {

	_, err := newcreds.iamclient.UpdateAccessKey(&iam.UpdateAccessKeyInput{
		AccessKeyId: rtc.Opts["key"],
		Status:      aws.String("Inactive"),
	})
	checkErr("Faiked to disable key!", err)
}

func main() {
	// get run time config
	rtc := Newrtc()

	checkArgs(*rtc)

	NewCredSes := rotateCreds(*rtc)

	if fileExists(*rtc.Opts["output"]) {
		sedCreds(*rtc, NewCredSes)
	} else {
		switch *rtc.Opts["output"] {
		case "json":
			fmt.Printf("{ \"AccessKeyId\": \"%s\", \"SecretAccessKey\": \"%s\" }\n", NewCredSes.key, NewCredSes.secret)
		default:
			fmt.Printf("%s %s\n", NewCredSes.key, NewCredSes.secret)
		}
	}

	if *rtc.Flag["disable"] {
		disableKey(*rtc, NewCredSes)
	}
	os.Exit(0)
}
