<script lang="ts" setup>
import { reactive, onMounted, onUnmounted, ref } from "vue";
import { Greet } from "../../wailsjs/go/main/App";
import {Add} from "../../wailsjs/go/clip/ClipboardManager";
import { QInput } from "quasar";
import { EventsOn, EventsOff } from "../../wailsjs/runtime/runtime.js";

const data = reactive({
    name: "",
    resultText: "Please enter your name below 👇",
});

function greet() {
    Greet(data.name).then((result) => {
        data.resultText = result;
    });
}

function addToClipboard() {
  Add(data.name)
  data.name = "";
}
import { useQuasar } from "quasar";
import Card from "./common/Card.vue";

const $q = useQuasar();
const toggleDarkMode = () => {
    $q.dark.toggle();
};

const message = ref("");

const messages = ref<string[]>([]);


EventsOn("new-content", (data: { content: string }) => {
    messages.value.unshift(data.content)
});
</script>

<template>
    <main>
        <div id="result" class="result">{{ data.resultText }}</div>
        <div class="row q-col-gutter-md q-mb-md">
            <div class="col-10">
                <q-input
                    outlined
                    dense
                    label="Greet Name"
                    v-model="data.name"
                    type="text"
                />
            </div>
            <div class="col-2">
                <!-- <q-btn color="primary" class="text-capitalize" @click="greet"
                    >Greet</q-btn
                > -->
                <q-btn color="primary" class="text-capitalize" @click="addToClipboard"
                    >Add To Clipboard</q-btn
                >
            </div>
        </div>

        <div class="row q-col-gutter-md q-mb-md" id="content-body">
            <div class="col-12" v-for="(msg, index) in messages" key="index">
                <Card :content="msg" />
            </div>
        </div>
    </main>
</template>

<style scoped></style>
