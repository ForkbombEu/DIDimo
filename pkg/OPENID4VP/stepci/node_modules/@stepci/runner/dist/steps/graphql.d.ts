import Ajv from 'ajv';
import { CookieJar } from 'tough-cookie';
import { CapturesStorage } from '../utils/runner';
import { WorkflowConfig, WorkflowOptions } from '..';
import { HTTPStepBase, HTTPStepGraphQL } from './http';
export declare type GraphQLStep = HTTPStepGraphQL & HTTPStepBase;
export default function (params: GraphQLStep, captures: CapturesStorage, cookies: CookieJar, schemaValidator: Ajv, options?: WorkflowOptions, config?: WorkflowConfig): Promise<import("..").StepRunResult>;
