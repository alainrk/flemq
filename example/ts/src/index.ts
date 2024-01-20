import * as net from "net";

type FlemQClientOptions = {
  host?: string;
  port: number;
};

class FlemQ {
  private client: net.Socket;
  private options: FlemQClientOptions;

  constructor(opt: FlemQClientOptions) {
    this.client = new net.Socket();
    this.options = opt;
  }

  // connect client to server
  async connect(): Promise<FlemQ> {
    return new Promise((resolve, reject) => {
      this.client.connect(
        this.options.port,
        this.options.host || "localhost",
        () => {
          console.log("Connection established");
          resolve(this);
        }
      );

      this.client.on("error", (err) => {
        reject(err);
      });
    });
  }

  // send data to server
  private async send(data: string): Promise<string> {
    return new Promise((resolve, reject) => {
      this.client.write(data, (err) => {
        if (err) {
          reject(err);
        } else {
          resolve("");
        }
      });

      this.client.on("data", (data) => {
        resolve(data.toString());
      });
    });
  }

  async push(topic: string, data: string): Promise<string> {
    return this.send(`PUSH ${topic} ${data}`);
  }
}

(async () => {
  const flemq = new FlemQ({
    port: 22123,
  });

  await flemq.connect();
  await flemq.push("x", "hello world");
})();
