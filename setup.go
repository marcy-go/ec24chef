package main

import (
  "flag"
  "os"
  "fmt"
)

type Setup struct {
  configure bool
  host      string
  key_path  string
  validator string
  bucket    string
  object    string
}

func (r *Setup) Help() string {
  return "ec24chef setup"
}

func (r *Setup) Synopsis() string {
  return "Setup Chef-Client Helper for Amazon EC2"
}

func (r *Setup) Run(args []string) int {

  f := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
  def  := ""
  desc := "Chef Server IP or Host"
  f.StringVar(&r.host, "s",       def, desc)
  f.StringVar(&r.host, "-server", def, desc)

  def  = "/etc/chef/validation.pem"
  desc = "Path to the validator key"
  f.StringVar(&r.key_path, "k",    def, desc)
  f.StringVar(&r.key_path, "-key", def, desc)

  def  = "chef-validator"
  desc = "Validation client name"
  f.StringVar(&r.validator, "v",          def, desc)
  f.StringVar(&r.validator, "-validator", def, desc)

  def  = ""
  desc = "S3 bucket name for the validator key"
  f.StringVar(&r.bucket, "b",       def, desc)
  f.StringVar(&r.object, "-bucket", def, desc)

  def  = ""
  desc = "S3 object key for the validator key"
  f.StringVar(&r.bucket, "o",       def, desc)
  f.StringVar(&r.object, "-object", def, desc)

  f.Parse(args)

  if r.host == "" {
    fmt.Fprintln(os.Stderr, "-s or --server is required.")
    return 1
  }

  return 0
}
