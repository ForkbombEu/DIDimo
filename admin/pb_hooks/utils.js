/// <reference path="../pb_data/types.d.ts" />

module.exports = {
    /**
     * @param {core.RecordCreateEvent} e
     * @returns {models.Record | undefined}
     */
    getUserFromContext: (c) => {
        return $apis.requestInfo(c).authRecord;
    },
    /**
     * @returns {models.Record | undefined}
     */
    getOwnerRole: () => {
        const ownerRole = $app
            .dao()
            .findFirstRecordByData("orgRoles", "name", "owner");
        if (!ownerRole) throw new Error("missing owner role!");
        return ownerRole;
    },
};
