<script setup lang="ts">
import { onMounted, ref } from "vue";
import videojs, { type VideoJsPlayer } from "video.js";
import { useRoute, useRouter } from "vue-router";
import {
  KecpSignal,
  KecpConnection,
  AnswerPeer,
  KecpEventType,
  KecpMessageType,
  type KecpMessage,
} from "kecp-webrtc";
import {
  ElMessage,
  ElMessageBox,
  ElScrollbar,
  type UploadFile,
  type UploadFiles,
} from "element-plus";
import { kecpEndpoint } from "@/config";

type KecpChat = KecpMessage & { chatID: number };
interface HTMLMediaElementWithCaputreStream extends HTMLMediaElement {
  captureStream(): MediaStream;
}

const videoPlayerEl = ref<HTMLMediaElementWithCaputreStream>();
const player = ref<VideoJsPlayer>();
const roomID = ref("");
const username = ref("");
const usernameInput = ref("");
const conn = ref<KecpConnection>();
const userList = ref<string[]>([]);
const targetUser = ref<string>("");
const chatList = ref<KecpChat[]>([]);
const sig = new KecpSignal(kecpEndpoint);
const connected = ref(false);
const sendValue = ref("");
const chatIDCouter = ref(0);
const chatBox = ref<HTMLElement[]>([]);
const fileURL = ref("");
const route = useRoute();
const router = useRouter();
const scrollBar = ref<InstanceType<typeof ElScrollbar>>();
const innerRef = ref<HTMLDivElement>();

function connect() {
  conn.value = sig.newConnection({
    roomID: roomID.value,
    name: username.value,
  });

  conn.value.on(KecpEventType.Open, () => {
    connected.value = true;
  });

  conn.value.on(KecpEventType.UserListInit, (evt: CustomEvent<string[]>) => {
    userList.value = evt.detail.slice();
  });

  conn.value.on(KecpEventType.UserJoin, (evt: CustomEvent<string>) => {
    userList.value.push(evt.detail);
  });

  conn.value.on(KecpEventType.UserLeave, (evt: CustomEvent<string>) => {
    if (userList.value.indexOf(evt.detail) !== -1) {
      userList.value.splice(userList.value.indexOf(evt.detail));
    }
  });

  conn.value.on(KecpEventType.Chat, (evt: CustomEvent<KecpMessage>) => {
    chatIDCouter.value += 1;
    chatList.value.push({
      chatID: chatIDCouter.value,
      type: evt.detail.type,
      name: evt.detail.name,
      target: evt.detail.target,
      payload: evt.detail.payload,
    });
    setTimeout(() => {
      scrollBar.value?.scrollTo(0, innerRef.value?.clientHeight);
    }, 150);
  });

  conn.value.on(
    KecpEventType.VideoOffer,
    async (evt: CustomEvent<AnswerPeer>) => {
      const stream = new MediaStream();
      evt.detail.setHandleTrackEvent((event) => {
        stream.addTrack(event.track);
        if (videoPlayerEl.value) {
          if (videoPlayerEl.value.srcObject !== stream) {
            videoPlayerEl.value.srcObject = stream;
            videoPlayerEl.value.autoplay = true;
          }
        }
      });
      await evt.detail.preAnswer();
      ElMessageBox.confirm(
        `${evt.detail.getTargetName()} will establish a connection with you. Confirm?`,
        "Warning",
        {
          confirmButtonText: "Confirm",
          cancelButtonText: "Cancel",
          type: "warning",
        }
      )
        .then(async () => {
          await evt.detail.preAnswerConfirm();
          ElMessage({
            type: "success",
            message: "Offer completed",
          });
        })
        .catch(() => {
          evt.detail.close();
          if (videoPlayerEl.value) {
            videoPlayerEl.value.srcObject = null;
            videoPlayerEl.value.autoplay = false;
          }
          ElMessage({
            type: "info",
            message: "Offer canceled",
          });
        });
    }
  );

  conn.value.on(KecpEventType.VideoAnswer, (evt: CustomEvent<KecpMessage>) => {
    if (evt.detail.payload.type === "answer") {
      ElMessage({
        type: "success",
        message: "Offer completed",
      });
    }
  });

  conn.value.on(KecpEventType.Error, (evt: CustomEvent<KecpMessage>) => {
    if (evt.detail.payload === "cannot join the room") {
      router.replace("/create");
    }
    switch (evt.detail.payload) {
      case "name is already in use":
        ElMessage({
          type: "error",
          message: "Name is already in use. Try another one.",
        });
        username.value = "";
        connected.value = false;
        break;
      case "not a valid name":
        ElMessage({
          type: "error",
          message: "Not a valid name. Try another one.",
        });
        username.value = "";
        connected.value = false;
        break;
    }
  });
}

onMounted(() => {
  roomID.value = route.params["roomID"] as string;
  username.value = (route.params["username"] ?? "") as string;
  if (roomID.value.length !== 16) {
    router.replace("/create");
  }
  if (videoPlayerEl.value !== undefined) {
    player.value = videojs(videoPlayerEl.value, {
      controls: true,
      fill: true,
    });
  }
  if (!connected.value && username.value !== "") {
    connect();
  }
});

