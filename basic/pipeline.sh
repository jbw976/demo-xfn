# Try running ./pipeline.sh basic.json

DATE=$(date +"%Y-%m-%dT%H:%M:%S%:z")
cat $1 \
  | jq '.desired.composite.resource.labels |= {"labelizer.xfn.crossplane.io/crossplane": "rocks"} + .' \
  | jq --arg date "$DATE" '.desired.composite.resource.annotations |= {"pipeline.crossplane.io/date": $date} + .'

