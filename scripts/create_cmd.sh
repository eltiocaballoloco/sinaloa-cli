#!/bin/bash

# Usage: ./create_cmd.sh cmd_name subcmd_name package_path "flag1|flag2|..."

set -e  # Exit on error

CMD=$1
SUBCMD=$2
PACKAGE_PATH=$3
FLAGS=$4

CMD_CAP="$(tr '[:lower:]' '[:upper:]' <<< ${CMD:0:1})${CMD:1}"
SUBCMD_CAP="$(tr '[:lower:]' '[:upper:]' <<< ${SUBCMD:0:1})${SUBCMD:1}"

CMD_DIR=./cmd/$CMD
SUB_DIR=$CMD_DIR/sub
CMD_FILE=$CMD_DIR/$CMD.go
SUB_FILE=$SUB_DIR/$SUBCMD.go
ROOT_FILE=./cmd/root.go

mkdir -p $CMD_DIR $SUB_DIR

# Create the main command file
cat <<EOF > $CMD_FILE
package $CMD

import (
    "github.com/spf13/cobra"
    "$PACKAGE_PATH/sub"
)

var ${CMD_CAP}Cmd = &cobra.Command{
    Use:   "$CMD",
    Short: "Short description of $CMD",
    Long:  "Long description of $CMD",
    Run: func(cmd *cobra.Command, args []string) {
        cmd.Help()
    },
}

func init() {
    ${CMD_CAP}Cmd.AddCommand(sub.${SUBCMD_CAP}${CMD_CAP}Cmd)
}
EOF

# Create the subcommand file
cat <<EOF > $SUB_FILE
package sub

import (
    "fmt"
    "github.com/spf13/cobra"
)

var (
EOF

IFS='|' read -ra ADDR <<< "$FLAGS"
for i in "${ADDR[@]}"; do
    IFS=':' read -r varname flagname shortname desc required default <<< "$i"
    echo "    $varname string" >> $SUB_FILE
done

cat <<EOF >> $SUB_FILE
)

var ${SUBCMD_CAP}${CMD_CAP}Cmd = &cobra.Command{
    Use:   "$SUBCMD",
    Short: "Short description of $SUBCMD",
    Long:  "Long description of $SUBCMD",
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("Executing $SUBCMD in $CMD")
    },
}

func init() {
EOF

for i in "${ADDR[@]}"; do
    IFS=':' read -r varname flagname shortname desc required default <<< "$i"
    echo "    ${SUBCMD_CAP}${CMD_CAP}Cmd.Flags().StringVarP(&$varname, \"$flagname\", \"$shortname\", \"$default\", \"$desc\")" >> $SUB_FILE
    if [ "$required" = "true" ]; then
        echo "    ${SUBCMD_CAP}${CMD_CAP}Cmd.MarkFlagRequired(\"$flagname\")" >> $SUB_FILE
    fi
done

cat <<EOF >> $SUB_FILE
}

EOF

# Update root.go to import the new command and add it to the root command
sed -i '' "/\"github.com\/spf13\/cobra\"/a \\
    \"$PACKAGE_PATH\"
" $ROOT_FILE

sed -i '' "/func addSubcommandPalettes() {/a \\
    rootCmd.AddCommand($CMD.${CMD_CAP}Cmd)
" $ROOT_FILE
