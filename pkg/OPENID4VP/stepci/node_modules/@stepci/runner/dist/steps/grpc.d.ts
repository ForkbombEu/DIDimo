import Ajv from 'ajv';
import { gRPCRequestMetadata } from 'cool-grpc';
import { StepCheckCaptures, StepCheckJSONPath, StepCheckMatcher, StepCheckPerformance } from '..';
import { CapturesStorage } from './../utils/runner';
import { Credential } from './../utils/auth';
import { StepRunResult, WorkflowConfig, WorkflowOptions } from '..';
import { Matcher } from '../matcher';
export declare type gRPCStep = {
    proto: string | string[];
    host: string;
    service: string;
    method: string;
    data?: object | object[];
    timeout?: string | number;
    metadata?: gRPCRequestMetadata;
    auth?: gRPCStepAuth;
    captures?: gRPCStepCaptures;
    check?: gRPCStepCheck;
};
export declare type gRPCStepAuth = {
    tls?: Credential['tls'];
};
export declare type gRPCStepCaptures = {
    [key: string]: gRPCStepCapture;
};
export declare type gRPCStepCapture = {
    jsonpath?: string;
};
export declare type gRPCStepCheck = {
    json?: object;
    schema?: object;
    jsonpath?: StepCheckJSONPath | StepCheckMatcher;
    captures?: StepCheckCaptures;
    performance?: StepCheckPerformance | StepCheckMatcher;
    size?: number | Matcher[];
    co2?: number | Matcher[];
};
export declare type gRPCStepRequest = {
    proto?: string | string[];
    host: string;
    service: string;
    method: string;
    metadata?: gRPCRequestMetadata;
    data?: object | object[];
    tls?: Credential['tls'];
    size?: number;
};
export declare type gRPCStepResponse = {
    body: object | object[];
    duration: number;
    co2: number;
    size: number;
    status?: number;
    statusText?: string;
    metadata?: object;
};
export default function (params: gRPCStep, captures: CapturesStorage, schemaValidator: Ajv, options?: WorkflowOptions, config?: WorkflowConfig): Promise<StepRunResult>;
