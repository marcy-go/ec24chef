package main

import (
  "flag"
  "os"
  "fmt"
  "io"
  "github.com/awslabs/aws-sdk-go/aws"
  // "github.com/awslabs/aws-sdk-go/gen/ec2"
  "github.com/awslabs/aws-sdk-go/gen/s3"
  "github.com/marcy-go/ec2meta"
)

type Configure struct {
  server    string
  key_path  string
  validator string
  bucket    string
  object    string
  node      string
  role      string
  env       string
}

func (c *Configure) Help() string {
  return "ec24chef setup"
}

func (c *Configure) Synopsis() string {
  return "Configure Chef-Client Helper for Amazon EC2"
}

func (c *Configure) Run(args []string) int {

  f := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
  def  := ""
  desc := "Chef Server IP or Host"
  f.StringVar(&c.server, "s",       def, desc)
  f.StringVar(&c.server, "-server", def, desc)

  def  = "/etc/chef/validation.pem"
  desc = "Path to the validator key"
  f.StringVar(&c.key_path, "k",    def, desc)
  f.StringVar(&c.key_path, "-key", def, desc)

  def  = "chef-validator"
  desc = "Validation client name"
  f.StringVar(&c.validator, "v",          def, desc)
  f.StringVar(&c.validator, "-validator", def, desc)

  def  = ""
  desc = "S3 bucket name for the validator key"
  f.StringVar(&c.bucket, "b",       def, desc)
  f.StringVar(&c.bucket, "-bucket", def, desc)

  def  = ""
  desc = "S3 object key for the validator key"
  f.StringVar(&c.object, "o",       def, desc)
  f.StringVar(&c.object, "-object", def, desc)

  f.Parse(args)

  if c.server == "" {
    fmt.Fprintln(os.Stderr, "-s or --server is required.")
    return 1
  }

  if c.bucket != "" && c.object != "" {
    src, err := c.getKeyObject()
    if err != nil {
      fmt.Fprintln(os.Stderr, err.Error())
      return 1
    }
    err = c.writeKey(src)
    if err != nil {
      fmt.Fprintln(os.Stderr, err.Error())
      return 1
    }
  }

  id, err := ec2meta.GetInstanceId()
  if err != nil {
    fmt.Fprintln(os.Stderr, err.Error())
    return 1
  }
  c.node = id



  return 0
}

func (c *Configure) getKeyObject() (io.Reader, error) {
  var rc io.ReadCloser
  region, err := ec2meta.GetRegion()
  if err != nil {
    return rc, err
  }

  s3_cli := s3.New(aws.IAMCreds(), region, nil)
  req := s3.GetObjectRequest{}
  req.Bucket = aws.String(c.bucket)
  req.Key    = aws.String(c.object)
  res, err := s3_cli.GetObject(&req)
  if err != nil {
    return rc, err
  }
  return res.Body, nil
}

func (c *Configure) writeKey(src io.Reader) error {
  dst, err := os.Create(c.key_path)
  if err != nil {
    return err
  }
  _, err = io.Copy(dst, src)
  if err != nil {
    return err
  }
  return nil
}
