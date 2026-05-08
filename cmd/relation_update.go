/**
 * [INPUT]: 依赖 cmd/client（newClientFromProfile）、cmd/relation_create（loadRelationProperties）、fmt、github.com/spf13/cobra
 * [OUTPUT]: 对外提供 newRelationUpdateCmd 函数
 * [POS]: cmd/relation 的 update 子命令，从 JSON 文件加载 from/to 配置，调用 Meta Server API 更新 Relation
 * [PROTOCOL]: 变更时更新此头部，然后检查 CLAUDE.md
 */

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newRelationUpdateCmd() *cobra.Command {
	var jsonFile string

	cmd := &cobra.Command{
		Use:          "update <name>",
		Short:        "Update an existing relation on Make",
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			app, _ := cmd.Parent().Flags().GetString("app")
			return runRelationUpdate(args[0], app, jsonFile)
		},
	}

	cmd.Flags().StringVar(&jsonFile, "json", "", "path to JSON file containing relation properties (required)")
	_ = cmd.MarkFlagRequired("json")
	return cmd
}

func runRelationUpdate(name, app, jsonFile string) error {
	client, err := newClientFromProfile()
	if err != nil {
		return err
	}

	props, err := loadRelationProperties(jsonFile)
	if err != nil {
		return err
	}

	if err := client.UpdateRelation(name, app, props); err != nil {
		return err
	}

	fmt.Printf("Relation '%s' updated successfully in app '%s'\n", name, app)
	return nil
}
