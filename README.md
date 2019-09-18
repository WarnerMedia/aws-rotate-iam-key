# AWS-Rotate-IAM-Key

aws-rotate-iam-key makes it easy to rotate your IAM keys whether they be in your ~/.aws/credentials file or else where.

# Features!
  - Single binary with no dependencies. 
  - Runs on Linux, Windows and Mac Os
  - Can replace rotated keys in any file - using sed like methods.
  - Optionaly disables the rotated key.
  
### Requirements to compile - binaries available soon.
    - Go
    - Make

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
Wrote new key pair to /Users/dhamm/.aws/credentials
```
### Key and secret provided on command line and output as text. ( ideal for shell scripting )
```sh
$ ./aws-rotate-iam-key -k AKIAXHGHOG2E5ZZ7W7XN -s 3OwEwnWJJqjODdpNTHn6QbN6HiQPvvOUEX8cFIVK
AKIAXHGHOG2E7ZBSISP5 0tujNNpqUt8Fibw0I/TPf6RY2RiFWSzuO18YZpS9
```
### Key and secret provided on command line and output as json. ( handy for use with in languages like python and ruby)
```
$ ./aws-rotate-iam-key -k AKIAXHGHOG2E7ZBSISP5 -s 0tujNNpqUt8Fibw0I/TPf6RY2RiFWSzuO18YZpS9 -o json
{ "AccessKeyId": "AKIAXHGHOG2E3UBGCZC7", "SecretAccessKey": "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" }
```
### Rotate and write new creds to any file format. ( may have limitations on file size.  please limit to a few megs )
```
$ ./aws-rotate-iam-key -k AKIAXHGHOG2E7ZBSISP5 -s 0tujNNpqUt8Fibw0I/TPf6RY2RiFWSzuO18YZpS9 -o /path/to/config.json
```
### Rotate and diable. 
```
$ ./aws-rotate-iam-key -k AKIAXHGHOG2E7ZBSISP5 -s 0tujNNpqUt8Fibw0I/TPf6RY2RiFWSzuO18YZpS9 -d
AKIAXHGHOG2E7ZBSISP5 0tujNNpqUt8Fibw0I/TPf6RY2RiFWSzuO18YZpS9
```