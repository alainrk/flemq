export type FlemQSerDer = "base64";
export type FlemQClientOptions = {
    host?: string;
    port: number;
    serder: FlemQSerDer;
};
export declare class FlemQ {
    private client;
    private options;
    private currentHandler;
    constructor(opt: FlemQClientOptions);
    private handleResponse;
    connect(): Promise<FlemQ>;
    send(data: string): Promise<any>;
    push(topic: string, data: string): Promise<string>;
}
