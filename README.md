# AWS-Rotate-IAM-Key

aws-rotate-iam-key makes it easy to rotate your IAM keys whether they be in your ~/.aws/credentials file or else where.

This work is a deritive of https://github.com/stefansundin/aws-rotate-key.  Added flexibility to simplify key rotation from within a shell script. 

# Features!
  - Single binary with no dependencies. 
  - Runs on Linux, Windows and Mac Os
  - Can replace rotated keys in any file - using sed like methods.
  - Optionaly disables the rotated key.
  
# Requirements 
#### to compile - binaries available soon.
    - Go
    - Make
#### AWS Policy to apply to IAM user
Please note, the policy below specifies a condition liming access from specific IP addresses.  The IP addresses will have to be
changed before the policy can work in your environment.

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Action": [
                "iam:*AccessKey*",
                "iam:ChangePassword",
                "iam:GetUser",
                "iam:*ServiceSpecificCredential*",
                "iam:*SigningCertificate*"
            ],
            "Resource": [
                "arn:aws:iam::AwsAccountIDGoesHere:user/*${aws:username}"
            ],
            "Effect": "Allow",
            "Sid": "AllowRotateOwnKey",
            "Condition": {
                "IpAddress": { "aws:SourceIp": [ "52.72.140.215/32", "52.200.112.43/32", "52.200.113.35/32" ] }
            }
        }
    ]
}
``` 

# Installation
    1.  Download or clone repo.  
    2.  cd into repo
    3.  type make


# Usage:
```sh
Usage of ./aws-rotate-iam-key:
  -c string
    	AWS credentials file
  -d	Disable key after rotation.
  -k string
    	AWS IAM key.
  -o string
    	Output format - default is text, option json is json string, /path/to/file runs a regex on the file specified.
  -profile string
    	Named profile within AWS credentials file.
  -s string
    	AWS IAM secret
  -v	version 1.0.0 
    	built 2019-09-18T16:13:25-0400 
    	git repo = https://atom-git.turner.com/best-meta-aws/cloudutil/aws-rotate-iam-key
```
    
# Examples of use:	
### Updating a key within ~/.aws/credentials, referenced by profile
```sh
$ ./aws-rotate-iam-key -profile dch
Wrote new key pair to /Users/$UserName/.aws/credentials
```
### Key and secret provided on command line and output as text. ( ideal for shell scripting )
```sh
$ ./aws-rotate-iam-key -k AKIAXXXXX -s xxxxxxxxxxxxxx
AKIAXXXXX xxxxxxxxxxxxxx
```
### Key and secret provided on command line and output as json. ( handy for use with in languages like python and ruby)
```
$ ./aws-rotate-iam-key -k AKIAXXXXX -s xxxxxxxxxxxxxx -o json
{ "AccessKeyId": "AKIAXXXXX", "SecretAccessKey": "xxxxxxxxxxxx" }
```
### Rotate and write new creds to any file format. ( may have limitations on file size.  please limit to a few megs )
```
$ ./aws-rotate-iam-key -k AKIAXXXXX -s xxxxxxxxxxxxxx -o /path/to/config.json
```
### Rotate and diable. 
```
$ ./aws-rotate-iam-key -k AKIAXXXXX -s xxxxxxxxxxxxxx -d
AKIAXXXXX xxxxxxxxxxxxxx
```
### Rotate credentials held in MySQL in a cron job
```
ORIGCREDS=`echo "use mydb; select awskey,awssecret from users where u_login like 'mickeymouse'" | mysql | tail -n 1`
AWSKEY=`echo $ORIGCREDS | awk '{ print $1 }'`
AWSSEC=`echo $ORIGCREDS | awk '{ print $2 }'`
NEWCREDS=$(aws-rotate-iam-key -k envAWSKEY -s envAWSSEC)
NEWKEY=`echo $NEWCREDS | awk '{ print $1 }'
NEWSEC=`echo $NEWCREDS | awk '{print $2 }'
echo "use mydb; update users set awskey=$NEWKEY,awssecret=$NEWSEC where u_login like 'mickeymouse'" | mysql
```