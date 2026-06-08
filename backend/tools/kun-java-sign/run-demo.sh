#!/usr/bin/env bash
# Run KUN sign demo (no Maven 3.6+ required)
set -euo pipefail
cd "$(dirname "$0")"

if [ -z "${JAVA_HOME:-}" ] && [ -x /usr/libexec/java_home ]; then
  export JAVA_HOME="$(/usr/libexec/java_home 2>/dev/null || true)"
fi

mvn -q dependency:copy-dependencies -DoutputDirectory=target/lib
mvn -q compile

CP="target/classes:target/lib/*"
exec java -cp "$CP" com.motewallet.kun.KunSignDemo "$@"
