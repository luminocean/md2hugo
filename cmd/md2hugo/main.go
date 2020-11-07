package main

import (
    "fmt"
    "github.com/luminocean/md2hugo"
    "github.com/urfave/cli/v2"
    "log"
    "os"
)

func main() {
    app := &cli.App{
        Name:  "md2hugo",
        Usage: "converts your markdown documents into files in Hugo format with front matter",
        Authors: []*cli.Author{
			{
				Name:  "luMinO",
				Email: "luminocean@foxmail.com",
			},
		},
		Flags: []cli.Flag{
            &cli.StringFlag{
            	Name: "tag-base",
            	Aliases: []string{"T"},
				Value: md2hugo.TagBase,
            	Destination: &md2hugo.TagBase,
            },
        },
        Action: func(c *cli.Context) error {
            srcDir := c.Args().Get(0)
            if srcDir == "" {
                return fmt.Errorf("missing source directory as the first argument")
            }
            dstDir := c.Args().Get(1)
            if dstDir == "" {
                return fmt.Errorf("missing source directory as the first argument")
            }
            return md2hugo.ConvertAll(srcDir, dstDir)
        },
    }
    err := app.Run(os.Args)
    if err != nil {
        log.Fatal(err)
    }
    return
}
