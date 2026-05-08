/**
 * [INPUT]: 依赖 cmd/client（newClientFromProfile）、internal/api（RelationProperties/RelationEnd）、encoding/json、fmt、os、github.com/spf13/cobra
 * [OUTPUT]: 对外提供 newRelationCreateCmd 函数
 * [POS]: cmd/relation 的 create 子命令，从 JSON 文件加载 from/to 配置，调用 Meta Server API 创建 Relation
 * [PROTOCOL]: 变更时更新此头部，然后检查 CLAUDE.md
 */

package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/qfeius/makecli/internal/api"
	"github.com/spf13/cobra"
)

func newRelationCreateCmd() *cobra.Command {
	var jsonFile string

	cmd := &cobra.Command{
		Use:          "create <name>",
		Short:        "Create a new relation on Make",
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			app, _ := cmd.Parent().Flags().GetString("app")
			return runRelationCreate(args[0], app, jsonFile)
		},
	}

	cmd.Flags().StringVar(&jsonFile, "json", "", "path to JSON file containing relation properties (required)")
	_ = cmd.MarkFlagRequired("json")
	return cmd
}

func runRelationCreate(name, app, jsonFile string) error {
	client, err := newClientFromProfile()
	if err != nil {
		return err
	}

	props, err := loadRelationProperties(jsonFile)
	if err != nil {
		return err
	}

	if err := client.CreateRelation(name, app, props); err != nil {
		return err
	}

	fmt.Printf("Relation '%s' created successfully in app '%s'\n", name, app)
	return nil
}

// loadRelationProperties 读取 JSON 文件并解析为 RelationProperties
func loadRelationProperties(path string) (api.RelationProperties, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return api.RelationProperties{}, fmt.Errorf("读取 JSON 文件失败: %w", err)
	}
	var props api.RelationProperties
	if err := json.Unmarshal(data, &props); err != nil {
		return api.RelationProperties{}, fmt.Errorf("JSON 文件格式错误（需包含 from/to 对象）: %w", err)
	}
	return props, nil
}
