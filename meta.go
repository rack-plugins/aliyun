package aliyun

import (
	"github.com/fimreal/rack/module"
	"github.com/spf13/cobra"
)

const (
	ID            = "aliyun"
	Comment       = "aliyun api"
	RoutePrefix   = "/" + ID
	DefaultEnable = false
)

var Module = module.Module{
	ID:      ID,
	Comment: Comment,
	// gin route
	RouteFunc:   AddRoute,
	RoutePrefix: RoutePrefix,
	// cobra flag
	FlagFunc: ServeFlag,
}

func ServeFlag(serveCmd *cobra.Command) {
	serveCmd.Flags().Bool(ID, DefaultEnable, Comment)

	serveCmd.Flags().String("aliyun.akid", "", "ACCESS_KEY_ID")
	serveCmd.Flags().String("aliyun.aksecret", "", "ACCESS_KEY_SECRET")
	serveCmd.Flags().String("aliyun.regionid", "", "REGION_ID")
	serveCmd.Flags().Bool("aliyun.insecureskipverify", false, "是否跳过证书验证(小容器没有证书会遇到 https 连接证书验证失败)")
}
