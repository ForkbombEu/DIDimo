import fs from 'fs';
import path from 'node:path';
import { formatCode, GENERATED, initAdminPocketbase, logCodegenResult } from '@/utils/codegen';

//

const TYPE_NAME = 'Feature';
const OBJECT_NAME = `${TYPE_NAME}s`;

const pb = await initAdminPocketbase();
const featuresRecords = await pb.collection('flags').getFullList();
const featuresEntries = featuresRecords.map((f) => `${f.name.toUpperCase()}: '${f.name}'`);

const code = `
export const ${OBJECT_NAME} = {
	${featuresEntries.join(',\n')}
} as const

export type ${TYPE_NAME} = typeof ${OBJECT_NAME} [keyof typeof ${OBJECT_NAME}];
`;

const formattedCode = await formatCode(code);
const filePath = path.join(import.meta.dirname, `features-list.${GENERATED}.ts`);
fs.writeFileSync(filePath, formattedCode);
logCodegenResult('features list', filePath);
