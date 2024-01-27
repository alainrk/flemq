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

type Handler = (response: string) => void;

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

  // private deserialize(line: string): FlepResponse {
  private deserialize(line: string): string {
    let [flepType, data] = [line[0], line.substring(1)];
    let response: FlepResponse = { type: flepType, data: null };

    if ([FlepTypeError, FlepTypeString].indexOf(flepType) !== -1) {
      if (this.options.serder === "base64") {
        response.data = Buffer.from(data, "base64").toString();
      }
    }

    // return response;
    return response.data || "";
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
  private handleResponse(data: any) {
    if (this.handler == null) {
      throw new Error("Handler not set");
    }
    this.handler(data.toString());
  }

  // TODO: Handle multiple commands waiting for a response at the same time (if needed)
  private setHandler(handler: Handler) {
    this.handler = handler;
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
        this.handleResponse(data);
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

  async push(topic: string, data: string): Promise<string> {
    data = this.serialize(data);

    return new Promise((resolve, reject) => {
      this.setHandler(resolve.bind(this));
      this.client.write(`PUSH ${topic} ${data}\r\n`, (error) => {
        if (error) {
          reject(error);
        }
      });
    });
  }

  async pick(topic: string, offset: number): Promise<string> {
    assertPositiveInteger(offset);

    return new Promise((resolve, reject) => {
      const handler = (data: string) => {
        data = this.deserialize(data.substring(1));
        resolve(data);
      };
      this.handler = handler.bind(this);

      this.client.write(`PICK ${topic} ${offset}\r\n`, (error) => {
        if (error) {
          reject(error);
        }
      });
    });
  }

  async subscribe(topic: string, handler: Handler, offset = 0): Promise<void> {
    assertPositiveInteger(offset);

    const handleSubscribeResponse = (data: string) => {
      const lines = data.split("\r\n");
      for (let line of lines) {
        if (line.trim().length === 0) {
          continue;
        }
        // Remove the first character (the type), and deserialize the data
        handler(this.deserialize(line));
      }
    };

    return new Promise((resolve, reject) => {
      this.setHandler(handleSubscribeResponse.bind(this));

      this.client.write(`SUBSCRIBE ${topic} ${offset}\r\n`, (error) => {
        if (error) {
          reject(error);
        }
        resolve();
      });
    });
  }
}
