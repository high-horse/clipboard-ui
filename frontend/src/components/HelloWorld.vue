<script lang="ts" setup>
import { reactive, onMounted, onUnmounted, ref } from "vue";
import { Greet, Hide } from "../../wailsjs/go/main/App";
import { Add, GetAll, Remove, ClearHistory } from "../../wailsjs/go/clip/ClipboardManager";
import { QInput, Notify } from "quasar";
import {
    EventsOn,
    EventsOff,
    WindowHide,
} from "../../wailsjs/runtime/runtime.js";
import { ClipboardSetText } from "../../wailsjs/runtime";
import { CopiedContent } from "../type/types";

import { useQuasar } from "quasar";
import Card from "./common/Card.vue";

const $q = useQuasar();
const toggleDarkMode = () => {
    $q.dark.toggle();
};

const data = reactive({
    name: "",
    resultText: "Enter below to add new content to Clipboard 👇",
});

function greet() {
    Greet(data.name).then((result) => {
        data.resultText = result;
    });
}
async function addToClipboard() {
    try {
        const success = await ClipboardSetText(data.name);
        if (success) {
            data.name = "";
            await getHistory(); // Optional, if you're refreshing clipboard history
        } else {
            $q.notify({ message: "Failed to set clipboard", color: "red" });
        }
    } catch (err) {
        $q.notify({ message: "Clipboard error: " + err, color: "red" });
    }
}

async function handleDeleteClipboardContent(id: number) {
    try {
        await Remove(id);
        getHistory();
    } catch (err) {
        const errorMessage = err instanceof Error  ? err.message : "Error deleting clipboard content.";
        $q.notify({
            message: errorMessage,
            color: "red",
        });
    }
}

async function handleSetDefault(id: number) {
    try {
      // store locally
      const tempItem = messages.value?.find(item => item.key === id) || null;
      if (tempItem == null) {
        return;
      }
      // delete
      await handleDeleteClipboardContent(tempItem.key)
      // add to clipoard again
      await ClipboardSetText(tempItem.value);
    } catch (err) {
        const errorMessage = err instanceof Error ? err.message : "Unexpected Error Occured.";
        $q.notify({
            message: errorMessage,
            color: "red",
        });
    }
}

const message = ref<CopiedContent | null>(null);

const messages = ref<CopiedContent[]>([]);

EventsOn("new-content", async (data: { content: CopiedContent }) => {
  try{
    
    console.log("new cntent event ", data);
    messages.value.unshift(data.content);
    await getHistory();
  }catch(err) {
     await getHistory();
  }
});

EventsOn("history-cleared", () => {
  getHistory();
})

const getHistory = async () => {
    const response = await GetAll(); // response is array of strng
    // messages.value = response.reverse();
    messages.value = response;

    console.log("fetch history ", messages.value);
};

async function clearHistory() {
  try{
    await ClearHistory();
  }catch(err) {
    const errorMessage = err instanceof Error ? err.message : "Unexpected Error Occured.";
    $q.notify({
      message: errorMessage,
      color: "red",
    })
  }
}

onMounted(async () => {
    await getHistory();
});

const hide = async () => {
    // WindowHide();
    await Hide();
};
</script>

<template>
    <main>
        <!-- <div>
            <q-btn @click="hide()">Hide</q-btn>
        </div> -->
        <div id="result" class="result">{{ data.resultText }}</div>
        <q-form @submit.prevent="addToClipboard">
          <div class="row q-col-gutter-md q-mb-md">
              <div class="col-1">
                  <q-btn 
                  icon="refresh"
                  color="primary"
                  @click="getHistory"
                  />
              </div>
            <div class="col-9">
              <q-input
                outlined
                dense
                label="New Clipboard Content"
                v-model="data.name"
                type="text"
              />
            </div>
            <div class="col-1 flex items-center">
              <q-btn
                type="submit"
                color="primary"
                class="full-width text-capitalize"
              >
                  Submit
              </q-btn>
            </div>
            <div class="col-1">
              <q-btn-dropdown color="primary" icon="menu" dense flat>
                <q-list style="min-width: 150px">
                  <q-item clickable v-close-popup @click="clearHistory">
                    <q-item-section>Clear History</q-item-section>
                  </q-item>
                  <q-item clickable v-close-popup @click="hide">
                    <q-item-section>Hide</q-item-section>
                  </q-item>
                </q-list>
              </q-btn-dropdown>
            </div>

          </div>
        </q-form>


        <div class="row q-col-gutter-md q-mb-md" id="content-body">
            <!-- <div class="col-12" v-for="(msg, index) in messages" key="index">
                <Card :content="msg" />
            </div> -->
            <div class="col-12" v-for="(msg, index) in messages" :key="index">
                <Card
                    :content="msg.value"
                    :content-id="msg.key"
                    @set-default="handleSetDefault"
                    @delete-item="handleDeleteClipboardContent"
                />
            </div>
        </div>
    </main>
</template>

<style scoped></style>
