// SPDX-FileCopyrightText: 2024 The Forkbomb Company
//
// SPDX-License-Identifier: AGPL-3.0-or-later

import { z } from 'zod';
import _ from 'lodash';

//

const ZENROOM_NULL = 'null';

export const userChallengesSchema = z
	.object({
		whereParentsMet: z.string().default(''),
		nameFirstPet: z.string().default(''),
		whereHomeTown: z.string().default(''),
		nameFirstTeacher: z.string().default(''),
		nameMotherMaid: z.string().default('')
	})
	.partial()
	.transform(
		(obj) =>
			_.mapValues(obj, (v) => {
				if (typeof v == 'string' && v.length >= 1) return v;
				else return ZENROOM_NULL;
			}) as Required<typeof obj>
	)
	.refine((v) => {
		return Object.values(v).filter((v) => v === ZENROOM_NULL).length <= 2;
	}, 'AT_LEAST_THREE_QUESTIONS');

export type UserChallenges = z.infer<typeof userChallengesSchema>;
export type UserChallenge = keyof UserChallenges;

export const userChallenges: Array<{ id: UserChallenge; text: string }> = [
	{ id: 'whereParentsMet', text: 'Where did your parents meet?' },
	{ id: 'nameFirstPet', text: 'What is the name of your first pet?' },
	{ id: 'whereHomeTown', text: 'What is your home town?' },
	{ id: 'nameFirstTeacher', text: 'What is the name of your first teacher?' },
	{ id: 'nameMotherMaid', text: 'What is the surname of your mother before wedding?' }
];
