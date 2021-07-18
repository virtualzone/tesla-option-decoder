#!/bin/bash
rm -f optioncodes.md optioncodes.json
echo "Downloading option code markdown..."
curl -s "https://raw.githubusercontent.com/timdorr/tesla-api/master/docs/vehicle/optioncodes.md" > optioncodes.md
echo "Converting markdown to JSON..."
line=0
echo "{" > optioncodes.json
cat optioncodes.md | grep "|" | while read -r a; do
    line=$((line+1))
    if [[ $line -gt 2 ]]; then
        IFS='|' read -r -a array <<< "$a";
        declare -a row
        for i in "${!array[@]}"; do
            #e=`echo "${array[i]}" | xargs -0`;
            e="$(echo "${array[i]}" | sed -e 's/^[[:space:]]*//' -e 's/[[:space:]]*$//' -e 's/"/\\"/')"
            row[$i]=$e
        done
        echo "\"${row[1]}\": {\"title\": \"${row[2]}\", \"description\": \"${row[3]}\"}," >> optioncodes.json
        unset row
    fi
done
echo '"__": {}' >> optioncodes.json
echo "}" >> optioncodes.json
rm -f optioncodes.md
echo "Done."