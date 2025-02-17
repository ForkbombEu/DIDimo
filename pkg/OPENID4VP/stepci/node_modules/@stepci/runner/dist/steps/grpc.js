"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
const node_path_1 = __importDefault(require("node:path"));
const jsonpath_plus_1 = require("jsonpath-plus");
const parse_duration_1 = __importDefault(require("parse-duration"));
const cool_grpc_1 = require("cool-grpc");
const { co2 } = require('@tgwf/co2');
const auth_1 = require("./../utils/auth");
const matcher_1 = require("../matcher");
async function default_1(params, captures, schemaValidator, options, config) {
    const stepResult = {
        type: 'grpc',
    };
    const ssw = new co2();
    // Load TLS configuration from file or string
    let tlsConfig;
    if (params.auth) {
        tlsConfig = await (0, auth_1.getTLSCertificate)(params.auth.tls, {
            workflowPath: options?.path,
        });
    }
    const protos = [];
    if (config?.grpc?.proto) {
        protos.push(...config.grpc.proto);
    }
    if (params.proto) {
        protos.push(...(Array.isArray(params.proto) ? params.proto : [params.proto]));
    }
    const proto = protos.map((p) => node_path_1.default.join(node_path_1.default.dirname(options?.path || __dirname), p));
    const request = {
        proto,
        host: params.host,
        metadata: params.metadata,
        service: params.service,
        method: params.method,
        data: params.data,
    };
    const { metadata, statusCode, statusMessage, data, size } = await (0, cool_grpc_1.makeRequest)(proto, {
        ...request,
        tls: tlsConfig,
        beforeRequest: (req) => {
            options?.ee?.emit('step:grpc_request', request);
        },
        afterResponse: (res) => {
            options?.ee?.emit('step:grpc_response', res);
        },
        options: {
            deadline: typeof params.timeout === 'string' ? (0, parse_duration_1.default)(params.timeout) : params.timeout
        }
    });
    stepResult.request = request;
    stepResult.response = {
        body: data,
        co2: ssw.perByte(size),
        size: size,
        status: statusCode,
        statusText: statusMessage,
        metadata,
    };
    // Captures
    if (params.captures) {
        for (const name in params.captures) {
            const capture = params.captures[name];
            if (capture.jsonpath) {
                captures[name] = (0, jsonpath_plus_1.JSONPath)({ path: capture.jsonpath, json: data })[0];
            }
        }
    }
    if (params.check) {
        stepResult.checks = {};
        // Check JSON
        if (params.check.json) {
            stepResult.checks.json = (0, matcher_1.checkResult)(data, params.check.json);
        }
        // Check Schema
        if (params.check.schema) {
            const validate = schemaValidator.compile(params.check.schema);
            stepResult.checks.schema = {
                expected: params.check.schema,
                given: data,
                passed: validate(data),
            };
        }
        // Check JSONPath
        if (params.check.jsonpath) {
            stepResult.checks.jsonpath = {};
            for (const path in params.check.jsonpath) {
                const result = (0, jsonpath_plus_1.JSONPath)({ path, json: data });
                stepResult.checks.jsonpath[path] = (0, matcher_1.checkResult)(result[0], params.check.jsonpath[path]);
            }
        }
        // Check captures
        if (params.check.captures) {
            stepResult.checks.captures = {};
            for (const capture in params.check.captures) {
                stepResult.checks.captures[capture] = (0, matcher_1.checkResult)(captures[capture], params.check.captures[capture]);
            }
        }
        // Check performance
        if (params.check.performance) {
            stepResult.checks.performance = {};
            if (params.check.performance.total) {
                stepResult.checks.performance.total = (0, matcher_1.checkResult)(stepResult.response?.duration, params.check.performance.total);
            }
        }
        // Check byte size
        if (params.check.size) {
            stepResult.checks.size = (0, matcher_1.checkResult)(size, params.check.size);
        }
        // Check co2 emissions
        if (params.check.co2) {
            stepResult.checks.co2 = (0, matcher_1.checkResult)(stepResult.response?.co2, params.check.co2);
        }
    }
    return stepResult;
}
exports.default = default_1;
