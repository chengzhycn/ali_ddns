/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"
	"strings"
	"time"

	"github.com/chengzhycn/ali_ddns/pkg/ali"
	"github.com/spf13/cobra"
	"github.com/vishvananda/netlink"
	"go.uber.org/zap"
)

var (
	accessKeyId     = os.Getenv("ALI_ACCESS_KEY_ID")
	accessKeySecret = os.Getenv("ALI_ACCESS_KEY_SECRET")
	domainName      = os.Getenv("ALI_DOMAIN_NAME")
	devName         = os.Getenv("DEV_NAME")
	v4RR            = os.Getenv("V4_RR")
	v6RR            = os.Getenv("V6_RR")

	checkInterval = 10 * time.Minute
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ali_ddns",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		logger, _ := zap.NewProduction()
		defer logger.Sync()
		suger := logger.Sugar()

		ticker := time.NewTicker(checkInterval)
		defer ticker.Stop()

		for {
			v4Addr, v6Addr, err := GetLocalIPAddr()
			if err != nil {
				suger.Errorf("get local ip addr error: %v", err)
				continue
			}

			client, err := ali.CreateClient(accessKeyId, accessKeySecret)
			if err != nil {
				suger.Errorf("create ali client error: %v", err)
				continue
			}

			if v4RR == "" && v6RR == "" {
				suger.Panicf("no dns record specified")
			}

			records, err := client.DescribeDNSRecord(domainName)
			if err != nil {
				suger.Errorf("describe dns record error: %v", err)
				continue
			}

			for _, record := range records {
				switch record.RR {
				case v4RR:
					suger.Infof("record: %v", record)
					if record.Value != v4Addr {
						prevRecord := record.Value
						record.Value = v4Addr
						if err = client.UpdateDNSRecord(record); err != nil {
							suger.Errorf("update dns record error: %v", err)
							continue
						}

						suger.Infof("update dns record %s from %s to %s", record.RR, prevRecord, v4Addr)
					}
				case v6RR:
					suger.Infof("record: %v", record)
					if record.Value != v6Addr {
						prevRecord := record.Value
						record.Value = v6Addr
						if err = client.UpdateDNSRecord(record); err != nil {
							suger.Errorf("update dns record error: %v", err)
							continue
						}

						suger.Infof("update dns record %s from %s to %s", record.RR, prevRecord, v6Addr)
					}
				}
			}

			<-ticker.C
		}
	},
}

func GetLocalIPAddr() (string, string, error) {
	var v4Addr, v6Addr string

	link, err := netlink.LinkByName(devName)
	if err != nil {
		return "", "", err
	}

	v4AddrList, err := netlink.AddrList(link, netlink.FAMILY_V4)
	if err != nil {
		return "", "", err
	}

	for _, addr := range v4AddrList {
		if addr.IPNet.IP.IsPrivate() {
			v4Addr = addr.IPNet.IP.String()
		}
	}

	v6AddrList, err := netlink.AddrList(link, netlink.FAMILY_V6)
	if err != nil {
		return "", "", err
	}

	for _, addr := range v6AddrList {
		if addr.IPNet.IP.IsGlobalUnicast() && !strings.HasPrefix(addr.IPNet.IP.String(), "fd") {
			v6Addr = addr.IPNet.IP.String()
		}
	}

	return v4Addr, v6Addr, nil
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ali_ddns.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
