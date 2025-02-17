import { code, parenless, noopTag as noop } from "./utils.mjs";
import { Jison } from "./lib/jison.mjs";

const _ = String.raw;

const bool = "std.coerceBoolean";
const operatorCode = code`ops['${2}'](${1}, ${3})`;

const grammar = {
  // Lexical tokens
  lex: {
    rules: [
      [_`\*`, `return "*" ;`],
      [_`\/`, `return "/" ;`],
      [_`-`, `return "-" ;`],
      [_`\+`, `return "+" ;`],
      [_`\^`, `return "^" ;`],
      [_`\(`, `return "(" ;`],
      [_`\)`, `return ")" ;`],
      [_`\,`, `return "," ;`],
      [_`==`, `return "==";`],
      [_`\!=`, `return "!=";`],
      [_`\~=`, `return "~=";`],
      [_`>=`, `return ">=";`],
      [_`<=`, `return "<=";`],
      [_`<`, `return "<" ;`],
      [_`>`, `return ">" ;`],
      [_`not\s+in[^\w]`, `return "notIn";`],
      [_`and[^\w]`, `return "and" ;`],
      [_`or[^\w]`, `return "or"  ;`],
      [_`not[^\w]`, `return "not" ;`],
      [_`in[^\w]`, `return "in"  ;`],
      [_`of[^\w]`, `return "of"  ;`],
      [_`if[^\w]`, `return "if"  ;`],
      [_`then[^\w]`, `return "then";`],
      [_`else[^\w]`, `return "else";`],
      [_`mod[^\w]`, `return "mod" ;`],

      [_`\s+`, ""], // skip whitespace
      [_`[0-9]+(?:\.[0-9]+)?(?![0-9\.])`, `return "Number";`], // 212.321

      [
        _`[a-zA-Z$_][\.a-zA-Z0-9$_]*`,
        `yytext = JSON.stringify({
                    name: yytext,
                    type: 'unescaped'
                });
                return "Symbol";`,
      ], // some.Symbol22

      [
        _`'(?:\\'|\\\\|[^'\\])*'`,
        `yytext = JSON.stringify({
                    name: yy.buildString("'", yytext),
                    type: 'single-quoted'
                });
                return "Symbol";`,
      ], // 'any \'escaped\' symbol'

      [
        _`"(?:\\"|\\\\|[^"\\])*"`,
        `yytext = JSON.stringify(yy.buildString('"', yytext));
                return "String";`,
      ], // "any \"escaped\" string"

      // Deprecated syntax
      [_`\%`, `return "%" ;`],
      [_`\?`, `return "?" ;`],
      [_`:`, `return ":" ;`],

      // End
      [_`$`, 'return "EndOfExpression";'],
    ],
  },
  // Operator precedence - lowest precedence first.
  // See http://www.gnu.org/software/bison/manual/html_node/Precedence.html
  // for a good explanation of how it works in Bison (and hence, Jison).
  // Different languages have different rules, but this seems a good starting
  // point: http://en.wikipedia.org/wiki/Order_of_operations#Programming_languages
  operators: [
    ["right", "if", "then", "else", /* deprecated: */ "?", ":"],
    ["left", "or"],
    ["left", "and"],
    ["left", "in", "notIn"],
    ["left", "NONCHAINEDREL"],
    ["left", "==", "!=", "<", "<=", ">", ">=", "~="],
    ["left", "CHAINEDREL"],
    ["left", "+", "-"],
    ["left", "*", "/", "mod", /* deprecated: */ "%"],
    ["left", "not", "UMINUS"],
    ["right", "^"],
    ["left", "of"],
  ],
  // Grammar
  bnf: {
    expressions: [
      // Entry point
      ["e EndOfExpression", "return $1;"],
    ],
    e: [
      ["- e", code`ops['-'](${2})`, { prec: "UMINUS" }],
      ["e + e", operatorCode],
      ["e - e", operatorCode],
      ["e * e", operatorCode],
      ["e / e", operatorCode],
      ["e ^ e", operatorCode],
      ["e mod e", code`ops.mod(${1}, ${3})`],

      ["e and e", code`${bool}(${1}) && ${bool}(${3})`],
      ["e or e", code`${bool}(${1}) || ${bool}(${3})`],
      ["not e", code`! ${bool}(${2})`],

      ["if e then e else e", code`${bool}(${2}) ? ${4} : ${6}`],
      ["e in e", code`std.isSubset(${1}, ${3})`],
      ["e notIn e", code`!std.isSubset(${1}, ${3})`],

      ["( e )", code`${2}`],
      ["( Arguments , e )", code`[ ${2}, ${4} ]`],

      ["Number", parenless`${1}`],
      ["Symbol", parenless`prop(${1}, data)`],
      ["String", parenless`${1}`],
      ["Symbol of e", parenless`prop(${1}, ${3})`],

      ["Symbol ( )", parenless`call(${1})`],
      ["Symbol ( Arguments )", parenless`call(${1}, ${3})`],

      ["Relation", `$$ = yy.reduceRelation($1);`, { prec: "==" }],

      // Deprecated
      [
        "e % e",
        parenless`std.warnDeprecated('modulo', ops['mod'](${1}, ${3}))`,
      ],
      [
        "e ? e : e",
        parenless`std.warnDeprecated('ternary', ${bool}(${1}) ? ${3} : ${5})`,
      ],
    ],
    RelationalOperator: [
      noop`==`,
      noop`!=`,
      noop`~=`,
      noop`<`,
      noop`<=`,
      noop`>=`,
      noop`>`,
    ],
    Relation: [
      [
        "e RelationalOperator Relation",
        `$$ = [$1, $2, ...$3]`,
        { prec: "CHAINEDREL" },
      ],
      [
        "e RelationalOperator e",
        `$$ = [$1, $2, $3];`,
        { prec: "NONCHAINEDREL" },
      ],
    ],
    Arguments: [
      ["e", parenless`${1}`],
      ["Arguments , e", parenless`${1}, ${3}`],
    ],
  },
};

export const parser = new Jison.Parser(grammar);
