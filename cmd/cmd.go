package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/violin0622/ink-scroll/cmd/contents"
)

func New() *cobra.Command {
	root := &cobra.Command{
		Use:     `github.com/violin0622/ink-scroll`,
		Short:   `Play DDD Example Program.`,
		Long:    slogan,
		Version: fmt.Sprintf(versionTpl, Version, GitSHA, GoVersion, BuiltAt),
	}

	root.AddCommand(contents.Contents())
	return root
}

const (
	slogan = `
     ____  __               ____  ____  ____
    / __ \/ /___ ___  __   / __ \/ __ \/ __ \
   / /_/ / / __ '/ / / /  / / / / / / / / / /
  / ____/ / /_/ / /_/ /  / /_/ / /_/ / /_/ /
 /_/   /_/\__,_/\__, /  /_____/_____/_____/
               /____/

Play DDD Example Program.
`

	versionTpl = `

Version: %s
Git SHA: %s
Go Version: %s
Built At: %s
`
)

var Version, GitSHA, BuiltAt, GoVersion string
