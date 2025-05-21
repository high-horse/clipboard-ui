<script setup lang="ts">
defineProps<{
  content: String;
}>();
</script>

<template>

<q-card class="q-md-sm" bordered>
    <q-card-section>
        <div class="text-body-1">{{content}}</div>
    </q-card-section>
</q-card>
</template>