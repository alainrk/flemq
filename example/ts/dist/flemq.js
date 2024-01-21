"use strict";
var __createBinding = (this && this.__createBinding) || (Object.create ? (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    var desc = Object.getOwnPropertyDescriptor(m, k);
    if (!desc || ("get" in desc ? !m.__esModule : desc.writable || desc.configurable)) {
      desc = { enumerable: true, get: function() { return m[k]; } };
    }
    Object.defineProperty(o, k2, desc);
}) : (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    o[k2] = m[k];
}));
var __setModuleDefault = (this && this.__setModuleDefault) || (Object.create ? (function(o, v) {
    Object.defineProperty(o, "default", { enumerable: true, value: v });
}) : function(o, v) {
    o["default"] = v;
});
var __importStar = (this && this.__importStar) || function (mod) {
    if (mod && mod.__esModule) return mod;
    var result = {};
    if (mod != null) for (var k in mod) if (k !== "default" && Object.prototype.hasOwnProperty.call(mod, k)) __createBinding(result, mod, k);
    __setModuleDefault(result, mod);
    return result;
};
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
exports.FlemQ = void 0;
const net = __importStar(require("net"));
class FlemQ {
    constructor(opt) {
        this.client = new net.Socket();
        this.options = opt;
    }
    // connect client to server
    connect() {
        return __awaiter(this, void 0, void 0, function* () {
            return new Promise((resolve, reject) => {
                this.client.connect(this.options.port, this.options.host || "localhost", () => {
                    resolve(this);
                });
                this.client.on("error", (err) => {
                    reject(err);
                });
            });
        });
    }
    // send data to server
    send(data) {
        return __awaiter(this, void 0, void 0, function* () {
            return new Promise((resolve, reject) => {
                this.client.write(data, (err) => {
                    if (err) {
                        reject(err);
                    }
                });
                this.client.on("data", (data) => {
                    resolve(data.toString());
                });
                this.client.on("error", (err) => {
                    reject(err);
                });
            });
        });
    }
    push(topic, data) {
        return __awaiter(this, void 0, void 0, function* () {
            if (this.options.serder === "base64") {
                data = btoa(data);
            }
            return this.send(`PUSH ${topic} ${data}`);
        });
    }
}
exports.FlemQ = FlemQ;
//# sourceMappingURL=flemq.js.map