import boto3
import logging

logger = logging.getLogger()

def deploy_webserver_ec2(aws_access_key_id,aws_secret_access_key):

    try:
        logging.basicConfig(level=logging.INFO)
        logger.info("Creating EC2 Instance")
        ec2 = boto3.resource('ec2',aws_access_key_id = aws_access_key_id,
            aws_secret_access_key = aws_secret_access_key)
        inst = ec2.create_instances(
            ImageId = 'ami-05df9ea8bb34958cc', 
            MinCount = 1,
            MaxCount = 1,
            InstanceType = 't2.micro',
            KeyName = 'my_key'
        )
        inst[0].wait_until_running()
        inst[0].load()
        inst[0].reload()
        ip_address = inst[0].public_ip_address
        logger.info("Deploying websites to EC2 instance "+ ip_address)
        #once the Instance is created and running, now we can proceed with the Terraform script
        #Create a working directory
        import os
        os.mkdir('Terraform_script')
        #change directory to Terraform_script
        os.chdir('Terraform_script') 
        #Run the Terraform init command 
        os.system('terraform init') 
        #Updating the Terraform configuration
        with open('main.tf', 'a') as tf:
            data  = """
            provider "aws" {
                access_key = "{}"
                secret_key = "{}"
                region     = "us-east-1"
            }

            resource "aws_instance" "example" {
                  ami           = "ami-05df9ea8bb34958cc"
                  instance_type = "t2.micro"
            }""".format(aws_access_key_id, aws_secret_access_key)
            tf.write(data)
            logging.info("Terraform Deployment is successful. Website has been Deployed in EC2 Instance")

    except Exception as e:
        logger.error("Deployment Failed!")
        logger.error(e)
