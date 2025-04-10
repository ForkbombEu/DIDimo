<script setup lang="ts">
import { footer } from "./data";
</script>

<template>
    <div class="custom-footer">
        <div class="container">
            <div class="columns">
                <div
                    v-for="(links, category) in footer"
                    :key="category"
                    class="column"
                >
                    <h2>{{ category }}</h2>
                    <ul>
                        <li v-for="link in links" :key="link.url">
                            <a :href="link.url">{{ link.label }}</a>
                        </li>
                    </ul>
                </div>
            </div>
        </div>
    </div>
</template>

<style scoped>
.custom-footer {
    padding: 32px 12px;
}

.container {
    margin: 0 auto;
    max-width: var(--vp-layout-max-width);
}

.columns {
    display: grid;
    gap: 24px;
    grid-template-columns: repeat(1, 1fr);
}

@media (min-width: 640px) {
    .columns {
        grid-template-columns: repeat(2, 1fr);
    }
}

@media (min-width: 960px) {
    .columns {
        grid-template-columns: repeat(3, 1fr);
    }
}

@media (min-width: 1280px) {
    .columns {
        grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
    }
}

.column h2 {
    font-size: 1rem;
    font-weight: 600;
    margin-bottom: 12px;
    color: var(--vp-c-text-1);
}

.column ul {
    list-style: none;
    padding: 0;
    margin: 0;
}

.column li {
    margin-bottom: 8px;
}

.column a {
    font-size: 0.875rem;
    color: var(--vp-c-text-2);
    transition: color 0.25s;
    text-decoration: none;
}

.column a:hover {
    color: var(--vp-c-brand-1);
}
</style>
