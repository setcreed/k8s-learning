package utils

import (
	"io/ioutil"
	"os"

	"github.com/olekukonko/tablewriter"
	"gopkg.in/yaml.v3"
)

// 从文件中读取内容  并反序列为 struct
func YamlFile2Struct(path string, obj interface{}) error {
	b, err := GetFileContent(path)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(b, obj)
	if err != nil {
		return err
	}
	return nil
}

// 单独封装的 文件读取函数
func GetFileContent(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// item is in []string{}
func InArray(arr []string, item string) bool {
	for _, p := range arr {
		if p == item {
			return true
		}
	}
	return false
}

// 设置table的样式
func SetTable(table *tablewriter.Table) {
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("\t") // pad with tabs
	table.SetNoWhiteSpace(true)
}
