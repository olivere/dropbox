package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox/files"
	"github.com/pkg/errors"

	"github.com/olivere/dropbox/x/env"
)

// uploadCommand uploads a file to a Dropbox folder.
type uploadCommand struct {
	appKey    string
	appSecret string
	verbose   bool
	input     string
	outpath   string
}

func init() {
	RegisterCommand("upload", func(flags *flag.FlagSet) Command {
		cmd := new(uploadCommand)
		flags.StringVar(&cmd.appKey, "key", env.String("", "DROPBOX_KEY"), "Dropbox API Key (DROPBOX_KEY)")
		flags.StringVar(&cmd.appSecret, "secret", env.String("", "DROPBOX_SECRET"), "Dropbox API Secret (DROPBOX_SECRET)")
		flags.BoolVar(&cmd.verbose, "v", false, "Verbose output")
		return cmd
	})
}

func (cmd *uploadCommand) Describe() string {
	return "Upload one or more files to a Dropbox folder"
}

func (cmd *uploadCommand) Usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s upload <in1> <in2> <folder>\n", os.Args[0])
}

func (cmd *uploadCommand) Run(args []string) error {
	if cmd.appKey == "" {
		return errors.New("Dropbox API key is missing")
	}
	if cmd.appSecret == "" {
		return errors.New("Dropbox API secret is missing")
	}

	if len(args) < 2 {
		return errors.New("at least one input and one output file/folder is required")
	}

	_, token, err := createClient(cmd.appKey, cmd.appSecret)
	if err != nil {
		return err
	}

	config := dropbox.Config{
		Token: token.AccessToken,
		// AsMemberID: "",
		// Domain: "",
		Verbose: cmd.verbose,
	}

	dbx := files.New(config)
	for _, input := range args[:len(args)-1] {
		r, err := os.Open(input)
		if err != nil {
			return errors.Wrapf(err, "unable to open file %s", input)
		}
		arg := files.NewCommitInfoWithProperties(args[len(args)-1])
		if _, err := dbx.AlphaUpload(arg, r); err != nil {
			r.Close()
			return errors.Wrap(err, "unable to copy")
		}
		r.Close()
	}

	return nil
}
