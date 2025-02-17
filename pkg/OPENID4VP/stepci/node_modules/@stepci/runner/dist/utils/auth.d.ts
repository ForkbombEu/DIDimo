/// <reference types="node" />
import { StepFile, TryFileOptions } from './files';
export declare type Credential = {
    basic?: {
        username: string;
        password: string;
    };
    bearer?: {
        token: string;
    };
    oauth?: {
        endpoint: string;
        client_id: string;
        client_secret: string;
        audience?: string;
    };
    certificate?: {
        ca?: string | StepFile;
        cert?: string | StepFile;
        key?: string | StepFile;
        passphrase?: string;
    };
    tls?: {
        rootCerts?: string | StepFile;
        privateKey?: string | StepFile;
        certChain?: string | StepFile;
    };
};
export declare type CredentialsStorage = {
    [key: string]: Credential;
};
declare type OAuthClientConfig = {
    endpoint: string;
    client_id: string;
    client_secret: string;
    audience?: string;
};
export declare type OAuthResponse = {
    access_token: string;
    expires_in: number;
    token_type: string;
};
export declare type HTTPCertificate = {
    certificate?: string | Buffer;
    key?: string | Buffer;
    certificateAuthority?: string | Buffer;
    passphrase?: string;
};
export declare type TLSCertificate = {
    rootCerts?: string | Buffer;
    privateKey?: string | Buffer;
    certChain?: string | Buffer;
};
export declare function getOAuthToken(clientConfig: OAuthClientConfig): Promise<OAuthResponse>;
export declare function getAuthHeader(credential: Credential): Promise<string | undefined>;
export declare function getClientCertificate(certificate: Credential['certificate'], options?: TryFileOptions): Promise<HTTPCertificate | undefined>;
export declare function getTLSCertificate(certificate: Credential['tls'], options?: TryFileOptions): Promise<TLSCertificate | undefined>;
export {};
