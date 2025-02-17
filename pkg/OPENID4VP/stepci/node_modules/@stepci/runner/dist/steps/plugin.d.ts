import Ajv from 'ajv';
import { CookieJar } from 'tough-cookie';
import { CapturesStorage } from '../utils/runner';
import { WorkflowConfig, WorkflowOptions } from '..';
export declare type PluginStep = {
    id: string;
    params?: object;
    check?: object;
};
export default function (params: PluginStep, captures: CapturesStorage, cookies: CookieJar, schemaValidator: Ajv, options?: WorkflowOptions, config?: WorkflowConfig): Promise<any>;
