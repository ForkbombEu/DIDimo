# Organizations

This docs page explain the information flows relative to:

- `organization membership request`
- `organization invite`

## User requests membership

### Business logic

- `🟢 user` sends membership request to `💾 pb`

  - `🖥️ ui` makes a `POST` to `api/collections/orgJoinRequests`
  - `body`: `{organizationId: string, userId: string}`

- This triggers a `hook` on `💾 pb` that sends an `✉️ email` to the `🧑‍💼 admins`

  - `hook`: `onRecordAfterCreateRequest` on `orgJoinRequests` collection.

- `🧑‍💼 admins` from the `✉️ email`, are redirected to `🖥️ ui`: `/my/organizations/[id]/members`

- `🧑‍💼 admins` in the `🖥️ ui` have two options:

  - accept invite
    - send `POST` to `<pb>/api/collections/orgJoinRequests`
    - `body`: `{status: "accepted"}`
  - decline invite
    - send `POST` to `<pb>/api/collections/orgJoinRequests`
    - `body`: `{status: "declined"}`

- `hook` on `💾 pb` sends a notification `✉️ email` to the `🟢 user` about the status change

  - `hook`: `onRecordAfterUpdateRequest` on `orgJoinRequests` collection

### Relevant files

- `admin/pb_hooks/organizations_request_membership.pb.js`
- `webapp/src/routes/[[lang]]/my/organizations/join/+page.svelte`
- `webapp/src/routes/[[lang]]/my/organizations/[id]/members/**/*`

## Admin invites users

### Business logic

- `🧑‍💼 admins` invite `🟢 users`

  - `🖥️ ui` sends `POST` to `<pb>/organizations/invite`
  - `body`: `{organizationId: string, emails: string[]}`

- `💾 pb` handles the route:

  - `org_invite` `records` are created
  - `✉️ email` is sent to `🟢 users` with `url` `<ui>/organization-invite-[orgId]-[inviteId]-[email]-[[userId]]`

    - `[[userId]]` is optional as `🟢 user` may be not registered on the platform

- `🟢 user` clicks on the link in the `✉️ email` and is redirected to the `🖥️ ui`

- `🖥️ ui` starts an `OrganizationInviteSession` (saved in sessionStorage)

- Two things can happen now:

  - `Case A`: the `🟢 user` has already an `account`

    - `A.1`: is signed in
      - (pass)
    - `A.2`: is `not` signed in
      - `redirect` to `<ui>/login`
      - `🟢 user` completes login

  - `Case B`, the `🟢 user` is `not registered`

    - `redirect` to `<ui>/register`
    - `🟢 user` completes registration

- `🖥️ ui` ends the `OrganizationInviteSession`

- `🖥️ ui` makes a `redirect` to `<ui>/my/organizations`

  - Invite is displayed on screen

- `🟢 user` can now:

  - accept the invitation
    - `POST` to `/organizations/invite/accept`
    - `body`: `{inviteId:string}`
  - decline the invitation
    - `POST` to `/organizations/invite/decline`
    - `body`: `{inviteId:string}`

### Relevant files

- `admin/pb_hooks/organizations_invites.pb.js`
- `webapp/src/routes/[[lang]]/(nru)/organization-invite-[orgId]-[inviteId]-[email]-[[userId]]/+page.ts`
- `webapp/src/routes/[[lang]]/my/+layout.ts`
- `webapp/src/routes/[[lang]]/my/organizations/+page.svelte`
- `webapp/src/routes/[[lang]]/my/organizations/[id]/members/+page.svelte`
