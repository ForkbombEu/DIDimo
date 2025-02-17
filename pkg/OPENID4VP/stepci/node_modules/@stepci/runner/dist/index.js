"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.run = exports.runFromFile = exports.runFromYAML = void 0;
const tough_cookie_1 = require("tough-cookie");
const liquidless_1 = require("liquidless");
const liquidless_faker_1 = require("liquidless-faker");
const liquidless_naughtystrings_1 = require("liquidless-naughtystrings");
const fs_1 = __importDefault(require("fs"));
const js_yaml_1 = __importDefault(require("js-yaml"));
const json_schema_ref_parser_1 = __importDefault(require("@apidevtools/json-schema-ref-parser"));
const ajv_1 = __importDefault(require("ajv"));
const ajv_formats_1 = __importDefault(require("ajv-formats"));
const p_limit_1 = __importDefault(require("p-limit"));
const node_path_1 = __importDefault(require("node:path"));
const testdata_1 = require("./utils/testdata");
const runner_1 = require("./utils/runner");
const http_1 = __importDefault(require("./steps/http"));
const grpc_1 = __importDefault(require("./steps/grpc"));
const sse_1 = __importDefault(require("./steps/sse"));
const delay_1 = __importDefault(require("./steps/delay"));
const plugin_1 = __importDefault(require("./steps/plugin"));
const trpc_1 = __importDefault(require("./steps/trpc"));
const graphql_1 = __importDefault(require("./steps/graphql"));
const parse_duration_1 = __importDefault(require("parse-duration"));
const schema_1 = require("./utils/schema");
const templateDelimiters = ['${{', '}}'];
function renderObject(object, props) {
    return (0, liquidless_1.renderObject)(object, props, {
        filters: {
            fake: liquidless_faker_1.fake,
            naughtystring: liquidless_naughtystrings_1.naughtystring
        },
        delimiters: templateDelimiters
    });
}
// Run from test file
async function runFromYAML(yamlString, options) {
    const workflow = js_yaml_1.default.load(yamlString);
    const dereffed = await json_schema_ref_parser_1.default.dereference(workflow, {
        dereference: {
            circular: 'ignore'
        }
    });
    return run(dereffed, options);
}
exports.runFromYAML = runFromYAML;
// Run from test file
async function runFromFile(path, options) {
    const testFile = await fs_1.default.promises.readFile(path);
    return runFromYAML(testFile.toString(), { ...options, path });
}
exports.runFromFile = runFromFile;
// Run workflow
async function run(workflow, options) {
    const timestamp = new Date();
    const schemaValidator = new ajv_1.default({ strictSchema: false });
    (0, ajv_formats_1.default)(schemaValidator);
    // Templating for env, components, config
    let env = { ...workflow.env, ...options?.env };
    if (workflow.env) {
        env = renderObject(env, { env, secrets: options?.secrets });
    }
    if (workflow.components) {
        workflow.components = renderObject(workflow.components, { env, secrets: options?.secrets });
    }
    if (workflow.components?.schemas) {
        (0, schema_1.addCustomSchemas)(schemaValidator, workflow.components.schemas);
    }
    if (workflow.config) {
        workflow.config = renderObject(workflow.config, { env, secrets: options?.secrets });
    }
    if (workflow.include) {
        for (const workflowPath of workflow.include) {
            const testFile = await fs_1.default.promises.readFile(node_path_1.default.join(node_path_1.default.dirname(options?.path || __dirname), workflowPath));
            const test = js_yaml_1.default.load(testFile.toString());
            workflow.tests = { ...workflow.tests, ...test.tests };
        }
    }
    const concurrency = options?.concurrency || workflow.config?.concurrency || Object.keys(workflow.tests).length;
    const limit = (0, p_limit_1.default)(concurrency <= 0 ? 1 : concurrency);
    const testResults = [];
    const captures = {};
    // Run `before` section
    if (workflow.before) {
        const beforeResult = await runTest('before', workflow.before, schemaValidator, options, workflow.config, env, captures);
        testResults.push(beforeResult);
    }
    // Run `tests` section
    const input = [];
    Object.entries(workflow.tests).map(([id, test]) => input.push(limit(() => runTest(id, test, schemaValidator, options, workflow.config, env, { ...captures }))));
    testResults.push(...await Promise.all(input));
    // Run `after` section
    if (workflow.after) {
        const afterResult = await runTest('after', workflow.after, schemaValidator, options, workflow.config, env, captures);
        testResults.push(afterResult);
    }
    const workflowResult = {
        workflow,
        result: {
            tests: testResults,
            timestamp,
            passed: testResults.every(test => test.passed),
            duration: Date.now() - timestamp.valueOf(),
            co2: testResults.map(test => test.co2).reduce((a, b) => a + b),
            bytesSent: testResults.map(test => test.bytesSent).reduce((a, b) => a + b),
            bytesReceived: testResults.map(test => test.bytesReceived).reduce((a, b) => a + b),
        },
        path: options?.path
    };
    options?.ee?.emit('workflow:result', workflowResult);
    return workflowResult;
}
exports.run = run;
async function runTest(id, test, schemaValidator, options, config, env, capturesStorage) {
    const testResult = {
        id,
        name: test.name,
        steps: [],
        passed: true,
        timestamp: new Date(),
        duration: 0,
        co2: 0,
        bytesSent: 0,
        bytesReceived: 0
    };
    const captures = capturesStorage ?? {};
    const cookies = new tough_cookie_1.CookieJar();
    let previous;
    let testData = {};
    // Load test data
    if (test.testdata) {
        const parsedCSV = await (0, testdata_1.parseCSV)(test.testdata, { ...test.testdata.options, workflowPath: options?.path });
        testData = parsedCSV[Math.floor(Math.random() * parsedCSV.length)];
    }
    for (let step of test.steps) {
        const tryStep = async () => runStep(previous, step, id, test, captures, cookies, schemaValidator, testData, options, config, env);
        let stepResult = await tryStep();
        // Retries
        if ((stepResult.errored || (!stepResult.passed && !stepResult.skipped)) && step.retries && step.retries.count > 0) {
            for (let i = 0; i < step.retries.count; i++) {
                await new Promise(resolve => {
                    setTimeout(resolve, typeof step.retries?.interval === 'string' ? (0, parse_duration_1.default)(step.retries?.interval) : step.retries?.interval);
                });
                stepResult = await tryStep();
                if (stepResult.passed)
                    break;
            }
        }
        testResult.steps.push(stepResult);
        previous = stepResult;
        options?.ee?.emit('step:result', stepResult);
    }
    testResult.duration = Date.now() - testResult.timestamp.valueOf();
    testResult.co2 = testResult.steps.map(step => step.co2).reduce((a, b) => a + b);
    testResult.bytesSent = testResult.steps.map(step => step.bytesSent).reduce((a, b) => a + b);
    testResult.bytesReceived = testResult.steps.map(step => step.bytesReceived).reduce((a, b) => a + b);
    testResult.passed = testResult.steps.every(step => step.passed);
    options?.ee?.emit('test:result', testResult);
    return testResult;
}
async function runStep(previous, step, id, test, captures, cookies, schemaValidator, testData, options, config, env) {
    let stepResult = {
        id: step.id,
        testId: id,
        name: step.name,
        timestamp: new Date(),
        passed: true,
        errored: false,
        skipped: false,
        duration: 0,
        responseTime: 0,
        bytesSent: 0,
        bytesReceived: 0,
        co2: 0
    };
    let runResult;
    // Skip current step is the previous one failed or condition was unmet
    if (!config?.continueOnFail && (previous && !previous.passed)) {
        stepResult.passed = false;
        stepResult.errorMessage = 'Step was skipped because previous one failed';
        stepResult.skipped = true;
    }
    else if (step.if && !(0, runner_1.checkCondition)(step.if, { captures, env: { ...env, ...test.env } })) {
        stepResult.skipped = true;
        stepResult.errorMessage = 'Step was skipped because the condition was unmet';
    }
    else {
        try {
            step = renderObject(step, {
                captures,
                env: { ...env, ...test.env },
                secrets: options?.secrets,
                testdata: testData
            });
            if (step.http) {
                runResult = await (0, http_1.default)(step.http, captures, cookies, schemaValidator, options, config);
            }
            if (step.trpc) {
                runResult = await (0, trpc_1.default)(step.trpc, captures, cookies, schemaValidator, options, config);
            }
            if (step.graphql) {
                runResult = await (0, graphql_1.default)(step.graphql, captures, cookies, schemaValidator, options, config);
            }
            if (step.grpc) {
                runResult = await (0, grpc_1.default)(step.grpc, captures, schemaValidator, options, config);
            }
            if (step.sse) {
                runResult = await (0, sse_1.default)(step.sse, captures, schemaValidator, options, config);
            }
            if (step.delay) {
                runResult = await (0, delay_1.default)(step.delay);
            }
            if (step.plugin) {
                runResult = await (0, plugin_1.default)(step.plugin, captures, cookies, schemaValidator, options, config);
            }
            stepResult.passed = (0, runner_1.didChecksPass)(runResult?.checks);
        }
        catch (error) {
            stepResult.passed = false;
            stepResult.errored = true;
            stepResult.errorMessage = error.message;
            options?.ee?.emit('step:error', error);
        }
    }
    stepResult.type = runResult?.type;
    stepResult.request = runResult?.request;
    stepResult.response = runResult?.response;
    stepResult.checks = runResult?.checks;
    stepResult.responseTime = runResult?.response?.duration || 0;
    stepResult.co2 = runResult?.response?.co2 || 0;
    stepResult.bytesSent = runResult?.request?.size || 0;
    stepResult.bytesReceived = runResult?.response?.size || 0;
    stepResult.duration = Date.now() - stepResult.timestamp.valueOf();
    stepResult.captures = Object.keys(captures).length > 0 ? captures : undefined;
    stepResult.cookies = Object.keys(cookies.toJSON().cookies).length > 0 ? cookies.toJSON().cookies : undefined;
    return stepResult;
}
