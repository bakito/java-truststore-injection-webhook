#!/bin/bash
set -e

echo "## Pod adguardhome-origin logs" >> $GITHUB_STEP_SUMMARY
echo '```' >> $GITHUB_STEP_SUMMARY
kubectl logs -l app.kubernetes.io/instance=java-truststore-injection-webhook >> $GITHUB_STEP_SUMMARY
echo '```' >> $GITHUB_STEP_SUMMARY
