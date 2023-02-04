//create the aws vpc
package main

import (
"github.com/aws/aws-sdk-go/aws"
"github.com/aws/aws-sdk-go/aws/session"
"github.com/aws/aws-sdk-go/service/ec2"
)

func main() {
sess := session.Must(session.NewSessionWithOptions(session.Options { 
    SharedConfigState: session.SharedConfigEnable,
}))

svc := ec2.New(sess) // Create EC2 service

// Create the VPC
createVpcOut, err := svc.CreateVpc(&ec2.CreateVpcInput {
    CidrBlock: aws.String("10.0.0.0/16"),
    InstanceTenancy: aws.String("default"),
})
if err != nil {
    fmt.Println(err.Error())
    return
}
vpcID := *createVpcOut.Vpc.VpcId

// Add name tag to VPC
_, err := svc.CreateTags(&ec2.CreateTagsInput {
    Resources: []*string{vpcID}, 
    Tags: []*ec2.Tag {
        {
            Key:   aws.String("Name"),
            Value: aws.String("development"),
        },
    },
})
if err != nil {
    fmt.Println(err.Error())
    return
}

fmt.Println("Successfully created vpc " + vpcID)

//create the EKS cluster
package main

import (
    "context"
    "fmt"
    eks "github.com/aws/aws-sdk-go/service/eks"
)

func main() {
// Create EKS client
sess := session.Must(session.NewSession(&aws.Config{
    Region: aws.String("us-east-1"),
}))

svc := eks.New(sess)

// Create the EKS cluster
createClusterOut, err := svc.CreateCluster(&eks.CreateClusterInput{
    Name: aws.String("development"),
    ResourcesVpcConfig: &eks.VpcConfigRequest{
        SubnetIds:     [],
        SecurityGroupIds: [],
        EndpointPublicAccess:   aws.Bool(true),
        EndpointPrivateAccess:  aws.Bool(false),
    },
    RoleArn: aws.String("arn:aws:iam::0000000000000:role/role_name"),
})

if err != nil {
    fmt.Println(err.Error())
    return
}
clusterName := *createClusterOut.Cluster.ClusterName

//wait for the cluster to be created
err = svc.WaitUntilClusterActive(context.Background(), &eks.DescribeClusterInput{
    Name: clusterName,
})
if err != nil {
    fmt.Println(err.Error())
    return
}

fmt.Println("Successfully created cluster " + clusterName)

//use helm and bitnami chart to deploy harbor registry
package main

import (
    "fmt"
    "os/exec"
)

func main() {
    cmd := exec.Command("helm", "install", "harbor", "--set", "global.ingress.enabled=true", "--values", "values.yaml!", "bitnami/harbor")
    out, err := cmd.CombinedOutput()
    if err ! = nil {
        fmt.Println("error: "+err.Error())
        return
    }
    fmt.Println("Successfully installed Harbor Registry")
    fmt.Printf("%!s(MISSING)\n", out)
}
