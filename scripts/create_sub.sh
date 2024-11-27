#!/bin/bash

# Usage: ./create_sub.sh cmd_name subcmd_name package "flag1:flag2:..."
CMD=$1
SUBCMD=$2
PACKAGE_REF=$3
FLAGS=$4

# Create the sub-command
mkdir -p "$CMD/sub"
SUBCMD_FILE="$CMD/sub/$SUBCMD.go"
touch "$SUBCMD_FILE"

# Construct the sub-command file
cat <<EOF > "$SUBCMD_FILE"
package sub

import (
    "fmt"
    "github.com/spf13/cobra"
)

var (
EOF

IFS=':' read -ra FLAG_ARRAY <<< "$FLAGS"
for FLAG in "${FLAG_ARRAY[@]}"; do
    IFS='|' read -ra FLAG_PARTS <<< "$FLAG"
    VAR_NAME="${FLAG_PARTS[0]}"
    FLAG_NAME="${FLAG_PARTS[1]}"
    SHORT_NAME="${FLAG_PARTS[2]}"
    DESC="${FLAG_PARTS[3]}"
    REQUIRED="${FLAG_PARTS[4]}"

    # Create a string variable for each flag
    echo "    $VAR_NAME string" >> "$SUBCMD_FILE"
    
    # Create the flag
    echo "    ${SUBCMD^}.Flags().StringVarP(&$VAR_NAME, \"$FLAG_NAME\", \"$SHORT_NAME\", \"\", \"$DESC\")" >> "$SUBCMD_FILE"
    
    # Mark the flag as required if needed
    if [ "$REQUIRED" = "true" ]; then
        echo "    ${SUBCMD^}.Flags().MarkFlagRequired(\"$FLAG_NAME\")" >> "$SUBCMD_FILE"
    fi
done

# Complete the sub-command file
cat <<EOF >> "$SUBCMD_FILE"
)

var ${SUBCMD^}Cmd = &cobra.Command{
    Use:   "$SUBCMD",
    Short: "Short description of $SUBCMD",
    Long:  "Long description of $SUBCMD",
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("Executing $SUBCMD in $CMD")
    },
}

func init() {
    ${CMD^}Cmd.AddCommand(${SUBCMD^}Cmd)
}
EOF

# Print a message indicating success
echo "Created sub-command $SUBCMD for command $CMD in package $PACKAGE_REF"
