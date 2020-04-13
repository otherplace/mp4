package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	cli "github.com/jawher/mow.cli"
	"github.com/otherplace/mp4"
	"github.com/otherplace/mp4/filter"
)

func main() {
	cmd := cli.App("mp4tool", "MP4 command line tool")

	cmd.Command("info", "Displays information about a media", func(cmd *cli.Cmd) {
		file := cmd.StringArg("FILE", "", "the file to display")
		isJson := cmd.BoolArg("JSON", false, "display as JSON")
		cmd.Action = func() {
			fd, err := os.Open(*file)
			defer fd.Close()
			v, err := mp4.Decode(fd)
			if err != nil {
				panic(err)
			}
			if *isJson {
				jsonBytes, err := json.Marshal(v)
				if err != nil {
					panic(err)
				}
				fmt.Println(string(jsonBytes))
			} else {

				v.Dump()
			}
		}
	})

	cmd.Command("clip", "Generates a clip", func(cmd *cli.Cmd) {
		start := cmd.IntOpt("s start", 0, "start time (sec)")
		duration := cmd.IntOpt("d duration", 10, "duration (sec)")
		src := cmd.StringArg("SRC", "", "the source file name")
		dst := cmd.StringArg("DST", "", "the destination file name")
		cmd.Action = func() {
			in, err := os.Open(*src)
			if err != nil {
				fmt.Println(err)
			}
			defer in.Close()
			v, err := mp4.Decode(in)
			if err != nil {
				fmt.Println(err)
			}
			out, err := os.Create(*dst)
			if err != nil {
				fmt.Println(err)
			}
			defer out.Close()
			filter.EncodeFiltered(out, v, filter.Clip(time.Duration(*start)*time.Second, time.Duration(*duration)*time.Second))
		}
	})

	cmd.Command("copy", "Decodes a media and reencodes it to another file", func(cmd *cli.Cmd) {
		src := cmd.StringArg("SRC", "", "the source file name")
		dst := cmd.StringArg("DST", "", "the destination file name")
		cmd.Action = func() {
			in, err := os.Open(*src)
			if err != nil {
				fmt.Println(err)
			}
			defer in.Close()
			v, err := mp4.Decode(in)
			if err != nil {
				fmt.Println(err)
			}
			out, err := os.Create(*dst)
			if err != nil {
				fmt.Println(err)
			}
			defer out.Close()
			v.Encode(out)
		}
	})
	cmd.Run(os.Args)
}
