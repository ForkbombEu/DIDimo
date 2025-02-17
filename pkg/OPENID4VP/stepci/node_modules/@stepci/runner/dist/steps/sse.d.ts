/// <reference types="node" />
import { CapturesStorage } from './../utils/runner';
import { Matcher } from '../matcher';
import Ajv from 'ajv';
import { StepCheckJSONPath, StepCheckMatcher, StepRunResult, WorkflowConfig, WorkflowOptions } from '..';
import { Credential } from './../utils/auth';
import { HTTPStepHeaders, HTTPStepParams } from './http';
export declare type SSEStep = {
    url: string;
    headers?: HTTPStepHeaders;
    params?: HTTPStepParams;
    auth?: Credential;
    json?: object;
    check?: {
        messages?: SSEStepCheck[];
    };
    timeout?: number;
};
export declare type SSEStepCheck = {
    id: string;
    json?: object;
    schema?: object;
    jsonpath?: StepCheckJSONPath | StepCheckMatcher;
    body?: string | Matcher[];
};
export declare type SSEStepRequest = {
    url?: string;
    headers?: HTTPStepHeaders;
    size?: number;
};
export declare type SSEStepResponse = {
    contentType?: string;
    duration?: number;
    body: Buffer;
    size?: number;
    bodySize?: number;
    co2: number;
};
export default function (params: SSEStep, captures: CapturesStorage, schemaValidator: Ajv, options?: WorkflowOptions, config?: WorkflowConfig): Promise<StepRunResult>;
