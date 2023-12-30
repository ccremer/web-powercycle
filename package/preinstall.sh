#!/bin/sh

user_name="web-powercycle"

cleanInstall() {
    # Step 3 (clean install), enable the service in the proper way for this platform
    adduser --system --no-create-home "${user_name}"
}

upgrade() {
    # Step 3(upgrade), do what you need
    cleanInstall
}

# Step 2, check if this is a clean install or an upgrade
action="$1"
if  [ "$1" = "configure" ] && [ -z "$2" ]; then
  # Alpine linux does not pass args, and deb passes $1=configure
  action="install"
elif [ "$1" = "configure" ] && [ -n "$2" ]; then
    # deb passes $1=configure $2=<current version>
    action="upgrade"
fi

case "$action" in
  "1" | "install")
    cleanInstall
    ;;
  "2" | "upgrade")
    upgrade
    ;;
  *)
    # Alpine: $1 == version being installed
    cleanInstall
    ;;
esac
