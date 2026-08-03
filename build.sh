#! /bin/bash

ScriptPath=$0
Dir=$(cd $(dirname "$ScriptPath"); pwd)
ProjectNameFile="$Dir/.sis/project_name.txt"
ProjectName=$(tr -d '[:space:]' < "$ProjectNameFile")


# ##########################################################
# command-line handling

while [[ $# -gt 0 ]]; do

  case $1 in
    --help)

      [ -f "$Dir/.sis/script_info_lines.txt" ] && cat "$Dir/.sis/script_info_lines.txt"
      cat << EOF
Builds the ${ProjectName} library and (if present) examples

$ScriptPath [ ... flags/options ... ]

Flags/options:

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

mkdir -p "$Dir/scratch/build-artefacts/"


# library
echo "building ${ProjectName} library ..."
go build -v "$Dir"


# examples
if [ -d "$Dir/examples" ]; then

  echo "building ${ProjectName} examples ..."

  # subdirectory examples (examples/<name>/main.go)
  for example_dir in "$Dir"/examples/*/; do

    [ -d "$example_dir" ] || continue

    if [ -f "${example_dir}main.go" ]; then

      example_name=$(basename "$example_dir")
      echo "building example ${example_name} ..."
      go build -v -o "$Dir/scratch/build-artefacts/${example_name}" "$example_dir"
    fi
  done

  # loose single-file examples (examples/*.go)
  for example_file in "$Dir"/examples/*.go; do

    [ -f "$example_file" ] || continue

    example_name=$(basename "$example_file" .go)
    echo "building example ${example_name} ..."
    go build -v -o "$Dir/scratch/build-artefacts/${example_name}" "$example_file"
  done
fi


# tests (standalone programs under tests/, if present)
if [ -d "$Dir/tests" ]; then

  echo "building ${ProjectName} tests ..."

  for test_dir in "$Dir"/tests/*/; do

    [ -d "$test_dir" ] || continue

    if [ -f "${test_dir}main.go" ]; then

      test_name=$(basename "$test_dir")
      echo "building test ${test_name} ..."
      go build -v -o "$Dir/scratch/build-artefacts/${test_name}" "$test_dir"
    fi
  done

  for test_file in "$Dir"/tests/*.go; do

    [ -f "$test_file" ] || continue

    test_name=$(basename "$test_file" .go)
    echo "building test ${test_name} ..."
    go build -v -o "$Dir/scratch/build-artefacts/${test_name}" "$test_file"
  done
fi


# ############################## end of file ############################# #
