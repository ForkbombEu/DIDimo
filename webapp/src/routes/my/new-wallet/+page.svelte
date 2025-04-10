<!--
SPDX-FileCopyrightText: 2025 Forkbomb BV

SPDX-License-Identifier: AGPL-3.0-or-later
-->

<script lang="ts">
	import { createForm, Form } from "@/forms";
	import { Field } from "@/forms/fields";
	import SelectField from "@/forms/fields/selectField.svelte";
	import { pb } from "@/pocketbase/index.js";
	import type { GenericRecord } from "@/utils/types";
	import { zod } from "sveltekit-superforms/adapters";
    import { z } from "zod";

    const schema = z.object({
       name: z.string().min(1, "Name is required"),
       description: z.string(),
       playstore_url: z.string().url("Invalid URL"),
       appstore_url: z.string().url("Invalid URL"),
       repository: z.string().url("Invalid URL"),
       home_url: z.string().url("Invalid URL"),
       conformance_checks: z.array(z.string())
    });


    const form = createForm<GenericRecord>({
        adapter: zod(schema),
        onSubmit: async ({ form }) => {
            const { data } = form;
            // Handle form submission
            console.log("Form submitted:", data);
            pb.collection("wallets").create({...data, conformace_checks: {runs: data.conformance_checks}});
        },
        initialData: {
            name: "",
            description: "",
            playstore_url: "",
            appstore_url: "",
            repository: "",
            home_url: "",
            conformance_checks: []
        }
    });

    const { form: formData } = form;
    let { data } = $props();
	const { executions } = $derived(data);

    const items = executions.map((execution) => ({
        value: execution.execution.runId,
        label: execution.execution.workflowId
    }));


</script>

<Form {form}>
    <Field
        {form}
        name="name"
        options={{
            type: "text",
            label: "App Name",
            placeholder: "Enter app name"
        }}
    />
    <Field
        {form}
        name="description"
        options={{
            type: "textarea",
            label: "Description",
            placeholder: "Enter app description"
        }}
    />
    <Field
        {form}
        name="playstore_url"
        options={{
            type: "url",
            label: "Play Store URL",
            placeholder: "Enter Play Store URL"
        }}
    />
    <Field
        {form}
        name="appstore_url"
        options={{
            type: "url",
            label: "App Store URL",
            placeholder: "Enter App Store URL"
        }}
    />
    <Field
        {form}
        name="repository"
        options={{
            type: "url",
            label: "Repository URL",
            placeholder: "Enter repository URL"
        }}
    />
    <Field
        {form}
        name="home_url"
        options={{
            type: "url",
            label: "Home URL",
            placeholder: "Enter home URL"
        }}
    />
    <SelectField
        {form}
        name="conformance_checks"
        options={{
            type: "multiple",
            label: "Conformance Checks",
            items: items,
        }}
    />
</Form>

