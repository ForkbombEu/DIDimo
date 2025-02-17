/// <reference types="node" />
import { Cookie } from 'tough-cookie';
import { EventEmitter } from 'node:events';
import { Phase } from 'phasic';
import { Matcher, CheckResult, CheckResults } from './matcher';
import { LoadTestCheck } from './loadtesting';
import { TestData } from './utils/testdata';
import { CapturesStorage } from './utils/runner';
import { CredentialsStorage } from './utils/auth';
import { HTTPStep, HTTPStepRequest, HTTPStepResponse } from './steps/http';
import { gRPCStep, gRPCStepRequest, gRPCStepResponse } from './steps/grpc';
import { SSEStep, SSEStepRequest, SSEStepResponse } from './steps/sse';
import { PluginStep } from './steps/plugin';
import { tRPCStep } from './steps/trpc';
import { GraphQLStep } from './steps/graphql';
export declare type Workflow = {
    version: string;
    name: string;
    env?: WorkflowEnv;
    /**
     * @deprecated Import files using `$refs` instead.
    */
    include?: string[];
    before?: Test;
    tests: Tests;
    after?: Test;
    components?: WorkflowComponents;
    config?: WorkflowConfig;
};
export declare type WorkflowEnv = {
    [key: string]: string;
};
export declare type WorkflowComponents = {
    schemas?: {
        [key: string]: any;
    };
    credentials?: CredentialsStorage;
};
export declare type WorkflowConfig = {
    loadTest?: {
        phases: Phase[];
        check?: LoadTestCheck;
    };
    continueOnFail?: boolean;
    http?: {
        baseURL?: string;
        rejectUnauthorized?: boolean;
        http2?: boolean;
    };
    grpc?: {
        proto: string | string[];
    };
    concurrency?: number;
};
export declare type WorkflowOptions = {
    path?: string;
    secrets?: WorkflowOptionsSecrets;
    ee?: EventEmitter;
    env?: WorkflowEnv;
    concurrency?: number;
};
declare type WorkflowOptionsSecrets = {
    [key: string]: string;
};
export declare type WorkflowResult = {
    workflow: Workflow;
    result: {
        tests: TestResult[];
        passed: boolean;
        timestamp: Date;
        duration: number;
        bytesSent: number;
        bytesReceived: number;
        co2: number;
    };
    path?: string;
};
export declare type Test = {
    name?: string;
    env?: object;
    steps: Step[];
    testdata?: TestData;
};
export declare type Tests = {
    [key: string]: Test;
};
export declare type Step = {
    id?: string;
    name?: string;
    retries?: {
        count: number;
        interval?: string | number;
    };
    if?: string;
    http?: HTTPStep;
    trpc?: tRPCStep;
    graphql?: GraphQLStep;
    grpc?: gRPCStep;
    sse?: SSEStep;
    delay?: string;
    plugin?: PluginStep;
};
export declare type StepCheckValue = {
    [key: string]: string;
};
export declare type StepCheckJSONPath = {
    [key: string]: any;
};
export declare type StepCheckPerformance = {
    [key: string]: number;
};
export declare type StepCheckCaptures = {
    [key: string]: any;
};
export declare type StepCheckMatcher = {
    [key: string]: Matcher[];
};
export declare type TestResult = {
    id: string;
    name?: string;
    steps: StepResult[];
    passed: boolean;
    timestamp: Date;
    duration: number;
    co2: number;
    bytesSent: number;
    bytesReceived: number;
};
export declare type StepResult = {
    id?: string;
    testId: string;
    name?: string;
    retries?: number;
    captures?: CapturesStorage;
    cookies?: Cookie.Serialized[];
    errored: boolean;
    errorMessage?: string;
    passed: boolean;
    skipped: boolean;
    timestamp: Date;
    responseTime: number;
    duration: number;
    co2: number;
    bytesSent: number;
    bytesReceived: number;
} & StepRunResult;
export declare type StepRunResult = {
    type?: string;
    checks?: StepCheckResult;
    request?: HTTPStepRequest | gRPCStepRequest | SSEStepRequest | any;
    response?: HTTPStepResponse | gRPCStepResponse | SSEStepResponse | any;
};
export declare type StepCheckResult = {
    [key: string]: CheckResult | CheckResults;
};
export declare function runFromYAML(yamlString: string, options?: WorkflowOptions): Promise<WorkflowResult>;
export declare function runFromFile(path: string, options?: WorkflowOptions): Promise<WorkflowResult>;
export declare function run(workflow: Workflow, options?: WorkflowOptions): Promise<WorkflowResult>;
export {};
