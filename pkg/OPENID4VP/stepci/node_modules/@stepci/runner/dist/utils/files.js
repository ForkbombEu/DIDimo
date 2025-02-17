"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.tryFile = void 0;
const fs_1 = __importDefault(require("fs"));
const path_1 = __importDefault(require("path"));
async function tryFile(input, options) {
    if (input.file) {
        return await fs_1.default.promises.readFile(path_1.default.join(path_1.default.dirname(options?.workflowPath || __dirname), input.file));
    }
    else {
        return input;
    }
}
exports.tryFile = tryFile;
