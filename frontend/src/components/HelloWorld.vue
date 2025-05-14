<script lang="ts" setup>
import { reactive } from "vue";
import { Greet } from "../../wailsjs/go/main/App";
import { QInput } from "quasar";

const data = reactive({
    name: "",
    resultText: "Please enter your name below 👇",
});

function greet() {
    Greet(data.name).then((result) => {
        data.resultText = result;
    });
}
import { useQuasar } from 'quasar'

const $q = useQuasar()
const toggleDarkMode = () => {
  $q.dark.toggle()
}
</script>

<template>
    <main>
        <div id="result" class="result">{{ data.resultText }}</div>
        <div class="row q-col-gutter-md q-mb-md">
            <div class="col-11">
                 <q-input outlined dense label="Greet Name" v-model="data.name" type="text"/>
            </div>
            <div class="col-1">
                <q-btn color="primary" class="text-capitalize" @click="greet">Greet</q-btn>
            </div>
        </div>
        <!-- <q-btn label="Toggle Mode" class="text-capitalize " color="primary" @click.prevent="toggleDarkMode"></q-btn> -->
    </main>
</template>

<style scoped>
</style>
