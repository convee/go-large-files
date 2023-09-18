package cmd

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/json-iterator/go/extra"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/time/rate"
)

var globalCtx = context.Background()

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "large-files",
	Short: "",
	Long:  "",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rand.Seed(time.Now().UnixNano())
	extra.RegisterFuzzyDecoders()
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&root.cfgFile, "config", "config.yaml", "config.yaml")
	rootCmd.PersistentFlags().Int64Var(&root.minutesCnt, "qps", 10000000, "")
	rootCmd.PersistentFlags().IntVarP(&root.threads, "concurrent", "c", 100, "")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if root.cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(root.cfgFile)
	}
	viper.AutomaticEnv() // read in environment variables that match
	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		fmt.Println("read config file failed:", root.cfgFile, err)
		os.Exit(0)
	}
}

var root struct {
	limiter    *rate.Limiter
	minutesCnt int64
	cfgFile    string
	reportOnce sync.Once
	threads    int
	test       bool
}

// CheckAndReport qps 限流
func CheckAndReport(args ...interface{}) error {
	root.reportOnce.Do(func() {
		qps := root.minutesCnt
		root.minutesCnt = 0
		if qps < 1 {
			qps = 1e6
		}
		root.limiter = rate.NewLimiter(rate.Limit(qps), int(qps))
		go func() {
			for range time.Tick(time.Minute) {
				cnt := atomic.SwapInt64(&root.minutesCnt, 0)
				if cnt != 0 {
					fmt.Println(time.Now().String(), "QPS:", cnt/60)
				}
			}
		}()
	})
	err := root.limiter.Wait(globalCtx)
	if err != nil {
		return err
	}
	if atomic.AddInt64(&root.minutesCnt, 1) == 1 {
		fmt.Println(args...)
	}
	return nil
}
