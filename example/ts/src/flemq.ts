import * as net from "net";
import { assertPositiveInteger } from "./common";

export type FlemQSerDer = "base64";
export type FlepType = string;

export type FlemQClientOptions = {
  host?: string;
  port: number;
  serder: FlemQSerDer;
};

export type FlemQHandler = (response: FlepResponse) => void;
export type FlepResponse = {
  type: string;
  data: string | null;
};

type Handler = (response: FlepResponse) => void;

export const FlepTypeString = "+";
export const FlepTypeError = "-";

export class FlemQ {
  private client: net.Socket;
  private options: FlemQClientOptions;
  private handler: Handler | null = null;

  constructor(opt: FlemQClientOptions) {
    this.client = new net.Socket();
    this.options = opt;
  }

  private deserialize(line: string): FlepResponse {
    let [flepType, data] = [line[0], line.substring(1)];
    let response: FlepResponse = { type: flepType, data: null };

    if ([FlepTypeString].indexOf(flepType) !== -1) {
      if (this.options.serder === "base64") {
        response.data = Buffer.from(data, "base64").toString();
      }
    } else if ([FlepTypeError].indexOf(flepType) !== -1) {
      // No need to deserialize errors coming from the server
      response.data = data;
    }

    return response;
  }

  private serialize(data: string): string {
    if (this.options.serder === "base64") {
      data = Buffer.from(data).toString("base64");
    }
    return data;
  }

  /**
   * Handle the response from the server, checking if the handler is set
   * @throws Error if handler is already set
   */
  private handleResponse(data: string) {
    if (this.handler == null) {
      throw new Error("Handler not set");
    }

    data = data.toString();

    const lines = data.split("\r\n");
    for (let line of lines) {
      if (line.trim().length === 0) {
        continue;
      }
      // Remove the first character (the type), and deserialize the data
      this.handler(this.deserialize(line));
    }
  }

  // connect client to server
  async connect(): Promise<FlemQ> {
    return new Promise((resolve, reject) => {
      this.client.connect(
        this.options.port,
        this.options.host || "localhost",
        () => {
          resolve(this);
        }
      );

      this.client.on("data", (data) => {
        this.handleResponse(data.toString());
      });

      this.client.on("close", () => {
        // TODO: Handle close
        console.log("Connection closed");
      });

      this.client.on("error", (error) => {
        reject(error);
      });
    });
  }

  push(topic: string, data: string, handler: Handler) {
    data = this.serialize(data);

    this.handler = handler.bind(this);
    this.client.write(`PUSH ${topic} ${data}\r\n`, (error) => {
      if (error) {
        throw error;
      }
    });
  }

  pick(topic: string, offset: number, handler: Handler) {
    assertPositiveInteger(offset);
    this.handler = handler.bind(this);

    this.client.write(`PICK ${topic} ${offset}\r\n`, (error) => {
      if (error) {
        throw error;
      }
    });
  }

  subscribe(topic: string, handler: Handler, offset = 0) {
    assertPositiveInteger(offset);
    this.handler = handler.bind(this);

    this.client.write(`SUBSCRIBE ${topic} ${offset}\r\n`, (error) => {
      if (error) {
        throw error;
      }
    });
  }
}
