<script setup lang="ts">
defineProps<{
  content: String;
  contentId: number
}>();

</script>

<template>

<q-card class="q-md-sm card-hover" bordered style="max-height: 100px; overflow-y: scroll;">
    <q-card-section>
      <div class="row justify-between items-center" @click.prevent="$emit('set-default', contentId)">
        <div class="text-body-1">{{ content }} </div>
        <q-btn icon="delete" color="red" flat round @click.prevent="$emit('delete-item', contentId)" />
      </div>
    </q-card-section>
</q-card>
</template>

<style scoped>
.card-hover {
  transition: transform 0.2s ease, box-shadow 0.2s ease;
}
.card-hover:hover {
  transform: scale(1.03);
  box-shadow: 0 8px 20px rgba(0, 0, 0, 0.2);
  z-index: 10;
}
</style>
