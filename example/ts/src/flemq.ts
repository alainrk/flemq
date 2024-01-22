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
        console.log(response);
      },
      rejecter: (error: any) => {
        console.error(error);
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
        console.log("Received:", data.toString());
        this.handleResponse(data);
      });

      this.client.on("close", () => {
        console.log("Connection closed");
      });

      this.client.on("error", (error) => {
        console.error("An error occurred:", error);
        reject(error);
      });
    });
  }

  public send(data: string): Promise<any> {
    return new Promise((resolve, reject) => {
      this.currentHandler = {
        resolver: resolve.bind(this),
        // TODO: Timeout to reject the promise (excl. subscribe)
        rejecter: reject.bind(this),
      };
      this.client.write(data, (error) => {
        if (error) {
          reject(error);
        }
      });
    });
  }

  async push(topic: string, data: string): Promise<string> {
    if (this.options.serder === "base64") {
      data = Buffer.from(data).toString("base64");
    }

    return this.send(`PUSH ${topic} ${data}\r\n`);
  }
}