function uploadFile(uploadFile: UploadFile, uploadFiles: UploadFiles) {
  if (uploadFiles === null) {
    return;
  }
  const file = uploadFile.raw;
  if (file !== undefined) {
    if (fileURL.value !== "") {
      URL.revokeObjectURL(fileURL.value);
    }
    fileURL.value = URL.createObjectURL(file);
    if (videoPlayerEl.value !== null && videoPlayerEl.value !== undefined) {
      videoPlayerEl.value.srcObject = null;
      player.value?.src({ type: file.type, src: fileURL.value });
    }
  }
}

function handleSendButton() {
  if (sendValue.value === "") {
    return;
  }
  conn.value?.send(
    JSON.stringify({
      type: KecpMessageType.Chat,
      name: conn.value.getName(),
      payload: sendValue.value,
    })
  );
}

function handleJoinButton() {
  username.value = usernameInput.value;
  connect();
}

function handleDialButton() {
  if (targetUser.value === conn.value?.getName()) {
    return;
  }

  const stream = videoPlayerEl.value?.captureStream();
  if (stream) {
    const offer = conn.value?.newOffer(targetUser.value);
    offer?.addStream(stream);
    offer?.setBandwidth(10000);
  }
}

function playerRewind() {
  if (player.value !== undefined) {
    const target = player.value.currentTime() - 5;
    player.value.currentTime(target);
  }
}

function playerForward() {
  if (player.value !== undefined) {
    const target = player.value.currentTime() + 5;
    player.value.currentTime(target);
  }
}

function playerPause() {
  if (player.value !== undefined) {
    if (player.value.paused()) {
      player.value.play();
    } else {
      player.value.pause();
    }
  }
}
</script>

<template>
  <div
    class="room-grid-container"
    @keydown.arrow-left="playerRewind"
    @keydown.arrow-right="playerForward"
    @keydown.space="playerPause"
  >
    <div class="left flex flex-col">
      <div class="video-container">
        <video
          ref="videoPlayerEl"
          class="video-js vjs-big-play-centered vjs-show-big-play-button-on-pause"
        ></video>
      </div>
      <div class="chat-input-container" @keyup.enter="handleSendButton()">
        <span>Chat:</span>
        <el-input class="chat-input px-2" v-model="sendValue" />
        <el-button id="send" :disabled="!connected" @click="handleSendButton()">
          Send
        </el-button>
      </div>
    </div>
    <div class="right">
      <div class="ctrl-container flex-shrink-0">
        <div class="flex justify-center" v-if="username !== ''">
          <el-select v-model="targetUser" placeholder="Select">
            <el-option
              v-for="item in userList"
              :key="item"
              :label="item"
              :value="item"
            />
          </el-select>
          <el-button
            id="send"
            :disabled="!connected || targetUser === '' || fileURL == ''"
            @click="handleDialButton()"
          >
            Dial
          </el-button>
        </div>
        <div class="flex justify-center" v-else>
          <el-input class="chat-input" v-model="usernameInput" />
          <el-button
            id="join"
            :disabled="connected"
            @click="handleJoinButton()"
          >
            Join
          </el-button>
        </div>
        <div class="mt-2">
          <el-upload multiple :auto-upload="false" :on-change="uploadFile">
            <el-button type="primary">Click to upload video files</el-button>
          </el-upload>
        </div>
      </div>
      <el-scrollbar class="chat-container mt-5" ref="scrollBar">
        <div ref="innerRef">
          <span
            class="leading-snug"
            v-for="chat in chatList"
            ref="chatBox"
            :key="chat.chatID"
          >
            {{ chat.name ?? "" }}: {{ chat.payload }} <br />
          </span>
        </div>
      </el-scrollbar>
    </div>
  </div>
</template>

<style>
@import "video.js/dist/video-js.css";
</style>

<style scoped>
.room-grid-container {
  padding: 30px 70px;
  display: grid;
  height: 100%;
  grid-template-areas: "left right";
  grid-template-columns: auto 300px;
  grid-auto-rows: 1fr;
  gap: 0 50px;
}

@media (max-width: 1024px) {
  .room-grid-container {
    padding: 30px 4vw;
    grid-template-areas:
      "left"
      "right";
    grid-template-columns: auto;
  }
}

.left {
  grid-area: left;
}

.video-container {
  aspect-ratio: 16/9;
  box-shadow: var(--el-box-shadow-light);
}

.chat-input-container {
  box-shadow: var(--el-box-shadow-light);
  border-bottom-left-radius: var(--el-border-radius-base);
  border-bottom-right-radius: var(--el-border-radius-base);
  display: flex;
  padding: 0 20px;
  align-items: center;
  height: 50px;
}

.chat-input {
  width: 100%;
}

.ctrl-container {
  grid-area: ctrl-container;
  display: flex;
  flex-direction: column;
  padding: 16px;
  border-radius: var(--el-border-radius-base);
  box-shadow: var(--el-box-shadow-light);
  min-height: 0;
}

@media (max-width: 1024px) {
  .ctrl-container {
    margin-top: 30px;
  }
}

.right {
  grid-area: right;
  display: flex;
  flex-direction: column;
}

.chat-container {
  padding: 16px;
  border-radius: var(--el-border-radius-base);
  box-shadow: var(--el-box-shadow-light);
  max-height: 65vh;
}
</style>
