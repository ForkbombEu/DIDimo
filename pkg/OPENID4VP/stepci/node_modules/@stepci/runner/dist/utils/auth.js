"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.getTLSCertificate = exports.getClientCertificate = exports.getAuthHeader = exports.getOAuthToken = void 0;
const got_1 = __importDefault(require("got"));
const files_1 = require("./files");
async function getOAuthToken(clientConfig) {
    return await got_1.default.post(clientConfig.endpoint, {
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({
            grant_type: 'client_credentials',
            client_id: clientConfig.client_id,
            client_secret: clientConfig.client_secret,
            audience: clientConfig.audience
        })
    })
        .json();
}
exports.getOAuthToken = getOAuthToken;
async function getAuthHeader(credential) {
    if (credential.basic) {
        return 'Basic ' + Buffer.from(credential.basic.username + ':' + credential.basic.password).toString('base64');
    }
    if (credential.bearer) {
        return 'Bearer ' + credential.bearer.token;
    }
    if (credential.oauth) {
        const { access_token } = await getOAuthToken(credential.oauth);
        return 'Bearer ' + access_token;
    }
}
exports.getAuthHeader = getAuthHeader;
async function getClientCertificate(certificate, options) {
    if (certificate) {
        const cert = {};
        if (certificate.cert) {
            cert.certificate = await (0, files_1.tryFile)(certificate.cert, options);
        }
        if (certificate.key) {
            cert.key = await (0, files_1.tryFile)(certificate.key, options);
        }
        if (certificate.ca) {
            cert.certificateAuthority = await (0, files_1.tryFile)(certificate.ca, options);
        }
        if (certificate.passphrase) {
            cert.passphrase = certificate.passphrase;
        }
        return cert;
    }
}
exports.getClientCertificate = getClientCertificate;
async function getTLSCertificate(certificate, options) {
    if (certificate) {
        const tlsConfig = {};
        if (certificate.rootCerts) {
            tlsConfig.rootCerts = await (0, files_1.tryFile)(certificate.rootCerts, options);
        }
        if (certificate.privateKey) {
            tlsConfig.privateKey = await (0, files_1.tryFile)(certificate.privateKey, options);
        }
        if (certificate.certChain) {
            tlsConfig.certChain = await (0, files_1.tryFile)(certificate.certChain, options);
        }
        return tlsConfig;
    }
}
exports.getTLSCertificate = getTLSCertificate;
