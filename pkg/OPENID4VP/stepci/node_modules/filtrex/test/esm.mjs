import { compileExpression } from "../dist/esm/filtrex.mjs";
import { describe, it } from "mocha";
import { expect } from "chai";

const evl = (str, obj) => compileExpression(str)(obj);

describe("ES Modules", () => {
  it("works", () => {
    expect(evl("(1 + 2) * 3")).equals(9);
  });
});
