"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
async function default_1(params, captures, cookies, schemaValidator, options, config) {
    const plugin = require(params.id);
    return plugin.default(params, captures, cookies, schemaValidator, options, config);
}
exports.default = default_1;
