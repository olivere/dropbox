package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox/files"
	"github.com/pkg/errors"

	"github.com/olivere/dropbox/x/env"
)

// lsCommand lists the contents of a Dropbox folder.
type lsCommand struct {
	appKey    string
	appSecret string
	domain    string
	verbose   bool
	path      string
	recursive bool
}

func init() {
	RegisterCommand("ls", func(flags *flag.FlagSet) Command {
		cmd := new(lsCommand)
		flags.StringVar(&cmd.appKey, "key", env.String("", "DROPBOX_KEY"), "Dropbox API Key (DROPBOX_KEY)")
		flags.StringVar(&cmd.appSecret, "secret", env.String("", "DROPBOX_SECRET"), "Dropbox API Secret (DROPBOX_SECRET)")
		flags.StringVar(&cmd.domain, "domain", env.String("", "DROPBOX_DOMAIN"), "Dropbox API Domain (DROPBOX_DOMAIN)")
		flags.BoolVar(&cmd.recursive, "r", false, "List recursively")
		flags.BoolVar(&cmd.verbose, "v", false, "Verbose output")
		return cmd
	})
}

func (cmd *lsCommand) Describe() string {
	return "List files of a Dropbox folder"
}

func (cmd *lsCommand) Usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s ls [-r] <folder>\n", os.Args[0])
}

func (cmd *lsCommand) Run(args []string) error {
	if cmd.appKey == "" {
		return errors.New("Dropbox API key is missing")
	}
	if cmd.appSecret == "" {
		return errors.New("Dropbox API secret is missing")
	}

	if len(args) < 1 {
		return errors.New("please specify a folder to list")
	}
	cmd.path = pathify(args[0])

	_, token, err := createClient(cmd.appKey, cmd.appSecret, cmd.domain)
	if err != nil {
		return err
	}
	if token == nil {
		return errors.Wrap(err, "unable to get token when creating client")
	}

	config := dropbox.Config{
		Token: token.AccessToken,
		// AsMemberID: "",
		// Domain: "",
		// LogLevel: dropbox.LogOff,
	}
	if cmd.verbose {
		config.LogLevel = dropbox.LogInfo
	}

	dbx := files.New(config)
	listArg := files.NewListFolderArg(cmd.path)
	listArg.Recursive = cmd.recursive
	listRes, err := dbx.ListFolder(listArg)
	if err != nil {
		return errors.Wrapf(err, "unable to list folder at %s", cmd.path)
	}
	cmd.printFiles(listRes.Entries)
	cursor := listRes.Cursor
	hasMore := listRes.HasMore
	for hasMore {
		arg := files.NewListFolderContinueArg(cursor)
		res, err := dbx.ListFolderContinue(arg)
		if err != nil {
			return errors.Wrapf(err, "unable to continue listing folder at %s", cmd.path)
		}
		cmd.printFiles(res.Entries)
		hasMore = res.HasMore
		cursor = res.Cursor
	}
	return nil
}

func (cmd *lsCommand) printFiles(entries []files.IsMetadata) {
	for _, entry := range entries {
		switch f := entry.(type) {
		case *files.FileMetadata:
			cmd.printFile(f)
		case *files.FolderMetadata:
			cmd.printFolder(f)
		}
	}
}

func (cmd *lsCommand) printFile(f *files.FileMetadata) {
	fmt.Printf("-rw-r--r--\t%d\t%s\t%s\n",
		f.Size,
		f.ServerModified.Format(time.RubyDate),
		f.PathDisplay)
}

func (cmd *lsCommand) printFolder(f *files.FolderMetadata) {
	fmt.Printf("drw-r--r--\t%d\t%s\t%s\n",
		0,
		"",
		f.PathDisplay)
}
