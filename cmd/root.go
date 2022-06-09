package cmd

import (
	"fmt"
	"gscript/complier"
	"gscript/proto"
	"gscript/std"
	"gscript/vm"
	"os"
	"path"
	"strings"

	"github.com/spf13/cobra"
)

var (
	Flag_Output string
	Flag_Asm    bool
)

var rootCmd = &cobra.Command{
	Use:  "gsc",
	Args: cobra.MinimumNArgs(1),
}

var versionCmd = &cobra.Command{
	Use:                   "version",
	Short:                 "Show Version",
	Args:                  cobra.NoArgs,
	DisableFlagParsing:    true,
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(Version)
	},
}

var runCmd = &cobra.Command{
	Use:                   "run <file>",
	Short:                 "Excute script file. Usage: gsc run <file>",
	Args:                  cobra.ExactArgs(1),
	DisableFlagParsing:    true,
	DisableFlagsInUseLine: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		vm, err := initVM(args[0])
		if err != nil {
			return err
		}
		vm.Run()
		return nil
	},
}

var debugCmd = &cobra.Command{
	Use:                   "debug <file>",
	Short:                 "Debug script file. Usage: gsc debug <file>",
	Args:                  cobra.ExactArgs(1),
	DisableFlagsInUseLine: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		vm, err := initVM(args[0])
		if err != nil {
			return err
		}
		vm.Debug()
		return nil
	},
}

var buildCmd = &cobra.Command{
	Use:                   "build [flags] <file>",
	Short:                 "complie srcipt file. Usage: gsc build [flags] <file>",
	Args:                  cobra.ExactArgs(1),
	DisableFlagsInUseLine: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO handle Flag_Asm
		src := args[0]
		if Flag_Output == "" {
			Flag_Output = strings.TrimSuffix(src, path.Ext(path.Base(src))) + std.ProtoSuffix
		}
		protos, err := complier.ComplieWithSrcFile(src)
		if err != nil {
			return err
		}
		return proto.WriteProtosToFile(Flag_Output, protos)
	},
}

func initVM(src string) (*vm.VM, error) {
	var protos []proto.Proto
	var err error
	if proto.IsProtoFile(src) {
		_, protos, err = proto.ReadProtosFromFile(src)
	} else {
		protos, err = complier.ComplieWithSrcFile(src)
	}
	if err != nil {
		return nil, err
	}

	stdlibs, err := std.ReadProtos()
	if err != nil {
		return nil, err
	}
	v := vm.NewVM(protos, stdlibs)
	return v, nil
}

func Execute() {
	buildCmd.Flags().StringVarP(&Flag_Output, "output", "o", "", "output file")
	buildCmd.Flags().BoolVarP(&Flag_Asm, "asm", "a", false, "output human-readable assembly code")
	rootCmd.AddCommand(runCmd, versionCmd, debugCmd, buildCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
