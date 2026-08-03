#! /bin/bash

ScriptPath=$0
Dir=$(cd $(dirname "$ScriptPath"); pwd)
ProjectNameFile="$Dir/.sis/project_name.txt"
ProjectName=$(tr -d '[:space:]' < "$ProjectNameFile")

ListOnly=0
Verbosity=${XTESTS_VERBOSITY:-${TEST_VERBOSITY:-3}}


# ##########################################################
# command-line handling

while [[ $# -gt 0 ]]; do

  case $1 in
    -l|--list-only)

      ListOnly=1
      ;;
    --verbosity)

      shift
      Verbosity=$1
      ;;
    --help)

      [ -f "$Dir/.sis/script_info_lines.txt" ] && cat "$Dir/.sis/script_info_lines.txt"
      cat << EOF
Runs all (matching) component and unit test packages

$ScriptPath [ ... flags/options ... ]

Flags/options:

    behaviour:

    -l
    --list-only
        lists the target packages but does not execute them

    --verbosity <verbosity>
        specifies an explicit verbosity for the unit-test(s)


    standard flags:

    --help
        displays this help and terminates

EOF

      exit 0
      ;;
    *)

      >&2 echo "$ScriptPath: unrecognised argument '$1'; use --help for usage"

      exit 1
      ;;
  esac

  shift
done


# ##########################################################
# main()

cd "$Dir" || exit 1

# Only packages that contain tests (avoids broken multi-main examples/ dirs)
Packages=$(go list -f '{{if or .TestGoFiles .XTestGoFiles}}{{.ImportPath}}{{end}}' ./...)

if [ -z "$Packages" ]; then

  echo "${ScriptPath}: no test packages found for ${ProjectName}"

  exit 0
fi

if [ $ListOnly -ne 0 ]; then

  echo "Listing all ${ProjectName} test packages"

  printf '%s\n' $Packages

  exit 0
fi

if [ $Verbosity -ge 2 ]; then

  echo "Running all ${ProjectName} unit-test packages"
fi

if [ $Verbosity -ge 2 ]; then

  go test -v $Packages
else

  go test $Packages
fi


# ############################## end of file ############################# #
