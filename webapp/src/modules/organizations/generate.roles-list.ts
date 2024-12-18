import fs from 'fs';
import path from 'node:path';
import { formatCode, GENERATED, initAdminPocketbase, logCodegenResult } from '@/utils/codegen';

//

const TYPE_NAME = 'OrgRole';
const OBJECT_NAME = `${TYPE_NAME}s`;

const pb = await initAdminPocketbase();
const rolesRecords = await pb.collection('orgRoles').getFullList();
const rolesEntries = rolesRecords.map((r) => `${r.name.toUpperCase()}: '${r.name}'`);

const code = `
export const ${OBJECT_NAME} = {
	${rolesEntries.join(',\n')}
} as const

export type ${TYPE_NAME} = typeof ${OBJECT_NAME} [keyof typeof ${OBJECT_NAME}];
`;

const formattedCode = await formatCode(code);
const filePath = path.join(import.meta.dirname, `roles-list.${GENERATED}.ts`);
fs.writeFileSync(filePath, formattedCode);
logCodegenResult('organization roles list', filePath);
