"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
const parse_duration_1 = __importDefault(require("parse-duration"));
async function default_1(params) {
    const stepResult = {
        type: 'delay',
    };
    stepResult.type = 'delay';
    await new Promise((resolve) => setTimeout(resolve, typeof params === 'string' ? (0, parse_duration_1.default)(params) : params));
    return stepResult;
}
exports.default = default_1;
