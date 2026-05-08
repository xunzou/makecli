/**
 * [INPUT]: 依赖 cmd/client（newClientFromProfile）、fmt、github.com/spf13/cobra
 * [OUTPUT]: 对外提供 newRelationDeleteCmd 函数
 * [POS]: cmd/relation 的 delete 子命令，调用 Meta Server API 删除指定 Relation
 * [PROTOCOL]: 变更时更新此头部，然后检查 CLAUDE.md
 */

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newRelationDeleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "delete <name>",
		Short:        "Delete a relation on Make",
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			app, _ := cmd.Parent().Flags().GetString("app")
			return runRelationDelete(args[0], app)
		},
	}

	return cmd
}

func runRelationDelete(name, app string) error {
	client, err := newClientFromProfile()
	if err != nil {
		return err
	}

	if err := client.DeleteRelation(name, app); err != nil {
		return err
	}

	fmt.Printf("Relation '%s' deleted successfully from app '%s'\n", name, app)
	return nil
}
