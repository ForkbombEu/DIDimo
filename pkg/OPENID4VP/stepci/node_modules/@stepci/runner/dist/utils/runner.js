"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.didChecksPass = exports.getCookie = exports.checkCondition = void 0;
const filtrex_1 = require("filtrex");
const flat_1 = __importDefault(require("flat"));
// Check if expression
function checkCondition(expression, data) {
    const filter = (0, filtrex_1.compileExpression)(expression);
    return filter((0, flat_1.default)(data));
}
exports.checkCondition = checkCondition;
// Get cookie
function getCookie(store, name, url) {
    return store.getCookiesSync(url).filter(cookie => cookie.key === name)[0]?.value;
}
exports.getCookie = getCookie;
// Did all checks pass?
function didChecksPass(checks) {
    if (!checks)
        return true;
    return Object.values(checks).map(check => {
        return check['passed'] ? check.passed : Object.values(check).map((c) => c.passed).every(passed => passed);
    })
        .every(passed => passed);
}
exports.didChecksPass = didChecksPass;
