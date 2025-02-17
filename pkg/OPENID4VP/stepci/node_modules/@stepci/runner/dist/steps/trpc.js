"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
const http_1 = __importDefault(require("./http"));
async function default_1(params, captures, cookies, schemaValidator, options, config) {
    return (0, http_1.default)({
        trpc: {
            query: params.query,
            mutation: params.mutation,
        },
        ...params,
    }, captures, cookies, schemaValidator, options, config);
}
exports.default = default_1;
