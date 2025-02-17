/// <reference types="node" />
import { Headers, PlainResponse } from 'got';
import FormData from 'form-data';
import Ajv from 'ajv';
import { CookieJar } from 'tough-cookie';
import { StepFile } from './../utils/files';
import { CapturesStorage } from './../utils/runner';
import { Credential } from './../utils/auth';
import { StepCheckCaptures, StepCheckJSONPath, StepCheckMatcher, StepCheckPerformance, StepCheckValue, StepRunResult, WorkflowConfig, WorkflowOptions } from '..';
import { Matcher } from '../matcher';
export declare type HTTPStepBase = {
    url: string;
    method: string;
    headers?: HTTPStepHeaders;
    params?: HTTPStepParams;
    cookies?: HTTPStepCookies;
    auth?: Credential;
    captures?: HTTPStepCaptures;
    check?: HTTPStepCheck;
    followRedirects?: boolean;
    timeout?: string | number;
    retries?: number;
};
export declare type HTTPStep = {
    body?: string | StepFile;
    form?: HTTPStepForm;
    formData?: HTTPStepMultiPartForm;
    json?: object;
    graphql?: HTTPStepGraphQL;
    trpc?: HTTPStepTRPC;
} & HTTPStepBase;
export declare type HTTPStepTRPC = {
    query?: {
        [key: string]: object;
    } | {
        [key: string]: object;
    }[];
    mutation?: {
        [key: string]: object;
    };
};
export declare type HTTPStepHeaders = {
    [key: string]: string;
};
export declare type HTTPStepParams = {
    [key: string]: string;
};
export declare type HTTPStepCookies = {
    [key: string]: string;
};
export declare type HTTPStepForm = {
    [key: string]: string;
};
export declare type HTTPRequestPart = {
    type?: string;
    value?: string;
    json?: object;
};
export declare type HTTPStepMultiPartForm = {
    [key: string]: string | StepFile | HTTPRequestPart;
};
export declare type HTTPStepGraphQL = {
    query: string;
    variables: object;
};
export declare type HTTPStepCaptures = {
    [key: string]: HTTPStepCapture;
};
export declare type HTTPStepCapture = {
    xpath?: string;
    jsonpath?: string;
    header?: string;
    selector?: string;
    cookie?: string;
    regex?: string;
    body?: boolean;
};
export declare type HTTPStepCheck = {
    status?: string | number | Matcher[];
    statusText?: string | Matcher[];
    redirected?: boolean;
    redirects?: string[];
    headers?: StepCheckValue | StepCheckMatcher;
    body?: string | Matcher[];
    json?: object;
    schema?: object;
    jsonpath?: StepCheckJSONPath | StepCheckMatcher;
    xpath?: StepCheckValue | StepCheckMatcher;
    selectors?: StepCheckValue | StepCheckMatcher;
    cookies?: StepCheckValue | StepCheckMatcher;
    captures?: StepCheckCaptures;
    sha256?: string;
    md5?: string;
    performance?: StepCheckPerformance | StepCheckMatcher;
    ssl?: StepCheckSSL;
    size?: number | Matcher[];
    requestSize?: number | Matcher[];
    bodySize?: number | Matcher[];
    co2?: number | Matcher[];
};
export declare type StepCheckSSL = {
    valid?: boolean;
    signed?: boolean;
    daysUntilExpiration?: number | Matcher[];
};
export declare type HTTPStepRequest = {
    protocol: string;
    url: string;
    method?: string;
    headers?: HTTPStepHeaders;
    body?: string | Buffer | FormData;
    size?: number;
};
export declare type HTTPStepResponse = {
    protocol: string;
    status: number;
    statusText?: string;
    duration?: number;
    contentType?: string;
    timings: PlainResponse['timings'];
    headers?: Headers;
    ssl?: StepResponseSSL;
    body: Buffer;
    co2: number;
    size?: number;
    bodySize?: number;
};
export declare type StepResponseSSL = {
    valid: boolean;
    signed: boolean;
    validUntil: Date;
    daysUntilExpiration: number;
};
export default function (params: HTTPStep, captures: CapturesStorage, cookies: CookieJar, schemaValidator: Ajv, options?: WorkflowOptions, config?: WorkflowConfig): Promise<StepRunResult>;
