import * as net from "net";

export type FlemQSerDer = "base64";

export type FlemQClientOptions = {
  host?: string;
  port: number;
  serder: FlemQSerDer;
};

type Handler = {
  resolver: (response: any) => void;
  rejecter: (error: any) => void;
};

export class FlemQ {
  private client: net.Socket;
  private options: FlemQClientOptions;
  private currentHandler: Handler;

  constructor(opt: FlemQClientOptions) {
    this.client = new net.Socket();
    this.options = opt;
    // TODO: Handle concurrency, only one handler at a time and per client can be active (i.e. command await for a response)
    this.currentHandler = {
      resolver: (response: any) => {
        console.log("Default resolve handler:", response);
      },
      rejecter: (error: any) => {
        console.error("Default reject handler:", error);
      },
    };
  }

  private handleResponse(data: any) {
    this.currentHandler.resolver(data.toString());
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

  serialize(data: string): string {
    if (this.options.serder === "base64") {
      data = Buffer.from(data).toString("base64");
    }
    return data;
  }

  deserialize(data: string): string {
    if (this.options.serder === "base64") {
      data = Buffer.from(data, "base64").toString();
    }
    return data;
  }

  async push(topic: string, data: string): Promise<string> {
    data = this.serialize(data);

    return new Promise((resolve, reject) => {
      this.currentHandler = {
        resolver: resolve.bind(this),
        // TODO: Timeout to reject the promise (excl. subscribe)
        rejecter: reject.bind(this),
      };
      this.client.write(`PUSH ${topic} ${data}\r\n`, (error) => {
        if (error) {
          reject(error);
        }
      });
    });
  }

  handleSubscribeResponse(data: string) {
    // Loop over lines
    const lines = data.split("\r\n");
    for (let line of lines) {
      if (line.trim().length === 0) {
        continue;
      }
      // Remove the first character (the type), and deserialize the data
      console.log("Resolved:", this.deserialize(line.substring(1)));
    }
  }

  // TODO: This will get a handler that processes the data
  async subscribe(topic: string, offset = 0): Promise<void> {
    return new Promise((resolve, reject) => {
      this.currentHandler = {
        resolver: this.handleSubscribeResponse.bind(this),
        rejecter: reject.bind(this),
      };

      this.client.write(`SUBSCRIBE ${topic} ${offset}\r\n`, (error) => {
        if (error) {
          reject(error);
        }
        resolve();
      });
    });
  }
}
