package main

import (
  "os"
  "fmt"
  "github.com/mitchellh/cli"
)

func main() {

  c := cli.NewCLI("ec24chef", "0.0.1")
  c.Args = os.Args[1:]
  c.Commands = map[string]cli.CommandFactory{
    "configure": func() (cli.Command, error) {
      return &Configure{}, nil
    },
    // "run": func() (cli.Command, error) {
    //   return &Run{}, nil
    // },
  }

  ret, err := c.Run()
  if err != nil {
    fmt.Fprintf(os.Stderr, err.Error())
  }
  os.Exit(ret)

}
