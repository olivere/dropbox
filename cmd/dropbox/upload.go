package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox/files"
	"github.com/pkg/errors"

	"github.com/olivere/dropbox/x/env"
)

const (
	chunkSize int64 = 1 << 24
)

// uploadCommand uploads a file to a Dropbox folder.
type uploadCommand struct {
	appKey    string
	appSecret string
	domain    string
	verbose   bool
	input     string
	outpath   string
}

func init() {
	RegisterCommand("upload", func(flags *flag.FlagSet) Command {
		cmd := new(uploadCommand)
		flags.StringVar(&cmd.appKey, "key", env.String("", "DROPBOX_KEY"), "Dropbox API Key (DROPBOX_KEY)")
		flags.StringVar(&cmd.appSecret, "secret", env.String("", "DROPBOX_SECRET"), "Dropbox API Secret (DROPBOX_SECRET)")
		flags.StringVar(&cmd.domain, "domain", env.String("", "DROPBOX_DOMAIN"), "Dropbox API Domain (DROPBOX_DOMAIN)")
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

	path := pathify(args[len(args)-1])

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
	for _, input := range args[:len(args)-1] {
		r, err := os.Open(input)
		if err != nil {
			return errors.Wrapf(err, "unable to open file %s", input)
		}
		fi, err := r.Stat()
		if err != nil {
			return errors.Wrapf(err, "unable to stat file %s", input)
		}

		arg := files.NewCommitInfo(path)
		arg.Mode.Tag = "overwrite"
		arg.ClientModified = time.Now().UTC().Round(time.Second)
		if fi.Size() > chunkSize {
			if err = cmd.chunkedUpload(dbx, r, arg, fi.Size()); err != nil {
				return errors.Wrapf(err, "unable to chunk upload %s", input)
			}
		} else {
			if _, err = dbx.Upload(arg, r); err != nil {
				return errors.Wrapf(err, "unable to upload %s", input)
			}
		}
	}

	/*
		dbx := files.New(config)
		for _, input := range args[:len(args)-1] {
			r, err := os.Open(input)
			if err != nil {
				return errors.Wrapf(err, "unable to open file %s", input)
			}
			arg := files.NewCommitInfoWithProperties(path)
			if _, err := dbx.AlphaUpload(arg, r); err != nil {
				_ = r.Close()
				return errors.Wrapf(err, "unable to copy %s", input)
			}
			_ = r.Close()
		}
	*/

	return nil
}

func (cmd *uploadCommand) chunkedUpload(dbx files.Client, r io.Reader, commitInfo *files.CommitInfo, sizeTotal int64) error {
	res, err := dbx.UploadSessionStart(
		files.NewUploadSessionStartArg(),
		&io.LimitedReader{R: r, N: chunkSize},
	)
	if err != nil {
		return err
	}

	written := chunkSize

	for (sizeTotal - written) > chunkSize {
		cursor := files.NewUploadSessionCursor(res.SessionId, uint64(written))
		args := files.NewUploadSessionAppendArg(cursor)

		err = dbx.UploadSessionAppendV2(args, &io.LimitedReader{R: r, N: chunkSize})
		if err != nil {
			return err
		}
		written += chunkSize
	}

	cursor := files.NewUploadSessionCursor(res.SessionId, uint64(written))
	args := files.NewUploadSessionFinishArg(cursor, commitInfo)

	if _, err = dbx.UploadSessionFinish(args, r); err != nil {
		return err
	}
	return nil
}
