// SPDX-FileCopyrightText: 2024 The Forkbomb Company
//
// SPDX-License-Identifier: AGPL-3.0-or-later

import qrcode from 'qrcode-generator';

//

export function generateQrCode(text: string, cellSize = 20) {
	const qr = qrcode(0, 'L');
	qr.addData(text);
	qr.make();
	return qr.createDataURL(cellSize);
}
