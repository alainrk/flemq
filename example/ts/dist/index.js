"use strict";
var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
Object.defineProperty(exports, "__esModule", { value: true });
const flemq_1 = require("../src/flemq");
const sleep = (msec) => __awaiter(void 0, void 0, void 0, function* () {
    return new Promise((resolve) => {
        setTimeout(resolve, msec);
    });
});
// Publisher
(() => __awaiter(void 0, void 0, void 0, function* () {
    const flemq = new flemq_1.FlemQ({
        port: 22123,
        serder: "base64",
    });
    yield flemq.connect();
    for (let i = 0; i < 100; i++) {
        console.log(`Sending ${i}`);
        const res = yield flemq.push("ts_tests", `Hello from TS ${i}`);
        console.log("Res:", res);
        yield sleep(500);
    }
}))();
// Subscriber
(() => __awaiter(void 0, void 0, void 0, function* () {
    const flemq = new flemq_1.FlemQ({
        port: 22123,
        serder: "base64",
    });
    yield flemq.connect();
}))();
//# sourceMappingURL=index.js.map