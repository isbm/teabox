# Common module functions for a module, written in Bash


# Snippet of loading this library
# ===============================
# Copypaste it to your bash script at the top of it and adjust the amount of
# back-path elements (../) depending where your script is located in the module tree:
#
#     SELF_PATH="$( cd -- "$(dirname "$0")" >/dev/null 2>&1 ; pwd -P )"
#     source $SELF_PATH/../../../lib/shell/bash-common.sh
#
# This should load this library and all its functions will be available in your script.


# Snippet of parsing command line
# ===============================
# Optionally, copypaste this at the end of your Bash script to parse command line
# and adjust it to your needs:
#
#     function run() {
#         echo "Your stuff here"
#     }
#
#     if [[ $# -eq 0 ]]; then
#         run $@
#     else
#         while [[ $# -gt 0 ]]; do
#	      case $1 in
#	          -s|--setup)
#             # setup_function
#		      shift
#		      ;;
#
#	          -h|--help)
#             # help_display_function
#		      shift
#		      ;;
#
#	          *)
#		      run $@
#		      break
#		      ;;
#
#	      esac
#         done
#     fi
#
# Don't forget to uncomment it. :-P


# Standard Socket Location
# ------------------------
# You can surely redefine it here for everything or overwrite it in your module.
SOCK="/tmp/teabox.sock"


# Call socket with an API call
# ----------------------------
# This requires OpenBSD's netcat installed and availble as "nc".
# It takes one required parameter and two optional:
#
#     api <API> [VALUE] [TYPE]
#
# Example usage:
#
#     api logger-status "Hello world"
#     api init-alloc-progress 3 int
#     api field-set-by-ord "{3}false" bool
#
# For more info about API, refer to the documentation.
#
function api() {
    cls=$1
    msg=$2
    typ=$3
    $(echo "$cls:$typ:$msg" | nc -w0 -U $SOCK)
    if [[ "$?" == "1" ]]; then
        exit 1
    fi
}

