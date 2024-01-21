export type FlemQSerDer = "base64";
export type FlemQClientOptions = {
    host?: string;
    port: number;
    serder: FlemQSerDer;
};
export declare class FlemQ {
    private client;
    private options;
    constructor(opt: FlemQClientOptions);
    connect(): Promise<FlemQ>;
    private send;
    push(topic: string, data: string): Promise<string>;
}
