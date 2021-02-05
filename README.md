# TcTool

**Dev Environment**
Create an EC2 Env in Cloud9
Open ~/.bashrc and update the following 2 lines
    GOPATH=~/environment/go
    export GOPATH
To increase the EBS size, run the cmd - sh resize.sh 20
Run the following to install aws-sdk
    go get -u github.com/aws/aws-sdk-go/...
To download some examples, run the following cmd
    git clone https://github.com/awsdocs/aws-doc-sdk-examples.git

**Features**
- CIS AWS Foundations Benchmark controls - config tests
https://docs.aws.amazon.com/securityhub/latest/userguide/securityhub-cis-controls.html
https://d1.awsstatic.com/whitepapers/compliance/AWS_CIS_Foundations_Benchmark.pdf

- AWS Foundational Security Best Practices controls - best practices
https://docs.aws.amazon.com/securityhub/latest/userguide/securityhub-standards-fsbp-controls.html

- AWS Inspector run - EC2 vuln tests

- AWS Trust Advisor Results

**Build the tool**
git checkout https://github.com/aseemsethi/tctool
go build src/main.go

**Tool Prerequisites**
1) Set the following env variable and the Region in the Config file
$ export AWS_SDK_LOAD_CONFIG="true"
$ cd ~/.aws 
$ more config 
region = us-east-1
The above setting of region is done so that the call to create a new session passes with the right region.
sess, err := session.NewSessionWithOptions

2) Create a AWS Role in your account, and associate to the EC2 that will run the tool.
The Role is described as follows:
                AWSCloudTrailReadOnlyAccess
                IAMReadOnlyAccess
                IAMAccessAdvisorReadOnly
                AmazonInspectorFullAccess
                AmazonVPCReadOnlyAccess
                AmazonSSMManagedInstanceCore
                AmazonS3ReadOnlyAccess
                AmazonVPCReadOnlyAccess
and
                "kms:TagResource",
                "kms:ScheduleKeyDeletion",
                "kms:PutKeyPolicy",
                "kms:CreateKey",
                "kms:ListResourceTags",
                "kms:CreateGrant"
    
3) For AWS Inspector run, attach the following Role to every EC2 instance. This
allows access of the SSM agent on EC2 to communicate with EC2 Systems Manager..
Also added to the Policies above, so no additional step is needed.
_AmazonSSMManagedInstanceCore_
_AmazonInspectorFullAccess_

Also, tag all EC2s where you want inspector to run with tag "inspector" : "true"
Note that all following rules are run - 
Common Vulnerabilities and Exposures-1.1
CIS Operating System Security Configuration Benchmarks-1.0
Network Reachability-1.1
Security Best Practices-1.0

- Need ssm-agent and inspector-agent to be installed in all EC2s.
Run cmd - sudo systemctl status amazon-ssm-agent - to check if ssm-agent is installed and running
Else, install ssm-agent on all EC2s.
- inspector agent is installed automatically by Sytems Nanager, when it is run
- You can do this manually too; chose a 15 min run to see some results

