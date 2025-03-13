/// <reference path="../pb_data/types.d.ts" />
/// <reference path="../webapp/static/emails/_emailData.d.ts" />

interface RecordModel<R extends Record<string, unknown>> extends core.Record {
    get: <K extends keyof R & string>(key: K) => R[K];
}

type EmailName = keyof EmailData;
type EmailAttributes<E extends EmailName> = EmailData[E][number];

type EmailRenderer = <E extends EmailName>(
    emailName: E,
    emailData: {
        [Key in EmailAttributes<E>]: string;
    }
) => { html: string; subject: string };
