"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
const eventsource_1 = __importDefault(require("eventsource"));
const jsonpath_plus_1 = require("jsonpath-plus");
const { co2 } = require('@tgwf/co2');
const matcher_1 = require("../matcher");
const auth_1 = require("./../utils/auth");
async function default_1(params, captures, schemaValidator, options, config) {
    const stepResult = {
        type: 'sse',
    };
    const ssw = new co2();
    stepResult.type = 'sse';
    if (params.auth) {
        const authHeader = await (0, auth_1.getAuthHeader)(params.auth);
        if (authHeader) {
            if (!params.headers)
                params.headers = {};
            params.headers['Authorization'] = authHeader;
        }
    }
    await new Promise((resolve, reject) => {
        const ev = new eventsource_1.default(params.url || '', {
            headers: params.headers,
            rejectUnauthorized: config?.http?.rejectUnauthorized ?? false,
        });
        const messages = [];
        const timeout = setTimeout(() => {
            ev.close();
            const messagesBuffer = Buffer.from(messages.map((m) => m.data).join('\n'));
            stepResult.request = {
                url: params.url,
                headers: params.headers,
                size: 0,
            };
            stepResult.response = {
                contentType: 'text/event-stream',
                body: messagesBuffer,
                size: messagesBuffer.length,
                bodySize: messagesBuffer.length,
                co2: ssw.perByte(messagesBuffer.length),
                duration: params.timeout,
            };
            resolve(true);
        }, params.timeout || 10000);
        ev.onerror = (error) => {
            clearTimeout(timeout);
            ev.close();
            reject(error);
        };
        if (params.check) {
            if (!stepResult.checks)
                stepResult.checks = {};
            if (!stepResult.checks.messages)
                stepResult.checks.messages = {};
            params.check.messages?.forEach((check) => {
                ;
                (stepResult.checks?.messages)[check.id] = {
                    expected: check.body || check.json || check.jsonpath || check.schema,
                    given: undefined,
                    passed: false,
                };
            });
        }
        ev.onmessage = (message) => {
            messages.push(message);
            if (params.check) {
                params.check.messages?.forEach((check, id) => {
                    if (check.body) {
                        const result = (0, matcher_1.checkResult)(message.data, check.body);
                        if (result.passed && stepResult.checks?.messages)
                            stepResult.checks.messages[check.id] = result;
                    }
                    if (check.json) {
                        try {
                            const result = (0, matcher_1.checkResult)(JSON.parse(message.data), check.json);
                            if (result.passed && stepResult.checks?.messages)
                                stepResult.checks.messages[check.id] = result;
                        }
                        catch (e) {
                            reject(e);
                        }
                    }
                    if (check.schema) {
                        try {
                            const sample = JSON.parse(message.data);
                            const validate = schemaValidator.compile(check.schema);
                            const result = {
                                expected: check.schema,
                                given: sample,
                                passed: validate(sample),
                            };
                            if (result.passed && stepResult.checks?.messages)
                                stepResult.checks.messages[check.id] = result;
                        }
                        catch (e) {
                            reject(e);
                        }
                    }
                    if (check.jsonpath) {
                        try {
                            let jsonpathResult = {};
                            const json = JSON.parse(message.data);
                            for (const path in check.jsonpath) {
                                const result = (0, jsonpath_plus_1.JSONPath)({ path, json });
                                jsonpathResult[path] = (0, matcher_1.checkResult)(result[0], check.jsonpath[path]);
                            }
                            const passed = Object.values(jsonpathResult)
                                .map((c) => c.passed)
                                .every((passed) => passed);
                            if (passed && stepResult.checks?.messages)
                                stepResult.checks.messages[check.id] = {
                                    expected: check.jsonpath,
                                    given: jsonpathResult,
                                    passed,
                                };
                        }
                        catch (e) {
                            reject(e);
                        }
                    }
                });
            }
        };
    });
    return stepResult;
}
exports.default = default_1;
