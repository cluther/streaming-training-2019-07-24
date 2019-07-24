#!/usr/bin/env bash

API_KEY="<ZENOSS_API_KEY>"
SOURCE_TYPE="com.zenoss.example-json.sh"
APP="example-json.sh"
SOURCE=$(hostname)
TIME_NOW_MS=$(date +%s000)

# API helper for bash scripts.
data-receiver-api () {
    SERVICE=$1
    curl -s \
        https://api.zenoss.io/zenoss.cloud.DataReceiverService/$SERVICE \
        -H "content-type: application/json" \
        -H "zenoss-api-key: $API_KEY" \
        -X POST \
        -d @- | jq
}

# Send a model.
echo "Sending model for $APP."
cat << EOF | data-receiver-api PutModels
{
    "models": [
        {
            "timestamp": $TIME_NOW_MS,
            "dimensions": {
                "app": "$APP",
                "source": "$SOURCE"
            },
            "metadataFields": {
                "source-type": "$SOURCE_TYPE",
                "source": "$SOURCE",
                "name": "$APP"
            }
        }
    ]
}
EOF

echo

# Send a metric for the entity in the model above.
echo "Sending random.number metric for $APP."
cat << EOF | data-receiver-api PutMetrics
{
    "metrics": [
        {
            "metric": "random.number",
            "timestamp": $TIME_NOW_MS,
            "value": $RANDOM,
            "dimensions": {
                "app": "$APP"
            },
            "metadataFields": {
                "source-type": "$SOURCE_TYPE",
                "source": "$SOURCE"
            }
        }
    ]
}
EOF
