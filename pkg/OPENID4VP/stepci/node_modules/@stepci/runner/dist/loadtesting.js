"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.loadTest = exports.loadTestFromFile = void 0;
const fs_1 = __importDefault(require("fs"));
const js_yaml_1 = __importDefault(require("js-yaml"));
const json_schema_ref_parser_1 = __importDefault(require("@apidevtools/json-schema-ref-parser"));
const phasic_1 = require("phasic");
const simple_statistics_1 = require("simple-statistics");
const index_1 = require("./index");
const matcher_1 = require("./matcher");
function metricsResult(numbers) {
    return {
        min: (0, simple_statistics_1.min)(numbers),
        max: (0, simple_statistics_1.max)(numbers),
        avg: (0, simple_statistics_1.mean)(numbers),
        med: (0, simple_statistics_1.median)(numbers),
        p95: (0, simple_statistics_1.quantile)(numbers, 0.95),
        p99: (0, simple_statistics_1.quantile)(numbers, 0.99),
    };
}
async function loadTestFromFile(path, options) {
    const testFile = await fs_1.default.promises.readFile(path);
    const workflow = js_yaml_1.default.load(testFile.toString());
    const dereffed = await json_schema_ref_parser_1.default.dereference(workflow, {
        dereference: {
            circular: 'ignore'
        }
    });
    return loadTest(dereffed, { ...options, path });
}
exports.loadTestFromFile = loadTestFromFile;
// Load-testing functionality
async function loadTest(workflow, options) {
    if (!workflow.config?.loadTest?.phases)
        throw Error('No load test config detected');
    const start = new Date();
    const resultList = await (0, phasic_1.runPhases)(workflow.config?.loadTest?.phases, () => (0, index_1.run)(workflow, options));
    const results = resultList.map(result => result.value.result);
    // Tests metrics
    const testsPassed = results.filter((r) => r.passed === true).length;
    const testsFailed = results.filter((r) => r.passed === false).length;
    // Steps metrics
    const steps = results.map(r => r.tests).map(test => test.map(test => test.steps)).flat(2);
    const stepsPassed = steps.filter(step => step.passed === true).length;
    const stepsFailed = steps.filter(step => step.passed === false).length;
    const stepsSkipped = steps.filter(step => step.skipped === true).length;
    const stepsErrored = steps.filter(step => step.errored === true).length;
    // Response metrics
    const responseTime = metricsResult(steps.map(step => step.responseTime));
    // Size Metrics
    const bytesSent = results.map(result => result.bytesSent).reduce((a, b) => a + b);
    const bytesReceived = results.map(result => result.bytesReceived).reduce((a, b) => a + b);
    const co2 = results.map(result => result.co2).reduce((a, b) => a + b);
    // Checks
    let checks;
    if (workflow.config?.loadTest?.check) {
        checks = {};
        if (workflow.config?.loadTest?.check.min) {
            checks.min = (0, matcher_1.checkResult)(responseTime.min, workflow.config?.loadTest?.check.min);
        }
        if (workflow.config?.loadTest?.check.max) {
            checks.max = (0, matcher_1.checkResult)(responseTime.max, workflow.config?.loadTest?.check.max);
        }
        if (workflow.config?.loadTest?.check.avg) {
            checks.avg = (0, matcher_1.checkResult)(responseTime.avg, workflow.config?.loadTest?.check.avg);
        }
        if (workflow.config?.loadTest?.check.med) {
            checks.med = (0, matcher_1.checkResult)(responseTime.med, workflow.config?.loadTest?.check.med);
        }
        if (workflow.config?.loadTest?.check.p95) {
            checks.p95 = (0, matcher_1.checkResult)(responseTime.p95, workflow.config?.loadTest?.check.p95);
        }
        if (workflow.config?.loadTest?.check.p99) {
            checks.p99 = (0, matcher_1.checkResult)(responseTime.p99, workflow.config?.loadTest?.check.p99);
        }
    }
    const result = {
        workflow,
        result: {
            stats: {
                steps: {
                    failed: stepsFailed,
                    passed: stepsPassed,
                    skipped: stepsSkipped,
                    errored: stepsErrored,
                    total: steps.length
                },
                tests: {
                    failed: testsFailed,
                    passed: testsPassed,
                    total: results.length
                },
            },
            responseTime,
            bytesSent,
            bytesReceived,
            co2,
            rps: steps.length / ((Date.now() - start.valueOf()) / 1000),
            iterations: results.length,
            duration: Date.now() - start.valueOf(),
            checks,
            passed: checks ? Object.entries(checks).map(([i, check]) => check.passed).every(passed => passed) : true
        }
    };
    options?.ee?.emit('loadtest:result', result);
    return result;
}
exports.loadTest = loadTest;
