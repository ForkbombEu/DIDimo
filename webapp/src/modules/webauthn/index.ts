import { pb } from '@/pocketbase';
import { log } from '@/utils/other';

//

export async function registerUser(username: string, description = '') {
	const credentialCreationOptions = await pb.send('/api/webauthn/register/begin/' + username, {});

	credentialCreationOptions.publicKey.challenge = bufferDecode(
		credentialCreationOptions.publicKey.challenge
	);
	credentialCreationOptions.publicKey.user.id = bufferDecode(
		credentialCreationOptions.publicKey.user.id
	);
	const credential = await navigator.credentials.create({
		publicKey: credentialCreationOptions.publicKey
	});

	if (credential != null) {
		const pkCred = credential as PublicKeyCredential;
		const attestationObject = (pkCred.response as AuthenticatorAttestationResponse)
			.attestationObject;
		const clientDataJSON = pkCred.response.clientDataJSON;
		const rawId = pkCred.rawId;

		await pb.send('/api/webauthn/register/finish/' + username, {
			body: {
				id: credential.id,
				rawId: bufferEncode(rawId),
				type: credential.type,
				response: {
					attestationObject: bufferEncode(attestationObject),
					clientDataJSON: bufferEncode(clientDataJSON)
				},
				description
			},
			method: 'POST'
		});
	}
}

export async function loginUser(username: string) {
	const credentialRequestOptions = await pb.send('/api/webauthn/login/begin/' + username, {});

	credentialRequestOptions.publicKey.challenge = bufferDecode(
		credentialRequestOptions.publicKey.challenge
	);
	credentialRequestOptions.publicKey.allowCredentials.forEach(function (listItem: {
		id: string | Uint8Array;
	}) {
		listItem.id = bufferDecode(listItem.id as string);
	});
	const credential = await navigator.credentials.get({
		publicKey: credentialRequestOptions.publicKey
	});

	if (credential != null) {
		const assertion = credential as PublicKeyCredential;
		const resp = assertion.response as AuthenticatorAssertionResponse;

		const authData = resp.authenticatorData;
		const clientDataJSON = resp.clientDataJSON;
		const rawId = assertion.rawId;
		const sig = resp.signature;
		const userHandle = resp.userHandle;
		if (!userHandle) throw new Error('WEBAUTHN: Missing user handle');

		const token = await pb.send('/api/webauthn/login/finish/' + username, {
			body: {
				id: assertion.id,
				rawId: bufferEncode(rawId),
				type: assertion.type,
				response: {
					authenticatorData: bufferEncode(authData),
					clientDataJSON: bufferEncode(clientDataJSON),
					signature: bufferEncode(sig),
					userHandle: bufferEncode(userHandle)
				}
			},
			method: 'POST'
		});
		log(token);
		pb.authStore.save(token.token, token.user);
	}
}

export function isWebauthnSupported() {
	return Boolean(window.PublicKeyCredential);
}

export async function isPlatformAuthenticatorAvailable() {
	if (isWebauthnSupported())
		return await PublicKeyCredential.isUserVerifyingPlatformAuthenticatorAvailable();
	else return false;
}

/* Utils */

// Base64 to ArrayBuffer
function bufferDecode(value: string) {
	return Uint8Array.from(atob(value.replace(/-/g, '+').replace(/_/g, '/')), (c) =>
		c.charCodeAt(0)
	);
}

function bufferEncode(buffer: ArrayBuffer) {
	let binary = '';
	const bytes = new Uint8Array(buffer);
	const len = bytes.byteLength;
	for (let i = 0; i < len; i++) {
		binary += String.fromCharCode(bytes[i]);
	}
	return btoa(binary).replace(/\+/g, '-').replace(/\//g, '_').replace(/=/g, '');
}
