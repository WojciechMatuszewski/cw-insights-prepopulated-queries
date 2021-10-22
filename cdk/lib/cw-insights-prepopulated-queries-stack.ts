import * as cdk from "@aws-cdk/core";
import * as customResources from "@aws-cdk/custom-resources";
import * as iam from "@aws-cdk/aws-iam";
import * as apigwv2 from "@aws-cdk/aws-apigatewayv2";
import * as apigwv2Integrations from "@aws-cdk/aws-apigatewayv2-integrations";
import * as lambdaNodejs from "@aws-cdk/aws-lambda-nodejs";
import * as lambda from "@aws-cdk/aws-lambda";
import { join } from "path";

export class CWInsightsPrepopulatedQueriesStack extends cdk.Stack {
  constructor(scope: cdk.Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    const hiHandlerFunction = new lambdaNodejs.NodejsFunction(
      this,
      "hiHandler",
      {
        entry: join(__dirname, "handler.ts"),
        handler: "handler",
        tracing: lambda.Tracing.ACTIVE
      }
    );

    const httpApi = new apigwv2.HttpApi(this, "api");
    httpApi.addRoutes({
      integration: new apigwv2Integrations.LambdaProxyIntegration({
        handler: hiHandlerFunction
      }),
      path: "/",
      methods: [apigwv2.HttpMethod.GET]
    });

    new cdk.CfnOutput(this, "apiEndpoint", {
      value: `${httpApi.apiEndpoint}/`
    });

    new InsightsQuery(this, "byAPIGWRequestId", {
      name: "By APIGW RequestId",
      queryString: BY_APIGW_REQUEST_ID_QUERY
    });

    new InsightsQuery(this, "byXRayTraceId", {
      name: "By X-Ray TraceId",
      queryString: BY_XRAY_TRACE_ID_QUERY
    });

    new InsightsQuery(this, "byLambdaTimeout", {
      name: "Find AWS Lambda timeouts",
      queryString: BY_LAMBDA_TIMEOUT_QUERY
    });

    new InsightsQuery(this, "listLogs", {
      name: "List logs",
      queryString: LIST_LOGS_QUERY
    });
  }
}

interface InsightsQueryProps {
  name: string;
  queryString: string;
}

class InsightsQuery extends cdk.Construct {
  public resource: customResources.AwsCustomResource;
  constructor(scope: cdk.Construct, id: string, props: InsightsQueryProps) {
    super(scope, id);

    this.resource = new customResources.AwsCustomResource(
      this,
      "insightsQuery",
      {
        onCreate: {
          action: "putQueryDefinition",
          service: "CloudWatchLogs",
          parameters: {
            name: props.name,
            queryString: props.queryString
          },
          physicalResourceId:
            customResources.PhysicalResourceId.fromResponse("queryDefinitionId")
        },
        policy: customResources.AwsCustomResourcePolicy.fromStatements([
          new iam.PolicyStatement({
            effect: iam.Effect.ALLOW,
            actions: ["logs:PutQueryDefinition"],
            resources: ["arn:aws:logs:*:*:*"]
          })
        ]),
        onUpdate: {
          action: "putQueryDefinition",
          service: "CloudWatchLogs",
          parameters: {
            name: props.name,
            queryDefinitionId:
              new customResources.PhysicalResourceIdReference(),
            queryString: props.queryString
          },
          physicalResourceId:
            customResources.PhysicalResourceId.fromResponse("queryDefinitionId")
        },
        onDelete: {
          action: "deleteQueryDefinition",
          service: "CloudWatchLogs",
          parameters: {
            queryDefinitionId: new customResources.PhysicalResourceIdReference()
          }
        }
      }
    );
  }
}

const LIST_LOGS_QUERY = `fields @timestamp, @logStream, @message
| sort @timestamp desc`;

const BY_APIGW_REQUEST_ID_QUERY = `fields @timestamp, @logStream, @message
| sort @timestamp desc
| filter @requestId = "PASTE_REQUEST_ID_HERE"`;

const BY_XRAY_TRACE_ID_QUERY = `fields @timestamp, @logStream, @message
| sort @timestamp desc
| filter @xrayTraceId = "PASTE_TRACE_ID_HERE"`;

const BY_LAMBDA_TIMEOUT_QUERY = `fields @timestamp, @logStream, @message
| sort @timestamp desc
| filter @message like /task timed out/`;
