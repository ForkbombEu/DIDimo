import Ajv from 'ajv';
import { CookieJar } from 'tough-cookie';
import { CapturesStorage } from '../utils/runner';
import { WorkflowConfig, WorkflowOptions } from '..';
import { HTTPStepBase, HTTPStepTRPC } from './http';
export declare type tRPCStep = HTTPStepTRPC & HTTPStepBase;
export default function (params: tRPCStep, captures: CapturesStorage, cookies: CookieJar, schemaValidator: Ajv, options?: WorkflowOptions, config?: WorkflowConfig): Promise<import("..").StepRunResult>;
