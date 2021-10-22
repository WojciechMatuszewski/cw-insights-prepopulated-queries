#!/usr/bin/env node
import "source-map-support/register";
import * as cdk from "@aws-cdk/core";
import { CWInsightsPrepopulatedQueriesStack } from "../lib/cw-insights-prepopulated-queries-stack";

const app = new cdk.App();
new CWInsightsPrepopulatedQueriesStack(
  app,
  "CWInsightsPrePopulatedQueriesStack"
);
