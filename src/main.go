package main

import (
	"./utils"
	"flag"
	"os"
)

func init(){
	flag.Parse()
}

func main(){
	cli := utils.NewEc2Client(*utils.ArgRegion)
	if *utils.ArgShowTags {
		utils.Unwrap("Failed to describe tag.", utils.DescribeTag(cli, *utils.ArgShowTagsFilter))
	} else {
		if *utils.ArgShowTagsFilter != "" {
			utils.ExitErrorf("Require to set -showtags flag.")
		} else if *utils.ArgTags == "" {
			utils.ExitErrorf("At least one tag setting is required.")
			os.Exit(1)
		}
		tags, err := utils.GenerateTags(*utils.ArgTags)
		utils.Unwrap("Failed to generate tags.", err)

		if *utils.ArgAdd {
			utils.Unwrap("Failed to create tags.", utils.CreateTag(cli, *utils.ArgInstances, tags))
		} else if *utils.ArgDel {
			utils.Unwrap("Failed to delete tags.", utils.DeleteTag(cli, *utils.ArgInstances, tags))
		} else {
			utils.ExitErrorf("You must specify at least `-add` or` -del`.")
		}
	}
}
