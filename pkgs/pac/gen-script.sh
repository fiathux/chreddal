#!/bin/bash

cat > buildinscript.go <<EOF
// This file is generated by gen-script.sh
// `date +"%Y-%m-%d, %H:%M:%S.%N"`

package pac

var buildinScript = mustLoadScript(
\``cat buildin-script.js |base64`\`,
)
EOF
